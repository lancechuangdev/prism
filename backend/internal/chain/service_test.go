package chain

import (
	"context"
	"testing"

	"github.com/lancechuangdev/prism/backend/internal/store"
)

func TestServiceListsPoolsAndTokens(t *testing.T) {
	ctx := context.Background()
	repo := store.NewMemoryStore()
	if err := SyncPools(ctx, NewDemoReader(), repo, "97"); err != nil {
		t.Fatalf("sync pools: %v", err)
	}

	service := NewService(repo)

	poolBases, err := service.ListPoolBases(ctx, "97")
	if err != nil {
		t.Fatalf("list pool bases: %v", err)
	}
	if len(poolBases) != 1 {
		t.Fatalf("expected one pool base, got %d", len(poolBases))
	}

	poolData, err := service.ListPoolData(ctx, "97")
	if err != nil {
		t.Fatalf("list pool data: %v", err)
	}
	if len(poolData) != 1 {
		t.Fatalf("expected one pool data record, got %d", len(poolData))
	}

	tokens, err := service.ListTokens(ctx, "97")
	if err != nil {
		t.Fatalf("list tokens: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected two tokens, got %d", len(tokens))
	}
}
