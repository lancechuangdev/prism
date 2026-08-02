import { defineChain } from 'viem'

import type { AppConfig } from './env'

export function createPrismChain(config: AppConfig) {
  return defineChain({
    id: config.chain.id,
    name: config.chain.name,
    nativeCurrency: config.chain.nativeCurrency,
    rpcUrls: {
      default: { http: [config.chain.rpcUrl] },
    },
    blockExplorers: config.chain.explorerUrl
      ? {
          default: {
            name: `${config.chain.name} Explorer`,
            url: config.chain.explorerUrl,
          },
        }
      : undefined,
  })
}
