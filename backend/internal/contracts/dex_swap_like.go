// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DexSwapLikeMetaData contains all meta data concerning the DexSwapLike contract.
var DexSwapLikeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"getAmountIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DexSwapLikeABI is the input ABI used to generate the binding from.
// Deprecated: Use DexSwapLikeMetaData.ABI instead.
var DexSwapLikeABI = DexSwapLikeMetaData.ABI

// DexSwapLike is an auto generated Go binding around an Ethereum contract.
type DexSwapLike struct {
	DexSwapLikeCaller     // Read-only binding to the contract
	DexSwapLikeTransactor // Write-only binding to the contract
	DexSwapLikeFilterer   // Log filterer for contract events
}

// DexSwapLikeCaller is an auto generated read-only Go binding around an Ethereum contract.
type DexSwapLikeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexSwapLikeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DexSwapLikeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexSwapLikeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DexSwapLikeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexSwapLikeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DexSwapLikeSession struct {
	Contract     *DexSwapLike      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DexSwapLikeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DexSwapLikeCallerSession struct {
	Contract *DexSwapLikeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DexSwapLikeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DexSwapLikeTransactorSession struct {
	Contract     *DexSwapLikeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DexSwapLikeRaw is an auto generated low-level Go binding around an Ethereum contract.
type DexSwapLikeRaw struct {
	Contract *DexSwapLike // Generic contract binding to access the raw methods on
}

// DexSwapLikeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DexSwapLikeCallerRaw struct {
	Contract *DexSwapLikeCaller // Generic read-only contract binding to access the raw methods on
}

// DexSwapLikeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DexSwapLikeTransactorRaw struct {
	Contract *DexSwapLikeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDexSwapLike creates a new instance of DexSwapLike, bound to a specific deployed contract.
func NewDexSwapLike(address common.Address, backend bind.ContractBackend) (*DexSwapLike, error) {
	contract, err := bindDexSwapLike(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DexSwapLike{DexSwapLikeCaller: DexSwapLikeCaller{contract: contract}, DexSwapLikeTransactor: DexSwapLikeTransactor{contract: contract}, DexSwapLikeFilterer: DexSwapLikeFilterer{contract: contract}}, nil
}

// NewDexSwapLikeCaller creates a new read-only instance of DexSwapLike, bound to a specific deployed contract.
func NewDexSwapLikeCaller(address common.Address, caller bind.ContractCaller) (*DexSwapLikeCaller, error) {
	contract, err := bindDexSwapLike(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DexSwapLikeCaller{contract: contract}, nil
}

// NewDexSwapLikeTransactor creates a new write-only instance of DexSwapLike, bound to a specific deployed contract.
func NewDexSwapLikeTransactor(address common.Address, transactor bind.ContractTransactor) (*DexSwapLikeTransactor, error) {
	contract, err := bindDexSwapLike(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DexSwapLikeTransactor{contract: contract}, nil
}

// NewDexSwapLikeFilterer creates a new log filterer instance of DexSwapLike, bound to a specific deployed contract.
func NewDexSwapLikeFilterer(address common.Address, filterer bind.ContractFilterer) (*DexSwapLikeFilterer, error) {
	contract, err := bindDexSwapLike(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DexSwapLikeFilterer{contract: contract}, nil
}

// bindDexSwapLike binds a generic wrapper to an already deployed contract.
func bindDexSwapLike(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DexSwapLikeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DexSwapLike *DexSwapLikeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DexSwapLike.Contract.DexSwapLikeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DexSwapLike *DexSwapLikeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DexSwapLike.Contract.DexSwapLikeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DexSwapLike *DexSwapLikeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DexSwapLike.Contract.DexSwapLikeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DexSwapLike *DexSwapLikeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DexSwapLike.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DexSwapLike *DexSwapLikeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DexSwapLike.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DexSwapLike *DexSwapLikeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DexSwapLike.Contract.contract.Transact(opts, method, params...)
}

// GetAmountIn is a paid mutator transaction binding the contract method 0x53b609b5.
//
// Solidity: function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) returns(uint256 amountIn)
func (_DexSwapLike *DexSwapLikeTransactor) GetAmountIn(opts *bind.TransactOpts, tokenIn common.Address, tokenOut common.Address, amountOut *big.Int) (*types.Transaction, error) {
	return _DexSwapLike.contract.Transact(opts, "getAmountIn", tokenIn, tokenOut, amountOut)
}

// GetAmountIn is a paid mutator transaction binding the contract method 0x53b609b5.
//
// Solidity: function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) returns(uint256 amountIn)
func (_DexSwapLike *DexSwapLikeSession) GetAmountIn(tokenIn common.Address, tokenOut common.Address, amountOut *big.Int) (*types.Transaction, error) {
	return _DexSwapLike.Contract.GetAmountIn(&_DexSwapLike.TransactOpts, tokenIn, tokenOut, amountOut)
}

// GetAmountIn is a paid mutator transaction binding the contract method 0x53b609b5.
//
// Solidity: function getAmountIn(address tokenIn, address tokenOut, uint256 amountOut) returns(uint256 amountIn)
func (_DexSwapLike *DexSwapLikeTransactorSession) GetAmountIn(tokenIn common.Address, tokenOut common.Address, amountOut *big.Int) (*types.Transaction, error) {
	return _DexSwapLike.Contract.GetAmountIn(&_DexSwapLike.TransactOpts, tokenIn, tokenOut, amountOut)
}
