// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.28;

interface IChainlinkAggregatorV3 {
    function decimals() external view returns (uint8);
    function latestRoundData()
        external
        view
        returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound);
}

contract ChainlinkOracle {
    uint8 private constant PRICE_DECIMALS = 18;

    struct FeedConfig {
        address feed;
        uint48 maxStaleness;
    }

    address public owner;
    mapping(address => FeedConfig) public feeds;

    event OwnerChanged(address indexed oldOwner, address indexed newOwner);
    event FeedConfigured(address indexed token, address indexed feed, uint48 maxStaleness);

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }

    constructor(address initialOwner) {
        require(initialOwner != address(0), "Invalid owner");
        owner = initialOwner;
        emit OwnerChanged(address(0), initialOwner);
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "Invalid owner");
        emit OwnerChanged(owner, newOwner);
        owner = newOwner;
    }

    function setFeed(address token, address feed, uint48 maxStaleness) external onlyOwner {
        require(token != address(0) && token.code.length > 0, "Invalid token");
        require(feed != address(0) && feed.code.length > 0, "Invalid feed");
        require(maxStaleness > 0, "Invalid staleness");

        uint8 feedDecimals = IChainlinkAggregatorV3(feed).decimals();
        require(feedDecimals <= 36, "Unsupported feed decimals");

        feeds[token] = FeedConfig(feed, maxStaleness);
        emit FeedConfigured(token, feed, maxStaleness);
    }

    function getPrice(address token) external view returns (uint256) {
        FeedConfig memory config = feeds[token];
        require(config.feed != address(0), "Feed not configured");

        (uint80 roundId, int256 answer,, uint256 updatedAt, uint80 answeredInRound) =
            IChainlinkAggregatorV3(config.feed).latestRoundData();
        require(answer > 0, "Invalid feed answer");
        require(updatedAt != 0 && updatedAt <= block.timestamp, "Invalid feed timestamp");
        require(block.timestamp - updatedAt <= config.maxStaleness, "Stale feed answer");
        require(answeredInRound >= roundId, "Incomplete feed round");

        uint8 feedDecimals = IChainlinkAggregatorV3(config.feed).decimals();
        uint256 unsignedAnswer = uint256(answer);
        if (feedDecimals < PRICE_DECIMALS) {
            return unsignedAnswer * (10 ** (PRICE_DECIMALS - feedDecimals));
        }
        if (feedDecimals > PRICE_DECIMALS) {
            return unsignedAnswer / (10 ** (feedDecimals - PRICE_DECIMALS));
        }
        return unsignedAnswer;
    }
}
