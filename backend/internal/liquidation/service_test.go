package liquidation

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"testing"
)

type fakeChain struct {
	under          map[int64]bool
	collateral     map[int64]*big.Int
	liquidated     []int64
	liquidationErr error
}

func (f *fakeChain) PoolCount(context.Context) (int64, error) {
	return int64(len(f.under)), nil
}

func (f *fakeChain) IsUndercollateralized(_ context.Context, poolID int64) (bool, error) {
	return f.under[poolID], nil
}

func (f *fakeChain) MaxCollateral(_ context.Context, poolID int64) (*big.Int, error) {
	return f.collateral[poolID], nil
}

func (f *fakeChain) Liquidate(_ context.Context, poolID int64, _ *big.Int) (string, error) {
	if f.liquidationErr != nil {
		return "", f.liquidationErr
	}
	f.liquidated = append(f.liquidated, poolID)
	return "0xtest", nil
}

func TestCheckAndLiquidateOnlyTriggersUnsafePools(t *testing.T) {
	chain := &fakeChain{
		under:      map[int64]bool{0: false, 1: true},
		collateral: map[int64]*big.Int{1: big.NewInt(100)},
	}
	service := NewService(chain, slog.New(slog.NewTextHandler(io.Discard, nil)))

	if err := service.CheckAndLiquidate(context.Background()); err != nil {
		t.Fatalf("check and liquidate: %v", err)
	}
	if len(chain.liquidated) != 1 || chain.liquidated[0] != 1 {
		t.Fatalf("unexpected liquidations: %v", chain.liquidated)
	}
}

func TestCheckAndLiquidateReportsSubmissionFailure(t *testing.T) {
	chain := &fakeChain{
		under:          map[int64]bool{0: true},
		collateral:     map[int64]*big.Int{0: big.NewInt(100)},
		liquidationErr: errors.New("submission failed"),
	}
	service := NewService(chain, slog.New(slog.NewTextHandler(io.Discard, nil)))

	if err := service.CheckAndLiquidate(context.Background()); err == nil {
		t.Fatal("expected liquidation failure")
	}
}

func TestParsePrivateKeyAcceptsOptionalHexPrefix(t *testing.T) {
	const key = "4f3edf983ac63ad7c6a17ad38b2b4f8f6f0f6f2d14b8f4c9f5f5ecf7a2f5a111"
	withoutPrefix, err := parsePrivateKey(key)
	if err != nil {
		t.Fatalf("parse key without prefix: %v", err)
	}
	withPrefix, err := parsePrivateKey("0x" + key)
	if err != nil {
		t.Fatalf("parse key with prefix: %v", err)
	}
	if withoutPrefix.D.Cmp(withPrefix.D) != 0 {
		t.Fatal("expected both forms to parse to the same key")
	}
}
