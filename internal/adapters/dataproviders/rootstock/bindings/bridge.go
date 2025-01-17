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
	_ = abi.ConvertType
)

// RskBridgeMetaData contains all meta data concerning the RskBridge contract.
var RskBridgeMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestChainHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForBtcReleaseClient\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForDebugging\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainInitialBlockHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"depth\",\"type\":\"int256\"}],\"name\":\"getBtcBlockchainBlockHashAtDepth\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"getBtcTxHashProcessedHeight\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"isBtcTxHashAlreadyProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"atx\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"height\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"pmt\",\"type\":\"bytes\"}],\"name\":\"registerBtcTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"txhash\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"blocks\",\"type\":\"bytes[]\"}],\"name\":\"receiveHeaders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ablock\",\"type\":\"bytes\"}],\"name\":\"receiveHeader\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getRetiringFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getRetiringFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rskKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"mstKey\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKeyMultikey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"hash\",\"type\":\"bytes\"}],\"name\":\"commitFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollbackFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getPendingFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getPendingFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockWhitelistSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"getLockWhitelistEntryByAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addOneOffLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"addUnlimitedLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"removeLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"disableDelay\",\"type\":\"int256\"}],\"name\":\"setLockWhitelistDisableBlockDelay\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePerKb\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"feePerKb\",\"type\":\"int256\"}],\"name\":\"voteFeePerKbChange\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateCollections\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLockTxValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"getBtcTransactionConfirmations\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockingCap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"newLockingCap\",\"type\":\"int256\"}],\"name\":\"increaseLockingCap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"witnessMerkleRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"witnessReservedValue\",\"type\":\"bytes32\"}],\"name\":\"registerBtcCoinbaseTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"hasBtcBlockCoinbaseTransactionInformation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"shouldTransferToContract\",\"type\":\"bool\"}],\"name\":\"registerFastBridgeBtcTransaction\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveFederationCreationBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActivePowpegRedeemScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestBlockHeader\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"btcBlockHeight\",\"type\":\"uint256\"}],\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainParentBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
	parsed, err := RskBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
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
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_RskBridge *RskBridgeSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetActiveFederationCreationBlockHeight(&_RskBridge.CallOpts)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_RskBridge *RskBridgeCallerSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetActiveFederationCreationBlockHeight(&_RskBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
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
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (_RskBridge *RskBridgeSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _RskBridge.Contract.GetActivePowpegRedeemScript(&_RskBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _RskBridge.Contract.GetActivePowpegRedeemScript(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainBestBlockHeader(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainBestBlockHeader")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBestBlockHeader(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBestBlockHeader(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainBestChainHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainBestChainHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetBtcBlockchainBestChainHeight(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetBtcBlockchainBestChainHeight(&_RskBridge.CallOpts)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainBlockHashAtDepth(opts *bind.CallOpts, depth *big.Int) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHashAtDepth", depth)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_RskBridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_RskBridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_RskBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_RskBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
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
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_RskBridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_RskBridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainInitialBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainInitialBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetBtcBlockchainInitialBlockHeight(&_RskBridge.CallOpts)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _RskBridge.Contract.GetBtcBlockchainInitialBlockHeight(&_RskBridge.CallOpts)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetBtcBlockchainParentBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcBlockchainParentBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_RskBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _RskBridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_RskBridge.CallOpts, btcBlockHash)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_RskBridge *RskBridgeCaller) GetBtcTransactionConfirmations(opts *bind.CallOpts, txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcTransactionConfirmations", txHash, blockHash, merkleBranchPath, merkleBranchHashes)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_RskBridge *RskBridgeSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _RskBridge.Contract.GetBtcTransactionConfirmations(&_RskBridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _RskBridge.Contract.GetBtcTransactionConfirmations(&_RskBridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_RskBridge *RskBridgeCaller) GetBtcTxHashProcessedHeight(opts *bind.CallOpts, hash string) (int64, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getBtcTxHashProcessedHeight", hash)

	if err != nil {
		return *new(int64), err
	}

	out0 := *abi.ConvertType(out[0], new(int64)).(*int64)

	return out0, err

}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_RskBridge *RskBridgeSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _RskBridge.Contract.GetBtcTxHashProcessedHeight(&_RskBridge.CallOpts, hash)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_RskBridge *RskBridgeCallerSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _RskBridge.Contract.GetBtcTxHashProcessedHeight(&_RskBridge.CallOpts, hash)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
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
// Solidity: function getFederationAddress() view returns(string)
func (_RskBridge *RskBridgeSession) GetFederationAddress() (string, error) {
	return _RskBridge.Contract.GetFederationAddress(&_RskBridge.CallOpts)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_RskBridge *RskBridgeCallerSession) GetFederationAddress() (string, error) {
	return _RskBridge.Contract.GetFederationAddress(&_RskBridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationCreationBlockNumber(&_RskBridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationCreationBlockNumber(&_RskBridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationCreationTime() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationCreationTime(&_RskBridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationCreationTime() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationCreationTime(&_RskBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
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
// Solidity: function getFederationSize() view returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationSize(&_RskBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationSize(&_RskBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
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
// Solidity: function getFederationThreshold() view returns(int256)
func (_RskBridge *RskBridgeSession) GetFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationThreshold(&_RskBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetFederationThreshold(&_RskBridge.CallOpts)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
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
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetFeePerKb(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getFeePerKb")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_RskBridge *RskBridgeSession) GetFeePerKb() (*big.Int, error) {
	return _RskBridge.Contract.GetFeePerKb(&_RskBridge.CallOpts)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetFeePerKb() (*big.Int, error) {
	return _RskBridge.Contract.GetFeePerKb(&_RskBridge.CallOpts)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_RskBridge *RskBridgeCaller) GetLockWhitelistAddress(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getLockWhitelistAddress", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_RskBridge *RskBridgeSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _RskBridge.Contract.GetLockWhitelistAddress(&_RskBridge.CallOpts, index)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_RskBridge *RskBridgeCallerSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _RskBridge.Contract.GetLockWhitelistAddress(&_RskBridge.CallOpts, index)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_RskBridge *RskBridgeCaller) GetLockWhitelistEntryByAddress(opts *bind.CallOpts, aaddress string) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getLockWhitelistEntryByAddress", aaddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_RskBridge *RskBridgeSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _RskBridge.Contract.GetLockWhitelistEntryByAddress(&_RskBridge.CallOpts, aaddress)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _RskBridge.Contract.GetLockWhitelistEntryByAddress(&_RskBridge.CallOpts, aaddress)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetLockWhitelistSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getLockWhitelistSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_RskBridge *RskBridgeSession) GetLockWhitelistSize() (*big.Int, error) {
	return _RskBridge.Contract.GetLockWhitelistSize(&_RskBridge.CallOpts)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetLockWhitelistSize() (*big.Int, error) {
	return _RskBridge.Contract.GetLockWhitelistSize(&_RskBridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetLockingCap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getLockingCap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_RskBridge *RskBridgeSession) GetLockingCap() (*big.Int, error) {
	return _RskBridge.Contract.GetLockingCap(&_RskBridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetLockingCap() (*big.Int, error) {
	return _RskBridge.Contract.GetLockingCap(&_RskBridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
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
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_RskBridge *RskBridgeSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _RskBridge.Contract.GetMinimumLockTxValue(&_RskBridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _RskBridge.Contract.GetMinimumLockTxValue(&_RskBridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetPendingFederationHash(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getPendingFederationHash")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_RskBridge *RskBridgeSession) GetPendingFederationHash() ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederationHash(&_RskBridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetPendingFederationHash() ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederationHash(&_RskBridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetPendingFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getPendingFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_RskBridge *RskBridgeSession) GetPendingFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetPendingFederationSize(&_RskBridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetPendingFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetPendingFederationSize(&_RskBridge.CallOpts)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetPendingFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getPendingFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetPendingFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getPendingFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetPendingFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_RskBridge *RskBridgeCaller) GetRetiringFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_RskBridge *RskBridgeSession) GetRetiringFederationAddress() (string, error) {
	return _RskBridge.Contract.GetRetiringFederationAddress(&_RskBridge.CallOpts)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederationAddress() (string, error) {
	return _RskBridge.Contract.GetRetiringFederationAddress(&_RskBridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetRetiringFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationCreationBlockNumber(&_RskBridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationCreationBlockNumber(&_RskBridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetRetiringFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationCreationTime(&_RskBridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationCreationTime(&_RskBridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetRetiringFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_RskBridge *RskBridgeSession) GetRetiringFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationSize(&_RskBridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederationSize() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationSize(&_RskBridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_RskBridge *RskBridgeCaller) GetRetiringFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_RskBridge *RskBridgeSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationThreshold(&_RskBridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _RskBridge.Contract.GetRetiringFederationThreshold(&_RskBridge.CallOpts)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetRetiringFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetRetiringFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _RskBridge.Contract.GetRetiringFederatorPublicKey(&_RskBridge.CallOpts, index)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetRetiringFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getRetiringFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetRetiringFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _RskBridge.Contract.GetRetiringFederatorPublicKeyOfType(&_RskBridge.CallOpts, index, atype)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetStateForBtcReleaseClient(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getStateForBtcReleaseClient")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_RskBridge *RskBridgeSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _RskBridge.Contract.GetStateForBtcReleaseClient(&_RskBridge.CallOpts)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _RskBridge.Contract.GetStateForBtcReleaseClient(&_RskBridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_RskBridge *RskBridgeCaller) GetStateForDebugging(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "getStateForDebugging")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_RskBridge *RskBridgeSession) GetStateForDebugging() ([]byte, error) {
	return _RskBridge.Contract.GetStateForDebugging(&_RskBridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_RskBridge *RskBridgeCallerSession) GetStateForDebugging() ([]byte, error) {
	return _RskBridge.Contract.GetStateForDebugging(&_RskBridge.CallOpts)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_RskBridge *RskBridgeCaller) IsBtcTxHashAlreadyProcessed(opts *bind.CallOpts, hash string) (bool, error) {
	var out []interface{}
	err := _RskBridge.contract.Call(opts, &out, "isBtcTxHashAlreadyProcessed", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_RskBridge *RskBridgeSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _RskBridge.Contract.IsBtcTxHashAlreadyProcessed(&_RskBridge.CallOpts, hash)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_RskBridge *RskBridgeCallerSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _RskBridge.Contract.IsBtcTxHashAlreadyProcessed(&_RskBridge.CallOpts, hash)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_RskBridge *RskBridgeTransactor) AddFederatorPublicKey(opts *bind.TransactOpts, key []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addFederatorPublicKey", key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_RskBridge *RskBridgeSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddFederatorPublicKey(&_RskBridge.TransactOpts, key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddFederatorPublicKey(&_RskBridge.TransactOpts, key)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_RskBridge *RskBridgeTransactor) AddFederatorPublicKeyMultikey(opts *bind.TransactOpts, btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addFederatorPublicKeyMultikey", btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_RskBridge *RskBridgeSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddFederatorPublicKeyMultikey(&_RskBridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddFederatorPublicKeyMultikey(&_RskBridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeTransactor) AddLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.AddLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.AddLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeTransactor) AddOneOffLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addOneOffLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.AddOneOffLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.AddOneOffLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_RskBridge *RskBridgeTransactor) AddSignature(opts *bind.TransactOpts, pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addSignature", pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_RskBridge *RskBridgeSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddSignature(&_RskBridge.TransactOpts, pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_RskBridge *RskBridgeTransactorSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.AddSignature(&_RskBridge.TransactOpts, pubkey, signatures, txhash)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeTransactor) AddUnlimitedLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "addUnlimitedLockWhitelistAddress", aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _RskBridge.Contract.AddUnlimitedLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _RskBridge.Contract.AddUnlimitedLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_RskBridge *RskBridgeTransactor) CommitFederation(opts *bind.TransactOpts, hash []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "commitFederation", hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_RskBridge *RskBridgeSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.CommitFederation(&_RskBridge.TransactOpts, hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.CommitFederation(&_RskBridge.TransactOpts, hash)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_RskBridge *RskBridgeTransactor) CreateFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "createFederation")
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_RskBridge *RskBridgeSession) CreateFederation() (*types.Transaction, error) {
	return _RskBridge.Contract.CreateFederation(&_RskBridge.TransactOpts)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_RskBridge *RskBridgeTransactorSession) CreateFederation() (*types.Transaction, error) {
	return _RskBridge.Contract.CreateFederation(&_RskBridge.TransactOpts)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_RskBridge *RskBridgeTransactor) HasBtcBlockCoinbaseTransactionInformation(opts *bind.TransactOpts, blockHash [32]byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "hasBtcBlockCoinbaseTransactionInformation", blockHash)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_RskBridge *RskBridgeSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_RskBridge.TransactOpts, blockHash)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_RskBridge *RskBridgeTransactorSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_RskBridge.TransactOpts, blockHash)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_RskBridge *RskBridgeTransactor) IncreaseLockingCap(opts *bind.TransactOpts, newLockingCap *big.Int) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "increaseLockingCap", newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_RskBridge *RskBridgeSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.IncreaseLockingCap(&_RskBridge.TransactOpts, newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_RskBridge *RskBridgeTransactorSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.IncreaseLockingCap(&_RskBridge.TransactOpts, newLockingCap)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_RskBridge *RskBridgeTransactor) ReceiveHeader(opts *bind.TransactOpts, ablock []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "receiveHeader", ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_RskBridge *RskBridgeSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.ReceiveHeader(&_RskBridge.TransactOpts, ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.ReceiveHeader(&_RskBridge.TransactOpts, ablock)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_RskBridge *RskBridgeTransactor) ReceiveHeaders(opts *bind.TransactOpts, blocks [][]byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "receiveHeaders", blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_RskBridge *RskBridgeSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.ReceiveHeaders(&_RskBridge.TransactOpts, blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_RskBridge *RskBridgeTransactorSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.ReceiveHeaders(&_RskBridge.TransactOpts, blocks)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_RskBridge *RskBridgeTransactor) RegisterBtcCoinbaseTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "registerBtcCoinbaseTransaction", btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_RskBridge *RskBridgeSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterBtcCoinbaseTransaction(&_RskBridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_RskBridge *RskBridgeTransactorSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterBtcCoinbaseTransaction(&_RskBridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_RskBridge *RskBridgeTransactor) RegisterBtcTransaction(opts *bind.TransactOpts, atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "registerBtcTransaction", atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_RskBridge *RskBridgeSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterBtcTransaction(&_RskBridge.TransactOpts, atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_RskBridge *RskBridgeTransactorSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterBtcTransaction(&_RskBridge.TransactOpts, atx, height, pmt)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeTransactor) RegisterFastBridgeBtcTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterFastBridgeBtcTransaction(&_RskBridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _RskBridge.Contract.RegisterFastBridgeBtcTransaction(&_RskBridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeTransactor) RemoveLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "removeLockWhitelistAddress", aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _RskBridge.Contract.RemoveLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _RskBridge.Contract.RemoveLockWhitelistAddress(&_RskBridge.TransactOpts, aaddress)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_RskBridge *RskBridgeTransactor) RollbackFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "rollbackFederation")
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_RskBridge *RskBridgeSession) RollbackFederation() (*types.Transaction, error) {
	return _RskBridge.Contract.RollbackFederation(&_RskBridge.TransactOpts)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_RskBridge *RskBridgeTransactorSession) RollbackFederation() (*types.Transaction, error) {
	return _RskBridge.Contract.RollbackFederation(&_RskBridge.TransactOpts)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_RskBridge *RskBridgeTransactor) SetLockWhitelistDisableBlockDelay(opts *bind.TransactOpts, disableDelay *big.Int) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "setLockWhitelistDisableBlockDelay", disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_RskBridge *RskBridgeSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.SetLockWhitelistDisableBlockDelay(&_RskBridge.TransactOpts, disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.SetLockWhitelistDisableBlockDelay(&_RskBridge.TransactOpts, disableDelay)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_RskBridge *RskBridgeTransactor) UpdateCollections(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "updateCollections")
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_RskBridge *RskBridgeSession) UpdateCollections() (*types.Transaction, error) {
	return _RskBridge.Contract.UpdateCollections(&_RskBridge.TransactOpts)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_RskBridge *RskBridgeTransactorSession) UpdateCollections() (*types.Transaction, error) {
	return _RskBridge.Contract.UpdateCollections(&_RskBridge.TransactOpts)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_RskBridge *RskBridgeTransactor) VoteFeePerKbChange(opts *bind.TransactOpts, feePerKb *big.Int) (*types.Transaction, error) {
	return _RskBridge.contract.Transact(opts, "voteFeePerKbChange", feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_RskBridge *RskBridgeSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.VoteFeePerKbChange(&_RskBridge.TransactOpts, feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_RskBridge *RskBridgeTransactorSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _RskBridge.Contract.VoteFeePerKbChange(&_RskBridge.TransactOpts, feePerKb)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RskBridge *RskBridgeTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RskBridge.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RskBridge *RskBridgeSession) Receive() (*types.Transaction, error) {
	return _RskBridge.Contract.Receive(&_RskBridge.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_RskBridge *RskBridgeTransactorSession) Receive() (*types.Transaction, error) {
	return _RskBridge.Contract.Receive(&_RskBridge.TransactOpts)
}
