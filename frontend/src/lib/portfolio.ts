import type { Address, Hash } from 'viem'

import type { PoolState } from './api/types'
import type { PoolRecord } from './pools'

export type UserPoolSide = {
  stakeAmount: bigint
  refundAmount: bigint
  hasRefunded: boolean
  hasClaimed: boolean
  positionBalance: bigint
}

export type UserPoolPosition = {
  pool: PoolRecord
  liveState: PoolState
  lender: UserPoolSide
  borrower: UserPoolSide
}

export type ActivityKind =
  | 'deposit-lend'
  | 'deposit-borrow'
  | 'refund-lend'
  | 'refund-borrow'
  | 'claim-lend'
  | 'claim-borrow'
  | 'withdraw-lend'
  | 'withdraw-borrow'

export type PortfolioActivity = {
  kind: ActivityKind
  poolIndex: number
  amount: bigint
  token?: Address
  transactionHash: Hash
  blockNumber: bigint
}

export const activityLabels: Record<ActivityKind, string> = {
  'deposit-lend': 'Lending deposit',
  'deposit-borrow': 'Collateral deposit',
  'refund-lend': 'Lending refund',
  'refund-borrow': 'Collateral refund',
  'claim-lend': 'Lender position claimed',
  'claim-borrow': 'Borrower position claimed',
  'withdraw-lend': 'Lender proceeds withdrawn',
  'withdraw-borrow': 'Borrower collateral withdrawn',
}

export function lenderRefundAvailable(position: UserPoolPosition) {
  const { pool, lender } = position
  if (
    position.liveState !== '1' ||
    lender.stakeAmount === 0n ||
    lender.hasRefunded ||
    !pool.data
  )
    return 0n
  const remaining =
    BigInt(pool.base.totalLendDeposited) - BigInt(pool.data.settleAmountLend)
  const total = BigInt(pool.base.totalLendDeposited)
  if (remaining <= 0n || total === 0n) return 0n
  return (lender.stakeAmount * remaining) / total
}

export function borrowerRefundAvailable(position: UserPoolPosition) {
  const { pool, borrower } = position
  if (
    position.liveState !== '1' ||
    borrower.stakeAmount === 0n ||
    borrower.hasRefunded ||
    !pool.data
  )
    return 0n
  const remaining =
    BigInt(pool.base.totalCollateralDeposited) -
    BigInt(pool.data.settleAmountBorrow)
  const total = BigInt(pool.base.totalCollateralDeposited)
  if (remaining <= 0n || total === 0n) return 0n
  return (borrower.stakeAmount * remaining) / total
}

export function lenderClaimAvailable(position: UserPoolPosition) {
  const { pool, lender } = position
  if (
    position.liveState !== '1' ||
    lender.stakeAmount === 0n ||
    lender.hasClaimed ||
    !pool.data
  )
    return 0n
  const total = BigInt(pool.base.totalLendDeposited)
  if (total === 0n) return 0n
  return (BigInt(pool.data.settleAmountLend) * lender.stakeAmount) / total
}

export function borrowerClaimAvailable(position: UserPoolPosition) {
  const { pool, borrower } = position
  if (
    position.liveState !== '1' ||
    borrower.stakeAmount === 0n ||
    borrower.hasClaimed ||
    !pool.data
  )
    return { position: 0n, loan: 0n }
  const total = BigInt(pool.base.totalCollateralDeposited)
  if (total === 0n) return { position: 0n, loan: 0n }
  return {
    position:
      (BigInt(pool.data.settleAmountBorrow) * borrower.stakeAmount) / total,
    loan: (BigInt(pool.data.settleAmountLend) * borrower.stakeAmount) / total,
  }
}

export function actionableCount(positions: UserPoolPosition[]) {
  return positions.reduce(
    (total, position) =>
      total +
      Number(lenderRefundAvailable(position) > 0n) +
      Number(borrowerRefundAvailable(position) > 0n) +
      Number(lenderClaimAvailable(position) > 0n) +
      Number(borrowerClaimAvailable(position).position > 0n),
    0,
  )
}
