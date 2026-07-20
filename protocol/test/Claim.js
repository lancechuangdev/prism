import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("PrismPool position claims", function () {
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
    await lenderPositionToken.setMinter(await pool.getAddress(), true);
    await borrowerPositionToken.setMinter(await pool.getAddress(), true);

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

  it("lets lenders claim position tokens for matched lend amounts", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("60000"));
    await pool.connect(carol).depositLend(0, ethers.parseEther("40000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    await expect(pool.connect(alice).claimLenderPosition(0))
      .to.emit(pool, "ClaimLend")
      .withArgs(
        alice.address,
        0,
        await lenderPositionToken.getAddress(),
        ethers.parseEther("30000"),
      );
    await pool.connect(carol).claimLenderPosition(0);

    const aliceInfo = await pool.userLendInfo(alice.address, 0);
    const carolInfo = await pool.userLendInfo(carol.address, 0);
    expect(await lenderPositionToken.balanceOf(alice.address)).to.equal(
      ethers.parseEther("30000"),
    );
    expect(await lenderPositionToken.balanceOf(carol.address)).to.equal(
      ethers.parseEther("20000"),
    );
    expect(await lenderPositionToken.totalSupply()).to.equal(
      ethers.parseEther("50000"),
    );
    expect(aliceInfo.hasClaimed).to.equal(true);
    expect(carolInfo.hasClaimed).to.equal(true);
  });

  it("lets borrowers claim position tokens and matched loans", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("1.2"));
    await pool.connect(carol).depositBorrow(0, ethers.parseEther("1.8"));
    await moveToSettleAndSettle();

    await expect(pool.connect(bob).claimBorrowerPositionAndLoan(0))
      .to.emit(pool, "ClaimBorrow")
      .withArgs(
        bob.address,
        0,
        await borrowerPositionToken.getAddress(),
        ethers.parseEther("0.4"),
        ethers.parseEther("10000"),
      );
    await pool.connect(carol).claimBorrowerPositionAndLoan(0);

    const bobInfo = await pool.userBorrowInfo(bob.address, 0);
    const carolInfo = await pool.userBorrowInfo(carol.address, 0);
    expect(await borrowerPositionToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("0.4"),
    );
    expect(await borrowerPositionToken.balanceOf(carol.address)).to.equal(
      ethers.parseEther("0.6"),
    );
    expect(await borrowerPositionToken.totalSupply()).to.equal(
      ethers.parseEther("1"),
    );
    expect(await lendToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("10000"),
    );
    expect(await lendToken.balanceOf(carol.address)).to.equal(
      ethers.parseEther("115000"),
    );
    expect(bobInfo.hasClaimed).to.equal(true);
    expect(carolInfo.hasClaimed).to.equal(true);
  });

  it("allows refunds and claims in either order", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("100000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle();

    await pool.connect(alice).refundExcessLend(0);
    await pool.connect(alice).claimLenderPosition(0);
    await pool.connect(bob).claimBorrowerPositionAndLoan(0);

    expect(await lendToken.balanceOf(alice.address)).to.equal(
      ethers.parseEther("50000"),
    );
    expect(await lendToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("50000"),
    );
    expect(await lenderPositionToken.balanceOf(alice.address)).to.equal(
      ethers.parseEther("50000"),
    );
    expect(await borrowerPositionToken.balanceOf(bob.address)).to.equal(
      ethers.parseEther("2"),
    );
  });

  it("rejects claims before settlement, without stakes, and after claiming", async function () {
    await createPoolAndApproveAll();
    await pool.connect(alice).depositLend(0, ethers.parseEther("25000"));
    await pool.connect(bob).depositBorrow(0, ethers.parseEther("2"));

    await expect(pool.connect(alice).claimLenderPosition(0)).to.be.revertedWith(
      "Invalid pool state",
    );
    await expect(
      pool.connect(bob).claimBorrowerPositionAndLoan(0),
    ).to.be.revertedWith("Invalid pool state");

    await moveToSettleAndSettle();

    await expect(pool.connect(carol).claimLenderPosition(0)).to.be.revertedWith(
      "No lender position",
    );
    await expect(
      pool.connect(carol).claimBorrowerPositionAndLoan(0),
    ).to.be.revertedWith("No borrower position");

    await pool.connect(alice).claimLenderPosition(0);
    await pool.connect(bob).claimBorrowerPositionAndLoan(0);

    await expect(pool.connect(alice).claimLenderPosition(0)).to.be.revertedWith(
      "Already claimed",
    );
    await expect(
      pool.connect(bob).claimBorrowerPositionAndLoan(0),
    ).to.be.revertedWith("Already claimed");
  });

  it("requires the pool to be a position-token minter", async function () {
    const poolWithoutMinterRole = await deployPool();
    await createPoolAndApproveAll(poolWithoutMinterRole);
    await poolWithoutMinterRole
      .connect(alice)
      .depositLend(0, ethers.parseEther("25000"));
    await poolWithoutMinterRole
      .connect(bob)
      .depositBorrow(0, ethers.parseEther("2"));
    await moveToSettleAndSettle(poolWithoutMinterRole);

    await expect(
      poolWithoutMinterRole.connect(alice).claimLenderPosition(0),
    ).to.be.revertedWith("caller is not minter");
    await expect(
      poolWithoutMinterRole.connect(bob).claimBorrowerPositionAndLoan(0),
    ).to.be.revertedWith("caller is not minter");
  });
});
