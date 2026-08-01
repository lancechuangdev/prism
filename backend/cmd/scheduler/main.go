package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lancechuangdev/prism/backend/internal/cache"
	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/liquidation"
	"github.com/lancechuangdev/prism/backend/internal/logging"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/scheduler"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)
	if err := cfg.Validate(config.ComponentScheduler); err != nil {
		logger.Error("invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

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

	reader, err := chain.NewRPCReader(ctx, cfg.ChainRPCURL, cfg.PoolAddress)
	if err != nil {
		logger.Error("open chain RPC reader failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer reader.Close()

	upstreamPriceProvider, err := price.NewConfiguredQuoteProvider(
		ctx,
		cfg.Env,
		cfg.PriceProvider,
		cfg.ChainRPCURL,
		cfg.OracleAddress,
		cfg.PriceTokenAddresses,
	)
	if err != nil {
		logger.Error("configure price provider failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer price.CloseQuoteProvider(upstreamPriceProvider)
	if _, err := upstreamPriceProvider.Latest(ctx, cfg.PriceSymbol); err != nil {
		logger.Error("verify configured price failed", slog.String("symbol", cfg.PriceSymbol), slog.Any("error", err))
		os.Exit(1)
	}
	priceProvider := price.NewCachedQuoteProvider(upstreamPriceProvider, cacheStore, cfg.PriceCacheTTL)
	priceService := price.NewService(priceProvider)

	var liquidationChecker scheduler.LiquidationChecker
	if cfg.LiquidationEnabled {
		liquidationChain, err := liquidation.NewRPCChain(
			ctx,
			cfg.ChainRPCURL,
			cfg.PoolAddress,
			cfg.ChainID,
			cfg.LiquidationKey,
			cfg.LiquidationSlippageBPS,
		)
		if err != nil {
			logger.Error("configure liquidation keeper failed", slog.Any("error", err))
			os.Exit(1)
		}
		defer liquidationChain.Close()
		liquidationChecker = liquidation.NewService(liquidationChain, logger)
	}
	syncer := scheduler.NewPoolSyncer(reader, repo, cfg.ChainID, priceService, cfg.PriceSymbol, liquidationChecker, logger)

	logger.Info(
		"scheduler starting",
		slog.String("chainID", cfg.ChainID),
		slog.Duration("interval", cfg.SyncInterval),
		slog.String("priceSymbol", cfg.PriceSymbol),
		slog.Bool("liquidationEnabled", cfg.LiquidationEnabled),
		slog.Uint64("liquidationSlippageBPS", cfg.LiquidationSlippageBPS),
	)

	if err := syncer.Run(ctx, cfg.SyncInterval); err != nil {
		logger.Error("scheduler stopped with error", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("scheduler stopped")
}

func openCache(ctx context.Context, cfg config.Config) (cache.Cache, func(), error) {
	redisCache, err := cache.OpenRedis(ctx, cache.RedisConfig{
		Address:       cfg.RedisAddress,
		Password:      cfg.RedisPassword,
		DB:            cfg.RedisDB,
		TLSEnabled:    cfg.RedisTLS,
		TLSServerName: cfg.RedisTLSServerName,
		RequireTLS:    strings.EqualFold(cfg.Env, "production"),
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
