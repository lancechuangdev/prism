package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/auth"
	"github.com/lancechuangdev/prism/backend/internal/cache"
	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/httpserver"
	"github.com/lancechuangdev/prism/backend/internal/logging"
	"github.com/lancechuangdev/prism/backend/internal/multisig"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)
	if err := cfg.Validate(config.ComponentAPI); err != nil {
		logger.Error("invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

	repo, closeStore, err := openStore(context.Background(), cfg)
	if err != nil {
		logger.Error("open store failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeStore()

	reader, err := chain.NewRPCReader(
		context.Background(),
		cfg.ChainRPCURL,
		cfg.PoolAddress,
	)
	if err != nil {
		logger.Error("open chain RPC reader failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer reader.Close()

	multisigReader, err := multisig.NewRPCReader(
		context.Background(),
		cfg.ChainRPCURL,
		cfg.MultisigAddress,
	)
	if err != nil {
		logger.Error("open multisig RPC reader failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer multisigReader.Close()
	multisigConfig, err := multisigReader.Config(context.Background())
	if err != nil {
		logger.Error("read multisig config failed", slog.Any("error", err))
		os.Exit(1)
	}
	if multisigConfig.ChainID != cfg.ChainID {
		logger.Error("multisig chain ID mismatch", slog.String("rpcChainID", multisigConfig.ChainID), slog.String("configuredChainID", cfg.ChainID))
		os.Exit(1)
	}
	poolOwner, err := reader.PoolOwner(context.Background(), cfg.ChainID)
	if err != nil {
		logger.Error("read PrismPool owner failed", slog.Any("error", err))
		os.Exit(1)
	}
	if !strings.EqualFold(poolOwner, multisigConfig.ContractAddress) {
		logger.Error(
			"PrismPool owner is not the configured multisig",
			slog.String("poolOwner", poolOwner),
			slog.String("multisigAddress", multisigConfig.ContractAddress),
		)
		os.Exit(1)
	}

	if err := chain.SyncPools(context.Background(), reader, repo, cfg.ChainID); err != nil {
		logger.Error("sync contract data failed", slog.Any("error", err))
		os.Exit(1)
	}

	poolTransactions, err := chain.NewPoolTransactionBuilder(cfg.ChainID, cfg.PoolAddress)
	if err != nil {
		logger.Error("configure pool transaction builder failed", slog.Any("error", err))
		os.Exit(1)
	}

	multisigTransactions, err := multisig.NewTransactionBuilder(cfg.ChainID)
	if err != nil {
		logger.Error("configure multisig transaction builder failed", slog.Any("error", err))
		os.Exit(1)
	}

	cacheStore, closeCache, err := openCache(context.Background(), cfg)
	if err != nil {
		logger.Error("open cache failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeCache()

	authService := auth.NewService(auth.Config{
		AdminUsername: cfg.AdminUsername,
		AdminPassword: cfg.AdminPassword,
		TokenSecret:   cfg.TokenSecret,
		TokenTTL:      cfg.TokenTTL,
	}, auth.NewCacheSessionStore(cacheStore))

	chainQueryService := chain.NewQueryService(repo)
	upstreamPriceProvider, err := price.NewConfiguredQuoteProvider(
		cfg.Env, cfg.PriceProvider, cfg.PriceProviderURL, cfg.PriceProviderToken,
	)
	if err != nil {
		logger.Error("configure price provider failed", slog.Any("error", err))
		os.Exit(1)
	}
	priceProvider := price.NewCachedQuoteProvider(upstreamPriceProvider, cacheStore, cfg.PriceCacheTTL)
	priceService := price.NewService(priceProvider)
	server := httpserver.New(cfg, logger, chainQueryService, poolTransactions, multisigTransactions, multisigReader, authService, priceService)

	go func() {
		logger.Info("api server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	waitForShutdown(server, logger)
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

func waitForShutdown(server *http.Server, logger *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("api server shutting down")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("api server shutdown failed", slog.Any("error", err))
	}
}
