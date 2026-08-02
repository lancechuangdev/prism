import { useEffect, useState } from 'react'

function describeDuration(targetSeconds: string, now: number) {
  const difference = Number(targetSeconds) * 1000 - now
  if (difference <= 0) return 'Reached'
  const days = Math.floor(difference / 86_400_000)
  const hours = Math.floor((difference % 86_400_000) / 3_600_000)
  const minutes = Math.floor((difference % 3_600_000) / 60_000)
  if (days > 0) return `${days}d ${hours}h`
  if (hours > 0) return `${hours}h ${minutes}m`
  return `${Math.max(1, minutes)}m`
}

export function Countdown({
  target,
  label,
}: {
  target: string
  label: string
}) {
  const [now, setNow] = useState(() => Date.now())

  useEffect(() => {
    const timer = window.setInterval(() => setNow(Date.now()), 60_000)
    return () => window.clearInterval(timer)
  }, [])

  return (
    <span title={new Date(Number(target) * 1000).toLocaleString()}>
      {label} · {describeDuration(target, now)}
    </span>
  )
}
