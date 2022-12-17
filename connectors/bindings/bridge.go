// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// RskBridgeMetaData contains all meta data concerning the RskBridge contract.
var RskBridgeMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"type\":\"function\",\"constant\":true,\"inputs\":[{\"name\":\"btcBlockHeight\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}]},{\"name\":\"registerFastBridgeBtcTransaction\",\"type\":\"function\",\"constant\":true,\"inputs\":[{\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\"},{\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\"},{\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\"},{\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"name\":\"shouldTransferToContract\",\"type\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}]},{\"name\":\"getActiveFederationCreationBlockHeight\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}]},{\"name\":\"getFederatorPublicKeyOfType\",\"type\":\"function\",\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}]},{\"name\":\"getFederatorPublicKey\",\"type\":\"function\",\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}]},{\"name\":\"getFederationThreshold\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}]},{\"name\":\"getFederationSize\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}]},{\"name\":\"getFederationAddress\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\"}]},{\"name\":\"getMinimumLockTxValue\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}]},{\"name\":\"getActivePowpegRedeemScript\",\"type\":\"function\",\"constant\":true,\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}]}]",
}

// RskBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use RskBridgeMetaData.ABI instead.
var RskBridgeABI = RskBridgeMetaData.ABI

// RskBridge is an auto generated Go binding around an Ethereum contract.
type RskBridge struct {
	RskBridgeCaller     // Read-only binding to the contract
	RskBridgeTransactor // Write-only binding to the contract
	RskBridgeFilterer   // Log filterer for contract events
}

// RskBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type RskBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RskBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RskBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RskBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RskBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RskBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RskBridgeSession struct {
	Contract     *RskBridge        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RskBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RskBridgeCallerSession struct {
	Contract *RskBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RskBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RskBridgeTransactorSession struct {
	Contract     *RskBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RskBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type RskBridgeRaw struct {
	Contract *RskBridge // Generic contract binding to access the raw methods on
}

// RskBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RskBridgeCallerRaw struct {
	Contract *RskBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// RskBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RskBridgeTransactorRaw struct {
	Contract *RskBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRskBridge creates a new instance of RskBridge, bound to a specific deployed contract.
func NewRskBridge(address common.Address, backend bind.ContractBackend) (*RskBridge, error) {
	contract, err := bindRskBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RskBridge{RskBridgeCaller: RskBridgeCaller{contract: contract}, RskBridgeTransactor: RskBridgeTransactor{contract: contract}, RskBridgeFilterer: RskBridgeFilterer{contract: contract}}, nil
}

// NewRskBridgeCaller creates a new read-only instance of RskBridge, bound to a specific deployed contract.
func NewRskBridgeCaller(address common.Address, caller bind.ContractCaller) (*RskBridgeCaller, error) {
	contract, err := bindRskBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RskBridgeCaller{contract: contract}, nil
}

// NewRskBridgeTransactor creates a new write-only instance of RskBridge, bound to a specific deployed contract.
func NewRskBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*RskBridgeTransactor, error) {
	contract, err := bindRskBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RskBridgeTransactor{contract: contract}, nil
}

// NewRskBridgeFilterer creates a new log filterer instance of RskBridge, bound to a specific deployed contract.
func NewRskBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*RskBridgeFilterer, error) {
	contract, err := bindRskBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RskBridgeFilterer{contract: contract}, nil
}

// bindRskBridge binds a generic wrapper to an already deployed contract.
func bindRskBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RskBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RskBridge *RskBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RskBridge.Contract.RskBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RskBridge *RskBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.Contract.RskBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RskBridge *RskBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RskBridge.Contract.RskBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RskBridge *RskBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RskBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RskBridge *RskBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RskBridge *RskBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RskBridge.Contract.contract.Transact(opts, method, params...)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() returns(uint256)
func (_RskBridge *RskBridgeCaller) GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getActiveFederationCreationBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() returns(uint256)
func (_RskBridge *RskBridgeSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetActiveFederationCreationBlockHeight(&_RskBridge.CallOpts)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() returns(uint256)
func (_RskBridge *RskBridgeCallerSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetActiveFederationCreationBlockHeight(&_RskBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() returns(bytes)
func (_RskBridge *RskBridgeCaller) GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getActivePowpegRedeemScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() returns(bytes)
func (_RskBridge *RskBridgeSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _RskBridge.Contract.GetActivePowpegRedeemScript(&_RskBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _RskBridge.Contract.GetActivePowpegRedeemScript(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) returns(bytes)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainBlockHeaderByHeight(opts *bind.CallOpts, btcBlockHeight *big.Int) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHeight", btcBlockHeight)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_RskBridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_RskBridge.CallOpts, btcBlockHeight)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() returns(string)
func (_RskBridge *RskBridgeCaller) GetFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() returns(string)
func (_RskBridge *RskBridgeSession) GetFederationAddress() (string, error) {
	return _RskBridge.Contract.GetFederationAddress(&_RskBridge.CallOpts)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() returns(string)
func (_RskBridge *RskBridgeCallerSession) GetFederationAddress() (string, error) {
	return _RskBridge.Contract.GetFederationAddress(&_RskBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() returns(int256)
func (_RskBridge *RskBridgeCaller) GetFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationSize(&_RskBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationSize(&_RskBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() returns(int256)
func (_RskBridge *RskBridgeCaller) GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationThreshold(&_RskBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationThreshold(&_RskBridge.CallOpts)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) returns(bytes)
func (_RskBridge *RskBridgeCaller) GetFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) returns(bytes)
func (_RskBridge *RskBridgeSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string ) returns(bytes)
func (_RskBridge *RskBridgeCaller) GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, arg1 string) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederatorPublicKeyOfType", index, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string ) returns(bytes)
func (_RskBridge *RskBridgeSession) GetFederatorPublicKeyOfType(index *big.Int, arg1 string) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, arg1)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string ) returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetFederatorPublicKeyOfType(index *big.Int, arg1 string) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, arg1)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() returns(int256)
func (_RskBridge *RskBridgeCaller) GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getMinimumLockTxValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() returns(int256)
func (_RskBridge *RskBridgeSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _RskBridge.Contract.GetMinimumLockTxValue(&_RskBridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _RskBridge.Contract.GetMinimumLockTxValue(&_RskBridge.CallOpts)
}

// RegisterFastBridgeBtcTransaction is a free data retrieval call binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeCaller) RegisterFastBridgeBtcTransaction(opts *bind.CallOpts, btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RegisterFastBridgeBtcTransaction is a free data retrieval call binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*big.Int, error) {
	return _RskBridge.Contract.RegisterFastBridgeBtcTransaction(&_RskBridge.CallOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a free data retrieval call binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeCallerSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*big.Int, error) {
	return _RskBridge.Contract.RegisterFastBridgeBtcTransaction(&_RskBridge.CallOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}
