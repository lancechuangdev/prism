import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("Create pool", function () {
  const INTEREST_RATE = 1_000_000;
  const MAX_SUPPLY = ethers.parseEther("100000");
  const COLLATERALIZATION_RATIO = 200_000_000;
  const LIQUIDATE_THRESHOLD = 20_000_000;

  let owner;
  let alice;
  let feeRecipient;
  let lendToken;
  let borrowToken;
  let oracle;
  let swap;
  let lenderPositionToken;
  let borrowerPositionToken;
  let pool;

  beforeEach(async function () {
    [owner, alice, feeRecipient, lendToken, borrowToken] =
      await ethers.getSigners();

    const mockOracle = await ethers.getContractFactory("MockOracle");
    oracle = await mockOracle.deploy();
    await oracle.waitForDeployment();

    const positionToken = await ethers.getContractFactory("PositionToken");
    lenderPositionToken = await positionToken.deploy("Lender BUSD", "lBUSD");
    borrowerPositionToken = await positionToken.deploy("Borrower BTC", "bBTC");
    await lenderPositionToken.waitForDeployment();
    await borrowerPositionToken.waitForDeployment();

    swap = await ethers.deployContract("FixedRateSwap");

    const prismPool = await ethers.getContractFactory("PrismPool");
    pool = await prismPool.deploy(
      await oracle.getAddress(),
      await swap.getAddress(),
      feeRecipient.address,
    );
    await pool.waitForDeployment();
  });

  async function buildCreateParams(overrides = {}) {
    const latestBlock = await ethers.provider.getBlock("latest");
    const settleTime = latestBlock.timestamp + 3600;
    const params = {
      settleTime,
      maturityTime: settleTime + 7 * 24 * 60 * 60,
      interestRate: INTEREST_RATE,
      maxLendSupply: MAX_SUPPLY,
      collateralizationRatio: COLLATERALIZATION_RATIO,
      lendToken: lendToken.address,
      collateralToken: borrowToken.address,
      lenderPositionToken: await lenderPositionToken.getAddress(),
      borrowerPositionToken: await borrowerPositionToken.getAddress(),
      liquidateRate: LIQUIDATE_THRESHOLD,
    };

    return { ...params, ...overrides };
  }

  it("deploys with owner, oracle, fee address, and no pools", async function () {
    expect(await pool.owner()).to.equal(owner.address);
    expect(await pool.oracle()).to.equal(await oracle.getAddress());
    expect(await pool.dexSwap()).to.equal(await swap.getAddress());
    expect(await pool.feeAddress()).to.equal(feeRecipient.address);
    expect(await pool.globalPaused()).to.equal(false);
    expect(await pool.poolCount()).to.equal(0n);
  });

  it("creates a pool in FUNDING state", async function () {
    const params = await buildCreateParams();

    await pool.createPool(params);

    expect(await pool.poolCount()).to.equal(1n);
    expect(await pool.getPoolState(0)).to.equal(0n);
    expect(await pool.isBeforeSettleTime(0)).to.equal(true);

    const createdPool = await pool.getPool(0);
    expect(createdPool.settleTime).to.equal(BigInt(params.settleTime));
    expect(createdPool.maturityTime).to.equal(BigInt(params.maturityTime));
    expect(createdPool.interestRate).to.equal(BigInt(params.interestRate));
    expect(createdPool.maxLendSupply).to.equal(params.maxLendSupply);
    expect(createdPool.totalLendDeposited).to.equal(0n);
    expect(createdPool.totalCollateralDeposited).to.equal(0n);
    expect(createdPool.collateralizationRatio).to.equal(
      BigInt(params.collateralizationRatio),
    );
    expect(createdPool.lendToken).to.equal(params.lendToken);
    expect(createdPool.collateralToken).to.equal(params.collateralToken);
    expect(createdPool.state).to.equal(0n);
    expect(createdPool.lenderPositionToken).to.equal(
      params.lenderPositionToken,
    );
    expect(createdPool.borrowerPositionToken).to.equal(
      params.borrowerPositionToken,
    );
    expect(createdPool.liquidateRate).to.equal(BigInt(params.liquidateRate));
  });

  it("creates multiple pools with increasing ids", async function () {
    await pool.createPool(await buildCreateParams());
    await pool.createPool(
      await buildCreateParams({
        maxLendSupply: ethers.parseEther("50000"),
      }),
    );

    expect(await pool.poolCount()).to.equal(2n);

    const secondPool = await pool.getPool(1);
    expect(secondPool.maxLendSupply).to.equal(ethers.parseEther("50000"));
  });

  it("blocks non-owner pool creation", async function () {
    await expect(
      pool.connect(alice).createPool(await buildCreateParams()),
    ).to.be.revertedWith("Not the owner");
  });

  it("validates required create pool fields", async function () {
    const params = await buildCreateParams();

    await expect(
      pool.createPool({ ...params, settleTime: 1 }),
    ).to.be.revertedWith("Settle time must be in the future");
    await expect(
      pool.createPool({ ...params, maturityTime: params.settleTime }),
    ).to.be.revertedWith("End time must be after settle time");
    await expect(
      pool.createPool({ ...params, maxLendSupply: 0 }),
    ).to.be.revertedWith("Max supply must be positive");
    await expect(
      pool.createPool({ ...params, interestRate: 0 }),
    ).to.be.revertedWith("Interest rate must be positive");
    await expect(
      pool.createPool({ ...params, collateralizationRatio: 0 }),
    ).to.be.revertedWith("Mortgage rate must be positive");
    await expect(
      pool.createPool({ ...params, liquidateRate: 0 }),
    ).to.be.revertedWith("Liquidate rate must be positive");
    await expect(
      pool.createPool({ ...params, lendToken: ethers.ZeroAddress }),
    ).to.be.revertedWith("Invalid lend token address");
    await expect(
      pool.createPool({ ...params, collateralToken: ethers.ZeroAddress }),
    ).to.be.revertedWith("Invalid borrow token address");
    await expect(
      pool.createPool({ ...params, collateralToken: params.lendToken }),
    ).to.be.revertedWith("Lend and borrow tokens must be different");
    await expect(
      pool.createPool({
        ...params,
        lenderPositionToken: ethers.ZeroAddress,
      }),
    ).to.be.revertedWith("Invalid sp token address");
    await expect(
      pool.createPool({
        ...params,
        borrowerPositionToken: ethers.ZeroAddress,
      }),
    ).to.be.revertedWith("Invalid jp token address");
    await expect(
      pool.createPool({
        ...params,
        borrowerPositionToken: params.lenderPositionToken,
      }),
    ).to.be.revertedWith("SP and JP tokens must be different");
  });

  it("rejects reads for pools that do not exist", async function () {
    await expect(pool.getPool(0)).to.be.revertedWith("Invalid pool ID");
    await expect(pool.getPoolData(0)).to.be.revertedWith("Invalid pool ID");
    await expect(pool.getPoolState(0)).to.be.revertedWith("Invalid pool ID");
    await expect(pool.isBeforeSettleTime(0)).to.be.revertedWith(
      "Invalid pool ID",
    );
  });
});
