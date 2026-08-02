import { useCallback, useEffect, useState } from 'react'

import { config } from '../config/env'
import { createApiClient } from '../lib/api/client'
import { mergePools, type PoolRecord } from '../lib/pools'

type PoolsState = {
  pools: PoolRecord[]
  loading: boolean
  error?: string
  indexedAt?: Date
  stale?: boolean
}

const api = createApiClient(config)

export function usePools() {
  const [state, setState] = useState<PoolsState>({ pools: [], loading: true })

  const load = useCallback(async (signal?: AbortSignal) => {
    setState((current) => ({ ...current, loading: true, error: undefined }))
    try {
      const [bases, data] = await Promise.all([
        api.listPoolBases(signal),
        api.listPoolData(signal),
      ])
      const pools = mergePools(bases.data, data.data)
      const newestTimestamp = pools.reduce(
        (latest, pool) =>
          Math.max(
            latest,
            Date.parse(pool.base.updatedAt),
            pool.data ? Date.parse(pool.data.updatedAt) : 0,
          ),
        0,
      )
      setState({
        pools,
        loading: false,
        indexedAt: newestTimestamp ? new Date(newestTimestamp) : undefined,
        stale: newestTimestamp
          ? Date.now() - newestTimestamp > 2 * 60_000
          : undefined,
      })
    } catch (cause) {
      if (signal?.aborted) return
      setState((current) => ({
        ...current,
        loading: false,
        error:
          cause instanceof Error
            ? cause.message
            : 'Pool data could not be loaded.',
      }))
    }
  }, [])

  useEffect(() => {
    const controller = new AbortController()
    void load(controller.signal)
    return () => controller.abort()
  }, [load])

  return { ...state, refresh: () => load() }
}
