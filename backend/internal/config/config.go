package config

import "os"

const (
	defaultEnv        = "local"
	defaultPort       = "8080"
	defaultAPIVersion = "1"
)

type Config struct {
	Env        string
	Port       string
	APIVersion string
}

func Load() Config {
	return Config{
		Env:        readEnv("PRISM_ENV", defaultEnv),
		Port:       readEnv("PRISM_API_PORT", defaultPort),
		APIVersion: readEnv("PRISM_API_VERSION", defaultAPIVersion),
	}
}

func readEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
