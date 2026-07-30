package price

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lancechuangdev/prism/backend/internal/contracts"
	"github.com/lancechuangdev/prism/backend/internal/resilience"
)

const chainlinkPriceDecimals = 18

var chainlinkReadRetryPolicy = resilience.Policy{
	Attempts:       3,
	AttemptTimeout: 5 * time.Second,
	InitialBackoff: 100 * time.Millisecond,
	MaxBackoff:     time.Second,
}

type chainlinkOracleReader interface {
	Latest(ctx context.Context, token common.Address) (*big.Int, time.Time, error)
}

// blockchain-specific low-level Ethereum adapter
type rpcChainlinkOracle struct {
	client *ethclient.Client
	oracle *contracts.ChainlinkOracleCaller
}

func (r *rpcChainlinkOracle) Latest(ctx context.Context, token common.Address) (*big.Int, time.Time, error) {
	opts := &bind.CallOpts{Context: ctx}
	price, err := r.oracle.GetPrice(opts, token)
	if err != nil {
		return nil, time.Time{}, err
	}
	config, err := r.oracle.Feeds(opts, token)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("read configured feed: %w", err)
	}
	if config.Feed == (common.Address{}) {
		return nil, time.Time{}, fmt.Errorf("feed is not configured")
	}
	aggregator, err := contracts.NewChainlinkAggregatorV3Caller(config.Feed, r.client)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("bind configured feed: %w", err)
	}
	round, err := aggregator.LatestRoundData(opts)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("read configured feed round: %w", err)
	}
	if round.UpdatedAt == nil || !round.UpdatedAt.IsInt64() || round.UpdatedAt.Sign() <= 0 {
		return nil, time.Time{}, fmt.Errorf("configured feed returned an invalid timestamp")
	}
	return price, time.Unix(round.UpdatedAt.Int64(), 0).UTC(), nil
}

type ChainlinkQuoteProvider struct {
	client *ethclient.Client
	oracle chainlinkOracleReader
	tokens map[string]common.Address
}

func NewChainlinkQuoteProvider(
	ctx context.Context,
	rpcURL string,
	oracleAddress string,
	tokenAddressesJSON string,
) (*ChainlinkQuoteProvider, error) {
	if !common.IsHexAddress(oracleAddress) || common.HexToAddress(oracleAddress) == (common.Address{}) {
		return nil, fmt.Errorf("invalid ChainlinkOracle address %q", oracleAddress)
	}
	tokens, err := parsePriceTokenAddresses(tokenAddressesJSON)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect to ChainlinkOracle RPC: %w", err)
	}
	keepClient := false
	defer func() {
		if !keepClient {
			client.Close()
		}
	}()

	address := common.HexToAddress(oracleAddress)
	code, err := client.CodeAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("read ChainlinkOracle bytecode: %w", err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract deployed at ChainlinkOracle address %s", address.Hex())
	}
	oracle, err := contracts.NewChainlinkOracleCaller(address, client)
	if err != nil {
		return nil, fmt.Errorf("bind ChainlinkOracle: %w", err)
	}

	keepClient = true
	return &ChainlinkQuoteProvider{
		client: client,
		oracle: &rpcChainlinkOracle{client: client, oracle: oracle},
		tokens: tokens,
	}, nil
}

func parsePriceTokenAddresses(value string) (map[string]common.Address, error) {
	var configured map[string]string
	if err := json.Unmarshal([]byte(value), &configured); err != nil {
		return nil, fmt.Errorf("PRISM_PRICE_TOKEN_ADDRESSES must be a JSON object: %w", err)
	}
	if len(configured) == 0 {
		return nil, fmt.Errorf("PRISM_PRICE_TOKEN_ADDRESSES must not be empty")
	}

	tokens := make(map[string]common.Address, len(configured))
	for symbol, addressValue := range configured {
		normalizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
		if normalizedSymbol == "" {
			return nil, fmt.Errorf("PRISM_PRICE_TOKEN_ADDRESSES contains an empty symbol")
		}
		if !common.IsHexAddress(addressValue) || common.HexToAddress(addressValue) == (common.Address{}) {
			return nil, fmt.Errorf("invalid token address for symbol %s", normalizedSymbol)
		}
		if _, exists := tokens[normalizedSymbol]; exists {
			return nil, fmt.Errorf("duplicate token symbol %s", normalizedSymbol)
		}
		tokens[normalizedSymbol] = common.HexToAddress(addressValue)
	}
	return tokens, nil
}

func (p *ChainlinkQuoteProvider) Latest(ctx context.Context, symbol string) (Quote, error) {
	normalizedSymbol := strings.ToUpper(strings.TrimSpace(symbol))
	token, ok := p.tokens[normalizedSymbol]
	if !ok {
		return Quote{}, fmt.Errorf("price token address for %s is not configured", normalizedSymbol)
	}

	type observation struct {
		value     *big.Int
		updatedAt time.Time
	}
	result, err := resilience.Value(ctx, chainlinkReadRetryPolicy, func(attemptCtx context.Context) (observation, error) {
		value, updatedAt, err := p.oracle.Latest(attemptCtx, token)
		return observation{value: value, updatedAt: updatedAt}, err
	})
	if err != nil {
		return Quote{}, fmt.Errorf("read ChainlinkOracle price for %s: %w", normalizedSymbol, err)
	}
	if result.value == nil || result.value.Sign() <= 0 {
		return Quote{}, fmt.Errorf("ChainlinkOracle returned a non-positive price for %s", normalizedSymbol)
	}
	if result.updatedAt.IsZero() {
		return Quote{}, fmt.Errorf("ChainlinkOracle returned an invalid update time for %s", normalizedSymbol)
	}

	return Quote{
		Symbol:    normalizedSymbol,
		Currency:  "USD",
		Price:     formatChainlinkPrice(result.value),
		Source:    "chainlink-oracle",
		UpdatedAt: result.updatedAt.UTC(),
	}, nil
}

func formatChainlinkPrice(value *big.Int) string {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(chainlinkPriceDecimals), nil)
	whole, fraction := new(big.Int), new(big.Int)
	whole.QuoRem(value, scale, fraction)
	if fraction.Sign() == 0 {
		return whole.String()
	}
	fractionText := fmt.Sprintf("%018s", fraction.String())
	return whole.String() + "." + strings.TrimRight(fractionText, "0")
}

func (p *ChainlinkQuoteProvider) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
