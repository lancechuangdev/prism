import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("UniswapV3SwapAdapter", function () {
  const fee = 3000;
  const rate = ethers.parseUnits("2", 18);
  let owner;
  let trader;
  let recipient;
  let tokenIn;
  let tokenOut;
  let uniswap;
  let adapter;

  beforeEach(async function () {
    [owner, trader, recipient] = await ethers.getSigners();
    tokenIn = await ethers.deployContract("PositionToken", ["Input", "IN"]);
    tokenOut = await ethers.deployContract("PositionToken", ["Output", "OUT"]);
    uniswap = await ethers.deployContract("MockUniswapV3");
    adapter = await ethers.deployContract("UniswapV3SwapAdapter", [
      owner.address,
      await uniswap.getAddress(),
      await uniswap.getAddress(),
    ]);

    await tokenIn.setMinter(owner.address, true);
    await tokenOut.setMinter(owner.address, true);
    await tokenIn.mint(trader.address, ethers.parseEther("100"));
    await tokenOut.mint(await uniswap.getAddress(), ethers.parseEther("1000"));
    await uniswap.setRate(
      await tokenIn.getAddress(),
      await tokenOut.getAddress(),
      rate,
    );
    await adapter.setPoolFee(
      await tokenIn.getAddress(),
      await tokenOut.getAddress(),
      fee,
    );
  });

  it("uses the configured fee tier for exact-input and exact-output quotes", async function () {
    expect(
      await adapter.getAmountOut.staticCall(
        await tokenIn.getAddress(),
        await tokenOut.getAddress(),
        ethers.parseEther("5"),
      ),
    ).to.equal(ethers.parseEther("10"));
    expect(
      await adapter.getAmountIn.staticCall(
        await tokenIn.getAddress(),
        await tokenOut.getAddress(),
        ethers.parseEther("10"),
      ),
    ).to.equal(ethers.parseEther("5"));
  });

  it("executes exact-input swaps and enforces minimum output", async function () {
    const amountIn = ethers.parseEther("5");
    await tokenIn.connect(trader).approve(await adapter.getAddress(), amountIn);

    await adapter
      .connect(trader)
      .swapExactTokensForTokens(
        await tokenIn.getAddress(),
        await tokenOut.getAddress(),
        amountIn,
        ethers.parseEther("10"),
        recipient.address,
      );
    expect(await tokenOut.balanceOf(recipient.address)).to.equal(
      ethers.parseEther("10"),
    );
  });

  it("executes exact-output swaps and enforces maximum input", async function () {
    const amountOut = ethers.parseEther("10");
    const amountIn = ethers.parseEther("5");
    await tokenIn.connect(trader).approve(await adapter.getAddress(), amountIn);

    await adapter
      .connect(trader)
      .swapTokensForExactTokens(
        await tokenIn.getAddress(),
        await tokenOut.getAddress(),
        amountOut,
        amountIn,
        recipient.address,
      );
    expect(await tokenIn.balanceOf(trader.address)).to.equal(
      ethers.parseEther("95"),
    );

    await expect(
      adapter
        .connect(trader)
        .swapTokensForExactTokens(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          amountOut,
          amountIn - 1n,
          recipient.address,
        ),
    ).to.be.revertedWith("Excessive input amount");
  });

  it("restricts fee configuration and rejects unconfigured routes", async function () {
    await expect(
      adapter
        .connect(trader)
        .setPoolFee(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          fee,
        ),
    ).to.be.revertedWith("Not the owner");
    await expect(
      adapter.getAmountOut.staticCall(
        await tokenOut.getAddress(),
        await tokenIn.getAddress(),
        1,
      ),
    ).to.be.revertedWith("Pool fee not configured");
  });
});
