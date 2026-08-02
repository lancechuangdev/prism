import { useCallback, useEffect, useMemo, useState } from 'react'
import { formatUnits, getAddress, type Hash } from 'viem'

import { config } from '../config/env'
import {
  depositValidation,
  estimatedBorrowAmount,
  minimumTokenAmount,
  parseDepositAmount,
  projectedLenderInterest,
  type DepositSide,
} from '../lib/deposits'
import { erc20Abi, prismPoolAbi } from '../lib/contracts/abis'
import { publicClient } from '../lib/contracts/client'
import { formatRate, formatTokenAmount } from '../lib/format'
import type { PoolRecord } from '../lib/pools'
import { useWallet } from '../wallet/WalletProvider'
import { Button } from './Button'

type ChainState = {
  balance?: bigint
  allowance?: bigint
  minimum?: bigint
  paused?: boolean
  poolState?: number
  loading: boolean
  error?: string
}
type TransactionState = {
  stage: 'idle' | 'wallet' | 'confirming' | 'success' | 'error'
  action?: 'approval' | 'deposit'
  hash?: Hash
  message?: string
}

function errorMessage(cause: unknown) {
  if (
    typeof cause === 'object' &&
    cause &&
    'shortMessage' in cause &&
    typeof cause.shortMessage === 'string'
  )
    return cause.shortMessage
  return cause instanceof Error
    ? cause.message
    : 'The transaction could not be completed.'
}

export function DepositPanel({
  pool,
  onConfirmed,
}: {
  pool: PoolRecord
  onConfirmed: () => void
}) {
  const wallet = useWallet()
  const [side, setSide] = useState<DepositSide>('lend')
  const [value, setValue] = useState('')
  const [now, setNow] = useState(() => Date.now())
  const [chainState, setChainState] = useState<ChainState>({ loading: false })
  const [transaction, setTransaction] = useState<TransactionState>({
    stage: 'idle',
  })
  const token =
    side === 'lend' ? pool.base.lendToken : pool.base.collateralToken
  const tokenAddress = getAddress(token.address)
  const poolAddress = config.contracts.pool
  const amount = useMemo(
    () => parseDepositAmount(value, token.decimals),
    [token.decimals, value],
  )
  const remainingSupply =
    BigInt(pool.base.maxLendSupply) - BigInt(pool.base.totalLendDeposited)
  const tokenMinimum =
    chainState.minimum === undefined
      ? undefined
      : minimumTokenAmount(chainState.minimum, token.decimals)
  const deadlineReached = Math.floor(now / 1000) >= Number(pool.base.settleTime)
  const validationError = depositValidation({
    value,
    amount,
    balance: chainState.balance,
    minimum: tokenMinimum,
    remainingSupply,
    side,
    deadlineReached,
    paused: chainState.paused,
  })
  const needsApproval =
    amount !== undefined && (chainState.allowance ?? 0n) < amount

  const refreshChainState = useCallback(async () => {
    if (!wallet.account || wallet.status !== 'connected') {
      setChainState({ loading: false })
      return
    }
    setChainState((current) => ({
      ...current,
      loading: true,
      error: undefined,
    }))
    try {
      const [balance, allowance, minimum, paused, livePoolState] =
        await Promise.all([
          publicClient.readContract({
            address: tokenAddress,
            abi: erc20Abi,
            functionName: 'balanceOf',
            args: [wallet.account],
          }),
          publicClient.readContract({
            address: tokenAddress,
            abi: erc20Abi,
            functionName: 'allowance',
            args: [wallet.account, poolAddress],
          }),
          publicClient.readContract({
            address: poolAddress,
            abi: prismPoolAbi,
            functionName: side === 'lend' ? 'minLendAmount' : 'minBorrowAmount',
          }),
          publicClient.readContract({
            address: poolAddress,
            abi: prismPoolAbi,
            functionName: 'globalPaused',
          }),
          publicClient.readContract({
            address: poolAddress,
            abi: prismPoolAbi,
            functionName: 'getPoolState',
            args: [BigInt(pool.index)],
          }),
        ])
      setChainState({
        balance,
        allowance,
        minimum,
        paused,
        poolState: Number(livePoolState),
        loading: false,
      })
    } catch (cause) {
      setChainState({ loading: false, error: errorMessage(cause) })
    }
  }, [
    pool.index,
    poolAddress,
    side,
    tokenAddress,
    wallet.account,
    wallet.status,
  ])

  useEffect(() => {
    queueMicrotask(() => void refreshChainState())
  }, [refreshChainState])

  useEffect(() => {
    const timer = window.setInterval(() => setNow(Date.now()), 60_000)
    return () => window.clearInterval(timer)
  }, [])

  function selectSide(nextSide: DepositSide) {
    setSide(nextSide)
    setValue('')
    setTransaction({ stage: 'idle' })
  }

  async function submit(action: 'approval' | 'deposit') {
    if (
      !wallet.account ||
      !wallet.walletClient ||
      amount === undefined ||
      validationError
    )
      return
    setTransaction({ stage: 'wallet', action })
    try {
      let hash: Hash
      if (action === 'approval') {
        const simulation = await publicClient.simulateContract({
          account: wallet.account,
          address: tokenAddress,
          abi: erc20Abi,
          functionName: 'approve',
          args: [poolAddress, amount],
        })
        hash = await wallet.walletClient.writeContract(simulation.request)
      } else if (side === 'lend') {
        const simulation = await publicClient.simulateContract({
          account: wallet.account,
          address: poolAddress,
          abi: prismPoolAbi,
          functionName: 'depositLend',
          args: [BigInt(pool.index), amount],
        })
        hash = await wallet.walletClient.writeContract(simulation.request)
      } else {
        const simulation = await publicClient.simulateContract({
          account: wallet.account,
          address: poolAddress,
          abi: prismPoolAbi,
          functionName: 'depositBorrow',
          args: [BigInt(pool.index), amount],
        })
        hash = await wallet.walletClient.writeContract(simulation.request)
      }
      setTransaction({ stage: 'confirming', action, hash })
      const receipt = await publicClient.waitForTransactionReceipt({
        hash,
        confirmations: 1,
        onReplaced: (replacement) =>
          setTransaction({
            stage: 'confirming',
            action,
            hash: replacement.transaction.hash,
            message:
              'The wallet replaced this transaction. Tracking the replacement…',
          }),
      })
      if (receipt.status !== 'success')
        throw new Error('The transaction reverted onchain.')
      setTransaction({
        stage: 'success',
        action,
        hash: receipt.transactionHash,
        message:
          action === 'approval'
            ? `Approved exactly ${value} ${token.symbol}. You can now deposit.`
            : `Deposited ${value} ${token.symbol} into Pool ${pool.index}.`,
      })
      await refreshChainState()
      if (action === 'deposit') {
        setValue('')
        onConfirmed()
      }
    } catch (cause) {
      setTransaction({ stage: 'error', action, message: errorMessage(cause) })
    }
  }

  const projectedInterest =
    side === 'lend' && amount !== undefined
      ? projectedLenderInterest(amount, pool.base)
      : undefined
  const estimatedLoan =
    side === 'borrow' && amount !== undefined
      ? estimatedBorrowAmount(amount, pool.base)
      : undefined
  const liveStateError =
    chainState.poolState !== undefined && chainState.poolState !== 0
      ? 'Live contract state is no longer Funding.'
      : undefined
  const blockingError = validationError ?? liveStateError
  const busy =
    transaction.stage === 'wallet' || transaction.stage === 'confirming'

  return (
    <section className="deposit-panel" aria-labelledby="deposit-panel-title">
      <div
        className="deposit-tabs"
        role="tablist"
        aria-label="Pool participation side"
      >
        <button
          className={side === 'lend' ? 'active' : ''}
          role="tab"
          aria-selected={side === 'lend'}
          onClick={() => selectSide('lend')}
        >
          Lend {pool.base.lendToken.symbol}
        </button>
        <button
          className={side === 'borrow' ? 'active' : ''}
          role="tab"
          aria-selected={side === 'borrow'}
          onClick={() => selectSide('borrow')}
        >
          Supply {pool.base.collateralToken.symbol}
        </button>
      </div>
      <div className="deposit-panel__body">
        <p className="eyebrow">
          {side === 'lend'
            ? 'Lend at a fixed rate'
            : 'Borrow against collateral'}
        </p>
        <h2 id="deposit-panel-title">
          {side === 'lend'
            ? `Supply ${token.symbol}`
            : `Deposit ${token.symbol}`}
        </h2>
        <p className="deposit-description">
          {side === 'lend'
            ? `Funds are matched at settlement. The displayed ${formatRate(pool.base.interestRate)} APR is prorated for this pool’s term.`
            : `Collateral is matched at settlement at the pool’s ${formatRate(pool.base.collateralizationRatio, 100_000_000n, 0)} requirement.`}
        </p>
        {!wallet.account ? (
          <Button onClick={() => void wallet.connect()}>Connect wallet</Button>
        ) : wallet.status === 'wrong-network' ? (
          <Button onClick={() => void wallet.switchNetwork()}>
            Switch to {config.chain.name}
          </Button>
        ) : (
          <>
            <div className="amount-field">
              <div>
                <label htmlFor="deposit-amount">Amount</label>
                <button
                  type="button"
                  onClick={() =>
                    chainState.balance !== undefined &&
                    setValue(formatUnits(chainState.balance, token.decimals))
                  }
                >
                  Max
                </button>
              </div>
              <div>
                <input
                  id="deposit-amount"
                  inputMode="decimal"
                  autoComplete="off"
                  placeholder="0.00"
                  value={value}
                  onChange={(event) => {
                    setValue(event.target.value)
                    setTransaction({ stage: 'idle' })
                  }}
                  aria-describedby="amount-help"
                />
                <span>{token.symbol}</span>
              </div>
              <small id="amount-help">
                Balance{' '}
                {chainState.balance === undefined
                  ? '—'
                  : formatTokenAmount(chainState.balance, token.decimals)}{' '}
                {token.symbol}
                {tokenMinimum !== undefined && (
                  <>
                    {' '}
                    · Minimum {formatTokenAmount(tokenMinimum, token.decimals)}
                  </>
                )}
              </small>
            </div>
            <dl className="deposit-preview">
              {side === 'lend' ? (
                <>
                  <div>
                    <dt>Projected interest</dt>
                    <dd>
                      {projectedInterest === undefined
                        ? '—'
                        : `${formatTokenAmount(projectedInterest, token.decimals)} ${token.symbol}`}
                    </dd>
                  </div>
                  <div>
                    <dt>Term</dt>
                    <dd>
                      {Math.round(
                        (Number(pool.base.maturityTime) -
                          Number(pool.base.settleTime)) /
                          86_400,
                      )}{' '}
                      days
                    </dd>
                  </div>
                </>
              ) : (
                <>
                  <div>
                    <dt>Estimated loan</dt>
                    <dd>
                      {estimatedLoan === undefined
                        ? 'Oracle price required'
                        : `${formatTokenAmount(estimatedLoan, pool.base.lendToken.decimals)} ${pool.base.lendToken.symbol}`}
                    </dd>
                  </div>
                  <div>
                    <dt>Required collateral</dt>
                    <dd>
                      {formatRate(
                        pool.base.collateralizationRatio,
                        100_000_000n,
                        0,
                      )}
                    </dd>
                  </div>
                </>
              )}
            </dl>
            {(blockingError || chainState.error) && (
              <p className="form-error" role="alert">
                {blockingError ??
                  `Live chain data unavailable: ${chainState.error}`}
              </p>
            )}
            {transaction.message && (
              <div
                className={`transaction-message transaction-message--${transaction.stage}`}
                role="status"
              >
                <p>{transaction.message}</p>
                {transaction.hash && config.chain.explorerUrl && (
                  <a
                    href={`${config.chain.explorerUrl}/tx/${transaction.hash}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    View transaction ↗
                  </a>
                )}
              </div>
            )}
            {transaction.stage === 'confirming' && (
              <p className="transaction-progress" role="status">
                Waiting for onchain confirmation…
              </p>
            )}
            <Button
              className="deposit-submit"
              disabled={
                busy ||
                chainState.loading ||
                !value ||
                Boolean(blockingError) ||
                Boolean(chainState.error)
              }
              onClick={() =>
                void submit(needsApproval ? 'approval' : 'deposit')
              }
            >
              {transaction.stage === 'wallet'
                ? 'Confirm in wallet…'
                : transaction.stage === 'confirming'
                  ? 'Confirming…'
                  : needsApproval
                    ? `Approve exactly ${value || '0'} ${token.symbol}`
                    : `Deposit ${token.symbol}`}
            </Button>
            <p className="deposit-safety">
              Approval and deposit are separate transactions. Prism requests
              only the entered amount.
            </p>
          </>
        )}
      </div>
    </section>
  )
}
