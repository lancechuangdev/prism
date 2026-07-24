package chain

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

const (
	testPoolAddress   = "0x1000000000000000000000000000000000000001"
	testLendAddress   = "0x2000000000000000000000000000000000000002"
	testBorrowAddress = "0x3000000000000000000000000000000000000003"
	testSPAddress     = "0x4000000000000000000000000000000000000004"
	testJPAddress     = "0x5000000000000000000000000000000000000005"
)

type testEthAPI struct {
	poolABI  abi.ABI
	tokenABI abi.ABI
}

type testCallArgs struct {
	To    *common.Address `json:"to"`
	Input *hexutil.Bytes  `json:"input"`
}

func (a *testEthAPI) ChainId() *hexutil.Big {
	value := big.NewInt(31337)
	return (*hexutil.Big)(value)
}

func (a *testEthAPI) Call(_ context.Context, args testCallArgs, _ string) (hexutil.Bytes, error) {
	if args.To == nil || args.Input == nil || len(*args.Input) < 4 {
		return nil, fmt.Errorf("invalid eth_call")
	}
	data := []byte(*args.Input)
	if method, err := a.poolABI.MethodById(data[:4]); err == nil {
		switch method.Name {
		case "poolCount":
			return method.Outputs.Pack(big.NewInt(1))
		case "getPool":
			return method.Outputs.Pack(contracts.PrismPoolPoolBaseInfo{
				SettleTime: big.NewInt(100), MaturityTime: big.NewInt(200),
				InterestRate: big.NewInt(1_000_000), MaxLendSupply: big.NewInt(1_000),
				TotalLendDeposited: big.NewInt(500), TotalCollateralDeposited: big.NewInt(2),
				CollateralizationRatio: big.NewInt(200_000_000),
				LendToken:              common.HexToAddress(testLendAddress), CollateralToken: common.HexToAddress(testBorrowAddress),
				State: 1, LenderPositionToken: common.HexToAddress(testSPAddress),
				BorrowerPositionToken: common.HexToAddress(testJPAddress), LiquidateRate: big.NewInt(20_000_000),
			})
		case "getPoolData":
			return method.Outputs.Pack(contracts.PrismPoolPoolDataInfo{
				SettleAmountLend: big.NewInt(400), SettleAmountBorrow: big.NewInt(1),
				FinishAmountLend: big.NewInt(0), FinishAmountBorrow: big.NewInt(0),
				LiquidationAmountLend: big.NewInt(0), LiquidationAmountBorrow: big.NewInt(0),
			})
		}
	}
	if method, err := a.tokenABI.MethodById(data[:4]); err == nil {
		switch method.Name {
		case "symbol":
			if *args.To == common.HexToAddress(testLendAddress) {
				return method.Outputs.Pack("pUSD")
			}
			return method.Outputs.Pack("pETH")
		case "decimals":
			return method.Outputs.Pack(uint8(18))
		}
	}
	return nil, fmt.Errorf("unknown contract call")
}

func TestRPCReaderReadsContractSnapshots(t *testing.T) {
	ctx := context.Background()
	server := rpc.NewServer()
	poolABI, err := contracts.PrismPoolMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	tokenABI, err := contracts.ERC20TokenMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	if err := server.RegisterName("eth", &testEthAPI{poolABI: *poolABI, tokenABI: *tokenABI}); err != nil {
		t.Fatal(err)
	}

	reader, err := newInProcessRPCReader(ctx, rpc.DialInProc(server), testPoolAddress)
	if err != nil {
		t.Fatalf("create RPC reader: %v", err)
	}
	defer reader.Close()

	count, err := reader.PoolLength(ctx, "31337")
	if err != nil || count != 1 {
		t.Fatalf("pool count = %d, err = %v", count, err)
	}
	base, err := reader.PoolBaseInfo(ctx, "31337", 0)
	if err != nil {
		t.Fatalf("read pool base: %v", err)
	}
	if base.SettleTime != "100" ||
		base.MaturityTime != "200" ||
		base.MaxLendSupply != "1000" ||
		base.TotalLendDeposited != "500" ||
		base.TotalCollateralDeposited != "2" ||
		base.CollateralizationRatio != "200000000" ||
		base.LendTokenAddress != common.HexToAddress(testLendAddress).Hex() ||
		base.CollateralTokenAddress != common.HexToAddress(testBorrowAddress).Hex() ||
		base.State != "1" ||
		base.LenderPositionToken != common.HexToAddress(testSPAddress).Hex() ||
		base.BorrowerPositionToken != common.HexToAddress(testJPAddress).Hex() ||
		base.LiquidateRate != "20000000" {
		t.Fatalf("unexpected pool base: %+v", base)
	}
	data, err := reader.PoolDataInfo(ctx, "31337", 0)
	if err != nil || data.SettleAmountLend != "400" {
		t.Fatalf("unexpected pool data: %+v, err = %v", data, err)
	}
	token, err := reader.TokenInfo(ctx, "31337", testLendAddress)
	if err != nil {
		t.Fatalf("read token: %v", err)
	}
	if token.Symbol != "pUSD" || token.Decimals != 18 {
		t.Fatalf("unexpected token: %+v", token)
	}
}

func TestRPCReaderRejectsWrongChain(t *testing.T) {
	reader := &RPCReader{chainID: "31337"}
	if _, err := reader.PoolLength(context.Background(), "97"); err == nil {
		t.Fatal("expected chain ID mismatch")
	}
}
