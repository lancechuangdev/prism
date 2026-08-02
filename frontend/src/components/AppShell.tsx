import type { PropsWithChildren } from 'react'

import { config } from '../config/env'
import { Link, usePathname } from '../routing'
import { useWallet } from '../wallet/WalletProvider'
import { WalletButton } from './WalletButton'

export function AppShell({ children }: PropsWithChildren) {
  const wallet = useWallet()
  const pathname = usePathname()

  return (
    <div className="app-shell">
      <a className="skip-link" href="#main-content">
        Skip to content
      </a>
      <header className="site-header">
        <Link className="brand" to="/" aria-label="Prism home">
          <span className="brand__mark" aria-hidden="true">
            P
          </span>
          <span>Prism</span>
        </Link>
        <nav aria-label="Primary navigation">
          <Link
            className={pathname === '/pools' ? 'active' : undefined}
            to="/pools"
          >
            Pools
          </Link>
          <Link
            className={pathname === '/portfolio' ? 'active' : undefined}
            to="/portfolio"
          >
            Portfolio
          </Link>
          <Link
            className={pathname === '/governance' ? 'active' : undefined}
            to="/governance"
          >
            Governance
          </Link>
        </nav>
        <WalletButton />
      </header>
      {(wallet.status === 'wrong-network' || wallet.error) && (
        <div className="network-banner" role="status">
          {wallet.error ??
            `Your wallet is not connected to ${config.chain.name}.`}
        </div>
      )}
      <main id="main-content">{children}</main>
      <footer>
        <span>Prism protocol</span>
        <span>{config.chain.name}</span>
      </footer>
    </div>
  )
}
