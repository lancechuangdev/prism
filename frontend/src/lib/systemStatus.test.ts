import { describe, expect, it } from 'vitest'

import { keeperLabel, overallLevel, type SystemSnapshot } from './systemStatus'

const healthy: SystemSnapshot = {
  api: 'healthy',
  rpc: 'healthy',
  contract: 'healthy',
  dependencies: {},
  checkedAt: new Date(0),
}

describe('system status', () => {
  it('degrades the product when the contract is paused', () => {
    expect(overallLevel({ ...healthy, paused: true })).toBe('degraded')
  })

  it('reports any unavailable dependency as unavailable', () => {
    expect(overallLevel({ ...healthy, api: 'unavailable' })).toBe('unavailable')
  })

  it('does not claim that a configured keeper is healthy', () => {
    expect(keeperLabel('0x0000000000000000000000000000000000000001')).toBe(
      'Address configured',
    )
    expect(keeperLabel('0x0000000000000000000000000000000000000000')).toBe(
      'Not configured',
    )
  })
})
