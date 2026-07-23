package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/auth"
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

	repo := store.NewMemoryStore()
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

	priceService := price.NewService(price.NewDemoProvider())
	multisigService := multisig.NewService(repo)
	server := httpserver.New(cfg, logger, repo, authService, priceService, multisigService)

	go func() {
		logger.Info("api server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	waitForShutdown(server, logger)
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
