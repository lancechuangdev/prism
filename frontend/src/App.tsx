import { AppShell } from './components/AppShell'
import { HomePage } from './pages/HomePage'
import { GovernancePage } from './pages/GovernancePage'
import { NotFoundPage } from './pages/NotFoundPage'
import { PortfolioPage } from './pages/PortfolioPage'
import { PoolDetailPage } from './pages/PoolDetailPage'
import { PoolMarketplacePage } from './pages/PoolMarketplacePage'
import { usePathname } from './routing'

function CurrentPage() {
  const pathname = usePathname()

  switch (pathname) {
    case '/':
      return <HomePage />
    case '/pools':
      return <PoolMarketplacePage />
    case '/portfolio':
      return <PortfolioPage />
    case '/governance':
      return <GovernancePage />
    default: {
      const poolMatch = pathname.match(/^\/pools\/(\d+)$/)
      return poolMatch ? (
        <PoolDetailPage poolIndex={Number(poolMatch[1])} />
      ) : (
        <NotFoundPage />
      )
    }
  }
}

export function App() {
  return (
    <AppShell>
      <CurrentPage />
    </AppShell>
  )
}
