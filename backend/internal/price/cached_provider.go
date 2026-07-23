package price

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/cache"
)

type CachedProvider struct {
	next  Provider
	cache cache.Cache
	ttl   time.Duration
}

func NewCachedProvider(next Provider, cache cache.Cache, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		next:  next,
		cache: cache,
		ttl:   ttl,
	}
}

func (p *CachedProvider) Latest(ctx context.Context, symbol string) (Quote, error) {
	key := "price:" + symbol
	cachedValue, err := p.cache.Get(ctx, key)
	if err == nil {
		quote := Quote{}
		if err := json.Unmarshal(cachedValue, &quote); err == nil {
			return quote, nil
		}
	}
	if err != nil && errors.Is(err, cache.ErrMiss) {
		quote, err := p.next.Latest(ctx, symbol)
		if err != nil {
			return Quote{}, err
		}

		value, err := json.Marshal(quote)
		if err != nil {
			return Quote{}, err
		}

		if err := p.cache.Set(ctx, key, value, p.ttl); err != nil {
			return Quote{}, err
		}

		return quote, nil
	}

	return Quote{}, err
}
