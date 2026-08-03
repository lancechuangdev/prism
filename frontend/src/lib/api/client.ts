import type { AppConfig } from '../../config/env'
import type {
  DataResponse,
  IndexedPool,
  ListResponse,
  PoolBase,
  PoolData,
  PriceQuote,
  MultisigConfig,
  PreparedProposal,
  ProposalOperation,
  ProposalStatus,
  TokenInfo,
} from './types'

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export class PrismApi {
  constructor(
    private readonly baseUrl: string,
    private readonly chainId: number,
  ) {}

  private async get<T>(path: string, signal?: AbortSignal): Promise<T> {
    const response = await fetch(`${this.baseUrl}${path}`, {
      headers: { Accept: 'application/json' },
      signal,
    })

    if (!response.ok) {
      throw new ApiError(
        `Prism API request failed with status ${response.status}`,
        response.status,
      )
    }

    return response.json() as Promise<T>
  }

  private async send<T>(
    path: string,
    method: 'POST',
    body: unknown,
    token?: string,
  ) {
    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify(body),
    })
    if (!response.ok) {
      const payload = (await response.json().catch(() => undefined)) as
        { error?: string } | undefined
      throw new ApiError(
        payload?.error ??
          `Prism API request failed with status ${response.status}`,
        response.status,
      )
    }
    return response.status === 204
      ? (undefined as T)
      : (response.json() as Promise<T>)
  }

  listPoolBases(signal?: AbortSignal) {
    return this.get<ListResponse<IndexedPool<PoolBase>>>(
      `/api/v1/poolBaseInfo?chainId=${this.chainId}`,
      signal,
    )
  }

  listPoolData(signal?: AbortSignal) {
    return this.get<ListResponse<IndexedPool<PoolData>>>(
      `/api/v1/poolDataInfo?chainId=${this.chainId}`,
      signal,
    )
  }

  listTokens(signal?: AbortSignal) {
    return this.get<ListResponse<TokenInfo>>(
      `/api/v1/token?chainId=${this.chainId}`,
      signal,
    )
  }

  getPrice(symbol: string, signal?: AbortSignal) {
    return this.get<DataResponse<PriceQuote>>(
      `/api/v1/price?symbol=${encodeURIComponent(symbol)}`,
      signal,
    )
  }

  getReadiness(signal?: AbortSignal) {
    return this.get<{ status: string; dependencies: Record<string, string> }>(
      '/readyz',
      signal,
    )
  }

  login(name: string, password: string) {
    return this.send<{ tokenId: string }>('/api/v1/user/login', 'POST', {
      name,
      password,
    })
  }

  logout(token: string) {
    return this.send<Record<string, never>>(
      '/api/v1/user/logout',
      'POST',
      {},
      token,
    )
  }

  getMultisig(signal?: AbortSignal) {
    return this.get<DataResponse<MultisigConfig>>('/api/v1/multisig', signal)
  }

  getProposalStatus(hash: string, signal?: AbortSignal) {
    return this.get<DataResponse<ProposalStatus>>(
      `/api/v1/multisig/proposals/${encodeURIComponent(hash)}`,
      signal,
    )
  }

  prepareProposal(
    chainId: string,
    nonce: string,
    operation: ProposalOperation,
    token: string,
  ) {
    return this.send<DataResponse<PreparedProposal>>(
      '/api/v1/multisig/proposals',
      'POST',
      { chain_id: chainId, nonce, operation },
      token,
    )
  }
}

export function createApiClient(config: AppConfig) {
  return new PrismApi(config.apiUrl, config.chain.id)
}
