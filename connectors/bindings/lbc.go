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

// LiquidityBridgeContractQuote is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractQuote struct {
	FedBtcAddress               [20]byte
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	BtcRefundAddress            []byte
	RskRefundAddress            common.Address
	LiquidityProviderBtcAddress []byte
	CallFee                     uint64
	PenaltyFee                  uint64
	ContractAddress             common.Address
	Data                        []byte
	GasLimit                    uint32
	Nonce                       int64
	Value                       uint64
	AgreementTimestamp          uint32
	TimeForDeposit              uint32
	CallTime                    uint32
	DepositConfirmations        uint16
	CallOnRegister              bool
}

// LBCMetaData contains all meta data concerning the LBC contract.
var LBCMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"minimumCollateral\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"rewardPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"resignDelayBlocks\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"dustThreshold\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"callFee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDustThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"callFee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"callFee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// LBCABI is the input ABI used to generate the binding from.
// Deprecated: Use LBCMetaData.ABI instead.
var LBCABI = LBCMetaData.ABI

// LBC is an auto generated Go binding around an Ethereum contract.
type LBC struct {
	LBCCaller     // Read-only binding to the contract
	LBCTransactor // Write-only binding to the contract
	LBCFilterer   // Log filterer for contract events
}

// LBCCaller is an auto generated read-only Go binding around an Ethereum contract.
type LBCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LBCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LBCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LBCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LBCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LBCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LBCSession struct {
	Contract     *LBC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LBCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LBCCallerSession struct {
	Contract *LBCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LBCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LBCTransactorSession struct {
	Contract     *LBCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LBCRaw is an auto generated low-level Go binding around an Ethereum contract.
type LBCRaw struct {
	Contract *LBC // Generic contract binding to access the raw methods on
}

// LBCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LBCCallerRaw struct {
	Contract *LBCCaller // Generic read-only contract binding to access the raw methods on
}

// LBCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LBCTransactorRaw struct {
	Contract *LBCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLBC creates a new instance of LBC, bound to a specific deployed contract.
func NewLBC(address common.Address, backend bind.ContractBackend) (*LBC, error) {
	contract, err := bindLBC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LBC{LBCCaller: LBCCaller{contract: contract}, LBCTransactor: LBCTransactor{contract: contract}, LBCFilterer: LBCFilterer{contract: contract}}, nil
}

// NewLBCCaller creates a new read-only instance of LBC, bound to a specific deployed contract.
func NewLBCCaller(address common.Address, caller bind.ContractCaller) (*LBCCaller, error) {
	contract, err := bindLBC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LBCCaller{contract: contract}, nil
}

// NewLBCTransactor creates a new write-only instance of LBC, bound to a specific deployed contract.
func NewLBCTransactor(address common.Address, transactor bind.ContractTransactor) (*LBCTransactor, error) {
	contract, err := bindLBC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LBCTransactor{contract: contract}, nil
}

// NewLBCFilterer creates a new log filterer instance of LBC, bound to a specific deployed contract.
func NewLBCFilterer(address common.Address, filterer bind.ContractFilterer) (*LBCFilterer, error) {
	contract, err := bindLBC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LBCFilterer{contract: contract}, nil
}

// bindLBC binds a generic wrapper to an already deployed contract.
func bindLBC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LBCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LBC *LBCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LBC.Contract.LBCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LBC *LBCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.Contract.LBCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LBC *LBCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LBC.Contract.LBCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LBC *LBCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LBC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LBC *LBCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LBC *LBCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LBC.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LBC *LBCCaller) GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LBC *LBCSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LBC.Contract.GetBalance(&_LBC.CallOpts, addr)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LBC *LBCCallerSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LBC.Contract.GetBalance(&_LBC.CallOpts, addr)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LBC *LBCCaller) GetBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LBC *LBCSession) GetBridgeAddress() (common.Address, error) {
	return _LBC.Contract.GetBridgeAddress(&_LBC.CallOpts)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LBC *LBCCallerSession) GetBridgeAddress() (common.Address, error) {
	return _LBC.Contract.GetBridgeAddress(&_LBC.CallOpts)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LBC *LBCCaller) GetCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LBC *LBCSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LBC.Contract.GetCollateral(&_LBC.CallOpts, addr)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LBC *LBCCallerSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LBC.Contract.GetCollateral(&_LBC.CallOpts, addr)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(int256)
func (_LBC *LBCCaller) GetDustThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getDustThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(int256)
func (_LBC *LBCSession) GetDustThreshold() (*big.Int, error) {
	return _LBC.Contract.GetDustThreshold(&_LBC.CallOpts)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(int256)
func (_LBC *LBCCallerSession) GetDustThreshold() (*big.Int, error) {
	return _LBC.Contract.GetDustThreshold(&_LBC.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LBC *LBCCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LBC *LBCSession) GetMinCollateral() (*big.Int, error) {
	return _LBC.Contract.GetMinCollateral(&_LBC.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LBC *LBCCallerSession) GetMinCollateral() (*big.Int, error) {
	return _LBC.Contract.GetMinCollateral(&_LBC.CallOpts)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LBC *LBCCaller) GetResignDelayBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getResignDelayBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LBC *LBCSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LBC.Contract.GetResignDelayBlocks(&_LBC.CallOpts)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LBC *LBCCallerSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LBC.Contract.GetResignDelayBlocks(&_LBC.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LBC *LBCCaller) GetRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "getRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LBC *LBCSession) GetRewardPercentage() (*big.Int, error) {
	return _LBC.Contract.GetRewardPercentage(&_LBC.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LBC *LBCCallerSession) GetRewardPercentage() (*big.Int, error) {
	return _LBC.Contract.GetRewardPercentage(&_LBC.CallOpts)
}

// HashQuote is a free data retrieval call binding the contract method 0x71c6e1de.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) pure returns(bytes32)
func (_LBC *LBCCaller) HashQuote(opts *bind.CallOpts, quote LiquidityBridgeContractQuote) ([32]byte, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "hashQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashQuote is a free data retrieval call binding the contract method 0x71c6e1de.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) pure returns(bytes32)
func (_LBC *LBCSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LBC.Contract.HashQuote(&_LBC.CallOpts, quote)
}

// HashQuote is a free data retrieval call binding the contract method 0x71c6e1de.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) pure returns(bytes32)
func (_LBC *LBCCallerSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LBC.Contract.HashQuote(&_LBC.CallOpts, quote)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LBC *LBCCaller) IsOperational(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _LBC.contract.Call(opts, &out, "isOperational", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LBC *LBCSession) IsOperational(addr common.Address) (bool, error) {
	return _LBC.Contract.IsOperational(&_LBC.CallOpts, addr)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LBC *LBCCallerSession) IsOperational(addr common.Address) (bool, error) {
	return _LBC.Contract.IsOperational(&_LBC.CallOpts, addr)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LBC *LBCTransactor) AddCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "addCollateral")
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LBC *LBCSession) AddCollateral() (*types.Transaction, error) {
	return _LBC.Contract.AddCollateral(&_LBC.TransactOpts)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LBC *LBCTransactorSession) AddCollateral() (*types.Transaction, error) {
	return _LBC.Contract.AddCollateral(&_LBC.TransactOpts)
}

// CallForUser is a paid mutator transaction binding the contract method 0xe291c81e.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LBC *LBCTransactor) CallForUser(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "callForUser", quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xe291c81e.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LBC *LBCSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LBC.Contract.CallForUser(&_LBC.TransactOpts, quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xe291c81e.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LBC *LBCTransactorSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LBC.Contract.CallForUser(&_LBC.TransactOpts, quote)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LBC *LBCTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LBC *LBCSession) Deposit() (*types.Transaction, error) {
	return _LBC.Contract.Deposit(&_LBC.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LBC *LBCTransactorSession) Deposit() (*types.Transaction, error) {
	return _LBC.Contract.Deposit(&_LBC.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LBC *LBCTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LBC *LBCSession) Register() (*types.Transaction, error) {
	return _LBC.Contract.Register(&_LBC.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LBC *LBCTransactorSession) Register() (*types.Transaction, error) {
	return _LBC.Contract.Register(&_LBC.TransactOpts)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0xf4318e56.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LBC *LBCTransactor) RegisterPegIn(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0xf4318e56.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LBC *LBCSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LBC.Contract.RegisterPegIn(&_LBC.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0xf4318e56.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint64,uint64,address,bytes,uint32,int64,uint64,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LBC *LBCTransactorSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LBC.Contract.RegisterPegIn(&_LBC.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LBC *LBCTransactor) Resign(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "resign")
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LBC *LBCSession) Resign() (*types.Transaction, error) {
	return _LBC.Contract.Resign(&_LBC.TransactOpts)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LBC *LBCTransactorSession) Resign() (*types.Transaction, error) {
	return _LBC.Contract.Resign(&_LBC.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LBC *LBCTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LBC *LBCSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LBC.Contract.Withdraw(&_LBC.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LBC *LBCTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LBC.Contract.Withdraw(&_LBC.TransactOpts, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LBC *LBCTransactor) WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.Transact(opts, "withdrawCollateral")
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LBC *LBCSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LBC.Contract.WithdrawCollateral(&_LBC.TransactOpts)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LBC *LBCTransactorSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LBC.Contract.WithdrawCollateral(&_LBC.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LBC *LBCTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LBC.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LBC *LBCSession) Receive() (*types.Transaction, error) {
	return _LBC.Contract.Receive(&_LBC.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LBC *LBCTransactorSession) Receive() (*types.Transaction, error) {
	return _LBC.Contract.Receive(&_LBC.TransactOpts)
}

// LBCBalanceDecreaseIterator is returned from FilterBalanceDecrease and is used to iterate over the raw logs and unpacked data for BalanceDecrease events raised by the LBC contract.
type LBCBalanceDecreaseIterator struct {
	Event *LBCBalanceDecrease // Event containing the contract specifics and raw log

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
func (it *LBCBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCBalanceDecrease)
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
		it.Event = new(LBCBalanceDecrease)
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
func (it *LBCBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCBalanceDecrease represents a BalanceDecrease event raised by the LBC contract.
type LBCBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceDecrease is a free log retrieval operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) FilterBalanceDecrease(opts *bind.FilterOpts) (*LBCBalanceDecreaseIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return &LBCBalanceDecreaseIterator{contract: _LBC.contract, event: "BalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchBalanceDecrease is a free log subscription operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) WatchBalanceDecrease(opts *bind.WatchOpts, sink chan<- *LBCBalanceDecrease) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCBalanceDecrease)
				if err := _LBC.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
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

// ParseBalanceDecrease is a log parse operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) ParseBalanceDecrease(log types.Log) (*LBCBalanceDecrease, error) {
	event := new(LBCBalanceDecrease)
	if err := _LBC.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCBalanceIncreaseIterator is returned from FilterBalanceIncrease and is used to iterate over the raw logs and unpacked data for BalanceIncrease events raised by the LBC contract.
type LBCBalanceIncreaseIterator struct {
	Event *LBCBalanceIncrease // Event containing the contract specifics and raw log

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
func (it *LBCBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCBalanceIncrease)
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
		it.Event = new(LBCBalanceIncrease)
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
func (it *LBCBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCBalanceIncrease represents a BalanceIncrease event raised by the LBC contract.
type LBCBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceIncrease is a free log retrieval operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) FilterBalanceIncrease(opts *bind.FilterOpts) (*LBCBalanceIncreaseIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return &LBCBalanceIncreaseIterator{contract: _LBC.contract, event: "BalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchBalanceIncrease is a free log subscription operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) WatchBalanceIncrease(opts *bind.WatchOpts, sink chan<- *LBCBalanceIncrease) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCBalanceIncrease)
				if err := _LBC.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
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

// ParseBalanceIncrease is a log parse operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LBC *LBCFilterer) ParseBalanceIncrease(log types.Log) (*LBCBalanceIncrease, error) {
	event := new(LBCBalanceIncrease)
	if err := _LBC.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCBridgeCapExceededIterator is returned from FilterBridgeCapExceeded and is used to iterate over the raw logs and unpacked data for BridgeCapExceeded events raised by the LBC contract.
type LBCBridgeCapExceededIterator struct {
	Event *LBCBridgeCapExceeded // Event containing the contract specifics and raw log

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
func (it *LBCBridgeCapExceededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCBridgeCapExceeded)
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
		it.Event = new(LBCBridgeCapExceeded)
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
func (it *LBCBridgeCapExceededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCBridgeCapExceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCBridgeCapExceeded represents a BridgeCapExceeded event raised by the LBC contract.
type LBCBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeCapExceeded is a free log retrieval operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) FilterBridgeCapExceeded(opts *bind.FilterOpts) (*LBCBridgeCapExceededIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return &LBCBridgeCapExceededIterator{contract: _LBC.contract, event: "BridgeCapExceeded", logs: logs, sub: sub}, nil
}

// WatchBridgeCapExceeded is a free log subscription operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) WatchBridgeCapExceeded(opts *bind.WatchOpts, sink chan<- *LBCBridgeCapExceeded) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCBridgeCapExceeded)
				if err := _LBC.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
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

// ParseBridgeCapExceeded is a log parse operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) ParseBridgeCapExceeded(log types.Log) (*LBCBridgeCapExceeded, error) {
	event := new(LBCBridgeCapExceeded)
	if err := _LBC.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCBridgeErrorIterator is returned from FilterBridgeError and is used to iterate over the raw logs and unpacked data for BridgeError events raised by the LBC contract.
type LBCBridgeErrorIterator struct {
	Event *LBCBridgeError // Event containing the contract specifics and raw log

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
func (it *LBCBridgeErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCBridgeError)
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
		it.Event = new(LBCBridgeError)
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
func (it *LBCBridgeErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCBridgeErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCBridgeError represents a BridgeError event raised by the LBC contract.
type LBCBridgeError struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeError is a free log retrieval operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) FilterBridgeError(opts *bind.FilterOpts) (*LBCBridgeErrorIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return &LBCBridgeErrorIterator{contract: _LBC.contract, event: "BridgeError", logs: logs, sub: sub}, nil
}

// WatchBridgeError is a free log subscription operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) WatchBridgeError(opts *bind.WatchOpts, sink chan<- *LBCBridgeError) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCBridgeError)
				if err := _LBC.contract.UnpackLog(event, "BridgeError", log); err != nil {
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

// ParseBridgeError is a log parse operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LBC *LBCFilterer) ParseBridgeError(log types.Log) (*LBCBridgeError, error) {
	event := new(LBCBridgeError)
	if err := _LBC.contract.UnpackLog(event, "BridgeError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCCallForUserIterator is returned from FilterCallForUser and is used to iterate over the raw logs and unpacked data for CallForUser events raised by the LBC contract.
type LBCCallForUserIterator struct {
	Event *LBCCallForUser // Event containing the contract specifics and raw log

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
func (it *LBCCallForUserIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCCallForUser)
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
		it.Event = new(LBCCallForUser)
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
func (it *LBCCallForUserIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCCallForUserIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCCallForUser represents a CallForUser event raised by the LBC contract.
type LBCCallForUser struct {
	From      common.Address
	Dest      common.Address
	GasLimit  *big.Int
	Value     *big.Int
	Data      []byte
	Success   bool
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCallForUser is a free log retrieval operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address from, address dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LBC *LBCFilterer) FilterCallForUser(opts *bind.FilterOpts) (*LBCCallForUserIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "CallForUser")
	if err != nil {
		return nil, err
	}
	return &LBCCallForUserIterator{contract: _LBC.contract, event: "CallForUser", logs: logs, sub: sub}, nil
}

// WatchCallForUser is a free log subscription operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address from, address dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LBC *LBCFilterer) WatchCallForUser(opts *bind.WatchOpts, sink chan<- *LBCCallForUser) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "CallForUser")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCCallForUser)
				if err := _LBC.contract.UnpackLog(event, "CallForUser", log); err != nil {
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

// ParseCallForUser is a log parse operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address from, address dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LBC *LBCFilterer) ParseCallForUser(log types.Log) (*LBCCallForUser, error) {
	event := new(LBCCallForUser)
	if err := _LBC.contract.UnpackLog(event, "CallForUser", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCCollateralIncreaseIterator is returned from FilterCollateralIncrease and is used to iterate over the raw logs and unpacked data for CollateralIncrease events raised by the LBC contract.
type LBCCollateralIncreaseIterator struct {
	Event *LBCCollateralIncrease // Event containing the contract specifics and raw log

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
func (it *LBCCollateralIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCCollateralIncrease)
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
		it.Event = new(LBCCollateralIncrease)
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
func (it *LBCCollateralIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCCollateralIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCCollateralIncrease represents a CollateralIncrease event raised by the LBC contract.
type LBCCollateralIncrease struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralIncrease is a free log retrieval operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LBC *LBCFilterer) FilterCollateralIncrease(opts *bind.FilterOpts) (*LBCCollateralIncreaseIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return &LBCCollateralIncreaseIterator{contract: _LBC.contract, event: "CollateralIncrease", logs: logs, sub: sub}, nil
}

// WatchCollateralIncrease is a free log subscription operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LBC *LBCFilterer) WatchCollateralIncrease(opts *bind.WatchOpts, sink chan<- *LBCCollateralIncrease) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCCollateralIncrease)
				if err := _LBC.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
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

// ParseCollateralIncrease is a log parse operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LBC *LBCFilterer) ParseCollateralIncrease(log types.Log) (*LBCCollateralIncrease, error) {
	event := new(LBCCollateralIncrease)
	if err := _LBC.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the LBC contract.
type LBCDepositIterator struct {
	Event *LBCDeposit // Event containing the contract specifics and raw log

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
func (it *LBCDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCDeposit)
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
		it.Event = new(LBCDeposit)
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
func (it *LBCDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCDeposit represents a Deposit event raised by the LBC contract.
type LBCDeposit struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LBC *LBCFilterer) FilterDeposit(opts *bind.FilterOpts) (*LBCDepositIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &LBCDepositIterator{contract: _LBC.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LBC *LBCFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *LBCDeposit) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCDeposit)
				if err := _LBC.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LBC *LBCFilterer) ParseDeposit(log types.Log) (*LBCDeposit, error) {
	event := new(LBCDeposit)
	if err := _LBC.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCPenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the LBC contract.
type LBCPenalizedIterator struct {
	Event *LBCPenalized // Event containing the contract specifics and raw log

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
func (it *LBCPenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCPenalized)
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
		it.Event = new(LBCPenalized)
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
func (it *LBCPenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCPenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCPenalized represents a Penalized event raised by the LBC contract.
type LBCPenalized struct {
	LiquidityProvider common.Address
	Penalty           *big.Int
	QuoteHash         [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LBC *LBCFilterer) FilterPenalized(opts *bind.FilterOpts) (*LBCPenalizedIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return &LBCPenalizedIterator{contract: _LBC.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LBC *LBCFilterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *LBCPenalized) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCPenalized)
				if err := _LBC.contract.UnpackLog(event, "Penalized", log); err != nil {
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

// ParsePenalized is a log parse operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LBC *LBCFilterer) ParsePenalized(log types.Log) (*LBCPenalized, error) {
	event := new(LBCPenalized)
	if err := _LBC.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCRefundIterator is returned from FilterRefund and is used to iterate over the raw logs and unpacked data for Refund events raised by the LBC contract.
type LBCRefundIterator struct {
	Event *LBCRefund // Event containing the contract specifics and raw log

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
func (it *LBCRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCRefund)
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
		it.Event = new(LBCRefund)
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
func (it *LBCRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCRefund represents a Refund event raised by the LBC contract.
type LBCRefund struct {
	Dest      common.Address
	Amount    *big.Int
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRefund is a free log retrieval operation binding the contract event 0x30ecd700ed3e55b0cd5fc1a997fedf95b19ea877119957e5777dacdcb1c0aa28.
//
// Solidity: event Refund(address dest, int256 amount, bytes32 quoteHash)
func (_LBC *LBCFilterer) FilterRefund(opts *bind.FilterOpts) (*LBCRefundIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return &LBCRefundIterator{contract: _LBC.contract, event: "Refund", logs: logs, sub: sub}, nil
}

// WatchRefund is a free log subscription operation binding the contract event 0x30ecd700ed3e55b0cd5fc1a997fedf95b19ea877119957e5777dacdcb1c0aa28.
//
// Solidity: event Refund(address dest, int256 amount, bytes32 quoteHash)
func (_LBC *LBCFilterer) WatchRefund(opts *bind.WatchOpts, sink chan<- *LBCRefund) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCRefund)
				if err := _LBC.contract.UnpackLog(event, "Refund", log); err != nil {
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

// ParseRefund is a log parse operation binding the contract event 0x30ecd700ed3e55b0cd5fc1a997fedf95b19ea877119957e5777dacdcb1c0aa28.
//
// Solidity: event Refund(address dest, int256 amount, bytes32 quoteHash)
func (_LBC *LBCFilterer) ParseRefund(log types.Log) (*LBCRefund, error) {
	event := new(LBCRefund)
	if err := _LBC.contract.UnpackLog(event, "Refund", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCRegisterIterator is returned from FilterRegister and is used to iterate over the raw logs and unpacked data for Register events raised by the LBC contract.
type LBCRegisterIterator struct {
	Event *LBCRegister // Event containing the contract specifics and raw log

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
func (it *LBCRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCRegister)
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
		it.Event = new(LBCRegister)
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
func (it *LBCRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCRegister represents a Register event raised by the LBC contract.
type LBCRegister struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
func (_LBC *LBCFilterer) FilterRegister(opts *bind.FilterOpts) (*LBCRegisterIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return &LBCRegisterIterator{contract: _LBC.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
func (_LBC *LBCFilterer) WatchRegister(opts *bind.WatchOpts, sink chan<- *LBCRegister) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCRegister)
				if err := _LBC.contract.UnpackLog(event, "Register", log); err != nil {
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

// ParseRegister is a log parse operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
func (_LBC *LBCFilterer) ParseRegister(log types.Log) (*LBCRegister, error) {
	event := new(LBCRegister)
	if err := _LBC.contract.UnpackLog(event, "Register", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCResignedIterator is returned from FilterResigned and is used to iterate over the raw logs and unpacked data for Resigned events raised by the LBC contract.
type LBCResignedIterator struct {
	Event *LBCResigned // Event containing the contract specifics and raw log

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
func (it *LBCResignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCResigned)
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
		it.Event = new(LBCResigned)
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
func (it *LBCResignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCResignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCResigned represents a Resigned event raised by the LBC contract.
type LBCResigned struct {
	From common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LBC *LBCFilterer) FilterResigned(opts *bind.FilterOpts) (*LBCResignedIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return &LBCResignedIterator{contract: _LBC.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LBC *LBCFilterer) WatchResigned(opts *bind.WatchOpts, sink chan<- *LBCResigned) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCResigned)
				if err := _LBC.contract.UnpackLog(event, "Resigned", log); err != nil {
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

// ParseResigned is a log parse operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LBC *LBCFilterer) ParseResigned(log types.Log) (*LBCResigned, error) {
	event := new(LBCResigned)
	if err := _LBC.contract.UnpackLog(event, "Resigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCWithdrawCollateralIterator is returned from FilterWithdrawCollateral and is used to iterate over the raw logs and unpacked data for WithdrawCollateral events raised by the LBC contract.
type LBCWithdrawCollateralIterator struct {
	Event *LBCWithdrawCollateral // Event containing the contract specifics and raw log

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
func (it *LBCWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCWithdrawCollateral)
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
		it.Event = new(LBCWithdrawCollateral)
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
func (it *LBCWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCWithdrawCollateral represents a WithdrawCollateral event raised by the LBC contract.
type LBCWithdrawCollateral struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCollateral is a free log retrieval operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LBC *LBCFilterer) FilterWithdrawCollateral(opts *bind.FilterOpts) (*LBCWithdrawCollateralIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &LBCWithdrawCollateralIterator{contract: _LBC.contract, event: "WithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchWithdrawCollateral is a free log subscription operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LBC *LBCFilterer) WatchWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *LBCWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCWithdrawCollateral)
				if err := _LBC.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
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

// ParseWithdrawCollateral is a log parse operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LBC *LBCFilterer) ParseWithdrawCollateral(log types.Log) (*LBCWithdrawCollateral, error) {
	event := new(LBCWithdrawCollateral)
	if err := _LBC.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LBCWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the LBC contract.
type LBCWithdrawalIterator struct {
	Event *LBCWithdrawal // Event containing the contract specifics and raw log

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
func (it *LBCWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LBCWithdrawal)
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
		it.Event = new(LBCWithdrawal)
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
func (it *LBCWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LBCWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LBCWithdrawal represents a Withdrawal event raised by the LBC contract.
type LBCWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LBC *LBCFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*LBCWithdrawalIterator, error) {

	logs, sub, err := _LBC.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &LBCWithdrawalIterator{contract: _LBC.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LBC *LBCFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *LBCWithdrawal) (event.Subscription, error) {

	logs, sub, err := _LBC.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LBCWithdrawal)
				if err := _LBC.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LBC *LBCFilterer) ParseWithdrawal(log types.Log) (*LBCWithdrawal, error) {
	event := new(LBCWithdrawal)
	if err := _LBC.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
