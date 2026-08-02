import { createPublicClient, http } from 'viem'

import { createPrismChain } from '../../config/chain'
import { config } from '../../config/env'

export const prismChain = createPrismChain(config)
export const publicClient = createPublicClient({
  chain: prismChain,
  transport: http(config.chain.rpcUrl),
})
