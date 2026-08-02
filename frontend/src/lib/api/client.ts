import type { AppConfig } from '../../config/env'
import type {
  DataResponse,
  IndexedPool,
  ListResponse,
  PoolBase,
  PoolData,
  PriceQuote,
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
}

export function createApiClient(config: AppConfig) {
  return new PrismApi(config.apiUrl, config.chain.id)
}
