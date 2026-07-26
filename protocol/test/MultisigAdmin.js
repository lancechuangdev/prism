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

  it("adds an owner through multisig approval and execution", async function () {
    const multiSigAddress = await multiSig.getAddress();
    const data = multiSig.interface.encodeFunctionData("addOwner", [
      carol.address,
    ]);

    await expect(approveAndExecute(owner, alice, multiSigAddress, 0, data, 5))
      .to.emit(multiSig, "OwnerAdded")
      .withArgs(carol.address);

    expect(await multiSig.isOwner(carol.address)).to.equal(true);
    expect(await multiSig.ownerCount()).to.equal(4n);
    expect(await multiSig.getOwner(3)).to.equal(carol.address);
  });

  it("removes an owner through multisig approval and execution", async function () {
    const multiSigAddress = await multiSig.getAddress();
    const data = multiSig.interface.encodeFunctionData("removeOwner", [
      bob.address,
    ]);

    await expect(approveAndExecute(owner, alice, multiSigAddress, 0, data, 6))
      .to.emit(multiSig, "OwnerRemoved")
      .withArgs(bob.address);

    expect(await multiSig.isOwner(bob.address)).to.equal(false);
    expect(await multiSig.ownerCount()).to.equal(2n);
  });

  it("replaces an owner through multisig approval and execution", async function () {
    const multiSigAddress = await multiSig.getAddress();
    const data = multiSig.interface.encodeFunctionData("replaceOwner", [
      bob.address,
      carol.address,
    ]);

    await expect(approveAndExecute(owner, alice, multiSigAddress, 0, data, 7))
      .to.emit(multiSig, "OwnerReplaced")
      .withArgs(bob.address, carol.address);

    expect(await multiSig.isOwner(bob.address)).to.equal(false);
    expect(await multiSig.isOwner(carol.address)).to.equal(true);
    expect(await multiSig.ownerCount()).to.equal(3n);
  });

  it("changes the threshold through multisig approval and execution", async function () {
    const multiSigAddress = await multiSig.getAddress();
    const data = multiSig.interface.encodeFunctionData("changeThreshold", [3]);

    await expect(approveAndExecute(owner, alice, multiSigAddress, 0, data, 8))
      .to.emit(multiSig, "ThresholdChanged")
      .withArgs(2n, 3n);

    expect(await multiSig.threshold()).to.equal(3n);
  });

  it("prevents owners from bypassing multisig administration", async function () {
    await expect(
      multiSig.connect(owner).addOwner(carol.address),
    ).to.be.revertedWith("Only multisig");
    await expect(
      multiSig.connect(owner).removeOwner(bob.address),
    ).to.be.revertedWith("Only multisig");
    await expect(
      multiSig.connect(owner).replaceOwner(bob.address, carol.address),
    ).to.be.revertedWith("Only multisig");
    await expect(multiSig.connect(owner).changeThreshold(1)).to.be.revertedWith(
      "Only multisig",
    );
  });

  it("keeps the threshold valid when removing owners", async function () {
    const multiSigAddress = await multiSig.getAddress();
    const thresholdData = multiSig.interface.encodeFunctionData(
      "changeThreshold",
      [3],
    );
    await approveAndExecute(owner, alice, multiSigAddress, 0, thresholdData, 9);

    const removeData = multiSig.interface.encodeFunctionData("removeOwner", [
      bob.address,
    ]);
    await multiSig
      .connect(owner)
      .approveTransaction(multiSigAddress, 0, removeData, 10);
    await multiSig
      .connect(alice)
      .approveTransaction(multiSigAddress, 0, removeData, 10);
    await multiSig
      .connect(bob)
      .approveTransaction(multiSigAddress, 0, removeData, 10);

    await expect(
      multiSig
        .connect(owner)
        .executeTransaction(multiSigAddress, 0, removeData, 10),
    ).to.be.revertedWith("Transaction execution failed");
    expect(await multiSig.isOwner(bob.address)).to.equal(true);
  });

  it("invalidates pending approvals after the owner configuration changes", async function () {
    const oracleAddress = await oracle.getAddress();
    const oracleData = oracle.interface.encodeFunctionData("setPrice", [
      asset.address,
      100_000_000,
    ]);
    const oldHash = await multiSig.getTransactionHash(
      oracleAddress,
      0,
      oracleData,
      12,
    );
    await multiSig
      .connect(owner)
      .approveTransaction(oracleAddress, 0, oracleData, 12);
    await multiSig
      .connect(alice)
      .approveTransaction(oracleAddress, 0, oracleData, 12);

    const multiSigAddress = await multiSig.getAddress();
    const addOwnerData = multiSig.interface.encodeFunctionData("addOwner", [
      carol.address,
    ]);
    await approveAndExecute(owner, alice, multiSigAddress, 0, addOwnerData, 13);

    expect(await multiSig.configurationVersion()).to.equal(1n);
    expect(await multiSig.transactionConfigurationVersion(oldHash)).to.equal(
      1n,
    );
    await expect(
      multiSig
        .connect(owner)
        .executeTransaction(oracleAddress, 0, oracleData, 12),
    ).to.be.revertedWith("Not enough approvals");
  });

  it("wraps a failed target call without marking it executed", async function () {
    const oracleAddress = await oracle.getAddress();
    const data = oracle.interface.encodeFunctionData("setPrice", [
      ethers.ZeroAddress,
      100_000_000,
    ]);
    const nonce = 11;

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
