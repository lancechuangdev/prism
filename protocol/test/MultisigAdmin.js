import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("ThresholdMultiSig administration", function () {
  let owner;
  let alice;
  let bob;
  let carol;
  let asset;
  let oracle;
  let positionToken;
  let multiSig;

  beforeEach(async function () {
    [owner, alice, bob, carol, asset] = await ethers.getSigners();

    oracle = await ethers.deployContract("MockOracle");
    positionToken = await ethers.deployContract("PositionToken", [
      "Lender Position",
      "LPOS",
    ]);
    multiSig = await ethers.deployContract("ThresholdMultiSig", [
      [owner.address, alice.address, bob.address],
      2,
    ]);

    await oracle.transferOwnership(await multiSig.getAddress());
    await positionToken.transferOwnership(await multiSig.getAddress());
  });

  async function approveAndExecute(
    signerA,
    signerB,
    target,
    value,
    data,
    nonce,
  ) {
    await multiSig
      .connect(signerA)
      .approveTransaction(target, value, data, nonce);
    await multiSig
      .connect(signerB)
      .approveTransaction(target, value, data, nonce);
    return multiSig
      .connect(signerB)
      .executeTransaction(target, value, data, nonce);
  }

  it("stores owners and threshold", async function () {
    expect(await multiSig.threshold()).to.equal(2n);
    expect(await multiSig.ownerCount()).to.equal(3n);
    expect(await multiSig.getOwner(0)).to.equal(owner.address);
    expect(await multiSig.getOwner(1)).to.equal(alice.address);
    expect(await multiSig.getOwner(2)).to.equal(bob.address);
    expect(await multiSig.isOwner(carol.address)).to.equal(false);
  });

  it("executes an admin call only after enough approvals", async function () {
    const oracleAddress = await oracle.getAddress();
    const data = oracle.interface.encodeFunctionData("setPrice", [
      asset.address,
      100_000_000,
    ]);
    const nonce = 1;

    await expect(
      oracle.setPrice(asset.address, 100_000_000),
    ).to.be.revertedWith("Not the owner");

    const txHash = await multiSig.getTransactionHash(
      oracleAddress,
      0,
      data,
      nonce,
    );

    await expect(
      multiSig.connect(owner).approveTransaction(oracleAddress, 0, data, nonce),
    )
      .to.emit(multiSig, "TransactionApproved")
      .withArgs(txHash, owner.address, 1);

    await expect(
      multiSig.connect(owner).executeTransaction(oracleAddress, 0, data, nonce),
    ).to.be.revertedWith("Not enough approvals");

    await multiSig
      .connect(alice)
      .approveTransaction(oracleAddress, 0, data, nonce);

    await expect(
      multiSig.connect(bob).executeTransaction(oracleAddress, 0, data, nonce),
    )
      .to.emit(multiSig, "TransactionExecuted")
      .withArgs(txHash, bob.address, oracleAddress);

    expect(await oracle.getPrice(asset.address)).to.equal(100_000_000n);
  });

  it("binds approvals to the exact target, calldata, value, chain, and nonce", async function () {
    const oracleAddress = await oracle.getAddress();
    const approvedData = oracle.interface.encodeFunctionData("setPrice", [
      asset.address,
      100_000_000,
    ]);
    const differentData = oracle.interface.encodeFunctionData("setPrice", [
      asset.address,
      200_000_000,
    ]);
    const nonce = 2;

    await multiSig
      .connect(owner)
      .approveTransaction(oracleAddress, 0, approvedData, nonce);
    await multiSig
      .connect(alice)
      .approveTransaction(oracleAddress, 0, approvedData, nonce);

    await expect(
      multiSig
        .connect(bob)
        .executeTransaction(oracleAddress, 0, differentData, nonce),
    ).to.be.revertedWith("Not enough approvals");

    await multiSig
      .connect(owner)
      .executeTransaction(oracleAddress, 0, approvedData, nonce);
    expect(await oracle.getPrice(asset.address)).to.equal(100_000_000n);
  });

  it("rejects duplicate approvals, non-owner approvals, and replay", async function () {
    const oracleAddress = await oracle.getAddress();
    const data = oracle.interface.encodeFunctionData("setPrice", [
      asset.address,
      100_000_000,
    ]);
    const nonce = 3;

    await multiSig
      .connect(owner)
      .approveTransaction(oracleAddress, 0, data, nonce);

    await expect(
      multiSig.connect(owner).approveTransaction(oracleAddress, 0, data, nonce),
    ).to.be.revertedWith("Already approved by this owner");
    await expect(
      multiSig.connect(carol).approveTransaction(oracleAddress, 0, data, nonce),
    ).to.be.revertedWith("Not an owner");

    await multiSig
      .connect(alice)
      .approveTransaction(oracleAddress, 0, data, nonce);
    await multiSig
      .connect(alice)
      .executeTransaction(oracleAddress, 0, data, nonce);

    await expect(
      multiSig.connect(bob).executeTransaction(oracleAddress, 0, data, nonce),
    ).to.be.revertedWith("Transaction already executed");
  });

  it("administers a PositionToken after ownership transfer", async function () {
    const tokenAddress = await positionToken.getAddress();
    const data = positionToken.interface.encodeFunctionData("setMinter", [
      carol.address,
      true,
    ]);

    await approveAndExecute(owner, alice, tokenAddress, 0, data, 4);

    expect(await positionToken.isMinter(carol.address)).to.equal(true);
    expect(await positionToken.owner()).to.equal(await multiSig.getAddress());
  });

  it("wraps a failed target call without marking it executed", async function () {
    const oracleAddress = await oracle.getAddress();
    const data = oracle.interface.encodeFunctionData("setPrice", [
      ethers.ZeroAddress,
      100_000_000,
    ]);
    const nonce = 5;

    await multiSig
      .connect(owner)
      .approveTransaction(oracleAddress, 0, data, nonce);
    await multiSig
      .connect(alice)
      .approveTransaction(oracleAddress, 0, data, nonce);

    await expect(
      multiSig.connect(owner).executeTransaction(oracleAddress, 0, data, nonce),
    ).to.be.revertedWith("Transaction execution failed");

    const txHash = await multiSig.getTransactionHash(
      oracleAddress,
      0,
      data,
      nonce,
    );
    expect(await multiSig.executed(txHash)).to.equal(false);
  });
});
