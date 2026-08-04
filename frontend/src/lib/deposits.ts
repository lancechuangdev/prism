import { parseUnits } from 'viem'

import type { PoolBase, TokenSnapshot } from './api/types'
import { parseUsdPrice } from './prices'

const RATE_SCALE = 100_000_000n
const SECONDS_PER_YEAR = 365n * 24n * 60n * 60n

export type DepositSide = 'lend' | 'borrow'

export function parseDepositAmount(value: string, decimals: number) {
  if (!/^\d*(\.\d*)?$/.test(value) || value === '' || value === '.')
    return undefined
  const fraction = value.split('.')[1]
  if (fraction && fraction.length > decimals) return undefined
  try {
    return parseUnits(value, decimals)
  } catch {
    return undefined
  }
}

export function normalizeTo18Decimals(amount: bigint, decimals: number) {
  if (decimals === 18) return amount
  if (decimals < 18) return amount * 10n ** BigInt(18 - decimals)
  return amount / 10n ** BigInt(decimals - 18)
}

export function minimumTokenAmount(
  normalizedMinimum: bigint,
  tokenDecimals: number,
) {
  if (tokenDecimals === 18) return normalizedMinimum
  if (tokenDecimals < 18)
    return normalizedMinimum / 10n ** BigInt(18 - tokenDecimals)
  return normalizedMinimum * 10n ** BigInt(tokenDecimals - 18)
}

export function projectedLenderInterest(amount: bigint, pool: PoolBase) {
  const term = BigInt(pool.maturityTime) - BigInt(pool.settleTime)
  return (
    (amount * BigInt(pool.interestRate) * term) /
    (RATE_SCALE * SECONDS_PER_YEAR)
  )
}

function tokenValue(amount: bigint, token: TokenSnapshot) {
  const price = parseUsdPrice(token.price)
  if (price === undefined) return undefined
  return (amount * price) / 10n ** BigInt(token.decimals)
}

export function estimatedBorrowAmount(
  collateralAmount: bigint,
  pool: PoolBase,
) {
  const collateralValue = tokenValue(collateralAmount, pool.collateralToken)
  const lendPrice = parseUsdPrice(pool.lendToken.price)
  if (collateralValue === undefined || lendPrice === undefined) return undefined
  const lendValue =
    (collateralValue * RATE_SCALE) / BigInt(pool.collateralizationRatio)
  return (lendValue * 10n ** BigInt(pool.lendToken.decimals)) / lendPrice
}

export function depositValidation(input: {
  value: string
  amount?: bigint
  balance?: bigint
  minimum?: bigint
  remainingSupply?: bigint
  side: DepositSide
  deadlineReached: boolean
  paused?: boolean
}) {
  if (!input.value) return undefined
  if (input.amount === undefined) return 'Enter a valid token amount.'
  if (input.amount <= 0n) return 'Amount must be greater than zero.'
  if (input.deadlineReached)
    return 'This pool has reached its settlement deadline.'
  if (input.paused) return 'The Prism pool contract is paused.'
  if (input.minimum !== undefined && input.amount < input.minimum)
    return 'Amount is below the protocol minimum.'
  if (input.balance !== undefined && input.amount > input.balance)
    return 'Amount exceeds your wallet balance.'
  if (
    input.side === 'lend' &&
    input.remainingSupply !== undefined &&
    input.amount > input.remainingSupply
  )
    return 'Amount exceeds the pool’s remaining lending capacity.'
  return undefined
}
