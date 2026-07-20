// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

// Asset tokens for lending and borrowing, for example:
// lendToken (stablecoins): USDT / USDC / BUSD
// borrowToken (volatile assets): WBTC / WETH / DAI
// Alice deposits 10,000 USDC into the lending pool.
// Bob deposits 1 BTC as collateral.
// Bob borrows 5,000 USDC.
interface IERC20Like {
    function approve(address spender, uint256 amount) external returns (bool);
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}


// Debt tokens are minted when users borrow from the lending pool and burned when they repay their loans.
interface IDebtTokenLike {
    function mint(address to, uint256 amount) external returns (bool);
    function burn(address from, uint256 amount) external returns (bool);
}

interface IOracleLike {
    function getPrice(address token) external view returns (uint256);
}

interface IDexSwapLike {
    function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) external view returns (uint256);
    function getAmountOut(address tokenIn, address tokenOut, uint256 amountIn) external view returns (uint256);
    function swapExactTokensForTokens(
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        uint256 amountOutMin,
        address recipient
    ) external returns (uint256 amountOut);
    function swapTokensForExactTokens(
        address tokenIn,
        address tokenOut,
        uint256 amountOut,
        uint256 amountInMax,
        address recipient
    ) external returns (uint256 amountIn);
}

contract PrismPool {
    uint256 private constant RATE_SCALE = 1e8;
    uint256 private constant PRICE_SCALE = 1e18; // Oracle price scale
    uint256 private constant SECONDS_PER_YEAR = 365 days; // Number of seconds in a year

    enum PoolState {
        FUNDING, // Accepting lender funds and borrower collateral
        ACTIVE, // Settlement completed and loan is active
        REPAID, // Loan completed and lender proceeds are available
        LIQUIDATED, // Collateral has already been liquidated
        CANCELLED // Pool could not be matched and deposits can be recovered
    }

    struct CreatePoolParams {
        uint256 settleTime;
        uint256 maturityTime;
        uint256 interestRate; // Annual interest rate in basis points (1% = 100 bps)
        uint256 maxLendSupply;
        uint256 collateralizationRatio; // Collateralization ratio in basis points (1% = 100 bps)
        address lendToken;
        address collateralToken;
        address lenderPositionToken;
        address borrowerPositionToken;
        uint256 liquidateRate; // Liquidation threshold in basis points (1% = 100 bps)
    }

    struct PoolBaseInfo {
        uint256 settleTime;
        uint256 maturityTime;
        uint256 interestRate;
        uint256 maxLendSupply;
        uint256 totalLendDeposited;
        uint256 totalCollateralDeposited;
        uint256 collateralizationRatio;
        address lendToken;
        address collateralToken;
        PoolState state;
        address lenderPositionToken;
        address borrowerPositionToken;
        uint256 liquidateRate;
    }

    struct PoolDataInfo {
        uint256 settleAmountLend; // required lend amount
        uint256 settleAmountBorrow; // required collateral amount
        uint256 finishAmountLend;
        uint256 finishAmountBorrow;
        uint256 liquidationAmountLend;
        uint256 liquidationAmountBorrow;
    }

    struct LendInfo {
        uint256 stakeAmount;
        uint256 refundAmount;
        bool hasRefunded;
        bool hasClaimed;
    }

    struct BorrowInfo {
        uint256 stakeAmount;
        uint256 refundAmount;
        bool hasRefunded;
        bool hasClaimed;
    }

    address public owner;
    address public oracle;
    address public dexSwap;
    address payable public feeAddress;
    bool public globalPaused;
    uint256 public minLendAmount = 100 ether;
    uint256 public minBorrowAmount = 1 ether;

    PoolBaseInfo[] private pools;
    PoolDataInfo[] private poolData;
    mapping(address => mapping(uint256 => LendInfo)) public userLendInfo; // user => poolId => LendInfo
    mapping(address => mapping(uint256 => BorrowInfo)) public userBorrowInfo; // user => poolId => BorrowInfo

    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event PoolCreated(
        uint256 indexed poolId,
        address indexed lendToken,
        address indexed borrowToken,
        address spToken,
        address jpToken,
        uint256 settleTime,
        uint256 endTime
    );
    event FeeAddressChanged(address indexed oldFeeAddress, address indexed newFeeAddress);
    event DexSwapChanged(address indexed oldDexSwap, address indexed newDexSwap);
    event MinLendAmountChanged(uint256 oldMinAmount, uint256 newMinAmount);
    event MinBorrowAmountChanged(uint256 oldMinAmount, uint256 newMinAmount);
    event PauseChanged(bool paused);
    event DepositLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount);
    event DepositBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount);
    event RefundLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount);
    event RefundBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount);
    event ClaimLend(address indexed lender, uint256 indexed poolId, address indexed spToken, uint256 spAmount);
    event ClaimBorrow(
        address indexed borrower,
        uint256 indexed poolId,
        address indexed jpToken,
        uint256 jpAmount,
        uint256 loanAmount
    );
    event PoolRepaid(
        uint256 indexed poolId,
        address indexed router,
        uint256 collateralSold,
        uint256 repaymentAmount,
        uint256 remainingCollateralAmount
    );
    event PoolLiquidated(
        uint256 indexed poolId,
        address indexed router,
        uint256 collateralSold,
        uint256 lendTokenRecovered,
        uint256 remainingCollateralAmount
    );
    event WithdrawLend(address indexed lender, uint256 indexed poolId, uint256 spAmount, uint256 lendAmount);
    event WithdrawBorrow(address indexed borrower, uint256 indexed poolId, uint256 jpAmount, uint256 collateralAmount);
    event StateChanged(uint256 indexed poolId, PoolState oldState, PoolState newState);

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    modifier whenNotPaused() {
        require(!globalPaused, "Contract is paused");
        _;
    }

    modifier isState(uint256 poolId, PoolState expectedState) {
        require(poolId < pools.length, "Invalid pool ID");
        require(pools[poolId].state == expectedState, "Invalid pool state");
        _;
    }

    modifier stateClosed(uint256 poolId) {
        require(
            pools[poolId].state == PoolState.REPAID || pools[poolId].state == PoolState.LIQUIDATED,
            "Pool not closed"
        );
        _;
    }

    modifier beforeSettle(uint256 poolId) {
        require(poolId < pools.length, "Invalid pool ID");
        require(block.timestamp < pools[poolId].settleTime, "Operation not allowed after settle time");
        _;
    }

    modifier afterSettle(uint256 poolId) {
        require(poolId < pools.length, "Invalid pool ID");
        require(block.timestamp >= pools[poolId].settleTime, "Operation allowed only after settle time");
        _;
    }

    modifier afterEnd(uint256 poolId) {
        require(block.timestamp >= pools[poolId].maturityTime, "Operation allowed only after end time");
        _;
    }

    constructor(address oracle_, address dexSwap_, address payable feeAddress_) {
        require(oracle_ != address(0), "Invalid oracle address");
        require(dexSwap_ != address(0), "Invalid dexSwap address");
        require(feeAddress_ != address(0), "Invalid feeAddress");

        owner = msg.sender;
        oracle = oracle_;
        dexSwap = dexSwap_;
        feeAddress = feeAddress_;

        emit OwnerChanged(address(0), owner);
        emit DexSwapChanged(address(0), dexSwap);
        emit FeeAddressChanged(address(0), feeAddress);
    }

    function createPool(CreatePoolParams calldata params) external onlyOwner returns (uint256 poolId) {
        require(params.settleTime > block.timestamp, "Settle time must be in the future");
        require(params.maturityTime > params.settleTime, "End time must be after settle time");
        require(params.interestRate > 0, "Interest rate must be positive");
        require(params.maxLendSupply > 0, "Max supply must be positive");
        require(params.collateralizationRatio > 0, "Mortgage rate must be positive");
        require(params.liquidateRate > 0, "Liquidate rate must be positive");
        require(params.lendToken != address(0), "Invalid lend token address");
        require(params.collateralToken != address(0), "Invalid borrow token address");
        require(params.lendToken != params.collateralToken, "Lend and borrow tokens must be different");
        require(params.lenderPositionToken != address(0), "Invalid sp token address");
        require(params.borrowerPositionToken != address(0), "Invalid jp token address"); 
        require(params.lenderPositionToken != params.borrowerPositionToken, "SP and JP tokens must be different");

        poolId = pools.length;

        pools.push(
            PoolBaseInfo({
                settleTime: params.settleTime,
                maturityTime: params.maturityTime,
                interestRate: params.interestRate,
                maxLendSupply: params.maxLendSupply,
                totalLendDeposited: 0,
                totalCollateralDeposited: 0,
                collateralizationRatio: params.collateralizationRatio,
                lendToken: params.lendToken,
                collateralToken: params.collateralToken,
                state: PoolState.FUNDING,
                lenderPositionToken: params.lenderPositionToken,
                borrowerPositionToken: params.borrowerPositionToken,
                liquidateRate: params.liquidateRate
            })
        );

        poolData.push(
            PoolDataInfo({
                settleAmountLend: 0,
                settleAmountBorrow: 0,
                finishAmountLend: 0,
                finishAmountBorrow: 0,
                liquidationAmountLend: 0,
                liquidationAmountBorrow: 0
            })
        );

        emit PoolCreated(
            poolId, 
            params.lendToken, 
            params.collateralToken, 
            params.lenderPositionToken, 
            params.borrowerPositionToken, 
            params.settleTime, 
            params.maturityTime
        );
    }

    function poolCount() external view returns (uint256) {
        return pools.length;
    }

    function getPool(uint256 poolId) external view returns (PoolBaseInfo memory poolInfo) {
        require(poolId < pools.length, "Invalid pool ID");
        poolInfo = pools[poolId];
    }

    function getPoolData(uint256 poolId) external view returns (PoolDataInfo memory poolDataInfo) {
        require(poolId < poolData.length, "Invalid pool ID");
        poolDataInfo = poolData[poolId];
    }

    function getPoolState(uint256 poolId) external view returns (PoolState state) {
        require(poolId < pools.length, "Invalid pool ID");
        state = pools[poolId].state;
    }

    function isBeforeSettleTime(uint256 poolId) external view returns (bool) {
        require(poolId < pools.length, "Invalid pool ID");
        return block.timestamp < pools[poolId].settleTime;
    }

    function getRequiredRepayment(uint256 poolId) public view returns (uint256) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];

        uint256 term = pool.maturityTime - pool.settleTime;
        uint256 interest = (data.settleAmountBorrow * pool.interestRate * term) / (RATE_SCALE * SECONDS_PER_YEAR);
        return data.settleAmountBorrow + interest;
    }

    function isUndercollateralized(uint256 poolId) public view returns (bool) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];

        if (pool.state != PoolState.ACTIVE || data.settleAmountLend == 0 || data.settleAmountBorrow == 0) {
            return false;
        }

        uint256 lendPrice = IOracleLike(oracle).getPrice(pool.lendToken);
        uint256 borrowPrice = IOracleLike(oracle).getPrice(pool.collateralToken);
        require(lendPrice > 0 && borrowPrice > 0, "Invalid price from oracle");

        uint256 borrowToLendRatio = (borrowPrice * PRICE_SCALE) / lendPrice;
        uint256 collateralValueInLend = (data.settleAmountBorrow * borrowToLendRatio) / PRICE_SCALE;
        uint256 liquidationThreshold = (data.settleAmountLend * (RATE_SCALE + pool.liquidateRate)) / RATE_SCALE;

        return collateralValueInLend < liquidationThreshold;
    }

    function depositLend(uint256 poolId, uint256 amount)
        external
        whenNotPaused
        isState(poolId, PoolState.FUNDING)
        beforeSettle(poolId)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(amount > 0, "Invalid amount");

        PoolBaseInfo storage pool = pools[poolId];
        LendInfo storage lendInfo = userLendInfo[msg.sender][poolId];

        require(amount >= minLendAmount, "Amount below minimum lend amount");
        require(pool.totalLendDeposited + amount <= pool.maxLendSupply, "Exceeds max supply");

        bool success = IERC20Like(pool.lendToken).transferFrom(msg.sender, address(this), amount);
        require(success, "Token transfer failed");

        lendInfo.stakeAmount += amount;
        lendInfo.hasRefunded = false;
        lendInfo.hasClaimed = false;
        pool.totalLendDeposited += amount;

        emit DepositLend(msg.sender, poolId, pool.lendToken, amount);
    }

    function depositBorrow(uint256 poolId, uint256 amount)
        external
        whenNotPaused
        isState(poolId, PoolState.FUNDING)
        beforeSettle(poolId)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(amount >= minBorrowAmount, "Amount below minimum borrow amount");

        PoolBaseInfo storage pool = pools[poolId];
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][poolId];

        bool success = IERC20Like(pool.collateralToken).transferFrom(msg.sender, address(this), amount);
        require(success, "Token transfer failed");

        borrowInfo.stakeAmount += amount;
        borrowInfo.hasRefunded = false;
        borrowInfo.hasClaimed = false;
        pool.totalCollateralDeposited += amount;

        emit DepositBorrow(msg.sender, poolId, pool.collateralToken, amount);
    }

    function settle(uint256 poolId) external onlyOwner whenNotPaused isState(poolId, PoolState.FUNDING) afterSettle(poolId) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];

        if (pool.totalLendDeposited == 0 || pool.totalCollateralDeposited == 0) {
            pool.state = PoolState.CANCELLED;
            emit StateChanged(poolId, PoolState.FUNDING, PoolState.CANCELLED);
            return;
        }

        uint256 lendPrice = IOracleLike(oracle).getPrice(pool.lendToken);
        uint256 borrowPrice = IOracleLike(oracle).getPrice(pool.collateralToken);
        require(lendPrice > 0 && borrowPrice > 0, "Invalid price from oracle");

        uint256 borrowToLendRatio = (borrowPrice * PRICE_SCALE) / lendPrice;
        uint256 collateralValueInLend = (pool.totalCollateralDeposited * borrowToLendRatio) / PRICE_SCALE;
        uint256 maxSettleLend = (collateralValueInLend * RATE_SCALE) / pool.collateralizationRatio;

        if (pool.totalLendDeposited > maxSettleLend) {
            data.settleAmountLend = maxSettleLend;
            data.settleAmountBorrow = pool.totalCollateralDeposited;
        } else {
            data.settleAmountLend = pool.totalLendDeposited;
            data.settleAmountBorrow = (pool.totalLendDeposited * pool.collateralizationRatio * lendPrice) / (borrowPrice * RATE_SCALE);
        }

        _setPoolState(poolId, PoolState.ACTIVE);
    }

    function refundExcessLend(uint256 poolId) external whenNotPaused isState(poolId, PoolState.ACTIVE) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        LendInfo storage lendInfo = userLendInfo[msg.sender][poolId];

        require(lendInfo.stakeAmount > 0, "No refund");
        require(!lendInfo.hasRefunded, "Already refunded");

        // remaining = TOTAL lender money that was NOT used in settlement
        uint256 remaining = pool.totalLendDeposited - data.settleAmountLend;
        require(remaining > 0, "No refund");

        // refund = THIS lender's amount * THIS lender's share of the pool
        uint256 refundAmount = (lendInfo.stakeAmount * remaining) / pool.totalLendDeposited;
        lendInfo.refundAmount += refundAmount;
        lendInfo.hasRefunded = true;

        bool success = IERC20Like(pool.lendToken).transfer(msg.sender, refundAmount);
        require(success, "Refund transfer failed");

        emit RefundLend(msg.sender, poolId, pool.lendToken, refundAmount);
    }

    function refundExcessCollateral(uint256 poolId) external whenNotPaused isState(poolId, PoolState.ACTIVE) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][poolId];

        require(borrowInfo.stakeAmount > 0, "No borrow stake");
        require(!borrowInfo.hasRefunded, "Already refunded");

        uint256 remaining = pool.totalCollateralDeposited - data.settleAmountBorrow;
        require(remaining > 0, "No borrow refund");

        uint256 refundAmount = (borrowInfo.stakeAmount * remaining) / pool.totalCollateralDeposited;
        borrowInfo.refundAmount += refundAmount;
        borrowInfo.hasRefunded = true;

        bool success = IERC20Like(pool.collateralToken).transfer(msg.sender, refundAmount);
        require(success, "Borrow refund transfer failed");

        emit RefundBorrow(msg.sender, poolId, pool.collateralToken, refundAmount);
    }

    function claimLenderPosition(uint256 poolId) external whenNotPaused isState(poolId, PoolState.ACTIVE) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        LendInfo storage lendInfo = userLendInfo[msg.sender][poolId];

        require(lendInfo.stakeAmount > 0, "No lender position");
        require(!lendInfo.hasClaimed, "Already claimed");

        uint256 lenderPosition = (data.settleAmountLend * lendInfo.stakeAmount) / pool.totalLendDeposited;
        require(lenderPosition > 0, "No claim");

        lendInfo.hasClaimed = true;

        bool success = IDebtTokenLike(pool.lenderPositionToken).mint(msg.sender, lenderPosition);
        require(success, "Mint failed");

        emit ClaimLend(msg.sender, poolId, pool.lenderPositionToken, lenderPosition);
    }

    function claimBorrowerPositionAndLoan(uint256 poolId) external whenNotPaused isState(poolId, PoolState.ACTIVE) {
        require(poolId < pools.length, "Invalid pool ID");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        BorrowInfo storage borrowInfo = userBorrowInfo[msg.sender][poolId];

        require(borrowInfo.stakeAmount > 0, "No borrower position");
        require(!borrowInfo.hasClaimed, "Already claimed");

        uint256 borrowerPosition = (data.settleAmountBorrow * borrowInfo.stakeAmount) / pool.totalCollateralDeposited;
        uint256 loanAmount = (data.settleAmountLend * borrowInfo.stakeAmount) / pool.totalCollateralDeposited;
        require(borrowerPosition > 0, "No borrower position");
        require(loanAmount > 0, "No loan");

        borrowInfo.hasClaimed = true;

        bool minted = IDebtTokenLike(pool.borrowerPositionToken).mint(msg.sender, borrowerPosition);
        require(minted, "Mint failed");

        bool transferred = IERC20Like(pool.lendToken).transfer(msg.sender, loanAmount);
        require(transferred, "Loan transfer failed");

        emit ClaimBorrow(msg.sender, poolId, pool.borrowerPositionToken, borrowerPosition, loanAmount);
    }

    function repayPool(uint256 poolId, uint256 maxCollateralAmount)
        external
        onlyOwner
        whenNotPaused
        isState(poolId, PoolState.ACTIVE)
        afterEnd(poolId)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(dexSwap != address(0), "Dex swap not available");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        uint256 requiredRepayment = getRequiredRepayment(poolId);
        uint256 collateralToSell = IDexSwapLike(dexSwap).getAmountIn(
            pool.collateralToken,
            pool.lendToken,
            requiredRepayment
        );

        require(collateralToSell <= maxCollateralAmount, "Swap slippage too high");
        require(collateralToSell <= data.settleAmountBorrow, "Insufficient collateral");

        // approve the DEX router to spend the collateral tokens first, 
        // then the DEX will pull the collateral and swap them for lend tokens.
        bool approved = IERC20Like(pool.collateralToken).approve(dexSwap, collateralToSell);
        require(approved, "Collateral approve failed");

        uint256 soldAmount = IDexSwapLike(dexSwap).swapTokensForExactTokens(
            pool.collateralToken,
            pool.lendToken,
            requiredRepayment,
            maxCollateralAmount,
            address(this)
        );

        data.finishAmountLend = requiredRepayment;
        data.finishAmountBorrow = data.settleAmountBorrow - soldAmount;

        _setPoolState(poolId, PoolState.REPAID);

        emit PoolRepaid(poolId, dexSwap, soldAmount, requiredRepayment, data.finishAmountBorrow);
    }

    function liquidate(uint256 poolId, uint256 maxCollateralAmount)
        external
        onlyOwner
        whenNotPaused
        isState(poolId, PoolState.ACTIVE)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(dexSwap != address(0), "Dex swap not available");
        require(isUndercollateralized(poolId), "Pool is sufficiently collateralized");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];
        uint256 requiredRepayment = getRequiredRepayment(poolId);
        uint256 collateralToSell = IDexSwapLike(dexSwap).getAmountIn(
            pool.collateralToken,
            pool.lendToken,
            requiredRepayment
        );

        uint256 soldAmount;
        uint256 recoveredAmount;

        if (collateralToSell <= data.settleAmountBorrow) {
            require(collateralToSell <= maxCollateralAmount, "Dex slippage too high");

            bool approved = IERC20Like(pool.collateralToken).approve(dexSwap, collateralToSell);
            require(approved, "Collateral approve failed");

            soldAmount = IDexSwapLike(dexSwap).swapTokensForExactTokens(
                pool.collateralToken,
                pool.lendToken,
                requiredRepayment,
                maxCollateralAmount,
                address(this)
            );
            recoveredAmount = requiredRepayment;
        } else {
            soldAmount = data.settleAmountBorrow;
            require(soldAmount <= maxCollateralAmount, "Dex slippage too high");

            bool approved = IERC20Like(pool.collateralToken).approve(dexSwap, soldAmount);
            require(approved, "Collateral approve failed");

            recoveredAmount = IDexSwapLike(dexSwap).swapExactTokensForTokens(
                pool.collateralToken,
                pool.lendToken,
                soldAmount,
                0,
                address(this)
            );
        }

        data.liquidationAmountLend = recoveredAmount;
        data.liquidationAmountBorrow = data.settleAmountBorrow - soldAmount;

        _setPoolState(poolId, PoolState.LIQUIDATED);

        emit PoolLiquidated(poolId, dexSwap, soldAmount, recoveredAmount, data.liquidationAmountBorrow);
    }

    function withdrawLend(uint256 poolId, uint256 lenderPosition)
        external
        whenNotPaused
        stateClosed(poolId)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(lenderPosition > 0, "No lender position");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];

        uint256 poolAvailableLendAmount =
            pool.state == PoolState.REPAID ? data.finishAmountLend : data.liquidationAmountLend;
        uint256 lendAmount = (poolAvailableLendAmount * lenderPosition) / data.settleAmountLend;

        bool burned = IDebtTokenLike(pool.lenderPositionToken).burn(msg.sender, lenderPosition);
        require(burned, "Burn lender position token failed");

        bool transferred = IERC20Like(pool.lendToken).transfer(msg.sender, lendAmount);
        require(transferred, "Lender withdraw transfer failed");

        emit WithdrawLend(msg.sender, poolId, lenderPosition, lendAmount);
    }

    function withdrawBorrow(uint256 poolId, uint256 borrowerPosition)
        external
        whenNotPaused
        stateClosed(poolId)
    {
        require(poolId < pools.length, "Invalid pool ID");
        require(borrowerPosition > 0, "No borrower position");

        PoolBaseInfo storage pool = pools[poolId];
        PoolDataInfo storage data = poolData[poolId];

        uint256 poolAvailableCollateralAmount =
            pool.state == PoolState.REPAID ? data.finishAmountBorrow : data.liquidationAmountBorrow;
        uint256 collateralAmount = (poolAvailableCollateralAmount * borrowerPosition) / data.settleAmountBorrow;

        bool burned = IDebtTokenLike(pool.borrowerPositionToken).burn(msg.sender, borrowerPosition);
        require(burned, "Burn borrower position failed");

        bool transferred = IERC20Like(pool.collateralToken).transfer(msg.sender, collateralAmount);
        require(transferred, "Borrower withdraw transfer failed");

        emit WithdrawBorrow(msg.sender, poolId, borrowerPosition, collateralAmount);
    }

    function _setPoolState(uint256 poolId, PoolState newState) internal {
        require(poolId < pools.length, "Invalid pool ID");
        PoolBaseInfo storage pool = pools[poolId];
        PoolState oldState = pool.state;
        pool.state = newState;
        emit StateChanged(poolId, oldState, newState);
    }
}
