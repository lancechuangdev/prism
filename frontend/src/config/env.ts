import { getAddress, isAddress } from 'viem'
import { z } from 'zod'

const address = z
  .string()
  .refine(isAddress, 'must be a valid EVM address')
  .transform((value) => getAddress(value))
const url = z.string().url()
const apiUrl = z
  .string()
  .refine(
    (value) => value.startsWith('/') || URL.canParse(value),
    'must be an absolute URL or a root-relative path',
  )

const environmentSchema = z.object({
  VITE_PRISM_API_URL: apiUrl,
  VITE_PRISM_CHAIN_ID: z.coerce.number().int().positive(),
  VITE_PRISM_CHAIN_NAME: z.string().min(1),
  VITE_PRISM_RPC_URL: url,
  VITE_PRISM_EXPLORER_URL: z.union([url, z.literal('')]).optional(),
  VITE_PRISM_NATIVE_CURRENCY_NAME: z.string().min(1).default('Ether'),
  VITE_PRISM_NATIVE_CURRENCY_SYMBOL: z.string().min(1).default('ETH'),
  VITE_PRISM_NATIVE_CURRENCY_DECIMALS: z.coerce
    .number()
    .int()
    .min(0)
    .max(255)
    .default(18),
  VITE_PRISM_POOL_ADDRESS: address,
  VITE_PRISM_MULTISIG_ADDRESS: address,
  VITE_PRISM_DEPLOYMENT_BLOCK: z.string().regex(/^\d+$/).default('0'),
  VITE_PRISM_AUTH_MODE: z.enum(['local', 'cognito']).default('local'),
  VITE_PRISM_COGNITO_DOMAIN: z.string().optional(),
  VITE_PRISM_COGNITO_CLIENT_ID: z.string().optional(),
  VITE_PRISM_COGNITO_REDIRECT_URI: z.union([url, z.literal('')]).optional(),
  VITE_PRISM_COGNITO_LOGOUT_URI: z.union([url, z.literal('')]).optional(),
})

export type AppConfig = {
  apiUrl: string
  chain: {
    id: number
    name: string
    rpcUrl: string
    explorerUrl?: string
    nativeCurrency: { name: string; symbol: string; decimals: number }
  }
  contracts: {
    pool: `0x${string}`
    multisig: `0x${string}`
    deploymentBlock: bigint
  }
  auth:
    | { mode: 'local' }
    | {
        mode: 'cognito'
        domain: string
        clientId: string
        redirectUri: string
        logoutUri: string
        scopes: string[]
      }
}

export class ConfigurationError extends Error {
  constructor(public readonly issues: string[]) {
    super('Prism frontend configuration is invalid')
    this.name = 'ConfigurationError'
  }
}

export function parseConfig(environment: Record<string, unknown>): AppConfig {
  const result = environmentSchema.safeParse(environment)

  if (!result.success) {
    throw new ConfigurationError(
      result.error.issues.map(
        (issue) => `${issue.path.join('.')}: ${issue.message}`,
      ),
    )
  }

  const values = result.data
  const cognitoValues = [
    values.VITE_PRISM_COGNITO_DOMAIN,
    values.VITE_PRISM_COGNITO_CLIENT_ID,
    values.VITE_PRISM_COGNITO_REDIRECT_URI,
    values.VITE_PRISM_COGNITO_LOGOUT_URI,
  ]
  if (
    values.VITE_PRISM_AUTH_MODE === 'cognito' &&
    cognitoValues.some((value) => !value)
  ) {
    throw new ConfigurationError([
      'Cognito domain, client ID, redirect URI, and logout URI are required when VITE_PRISM_AUTH_MODE=cognito',
    ])
  }

  return {
    apiUrl: values.VITE_PRISM_API_URL.replace(/\/$/, ''),
    chain: {
      id: values.VITE_PRISM_CHAIN_ID,
      name: values.VITE_PRISM_CHAIN_NAME,
      rpcUrl: values.VITE_PRISM_RPC_URL,
      explorerUrl: values.VITE_PRISM_EXPLORER_URL || undefined,
      nativeCurrency: {
        name: values.VITE_PRISM_NATIVE_CURRENCY_NAME,
        symbol: values.VITE_PRISM_NATIVE_CURRENCY_SYMBOL,
        decimals: values.VITE_PRISM_NATIVE_CURRENCY_DECIMALS,
      },
    },
    contracts: {
      pool: values.VITE_PRISM_POOL_ADDRESS,
      multisig: values.VITE_PRISM_MULTISIG_ADDRESS,
      deploymentBlock: BigInt(values.VITE_PRISM_DEPLOYMENT_BLOCK),
    },
    auth:
      values.VITE_PRISM_AUTH_MODE === 'local'
        ? { mode: 'local' }
        : {
            mode: 'cognito',
            domain: values.VITE_PRISM_COGNITO_DOMAIN!,
            clientId: values.VITE_PRISM_COGNITO_CLIENT_ID!,
            redirectUri: values.VITE_PRISM_COGNITO_REDIRECT_URI!,
            logoutUri: values.VITE_PRISM_COGNITO_LOGOUT_URI!,
            scopes: [
              'openid',
              'profile',
              'prism/proposals.write',
              'prism/admin.read',
            ],
          },
  }
}

export const config = parseConfig(import.meta.env)
