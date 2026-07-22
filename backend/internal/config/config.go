package config

import (
	"os"
	"time"
)

const (
	defaultEnv        = "local"
	defaultPort       = "8080"
	defaultAPIVersion = "1"
	defaultChainId    = "97"
	defaultSyncEvery  = 2 * time.Minute
)

type Config struct {
	Env         string
	Port        string
	APIVersion  string
	ChainID     string
	SyncInteral time.Duration
}

func Load() Config {
	return Config{
		Env:         readEnv("PRISM_ENV", defaultEnv),
		Port:        readEnv("PRISM_API_PORT", defaultPort),
		APIVersion:  readEnv("PRISM_API_VERSION", defaultAPIVersion),
		ChainID:     readEnv("PRISM_CHAIN_ID", defaultChainId),
		SyncInteral: readDurationEnv("PRISM_SYNC_INTERVAL", defaultSyncEvery),
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
