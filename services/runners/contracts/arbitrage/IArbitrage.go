// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrage

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ArbitrageABI is the input ABI used to generate the binding from.
const ArbitrageABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router2\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token2\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount2\",\"type\":\"uint256\"}],\"name\":\"startArbitrage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// Arbitrage is an auto generated Go binding around an Ethereum contract.
type Arbitrage struct {
	ArbitrageCaller     // Read-only binding to the contract
	ArbitrageTransactor // Write-only binding to the contract
	ArbitrageFilterer   // Log filterer for contract events
}

// ArbitrageCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArbitrageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArbitrageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArbitrageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbitrageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArbitrageSession struct {
	Contract     *Arbitrage        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArbitrageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArbitrageCallerSession struct {
	Contract *ArbitrageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ArbitrageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArbitrageTransactorSession struct {
	Contract     *ArbitrageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ArbitrageRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArbitrageRaw struct {
	Contract *Arbitrage // Generic contract binding to access the raw methods on
}

// ArbitrageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArbitrageCallerRaw struct {
	Contract *ArbitrageCaller // Generic read-only contract binding to access the raw methods on
}

// ArbitrageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArbitrageTransactorRaw struct {
	Contract *ArbitrageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArbitrage creates a new instance of Arbitrage, bound to a specific deployed contract.
func NewArbitrage(address common.Address, backend bind.ContractBackend) (*Arbitrage, error) {
	contract, err := bindArbitrage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Arbitrage{ArbitrageCaller: ArbitrageCaller{contract: contract}, ArbitrageTransactor: ArbitrageTransactor{contract: contract}, ArbitrageFilterer: ArbitrageFilterer{contract: contract}}, nil
}

// NewArbitrageCaller creates a new read-only instance of Arbitrage, bound to a specific deployed contract.
func NewArbitrageCaller(address common.Address, caller bind.ContractCaller) (*ArbitrageCaller, error) {
	contract, err := bindArbitrage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrageCaller{contract: contract}, nil
}

// NewArbitrageTransactor creates a new write-only instance of Arbitrage, bound to a specific deployed contract.
func NewArbitrageTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrageTransactor, error) {
	contract, err := bindArbitrage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrageTransactor{contract: contract}, nil
}

// NewArbitrageFilterer creates a new log filterer instance of Arbitrage, bound to a specific deployed contract.
func NewArbitrageFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrageFilterer, error) {
	contract, err := bindArbitrage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrageFilterer{contract: contract}, nil
}

// bindArbitrage binds a generic wrapper to an already deployed contract.
func bindArbitrage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArbitrageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arbitrage *ArbitrageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Arbitrage.Contract.ArbitrageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arbitrage *ArbitrageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arbitrage.Contract.ArbitrageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arbitrage *ArbitrageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arbitrage.Contract.ArbitrageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arbitrage *ArbitrageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Arbitrage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arbitrage *ArbitrageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arbitrage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arbitrage *ArbitrageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arbitrage.Contract.contract.Transact(opts, method, params...)
}

// StartArbitrage is a paid mutator transaction binding the contract method 0xbe493263.
//
// Solidity: function startArbitrage(address router1, address router2, address token0, address token1, address token2, uint256 amount0, uint256 amount1, uint256 amount2) payable returns()
func (_Arbitrage *ArbitrageTransactor) StartArbitrage(opts *bind.TransactOpts, router1 common.Address, router2 common.Address, token0 common.Address, token1 common.Address, token2 common.Address, amount0 *big.Int, amount1 *big.Int, amount2 *big.Int) (*types.Transaction, error) {
	return _Arbitrage.contract.Transact(opts, "startArbitrage", router1, router2, token0, token1, token2, amount0, amount1, amount2)
}

// StartArbitrage is a paid mutator transaction binding the contract method 0xbe493263.
//
// Solidity: function startArbitrage(address router1, address router2, address token0, address token1, address token2, uint256 amount0, uint256 amount1, uint256 amount2) payable returns()
func (_Arbitrage *ArbitrageSession) StartArbitrage(router1 common.Address, router2 common.Address, token0 common.Address, token1 common.Address, token2 common.Address, amount0 *big.Int, amount1 *big.Int, amount2 *big.Int) (*types.Transaction, error) {
	return _Arbitrage.Contract.StartArbitrage(&_Arbitrage.TransactOpts, router1, router2, token0, token1, token2, amount0, amount1, amount2)
}

// StartArbitrage is a paid mutator transaction binding the contract method 0xbe493263.
//
// Solidity: function startArbitrage(address router1, address router2, address token0, address token1, address token2, uint256 amount0, uint256 amount1, uint256 amount2) payable returns()
func (_Arbitrage *ArbitrageTransactorSession) StartArbitrage(router1 common.Address, router2 common.Address, token0 common.Address, token1 common.Address, token2 common.Address, amount0 *big.Int, amount1 *big.Int, amount2 *big.Int) (*types.Transaction, error) {
	return _Arbitrage.Contract.StartArbitrage(&_Arbitrage.TransactOpts, router1, router2, token0, token1, token2, amount0, amount1, amount2)
}
