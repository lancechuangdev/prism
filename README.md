# Prism

Prism is a fixed-rate lending protocol with an accompanying indexing and API
backend. The repository contains two projects:

- [`protocol`](./protocol) contains the Solidity contracts, Hardhat
  configuration, and contract tests.
- [`backend`](./backend) contains the Go API and scheduler, plus MySQL and Redis
  infrastructure for indexed data and price caching.

## Current integration status

The protocol and backend are currently developed in the same repository but
are not connected at runtime.

- The contracts implement and test the lending lifecycle on Hardhat networks.
- The backend's `chain.Reader` abstraction is ready for a contract adapter, but
  both backend executables currently use an in-process `DemoReader`.
- Backend prices currently come from `DemoProvider`, wrapped by a Redis-backed
  cache.

A production integration still needs a chain RPC reader, deployed contract
addresses and a real oracle or price-provider adapter.

## System architecture

```mermaid
flowchart LR
  Users[Protocol users] -->|transactions| Pool[PrismPool]
  Admins[Protocol admins] --> MultiSig[ThresholdMultiSig]
  MultiSig -->|approved calls| Pool
  Pool --> Positions[PositionToken contracts]
  Pool --> Oracle[Price oracle]
  Pool --> DEX[DEX swap adapter]

  Client[Frontend, admin, or curl] -->|HTTP| API[Backend API]

  subgraph Backend[Backend runtime]
    API --> Services[Auth, chain, price, and multisig services]
    API -->|startup| Sync[Pool synchronization]
    Scheduler[Scheduler] --> Sync
    Services --> Repository[Repository interface]
    Sync --> Repository
    Services --> APICachedPrice[API cached price provider]
    Scheduler --> SchedulerCachedPrice[Scheduler cached price provider]
  end

  Repository --> MySQL[(MySQL)]
  APICachedPrice <--> Redis[(Redis)]
  SchedulerCachedPrice <--> Redis
  APICachedPrice --> DemoPrice[Demo price provider]
  SchedulerCachedPrice --> DemoPrice
  DemoReader[Demo chain reader] --> Sync

  Pool -.->|future RPC reader| Sync
```

The dotted path represents the intended contract integration. In the current
backend, `DemoReader` supplies pool and token snapshots instead.

## Protocol

The protocol is built with Solidity 0.8.28, Hardhat 3, ethers 6, and
OpenZeppelin Contracts.

### Contracts

| Contract | Responsibility |
| --- | --- |
| `PrismPool` | Creates lending pools and manages deposits, settlement, claims, repayment, liquidation, and withdrawals. |
| `PositionToken` | ERC-20 receipt token minted and burned by approved pool contracts. |
| `MockOracle` | Owner-controlled token prices for local development and tests. |
| `FixedRateSwap` | Deterministic exact-input and exact-output swaps for tests. |
| `ThresholdMultiSig` | Threshold approval and execution of administrative calls. |

### Lending lifecycle

```mermaid
stateDiagram-v2
  [*] --> FUNDING: createPool
  FUNDING --> FUNDING: lender and collateral deposits
  FUNDING --> ACTIVE: settle with matched liquidity
  FUNDING --> CANCELLED: settlement cannot match both sides
  ACTIVE --> ACTIVE: refunds and position claims
  ACTIVE --> REPAID: repay at maturity through DEX
  ACTIVE --> LIQUIDATED: liquidate an undercollateralized pool
  REPAID --> [*]: position holders withdraw
  LIQUIDATED --> [*]: position holders withdraw
```

During funding, lenders deposit the lend token and borrowers deposit
collateral. Settlement matches the two sides using oracle prices. Lenders and
borrowers can then claim position tokens, with borrowers also receiving the
matched loan. At maturity, repayment swaps enough collateral for the required
lend tokens. An undercollateralized active pool can instead be liquidated.

### Run the contract tests

```bash
cd protocol
npm install
npm test
```

The test command uses the optimized Hardhat production build profile. The test
suite covers pool creation, deposits, settlement, refunds, claims, repayment,
liquidation, position tokens, the mock oracle, fixed-rate swaps, and multisig
administration.

Hardhat also defines simulated L1 and OP networks and an HTTP Sepolia network.
Sepolia requires `SEPOLIA_RPC_URL` and `SEPOLIA_PRIVATE_KEY`, but a Prism
deployment module still needs to be added before deploying the system.

See [`protocol/README.md`](./protocol/README.md) for the contract-focused
lifecycle notes.

## Backend

The backend is a Go module with two executables:

- `cmd/api` performs an initial pool sync, serves public and protected HTTP
  endpoints, and shuts down gracefully on `SIGINT` or `SIGTERM`.
- `cmd/scheduler` performs an initial sync and repeats it according to
  `PRISM_SYNC_INTERVAL`.

Both processes can use an in-memory repository or MySQL. They use Redis for
price caching and currently fetch cache misses from the demo price provider.
When both processes use MySQL, the scheduler's indexed snapshots are visible
to the API.

### Run the complete backend stack

Docker Compose is the shortest path because it supplies MySQL and Redis:

```bash
cd backend
docker compose up --build
```

This starts four containers:

| Service | Purpose | Host port |
| --- | --- | --- |
| `api` | HTTP API | `8080` |
| `scheduler` | Periodic pool and price synchronization | none |
| `mysql` | Persistent indexed state | `3306` |
| `redis` | Price cache | `6379` |

Check the running API:

```bash
curl http://localhost:8080/healthz
curl "http://localhost:8080/api/v1/poolBaseInfo?chainId=97"
curl "http://localhost:8080/api/v1/price?symbol=PRM"
```

Stop the stack while preserving MySQL data:

```bash
docker compose down
```

To remove the MySQL volume as well:

```bash
docker compose down -v
```

### API summary

| Method | Path | Access |
| --- | --- | --- |
| `GET` | `/healthz` | Public |
| `GET` | `/api/v1/poolBaseInfo?chainId=97` | Public |
| `GET` | `/api/v1/poolDataInfo?chainId=97` | Public |
| `GET` | `/api/v1/token?chainId=97` | Public |
| `GET` | `/api/v1/price?symbol=PRM` | Public |
| `POST` | `/api/v1/user/login` | Public |
| `POST` | `/api/v1/user/logout` | Bearer token |
| `GET` | `/api/v1/admin/session` | Bearer token |
| `POST` | `/api/v1/pool/setMultiSign` | Bearer token |
| `POST` | `/api/v1/pool/getMultiSign` | Bearer token |

The Compose configuration uses development credentials:

```text
username: admin
password: password
```

Pass the token returned by login as:

```text
Authorization: Bearer <tokenId>
```

Do not use the Compose credentials or token secret in a public environment.

### Backend configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PRISM_ENV` | `local` | Logging/runtime environment. |
| `PRISM_API_PORT` | `8080` | API listen port. |
| `PRISM_API_VERSION` | `1` | Version used in the `/api/v{n}` route prefix. |
| `PRISM_CHAIN_ID` | `97` | Chain identifier used for indexed records. |
| `PRISM_SYNC_INTERVAL` | `2m` | Scheduler synchronization interval. |
| `PRISM_STORE` | `memory` | Repository driver: `memory` or `mysql`. |
| `PRISM_MYSQL_DSN` | empty | MySQL connection string. |
| `PRISM_REDIS_ADDR` | `127.0.0.1:6379` | Redis address. |
| `PRISM_REDIS_PASSWORD` | empty | Redis password. |
| `PRISM_REDIS_DB` | `0` | Redis database number. |
| `PRISM_PRICE_SYMBOL` | `PRM` | Symbol refreshed by the scheduler. |
| `PRISM_PRICE_CACHE_TTL` | `30s` | Price-cache lifetime. |
| `PRISM_ADMIN_USERNAME` | `admin` | Development admin username. |
| `PRISM_ADMIN_PASSWORD` | `password` | Development admin password. |
| `PRISM_TOKEN_SECRET` | `local-development-secret` | HMAC token-signing secret. |
| `PRISM_TOKEN_TTL` | `1h` | Authentication token lifetime. |

See [`backend/README.md`](./backend/README.md) for endpoint payloads, storage
details, and the backend's implementation history.

### Run backend tests

```bash
cd backend
go test ./...
```

## Repository layout

```text
.
├── protocol/
│   ├── contracts/       Solidity protocol contracts
│   ├── test/            Hardhat integration and unit tests
│   ├── ignition/        Hardhat Ignition modules
│   └── scripts/         Network interaction scripts
└── backend/
    ├── cmd/api/         API executable
    ├── cmd/scheduler/   Scheduler executable
    ├── internal/        Services, adapters, storage, cache, and HTTP server
    ├── Dockerfile       Multi-stage Go image
    └── docker-compose.yml
```

## Development checks

Run both project test suites before submitting changes:

```bash
cd protocol
npm test

cd ../backend
go test ./...
```
