import type { PoolState } from '../lib/api/types'
import { poolStateLabels } from '../lib/pools'

export function PoolStateBadge({ state }: { state: PoolState }) {
  return (
    <span className={`pool-state pool-state--${state}`}>
      {poolStateLabels[state]}
    </span>
  )
}
