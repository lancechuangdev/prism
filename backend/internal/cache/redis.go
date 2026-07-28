package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address       string
	Password      string
	DB            int
	TLSEnabled    bool
	TLSServerName string
	RequireTLS    bool
}

type RedisCache struct {
	client *redis.Client
}

func OpenRedis(ctx context.Context, cfg RedisConfig) (*RedisCache, error) {
	options, err := redisOptions(cfg)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func redisOptions(cfg RedisConfig) (*redis.Options, error) {
	if cfg.RequireTLS && !cfg.TLSEnabled {
		return nil, fmt.Errorf("Redis TLS is required in production")
	}
	options := &redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	if cfg.TLSEnabled {
		serverName := cfg.TLSServerName
		if serverName == "" {
			host, _, err := net.SplitHostPort(cfg.Address)
			if err != nil {
				return nil, fmt.Errorf("derive Redis TLS server name from address %q: %w", cfg.Address, err)
			}
			serverName = host
		}
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: serverName,
		}
	}
	return options, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrMiss
	}
	return value, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
