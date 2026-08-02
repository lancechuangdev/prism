import { AppShell } from './components/AppShell'
import { HomePage } from './pages/HomePage'
import { NotFoundPage } from './pages/NotFoundPage'
import { PlaceholderPage } from './pages/PlaceholderPage'
import { usePathname } from './routing'

function CurrentPage() {
  const pathname = usePathname()

  switch (pathname) {
    case '/':
      return <HomePage />
    case '/pools':
      return (
        <PlaceholderPage
          eyebrow="Marketplace"
          title="Pool discovery is next"
          description="The Phase 1 marketplace will make every indexed Prism pool readable and comparable."
        />
      )
    case '/portfolio':
      return (
        <PlaceholderPage
          eyebrow="Your positions"
          title="A single view of your capital"
          description="Connect a wallet to prepare for wallet-specific balances, refunds, and claims."
          requiresWallet
        />
      )
    case '/governance':
      return (
        <PlaceholderPage
          eyebrow="Multisig"
          title="Governance without raw calldata"
          description="The operator console will turn Prism proposals into clear, reviewable actions."
          requiresWallet
        />
      )
    default:
      return <NotFoundPage />
  }
}

export function App() {
  return (
    <AppShell>
      <CurrentPage />
    </AppShell>
  )
}
