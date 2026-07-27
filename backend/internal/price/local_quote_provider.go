package price

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type LocalQuoteProvider struct {
	quotes map[string]Quote
	now    func() time.Time
}

func NewLocalQuoteProvider() *LocalQuoteProvider {
	return &LocalQuoteProvider{
		quotes: map[string]Quote{
			"PRM": {
				Symbol:   "PRM",
				Currency: "USDT",
				Price:    "0.0027",
				Source:   "local",
			},
		},
		now: time.Now,
	}
}

func (p *LocalQuoteProvider) Latest(_ context.Context, symbol string) (Quote, error) {
	quote, ok := p.quotes[strings.ToUpper(symbol)]
	if !ok {
		return Quote{}, fmt.Errorf("price for %s not found", symbol)
	}

	quote.UpdatedAt = p.now().UTC()
	return quote, nil
}
