import type { IndexedPool, PoolBase, PoolData, PoolState } from './api/types'

export const poolStateLabels: Record<PoolState, string> = {
  '0': 'Funding',
  '1': 'Active',
  '2': 'Repaid',
  '3': 'Liquidated',
  '4': 'Cancelled',
}

export type PoolRecord = { index: number; base: PoolBase; data?: PoolData }
export type PoolSort = 'yield-desc' | 'maturity-asc' | 'funding-desc' | 'newest'
export type PoolFilters = {
  query: string
  state: PoolState | 'all'
  sort: PoolSort
}

export function mergePools(
  bases: IndexedPool<PoolBase>[],
  data: IndexedPool<PoolData>[],
): PoolRecord[] {
  const dataByPoolId = new Map(
    data.map((item) => [item.pool_data.key.PoolID, item.pool_data]),
  )
  return bases.map((item) => ({
    index: item.index,
    base: item.pool_data,
    data: dataByPoolId.get(item.pool_data.key.PoolID),
  }))
}

export function filterAndSortPools(pools: PoolRecord[], filters: PoolFilters) {
  const query = filters.query.trim().toLowerCase()
  return pools
    .filter(({ base, index }) => {
      const matchesState =
        filters.state === 'all' || base.state === filters.state
      const matchesQuery =
        !query ||
        base.lendToken.symbol.toLowerCase().includes(query) ||
        base.collateralToken.symbol.toLowerCase().includes(query) ||
        String(index).includes(query)
      return matchesState && matchesQuery
    })
    .sort((left, right) => {
      switch (filters.sort) {
        case 'maturity-asc':
          return compareBigInt(left.base.maturityTime, right.base.maturityTime)
        case 'funding-desc':
          return compareRatio(
            right.base.totalLendDeposited,
            right.base.maxLendSupply,
            left.base.totalLendDeposited,
            left.base.maxLendSupply,
          )
        case 'newest':
          return right.index - left.index
        default:
          return compareBigInt(right.base.interestRate, left.base.interestRate)
      }
    })
}

function compareBigInt(left: string, right: string) {
  const a = BigInt(left)
  const b = BigInt(right)
  return a === b ? 0 : a > b ? 1 : -1
}

function compareRatio(
  aNumerator: string,
  aDenominator: string,
  bNumerator: string,
  bDenominator: string,
) {
  const left = BigInt(aNumerator) * BigInt(bDenominator)
  const right = BigInt(bNumerator) * BigInt(aDenominator)
  return left === right ? 0 : left > right ? 1 : -1
}

export function fundingPercentage(deposited: string, maximum: string) {
  const max = BigInt(maximum)
  if (max === 0n) return 0
  return Math.min(100, Number((BigInt(deposited) * 1_000n) / max) / 10)
}

export function newestPoolUpdate(pool: PoolRecord) {
  const timestamps = [pool.base.updatedAt, pool.data?.updatedAt].filter(
    Boolean,
  ) as string[]
  return timestamps.sort().at(-1) ?? pool.base.updatedAt
}

export function collateralHealth(pool: PoolRecord) {
  const lendPrice = BigInt(pool.base.lendToken.price || '0')
  const collateralPrice = BigInt(pool.base.collateralToken.price || '0')
  if (lendPrice === 0n || collateralPrice === 0n) return undefined

  const lendAmount = BigInt(
    pool.data?.settleAmountLend || pool.base.totalLendDeposited,
  )
  const collateralAmount = BigInt(
    pool.data?.settleAmountBorrow || pool.base.totalCollateralDeposited,
  )
  if (lendAmount === 0n) return undefined

  const lendValue =
    (lendAmount * lendPrice * 10n ** 18n) /
    10n ** BigInt(pool.base.lendToken.decimals)
  const collateralValue =
    (collateralAmount * collateralPrice * 10n ** 18n) /
    10n ** BigInt(pool.base.collateralToken.decimals)
  return Number((collateralValue * 1_000n) / lendValue) / 10
}
