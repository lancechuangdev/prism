package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	repo, closeStore, err := openStore(context.Background(), cfg)
	if err != nil {
		logger.Error("open store failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeStore()

	reader := chain.NewDemoReader()

	if err := chain.SyncPools(context.Background(), reader, repo, cfg.ChainID); err != nil {
		logger.Error("sync demo contract data failed", slog.Any("error", err))
		os.Exit(1)
	}

	authService := auth.NewService(auth.Config{
		AdminUsername: cfg.AdminUsername,
		AdminPassword: cfg.AdminPassword,
		TokenSecret:   cfg.TokenSecret,
		TokenTTL:      cfg.TokenTTL,
	})

	cacheStore, closeCache, err := openCache(context.Background(), cfg)
	if err != nil {
		logger.Error("open cache failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer closeCache()

	chainService := chain.NewService(repo)
	priceProvider := price.NewCachedProvider(price.NewDemoProvider(), cacheStore, cfg.PriceCacheTTL)
	priceService := price.NewService(priceProvider)
	multisigService := multisig.NewService(repo)
	server := httpserver.New(cfg, logger, chainService, authService, priceService, multisigService)

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
		Address:  cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
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
