import { expect } from "chai";
import { network } from "hardhat";

const { ethers } = await network.create();

describe("Position Token", function () {
  let owner;
  let alice;
  let bob;
  let busd;
  let btc;
  let lenderPositionToken;
  let borrowerPositionToken;

  beforeEach(async function () {
    [owner, alice, bob, busd, btc] = await ethers.getSigners();

    const positionToken = await ethers.getContractFactory("PositionToken");
    lenderPositionToken = await positionToken.deploy("Lending BUSD", "lBUSD");
    borrowerPositionToken = await positionToken.deploy(
      "Collateral BTC",
      "cBTC",
    );
    await lenderPositionToken.waitForDeployment();
    await borrowerPositionToken.waitForDeployment();
  });

  it("starts with ERC20 metadata and zero supply", async function () {
    expect(await lenderPositionToken.name()).to.equal("Lending BUSD");
    expect(await lenderPositionToken.symbol()).to.equal("lBUSD");
    expect(await lenderPositionToken.decimals()).to.equal(18n);
    expect(await lenderPositionToken.totalSupply()).to.equal(0n);
  });

  it("lets the owner add and remove minters", async function () {
    expect(await lenderPositionToken.isMinter(alice.address)).to.equal(false);

    await lenderPositionToken.setMinter(alice.address, true);

    expect(await lenderPositionToken.isMinter(alice.address)).to.equal(true);

    await lenderPositionToken.setMinter(alice.address, false);

    expect(await lenderPositionToken.isMinter(alice.address)).to.equal(false);
  });

  it("blocks non-owners from managing minters", async function () {
    await expect(
      lenderPositionToken.connect(alice).setMinter(alice.address, true),
    )
      .to.be.revertedWithCustomError(
        lenderPositionToken,
        "OwnableUnauthorizedAccount",
      )
      .withArgs(alice.address);
    await expect(
      lenderPositionToken.connect(alice).setMinter(owner.address, false),
    )
      .to.be.revertedWithCustomError(
        lenderPositionToken,
        "OwnableUnauthorizedAccount",
      )
      .withArgs(alice.address);
  });

  it("lets minters mint and burn receipt tokens", async function () {
    await lenderPositionToken.setMinter(alice.address, true);

    await lenderPositionToken.connect(alice).mint(bob.address, 1000);
    expect(await lenderPositionToken.balanceOf(bob.address)).to.equal(1000n);
    expect(await lenderPositionToken.totalSupply()).to.equal(1000n);

    await lenderPositionToken.connect(alice).burn(bob.address, 400);
    expect(await lenderPositionToken.balanceOf(bob.address)).to.equal(600n);
    expect(await lenderPositionToken.totalSupply()).to.equal(600n);
  });

  it("blocks accounts that are not minters from minting or burning", async function () {
    await expect(
      borrowerPositionToken.connect(alice).mint(bob.address, 1000),
    ).to.be.revertedWith("caller is not minter");
    await expect(
      borrowerPositionToken.connect(alice).burn(bob.address, 1000),
    ).to.be.revertedWith("caller is not minter");
  });

  it("supports normal ERC20 transfer and allowance behavior", async function () {
    await lenderPositionToken.setMinter(owner.address, true);
    await lenderPositionToken.mint(alice.address, 1000);

    await lenderPositionToken.connect(alice).transfer(bob.address, 250);
    expect(await lenderPositionToken.balanceOf(alice.address)).to.equal(750n);
    expect(await lenderPositionToken.balanceOf(bob.address)).to.equal(250n);

    await lenderPositionToken.connect(alice).approve(owner.address, 300);
    await lenderPositionToken.transferFrom(alice.address, bob.address, 300);

    expect(await lenderPositionToken.balanceOf(alice.address)).to.equal(450n);
    expect(await lenderPositionToken.balanceOf(bob.address)).to.equal(550n);
    expect(
      await lenderPositionToken.allowance(alice.address, owner.address),
    ).to.equal(0n);
  });
});
