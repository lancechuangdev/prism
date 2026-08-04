import { describe, expect, it } from 'vitest'

import type { PoolBase } from './api/types'
import {
  depositValidation,
  estimatedBorrowAmount,
  minimumTokenAmount,
  normalizeTo18Decimals,
  parseDepositAmount,
  projectedLenderInterest,
} from './deposits'

const pool: PoolBase = {
  key: { ChainID: '31337', PoolID: 1 },
  settleTime: '2000000000',
  maturityTime: String(2000000000 + 30 * 24 * 60 * 60),
  interestRate: '12000000',
  maxLendSupply: '100000000000',
  totalLendDeposited: '0',
  totalCollateralDeposited: '0',
  collateralizationRatio: '200000000',
  lendToken: {
    address: '0x1',
    symbol: 'USDC',
    logoUrl: '',
    price: '1.00',
    fee: '',
    decimals: 6,
  },
  collateralToken: {
    address: '0x2',
    symbol: 'WETH',
    logoUrl: '',
    price: '3000.00',
    fee: '',
    decimals: 18,
  },
  state: '0',
  lenderPositionToken: '0x3',
  borrowerPositionToken: '0x4',
  liquidateRate: '20000000',
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-01-01T00:00:00Z',
}

describe('deposit calculations and validation', () => {
  it('parses token amounts without floating point arithmetic', () => {
    expect(parseDepositAmount('12.345678', 6)).toBe(12_345_678n)
    expect(parseDepositAmount('12.3456789', 6)).toBeUndefined()
  })

  it('normalizes different token decimals and converts normalized minimums', () => {
    expect(normalizeTo18Decimals(1_000_000n, 6)).toBe(10n ** 18n)
    expect(minimumTokenAmount(100n * 10n ** 18n, 6)).toBe(100_000_000n)
  })

  it('prorates lender interest for the configured term', () => {
    expect(projectedLenderInterest(1_000_000_000n, pool)).toBe(9_863_013n)
  })

  it('estimates a loan from indexed prices and collateralization', () => {
    expect(estimatedBorrowAmount(1n * 10n ** 18n, pool)).toBe(1_500_000_000n)
  })

  it('blocks amounts over balance and remaining pool capacity', () => {
    expect(
      depositValidation({
        value: '11',
        amount: 11n,
        balance: 10n,
        side: 'borrow',
        deadlineReached: false,
      }),
    ).toBe('Amount exceeds your wallet balance.')
    expect(
      depositValidation({
        value: '11',
        amount: 11n,
        balance: 20n,
        remainingSupply: 10n,
        side: 'lend',
        deadlineReached: false,
      }),
    ).toBe('Amount exceeds the pool’s remaining lending capacity.')
  })

  it('blocks expired and paused pools before wallet submission', () => {
    expect(
      depositValidation({
        value: '1',
        amount: 1n,
        side: 'lend',
        deadlineReached: true,
      }),
    ).toBe('This pool has reached its settlement deadline.')
    expect(
      depositValidation({
        value: '1',
        amount: 1n,
        side: 'lend',
        deadlineReached: false,
        paused: true,
      }),
    ).toBe('The Prism pool contract is paused.')
  })
})
