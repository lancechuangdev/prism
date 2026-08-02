type DataStatusProps = { updatedAt?: Date; loading?: boolean; stale?: boolean }

export function DataStatus({ updatedAt, loading, stale }: DataStatusProps) {
  return (
    <div className="data-status" role="status">
      <span
        className={`data-status__dot ${loading ? 'data-status__dot--loading' : ''} ${stale ? 'data-status__dot--stale' : ''}`}
        aria-hidden="true"
      />
      <span>
        {loading
          ? 'Refreshing indexed data'
          : stale
            ? 'Indexed data may be stale'
            : 'Indexed API data'}
      </span>
      {updatedAt && (
        <time dateTime={updatedAt.toISOString()}>
          Updated{' '}
          {updatedAt.toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
          })}
        </time>
      )}
    </div>
  )
}
