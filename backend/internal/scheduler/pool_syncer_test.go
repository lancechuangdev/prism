package scheduler

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/chain"
	"github.com/lancechuangdev/prism/backend/internal/price"
	"github.com/lancechuangdev/prism/backend/internal/store"
)

type fakeQuoteProvider struct{}

func (fakeQuoteProvider) Latest(_ context.Context, symbol string) (price.Quote, error) {
	prices := map[string]string{"BUSD": "1.00", "BTC": "42000.00"}
	return price.Quote{
		Symbol: symbol, Currency: "USD", Price: prices[symbol], Source: "test", UpdatedAt: time.Now(),
	}, nil
}

type fakeLiquidationChecker struct {
	called bool
}

func (f *fakeLiquidationChecker) CheckAndLiquidate(context.Context) error {
	f.called = true
	return nil
}

func TestPoolSyncerRunOnce(t *testing.T) {
	ctx := context.Background()
	repo := store.NewMemoryStore()
	reader := chain.NewFakeReader()
	prices := price.NewService(fakeQuoteProvider{})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	liquidations := &fakeLiquidationChecker{}
	syncer := NewPoolSyncer(reader, repo, "31337", prices, liquidations, logger)

	if err := syncer.RunOnce(ctx); err != nil {
		t.Fatalf("run once: %v", err)
	}

	pools, err := repo.ListPoolBases(ctx, "31337")
	if err != nil {
		t.Fatalf("list pool bases: %v", err)
	}
	if len(pools) != 1 {
		t.Fatalf("expected one synced pool, got %d", len(pools))
	}

	tokens, err := repo.ListTokens(ctx, "31337")
	if err != nil {
		t.Fatalf("list tokens: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected two synced tokens, got %d", len(tokens))
	}
	for _, token := range tokens {
		if token.Price == "" {
			t.Fatalf("expected %s price to be persisted", token.Symbol)
		}
	}
	if pools[0].LendToken.Price != "1.00" {
		t.Fatalf("lend token price = %q, want 1.00", pools[0].LendToken.Price)
	}
	if pools[0].CollateralToken.Price != "42000.00" {
		t.Fatalf("collateral token price = %q, want 42000.00", pools[0].CollateralToken.Price)
	}
	if !liquidations.called {
		t.Fatal("expected liquidation checker to run")
	}
}
