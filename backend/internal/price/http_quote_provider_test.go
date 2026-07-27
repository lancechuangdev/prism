package price

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHTTPQuoteProviderFetchesQuote(t *testing.T) {
	updatedAt := time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC)
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("symbol") != "PRM" {
			t.Fatalf("unexpected symbol query %q", r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer secret-token" {
			t.Fatalf("unexpected authorization header %q", r.Header.Get("Authorization"))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{
			"symbol":"prm",
			"currency":"USD",
			"price":"0.0027",
			"source":"test-provider",
			"updatedAt":"2026-07-26T12:00:00Z"
		}`)),
			Request: r,
		}, nil
	})

	provider, err := NewHTTPQuoteProvider("https://prices.example/v1/quote", "secret-token")
	if err != nil {
		t.Fatal(err)
	}
	provider.client.Transport = transport
	quote, err := provider.Latest(context.Background(), "PRM")
	if err != nil {
		t.Fatalf("fetch quote: %v", err)
	}
	if quote.Symbol != "PRM" ||
		quote.Currency != "USD" ||
		quote.Price != "0.0027" ||
		quote.Source != "test-provider" ||
		!quote.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected quote: %+v", quote)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestConfiguredProviderRejectsDemoAndInsecureHTTPInProduction(t *testing.T) {
	if _, err := NewConfiguredQuoteProvider("production", "local", "", ""); err == nil {
		t.Fatal("expected production local provider to be rejected")
	}
	if _, err := NewConfiguredQuoteProvider("production", "http", "http://prices.example", ""); err == nil {
		t.Fatal("expected insecure production provider URL to be rejected")
	}
}

func TestConfiguredProviderAllowsDemoOutsideProduction(t *testing.T) {
	provider, err := NewConfiguredQuoteProvider("local", "local", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := provider.(*LocalQuoteProvider); !ok {
		t.Fatalf("expected LocalQuoteProvider, got %T", provider)
	}
}
