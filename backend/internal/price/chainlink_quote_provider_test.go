package price

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type stubChainlinkOracle struct {
	price     *big.Int
	updatedAt time.Time
	err       error
	token     common.Address
}

func (s *stubChainlinkOracle) Latest(_ context.Context, token common.Address) (*big.Int, time.Time, error) {
	s.token = token
	return s.price, s.updatedAt, s.err
}

func TestChainlinkQuoteProviderLatest(t *testing.T) {
	token := common.HexToAddress("0x1000000000000000000000000000000000000001")
	price, ok := new(big.Int).SetString("1915730000000000000000", 10)
	if !ok {
		t.Fatal("invalid test price")
	}
	now := time.Date(2026, 7, 30, 20, 0, 0, 0, time.UTC)
	oracle := &stubChainlinkOracle{price: price, updatedAt: now}
	provider := &ChainlinkQuoteProvider{
		oracle: oracle,
		tokens: map[string]common.Address{"WETH": token},
	}

	quote, err := provider.Latest(context.Background(), "weth")
	if err != nil {
		t.Fatal(err)
	}
	if oracle.token != token {
		t.Fatalf("oracle token = %s", oracle.token)
	}
	if quote.Symbol != "WETH" || quote.Currency != "USD" ||
		quote.Price != "1915.73" || quote.Source != "chainlink-oracle" ||
		!quote.UpdatedAt.Equal(now) {
		t.Fatalf("unexpected quote: %#v", quote)
	}
}

func TestChainlinkQuoteProviderRejectsUnknownSymbolAndOracleFailure(t *testing.T) {
	token := common.HexToAddress("0x1000000000000000000000000000000000000001")
	provider := &ChainlinkQuoteProvider{
		oracle: &stubChainlinkOracle{price: big.NewInt(1), updatedAt: time.Now()},
		tokens: map[string]common.Address{"USDC": token},
	}
	if _, err := provider.Latest(context.Background(), "WETH"); err == nil ||
		!strings.Contains(err.Error(), "not configured") {
		t.Fatalf("unknown-symbol error = %v", err)
	}

	provider.oracle = &stubChainlinkOracle{err: errors.New("stale feed")}
	if _, err := provider.Latest(context.Background(), "USDC"); err == nil ||
		!strings.Contains(err.Error(), "stale feed") {
		t.Fatalf("oracle error = %v", err)
	}

	provider.oracle = &stubChainlinkOracle{price: big.NewInt(0), updatedAt: time.Now()}
	if _, err := provider.Latest(context.Background(), "USDC"); err == nil ||
		!strings.Contains(err.Error(), "non-positive") {
		t.Fatalf("non-positive error = %v", err)
	}
}

func TestParsePriceTokenAddresses(t *testing.T) {
	tokens, err := parsePriceTokenAddresses(`{
		"usdc":"0x1000000000000000000000000000000000000001",
		"WETH":"0x2000000000000000000000000000000000000002"
	}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 ||
		tokens["USDC"] != common.HexToAddress("0x1000000000000000000000000000000000000001") ||
		tokens["WETH"] != common.HexToAddress("0x2000000000000000000000000000000000000002") {
		t.Fatalf("unexpected tokens: %#v", tokens)
	}

	for _, value := range []string{
		``,
		`[]`,
		`{}`,
		`{"WETH":"invalid"}`,
		`{"":"0x1000000000000000000000000000000000000001"}`,
		`{"weth":"0x1000000000000000000000000000000000000001","WETH":"0x2000000000000000000000000000000000000002"}`,
	} {
		if _, err := parsePriceTokenAddresses(value); err == nil {
			t.Errorf("expected %q to fail", value)
		}
	}
}

func TestFormatChainlinkPrice(t *testing.T) {
	for _, test := range []struct {
		value string
		want  string
	}{
		{value: "1000000000000000000", want: "1"},
		{value: "999758670000000000", want: "0.99975867"},
		{value: "1915730000000000000000", want: "1915.73"},
		{value: "1", want: "0.000000000000000001"},
	} {
		value, ok := new(big.Int).SetString(test.value, 10)
		if !ok {
			t.Fatal("invalid test value")
		}
		if got := formatChainlinkPrice(value); got != test.want {
			t.Errorf("formatChainlinkPrice(%s) = %q, want %q", test.value, got, test.want)
		}
	}
}

func TestConfiguredProviderAllowsLocalOutsideProduction(t *testing.T) {
	provider, err := NewConfiguredQuoteProvider(
		context.Background(), "local", "local", "", "", "",
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := provider.(*LocalQuoteProvider); !ok {
		t.Fatalf("expected LocalQuoteProvider, got %T", provider)
	}
}

func TestConfiguredProviderRejectsLocalInProduction(t *testing.T) {
	if _, err := NewConfiguredQuoteProvider(
		context.Background(), "production", "local", "", "", "",
	); err == nil {
		t.Fatal("expected production local provider to be rejected")
	}
}
