// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.24;

contract MockOracle {
    address public owner;

    mapping(address => uint256) private prices; // token address => price

    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event PriceUpdated(address indexed token, uint256 price);

    constructor() {
        owner = msg.sender;
        emit OwnerChanged(address(0), owner);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid new owner");
        emit OwnerChanged(owner, newOwner);
        owner = newOwner;
    }

    function setPrice(address token, uint256 price) external onlyOwner {
        _setPrice(token, price);
    }

    function setPrices(address[] calldata tokens, uint256[] calldata newPrices) external onlyOwner {
        require(tokens.length == newPrices.length, "Mismatched array lengths");

        for (uint256 i = 0; i < tokens.length; i++) {
            _setPrice(tokens[i], newPrices[i]);
        }
    }

    function getPrice(address token) external view returns (uint256) {
        require(token != address(0), "Invalid token address");

        uint256 price = prices[token];
        require(price > 0, "Price not set for this token");
        return price;
    }

    function getPrices(address[] calldata tokens) external view returns (uint256[] memory result) {
        result = new uint256[](tokens.length);

        for (uint256 i = 0; i < tokens.length; i++) {
            require(tokens[i] != address(0), "Invalid token address");
            uint256 price = prices[tokens[i]];

            require(price > 0, "Price not set for this token");
            result[i] = price;
        }
    }

    function _setPrice(address token, uint256 price) internal {
        require(token != address(0), "Invalid token address");
        require(price > 0, "Price must be positive");

        prices[token] = price;
        
        emit PriceUpdated(token, price);
    }
}