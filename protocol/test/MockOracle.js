import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

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

  it("lets the owner set and read one asset price", async function () {
    await oracle.setPrice(busd.address, 100000000);

    expect(await oracle.getPrice(busd.address)).to.equal(100000000n);
    await expect(oracle.getPrice(btc.address)).to.be.revertedWith(
      "Price not set for this token",
    );
  });

  it("sets batch prices by asset address", async function () {
    await oracle.setPrices(
      [busd.address, btc.address],
      [100000000, 5000000000000],
    );

    expect(await oracle.getPrice(busd.address)).to.equal(100000000n);
    expect(await oracle.getPrice(btc.address)).to.equal(5000000000000n);

    const prices = await oracle.getPrices([busd.address, btc.address]);
    expect(prices[0]).to.equal(100000000n);
    expect(prices[1]).to.equal(5000000000000n);
  });

  it("blocks non-owners and invalid prices", async function () {
    await expect(
      oracle.connect(alice).setPrice(busd.address, 100000000),
    ).to.be.revertedWith("Not the owner");
    await expect(
      oracle.setPrice(ethers.ZeroAddress, 100000000),
    ).to.be.revertedWith("Invalid token address");
    await expect(oracle.setPrice(busd.address, 0)).to.be.revertedWith(
      "Price must be positive",
    );
    await expect(
      oracle.setPrices([busd.address], [100000000, 5000000000000]),
    ).to.be.revertedWith("Mismatched array lengths");
  });
});
