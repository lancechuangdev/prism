import { Button } from './Button'
import { useSystemStatus } from '../hooks/useSystemStatus'
import { formatAddress } from '../lib/format'
import { keeperLabel, overallLevel } from '../lib/systemStatus'

export function SystemStatus() {
  const { snapshot, refreshing, refresh } = useSystemStatus()
  const level = snapshot ? overallLevel(snapshot) : 'unavailable'

  return (
    <section
      className={`system-status system-status--${level}`}
      aria-labelledby="system-status-title"
      aria-live="polite"
    >
      <details>
        <summary>
          <span className="system-status__dot" aria-hidden="true" />
          <span id="system-status-title">
            {refreshing && !snapshot
              ? 'Checking system status'
              : level === 'healthy'
                ? 'Systems operational'
                : level === 'degraded'
                  ? 'Prism is operating with restrictions'
                  : 'Some Prism services are unavailable'}
          </span>
        </summary>
        {snapshot && (
          <div className="system-status__details">
            <dl>
              <div>
                <dt>API dependencies</dt>
                <dd>{snapshot.api}</dd>
              </div>
              <div>
                <dt>Chain RPC</dt>
                <dd>{snapshot.rpc}</dd>
              </div>
              <div>
                <dt>Pool contract</dt>
                <dd>{snapshot.contract}</dd>
              </div>
              <div>
                <dt>Transactions</dt>
                <dd>{snapshot.paused ? 'Paused' : 'Enabled'}</dd>
              </div>
              <div>
                <dt>Keeper</dt>
                <dd>{keeperLabel(snapshot.liquidator)}</dd>
              </div>
            </dl>
            {snapshot.liquidator && (
              <p>
                Configured liquidator:{' '}
                <code>{formatAddress(snapshot.liquidator)}</code>. This confirms
                authorization only; scheduler heartbeat is not exposed.
              </p>
            )}
            {Object.keys(snapshot.dependencies).length > 0 && (
              <p>
                Readiness:{' '}
                {Object.entries(snapshot.dependencies)
                  .map(([name, status]) => `${name} ${status}`)
                  .join(' · ')}
              </p>
            )}
            <div className="system-status__footer">
              <time dateTime={snapshot.checkedAt.toISOString()}>
                Checked {snapshot.checkedAt.toLocaleTimeString()}
              </time>
              <Button
                variant="quiet"
                disabled={refreshing}
                onClick={() => void refresh()}
              >
                {refreshing ? 'Checking…' : 'Check again'}
              </Button>
            </div>
          </div>
        )}
      </details>
    </section>
  )
}
