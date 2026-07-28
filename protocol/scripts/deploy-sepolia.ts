import { mkdir, rename, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

const { ethers } = await network.create();
const EXPECTED_CHAIN_ID = 11155111n;

type FeedConfig = {
  token: string;
  feed: string;
  maxStaleness: number;
};

type PoolConfig = {
  tokenIn: string;
  tokenOut: string;
  fee: number;
};

function requiredAddress(name: string): string {
  const value = process.env[name];
  if (!value || !ethers.isAddress(value) || value === ethers.ZeroAddress) {
    throw new Error(`${name} must be a non-zero address`);
  }
  return ethers.getAddress(value);
}

function requiredConfig<T>(name: string): T[] {
  const value = process.env[name];
  if (!value) {
    throw new Error(`${name} is required`);
  }
  const parsed: unknown = JSON.parse(value);
  if (!Array.isArray(parsed) || parsed.length === 0) {
    throw new Error(`${name} must be a non-empty JSON array`);
  }
  return parsed as T[];
}

async function requireContract(name: string, address: string) {
  if ((await ethers.provider.getCode(address)) === "0x") {
    throw new Error(`${name} has no deployed code at ${address}`);
  }
}

const networkInfo = await ethers.provider.getNetwork();
if (networkInfo.chainId !== EXPECTED_CHAIN_ID) {
  throw new Error(
    `production deployment requires Sepolia chain ${EXPECTED_CHAIN_ID}, connected to ${networkInfo.chainId}`,
  );
}

const [deployer] = await ethers.getSigners();
const multisig = requiredAddress("PRISM_MULTISIG_ADDRESS");
const feeAddress = requiredAddress("PRISM_FEE_ADDRESS");
const router = requiredAddress("PRISM_UNISWAP_V3_ROUTER");
const quoter = requiredAddress("PRISM_UNISWAP_V3_QUOTER");
const feeds = requiredConfig<FeedConfig>("PRISM_CHAINLINK_FEEDS");
const pools = requiredConfig<PoolConfig>("PRISM_UNISWAP_V3_POOLS");

await Promise.all([
  requireContract("PRISM_MULTISIG_ADDRESS", multisig),
  requireContract("PRISM_UNISWAP_V3_ROUTER", router),
  requireContract("PRISM_UNISWAP_V3_QUOTER", quoter),
]);

const oracle = await ethers.deployContract("ChainlinkOracle", [
  deployer.address,
]);
const swap = await ethers.deployContract("UniswapV3SwapAdapter", [
  deployer.address,
  router,
  quoter,
]);
await Promise.all([oracle.waitForDeployment(), swap.waitForDeployment()]);

for (const config of feeds) {
  const token = ethers.getAddress(config.token);
  const feed = ethers.getAddress(config.feed);
  if (!Number.isSafeInteger(config.maxStaleness) || config.maxStaleness <= 0) {
    throw new Error(`invalid maxStaleness for token ${token}`);
  }
  await Promise.all([
    requireContract("feed token", token),
    requireContract("Chainlink feed", feed),
  ]);
  await (await oracle.setFeed(token, feed, config.maxStaleness)).wait();
  await oracle.getPrice(token);
}

for (const config of pools) {
  const tokenIn = ethers.getAddress(config.tokenIn);
  const tokenOut = ethers.getAddress(config.tokenOut);
  if (!Number.isSafeInteger(config.fee) || config.fee <= 0) {
    throw new Error(`invalid Uniswap fee for ${tokenIn}/${tokenOut}`);
  }
  await Promise.all([
    requireContract("pool input token", tokenIn),
    requireContract("pool output token", tokenOut),
  ]);
  await (await swap.setPoolFee(tokenIn, tokenOut, config.fee)).wait();
}

await (await oracle.transferOwnership(multisig)).wait();
await (await swap.transferOwnership(multisig)).wait();

const prismPool = await ethers.deployContract("PrismPool", [
  multisig,
  await oracle.getAddress(),
  await swap.getAddress(),
  feeAddress,
]);
await prismPool.waitForDeployment();

const deployment = {
  schemaVersion: 1,
  environment: "production",
  network: "sepolia",
  chainId: networkInfo.chainId.toString(),
  deployedAt: new Date().toISOString(),
  deploymentBlock: await ethers.provider.getBlockNumber(),
  deployer: deployer.address,
  multisig,
  feeAddress,
  chainlinkOracle: await oracle.getAddress(),
  uniswapV3SwapAdapter: await swap.getAddress(),
  uniswapV3Router: router,
  uniswapV3Quoter: quoter,
  prismPool: await prismPool.getAddress(),
  feeds,
  pools,
};

const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deploymentPath = path.join(protocolRoot, "deployments", "sepolia.json");
const temporaryPath = `${deploymentPath}.tmp`;
await mkdir(path.dirname(deploymentPath), { recursive: true });
await writeFile(
  temporaryPath,
  `${JSON.stringify(deployment, null, 2)}\n`,
  "utf8",
);
await rename(temporaryPath, deploymentPath);

console.log(JSON.stringify(deployment, null, 2));
console.log(`Deployment metadata saved to ${deploymentPath}`);
