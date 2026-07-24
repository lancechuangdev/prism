package chain

import (
	"context"

	"github.com/lancechuangdev/prism/backend/internal/store"
)

// QueryService provides read access to chain-indexed pool and token data.
type QueryService struct {
	repo QueryRepository
}

type QueryRepository interface {
	store.PoolRepository
	store.TokenRepository
}

func NewQueryService(repo QueryRepository) *QueryService {
	return &QueryService{repo: repo}
}

func (s *QueryService) ListPoolBases(ctx context.Context, chainID string) ([]store.PoolBase, error) {
	return s.repo.ListPoolBases(ctx, chainID)
}

func (s *QueryService) ListPoolData(ctx context.Context, chainID string) ([]store.PoolData, error) {
	return s.repo.ListPoolData(ctx, chainID)
}

func (s *QueryService) ListTokens(ctx context.Context, chainID string) ([]store.TokenInfo, error) {
	return s.repo.ListTokens(ctx, chainID)
}
