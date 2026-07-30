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

type FeedCheck = {
  token: string;
  tokenSymbol: string;
  tokenDecimals: number;
  feed: string;
  feedDescription: string;
  feedDecimals: number;
  latestAnswer: string;
  updatedAt: number;
};

const ERC20_METADATA_ABI = [
  "function symbol() view returns (string)",
  "function decimals() view returns (uint8)",
];
const CHAINLINK_FEED_ABI = [
  "function description() view returns (string)",
  "function decimals() view returns (uint8)",
  "function latestRoundData() view returns (uint80,int256,uint256,uint256,uint80)",
];

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

function requiredPositiveInteger(name: string): number {
  const value = process.env[name];
  const parsed = Number(value);
  if (!value || !Number.isSafeInteger(parsed) || parsed <= 0) {
    throw new Error(`${name} must be a positive integer`);
  }
  return parsed;
}

async function requireContract(name: string, address: string) {
  if ((await ethers.provider.getCode(address)) === "0x") {
    throw new Error(`${name} has no deployed code at ${address}`);
  }
}

async function checkFeedConfig(
  config: FeedConfig,
  blockTimestamp: number,
): Promise<FeedCheck> {
  const token = ethers.getAddress(config.token);
  const feed = ethers.getAddress(config.feed);
  if (!Number.isSafeInteger(config.maxStaleness) || config.maxStaleness <= 0) {
    throw new Error(`invalid maxStaleness for token ${token}`);
  }

  await Promise.all([
    requireContract("feed token", token),
    requireContract("Chainlink feed", feed),
  ]);

  const tokenContract = new ethers.Contract(
    token,
    ERC20_METADATA_ABI,
    ethers.provider,
  );
  const feedContract = new ethers.Contract(
    feed,
    CHAINLINK_FEED_ABI,
    ethers.provider,
  );
  const [tokenSymbol, tokenDecimalsValue, feedDescription, feedDecimalsValue] =
    await Promise.all([
      tokenContract.symbol() as Promise<string>,
      tokenContract.decimals() as Promise<bigint>,
      feedContract.description() as Promise<string>,
      feedContract.decimals() as Promise<bigint>,
    ]);
  const tokenDecimals = Number(tokenDecimalsValue);
  const feedDecimals = Number(feedDecimalsValue);
  if (!tokenSymbol.trim()) {
    throw new Error(`token ${token} returned an empty symbol`);
  }
  if (!Number.isSafeInteger(tokenDecimals) || tokenDecimals > 36) {
    throw new Error(`token ${token} has unsupported decimals ${tokenDecimals}`);
  }
  if (!feedDescription.trim()) {
    throw new Error(`Chainlink feed ${feed} returned an empty description`);
  }
  if (!Number.isSafeInteger(feedDecimals) || feedDecimals > 36) {
    throw new Error(
      `Chainlink feed ${feed} has unsupported decimals ${feedDecimals}`,
    );
  }

  const [roundId, answer, , updatedAtValue, answeredInRound] =
    await feedContract.latestRoundData();
  const updatedAt = Number(updatedAtValue);
  if (answer <= 0n) {
    throw new Error(`Chainlink feed ${feed} returned a non-positive answer`);
  }
  if (
    !Number.isSafeInteger(updatedAt) ||
    updatedAt <= 0 ||
    updatedAt > blockTimestamp
  ) {
    throw new Error(`Chainlink feed ${feed} returned an invalid timestamp`);
  }
  if (blockTimestamp - updatedAt > config.maxStaleness) {
    throw new Error(
      `Chainlink feed ${feed} is older than maxStaleness for token ${token}`,
    );
  }
  if (answeredInRound < roundId) {
    throw new Error(`Chainlink feed ${feed} returned an incomplete round`);
  }

  return {
    token,
    tokenSymbol,
    tokenDecimals,
    feed,
    feedDescription,
    feedDecimals,
    latestAnswer: answer.toString(),
    updatedAt,
  };
}

const networkInfo = await ethers.provider.getNetwork();
if (networkInfo.chainId !== EXPECTED_CHAIN_ID) {
  throw new Error(
    `production deployment requires Sepolia chain ${EXPECTED_CHAIN_ID}, connected to ${networkInfo.chainId}`,
  );
}

const [deployer] = await ethers.getSigners();
const multisigOwners = requiredConfig<string>("PRISM_MULTISIG_OWNERS").map(
  (owner, index) => {
    if (!ethers.isAddress(owner) || owner === ethers.ZeroAddress) {
      throw new Error(
        `PRISM_MULTISIG_OWNERS[${index}] must be a non-zero address`,
      );
    }
    return ethers.getAddress(owner);
  },
);
if (new Set(multisigOwners).size !== multisigOwners.length) {
  throw new Error("PRISM_MULTISIG_OWNERS must contain unique addresses");
}
const multisigThreshold = requiredPositiveInteger("PRISM_MULTISIG_THRESHOLD");
if (multisigThreshold > multisigOwners.length) {
  throw new Error(
    "PRISM_MULTISIG_THRESHOLD cannot exceed the number of owners",
  );
}
const feeAddress = requiredAddress("PRISM_FEE_ADDRESS");
const router = requiredAddress("PRISM_UNISWAP_V3_ROUTER");
const quoter = requiredAddress("PRISM_UNISWAP_V3_QUOTER");
const feeds = requiredConfig<FeedConfig>("PRISM_CHAINLINK_FEEDS");
const pools = requiredConfig<PoolConfig>("PRISM_UNISWAP_V3_POOLS");

await Promise.all([
  requireContract("PRISM_UNISWAP_V3_ROUTER", router),
  requireContract("PRISM_UNISWAP_V3_QUOTER", quoter),
]);

const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest Sepolia block is unavailable");
}
const feedChecks = await Promise.all(
  feeds.map((config) => checkFeedConfig(config, latestBlock.timestamp)),
);
if (new Set(feedChecks.map(({ token }) => token)).size !== feedChecks.length) {
  throw new Error("PRISM_CHAINLINK_FEEDS must configure each token once");
}

const multisigContract = await ethers.deployContract("ThresholdMultiSig", [
  multisigOwners,
  multisigThreshold,
]);
const oracle = await ethers.deployContract("ChainlinkOracle", [
  deployer.address,
]);
const swap = await ethers.deployContract("UniswapV3SwapAdapter", [
  deployer.address,
  router,
  quoter,
]);
await Promise.all([
  multisigContract.waitForDeployment(),
  oracle.waitForDeployment(),
  swap.waitForDeployment(),
]);
const multisig = await multisigContract.getAddress();

for (const config of feeds) {
  const token = ethers.getAddress(config.token);
  const feed = ethers.getAddress(config.feed);
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
  multisigOwners,
  multisigThreshold: (await multisigContract.threshold()).toString(),
  feeAddress,
  chainlinkOracle: await oracle.getAddress(),
  uniswapV3SwapAdapter: await swap.getAddress(),
  uniswapV3Router: router,
  uniswapV3Quoter: quoter,
  prismPool: await prismPool.getAddress(),
  feeds,
  feedChecks,
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
