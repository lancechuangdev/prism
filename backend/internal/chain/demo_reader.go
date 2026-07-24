package chain

import (
	"context"
	"fmt"
)

type DemoReader struct {
	pools  map[string][]ContractPoolBase       // chainID -> poolIndex -> poolBase
	data   map[string][]ContractPoolData       // chainID -> dataIndex -> poolData
	tokens map[string]map[string]ContractToken // chainID -> tokenAddress -> tokenInfo
}

func NewDemoReader() *DemoReader {
	busd := ContractToken{
		Address:  "0xbusd",
		Symbol:   "BUSD",
		LogoURL:  "/storage/img/BUSD.png",
		Price:    "100000000",
		Fee:      "250000",
		Decimals: 18,
	}
	btc := ContractToken{
		Address:  "0xbtc",
		Symbol:   "BTC",
		LogoURL:  "/storage/img/BTC.png",
		Price:    "4200000000000",
		Fee:      "250000",
		Decimals: 18,
	}

	return &DemoReader{
		pools: map[string][]ContractPoolBase{
			"97": {
				{
					SettleTime:               "1767225600",
					MaturityTime:             "1767830400",
					InterestRate:             "10000000",
					MaxLendSupply:            "1000000000000000000000",
					TotalLendDeposited:       "500000000000000000000",
					TotalCollateralDeposited: "12500000000000000",
					CollateralizationRatio:   "200000000",
					LendTokenAddress:         busd.Address,
					CollateralTokenAddress:   btc.Address,
					State:                    "1",
					LenderPositionToken:      "0xsp",
					BorrowerPositionToken:    "0xjp",
					LiquidateRate:            "20000000",
				},
			},
		},
		data: map[string][]ContractPoolData{
			"97": {
				{
					SettleAmountLend:        "0",
					SettleAmountBorrow:      "0",
					FinishAmountLend:        "0",
					FinishAmountBorrow:      "0",
					LiquidationAmountLend:   "0",
					LiquidationAmountBorrow: "0",
				},
			},
		},
		tokens: map[string]map[string]ContractToken{
			"97": {
				busd.Address: busd,
				btc.Address:  btc,
			},
		},
	}
}

func (r *DemoReader) PoolLength(_ context.Context, chainID string) (int64, error) {
	return int64(len(r.pools[chainID])), nil
}

func (r *DemoReader) PoolBaseInfo(_ context.Context, chainID string, contractIndex int64) (ContractPoolBase, error) {
	pools := r.pools[chainID]
	if contractIndex < 0 || contractIndex >= int64(len(pools)) {
		return ContractPoolBase{}, fmt.Errorf("pool base index %d not found on chain %s", contractIndex, chainID)
	}
	return pools[contractIndex], nil
}

func (r *DemoReader) PoolDataInfo(_ context.Context, chainID string, contractIndex int64) (ContractPoolData, error) {
	pools := r.data[chainID]
	if contractIndex < 0 || contractIndex >= int64(len(pools)) {
		return ContractPoolData{}, fmt.Errorf("pool data index %d not found on chain %s", contractIndex, chainID)
	}
	return pools[contractIndex], nil
}

func (r *DemoReader) TokenInfo(_ context.Context, chainID string, tokenAddress string) (ContractToken, error) {
	tokens := r.tokens[chainID]
	token, ok := tokens[tokenAddress]
	if !ok {
		return ContractToken{}, fmt.Errorf("token %s not found on chain %s", tokenAddress, chainID)
	}
	return token, nil
}
