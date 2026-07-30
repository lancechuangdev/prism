package liquidation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
)

type Chain interface {
	PoolCount(ctx context.Context) (int64, error)
	IsUndercollateralized(ctx context.Context, poolID int64) (bool, error)
	MaxCollateral(ctx context.Context, poolID int64) (*big.Int, error)
	Liquidate(ctx context.Context, poolID int64, maxCollateral *big.Int) (string, error)
}

type Service struct {
	chain  Chain
	logger *slog.Logger
}

func NewService(chain Chain, logger *slog.Logger) *Service {
	return &Service{chain: chain, logger: logger}
}

func (s *Service) CheckAndLiquidate(ctx context.Context) error {
	count, err := s.chain.PoolCount(ctx)
	if err != nil {
		return fmt.Errorf("read liquidation pool count: %w", err)
	}

	var failures []error
	for poolID := int64(0); poolID < count; poolID++ {
		undercollateralized, err := s.chain.IsUndercollateralized(ctx, poolID)
		if err != nil {
			failures = append(failures, fmt.Errorf("check pool %d: %w", poolID, err))
			continue
		}
		if !undercollateralized {
			continue
		}

		maxCollateral, err := s.chain.MaxCollateral(ctx, poolID)
		if err != nil {
			failures = append(failures, fmt.Errorf("read pool %d collateral: %w", poolID, err))
			continue
		}
		if maxCollateral == nil || maxCollateral.Sign() <= 0 {
			failures = append(failures, fmt.Errorf("pool %d has no settled collateral", poolID))
			continue
		}

		txHash, err := s.chain.Liquidate(ctx, poolID, maxCollateral)
		if err != nil {
			failures = append(failures, fmt.Errorf("liquidate pool %d: %w", poolID, err))
			continue
		}
		s.logger.Info(
			"pool liquidation confirmed",
			slog.String("event", "pool_liquidated"),
			slog.Int64("poolID", poolID),
			slog.String("transactionHash", txHash),
		)
	}
	return errors.Join(failures...)
}
