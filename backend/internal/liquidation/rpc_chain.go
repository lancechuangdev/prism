package liquidation

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lancechuangdev/prism/backend/internal/contracts"
)

type RPCChain struct {
	client      *ethclient.Client
	pool        *contracts.PrismPool
	auth        *bind.TransactOpts
	slippageBPS uint64
}

func NewRPCChain(ctx context.Context, rpcURL, poolAddress, expectedChainID, privateKey string, slippageBPS uint64) (*RPCChain, error) {
	if !common.IsHexAddress(poolAddress) || common.HexToAddress(poolAddress) == (common.Address{}) {
		return nil, fmt.Errorf("invalid PrismPool address %q", poolAddress)
	}
	if slippageBPS > 10_000 {
		return nil, fmt.Errorf("liquidation slippage BPS must not exceed 10000")
	}
	key, err := parsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	rpcClient, err := rpc.DialOptions(ctx, rpcURL, rpc.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}))
	if err != nil {
		return nil, fmt.Errorf("connect liquidation RPC: %w", err)
	}
	client := ethclient.NewClient(rpcClient)
	keepClient := false
	defer func() {
		if !keepClient {
			client.Close()
		}
	}()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("read liquidation RPC chain ID: %w", err)
	}
	if chainID.String() != strings.TrimSpace(expectedChainID) {
		return nil, fmt.Errorf("liquidation RPC chain ID %s does not match configured chain ID %s", chainID, expectedChainID)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, fmt.Errorf("create liquidation signer: %w", err)
	}
	pool, err := contracts.NewPrismPool(common.HexToAddress(poolAddress), client)
	if err != nil {
		return nil, fmt.Errorf("bind PrismPool liquidator: %w", err)
	}
	authorized, err := pool.Liquidator(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("read authorized PrismPool liquidator: %w", err)
	}
	if authorized != auth.From {
		return nil, fmt.Errorf("scheduler signer %s is not the authorized PrismPool liquidator %s", auth.From.Hex(), authorized.Hex())
	}

	keepClient = true
	return &RPCChain{client: client, pool: pool, auth: auth, slippageBPS: slippageBPS}, nil
}

func parsePrivateKey(value string) (*ecdsa.PrivateKey, error) {
	normalized := strings.TrimPrefix(strings.TrimSpace(value), "0x")
	if normalized == "" {
		return nil, fmt.Errorf("liquidation private key is required")
	}
	key, err := crypto.HexToECDSA(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid liquidation private key: %w", err)
	}
	return key, nil
}

func (r *RPCChain) PoolCount(ctx context.Context) (int64, error) {
	count, err := r.pool.PoolCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return 0, err
	}
	if !count.IsInt64() {
		return 0, fmt.Errorf("pool count exceeds int64")
	}
	return count.Int64(), nil
}

func (r *RPCChain) IsUndercollateralized(ctx context.Context, poolID int64) (bool, error) {
	return r.pool.IsUndercollateralized(&bind.CallOpts{Context: ctx}, big.NewInt(poolID))
}

func (r *RPCChain) MaxCollateral(ctx context.Context, poolID int64) (*big.Int, error) {
	base, err := r.pool.GetPool(&bind.CallOpts{Context: ctx}, big.NewInt(poolID))
	if err != nil {
		return nil, err
	}
	data, err := r.pool.GetPoolData(&bind.CallOpts{Context: ctx}, big.NewInt(poolID))
	if err != nil {
		return nil, err
	}
	repayment, err := r.pool.GetRequiredRepayment(&bind.CallOpts{Context: ctx}, big.NewInt(poolID))
	if err != nil {
		return nil, err
	}
	dexAddress, err := r.pool.DexSwap(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, err
	}
	dex, err := contracts.NewDexSwapLikeCaller(dexAddress, r.client)
	if err != nil {
		return nil, err
	}
	var output []interface{}
	raw := contracts.DexSwapLikeCallerRaw{Contract: dex}
	err = raw.Call(
		&bind.CallOpts{Context: ctx, From: r.auth.From},
		&output,
		"getAmountIn",
		base.CollateralToken,
		base.LendToken,
		repayment,
	)
	if err != nil {
		return nil, fmt.Errorf("quote liquidation collateral: %w", err)
	}
	if len(output) != 1 {
		return nil, fmt.Errorf("quote liquidation collateral returned %d values", len(output))
	}
	quoted := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	maxCollateral := new(big.Int).Mul(quoted, new(big.Int).SetUint64(10_000+r.slippageBPS))
	maxCollateral.Add(maxCollateral, big.NewInt(9_999))
	maxCollateral.Quo(maxCollateral, big.NewInt(10_000))
	if maxCollateral.Cmp(data.SettleAmountBorrow) > 0 {
		maxCollateral.Set(data.SettleAmountBorrow)
	}
	return maxCollateral, nil
}

func (r *RPCChain) Liquidate(ctx context.Context, poolID int64, maxCollateral *big.Int) (string, error) {
	opts := *r.auth
	opts.Context = ctx
	tx, err := r.pool.Liquidate(&opts, big.NewInt(poolID), maxCollateral)
	if err != nil {
		return "", fmt.Errorf("simulate and submit transaction: %w", err)
	}
	receipt, err := bind.WaitMined(ctx, r.client, tx)
	if err != nil {
		return tx.Hash().Hex(), fmt.Errorf("wait for transaction %s: %w", tx.Hash().Hex(), err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return tx.Hash().Hex(), fmt.Errorf("transaction %s reverted", tx.Hash().Hex())
	}
	return tx.Hash().Hex(), nil
}

func (r *RPCChain) Close() {
	r.client.Close()
}
