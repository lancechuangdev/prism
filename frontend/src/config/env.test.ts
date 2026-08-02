import { describe, expect, it } from 'vitest'

import { ConfigurationError, parseConfig } from './env'

const validEnvironment = {
  VITE_PRISM_API_URL: 'http://localhost:8080/',
  VITE_PRISM_CHAIN_ID: '31337',
  VITE_PRISM_CHAIN_NAME: 'Hardhat Local',
  VITE_PRISM_RPC_URL: 'http://127.0.0.1:8545',
  VITE_PRISM_EXPLORER_URL: '',
  VITE_PRISM_POOL_ADDRESS: '0x0000000000000000000000000000000000000001',
  VITE_PRISM_MULTISIG_ADDRESS: '0x0000000000000000000000000000000000000002',
}

describe('parseConfig', () => {
  it('normalizes valid frontend configuration', () => {
    const config = parseConfig(validEnvironment)
    expect(config.apiUrl).toBe('http://localhost:8080')
    expect(config.chain.id).toBe(31337)
    expect(config.chain.nativeCurrency.symbol).toBe('ETH')
  })

  it('reports invalid fields before the application starts', () => {
    expect(() =>
      parseConfig({ ...validEnvironment, VITE_PRISM_CHAIN_ID: 'invalid' }),
    ).toThrow(ConfigurationError)
  })
})
