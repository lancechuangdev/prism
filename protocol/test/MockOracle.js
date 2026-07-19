import assert from "node:assert/strict";
import { network } from "hardhat";

const { ethers } = await network.create();

async function expectRevert(action, message) {
  await assert.rejects(action, (error) => error.message.includes(message));
}

describe("Mock Oracle", function () {
  let busd;
  let btc;
  let alice;
  let oracle;

  beforeEach(async function () {
    [, alice, busd, btc] = await ethers.getSigners();

    const mockOracle = await ethers.getContractFactory("MockOracle");
    oracle = await mockOracle.deploy();
    await oracle.waitForDeployment();
  });

  describe("MockOracle", function () {
    it("lets the owner set and read one asset price", async function () {
      await oracle.setPrice(busd.address, 100000000);

      assert.equal(await oracle.getPrice(busd.address), 100000000n);
      await expectRevert(
        oracle.getPrice(btc.address),
        "Price not set for this token",
      );
    });

    it("sets batch prices by asset address", async function () {
      await oracle.setPrices(
        [busd.address, btc.address],
        [100000000, 5000000000000],
      );

      assert.equal(await oracle.getPrice(busd.address), 100000000n);
      assert.equal(await oracle.getPrice(btc.address), 5000000000000n);

      const prices = await oracle.getPrices([busd.address, btc.address]);
      assert.equal(prices[0], 100000000n);
      assert.equal(prices[1], 5000000000000n);
    });

    it("blocks non-owners and invalid prices", async function () {
      await expectRevert(
        oracle.connect(alice).setPrice(busd.address, 100000000),
        "Not the owner",
      );
      await expectRevert(
        oracle.setPrice(ethers.ZeroAddress, 100000000),
        "Invalid token address",
      );
      await expectRevert(
        oracle.setPrice(busd.address, 0),
        "Price must be positive",
      );
      await expectRevert(
        oracle.setPrices([busd.address], [100000000, 5000000000000]),
        "Mismatched array lengths",
      );
    });
  });
});
