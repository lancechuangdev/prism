package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	defaultEnv           = "local"
	defaultPort          = "8080"
	defaultAPIVersion    = "1"
	defaultChainId       = "31337"
	defaultChainRPCURL   = "http://127.0.0.1:8545"
	defaultSyncEvery     = 2 * time.Minute
	defaultAdminUser     = "admin"
	defaultAdminPass     = "password"
	defaultTokenTTL      = time.Hour
	defaultTokenSecret   = "local-development-secret"
	defaultAuthMode      = "local"
	defaultProposalScope = "prism/proposals.write"
	defaultAdminScope    = "prism/admin.read"
	defaultPriceSymbol   = "PRM"
	defaultPriceProvider = "local"
	defaultStoreDriver   = "memory"
	defaultRedisAddr     = "127.0.0.1:6379"
	defaultPriceTTL      = 30 * time.Second

	ComponentAPI       = "api"
	ComponentScheduler = "scheduler"
	ComponentMigration = "migration"
)

type Config struct {
	Env                    string
	Port                   string
	APIVersion             string
	ChainID                string
	ChainRPCURL            string
	PoolAddress            string
	MultisigAddress        string
	SyncInterval           time.Duration
	LiquidationEnabled     bool
	LiquidationKey         string
	LiquidationSlippageBPS uint64
	AdminUsername          string
	AdminPassword          string
	TokenSecret            string
	TokenTTL               time.Duration
	AuthMode               string
	CognitoRegion          string
	CognitoUserPoolID      string
	CognitoClientID        string
	CORSAllowedOrigins     string
	ProposalWriteScope     string
	AdminReadScope         string
	PriceSymbol            string
	PriceProvider          string
	OracleAddress          string
	PriceTokenAddresses    string
	StoreDriver            string
	MySQLDSN               string
	RedisAddress           string
	RedisPassword          string
	RedisDB                int
	RedisTLS               bool
	RedisTLSServerName     string
	PriceCacheTTL          time.Duration
}

func Load() Config {
	return Config{
		Env:                    readEnv("PRISM_ENV", defaultEnv),
		Port:                   readEnv("PRISM_API_PORT", defaultPort),
		APIVersion:             readEnv("PRISM_API_VERSION", defaultAPIVersion),
		ChainID:                readEnv("PRISM_CHAIN_ID", defaultChainId),
		ChainRPCURL:            readEnv("PRISM_CHAIN_RPC_URL", defaultChainRPCURL),
		PoolAddress:            readEnv("PRISM_POOL_ADDRESS", ""),
		MultisigAddress:        readEnv("PRISM_MULTISIG_ADDRESS", ""),
		SyncInterval:           readDurationEnv("PRISM_SYNC_INTERVAL", defaultSyncEvery),
		LiquidationEnabled:     readBoolEnv("PRISM_LIQUIDATION_ENABLED", false),
		LiquidationKey:         readEnv("PRISM_LIQUIDATION_PRIVATE_KEY", ""),
		LiquidationSlippageBPS: readUint64Env("PRISM_LIQUIDATION_SLIPPAGE_BPS", 100),
		AdminUsername:          readEnv("PRISM_ADMIN_USERNAME", defaultAdminUser),
		AdminPassword:          readEnv("PRISM_ADMIN_PASSWORD", defaultAdminPass),
		TokenSecret:            readEnv("PRISM_TOKEN_SECRET", defaultTokenSecret),
		TokenTTL:               readDurationEnv("PRISM_TOKEN_TTL", defaultTokenTTL),
		AuthMode:               strings.ToLower(readEnv("PRISM_AUTH_MODE", defaultAuthMode)),
		CognitoRegion:          readEnv("PRISM_COGNITO_REGION", ""),
		CognitoUserPoolID:      readEnv("PRISM_COGNITO_USER_POOL_ID", ""),
		CognitoClientID:        readEnv("PRISM_COGNITO_CLIENT_ID", ""),
		CORSAllowedOrigins:     readEnv("PRISM_CORS_ALLOWED_ORIGINS", ""),
		ProposalWriteScope:     readEnv("PRISM_COGNITO_PROPOSAL_SCOPE", defaultProposalScope),
		AdminReadScope:         readEnv("PRISM_COGNITO_ADMIN_SCOPE", defaultAdminScope),
		PriceSymbol:            readEnv("PRISM_PRICE_SYMBOL", defaultPriceSymbol),
		PriceProvider:          readEnv("PRISM_PRICE_PROVIDER", defaultPriceProvider),
		OracleAddress:          readEnv("PRISM_ORACLE_ADDRESS", ""),
		PriceTokenAddresses:    readEnv("PRISM_PRICE_TOKEN_ADDRESSES", ""),
		StoreDriver:            strings.ToLower(readEnv("PRISM_STORE", defaultStoreDriver)),
		MySQLDSN:               readMySQLDSN(),
		RedisAddress:           readEnv("PRISM_REDIS_ADDR", defaultRedisAddr),
		RedisPassword:          readEnv("PRISM_REDIS_PASSWORD", ""),
		RedisDB:                readIntEnv("PRISM_REDIS_DB", 0),
		RedisTLS:               readBoolEnv("PRISM_REDIS_TLS", false),
		RedisTLSServerName:     readEnv("PRISM_REDIS_TLS_SERVER_NAME", ""),
		PriceCacheTTL:          readDurationEnv("PRISM_PRICE_CACHE_TTL", defaultPriceTTL),
	}
}

func readMySQLDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("PRISM_MYSQL_DSN")); dsn != "" {
		return dsn
	}

	host := strings.TrimSpace(os.Getenv("PRISM_MYSQL_HOST"))
	username := strings.TrimSpace(os.Getenv("PRISM_MYSQL_USERNAME"))
	password := os.Getenv("PRISM_MYSQL_PASSWORD")
	if host == "" && username == "" && password == "" {
		return ""
	}

	port := readEnv("PRISM_MYSQL_PORT", "3306")
	database := readEnv("PRISM_MYSQL_DATABASE", "prism")
	return (&mysql.Config{
		User:                 username,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		DBName:               database,
		ParseTime:            true,
		Loc:                  time.UTC,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}).FormatDSN()
}

// Validate checks configuration required by a specific executable. Production
// validation rejects development fallbacks instead of silently starting with
// ephemeral storage or known credentials.
func (c Config) Validate(component string) error {
	if component != ComponentAPI && component != ComponentScheduler && component != ComponentMigration {
		return fmt.Errorf("unknown component %q", component)
	}

	var problems []error
	if c.StoreDriver != "memory" && c.StoreDriver != "mysql" {
		problems = append(problems, fmt.Errorf("PRISM_STORE must be memory or mysql"))
	}
	if component == ComponentAPI && c.AuthMode != "" && c.AuthMode != "local" && c.AuthMode != "cognito" {
		problems = append(problems, fmt.Errorf("PRISM_AUTH_MODE must be local or cognito"))
	}
	if c.StoreDriver == "mysql" && strings.TrimSpace(c.MySQLDSN) == "" {
		problems = append(problems, fmt.Errorf("PRISM_MYSQL_DSN or complete PRISM_MYSQL_* connection fields are required when PRISM_STORE=mysql"))
	}
	if component == ComponentMigration && c.StoreDriver != "mysql" {
		problems = append(problems, fmt.Errorf("PRISM_STORE must be mysql for migrations"))
	}
	if component == ComponentScheduler && c.LiquidationEnabled && strings.TrimSpace(c.LiquidationKey) == "" {
		problems = append(problems, fmt.Errorf("PRISM_LIQUIDATION_PRIVATE_KEY is required when PRISM_LIQUIDATION_ENABLED=true"))
	}
	if c.LiquidationSlippageBPS > 10_000 {
		problems = append(problems, fmt.Errorf("PRISM_LIQUIDATION_SLIPPAGE_BPS must not exceed 10000"))
	}

	if !strings.EqualFold(c.Env, "production") {
		return errors.Join(problems...)
	}

	if c.StoreDriver != "mysql" {
		problems = append(problems, fmt.Errorf("PRISM_STORE must be mysql in production"))
	}
	if component == ComponentMigration {
		return errors.Join(problems...)
	}

	if strings.TrimSpace(c.PoolAddress) == "" {
		problems = append(problems, fmt.Errorf("PRISM_POOL_ADDRESS is required in production"))
	}
	if !c.RedisTLS {
		problems = append(problems, fmt.Errorf("PRISM_REDIS_TLS must be true in production"))
	}
	if !strings.EqualFold(c.PriceProvider, "chainlink") {
		problems = append(problems, fmt.Errorf("PRISM_PRICE_PROVIDER must be chainlink in production"))
	}
	if strings.TrimSpace(c.OracleAddress) == "" {
		problems = append(problems, fmt.Errorf("PRISM_ORACLE_ADDRESS is required in production"))
	}
	if strings.TrimSpace(c.PriceTokenAddresses) == "" {
		problems = append(problems, fmt.Errorf("PRISM_PRICE_TOKEN_ADDRESSES is required in production"))
	}
	if component == ComponentAPI {
		if strings.TrimSpace(c.MultisigAddress) == "" {
			problems = append(problems, fmt.Errorf("PRISM_MULTISIG_ADDRESS is required in production"))
		}
		if strings.TrimSpace(c.CORSAllowedOrigins) == "" {
			problems = append(problems, fmt.Errorf("PRISM_CORS_ALLOWED_ORIGINS is required in production"))
		}
		if c.AuthMode == "cognito" {
			if strings.TrimSpace(c.CognitoRegion) == "" {
				problems = append(problems, fmt.Errorf("PRISM_COGNITO_REGION is required with Cognito authentication"))
			}
			if strings.TrimSpace(c.CognitoUserPoolID) == "" {
				problems = append(problems, fmt.Errorf("PRISM_COGNITO_USER_POOL_ID is required with Cognito authentication"))
			}
			if strings.TrimSpace(c.CognitoClientID) == "" {
				problems = append(problems, fmt.Errorf("PRISM_COGNITO_CLIENT_ID is required with Cognito authentication"))
			}
			if strings.TrimSpace(c.ProposalWriteScope) == "" {
				problems = append(problems, fmt.Errorf("PRISM_COGNITO_PROPOSAL_SCOPE is required with Cognito authentication"))
			}
			if strings.TrimSpace(c.AdminReadScope) == "" {
				problems = append(problems, fmt.Errorf("PRISM_COGNITO_ADMIN_SCOPE is required with Cognito authentication"))
			}
		} else {
			if strings.TrimSpace(c.AdminUsername) == "" || strings.TrimSpace(c.AdminUsername) == defaultAdminUser {
				problems = append(problems, fmt.Errorf("PRISM_ADMIN_USERNAME must not use the development default in production"))
			}
			if len(c.AdminPassword) < 12 || strings.TrimSpace(c.AdminPassword) == defaultAdminPass {
				problems = append(problems, fmt.Errorf("PRISM_ADMIN_PASSWORD must contain at least 12 characters and not use the development default in production"))
			}
			if len(c.TokenSecret) < 32 || strings.TrimSpace(c.TokenSecret) == defaultTokenSecret {
				problems = append(problems, fmt.Errorf("PRISM_TOKEN_SECRET must contain at least 32 characters and not use the development default in production"))
			}
		}
	}

	return errors.Join(problems...)
}

func readBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func readUint64Env(key string, fallback uint64) uint64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
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
