package config

import (
	"strings"
	"testing"
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

func TestValidateProductionAPI(t *testing.T) {
	cfg := validProductionConfig()

	if err := cfg.Validate(ComponentAPI); err != nil {
		t.Fatalf("valid production API configuration: %v", err)
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
		"PRISM_REDIS_TLS must be true",
		"PRISM_PRICE_PROVIDER must be http",
		"PRISM_PRICE_PROVIDER_URL must use HTTPS",
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
		Env:              "production",
		PoolAddress:      "0x1000000000000000000000000000000000000001",
		MultisigAddress:  "0x2000000000000000000000000000000000000002",
		AdminUsername:    "prism-admin",
		AdminPassword:    "non-default-password",
		TokenSecret:      "non-default-token-secret-at-least-32-bytes",
		StoreDriver:      "mysql",
		MySQLDSN:         "user:password@tcp(mysql:3306)/prism",
		RedisTLS:         true,
		PriceProvider:    "http",
		PriceProviderURL: "https://quotes.example.com",
	}
}
