package config

import "testing"

func TestLoadChainRPCConfiguration(t *testing.T) {
	t.Setenv("PRISM_CHAIN_ID", "31337")
	t.Setenv("PRISM_CHAIN_RPC_URL", "http://node.example:8545")
	t.Setenv("PRISM_POOL_ADDRESS", "0x1000000000000000000000000000000000000001")

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
}
