import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type FormEvent,
} from 'react'
import { type Address, type Hex } from 'viem'

import { useAuth } from '../auth/AuthProvider'
import { Button } from '../components/Button'
import { config } from '../config/env'
import { createApiClient } from '../lib/api/client'
import type {
  MultisigConfig,
  PreparedProposal,
  ProposalStatus,
  PreparedTransaction,
} from '../lib/api/types'
import { publicClient } from '../lib/contracts/client'
import { formatAddress } from '../lib/format'
import {
  normalizeOperation,
  operationLabels,
  operationSummary,
  validateOperation,
} from '../lib/governance'
import { useWallet } from '../wallet/WalletProvider'

const api = createApiClient(config)
const operations = Object.keys(operationLabels)
const createPoolFields = [
  'settleTime',
  'maturityTime',
  'interestRate',
  'maxLendSupply',
  'collateralizationRatio',
  'lendToken',
  'collateralToken',
  'lenderPositionToken',
  'borrowerPositionToken',
  'liquidateRate',
]
const fieldLabels: Record<string, string> = {
  owner: 'Owner address',
  old_owner: 'Current owner address',
  new_owner: 'Replacement owner address',
  threshold: 'Required signatures',
  settleTime: 'Settlement time (Unix seconds)',
  maturityTime: 'Maturity time (Unix seconds)',
  interestRate: 'Fixed interest rate (contract units)',
  maxLendSupply: 'Maximum lending supply (smallest token unit)',
  collateralizationRatio: 'Collateralization ratio (contract units)',
  lendToken: 'Lending token address',
  collateralToken: 'Collateral token address',
  lenderPositionToken: 'Lender position token address',
  borrowerPositionToken: 'Borrower position token address',
  liquidateRate: 'Liquidation margin (contract units)',
  poolId: 'Pool ID',
  maxCollateralAmount: 'Maximum collateral (smallest token unit)',
}

function initialParams(type: string): Record<string, string | number> {
  if (type === 'replace_owner') return { old_owner: '', new_owner: '' }
  if (type === 'change_threshold') return { threshold: 1 }
  if (type === 'create_pool')
    return Object.fromEntries(createPoolFields.map((field) => [field, '']))
  if (['settle_pool'].includes(type)) return { poolId: '' }
  if (['repay_pool', 'liquidate_pool'].includes(type))
    return { poolId: '', maxCollateralAmount: '' }
  return { owner: '' }
}

function LoginPanel() {
  const auth = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  if (auth.status === 'authenticated')
    return (
      <div className="governance-login">
        <span>Signed in{auth.username ? ` as ${auth.username}` : ''}</span>
        <Button variant="quiet" onClick={() => void auth.logout()}>
          Sign out
        </Button>
      </div>
    )
  if (config.auth.mode === 'cognito')
    return (
      <div className="governance-login">
        <span>Operator authentication is required to prepare proposals.</span>
        <Button
          disabled={auth.status === 'authenticating'}
          onClick={() => void auth.loginCognito()}
        >
          Sign in with Cognito
        </Button>
        {auth.error && <p className="form-error">{auth.error}</p>}
      </div>
    )
  return (
    <form
      className="governance-login"
      onSubmit={(event) => {
        event.preventDefault()
        void auth.loginLocal(username, password)
      }}
    >
      <label>
        Username
        <input
          autoComplete="username"
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          required
        />
      </label>
      <label>
        Password
        <input
          type="password"
          autoComplete="current-password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          required
        />
      </label>
      <Button disabled={auth.status === 'authenticating'}>
        {auth.status === 'authenticating' ? 'Signing in…' : 'Operator sign in'}
      </Button>
      {auth.error && <p className="form-error">{auth.error}</p>}
    </form>
  )
}

function ProposalStatusView({ status }: { status: ProposalStatus }) {
  return (
    <section className="governance-panel">
      <div className="proposal-heading">
        <div>
          <p className="eyebrow">On-chain status</p>
          <h2>
            {status.executed
              ? 'Executed'
              : status.readyToExecute
                ? 'Ready to execute'
                : 'Collecting approvals'}
          </h2>
        </div>
        <strong>
          {status.approvalCount} / {status.threshold}
        </strong>
      </div>
      {!status.configurationValid && (
        <p className="warning-box">
          The owner configuration changed after this proposal was created.
          Prepare a new proposal.
        </p>
      )}
      <div className="owner-list">
        {status.owners.map((owner) => (
          <div key={owner.address}>
            <code>{formatAddress(owner.address)}</code>
            <span>{owner.approved ? 'Approved' : 'Pending'}</span>
          </div>
        ))}
      </div>
    </section>
  )
}

export function GovernancePage() {
  const auth = useAuth()
  const wallet = useWallet()
  const [multisig, setMultisig] = useState<MultisigConfig>()
  const [type, setType] = useState('add_owner')
  const [params, setParams] = useState<Record<string, string | number>>(
    initialParams(type),
  )
  const [nonce, setNonce] = useState(() => Date.now().toString())
  const [prepared, setPrepared] = useState<PreparedProposal>()
  const [proposalStatus, setProposalStatus] = useState<ProposalStatus>()
  const [lookupHash, setLookupHash] = useState('')
  const [busy, setBusy] = useState<string>()
  const [error, setError] = useState<string>()
  const operation = useMemo(
    () => normalizeOperation({ type, params }),
    [params, type],
  )
  const validation = useMemo(() => validateOperation(operation), [operation])

  const loadConfig = useCallback(
    () =>
      api
        .getMultisig()
        .then(({ data }) => setMultisig(data))
        .catch((cause: unknown) =>
          setError(
            cause instanceof Error
              ? cause.message
              : 'Could not load multisig configuration.',
          ),
        ),
    [],
  )
  const loadStatus = useCallback(async (hash: string) => {
    const { data } = await api.getProposalStatus(hash)
    setProposalStatus(data)
    return data
  }, [])
  useEffect(() => {
    void loadConfig()
  }, [loadConfig])

  async function prepare(event: FormEvent) {
    event.preventDefault()
    if (!auth.token || validation.length) return
    setBusy('prepare')
    setError(undefined)
    try {
      const { data } = await api.prepareProposal(
        String(config.chain.id),
        nonce,
        operation,
        auth.token,
      )
      setPrepared(data)
      setLookupHash(data.proposal.transactionHash)
      await loadStatus(data.proposal.transactionHash).catch(() => undefined)
    } catch (cause) {
      setError(
        cause instanceof Error ? cause.message : 'Proposal preparation failed.',
      )
    } finally {
      setBusy(undefined)
    }
  }

  async function broadcast(transaction: PreparedTransaction, action: string) {
    if (!wallet.account || !wallet.walletClient) return
    if (wallet.status !== 'connected') {
      await wallet.switchNetwork()
      return
    }
    if (
      !multisig?.owners.some(
        (owner) => owner.toLowerCase() === wallet.account?.toLowerCase(),
      )
    ) {
      setError('The connected wallet is not a multisig owner.')
      return
    }
    setBusy(action)
    setError(undefined)
    try {
      await publicClient.call({
        account: wallet.account,
        to: transaction.to as Address,
        data: transaction.data as Hex,
        value: BigInt(transaction.value),
      })
      const hash = await wallet.walletClient.sendTransaction({
        account: wallet.account,
        chain: undefined,
        to: transaction.to as Address,
        data: transaction.data as Hex,
        value: BigInt(transaction.value),
      })
      await publicClient.waitForTransactionReceipt({ hash })
      if (prepared) await loadStatus(prepared.proposal.transactionHash)
      await loadConfig()
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : `${action} failed.`)
    } finally {
      setBusy(undefined)
    }
  }

  function exportProposal() {
    if (!prepared) return
    const url = URL.createObjectURL(
      new Blob([JSON.stringify(prepared, null, 2)], {
        type: 'application/json',
      }),
    )
    const link = document.createElement('a')
    link.href = url
    link.download = `prism-proposal-${prepared.proposal.transactionHash}.json`
    link.click()
    URL.revokeObjectURL(url)
  }

  const ownerConnected = Boolean(
    wallet.account &&
    multisig?.owners.some(
      (owner) => owner.toLowerCase() === wallet.account?.toLowerCase(),
    ),
  )
  const currentOwner = proposalStatus?.owners.find(
    (owner) => owner.address.toLowerCase() === wallet.account?.toLowerCase(),
  )

  return (
    <section className="page governance-page">
      <header className="governance-hero">
        <div>
          <p className="eyebrow">Multisig governance</p>
          <h1>Operator console</h1>
          <p>
            Prepare human-readable protocol actions, collect owner approvals,
            and execute only after quorum.
          </p>
        </div>
        <LoginPanel />
      </header>
      {error && (
        <div className="warning-box" role="alert">
          {error}
        </div>
      )}
      <section className="governance-panel">
        <div className="proposal-heading">
          <div>
            <p className="eyebrow">Configuration</p>
            <h2>
              {multisig
                ? `${multisig.threshold} of ${multisig.owners.length} owners`
                : 'Loading multisig…'}
            </h2>
          </div>
          {multisig && <code>{formatAddress(multisig.contract_address)}</code>}
        </div>
        <div className="owner-list">
          {multisig?.owners.map((owner) => (
            <div key={owner}>
              <code>{formatAddress(owner)}</code>
              <span>
                {wallet.account?.toLowerCase() === owner.toLowerCase()
                  ? 'Connected'
                  : 'Owner'}
              </span>
            </div>
          ))}
        </div>
      </section>
      <div className="governance-grid">
        <form className="governance-panel governance-form" onSubmit={prepare}>
          <p className="eyebrow">New proposal</p>
          <label>
            Operation
            <select
              value={type}
              onChange={(event) => {
                setType(event.target.value)
                setParams(initialParams(event.target.value))
                setPrepared(undefined)
                setProposalStatus(undefined)
              }}
            >
              {operations.map((value) => (
                <option key={value} value={value}>
                  {operationLabels[value]}
                </option>
              ))}
            </select>
          </label>
          {Object.keys(params).map((field) => (
            <label key={field}>
              {fieldLabels[field] ?? field}
              <input
                value={params[field]}
                inputMode={
                  field.includes('owner') || field.includes('Token')
                    ? 'text'
                    : 'numeric'
                }
                placeholder={
                  field.includes('Token') || field.includes('owner')
                    ? '0x…'
                    : 'Decimal integer'
                }
                onChange={(event) =>
                  setParams((current) => ({
                    ...current,
                    [field]:
                      field === 'threshold'
                        ? Number(event.target.value)
                        : event.target.value,
                  }))
                }
              />
            </label>
          ))}
          <label>
            Unique nonce
            <input
              value={nonce}
              inputMode="numeric"
              onChange={(event) => setNonce(event.target.value)}
            />
          </label>
          <div className="proposal-review">
            <span>Review</span>
            <strong>{operationSummary(operation)}</strong>
            <small>
              All protocol amounts and rates use contract integer units. Token
              amounts use the token’s smallest unit.
            </small>
          </div>
          {validation.length > 0 && (
            <ul className="form-errors">
              {validation.map((message) => (
                <li key={message}>{message}</li>
              ))}
            </ul>
          )}
          <Button
            disabled={
              !auth.token ||
              !auth.hasScope('prism/proposals.write') ||
              validation.length > 0 ||
              !/^\d+$/.test(nonce) ||
              busy === 'prepare'
            }
          >
            {busy === 'prepare'
              ? 'Preparing…'
              : 'Prepare unsigned transactions'}
          </Button>
          {!auth.token && (
            <small>
              Operator sign-in is required. Wallet connection is not
              authentication.
            </small>
          )}
          {!/^\d+$/.test(nonce) && (
            <small className="form-error">
              Nonce must be a non-negative decimal integer.
            </small>
          )}
        </form>
        <section className="governance-panel">
          <p className="eyebrow">Proposal lookup</p>
          <form
            className="lookup-form"
            onSubmit={(event) => {
              event.preventDefault()
              setError(undefined)
              void loadStatus(lookupHash).catch((cause: unknown) =>
                setError(
                  cause instanceof Error ? cause.message : 'Lookup failed.',
                ),
              )
            }}
          >
            <input
              aria-label="Proposal transaction hash"
              placeholder="0x transaction hash"
              value={lookupHash}
              onChange={(event) => setLookupHash(event.target.value)}
            />
            <Button
              variant="secondary"
              disabled={!/^0x[0-9a-fA-F]{64}$/.test(lookupHash)}
            >
              Look up
            </Button>
          </form>
          {prepared ? (
            <div className="prepared-proposal">
              <h2>{operationSummary(operation)}</h2>
              <dl>
                <div>
                  <dt>Proposal hash</dt>
                  <dd>
                    <code>
                      {formatAddress(prepared.proposal.transactionHash)}
                    </code>
                  </dd>
                </div>
                <div>
                  <dt>Target</dt>
                  <dd>
                    <code>{formatAddress(prepared.proposal.target)}</code>
                  </dd>
                </div>
                <div>
                  <dt>Nonce</dt>
                  <dd>{prepared.proposal.nonce}</dd>
                </div>
                <div>
                  <dt>Value</dt>
                  <dd>{prepared.proposal.value}</dd>
                </div>
              </dl>
              <div className="proposal-actions">
                <Button
                  variant="secondary"
                  onClick={() =>
                    void navigator.clipboard.writeText(
                      JSON.stringify(prepared, null, 2),
                    )
                  }
                >
                  Copy JSON
                </Button>
                <Button variant="secondary" onClick={exportProposal}>
                  Export JSON
                </Button>
              </div>
              <p className="transaction-note">
                Wallet simulation runs immediately before broadcast. The backend
                prepares calldata but never signs or submits it.
              </p>
              <div className="proposal-actions">
                <Button
                  disabled={
                    !ownerConnected ||
                    currentOwner?.approved ||
                    proposalStatus?.executed ||
                    Boolean(busy)
                  }
                  onClick={() =>
                    void broadcast(prepared.approvalTransaction, 'approval')
                  }
                >
                  {busy === 'approval'
                    ? 'Confirming approval…'
                    : currentOwner?.approved
                      ? 'Already approved'
                      : 'Approve with wallet'}
                </Button>
                <Button
                  disabled={
                    !ownerConnected ||
                    !proposalStatus?.readyToExecute ||
                    proposalStatus.executed ||
                    Boolean(busy)
                  }
                  onClick={() =>
                    void broadcast(prepared.executionTransaction, 'execution')
                  }
                >
                  {busy === 'execution'
                    ? 'Confirming execution…'
                    : 'Execute proposal'}
                </Button>
              </div>
              {!wallet.account && (
                <Button variant="quiet" onClick={() => void wallet.connect()}>
                  Connect owner wallet
                </Button>
              )}
            </div>
          ) : (
            <p className="empty-copy">
              Prepare a proposal or look up its hash. Signing actions appear
              only when the unsigned transaction bundle is available in this
              browser session.
            </p>
          )}
        </section>
      </div>
      {proposalStatus && <ProposalStatusView status={proposalStatus} />}
      <p className="governance-disclosure">
        Always verify the target, operation, values, and network in your wallet.
        Authentication authorizes proposal preparation; only multisig owner
        signatures authorize on-chain changes.
      </p>
    </section>
  )
}
