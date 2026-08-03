import { config } from '../config/env'

type TelemetryEvent =
  'page_view' | 'wallet_connect_requested' | 'proposal_prepared'

export function track(
  event: TelemetryEvent,
  properties: Record<string, string | number | boolean> = {},
) {
  if (!config.telemetryEndpoint || navigator.doNotTrack === '1') return false
  const safeProperties = Object.fromEntries(
    Object.entries(properties).filter(
      ([key, value]) =>
        /^(page|operation|result|chain_id)$/.test(key) &&
        ['string', 'number', 'boolean'].includes(typeof value),
    ),
  )
  const body = JSON.stringify({
    event,
    properties: safeProperties,
    occurred_at: new Date().toISOString(),
  })
  if (navigator.sendBeacon)
    return navigator.sendBeacon(
      config.telemetryEndpoint,
      new Blob([body], { type: 'application/json' }),
    )
  void fetch(config.telemetryEndpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
    keepalive: true,
    credentials: 'omit',
  })
  return true
}
