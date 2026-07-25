import { mkdir, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

const { ethers } = await network.create();
const [deployer, feeRecipient, thirdMultisigOwner] = await ethers.getSigners();

const oracle = await ethers.deployContract("MockOracle");
const swap = await ethers.deployContract("FixedRateSwap");
const multisig = await ethers.deployContract("ThresholdMultiSig", [
  [deployer.address, feeRecipient.address, thirdMultisigOwner.address],
  2,
]);
const lendToken = await ethers.deployContract("PositionToken", [
  "Prism USD",
  "pUSD",
]);
const collateralToken = await ethers.deployContract("PositionToken", [
  "Prism ETH",
  "pETH",
]);
const lenderPositionToken = await ethers.deployContract("PositionToken", [
  "Prism Lender Position",
  "pLEND",
]);
const borrowerPositionToken = await ethers.deployContract("PositionToken", [
  "Prism Borrower Position",
  "pBORROW",
]);

const pool = await ethers.deployContract("PrismPool", [
  await oracle.getAddress(),
  await swap.getAddress(),
  feeRecipient.address,
]);

await Promise.all([
  oracle.waitForDeployment(),
  swap.waitForDeployment(),
  multisig.waitForDeployment(),
  lendToken.waitForDeployment(),
  collateralToken.waitForDeployment(),
  lenderPositionToken.waitForDeployment(),
  borrowerPositionToken.waitForDeployment(),
  pool.waitForDeployment(),
]);

const poolAddress = await pool.getAddress();
await (await lenderPositionToken.setMinter(poolAddress, true)).wait();
await (await borrowerPositionToken.setMinter(poolAddress, true)).wait();

const lendTokenAddress = await lendToken.getAddress();
const collateralTokenAddress = await collateralToken.getAddress();
await (
  await oracle.setPrices(
    [lendTokenAddress, collateralTokenAddress],
    [ethers.parseEther("1"), ethers.parseEther("3000")],
  )
).wait();

const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest block is unavailable");
}

const settleTime = latestBlock.timestamp + 24 * 60 * 60;
await (
  await pool.createPool({
    settleTime,
    maturityTime: settleTime + 7 * 24 * 60 * 60,
    interestRate: 1_000_000,
    maxLendSupply: ethers.parseEther("100000"),
    collateralizationRatio: 200_000_000,
    lendToken: lendTokenAddress,
    collateralToken: collateralTokenAddress,
    lenderPositionToken: await lenderPositionToken.getAddress(),
    borrowerPositionToken: await borrowerPositionToken.getAddress(),
    liquidateRate: 20_000_000,
  })
).wait();

const chain = await ethers.provider.getNetwork();
const deployment = {
  rpcUrl: "http://127.0.0.1:8545",
  chainId: chain.chainId.toString(),
  deployer: deployer.address,
  prismPool: poolAddress,
  oracle: await oracle.getAddress(),
  dexSwap: await swap.getAddress(),
  multisig: await multisig.getAddress(),
  multisigOwners: [
    deployer.address,
    feeRecipient.address,
    thirdMultisigOwner.address,
  ],
  multisigThreshold: (await multisig.threshold()).toString(),
  lendToken: lendTokenAddress,
  collateralToken: collateralTokenAddress,
  lenderPositionToken: await lenderPositionToken.getAddress(),
  borrowerPositionToken: await borrowerPositionToken.getAddress(),
  poolCount: (await pool.poolCount()).toString(),
};

const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deploymentPath = path.join(protocolRoot, "deployments", "local.json");
await mkdir(path.dirname(deploymentPath), { recursive: true });
await writeFile(
  deploymentPath,
  `${JSON.stringify(deployment, null, 2)}\n`,
  "utf8",
);

console.log(JSON.stringify(deployment, null, 2));
console.log(`Deployment addresses saved to ${deploymentPath}`);
