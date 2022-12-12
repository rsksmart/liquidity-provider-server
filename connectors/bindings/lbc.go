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

// LiquidityBridgeContractProvider is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractProvider struct {
	Id       *big.Int
	Provider common.Address
}

// LiquidityBridgeContractQuote is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractQuote struct {
	FedBtcAddress               [20]byte
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	BtcRefundAddress            []byte
	RskRefundAddress            common.Address
	LiquidityProviderBtcAddress []byte
	CallFee                     *big.Int
	PenaltyFee                  *big.Int
	ContractAddress             common.Address
	Data                        []byte
	GasLimit                    uint32
	Nonce                       int64
	Value                       *big.Int
	AgreementTimestamp          uint32
	TimeForDeposit              uint32
	CallTime                    uint32
	DepositConfirmations        uint16
	CallOnRegister              bool
}

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rskKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"mstKey\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKeyMultikey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addOneOffLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"txhash\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"addUnlimitedLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"hash\",\"type\":\"bytes\"}],\"name\":\"commitFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveFederationCreationBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestBlockHeader\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestChainHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"depth\",\"type\":\"int256\"}],\"name\":\"getBtcBlockchainBlockHashAtDepth\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"btcBlockHeight\",\"type\":\"uint256\"}],\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainInitialBlockHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainParentBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"getBtcTransactionConfirmations\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"getBtcTxHashProcessedHeight\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePerKb\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"getLockWhitelistEntryByAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockWhitelistSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockingCap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLockTxValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getPendingFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getPendingFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getRetiringFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getRetiringFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForBtcReleaseClient\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForDebugging\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"hasBtcBlockCoinbaseTransactionInformation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"newLockingCap\",\"type\":\"int256\"}],\"name\":\"increaseLockingCap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"isBtcTxHashAlreadyProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ablock\",\"type\":\"bytes\"}],\"name\":\"receiveHeader\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"blocks\",\"type\":\"bytes[]\"}],\"name\":\"receiveHeaders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"witnessMerkleRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"witnessReservedValue\",\"type\":\"bytes32\"}],\"name\":\"registerBtcCoinbaseTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"atx\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"height\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"pmt\",\"type\":\"bytes\"}],\"name\":\"registerBtcTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"shouldTransferToContract\",\"type\":\"bool\"}],\"name\":\"registerFastBridgeBtcTransaction\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"removeLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollbackFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"disableDelay\",\"type\":\"int256\"}],\"name\":\"setLockWhitelistDisableBlockDelay\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateCollections\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"feePerKb\",\"type\":\"int256\"}],\"name\":\"voteFeePerKbChange\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"ecefd339": "addFederatorPublicKey(bytes)",
		"444ff9da": "addFederatorPublicKeyMultikey(bytes,bytes,bytes)",
		"502bbbce": "addLockWhitelistAddress(string,int256)",
		"848206d9": "addOneOffLockWhitelistAddress(string,int256)",
		"f10b9c59": "addSignature(bytes,bytes[],bytes)",
		"b906c938": "addUnlimitedLockWhitelistAddress(string)",
		"1533330f": "commitFederation(bytes)",
		"1183d5d1": "createFederation()",
		"177d6e18": "getActiveFederationCreationBlockHeight()",
		"f0b2424b": "getBtcBlockchainBestBlockHeader()",
		"14c89c01": "getBtcBlockchainBestChainHeight()",
		"efd44418": "getBtcBlockchainBlockHashAtDepth(int256)",
		"739e364a": "getBtcBlockchainBlockHeaderByHash(bytes32)",
		"bd0c1fff": "getBtcBlockchainBlockHeaderByHeight(uint256)",
		"4897193f": "getBtcBlockchainInitialBlockHeight()",
		"e0236724": "getBtcBlockchainParentBlockHeaderByHash(bytes32)",
		"5b644587": "getBtcTransactionConfirmations(bytes32,bytes32,uint256,bytes32[])",
		"97fcca7d": "getBtcTxHashProcessedHeight(string)",
		"6923fa85": "getFederationAddress()",
		"1b2045ee": "getFederationCreationBlockNumber()",
		"5e2db9d4": "getFederationCreationTime()",
		"802ad4b6": "getFederationSize()",
		"0fd47456": "getFederationThreshold()",
		"6b89a1af": "getFederatorPublicKey(int256)",
		"549cfd1c": "getFederatorPublicKeyOfType(int256,string)",
		"724ec886": "getFeePerKb()",
		"93988b76": "getLockWhitelistAddress(int256)",
		"251c5f7b": "getLockWhitelistEntryByAddress(string)",
		"e9e658dc": "getLockWhitelistSize()",
		"3f9db977": "getLockingCap()",
		"2f8d158f": "getMinimumLockTxValue()",
		"6ce0ed5a": "getPendingFederationHash()",
		"3ac72b38": "getPendingFederationSize()",
		"492f7c44": "getPendingFederatorPublicKey(int256)",
		"c61295d9": "getPendingFederatorPublicKeyOfType(int256,string)",
		"47227286": "getRetiringFederationAddress()",
		"d905153f": "getRetiringFederationCreationBlockNumber()",
		"3f0ce9b1": "getRetiringFederationCreationTime()",
		"d970b0fd": "getRetiringFederationSize()",
		"07bbdfc4": "getRetiringFederationThreshold()",
		"4675d6de": "getRetiringFederatorPublicKey(int256)",
		"68bc2b2b": "getRetiringFederatorPublicKeyOfType(int256,string)",
		"c4fbca20": "getStateForBtcReleaseClient()",
		"0d0cee93": "getStateForDebugging()",
		"253b944b": "hasBtcBlockCoinbaseTransactionInformation(bytes32)",
		"2910aeb2": "increaseLockingCap(int256)",
		"248a982d": "isBtcTxHashAlreadyProcessed(string)",
		"884bdd86": "receiveHeader(bytes)",
		"e5400e7b": "receiveHeaders(bytes[])",
		"ccf417ae": "registerBtcCoinbaseTransaction(bytes,bytes32,bytes,bytes32,bytes32)",
		"43dc0656": "registerBtcTransaction(bytes,int256,bytes)",
		"6adc0133": "registerFastBridgeBtcTransaction(bytes,uint256,bytes,bytes32,bytes,address,bytes,bool)",
		"fcdeb46f": "removeLockWhitelistAddress(string)",
		"8dec3d32": "rollbackFederation()",
		"c1cc54f5": "setLockWhitelistDisableBlockDelay(int256)",
		"0c5a9990": "updateCollections()",
		"0461313e": "voteFeePerKbChange(int256)",
	},
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// Deprecated: Use BridgeMetaData.Sigs instead.
// BridgeFuncSigs maps the 4-byte function signature to its string representation.
var BridgeFuncSigs = BridgeMetaData.Sigs

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_Bridge *BridgeCaller) GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getActiveFederationCreationBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_Bridge *BridgeSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _Bridge.Contract.GetActiveFederationCreationBlockHeight(&_Bridge.CallOpts)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_Bridge *BridgeCallerSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _Bridge.Contract.GetActiveFederationCreationBlockHeight(&_Bridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_Bridge *BridgeCaller) GetBtcBlockchainBestBlockHeader(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainBestBlockHeader")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_Bridge *BridgeSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBestBlockHeader(&_Bridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBestBlockHeader(&_Bridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_Bridge *BridgeCaller) GetBtcBlockchainBestChainHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainBestChainHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_Bridge *BridgeSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _Bridge.Contract.GetBtcBlockchainBestChainHeight(&_Bridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _Bridge.Contract.GetBtcBlockchainBestChainHeight(&_Bridge.CallOpts)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_Bridge *BridgeCaller) GetBtcBlockchainBlockHashAtDepth(opts *bind.CallOpts, depth *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainBlockHashAtDepth", depth)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_Bridge *BridgeSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_Bridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_Bridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeCaller) GetBtcBlockchainBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_Bridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_Bridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_Bridge *BridgeCaller) GetBtcBlockchainBlockHeaderByHeight(opts *bind.CallOpts, btcBlockHeight *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHeight", btcBlockHeight)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_Bridge *BridgeSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_Bridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_Bridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_Bridge *BridgeCaller) GetBtcBlockchainInitialBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainInitialBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_Bridge *BridgeSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _Bridge.Contract.GetBtcBlockchainInitialBlockHeight(&_Bridge.CallOpts)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _Bridge.Contract.GetBtcBlockchainInitialBlockHeight(&_Bridge.CallOpts)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeCaller) GetBtcBlockchainParentBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcBlockchainParentBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_Bridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _Bridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_Bridge.CallOpts, btcBlockHash)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_Bridge *BridgeCaller) GetBtcTransactionConfirmations(opts *bind.CallOpts, txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcTransactionConfirmations", txHash, blockHash, merkleBranchPath, merkleBranchHashes)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_Bridge *BridgeSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _Bridge.Contract.GetBtcTransactionConfirmations(&_Bridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_Bridge *BridgeCallerSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _Bridge.Contract.GetBtcTransactionConfirmations(&_Bridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_Bridge *BridgeCaller) GetBtcTxHashProcessedHeight(opts *bind.CallOpts, hash string) (int64, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getBtcTxHashProcessedHeight", hash)

	if err != nil {
		return *new(int64), err
	}

	out0 := *abi.ConvertType(out[0], new(int64)).(*int64)

	return out0, err

}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_Bridge *BridgeSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _Bridge.Contract.GetBtcTxHashProcessedHeight(&_Bridge.CallOpts, hash)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_Bridge *BridgeCallerSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _Bridge.Contract.GetBtcTxHashProcessedHeight(&_Bridge.CallOpts, hash)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_Bridge *BridgeCaller) GetFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_Bridge *BridgeSession) GetFederationAddress() (string, error) {
	return _Bridge.Contract.GetFederationAddress(&_Bridge.CallOpts)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_Bridge *BridgeCallerSession) GetFederationAddress() (string, error) {
	return _Bridge.Contract.GetFederationAddress(&_Bridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeCaller) GetFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _Bridge.Contract.GetFederationCreationBlockNumber(&_Bridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeCallerSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _Bridge.Contract.GetFederationCreationBlockNumber(&_Bridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_Bridge *BridgeCaller) GetFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_Bridge *BridgeSession) GetFederationCreationTime() (*big.Int, error) {
	return _Bridge.Contract.GetFederationCreationTime(&_Bridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_Bridge *BridgeCallerSession) GetFederationCreationTime() (*big.Int, error) {
	return _Bridge.Contract.GetFederationCreationTime(&_Bridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_Bridge *BridgeCaller) GetFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_Bridge *BridgeSession) GetFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetFederationSize(&_Bridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_Bridge *BridgeCallerSession) GetFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetFederationSize(&_Bridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_Bridge *BridgeCaller) GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_Bridge *BridgeSession) GetFederationThreshold() (*big.Int, error) {
	return _Bridge.Contract.GetFederationThreshold(&_Bridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_Bridge *BridgeCallerSession) GetFederationThreshold() (*big.Int, error) {
	return _Bridge.Contract.GetFederationThreshold(&_Bridge.CallOpts)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCaller) GetFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_Bridge *BridgeCaller) GetFeePerKb(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getFeePerKb")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_Bridge *BridgeSession) GetFeePerKb() (*big.Int, error) {
	return _Bridge.Contract.GetFeePerKb(&_Bridge.CallOpts)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_Bridge *BridgeCallerSession) GetFeePerKb() (*big.Int, error) {
	return _Bridge.Contract.GetFeePerKb(&_Bridge.CallOpts)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_Bridge *BridgeCaller) GetLockWhitelistAddress(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getLockWhitelistAddress", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_Bridge *BridgeSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _Bridge.Contract.GetLockWhitelistAddress(&_Bridge.CallOpts, index)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_Bridge *BridgeCallerSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _Bridge.Contract.GetLockWhitelistAddress(&_Bridge.CallOpts, index)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_Bridge *BridgeCaller) GetLockWhitelistEntryByAddress(opts *bind.CallOpts, aaddress string) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getLockWhitelistEntryByAddress", aaddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_Bridge *BridgeSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _Bridge.Contract.GetLockWhitelistEntryByAddress(&_Bridge.CallOpts, aaddress)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_Bridge *BridgeCallerSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _Bridge.Contract.GetLockWhitelistEntryByAddress(&_Bridge.CallOpts, aaddress)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_Bridge *BridgeCaller) GetLockWhitelistSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getLockWhitelistSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_Bridge *BridgeSession) GetLockWhitelistSize() (*big.Int, error) {
	return _Bridge.Contract.GetLockWhitelistSize(&_Bridge.CallOpts)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_Bridge *BridgeCallerSession) GetLockWhitelistSize() (*big.Int, error) {
	return _Bridge.Contract.GetLockWhitelistSize(&_Bridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_Bridge *BridgeCaller) GetLockingCap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getLockingCap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_Bridge *BridgeSession) GetLockingCap() (*big.Int, error) {
	return _Bridge.Contract.GetLockingCap(&_Bridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_Bridge *BridgeCallerSession) GetLockingCap() (*big.Int, error) {
	return _Bridge.Contract.GetLockingCap(&_Bridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_Bridge *BridgeCaller) GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getMinimumLockTxValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_Bridge *BridgeSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _Bridge.Contract.GetMinimumLockTxValue(&_Bridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_Bridge *BridgeCallerSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _Bridge.Contract.GetMinimumLockTxValue(&_Bridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_Bridge *BridgeCaller) GetPendingFederationHash(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getPendingFederationHash")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_Bridge *BridgeSession) GetPendingFederationHash() ([]byte, error) {
	return _Bridge.Contract.GetPendingFederationHash(&_Bridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_Bridge *BridgeCallerSession) GetPendingFederationHash() ([]byte, error) {
	return _Bridge.Contract.GetPendingFederationHash(&_Bridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_Bridge *BridgeCaller) GetPendingFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getPendingFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_Bridge *BridgeSession) GetPendingFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetPendingFederationSize(&_Bridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_Bridge *BridgeCallerSession) GetPendingFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetPendingFederationSize(&_Bridge.CallOpts)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCaller) GetPendingFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getPendingFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetPendingFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetPendingFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeCaller) GetPendingFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getPendingFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _Bridge.Contract.GetPendingFederatorPublicKeyOfType(&_Bridge.CallOpts, index, atype)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _Bridge.Contract.GetPendingFederatorPublicKeyOfType(&_Bridge.CallOpts, index, atype)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_Bridge *BridgeCaller) GetRetiringFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_Bridge *BridgeSession) GetRetiringFederationAddress() (string, error) {
	return _Bridge.Contract.GetRetiringFederationAddress(&_Bridge.CallOpts)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_Bridge *BridgeCallerSession) GetRetiringFederationAddress() (string, error) {
	return _Bridge.Contract.GetRetiringFederationAddress(&_Bridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeCaller) GetRetiringFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationCreationBlockNumber(&_Bridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_Bridge *BridgeCallerSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationCreationBlockNumber(&_Bridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_Bridge *BridgeCaller) GetRetiringFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_Bridge *BridgeSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationCreationTime(&_Bridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_Bridge *BridgeCallerSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationCreationTime(&_Bridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_Bridge *BridgeCaller) GetRetiringFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_Bridge *BridgeSession) GetRetiringFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationSize(&_Bridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_Bridge *BridgeCallerSession) GetRetiringFederationSize() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationSize(&_Bridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_Bridge *BridgeCaller) GetRetiringFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_Bridge *BridgeSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationThreshold(&_Bridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_Bridge *BridgeCallerSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _Bridge.Contract.GetRetiringFederationThreshold(&_Bridge.CallOpts)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCaller) GetRetiringFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetRetiringFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _Bridge.Contract.GetRetiringFederatorPublicKey(&_Bridge.CallOpts, index)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeCaller) GetRetiringFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRetiringFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _Bridge.Contract.GetRetiringFederatorPublicKeyOfType(&_Bridge.CallOpts, index, atype)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_Bridge *BridgeCallerSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _Bridge.Contract.GetRetiringFederatorPublicKeyOfType(&_Bridge.CallOpts, index, atype)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_Bridge *BridgeCaller) GetStateForBtcReleaseClient(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getStateForBtcReleaseClient")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_Bridge *BridgeSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _Bridge.Contract.GetStateForBtcReleaseClient(&_Bridge.CallOpts)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_Bridge *BridgeCallerSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _Bridge.Contract.GetStateForBtcReleaseClient(&_Bridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_Bridge *BridgeCaller) GetStateForDebugging(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getStateForDebugging")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_Bridge *BridgeSession) GetStateForDebugging() ([]byte, error) {
	return _Bridge.Contract.GetStateForDebugging(&_Bridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_Bridge *BridgeCallerSession) GetStateForDebugging() ([]byte, error) {
	return _Bridge.Contract.GetStateForDebugging(&_Bridge.CallOpts)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_Bridge *BridgeCaller) IsBtcTxHashAlreadyProcessed(opts *bind.CallOpts, hash string) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "isBtcTxHashAlreadyProcessed", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_Bridge *BridgeSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _Bridge.Contract.IsBtcTxHashAlreadyProcessed(&_Bridge.CallOpts, hash)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_Bridge *BridgeCallerSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _Bridge.Contract.IsBtcTxHashAlreadyProcessed(&_Bridge.CallOpts, hash)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_Bridge *BridgeTransactor) AddFederatorPublicKey(opts *bind.TransactOpts, key []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addFederatorPublicKey", key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_Bridge *BridgeSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddFederatorPublicKey(&_Bridge.TransactOpts, key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_Bridge *BridgeTransactorSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddFederatorPublicKey(&_Bridge.TransactOpts, key)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_Bridge *BridgeTransactor) AddFederatorPublicKeyMultikey(opts *bind.TransactOpts, btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addFederatorPublicKeyMultikey", btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_Bridge *BridgeSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddFederatorPublicKeyMultikey(&_Bridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_Bridge *BridgeTransactorSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddFederatorPublicKeyMultikey(&_Bridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeTransactor) AddLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AddLockWhitelistAddress(&_Bridge.TransactOpts, aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeTransactorSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AddLockWhitelistAddress(&_Bridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeTransactor) AddOneOffLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addOneOffLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AddOneOffLockWhitelistAddress(&_Bridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_Bridge *BridgeTransactorSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.AddOneOffLockWhitelistAddress(&_Bridge.TransactOpts, aaddress, maxTransferValue)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_Bridge *BridgeTransactor) AddSignature(opts *bind.TransactOpts, pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addSignature", pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_Bridge *BridgeSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddSignature(&_Bridge.TransactOpts, pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_Bridge *BridgeTransactorSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _Bridge.Contract.AddSignature(&_Bridge.TransactOpts, pubkey, signatures, txhash)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeTransactor) AddUnlimitedLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "addUnlimitedLockWhitelistAddress", aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _Bridge.Contract.AddUnlimitedLockWhitelistAddress(&_Bridge.TransactOpts, aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeTransactorSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _Bridge.Contract.AddUnlimitedLockWhitelistAddress(&_Bridge.TransactOpts, aaddress)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_Bridge *BridgeTransactor) CommitFederation(opts *bind.TransactOpts, hash []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "commitFederation", hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_Bridge *BridgeSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _Bridge.Contract.CommitFederation(&_Bridge.TransactOpts, hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_Bridge *BridgeTransactorSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _Bridge.Contract.CommitFederation(&_Bridge.TransactOpts, hash)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_Bridge *BridgeTransactor) CreateFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "createFederation")
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_Bridge *BridgeSession) CreateFederation() (*types.Transaction, error) {
	return _Bridge.Contract.CreateFederation(&_Bridge.TransactOpts)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_Bridge *BridgeTransactorSession) CreateFederation() (*types.Transaction, error) {
	return _Bridge.Contract.CreateFederation(&_Bridge.TransactOpts)
}

// GetFederatorPublicKeyOfType is a paid mutator transaction binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) returns(bytes)
func (_Bridge *BridgeTransactor) GetFederatorPublicKeyOfType(opts *bind.TransactOpts, index *big.Int, atype string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "getFederatorPublicKeyOfType", index, atype)
}

// GetFederatorPublicKeyOfType is a paid mutator transaction binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) returns(bytes)
func (_Bridge *BridgeSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) (*types.Transaction, error) {
	return _Bridge.Contract.GetFederatorPublicKeyOfType(&_Bridge.TransactOpts, index, atype)
}

// GetFederatorPublicKeyOfType is a paid mutator transaction binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) returns(bytes)
func (_Bridge *BridgeTransactorSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) (*types.Transaction, error) {
	return _Bridge.Contract.GetFederatorPublicKeyOfType(&_Bridge.TransactOpts, index, atype)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_Bridge *BridgeTransactor) HasBtcBlockCoinbaseTransactionInformation(opts *bind.TransactOpts, blockHash [32]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "hasBtcBlockCoinbaseTransactionInformation", blockHash)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_Bridge *BridgeSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_Bridge.TransactOpts, blockHash)
}

// HasBtcBlockCoinbaseTransactionInformation is a paid mutator transaction binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) returns(bool)
func (_Bridge *BridgeTransactorSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_Bridge.TransactOpts, blockHash)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_Bridge *BridgeTransactor) IncreaseLockingCap(opts *bind.TransactOpts, newLockingCap *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "increaseLockingCap", newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_Bridge *BridgeSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.IncreaseLockingCap(&_Bridge.TransactOpts, newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_Bridge *BridgeTransactorSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.IncreaseLockingCap(&_Bridge.TransactOpts, newLockingCap)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_Bridge *BridgeTransactor) ReceiveHeader(opts *bind.TransactOpts, ablock []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "receiveHeader", ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_Bridge *BridgeSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _Bridge.Contract.ReceiveHeader(&_Bridge.TransactOpts, ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_Bridge *BridgeTransactorSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _Bridge.Contract.ReceiveHeader(&_Bridge.TransactOpts, ablock)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_Bridge *BridgeTransactor) ReceiveHeaders(opts *bind.TransactOpts, blocks [][]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "receiveHeaders", blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_Bridge *BridgeSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.ReceiveHeaders(&_Bridge.TransactOpts, blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_Bridge *BridgeTransactorSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _Bridge.Contract.ReceiveHeaders(&_Bridge.TransactOpts, blocks)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_Bridge *BridgeTransactor) RegisterBtcCoinbaseTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerBtcCoinbaseTransaction", btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_Bridge *BridgeSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterBtcCoinbaseTransaction(&_Bridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_Bridge *BridgeTransactorSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterBtcCoinbaseTransaction(&_Bridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_Bridge *BridgeTransactor) RegisterBtcTransaction(opts *bind.TransactOpts, atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerBtcTransaction", atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_Bridge *BridgeSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterBtcTransaction(&_Bridge.TransactOpts, atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_Bridge *BridgeTransactorSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterBtcTransaction(&_Bridge.TransactOpts, atx, height, pmt)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_Bridge *BridgeTransactor) RegisterFastBridgeBtcTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_Bridge *BridgeSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterFastBridgeBtcTransaction(&_Bridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_Bridge *BridgeTransactorSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _Bridge.Contract.RegisterFastBridgeBtcTransaction(&_Bridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeTransactor) RemoveLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "removeLockWhitelistAddress", aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _Bridge.Contract.RemoveLockWhitelistAddress(&_Bridge.TransactOpts, aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_Bridge *BridgeTransactorSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _Bridge.Contract.RemoveLockWhitelistAddress(&_Bridge.TransactOpts, aaddress)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_Bridge *BridgeTransactor) RollbackFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "rollbackFederation")
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_Bridge *BridgeSession) RollbackFederation() (*types.Transaction, error) {
	return _Bridge.Contract.RollbackFederation(&_Bridge.TransactOpts)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_Bridge *BridgeTransactorSession) RollbackFederation() (*types.Transaction, error) {
	return _Bridge.Contract.RollbackFederation(&_Bridge.TransactOpts)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_Bridge *BridgeTransactor) SetLockWhitelistDisableBlockDelay(opts *bind.TransactOpts, disableDelay *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setLockWhitelistDisableBlockDelay", disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_Bridge *BridgeSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.SetLockWhitelistDisableBlockDelay(&_Bridge.TransactOpts, disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_Bridge *BridgeTransactorSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.SetLockWhitelistDisableBlockDelay(&_Bridge.TransactOpts, disableDelay)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_Bridge *BridgeTransactor) UpdateCollections(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "updateCollections")
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_Bridge *BridgeSession) UpdateCollections() (*types.Transaction, error) {
	return _Bridge.Contract.UpdateCollections(&_Bridge.TransactOpts)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_Bridge *BridgeTransactorSession) UpdateCollections() (*types.Transaction, error) {
	return _Bridge.Contract.UpdateCollections(&_Bridge.TransactOpts)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_Bridge *BridgeTransactor) VoteFeePerKbChange(opts *bind.TransactOpts, feePerKb *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "voteFeePerKbChange", feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_Bridge *BridgeSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.VoteFeePerKbChange(&_Bridge.TransactOpts, feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_Bridge *BridgeTransactorSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.VoteFeePerKbChange(&_Bridge.TransactOpts, feePerKb)
}

// LiquidityBridgeContractMetaData contains all meta data concerning the LiquidityBridgeContract contract.
var LiquidityBridgeContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minimumPegIn\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"rewardPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"resignDelayBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"dustThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"getBtcBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDustThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinPegIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProviders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"provider\",\"type\":\"address\"}],\"internalType\":\"structLiquidityBridgeContract.Provider[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"9e816999": "addCollateral()",
		"ac29d744": "callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool))",
		"d0e30db0": "deposit()",
		"f8b2cb4f": "getBalance(address)",
		"fb32c508": "getBridgeAddress()",
		"a0cd70fc": "getBtcBlockTimestamp(bytes)",
		"9b56d6c9": "getCollateral(address)",
		"33f07ad3": "getDustThreshold()",
		"e830b690": "getMinCollateral()",
		"fa88dcde": "getMinPegIn()",
		"edc922a9": "getProviders()",
		"bd5798c3": "getResignDelayBlocks()",
		"c7213163": "getRewardPercentage()",
		"1b032188": "hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool))",
		"457385f2": "isOperational(address)",
		"1aa3a008": "register()",
		"6e2e8c70": "registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool),bytes,bytes,bytes,uint256)",
		"69652fcf": "resign()",
		"2e1a7d4d": "withdraw(uint256)",
		"59c153be": "withdrawCollateral()",
	},
	Bin: "0x60c06040523480156200001157600080fd5b5060405162002cf838038062002cf883398101604081905262000034916200010b565b60648363ffffffff161115620000905760405162461bcd60e51b815260206004820152601960248201527f496e76616c6964207265776172642070657263656e7461676500000000000000604482015260640160405180910390fd5b600080546001600160a01b039097166001600160a01b03199097169690961790955560809390935260a0919091526006805463ffffffff938416640100000000026001600160401b031990911693909216929092171790556007556200017f565b805163ffffffff811681146200010657600080fd5b919050565b60008060008060008060c087890312156200012557600080fd5b86516001600160a01b03811681146200013d57600080fd5b60208801516040890151919750955093506200015c60608801620000f1565b92506200016c60808801620000f1565b915060a087015190509295509295509295565b60805160a051612b3e620001ba600039600081816103b70152611ad501526000818161032c0152818161049001526107900152612b3e6000f3fe6080604052600436106101235760003560e01c8063a0cd70fc116100a0578063e830b69011610064578063e830b6901461031d578063edc922a914610350578063f8b2cb4f14610372578063fa88dcde146103a8578063fb32c508146103db57600080fd5b8063a0cd70fc146102a4578063ac29d744146102c4578063bd5798c3146102d7578063c7213163146102fa578063d0e30db01461031557600080fd5b806359c153be116100e757806359c153be1461021c57806369652fcf146102315780636e2e8c70146102465780639b56d6c9146102665780639e8169991461029c57600080fd5b80631aa3a0081461017c5780631b032188146101845780632e1a7d4d146101b757806333f07ad3146101d7578063457385f2146101ec57600080fd5b36610177576000546001600160a01b031633146101755760405162461bcd60e51b815260206004820152600b60248201526a139bdd08185b1b1bddd95960aa1b60448201526064015b60405180910390fd5b005b600080fd5b610175610403565b34801561019057600080fd5b506101a461019f36600461253f565b61060d565b6040519081526020015b60405180910390f35b3480156101c357600080fd5b506101756101d2366004612574565b61061e565b3480156101e357600080fd5b506007546101a4565b3480156101f857600080fd5b5061020c61020736600461258d565b610763565b60405190151581526020016101ae565b34801561022857600080fd5b506101756107b4565b34801561023d57600080fd5b5061017561095a565b34801561025257600080fd5b506101a46102613660046125aa565b610a12565b34801561027257600080fd5b506101a461028136600461258d565b6001600160a01b031660009081526002602052604090205490565b610175611370565b3480156102b057600080fd5b506101a46102bf36600461265f565b6113ee565b61020c6102d236600461253f565b611450565b3480156102e357600080fd5b50600654640100000000900463ffffffff166101a4565b34801561030657600080fd5b5060065463ffffffff166101a4565b61017561180b565b34801561032957600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101a4565b34801561035c57600080fd5b5061036561183c565b6040516101ae9190612694565b34801561037e57600080fd5b506101a461038d36600461258d565b6001600160a01b031660009081526001602052604090205490565b3480156103b457600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101a4565b3480156103e757600080fd5b506000546040516001600160a01b0390911681526020016101ae565b32331461043c5760405162461bcd60e51b81526020600482015260076024820152664e6f7420454f4160c81b604482015260640161016c565b336000908152600260205260409020541561048e5760405162461bcd60e51b8152602060048201526012602482015271105b1c9958591e481c9959da5cdd195c995960721b604482015260640161016c565b7f00000000000000000000000000000000000000000000000000000000000000003410156104f65760405162461bcd60e51b8152602060048201526015602482015274139bdd08195b9bdd59da0818dbdb1b185d195c985b605a1b604482015260640161016c565b33600090815260056020526040902054156105535760405162461bcd60e51b815260206004820152601960248201527f576974686472617720636f6c6c61746572616c20666972737400000000000000604482015260640161016c565b336000908152600260205260408120349055600880549161057383612702565b90915550506040805180820182526008805480835233602080850182815260009384526003825292869020945185559151600190940180546001600160a01b0319166001600160a01b0390951694909417909355905483519081529081019190915234918101919091527fa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e906060015b60405180910390a1565b60006106188261193b565b92915050565b336000908152600160205260409020548111156106725760405162461bcd60e51b8152602060048201526012602482015271496e73756666696369656e742066756e647360701b604482015260640161016c565b336000908152600160205260408120805483929061069190849061271b565b9091555050604051600090339083908381818185875af1925050503d80600081146106d8576040519150601f19603f3d011682016040523d82523d6000602084013e6106dd565b606091505b50509050806107255760405162461bcd60e51b815260206004820152601460248201527314d95b991a5b99c8199d5b991cc819985a5b195960621b604482015260640161016c565b60408051338152602081018490527f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b6591015b60405180910390a15050565b600061076e82611b66565b80156106185750506001600160a01b03166000908152600260205260409020547f0000000000000000000000000000000000000000000000000000000000000000111590565b336000908152600560205260409020546108075760405162461bcd60e51b815260206004820152601460248201527313995959081d1bc81c995cda59db88199a5c9cdd60621b604482015260640161016c565b6006543360009081526005602052604090205464010000000090910463ffffffff1690610834904361271b565b10156108765760405162461bcd60e51b81526020600482015260116024820152704e6f7420656e6f75676820626c6f636b7360781b604482015260640161016c565b33600081815260026020908152604080832080549084905560059092528083208390555190929083908381818185875af1925050503d80600081146108d7576040519150601f19603f3d011682016040523d82523d6000602084013e6108dc565b606091505b50509050806109245760405162461bcd60e51b815260206004820152601460248201527314d95b991a5b99c8199d5b991cc819985a5b195960621b604482015260640161016c565b60408051338152602081018490527fa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d29101610757565b61096333611b66565b61097f5760405162461bcd60e51b815260040161016c9061272e565b33600090815260056020526040902054156109cf5760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e481c995cda59db995960821b604482015260640161016c565b3360008181526005602090815260409182902043905590519182527fa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d869101610603565b60095460009060ff1615610a595760405162461bcd60e51b815260206004820152600e60248201526d1499595b9d1c985b9d0818d85b1b60921b604482015260640161016c565b6009805460ff191660011790556000610a718761193b565b6000818152600a6020526040902054909150600160ff9091161115610ad85760405162461bcd60e51b815260206004820152601860248201527f51756f746520616c726561647920726567697374657265640000000000000000604482015260640161016c565b60408088015190516301a86b5560e41b815273__$a15b1ee43228f2a7b300d8dd4e3fc2a5c7$__91631a86b55091610b17919085908b906004016127a6565b602060405180830381865af4158015610b34573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b5891906127d6565b610b985760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b604482015260640161016c565b637fffffff8310610beb5760405162461bcd60e51b815260206004820152601e60248201527f486569676874206d757374206265206c6f776572207468616e20325e33310000604482015260640161016c565b6000610bfa8887878786611ba5565b905061012f8101610c645760405162461bcd60e51b815260206004820152602e60248201527f4572726f72202d3330333a204661696c656420746f2076616c6964617465204260448201526d2a21903a3930b739b0b1ba34b7b760911b606482015260840161016c565b61012e8101610cc75760405162461bcd60e51b815260206004820152602960248201527f4572726f72202d3330323a205472616e73616374696f6e20616c7265616479206044820152681c1c9bd8d95cdcd95960ba1b606482015260840161016c565b6101308101610d265760405162461bcd60e51b815260206004820152602560248201527f4572726f72202d3330343a205472616e73616374696f6e2076616c7565206973604482015264207a65726f60d81b606482015260840161016c565b6101318101610d9d5760405162461bcd60e51b815260206004820152603760248201527f4572726f72202d3330353a205472616e73616374696f6e205554584f2076616c60448201527f75652069732062656c6f7720746865206d696e696d756d000000000000000000606482015260840161016c565b6103848101610dee5760405162461bcd60e51b815260206004820152601860248201527f4572726f72202d3930303a20427269646765206572726f720000000000000000604482015260640161016c565b6000811380610dfe575060c71981145b80610e0a575060631981145b610e4d5760405162461bcd60e51b81526020600482015260146024820152732ab735b737bbb710213934b233b29032b93937b960611b604482015260640161016c565b600082815260046020526040902054610e70908990839063ffffffff1687611c4e565b15610f635760e08801516040808a01516001600160a01b031660009081526002602052908120549091610ea291611ed3565b905080600260008b604001516001600160a01b03166001600160a01b031681526020019081526020016000206000828254610edd919061271b565b909155505060408981015181516001600160a01b039091168152602081018390529081018490527f9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f49060600160405180910390a1600654600090606490610f4a9063ffffffff16846127f3565b610f54919061280a565b9050610f603382611eeb565b50505b60c719811480610f74575060631981145b15610fe8576000828152600a60209081526040808320805460ff191660021790556004825291829020805464ffffffffff1916905581518481529081018390527ffb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe910160405180910390a1915061135d9050565b80610ff38982611f57565b60008381526004602052604090205463ffffffff161561115957600083815260046020526040812054640100000000900460ff16156110515761104a828b60c001518c6101800151611045919061282c565b611ed3565b9050611062565b61105f828b60c00151611ed3565b90505b6110708a6040015182611eeb565b600061107c828461271b565b90506007548111156111525760808b01516040516000916001600160a01b0316906108fc90849084818181858888f193505050503d80600081146110dc576040519150601f19603f3d011682016040523d82523d6000602084013e6110e1565b606091505b50506080808e0151604080516001600160a01b0390921682526020820186905283151590820152606081018990529192507f3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6910160405180910390a180611150576111508c6040015183611eeb565b505b5050611329565b6102208901518190801561117257508961018001518110155b156112665760008a61010001516001600160a01b03168b610140015163ffffffff168c61018001518d61012001516040516111ad919061283f565b600060405180830381858888f193505050503d80600081146111eb576040519150601f19603f3d011682016040523d82523d6000602084013e6111f0565b606091505b505090507fbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d338c61010001518d61014001518e61018001518f6101200151868b604051611243979695949392919061285b565b60405180910390a18015611264576101808b0151611261908361271b565b91505b505b6007548111156113275760808a01516040516000916001600160a01b0316906108fc90849084818181858888f193505050503d80600081146112c4576040519150601f19603f3d011682016040523d82523d6000602084013e6112c9565b606091505b50506080808d0151604080516001600160a01b0390921682526020820186905283151590820152606081018890529192507f3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6910160405180910390a1505b505b506000918252600a60209081526040808420805460ff191660021790556004909152909120805464ffffffffff1916905590505b6009805460ff1916905595945050505050565b61137933611b66565b6113955760405162461bcd60e51b815260040161016c9061272e565b33600090815260026020526040812080543492906113b490849061282c565b9091555050604080513381523460208201527f456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af849101610603565b600081516050146114395760405162461bcd60e51b81526020600482015260156024820152740d2dcecc2d8d2c840d0cac2c8cae440d8cadccee8d605b1b604482015260640161016c565b611444826044611fdf565b63ffffffff1692915050565b600061145b33611b66565b6114775760405162461bcd60e51b815260040161016c9061272e565b60095460ff16156114bb5760405162461bcd60e51b815260206004820152600e60248201526d1499595b9d1c985b9d0818d85b1b60921b604482015260640161016c565b6009805460ff191660011790556040820151336001600160a01b03909116146115155760405162461bcd60e51b815260206004820152600c60248201526b155b985d5d1a1bdc9a5e995960a21b604482015260640161016c565b6101808201516040808401516001600160a01b031660009081526001602052205461154190349061282c565b10156115845760405162461bcd60e51b8152602060048201526012602482015271496e73756666696369656e742066756e647360701b604482015260640161016c565b600061158f8361193b565b6000818152600a602052604090205490915060ff16156115f15760405162461bcd60e51b815260206004820152601760248201527f51756f746520616c72656164792070726f636573736564000000000000000000604482015260640161016c565b6115ff836040015134611eeb565b610140830151611612906188b8906128b0565b63ffffffff165a101561165a5760405162461bcd60e51b815260206004820152601060248201526f496e73756666696369656e742067617360801b604482015260640161016c565b60008361010001516001600160a01b031684610140015163ffffffff16856101800151866101200151604051611690919061283f565b600060405180830381858888f193505050503d80600081146116ce576040519150601f19603f3d011682016040523d82523d6000602084013e6116d3565b606091505b509091505063ffffffff42111561172c5760405162461bcd60e51b815260206004820152601860248201527f426c6f636b2074696d657374616d70206f766572666c6f770000000000000000604482015260640161016c565b6000828152600460205260409020805463ffffffff19164263ffffffff1617905580156117895760008281526004602052604090819020805464ff00000000191664010000000017905584015161018085015161178991906120d7565b7fbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d3385610100015186610140015187610180015188610120015186886040516117d8979695949392919061285b565b60405180910390a16000918252600a6020526040909120805460ff1916600117905590506009805460ff19169055919050565b61181433611b66565b6118305760405162461bcd60e51b815260040161016c9061272e565b61183a3334611eeb565b565b6060600060085467ffffffffffffffff81111561185b5761185b61222e565b6040519080825280602002602001820160405280156118a057816020015b60408051808201909152600080825260208201528152602001906001900390816118795790505b5090506000805b6008548111611933576000818152600360205260409020541561192157600081815260036020908152604091829020825180840190935280548352600101546001600160a01b0316908201528351849084908110611907576119076128d4565b6020026020010181905250818061191d90612702565b9250505b8061192b81612702565b9150506118a7565b509092915050565b600081602001516001600160a01b0316306001600160a01b0316146119965760405162461bcd60e51b815260206004820152601160248201527057726f6e67204c4243206164647265737360781b604482015260640161016c565b6101008201516000546001600160a01b03918216911603611a0c5760405162461bcd60e51b815260206004820152602a60248201527f427269646765206973206e6f7420616e20616363657074656420636f6e7472616044820152696374206164647265737360b01b606482015260840161016c565b816060015151601514611a725760405162461bcd60e51b815260206004820152602860248201527f42544320726566756e642061646472657373206d757374206265203231206279604482015267746573206c6f6e6760c01b606482015260840161016c565b8160a0015151601514611ad35760405162461bcd60e51b8152602060048201526024808201527f425443204c502061646472657373206d757374206265203231206279746573206044820152636c6f6e6760e01b606482015260840161016c565b7f00000000000000000000000000000000000000000000000000000000000000008260c00151836101800151611b09919061282c565b1015611b4f5760405162461bcd60e51b8152602060048201526015602482015274151bdbc81b1bddc81859dc99595908185b5bdd5b9d605a1b604482015260640161016c565b611b5882612143565b805190602001209050919050565b6001600160a01b038116600090815260026020526040812054158015906106185750506001600160a01b03166000908152600560205260409020541590565b60008054606087015160a0880151848452600460208190526040808620549051636adc013360e01b81526001600160a01b0390951694636adc013394611c01948c948b948d948c94933093909263ffffffff16151591016128ea565b6020604051808303816000875af1158015611c20573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c44919061296d565b9695505050505050565b60008084138015611c7257508460c00151856101800151611c6f919061282c565b84105b15611c7f57506000611ecb565b6000805460405163bd0c1fff60e01b8152600481018590526001600160a01b039091169063bd0c1fff90602401600060405180830381865afa158015611cc9573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611cf19190810190612986565b90506000815111611d3b5760405162461bcd60e51b8152602060048201526014602482015273125b9d985b1a5908189b1bd8dac81a195a59da1d60621b604482015260640161016c565b6000611d46826113ee565b90506000611d73886101c0015163ffffffff16896101a0015163ffffffff1661217e90919063ffffffff16565b905080821115611d895760009350505050611ecb565b85600003611d9d5760019350505050611ecb565b600080546102008a01516001600160a01b039091169063bd0c1fff90600190611dca9061ffff168a61282c565b611dd4919061271b565b6040518263ffffffff1660e01b8152600401611df291815260200190565b600060405180830381865afa158015611e0f573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611e379190810190612986565b90506000815111611e815760405162461bcd60e51b8152602060048201526014602482015273125b9d985b1a5908189b1bd8dac81a195a59da1d60621b604482015260640161016c565b6000611e8c826113ee565b9050611eac8a6101e0015163ffffffff168261217e90919063ffffffff16565b881115611ec157600195505050505050611ecb565b6000955050505050505b949350505050565b6000818310611ee25781611ee4565b825b9392505050565b6001600160a01b03821660009081526001602052604081208054839290611f1390849061282c565b9091555050604080516001600160a01b0384168152602081018390527f42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f539101610757565b60008260c00151836101800151611f6e919061282c565b90506000611f7e6127108361280a565b9050611f8a818361271b565b831015611fd95760405162461bcd60e51b815260206004820152601a60248201527f546f6f206c6f77207472616e7366657272656420616d6f756e74000000000000604482015260640161016c565b50505050565b6000611fec82600461282c565b835110156120335760405162461bcd60e51b8152602060048201526014602482015273736c6963696e67206f7574206f662072616e676560601b604482015260640161016c565b60188361204184600361282c565b81518110612051576120516128d4565b016020015160f81c901b60108461206985600261282c565b81518110612079576120796128d4565b016020015160f81c901b60088561209186600161282c565b815181106120a1576120a16128d4565b0160200151865160f89190911c90911b908690869081106120c4576120c46128d4565b016020015160f81c171717905092915050565b6001600160a01b038216600090815260016020526040812080548392906120ff90849061271b565b9091555050604080516001600160a01b0384168152602081018390527f8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc00645879101610757565b606061214e82612196565b612157836121de565b6040516020016121689291906129f4565b6040516020818303038152906040529050919050565b600082820183811015611ee457600019915050610618565b6060816000015182602001518360400151846060015185608001518660a001518760c001518860e0015189610100015160405160200161216899989796959493929190612a19565b6060816101200151826101400151836101600151846101800151856101a00151866101c00151876101e0015188610200015189610220015160405160200161216899989796959493929190612a9d565b634e487b7160e01b600052604160045260246000fd5b604051610240810167ffffffffffffffff811182821017156122685761226861222e565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156122975761229761222e565b604052919050565b80356bffffffffffffffffffffffff19811681146122bc57600080fd5b919050565b6001600160a01b03811681146122d657600080fd5b50565b80356122bc816122c1565b600067ffffffffffffffff8211156122fe576122fe61222e565b50601f01601f191660200190565b600082601f83011261231d57600080fd5b813561233061232b826122e4565b61226e565b81815284602083860101111561234557600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff811681146122bc57600080fd5b8035600781900b81146122bc57600080fd5b803561ffff811681146122bc57600080fd5b80151581146122d657600080fd5b80356122bc8161239a565b600061024082840312156123c657600080fd5b6123ce612244565b90506123d98261229f565b81526123e7602083016122d9565b60208201526123f8604083016122d9565b6040820152606082013567ffffffffffffffff8082111561241857600080fd5b6124248583860161230c565b6060840152612435608085016122d9565b608084015260a084013591508082111561244e57600080fd5b61245a8583860161230c565b60a084015260c084013560c084015260e084013560e084015261010091506124838285016122d9565b828401526101209150818401358181111561249d57600080fd5b6124a98682870161230c565b838501525050506101406124be818401612362565b908201526101606124d0838201612376565b9082015261018082810135908201526101a06124ed818401612362565b908201526101c06124ff838201612362565b908201526101e0612511838201612362565b90820152610200612523838201612388565b908201526102206125358382016123a8565b9082015292915050565b60006020828403121561255157600080fd5b813567ffffffffffffffff81111561256857600080fd5b611ecb848285016123b3565b60006020828403121561258657600080fd5b5035919050565b60006020828403121561259f57600080fd5b8135611ee4816122c1565b600080600080600060a086880312156125c257600080fd5b853567ffffffffffffffff808211156125da57600080fd5b6125e689838a016123b3565b965060208801359150808211156125fc57600080fd5b61260889838a0161230c565b9550604088013591508082111561261e57600080fd5b61262a89838a0161230c565b9450606088013591508082111561264057600080fd5b5061264d8882890161230c565b95989497509295608001359392505050565b60006020828403121561267157600080fd5b813567ffffffffffffffff81111561268857600080fd5b611ecb8482850161230c565b602080825282518282018190526000919060409081850190868401855b828110156126df578151805185528601516001600160a01b03168685015292840192908501906001016126b1565b5091979650505050505050565b634e487b7160e01b600052601160045260246000fd5b600060018201612714576127146126ec565b5060010190565b81810381811115610618576106186126ec565b6020808252600e908201526d139bdd081c9959da5cdd195c995960921b604082015260600190565b60005b83811015612771578181015183820152602001612759565b50506000910152565b60008151808452612792816020860160208601612756565b601f01601f19169290920160200192915050565b60018060a01b03841681528260208201526060604082015260006127cd606083018461277a565b95945050505050565b6000602082840312156127e857600080fd5b8151611ee48161239a565b8082028115828204841417610618576106186126ec565b60008261282757634e487b7160e01b600052601260045260246000fd5b500490565b80820180821115610618576106186126ec565b60008251612851818460208701612756565b9190910192915050565b6001600160a01b0388811682528716602082015263ffffffff861660408201526060810185905260e06080820181905260009061289a9083018661277a565b93151560a08301525060c0015295945050505050565b63ffffffff8181168382160190808211156128cd576128cd6126ec565b5092915050565b634e487b7160e01b600052603260045260246000fd5b60006101008083526128fe8184018c61277a565b90508960208401528281036040840152612918818a61277a565b90508760608401528281036080840152612932818861277a565b6001600160a01b03871660a085015283810360c08501529050612955818661277a565b91505082151560e08301529998505050505050505050565b60006020828403121561297f57600080fd5b5051919050565b60006020828403121561299857600080fd5b815167ffffffffffffffff8111156129af57600080fd5b8201601f810184136129c057600080fd5b80516129ce61232b826122e4565b8181528560208385010111156129e357600080fd5b6127cd826020830160208601612756565b604081526000612a07604083018561277a565b82810360208401526127cd818561277a565b6bffffffffffffffffffffffff198a1681526001600160a01b038981166020830152888116604083015261012060608301819052600091612a5c8483018b61277a565b9150808916608085015283820360a0850152612a78828961277a565b60c085019790975260e084019590955250509116610100909101529695505050505050565b6000610120808352612ab18184018d61277a565b63ffffffff9b8c16602085015260079a909a0b604084015250506060810196909652938716608086015291861660a085015290941660c083015261ffff90931660e08201529115156101009092019190915291905056fea26469706673582212200519381d08b5812c2d1c56fe0accff1d1577ddda566d834b0900a013694a51bf64736f6c63430008110033",
}

// LiquidityBridgeContractABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidityBridgeContractMetaData.ABI instead.
var LiquidityBridgeContractABI = LiquidityBridgeContractMetaData.ABI

// Deprecated: Use LiquidityBridgeContractMetaData.Sigs instead.
// LiquidityBridgeContractFuncSigs maps the 4-byte function signature to its string representation.
var LiquidityBridgeContractFuncSigs = LiquidityBridgeContractMetaData.Sigs

// LiquidityBridgeContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidityBridgeContractMetaData.Bin instead.
var LiquidityBridgeContractBin = LiquidityBridgeContractMetaData.Bin

// DeployLiquidityBridgeContract deploys a new Ethereum contract, binding an instance of LiquidityBridgeContract to it.
func DeployLiquidityBridgeContract(auth *bind.TransactOpts, backend bind.ContractBackend, bridgeAddress common.Address, minimumCollateral *big.Int, minimumPegIn *big.Int, rewardPercentage uint32, resignDelayBlocks uint32, dustThreshold *big.Int) (common.Address, *types.Transaction, *LiquidityBridgeContract, error) {
	parsed, err := LiquidityBridgeContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	signatureValidatorAddr, _, _, _ := DeploySignatureValidator(auth, backend)
	LiquidityBridgeContractBin = strings.Replace(LiquidityBridgeContractBin, "__$a15b1ee43228f2a7b300d8dd4e3fc2a5c7$__", signatureValidatorAddr.String()[2:], -1)

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidityBridgeContractBin), backend, bridgeAddress, minimumCollateral, minimumPegIn, rewardPercentage, resignDelayBlocks, dustThreshold)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LiquidityBridgeContract{LiquidityBridgeContractCaller: LiquidityBridgeContractCaller{contract: contract}, LiquidityBridgeContractTransactor: LiquidityBridgeContractTransactor{contract: contract}, LiquidityBridgeContractFilterer: LiquidityBridgeContractFilterer{contract: contract}}, nil
}

// LiquidityBridgeContract is an auto generated Go binding around an Ethereum contract.
type LiquidityBridgeContract struct {
	LiquidityBridgeContractCaller     // Read-only binding to the contract
	LiquidityBridgeContractTransactor // Write-only binding to the contract
	LiquidityBridgeContractFilterer   // Log filterer for contract events
}

// LiquidityBridgeContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidityBridgeContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidityBridgeContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidityBridgeContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityBridgeContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidityBridgeContractSession struct {
	Contract     *LiquidityBridgeContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LiquidityBridgeContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidityBridgeContractCallerSession struct {
	Contract *LiquidityBridgeContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// LiquidityBridgeContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidityBridgeContractTransactorSession struct {
	Contract     *LiquidityBridgeContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// LiquidityBridgeContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidityBridgeContractRaw struct {
	Contract *LiquidityBridgeContract // Generic contract binding to access the raw methods on
}

// LiquidityBridgeContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidityBridgeContractCallerRaw struct {
	Contract *LiquidityBridgeContractCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidityBridgeContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidityBridgeContractTransactorRaw struct {
	Contract *LiquidityBridgeContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidityBridgeContract creates a new instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContract(address common.Address, backend bind.ContractBackend) (*LiquidityBridgeContract, error) {
	contract, err := bindLiquidityBridgeContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContract{LiquidityBridgeContractCaller: LiquidityBridgeContractCaller{contract: contract}, LiquidityBridgeContractTransactor: LiquidityBridgeContractTransactor{contract: contract}, LiquidityBridgeContractFilterer: LiquidityBridgeContractFilterer{contract: contract}}, nil
}

// NewLiquidityBridgeContractCaller creates a new read-only instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractCaller(address common.Address, caller bind.ContractCaller) (*LiquidityBridgeContractCaller, error) {
	contract, err := bindLiquidityBridgeContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCaller{contract: contract}, nil
}

// NewLiquidityBridgeContractTransactor creates a new write-only instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidityBridgeContractTransactor, error) {
	contract, err := bindLiquidityBridgeContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractTransactor{contract: contract}, nil
}

// NewLiquidityBridgeContractFilterer creates a new log filterer instance of LiquidityBridgeContract, bound to a specific deployed contract.
func NewLiquidityBridgeContractFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidityBridgeContractFilterer, error) {
	contract, err := bindLiquidityBridgeContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractFilterer{contract: contract}, nil
}

// bindLiquidityBridgeContract binds a generic wrapper to an already deployed contract.
func bindLiquidityBridgeContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidityBridgeContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidityBridgeContract *LiquidityBridgeContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.LiquidityBridgeContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidityBridgeContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBridgeAddress() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.GetBridgeAddress(&_LiquidityBridgeContract.CallOpts)
}

// GetBridgeAddress is a free data retrieval call binding the contract method 0xfb32c508.
//
// Solidity: function getBridgeAddress() view returns(address)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBridgeAddress() (common.Address, error) {
	return _LiquidityBridgeContract.Contract.GetBridgeAddress(&_LiquidityBridgeContract.CallOpts)
}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetBtcBlockTimestamp(opts *bind.CallOpts, header []byte) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getBtcBlockTimestamp", header)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetBtcBlockTimestamp(header []byte) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBtcBlockTimestamp(&_LiquidityBridgeContract.CallOpts, header)
}

// GetBtcBlockTimestamp is a free data retrieval call binding the contract method 0xa0cd70fc.
//
// Solidity: function getBtcBlockTimestamp(bytes header) pure returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetBtcBlockTimestamp(header []byte) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetBtcBlockTimestamp(&_LiquidityBridgeContract.CallOpts, header)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetCollateral is a free data retrieval call binding the contract method 0x9b56d6c9.
//
// Solidity: function getCollateral(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetCollateral(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetCollateral(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetDustThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getDustThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetDustThreshold() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetDustThreshold(&_LiquidityBridgeContract.CallOpts)
}

// GetDustThreshold is a free data retrieval call binding the contract method 0x33f07ad3.
//
// Solidity: function getDustThreshold() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetDustThreshold() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetDustThreshold(&_LiquidityBridgeContract.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetMinCollateral() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinCollateral(&_LiquidityBridgeContract.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetMinCollateral() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinCollateral(&_LiquidityBridgeContract.CallOpts)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetMinPegIn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getMinPegIn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetMinPegIn() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinPegIn(&_LiquidityBridgeContract.CallOpts)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetMinPegIn() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetMinPegIn(&_LiquidityBridgeContract.CallOpts)
}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetProviders(opts *bind.CallOpts) ([]LiquidityBridgeContractProvider, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getProviders")

	if err != nil {
		return *new([]LiquidityBridgeContractProvider), err
	}

	out0 := *abi.ConvertType(out[0], new([]LiquidityBridgeContractProvider)).(*[]LiquidityBridgeContractProvider)

	return out0, err

}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetProviders() ([]LiquidityBridgeContractProvider, error) {
	return _LiquidityBridgeContract.Contract.GetProviders(&_LiquidityBridgeContract.CallOpts)
}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address)[])
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetProviders() ([]LiquidityBridgeContractProvider, error) {
	return _LiquidityBridgeContract.Contract.GetProviders(&_LiquidityBridgeContract.CallOpts)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetResignDelayBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getResignDelayBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetResignDelayBlocks(&_LiquidityBridgeContract.CallOpts)
}

// GetResignDelayBlocks is a free data retrieval call binding the contract method 0xbd5798c3.
//
// Solidity: function getResignDelayBlocks() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetResignDelayBlocks() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetResignDelayBlocks(&_LiquidityBridgeContract.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetRewardPercentage() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetRewardPercentage(&_LiquidityBridgeContract.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetRewardPercentage() (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetRewardPercentage(&_LiquidityBridgeContract.CallOpts)
}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) HashQuote(opts *bind.CallOpts, quote LiquidityBridgeContractQuote) ([32]byte, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "hashQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// HashQuote is a free data retrieval call binding the contract method 0x1b032188.
//
// Solidity: function hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) HashQuote(quote LiquidityBridgeContractQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) IsOperational(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "isOperational", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) IsOperational(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperational(&_LiquidityBridgeContract.CallOpts, addr)
}

// IsOperational is a free data retrieval call binding the contract method 0x457385f2.
//
// Solidity: function isOperational(address addr) view returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) IsOperational(addr common.Address) (bool, error) {
	return _LiquidityBridgeContract.Contract.IsOperational(&_LiquidityBridgeContract.CallOpts, addr)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) AddCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "addCollateral")
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) AddCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// AddCollateral is a paid mutator transaction binding the contract method 0x9e816999.
//
// Solidity: function addCollateral() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) AddCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.AddCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) CallForUser(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "callForUser", quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.CallForUser(&_LiquidityBridgeContract.TransactOpts, quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xac29d744.
//
// Solidity: function callForUser((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote) payable returns(bool)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) CallForUser(quote LiquidityBridgeContractQuote) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.CallForUser(&_LiquidityBridgeContract.TransactOpts, quote)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Deposit() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Deposit(&_LiquidityBridgeContract.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Deposit() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Deposit(&_LiquidityBridgeContract.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Register() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Register(&_LiquidityBridgeContract.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Register() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Register(&_LiquidityBridgeContract.TransactOpts)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RegisterPegIn(opts *bind.TransactOpts, quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegIn(&_LiquidityBridgeContract.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x6e2e8c70.
//
// Solidity: function registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RegisterPegIn(quote LiquidityBridgeContractQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegIn(&_LiquidityBridgeContract.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Resign(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "resign")
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Resign() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Resign(&_LiquidityBridgeContract.TransactOpts)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Resign() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Resign(&_LiquidityBridgeContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Withdraw(&_LiquidityBridgeContract.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Withdraw(&_LiquidityBridgeContract.TransactOpts, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "withdrawCollateral")
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) WithdrawCollateral() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.WithdrawCollateral(&_LiquidityBridgeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) Receive() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Receive(&_LiquidityBridgeContract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) Receive() (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.Receive(&_LiquidityBridgeContract.TransactOpts)
}

// LiquidityBridgeContractBalanceDecreaseIterator is returned from FilterBalanceDecrease and is used to iterate over the raw logs and unpacked data for BalanceDecrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceDecreaseIterator struct {
	Event *LiquidityBridgeContractBalanceDecrease // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBalanceDecrease)
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
		it.Event = new(LiquidityBridgeContractBalanceDecrease)
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
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBalanceDecrease represents a BalanceDecrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceDecrease is a free log retrieval operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBalanceDecrease(opts *bind.FilterOpts) (*LiquidityBridgeContractBalanceDecreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBalanceDecreaseIterator{contract: _LiquidityBridgeContract.contract, event: "BalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchBalanceDecrease is a free log subscription operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBalanceDecrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBalanceDecrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BalanceDecrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBalanceDecrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBalanceDecrease(log types.Log) (*LiquidityBridgeContractBalanceDecrease, error) {
	event := new(LiquidityBridgeContractBalanceDecrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBalanceIncreaseIterator is returned from FilterBalanceIncrease and is used to iterate over the raw logs and unpacked data for BalanceIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceIncreaseIterator struct {
	Event *LiquidityBridgeContractBalanceIncrease // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBalanceIncrease)
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
		it.Event = new(LiquidityBridgeContractBalanceIncrease)
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
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBalanceIncrease represents a BalanceIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceIncrease is a free log retrieval operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBalanceIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractBalanceIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBalanceIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "BalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchBalanceIncrease is a free log subscription operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBalanceIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBalanceIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BalanceIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBalanceIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBalanceIncrease(log types.Log) (*LiquidityBridgeContractBalanceIncrease, error) {
	event := new(LiquidityBridgeContractBalanceIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBridgeCapExceededIterator is returned from FilterBridgeCapExceeded and is used to iterate over the raw logs and unpacked data for BridgeCapExceeded events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeCapExceededIterator struct {
	Event *LiquidityBridgeContractBridgeCapExceeded // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBridgeCapExceeded)
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
		it.Event = new(LiquidityBridgeContractBridgeCapExceeded)
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
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBridgeCapExceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBridgeCapExceeded represents a BridgeCapExceeded event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeCapExceeded is a free log retrieval operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBridgeCapExceeded(opts *bind.FilterOpts) (*LiquidityBridgeContractBridgeCapExceededIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBridgeCapExceededIterator{contract: _LiquidityBridgeContract.contract, event: "BridgeCapExceeded", logs: logs, sub: sub}, nil
}

// WatchBridgeCapExceeded is a free log subscription operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBridgeCapExceeded(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBridgeCapExceeded) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BridgeCapExceeded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBridgeCapExceeded)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBridgeCapExceeded(log types.Log) (*LiquidityBridgeContractBridgeCapExceeded, error) {
	event := new(LiquidityBridgeContractBridgeCapExceeded)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractBridgeErrorIterator is returned from FilterBridgeError and is used to iterate over the raw logs and unpacked data for BridgeError events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeErrorIterator struct {
	Event *LiquidityBridgeContractBridgeError // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractBridgeErrorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractBridgeError)
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
		it.Event = new(LiquidityBridgeContractBridgeError)
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
func (it *LiquidityBridgeContractBridgeErrorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractBridgeErrorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractBridgeError represents a BridgeError event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractBridgeError struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeError is a free log retrieval operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterBridgeError(opts *bind.FilterOpts) (*LiquidityBridgeContractBridgeErrorIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractBridgeErrorIterator{contract: _LiquidityBridgeContract.contract, event: "BridgeError", logs: logs, sub: sub}, nil
}

// WatchBridgeError is a free log subscription operation binding the contract event 0xa0f8bae2e63548ef07d0f252b12cda04ea27800c1e2605af7b822cdef64e756f.
//
// Solidity: event BridgeError(bytes32 quoteHash, int256 errorCode)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchBridgeError(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractBridgeError) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "BridgeError")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractBridgeError)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeError", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseBridgeError(log types.Log) (*LiquidityBridgeContractBridgeError, error) {
	event := new(LiquidityBridgeContractBridgeError)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "BridgeError", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractCallForUserIterator is returned from FilterCallForUser and is used to iterate over the raw logs and unpacked data for CallForUser events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCallForUserIterator struct {
	Event *LiquidityBridgeContractCallForUser // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractCallForUserIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractCallForUser)
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
		it.Event = new(LiquidityBridgeContractCallForUser)
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
func (it *LiquidityBridgeContractCallForUserIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractCallForUserIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractCallForUser represents a CallForUser event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCallForUser struct {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterCallForUser(opts *bind.FilterOpts) (*LiquidityBridgeContractCallForUserIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "CallForUser")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCallForUserIterator{contract: _LiquidityBridgeContract.contract, event: "CallForUser", logs: logs, sub: sub}, nil
}

// WatchCallForUser is a free log subscription operation binding the contract event 0xbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d.
//
// Solidity: event CallForUser(address from, address dest, uint256 gasLimit, uint256 value, bytes data, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchCallForUser(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractCallForUser) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "CallForUser")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractCallForUser)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CallForUser", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseCallForUser(log types.Log) (*LiquidityBridgeContractCallForUser, error) {
	event := new(LiquidityBridgeContractCallForUser)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CallForUser", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractCollateralIncreaseIterator is returned from FilterCollateralIncrease and is used to iterate over the raw logs and unpacked data for CollateralIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCollateralIncreaseIterator struct {
	Event *LiquidityBridgeContractCollateralIncrease // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractCollateralIncrease)
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
		it.Event = new(LiquidityBridgeContractCollateralIncrease)
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
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractCollateralIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractCollateralIncrease represents a CollateralIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractCollateralIncrease struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCollateralIncrease is a free log retrieval operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterCollateralIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractCollateralIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractCollateralIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "CollateralIncrease", logs: logs, sub: sub}, nil
}

// WatchCollateralIncrease is a free log subscription operation binding the contract event 0x456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84.
//
// Solidity: event CollateralIncrease(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchCollateralIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractCollateralIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "CollateralIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractCollateralIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseCollateralIncrease(log types.Log) (*LiquidityBridgeContractCollateralIncrease, error) {
	event := new(LiquidityBridgeContractCollateralIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "CollateralIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractDepositIterator struct {
	Event *LiquidityBridgeContractDeposit // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractDeposit)
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
		it.Event = new(LiquidityBridgeContractDeposit)
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
func (it *LiquidityBridgeContractDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractDeposit represents a Deposit event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractDeposit struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterDeposit(opts *bind.FilterOpts) (*LiquidityBridgeContractDepositIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractDepositIterator{contract: _LiquidityBridgeContract.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractDeposit) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractDeposit)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Deposit", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseDeposit(log types.Log) (*LiquidityBridgeContractDeposit, error) {
	event := new(LiquidityBridgeContractDeposit)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPenalizedIterator struct {
	Event *LiquidityBridgeContractPenalized // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractPenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPenalized)
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
		it.Event = new(LiquidityBridgeContractPenalized)
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
func (it *LiquidityBridgeContractPenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPenalized represents a Penalized event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPenalized struct {
	LiquidityProvider common.Address
	Penalty           *big.Int
	QuoteHash         [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPenalized(opts *bind.FilterOpts) (*LiquidityBridgeContractPenalizedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPenalizedIterator{contract: _LiquidityBridgeContract.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f4.
//
// Solidity: event Penalized(address liquidityProvider, uint256 penalty, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPenalized) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Penalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPenalized)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Penalized", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePenalized(log types.Log) (*LiquidityBridgeContractPenalized, error) {
	event := new(LiquidityBridgeContractPenalized)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractRefundIterator is returned from FilterRefund and is used to iterate over the raw logs and unpacked data for Refund events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRefundIterator struct {
	Event *LiquidityBridgeContractRefund // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractRefund)
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
		it.Event = new(LiquidityBridgeContractRefund)
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
func (it *LiquidityBridgeContractRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractRefund represents a Refund event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRefund struct {
	Dest      common.Address
	Amount    *big.Int
	Success   bool
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRefund is a free log retrieval operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterRefund(opts *bind.FilterOpts) (*LiquidityBridgeContractRefundIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractRefundIterator{contract: _LiquidityBridgeContract.contract, event: "Refund", logs: logs, sub: sub}, nil
}

// WatchRefund is a free log subscription operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchRefund(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractRefund) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Refund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractRefund)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Refund", log); err != nil {
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

// ParseRefund is a log parse operation binding the contract event 0x3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6.
//
// Solidity: event Refund(address dest, uint256 amount, bool success, bytes32 quoteHash)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseRefund(log types.Log) (*LiquidityBridgeContractRefund, error) {
	event := new(LiquidityBridgeContractRefund)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Refund", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractRegisterIterator is returned from FilterRegister and is used to iterate over the raw logs and unpacked data for Register events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRegisterIterator struct {
	Event *LiquidityBridgeContractRegister // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractRegister)
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
		it.Event = new(LiquidityBridgeContractRegister)
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
func (it *LiquidityBridgeContractRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractRegister represents a Register event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractRegister struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 id, address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterRegister(opts *bind.FilterOpts) (*LiquidityBridgeContractRegisterIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractRegisterIterator{contract: _LiquidityBridgeContract.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 id, address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchRegister(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractRegister) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractRegister)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Register", log); err != nil {
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

// ParseRegister is a log parse operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 id, address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseRegister(log types.Log) (*LiquidityBridgeContractRegister, error) {
	event := new(LiquidityBridgeContractRegister)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Register", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractResignedIterator is returned from FilterResigned and is used to iterate over the raw logs and unpacked data for Resigned events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractResignedIterator struct {
	Event *LiquidityBridgeContractResigned // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractResignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractResigned)
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
		it.Event = new(LiquidityBridgeContractResigned)
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
func (it *LiquidityBridgeContractResignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractResignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractResigned represents a Resigned event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractResigned struct {
	From common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterResigned(opts *bind.FilterOpts) (*LiquidityBridgeContractResignedIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractResignedIterator{contract: _LiquidityBridgeContract.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address from)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchResigned(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractResigned) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Resigned")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractResigned)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Resigned", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseResigned(log types.Log) (*LiquidityBridgeContractResigned, error) {
	event := new(LiquidityBridgeContractResigned)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Resigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractWithdrawCollateralIterator is returned from FilterWithdrawCollateral and is used to iterate over the raw logs and unpacked data for WithdrawCollateral events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawCollateralIterator struct {
	Event *LiquidityBridgeContractWithdrawCollateral // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractWithdrawCollateral)
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
		it.Event = new(LiquidityBridgeContractWithdrawCollateral)
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
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractWithdrawCollateral represents a WithdrawCollateral event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawCollateral struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCollateral is a free log retrieval operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterWithdrawCollateral(opts *bind.FilterOpts) (*LiquidityBridgeContractWithdrawCollateralIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractWithdrawCollateralIterator{contract: _LiquidityBridgeContract.contract, event: "WithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchWithdrawCollateral is a free log subscription operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "WithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractWithdrawCollateral)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseWithdrawCollateral(log types.Log) (*LiquidityBridgeContractWithdrawCollateral, error) {
	event := new(LiquidityBridgeContractWithdrawCollateral)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawalIterator struct {
	Event *LiquidityBridgeContractWithdrawal // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractWithdrawal)
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
		it.Event = new(LiquidityBridgeContractWithdrawal)
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
func (it *LiquidityBridgeContractWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractWithdrawal represents a Withdrawal event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*LiquidityBridgeContractWithdrawalIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractWithdrawalIterator{contract: _LiquidityBridgeContract.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractWithdrawal) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractWithdrawal)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseWithdrawal(log types.Log) (*LiquidityBridgeContractWithdrawal, error) {
	event := new(LiquidityBridgeContractWithdrawal)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SafeMathMetaData contains all meta data concerning the SafeMath contract.
var SafeMathMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220462a3238c6542b1ec5b1d70f766b991a89b511e96d0c38332223b2a45f80e3a864736f6c63430008110033",
}

// SafeMathABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeMathMetaData.ABI instead.
var SafeMathABI = SafeMathMetaData.ABI

// SafeMathBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SafeMathMetaData.Bin instead.
var SafeMathBin = SafeMathMetaData.Bin

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := SafeMathMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// SignatureValidatorMetaData contains all meta data concerning the SignatureValidator contract.
var SignatureValidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1a86b550": "verify(address,bytes32,bytes)",
	},
	Bin: "0x6102a961003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c80631a86b5501461003a575b600080fd5b61004d610048366004610168565b610061565b604051901515815260200160405180910390f35b602081810151604080840151606085015182518084018452601c81527f19457468657265756d205369676e6564204d6573736167653a0a333200000000818701529251600095929391861a9286916100bd9184918b9101610241565b60408051601f1981840301815282825280516020918201206000845290830180835281905260ff861691830191909152606082018790526080820186905291506001600160a01b038a169060019060a0016020604051602081039080840390855afa158015610130573d6000803e3d6000fd5b505050602060405103516001600160a01b031614955050505050509392505050565b634e487b7160e01b600052604160045260246000fd5b60008060006060848603121561017d57600080fd5b83356001600160a01b038116811461019457600080fd5b925060208401359150604084013567ffffffffffffffff808211156101b857600080fd5b818601915086601f8301126101cc57600080fd5b8135818111156101de576101de610152565b604051601f8201601f19908116603f0116810190838211818310171561020657610206610152565b8160405282815289602084870101111561021f57600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b6000835160005b818110156102625760208187018101518583015201610248565b50919091019182525060200191905056fea264697066735822122091317b02258fbb265511a84e53a368712850952762587f61efcbd66500258d0064736f6c63430008110033",
}

// SignatureValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use SignatureValidatorMetaData.ABI instead.
var SignatureValidatorABI = SignatureValidatorMetaData.ABI

// Deprecated: Use SignatureValidatorMetaData.Sigs instead.
// SignatureValidatorFuncSigs maps the 4-byte function signature to its string representation.
var SignatureValidatorFuncSigs = SignatureValidatorMetaData.Sigs

// SignatureValidatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SignatureValidatorMetaData.Bin instead.
var SignatureValidatorBin = SignatureValidatorMetaData.Bin

// DeploySignatureValidator deploys a new Ethereum contract, binding an instance of SignatureValidator to it.
func DeploySignatureValidator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SignatureValidator, error) {
	parsed, err := SignatureValidatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SignatureValidatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SignatureValidator{SignatureValidatorCaller: SignatureValidatorCaller{contract: contract}, SignatureValidatorTransactor: SignatureValidatorTransactor{contract: contract}, SignatureValidatorFilterer: SignatureValidatorFilterer{contract: contract}}, nil
}

// SignatureValidator is an auto generated Go binding around an Ethereum contract.
type SignatureValidator struct {
	SignatureValidatorCaller     // Read-only binding to the contract
	SignatureValidatorTransactor // Write-only binding to the contract
	SignatureValidatorFilterer   // Log filterer for contract events
}

// SignatureValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SignatureValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SignatureValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SignatureValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SignatureValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SignatureValidatorSession struct {
	Contract     *SignatureValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SignatureValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SignatureValidatorCallerSession struct {
	Contract *SignatureValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// SignatureValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SignatureValidatorTransactorSession struct {
	Contract     *SignatureValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// SignatureValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SignatureValidatorRaw struct {
	Contract *SignatureValidator // Generic contract binding to access the raw methods on
}

// SignatureValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SignatureValidatorCallerRaw struct {
	Contract *SignatureValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// SignatureValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SignatureValidatorTransactorRaw struct {
	Contract *SignatureValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSignatureValidator creates a new instance of SignatureValidator, bound to a specific deployed contract.
func NewSignatureValidator(address common.Address, backend bind.ContractBackend) (*SignatureValidator, error) {
	contract, err := bindSignatureValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SignatureValidator{SignatureValidatorCaller: SignatureValidatorCaller{contract: contract}, SignatureValidatorTransactor: SignatureValidatorTransactor{contract: contract}, SignatureValidatorFilterer: SignatureValidatorFilterer{contract: contract}}, nil
}

// NewSignatureValidatorCaller creates a new read-only instance of SignatureValidator, bound to a specific deployed contract.
func NewSignatureValidatorCaller(address common.Address, caller bind.ContractCaller) (*SignatureValidatorCaller, error) {
	contract, err := bindSignatureValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureValidatorCaller{contract: contract}, nil
}

// NewSignatureValidatorTransactor creates a new write-only instance of SignatureValidator, bound to a specific deployed contract.
func NewSignatureValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*SignatureValidatorTransactor, error) {
	contract, err := bindSignatureValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SignatureValidatorTransactor{contract: contract}, nil
}

// NewSignatureValidatorFilterer creates a new log filterer instance of SignatureValidator, bound to a specific deployed contract.
func NewSignatureValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*SignatureValidatorFilterer, error) {
	contract, err := bindSignatureValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SignatureValidatorFilterer{contract: contract}, nil
}

// bindSignatureValidator binds a generic wrapper to an already deployed contract.
func bindSignatureValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SignatureValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureValidator *SignatureValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureValidator.Contract.SignatureValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureValidator *SignatureValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureValidator.Contract.SignatureValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureValidator *SignatureValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureValidator.Contract.SignatureValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SignatureValidator *SignatureValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SignatureValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SignatureValidator *SignatureValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SignatureValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SignatureValidator *SignatureValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SignatureValidator.Contract.contract.Transact(opts, method, params...)
}

// Verify is a free data retrieval call binding the contract method 0x1a86b550.
//
// Solidity: function verify(address addr, bytes32 quoteHash, bytes signature) pure returns(bool)
func (_SignatureValidator *SignatureValidatorCaller) Verify(opts *bind.CallOpts, addr common.Address, quoteHash [32]byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _SignatureValidator.contract.Call(opts, &out, "verify", addr, quoteHash, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Verify is a free data retrieval call binding the contract method 0x1a86b550.
//
// Solidity: function verify(address addr, bytes32 quoteHash, bytes signature) pure returns(bool)
func (_SignatureValidator *SignatureValidatorSession) Verify(addr common.Address, quoteHash [32]byte, signature []byte) (bool, error) {
	return _SignatureValidator.Contract.Verify(&_SignatureValidator.CallOpts, addr, quoteHash, signature)
}

// Verify is a free data retrieval call binding the contract method 0x1a86b550.
//
// Solidity: function verify(address addr, bytes32 quoteHash, bytes signature) pure returns(bool)
func (_SignatureValidator *SignatureValidatorCallerSession) Verify(addr common.Address, quoteHash [32]byte, signature []byte) (bool, error) {
	return _SignatureValidator.Contract.Verify(&_SignatureValidator.CallOpts, addr, quoteHash, signature)
}