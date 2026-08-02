import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react'
import {
  createWalletClient,
  custom,
  type Address,
  type EIP1193Provider,
  type WalletClient,
} from 'viem'

import { createPrismChain } from '../config/chain'
import { config } from '../config/env'

type WalletStatus =
  'read-only' | 'connecting' | 'connected' | 'wrong-network' | 'error'

type WalletContextValue = {
  account?: Address
  status: WalletStatus
  error?: string
  walletClient?: WalletClient
  connect: () => Promise<void>
  disconnect: () => void
  switchNetwork: () => Promise<void>
}

type EthereumProvider = EIP1193Provider & {
  on?: (event: string, listener: (...args: unknown[]) => void) => void
  removeListener?: (
    event: string,
    listener: (...args: unknown[]) => void,
  ) => void
}

declare global {
  interface Window {
    ethereum?: EthereumProvider
  }
}

const WalletContext = createContext<WalletContextValue | null>(null)
const chain = createPrismChain(config)

function getProvider() {
  return window.ethereum
}

export function WalletProvider({ children }: PropsWithChildren) {
  const [account, setAccount] = useState<Address>()
  const [chainId, setChainId] = useState<number>()
  const [status, setStatus] = useState<WalletStatus>('read-only')
  const [error, setError] = useState<string>()

  const walletClient = useMemo(() => {
    const provider = getProvider()
    return provider
      ? createWalletClient({ account, chain, transport: custom(provider) })
      : undefined
  }, [account])

  const refreshProviderState = useCallback(async () => {
    const provider = getProvider()
    if (!provider) return

    const [accounts, providerChainId] = await Promise.all([
      provider.request({ method: 'eth_accounts' }) as Promise<Address[]>,
      provider.request({ method: 'eth_chainId' }) as Promise<string>,
    ])
    const nextAccount = accounts[0]
    const nextChainId = Number(providerChainId)
    setAccount(nextAccount)
    setChainId(nextChainId)
    setStatus(
      nextAccount
        ? nextChainId === config.chain.id
          ? 'connected'
          : 'wrong-network'
        : 'read-only',
    )
  }, [])

  useEffect(() => {
    const provider = getProvider()
    if (!provider) return

    const handleChange = () => void refreshProviderState()
    provider.on?.('accountsChanged', handleChange)
    provider.on?.('chainChanged', handleChange)
    queueMicrotask(
      () => void refreshProviderState().catch(() => setStatus('read-only')),
    )

    return () => {
      provider.removeListener?.('accountsChanged', handleChange)
      provider.removeListener?.('chainChanged', handleChange)
    }
  }, [refreshProviderState])

  const connect = useCallback(async () => {
    const provider = getProvider()
    if (!provider) {
      setError(
        'No injected wallet was found. Install a compatible wallet or continue in read-only mode.',
      )
      setStatus('error')
      return
    }

    setStatus('connecting')
    setError(undefined)
    try {
      await provider.request({ method: 'eth_requestAccounts' })
      await refreshProviderState()
    } catch (cause) {
      setError(
        cause instanceof Error
          ? cause.message
          : 'The wallet connection was not completed.',
      )
      setStatus('error')
    }
  }, [refreshProviderState])

  const switchNetwork = useCallback(async () => {
    const provider = getProvider()
    if (!provider) return

    setError(undefined)
    try {
      await provider.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: `0x${config.chain.id.toString(16)}` }],
      })
      await refreshProviderState()
    } catch (cause) {
      const switchError = cause as { code?: number; message?: string }
      if (switchError.code === 4902) {
        try {
          await provider.request({
            method: 'wallet_addEthereumChain',
            params: [
              {
                chainId: `0x${config.chain.id.toString(16)}`,
                chainName: config.chain.name,
                nativeCurrency: config.chain.nativeCurrency,
                rpcUrls: [config.chain.rpcUrl],
                blockExplorerUrls: config.chain.explorerUrl
                  ? [config.chain.explorerUrl]
                  : undefined,
              },
            ],
          })
          await refreshProviderState()
          return
        } catch (addError) {
          setError(
            addError instanceof Error
              ? addError.message
              : 'The network could not be added to the wallet.',
          )
        }
      } else {
        setError(switchError.message ?? 'The wallet did not switch networks.')
      }
      setStatus('error')
    }
  }, [refreshProviderState])

  const disconnect = useCallback(() => {
    setAccount(undefined)
    setChainId(undefined)
    setError(undefined)
    setStatus('read-only')
  }, [])

  const value = useMemo(
    () => ({
      account,
      status: account && chainId !== config.chain.id ? 'wrong-network' : status,
      error,
      walletClient,
      connect,
      disconnect,
      switchNetwork,
    }),
    [
      account,
      chainId,
      connect,
      disconnect,
      error,
      status,
      switchNetwork,
      walletClient,
    ],
  )

  return (
    <WalletContext.Provider value={value}>{children}</WalletContext.Provider>
  )
}

// The provider and its hook intentionally share this module.
// eslint-disable-next-line react-refresh/only-export-components
export function useWallet() {
  const context = useContext(WalletContext)
  if (!context) throw new Error('useWallet must be used within WalletProvider')
  return context
}
