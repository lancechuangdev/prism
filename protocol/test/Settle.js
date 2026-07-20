import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool settle", function () {
  const INTEREST_RATE = 1_000_000;
  const MAX_LEND_SUPPLY = ethers.parseEther("100000");
  const COLLATERALIZATION_RATIO = 200_000_000; // 200_000_000 / RATE_SCALE = 2 = 200% means every $1 lent requires $2 of collateral
  const LIQUIDATE_RATE = 20_000_000;
  const LEND_TOKEN_PRICE = 100_000_000;
  const COLLATERAL_TOKEN_PRICE = 5_000_000_000_000;

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

    await lendToken.setMinter(owner.address, true);
    await collateralToken.setMinter(owner.address, true);
    await lendToken.mint(alice.address, ethers.parseEther("100000"));
    await collateralToken.mint(bob.address, ethers.parseEther("10"));

    pool = await deployPool(oracle);

    await oracle.setPrice(await lendToken.getAddress(), LEND_TOKEN_PRICE);
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_TOKEN_PRICE,
    );
  });

  async function deployPool(poolOracle) {
    return ethers.deployContract("PrismPool", [
      await poolOracle.getAddress(),
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

  async function createPoolAndApprove(targetPool = pool) {
    await targetPool.createPool(await buildCreateParams());
    const poolAddress = await targetPool.getAddress();
    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("100000"));
    await collateralToken
      .connect(bob)
      .approve(poolAddress, ethers.parseEther("10"));
  }

  async function moveToSettleTime() {
    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);
  }

  it("becomes ACTIVE when collateral limits lender demand", async function () {
    await createPoolAndApprove();
    await pool.connect(alice).depositLend(0, ethers.parseEther("100000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleTime();

    await expect(pool.settle(0))
      .to.emit(pool, "StateChanged")
      .withArgs(0, 0, 1);

    const createdPool = await pool.getPool(0);
    const data = await pool.getPoolData(0);

    expect(createdPool.state).to.equal(1n);
    expect(data.settleAmountLend).to.equal(ethers.parseEther("50000"));
    expect(data.settleAmountBorrow).to.equal(ethers.parseEther("2"));
  });

  it("becomes ACTIVE when lender demand limits collateral", async function () {
    await createPoolAndApprove();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleTime();

    await pool.settle(0);

    const data = await pool.getPoolData(0);
    expect(await pool.getPoolState(0)).to.equal(1n);
    expect(data.settleAmountLend).to.equal(ethers.parseEther("25000"));
    expect(data.settleAmountBorrow).to.equal(ethers.parseEther("1"));
  });

  it("becomes CANCELLED when either side is empty", async function () {
    await pool.createPool(await buildCreateParams());
    const poolAddress = await pool.getAddress();
    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("1000"));
    await pool.connect(alice).depositLend(0, ethers.parseEther("1000"));
    await moveToSettleTime();

    await expect(pool.settle(0))
      .to.emit(pool, "StateChanged")
      .withArgs(0, 0, 4);

    const data = await pool.getPoolData(0);
    expect(await pool.getPoolState(0)).to.equal(4n);
    expect(data.settleAmountLend).to.equal(0n);
    expect(data.settleAmountBorrow).to.equal(0n);
  });

  it("rejects early, non-owner, and repeated settlement", async function () {
    await createPoolAndApprove();
    await pool.connect(alice).depositLend(0, ethers.parseEther("1000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("1"));

    await expect(pool.settle(0)).to.be.revertedWith(
      "Operation allowed only after settle time",
    );

    await moveToSettleTime();
    await expect(pool.connect(alice).settle(0)).to.be.revertedWith(
      "Not the owner",
    );

    await pool.settle(0);

    await expect(pool.settle(0)).to.be.revertedWith("Invalid pool state");
  });

  it("rejects settlement when oracle prices are missing", async function () {
    const emptyOracle = await ethers.deployContract("MockOracle");
    const secondPool = await deployPool(emptyOracle);
    await createPoolAndApprove(secondPool);

    await secondPool.connect(alice).depositLend(0, ethers.parseEther("1000"));
    await secondPool.connect(bob).depositBorrow(0, ethers.parseEther("1"));
    await moveToSettleTime();

    await expect(secondPool.settle(0)).to.be.revertedWith(
      "Price not set for this token",
    );
  });
});
