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

// ThresholdMultiSigMetaData contains all meta data concerning the ThresholdMultiSig contract.
var ThresholdMultiSigMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"owners_\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"threshold_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerReplaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"ThresholdChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"approvalCount\",\"type\":\"uint256\"}],\"name\":\"TransactionApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"TransactionExecuted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"addOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"approvalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"approveTransaction\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"changeThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"configurationVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"executeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"executed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"getTransactionHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasApproved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ownerCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"removeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"replaceOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactionConfigurationVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ThresholdMultiSigABI is the input ABI used to generate the binding from.
// Deprecated: Use ThresholdMultiSigMetaData.ABI instead.
var ThresholdMultiSigABI = ThresholdMultiSigMetaData.ABI

// ThresholdMultiSig is an auto generated Go binding around an Ethereum contract.
type ThresholdMultiSig struct {
	ThresholdMultiSigCaller     // Read-only binding to the contract
	ThresholdMultiSigTransactor // Write-only binding to the contract
	ThresholdMultiSigFilterer   // Log filterer for contract events
}

// ThresholdMultiSigCaller is an auto generated read-only Go binding around an Ethereum contract.
type ThresholdMultiSigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThresholdMultiSigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ThresholdMultiSigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThresholdMultiSigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ThresholdMultiSigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ThresholdMultiSigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ThresholdMultiSigSession struct {
	Contract     *ThresholdMultiSig // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ThresholdMultiSigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ThresholdMultiSigCallerSession struct {
	Contract *ThresholdMultiSigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ThresholdMultiSigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ThresholdMultiSigTransactorSession struct {
	Contract     *ThresholdMultiSigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ThresholdMultiSigRaw is an auto generated low-level Go binding around an Ethereum contract.
type ThresholdMultiSigRaw struct {
	Contract *ThresholdMultiSig // Generic contract binding to access the raw methods on
}

// ThresholdMultiSigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ThresholdMultiSigCallerRaw struct {
	Contract *ThresholdMultiSigCaller // Generic read-only contract binding to access the raw methods on
}

// ThresholdMultiSigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ThresholdMultiSigTransactorRaw struct {
	Contract *ThresholdMultiSigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewThresholdMultiSig creates a new instance of ThresholdMultiSig, bound to a specific deployed contract.
func NewThresholdMultiSig(address common.Address, backend bind.ContractBackend) (*ThresholdMultiSig, error) {
	contract, err := bindThresholdMultiSig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSig{ThresholdMultiSigCaller: ThresholdMultiSigCaller{contract: contract}, ThresholdMultiSigTransactor: ThresholdMultiSigTransactor{contract: contract}, ThresholdMultiSigFilterer: ThresholdMultiSigFilterer{contract: contract}}, nil
}

// NewThresholdMultiSigCaller creates a new read-only instance of ThresholdMultiSig, bound to a specific deployed contract.
func NewThresholdMultiSigCaller(address common.Address, caller bind.ContractCaller) (*ThresholdMultiSigCaller, error) {
	contract, err := bindThresholdMultiSig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigCaller{contract: contract}, nil
}

// NewThresholdMultiSigTransactor creates a new write-only instance of ThresholdMultiSig, bound to a specific deployed contract.
func NewThresholdMultiSigTransactor(address common.Address, transactor bind.ContractTransactor) (*ThresholdMultiSigTransactor, error) {
	contract, err := bindThresholdMultiSig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigTransactor{contract: contract}, nil
}

// NewThresholdMultiSigFilterer creates a new log filterer instance of ThresholdMultiSig, bound to a specific deployed contract.
func NewThresholdMultiSigFilterer(address common.Address, filterer bind.ContractFilterer) (*ThresholdMultiSigFilterer, error) {
	contract, err := bindThresholdMultiSig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigFilterer{contract: contract}, nil
}

// bindThresholdMultiSig binds a generic wrapper to an already deployed contract.
func bindThresholdMultiSig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ThresholdMultiSigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ThresholdMultiSig *ThresholdMultiSigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ThresholdMultiSig.Contract.ThresholdMultiSigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ThresholdMultiSig *ThresholdMultiSigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ThresholdMultiSigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ThresholdMultiSig *ThresholdMultiSigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ThresholdMultiSigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ThresholdMultiSig *ThresholdMultiSigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ThresholdMultiSig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ThresholdMultiSig *ThresholdMultiSigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ThresholdMultiSig *ThresholdMultiSigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.contract.Transact(opts, method, params...)
}

// ApprovalCount is a free data retrieval call binding the contract method 0x8d1ba346.
//
// Solidity: function approvalCount(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) ApprovalCount(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "approvalCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApprovalCount is a free data retrieval call binding the contract method 0x8d1ba346.
//
// Solidity: function approvalCount(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigSession) ApprovalCount(arg0 [32]byte) (*big.Int, error) {
	return _ThresholdMultiSig.Contract.ApprovalCount(&_ThresholdMultiSig.CallOpts, arg0)
}

// ApprovalCount is a free data retrieval call binding the contract method 0x8d1ba346.
//
// Solidity: function approvalCount(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) ApprovalCount(arg0 [32]byte) (*big.Int, error) {
	return _ThresholdMultiSig.Contract.ApprovalCount(&_ThresholdMultiSig.CallOpts, arg0)
}

// ConfigurationVersion is a free data retrieval call binding the contract method 0x2dc9afb4.
//
// Solidity: function configurationVersion() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) ConfigurationVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "configurationVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConfigurationVersion is a free data retrieval call binding the contract method 0x2dc9afb4.
//
// Solidity: function configurationVersion() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigSession) ConfigurationVersion() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.ConfigurationVersion(&_ThresholdMultiSig.CallOpts)
}

// ConfigurationVersion is a free data retrieval call binding the contract method 0x2dc9afb4.
//
// Solidity: function configurationVersion() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) ConfigurationVersion() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.ConfigurationVersion(&_ThresholdMultiSig.CallOpts)
}

// Executed is a free data retrieval call binding the contract method 0xa9fcfb33.
//
// Solidity: function executed(bytes32 ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) Executed(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "executed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Executed is a free data retrieval call binding the contract method 0xa9fcfb33.
//
// Solidity: function executed(bytes32 ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigSession) Executed(arg0 [32]byte) (bool, error) {
	return _ThresholdMultiSig.Contract.Executed(&_ThresholdMultiSig.CallOpts, arg0)
}

// Executed is a free data retrieval call binding the contract method 0xa9fcfb33.
//
// Solidity: function executed(bytes32 ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) Executed(arg0 [32]byte) (bool, error) {
	return _ThresholdMultiSig.Contract.Executed(&_ThresholdMultiSig.CallOpts, arg0)
}

// GetOwner is a free data retrieval call binding the contract method 0xc41a360a.
//
// Solidity: function getOwner(uint256 index) view returns(address)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) GetOwner(opts *bind.CallOpts, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "getOwner", index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0xc41a360a.
//
// Solidity: function getOwner(uint256 index) view returns(address)
func (_ThresholdMultiSig *ThresholdMultiSigSession) GetOwner(index *big.Int) (common.Address, error) {
	return _ThresholdMultiSig.Contract.GetOwner(&_ThresholdMultiSig.CallOpts, index)
}

// GetOwner is a free data retrieval call binding the contract method 0xc41a360a.
//
// Solidity: function getOwner(uint256 index) view returns(address)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) GetOwner(index *big.Int) (common.Address, error) {
	return _ThresholdMultiSig.Contract.GetOwner(&_ThresholdMultiSig.CallOpts, index)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0xb98a34de.
//
// Solidity: function getTransactionHash(address target, uint256 value, bytes data, uint256 nonce) view returns(bytes32)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) GetTransactionHash(opts *bind.CallOpts, target common.Address, value *big.Int, data []byte, nonce *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "getTransactionHash", target, value, data, nonce)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTransactionHash is a free data retrieval call binding the contract method 0xb98a34de.
//
// Solidity: function getTransactionHash(address target, uint256 value, bytes data, uint256 nonce) view returns(bytes32)
func (_ThresholdMultiSig *ThresholdMultiSigSession) GetTransactionHash(target common.Address, value *big.Int, data []byte, nonce *big.Int) ([32]byte, error) {
	return _ThresholdMultiSig.Contract.GetTransactionHash(&_ThresholdMultiSig.CallOpts, target, value, data, nonce)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0xb98a34de.
//
// Solidity: function getTransactionHash(address target, uint256 value, bytes data, uint256 nonce) view returns(bytes32)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) GetTransactionHash(target common.Address, value *big.Int, data []byte, nonce *big.Int) ([32]byte, error) {
	return _ThresholdMultiSig.Contract.GetTransactionHash(&_ThresholdMultiSig.CallOpts, target, value, data, nonce)
}

// HasApproved is a free data retrieval call binding the contract method 0x23cdc3f1.
//
// Solidity: function hasApproved(bytes32 , address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) HasApproved(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "hasApproved", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasApproved is a free data retrieval call binding the contract method 0x23cdc3f1.
//
// Solidity: function hasApproved(bytes32 , address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigSession) HasApproved(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _ThresholdMultiSig.Contract.HasApproved(&_ThresholdMultiSig.CallOpts, arg0, arg1)
}

// HasApproved is a free data retrieval call binding the contract method 0x23cdc3f1.
//
// Solidity: function hasApproved(bytes32 , address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) HasApproved(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _ThresholdMultiSig.Contract.HasApproved(&_ThresholdMultiSig.CallOpts, arg0, arg1)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) IsOwner(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "isOwner", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigSession) IsOwner(arg0 common.Address) (bool, error) {
	return _ThresholdMultiSig.Contract.IsOwner(&_ThresholdMultiSig.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address ) view returns(bool)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) IsOwner(arg0 common.Address) (bool, error) {
	return _ThresholdMultiSig.Contract.IsOwner(&_ThresholdMultiSig.CallOpts, arg0)
}

// OwnerCount is a free data retrieval call binding the contract method 0x0db02622.
//
// Solidity: function ownerCount() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) OwnerCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "ownerCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnerCount is a free data retrieval call binding the contract method 0x0db02622.
//
// Solidity: function ownerCount() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigSession) OwnerCount() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.OwnerCount(&_ThresholdMultiSig.CallOpts)
}

// OwnerCount is a free data retrieval call binding the contract method 0x0db02622.
//
// Solidity: function ownerCount() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) OwnerCount() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.OwnerCount(&_ThresholdMultiSig.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigSession) Threshold() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.Threshold(&_ThresholdMultiSig.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) Threshold() (*big.Int, error) {
	return _ThresholdMultiSig.Contract.Threshold(&_ThresholdMultiSig.CallOpts)
}

// TransactionConfigurationVersion is a free data retrieval call binding the contract method 0xb67766ce.
//
// Solidity: function transactionConfigurationVersion(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCaller) TransactionConfigurationVersion(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ThresholdMultiSig.contract.Call(opts, &out, "transactionConfigurationVersion", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransactionConfigurationVersion is a free data retrieval call binding the contract method 0xb67766ce.
//
// Solidity: function transactionConfigurationVersion(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigSession) TransactionConfigurationVersion(arg0 [32]byte) (*big.Int, error) {
	return _ThresholdMultiSig.Contract.TransactionConfigurationVersion(&_ThresholdMultiSig.CallOpts, arg0)
}

// TransactionConfigurationVersion is a free data retrieval call binding the contract method 0xb67766ce.
//
// Solidity: function transactionConfigurationVersion(bytes32 ) view returns(uint256)
func (_ThresholdMultiSig *ThresholdMultiSigCallerSession) TransactionConfigurationVersion(arg0 [32]byte) (*big.Int, error) {
	return _ThresholdMultiSig.Contract.TransactionConfigurationVersion(&_ThresholdMultiSig.CallOpts, arg0)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) AddOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "addOwner", owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigSession) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.AddOwner(&_ThresholdMultiSig.TransactOpts, owner)
}

// AddOwner is a paid mutator transaction binding the contract method 0x7065cb48.
//
// Solidity: function addOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) AddOwner(owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.AddOwner(&_ThresholdMultiSig.TransactOpts, owner)
}

// ApproveTransaction is a paid mutator transaction binding the contract method 0x1c91d549.
//
// Solidity: function approveTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes32 txHash)
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) ApproveTransaction(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "approveTransaction", target, value, data, nonce)
}

// ApproveTransaction is a paid mutator transaction binding the contract method 0x1c91d549.
//
// Solidity: function approveTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes32 txHash)
func (_ThresholdMultiSig *ThresholdMultiSigSession) ApproveTransaction(target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ApproveTransaction(&_ThresholdMultiSig.TransactOpts, target, value, data, nonce)
}

// ApproveTransaction is a paid mutator transaction binding the contract method 0x1c91d549.
//
// Solidity: function approveTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes32 txHash)
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) ApproveTransaction(target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ApproveTransaction(&_ThresholdMultiSig.TransactOpts, target, value, data, nonce)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 newThreshold) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) ChangeThreshold(opts *bind.TransactOpts, newThreshold *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "changeThreshold", newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 newThreshold) returns()
func (_ThresholdMultiSig *ThresholdMultiSigSession) ChangeThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ChangeThreshold(&_ThresholdMultiSig.TransactOpts, newThreshold)
}

// ChangeThreshold is a paid mutator transaction binding the contract method 0x694e80c3.
//
// Solidity: function changeThreshold(uint256 newThreshold) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) ChangeThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ChangeThreshold(&_ThresholdMultiSig.TransactOpts, newThreshold)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x06a41d09.
//
// Solidity: function executeTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes result)
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) ExecuteTransaction(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "executeTransaction", target, value, data, nonce)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x06a41d09.
//
// Solidity: function executeTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes result)
func (_ThresholdMultiSig *ThresholdMultiSigSession) ExecuteTransaction(target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ExecuteTransaction(&_ThresholdMultiSig.TransactOpts, target, value, data, nonce)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x06a41d09.
//
// Solidity: function executeTransaction(address target, uint256 value, bytes data, uint256 nonce) returns(bytes result)
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) ExecuteTransaction(target common.Address, value *big.Int, data []byte, nonce *big.Int) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ExecuteTransaction(&_ThresholdMultiSig.TransactOpts, target, value, data, nonce)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) RemoveOwner(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "removeOwner", owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigSession) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.RemoveOwner(&_ThresholdMultiSig.TransactOpts, owner)
}

// RemoveOwner is a paid mutator transaction binding the contract method 0x173825d9.
//
// Solidity: function removeOwner(address owner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) RemoveOwner(owner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.RemoveOwner(&_ThresholdMultiSig.TransactOpts, owner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(address oldOwner, address newOwner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) ReplaceOwner(opts *bind.TransactOpts, oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.Transact(opts, "replaceOwner", oldOwner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(address oldOwner, address newOwner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigSession) ReplaceOwner(oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ReplaceOwner(&_ThresholdMultiSig.TransactOpts, oldOwner, newOwner)
}

// ReplaceOwner is a paid mutator transaction binding the contract method 0xe20056e6.
//
// Solidity: function replaceOwner(address oldOwner, address newOwner) returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) ReplaceOwner(oldOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.ReplaceOwner(&_ThresholdMultiSig.TransactOpts, oldOwner, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ThresholdMultiSig.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ThresholdMultiSig *ThresholdMultiSigSession) Receive() (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.Receive(&_ThresholdMultiSig.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ThresholdMultiSig *ThresholdMultiSigTransactorSession) Receive() (*types.Transaction, error) {
	return _ThresholdMultiSig.Contract.Receive(&_ThresholdMultiSig.TransactOpts)
}

// ThresholdMultiSigOwnerAddedIterator is returned from FilterOwnerAdded and is used to iterate over the raw logs and unpacked data for OwnerAdded events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerAddedIterator struct {
	Event *ThresholdMultiSigOwnerAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigOwnerAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigOwnerAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigOwnerAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigOwnerAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigOwnerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigOwnerAdded represents a OwnerAdded event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerAdded struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerAdded is a free log retrieval operation binding the contract event 0x994a936646fe87ffe4f1e469d3d6aa417d6b855598397f323de5b449f765f0c3.
//
// Solidity: event OwnerAdded(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterOwnerAdded(opts *bind.FilterOpts, owner []common.Address) (*ThresholdMultiSigOwnerAddedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "OwnerAdded", ownerRule)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigOwnerAddedIterator{contract: _ThresholdMultiSig.contract, event: "OwnerAdded", logs: logs, sub: sub}, nil
}

// WatchOwnerAdded is a free log subscription operation binding the contract event 0x994a936646fe87ffe4f1e469d3d6aa417d6b855598397f323de5b449f765f0c3.
//
// Solidity: event OwnerAdded(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchOwnerAdded(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigOwnerAdded, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "OwnerAdded", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigOwnerAdded)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerAdded is a log parse operation binding the contract event 0x994a936646fe87ffe4f1e469d3d6aa417d6b855598397f323de5b449f765f0c3.
//
// Solidity: event OwnerAdded(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseOwnerAdded(log types.Log) (*ThresholdMultiSigOwnerAdded, error) {
	event := new(ThresholdMultiSigOwnerAdded)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThresholdMultiSigOwnerRemovedIterator is returned from FilterOwnerRemoved and is used to iterate over the raw logs and unpacked data for OwnerRemoved events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerRemovedIterator struct {
	Event *ThresholdMultiSigOwnerRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigOwnerRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigOwnerRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigOwnerRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigOwnerRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigOwnerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigOwnerRemoved represents a OwnerRemoved event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerRemoved struct {
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOwnerRemoved is a free log retrieval operation binding the contract event 0x58619076adf5bb0943d100ef88d52d7c3fd691b19d3a9071b555b651fbf418da.
//
// Solidity: event OwnerRemoved(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterOwnerRemoved(opts *bind.FilterOpts, owner []common.Address) (*ThresholdMultiSigOwnerRemovedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "OwnerRemoved", ownerRule)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigOwnerRemovedIterator{contract: _ThresholdMultiSig.contract, event: "OwnerRemoved", logs: logs, sub: sub}, nil
}

// WatchOwnerRemoved is a free log subscription operation binding the contract event 0x58619076adf5bb0943d100ef88d52d7c3fd691b19d3a9071b555b651fbf418da.
//
// Solidity: event OwnerRemoved(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchOwnerRemoved(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigOwnerRemoved, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "OwnerRemoved", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigOwnerRemoved)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerRemoved is a log parse operation binding the contract event 0x58619076adf5bb0943d100ef88d52d7c3fd691b19d3a9071b555b651fbf418da.
//
// Solidity: event OwnerRemoved(address indexed owner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseOwnerRemoved(log types.Log) (*ThresholdMultiSigOwnerRemoved, error) {
	event := new(ThresholdMultiSigOwnerRemoved)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThresholdMultiSigOwnerReplacedIterator is returned from FilterOwnerReplaced and is used to iterate over the raw logs and unpacked data for OwnerReplaced events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerReplacedIterator struct {
	Event *ThresholdMultiSigOwnerReplaced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigOwnerReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigOwnerReplaced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigOwnerReplaced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigOwnerReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigOwnerReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigOwnerReplaced represents a OwnerReplaced event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigOwnerReplaced struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerReplaced is a free log retrieval operation binding the contract event 0x0647e02c730b94c30d2e3d2002bcbc052266c1c2ae828a4bafa8a46f7bd415ba.
//
// Solidity: event OwnerReplaced(address indexed oldOwner, address indexed newOwner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterOwnerReplaced(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*ThresholdMultiSigOwnerReplacedIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "OwnerReplaced", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigOwnerReplacedIterator{contract: _ThresholdMultiSig.contract, event: "OwnerReplaced", logs: logs, sub: sub}, nil
}

// WatchOwnerReplaced is a free log subscription operation binding the contract event 0x0647e02c730b94c30d2e3d2002bcbc052266c1c2ae828a4bafa8a46f7bd415ba.
//
// Solidity: event OwnerReplaced(address indexed oldOwner, address indexed newOwner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchOwnerReplaced(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigOwnerReplaced, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "OwnerReplaced", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigOwnerReplaced)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerReplaced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnerReplaced is a log parse operation binding the contract event 0x0647e02c730b94c30d2e3d2002bcbc052266c1c2ae828a4bafa8a46f7bd415ba.
//
// Solidity: event OwnerReplaced(address indexed oldOwner, address indexed newOwner)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseOwnerReplaced(log types.Log) (*ThresholdMultiSigOwnerReplaced, error) {
	event := new(ThresholdMultiSigOwnerReplaced)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "OwnerReplaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThresholdMultiSigThresholdChangedIterator is returned from FilterThresholdChanged and is used to iterate over the raw logs and unpacked data for ThresholdChanged events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigThresholdChangedIterator struct {
	Event *ThresholdMultiSigThresholdChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigThresholdChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigThresholdChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigThresholdChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigThresholdChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigThresholdChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigThresholdChanged represents a ThresholdChanged event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigThresholdChanged struct {
	OldThreshold *big.Int
	NewThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterThresholdChanged is a free log retrieval operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterThresholdChanged(opts *bind.FilterOpts) (*ThresholdMultiSigThresholdChangedIterator, error) {

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "ThresholdChanged")
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigThresholdChangedIterator{contract: _ThresholdMultiSig.contract, event: "ThresholdChanged", logs: logs, sub: sub}, nil
}

// WatchThresholdChanged is a free log subscription operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchThresholdChanged(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigThresholdChanged) (event.Subscription, error) {

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "ThresholdChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigThresholdChanged)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "ThresholdChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseThresholdChanged is a log parse operation binding the contract event 0x3164947cf0f49f08dd0cd80e671535b1e11590d347c55dcaa97ba3c24a96b33a.
//
// Solidity: event ThresholdChanged(uint256 oldThreshold, uint256 newThreshold)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseThresholdChanged(log types.Log) (*ThresholdMultiSigThresholdChanged, error) {
	event := new(ThresholdMultiSigThresholdChanged)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "ThresholdChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThresholdMultiSigTransactionApprovedIterator is returned from FilterTransactionApproved and is used to iterate over the raw logs and unpacked data for TransactionApproved events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigTransactionApprovedIterator struct {
	Event *ThresholdMultiSigTransactionApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigTransactionApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigTransactionApproved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigTransactionApproved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigTransactionApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigTransactionApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigTransactionApproved represents a TransactionApproved event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigTransactionApproved struct {
	TxHash        [32]byte
	Owner         common.Address
	ApprovalCount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTransactionApproved is a free log retrieval operation binding the contract event 0x36088d8c95f10317cb92d7a374fc7ca80d20040c9e5d5a154dc30455cf41dc87.
//
// Solidity: event TransactionApproved(bytes32 indexed txHash, address indexed owner, uint256 approvalCount)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterTransactionApproved(opts *bind.FilterOpts, txHash [][32]byte, owner []common.Address) (*ThresholdMultiSigTransactionApprovedIterator, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "TransactionApproved", txHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigTransactionApprovedIterator{contract: _ThresholdMultiSig.contract, event: "TransactionApproved", logs: logs, sub: sub}, nil
}

// WatchTransactionApproved is a free log subscription operation binding the contract event 0x36088d8c95f10317cb92d7a374fc7ca80d20040c9e5d5a154dc30455cf41dc87.
//
// Solidity: event TransactionApproved(bytes32 indexed txHash, address indexed owner, uint256 approvalCount)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchTransactionApproved(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigTransactionApproved, txHash [][32]byte, owner []common.Address) (event.Subscription, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "TransactionApproved", txHashRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigTransactionApproved)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "TransactionApproved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionApproved is a log parse operation binding the contract event 0x36088d8c95f10317cb92d7a374fc7ca80d20040c9e5d5a154dc30455cf41dc87.
//
// Solidity: event TransactionApproved(bytes32 indexed txHash, address indexed owner, uint256 approvalCount)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseTransactionApproved(log types.Log) (*ThresholdMultiSigTransactionApproved, error) {
	event := new(ThresholdMultiSigTransactionApproved)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "TransactionApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ThresholdMultiSigTransactionExecutedIterator is returned from FilterTransactionExecuted and is used to iterate over the raw logs and unpacked data for TransactionExecuted events raised by the ThresholdMultiSig contract.
type ThresholdMultiSigTransactionExecutedIterator struct {
	Event *ThresholdMultiSigTransactionExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ThresholdMultiSigTransactionExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ThresholdMultiSigTransactionExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ThresholdMultiSigTransactionExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ThresholdMultiSigTransactionExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ThresholdMultiSigTransactionExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ThresholdMultiSigTransactionExecuted represents a TransactionExecuted event raised by the ThresholdMultiSig contract.
type ThresholdMultiSigTransactionExecuted struct {
	TxHash   [32]byte
	Executor common.Address
	Target   common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransactionExecuted is a free log retrieval operation binding the contract event 0x50204a8f660331b3d1a0869d18a8a64a9cdc964b78d818adcd6b9627b90d44f4.
//
// Solidity: event TransactionExecuted(bytes32 indexed txHash, address indexed executor, address indexed target)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) FilterTransactionExecuted(opts *bind.FilterOpts, txHash [][32]byte, executor []common.Address, target []common.Address) (*ThresholdMultiSigTransactionExecutedIterator, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.FilterLogs(opts, "TransactionExecuted", txHashRule, executorRule, targetRule)
	if err != nil {
		return nil, err
	}
	return &ThresholdMultiSigTransactionExecutedIterator{contract: _ThresholdMultiSig.contract, event: "TransactionExecuted", logs: logs, sub: sub}, nil
}

// WatchTransactionExecuted is a free log subscription operation binding the contract event 0x50204a8f660331b3d1a0869d18a8a64a9cdc964b78d818adcd6b9627b90d44f4.
//
// Solidity: event TransactionExecuted(bytes32 indexed txHash, address indexed executor, address indexed target)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) WatchTransactionExecuted(opts *bind.WatchOpts, sink chan<- *ThresholdMultiSigTransactionExecuted, txHash [][32]byte, executor []common.Address, target []common.Address) (event.Subscription, error) {

	var txHashRule []interface{}
	for _, txHashItem := range txHash {
		txHashRule = append(txHashRule, txHashItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _ThresholdMultiSig.contract.WatchLogs(opts, "TransactionExecuted", txHashRule, executorRule, targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ThresholdMultiSigTransactionExecuted)
				if err := _ThresholdMultiSig.contract.UnpackLog(event, "TransactionExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionExecuted is a log parse operation binding the contract event 0x50204a8f660331b3d1a0869d18a8a64a9cdc964b78d818adcd6b9627b90d44f4.
//
// Solidity: event TransactionExecuted(bytes32 indexed txHash, address indexed executor, address indexed target)
func (_ThresholdMultiSig *ThresholdMultiSigFilterer) ParseTransactionExecuted(log types.Log) (*ThresholdMultiSigTransactionExecuted, error) {
	event := new(ThresholdMultiSigTransactionExecuted)
	if err := _ThresholdMultiSig.contract.UnpackLog(event, "TransactionExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
