import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { App } from './App'
import { AuthProvider } from './auth/AuthProvider'
import { ErrorBoundary } from './components/ErrorBoundary'
import { WalletProvider } from './wallet/WalletProvider'
import './styles.css'

const root = document.getElementById('root')
if (!root) throw new Error('Application root element was not found')

createRoot(root).render(
  <StrictMode>
    <ErrorBoundary>
      <WalletProvider>
        <AuthProvider>
          <App />
        </AuthProvider>
      </WalletProvider>
    </ErrorBoundary>
  </StrictMode>,
)
