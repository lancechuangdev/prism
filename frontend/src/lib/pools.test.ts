import { describe, expect, it } from 'vitest'

import type { IndexedPool, PoolBase, PoolData } from './api/types'
import {
  collateralHealth,
  filterAndSortPools,
  fundingPercentage,
  mergePools,
} from './pools'

function base(
  index: number,
  overrides: Partial<PoolBase> = {},
): IndexedPool<PoolBase> {
  return {
    index,
    pool_data: {
      key: { ChainID: '31337', PoolID: index + 1 },
      settleTime: '2000000000',
      maturityTime: String(2000600000 + index),
      interestRate: String(1000000 + index),
      maxLendSupply: '1000000000000000000000',
      totalLendDeposited: '500000000000000000000',
      totalCollateralDeposited: '2000000000000000000',
      collateralizationRatio: '200000000',
      lendToken: {
        address: '0x1',
        symbol: index ? 'USDC' : 'PRM',
        logoUrl: '',
        price: '0.99978414',
        fee: '0',
        decimals: 18,
      },
      collateralToken: {
        address: '0x2',
        symbol: 'WETH',
        logoUrl: '',
        price: '2000.50',
        fee: '0',
        decimals: 18,
      },
      state: '0',
      lenderPositionToken: '0x3',
      borrowerPositionToken: '0x4',
      liquidateRate: '20000000',
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
      ...overrides,
    },
  }
}

function poolData(poolId: number): IndexedPool<PoolData> {
  return {
    index: poolId - 1,
    pool_data: {
      key: { ChainID: '31337', PoolID: poolId },
      settleAmountLend: '500000000000000000000',
      settleAmountBorrow: '1000000000000000000',
      finishAmountLend: '0',
      finishAmountBorrow: '0',
      liquidationAmountLend: '0',
      liquidationAmountBorrow: '0',
      createdAt: '2026-01-01T00:00:00Z',
      updatedAt: '2026-01-01T00:00:00Z',
    },
  }
}

describe('pool domain utilities', () => {
  it('joins settlement data by persisted pool ID instead of array order', () => {
    const pools = mergePools([base(0), base(1)], [poolData(2), poolData(1)])
    expect(pools[0].data?.key.PoolID).toBe(1)
    expect(pools[1].data?.key.PoolID).toBe(2)
  })

  it('filters by asset and lifecycle state', () => {
    const pools = mergePools([base(0), base(1, { state: '1' })], [])
    expect(
      filterAndSortPools(pools, {
        query: 'usdc',
        state: '1',
        sort: 'newest',
      }).map((pool) => pool.index),
    ).toEqual([1])
  })

  it('sorts yield with integer precision', () => {
    const pools = mergePools([base(0), base(1)], [])
    expect(
      filterAndSortPools(pools, {
        query: '',
        state: 'all',
        sort: 'yield-desc',
      })[0].index,
    ).toBe(1)
  })

  it('calculates bounded funding progress', () => {
    expect(fundingPercentage('505', '1000')).toBe(50.5)
    expect(fundingPercentage('1200', '1000')).toBe(100)
  })

  it('derives collateral health only when prices are indexed', () => {
    const [pool] = mergePools([base(0)], [poolData(1)])
    expect(collateralHealth(pool)).toBe(400.1)
    pool.base.lendToken.price = ''
    expect(collateralHealth(pool)).toBeUndefined()
  })
})
