import assert from "node:assert/strict";
import { network } from "hardhat";

const { ethers } = await network.create();

async function expectRevert(action, message) {
  try {
    await action;
  } catch (error) {
    assert.ok(
      error.message.includes(message),
      `Expected revert message "${message}", got "${error.message}"`,
    );
    return;
  }

  assert.fail(`Expected transaction to revert with "${message}"`);
}

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

  describe("PositionToken", function () {
    it("starts with ERC20 metadata and zero supply", async function () {
      assert.equal(await lenderPositionToken.name(), "Lending BUSD");
      assert.equal(await lenderPositionToken.symbol(), "lBUSD");
      assert.equal(await lenderPositionToken.decimals(), 18n);
      assert.equal(await lenderPositionToken.totalSupply(), 0n);
    });

    it("lets the owner add and remove minters", async function () {
      assert.equal(await lenderPositionToken.isMinter(alice.address), false);

      await lenderPositionToken.setMinter(alice.address, true);

      assert.equal(await lenderPositionToken.isMinter(alice.address), true);

      await lenderPositionToken.setMinter(alice.address, false);

      assert.equal(await lenderPositionToken.isMinter(alice.address), false);
    });

    it("blocks non-owners from managing minters", async function () {
      await expectRevert(
        lenderPositionToken.connect(alice).setMinter(alice.address, true),
        "OwnableUnauthorizedAccount",
      );
      await expectRevert(
        lenderPositionToken.connect(alice).setMinter(owner.address, false),
        "OwnableUnauthorizedAccount",
      );
    });

    it("lets minters mint and burn receipt tokens", async function () {
      await lenderPositionToken.setMinter(alice.address, true);

      await lenderPositionToken.connect(alice).mint(bob.address, 1000);
      assert.equal(await lenderPositionToken.balanceOf(bob.address), 1000n);
      assert.equal(await lenderPositionToken.totalSupply(), 1000n);

      await lenderPositionToken.connect(alice).burn(bob.address, 400);
      assert.equal(await lenderPositionToken.balanceOf(bob.address), 600n);
      assert.equal(await lenderPositionToken.totalSupply(), 600n);
    });

    it("blocks accounts that are not minters from minting or burning", async function () {
      await expectRevert(
        borrowerPositionToken.connect(alice).mint(bob.address, 1000),
        "caller is not minter",
      );
      await expectRevert(
        borrowerPositionToken.connect(alice).burn(bob.address, 1000),
        "caller is not minter",
      );
    });

    it("supports normal ERC20 transfer and allowance behavior", async function () {
      await lenderPositionToken.setMinter(owner.address, true);
      await lenderPositionToken.mint(alice.address, 1000);

      await lenderPositionToken.connect(alice).transfer(bob.address, 250);
      assert.equal(await lenderPositionToken.balanceOf(alice.address), 750n);
      assert.equal(await lenderPositionToken.balanceOf(bob.address), 250n);

      await lenderPositionToken.connect(alice).approve(owner.address, 300);
      await lenderPositionToken.transferFrom(alice.address, bob.address, 300);

      assert.equal(await lenderPositionToken.balanceOf(alice.address), 450n);
      assert.equal(await lenderPositionToken.balanceOf(bob.address), 550n);
      assert.equal(
        await lenderPositionToken.allowance(alice.address, owner.address),
        0n,
      );
    });
  });
});
