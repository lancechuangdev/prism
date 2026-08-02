import { describe, expect, it } from 'vitest'

import { formatAddress, formatRate, formatTokenAmount } from './format'

describe('formatting utilities', () => {
  it('formats token amounts without floating-point arithmetic', () => {
    expect(formatTokenAmount(1_234_567n, 6)).toBe('1.2345')
  })

  it('formats Prism scaled rates', () => {
    expect(formatRate('5000000')).toBe('5%')
  })

  it('shortens valid addresses', () => {
    expect(formatAddress('0x0000000000000000000000000000000000000001')).toBe(
      '0x0000…0001',
    )
  })
})
