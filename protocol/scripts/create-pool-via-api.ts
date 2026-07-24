import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type LocalDeployment = {
  chainId: string;
  deployer: string;
  prismPool: string;
  lendToken: string;
  collateralToken: string;
  lenderPositionToken: string;
  borrowerPositionToken: string;
};

type PreparedTransaction = {
  to: string;
  data: string;
  value: string;
  chainId: string;
};

type LoginResponse = {
  tokenId: string;
};

type DataResponse<T> = {
  data: T;
};

type IndexedPool = {
  index: number;
};

const apiUrl = (process.env.PRISM_API_URL ?? "http://127.0.0.1:8080").replace(
  /\/$/,
  "",
);
const username = process.env.PRISM_ADMIN_USERNAME ?? "admin";
const password = process.env.PRISM_ADMIN_PASSWORD ?? "password";
const indexTimeoutMs = parsePositiveInteger(
  "PRISM_INDEX_TIMEOUT_MS",
  process.env.PRISM_INDEX_TIMEOUT_MS ?? "90000",
);

const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deploymentPath = path.join(protocolRoot, "deployments", "local.json");
const deployment = await readLocalDeployment();

const { ethers } = await network.create();
const [owner] = await ethers.getSigners();
const chain = await ethers.provider.getNetwork();
const pool = await ethers.getContractAt(
  "PrismPool",
  deployment.prismPool,
  owner,
);

if (chain.chainId.toString() !== deployment.chainId) {
  throw new Error(
    `deployment chain ${deployment.chainId} does not match connected chain ${chain.chainId}`,
  );
}
if (
  ethers.getAddress(await pool.owner()) !== ethers.getAddress(owner.address) ||
  ethers.getAddress(deployment.deployer) !== ethers.getAddress(owner.address)
) {
  throw new Error(
    `local signer ${owner.address} is not the deployed PrismPool owner`,
  );
}

const login = await postJSON<LoginResponse>(`${apiUrl}/api/v1/user/login`, {
  name: username,
  password,
});

const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest block is unavailable");
}

const poolCountBefore = await pool.poolCount();
const settleTime = latestBlock.timestamp + 24 * 60 * 60;
const prepared = await postJSON<DataResponse<PreparedTransaction>>(
  `${apiUrl}/api/v1/pools`,
  {
    settleTime: settleTime.toString(),
    maturityTime: (settleTime + 7 * 24 * 60 * 60).toString(),
    interestRate: "1000000",
    maxLendSupply: ethers.parseEther("100000").toString(),
    collateralizationRatio: "200000000",
    lendToken: deployment.lendToken,
    collateralToken: deployment.collateralToken,
    lenderPositionToken: deployment.lenderPositionToken,
    borrowerPositionToken: deployment.borrowerPositionToken,
    liquidateRate: "20000000",
  },
  login.tokenId,
);

validatePreparedTransaction(prepared.data);

const transaction = await owner.sendTransaction({
  to: prepared.data.to,
  data: prepared.data.data,
  value: BigInt(prepared.data.value),
});
console.log(`Broadcast transaction ${transaction.hash}`);

const receipt = await transaction.wait();
if (receipt === null || receipt.status !== 1) {
  throw new Error(`createPool transaction ${transaction.hash} failed`);
}

const poolCountAfter = await pool.poolCount();
if (poolCountAfter !== poolCountBefore + 1n) {
  throw new Error(
    `expected poolCount ${poolCountBefore + 1n}, received ${poolCountAfter}`,
  );
}

const poolId = Number(poolCountBefore);
console.log(`Created on-chain pool ${poolId} in block ${receipt.blockNumber}`);
await waitForIndexedPool(poolId, deployment.chainId);
console.log(`Backend indexed pool ${poolId}`);

function validatePreparedTransaction(preparedTransaction: PreparedTransaction) {
  if (preparedTransaction.chainId !== deployment.chainId) {
    throw new Error(
      `API returned chain ${preparedTransaction.chainId}, expected ${deployment.chainId}`,
    );
  }
  if (
    ethers.getAddress(preparedTransaction.to) !==
    ethers.getAddress(deployment.prismPool)
  ) {
    throw new Error(
      `API returned target ${preparedTransaction.to}, expected ${deployment.prismPool}`,
    );
  }
  if (BigInt(preparedTransaction.value) !== 0n) {
    throw new Error("createPool transaction must not transfer native currency");
  }
}

async function postJSON<T>(
  url: string,
  body: unknown,
  token?: string,
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token !== undefined) {
    headers.Authorization = `Bearer ${token}`;
  }
  const response = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(body),
  });
  return readJSON<T>(response);
}

async function readJSON<T>(response: Response): Promise<T> {
  const text = await response.text();
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}: ${text}`);
  }
  return JSON.parse(text) as T;
}

async function waitForIndexedPool(poolId: number, chainId: string) {
  const deadline = Date.now() + indexTimeoutMs;
  let lastError = "";
  while (Date.now() < deadline) {
    try {
      const response = await fetch(
        `${apiUrl}/api/v1/poolBaseInfo?chainId=${encodeURIComponent(chainId)}`,
      );
      const result = await readJSON<DataResponse<IndexedPool[]>>(response);
      if (result.data.some((item) => item.index === poolId)) {
        return;
      }
      lastError = `pool ${poolId} is not present yet`;
    } catch (error) {
      lastError = error instanceof Error ? error.message : String(error);
    }
    await new Promise((resolve) => setTimeout(resolve, 2_000));
  }
  throw new Error(
    `the transaction succeeded, but the backend did not index pool ${poolId} within ${indexTimeoutMs}ms (${lastError}); ensure the API and scheduler are running and share the same database`,
  );
}

function parsePositiveInteger(name: string, value: string): number {
  const parsed = Number(value);
  if (!Number.isSafeInteger(parsed) || parsed <= 0) {
    throw new Error(`${name} must be a positive integer`);
  }
  return parsed;
}

async function readLocalDeployment(): Promise<LocalDeployment> {
  try {
    return JSON.parse(
      await readFile(deploymentPath, "utf8"),
    ) as LocalDeployment;
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === "ENOENT") {
      throw new Error(
        `local deployment file is missing at ${deploymentPath}; keep the Hardhat node running and run "npm run deploy:local" first`,
        { cause: error },
      );
    }
    throw error;
  }
}
