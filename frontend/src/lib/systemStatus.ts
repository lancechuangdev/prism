export type SystemLevel = 'healthy' | 'degraded' | 'unavailable'

export type SystemSnapshot = {
  api: SystemLevel
  rpc: SystemLevel
  contract: SystemLevel
  dependencies: Record<string, string>
  paused?: boolean
  liquidator?: string
  checkedAt: Date
}

export function overallLevel(snapshot: SystemSnapshot): SystemLevel {
  const levels = [snapshot.api, snapshot.rpc, snapshot.contract]
  if (levels.includes('unavailable')) return 'unavailable'
  if (levels.includes('degraded') || snapshot.paused) return 'degraded'
  return 'healthy'
}

export function keeperLabel(liquidator?: string) {
  if (!liquidator) return 'Status unavailable'
  return /^0x0{40}$/i.test(liquidator) ? 'Not configured' : 'Address configured'
}
