import { formatTokenAmount } from '../lib/format'
import {
  borrowerClaimAvailable,
  borrowerRedemption,
  borrowerRefundAvailable,
  lenderClaimAvailable,
  lenderRedemption,
  lenderRefundAvailable,
  type UserPoolPosition,
} from '../lib/portfolio'
import { Link } from '../routing'
import { LifecycleActions } from './LifecycleActions'
import { PoolStateBadge } from './PoolStateBadge'
import { TokenPair } from './TokenPair'

function Status({
  complete,
  children,
}: {
  complete: boolean
  children: string
}) {
  return (
    <span
      className={
        complete
          ? 'position-status position-status--complete'
          : 'position-status'
      }
    >
      {complete ? 'Complete' : children}
    </span>
  )
}

export function PositionCard({
  position,
  paused,
  onConfirmed,
}: {
  position: UserPoolPosition
  paused?: boolean
  onConfirmed: () => void
}) {
  const { pool, lender, borrower } = position
  const lenderRefund = lenderRefundAvailable(position)
  const borrowerRefund = borrowerRefundAvailable(position)
  const lenderClaim = lenderClaimAvailable(position)
  const borrowerClaim = borrowerClaimAvailable(position)
  const lenderProceeds = lenderRedemption(position)
  const borrowerProceeds = borrowerRedemption(position)

  return (
    <article className="position-card">
      <header>
        <div>
          <TokenPair
            lend={pool.base.lendToken}
            collateral={pool.base.collateralToken}
          />
          <span className="position-pool">Pool {pool.index}</span>
        </div>
        <PoolStateBadge state={position.liveState} />
      </header>
      {lender.stakeAmount > 0n || lender.positionBalance > 0n ? (
        <section>
          <div className="position-heading">
            <h3>Lender position</h3>
            <Status complete={lender.hasClaimed}>Pending claim</Status>
          </div>
          <dl className="position-values">
            <div>
              <dt>Deposited principal</dt>
              <dd>
                {formatTokenAmount(
                  lender.stakeAmount,
                  pool.base.lendToken.decimals,
                )}{' '}
                {pool.base.lendToken.symbol}
              </dd>
            </div>
            <div>
              <dt>Refund available</dt>
              <dd>
                {formatTokenAmount(lenderRefund, pool.base.lendToken.decimals)}{' '}
                {pool.base.lendToken.symbol}
              </dd>
            </div>
            <div>
              <dt>Claimable position</dt>
              <dd>
                {formatTokenAmount(lenderClaim, pool.base.lendToken.decimals)}
              </dd>
            </div>
            <div>
              <dt>Position-token balance</dt>
              <dd>
                {formatTokenAmount(
                  lender.positionBalance,
                  pool.base.lendToken.decimals,
                )}
              </dd>
            </div>
            <div>
              <dt>Redeemable proceeds</dt>
              <dd>
                {formatTokenAmount(
                  lenderProceeds,
                  pool.base.lendToken.decimals,
                )}{' '}
                {pool.base.lendToken.symbol}
              </dd>
            </div>
          </dl>
        </section>
      ) : null}
      {borrower.stakeAmount > 0n || borrower.positionBalance > 0n ? (
        <section>
          <div className="position-heading">
            <h3>Borrower position</h3>
            <Status complete={borrower.hasClaimed}>Pending claim</Status>
          </div>
          <dl className="position-values">
            <div>
              <dt>Collateral deposited</dt>
              <dd>
                {formatTokenAmount(
                  borrower.stakeAmount,
                  pool.base.collateralToken.decimals,
                )}{' '}
                {pool.base.collateralToken.symbol}
              </dd>
            </div>
            <div>
              <dt>Refund available</dt>
              <dd>
                {formatTokenAmount(
                  borrowerRefund,
                  pool.base.collateralToken.decimals,
                )}{' '}
                {pool.base.collateralToken.symbol}
              </dd>
            </div>
            <div>
              <dt>Claimable loan</dt>
              <dd>
                {formatTokenAmount(
                  borrowerClaim.loan,
                  pool.base.lendToken.decimals,
                )}{' '}
                {pool.base.lendToken.symbol}
              </dd>
            </div>
            <div>
              <dt>Position-token balance</dt>
              <dd>
                {formatTokenAmount(
                  borrower.positionBalance,
                  pool.base.collateralToken.decimals,
                )}
              </dd>
            </div>
            <div>
              <dt>Redeemable collateral</dt>
              <dd>
                {formatTokenAmount(
                  borrowerProceeds,
                  pool.base.collateralToken.decimals,
                )}{' '}
                {pool.base.collateralToken.symbol}
              </dd>
            </div>
          </dl>
        </section>
      ) : null}
      <LifecycleActions
        position={position}
        paused={paused}
        onConfirmed={onConfirmed}
      />
      <footer>
        <Link to={`/pools/${pool.index}`}>View pool details →</Link>
      </footer>
    </article>
  )
}
