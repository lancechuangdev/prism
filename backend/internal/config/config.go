package config

import (
	"os"
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
}

func Load() Config {
	return Config{
		Env:           readEnv("PRISM_ENV", defaultEnv),
		Port:          readEnv("PRISM_API_PORT", defaultPort),
		APIVersion:    readEnv("PRISM_API_VERSION", defaultAPIVersion),
		ChainID:       readEnv("PRISM_CHAIN_ID", defaultChainId),
		SyncInteral:   readDurationEnv("PRISM_SYNC_INTERVAL", defaultSyncEvery),
		AdminUsername: readEnv("PLEDGE_ADMIN_USERNAME", defaultAdminUser),
		AdminPassword: readEnv("PLEDGE_ADMIN_PASSWORD", defaultAdminPass),
		TokenSecret:   readEnv("PLEDGE_TOKEN_SECRET", defaultTokenSecret),
		TokenTTL:      readDurationEnv("PLEDGE_TOKEN_TTL", defaultTokenTTL),
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
