package price

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lancechuangdev/prism/backend/internal/resilience"
)

var quoteRetryPolicy = resilience.Policy{
	Attempts:       3,
	AttemptTimeout: 5 * time.Second,
	InitialBackoff: 100 * time.Millisecond,
	MaxBackoff:     time.Second,
}

type HTTPQuoteProvider struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewHTTPQuoteProvider(baseURL string, token string) (*HTTPQuoteProvider, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, fmt.Errorf("price provider URL must be an absolute HTTP or HTTPS URL")
	}
	return &HTTPQuoteProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (p *HTTPQuoteProvider) Latest(ctx context.Context, symbol string) (Quote, error) {
	endpoint, err := url.Parse(p.baseURL)
	if err != nil {
		return Quote{}, err
	}
	query := endpoint.Query()
	query.Set("symbol", symbol)
	endpoint.RawQuery = query.Encode()

	response, err := resilience.Value(ctx, quoteRetryPolicy, func(attemptCtx context.Context) (*http.Response, error) {
		request, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, endpoint.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("create price request: %w", err)
		}
		request.Header.Set("Accept", "application/json")
		if p.token != "" {
			request.Header.Set("Authorization", "Bearer "+p.token)
		}

		response, err := p.client.Do(request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= http.StatusInternalServerError {
			_ = response.Body.Close()
			return nil, fmt.Errorf("provider returned retryable status %s", response.Status)
		}
		return response, nil
	})
	if err != nil {
		return Quote{}, fmt.Errorf("fetch price for %s: %w", symbol, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Quote{}, fmt.Errorf("fetch price for %s: provider returned %s", symbol, response.Status)
	}

	var quote Quote
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&quote); err != nil {
		return Quote{}, fmt.Errorf("decode price for %s: %w", symbol, err)
	}
	if !strings.EqualFold(quote.Symbol, symbol) ||
		strings.TrimSpace(quote.Currency) == "" ||
		strings.TrimSpace(quote.Price) == "" ||
		strings.TrimSpace(quote.Source) == "" ||
		quote.UpdatedAt.IsZero() {
		return Quote{}, fmt.Errorf("price provider returned an invalid quote for %s", symbol)
	}
	numericPrice, ok := new(big.Rat).SetString(quote.Price)
	if !ok || numericPrice.Sign() <= 0 {
		return Quote{}, fmt.Errorf("price provider returned a non-positive price for %s", symbol)
	}
	quote.Symbol = strings.ToUpper(quote.Symbol)
	quote.Currency = strings.ToUpper(quote.Currency)
	return quote, nil
}
