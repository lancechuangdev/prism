import { useMemo, useState } from 'react'

import { Button } from '../components/Button'
import { DataStatus } from '../components/DataStatus'
import { PoolCard } from '../components/PoolCard'
import { usePools } from '../hooks/usePools'
import type { PoolState } from '../lib/api/types'
import {
  filterAndSortPools,
  poolStateLabels,
  type PoolSort,
} from '../lib/pools'

export function PoolMarketplacePage() {
  const { pools, loading, error, indexedAt, stale, refresh } = usePools()
  const [query, setQuery] = useState('')
  const [state, setState] = useState<PoolState | 'all'>('all')
  const [sort, setSort] = useState<PoolSort>('yield-desc')
  const visiblePools = useMemo(
    () => filterAndSortPools(pools, { query, state, sort }),
    [pools, query, sort, state],
  )

  return (
    <section className="page marketplace">
      <header className="page-heading">
        <div>
          <p className="eyebrow">Marketplace</p>
          <h1>
            Fixed terms,
            <br />
            visible risk.
          </h1>
        </div>
        <DataStatus updatedAt={indexedAt} loading={loading} stale={stale} />
      </header>
      <div className="market-controls">
        <label className="search-field">
          <span className="sr-only">Search pools</span>
          <input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search asset or pool"
          />
        </label>
        <label>
          <span className="sr-only">Pool state</span>
          <select
            value={state}
            onChange={(event) =>
              setState(event.target.value as PoolState | 'all')
            }
          >
            <option value="all">All states</option>
            {Object.entries(poolStateLabels).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span className="sr-only">Sort pools</span>
          <select
            value={sort}
            onChange={(event) => setSort(event.target.value as PoolSort)}
          >
            <option value="yield-desc">Highest yield</option>
            <option value="funding-desc">Most funded</option>
            <option value="maturity-asc">Maturity soonest</option>
            <option value="newest">Newest pool</option>
          </select>
        </label>
        <span className="result-count">
          {visiblePools.length} {visiblePools.length === 1 ? 'pool' : 'pools'}
        </span>
      </div>
      {error && (
        <div className="inline-message inline-message--error" role="alert">
          <div>
            <strong>Pool data is unavailable</strong>
            <p>
              {error} Check that the Prism backend is running and try again.
            </p>
          </div>
          <Button variant="secondary" onClick={() => void refresh()}>
            Try again
          </Button>
        </div>
      )}
      {loading && pools.length === 0 && (
        <div className="pool-grid" aria-label="Loading pools">
          {Array.from({ length: 3 }, (_, index) => (
            <div className="pool-card pool-card--skeleton" key={index} />
          ))}
        </div>
      )}
      {!loading && !error && pools.length === 0 && (
        <div className="empty-state">
          <span>00</span>
          <h2>No pools have been indexed</h2>
          <p>
            New Prism pools will appear here after the scheduler completes a
            chain sync.
          </p>
        </div>
      )}
      {!loading && pools.length > 0 && visiblePools.length === 0 && (
        <div className="empty-state">
          <span>—</span>
          <h2>No pools match these filters</h2>
          <p>Try another asset, lifecycle state, or clear the search.</p>
          <Button
            variant="secondary"
            onClick={() => {
              setQuery('')
              setState('all')
            }}
          >
            Clear filters
          </Button>
        </div>
      )}
      {visiblePools.length > 0 && (
        <div className="pool-grid">
          {visiblePools.map((pool) => (
            <PoolCard
              key={`${pool.base.key.ChainID}:${pool.index}`}
              pool={pool}
            />
          ))}
        </div>
      )}
    </section>
  )
}
