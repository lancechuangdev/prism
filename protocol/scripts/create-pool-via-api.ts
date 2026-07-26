import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type LocalDeployment = {
  chainId: string;
  prismPool: string;
  multisig: string;
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
const signers = await ethers.getSigners();
const chain = await ethers.provider.getNetwork();
const pool = await ethers.getContractAt("PrismPool", deployment.prismPool);
const multisig = await ethers.getContractAt(
  "ThresholdMultiSig",
  deployment.multisig,
);

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

const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest block is unavailable");
}

const poolCountBefore = await pool.poolCount();
const settleTime = latestBlock.timestamp + 24 * 60 * 60;
const nonce = Date.now().toString();
const prepared = await postJSON<DataResponse<PreparedProposal>>(
  `${apiUrl}/api/v1/multisig/proposals`,
  {
    chain_id: deployment.chainId,
    nonce,
    operation: {
      type: "create_pool",
      params: {
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

async function validatePreparedProposal(
  preparedProposal: PreparedProposal,
  expectedNonce: string,
) {
  if (preparedProposal.proposal.operation !== "create_pool") {
    throw new Error(
      `API returned operation ${preparedProposal.proposal.operation}, expected create_pool`,
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
    throw new Error("createPool proposal must not transfer native currency");
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
