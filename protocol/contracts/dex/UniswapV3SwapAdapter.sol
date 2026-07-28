// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

interface IUniswapV3QuoterV2 {
    struct QuoteExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint256 amountIn;
        uint24 fee;
        uint160 sqrtPriceLimitX96;
    }

    struct QuoteExactOutputSingleParams {
        address tokenIn;
        address tokenOut;
        uint256 amount;
        uint24 fee;
        uint160 sqrtPriceLimitX96;
    }

    function quoteExactInputSingle(QuoteExactInputSingleParams memory params)
        external
        returns (uint256 amountOut, uint160, uint32, uint256);

    function quoteExactOutputSingle(QuoteExactOutputSingleParams memory params)
        external
        returns (uint256 amountIn, uint160, uint32, uint256);
}

interface IUniswapV3SwapRouter {
    struct ExactInputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 amountIn;
        uint256 amountOutMinimum;
        uint160 sqrtPriceLimitX96;
    }

    struct ExactOutputSingleParams {
        address tokenIn;
        address tokenOut;
        uint24 fee;
        address recipient;
        uint256 amountOut;
        uint256 amountInMaximum;
        uint160 sqrtPriceLimitX96;
    }

    function exactInputSingle(ExactInputSingleParams calldata params) external payable returns (uint256 amountOut);
    function exactOutputSingle(ExactOutputSingleParams calldata params) external payable returns (uint256 amountIn);
}

interface IUniswapV3Deployment {
    function factory() external view returns (address);
}

interface IUniswapV3Factory {
    function getPool(address tokenA, address tokenB, uint24 fee) external view returns (address);
}

contract UniswapV3SwapAdapter {
    using SafeERC20 for IERC20;

    address public owner;
    IUniswapV3SwapRouter public immutable router;
    IUniswapV3QuoterV2 public immutable quoter;
    IUniswapV3Factory public immutable factory;
    mapping(address => mapping(address => uint24)) public poolFee;

    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event PoolFeeConfigured(address indexed tokenIn, address indexed tokenOut, uint24 fee);
    event Swap(
        address indexed sender,
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut,
        address recipient
    );

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    constructor(address initialOwner, address router_, address quoter_) {
        require(initialOwner != address(0), "Invalid owner");
        require(router_ != address(0) && router_.code.length > 0, "Invalid router");
        require(quoter_ != address(0) && quoter_.code.length > 0, "Invalid quoter");
        address routerFactory = IUniswapV3Deployment(router_).factory();
        require(routerFactory != address(0) && routerFactory.code.length > 0, "Invalid factory");
        require(IUniswapV3Deployment(quoter_).factory() == routerFactory, "Factory mismatch");
        owner = initialOwner;
        router = IUniswapV3SwapRouter(router_);
        quoter = IUniswapV3QuoterV2(quoter_);
        factory = IUniswapV3Factory(routerFactory);
        emit OwnerChanged(address(0), initialOwner);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid owner");
        emit OwnerChanged(owner, newOwner);
        owner = newOwner;
    }

    function setPoolFee(address tokenIn, address tokenOut, uint24 fee) external onlyOwner {
        require(tokenIn != address(0) && tokenIn.code.length > 0, "Invalid input token");
        require(tokenOut != address(0) && tokenOut.code.length > 0, "Invalid output token");
        require(tokenIn != tokenOut, "Identical tokens");
        require(fee > 0 && fee <= 1_000_000, "Invalid pool fee");
        address pool = factory.getPool(tokenIn, tokenOut, fee);
        require(pool != address(0) && pool.code.length > 0, "Uniswap pool not deployed");
        poolFee[tokenIn][tokenOut] = fee;
        emit PoolFeeConfigured(tokenIn, tokenOut, fee);
    }

    function getAmountOut(address tokenIn, address tokenOut, uint256 amountIn) public returns (uint256 amountOut) {
        uint24 fee = _configuredFee(tokenIn, tokenOut);
        (amountOut,,,) = quoter.quoteExactInputSingle(
            IUniswapV3QuoterV2.QuoteExactInputSingleParams(tokenIn, tokenOut, amountIn, fee, 0)
        );
    }

    function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) public returns (uint256 amountIn) {
        uint24 fee = _configuredFee(tokenIn, tokenOut);
        (amountIn,,,) = quoter.quoteExactOutputSingle(
            IUniswapV3QuoterV2.QuoteExactOutputSingleParams(tokenIn, tokenOut, amountOut, fee, 0)
        );
    }

    function swapExactTokensForTokens(
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        uint256 amountOutMin,
        address recipient
    ) external returns (uint256 amountOut) {
        require(recipient != address(0), "Invalid recipient");
        uint24 fee = _configuredFee(tokenIn, tokenOut);
        IERC20 input = IERC20(tokenIn);
        input.safeTransferFrom(msg.sender, address(this), amountIn);
        input.forceApprove(address(router), amountIn);

        amountOut = router.exactInputSingle(
            IUniswapV3SwapRouter.ExactInputSingleParams(tokenIn, tokenOut, fee, recipient, amountIn, amountOutMin, 0)
        );
        input.forceApprove(address(router), 0);
        emit Swap(msg.sender, tokenIn, tokenOut, amountIn, amountOut, recipient);
    }

    function swapTokensForExactTokens(
        address tokenIn,
        address tokenOut,
        uint256 amountOut,
        uint256 amountInMax,
        address recipient
    ) external returns (uint256 amountIn) {
        require(recipient != address(0), "Invalid recipient");
        uint24 fee = _configuredFee(tokenIn, tokenOut);
        amountIn = getAmountIn(tokenIn, tokenOut, amountOut);
        require(amountIn <= amountInMax, "Excessive input amount");

        IERC20 input = IERC20(tokenIn);
        input.safeTransferFrom(msg.sender, address(this), amountIn);
        input.forceApprove(address(router), amountIn);
        uint256 spent = router.exactOutputSingle(
            IUniswapV3SwapRouter.ExactOutputSingleParams(tokenIn, tokenOut, fee, recipient, amountOut, amountIn, 0)
        );
        input.forceApprove(address(router), 0);
        require(spent == amountIn, "Unexpected input amount");
        emit Swap(msg.sender, tokenIn, tokenOut, amountIn, amountOut, recipient);
    }

    function _configuredFee(address tokenIn, address tokenOut) private view returns (uint24 fee) {
        fee = poolFee[tokenIn][tokenOut];
        require(fee != 0, "Pool fee not configured");
    }
}
