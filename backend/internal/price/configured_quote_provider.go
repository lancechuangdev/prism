package price

import (
	"fmt"
	"strings"
)

func NewConfiguredQuoteProvider(environment string, providerType string, providerURL string, providerToken string) (QuoteProvider, error) {
	switch strings.ToLower(strings.TrimSpace(providerType)) {
	case "local":
		if strings.EqualFold(environment, "production") {
			return nil, fmt.Errorf("local price provider is not allowed in production")
		}
		return NewLocalQuoteProvider(), nil
	case "http":
		if strings.EqualFold(environment, "production") && !strings.HasPrefix(strings.ToLower(providerURL), "https://") {
			return nil, fmt.Errorf("production price provider URL must use HTTPS")
		}
		return NewHTTPQuoteProvider(providerURL, providerToken)
	default:
		return nil, fmt.Errorf("unsupported price provider %q", providerType)
	}
}
