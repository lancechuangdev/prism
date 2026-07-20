import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("FixedRateSwap", function () {
  const rate = ethers.parseUnits("2", 18);

  let owner;
  let trader;
  let recipient;
  let tokenIn;
  let tokenOut;
  let swap;

  beforeEach(async function () {
    [owner, trader, recipient] = await ethers.getSigners();

    tokenIn = await ethers.deployContract("PositionToken", [
      "Input Token",
      "TIN",
    ]);
    tokenOut = await ethers.deployContract("PositionToken", [
      "Output Token",
      "TOUT",
    ]);
    swap = await ethers.deployContract("FixedRateSwap");

    await tokenIn.setMinter(owner.address, true);
    await tokenOut.setMinter(owner.address, true);

    await tokenIn.mint(trader.address, ethers.parseEther("1000"));
    await tokenOut.mint(await swap.getAddress(), ethers.parseEther("10000"));

    await swap.setRate(
      await tokenIn.getAddress(),
      await tokenOut.getAddress(),
      rate,
    );
  });

  describe("rates and quotes", function () {
    it("sets a rate and emits RateChanged", async function () {
      const newRate = ethers.parseUnits("1.5", 18);

      await expect(
        swap.setRate(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          newRate,
        ),
      )
        .to.emit(swap, "RateChanged")
        .withArgs(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          newRate,
        );

      expect(
        await swap.rateOutPerTokenIn(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
        ),
      ).to.equal(newRate);
    });

    it("only lets the owner set valid rates", async function () {
      await expect(
        swap
          .connect(trader)
          .setRate(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            rate,
          ),
      ).to.be.revertedWith("Not the owner");

      await expect(
        swap.setRate(ethers.ZeroAddress, await tokenOut.getAddress(), rate),
      ).to.be.revertedWith("Invalid token address");

      await expect(
        swap.setRate(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          0,
        ),
      ).to.be.revertedWith("Rate must be positive");
    });

    it("quotes exact-input and exact-output swaps", async function () {
      expect(
        await swap.getAmountOut(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          ethers.parseEther("10"),
        ),
      ).to.equal(ethers.parseEther("20"));

      expect(
        await swap.getAmountIn(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          ethers.parseEther("20"),
        ),
      ).to.equal(ethers.parseEther("10"));
    });

    it("rounds exact-output input amounts up", async function () {
      await swap.setRate(
        await tokenIn.getAddress(),
        await tokenOut.getAddress(),
        ethers.parseUnits("3", 18),
      );

      expect(
        await swap.getAmountIn(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          10,
        ),
      ).to.equal(4n);
    });

    it("rejects quotes for pairs without a rate", async function () {
      await expect(
        swap.getAmountOut(
          await tokenOut.getAddress(),
          await tokenIn.getAddress(),
          1,
        ),
      ).to.be.revertedWith("Rate not set for this pair");
    });
  });

  describe("swaps", function () {
    it("executes an exact-input swap and updates balances", async function () {
      const amountIn = ethers.parseEther("10");
      const amountOut = ethers.parseEther("20");
      await tokenIn.connect(trader).approve(await swap.getAddress(), amountIn);

      await expect(
        swap
          .connect(trader)
          .swapExactTokensForTokens(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            amountIn,
            amountOut,
            recipient.address,
          ),
      )
        .to.emit(swap, "Swap")
        .withArgs(
          trader.address,
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          amountIn,
          amountOut,
          recipient.address,
        );

      expect(await tokenIn.balanceOf(trader.address)).to.equal(
        ethers.parseEther("990"),
      );
      expect(await tokenIn.balanceOf(await swap.getAddress())).to.equal(
        amountIn,
      );
      expect(await tokenOut.balanceOf(recipient.address)).to.equal(amountOut);
    });

    it("executes an exact-output swap and charges the quoted input", async function () {
      const amountOut = ethers.parseEther("25");
      const amountIn = ethers.parseEther("12.5");
      await tokenIn.connect(trader).approve(await swap.getAddress(), amountIn);

      await swap
        .connect(trader)
        .swapTokensForExactTokens(
          await tokenIn.getAddress(),
          await tokenOut.getAddress(),
          amountOut,
          amountIn,
          recipient.address,
        );

      expect(await tokenIn.balanceOf(await swap.getAddress())).to.equal(
        amountIn,
      );
      expect(await tokenOut.balanceOf(recipient.address)).to.equal(amountOut);
    });

    it("enforces exact-input slippage and recipient validation", async function () {
      const amountIn = ethers.parseEther("10");
      await tokenIn.connect(trader).approve(await swap.getAddress(), amountIn);

      await expect(
        swap
          .connect(trader)
          .swapExactTokensForTokens(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            amountIn,
            ethers.parseEther("21"),
            recipient.address,
          ),
      ).to.be.revertedWith("Insufficient output amount");

      await expect(
        swap
          .connect(trader)
          .swapExactTokensForTokens(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            amountIn,
            0,
            ethers.ZeroAddress,
          ),
      ).to.be.revertedWith("Invalid recipient");
    });

    it("enforces the exact-output maximum input", async function () {
      await expect(
        swap
          .connect(trader)
          .swapTokensForExactTokens(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            ethers.parseEther("20"),
            ethers.parseEther("9"),
            recipient.address,
          ),
      ).to.be.revertedWith("Excessive input amount");
    });

    it("requires the trader to approve tokenIn", async function () {
      await expect(
        swap
          .connect(trader)
          .swapExactTokensForTokens(
            await tokenIn.getAddress(),
            await tokenOut.getAddress(),
            ethers.parseEther("10"),
            0,
            recipient.address,
          ),
      ).to.be.revertedWithCustomError(tokenIn, "ERC20InsufficientAllowance");
    });
  });
});
