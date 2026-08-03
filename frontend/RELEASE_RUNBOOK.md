# Prism frontend release runbook

## Release prerequisites

Release only a reviewed commit with a clean working tree. Pin the target chain ID, RPC URL, deployment block, pool address, multisig address, API origin, Cognito domain, and public app-client ID to the same environment. Confirm the deployment manifest rather than copying addresses from a browser or chat message.

Run the local gate:

```bash
cd frontend
npm ci
npm run format:check
npm run lint
npm test
npm run security:check
npm run build
```

Review `npm audit` findings rather than applying force upgrades. A high or critical production dependency finding blocks release unless a written risk acceptance identifies exploitability, exposure, owner, and remediation date.

## Runtime configuration

Provide all required `VITE_PRISM_*` values at build time. Vite embeds these values in public JavaScript, so never put secrets, private keys, confidential app-client secrets, session tokens, or privileged API credentials in them.

Use a public Cognito app client without a secret. Register the exact production governance callback and logout URLs, use authorization code flow with PKCE, and allow `openid`, `profile`, `prism/proposals.write`, and `prism/admin.read`.

Telemetry remains off when `VITE_PRISM_TELEMETRY_ENDPOINT` is empty. When enabled, the receiver must accept anonymous first-party events, enforce a short retention period, avoid IP enrichment, and document its lawful basis and user-facing privacy notice. Prism sends only allowlisted coarse metadata and honors browser Do Not Track.

## Deployment

Build an immutable artifact once and promote the same artifact through environments when configuration permits. If environment values are compiled into separate artifacts, record the source commit, build identity, dependency lockfile hash, and complete non-secret configuration alongside each artifact.

Serve the application over HTTPS. Configure the host or CDN to return the CSP from `index.html` as an HTTP response header and add `frame-ancestors 'none'`, because `frame-ancestors` is not enforced from a meta tag. Also set `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, and an appropriate `Permissions-Policy`. Restrict production `connect-src` to the deployed API, RPC, Cognito, telemetry, and explorer origins instead of retaining development wildcards.

After deployment, verify `/`, `/pools`, `/portfolio`, and `/governance` through direct navigation and client-side navigation. Confirm the global status panel reports the intended chain, API dependencies, pool contract, pause state, and configured liquidator address.

## Target-network smoke test

Use dedicated test accounts and minimal-value assets. Do not test production releases with an unrestricted operator or treasury wallet.

1. Load the marketplace without a wallet and confirm indexed timestamps and stale warnings.
2. Connect a non-owner wallet, confirm the chain, and validate that ineligible governance signing is blocked.
3. Simulate one lending or collateral flow without broadcasting, then complete the approved minimal-value test transaction when authorized.
4. Verify the portfolio discovers the resulting position and activity from the configured deployment block.
5. Authenticate through Cognito, prepare a harmless reviewed proposal, export its JSON, and verify its target, calldata, nonce, and chain ID independently.
6. Approve with the required test owner wallets and execute only when the intended test action and quorum are confirmed.
7. Confirm refresh behavior after transaction confirmation and confirm no access token, wallet address, calldata, or transaction hash appears in telemetry payloads or application logs.

Contract-fork and end-to-end suites must run against the exact protocol artifact before enabling real-value actions. Until those suites exist and pass, treat the frontend as test-network software.

## Rollback

Keep the previous immutable frontend artifact and configuration available. Roll back immediately for incorrect contract addresses, wrong chain configuration, broken authentication, unsafe transaction construction, CSP regressions that disable required connections, or misleading health and pause state.

Rollback changes only the frontend artifact; it cannot reverse confirmed blockchain transactions. For an on-chain safety incident, pause the protocol through the audited multisig procedure when authorized and follow the protocol incident process.

After rollback, invalidate the CDN cache, verify the served artifact identity, repeat read-only smoke tests, and record the trigger, decision maker, timestamps, affected users, and follow-up owner.

## Incident response

Classify incidents as transaction safety, authentication, data integrity, availability, or privacy. Preserve browser console output only after redacting URLs, authorization headers, tokens, wallet addresses when unnecessary, calldata, and user-entered amounts. Never ask users to share seed phrases, private keys, complete JWTs, temporary Cognito authorization codes, or exported signed transactions.

For incorrect or stale indexed data, keep live contract validation authoritative for actions, display degraded status, and disable affected decisions if live preflight is unavailable. For RPC or API outages, avoid repeated transaction submission and reconcile submitted hashes on-chain before retrying. For suspected token exposure, revoke the Cognito session where supported, expire affected credentials, and investigate backend authorization logs using request IDs rather than tokens.

## Support checklist

Ask for the frontend release identifier, route, configured chain ID, approximate timestamp, browser and wallet version, and backend request ID when displayed. Users may share a public transaction hash when necessary, but support must explain that it reveals their public wallet activity.

Direct users to verify the selected network, contract addresses, wallet account, token decimals, and transaction status in a trusted explorer. For rejected wallet requests, distinguish user rejection from simulation failure and on-chain revert. Escalate any mismatch between the human-readable review and wallet calldata as a transaction-safety incident.

## Known observability limits

The frontend can verify that a nonzero liquidator address is authorized on-chain, but the backend does not expose scheduler heartbeat, last successful keeper cycle, or last liquidation attempt. “Address configured” must not be interpreted as keeper health. Indexed timestamps describe stored records, not a chain head or finalized indexer checkpoint.

Browser portfolio notifications run only while Prism is open and are deduplicated for the current tab session. There is no service worker, push subscription, notification backend, or guaranteed background delivery. Cancelled-pool refunds remain unavailable because the deployed protocol has no compatible refund function.
