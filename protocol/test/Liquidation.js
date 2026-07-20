import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool liquidation", function () {
  const INTEREST_RATE = 1_000_000;
  const MAX_LEND_SUPPLY = ethers.parseEther("100000");
  const COLLATERALIZATION_RATIO = 200_000_000;
  const LIQUIDATE_RATE = 20_000_000;
  const LEND_TOKEN_PRICE = 100_000_000;
  const COLLATERAL_TOKEN_PRICE = 5_000_000_000_000;
  const COLLATERAL_CRASH_PRICE = 1_000_000_000_000;
  const COLLATERAL_TO_LEND_RATE = ethers.parseEther("50000");

  let owner;
  let alice;
  let bob;
  let feeRecipient;
  let oracle;
  let swap;
  let lendToken;
  let collateralToken;
  let lenderPositionToken;
  let borrowerPositionToken;
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
    await lendToken.mint(await swap.getAddress(), ethers.parseEther("100000"));
    await collateralToken.mint(bob.address, ethers.parseEther("10"));

    pool = await ethers.deployContract("PrismPool", [
      await oracle.getAddress(),
      await swap.getAddress(),
      feeRecipient.address,
    ]);

    await lenderPositionToken.setMinter(await pool.getAddress(), true);
    await borrowerPositionToken.setMinter(await pool.getAddress(), true);
    await oracle.setPrice(await lendToken.getAddress(), LEND_TOKEN_PRICE);
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_TOKEN_PRICE,
    );
    await swap.setRate(
      await collateralToken.getAddress(),
      await lendToken.getAddress(),
      COLLATERAL_TO_LEND_RATE,
    );
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

  async function createActivePool() {
    await pool.createPool(await buildCreateParams());
    const poolAddress = await pool.getAddress();

    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("25000"));
    await collateralToken
      .connect(bob)
      .approve(poolAddress, ethers.parseEther("2"));
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));

    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);
    await pool.settle(0);
  }

  async function preparePositions() {
    await pool.connect(bob).refundExcessCollateral(0);
    await pool.connect(alice).claimLenderPosition(0);
    await pool.connect(bob).claimBorrowerPositionAndLoan(0);
  }

  it("becomes undercollateralized after the collateral price crashes", async function () {
    await pool.createPool(await buildCreateParams());
    expect(await pool.isUndercollateralized(0)).to.equal(false);

    const poolAddress = await pool.getAddress();
    await lendToken
      .connect(alice)
      .approve(poolAddress, ethers.parseEther("25000"));
    await collateralToken
      .connect(bob)
      .approve(poolAddress, ethers.parseEther("2"));
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await ethers.provider.send("evm_increaseTime", [3601]);
    await ethers.provider.send("evm_mine", []);
    await pool.settle(0);

    expect(await pool.isUndercollateralized(0)).to.equal(false);

    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );
    expect(await pool.isUndercollateralized(0)).to.equal(true);
  });

  it("sells only enough collateral when the DEX can cover repayment", async function () {
    await createActivePool();
    await preparePositions();
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );

    const requiredRepayment = await pool.getRequiredRepayment(0);
    const collateralToSell = await swap.getAmountIn(
      await collateralToken.getAddress(),
      await lendToken.getAddress(),
      requiredRepayment,
    );
    const remainingCollateral = ethers.parseEther("1") - collateralToSell;

    await expect(pool.liquidate(0, collateralToSell))
      .to.emit(pool, "PoolLiquidated")
      .withArgs(
        0,
        await swap.getAddress(),
        collateralToSell,
        requiredRepayment,
        remainingCollateral,
      );

    const data = await pool.getPoolData(0);
    expect(await pool.getPoolState(0)).to.equal(3n);
    expect(data.liquidationAmountLend).to.equal(requiredRepayment);
    expect(data.liquidationAmountBorrow).to.equal(remainingCollateral);
  });

  it("sells all matched collateral when it cannot fully cover repayment", async function () {
    await createActivePool();
    await preparePositions();
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );
    await swap.setRate(
      await collateralToken.getAddress(),
      await lendToken.getAddress(),
      ethers.parseEther("0.5"),
    );

    await expect(pool.liquidate(0, ethers.parseEther("1")))
      .to.emit(pool, "PoolLiquidated")
      .withArgs(
        0,
        await swap.getAddress(),
        ethers.parseEther("1"),
        ethers.parseEther("0.5"),
        0,
      );

    const data = await pool.getPoolData(0);
    expect(await pool.getPoolState(0)).to.equal(3n);
    expect(data.liquidationAmountLend).to.equal(ethers.parseEther("0.5"));
    expect(data.liquidationAmountBorrow).to.equal(0n);
  });

  it("lets borrower-position holders withdraw remaining collateral", async function () {
    await createActivePool();
    await preparePositions();
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );

    const requiredRepayment = await pool.getRequiredRepayment(0);
    const collateralToSell = await swap.getAmountIn(
      await collateralToken.getAddress(),
      await lendToken.getAddress(),
      requiredRepayment,
    );
    await pool.liquidate(0, collateralToSell);

    const borrowerPosition = await borrowerPositionToken.balanceOf(bob.address);
    await pool.connect(bob).withdrawBorrow(0, borrowerPosition);

    expect(await borrowerPositionToken.balanceOf(bob.address)).to.equal(0n);
    expect(await collateralToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("9") + (ethers.parseEther("1") - collateralToSell),
    );
  });

  it("rejects liquidation while healthy or when caller is not owner", async function () {
    await createActivePool();

    await expect(pool.liquidate(0, ethers.parseEther("1"))).to.be.revertedWith(
      "Pool is sufficiently collateralized",
    );

    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );
    await expect(
      pool.connect(alice).liquidate(0, ethers.parseEther("1")),
    ).to.be.revertedWith("Not the owner");
  });

  it("rejects liquidation when maximum collateral is too low", async function () {
    await createActivePool();
    await oracle.setPrice(
      await collateralToken.getAddress(),
      COLLATERAL_CRASH_PRICE,
    );

    const requiredRepayment = await pool.getRequiredRepayment(0);
    const collateralToSell = await swap.getAmountIn(
      await collateralToken.getAddress(),
      await lendToken.getAddress(),
      requiredRepayment,
    );

    await expect(pool.liquidate(0, collateralToSell - 1n)).to.be.revertedWith(
      "Dex slippage too high",
    );
  });
});
