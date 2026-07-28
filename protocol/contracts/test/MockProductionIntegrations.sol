// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockDecimalsToken is ERC20 {
    uint8 private immutable tokenDecimals;

    constructor(string memory name_, string memory symbol_, uint8 decimals_) ERC20(name_, symbol_) {
        tokenDecimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return tokenDecimals;
    }

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MockChainlinkAggregator {
    uint8 public immutable decimals;
    uint80 public roundId;
    int256 public answer;
    uint256 public updatedAt;
    uint80 public answeredInRound;

    constructor(uint8 decimals_) {
        decimals = decimals_;
    }

    function setRoundData(uint80 roundId_, int256 answer_, uint256 updatedAt_, uint80 answeredInRound_) external {
        roundId = roundId_;
        answer = answer_;
        updatedAt = updatedAt_;
        answeredInRound = answeredInRound_;
    }

    function latestRoundData()
        external
        view
        returns (uint80, int256, uint256, uint256, uint80)
    {
        return (roundId, answer, updatedAt, updatedAt, answeredInRound);
    }
}

contract MockUniswapV3 {
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

    uint256 public constant RATE_SCALE = 1e18;
    mapping(address => mapping(address => uint256)) public rate;

    function factory() external view returns (address) {
        return address(this);
    }

    function getPool(address tokenA, address tokenB, uint24) external view returns (address) {
        return rate[tokenA][tokenB] > 0 ? address(this) : address(0);
    }

    function setRate(address tokenIn, address tokenOut, uint256 rate_) external {
        rate[tokenIn][tokenOut] = rate_;
    }

    function quoteExactInputSingle(QuoteExactInputSingleParams memory params)
        external
        view
        returns (uint256 amountOut, uint160, uint32, uint256)
    {
        amountOut = _amountOut(params.tokenIn, params.tokenOut, params.amountIn);
        return (amountOut, 0, 0, 0);
    }

    function quoteExactOutputSingle(QuoteExactOutputSingleParams memory params)
        external
        view
        returns (uint256 amountIn, uint160, uint32, uint256)
    {
        amountIn = _amountIn(params.tokenIn, params.tokenOut, params.amount);
        return (amountIn, 0, 0, 0);
    }

    function exactInputSingle(ExactInputSingleParams calldata params) external payable returns (uint256 amountOut) {
        amountOut = _amountOut(params.tokenIn, params.tokenOut, params.amountIn);
        require(amountOut >= params.amountOutMinimum, "Too little output");
        require(IERC20(params.tokenIn).transferFrom(msg.sender, address(this), params.amountIn), "Input transfer failed");
        require(IERC20(params.tokenOut).transfer(params.recipient, amountOut), "Output transfer failed");
    }

    function exactOutputSingle(ExactOutputSingleParams calldata params) external payable returns (uint256 amountIn) {
        amountIn = _amountIn(params.tokenIn, params.tokenOut, params.amountOut);
        require(amountIn <= params.amountInMaximum, "Too much input");
        require(IERC20(params.tokenIn).transferFrom(msg.sender, address(this), amountIn), "Input transfer failed");
        require(IERC20(params.tokenOut).transfer(params.recipient, params.amountOut), "Output transfer failed");
    }

    function _amountOut(address tokenIn, address tokenOut, uint256 amountIn) private view returns (uint256) {
        uint256 configuredRate = rate[tokenIn][tokenOut];
        require(configuredRate > 0, "Rate missing");
        return (amountIn * configuredRate) / RATE_SCALE;
    }

    function _amountIn(address tokenIn, address tokenOut, uint256 amountOut) private view returns (uint256) {
        uint256 configuredRate = rate[tokenIn][tokenOut];
        require(configuredRate > 0, "Rate missing");
        return (amountOut * RATE_SCALE + configuredRate - 1) / configuredRate;
    }
}
