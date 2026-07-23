package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultEnv         = "local"
	defaultPort        = "8080"
	defaultAPIVersion  = "1"
	defaultChainId     = "97"
	defaultSyncEvery   = 2 * time.Minute
	defaultAdminUser   = "admin"
	defaultAdminPass   = "password"
	defaultTokenTTL    = time.Hour
	defaultTokenSecret = "local-development-secret"
	defaultPriceSymbol = "PRM"
	defaultStoreDriver = "memory"
	defaultRedisAddr   = "127.0.0.1:6379"
	defaultPriceTTL    = 30 * time.Second
)

type Config struct {
	Env           string
	Port          string
	APIVersion    string
	ChainID       string
	SyncInteral   time.Duration
	AdminUsername string
	AdminPassword string
	TokenSecret   string
	TokenTTL      time.Duration
	PriceSymbol   string
	StoreDriver   string
	MySQLDSN      string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	PriceCacheTTL time.Duration
}

func Load() Config {
	return Config{
		Env:           readEnv("PRISM_ENV", defaultEnv),
		Port:          readEnv("PRISM_API_PORT", defaultPort),
		APIVersion:    readEnv("PRISM_API_VERSION", defaultAPIVersion),
		ChainID:       readEnv("PRISM_CHAIN_ID", defaultChainId),
		SyncInteral:   readDurationEnv("PRISM_SYNC_INTERVAL", defaultSyncEvery),
		AdminUsername: readEnv("PRISM_ADMIN_USERNAME", defaultAdminUser),
		AdminPassword: readEnv("PRISM_ADMIN_PASSWORD", defaultAdminPass),
		TokenSecret:   readEnv("PRISM_TOKEN_SECRET", defaultTokenSecret),
		TokenTTL:      readDurationEnv("PRISM_TOKEN_TTL", defaultTokenTTL),
		PriceSymbol:   readEnv("PRISM_PRICE_SYMBOL", defaultPriceSymbol),
		StoreDriver:   strings.ToLower(readEnv("PRISM_STORE", defaultStoreDriver)),
		MySQLDSN:      readEnv("PRISM_MYSQL_DSN", ""),
		RedisAddress:  readEnv("PRISM_REDIS_ADDR", defaultRedisAddr),
		RedisPassword: readEnv("PRISM_REDIS_PASSWORD", ""),
		RedisDB:       readIntEnv("PRISM_REDIS_DB", 0),
		PriceCacheTTL: readDurationEnv("PRISM_PRICE_CACHE_TTL", defaultPriceTTL),
	}
}

func readEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func readDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}

func readIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return fallback
	}

	return parsed
}
