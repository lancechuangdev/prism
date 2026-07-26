import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type Deployment = {
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

type PreparedProposal = {
  proposal: {
    transactionHash: string;
    operation: string;
    target: string;
    value: string;
    data: string;
    nonce: string;
  };
  approvalTransaction: PreparedTransaction;
  executionTransaction: PreparedTransaction;
};

type IndexedPool = {
  index: number;
  pool_data: { state: string };
};

const FUNDING = 0n;
const ACTIVE = 1n;
const apiUrl = (process.env.PRISM_API_URL ?? "http://127.0.0.1:8080").replace(
  /\/$/,
  "",
);
const username = process.env.PRISM_ADMIN_USERNAME ?? "admin";
const password = process.env.PRISM_ADMIN_PASSWORD ?? "password";
const lendDeposit = parsePositiveAmount(
  "PRISM_SETUP_LEND_AMOUNT",
  process.env.PRISM_SETUP_LEND_AMOUNT ?? "1000000000000000000000",
);
const collateralDeposit = parsePositiveAmount(
  "PRISM_SETUP_COLLATERAL_AMOUNT",
  process.env.PRISM_SETUP_COLLATERAL_AMOUNT ?? "1000000000000000000",
);
const swapRate = parsePositiveAmount(
  "PRISM_SETUP_SWAP_RATE",
  process.env.PRISM_SETUP_SWAP_RATE ?? "3000000000000000000000",
);
const swapLiquidity = parsePositiveAmount(
  "PRISM_SETUP_SWAP_LIQUIDITY",
  process.env.PRISM_SETUP_SWAP_LIQUIDITY ?? "100000000000000000000000",
);
const indexTimeoutMs = Number(process.env.PRISM_INDEX_TIMEOUT_MS ?? "90000");

const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deployment = JSON.parse(
  await readFile(path.join(protocolRoot, "deployments", "local.json"), "utf8"),
) as Deployment;

const { ethers } = await network.create();
const signers = await ethers.getSigners();
const chain = await ethers.provider.getNetwork();
if (chain.chainId.toString() !== deployment.chainId) {
  throw new Error(
    `deployment chain ${deployment.chainId} does not match connected chain ${chain.chainId}`,
  );
}

const pool = await ethers.getContractAt("PrismPool", deployment.prismPool);
const multisig = await ethers.getContractAt(
  "ThresholdMultiSig",
  deployment.multisig,
);
const swap = await ethers.getContractAt("FixedRateSwap", deployment.dexSwap);
const poolCount = await pool.poolCount();
if (poolCount === 0n) {
  throw new Error("no pools exist; run npm run create-pool:api first");
}
const poolId =
  process.env.PRISM_POOL_ID === undefined
    ? poolCount - 1n
    : parseNonNegativeAmount("PRISM_POOL_ID", process.env.PRISM_POOL_ID);
if (poolId >= poolCount) {
  throw new Error(`pool ${poolId} does not exist`);
}
if ((await pool.getPoolState(poolId)) !== FUNDING) {
  throw new Error(`pool ${poolId} is not in the FUNDING state`);
}

const poolInfo = await pool.getPool(poolId);
const latestBlock = await ethers.provider.getBlock("latest");
if (latestBlock === null) {
  throw new Error("latest block is unavailable");
}
if (BigInt(latestBlock.timestamp) >= poolInfo.settleTime) {
  throw new Error(
    `pool ${poolId} has reached its settlement time; create a new funding pool first`,
  );
}

const lender = signers[3];
const borrower = signers[4];
if (lender === undefined || borrower === undefined) {
  throw new Error("the local node must expose at least five test signers");
}
const lendToken = await ethers.getContractAt(
  "PositionToken",
  poolInfo.lendToken,
);
const collateralToken = await ethers.getContractAt(
  "PositionToken",
  poolInfo.collateralToken,
);

const lendTokenOwner = await signerForAddress(
  await lendToken.owner(),
  signers,
  "lend token owner",
);
const collateralTokenOwner = await signerForAddress(
  await collateralToken.owner(),
  signers,
  "collateral token owner",
);
const swapOwner = await signerForAddress(
  await swap.owner(),
  signers,
  "FixedRateSwap owner",
);

const lendOwnerWasMinter = await lendToken.isMinter(lendTokenOwner.address);
const collateralOwnerWasMinter = await collateralToken.isMinter(
  collateralTokenOwner.address,
);
if (!lendOwnerWasMinter) {
  await (
    await lendToken
      .connect(lendTokenOwner)
      .setMinter(lendTokenOwner.address, true)
  ).wait();
}
if (!collateralOwnerWasMinter) {
  await (
    await collateralToken
      .connect(collateralTokenOwner)
      .setMinter(collateralTokenOwner.address, true)
  ).wait();
}
try {
  await (
    await lendToken.connect(lendTokenOwner).mint(lender.address, lendDeposit)
  ).wait();
  await (
    await collateralToken
      .connect(collateralTokenOwner)
      .mint(borrower.address, collateralDeposit)
  ).wait();
  await (
    await lendToken
      .connect(lendTokenOwner)
      .mint(lendTokenOwner.address, swapLiquidity)
  ).wait();
  await (
    await lendToken
      .connect(lendTokenOwner)
      .transfer(deployment.dexSwap, swapLiquidity)
  ).wait();
} finally {
  if (!lendOwnerWasMinter) {
    await (
      await lendToken
        .connect(lendTokenOwner)
        .setMinter(lendTokenOwner.address, false)
    ).wait();
  }
  if (!collateralOwnerWasMinter) {
    await (
      await collateralToken
        .connect(collateralTokenOwner)
        .setMinter(collateralTokenOwner.address, false)
    ).wait();
  }
}
await (
  await swap
    .connect(swapOwner)
    .setRate(poolInfo.collateralToken, poolInfo.lendToken, swapRate)
).wait();

await (
  await lendToken.connect(lender).approve(deployment.prismPool, lendDeposit)
).wait();
await (
  await collateralToken
    .connect(borrower)
    .approve(deployment.prismPool, collateralDeposit)
).wait();
await (await pool.connect(lender).depositLend(poolId, lendDeposit)).wait();
await (
  await pool.connect(borrower).depositBorrow(poolId, collateralDeposit)
).wait();
console.log(
  `Funded pool ${poolId}: lend=${lendDeposit}, collateral=${collateralDeposit}`,
);
console.log(
  `Configured swap rate=${swapRate}, lend liquidity=${swapLiquidity}`,
);

await ethers.provider.send("evm_setNextBlockTimestamp", [
  Number(poolInfo.settleTime),
]);
await ethers.provider.send("evm_mine", []);
console.log(`Advanced local chain to settle time ${poolInfo.settleTime}`);

const ownerSigners = await multisigOwnerSigners(multisig, signers);
const threshold = Number(await multisig.threshold());
if (ownerSigners.length < threshold) {
  throw new Error(
    `only ${ownerSigners.length} local multisig signers are available; ${threshold} are required`,
  );
}

const login = await postJSON<{ tokenId: string }>(
  `${apiUrl}/api/v1/user/login`,
  { name: username, password },
);
const nonce = Date.now().toString();
const prepared = await postJSON<{ data: PreparedProposal }>(
  `${apiUrl}/api/v1/multisig/proposals`,
  {
    chain_id: deployment.chainId,
    nonce,
    operation: {
      type: "settle_pool",
      params: { poolId: poolId.toString() },
    },
  },
  login.tokenId,
);
await validateProposal(prepared.data, nonce);

for (const signer of ownerSigners.slice(0, threshold)) {
  const transaction = await signer.sendTransaction({
    to: prepared.data.approvalTransaction.to,
    data: prepared.data.approvalTransaction.data,
    value: BigInt(prepared.data.approvalTransaction.value),
  });
  await requireSuccessfulReceipt(transaction);
}
const execution = await ownerSigners[0].sendTransaction({
  to: prepared.data.executionTransaction.to,
  data: prepared.data.executionTransaction.data,
  value: BigInt(prepared.data.executionTransaction.value),
});
const receipt = await requireSuccessfulReceipt(execution);
if ((await pool.getPoolState(poolId)) !== ACTIVE) {
  throw new Error(
    `pool ${poolId} did not become ACTIVE; check deposit sizes and oracle prices`,
  );
}
console.log(`Activated pool ${poolId} in block ${receipt.blockNumber}`);
await waitForIndexedState(poolId, ACTIVE.toString());
console.log(
  `Pool ${poolId} is ready: PRISM_POOL_ID=${poolId} npm run repay-pool:api`,
);

async function validateProposal(prepared: PreparedProposal, nonce: string) {
  if (
    prepared.proposal.operation !== "settle_pool" ||
    prepared.proposal.nonce !== nonce ||
    ethers.getAddress(prepared.proposal.target) !==
      ethers.getAddress(deployment.prismPool)
  ) {
    throw new Error("API returned an unexpected settlement proposal");
  }
  const decoded = pool.interface.decodeFunctionData(
    "settle",
    prepared.proposal.data,
  );
  if (decoded[0] !== poolId) {
    throw new Error(`API encoded settlement for pool ${decoded[0]}`);
  }
  if (BigInt(prepared.proposal.value) !== 0n) {
    throw new Error("settlement proposal must not transfer native currency");
  }
  for (const [name, transaction] of [
    ["approval", prepared.approvalTransaction],
    ["execution", prepared.executionTransaction],
  ] as const) {
    if (
      transaction.chainId !== deployment.chainId ||
      ethers.getAddress(transaction.to) !==
        ethers.getAddress(deployment.multisig) ||
      BigInt(transaction.value) !== 0n
    ) {
      throw new Error(`API returned an unexpected ${name} transaction`);
    }
  }
  const expectedHash = await multisig.getTransactionHash(
    prepared.proposal.target,
    BigInt(prepared.proposal.value),
    prepared.proposal.data,
    BigInt(prepared.proposal.nonce),
  );
  if (expectedHash !== prepared.proposal.transactionHash) {
    throw new Error("API returned an invalid settlement proposal hash");
  }
}

async function multisigOwnerSigners(
  contract: typeof multisig,
  availableSigners: typeof signers,
) {
  const owners = new Set<string>();
  for (let index = 0; index < Number(await contract.ownerCount()); index++) {
    owners.add((await contract.getOwner(index)).toLowerCase());
  }
  return availableSigners.filter((signer) =>
    owners.has(signer.address.toLowerCase()),
  );
}

async function signerForAddress(
  address: string,
  availableSigners: typeof signers,
  label: string,
) {
  const signer = availableSigners.find(
    (candidate) => candidate.address.toLowerCase() === address.toLowerCase(),
  );
  if (signer === undefined) {
    throw new Error(`${label} ${address} is not a local signer`);
  }
  return signer;
}

async function requireSuccessfulReceipt(
  transaction: Awaited<ReturnType<(typeof signers)[number]["sendTransaction"]>>,
) {
  console.log(`Broadcast transaction ${transaction.hash}`);
  const receipt = await transaction.wait();
  if (receipt === null || receipt.status !== 1) {
    throw new Error(`transaction ${transaction.hash} failed`);
  }
  return receipt;
}

async function postJSON<T>(url: string, body: unknown, token?: string) {
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
  const text = await response.text();
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}: ${text}`);
  }
  return JSON.parse(text) as T;
}

async function waitForIndexedState(expectedPoolId: bigint, state: string) {
  const deadline = Date.now() + indexTimeoutMs;
  let lastState = "missing";
  while (Date.now() < deadline) {
    const response = await fetch(
      `${apiUrl}/api/v1/poolBaseInfo?chainId=${deployment.chainId}`,
    );
    if (response.ok) {
      const result = (await response.json()) as { data: IndexedPool[] };
      const indexed = result.data.find(
        (item) => BigInt(item.index) === expectedPoolId,
      );
      lastState = indexed?.pool_data.state ?? "missing";
      if (lastState === state) {
        return;
      }
    }
    await new Promise((resolve) => setTimeout(resolve, 2_000));
  }
  throw new Error(
    `backend did not index pool ${expectedPoolId} state ${state} within ${indexTimeoutMs}ms; last state=${lastState}`,
  );
}

function parseNonNegativeAmount(name: string, value: string): bigint {
  if (!/^(0|[1-9][0-9]*)$/.test(value)) {
    throw new Error(`${name} must be a non-negative decimal integer`);
  }
  return BigInt(value);
}

function parsePositiveAmount(name: string, value: string): bigint {
  const parsed = parseNonNegativeAmount(name, value);
  if (parsed === 0n) {
    throw new Error(`${name} must be positive`);
  }
  return parsed;
}
