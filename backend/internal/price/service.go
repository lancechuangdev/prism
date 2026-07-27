package price

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Quote struct {
	Symbol    string    `json:"symbol"`
	Currency  string    `json:"currency"`
	Price     string    `json:"price"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type QuoteProvider interface {
	Latest(ctx context.Context, symbol string) (Quote, error)
}

type Service struct {
	provider QuoteProvider
}

func NewService(provider QuoteProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) Latest(ctx context.Context, symbol string) (Quote, error) {
	normalizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
	if normalizedSymbol == "" {
		return Quote{}, fmt.Errorf("symbol is required")
	}
	return s.provider.Latest(ctx, normalizedSymbol)
}
