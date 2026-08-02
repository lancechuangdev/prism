import { describe, expect, it } from 'vitest'

import type { PoolRecord } from './pools'
import {
  actionableCount,
  borrowerClaimAvailable,
  borrowerRefundAvailable,
  lenderClaimAvailable,
  lenderRefundAvailable,
  type UserPoolPosition,
} from './portfolio'

const pool: PoolRecord = {
  index: 0,
  base: {
    key: { ChainID: '31337', PoolID: 1 },
    settleTime: '100',
    maturityTime: '200',
    interestRate: '1000000',
    maxLendSupply: '1000',
    totalLendDeposited: '1000',
    totalCollateralDeposited: '100',
    collateralizationRatio: '200000000',
    lendToken: {
      address: '0x1',
      symbol: 'pUSD',
      logoUrl: '',
      price: '',
      fee: '',
      decimals: 18,
    },
    collateralToken: {
      address: '0x2',
      symbol: 'pETH',
      logoUrl: '',
      price: '',
      fee: '',
      decimals: 18,
    },
    state: '1',
    lenderPositionToken: '0x3',
    borrowerPositionToken: '0x4',
    liquidateRate: '20000000',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  },
  data: {
    key: { ChainID: '31337', PoolID: 1 },
    settleAmountLend: '800',
    settleAmountBorrow: '80',
    finishAmountLend: '0',
    finishAmountBorrow: '0',
    liquidationAmountLend: '0',
    liquidationAmountBorrow: '0',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
  },
}

function position(overrides: Partial<UserPoolPosition> = {}): UserPoolPosition {
  return {
    pool,
    liveState: '1',
    lender: {
      stakeAmount: 250n,
      refundAmount: 0n,
      hasRefunded: false,
      hasClaimed: false,
      positionBalance: 0n,
    },
    borrower: {
      stakeAmount: 50n,
      refundAmount: 0n,
      hasRefunded: false,
      hasClaimed: false,
      positionBalance: 0n,
    },
    ...overrides,
  }
}

describe('portfolio eligibility and amounts', () => {
  it('calculates proportional lender refunds and position claims', () => {
    expect(lenderRefundAvailable(position())).toBe(50n)
    expect(lenderClaimAvailable(position())).toBe(200n)
  })

  it('calculates proportional borrower refunds, positions, and loans', () => {
    expect(borrowerRefundAvailable(position())).toBe(10n)
    expect(borrowerClaimAvailable(position())).toEqual({
      position: 40n,
      loan: 400n,
    })
  })

  it('does not offer completed actions or actions outside the active state', () => {
    const completed = position({
      lender: { ...position().lender, hasRefunded: true, hasClaimed: true },
      borrower: { ...position().borrower, hasRefunded: true, hasClaimed: true },
    })
    expect(actionableCount([completed])).toBe(0)
    expect(actionableCount([position({ liveState: '2' })])).toBe(0)
  })

  it('counts each independently available lifecycle action', () => {
    expect(actionableCount([position()])).toBe(4)
  })
})
