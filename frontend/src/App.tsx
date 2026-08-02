import { AppShell } from './components/AppShell'
import { HomePage } from './pages/HomePage'
import { NotFoundPage } from './pages/NotFoundPage'
import { PlaceholderPage } from './pages/PlaceholderPage'
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
      return (
        <PlaceholderPage
          eyebrow="Multisig"
          title="Governance without raw calldata"
          description="The operator console will turn Prism proposals into clear, reviewable actions."
          requiresWallet
        />
      )
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
