import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool depositBorrow", function () {
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
      "Mock USDT",
      "mUSDT",
    ]);
    collateralToken = await ethers.deployContract("PositionToken", [
      "Mock WBTC",
      "mWBTC",
    ]);
    lenderPositionToken = await ethers.deployContract("PositionToken", [
      "Lender Position USDT",
      "lpUSDT",
    ]);
    borrowerPositionToken = await ethers.deployContract("PositionToken", [
      "Borrower Position WBTC",
      "bpWBTC",
    ]);

    await collateralToken.setMinter(owner.address, true);
    await collateralToken.mint(alice.address, ethers.parseEther("10"));
    await collateralToken.mint(bob.address, ethers.parseEther("10"));

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

  it("accepts borrower collateral before settlement", async function () {
    const depositAmount = ethers.parseEther("2");
    const poolAddress = await pool.getAddress();

    await collateralToken.connect(alice).approve(poolAddress, depositAmount);

    await expect(pool.connect(alice).depositBorrow(0, depositAmount))
      .to.emit(pool, "DepositBorrow")
      .withArgs(
        alice.address,
        0,
        await collateralToken.getAddress(),
        depositAmount,
      );

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userBorrowInfo(alice.address, 0);

    expect(createdPool.totalCollateralDeposited).to.equal(depositAmount);
    expect(createdPool.totalLendDeposited).to.equal(0n);
    expect(aliceInfo.stakeAmount).to.equal(depositAmount);
    expect(aliceInfo.refundAmount).to.equal(0n);
    expect(aliceInfo.hasRefunded).to.equal(false);
    expect(aliceInfo.hasClaimed).to.equal(false);
    expect(await collateralToken.balanceOf(alice.address)).to.equal(
      ethers.parseEther("8"),
    );
    expect(await collateralToken.balanceOf(poolAddress)).to.equal(
      depositAmount,
    );
  });

  it("aggregates multiple collateral deposits from one borrower", async function () {
    const poolAddress = await pool.getAddress();

    await collateralToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("5"));
    await pool.connect(alice).depositBorrow(0, ethers.parseEther("2"));
    await pool.connect(alice).depositBorrow(0, ethers.parseEther("3"));

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userBorrowInfo(alice.address, 0);

    expect(createdPool.totalCollateralDeposited).to.equal(
      ethers.parseEther("5"),
    );
    expect(aliceInfo.stakeAmount).to.equal(ethers.parseEther("5"));
    expect(await collateralToken.balanceOf(poolAddress)).to.equal(
      ethers.parseEther("5"),
    );
  });

  it("tracks collateral deposits from different borrowers separately", async function () {
    const poolAddress = await pool.getAddress();

    await collateralToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("2"));
    await collateralToken
      .connect(bob)
      .approve(poolAddress, ethers.parseEther("3"));

    await pool.connect(alice).depositBorrow(0, ethers.parseEther("2"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("3"));

    const createdPool = await pool.getPool(0);
    const aliceInfo = await pool.userBorrowInfo(alice.address, 0);
    const bobInfo = await pool.userBorrowInfo(bob.address, 0);

    expect(createdPool.totalCollateralDeposited).to.equal(
      ethers.parseEther("5"),
    );
    expect(aliceInfo.stakeAmount).to.equal(ethers.parseEther("2"));
    expect(bobInfo.stakeAmount).to.equal(ethers.parseEther("3"));
  });

  it("rejects collateral deposits below the minimum", async function () {
    const poolAddress = await pool.getAddress();
    await collateralToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("0.9"));

    await expect(
      pool.connect(alice).depositBorrow(0, ethers.parseEther("0.9")),
    ).to.be.revertedWith("Amount below minimum borrow amount");
  });

  it("rejects collateral deposits without enough approval", async function () {
    await expect(
      pool.connect(alice).depositBorrow(0, ethers.parseEther("2")),
    ).to.be.revertedWithCustomError(
      collateralToken,
      "ERC20InsufficientAllowance",
    );
  });

  it("rejects collateral deposits after settlement time", async function () {
    const poolAddress = await pool.getAddress();
    await collateralToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("2"));

    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);

    await expect(
      pool.connect(alice).depositBorrow(0, ethers.parseEther("2")),
    ).to.be.revertedWith("Operation not allowed after settle time");
  });

  it("rejects collateral deposits into pools that do not exist", async function () {
    await expect(
      pool.connect(alice).depositBorrow(1, ethers.parseEther("2")),
    ).to.be.revertedWith("Invalid pool ID");
  });
});
