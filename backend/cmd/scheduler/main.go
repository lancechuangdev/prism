package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/logging"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/scheduler"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)
	repo := store.NewMemoryStore()
	reader := chain.NewDemoReader()
	provider := price.NewDemoProvider()
	prices := price.NewService(provider)
	syncer := scheduler.NewPoolSyncer(reader, repo, cfg.ChainID, prices, cfg.PriceSymbol, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info(
		"scheduler starting",
		slog.String("chainID", cfg.ChainID),
		slog.Duration("interval", cfg.SyncInteral),
		slog.String("priceSymbol", cfg.PriceSymbol),
	)

	if err := syncer.Run(ctx, cfg.SyncInteral); err != nil {
		logger.Error("scheduler stopped with error", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("scheduler stopped")
}
