package multisig

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidConfig = errors.New("invalid multisig config")
	ErrNotFound      = errors.New("multisig config not found")
)

type Config struct {
	ChainID         string   `json:"chain_id"`
	ContractAddress string   `json:"contract_address"`
	Owners          []string `json:"owners"`
	Threshold       uint64   `json:"threshold"`
}

type MultiSigStore interface {
	Save(ctx context.Context, cfg Config) error
	Get(ctx context.Context, chainID string) (Config, error)
}

type Service struct {
	store MultiSigStore
}

func NewService(store MultiSigStore) *Service {
	return &Service{store: store}
}

func (s *Service) Set(ctx context.Context, cfg Config) error {
	if cfg.ChainID == "" {
		return fmt.Errorf("%w: chain_id is required", ErrInvalidConfig)
	}
	if cfg.ContractAddress == "" {
		return fmt.Errorf("%w: contract_address is required", ErrInvalidConfig)
	}
	if len(cfg.Owners) == 0 {
		return fmt.Errorf("%w: owners are required", ErrInvalidConfig)
	}
	if cfg.Threshold == 0 || cfg.Threshold > uint64(len(cfg.Owners)) {
		return fmt.Errorf("%w: threshold must be between 1 and the number of owners", ErrInvalidConfig)
	}
	return s.store.Save(ctx, cfg)
}

func (s *Service) Get(ctx context.Context, chainID string) (Config, error) {
	if chainID == "" {
		return Config{}, fmt.Errorf("%w: chain_id is required", ErrInvalidConfig)
	}
	return s.store.Get(ctx, chainID)
}
