import { readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { network } from "hardhat";

type LocalDeployment = {
  chainId: string;
  lendToken: string;
  collateralToken: string;
};

const LOCAL_CHAIN_ID = "31337";
const protocolRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
const deployment = JSON.parse(
  await readFile(path.join(protocolRoot, "deployments", "local.json"), "utf8"),
) as LocalDeployment;

const { ethers } = await network.create();
const chain = await ethers.provider.getNetwork();
if (
  deployment.chainId !== LOCAL_CHAIN_ID ||
  chain.chainId.toString() !== LOCAL_CHAIN_ID
) {
  throw new Error(
    `local wallet funding is restricted to chain ${LOCAL_CHAIN_ID}; deployment=${deployment.chainId}, connected=${chain.chainId}`,
  );
}

const recipients = parseRecipients(
  process.env.PRISM_WALLET_ADDRESSES ?? process.env.PRISM_WALLET_ADDRESS,
);
const pUSDAmount = parseTokenAmount(
  "PRISM_LOCAL_PUSD_AMOUNT",
  process.env.PRISM_LOCAL_PUSD_AMOUNT ?? "10000",
);
const pETHAmount = parseTokenAmount(
  "PRISM_LOCAL_PETH_AMOUNT",
  process.env.PRISM_LOCAL_PETH_AMOUNT ?? "100",
);
const nativeAmount = parseTokenAmount(
  "PRISM_LOCAL_ETH_AMOUNT",
  process.env.PRISM_LOCAL_ETH_AMOUNT ?? "10",
);

const signers = await ethers.getSigners();
const nativeFunder = signers[0];
if (nativeFunder === undefined) {
  throw new Error("the local node does not expose a signer for native funding");
}

const pUSD = await ethers.getContractAt("PositionToken", deployment.lendToken);
const pETH = await ethers.getContractAt(
  "PositionToken",
  deployment.collateralToken,
);

await mintToRecipients(pUSD, "pUSD", pUSDAmount);
await mintToRecipients(pETH, "pETH", pETHAmount);

if (nativeAmount > 0n) {
  for (const recipient of recipients) {
    const transaction = await nativeFunder.sendTransaction({
      to: recipient,
      value: nativeAmount,
    });
    await requireSuccessfulReceipt(transaction.hash);
    console.log(
      `Sent ${ethers.formatEther(nativeAmount)} ETH to ${recipient} for local gas`,
    );
  }
}

for (const recipient of recipients) {
  const [pUSDBalance, pETHBalance, nativeBalance] = await Promise.all([
    pUSD.balanceOf(recipient),
    pETH.balanceOf(recipient),
    ethers.provider.getBalance(recipient),
  ]);
  console.log(
    JSON.stringify(
      {
        wallet: recipient,
        pUSD: ethers.formatEther(pUSDBalance),
        pETH: ethers.formatEther(pETHBalance),
        nativeETH: ethers.formatEther(nativeBalance),
      },
      null,
      2,
    ),
  );
}

async function mintToRecipients(
  token: typeof pUSD,
  symbol: string,
  amount: bigint,
) {
  if (amount === 0n) return;

  const ownerAddress = ethers.getAddress(await token.owner());
  const owner = signers.find(
    (signer) => signer.address.toLowerCase() === ownerAddress.toLowerCase(),
  );
  if (owner === undefined) {
    throw new Error(
      `${symbol} owner ${ownerAddress} is not an unlocked local signer`,
    );
  }

  const ownerWasMinter = await token.isMinter(ownerAddress);
  if (!ownerWasMinter) {
    await requireSuccessfulReceipt(
      (await token.connect(owner).setMinter(ownerAddress, true)).hash,
    );
  }

  try {
    for (const recipient of recipients) {
      await requireSuccessfulReceipt(
        (await token.connect(owner).mint(recipient, amount)).hash,
      );
      console.log(
        `Minted ${ethers.formatEther(amount)} ${symbol} to ${recipient}`,
      );
    }
  } finally {
    if (!ownerWasMinter) {
      await requireSuccessfulReceipt(
        (await token.connect(owner).setMinter(ownerAddress, false)).hash,
      );
    }
  }
}

function parseRecipients(value: string | undefined) {
  if (value === undefined || value.trim() === "") {
    throw new Error(
      "set PRISM_WALLET_ADDRESS or comma-separated PRISM_WALLET_ADDRESSES",
    );
  }

  const recipients = [
    ...new Set(
      value
        .split(",")
        .map((address) => address.trim())
        .filter(Boolean)
        .map((address) => ethers.getAddress(address)),
    ),
  ];
  if (recipients.length === 0) {
    throw new Error("at least one wallet address is required");
  }
  if (recipients.some((address) => address === ethers.ZeroAddress)) {
    throw new Error("the zero address cannot receive local test funds");
  }
  return recipients;
}

function parseTokenAmount(name: string, value: string) {
  if (!/^\d+(\.\d+)?$/.test(value)) {
    throw new Error(`${name} must be a non-negative decimal token amount`);
  }
  const amount = ethers.parseEther(value);
  if (amount < 0n) {
    throw new Error(`${name} must not be negative`);
  }
  return amount;
}

async function requireSuccessfulReceipt(hash: string) {
  const receipt = await ethers.provider.waitForTransaction(hash);
  if (receipt === null || receipt.status !== 1) {
    throw new Error(`transaction ${hash} failed`);
  }
}
