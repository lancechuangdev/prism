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

// PrismPoolCreatePoolParams is an auto generated low-level Go binding around an user-defined struct.
type PrismPoolCreatePoolParams struct {
	SettleTime             *big.Int
	MaturityTime           *big.Int
	InterestRate           *big.Int
	MaxLendSupply          *big.Int
	CollateralizationRatio *big.Int
	LendToken              common.Address
	CollateralToken        common.Address
	LenderPositionToken    common.Address
	BorrowerPositionToken  common.Address
	LiquidateRate          *big.Int
}

// PrismPoolPoolBaseInfo is an auto generated low-level Go binding around an user-defined struct.
type PrismPoolPoolBaseInfo struct {
	SettleTime               *big.Int
	MaturityTime             *big.Int
	InterestRate             *big.Int
	MaxLendSupply            *big.Int
	TotalLendDeposited       *big.Int
	TotalCollateralDeposited *big.Int
	CollateralizationRatio   *big.Int
	LendToken                common.Address
	CollateralToken          common.Address
	State                    uint8
	LenderPositionToken      common.Address
	BorrowerPositionToken    common.Address
	LiquidateRate            *big.Int
}

// PrismPoolPoolDataInfo is an auto generated low-level Go binding around an user-defined struct.
type PrismPoolPoolDataInfo struct {
	SettleAmountLend        *big.Int
	SettleAmountBorrow      *big.Int
	FinishAmountLend        *big.Int
	FinishAmountBorrow      *big.Int
	LiquidationAmountLend   *big.Int
	LiquidationAmountBorrow *big.Int
}

// PrismPoolMetaData contains all meta data concerning the PrismPool contract.
var PrismPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dexSwap_\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"feeAddress_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"jpToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"jpAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"loanAmount\",\"type\":\"uint256\"}],\"name\":\"ClaimBorrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"spAmount\",\"type\":\"uint256\"}],\"name\":\"ClaimLend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositBorrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositLend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldDexSwap\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newDexSwap\",\"type\":\"address\"}],\"name\":\"DexSwapChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldFeeAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newFeeAddress\",\"type\":\"address\"}],\"name\":\"FeeAddressChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMinAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMinAmount\",\"type\":\"uint256\"}],\"name\":\"MinBorrowAmountChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldMinAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newMinAmount\",\"type\":\"uint256\"}],\"name\":\"MinLendAmountChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"PauseChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lendToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrowToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"spToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"jpToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"settleTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"collateralSold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lendTokenRecovered\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingCollateralAmount\",\"type\":\"uint256\"}],\"name\":\"PoolLiquidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"collateralSold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"repaymentAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"remainingCollateralAmount\",\"type\":\"uint256\"}],\"name\":\"PoolRepaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RefundBorrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RefundLend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumPrismPool.PoolState\",\"name\":\"oldState\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumPrismPool.PoolState\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"StateChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"borrower\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"jpAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"collateralAmount\",\"type\":\"uint256\"}],\"name\":\"WithdrawBorrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"spAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lendAmount\",\"type\":\"uint256\"}],\"name\":\"WithdrawLend\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"claimBorrowerPositionAndLoan\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"claimLenderPosition\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"settleTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturityTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxLendSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateralizationRatio\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lendToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"collateralToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lenderPositionToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrowerPositionToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidateRate\",\"type\":\"uint256\"}],\"internalType\":\"structPrismPool.CreatePoolParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"createPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositBorrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositLend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dexSwap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"getPool\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"settleTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturityTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxLendSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalLendDeposited\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCollateralDeposited\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateralizationRatio\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lendToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"collateralToken\",\"type\":\"address\"},{\"internalType\":\"enumPrismPool.PoolState\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"lenderPositionToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"borrowerPositionToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidateRate\",\"type\":\"uint256\"}],\"internalType\":\"structPrismPool.PoolBaseInfo\",\"name\":\"poolInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"getPoolData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"settleAmountLend\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"settleAmountBorrow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"finishAmountLend\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"finishAmountBorrow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationAmountLend\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationAmountBorrow\",\"type\":\"uint256\"}],\"internalType\":\"structPrismPool.PoolDataInfo\",\"name\":\"poolDataInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"getPoolState\",\"outputs\":[{\"internalType\":\"enumPrismPool.PoolState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"getRequiredRepayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalPaused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"isBeforeSettleTime\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"isUndercollateralized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCollateralAmount\",\"type\":\"uint256\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minBorrowAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minLendAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"refundExcessCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"refundExcessLend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCollateralAmount\",\"type\":\"uint256\"}],\"name\":\"repayPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userBorrowInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"hasRefunded\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"hasClaimed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userLendInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"refundAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"hasRefunded\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"hasClaimed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"borrowerPosition\",\"type\":\"uint256\"}],\"name\":\"withdrawBorrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"poolId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lenderPosition\",\"type\":\"uint256\"}],\"name\":\"withdrawLend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PrismPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use PrismPoolMetaData.ABI instead.
var PrismPoolABI = PrismPoolMetaData.ABI

// PrismPool is an auto generated Go binding around an Ethereum contract.
type PrismPool struct {
	PrismPoolCaller     // Read-only binding to the contract
	PrismPoolTransactor // Write-only binding to the contract
	PrismPoolFilterer   // Log filterer for contract events
}

// PrismPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type PrismPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrismPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrismPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrismPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrismPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrismPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PrismPoolSession struct {
	Contract     *PrismPool        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PrismPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PrismPoolCallerSession struct {
	Contract *PrismPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PrismPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PrismPoolTransactorSession struct {
	Contract     *PrismPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PrismPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type PrismPoolRaw struct {
	Contract *PrismPool // Generic contract binding to access the raw methods on
}

// PrismPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PrismPoolCallerRaw struct {
	Contract *PrismPoolCaller // Generic read-only contract binding to access the raw methods on
}

// PrismPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PrismPoolTransactorRaw struct {
	Contract *PrismPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPrismPool creates a new instance of PrismPool, bound to a specific deployed contract.
func NewPrismPool(address common.Address, backend bind.ContractBackend) (*PrismPool, error) {
	contract, err := bindPrismPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PrismPool{PrismPoolCaller: PrismPoolCaller{contract: contract}, PrismPoolTransactor: PrismPoolTransactor{contract: contract}, PrismPoolFilterer: PrismPoolFilterer{contract: contract}}, nil
}

// NewPrismPoolCaller creates a new read-only instance of PrismPool, bound to a specific deployed contract.
func NewPrismPoolCaller(address common.Address, caller bind.ContractCaller) (*PrismPoolCaller, error) {
	contract, err := bindPrismPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PrismPoolCaller{contract: contract}, nil
}

// NewPrismPoolTransactor creates a new write-only instance of PrismPool, bound to a specific deployed contract.
func NewPrismPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*PrismPoolTransactor, error) {
	contract, err := bindPrismPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PrismPoolTransactor{contract: contract}, nil
}

// NewPrismPoolFilterer creates a new log filterer instance of PrismPool, bound to a specific deployed contract.
func NewPrismPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*PrismPoolFilterer, error) {
	contract, err := bindPrismPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PrismPoolFilterer{contract: contract}, nil
}

// bindPrismPool binds a generic wrapper to an already deployed contract.
func bindPrismPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PrismPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PrismPool *PrismPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PrismPool.Contract.PrismPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PrismPool *PrismPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrismPool.Contract.PrismPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PrismPool *PrismPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PrismPool.Contract.PrismPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PrismPool *PrismPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PrismPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PrismPool *PrismPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PrismPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PrismPool *PrismPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PrismPool.Contract.contract.Transact(opts, method, params...)
}

// DexSwap is a free data retrieval call binding the contract method 0xb15a068d.
//
// Solidity: function dexSwap() view returns(address)
func (_PrismPool *PrismPoolCaller) DexSwap(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "dexSwap")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DexSwap is a free data retrieval call binding the contract method 0xb15a068d.
//
// Solidity: function dexSwap() view returns(address)
func (_PrismPool *PrismPoolSession) DexSwap() (common.Address, error) {
	return _PrismPool.Contract.DexSwap(&_PrismPool.CallOpts)
}

// DexSwap is a free data retrieval call binding the contract method 0xb15a068d.
//
// Solidity: function dexSwap() view returns(address)
func (_PrismPool *PrismPoolCallerSession) DexSwap() (common.Address, error) {
	return _PrismPool.Contract.DexSwap(&_PrismPool.CallOpts)
}

// FeeAddress is a free data retrieval call binding the contract method 0x41275358.
//
// Solidity: function feeAddress() view returns(address)
func (_PrismPool *PrismPoolCaller) FeeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "feeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeAddress is a free data retrieval call binding the contract method 0x41275358.
//
// Solidity: function feeAddress() view returns(address)
func (_PrismPool *PrismPoolSession) FeeAddress() (common.Address, error) {
	return _PrismPool.Contract.FeeAddress(&_PrismPool.CallOpts)
}

// FeeAddress is a free data retrieval call binding the contract method 0x41275358.
//
// Solidity: function feeAddress() view returns(address)
func (_PrismPool *PrismPoolCallerSession) FeeAddress() (common.Address, error) {
	return _PrismPool.Contract.FeeAddress(&_PrismPool.CallOpts)
}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint8,address,address,uint256) poolInfo)
func (_PrismPool *PrismPoolCaller) GetPool(opts *bind.CallOpts, poolId *big.Int) (PrismPoolPoolBaseInfo, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "getPool", poolId)

	if err != nil {
		return *new(PrismPoolPoolBaseInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(PrismPoolPoolBaseInfo)).(*PrismPoolPoolBaseInfo)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint8,address,address,uint256) poolInfo)
func (_PrismPool *PrismPoolSession) GetPool(poolId *big.Int) (PrismPoolPoolBaseInfo, error) {
	return _PrismPool.Contract.GetPool(&_PrismPool.CallOpts, poolId)
}

// GetPool is a free data retrieval call binding the contract method 0x068bcd8d.
//
// Solidity: function getPool(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,address,uint8,address,address,uint256) poolInfo)
func (_PrismPool *PrismPoolCallerSession) GetPool(poolId *big.Int) (PrismPoolPoolBaseInfo, error) {
	return _PrismPool.Contract.GetPool(&_PrismPool.CallOpts, poolId)
}

// GetPoolData is a free data retrieval call binding the contract method 0xd21c87ad.
//
// Solidity: function getPoolData(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256) poolDataInfo)
func (_PrismPool *PrismPoolCaller) GetPoolData(opts *bind.CallOpts, poolId *big.Int) (PrismPoolPoolDataInfo, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "getPoolData", poolId)

	if err != nil {
		return *new(PrismPoolPoolDataInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(PrismPoolPoolDataInfo)).(*PrismPoolPoolDataInfo)

	return out0, err

}

// GetPoolData is a free data retrieval call binding the contract method 0xd21c87ad.
//
// Solidity: function getPoolData(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256) poolDataInfo)
func (_PrismPool *PrismPoolSession) GetPoolData(poolId *big.Int) (PrismPoolPoolDataInfo, error) {
	return _PrismPool.Contract.GetPoolData(&_PrismPool.CallOpts, poolId)
}

// GetPoolData is a free data retrieval call binding the contract method 0xd21c87ad.
//
// Solidity: function getPoolData(uint256 poolId) view returns((uint256,uint256,uint256,uint256,uint256,uint256) poolDataInfo)
func (_PrismPool *PrismPoolCallerSession) GetPoolData(poolId *big.Int) (PrismPoolPoolDataInfo, error) {
	return _PrismPool.Contract.GetPoolData(&_PrismPool.CallOpts, poolId)
}

// GetPoolState is a free data retrieval call binding the contract method 0xb1597517.
//
// Solidity: function getPoolState(uint256 poolId) view returns(uint8 state)
func (_PrismPool *PrismPoolCaller) GetPoolState(opts *bind.CallOpts, poolId *big.Int) (uint8, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "getPoolState", poolId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPoolState is a free data retrieval call binding the contract method 0xb1597517.
//
// Solidity: function getPoolState(uint256 poolId) view returns(uint8 state)
func (_PrismPool *PrismPoolSession) GetPoolState(poolId *big.Int) (uint8, error) {
	return _PrismPool.Contract.GetPoolState(&_PrismPool.CallOpts, poolId)
}

// GetPoolState is a free data retrieval call binding the contract method 0xb1597517.
//
// Solidity: function getPoolState(uint256 poolId) view returns(uint8 state)
func (_PrismPool *PrismPoolCallerSession) GetPoolState(poolId *big.Int) (uint8, error) {
	return _PrismPool.Contract.GetPoolState(&_PrismPool.CallOpts, poolId)
}

// GetRequiredRepayment is a free data retrieval call binding the contract method 0xe628cf27.
//
// Solidity: function getRequiredRepayment(uint256 poolId) view returns(uint256)
func (_PrismPool *PrismPoolCaller) GetRequiredRepayment(opts *bind.CallOpts, poolId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "getRequiredRepayment", poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRequiredRepayment is a free data retrieval call binding the contract method 0xe628cf27.
//
// Solidity: function getRequiredRepayment(uint256 poolId) view returns(uint256)
func (_PrismPool *PrismPoolSession) GetRequiredRepayment(poolId *big.Int) (*big.Int, error) {
	return _PrismPool.Contract.GetRequiredRepayment(&_PrismPool.CallOpts, poolId)
}

// GetRequiredRepayment is a free data retrieval call binding the contract method 0xe628cf27.
//
// Solidity: function getRequiredRepayment(uint256 poolId) view returns(uint256)
func (_PrismPool *PrismPoolCallerSession) GetRequiredRepayment(poolId *big.Int) (*big.Int, error) {
	return _PrismPool.Contract.GetRequiredRepayment(&_PrismPool.CallOpts, poolId)
}

// GlobalPaused is a free data retrieval call binding the contract method 0x61a552dc.
//
// Solidity: function globalPaused() view returns(bool)
func (_PrismPool *PrismPoolCaller) GlobalPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "globalPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GlobalPaused is a free data retrieval call binding the contract method 0x61a552dc.
//
// Solidity: function globalPaused() view returns(bool)
func (_PrismPool *PrismPoolSession) GlobalPaused() (bool, error) {
	return _PrismPool.Contract.GlobalPaused(&_PrismPool.CallOpts)
}

// GlobalPaused is a free data retrieval call binding the contract method 0x61a552dc.
//
// Solidity: function globalPaused() view returns(bool)
func (_PrismPool *PrismPoolCallerSession) GlobalPaused() (bool, error) {
	return _PrismPool.Contract.GlobalPaused(&_PrismPool.CallOpts)
}

// IsBeforeSettleTime is a free data retrieval call binding the contract method 0x0e812252.
//
// Solidity: function isBeforeSettleTime(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolCaller) IsBeforeSettleTime(opts *bind.CallOpts, poolId *big.Int) (bool, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "isBeforeSettleTime", poolId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBeforeSettleTime is a free data retrieval call binding the contract method 0x0e812252.
//
// Solidity: function isBeforeSettleTime(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolSession) IsBeforeSettleTime(poolId *big.Int) (bool, error) {
	return _PrismPool.Contract.IsBeforeSettleTime(&_PrismPool.CallOpts, poolId)
}

// IsBeforeSettleTime is a free data retrieval call binding the contract method 0x0e812252.
//
// Solidity: function isBeforeSettleTime(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolCallerSession) IsBeforeSettleTime(poolId *big.Int) (bool, error) {
	return _PrismPool.Contract.IsBeforeSettleTime(&_PrismPool.CallOpts, poolId)
}

// IsUndercollateralized is a free data retrieval call binding the contract method 0x3f80c1df.
//
// Solidity: function isUndercollateralized(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolCaller) IsUndercollateralized(opts *bind.CallOpts, poolId *big.Int) (bool, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "isUndercollateralized", poolId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUndercollateralized is a free data retrieval call binding the contract method 0x3f80c1df.
//
// Solidity: function isUndercollateralized(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolSession) IsUndercollateralized(poolId *big.Int) (bool, error) {
	return _PrismPool.Contract.IsUndercollateralized(&_PrismPool.CallOpts, poolId)
}

// IsUndercollateralized is a free data retrieval call binding the contract method 0x3f80c1df.
//
// Solidity: function isUndercollateralized(uint256 poolId) view returns(bool)
func (_PrismPool *PrismPoolCallerSession) IsUndercollateralized(poolId *big.Int) (bool, error) {
	return _PrismPool.Contract.IsUndercollateralized(&_PrismPool.CallOpts, poolId)
}

// MinBorrowAmount is a free data retrieval call binding the contract method 0x0a5db8a8.
//
// Solidity: function minBorrowAmount() view returns(uint256)
func (_PrismPool *PrismPoolCaller) MinBorrowAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "minBorrowAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinBorrowAmount is a free data retrieval call binding the contract method 0x0a5db8a8.
//
// Solidity: function minBorrowAmount() view returns(uint256)
func (_PrismPool *PrismPoolSession) MinBorrowAmount() (*big.Int, error) {
	return _PrismPool.Contract.MinBorrowAmount(&_PrismPool.CallOpts)
}

// MinBorrowAmount is a free data retrieval call binding the contract method 0x0a5db8a8.
//
// Solidity: function minBorrowAmount() view returns(uint256)
func (_PrismPool *PrismPoolCallerSession) MinBorrowAmount() (*big.Int, error) {
	return _PrismPool.Contract.MinBorrowAmount(&_PrismPool.CallOpts)
}

// MinLendAmount is a free data retrieval call binding the contract method 0x7b567d6a.
//
// Solidity: function minLendAmount() view returns(uint256)
func (_PrismPool *PrismPoolCaller) MinLendAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "minLendAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinLendAmount is a free data retrieval call binding the contract method 0x7b567d6a.
//
// Solidity: function minLendAmount() view returns(uint256)
func (_PrismPool *PrismPoolSession) MinLendAmount() (*big.Int, error) {
	return _PrismPool.Contract.MinLendAmount(&_PrismPool.CallOpts)
}

// MinLendAmount is a free data retrieval call binding the contract method 0x7b567d6a.
//
// Solidity: function minLendAmount() view returns(uint256)
func (_PrismPool *PrismPoolCallerSession) MinLendAmount() (*big.Int, error) {
	return _PrismPool.Contract.MinLendAmount(&_PrismPool.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PrismPool *PrismPoolCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PrismPool *PrismPoolSession) Oracle() (common.Address, error) {
	return _PrismPool.Contract.Oracle(&_PrismPool.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_PrismPool *PrismPoolCallerSession) Oracle() (common.Address, error) {
	return _PrismPool.Contract.Oracle(&_PrismPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrismPool *PrismPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrismPool *PrismPoolSession) Owner() (common.Address, error) {
	return _PrismPool.Contract.Owner(&_PrismPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PrismPool *PrismPoolCallerSession) Owner() (common.Address, error) {
	return _PrismPool.Contract.Owner(&_PrismPool.CallOpts)
}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_PrismPool *PrismPoolCaller) PoolCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "poolCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_PrismPool *PrismPoolSession) PoolCount() (*big.Int, error) {
	return _PrismPool.Contract.PoolCount(&_PrismPool.CallOpts)
}

// PoolCount is a free data retrieval call binding the contract method 0xf525cb68.
//
// Solidity: function poolCount() view returns(uint256)
func (_PrismPool *PrismPoolCallerSession) PoolCount() (*big.Int, error) {
	return _PrismPool.Contract.PoolCount(&_PrismPool.CallOpts)
}

// UserBorrowInfo is a free data retrieval call binding the contract method 0x3c9fadc3.
//
// Solidity: function userBorrowInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolCaller) UserBorrowInfo(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "userBorrowInfo", arg0, arg1)

	outstruct := new(struct {
		StakeAmount  *big.Int
		RefundAmount *big.Int
		HasRefunded  bool
		HasClaimed   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakeAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RefundAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.HasRefunded = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.HasClaimed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// UserBorrowInfo is a free data retrieval call binding the contract method 0x3c9fadc3.
//
// Solidity: function userBorrowInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolSession) UserBorrowInfo(arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	return _PrismPool.Contract.UserBorrowInfo(&_PrismPool.CallOpts, arg0, arg1)
}

// UserBorrowInfo is a free data retrieval call binding the contract method 0x3c9fadc3.
//
// Solidity: function userBorrowInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolCallerSession) UserBorrowInfo(arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	return _PrismPool.Contract.UserBorrowInfo(&_PrismPool.CallOpts, arg0, arg1)
}

// UserLendInfo is a free data retrieval call binding the contract method 0xbb176a64.
//
// Solidity: function userLendInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolCaller) UserLendInfo(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	var out []interface{}
	err := _PrismPool.contract.Call(opts, &out, "userLendInfo", arg0, arg1)

	outstruct := new(struct {
		StakeAmount  *big.Int
		RefundAmount *big.Int
		HasRefunded  bool
		HasClaimed   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakeAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RefundAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.HasRefunded = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.HasClaimed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// UserLendInfo is a free data retrieval call binding the contract method 0xbb176a64.
//
// Solidity: function userLendInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolSession) UserLendInfo(arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	return _PrismPool.Contract.UserLendInfo(&_PrismPool.CallOpts, arg0, arg1)
}

// UserLendInfo is a free data retrieval call binding the contract method 0xbb176a64.
//
// Solidity: function userLendInfo(address , uint256 ) view returns(uint256 stakeAmount, uint256 refundAmount, bool hasRefunded, bool hasClaimed)
func (_PrismPool *PrismPoolCallerSession) UserLendInfo(arg0 common.Address, arg1 *big.Int) (struct {
	StakeAmount  *big.Int
	RefundAmount *big.Int
	HasRefunded  bool
	HasClaimed   bool
}, error) {
	return _PrismPool.Contract.UserLendInfo(&_PrismPool.CallOpts, arg0, arg1)
}

// ClaimBorrowerPositionAndLoan is a paid mutator transaction binding the contract method 0x5b793a52.
//
// Solidity: function claimBorrowerPositionAndLoan(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactor) ClaimBorrowerPositionAndLoan(opts *bind.TransactOpts, poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "claimBorrowerPositionAndLoan", poolId)
}

// ClaimBorrowerPositionAndLoan is a paid mutator transaction binding the contract method 0x5b793a52.
//
// Solidity: function claimBorrowerPositionAndLoan(uint256 poolId) returns()
func (_PrismPool *PrismPoolSession) ClaimBorrowerPositionAndLoan(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.ClaimBorrowerPositionAndLoan(&_PrismPool.TransactOpts, poolId)
}

// ClaimBorrowerPositionAndLoan is a paid mutator transaction binding the contract method 0x5b793a52.
//
// Solidity: function claimBorrowerPositionAndLoan(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactorSession) ClaimBorrowerPositionAndLoan(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.ClaimBorrowerPositionAndLoan(&_PrismPool.TransactOpts, poolId)
}

// ClaimLenderPosition is a paid mutator transaction binding the contract method 0x9ae4c4bb.
//
// Solidity: function claimLenderPosition(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactor) ClaimLenderPosition(opts *bind.TransactOpts, poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "claimLenderPosition", poolId)
}

// ClaimLenderPosition is a paid mutator transaction binding the contract method 0x9ae4c4bb.
//
// Solidity: function claimLenderPosition(uint256 poolId) returns()
func (_PrismPool *PrismPoolSession) ClaimLenderPosition(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.ClaimLenderPosition(&_PrismPool.TransactOpts, poolId)
}

// ClaimLenderPosition is a paid mutator transaction binding the contract method 0x9ae4c4bb.
//
// Solidity: function claimLenderPosition(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactorSession) ClaimLenderPosition(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.ClaimLenderPosition(&_PrismPool.TransactOpts, poolId)
}

// CreatePool is a paid mutator transaction binding the contract method 0x2d7af114.
//
// Solidity: function createPool((uint256,uint256,uint256,uint256,uint256,address,address,address,address,uint256) params) returns(uint256 poolId)
func (_PrismPool *PrismPoolTransactor) CreatePool(opts *bind.TransactOpts, params PrismPoolCreatePoolParams) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "createPool", params)
}

// CreatePool is a paid mutator transaction binding the contract method 0x2d7af114.
//
// Solidity: function createPool((uint256,uint256,uint256,uint256,uint256,address,address,address,address,uint256) params) returns(uint256 poolId)
func (_PrismPool *PrismPoolSession) CreatePool(params PrismPoolCreatePoolParams) (*types.Transaction, error) {
	return _PrismPool.Contract.CreatePool(&_PrismPool.TransactOpts, params)
}

// CreatePool is a paid mutator transaction binding the contract method 0x2d7af114.
//
// Solidity: function createPool((uint256,uint256,uint256,uint256,uint256,address,address,address,address,uint256) params) returns(uint256 poolId)
func (_PrismPool *PrismPoolTransactorSession) CreatePool(params PrismPoolCreatePoolParams) (*types.Transaction, error) {
	return _PrismPool.Contract.CreatePool(&_PrismPool.TransactOpts, params)
}

// DepositBorrow is a paid mutator transaction binding the contract method 0x16f941b5.
//
// Solidity: function depositBorrow(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolTransactor) DepositBorrow(opts *bind.TransactOpts, poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "depositBorrow", poolId, amount)
}

// DepositBorrow is a paid mutator transaction binding the contract method 0x16f941b5.
//
// Solidity: function depositBorrow(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolSession) DepositBorrow(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.DepositBorrow(&_PrismPool.TransactOpts, poolId, amount)
}

// DepositBorrow is a paid mutator transaction binding the contract method 0x16f941b5.
//
// Solidity: function depositBorrow(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolTransactorSession) DepositBorrow(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.DepositBorrow(&_PrismPool.TransactOpts, poolId, amount)
}

// DepositLend is a paid mutator transaction binding the contract method 0x90590da0.
//
// Solidity: function depositLend(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolTransactor) DepositLend(opts *bind.TransactOpts, poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "depositLend", poolId, amount)
}

// DepositLend is a paid mutator transaction binding the contract method 0x90590da0.
//
// Solidity: function depositLend(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolSession) DepositLend(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.DepositLend(&_PrismPool.TransactOpts, poolId, amount)
}

// DepositLend is a paid mutator transaction binding the contract method 0x90590da0.
//
// Solidity: function depositLend(uint256 poolId, uint256 amount) returns()
func (_PrismPool *PrismPoolTransactorSession) DepositLend(poolId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.DepositLend(&_PrismPool.TransactOpts, poolId, amount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xd296d1f1.
//
// Solidity: function liquidate(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolTransactor) Liquidate(opts *bind.TransactOpts, poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "liquidate", poolId, maxCollateralAmount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xd296d1f1.
//
// Solidity: function liquidate(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolSession) Liquidate(poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.Liquidate(&_PrismPool.TransactOpts, poolId, maxCollateralAmount)
}

// Liquidate is a paid mutator transaction binding the contract method 0xd296d1f1.
//
// Solidity: function liquidate(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolTransactorSession) Liquidate(poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.Liquidate(&_PrismPool.TransactOpts, poolId, maxCollateralAmount)
}

// RefundExcessCollateral is a paid mutator transaction binding the contract method 0x64ac0885.
//
// Solidity: function refundExcessCollateral(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactor) RefundExcessCollateral(opts *bind.TransactOpts, poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "refundExcessCollateral", poolId)
}

// RefundExcessCollateral is a paid mutator transaction binding the contract method 0x64ac0885.
//
// Solidity: function refundExcessCollateral(uint256 poolId) returns()
func (_PrismPool *PrismPoolSession) RefundExcessCollateral(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RefundExcessCollateral(&_PrismPool.TransactOpts, poolId)
}

// RefundExcessCollateral is a paid mutator transaction binding the contract method 0x64ac0885.
//
// Solidity: function refundExcessCollateral(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactorSession) RefundExcessCollateral(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RefundExcessCollateral(&_PrismPool.TransactOpts, poolId)
}

// RefundExcessLend is a paid mutator transaction binding the contract method 0x91996e7a.
//
// Solidity: function refundExcessLend(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactor) RefundExcessLend(opts *bind.TransactOpts, poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "refundExcessLend", poolId)
}

// RefundExcessLend is a paid mutator transaction binding the contract method 0x91996e7a.
//
// Solidity: function refundExcessLend(uint256 poolId) returns()
func (_PrismPool *PrismPoolSession) RefundExcessLend(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RefundExcessLend(&_PrismPool.TransactOpts, poolId)
}

// RefundExcessLend is a paid mutator transaction binding the contract method 0x91996e7a.
//
// Solidity: function refundExcessLend(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactorSession) RefundExcessLend(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RefundExcessLend(&_PrismPool.TransactOpts, poolId)
}

// RepayPool is a paid mutator transaction binding the contract method 0x3c02db14.
//
// Solidity: function repayPool(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolTransactor) RepayPool(opts *bind.TransactOpts, poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "repayPool", poolId, maxCollateralAmount)
}

// RepayPool is a paid mutator transaction binding the contract method 0x3c02db14.
//
// Solidity: function repayPool(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolSession) RepayPool(poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RepayPool(&_PrismPool.TransactOpts, poolId, maxCollateralAmount)
}

// RepayPool is a paid mutator transaction binding the contract method 0x3c02db14.
//
// Solidity: function repayPool(uint256 poolId, uint256 maxCollateralAmount) returns()
func (_PrismPool *PrismPoolTransactorSession) RepayPool(poolId *big.Int, maxCollateralAmount *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.RepayPool(&_PrismPool.TransactOpts, poolId, maxCollateralAmount)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactor) Settle(opts *bind.TransactOpts, poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "settle", poolId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 poolId) returns()
func (_PrismPool *PrismPoolSession) Settle(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.Settle(&_PrismPool.TransactOpts, poolId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 poolId) returns()
func (_PrismPool *PrismPoolTransactorSession) Settle(poolId *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.Settle(&_PrismPool.TransactOpts, poolId)
}

// WithdrawBorrow is a paid mutator transaction binding the contract method 0x1e107979.
//
// Solidity: function withdrawBorrow(uint256 poolId, uint256 borrowerPosition) returns()
func (_PrismPool *PrismPoolTransactor) WithdrawBorrow(opts *bind.TransactOpts, poolId *big.Int, borrowerPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "withdrawBorrow", poolId, borrowerPosition)
}

// WithdrawBorrow is a paid mutator transaction binding the contract method 0x1e107979.
//
// Solidity: function withdrawBorrow(uint256 poolId, uint256 borrowerPosition) returns()
func (_PrismPool *PrismPoolSession) WithdrawBorrow(poolId *big.Int, borrowerPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.WithdrawBorrow(&_PrismPool.TransactOpts, poolId, borrowerPosition)
}

// WithdrawBorrow is a paid mutator transaction binding the contract method 0x1e107979.
//
// Solidity: function withdrawBorrow(uint256 poolId, uint256 borrowerPosition) returns()
func (_PrismPool *PrismPoolTransactorSession) WithdrawBorrow(poolId *big.Int, borrowerPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.WithdrawBorrow(&_PrismPool.TransactOpts, poolId, borrowerPosition)
}

// WithdrawLend is a paid mutator transaction binding the contract method 0x38f2aa76.
//
// Solidity: function withdrawLend(uint256 poolId, uint256 lenderPosition) returns()
func (_PrismPool *PrismPoolTransactor) WithdrawLend(opts *bind.TransactOpts, poolId *big.Int, lenderPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.contract.Transact(opts, "withdrawLend", poolId, lenderPosition)
}

// WithdrawLend is a paid mutator transaction binding the contract method 0x38f2aa76.
//
// Solidity: function withdrawLend(uint256 poolId, uint256 lenderPosition) returns()
func (_PrismPool *PrismPoolSession) WithdrawLend(poolId *big.Int, lenderPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.WithdrawLend(&_PrismPool.TransactOpts, poolId, lenderPosition)
}

// WithdrawLend is a paid mutator transaction binding the contract method 0x38f2aa76.
//
// Solidity: function withdrawLend(uint256 poolId, uint256 lenderPosition) returns()
func (_PrismPool *PrismPoolTransactorSession) WithdrawLend(poolId *big.Int, lenderPosition *big.Int) (*types.Transaction, error) {
	return _PrismPool.Contract.WithdrawLend(&_PrismPool.TransactOpts, poolId, lenderPosition)
}

// PrismPoolClaimBorrowIterator is returned from FilterClaimBorrow and is used to iterate over the raw logs and unpacked data for ClaimBorrow events raised by the PrismPool contract.
type PrismPoolClaimBorrowIterator struct {
	Event *PrismPoolClaimBorrow // Event containing the contract specifics and raw log

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
func (it *PrismPoolClaimBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolClaimBorrow)
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
		it.Event = new(PrismPoolClaimBorrow)
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
func (it *PrismPoolClaimBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolClaimBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolClaimBorrow represents a ClaimBorrow event raised by the PrismPool contract.
type PrismPoolClaimBorrow struct {
	Borrower   common.Address
	PoolId     *big.Int
	JpToken    common.Address
	JpAmount   *big.Int
	LoanAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimBorrow is a free log retrieval operation binding the contract event 0xecb837876f9a3ff8111d2bf59181e0e5aea8cfd361ba8d26f3324be6785122b5.
//
// Solidity: event ClaimBorrow(address indexed borrower, uint256 indexed poolId, address indexed jpToken, uint256 jpAmount, uint256 loanAmount)
func (_PrismPool *PrismPoolFilterer) FilterClaimBorrow(opts *bind.FilterOpts, borrower []common.Address, poolId []*big.Int, jpToken []common.Address) (*PrismPoolClaimBorrowIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var jpTokenRule []interface{}
	for _, jpTokenItem := range jpToken {
		jpTokenRule = append(jpTokenRule, jpTokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "ClaimBorrow", borrowerRule, poolIdRule, jpTokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolClaimBorrowIterator{contract: _PrismPool.contract, event: "ClaimBorrow", logs: logs, sub: sub}, nil
}

// WatchClaimBorrow is a free log subscription operation binding the contract event 0xecb837876f9a3ff8111d2bf59181e0e5aea8cfd361ba8d26f3324be6785122b5.
//
// Solidity: event ClaimBorrow(address indexed borrower, uint256 indexed poolId, address indexed jpToken, uint256 jpAmount, uint256 loanAmount)
func (_PrismPool *PrismPoolFilterer) WatchClaimBorrow(opts *bind.WatchOpts, sink chan<- *PrismPoolClaimBorrow, borrower []common.Address, poolId []*big.Int, jpToken []common.Address) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var jpTokenRule []interface{}
	for _, jpTokenItem := range jpToken {
		jpTokenRule = append(jpTokenRule, jpTokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "ClaimBorrow", borrowerRule, poolIdRule, jpTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolClaimBorrow)
				if err := _PrismPool.contract.UnpackLog(event, "ClaimBorrow", log); err != nil {
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

// ParseClaimBorrow is a log parse operation binding the contract event 0xecb837876f9a3ff8111d2bf59181e0e5aea8cfd361ba8d26f3324be6785122b5.
//
// Solidity: event ClaimBorrow(address indexed borrower, uint256 indexed poolId, address indexed jpToken, uint256 jpAmount, uint256 loanAmount)
func (_PrismPool *PrismPoolFilterer) ParseClaimBorrow(log types.Log) (*PrismPoolClaimBorrow, error) {
	event := new(PrismPoolClaimBorrow)
	if err := _PrismPool.contract.UnpackLog(event, "ClaimBorrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolClaimLendIterator is returned from FilterClaimLend and is used to iterate over the raw logs and unpacked data for ClaimLend events raised by the PrismPool contract.
type PrismPoolClaimLendIterator struct {
	Event *PrismPoolClaimLend // Event containing the contract specifics and raw log

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
func (it *PrismPoolClaimLendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolClaimLend)
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
		it.Event = new(PrismPoolClaimLend)
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
func (it *PrismPoolClaimLendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolClaimLendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolClaimLend represents a ClaimLend event raised by the PrismPool contract.
type PrismPoolClaimLend struct {
	Lender   common.Address
	PoolId   *big.Int
	SpToken  common.Address
	SpAmount *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterClaimLend is a free log retrieval operation binding the contract event 0xdb740bc3fb95a9071e37f68a4998b8c37ed95dec008652113d5e51c4cc71b4db.
//
// Solidity: event ClaimLend(address indexed lender, uint256 indexed poolId, address indexed spToken, uint256 spAmount)
func (_PrismPool *PrismPoolFilterer) FilterClaimLend(opts *bind.FilterOpts, lender []common.Address, poolId []*big.Int, spToken []common.Address) (*PrismPoolClaimLendIterator, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var spTokenRule []interface{}
	for _, spTokenItem := range spToken {
		spTokenRule = append(spTokenRule, spTokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "ClaimLend", lenderRule, poolIdRule, spTokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolClaimLendIterator{contract: _PrismPool.contract, event: "ClaimLend", logs: logs, sub: sub}, nil
}

// WatchClaimLend is a free log subscription operation binding the contract event 0xdb740bc3fb95a9071e37f68a4998b8c37ed95dec008652113d5e51c4cc71b4db.
//
// Solidity: event ClaimLend(address indexed lender, uint256 indexed poolId, address indexed spToken, uint256 spAmount)
func (_PrismPool *PrismPoolFilterer) WatchClaimLend(opts *bind.WatchOpts, sink chan<- *PrismPoolClaimLend, lender []common.Address, poolId []*big.Int, spToken []common.Address) (event.Subscription, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var spTokenRule []interface{}
	for _, spTokenItem := range spToken {
		spTokenRule = append(spTokenRule, spTokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "ClaimLend", lenderRule, poolIdRule, spTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolClaimLend)
				if err := _PrismPool.contract.UnpackLog(event, "ClaimLend", log); err != nil {
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

// ParseClaimLend is a log parse operation binding the contract event 0xdb740bc3fb95a9071e37f68a4998b8c37ed95dec008652113d5e51c4cc71b4db.
//
// Solidity: event ClaimLend(address indexed lender, uint256 indexed poolId, address indexed spToken, uint256 spAmount)
func (_PrismPool *PrismPoolFilterer) ParseClaimLend(log types.Log) (*PrismPoolClaimLend, error) {
	event := new(PrismPoolClaimLend)
	if err := _PrismPool.contract.UnpackLog(event, "ClaimLend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolDepositBorrowIterator is returned from FilterDepositBorrow and is used to iterate over the raw logs and unpacked data for DepositBorrow events raised by the PrismPool contract.
type PrismPoolDepositBorrowIterator struct {
	Event *PrismPoolDepositBorrow // Event containing the contract specifics and raw log

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
func (it *PrismPoolDepositBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolDepositBorrow)
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
		it.Event = new(PrismPoolDepositBorrow)
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
func (it *PrismPoolDepositBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolDepositBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolDepositBorrow represents a DepositBorrow event raised by the PrismPool contract.
type PrismPoolDepositBorrow struct {
	Borrower common.Address
	PoolId   *big.Int
	Token    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositBorrow is a free log retrieval operation binding the contract event 0x82f800f3a5bbe0a1e7aa3c305036b0df94c09f1611385d14bdd3f1eb6b153146.
//
// Solidity: event DepositBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) FilterDepositBorrow(opts *bind.FilterOpts, borrower []common.Address, poolId []*big.Int, token []common.Address) (*PrismPoolDepositBorrowIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "DepositBorrow", borrowerRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolDepositBorrowIterator{contract: _PrismPool.contract, event: "DepositBorrow", logs: logs, sub: sub}, nil
}

// WatchDepositBorrow is a free log subscription operation binding the contract event 0x82f800f3a5bbe0a1e7aa3c305036b0df94c09f1611385d14bdd3f1eb6b153146.
//
// Solidity: event DepositBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) WatchDepositBorrow(opts *bind.WatchOpts, sink chan<- *PrismPoolDepositBorrow, borrower []common.Address, poolId []*big.Int, token []common.Address) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "DepositBorrow", borrowerRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolDepositBorrow)
				if err := _PrismPool.contract.UnpackLog(event, "DepositBorrow", log); err != nil {
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

// ParseDepositBorrow is a log parse operation binding the contract event 0x82f800f3a5bbe0a1e7aa3c305036b0df94c09f1611385d14bdd3f1eb6b153146.
//
// Solidity: event DepositBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) ParseDepositBorrow(log types.Log) (*PrismPoolDepositBorrow, error) {
	event := new(PrismPoolDepositBorrow)
	if err := _PrismPool.contract.UnpackLog(event, "DepositBorrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolDepositLendIterator is returned from FilterDepositLend and is used to iterate over the raw logs and unpacked data for DepositLend events raised by the PrismPool contract.
type PrismPoolDepositLendIterator struct {
	Event *PrismPoolDepositLend // Event containing the contract specifics and raw log

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
func (it *PrismPoolDepositLendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolDepositLend)
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
		it.Event = new(PrismPoolDepositLend)
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
func (it *PrismPoolDepositLendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolDepositLendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolDepositLend represents a DepositLend event raised by the PrismPool contract.
type PrismPoolDepositLend struct {
	Lender common.Address
	PoolId *big.Int
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDepositLend is a free log retrieval operation binding the contract event 0x1ca06ce3432a8725cf5d40aadde8cb0a311be5440eb78bf327f551f214fbbfa9.
//
// Solidity: event DepositLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) FilterDepositLend(opts *bind.FilterOpts, lender []common.Address, poolId []*big.Int, token []common.Address) (*PrismPoolDepositLendIterator, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "DepositLend", lenderRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolDepositLendIterator{contract: _PrismPool.contract, event: "DepositLend", logs: logs, sub: sub}, nil
}

// WatchDepositLend is a free log subscription operation binding the contract event 0x1ca06ce3432a8725cf5d40aadde8cb0a311be5440eb78bf327f551f214fbbfa9.
//
// Solidity: event DepositLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) WatchDepositLend(opts *bind.WatchOpts, sink chan<- *PrismPoolDepositLend, lender []common.Address, poolId []*big.Int, token []common.Address) (event.Subscription, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "DepositLend", lenderRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolDepositLend)
				if err := _PrismPool.contract.UnpackLog(event, "DepositLend", log); err != nil {
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

// ParseDepositLend is a log parse operation binding the contract event 0x1ca06ce3432a8725cf5d40aadde8cb0a311be5440eb78bf327f551f214fbbfa9.
//
// Solidity: event DepositLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) ParseDepositLend(log types.Log) (*PrismPoolDepositLend, error) {
	event := new(PrismPoolDepositLend)
	if err := _PrismPool.contract.UnpackLog(event, "DepositLend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolDexSwapChangedIterator is returned from FilterDexSwapChanged and is used to iterate over the raw logs and unpacked data for DexSwapChanged events raised by the PrismPool contract.
type PrismPoolDexSwapChangedIterator struct {
	Event *PrismPoolDexSwapChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolDexSwapChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolDexSwapChanged)
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
		it.Event = new(PrismPoolDexSwapChanged)
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
func (it *PrismPoolDexSwapChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolDexSwapChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolDexSwapChanged represents a DexSwapChanged event raised by the PrismPool contract.
type PrismPoolDexSwapChanged struct {
	OldDexSwap common.Address
	NewDexSwap common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDexSwapChanged is a free log retrieval operation binding the contract event 0xae3d09c74a15b2aa044a24e437d9510d24d2cbdf016806a5fb114929b41b4a97.
//
// Solidity: event DexSwapChanged(address indexed oldDexSwap, address indexed newDexSwap)
func (_PrismPool *PrismPoolFilterer) FilterDexSwapChanged(opts *bind.FilterOpts, oldDexSwap []common.Address, newDexSwap []common.Address) (*PrismPoolDexSwapChangedIterator, error) {

	var oldDexSwapRule []interface{}
	for _, oldDexSwapItem := range oldDexSwap {
		oldDexSwapRule = append(oldDexSwapRule, oldDexSwapItem)
	}
	var newDexSwapRule []interface{}
	for _, newDexSwapItem := range newDexSwap {
		newDexSwapRule = append(newDexSwapRule, newDexSwapItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "DexSwapChanged", oldDexSwapRule, newDexSwapRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolDexSwapChangedIterator{contract: _PrismPool.contract, event: "DexSwapChanged", logs: logs, sub: sub}, nil
}

// WatchDexSwapChanged is a free log subscription operation binding the contract event 0xae3d09c74a15b2aa044a24e437d9510d24d2cbdf016806a5fb114929b41b4a97.
//
// Solidity: event DexSwapChanged(address indexed oldDexSwap, address indexed newDexSwap)
func (_PrismPool *PrismPoolFilterer) WatchDexSwapChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolDexSwapChanged, oldDexSwap []common.Address, newDexSwap []common.Address) (event.Subscription, error) {

	var oldDexSwapRule []interface{}
	for _, oldDexSwapItem := range oldDexSwap {
		oldDexSwapRule = append(oldDexSwapRule, oldDexSwapItem)
	}
	var newDexSwapRule []interface{}
	for _, newDexSwapItem := range newDexSwap {
		newDexSwapRule = append(newDexSwapRule, newDexSwapItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "DexSwapChanged", oldDexSwapRule, newDexSwapRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolDexSwapChanged)
				if err := _PrismPool.contract.UnpackLog(event, "DexSwapChanged", log); err != nil {
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

// ParseDexSwapChanged is a log parse operation binding the contract event 0xae3d09c74a15b2aa044a24e437d9510d24d2cbdf016806a5fb114929b41b4a97.
//
// Solidity: event DexSwapChanged(address indexed oldDexSwap, address indexed newDexSwap)
func (_PrismPool *PrismPoolFilterer) ParseDexSwapChanged(log types.Log) (*PrismPoolDexSwapChanged, error) {
	event := new(PrismPoolDexSwapChanged)
	if err := _PrismPool.contract.UnpackLog(event, "DexSwapChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolFeeAddressChangedIterator is returned from FilterFeeAddressChanged and is used to iterate over the raw logs and unpacked data for FeeAddressChanged events raised by the PrismPool contract.
type PrismPoolFeeAddressChangedIterator struct {
	Event *PrismPoolFeeAddressChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolFeeAddressChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolFeeAddressChanged)
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
		it.Event = new(PrismPoolFeeAddressChanged)
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
func (it *PrismPoolFeeAddressChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolFeeAddressChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolFeeAddressChanged represents a FeeAddressChanged event raised by the PrismPool contract.
type PrismPoolFeeAddressChanged struct {
	OldFeeAddress common.Address
	NewFeeAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeeAddressChanged is a free log retrieval operation binding the contract event 0x12aeedbe105c518128d60ce786ade81525990783e6a01c9f849a47954c5043b7.
//
// Solidity: event FeeAddressChanged(address indexed oldFeeAddress, address indexed newFeeAddress)
func (_PrismPool *PrismPoolFilterer) FilterFeeAddressChanged(opts *bind.FilterOpts, oldFeeAddress []common.Address, newFeeAddress []common.Address) (*PrismPoolFeeAddressChangedIterator, error) {

	var oldFeeAddressRule []interface{}
	for _, oldFeeAddressItem := range oldFeeAddress {
		oldFeeAddressRule = append(oldFeeAddressRule, oldFeeAddressItem)
	}
	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "FeeAddressChanged", oldFeeAddressRule, newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolFeeAddressChangedIterator{contract: _PrismPool.contract, event: "FeeAddressChanged", logs: logs, sub: sub}, nil
}

// WatchFeeAddressChanged is a free log subscription operation binding the contract event 0x12aeedbe105c518128d60ce786ade81525990783e6a01c9f849a47954c5043b7.
//
// Solidity: event FeeAddressChanged(address indexed oldFeeAddress, address indexed newFeeAddress)
func (_PrismPool *PrismPoolFilterer) WatchFeeAddressChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolFeeAddressChanged, oldFeeAddress []common.Address, newFeeAddress []common.Address) (event.Subscription, error) {

	var oldFeeAddressRule []interface{}
	for _, oldFeeAddressItem := range oldFeeAddress {
		oldFeeAddressRule = append(oldFeeAddressRule, oldFeeAddressItem)
	}
	var newFeeAddressRule []interface{}
	for _, newFeeAddressItem := range newFeeAddress {
		newFeeAddressRule = append(newFeeAddressRule, newFeeAddressItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "FeeAddressChanged", oldFeeAddressRule, newFeeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolFeeAddressChanged)
				if err := _PrismPool.contract.UnpackLog(event, "FeeAddressChanged", log); err != nil {
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

// ParseFeeAddressChanged is a log parse operation binding the contract event 0x12aeedbe105c518128d60ce786ade81525990783e6a01c9f849a47954c5043b7.
//
// Solidity: event FeeAddressChanged(address indexed oldFeeAddress, address indexed newFeeAddress)
func (_PrismPool *PrismPoolFilterer) ParseFeeAddressChanged(log types.Log) (*PrismPoolFeeAddressChanged, error) {
	event := new(PrismPoolFeeAddressChanged)
	if err := _PrismPool.contract.UnpackLog(event, "FeeAddressChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolMinBorrowAmountChangedIterator is returned from FilterMinBorrowAmountChanged and is used to iterate over the raw logs and unpacked data for MinBorrowAmountChanged events raised by the PrismPool contract.
type PrismPoolMinBorrowAmountChangedIterator struct {
	Event *PrismPoolMinBorrowAmountChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolMinBorrowAmountChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolMinBorrowAmountChanged)
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
		it.Event = new(PrismPoolMinBorrowAmountChanged)
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
func (it *PrismPoolMinBorrowAmountChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolMinBorrowAmountChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolMinBorrowAmountChanged represents a MinBorrowAmountChanged event raised by the PrismPool contract.
type PrismPoolMinBorrowAmountChanged struct {
	OldMinAmount *big.Int
	NewMinAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinBorrowAmountChanged is a free log retrieval operation binding the contract event 0x41fc793f24f29b1ad689ec1bf73adebbb25e4cd1f72e453355b026bff3ec1bc5.
//
// Solidity: event MinBorrowAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) FilterMinBorrowAmountChanged(opts *bind.FilterOpts) (*PrismPoolMinBorrowAmountChangedIterator, error) {

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "MinBorrowAmountChanged")
	if err != nil {
		return nil, err
	}
	return &PrismPoolMinBorrowAmountChangedIterator{contract: _PrismPool.contract, event: "MinBorrowAmountChanged", logs: logs, sub: sub}, nil
}

// WatchMinBorrowAmountChanged is a free log subscription operation binding the contract event 0x41fc793f24f29b1ad689ec1bf73adebbb25e4cd1f72e453355b026bff3ec1bc5.
//
// Solidity: event MinBorrowAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) WatchMinBorrowAmountChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolMinBorrowAmountChanged) (event.Subscription, error) {

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "MinBorrowAmountChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolMinBorrowAmountChanged)
				if err := _PrismPool.contract.UnpackLog(event, "MinBorrowAmountChanged", log); err != nil {
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

// ParseMinBorrowAmountChanged is a log parse operation binding the contract event 0x41fc793f24f29b1ad689ec1bf73adebbb25e4cd1f72e453355b026bff3ec1bc5.
//
// Solidity: event MinBorrowAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) ParseMinBorrowAmountChanged(log types.Log) (*PrismPoolMinBorrowAmountChanged, error) {
	event := new(PrismPoolMinBorrowAmountChanged)
	if err := _PrismPool.contract.UnpackLog(event, "MinBorrowAmountChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolMinLendAmountChangedIterator is returned from FilterMinLendAmountChanged and is used to iterate over the raw logs and unpacked data for MinLendAmountChanged events raised by the PrismPool contract.
type PrismPoolMinLendAmountChangedIterator struct {
	Event *PrismPoolMinLendAmountChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolMinLendAmountChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolMinLendAmountChanged)
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
		it.Event = new(PrismPoolMinLendAmountChanged)
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
func (it *PrismPoolMinLendAmountChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolMinLendAmountChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolMinLendAmountChanged represents a MinLendAmountChanged event raised by the PrismPool contract.
type PrismPoolMinLendAmountChanged struct {
	OldMinAmount *big.Int
	NewMinAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinLendAmountChanged is a free log retrieval operation binding the contract event 0x6c5a4253d44de53b682c9267c472bef489a983331ad855c03fbad09e44216a53.
//
// Solidity: event MinLendAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) FilterMinLendAmountChanged(opts *bind.FilterOpts) (*PrismPoolMinLendAmountChangedIterator, error) {

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "MinLendAmountChanged")
	if err != nil {
		return nil, err
	}
	return &PrismPoolMinLendAmountChangedIterator{contract: _PrismPool.contract, event: "MinLendAmountChanged", logs: logs, sub: sub}, nil
}

// WatchMinLendAmountChanged is a free log subscription operation binding the contract event 0x6c5a4253d44de53b682c9267c472bef489a983331ad855c03fbad09e44216a53.
//
// Solidity: event MinLendAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) WatchMinLendAmountChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolMinLendAmountChanged) (event.Subscription, error) {

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "MinLendAmountChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolMinLendAmountChanged)
				if err := _PrismPool.contract.UnpackLog(event, "MinLendAmountChanged", log); err != nil {
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

// ParseMinLendAmountChanged is a log parse operation binding the contract event 0x6c5a4253d44de53b682c9267c472bef489a983331ad855c03fbad09e44216a53.
//
// Solidity: event MinLendAmountChanged(uint256 oldMinAmount, uint256 newMinAmount)
func (_PrismPool *PrismPoolFilterer) ParseMinLendAmountChanged(log types.Log) (*PrismPoolMinLendAmountChanged, error) {
	event := new(PrismPoolMinLendAmountChanged)
	if err := _PrismPool.contract.UnpackLog(event, "MinLendAmountChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolOwnerChangedIterator is returned from FilterOwnerChanged and is used to iterate over the raw logs and unpacked data for OwnerChanged events raised by the PrismPool contract.
type PrismPoolOwnerChangedIterator struct {
	Event *PrismPoolOwnerChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolOwnerChanged)
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
		it.Event = new(PrismPoolOwnerChanged)
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
func (it *PrismPoolOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolOwnerChanged represents a OwnerChanged event raised by the PrismPool contract.
type PrismPoolOwnerChanged struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerChanged is a free log retrieval operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_PrismPool *PrismPoolFilterer) FilterOwnerChanged(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*PrismPoolOwnerChangedIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolOwnerChangedIterator{contract: _PrismPool.contract, event: "OwnerChanged", logs: logs, sub: sub}, nil
}

// WatchOwnerChanged is a free log subscription operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_PrismPool *PrismPoolFilterer) WatchOwnerChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolOwnerChanged, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolOwnerChanged)
				if err := _PrismPool.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
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
func (_PrismPool *PrismPoolFilterer) ParseOwnerChanged(log types.Log) (*PrismPoolOwnerChanged, error) {
	event := new(PrismPoolOwnerChanged)
	if err := _PrismPool.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolPauseChangedIterator is returned from FilterPauseChanged and is used to iterate over the raw logs and unpacked data for PauseChanged events raised by the PrismPool contract.
type PrismPoolPauseChangedIterator struct {
	Event *PrismPoolPauseChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolPauseChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolPauseChanged)
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
		it.Event = new(PrismPoolPauseChanged)
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
func (it *PrismPoolPauseChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolPauseChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolPauseChanged represents a PauseChanged event raised by the PrismPool contract.
type PrismPoolPauseChanged struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPauseChanged is a free log retrieval operation binding the contract event 0x8fb6c181ee25a520cf3dd6565006ef91229fcfe5a989566c2a3b8c115570cec5.
//
// Solidity: event PauseChanged(bool paused)
func (_PrismPool *PrismPoolFilterer) FilterPauseChanged(opts *bind.FilterOpts) (*PrismPoolPauseChangedIterator, error) {

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "PauseChanged")
	if err != nil {
		return nil, err
	}
	return &PrismPoolPauseChangedIterator{contract: _PrismPool.contract, event: "PauseChanged", logs: logs, sub: sub}, nil
}

// WatchPauseChanged is a free log subscription operation binding the contract event 0x8fb6c181ee25a520cf3dd6565006ef91229fcfe5a989566c2a3b8c115570cec5.
//
// Solidity: event PauseChanged(bool paused)
func (_PrismPool *PrismPoolFilterer) WatchPauseChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolPauseChanged) (event.Subscription, error) {

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "PauseChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolPauseChanged)
				if err := _PrismPool.contract.UnpackLog(event, "PauseChanged", log); err != nil {
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

// ParsePauseChanged is a log parse operation binding the contract event 0x8fb6c181ee25a520cf3dd6565006ef91229fcfe5a989566c2a3b8c115570cec5.
//
// Solidity: event PauseChanged(bool paused)
func (_PrismPool *PrismPoolFilterer) ParsePauseChanged(log types.Log) (*PrismPoolPauseChanged, error) {
	event := new(PrismPoolPauseChanged)
	if err := _PrismPool.contract.UnpackLog(event, "PauseChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the PrismPool contract.
type PrismPoolPoolCreatedIterator struct {
	Event *PrismPoolPoolCreated // Event containing the contract specifics and raw log

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
func (it *PrismPoolPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolPoolCreated)
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
		it.Event = new(PrismPoolPoolCreated)
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
func (it *PrismPoolPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolPoolCreated represents a PoolCreated event raised by the PrismPool contract.
type PrismPoolPoolCreated struct {
	PoolId      *big.Int
	LendToken   common.Address
	BorrowToken common.Address
	SpToken     common.Address
	JpToken     common.Address
	SettleTime  *big.Int
	EndTime     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0xd11742d80e016020b8cc622078ca3f6b03c322fdd90d5879c731727654957cc4.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed lendToken, address indexed borrowToken, address spToken, address jpToken, uint256 settleTime, uint256 endTime)
func (_PrismPool *PrismPoolFilterer) FilterPoolCreated(opts *bind.FilterOpts, poolId []*big.Int, lendToken []common.Address, borrowToken []common.Address) (*PrismPoolPoolCreatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var lendTokenRule []interface{}
	for _, lendTokenItem := range lendToken {
		lendTokenRule = append(lendTokenRule, lendTokenItem)
	}
	var borrowTokenRule []interface{}
	for _, borrowTokenItem := range borrowToken {
		borrowTokenRule = append(borrowTokenRule, borrowTokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "PoolCreated", poolIdRule, lendTokenRule, borrowTokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolPoolCreatedIterator{contract: _PrismPool.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0xd11742d80e016020b8cc622078ca3f6b03c322fdd90d5879c731727654957cc4.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed lendToken, address indexed borrowToken, address spToken, address jpToken, uint256 settleTime, uint256 endTime)
func (_PrismPool *PrismPoolFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *PrismPoolPoolCreated, poolId []*big.Int, lendToken []common.Address, borrowToken []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var lendTokenRule []interface{}
	for _, lendTokenItem := range lendToken {
		lendTokenRule = append(lendTokenRule, lendTokenItem)
	}
	var borrowTokenRule []interface{}
	for _, borrowTokenItem := range borrowToken {
		borrowTokenRule = append(borrowTokenRule, borrowTokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "PoolCreated", poolIdRule, lendTokenRule, borrowTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolPoolCreated)
				if err := _PrismPool.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0xd11742d80e016020b8cc622078ca3f6b03c322fdd90d5879c731727654957cc4.
//
// Solidity: event PoolCreated(uint256 indexed poolId, address indexed lendToken, address indexed borrowToken, address spToken, address jpToken, uint256 settleTime, uint256 endTime)
func (_PrismPool *PrismPoolFilterer) ParsePoolCreated(log types.Log) (*PrismPoolPoolCreated, error) {
	event := new(PrismPoolPoolCreated)
	if err := _PrismPool.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolPoolLiquidatedIterator is returned from FilterPoolLiquidated and is used to iterate over the raw logs and unpacked data for PoolLiquidated events raised by the PrismPool contract.
type PrismPoolPoolLiquidatedIterator struct {
	Event *PrismPoolPoolLiquidated // Event containing the contract specifics and raw log

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
func (it *PrismPoolPoolLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolPoolLiquidated)
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
		it.Event = new(PrismPoolPoolLiquidated)
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
func (it *PrismPoolPoolLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolPoolLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolPoolLiquidated represents a PoolLiquidated event raised by the PrismPool contract.
type PrismPoolPoolLiquidated struct {
	PoolId                    *big.Int
	Router                    common.Address
	CollateralSold            *big.Int
	LendTokenRecovered        *big.Int
	RemainingCollateralAmount *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterPoolLiquidated is a free log retrieval operation binding the contract event 0x12cebae9b4631887732a49d48835cfc82f44ff34dc09d666127887f75ce832a6.
//
// Solidity: event PoolLiquidated(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 lendTokenRecovered, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) FilterPoolLiquidated(opts *bind.FilterOpts, poolId []*big.Int, router []common.Address) (*PrismPoolPoolLiquidatedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "PoolLiquidated", poolIdRule, routerRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolPoolLiquidatedIterator{contract: _PrismPool.contract, event: "PoolLiquidated", logs: logs, sub: sub}, nil
}

// WatchPoolLiquidated is a free log subscription operation binding the contract event 0x12cebae9b4631887732a49d48835cfc82f44ff34dc09d666127887f75ce832a6.
//
// Solidity: event PoolLiquidated(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 lendTokenRecovered, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) WatchPoolLiquidated(opts *bind.WatchOpts, sink chan<- *PrismPoolPoolLiquidated, poolId []*big.Int, router []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "PoolLiquidated", poolIdRule, routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolPoolLiquidated)
				if err := _PrismPool.contract.UnpackLog(event, "PoolLiquidated", log); err != nil {
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

// ParsePoolLiquidated is a log parse operation binding the contract event 0x12cebae9b4631887732a49d48835cfc82f44ff34dc09d666127887f75ce832a6.
//
// Solidity: event PoolLiquidated(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 lendTokenRecovered, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) ParsePoolLiquidated(log types.Log) (*PrismPoolPoolLiquidated, error) {
	event := new(PrismPoolPoolLiquidated)
	if err := _PrismPool.contract.UnpackLog(event, "PoolLiquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolPoolRepaidIterator is returned from FilterPoolRepaid and is used to iterate over the raw logs and unpacked data for PoolRepaid events raised by the PrismPool contract.
type PrismPoolPoolRepaidIterator struct {
	Event *PrismPoolPoolRepaid // Event containing the contract specifics and raw log

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
func (it *PrismPoolPoolRepaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolPoolRepaid)
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
		it.Event = new(PrismPoolPoolRepaid)
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
func (it *PrismPoolPoolRepaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolPoolRepaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolPoolRepaid represents a PoolRepaid event raised by the PrismPool contract.
type PrismPoolPoolRepaid struct {
	PoolId                    *big.Int
	Router                    common.Address
	CollateralSold            *big.Int
	RepaymentAmount           *big.Int
	RemainingCollateralAmount *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterPoolRepaid is a free log retrieval operation binding the contract event 0x6cde2ef5f275e21da9a00bf62648037748c7dee9d67c1b30099287af37fa2d66.
//
// Solidity: event PoolRepaid(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 repaymentAmount, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) FilterPoolRepaid(opts *bind.FilterOpts, poolId []*big.Int, router []common.Address) (*PrismPoolPoolRepaidIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "PoolRepaid", poolIdRule, routerRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolPoolRepaidIterator{contract: _PrismPool.contract, event: "PoolRepaid", logs: logs, sub: sub}, nil
}

// WatchPoolRepaid is a free log subscription operation binding the contract event 0x6cde2ef5f275e21da9a00bf62648037748c7dee9d67c1b30099287af37fa2d66.
//
// Solidity: event PoolRepaid(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 repaymentAmount, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) WatchPoolRepaid(opts *bind.WatchOpts, sink chan<- *PrismPoolPoolRepaid, poolId []*big.Int, router []common.Address) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var routerRule []interface{}
	for _, routerItem := range router {
		routerRule = append(routerRule, routerItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "PoolRepaid", poolIdRule, routerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolPoolRepaid)
				if err := _PrismPool.contract.UnpackLog(event, "PoolRepaid", log); err != nil {
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

// ParsePoolRepaid is a log parse operation binding the contract event 0x6cde2ef5f275e21da9a00bf62648037748c7dee9d67c1b30099287af37fa2d66.
//
// Solidity: event PoolRepaid(uint256 indexed poolId, address indexed router, uint256 collateralSold, uint256 repaymentAmount, uint256 remainingCollateralAmount)
func (_PrismPool *PrismPoolFilterer) ParsePoolRepaid(log types.Log) (*PrismPoolPoolRepaid, error) {
	event := new(PrismPoolPoolRepaid)
	if err := _PrismPool.contract.UnpackLog(event, "PoolRepaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolRefundBorrowIterator is returned from FilterRefundBorrow and is used to iterate over the raw logs and unpacked data for RefundBorrow events raised by the PrismPool contract.
type PrismPoolRefundBorrowIterator struct {
	Event *PrismPoolRefundBorrow // Event containing the contract specifics and raw log

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
func (it *PrismPoolRefundBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolRefundBorrow)
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
		it.Event = new(PrismPoolRefundBorrow)
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
func (it *PrismPoolRefundBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolRefundBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolRefundBorrow represents a RefundBorrow event raised by the PrismPool contract.
type PrismPoolRefundBorrow struct {
	Borrower common.Address
	PoolId   *big.Int
	Token    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRefundBorrow is a free log retrieval operation binding the contract event 0x2b55f629c76bbb6328ae918337679440ad5ca60d95f7193c6f43dd3aa9de8319.
//
// Solidity: event RefundBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) FilterRefundBorrow(opts *bind.FilterOpts, borrower []common.Address, poolId []*big.Int, token []common.Address) (*PrismPoolRefundBorrowIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "RefundBorrow", borrowerRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolRefundBorrowIterator{contract: _PrismPool.contract, event: "RefundBorrow", logs: logs, sub: sub}, nil
}

// WatchRefundBorrow is a free log subscription operation binding the contract event 0x2b55f629c76bbb6328ae918337679440ad5ca60d95f7193c6f43dd3aa9de8319.
//
// Solidity: event RefundBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) WatchRefundBorrow(opts *bind.WatchOpts, sink chan<- *PrismPoolRefundBorrow, borrower []common.Address, poolId []*big.Int, token []common.Address) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "RefundBorrow", borrowerRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolRefundBorrow)
				if err := _PrismPool.contract.UnpackLog(event, "RefundBorrow", log); err != nil {
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

// ParseRefundBorrow is a log parse operation binding the contract event 0x2b55f629c76bbb6328ae918337679440ad5ca60d95f7193c6f43dd3aa9de8319.
//
// Solidity: event RefundBorrow(address indexed borrower, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) ParseRefundBorrow(log types.Log) (*PrismPoolRefundBorrow, error) {
	event := new(PrismPoolRefundBorrow)
	if err := _PrismPool.contract.UnpackLog(event, "RefundBorrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolRefundLendIterator is returned from FilterRefundLend and is used to iterate over the raw logs and unpacked data for RefundLend events raised by the PrismPool contract.
type PrismPoolRefundLendIterator struct {
	Event *PrismPoolRefundLend // Event containing the contract specifics and raw log

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
func (it *PrismPoolRefundLendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolRefundLend)
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
		it.Event = new(PrismPoolRefundLend)
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
func (it *PrismPoolRefundLendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolRefundLendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolRefundLend represents a RefundLend event raised by the PrismPool contract.
type PrismPoolRefundLend struct {
	Lender common.Address
	PoolId *big.Int
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefundLend is a free log retrieval operation binding the contract event 0x16951a260cccf9377f46885646e0573078b0cb444b0670601031f2f34ff8557e.
//
// Solidity: event RefundLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) FilterRefundLend(opts *bind.FilterOpts, lender []common.Address, poolId []*big.Int, token []common.Address) (*PrismPoolRefundLendIterator, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "RefundLend", lenderRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolRefundLendIterator{contract: _PrismPool.contract, event: "RefundLend", logs: logs, sub: sub}, nil
}

// WatchRefundLend is a free log subscription operation binding the contract event 0x16951a260cccf9377f46885646e0573078b0cb444b0670601031f2f34ff8557e.
//
// Solidity: event RefundLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) WatchRefundLend(opts *bind.WatchOpts, sink chan<- *PrismPoolRefundLend, lender []common.Address, poolId []*big.Int, token []common.Address) (event.Subscription, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "RefundLend", lenderRule, poolIdRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolRefundLend)
				if err := _PrismPool.contract.UnpackLog(event, "RefundLend", log); err != nil {
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

// ParseRefundLend is a log parse operation binding the contract event 0x16951a260cccf9377f46885646e0573078b0cb444b0670601031f2f34ff8557e.
//
// Solidity: event RefundLend(address indexed lender, uint256 indexed poolId, address indexed token, uint256 amount)
func (_PrismPool *PrismPoolFilterer) ParseRefundLend(log types.Log) (*PrismPoolRefundLend, error) {
	event := new(PrismPoolRefundLend)
	if err := _PrismPool.contract.UnpackLog(event, "RefundLend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolStateChangedIterator is returned from FilterStateChanged and is used to iterate over the raw logs and unpacked data for StateChanged events raised by the PrismPool contract.
type PrismPoolStateChangedIterator struct {
	Event *PrismPoolStateChanged // Event containing the contract specifics and raw log

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
func (it *PrismPoolStateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolStateChanged)
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
		it.Event = new(PrismPoolStateChanged)
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
func (it *PrismPoolStateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolStateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolStateChanged represents a StateChanged event raised by the PrismPool contract.
type PrismPoolStateChanged struct {
	PoolId   *big.Int
	OldState uint8
	NewState uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStateChanged is a free log retrieval operation binding the contract event 0x18e785d11512625ba8e5486b7d963f337dfc0a628752fe18e284b0375eefe46c.
//
// Solidity: event StateChanged(uint256 indexed poolId, uint8 oldState, uint8 newState)
func (_PrismPool *PrismPoolFilterer) FilterStateChanged(opts *bind.FilterOpts, poolId []*big.Int) (*PrismPoolStateChangedIterator, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "StateChanged", poolIdRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolStateChangedIterator{contract: _PrismPool.contract, event: "StateChanged", logs: logs, sub: sub}, nil
}

// WatchStateChanged is a free log subscription operation binding the contract event 0x18e785d11512625ba8e5486b7d963f337dfc0a628752fe18e284b0375eefe46c.
//
// Solidity: event StateChanged(uint256 indexed poolId, uint8 oldState, uint8 newState)
func (_PrismPool *PrismPoolFilterer) WatchStateChanged(opts *bind.WatchOpts, sink chan<- *PrismPoolStateChanged, poolId []*big.Int) (event.Subscription, error) {

	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "StateChanged", poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolStateChanged)
				if err := _PrismPool.contract.UnpackLog(event, "StateChanged", log); err != nil {
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

// ParseStateChanged is a log parse operation binding the contract event 0x18e785d11512625ba8e5486b7d963f337dfc0a628752fe18e284b0375eefe46c.
//
// Solidity: event StateChanged(uint256 indexed poolId, uint8 oldState, uint8 newState)
func (_PrismPool *PrismPoolFilterer) ParseStateChanged(log types.Log) (*PrismPoolStateChanged, error) {
	event := new(PrismPoolStateChanged)
	if err := _PrismPool.contract.UnpackLog(event, "StateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolWithdrawBorrowIterator is returned from FilterWithdrawBorrow and is used to iterate over the raw logs and unpacked data for WithdrawBorrow events raised by the PrismPool contract.
type PrismPoolWithdrawBorrowIterator struct {
	Event *PrismPoolWithdrawBorrow // Event containing the contract specifics and raw log

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
func (it *PrismPoolWithdrawBorrowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolWithdrawBorrow)
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
		it.Event = new(PrismPoolWithdrawBorrow)
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
func (it *PrismPoolWithdrawBorrowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolWithdrawBorrowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolWithdrawBorrow represents a WithdrawBorrow event raised by the PrismPool contract.
type PrismPoolWithdrawBorrow struct {
	Borrower         common.Address
	PoolId           *big.Int
	JpAmount         *big.Int
	CollateralAmount *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterWithdrawBorrow is a free log retrieval operation binding the contract event 0xfc146a6bd65d1e7079b55c3d94caab9e851e63b1dcf74726cc276c8643682e1a.
//
// Solidity: event WithdrawBorrow(address indexed borrower, uint256 indexed poolId, uint256 jpAmount, uint256 collateralAmount)
func (_PrismPool *PrismPoolFilterer) FilterWithdrawBorrow(opts *bind.FilterOpts, borrower []common.Address, poolId []*big.Int) (*PrismPoolWithdrawBorrowIterator, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "WithdrawBorrow", borrowerRule, poolIdRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolWithdrawBorrowIterator{contract: _PrismPool.contract, event: "WithdrawBorrow", logs: logs, sub: sub}, nil
}

// WatchWithdrawBorrow is a free log subscription operation binding the contract event 0xfc146a6bd65d1e7079b55c3d94caab9e851e63b1dcf74726cc276c8643682e1a.
//
// Solidity: event WithdrawBorrow(address indexed borrower, uint256 indexed poolId, uint256 jpAmount, uint256 collateralAmount)
func (_PrismPool *PrismPoolFilterer) WatchWithdrawBorrow(opts *bind.WatchOpts, sink chan<- *PrismPoolWithdrawBorrow, borrower []common.Address, poolId []*big.Int) (event.Subscription, error) {

	var borrowerRule []interface{}
	for _, borrowerItem := range borrower {
		borrowerRule = append(borrowerRule, borrowerItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "WithdrawBorrow", borrowerRule, poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolWithdrawBorrow)
				if err := _PrismPool.contract.UnpackLog(event, "WithdrawBorrow", log); err != nil {
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

// ParseWithdrawBorrow is a log parse operation binding the contract event 0xfc146a6bd65d1e7079b55c3d94caab9e851e63b1dcf74726cc276c8643682e1a.
//
// Solidity: event WithdrawBorrow(address indexed borrower, uint256 indexed poolId, uint256 jpAmount, uint256 collateralAmount)
func (_PrismPool *PrismPoolFilterer) ParseWithdrawBorrow(log types.Log) (*PrismPoolWithdrawBorrow, error) {
	event := new(PrismPoolWithdrawBorrow)
	if err := _PrismPool.contract.UnpackLog(event, "WithdrawBorrow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PrismPoolWithdrawLendIterator is returned from FilterWithdrawLend and is used to iterate over the raw logs and unpacked data for WithdrawLend events raised by the PrismPool contract.
type PrismPoolWithdrawLendIterator struct {
	Event *PrismPoolWithdrawLend // Event containing the contract specifics and raw log

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
func (it *PrismPoolWithdrawLendIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PrismPoolWithdrawLend)
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
		it.Event = new(PrismPoolWithdrawLend)
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
func (it *PrismPoolWithdrawLendIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PrismPoolWithdrawLendIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PrismPoolWithdrawLend represents a WithdrawLend event raised by the PrismPool contract.
type PrismPoolWithdrawLend struct {
	Lender     common.Address
	PoolId     *big.Int
	SpAmount   *big.Int
	LendAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterWithdrawLend is a free log retrieval operation binding the contract event 0x04123ae46ef5abb125e8c9bfd72c4e8cf45f10facfc263b46d4b92a896491019.
//
// Solidity: event WithdrawLend(address indexed lender, uint256 indexed poolId, uint256 spAmount, uint256 lendAmount)
func (_PrismPool *PrismPoolFilterer) FilterWithdrawLend(opts *bind.FilterOpts, lender []common.Address, poolId []*big.Int) (*PrismPoolWithdrawLendIterator, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.FilterLogs(opts, "WithdrawLend", lenderRule, poolIdRule)
	if err != nil {
		return nil, err
	}
	return &PrismPoolWithdrawLendIterator{contract: _PrismPool.contract, event: "WithdrawLend", logs: logs, sub: sub}, nil
}

// WatchWithdrawLend is a free log subscription operation binding the contract event 0x04123ae46ef5abb125e8c9bfd72c4e8cf45f10facfc263b46d4b92a896491019.
//
// Solidity: event WithdrawLend(address indexed lender, uint256 indexed poolId, uint256 spAmount, uint256 lendAmount)
func (_PrismPool *PrismPoolFilterer) WatchWithdrawLend(opts *bind.WatchOpts, sink chan<- *PrismPoolWithdrawLend, lender []common.Address, poolId []*big.Int) (event.Subscription, error) {

	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}

	logs, sub, err := _PrismPool.contract.WatchLogs(opts, "WithdrawLend", lenderRule, poolIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PrismPoolWithdrawLend)
				if err := _PrismPool.contract.UnpackLog(event, "WithdrawLend", log); err != nil {
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

// ParseWithdrawLend is a log parse operation binding the contract event 0x04123ae46ef5abb125e8c9bfd72c4e8cf45f10facfc263b46d4b92a896491019.
//
// Solidity: event WithdrawLend(address indexed lender, uint256 indexed poolId, uint256 spAmount, uint256 lendAmount)
func (_PrismPool *PrismPoolFilterer) ParseWithdrawLend(log types.Log) (*PrismPoolWithdrawLend, error) {
	event := new(PrismPoolWithdrawLend)
	if err := _PrismPool.contract.UnpackLog(event, "WithdrawLend", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
