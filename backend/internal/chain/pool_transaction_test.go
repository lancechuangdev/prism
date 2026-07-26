package chain

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

func TestPrepareCreatePool(t *testing.T) {
	builder, err := NewPoolTransactionBuilder("31337", "0x1000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := builder.PrepareCreatePool(context.Background(), validCreatePoolParams())
	if err != nil {
		t.Fatalf("prepare transaction: %v", err)
	}
	if tx.To != "0x1000000000000000000000000000000000000001" ||
		tx.ChainID != "31337" ||
		tx.Value != "0x0" ||
		!strings.HasPrefix(tx.Data, "0x") ||
		len(tx.Data) <= 10 {
		t.Fatalf("unexpected transaction: %+v", tx)
	}
}

func TestPrepareCreatePoolRejectsInvalidValues(t *testing.T) {
	builder, err := NewPoolTransactionBuilder("31337", "0x1000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	params := validCreatePoolParams()
	params.SettleTime = "not-a-number"
	_, err = builder.PrepareCreatePool(context.Background(), params)
	if !errors.Is(err, ErrInvalidCreatePool) {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}

func TestPrepareSettlePool(t *testing.T) {
	builder, err := NewPoolTransactionBuilder("31337", "0x1000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := builder.PrepareSettlePool(context.Background(), "42")
	if err != nil {
		t.Fatalf("prepare transaction: %v", err)
	}
	if tx.To != "0x1000000000000000000000000000000000000001" ||
		tx.ChainID != "31337" ||
		tx.Value != "0x0" {
		t.Fatalf("unexpected transaction: %+v", tx)
	}

	contractABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	data, err := hexutil.Decode(tx.Data)
	if err != nil {
		t.Fatal(err)
	}
	method, err := contractABI.MethodById(data[:4])
	if err != nil {
		t.Fatal(err)
	}
	values, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		t.Fatal(err)
	}
	if method.Name != "settle" || len(values) != 1 || values[0].(*big.Int).Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("unexpected settle calldata: method=%s values=%v", method.Name, values)
	}
}

func TestPrepareSettlePoolRejectsInvalidPoolID(t *testing.T) {
	builder, err := NewPoolTransactionBuilder("31337", "0x1000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	for _, poolID := range []string{"", "-1", "1.5", new(big.Int).Lsh(big.NewInt(1), 256).String()} {
		_, err = builder.PrepareSettlePool(context.Background(), poolID)
		if !errors.Is(err, ErrInvalidSettlePool) {
			t.Fatalf("pool ID %q: expected invalid request error, got %v", poolID, err)
		}
	}
}

func validCreatePoolParams() CreatePoolParams {
	return CreatePoolParams{
		SettleTime: "2000000000", MaturityTime: "2000600000",
		InterestRate: "1000000", MaxLendSupply: "1000000000000000000000",
		CollateralizationRatio: "200000000",
		LendToken:              "0x1000000000000000000000000000000000000001",
		CollateralToken:        "0x2000000000000000000000000000000000000002",
		LenderPositionToken:    "0x3000000000000000000000000000000000000003",
		BorrowerPositionToken:  "0x4000000000000000000000000000000000000004",
		LiquidateRate:          "20000000",
	}
}
