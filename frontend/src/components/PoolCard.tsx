import { Countdown } from './Countdown'
import { PoolStateBadge } from './PoolStateBadge'
import { TokenPair } from './TokenPair'
import { formatRate, formatTokenAmount } from '../lib/format'
import { fundingPercentage, type PoolRecord } from '../lib/pools'
import { Link } from '../routing'

export function PoolCard({ pool }: { pool: PoolRecord }) {
  const { base } = pool
  const funding = fundingPercentage(base.totalLendDeposited, base.maxLendSupply)
  const deadline = base.state === '0' ? base.settleTime : base.maturityTime
  const deadlineLabel = base.state === '0' ? 'Settlement' : 'Maturity'

  return (
    <article className="pool-card">
      <div className="pool-card__top">
        <TokenPair lend={base.lendToken} collateral={base.collateralToken} />
        <PoolStateBadge state={base.state} />
      </div>
      <dl className="pool-card__metrics">
        <div>
          <dt>Fixed APR</dt>
          <dd className="metric-accent">{formatRate(base.interestRate)}</dd>
        </div>
        <div>
          <dt>Funded</dt>
          <dd>{funding.toFixed(1)}%</dd>
        </div>
        <div>
          <dt>Collateral</dt>
          <dd>{formatRate(base.collateralizationRatio)}</dd>
        </div>
      </dl>
      <div
        className="progress"
        role="progressbar"
        aria-label={`${base.lendToken.symbol} funding progress`}
        aria-valuenow={funding}
        aria-valuemin={0}
        aria-valuemax={100}
      >
        <span style={{ width: `${funding}%` }} />
      </div>
      <div className="pool-card__funding">
        <span>
          {formatTokenAmount(base.totalLendDeposited, base.lendToken.decimals)}{' '}
          {base.lendToken.symbol}
        </span>
        <span>
          of {formatTokenAmount(base.maxLendSupply, base.lendToken.decimals)}
        </span>
      </div>
      <div className="pool-card__footer">
        <span>Pool {pool.index}</span>
        <Countdown target={deadline} label={deadlineLabel} />
        <Link
          aria-label={`View pool ${pool.index}`}
          to={`/pools/${pool.index}`}
        >
          View pool <span aria-hidden="true">↗</span>
        </Link>
      </div>
    </article>
  )
}
