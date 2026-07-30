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

// ChainlinkOracleMetaData contains all meta data concerning the ChainlinkOracle contract.
var ChainlinkOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"maxStaleness\",\"type\":\"uint48\"}],\"name\":\"FeedConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"feeds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"maxStaleness\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feed\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"maxStaleness\",\"type\":\"uint48\"}],\"name\":\"setFeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ChainlinkOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainlinkOracleMetaData.ABI instead.
var ChainlinkOracleABI = ChainlinkOracleMetaData.ABI

// ChainlinkOracle is an auto generated Go binding around an Ethereum contract.
type ChainlinkOracle struct {
	ChainlinkOracleCaller     // Read-only binding to the contract
	ChainlinkOracleTransactor // Write-only binding to the contract
	ChainlinkOracleFilterer   // Log filterer for contract events
}

// ChainlinkOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainlinkOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainlinkOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainlinkOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainlinkOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainlinkOracleSession struct {
	Contract     *ChainlinkOracle  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainlinkOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainlinkOracleCallerSession struct {
	Contract *ChainlinkOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ChainlinkOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainlinkOracleTransactorSession struct {
	Contract     *ChainlinkOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ChainlinkOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainlinkOracleRaw struct {
	Contract *ChainlinkOracle // Generic contract binding to access the raw methods on
}

// ChainlinkOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainlinkOracleCallerRaw struct {
	Contract *ChainlinkOracleCaller // Generic read-only contract binding to access the raw methods on
}

// ChainlinkOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainlinkOracleTransactorRaw struct {
	Contract *ChainlinkOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChainlinkOracle creates a new instance of ChainlinkOracle, bound to a specific deployed contract.
func NewChainlinkOracle(address common.Address, backend bind.ContractBackend) (*ChainlinkOracle, error) {
	contract, err := bindChainlinkOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracle{ChainlinkOracleCaller: ChainlinkOracleCaller{contract: contract}, ChainlinkOracleTransactor: ChainlinkOracleTransactor{contract: contract}, ChainlinkOracleFilterer: ChainlinkOracleFilterer{contract: contract}}, nil
}

// NewChainlinkOracleCaller creates a new read-only instance of ChainlinkOracle, bound to a specific deployed contract.
func NewChainlinkOracleCaller(address common.Address, caller bind.ContractCaller) (*ChainlinkOracleCaller, error) {
	contract, err := bindChainlinkOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracleCaller{contract: contract}, nil
}

// NewChainlinkOracleTransactor creates a new write-only instance of ChainlinkOracle, bound to a specific deployed contract.
func NewChainlinkOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainlinkOracleTransactor, error) {
	contract, err := bindChainlinkOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracleTransactor{contract: contract}, nil
}

// NewChainlinkOracleFilterer creates a new log filterer instance of ChainlinkOracle, bound to a specific deployed contract.
func NewChainlinkOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainlinkOracleFilterer, error) {
	contract, err := bindChainlinkOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracleFilterer{contract: contract}, nil
}

// bindChainlinkOracle binds a generic wrapper to an already deployed contract.
func bindChainlinkOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainlinkOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkOracle *ChainlinkOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkOracle.Contract.ChainlinkOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkOracle *ChainlinkOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.ChainlinkOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkOracle *ChainlinkOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.ChainlinkOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ChainlinkOracle *ChainlinkOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ChainlinkOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ChainlinkOracle *ChainlinkOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ChainlinkOracle *ChainlinkOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.contract.Transact(opts, method, params...)
}

// Feeds is a free data retrieval call binding the contract method 0x2fba4aa9.
//
// Solidity: function feeds(address ) view returns(address feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleCaller) Feeds(opts *bind.CallOpts, arg0 common.Address) (struct {
	Feed         common.Address
	MaxStaleness *big.Int
}, error) {
	var out []interface{}
	err := _ChainlinkOracle.contract.Call(opts, &out, "feeds", arg0)

	outstruct := new(struct {
		Feed         common.Address
		MaxStaleness *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Feed = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.MaxStaleness = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Feeds is a free data retrieval call binding the contract method 0x2fba4aa9.
//
// Solidity: function feeds(address ) view returns(address feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleSession) Feeds(arg0 common.Address) (struct {
	Feed         common.Address
	MaxStaleness *big.Int
}, error) {
	return _ChainlinkOracle.Contract.Feeds(&_ChainlinkOracle.CallOpts, arg0)
}

// Feeds is a free data retrieval call binding the contract method 0x2fba4aa9.
//
// Solidity: function feeds(address ) view returns(address feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleCallerSession) Feeds(arg0 common.Address) (struct {
	Feed         common.Address
	MaxStaleness *big.Int
}, error) {
	return _ChainlinkOracle.Contract.Feeds(&_ChainlinkOracle.CallOpts, arg0)
}

// GetPrice is a free data retrieval call binding the contract method 0x41976e09.
//
// Solidity: function getPrice(address token) view returns(uint256)
func (_ChainlinkOracle *ChainlinkOracleCaller) GetPrice(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ChainlinkOracle.contract.Call(opts, &out, "getPrice", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x41976e09.
//
// Solidity: function getPrice(address token) view returns(uint256)
func (_ChainlinkOracle *ChainlinkOracleSession) GetPrice(token common.Address) (*big.Int, error) {
	return _ChainlinkOracle.Contract.GetPrice(&_ChainlinkOracle.CallOpts, token)
}

// GetPrice is a free data retrieval call binding the contract method 0x41976e09.
//
// Solidity: function getPrice(address token) view returns(uint256)
func (_ChainlinkOracle *ChainlinkOracleCallerSession) GetPrice(token common.Address) (*big.Int, error) {
	return _ChainlinkOracle.Contract.GetPrice(&_ChainlinkOracle.CallOpts, token)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChainlinkOracle *ChainlinkOracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ChainlinkOracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChainlinkOracle *ChainlinkOracleSession) Owner() (common.Address, error) {
	return _ChainlinkOracle.Contract.Owner(&_ChainlinkOracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ChainlinkOracle *ChainlinkOracleCallerSession) Owner() (common.Address, error) {
	return _ChainlinkOracle.Contract.Owner(&_ChainlinkOracle.CallOpts)
}

// SetFeed is a paid mutator transaction binding the contract method 0xc740e8f9.
//
// Solidity: function setFeed(address token, address feed, uint48 maxStaleness) returns()
func (_ChainlinkOracle *ChainlinkOracleTransactor) SetFeed(opts *bind.TransactOpts, token common.Address, feed common.Address, maxStaleness *big.Int) (*types.Transaction, error) {
	return _ChainlinkOracle.contract.Transact(opts, "setFeed", token, feed, maxStaleness)
}

// SetFeed is a paid mutator transaction binding the contract method 0xc740e8f9.
//
// Solidity: function setFeed(address token, address feed, uint48 maxStaleness) returns()
func (_ChainlinkOracle *ChainlinkOracleSession) SetFeed(token common.Address, feed common.Address, maxStaleness *big.Int) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.SetFeed(&_ChainlinkOracle.TransactOpts, token, feed, maxStaleness)
}

// SetFeed is a paid mutator transaction binding the contract method 0xc740e8f9.
//
// Solidity: function setFeed(address token, address feed, uint48 maxStaleness) returns()
func (_ChainlinkOracle *ChainlinkOracleTransactorSession) SetFeed(token common.Address, feed common.Address, maxStaleness *big.Int) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.SetFeed(&_ChainlinkOracle.TransactOpts, token, feed, maxStaleness)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChainlinkOracle *ChainlinkOracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ChainlinkOracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChainlinkOracle *ChainlinkOracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.TransferOwnership(&_ChainlinkOracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ChainlinkOracle *ChainlinkOracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ChainlinkOracle.Contract.TransferOwnership(&_ChainlinkOracle.TransactOpts, newOwner)
}

// ChainlinkOracleFeedConfiguredIterator is returned from FilterFeedConfigured and is used to iterate over the raw logs and unpacked data for FeedConfigured events raised by the ChainlinkOracle contract.
type ChainlinkOracleFeedConfiguredIterator struct {
	Event *ChainlinkOracleFeedConfigured // Event containing the contract specifics and raw log

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
func (it *ChainlinkOracleFeedConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainlinkOracleFeedConfigured)
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
		it.Event = new(ChainlinkOracleFeedConfigured)
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
func (it *ChainlinkOracleFeedConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainlinkOracleFeedConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainlinkOracleFeedConfigured represents a FeedConfigured event raised by the ChainlinkOracle contract.
type ChainlinkOracleFeedConfigured struct {
	Token        common.Address
	Feed         common.Address
	MaxStaleness *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFeedConfigured is a free log retrieval operation binding the contract event 0xd8c7b74c5d3aad3a2e584554bbc1f1822f9b057bb8ad568bd21f651cd824219d.
//
// Solidity: event FeedConfigured(address indexed token, address indexed feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleFilterer) FilterFeedConfigured(opts *bind.FilterOpts, token []common.Address, feed []common.Address) (*ChainlinkOracleFeedConfiguredIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var feedRule []interface{}
	for _, feedItem := range feed {
		feedRule = append(feedRule, feedItem)
	}

	logs, sub, err := _ChainlinkOracle.contract.FilterLogs(opts, "FeedConfigured", tokenRule, feedRule)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracleFeedConfiguredIterator{contract: _ChainlinkOracle.contract, event: "FeedConfigured", logs: logs, sub: sub}, nil
}

// WatchFeedConfigured is a free log subscription operation binding the contract event 0xd8c7b74c5d3aad3a2e584554bbc1f1822f9b057bb8ad568bd21f651cd824219d.
//
// Solidity: event FeedConfigured(address indexed token, address indexed feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleFilterer) WatchFeedConfigured(opts *bind.WatchOpts, sink chan<- *ChainlinkOracleFeedConfigured, token []common.Address, feed []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var feedRule []interface{}
	for _, feedItem := range feed {
		feedRule = append(feedRule, feedItem)
	}

	logs, sub, err := _ChainlinkOracle.contract.WatchLogs(opts, "FeedConfigured", tokenRule, feedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainlinkOracleFeedConfigured)
				if err := _ChainlinkOracle.contract.UnpackLog(event, "FeedConfigured", log); err != nil {
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

// ParseFeedConfigured is a log parse operation binding the contract event 0xd8c7b74c5d3aad3a2e584554bbc1f1822f9b057bb8ad568bd21f651cd824219d.
//
// Solidity: event FeedConfigured(address indexed token, address indexed feed, uint48 maxStaleness)
func (_ChainlinkOracle *ChainlinkOracleFilterer) ParseFeedConfigured(log types.Log) (*ChainlinkOracleFeedConfigured, error) {
	event := new(ChainlinkOracleFeedConfigured)
	if err := _ChainlinkOracle.contract.UnpackLog(event, "FeedConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainlinkOracleOwnerChangedIterator is returned from FilterOwnerChanged and is used to iterate over the raw logs and unpacked data for OwnerChanged events raised by the ChainlinkOracle contract.
type ChainlinkOracleOwnerChangedIterator struct {
	Event *ChainlinkOracleOwnerChanged // Event containing the contract specifics and raw log

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
func (it *ChainlinkOracleOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainlinkOracleOwnerChanged)
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
		it.Event = new(ChainlinkOracleOwnerChanged)
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
func (it *ChainlinkOracleOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainlinkOracleOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainlinkOracleOwnerChanged represents a OwnerChanged event raised by the ChainlinkOracle contract.
type ChainlinkOracleOwnerChanged struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerChanged is a free log retrieval operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_ChainlinkOracle *ChainlinkOracleFilterer) FilterOwnerChanged(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*ChainlinkOracleOwnerChangedIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ChainlinkOracle.contract.FilterLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ChainlinkOracleOwnerChangedIterator{contract: _ChainlinkOracle.contract, event: "OwnerChanged", logs: logs, sub: sub}, nil
}

// WatchOwnerChanged is a free log subscription operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_ChainlinkOracle *ChainlinkOracleFilterer) WatchOwnerChanged(opts *bind.WatchOpts, sink chan<- *ChainlinkOracleOwnerChanged, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ChainlinkOracle.contract.WatchLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainlinkOracleOwnerChanged)
				if err := _ChainlinkOracle.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
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

// ParseOwnerChanged is a log parse operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_ChainlinkOracle *ChainlinkOracleFilterer) ParseOwnerChanged(log types.Log) (*ChainlinkOracleOwnerChanged, error) {
	event := new(ChainlinkOracleOwnerChanged)
	if err := _ChainlinkOracle.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
