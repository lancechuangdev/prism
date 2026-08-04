package scheduler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

const syncDeadline = 2 * time.Minute

type LiquidationChecker interface {
	CheckAndLiquidate(ctx context.Context) error
}

type PoolSyncer struct {
	reader       chain.Reader
	repo         store.Repository
	chainID      string
	priceService *price.Service
	liquidations LiquidationChecker
	logger       *slog.Logger
}

func NewPoolSyncer(reader chain.Reader, repo store.Repository, chainID string, priceService *price.Service, liquidations LiquidationChecker, logger *slog.Logger) *PoolSyncer {
	return &PoolSyncer{
		reader:       reader,
		repo:         repo,
		chainID:      chainID,
		priceService: priceService,
		liquidations: liquidations,
		logger:       logger,
	}
}

func (s *PoolSyncer) RunOnce(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, syncDeadline)
	defer cancel()

	if err := chain.SyncPools(ctx, s.reader, s.repo, s.chainID); err != nil {
		s.logger.Error(
			"chain provider sync failed",
			slog.String("event", "provider_failure"),
			slog.String("provider", "chain_rpc"),
			slog.Any("error", err),
		)
		return err
	}

	pools, err := s.repo.ListPoolBases(ctx, s.chainID)
	if err != nil {
		return fmt.Errorf("count synced pools: %w", err)
	}

	tokens, err := s.repo.ListTokens(ctx, s.chainID)
	if err != nil {
		return fmt.Errorf("count synced tokens: %w", err)
	}

	s.logger.Info(
		"pool sync completed",
		slog.String("chainID", s.chainID),
		slog.Int("pools", len(pools)),
		slog.Int("tokens", len(tokens)),
	)

	if s.priceService != nil {
		if err := s.refreshTokenPrices(ctx, tokens, pools); err != nil {
			return err
		}
	}

	if s.liquidations != nil {
		if err := s.liquidations.CheckAndLiquidate(ctx); err != nil {
			s.logger.Error(
				"liquidation check failed",
				slog.String("event", "liquidation_failure"),
				slog.Any("error", err),
			)
			return fmt.Errorf("check liquidations: %w", err)
		}
	}

	s.logger.Info(
		"scheduler sync succeeded",
		slog.String("event", "scheduler_sync_success"),
		slog.Int64("completed_at_unix", time.Now().Unix()),
	)
	return nil
}

func (s *PoolSyncer) refreshTokenPrices(ctx context.Context, tokens []store.TokenInfo, pools []store.PoolBase) error {
	pricesByAddress := make(map[string]string, len(tokens))
	var failures []error
	for _, token := range tokens {
		quote, err := s.priceService.Latest(ctx, token.Symbol)
		if err != nil {
			s.logger.Error(
				"price provider refresh failed",
				slog.String("event", "provider_failure"),
				slog.String("provider", "chainlink_oracle"),
				slog.String("symbol", token.Symbol),
				slog.Any("error", err),
			)
			failures = append(failures, fmt.Errorf("refresh price %s: %w", token.Symbol, err))
			continue
		}

		token.Price = quote.Price
		if err := s.repo.UpsertToken(ctx, token); err != nil {
			failures = append(failures, fmt.Errorf("save price for token %s: %w", token.Symbol, err))
			continue
		}
		pricesByAddress[strings.ToLower(token.Key.Address)] = quote.Price
		s.logger.Info(
			"price refresh completed",
			slog.String("symbol", quote.Symbol),
			slog.String("currency", quote.Currency),
			slog.String("price", quote.Price),
			slog.String("source", quote.Source),
		)
	}

	for _, pool := range pools {
		if value, ok := pricesByAddress[strings.ToLower(pool.LendToken.Address)]; ok {
			pool.LendToken.Price = value
		}
		if value, ok := pricesByAddress[strings.ToLower(pool.CollateralToken.Address)]; ok {
			pool.CollateralToken.Price = value
		}
		if err := s.repo.UpsertPoolBase(ctx, pool); err != nil {
			failures = append(failures, fmt.Errorf("save token prices for pool %d: %w", pool.Key.PoolID, err))
		}
	}

	return errors.Join(failures...)
}

func (s *PoolSyncer) Run(ctx context.Context, interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("sync interval must be positive")
	}

	if err := s.RunOnce(ctx); err != nil {
		s.logger.Error(
			"initial pool sync failed",
			slog.String("event", "scheduler_sync_failure"),
			slog.Any("error", err),
		)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.RunOnce(ctx); err != nil {
				s.logger.Error(
					"scheduled pool sync failed",
					slog.String("event", "scheduler_sync_failure"),
					slog.Any("error", err),
				)
			}
		}
	}
}
