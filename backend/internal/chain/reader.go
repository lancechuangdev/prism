package chain

import "context"

type Reader interface {
	PoolLength(ctx context.Context, chainID string) (int64, error)
	PoolBaseInfo(ctx context.Context, chainID string, contractIndex int64) (ContractPoolBase, error)
	PoolDataInfo(ctx context.Context, chainID string, contractIndex int64) (ContractPoolData, error)
	TokenInfo(ctx context.Context, chainID string, tokenAddress string) (ContractToken, error)
}

type ContractPoolBase struct {
	SettleTime               string
	MaturityTime             string
	InterestRate             string
	MaxLendSupply            string
	TotalLendDeposited       string
	TotalCollateralDeposited string
	CollateralizationRatio   string
	LendTokenAddress         string
	CollateralTokenAddress   string
	State                    string
	LenderPositionToken      string
	BorrowerPositionToken    string
	LiquidateRate            string
}

type ContractPoolData struct {
	SettleAmountLend        string
	SettleAmountBorrow      string
	FinishAmountLend        string
	FinishAmountBorrow      string
	LiquidationAmountLend   string
	LiquidationAmountBorrow string
}

type ContractToken struct {
	Address  string
	Symbol   string
	LogoURL  string
	Price    string
	Fee      string
	Decimals int
}
