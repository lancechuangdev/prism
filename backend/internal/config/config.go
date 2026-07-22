package config

import "os"

const (
	defaultEnv        = "local"
	defaultPort       = "8080"
	defaultAPIVersion = "1"
	defaultChainId    = "97"
)

type Config struct {
	Env        string
	Port       string
	APIVersion string
	ChainID    string
}

func Load() Config {
	return Config{
		Env:        readEnv("PRISM_ENV", defaultEnv),
		Port:       readEnv("PRISM_API_PORT", defaultPort),
		APIVersion: readEnv("PRISM_API_VERSION", defaultAPIVersion),
		ChainID:    readEnv("PRISM_CHAIN_ID", defaultChainId),
	}
}

func readEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
