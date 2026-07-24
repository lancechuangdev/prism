package chain

import (
	"context"

	"github.com/lancechuangdev/prism/backend/internal/store"
)

// Service provides read access to chain-indexed pool and token data.
type Service struct {
	repo QueryRepository
}

type QueryRepository interface {
	store.PoolRepository
	store.TokenRepository
}

func NewService(repo QueryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListPoolBases(ctx context.Context, chainID string) ([]store.PoolBase, error) {
	return s.repo.ListPoolBases(ctx, chainID)
}

func (s *Service) ListPoolData(ctx context.Context, chainID string) ([]store.PoolData, error) {
	return s.repo.ListPoolData(ctx, chainID)
}

func (s *Service) ListTokens(ctx context.Context, chainID string) ([]store.TokenInfo, error) {
	return s.repo.ListTokens(ctx, chainID)
}
