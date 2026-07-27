package auth

import (
	"context"
	"errors"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/cache"
)

type CacheSessionStore struct {
	cache cache.Cache
}

func NewCacheSessionStore(cacheStore cache.Cache) *CacheSessionStore {
	return &CacheSessionStore{cache: cacheStore}
}

func (s *CacheSessionStore) Get(ctx context.Context, key string) (string, error) {
	value, err := s.cache.Get(ctx, key)
	if errors.Is(err, cache.ErrMiss) {
		return "", ErrSessionNotFound
	}
	return string(value), err
}

func (s *CacheSessionStore) Set(ctx context.Context, key string, username string, ttl time.Duration) error {
	return s.cache.Set(ctx, key, []byte(username), ttl)
}

func (s *CacheSessionStore) Delete(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, key)
}
