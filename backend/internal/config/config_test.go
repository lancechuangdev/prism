package config

import "testing"

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
