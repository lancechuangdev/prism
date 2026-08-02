import { getAddress, isAddress } from 'viem'
import { z } from 'zod'

const address = z
  .string()
  .refine(isAddress, 'must be a valid EVM address')
  .transform((value) => getAddress(value))
const url = z.string().url()

const environmentSchema = z.object({
  VITE_PRISM_API_URL: url,
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
    },
  }
}

export const config = parseConfig(import.meta.env)
