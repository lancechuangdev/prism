package price

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/cache"
)

func TestCachedProviderUsesCache(t *testing.T) {
	ctx := context.Background()
	cacheStore := newFakeCache()
	provider := NewDemoProvider()
	start := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	provider.now = func() time.Time { return start }

	cached := NewCachedProvider(provider, cacheStore, time.Minute)

	first, err := cached.Latest(ctx, "PRM")
	if err != nil {
		t.Fatalf("first latest: %v", err)
	}

	provider.quotes["PRM"] = Quote{
		Symbol:   "PRM",
		Currency: "USDT",
		Price:    "100",
		Source:   "demo",
	}

	second, err := cached.Latest(ctx, "PRM")
	if err != nil {
		t.Fatalf("second latest: %v", err)
	}

	if first.Price != second.Price {
		t.Fatalf("expected cached price %s, got %s", first.Price, second.Price)
	}
}

func TestCachedProviderFetchesNewValueAfterTTL(t *testing.T) {
	ctx := context.Background()
	cacheStore := newFakeCache()
	provider := NewDemoProvider()
	start := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	cacheStore.now = func() time.Time { return start }
	provider.now = func() time.Time { return start }

	cached := NewCachedProvider(provider, cacheStore, time.Minute)

	first, err := cached.Latest(ctx, "PRM")
	if err != nil {
		t.Fatalf("first latest: %v", err)
	}

	provider.quotes["PRM"] = Quote{
		Symbol:   "PRM",
		Currency: "USDT",
		Price:    "100",
		Source:   "demo",
	}
	cacheStore.now = func() time.Time { return start.Add(time.Minute) }
	provider.now = func() time.Time { return start.Add(time.Minute) }

	second, err := cached.Latest(ctx, "PRM")
	if err != nil {
		t.Fatalf("second latest: %v", err)
	}

	if second.Price != "100" {
		t.Fatalf("expected refreshed price 100, got %s", second.Price)
	}
	if second.Price == first.Price {
		t.Fatalf("expected price to change after TTL, got %s", second.Price)
	}
}

type fakeCacheRecord struct {
	value     []byte
	expiresAt time.Time
}

type fakeCache struct {
	records map[string]fakeCacheRecord
	now     func() time.Time
}

func newFakeCache() *fakeCache {
	return &fakeCache{
		records: make(map[string]fakeCacheRecord),
		now:     time.Now,
	}
}

func (c *fakeCache) Get(_ context.Context, key string) ([]byte, error) {
	record, ok := c.records[key]
	if !ok {
		return nil, cache.ErrMiss
	}
	if !c.now().Before(record.expiresAt) {
		delete(c.records, key)
		return nil, cache.ErrMiss
	}
	return append([]byte(nil), record.value...), nil
}

func (c *fakeCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.records[key] = fakeCacheRecord{
		value:     append([]byte(nil), value...),
		expiresAt: c.now().Add(ttl),
	}
	return nil
}

func (c *fakeCache) Delete(_ context.Context, key string) error {
	if _, ok := c.records[key]; !ok {
		return errors.New("missing key")
	}
	delete(c.records, key)
	return nil
}
