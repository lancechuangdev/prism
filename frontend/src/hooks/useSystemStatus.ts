import { useCallback, useEffect, useState } from 'react'
import type { Address } from 'viem'

import { config } from '../config/env'
import { createApiClient } from '../lib/api/client'
import { prismPoolAbi } from '../lib/contracts/abis'
import { publicClient } from '../lib/contracts/client'
import type { SystemSnapshot } from '../lib/systemStatus'

const api = createApiClient(config)

export function useSystemStatus() {
  const [snapshot, setSnapshot] = useState<SystemSnapshot>()
  const [refreshing, setRefreshing] = useState(true)

  const refresh = useCallback(async () => {
    setRefreshing(true)
    const [readiness, chainId, contract] = await Promise.allSettled([
      api.getReadiness(),
      publicClient.getChainId(),
      Promise.all([
        publicClient.readContract({
          address: config.contracts.pool,
          abi: prismPoolAbi,
          functionName: 'globalPaused',
        }),
        publicClient.readContract({
          address: config.contracts.pool,
          abi: prismPoolAbi,
          functionName: 'liquidator',
        }),
      ]),
    ])
    const chainMatches =
      chainId.status === 'fulfilled' && chainId.value === config.chain.id
    const contractValues =
      contract.status === 'fulfilled' ? contract.value : undefined
    setSnapshot({
      api:
        readiness.status === 'fulfilled' && readiness.value.status === 'ready'
          ? 'healthy'
          : readiness.status === 'fulfilled'
            ? 'degraded'
            : 'unavailable',
      rpc: chainMatches
        ? 'healthy'
        : chainId.status === 'fulfilled'
          ? 'degraded'
          : 'unavailable',
      contract: contractValues ? 'healthy' : 'unavailable',
      dependencies:
        readiness.status === 'fulfilled' ? readiness.value.dependencies : {},
      paused: contractValues?.[0],
      liquidator: contractValues?.[1] as Address | undefined,
      checkedAt: new Date(),
    })
    setRefreshing(false)
  }, [])

  useEffect(() => {
    queueMicrotask(() => void refresh())
    const timer = window.setInterval(() => void refresh(), 30_000)
    const onVisible = () => {
      if (document.visibilityState === 'visible') void refresh()
    }
    document.addEventListener('visibilitychange', onVisible)
    return () => {
      window.clearInterval(timer)
      document.removeEventListener('visibilitychange', onVisible)
    }
  }, [refresh])

  return { snapshot, refreshing, refresh }
}
