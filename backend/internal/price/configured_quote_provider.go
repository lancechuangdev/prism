package price

import (
	"context"
	"fmt"
	"strings"
)

func NewConfiguredQuoteProvider(
	ctx context.Context,
	environment string,
	providerType string,
	rpcURL string,
	oracleAddress string,
	tokenAddressesJSON string,
) (QuoteProvider, error) {
	switch strings.ToLower(strings.TrimSpace(providerType)) {
	case "local":
		if strings.EqualFold(environment, "production") {
			return nil, fmt.Errorf("local price provider is not allowed in production")
		}
		return NewLocalQuoteProvider(), nil
	case "chainlink":
		return NewChainlinkQuoteProvider(ctx, rpcURL, oracleAddress, tokenAddressesJSON)
	default:
		return nil, fmt.Errorf("unsupported price provider %q", providerType)
	}
}

func CloseQuoteProvider(provider QuoteProvider) {
	if closer, ok := provider.(interface{ Close() }); ok {
		closer.Close()
	}
}
