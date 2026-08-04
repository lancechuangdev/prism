import { formatUnits, parseUnits } from 'viem'

export const USD_PRICE_DECIMALS = 18

export function parseUsdPrice(value: string) {
  if (!value) return undefined
  try {
    const parsed = parseUnits(value, USD_PRICE_DECIMALS)
    return parsed > 0n ? parsed : undefined
  } catch {
    return undefined
  }
}

export function formatUsdPrice(value: string, maximumFractionDigits = 2) {
  const parsed = parseUsdPrice(value)
  if (parsed === undefined) return undefined
  const [integer, fraction = ''] = formatUnits(
    parsed,
    USD_PRICE_DECIMALS,
  ).split('.')
  const trimmedFraction = fraction
    .slice(0, maximumFractionDigits)
    .replace(/0+$/, '')
  return trimmedFraction ? `${integer}.${trimmedFraction}` : integer
}
