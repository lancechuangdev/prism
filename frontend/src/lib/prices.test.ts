import { describe, expect, it } from 'vitest'

import { formatUsdPrice, parseUsdPrice } from './prices'

describe('USD price utilities', () => {
  it('parses decimal Chainlink quote strings without floating point arithmetic', () => {
    expect(parseUsdPrice('0.99978414')).toBe(999_784_140_000_000_000n)
    expect(parseUsdPrice('1859.03459086')).toBe(1_859_034_590_860_000_000_000n)
  })

  it('formats decimal quotes and rejects missing or invalid prices', () => {
    expect(formatUsdPrice('0.99978414')).toBe('0.99')
    expect(formatUsdPrice('1859.03459086')).toBe('1859.03')
    expect(formatUsdPrice('')).toBeUndefined()
    expect(formatUsdPrice('invalid')).toBeUndefined()
  })
})
