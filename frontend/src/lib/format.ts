import { formatUnits, getAddress, isAddress } from 'viem'

export function formatTokenAmount(
  value: bigint | string,
  decimals: number,
  maximumFractionDigits = 4,
) {
  const amount = formatUnits(
    typeof value === 'bigint' ? value : BigInt(value),
    decimals,
  )
  const [integer, fraction = ''] = amount.split('.')
  const trimmedFraction = fraction
    .slice(0, maximumFractionDigits)
    .replace(/0+$/, '')
  return trimmedFraction ? `${integer}.${trimmedFraction}` : integer
}

export function formatRate(
  value: bigint | string,
  scale = 100_000_000n,
  maximumFractionDigits = 2,
) {
  const basisPoints = (BigInt(value) * 10_000n) / scale
  return `${formatTokenAmount(basisPoints, 2, maximumFractionDigits)}%`
}

export function formatDateTime(unixSeconds: bigint | string, locale?: string) {
  return new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(Number(unixSeconds) * 1000))
}

export function formatAddress(value: string) {
  if (!isAddress(value)) return value
  const address = getAddress(value)
  return `${address.slice(0, 6)}…${address.slice(-4)}`
}
