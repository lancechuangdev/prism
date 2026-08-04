import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type SepoliaDeployment = {
  chainId: string;
  prismPool: string;
  multisig: string;
};

const SEPOLIA_CHAIN_ID = "11155111";
const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deployment = JSON.parse(
  await readFile(
    path.join(protocolRoot, "deployments", "sepolia.json"),
    "utf8",
  ),
) as SepoliaDeployment;

const lenderName = requiredText("PRISM_LENDER_POSITION_TOKEN_NAME");
const lenderSymbol = requiredText("PRISM_LENDER_POSITION_TOKEN_SYMBOL");
const borrowerName = requiredText("PRISM_BORROWER_POSITION_TOKEN_NAME");
const borrowerSymbol = requiredText("PRISM_BORROWER_POSITION_TOKEN_SYMBOL");
if (lenderName === borrowerName || lenderSymbol === borrowerSymbol) {
  throw new Error(
    "lender and borrower position tokens must have distinct names and symbols",
  );
}

const { ethers } = await network.create();
const chain = await ethers.provider.getNetwork();
if (
  deployment.chainId !== SEPOLIA_CHAIN_ID ||
  chain.chainId.toString() !== SEPOLIA_CHAIN_ID
) {
  throw new Error(
    `position-token deployment is restricted to Sepolia chain ${SEPOLIA_CHAIN_ID}; manifest=${deployment.chainId}, connected=${chain.chainId}`,
  );
}

const prismPool = ethers.getAddress(deployment.prismPool);
const multisig = ethers.getAddress(deployment.multisig);
await Promise.all([
  requireContract("PrismPool", prismPool),
  requireContract("multisig", multisig),
]);

const lenderPositionToken = await deployPositionToken(
  "lender",
  lenderName,
  lenderSymbol,
);
const borrowerPositionToken = await deployPositionToken(
  "borrower",
  borrowerName,
  borrowerSymbol,
);

console.log(
  JSON.stringify(
    {
      chainId: SEPOLIA_CHAIN_ID,
      prismPool,
      owner: multisig,
      lenderPositionToken,
      borrowerPositionToken,
    },
    null,
    2,
  ),
);

async function deployPositionToken(
  role: "lender" | "borrower",
  name: string,
  symbol: string,
) {
  const token = await ethers.deployContract("PositionToken", [name, symbol]);
  await token.waitForDeployment();
  const address = await token.getAddress();

  const minterReceipt = await (await token.setMinter(prismPool, true)).wait();
  if (minterReceipt === null || minterReceipt.status !== 1) {
    throw new Error(`authorizing PrismPool on ${role} position token failed`);
  }

  const ownershipReceipt = await (
    await token.transferOwnership(multisig)
  ).wait();
  if (ownershipReceipt === null || ownershipReceipt.status !== 1) {
    throw new Error(`transferring ${role} position-token ownership failed`);
  }

  const [owner, poolIsMinter] = await Promise.all([
    token.owner(),
    token.isMinter(prismPool),
  ]);
  if (ethers.getAddress(owner) !== multisig || !poolIsMinter) {
    throw new Error(`${role} position-token post-deployment checks failed`);
  }

  console.log(`${role} position token: ${address}`);
  return address;
}

async function requireContract(name: string, address: string) {
  if ((await ethers.provider.getCode(address)) === "0x") {
    throw new Error(`${name} has no deployed code at ${address}`);
  }
}

function requiredText(name: string) {
  const value = process.env[name]?.trim();
  if (!value) throw new Error(`${name} is required`);
  if (value.length > 64)
    throw new Error(`${name} must not exceed 64 characters`);
  return value;
}
