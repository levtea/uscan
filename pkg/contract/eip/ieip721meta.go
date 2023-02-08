// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eip

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
)

// Ieip721metaMetaData contains all meta data concerning the Ieip721meta contract.
var Ieip721metaMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Ieip721metaABI is the input ABI used to generate the binding from.
// Deprecated: Use Ieip721metaMetaData.ABI instead.
var Ieip721metaABI = Ieip721metaMetaData.ABI

// Ieip721meta is an auto generated Go binding around an Ethereum contract.
type Ieip721meta struct {
	Ieip721metaCaller     // Read-only binding to the contract
	Ieip721metaTransactor // Write-only binding to the contract
	Ieip721metaFilterer   // Log filterer for contract events
}

// Ieip721metaCaller is an auto generated read-only Go binding around an Ethereum contract.
type Ieip721metaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ieip721metaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Ieip721metaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ieip721metaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Ieip721metaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Ieip721metaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Ieip721metaSession struct {
	Contract     *Ieip721meta      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Ieip721metaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Ieip721metaCallerSession struct {
	Contract *Ieip721metaCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// Ieip721metaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Ieip721metaTransactorSession struct {
	Contract     *Ieip721metaTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// Ieip721metaRaw is an auto generated low-level Go binding around an Ethereum contract.
type Ieip721metaRaw struct {
	Contract *Ieip721meta // Generic contract binding to access the raw methods on
}

// Ieip721metaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Ieip721metaCallerRaw struct {
	Contract *Ieip721metaCaller // Generic read-only contract binding to access the raw methods on
}

// Ieip721metaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Ieip721metaTransactorRaw struct {
	Contract *Ieip721metaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIeip721meta creates a new instance of Ieip721meta, bound to a specific deployed contract.
func NewIeip721meta(address common.Address, backend bind.ContractBackend) (*Ieip721meta, error) {
	contract, err := bindIeip721meta(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ieip721meta{Ieip721metaCaller: Ieip721metaCaller{contract: contract}, Ieip721metaTransactor: Ieip721metaTransactor{contract: contract}, Ieip721metaFilterer: Ieip721metaFilterer{contract: contract}}, nil
}

// NewIeip721metaCaller creates a new read-only instance of Ieip721meta, bound to a specific deployed contract.
func NewIeip721metaCaller(address common.Address, caller bind.ContractCaller) (*Ieip721metaCaller, error) {
	contract, err := bindIeip721meta(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Ieip721metaCaller{contract: contract}, nil
}

// NewIeip721metaTransactor creates a new write-only instance of Ieip721meta, bound to a specific deployed contract.
func NewIeip721metaTransactor(address common.Address, transactor bind.ContractTransactor) (*Ieip721metaTransactor, error) {
	contract, err := bindIeip721meta(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Ieip721metaTransactor{contract: contract}, nil
}

// NewIeip721metaFilterer creates a new log filterer instance of Ieip721meta, bound to a specific deployed contract.
func NewIeip721metaFilterer(address common.Address, filterer bind.ContractFilterer) (*Ieip721metaFilterer, error) {
	contract, err := bindIeip721meta(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Ieip721metaFilterer{contract: contract}, nil
}

// bindIeip721meta binds a generic wrapper to an already deployed contract.
func bindIeip721meta(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Ieip721metaABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ieip721meta *Ieip721metaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ieip721meta.Contract.Ieip721metaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ieip721meta *Ieip721metaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ieip721meta.Contract.Ieip721metaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ieip721meta *Ieip721metaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ieip721meta.Contract.Ieip721metaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ieip721meta *Ieip721metaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ieip721meta.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ieip721meta *Ieip721metaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ieip721meta.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ieip721meta *Ieip721metaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ieip721meta.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ieip721meta *Ieip721metaCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ieip721meta.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ieip721meta *Ieip721metaSession) Name() (string, error) {
	return _Ieip721meta.Contract.Name(&_Ieip721meta.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Ieip721meta *Ieip721metaCallerSession) Name() (string, error) {
	return _Ieip721meta.Contract.Name(&_Ieip721meta.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ieip721meta *Ieip721metaCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ieip721meta.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ieip721meta *Ieip721metaSession) Symbol() (string, error) {
	return _Ieip721meta.Contract.Symbol(&_Ieip721meta.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Ieip721meta *Ieip721metaCallerSession) Symbol() (string, error) {
	return _Ieip721meta.Contract.Symbol(&_Ieip721meta.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ieip721meta *Ieip721metaCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Ieip721meta.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ieip721meta *Ieip721metaSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Ieip721meta.Contract.TokenURI(&_Ieip721meta.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Ieip721meta *Ieip721metaCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Ieip721meta.Contract.TokenURI(&_Ieip721meta.CallOpts, tokenId)
}
