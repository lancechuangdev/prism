import { useCallback, useEffect, useMemo, useState } from 'react'

import { Button } from '../components/Button'
import { PositionCard } from '../components/PositionCard'
import { PortfolioNotifications } from '../components/PortfolioNotifications'
import { config } from '../config/env'
import { usePools } from '../hooks/usePools'
import { usePortfolio } from '../hooks/usePortfolio'
import { prismPoolAbi } from '../lib/contracts/abis'
import { publicClient } from '../lib/contracts/client'
import { formatAddress, formatTokenAmount } from '../lib/format'
import {
  activityLabels,
  actionableCount,
  borrowerClaimAvailable,
  borrowerRefundAvailable,
  lenderClaimAvailable,
  lenderRefundAvailable,
  lenderRedemption,
  borrowerRedemption,
  type PortfolioActivity,
} from '../lib/portfolio'
import { collateralHealth } from '../lib/pools'
import { useWallet } from '../wallet/WalletProvider'

function activityAsset(
  activity: PortfolioActivity,
  positions: ReturnType<typeof usePortfolio>['positions'],
) {
  const position = positions.find(
    (candidate) => candidate.pool.index === activity.poolIndex,
  )
  if (!position) return { symbol: 'tokens', decimals: 18 }
  if (
    activity.kind === 'deposit-borrow' ||
    activity.kind === 'refund-borrow' ||
    activity.kind === 'withdraw-borrow'
  )
    return {
      symbol: position.pool.base.collateralToken.symbol,
      decimals: position.pool.base.collateralToken.decimals,
    }
  if (activity.kind === 'claim-lend')
    return {
      symbol: 'lender position',
      decimals: position.pool.base.lendToken.decimals,
    }
  if (activity.kind === 'claim-borrow')
    return {
      symbol: 'borrower position',
      decimals: position.pool.base.collateralToken.decimals,
    }
  return {
    symbol: position.pool.base.lendToken.symbol,
    decimals: position.pool.base.lendToken.decimals,
  }
}

export function PortfolioPage() {
  const wallet = useWallet()
  const poolsState = usePools()
  const portfolio = usePortfolio(poolsState.pools, wallet.account)
  const [paused, setPaused] = useState<boolean>()
  const [now] = useState(() => Date.now())
  const [activityFilter, setActivityFilter] = useState('all')
  const [activityLimit, setActivityLimit] = useState(25)
  const refreshAll = useCallback(() => {
    void poolsState.refresh()
    void portfolio.refresh()
  }, [poolsState, portfolio])

  useEffect(() => {
    if (!wallet.account) return
    queueMicrotask(
      () =>
        void publicClient
          .readContract({
            address: config.contracts.pool,
            abi: prismPoolAbi,
            functionName: 'globalPaused',
          })
          .then(setPaused)
          .catch(() => setPaused(undefined)),
    )
  }, [wallet.account])

  const notifications = useMemo(
    () =>
      portfolio.positions
        .flatMap((position) => {
          const notices: string[] = []
          const available =
            lenderRefundAvailable(position) +
            borrowerRefundAvailable(position) +
            lenderClaimAvailable(position) +
            borrowerClaimAvailable(position).position +
            lenderRedemption(position) +
            borrowerRedemption(position)
          if (available > 0n)
            notices.push(
              `Pool ${position.pool.index} has a claim, refund, or redemption ready.`,
            )
          if (
            position.liveState === '1' &&
            Number(position.pool.base.maturityTime) * 1000 <= now
          )
            notices.push(
              `Pool ${position.pool.index} has matured and is awaiting repayment or liquidation.`,
            )
          const health = collateralHealth(position.pool)
          if (
            health !== undefined &&
            position.liveState === '1' &&
            health < 100 + Number(position.pool.base.liquidateRate) / 1_000_000
          )
            notices.push(
              `Pool ${position.pool.index} is below its liquidation threshold.`,
            )
          return notices
        })
        .concat(
          hasReusedPositionTokens(portfolio.positions)
            ? [
                'Some indexed pools reuse position-token contracts. Their displayed token balances are shared and redemption is not production-safe.',
              ]
            : [],
        ),
    [now, portfolio.positions],
  )
  const filteredActivity = useMemo(
    () =>
      portfolio.activity.filter((activity) => {
        if (activityFilter === 'all') return true
        if (activityFilter === 'lender') return activity.kind.endsWith('lend')
        if (activityFilter === 'borrower')
          return activity.kind.endsWith('borrow')
        return activity.kind === activityFilter
      }),
    [activityFilter, portfolio.activity],
  )

  if (!wallet.account)
    return (
      <section className="page portfolio-connect">
        <p className="eyebrow">Your positions</p>
        <h1>See your capital in one place.</h1>
        <p>
          Connect the wallet used to lend or supply collateral. Prism reads
          positions directly from the configured chain.
        </p>
        <Button onClick={() => void wallet.connect()}>Connect wallet</Button>
      </section>
    )
  if (wallet.status === 'wrong-network')
    return (
      <section className="page portfolio-connect">
        <p className="eyebrow">Wrong network</p>
        <h1>Switch to {config.chain.name}.</h1>
        <p>Your Prism positions are held on chain ID {config.chain.id}.</p>
        <Button onClick={() => void wallet.switchNetwork()}>
          Switch network
        </Button>
      </section>
    )

  const loading = poolsState.loading || portfolio.loading
  return (
    <section className="page portfolio-page">
      <header className="portfolio-header">
        <div>
          <p className="eyebrow">Your positions</p>
          <h1>Portfolio</h1>
          <p>{formatAddress(wallet.account)} · live chain balances</p>
        </div>
        <Button variant="secondary" disabled={loading} onClick={refreshAll}>
          {loading ? 'Refreshing…' : 'Refresh'}
        </Button>
      </header>
      {(poolsState.error || portfolio.error) && (
        <div className="inline-message inline-message--error" role="alert">
          <div>
            <strong>Portfolio data is incomplete</strong>
            <p>{poolsState.error ?? portfolio.error}</p>
          </div>
          <Button variant="secondary" onClick={refreshAll}>
            Try again
          </Button>
        </div>
      )}
      {portfolio.ignoredIndexedPools > 0 && (
        <div className="inline-message" role="status">
          <div>
            <strong>Stale indexed pools ignored</strong>
            <p>
              {portfolio.ignoredIndexedPools}{' '}
              {portfolio.ignoredIndexedPools === 1
                ? 'pool exists'
                : 'pools exist'}{' '}
              in the backend index but not in the current PrismPool contract.
              Reset the local MySQL volume to remove them permanently.
            </p>
          </div>
        </div>
      )}
      {portfolio.unindexedLivePools > 0 && (
        <div className="inline-message inline-message--error" role="alert">
          <div>
            <strong>Backend index is behind the contract</strong>
            <p>
              {portfolio.unindexedLivePools} live{' '}
              {portfolio.unindexedLivePools === 1 ? 'pool is' : 'pools are'}{' '}
              missing from indexed data. Redemption is disabled because
              position-token isolation cannot be proven across the complete pool
              set.
            </p>
          </div>
        </div>
      )}
      <div className="portfolio-summary">
        <article>
          <span>Positions</span>
          <strong>{portfolio.positions.length}</strong>
          <small>Across indexed pools</small>
        </article>
        <article>
          <span>Available actions</span>
          <strong>{actionableCount(portfolio.positions)}</strong>
          <small>Claims, refunds, and safe redemptions</small>
        </article>
        <article>
          <span>Protocol status</span>
          <strong className={paused ? 'danger-text' : ''}>
            {paused === undefined ? '—' : paused ? 'Paused' : 'Active'}
          </strong>
          <small>Live contract read</small>
        </article>
      </div>
      {notifications.length > 0 && (
        <section className="portfolio-notices" aria-labelledby="notices-title">
          <p className="eyebrow" id="notices-title">
            Needs attention
          </p>
          {notifications.map((notice) => (
            <p key={notice}>
              <span aria-hidden="true">→</span>
              {notice}
            </p>
          ))}
          <PortfolioNotifications notices={notifications} />
        </section>
      )}
      <section className="portfolio-section">
        <div className="portfolio-section__title">
          <div>
            <p className="eyebrow">Onchain positions</p>
            <h2>Supplied capital</h2>
          </div>
          <span>{portfolio.positions.length}</span>
        </div>
        {loading && portfolio.positions.length === 0 ? (
          <div className="position-grid">
            <div className="position-card position-card--skeleton" />
            <div className="position-card position-card--skeleton" />
          </div>
        ) : portfolio.positions.length === 0 ? (
          <div className="empty-state">
            <span>00</span>
            <h2>No positions found</h2>
            <p>
              This wallet has not deposited into any currently indexed Prism
              pool.
            </p>
          </div>
        ) : (
          <div className="position-grid">
            {portfolio.positions.map((position) => (
              <PositionCard
                key={position.pool.index}
                position={position}
                paused={paused}
                onConfirmed={refreshAll}
              />
            ))}
          </div>
        )}
      </section>
      <section className="portfolio-section">
        <div className="portfolio-section__title">
          <div>
            <p className="eyebrow">Contract events</p>
            <h2>Activity</h2>
          </div>
          <span>{portfolio.activity.length}</span>
        </div>
        <div className="activity-controls">
          <label>
            <span className="sr-only">Filter portfolio activity</span>
            <select
              value={activityFilter}
              onChange={(event) => {
                setActivityFilter(event.target.value)
                setActivityLimit(25)
              }}
            >
              <option value="all">All activity</option>
              <option value="lender">Lender activity</option>
              <option value="borrower">Borrower activity</option>
              <option value="withdraw-lend">Lender redemptions</option>
              <option value="withdraw-borrow">Borrower redemptions</option>
            </select>
          </label>
          <span>
            Live events · {filteredActivity.length} matching · newest first
          </span>
        </div>
        {portfolio.activity.length === 0 ? (
          <p className="activity-empty">
            No Prism activity was found from deployment block{' '}
            {config.contracts.deploymentBlock.toString()}.
          </p>
        ) : filteredActivity.length === 0 ? (
          <p className="activity-empty">No activity matches this filter.</p>
        ) : (
          <div className="activity-list">
            {filteredActivity.slice(0, activityLimit).map((activity) => {
              const asset = activityAsset(activity, portfolio.positions)
              return (
                <article key={`${activity.transactionHash}:${activity.kind}`}>
                  <span className="activity-icon" aria-hidden="true">
                    ↗
                  </span>
                  <div>
                    <strong>{activityLabels[activity.kind]}</strong>
                    <small>
                      Pool {activity.poolIndex} · Block{' '}
                      {activity.blockNumber.toString()}
                      {activity.timestamp && (
                        <>
                          {' '}
                          ·{' '}
                          <time dateTime={activity.timestamp.toISOString()}>
                            {activity.timestamp.toLocaleString()}
                          </time>
                        </>
                      )}
                      {activity.confirmations !== undefined &&
                        ` · ${activity.confirmations} confirmation${activity.confirmations === 1 ? '' : 's'}`}
                    </small>
                  </div>
                  <span>
                    {formatTokenAmount(activity.amount, asset.decimals)}{' '}
                    {asset.symbol}
                  </span>
                  {config.chain.explorerUrl ? (
                    <a
                      href={`${config.chain.explorerUrl}/tx/${activity.transactionHash}`}
                      target="_blank"
                      rel="noreferrer"
                    >
                      View ↗
                    </a>
                  ) : (
                    <code>{formatAddress(activity.transactionHash)}</code>
                  )}
                </article>
              )
            })}
            {filteredActivity.length > activityLimit && (
              <Button
                variant="secondary"
                onClick={() => setActivityLimit((current) => current + 25)}
              >
                Show 25 more
              </Button>
            )}
          </div>
        )}
      </section>
    </section>
  )
}

function hasReusedPositionTokens(
  positions: ReturnType<typeof usePortfolio>['positions'],
) {
  const addresses = positions.flatMap((position) => [
    position.pool.base.lenderPositionToken.toLowerCase(),
    position.pool.base.borrowerPositionToken.toLowerCase(),
  ])
  return new Set(addresses).size < addresses.length
}
