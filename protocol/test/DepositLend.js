import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool depositLend", function () {
  const INTEREST_RATE = 1_000_000;
  const MAX_LEND_SUPPLY = ethers.parseEther("1000");
  const COLLATERALIZATION_RATIO = 200_000_000;
  const LIQUIDATE_RATE = 20_000_000;

  let owner;
  let alice;
  let bob;
  let feeRecipient;
  let oracle;
  let lendToken;
  let collateralToken;
  let lenderPositionToken;
  let borrowerPositionToken;
  let swap;
  let pool;

  beforeEach(async function () {
    [owner, alice, bob, feeRecipient] = await ethers.getSigners();

    oracle = await ethers.deployContract("MockOracle");
    swap = await ethers.deployContract("FixedRateSwap");

    lendToken = await ethers.deployContract("PositionToken", [
      "Mock BUSD",
      "mBUSD",
    ]);
    collateralToken = await ethers.deployContract("PositionToken", [
      "Mock BTC",
      "mBTC",
    ]);
    lenderPositionToken = await ethers.deployContract("PositionToken", [
      "Lender Position BUSD",
      "lpBUSD",
    ]);
    borrowerPositionToken = await ethers.deployContract("PositionToken", [
      "Borrower Position BTC",
      "bpBTC",
    ]);

    await lendToken.setMinter(owner.address, true);
    await lendToken.mint(alice.address, ethers.parseEther("1000"));
    await lendToken.mint(bob.address, ethers.parseEther("1000"));

    pool = await ethers.deployContract("PrismPool", [
      await oracle.getAddress(),
      await swap.getAddress(),
      feeRecipient.address,
    ]);

    await pool.createPool(await buildCreateParams());
  });

  async function buildCreateParams(overrides = {}) {
    const latestBlock = await ethers.provider.getBlock("latest");
    const settleTime = latestBlock.timestamp + 3600;

    return {
      settleTime,
      maturityTime: settleTime + 7 * 24 * 60 * 60,
      interestRate: INTEREST_RATE,
      maxLendSupply: MAX_LEND_SUPPLY,
      collateralizationRatio: COLLATERALIZATION_RATIO,
      lendToken: await lendToken.getAddress(),
      collateralToken: await collateralToken.getAddress(),
      lenderPositionToken: await lenderPositionToken.getAddress(),
      borrowerPositionToken: await borrowerPositionToken.getAddress(),
      liquidateRate: LIQUIDATE_RATE,
      ...overrides,
    };
  }

  it("accepts a lender deposit before settlement", async function () {
    const depositAmount = ethers.parseEther("200");
    const poolAddress = await pool.getAddress();

    await lendToken.connect(alice).approve(poolAddress, depositAmount);

    await expect(pool.connect(alice).depositLend(0, depositAmount))
      .to.emit(pool, "DepositLend")
      .withArgs(alice.address, 0, await lendToken.getAddress(), depositAmount);

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userLendInfo(alice.address, 0);

    expect(createdPool.totalLendDeposited).to.equal(depositAmount);
    expect(aliceInfo.stakeAmount).to.equal(depositAmount);
    expect(aliceInfo.refundAmount).to.equal(0n);
    expect(aliceInfo.hasRefunded).to.equal(false);
    expect(aliceInfo.hasClaimed).to.equal(false);
    expect(await lendToken.balanceOf(alice.address)).to.equal(
      ethers.parseEther("800"),
    );
    expect(await lendToken.balanceOf(poolAddress)).to.equal(depositAmount);
  });

  it("aggregates multiple deposits from one lender", async function () {
    const poolAddress = await pool.getAddress();

    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("500"));
    await pool.connect(alice).depositLend(0, ethers.parseEther("200"));
    await pool.connect(alice).depositLend(0, ethers.parseEther("300"));

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userLendInfo(alice.address, 0);

    expect(createdPool.totalLendDeposited).to.equal(ethers.parseEther("500"));
    expect(aliceInfo.stakeAmount).to.equal(ethers.parseEther("500"));
    expect(await lendToken.balanceOf(poolAddress)).to.equal(
      ethers.parseEther("500"),
    );
  });

  it("tracks deposits from different lenders separately", async function () {
    const poolAddress = await pool.getAddress();

    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("200"));
    await lendToken.connect(bob).approve(poolAddress, ethers.parseEther("300"));

    await pool.connect(alice).depositLend(0, ethers.parseEther("200"));
    await pool.connect(bob).depositLend(0, ethers.parseEther("300"));

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userLendInfo(alice.address, 0);
    const bobInfo = await pool.userLendInfo(bob.address, 0);

    expect(createdPool.totalLendDeposited).to.equal(ethers.parseEther("500"));
    expect(aliceInfo.stakeAmount).to.equal(ethers.parseEther("200"));
    expect(bobInfo.stakeAmount).to.equal(ethers.parseEther("300"));
  });

  it("rejects deposits below the minimum or above pool capacity", async function () {
    const poolAddress = await pool.getAddress();
    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("2000"));

    await expect(
      pool.connect(alice).depositLend(0, ethers.parseEther("99")),
    ).to.be.revertedWith("Amount below minimum lend amount");

    await pool.connect(alice).depositLend(0, ethers.parseEther("900"));

    await expect(
      pool.connect(alice).depositLend(0, ethers.parseEther("101")),
    ).to.be.revertedWith("Exceeds max supply");
  });

  it("rejects a zero deposit", async function () {
    await expect(pool.connect(alice).depositLend(0, 0)).to.be.revertedWith(
      "Invalid amount",
    );
  });

  it("rejects deposits without enough approval", async function () {
    await expect(
      pool.connect(alice).depositLend(0, ethers.parseEther("200")),
    ).to.be.revertedWithCustomError(lendToken, "ERC20InsufficientAllowance");
  });

  it("rejects deposits after settlement time", async function () {
    const poolAddress = await pool.getAddress();
    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("200"));

    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);

    await expect(
      pool.connect(alice).depositLend(0, ethers.parseEther("200")),
    ).to.be.revertedWith("Operation not allowed after settle time");
  });

  it("rejects deposits into pools that do not exist", async function () {
    await expect(
      pool.connect(alice).depositLend(1, ethers.parseEther("200")),
    ).to.be.revertedWith("Invalid pool ID");
  });
});
