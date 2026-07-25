package multisig

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	ctx := context.Background()
	memoryStore := newFakeStore()
	now := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	memoryStore.now = func() time.Time { return now }
	service := NewService(memoryStore)

	err := service.Set(ctx, Config{
		ChainID:         "97",
		ContractAddress: "0xmultisig",
		Owners:          []string{"0xowner1", "0xowner2"},
		Threshold:       2,
	})
	if err != nil {
		t.Fatalf("set multisig: %v", err)
	}

	cfg, err := service.Get(ctx, "97")
	if err != nil {
		t.Fatalf("get multisig: %v", err)
	}
	if cfg.ContractAddress != "0xmultisig" {
		t.Fatalf("unexpected contract address: %s", cfg.ContractAddress)
	}
	if len(cfg.Owners) != 2 || cfg.Threshold != 2 {
		t.Fatalf("unexpected multisig config: %+v", cfg)
	}
}

func TestSetValidatesRequiredFieldsAndThreshold(t *testing.T) {
	service := NewService(newFakeStore())

	err := service.Set(context.Background(), Config{
		ChainID:         "97",
		ContractAddress: "0xmultisig",
		Owners:          []string{"0xowner1"},
		Threshold:       2,
	})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected invalid config, got %v", err)
	}
}

func TestGetReturnsNotFound(t *testing.T) {
	service := NewService(newFakeStore())

	_, err := service.Get(context.Background(), "97")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

type fakeStore struct {
	records map[string]Config
	now     func() time.Time
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		records: make(map[string]Config),
		now:     time.Now,
	}
}

func (s *fakeStore) Save(_ context.Context, cfg Config) error {
	s.records[cfg.ChainID] = cfg
	return nil
}

func (s *fakeStore) Get(_ context.Context, chainID string) (Config, error) {
	cfg, ok := s.records[chainID]
	if !ok {
		return Config{}, ErrNotFound
	}
	return cfg, nil
}
