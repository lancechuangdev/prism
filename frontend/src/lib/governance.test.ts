import { describe, expect, it } from 'vitest'

import { operationSummary, validateOperation } from './governance'

describe('governance operations', () => {
  it('validates owner addresses', () => {
    expect(
      validateOperation({ type: 'add_owner', params: { owner: 'invalid' } }),
    ).toContain('owner must be a valid address.')
  })

  it('validates pool dates without throwing on malformed input', () => {
    const params = {
      settleTime: '-',
      maturityTime: '1',
      interestRate: '1',
      maxLendSupply: '1',
      collateralizationRatio: '1',
      liquidateRate: '1',
      lendToken: '0x0000000000000000000000000000000000000001',
      collateralToken: '0x0000000000000000000000000000000000000002',
      lenderPositionToken: '0x0000000000000000000000000000000000000003',
      borrowerPositionToken: '0x0000000000000000000000000000000000000004',
    }
    expect(validateOperation({ type: 'create_pool', params })).toContain(
      'settleTime must be a non-negative decimal integer.',
    )
  })

  it('describes a threshold change in human-readable terms', () => {
    expect(
      operationSummary({
        type: 'change_threshold',
        params: { threshold: 2 },
      }),
    ).toBe('Change threshold to 2 signature(s)')
  })
})
