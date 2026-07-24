package store

import "time"

type TokenKey struct {
	ChainID string
	Address string
}

type TokenInfo struct {
	Key       TokenKey  `json:"key"`
	Symbol    string    `json:"symbol"`
	LogoURL   string    `json:"logoUrl"`
	Price     string    `json:"price"`
	Decimals  int       `json:"decimals"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TokenSnapshot struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	LogoURL  string `json:"logoUrl"`
	Price    string `json:"price"`
	Fee      string `json:"fee"`
	Decimals int    `json:"decimals"`
}

type PoolKey struct {
	ChainID string
	PoolID  int64
}

type PoolState string

const (
	PoolStateFunding    PoolState = "0"
	PoolStateActive     PoolState = "1"
	PoolStateREPAID     PoolState = "2"
	PoolStateLiquidated PoolState = "3"
	PoolStateCancelled  PoolState = "4"
)

type PoolBase struct {
	Key                      PoolKey       `json:"key"`
	SettleTime               string        `json:"settleTime"`
	MaturityTime             string        `json:"maturityTime"`
	InterestRate             string        `json:"interestRate"`
	MaxLendSupply            string        `json:"maxLendSupply"`
	TotalLendDeposited       string        `json:"totalLendDeposited"`
	TotalCollateralDeposited string        `json:"totalCollateralDeposited"`
	CollateralizationRatio   string        `json:"collateralizationRatio"`
	LendToken                TokenSnapshot `json:"lendToken"`
	CollateralToken          TokenSnapshot `json:"collateralToken"`
	State                    PoolState     `json:"state"`
	LenderPositionToken      string        `json:"lenderPositionToken"`
	BorrowerPositionToken    string        `json:"borrowerPositionToken"`
	LiquidateRate            string        `json:"liquidateRate"`
	CreatedAt                time.Time     `json:"createdAt"`
	UpdatedAt                time.Time     `json:"updatedAt"`
}

type PoolData struct {
	Key                     PoolKey   `json:"key"`
	SettleAmountLend        string    `json:"settleAmountLend"`
	SettleAmountBorrow      string    `json:"settleAmountBorrow"`
	FinishAmountLend        string    `json:"finishAmountLend"`
	FinishAmountBorrow      string    `json:"finishAmountBorrow"`
	LiquidationAmountLend   string    `json:"liquidationAmountLend"`
	LiquidationAmountBorrow string    `json:"liquidationAmountBorrow"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

type PoolSnapshot struct {
	Base PoolBase `json:"base"`
	Data PoolData `json:"data"`
}
