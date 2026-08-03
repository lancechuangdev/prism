/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_PRISM_API_URL: string
  readonly VITE_PRISM_CHAIN_ID: string
  readonly VITE_PRISM_CHAIN_NAME: string
  readonly VITE_PRISM_RPC_URL: string
  readonly VITE_PRISM_EXPLORER_URL?: string
  readonly VITE_PRISM_NATIVE_CURRENCY_NAME?: string
  readonly VITE_PRISM_NATIVE_CURRENCY_SYMBOL?: string
  readonly VITE_PRISM_NATIVE_CURRENCY_DECIMALS?: string
  readonly VITE_PRISM_POOL_ADDRESS: string
  readonly VITE_PRISM_MULTISIG_ADDRESS: string
  readonly VITE_PRISM_DEPLOYMENT_BLOCK?: string
  readonly VITE_PRISM_AUTH_MODE?: 'local' | 'cognito'
  readonly VITE_PRISM_COGNITO_DOMAIN?: string
  readonly VITE_PRISM_COGNITO_CLIENT_ID?: string
  readonly VITE_PRISM_COGNITO_REDIRECT_URI?: string
  readonly VITE_PRISM_COGNITO_LOGOUT_URI?: string
  readonly VITE_PRISM_TELEMETRY_ENDPOINT?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
