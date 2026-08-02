import { Button } from '../components/Button'
import { Countdown } from '../components/Countdown'
import { DataStatus } from '../components/DataStatus'
import { DepositPanel } from '../components/DepositPanel'
import { PoolStateBadge } from '../components/PoolStateBadge'
import { TokenPair } from '../components/TokenPair'
import { config } from '../config/env'
import { usePools } from '../hooks/usePools'
import {
  formatAddress,
  formatDateTime,
  formatRate,
  formatTokenAmount,
} from '../lib/format'
import {
  collateralHealth,
  fundingPercentage,
  newestPoolUpdate,
} from '../lib/pools'
import { Link } from '../routing'

function ExplorerAddress({ address }: { address: string }) {
  const content = formatAddress(address)
  if (!config.chain.explorerUrl) return <code>{content}</code>
  return (
    <a
      className="inline-link"
      href={`${config.chain.explorerUrl}/address/${address}`}
      target="_blank"
      rel="noreferrer"
    >
      {content} <span aria-hidden="true">↗</span>
    </a>
  )
}

function TokenPrice({ symbol, price }: { symbol: string; price: string }) {
  if (!price || price === '0')
    return <span className="unavailable">Not indexed</span>
  return (
    <>
      {formatTokenAmount(price, 8, 2)} USD <small>{symbol}</small>
    </>
  )
}

export function PoolDetailPage({ poolIndex }: { poolIndex: number }) {
  const { pools, loading, error, indexedAt, stale, refresh } = usePools()
  const pool = pools.find((candidate) => candidate.index === poolIndex)

  if (loading && !pool)
    return (
      <section className="page detail-loading">
        <div className="detail-hero detail-hero--skeleton" />
        <div className="detail-grid">
          <div />
          <div />
        </div>
      </section>
    )
  if (error && !pool)
    return (
      <section className="page placeholder">
        <p className="eyebrow">Pool {poolIndex}</p>
        <h1>Pool data is unavailable</h1>
        <p>{error}</p>
        <Button onClick={() => void refresh()}>Try again</Button>
      </section>
    )
  if (!pool)
    return (
      <section className="page placeholder">
        <p className="eyebrow">Pool {poolIndex}</p>
        <h1>Pool not found</h1>
        <p>
          This pool has not been indexed on {config.chain.name}, or the pool
          number is invalid.
        </p>
        <Link className="text-link" to="/pools">
          Return to marketplace →
        </Link>
      </section>
    )

  const { base, data } = pool
  const funding = fundingPercentage(base.totalLendDeposited, base.maxLendSupply)
  const health = collateralHealth(pool)
  const updatedAt = new Date(newestPoolUpdate(pool))

  return (
    <section className="page pool-detail">
      <div className="detail-back">
        <Link to="/pools">← All pools</Link>
        <DataStatus updatedAt={indexedAt} loading={loading} stale={stale} />
      </div>
      <header className="detail-hero">
        <div>
          <p className="eyebrow">Pool {pool.index}</p>
          <TokenPair lend={base.lendToken} collateral={base.collateralToken} />
          <div className="detail-state">
            <PoolStateBadge state={base.state} />
            <span>
              Indexed{' '}
              <time dateTime={updatedAt.toISOString()}>
                {updatedAt.toLocaleString()}
              </time>
            </span>
          </div>
        </div>
        <div className="detail-rate">
          <span>Fixed APR</span>
          <strong>{formatRate(base.interestRate)}</strong>
        </div>
      </header>
      {base.state === '0' && (
        <DepositPanel pool={pool} onConfirmed={() => void refresh()} />
      )}
      <div className="detail-grid">
        <div className="detail-main">
          <section className="detail-section">
            <div className="section-title">
              <div>
                <p className="eyebrow">Capital</p>
                <h2>Funding</h2>
              </div>
              <strong>{funding.toFixed(1)}%</strong>
            </div>
            <div
              className="progress progress--large"
              role="progressbar"
              aria-label="Funding progress"
              aria-valuenow={funding}
              aria-valuemin={0}
              aria-valuemax={100}
            >
              <span style={{ width: `${funding}%` }} />
            </div>
            <div className="funding-numbers">
              <div>
                <span>Deposited</span>
                <strong>
                  {formatTokenAmount(
                    base.totalLendDeposited,
                    base.lendToken.decimals,
                  )}{' '}
                  {base.lendToken.symbol}
                </strong>
              </div>
              <div>
                <span>Maximum supply</span>
                <strong>
                  {formatTokenAmount(
                    base.maxLendSupply,
                    base.lendToken.decimals,
                  )}{' '}
                  {base.lendToken.symbol}
                </strong>
              </div>
              <div>
                <span>Collateral deposited</span>
                <strong>
                  {formatTokenAmount(
                    base.totalCollateralDeposited,
                    base.collateralToken.decimals,
                  )}{' '}
                  {base.collateralToken.symbol}
                </strong>
              </div>
            </div>
          </section>
          <section className="detail-section">
            <p className="eyebrow">Lifecycle</p>
            <h2>Dates and settlement</h2>
            <div className="timeline">
              <div>
                <span>Settlement</span>
                <strong>{formatDateTime(base.settleTime)}</strong>
                <Countdown target={base.settleTime} label="Countdown" />
              </div>
              <div>
                <span>Maturity</span>
                <strong>{formatDateTime(base.maturityTime)}</strong>
                <Countdown target={base.maturityTime} label="Countdown" />
              </div>
            </div>
            {data && (
              <dl className="settlement-grid">
                <div>
                  <dt>Matched lending principal</dt>
                  <dd>
                    {formatTokenAmount(
                      data.settleAmountLend,
                      base.lendToken.decimals,
                    )}{' '}
                    {base.lendToken.symbol}
                  </dd>
                </div>
                <div>
                  <dt>Matched collateral</dt>
                  <dd>
                    {formatTokenAmount(
                      data.settleAmountBorrow,
                      base.collateralToken.decimals,
                    )}{' '}
                    {base.collateralToken.symbol}
                  </dd>
                </div>
                <div>
                  <dt>Lender repayment proceeds</dt>
                  <dd>
                    {formatTokenAmount(
                      data.finishAmountLend,
                      base.lendToken.decimals,
                    )}{' '}
                    {base.lendToken.symbol}
                  </dd>
                </div>
                <div>
                  <dt>Lender liquidation proceeds</dt>
                  <dd>
                    {formatTokenAmount(
                      data.liquidationAmountLend,
                      base.lendToken.decimals,
                    )}{' '}
                    {base.lendToken.symbol}
                  </dd>
                </div>
                <div>
                  <dt>Remaining borrower collateral</dt>
                  <dd>
                    {formatTokenAmount(
                      data.finishAmountBorrow,
                      base.collateralToken.decimals,
                    )}{' '}
                    {base.collateralToken.symbol}
                  </dd>
                </div>
                <div>
                  <dt>Remaining collateral after liquidation</dt>
                  <dd>
                    {formatTokenAmount(
                      data.liquidationAmountBorrow,
                      base.collateralToken.decimals,
                    )}{' '}
                    {base.collateralToken.symbol}
                  </dd>
                </div>
              </dl>
            )}
          </section>
          <section className="detail-section risk-disclosure">
            <p className="eyebrow">Before participating</p>
            <h2>Protocol limitations</h2>
            <ul>
              <li>
                Position-token redemption is unsafe if a position-token contract
                is reused across pools; production pools require pool-specific
                position accounting.
              </li>
              <li>Cancelled pools do not currently support refunds.</li>
              <li>
                Repayment-interest units remain under protocol review; these
                development pools must not be treated as production-ready.
              </li>
              <li>
                Prices and pool values shown here are indexed snapshots, not
                transaction simulations.
              </li>
            </ul>
          </section>
        </div>
        <aside className="detail-sidebar">
          <section>
            <p className="eyebrow">Risk parameters</p>
            <dl className="key-values">
              <div>
                <dt>Required collateral</dt>
                <dd>{formatRate(base.collateralizationRatio)}</dd>
              </div>
              <div>
                <dt>Liquidation margin</dt>
                <dd>{formatRate(base.liquidateRate)}</dd>
              </div>
              <div>
                <dt>Current collateral ratio</dt>
                <dd>
                  {health === undefined ? (
                    <span className="unavailable">Price unavailable</span>
                  ) : (
                    `${health.toFixed(1)}%`
                  )}
                </dd>
              </div>
            </dl>
          </section>
          <section>
            <p className="eyebrow">Indexed prices</p>
            <dl className="key-values">
              <div>
                <dt>{base.lendToken.symbol}</dt>
                <dd>
                  <TokenPrice
                    symbol={base.lendToken.symbol}
                    price={base.lendToken.price}
                  />
                </dd>
              </div>
              <div>
                <dt>{base.collateralToken.symbol}</dt>
                <dd>
                  <TokenPrice
                    symbol={base.collateralToken.symbol}
                    price={base.collateralToken.price}
                  />
                </dd>
              </div>
            </dl>
          </section>
          <section>
            <p className="eyebrow">Contracts</p>
            <dl className="key-values key-values--contracts">
              <div>
                <dt>Prism pool</dt>
                <dd>
                  <ExplorerAddress address={config.contracts.pool} />
                </dd>
              </div>
              <div>
                <dt>Lend token</dt>
                <dd>
                  <ExplorerAddress address={base.lendToken.address} />
                </dd>
              </div>
              <div>
                <dt>Collateral token</dt>
                <dd>
                  <ExplorerAddress address={base.collateralToken.address} />
                </dd>
              </div>
              <div>
                <dt>Lender position</dt>
                <dd>
                  <ExplorerAddress address={base.lenderPositionToken} />
                </dd>
              </div>
              <div>
                <dt>Borrower position</dt>
                <dd>
                  <ExplorerAddress address={base.borrowerPositionToken} />
                </dd>
              </div>
            </dl>
          </section>
        </aside>
      </div>
    </section>
  )
}
