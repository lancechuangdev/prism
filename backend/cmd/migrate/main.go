package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lancechuangdev/prism/backend/internal/config"
	"github.com/lancechuangdev/prism/backend/internal/logging"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)
	if err := cfg.Validate(config.ComponentMigration); err != nil {
		logger.Error("invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mysqlStore, err := store.OpenMySQL(ctx, cfg.MySQLDSN)
	if err != nil {
		logger.Error("open MySQL for migration failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer mysqlStore.Close()

	if err := mysqlStore.Migrate(ctx); err != nil {
		logger.Error("database migration failed", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("database migrations completed")
}
