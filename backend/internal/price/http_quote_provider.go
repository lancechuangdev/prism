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
)

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

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return Quote{}, fmt.Errorf("create price request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if p.token != "" {
		request.Header.Set("Authorization", "Bearer "+p.token)
	}

	response, err := p.client.Do(request)
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
