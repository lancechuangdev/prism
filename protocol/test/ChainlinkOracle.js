import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("ChainlinkOracle", function () {
  let owner;
  let stranger;
  let token;
  let feed;
  let oracle;

  beforeEach(async function () {
    [owner, stranger] = await ethers.getSigners();
    token = await ethers.deployContract("PositionToken", ["Token", "TOK"]);
    feed = await ethers.deployContract("MockChainlinkAggregator", [8]);
    oracle = await ethers.deployContract("ChainlinkOracle", [owner.address]);
  });

  it("normalizes valid feed answers to 18 decimals", async function () {
    const latest = await ethers.provider.getBlock("latest");
    await feed.setRoundData(10, 2_500n * 10n ** 8n, latest.timestamp, 10);
    await oracle.setFeed(
      await token.getAddress(),
      await feed.getAddress(),
      3600,
    );

    expect(await oracle.getPrice(await token.getAddress())).to.equal(
      2_500n * 10n ** 18n,
    );
  });

  it("rejects stale, non-positive, and incomplete rounds", async function () {
    const latest = await ethers.provider.getBlock("latest");
    await oracle.setFeed(await token.getAddress(), await feed.getAddress(), 60);

    await feed.setRoundData(10, 1, latest.timestamp - 61, 10);
    await expect(oracle.getPrice(await token.getAddress())).to.be.revertedWith(
      "Stale feed answer",
    );

    await feed.setRoundData(10, 0, latest.timestamp, 10);
    await expect(oracle.getPrice(await token.getAddress())).to.be.revertedWith(
      "Invalid feed answer",
    );

    await feed.setRoundData(10, 1, latest.timestamp, 9);
    await expect(oracle.getPrice(await token.getAddress())).to.be.revertedWith(
      "Incomplete feed round",
    );
  });

  it("restricts feed configuration to the owner", async function () {
    await expect(
      oracle
        .connect(stranger)
        .setFeed(await token.getAddress(), await feed.getAddress(), 3600),
    ).to.be.revertedWith("Not the owner");
  });
});
