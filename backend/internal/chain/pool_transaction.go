package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

var (
	ErrInvalidCreatePool    = errors.New("invalid create pool request")
	ErrInvalidSettlePool    = errors.New("invalid settle pool request")
	ErrInvalidRepayPool     = errors.New("invalid repay pool request")
	ErrInvalidLiquidatePool = errors.New("invalid liquidate pool request")
)

type CreatePoolParams struct {
	SettleTime             string
	MaturityTime           string
	InterestRate           string
	MaxLendSupply          string
	CollateralizationRatio string
	LendToken              string
	CollateralToken        string
	LenderPositionToken    string
	BorrowerPositionToken  string
	LiquidateRate          string
}

type RepayPoolParams struct {
	PoolID              string
	MaxCollateralAmount string
}

type LiquidatePoolParams struct {
	PoolID              string
	MaxCollateralAmount string
}

type PreparedTransaction struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainID string `json:"chainId"`
}

type PoolTransactionPreparer interface {
	PrepareCreatePool(ctx context.Context, params CreatePoolParams) (PreparedTransaction, error)
	PrepareSettlePool(ctx context.Context, poolID string) (PreparedTransaction, error)
	PrepareRepayPool(ctx context.Context, params RepayPoolParams) (PreparedTransaction, error)
	PrepareLiquidatePool(ctx context.Context, params LiquidatePoolParams) (PreparedTransaction, error)
}

type PoolTransactionBuilder struct {
	chainID     string
	poolAddress common.Address
}

func NewPoolTransactionBuilder(chainID string, poolAddress string) (*PoolTransactionBuilder, error) {
	parsedChainID, ok := new(big.Int).SetString(chainID, 10)
	if !ok || parsedChainID.Sign() <= 0 {
		return nil, fmt.Errorf("invalid chain ID %q", chainID)
	}
	if !common.IsHexAddress(poolAddress) || common.HexToAddress(poolAddress) == (common.Address{}) {
		return nil, fmt.Errorf("invalid PrismPool address %q", poolAddress)
	}
	return &PoolTransactionBuilder{
		chainID:     parsedChainID.String(),
		poolAddress: common.HexToAddress(poolAddress),
	}, nil
}

func (b *PoolTransactionBuilder) PrepareCreatePool(_ context.Context, params CreatePoolParams) (PreparedTransaction, error) {
	contractParams, err := validateCreatePoolParams(params)
	if err != nil {
		return PreparedTransaction{}, err
	}
	contractABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("load PrismPool ABI: %w", err)
	}
	data, err := contractABI.Pack("createPool", contractParams)
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("encode createPool: %w", err)
	}
	return PreparedTransaction{
		To:      b.poolAddress.Hex(),
		Data:    hexutil.Encode(data),
		Value:   "0x0",
		ChainID: b.chainID,
	}, nil
}

func (b *PoolTransactionBuilder) PrepareSettlePool(_ context.Context, poolID string) (PreparedTransaction, error) {
	parsedPoolID, ok := new(big.Int).SetString(poolID, 10)
	if !ok || parsedPoolID.Sign() < 0 || parsedPoolID.BitLen() > 256 {
		return PreparedTransaction{}, fmt.Errorf("%w: poolId must be a non-negative uint256 decimal integer", ErrInvalidSettlePool)
	}
	contractABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("load PrismPool ABI: %w", err)
	}
	data, err := contractABI.Pack("settle", parsedPoolID)
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("encode settle: %w", err)
	}
	return PreparedTransaction{
		To:      b.poolAddress.Hex(),
		Data:    hexutil.Encode(data),
		Value:   "0x0",
		ChainID: b.chainID,
	}, nil
}

func (b *PoolTransactionBuilder) PrepareRepayPool(_ context.Context, params RepayPoolParams) (PreparedTransaction, error) {
	parsedPoolID, ok := new(big.Int).SetString(params.PoolID, 10)
	if !ok || parsedPoolID.Sign() < 0 || parsedPoolID.BitLen() > 256 {
		return PreparedTransaction{}, fmt.Errorf("%w: poolId must be a non-negative uint256 decimal integer", ErrInvalidRepayPool)
	}
	maxCollateralAmount, ok := new(big.Int).SetString(params.MaxCollateralAmount, 10)
	if !ok || maxCollateralAmount.Sign() <= 0 || maxCollateralAmount.BitLen() > 256 {
		return PreparedTransaction{}, fmt.Errorf("%w: maxCollateralAmount must be a positive uint256 decimal integer", ErrInvalidRepayPool)
	}
	contractABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("load PrismPool ABI: %w", err)
	}
	data, err := contractABI.Pack("repayPool", parsedPoolID, maxCollateralAmount)
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("encode repayPool: %w", err)
	}
	return PreparedTransaction{
		To:      b.poolAddress.Hex(),
		Data:    hexutil.Encode(data),
		Value:   "0x0",
		ChainID: b.chainID,
	}, nil
}

func (b *PoolTransactionBuilder) PrepareLiquidatePool(_ context.Context, params LiquidatePoolParams) (PreparedTransaction, error) {
	parsedPoolID, ok := new(big.Int).SetString(params.PoolID, 10)
	if !ok || parsedPoolID.Sign() < 0 || parsedPoolID.BitLen() > 256 {
		return PreparedTransaction{}, fmt.Errorf("%w: poolId must be a non-negative uint256 decimal integer", ErrInvalidLiquidatePool)
	}
	maxCollateralAmount, ok := new(big.Int).SetString(params.MaxCollateralAmount, 10)
	if !ok || maxCollateralAmount.Sign() <= 0 || maxCollateralAmount.BitLen() > 256 {
		return PreparedTransaction{}, fmt.Errorf("%w: maxCollateralAmount must be a positive uint256 decimal integer", ErrInvalidLiquidatePool)
	}
	contractABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("load PrismPool ABI: %w", err)
	}
	data, err := contractABI.Pack("liquidate", parsedPoolID, maxCollateralAmount)
	if err != nil {
		return PreparedTransaction{}, fmt.Errorf("encode liquidate: %w", err)
	}
	return PreparedTransaction{
		To:      b.poolAddress.Hex(),
		Data:    hexutil.Encode(data),
		Value:   "0x0",
		ChainID: b.chainID,
	}, nil
}

func validateCreatePoolParams(params CreatePoolParams) (contracts.PrismPoolCreatePoolParams, error) {
	parsePositive := func(name string, value string) (*big.Int, error) {
		number, ok := new(big.Int).SetString(value, 10)
		if !ok || number.Sign() <= 0 {
			return nil, fmt.Errorf("%w: %s must be a positive decimal integer", ErrInvalidCreatePool, name)
		}
		return number, nil
	}

	settleTime, err := parsePositive("settleTime", params.SettleTime)
	if err != nil {
		return contracts.PrismPoolCreatePoolParams{}, err
	}
	maturityTime, err := parsePositive("maturityTime", params.MaturityTime)
	if err != nil {
		return contracts.PrismPoolCreatePoolParams{}, err
	}
	if maturityTime.Cmp(settleTime) <= 0 {
		return contracts.PrismPoolCreatePoolParams{}, fmt.Errorf("%w: maturityTime must be after settleTime", ErrInvalidCreatePool)
	}

	values := make([]*big.Int, 0, 4)
	for _, value := range []struct {
		name string
		raw  string
	}{
		{"interestRate", params.InterestRate},
		{"maxLendSupply", params.MaxLendSupply},
		{"collateralizationRatio", params.CollateralizationRatio},
		{"liquidateRate", params.LiquidateRate},
	} {
		number, err := parsePositive(value.name, value.raw)
		if err != nil {
			return contracts.PrismPoolCreatePoolParams{}, err
		}
		values = append(values, number)
	}

	addresses := []struct {
		name string
		raw  string
	}{
		{"lendToken", params.LendToken},
		{"collateralToken", params.CollateralToken},
		{"lenderPositionToken", params.LenderPositionToken},
		{"borrowerPositionToken", params.BorrowerPositionToken},
	}
	parsedAddresses := make([]common.Address, 0, len(addresses))
	for _, address := range addresses {
		if !common.IsHexAddress(address.raw) || common.HexToAddress(address.raw) == (common.Address{}) {
			return contracts.PrismPoolCreatePoolParams{}, fmt.Errorf("%w: %s must be a non-zero address", ErrInvalidCreatePool, address.name)
		}
		parsedAddresses = append(parsedAddresses, common.HexToAddress(address.raw))
	}

	return contracts.PrismPoolCreatePoolParams{
		SettleTime: settleTime, MaturityTime: maturityTime,
		InterestRate: values[0], MaxLendSupply: values[1],
		CollateralizationRatio: values[2], LendToken: parsedAddresses[0],
		CollateralToken: parsedAddresses[1], LenderPositionToken: parsedAddresses[2],
		BorrowerPositionToken: parsedAddresses[3], LiquidateRate: values[3],
	}, nil
}
