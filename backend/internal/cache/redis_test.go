package cache

import (
	"crypto/tls"
	"testing"
	"time"
)

func TestRedisOptionsEnableVerifiedTLS(t *testing.T) {
	options, err := redisOptions(RedisConfig{
		Address:    "master.prism.cache.amazonaws.com:6379",
		TLSEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.TLSConfig == nil {
		t.Fatal("expected TLS configuration")
	}
	if options.TLSConfig.MinVersion != tls.VersionTLS12 {
		t.Fatalf("unexpected minimum TLS version %d", options.TLSConfig.MinVersion)
	}
	if options.TLSConfig.ServerName != "master.prism.cache.amazonaws.com" {
		t.Fatalf("unexpected TLS server name %q", options.TLSConfig.ServerName)
	}
	if options.TLSConfig.InsecureSkipVerify {
		t.Fatal("certificate verification must remain enabled")
	}
	if options.DialTimeout != 5*time.Second ||
		options.ReadTimeout != 3*time.Second ||
		options.WriteTimeout != 3*time.Second ||
		options.MaxRetries != 2 ||
		options.MinRetryBackoff != 100*time.Millisecond ||
		options.MaxRetryBackoff != time.Second {
		t.Fatalf("unexpected Redis resilience options: %+v", options)
	}
}

func TestRedisOptionsRejectPlaintextWhenTLSRequired(t *testing.T) {
	if _, err := redisOptions(RedisConfig{
		Address:    "redis:6379",
		RequireTLS: true,
	}); err == nil {
		t.Fatal("expected plaintext Redis to be rejected")
	}
}

func TestRedisOptionsAllowPlaintextOutsideProduction(t *testing.T) {
	options, err := redisOptions(RedisConfig{Address: "redis:6379"})
	if err != nil {
		t.Fatal(err)
	}
	if options.TLSConfig != nil {
		t.Fatal("expected plaintext local Redis configuration")
	}
}
