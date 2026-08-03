import { useEffect, useState } from 'react'

import { Button } from './Button'

const deliveredKey = 'prism.portfolio.notifications.delivered'

function deliveredNotices() {
  try {
    return new Set<string>(
      JSON.parse(sessionStorage.getItem(deliveredKey) ?? '[]'),
    )
  } catch {
    return new Set<string>()
  }
}

export function PortfolioNotifications({ notices }: { notices: string[] }) {
  const supported = 'Notification' in window
  const [permission, setPermission] = useState<NotificationPermission>(
    supported ? Notification.permission : 'denied',
  )

  useEffect(() => {
    if (!supported || permission !== 'granted' || notices.length === 0) return
    const delivered = deliveredNotices()
    for (const notice of notices) {
      if (delivered.has(notice)) continue
      new Notification('Prism position needs attention', {
        body: notice,
        tag: `prism-${notice}`,
      })
      delivered.add(notice)
    }
    sessionStorage.setItem(deliveredKey, JSON.stringify([...delivered]))
  }, [notices, permission, supported])

  if (!supported)
    return (
      <p className="notification-status">
        Browser notifications are unavailable. In-page alerts remain active.
      </p>
    )
  if (permission === 'denied')
    return (
      <p className="notification-status">
        Browser notifications are blocked. Enable them in site settings to
        receive alerts while Prism is open.
      </p>
    )
  if (permission === 'granted')
    return (
      <p className="notification-status">
        Browser alerts are enabled while Prism is open. Background push delivery
        is not available.
      </p>
    )
  return (
    <div className="notification-opt-in">
      <p>
        Receive maturity, claim, refund, and collateral-risk alerts while Prism
        is open.
      </p>
      <Button
        variant="secondary"
        onClick={() =>
          void Notification.requestPermission().then(setPermission)
        }
      >
        Enable browser alerts
      </Button>
    </div>
  )
}
