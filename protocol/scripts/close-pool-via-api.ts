import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type LocalDeployment = {
  chainId: string;
  prismPool: string;
  multisig: string;
  dexSwap: string;
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

type Proposal = {
  transactionHash: string;
  operation: string;
  target: string;
  value: string;
  data: string;
  nonce: string;
};

type PreparedProposal = {
  proposal: Proposal;
  approvalTransaction: PreparedTransaction;
  executionTransaction: PreparedTransaction;
};

type IndexedPool = {
  index: number;
  pool_data: {
    state: string;
  };
};

const ACTIVE = 1n;
const REPAID = 2n;
const LIQUIDATED = 3n;
const operation =
  process.env.PRISM_POOL_OPERATION === "liquidate_pool"
    ? "liquidate_pool"
    : "repay_pool";
const contractMethod =
  operation === "liquidate_pool" ? "liquidate" : "repayPool";
const expectedState = operation === "liquidate_pool" ? LIQUIDATED : REPAID;
const expectedStateName =
  operation === "liquidate_pool" ? "LIQUIDATED" : "REPAID";

const apiUrl = (process.env.PRISM_API_URL ?? "http://127.0.0.1:8080").replace(
  /\/$/,
  "",
);
const username = process.env.PRISM_ADMIN_USERNAME ?? "admin";
const password = process.env.PRISM_ADMIN_PASSWORD ?? "password";
const poolId = parseNonNegativeInteger(
  "PRISM_POOL_ID",
  process.env.PRISM_POOL_ID ?? "0",
);
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
const signers = await ethers.getSigners();
const chain = await ethers.provider.getNetwork();
const pool = await ethers.getContractAt("PrismPool", deployment.prismPool);
const multisig = await ethers.getContractAt(
  "ThresholdMultiSig",
  deployment.multisig,
);
const swap = await ethers.getContractAt("FixedRateSwap", deployment.dexSwap);

if (chain.chainId.toString() !== deployment.chainId) {
  throw new Error(
    `deployment chain ${deployment.chainId} does not match connected chain ${chain.chainId}`,
  );
}
if (
  ethers.getAddress(await pool.owner()) !==
  ethers.getAddress(deployment.multisig)
) {
  throw new Error(
    `PrismPool owner ${await pool.owner()} does not match the deployed multisig ${deployment.multisig}`,
  );
}
if (poolId >= (await pool.poolCount())) {
  throw new Error(`pool ${poolId} does not exist`);
}
if ((await pool.getPoolState(poolId)) !== ACTIVE) {
  throw new Error(`pool ${poolId} is not in the ACTIVE state`);
}

const poolInfo = await pool.getPool(poolId);
const poolData = await pool.getPoolData(poolId);
const requiredRepayment = await pool.getRequiredRepayment(poolId);
let quotedCollateral: bigint;
try {
  quotedCollateral = await swap.getAmountIn(
    poolInfo.collateralToken,
    poolInfo.lendToken,
    requiredRepayment,
  );
} catch (error) {
  throw new Error(
    "FixedRateSwap cannot quote repayment; configure its collateral-to-lend rate before running this helper",
    { cause: error },
  );
}
const lendToken = await ethers.getContractAt(
  "PositionToken",
  poolInfo.lendToken,
);
const swapLendBalance = await lendToken.balanceOf(deployment.dexSwap);
if (swapLendBalance < requiredRepayment) {
  throw new Error(
    `FixedRateSwap has ${swapLendBalance} lend tokens but repayment requires ${requiredRepayment}; fund the swap before running this helper`,
  );
}

const configuredMaximum = process.env.PRISM_MAX_COLLATERAL_AMOUNT;
const maxCollateralAmount =
  configuredMaximum === undefined
    ? poolData.settleAmountBorrow
    : parsePositiveBigInt("PRISM_MAX_COLLATERAL_AMOUNT", configuredMaximum);
if (maxCollateralAmount <= 0n) {
  throw new Error("the pool has no settled collateral available for repayment");
}
if (quotedCollateral > maxCollateralAmount) {
  throw new Error(
    `repayment quote ${quotedCollateral} exceeds max collateral ${maxCollateralAmount}`,
  );
}

const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest block is unavailable");
}
if (
  operation === "repay_pool" &&
  BigInt(latestBlock.timestamp) < poolInfo.maturityTime
) {
  await ethers.provider.send("evm_setNextBlockTimestamp", [
    Number(poolInfo.maturityTime),
  ]);
  await ethers.provider.send("evm_mine", []);
  console.log(`Advanced local chain to maturity ${poolInfo.maturityTime}`);
}
if (
  operation === "liquidate_pool" &&
  !(await pool.isUndercollateralized(poolId))
) {
  throw new Error(`pool ${poolId} is not undercollateralized`);
}

const ownerCount = Number(await multisig.ownerCount());
const ownerAddresses = new Set<string>();
for (let index = 0; index < ownerCount; index++) {
  ownerAddresses.add(
    ethers.getAddress(await multisig.getOwner(index)).toLowerCase(),
  );
}
const ownerSigners = signers.filter((signer) =>
  ownerAddresses.has(signer.address.toLowerCase()),
);
const threshold = parsePositiveInteger(
  "on-chain multisig threshold",
  (await multisig.threshold()).toString(),
);
if (ownerSigners.length < threshold) {
  throw new Error(
    `only ${ownerSigners.length} local multisig owner signers are available, but ${threshold} approvals are required`,
  );
}

const login = await postJSON<LoginResponse>(`${apiUrl}/api/v1/user/login`, {
  name: username,
  password,
});
const nonce = Date.now().toString();
const prepared = await postJSON<DataResponse<PreparedProposal>>(
  `${apiUrl}/api/v1/multisig/proposals`,
  {
    chain_id: deployment.chainId,
    nonce,
    operation: {
      type: operation,
      params: {
        poolId: poolId.toString(),
        maxCollateralAmount: maxCollateralAmount.toString(),
      },
    },
  },
  login.tokenId,
);

await validatePreparedProposal(prepared.data, nonce);

for (const signer of ownerSigners.slice(0, threshold)) {
  const approval = await signer.sendTransaction({
    to: prepared.data.approvalTransaction.to,
    data: prepared.data.approvalTransaction.data,
    value: BigInt(prepared.data.approvalTransaction.value),
  });
  console.log(`Owner ${signer.address} broadcast approval ${approval.hash}`);
  const approvalReceipt = await approval.wait();
  if (approvalReceipt === null || approvalReceipt.status !== 1) {
    throw new Error(`approval transaction ${approval.hash} failed`);
  }
}

const execution = await ownerSigners[0].sendTransaction({
  to: prepared.data.executionTransaction.to,
  data: prepared.data.executionTransaction.data,
  value: BigInt(prepared.data.executionTransaction.value),
});
console.log(`Broadcast execution transaction ${execution.hash}`);
const receipt = await execution.wait();
if (receipt === null || receipt.status !== 1) {
  throw new Error(`execution transaction ${execution.hash} failed`);
}
if ((await pool.getPoolState(poolId)) !== expectedState) {
  throw new Error(
    `expected pool ${poolId} to enter the ${expectedStateName} state`,
  );
}
console.log(
  `${expectedStateName} on-chain pool ${poolId} in block ${receipt.blockNumber}`,
);
await waitForIndexedState(poolId, deployment.chainId, expectedState.toString());
console.log(`Backend indexed pool ${poolId} state ${expectedStateName}`);

async function validatePreparedProposal(
  preparedProposal: PreparedProposal,
  expectedNonce: string,
) {
  if (preparedProposal.proposal.operation !== operation) {
    throw new Error(
      `API returned operation ${preparedProposal.proposal.operation}, expected ${operation}`,
    );
  }
  if (
    ethers.getAddress(preparedProposal.proposal.target) !==
    ethers.getAddress(deployment.prismPool)
  ) {
    throw new Error(
      `API returned proposal target ${preparedProposal.proposal.target}, expected ${deployment.prismPool}`,
    );
  }
  if (preparedProposal.proposal.nonce !== expectedNonce) {
    throw new Error(
      `API returned nonce ${preparedProposal.proposal.nonce}, expected ${expectedNonce}`,
    );
  }
  if (BigInt(preparedProposal.proposal.value) !== 0n) {
    throw new Error(`${operation} proposal must not transfer native currency`);
  }

  const decoded = pool.interface.decodeFunctionData(
    contractMethod,
    preparedProposal.proposal.data,
  );
  if (decoded[0] !== poolId || decoded[1] !== maxCollateralAmount) {
    throw new Error(
      `API encoded ${contractMethod}(${decoded[0]}, ${decoded[1]}), expected ${contractMethod}(${poolId}, ${maxCollateralAmount})`,
    );
  }

  for (const [name, transaction] of [
    ["approval", preparedProposal.approvalTransaction],
    ["execution", preparedProposal.executionTransaction],
  ] as const) {
    if (transaction.chainId !== deployment.chainId) {
      throw new Error(
        `API returned ${name} chain ${transaction.chainId}, expected ${deployment.chainId}`,
      );
    }
    if (
      ethers.getAddress(transaction.to) !==
      ethers.getAddress(deployment.multisig)
    ) {
      throw new Error(
        `API returned ${name} target ${transaction.to}, expected ${deployment.multisig}`,
      );
    }
    if (BigInt(transaction.value) !== 0n) {
      throw new Error(`${name} transaction must not transfer native currency`);
    }
  }

  const expectedHash = await multisig.getTransactionHash(
    preparedProposal.proposal.target,
    BigInt(preparedProposal.proposal.value),
    preparedProposal.proposal.data,
    BigInt(preparedProposal.proposal.nonce),
  );
  if (expectedHash !== preparedProposal.proposal.transactionHash) {
    throw new Error(
      `API returned proposal hash ${preparedProposal.proposal.transactionHash}, expected ${expectedHash}`,
    );
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

async function waitForIndexedState(
  expectedPoolId: bigint,
  chainId: string,
  expectedState: string,
) {
  const deadline = Date.now() + indexTimeoutMs;
  let lastError = "";
  while (Date.now() < deadline) {
    try {
      const response = await fetch(
        `${apiUrl}/api/v1/poolBaseInfo?chainId=${encodeURIComponent(chainId)}`,
      );
      const result = await readJSON<DataResponse<IndexedPool[]>>(response);
      const indexedPool = result.data.find(
        (item) => BigInt(item.index) === expectedPoolId,
      );
      if (indexedPool?.pool_data.state === expectedState) {
        return;
      }
      lastError = `pool ${expectedPoolId} has indexed state ${indexedPool?.pool_data.state ?? "missing"}`;
    } catch (error) {
      lastError = error instanceof Error ? error.message : String(error);
    }
    await new Promise((resolve) => setTimeout(resolve, 2_000));
  }
  throw new Error(
    `the transaction succeeded, but the backend did not index pool ${expectedPoolId} state ${expectedState} within ${indexTimeoutMs}ms (${lastError}); ensure the API and scheduler are running and share the same database`,
  );
}

function parsePositiveInteger(name: string, value: string): number {
  const parsed = Number(value);
  if (!Number.isSafeInteger(parsed) || parsed <= 0) {
    throw new Error(`${name} must be a positive integer`);
  }
  return parsed;
}

function parseNonNegativeInteger(name: string, value: string): bigint {
  if (!/^(0|[1-9][0-9]*)$/.test(value)) {
    throw new Error(`${name} must be a non-negative decimal integer`);
  }
  return BigInt(value);
}

function parsePositiveBigInt(name: string, value: string): bigint {
  const parsed = parseNonNegativeInteger(name, value);
  if (parsed === 0n) {
    throw new Error(`${name} must be a positive decimal integer`);
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
