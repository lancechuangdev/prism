import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react'

import { config } from '../config/env'
import { createApiClient } from '../lib/api/client'

const TOKEN_KEY = 'prism.operator.token'
const VERIFIER_KEY = 'prism.oauth.verifier'
const STATE_KEY = 'prism.oauth.state'
const api = createApiClient(config)

type AuthContextValue = {
  token?: string
  username?: string
  scopes: string[]
  status: 'anonymous' | 'authenticating' | 'authenticated' | 'error'
  error?: string
  loginLocal: (username: string, password: string) => Promise<void>
  loginCognito: () => Promise<void>
  logout: () => Promise<void>
  hasScope: (scope: string) => boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

function base64Url(bytes: Uint8Array) {
  return btoa(String.fromCharCode(...bytes))
    .replaceAll('+', '-')
    .replaceAll('/', '_')
    .replaceAll('=', '')
}

function randomValue(size = 32) {
  return base64Url(crypto.getRandomValues(new Uint8Array(size)))
}

function tokenClaims(token: string) {
  try {
    const payload = token.split('.')[1]
    const decoded = JSON.parse(
      atob(payload.replaceAll('-', '+').replaceAll('_', '/')),
    ) as { username?: string; 'cognito:username'?: string; scope?: string }
    return {
      username: decoded.username ?? decoded['cognito:username'],
      scopes: decoded.scope?.split(' ').filter(Boolean) ?? [],
    }
  } catch {
    return { scopes: [] as string[] }
  }
}

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setToken] = useState(
    () => sessionStorage.getItem(TOKEN_KEY) ?? undefined,
  )
  const [localUsername, setLocalUsername] = useState<string>()
  const [status, setStatus] = useState<AuthContextValue['status']>(
    token ? 'authenticated' : 'anonymous',
  )
  const [error, setError] = useState<string>()
  const claims = token ? tokenClaims(token) : { scopes: [] as string[] }

  const rememberToken = useCallback((nextToken: string) => {
    sessionStorage.setItem(TOKEN_KEY, nextToken)
    setToken(nextToken)
    setStatus('authenticated')
  }, [])

  useEffect(() => {
    if (config.auth.mode !== 'cognito') return
    const query = new URLSearchParams(window.location.search)
    const code = query.get('code')
    const returnedState = query.get('state')
    if (!code) return

    const verifier = sessionStorage.getItem(VERIFIER_KEY)
    const expectedState = sessionStorage.getItem(STATE_KEY)
    if (!verifier || !returnedState || returnedState !== expectedState) {
      queueMicrotask(() => {
        setError(
          'The sign-in response could not be verified. Please try again.',
        )
        setStatus('error')
      })
      return
    }

    queueMicrotask(() => setStatus('authenticating'))
    const body = new URLSearchParams({
      grant_type: 'authorization_code',
      client_id: config.auth.clientId,
      code,
      redirect_uri: config.auth.redirectUri,
      code_verifier: verifier,
    })
    void fetch(`${config.auth.domain.replace(/\/$/, '')}/oauth2/token`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body,
    })
      .then(async (response) => {
        if (!response.ok)
          throw new Error('Cognito did not accept the authorization code.')
        return response.json() as Promise<{ access_token: string }>
      })
      .then(({ access_token }) => {
        rememberToken(access_token)
        sessionStorage.removeItem(VERIFIER_KEY)
        sessionStorage.removeItem(STATE_KEY)
        window.history.replaceState({}, '', window.location.pathname)
      })
      .catch((cause: unknown) => {
        setError(cause instanceof Error ? cause.message : 'Sign-in failed.')
        setStatus('error')
      })
  }, [rememberToken])

  const loginLocal = useCallback(
    async (username: string, password: string) => {
      setStatus('authenticating')
      setError(undefined)
      try {
        const result = await api.login(username, password)
        setLocalUsername(username)
        rememberToken(result.tokenId)
      } catch (cause) {
        setError(cause instanceof Error ? cause.message : 'Sign-in failed.')
        setStatus('error')
      }
    },
    [rememberToken],
  )

  const loginCognito = useCallback(async () => {
    if (config.auth.mode !== 'cognito') return
    const verifier = randomValue(64)
    const state = randomValue()
    const digest = await crypto.subtle.digest(
      'SHA-256',
      new TextEncoder().encode(verifier),
    )
    sessionStorage.setItem(VERIFIER_KEY, verifier)
    sessionStorage.setItem(STATE_KEY, state)
    const query = new URLSearchParams({
      client_id: config.auth.clientId,
      response_type: 'code',
      redirect_uri: config.auth.redirectUri,
      scope: config.auth.scopes.join(' '),
      state,
      code_challenge_method: 'S256',
      code_challenge: base64Url(new Uint8Array(digest)),
    })
    window.location.assign(
      `${config.auth.domain.replace(/\/$/, '')}/oauth2/authorize?${query}`,
    )
  }, [])

  const logout = useCallback(async () => {
    if (token && config.auth.mode === 'local')
      await api.logout(token).catch(() => undefined)
    sessionStorage.removeItem(TOKEN_KEY)
    setToken(undefined)
    setLocalUsername(undefined)
    setError(undefined)
    setStatus('anonymous')
    if (config.auth.mode === 'cognito') {
      const query = new URLSearchParams({
        client_id: config.auth.clientId,
        logout_uri: config.auth.logoutUri,
      })
      window.location.assign(
        `${config.auth.domain.replace(/\/$/, '')}/logout?${query}`,
      )
    }
  }, [token])

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      username: localUsername ?? claims.username,
      scopes: claims.scopes,
      status,
      error,
      loginLocal,
      loginCognito,
      logout,
      hasScope: (scope) =>
        config.auth.mode === 'local' || claims.scopes.includes(scope),
    }),
    [
      claims.scopes,
      claims.username,
      error,
      localUsername,
      loginCognito,
      loginLocal,
      logout,
      status,
      token,
    ],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

// The provider and hook intentionally share this module.
// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}
