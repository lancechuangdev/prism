import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool excess refunds", function () {
  const INTEREST_RATE = 1_000_000;
  const MAX_LEND_SUPPLY = ethers.parseEther("100000");
  const COLLATERALIZATION_RATIO = 200_000_000;
  const LIQUIDATE_RATE = 20_000_000;
  const LEND_TOKEN_PRICE = 100_000_000;
  const COLLATERAL_TOKEN_PRICE = 5_000_000_000_000;

  let owner;
  let alice;
  let bob;
  let carol;
  let feeRecipient;
  let oracle;
  let lendToken;
  let collateralToken;
  let lenderPositionToken;
  let borrowerPositionToken;
  let swap;
  let pool;

  beforeEach(async function () {
    [owner, alice, bob, carol, feeRecipient] = await ethers.getSigners();

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

    await lendToken.setMinter(owner.address, true);
    await collateralToken.setMinter(owner.address, true);
    await lendToken.mint(alice.address, ethers.parseEther("100000"));
    await lendToken.mint(carol.address, ethers.parseEther("100000"));
    await collateralToken.mint(bob.address, ethers.parseEther("10"));
    await collateralToken.mint(carol.address, ethers.parseEther("10"));

    pool = await deployPool();

    await oracle.setPrice(await lendToken.getAddress(), LEND_TOKEN_PRICE);
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_TOKEN_PRICE,
    );
  });

  async function deployPool() {
    return ethers.deployContract("PrismPool", [
      await oracle.getAddress(),
      await swap.getAddress(),
      feeRecipient.address,
    ]);
  }

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

  async function createPoolAndApproveAll(targetPool = pool) {
    await targetPool.createPool(await buildCreateParams());
    const poolAddress = await targetPool.getAddress();

    for (const account of [alice, carol]) {
      await lendToken
        .connect(account)
        .approve(poolAddress, ethers.parseEther("100000"));
    }
    for (const account of [bob, carol]) {
      await collateralToken
        .connect(account)
        .approve(poolAddress, ethers.parseEther("10"));
    }
  }

  async function moveToSettleAndSettle(targetPool = pool) {
    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);
    await targetPool.settle(0);
  }

  it("refunds excess lender tokens after collateral-limited settlement", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("100000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    const refundAmount = ethers.parseEther("50000");
    await expect(pool.connect(alice).refundExcessLend(0))
      .to.emit(pool, "RefundLend")
      .withArgs(alice.address, 0, await lendToken.getAddress(), refundAmount);

    const aliceInfo = await pool.userLendInfo(alice.address, 0);
    expect(aliceInfo.refundAmount).to.equal(refundAmount);
    expect(aliceInfo.hasRefunded).to.equal(true);
    expect(await lendToken.balanceOf(alice.address)).to.equal(refundAmount);
    expect(await lendToken.balanceOf(await pool.getAddress())).to.equal(
      refundAmount,
    );
  });

  it("refunds excess lender tokens proportionally", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("60000"));
    await pool.connect(carol).depositLend(0, ethers.parseEther("40000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    await pool.connect(alice).refundExcessLend(0);
    await pool.connect(carol).refundExcessLend(0);

    const aliceInfo = await pool.userLendInfo(alice.address, 0);
    const carolInfo = await pool.userLendInfo(carol.address, 0);
    expect(aliceInfo.refundAmount).to.equal(ethers.parseEther("30000"));
    expect(carolInfo.refundAmount).to.equal(ethers.parseEther("20000"));
  });

  it("refunds excess collateral after lender-limited settlement", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    const refundAmount = ethers.parseEther("1");
    await expect(pool.connect(bob).refundExcessCollateral(0))
      .to.emit(pool, "RefundBorrow")
      .withArgs(
        bob.address,
        0,
        await collateralToken.getAddress(),
        refundAmount,
      );

    const bobInfo = await pool.userBorrowInfo(bob.address, 0);
    expect(bobInfo.refundAmount).to.equal(refundAmount);
    expect(bobInfo.hasRefunded).to.equal(true);
    expect(await collateralToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("9"),
    );
    expect(await collateralToken.balanceOf(await pool.getAddress())).to.equal(
      ethers.parseEther("1"),
    );
  });

  it("refunds excess collateral proportionally", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("1.2"));
    await pool.connect(carol).depositBorrow(0, ethers.parseEther("1.8"));
    await moveToSettleAndSettle();

    await pool.connect(bob).refundExcessCollateral(0);
    await pool.connect(carol).refundExcessCollateral(0);

    const bobInfo = await pool.userBorrowInfo(bob.address, 0);
    const carolInfo = await pool.userBorrowInfo(carol.address, 0);
    expect(bobInfo.refundAmount).to.equal(ethers.parseEther("0.8"));
    expect(carolInfo.refundAmount).to.equal(ethers.parseEther("1.2"));
  });

  it("rejects duplicate refunds, missing stakes, and fully matched sides", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    await expect(pool.connect(alice).refundExcessLend(0)).to.be.revertedWith(
      "No refund",
    );
    await expect(
      pool.connect(carol).refundExcessCollateral(0),
    ).to.be.revertedWith("No borrow stake");

    await pool.connect(bob).refundExcessCollateral(0);
    await expect(
      pool.connect(bob).refundExcessCollateral(0),
    ).to.be.revertedWith("Already refunded");
  });

  it("rejects refunds before settlement and from CANCELLED pools", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("100000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));

    await expect(pool.connect(alice).refundExcessLend(0)).to.be.revertedWith(
      "Invalid pool state",
    );

    const cancelledPool = await deployPool();
    await createPoolAndApproveAll(cancelledPool);
    await cancelledPool
      .connect(carol)
      .depositLend(0, ethers.parseEther("1000"));
    await moveToSettleAndSettle(cancelledPool);

    expect(await cancelledPool.getPoolState(0)).to.equal(4n);
    await expect(
      cancelledPool.connect(carol).refundExcessLend(0),
    ).to.be.revertedWith("Invalid pool state");
  });
});
