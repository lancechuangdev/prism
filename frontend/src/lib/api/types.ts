export type PoolState = '0' | '1' | '2' | '3' | '4'

export type TokenKey = { chainID: string; address: string }
export type PoolKey = { chainID: string; poolID: number }

export type TokenSnapshot = {
  address: string
  symbol: string
  logoUrl: string
  price: string
  fee: string
  decimals: number
}

export type TokenInfo = TokenSnapshot & {
  key: TokenKey
  createdAt: string
  updatedAt: string
}

export type PoolBase = {
  key: PoolKey
  settleTime: string
  maturityTime: string
  interestRate: string
  maxLendSupply: string
  totalLendDeposited: string
  totalCollateralDeposited: string
  collateralizationRatio: string
  lendToken: TokenSnapshot
  collateralToken: TokenSnapshot
  state: PoolState
  lenderPositionToken: string
  borrowerPositionToken: string
  liquidateRate: string
  createdAt: string
  updatedAt: string
}

export type PoolData = {
  key: PoolKey
  settleAmountLend: string
  settleAmountBorrow: string
  finishAmountLend: string
  finishAmountBorrow: string
  liquidationAmountLend: string
  liquidationAmountBorrow: string
  createdAt: string
  updatedAt: string
}

export type IndexedPool<T> = { index: number; pool_data: T }
export type ListResponse<T> = { data: T[] }
export type DataResponse<T> = { data: T }
export type PriceQuote = { symbol: string; price: string; timestamp?: string }
