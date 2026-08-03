# Prism Frontend

The Prism frontend is the user and operator interface for the Prism fixed-rate lending protocol. It should combine indexed data from the Go API with live wallet and contract data so users can discover pools, lend or borrow, manage their positions, and participate in multisig administration.

## Development

The frontend uses React, TypeScript, Vite, Viem, and Zod. Copy `.env.example` to `.env.local` and replace its placeholder contract addresses with addresses from the deployment manifest for the chain you intend to use.

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

Run the complete local quality gate with:

```bash
npm run format:check
npm run lint
npm test
npm run build
```

Runtime configuration is validated before the application starts. The app supports injected EIP-1193 wallets and remains available in read-only mode when a wallet is absent or disconnected.

The application polls backend readiness and live chain state every 30 seconds and when a backgrounded tab becomes visible. The global status panel separates API dependencies, RPC chain identity, contract reachability, pause state, and the configured liquidator address. A configured liquidator is not reported as a healthy keeper because the backend does not expose scheduler heartbeat or last-action status.

During local development, `VITE_PRISM_API_URL=/prism-api` uses Vite's proxy to reach the backend at `http://localhost:8080` without requiring backend CORS headers. Use an absolute API URL for deployments where the frontend and API are hosted on different origins and the API explicitly allows the frontend origin.

Lending and collateral deposits require the connected wallet to hold the pool's local ERC-20 tokens and native ETH for gas. The frontend reads balances, allowances, protocol minimums, pause state, and pool state directly from the configured RPC; it requests an exact-amount approval, simulates the selected deposit, and waits for its receipt before reporting success. Indexed marketplace values refresh every 30 seconds to follow the backend scheduler.

Fund a connected local wallet before testing deposits:

```bash
cd ../protocol
PRISM_WALLET_ADDRESS=0x... npm run fund:local
```

The portfolio reads wallet positions and protocol events directly from the chain. Set `VITE_PRISM_DEPLOYMENT_BLOCK` to the PrismPool deployment block in hosted environments so activity queries do not scan from genesis; local development defaults to block `0`. Claim and refund buttons are shown only when the live pool state and wallet records make the action eligible, and each transaction is simulated before broadcast.

## Product areas

### Pool marketplace

The marketplace is the primary pool-discovery experience. It should show:

- lending and collateral assets;
- fixed interest rate and estimated lender return;
- funding progress against the maximum lending supply;
- collateralization and liquidation ratios;
- settlement and maturity times, including countdowns;
- pool state: `FUNDING`, `ACTIVE`, `REPAID`, `LIQUIDATED`, or `CANCELLED`;
- token prices, price freshness, and current collateral health; and
- filters for asset, state, maturity, and yield.

Each pool should have a detail page containing its terms, lifecycle, risk explanation, contract addresses, and relevant on-chain activity.

### Lender workflow

A lender should be able to:

- connect a wallet and switch to the configured network;
- enter a lending amount and review the projected return at maturity;
- approve the exact token amount and deposit it into a funding pool;
- track the portion of the deposit used during settlement;
- refund unused lending funds;
- claim lender position tokens; and
- track redemption availability when the protocol supports redemption.

The interface should present approval, deposit, confirmation, refund, and claim as separate transaction states.

### Borrower workflow

A borrower should be able to:

- enter a collateral amount and preview the estimated loan;
- review collateralization and liquidation thresholds;
- approve and deposit collateral;
- claim the borrower position token and loan after settlement;
- refund unused collateral; and
- monitor the position's collateral health as prices change.

Liquidation risk should be communicated with concrete prices and amounts in addition to a simple health indicator.

### Portfolio dashboard

After connecting a wallet, the user should have a consolidated view of:

- lending and collateral deposits;
- active lender and borrower position tokens;
- pending refunds and available claims;
- supplied principal and expected return;
- borrowed amounts and current collateral value;
- transaction history; and
- positions requiring action.

Wallet-specific position data is not currently exposed by the backend, so the frontend will need live contract reads or a future indexed positions API.

### Multisig governance console

The authenticated operator experience should support:

- viewing multisig owners and the current approval threshold;
- preparing pool creation, settlement, repayment, and liquidation proposals;
- preparing owner and threshold changes;
- inspecting proposal hashes and on-chain approval counts;
- signing and broadcasting approval transactions from connected owner wallets;
- executing proposals once they reach quorum; and
- copying or exporting unsigned transaction data for external and hardware wallet workflows.

Every proposal should include a human-readable summary of its effects, target contract, parameters, and expected state transition instead of showing only calldata.

### Protocol and service status

The frontend should surface operational information that materially affects a user's decisions:

- API, chain RPC, indexer, and price-feed health;
- the last data update or indexed block;
- price freshness;
- automatic liquidation keeper status;
- contract pause state; and
- clear degraded-data warnings.

### Transaction safety

Safety is part of the core user experience. The frontend should provide:

- chain ID and deployed-contract verification;
- exact token allowances by default;
- transaction simulation where supported;
- human-readable contract and wallet errors;
- pending, replaced, failed, and confirmed transaction states;
- block-explorer links;
- decimal-safe token and rate formatting;
- stale-price warnings; and
- prominent test-network and unsupported-network indicators.

The UI must distinguish indexed API data from live chain data and disclose when either source is stale or unavailable.

## Current integrations

The Go backend currently provides public endpoints for pool base data, pool settlement data, tokens, and prices. It also provides public multisig reads and protected session and proposal-preparation endpoints. See [`../backend/README.md`](../backend/README.md) for the complete API contract.

User transactions such as token approval, lending deposits, collateral deposits, refunds, and position claims should be submitted directly from the connected wallet to the deployed contracts. Deployment addresses and chain configuration should come from an environment-specific, validated manifest.

Production operator authentication uses Cognito access tokens with the `prism/proposals.write` and `prism/admin.read` scopes. Local development can use the backend's local login flow.

The governance console reads multisig owners, thresholds, and proposal status without authentication. Preparing proposals requires operator authentication through `VITE_PRISM_AUTH_MODE=local` or Cognito Hosted UI with PKCE. Cognito deployments must set `VITE_PRISM_COGNITO_DOMAIN`, `VITE_PRISM_COGNITO_CLIENT_ID`, `VITE_PRISM_COGNITO_REDIRECT_URI`, and `VITE_PRISM_COGNITO_LOGOUT_URI`. Access tokens are kept in session storage, and the backend remains responsible for validating every token and scope. A connected owner wallet is separately required to simulate, approve, or execute an unsigned multisig transaction.

Privacy-preserving product telemetry is disabled by default. Set `VITE_PRISM_TELEMETRY_ENDPOINT` only after the receiving system and privacy notice are approved. Events use an explicit allowlist, contain no wallet addresses, amounts, calldata, hashes, query strings, credentials, or tokens, and are suppressed when the browser sends Do Not Track. See [`RELEASE_RUNBOOK.md`](RELEASE_RUNBOOK.md) for deployment, target-network testing, rollback, incident, and support procedures.

## Known protocol limitations

The frontend must disclose protocol functionality that is incomplete or unsafe for production:

- position-token redemption exists but is unsafe when one position-token contract is reused across pools, so production pools require pool-specific position tokens or accounting;
- refunds from `CANCELLED` pools are not currently implemented;
- repayment-interest units require correction and testing before real-asset use; and
- preparing a multisig proposal does not guarantee that its eventual on-chain execution will succeed.

Affected actions should be hidden or disabled with a clear explanation. These limitations should also appear in pool risk disclosures before a user deposits.

## Roadmap

### Phase 0: Foundation

- Select the frontend stack and establish formatting, linting, type checking, unit testing, and CI.
- Define environment configuration for the API URL, chain ID, RPC URL, explorer URL, and verified contract addresses.
- Generate or import typed clients for the HTTP API and contract ABIs.
- Implement decimal-safe amount, rate, timestamp, and address utilities.
- Establish the application shell, routing, responsive layout, design tokens, error boundaries, and accessible UI primitives.
- Add wallet connection, network switching, and read-only fallback behavior.

**Exit criteria:** the application runs locally, connects to the configured chain and API, validates its configuration, and passes automated checks.

### Phase 1: Read-only marketplace

- [x] Build pool listing, filtering, sorting, loading, empty, and error states.
- [x] Build pool detail pages and combine pool base and settlement data.
- [x] Display token metadata, available indexed prices, data freshness, and funding progress.
- [x] Translate numeric pool states into user-facing lifecycle labels.
- [x] Add settlement and maturity countdowns, contract links, and risk disclosures.
- [x] Label indexed API values explicitly; live contract values are deferred until both sources are presented together.

**Exit criteria:** a user can understand and compare every indexed pool without connecting a wallet.

### Phase 2: Lending and borrowing

- [x] Add lender and borrower calculators with minimum, balance, allowance, supply, and deadline validation.
- [x] Implement exact-amount ERC-20 approval flows.
- [x] Implement `depositLend` and `depositBorrow` transactions.
- [x] Simulate approvals and deposits against live contract state before broadcast.
- [x] Track wallet prompts, submitted transactions, replacements, confirmations, failures, and explorer links.
- [x] Refresh balances and request indexed pool data after confirmation.

**Exit criteria:** a connected user can safely fund either side of an eligible pool and see the confirmed result.

### Phase 3: Portfolio and lifecycle actions

- [x] Discover wallet-specific lender and borrower positions from contracts.
- [x] Build the portfolio summary and per-position views.
- [x] Implement excess-lend and excess-collateral refunds.
- [x] Implement lender-position and borrower-position/loan claims.
- [x] Show action eligibility based on live pool state and the user's on-chain records.
- [x] Add contract-event history and actionable maturity, claim, refund, and collateral-risk notifications.

**Exit criteria:** users can find their positions and complete every currently supported post-settlement user action from one dashboard.

### Phase 4: Multisig governance

- [x] Add local and Cognito authentication flows with scope-aware action protection.
- [x] Build multisig owner, threshold, and proposal detail views.
- [x] Add guided forms for all supported proposal operation types.
- [x] Present decoded proposal intent and validate inputs before submission.
- [x] Integrate owner-wallet approval and execution transactions.
- [x] Show quorum progress and refresh proposal state after confirmations.
- [x] Support copying and exporting unsigned transactions.

**Exit criteria:** authorized multisig owners can prepare, review, approve, and execute supported protocol operations without manually constructing calldata.

### Phase 5: Safety, observability, and release readiness

- [x] Add API and dependency health indicators and stale-data handling.
- [x] Surface contract pause state and the configured liquidator while disclosing that keeper heartbeat is unavailable.
- [x] Add allowlisted, identifier-free, Do-Not-Track-aware product telemetry that is disabled by default.
- [x] Complete keyboard navigation, screen-reader status, focus, reduced-motion, and responsive navigation improvements.
- [ ] Add target-network contract-fork and end-to-end tests for every critical financial and governance flow; unit and component coverage is active in CI.
- [x] Add baseline CSP and production dependency checks, redact sensitive telemetry, and document authentication and transaction boundaries.
- [x] Publish deployment, rollback, incident, and support runbooks.

**Exit criteria:** critical user and governance journeys are tested against the target deployment and the frontend is safe to release to its intended network.

### Phase 6: Protocol-completeness follow-up

These items depend on protocol or backend work and should not be enabled until the underlying capability is implemented and audited:

- safe lender and borrower position-token redemption after the protocol accounting redesign;
- refunds for cancelled pools;
- indexed wallet positions and historical performance;
- richer event and transaction history APIs;
- keeper status and last-action reporting; and
- notification delivery for maturity, claims, and liquidation risk.

## Definition of done for financial actions

A frontend transaction feature is complete only when it:

- validates chain, account, pool state, deadline, balance, allowance, and amount;
- explains the assets and state changes before signature;
- handles rejection, revert, replacement, timeout, and confirmation;
- links to the submitted transaction on the configured explorer;
- refreshes affected API and live contract data; and
- has automated coverage for its successful path and important failure modes.
