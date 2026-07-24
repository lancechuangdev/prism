# prism smart contracts

Fixed-rate lending protocol distilled from previous DeFi experience.

## Lending pool transitions

```mermaid
flowchart LR
    F[FUNDING] -->|Owner settles successfully| A[ACTIVE]
    F -->|Settlement fails| C[CANCELLED]

    A -->|Owner liquidates| L[LIQUIDATED]
    A -->|Refund excess lend or collateral<br/>no state change| A
    A -->|Claim positions and borrower loan<br/>no state change| A

    A -->|Owner calls repayPool at maturity| P[Repayment execution]
    P -->|1. Approves and sends matched collateral| D[DEX Swap]
    D -->|2. Sends the exact required lend-token amount| P
    P -->|3. Records received lend tokens<br/>and remaining collateral| R[REPAID]
```

Refunding excess lender funds or borrower collateral does not change the pool
state. Claiming lender/borrower position tokens and the borrower loan also
leaves the pool `ACTIVE`. Refunds from a `CANCELLED` pool are not currently
implemented.

At maturity, the owner calls `repayPool()`. The pool asks the configured DEX
for the collateral required to obtain the exact lend-token repayment amount,
approves that collateral, and executes the swap. The recovered lend tokens and
remaining collateral are recorded before the pool moves to `REPAID`.

## Local deployment

Install the protocol dependencies if needed:

```bash
npm install
```

Start a persistent local Hardhat JSON-RPC node in the first terminal. For a backend running directly on the host:

```bash
npm run node
```

For a backend running in Docker, Hardhat must also accept connections arriving through the host's Docker bridge interface:

```bash
npx hardhat node --hostname 0.0.0.0
```

Docker's `host.docker.internal` hostname resolves to the host gateway, commonly `172.17.0.1` on Linux. A Hardhat node bound only to its default `127.0.0.1` address cannot accept connections through that gateway. Binding to `0.0.0.0` is intended for local development and may expose the development
node to the local network, depending on the host firewall.

Keep that process running. In a second terminal, deploy and seed the protocol:

```bash
npm run deploy:local
```

The deployment script uses the optimized production build and deploys:

- `MockOracle`
- `FixedRateSwap`
- lend and collateral tokens
- lender and borrower position tokens
- `PrismPool`

It authorizes `PrismPool` to mint the position tokens, configures local oracle prices, and creates one pool in the `FUNDING` state.

The command prints a JSON object containing the values needed by an RPC client and saves the same object to the generated `deployments/local.json` file:

```json
{
  "rpcUrl": "http://127.0.0.1:8545",
  "chainId": "31337",
  "prismPool": "0x..."
}
```

Contract addresses belong to the running local node. Restarting `npm run node` resets its chain state, so run `npm run deploy:local` again and use the newly generated addresses. The local deployment file is ignored by Git because its addresses are only valid for that running node.

To pass the generated pool address to the Docker backend:

```bash
export PRISM_POOL_ADDRESS="$(
  jq -r '.prismPool' protocol/deployments/local.json
)"
docker compose up --build
```

## Create a pool through the API

With the local node, API, and scheduler running, automate the complete non-custodial pool-creation workflow:

```bash
cd protocol
npm run create-pool:api
```

The script reads `deployments/local.json`, logs in to the backend, and calls `POST /api/v1/pools` to obtain unsigned `createPool` calldata. It verifies the target contract and chain, signs with the first local Hardhat account (the deployed pool owner), broadcasts the transaction, waits for confirmation, and then waits for the scheduler to make the new pool visible through `GET /api/v1/poolBaseInfo`.

The defaults match the local development configuration. Override them when needed:

```bash
PRISM_API_URL=http://127.0.0.1:8080 \
PRISM_ADMIN_USERNAME=admin \
PRISM_ADMIN_PASSWORD=password \
PRISM_INDEX_TIMEOUT_MS=90000 \
npm run create-pool:api
```

This script is for local development: the well-known Hardhat account owns the local deployment. In a real frontend, the user's wallet should validate and
sign the API's prepared transaction instead.

## Generate the Go contract bindings

The backend uses generated, typed Go bindings instead of maintaining contract ABIs by hand.

Run the following commands from the repository root.

Install the required command-line tools if they are not already available:

```bash
sudo apt install jq
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

Make sure the Go installation's binary directory is on `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Compile the contracts so Hardhat generates an up-to-date artifact:

```bash
cd protocol
npx hardhat build --build-profile production
cd ..
```

Extract the ABI arrays from the full Hardhat artifacts:

```bash
jq '.abi' \
  protocol/artifacts/contracts/pool/PrismPool.sol/PrismPool.json \
  > protocol/contracts/pool/PrismPool.abi.json
```

Create the Go package directory and generate both bindings:

```bash
mkdir -p backend/internal/contracts

abigen \
  --abi protocol/contracts/pool/PrismPool.abi.json \
  --pkg contracts \
  --type PrismPool \
  --out backend/internal/contracts/prism_pool.go
```

Format and verify the generated code:

```bash
go fmt -w backend/internal/contracts/prism_pool.go

cd backend
go test ./...
```

Repeat this process whenever either public contract ABI changes.
The `abigen` commands overwrite the existing generated Go files, which is expected; do not edit them manually.

Apply the same process to other contracts as needed.

## TODO: production asset support

Before accepting real assets such as USDT, WETH, or WBTC:

- normalize pool calculations for tokens with different decimals;
- replace the global 18-decimal minimums with per-token or per-pool minimums;
- use `SafeERC20` for non-standard ERC-20 return behavior;
- correct and test the repayment-interest units;
- replace `MockOracle` and `FixedRateSwap` with production integrations; and
- add explicit wrapped-asset UX, since native ETH and BTC are not ERC-20 tokens.
