package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lancechuangdev/prism/backend/internal/cache"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, closeStore, err := openStore(ctx, cfg)
	if err != nil {
		logger.Error("open store failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeStore()

	cacheStore, closeCache, err := openCache(ctx, cfg)
	if err != nil {
		logger.Error("open cache failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeCache()

	reader := chain.NewDemoReader()
	priceProvider := price.NewCachedProvider(price.NewDemoProvider(), cacheStore, cfg.PriceCacheTTL)
	priceService := price.NewService(priceProvider)
	syncer := scheduler.NewPoolSyncer(reader, repo, cfg.ChainID, priceService, cfg.PriceSymbol, logger)

	logger.Info(
		"scheduler starting",
		slog.String("chainID", cfg.ChainID),
		slog.Duration("interval", cfg.SyncInterval),
		slog.String("priceSymbol", cfg.PriceSymbol),
	)

	if err := syncer.Run(ctx, cfg.SyncInterval); err != nil {
		logger.Error("scheduler stopped with error", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("scheduler stopped")
}

func openCache(ctx context.Context, cfg config.Config) (cache.Cache, func(), error) {
	redisCache, err := cache.OpenRedis(ctx, cache.RedisConfig{
		Address:  cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		return nil, nil, err
	}
	return redisCache, func() { _ = redisCache.Close() }, nil
}

func openStore(ctx context.Context, cfg config.Config) (store.Repository, func(), error) {
	switch cfg.StoreDriver {
	case "memory":
		return store.NewMemoryStore(), func() {}, nil
	case "mysql":
		mysqlStore, err := store.OpenMySQL(ctx, cfg.MySQLDSN)
		if err != nil {
			return nil, nil, err
		}
		return mysqlStore, func() { _ = mysqlStore.Close() }, nil
	default:
		return nil, nil, fmt.Errorf("unsupported store driver %q", cfg.StoreDriver)
	}
}
