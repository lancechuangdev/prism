package config

import (
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestLoadChainRPCConfiguration(t *testing.T) {
	t.Setenv("PRISM_CHAIN_ID", "31337")
	t.Setenv("PRISM_CHAIN_RPC_URL", "http://node.example:8545")
	t.Setenv("PRISM_POOL_ADDRESS", "0x1000000000000000000000000000000000000001")
	t.Setenv("PRISM_MULTISIG_ADDRESS", "0x2000000000000000000000000000000000000002")

	cfg := Load()

	if cfg.ChainID != "31337" {
		t.Fatalf("chain ID = %q", cfg.ChainID)
	}
	if cfg.ChainRPCURL != "http://node.example:8545" {
		t.Fatalf("chain RPC URL = %q", cfg.ChainRPCURL)
	}
	if cfg.PoolAddress != "0x1000000000000000000000000000000000000001" {
		t.Fatalf("pool address = %q", cfg.PoolAddress)
	}
	if cfg.MultisigAddress != "0x2000000000000000000000000000000000000002" {
		t.Fatalf("multisig address = %q", cfg.MultisigAddress)
	}
}

func TestLoadRedisTLSConfiguration(t *testing.T) {
	t.Setenv("PRISM_REDIS_TLS", "true")
	t.Setenv("PRISM_REDIS_TLS_SERVER_NAME", "master.prism.cache.amazonaws.com")

	cfg := Load()

	if !cfg.RedisTLS {
		t.Fatal("expected Redis TLS to be enabled")
	}
	if cfg.RedisTLSServerName != "master.prism.cache.amazonaws.com" {
		t.Fatalf("Redis TLS server name = %q", cfg.RedisTLSServerName)
	}
}

func TestLoadCORSConfiguration(t *testing.T) {
	t.Setenv("PRISM_CORS_ALLOWED_ORIGINS", "https://prismapp.link")

	cfg := Load()

	if cfg.CORSAllowedOrigins != "https://prismapp.link" {
		t.Fatalf("CORS allowed origins = %q", cfg.CORSAllowedOrigins)
	}
}

func TestLoadBuildsMySQLDSNFromSecretFields(t *testing.T) {
	t.Setenv("PRISM_MYSQL_DSN", "")
	t.Setenv("PRISM_MYSQL_HOST", "database.internal")
	t.Setenv("PRISM_MYSQL_PORT", "3307")
	t.Setenv("PRISM_MYSQL_DATABASE", "backend")
	t.Setenv("PRISM_MYSQL_USERNAME", "prism")
	t.Setenv("PRISM_MYSQL_PASSWORD", "p@ss:word")

	cfg := Load()

	parsed, err := mysql.ParseDSN(cfg.MySQLDSN)
	if err != nil {
		t.Fatalf("parse generated MySQL DSN: %v", err)
	}
	if parsed.User != "prism" || parsed.Passwd != "p@ss:word" ||
		parsed.Addr != "database.internal:3307" || parsed.DBName != "backend" ||
		!parsed.ParseTime {
		t.Fatalf("unexpected generated MySQL configuration: %#v", parsed)
	}
}

func TestValidateProductionAPI(t *testing.T) {
	cfg := validProductionConfig()

	if err := cfg.Validate(ComponentAPI); err != nil {
		t.Fatalf("valid production API configuration: %v", err)
	}
}

func TestValidateProductionCognitoAPI(t *testing.T) {
	cfg := validProductionConfig()
	cfg.AuthMode = "cognito"
	cfg.AdminUsername = defaultAdminUser
	cfg.AdminPassword = defaultAdminPass
	cfg.TokenSecret = defaultTokenSecret
	cfg.CognitoRegion = "us-west-2"
	cfg.CognitoUserPoolID = "us-west-2_example"
	cfg.CognitoClientID = "client-id"
	cfg.ProposalWriteScope = "prism/proposals.write"
	cfg.AdminReadScope = "prism/admin.read"

	if err := cfg.Validate(ComponentAPI); err != nil {
		t.Fatalf("valid Cognito production API configuration: %v", err)
	}
}

func TestValidateProductionAPIRejectsUnsafeDefaults(t *testing.T) {
	cfg := Load()
	cfg.Env = "production"

	err := cfg.Validate(ComponentAPI)
	if err == nil {
		t.Fatal("expected unsafe production configuration to fail")
	}
	for _, expected := range []string{
		"PRISM_STORE must be mysql",
		"PRISM_POOL_ADDRESS is required",
		"PRISM_MULTISIG_ADDRESS is required",
		"PRISM_CORS_ALLOWED_ORIGINS is required",
		"PRISM_REDIS_TLS must be true",
		"PRISM_PRICE_PROVIDER must be chainlink",
		"PRISM_ORACLE_ADDRESS is required",
		"PRISM_PRICE_TOKEN_ADDRESSES is required",
		"PRISM_ADMIN_USERNAME",
		"PRISM_ADMIN_PASSWORD",
		"PRISM_TOKEN_SECRET",
	} {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("validation error %q does not contain %q", err, expected)
		}
	}
}

func TestValidateProductionSchedulerDoesNotRequireAPISecrets(t *testing.T) {
	cfg := validProductionConfig()
	cfg.AdminUsername = defaultAdminUser
	cfg.AdminPassword = defaultAdminPass
	cfg.TokenSecret = defaultTokenSecret
	cfg.MultisigAddress = ""

	if err := cfg.Validate(ComponentScheduler); err != nil {
		t.Fatalf("valid production scheduler configuration: %v", err)
	}
}

func TestValidateSchedulerRequiresKeyWhenLiquidationEnabled(t *testing.T) {
	cfg := Config{StoreDriver: "memory", LiquidationEnabled: true}

	err := cfg.Validate(ComponentScheduler)
	if err == nil || !strings.Contains(err.Error(), "PRISM_LIQUIDATION_PRIVATE_KEY") {
		t.Fatalf("expected missing liquidation key error, got %v", err)
	}

	cfg.LiquidationKey = "test-only-key"
	if err := cfg.Validate(ComponentScheduler); err != nil {
		t.Fatalf("valid local liquidation configuration: %v", err)
	}
}

func TestValidateMigrationOnlyRequiresMySQL(t *testing.T) {
	cfg := Config{
		Env:         "production",
		StoreDriver: "mysql",
		MySQLDSN:    "user:password@tcp(mysql:3306)/prism",
	}
	if err := cfg.Validate(ComponentMigration); err != nil {
		t.Fatalf("valid migration configuration: %v", err)
	}

	cfg.StoreDriver = "memory"
	if err := cfg.Validate(ComponentMigration); err == nil {
		t.Fatal("expected memory-backed migration configuration to fail")
	}
}

func TestValidateRejectsUnknownComponent(t *testing.T) {
	if err := (Config{}).Validate("worker"); err == nil {
		t.Fatal("expected unknown component to fail")
	}
}

func validProductionConfig() Config {
	return Config{
		Env:                 "production",
		PoolAddress:         "0x1000000000000000000000000000000000000001",
		MultisigAddress:     "0x2000000000000000000000000000000000000002",
		AdminUsername:       "prism-admin",
		AdminPassword:       "non-default-password",
		TokenSecret:         "non-default-token-secret-at-least-32-bytes",
		StoreDriver:         "mysql",
		MySQLDSN:            "user:password@tcp(mysql:3306)/prism",
		RedisTLS:            true,
		PriceProvider:       "chainlink",
		OracleAddress:       "0x3000000000000000000000000000000000000003",
		PriceTokenAddresses: `{"USDC":"0x4000000000000000000000000000000000000004"}`,
		CORSAllowedOrigins:  "https://prismapp.link",
	}
}
