// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

interface IERC20RouterLike {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

contract FixedRateSwap {
    uint256 private constant RATE_SCALE = 1e18;

    address public owner;
    mapping(address => mapping(address => uint256)) public rateOutPerTokenIn; // tokenIn => tokenOut => rateOutPerTokenIn

    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event RateChanged(address indexed tokenIn, address indexed tokenOut, uint256 rateOutPerTokenIn);
    event Swap(address indexed sender, address indexed tokenIn, address indexed tokenOut, uint256 amountIn, uint256 amountOut, address recipient);

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    constructor() {
        owner = msg.sender;
        emit OwnerChanged(address(0), owner);
    }

    function setRate(address tokenIn, address tokenOut, uint256 rate) external onlyOwner {
        require(tokenIn != address(0) && tokenOut != address(0), "Invalid token address");
        require(rate > 0, "Rate must be positive");

        rateOutPerTokenIn[tokenIn][tokenOut] = rate;
        emit RateChanged(tokenIn, tokenOut, rate);
    }

    // Exact-input swap: sell a fixed amount of tokenIn for as much tokenOut as possible, respecting the minimum amountOut specified.
    function swapExactTokensForTokens(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin, address recipient)
        external
        returns (uint256 amountOut)
    {
        require(recipient != address(0), "Invalid recipient");

        amountOut = getAmountOut(tokenIn, tokenOut, amountIn);
        require(amountOut >= amountOutMin, "Insufficient output amount");

        bool pulled = IERC20RouterLike(tokenIn).transferFrom(msg.sender, address(this), amountIn);
        require(pulled, "Token transfer failed");

        bool sent = IERC20RouterLike(tokenOut).transfer(recipient, amountOut);
        require(sent, "Token transfer failed");

        emit Swap(msg.sender, tokenIn, tokenOut, amountIn, amountOut, recipient);
    }

    function getAmountOut(address tokenIn, address tokenOut, uint256 amountIn) public view returns (uint256 amountOut) {
        uint256 rate = rateOutPerTokenIn[tokenIn][tokenOut];
        require(rate > 0, "Rate not set for this pair");
        amountOut = (amountIn * rate) / RATE_SCALE;
    }

    // Exact-output swap: buy a fixed amount of tokenOut for as little tokenIn as possible, respecting the maximum amountIn specified.
    function swapTokensForExactTokens(address tokenIn, address tokenOut, uint256 amountOut, uint256 amountInMax, address recipient)
        external
        returns (uint256 amountIn)
    {
        require(recipient != address(0), "Invalid recipient");

        amountIn = getAmountIn(tokenIn, tokenOut, amountOut);
        require(amountIn <= amountInMax, "Excessive input amount");

        bool pulled = IERC20RouterLike(tokenIn).transferFrom(msg.sender, address(this), amountIn);
        require(pulled, "Token transfer failed");

        bool sent = IERC20RouterLike(tokenOut).transfer(recipient, amountOut);
        require(sent, "Token transfer failed");

        emit Swap(msg.sender, tokenIn, tokenOut, amountIn, amountOut, recipient);
    }

    function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) public view returns (uint256 amountIn) {
        uint256 rate = rateOutPerTokenIn[tokenIn][tokenOut];
        require(rate > 0, "Rate not set for this pair");
        amountIn = (amountOut * RATE_SCALE + rate - 1) / rate; // Round up to ensure enough input tokens
    }
}