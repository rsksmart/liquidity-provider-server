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

// LiquidityBridgeContractPegOutQuote is an auto generated low-level Go binding around an user-defined struct.
type LiquidityBridgeContractPegOutQuote struct {
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	RskRefundAddress            common.Address
	Fee                         uint64
	PenaltyFee                  uint64
	Nonce                       int64
	ValueToTransfer             uint64
	AgreementTimestamp          uint32
	DepositDateLimit            uint32
	DepositConfirmations        uint16
	TransferConfirmations       uint16
	TransferTime                uint32
	ExpireDate                  uint32
	ExpireBlocks                uint32
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minimumCollateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minimumPegIn\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"rewardPercentage\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"resignDelayBlocks\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"dustThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeError\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CollateralIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quotehash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"processed\",\"type\":\"uint256\"}],\"name\":\"PegOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutBalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutBalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"send\",\"type\":\"address\"}],\"name\":\"Test\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"header\",\"type\":\"bytes\"}],\"name\":\"getBtcBlockTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDustThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinPegIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegOutBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"getProcessedQuote\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"valueToTransfer\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlocks\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegoutQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"valueToTransfer\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlocks\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"btcTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"btcBlockHeaderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"partialMerkleTree\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"refundPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"}],\"internalType\":\"structLiquidityBridgeContract.Quote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"fee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"penaltyFee\",\"type\":\"uint64\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint64\",\"name\":\"valueToTransfer\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlocks\",\"type\":\"uint32\"}],\"internalType\":\"structLiquidityBridgeContract.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"registerPegOut\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
		"644be07d": "getPegOutBalance(address)",
		"65fd4a34": "getProcessedQuote(bytes32)",
		"bd5798c3": "getResignDelayBlocks()",
		"c7213163": "getRewardPercentage()",
		"4e43529e": "hashPegoutQuote((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32))",
		"1b032188": "hashQuote((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool))",
		"457385f2": "isOperational(address)",
		"b0e1cb75": "refundPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32),bytes32,bytes32,uint256,bytes32[])",
		"1aa3a008": "register()",
		"6e2e8c70": "registerPegIn((bytes20,address,address,bytes,address,bytes,uint256,uint256,address,bytes,uint32,int64,uint256,uint32,uint32,uint32,uint16,bool),bytes,bytes,bytes,uint256)",
		"797e152d": "registerPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32),bytes)",
		"69652fcf": "resign()",
		"2e1a7d4d": "withdraw(uint256)",
		"59c153be": "withdrawCollateral()",
	},
	Bin: "0x60c06040523480156200001157600080fd5b50604051620037713803806200377183398101604081905262000034916200010b565b60648363ffffffff161115620000905760405162461bcd60e51b815260206004820152601960248201527f496e76616c6964207265776172642070657263656e7461676500000000000000604482015260640160405180910390fd5b600080546001600160a01b039097166001600160a01b03199097169690961790955560809390935260a0919091526006805463ffffffff938416640100000000026001600160401b031990911693909216929092171790556007556200017f565b805163ffffffff811681146200010657600080fd5b919050565b60008060008060008060c087890312156200012557600080fd5b86516001600160a01b03811681146200013d57600080fd5b60208801516040890151919750955093506200015c60608801620000f1565b92506200016c60808801620000f1565b915060a087015190509295509295509295565b60805160a0516135b7620001ba6000396000818161048c01526120510152600081816104230152818161056501526107f901526135b76000f3fe60806040526004361061014f5760003560e01c80639b56d6c9116100b6578063c72131631161006f578063c7213163146103f1578063d0e30db01461040c578063e830b69014610414578063f8b2cb4f14610447578063fa88dcde1461047d578063fb32c508146104b057600080fd5b80639b56d6c91461033d5780639e81699914610373578063a0cd70fc1461037b578063ac29d7441461039b578063b0e1cb75146103ae578063bd5798c3146103ce57600080fd5b806359c153be1161010857806359c153be14610268578063644be07d1461027d57806365fd4a34146102b357806369652fcf146102f55780636e2e8c701461030a578063797e152d1461032a57600080fd5b80631aa3a008146101a85780631b032188146101b05780632e1a7d4d146101e357806333f07ad314610203578063457385f2146102185780634e43529e1461024857600080fd5b366101a3576000546001600160a01b031633146101a15760405162461bcd60e51b815260206004820152600b60248201526a139bdd08185b1b1bddd95960aa1b60448201526064015b60405180910390fd5b005b600080fd5b6101a16104d8565b3480156101bc57600080fd5b506101d06101cb366004612d22565b610676565b6040519081526020015b60405180910390f35b3480156101ef57600080fd5b506101a16101fe366004612d56565b610687565b34801561020f57600080fd5b506007546101d0565b34801561022457600080fd5b50610238610233366004612d6f565b6107cc565b60405190151581526020016101da565b34801561025457600080fd5b506101d0610263366004612eaa565b61081d565b34801561027457600080fd5b506101a1610828565b34801561028957600080fd5b506101d0610298366004612d6f565b6001600160a01b031660009081526002602052604090205490565b3480156102bf57600080fd5b506102e36102ce366004612d56565b6000908152600a602052604090205460ff1690565b60405160ff90911681526020016101da565b34801561030157600080fd5b506101a16109ce565b34801561031657600080fd5b506101d0610325366004612ec7565b610a86565b6101a1610338366004612f7b565b6113c3565b34801561034957600080fd5b506101d0610358366004612d6f565b6001600160a01b031660009081526003602052604090205490565b6101a16116df565b34801561038757600080fd5b506101d0610396366004612fcb565b61175d565b6102386103a9366004612d22565b6117bf565b3480156103ba57600080fd5b506101a16103c9366004612fff565b611b59565b3480156103da57600080fd5b50600654640100000000900463ffffffff166101d0565b3480156103fd57600080fd5b5060065463ffffffff166101d0565b6101a1611e86565b34801561042057600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101d0565b34801561045357600080fd5b506101d0610462366004612d6f565b6001600160a01b031660009081526001602052604090205490565b34801561048957600080fd5b507f00000000000000000000000000000000000000000000000000000000000000006101d0565b3480156104bc57600080fd5b506000546040516001600160a01b0390911681526020016101da565b3233146105115760405162461bcd60e51b81526020600482015260076024820152664e6f7420454f4160c81b6044820152606401610198565b33600090815260036020526040902054156105635760405162461bcd60e51b8152602060048201526012602482015271105b1c9958591e481c9959da5cdd195c995960721b6044820152606401610198565b7f00000000000000000000000000000000000000000000000000000000000000003410156105cb5760405162461bcd60e51b8152602060048201526015602482015274139bdd08195b9bdd59da0818dbdb1b185d195c985b605a1b6044820152606401610198565b33600090815260056020526040902054156106285760405162461bcd60e51b815260206004820152601960248201527f576974686472617720636f6c6c61746572616c206669727374000000000000006044820152606401610198565b3360008181526003602090815260409182902034908190558251938452908301527e7dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a191015b60405180910390a1565b600061068182611eb7565b92915050565b336000908152600160205260409020548111156106db5760405162461bcd60e51b8152602060048201526012602482015271496e73756666696369656e742066756e647360701b6044820152606401610198565b33600090815260016020526040812080548392906106fa9084906130f1565b9091555050604051600090339083908381818185875af1925050503d8060008114610741576040519150601f19603f3d011682016040523d82523d6000602084013e610746565b606091505b505090508061078e5760405162461bcd60e51b815260206004820152601460248201527314d95b991a5b99c8199d5b991cc819985a5b195960621b6044820152606401610198565b60408051338152602081018490527f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b6591015b60405180910390a15050565b60006107d7826120e2565b80156106815750506001600160a01b03166000908152600360205260409020547f0000000000000000000000000000000000000000000000000000000000000000111590565b600061068182612121565b3360009081526005602052604090205461087b5760405162461bcd60e51b815260206004820152601460248201527313995959081d1bc81c995cda59db88199a5c9cdd60621b6044820152606401610198565b6006543360009081526005602052604090205464010000000090910463ffffffff16906108a890436130f1565b10156108ea5760405162461bcd60e51b81526020600482015260116024820152704e6f7420656e6f75676820626c6f636b7360781b6044820152606401610198565b33600081815260036020908152604080832080549084905560059092528083208390555190929083908381818185875af1925050503d806000811461094b576040519150601f19603f3d011682016040523d82523d6000602084013e610950565b606091505b50509050806109985760405162461bcd60e51b815260206004820152601460248201527314d95b991a5b99c8199d5b991cc819985a5b195960621b6044820152606401610198565b60408051338152602081018490527fa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d291016107c0565b6109d7336120e2565b6109f35760405162461bcd60e51b815260040161019890613104565b3360009081526005602052604090205415610a435760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e481c995cda59db995960821b6044820152606401610198565b3360008181526005602090815260409182902043905590519182527fa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86910161066c565b60085460009060ff1615610aac5760405162461bcd60e51b81526004016101989061312c565b6008805460ff191660011790556000610ac487611eb7565b600081815260096020526040902054909150600160ff9091161115610b2b5760405162461bcd60e51b815260206004820152601860248201527f51756f746520616c7265616479207265676973746572656400000000000000006044820152606401610198565b60408088015190516301a86b5560e41b815273__$a15b1ee43228f2a7b300d8dd4e3fc2a5c7$__91631a86b55091610b6a919085908b906004016131a4565b602060405180830381865af4158015610b87573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610bab91906131d4565b610beb5760405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b6044820152606401610198565b637fffffff8310610c3e5760405162461bcd60e51b815260206004820152601e60248201527f486569676874206d757374206265206c6f776572207468616e20325e333100006044820152606401610198565b6000610c4d888787878661217a565b905061012f8101610cb75760405162461bcd60e51b815260206004820152602e60248201527f4572726f72202d3330333a204661696c656420746f2076616c6964617465204260448201526d2a21903a3930b739b0b1ba34b7b760911b6064820152608401610198565b61012e8101610d1a5760405162461bcd60e51b815260206004820152602960248201527f4572726f72202d3330323a205472616e73616374696f6e20616c7265616479206044820152681c1c9bd8d95cdcd95960ba1b6064820152608401610198565b6101308101610d795760405162461bcd60e51b815260206004820152602560248201527f4572726f72202d3330343a205472616e73616374696f6e2076616c7565206973604482015264207a65726f60d81b6064820152608401610198565b6101318101610df05760405162461bcd60e51b815260206004820152603760248201527f4572726f72202d3330353a205472616e73616374696f6e205554584f2076616c60448201527f75652069732062656c6f7720746865206d696e696d756d0000000000000000006064820152608401610198565b6103848101610e415760405162461bcd60e51b815260206004820152601860248201527f4572726f72202d3930303a20427269646765206572726f7200000000000000006044820152606401610198565b6000811380610e51575060c71981145b80610e5d575060631981145b610ea05760405162461bcd60e51b81526020600482015260146024820152732ab735b737bbb710213934b233b29032b93937b960611b6044820152606401610198565b600082815260046020526040902054610ec3908990839063ffffffff1687612223565b15610fb65760e08801516040808a01516001600160a01b031660009081526003602052908120549091610ef5916124a8565b905080600360008b604001516001600160a01b03166001600160a01b031681526020019081526020016000206000828254610f3091906130f1565b909155505060408981015181516001600160a01b039091168152602081018390529081018490527f9685484093cc596fdaeab51abf645b1753dbb7d869bfd2eb21e2c646e47a36f49060600160405180910390a1600654600090606490610f9d9063ffffffff16846131f1565b610fa79190613208565b9050610fb333826124c0565b50505b60c719811480610fc7575060631981145b1561103b576000828152600960209081526040808320805460ff191660021790556004825291829020805464ffffffffff1916905581518481529081018390527ffb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe910160405180910390a191506113b09050565b80611046898261252c565b60008381526004602052604090205463ffffffff16156111ac57600083815260046020526040812054640100000000900460ff16156110a45761109d828b60c001518c6101800151611098919061322a565b6124a8565b90506110b5565b6110b2828b60c001516124a8565b90505b6110c38a60400151826124c0565b60006110cf82846130f1565b90506007548111156111a55760808b01516040516000916001600160a01b0316906108fc90849084818181858888f193505050503d806000811461112f576040519150601f19603f3d011682016040523d82523d6000602084013e611134565b606091505b50506080808e0151604080516001600160a01b0390921682526020820186905283151590820152606081018990529192507f3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6910160405180910390a1806111a3576111a38c60400151836124c0565b505b505061137c565b610220890151819080156111c557508961018001518110155b156112b95760008a61010001516001600160a01b03168b610140015163ffffffff168c61018001518d6101200151604051611200919061323d565b600060405180830381858888f193505050503d806000811461123e576040519150601f19603f3d011682016040523d82523d6000602084013e611243565b606091505b505090507fbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d338c61010001518d61014001518e61018001518f6101200151868b6040516112969796959493929190613259565b60405180910390a180156112b7576101808b01516112b490836130f1565b91505b505b60075481111561137a5760808a01516040516000916001600160a01b0316906108fc90849084818181858888f193505050503d8060008114611317576040519150601f19603f3d011682016040523d82523d6000602084013e61131c565b606091505b50506080808d0151604080516001600160a01b0390921682526020820186905283151590820152606081018890529192507f3052ea2f7e0d74fdc1c1e1f858ff1ae3d91ab1609717c3efedb95db603b255f6910160405180910390a1505b505b506000918252600960209081526040808420805460ff191660021790556004909152909120805464ffffffffff1916905590505b6008805460ff1916905595945050505050565b60085460ff16156113e65760405162461bcd60e51b81526004016101989061312c565b6008805460ff1916600117905560006113fe83612121565b60208401516040516301a86b5560e41b815291925073__$a15b1ee43228f2a7b300d8dd4e3fc2a5c7$__91631a86b5509161143f91859087906004016131a4565b602060405180830381865af415801561145c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061148091906131d4565b6114c55760405162461bcd60e51b81526020600482015260166024820152754c42433a20496e76616c6964207369676e617475726560501b6044820152606401610198565b4283610100015163ffffffff161061151f5760405162461bcd60e51b815260206004820152601b60248201527f4c42433a20426c6f636b20686569676874206f766572666c6f776e00000000006044820152606401610198565b6000818152600a602052604090205460ff166002036115805760405162461bcd60e51b815260206004820152601d60248201527f4c42433a2051756f746520616c726561647920706567676564206f75740000006044820152606401610198565b6000818152600a60205260408120805460ff19166002179055606084015160c08501516115ad91906132ae565b6001600160401b031690508034146116115760405162461bcd60e51b815260206004820152602160248201527f4c42433a206d73672076616c756520646f65736e74206d617463682071756f746044820152606560f81b6064820152608401610198565b804710156116595760405162461bcd60e51b81526020600482015260156024820152744c42433a204e6f7420656e6f7567682066756e647360581b6044820152606401610198565b6116678460400151826125b4565b60c08401516000838152600a60209081526040918290205482513381526001600160401b039094169184019190915282820185905260ff166060830152517fed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec9181900360800190a150506008805460ff191690555050565b6116e8336120e2565b6117045760405162461bcd60e51b815260040161019890613104565b336000908152600360205260408120805434929061172390849061322a565b9091555050604080513381523460208201527f456e0f4ea86ac283092c750200e8c877f6ad8901ae575f90e02081acd455af84910161066c565b600081516050146117a85760405162461bcd60e51b81526020600482015260156024820152740d2dcecc2d8d2c840d0cac2c8cae440d8cadccee8d605b1b6044820152606401610198565b6117b3826044612620565b63ffffffff1692915050565b60006117ca336120e2565b6117e65760405162461bcd60e51b815260040161019890613104565b60085460ff16156118095760405162461bcd60e51b81526004016101989061312c565b6008805460ff191660011790556040820151336001600160a01b03909116146118635760405162461bcd60e51b815260206004820152600c60248201526b155b985d5d1a1bdc9a5e995960a21b6044820152606401610198565b6101808201516040808401516001600160a01b031660009081526001602052205461188f90349061322a565b10156118d25760405162461bcd60e51b8152602060048201526012602482015271496e73756666696369656e742066756e647360701b6044820152606401610198565b60006118dd83611eb7565b60008181526009602052604090205490915060ff161561193f5760405162461bcd60e51b815260206004820152601760248201527f51756f746520616c72656164792070726f6365737365640000000000000000006044820152606401610198565b61194d8360400151346124c0565b610140830151611960906188b8906132d5565b63ffffffff165a10156119a85760405162461bcd60e51b815260206004820152601060248201526f496e73756666696369656e742067617360801b6044820152606401610198565b60008361010001516001600160a01b031684610140015163ffffffff168561018001518661012001516040516119de919061323d565b600060405180830381858888f193505050503d8060008114611a1c576040519150601f19603f3d011682016040523d82523d6000602084013e611a21565b606091505b509091505063ffffffff421115611a7a5760405162461bcd60e51b815260206004820152601860248201527f426c6f636b2074696d657374616d70206f766572666c6f7700000000000000006044820152606401610198565b6000828152600460205260409020805463ffffffff19164263ffffffff161790558015611ad75760008281526004602052604090819020805464ff000000001916640100000000179055840151610180850151611ad79190612718565b7fbfc7404e6fe464f0646fe2c6ab942b92d56be722bb39f8c6bc4830d2d32fb80d338561010001518661014001518761018001518861012001518688604051611b269796959493929190613259565b60405180910390a1600091825260096020526040909120805460ff1916600117905590506008805460ff19169055919050565b60085460ff1615611b7c5760405162461bcd60e51b81526004016101989061312c565b6008805460ff191660011790556000611b9486612121565b6000818152600a602052604090205490915060ff16600214611bf85760405162461bcd60e51b815260206004820152601860248201527f4c42433a2051756f7465206e6f742070726f63657373656400000000000000006044820152606401610198565b85610180015163ffffffff16421115611c535760405162461bcd60e51b815260206004820152601a60248201527f4c42433a2051756f7465206578706972656420627920646174650000000000006044820152606401610198565b856101a0015163ffffffff16431115611cae5760405162461bcd60e51b815260206004820152601c60248201527f4c42433a2051756f7465206578706972656420627920626c6f636b73000000006044820152606401610198565b85602001516001600160a01b0316336001600160a01b031614611d075760405162461bcd60e51b81526020600482015260116024820152702621219d102bb937b7339039b2b73232b960791b6044820152606401610198565b610140860151600054604051635b64458760e01b815261ffff909216916001600160a01b0390911690635b64458790611d4a9089908990899089906004016132f2565b602060405180830381865afa158015611d67573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611d8b919061334e565b1215611de85760405162461bcd60e51b815260206004820152602660248201527f4c42433a20446f6e2774206861766520726571756972656420636f6e6669726d6044820152656174696f6e7360d01b6064820152608401610198565b85602001516001600160a01b03166108fc87606001518860c00151611e0d91906132ae565b6001600160401b03169081150290604051600060405180830381858888f19350505050158015611e41573d6000803e3d6000fd5b50611e5d86604001518760c001516001600160401b0316612784565b6000908152600a60205260409020805460ff199081169091556008805490911690555050505050565b611e8f336120e2565b611eab5760405162461bcd60e51b815260040161019890613104565b611eb533346124c0565b565b600081602001516001600160a01b0316306001600160a01b031614611f125760405162461bcd60e51b815260206004820152601160248201527057726f6e67204c4243206164647265737360781b6044820152606401610198565b6101008201516000546001600160a01b03918216911603611f885760405162461bcd60e51b815260206004820152602a60248201527f427269646765206973206e6f7420616e20616363657074656420636f6e7472616044820152696374206164647265737360b01b6064820152608401610198565b816060015151601514611fee5760405162461bcd60e51b815260206004820152602860248201527f42544320726566756e642061646472657373206d757374206265203231206279604482015267746573206c6f6e6760c01b6064820152608401610198565b8160a001515160151461204f5760405162461bcd60e51b8152602060048201526024808201527f425443204c502061646472657373206d757374206265203231206279746573206044820152636c6f6e6760e01b6064820152608401610198565b7f00000000000000000000000000000000000000000000000000000000000000008260c00151836101800151612085919061322a565b10156120cb5760405162461bcd60e51b8152602060048201526015602482015274151bdbc81b1bddc81859dc99595908185b5bdd5b9d605a1b6044820152606401610198565b6120d4826127f0565b805190602001209050919050565b6001600160a01b038116600090815260036020526040812054158015906106815750506001600160a01b03166000908152600560205260409020541590565b80516000906001600160a01b031630146121715760405162461bcd60e51b815260206004820152601160248201527057726f6e67204c4243206164647265737360781b6044820152606401610198565b6120d48261282b565b60008054606087015160a0880151848452600460208190526040808620549051636adc013360e01b81526001600160a01b0390951694636adc0133946121d6948c948b948d948c94933093909263ffffffff1615159101613367565b6020604051808303816000875af11580156121f5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612219919061334e565b9695505050505050565b6000808413801561224757508460c00151856101800151612244919061322a565b84105b15612254575060006124a0565b6000805460405163bd0c1fff60e01b8152600481018590526001600160a01b039091169063bd0c1fff90602401600060405180830381865afa15801561229e573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526122c691908101906133ea565b905060008151116123105760405162461bcd60e51b8152602060048201526014602482015273125b9d985b1a5908189b1bd8dac81a195a59da1d60621b6044820152606401610198565b600061231b8261175d565b90506000612348886101c0015163ffffffff16896101a0015163ffffffff1661283f90919063ffffffff16565b90508082111561235e57600093505050506124a0565b8560000361237257600193505050506124a0565b600080546102008a01516001600160a01b039091169063bd0c1fff9060019061239f9061ffff168a61322a565b6123a991906130f1565b6040518263ffffffff1660e01b81526004016123c791815260200190565b600060405180830381865afa1580156123e4573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261240c91908101906133ea565b905060008151116124565760405162461bcd60e51b8152602060048201526014602482015273125b9d985b1a5908189b1bd8dac81a195a59da1d60621b6044820152606401610198565b60006124618261175d565b90506124818a6101e0015163ffffffff168261283f90919063ffffffff16565b881115612496576001955050505050506124a0565b6000955050505050505b949350505050565b60008183106124b757816124b9565b825b9392505050565b6001600160a01b038216600090815260016020526040812080548392906124e890849061322a565b9091555050604080516001600160a01b0384168152602081018390527f42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f5391016107c0565b60008260c00151836101800151612543919061322a565b9050600061255361271083613208565b905061255f81836130f1565b8310156125ae5760405162461bcd60e51b815260206004820152601a60248201527f546f6f206c6f77207472616e7366657272656420616d6f756e740000000000006044820152606401610198565b50505050565b6001600160a01b038216600090815260026020526040812080548392906125dc90849061322a565b9091555050604080516001600160a01b0384168152602081018390527f7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae53191016107c0565b600061262d82600461322a565b835110156126745760405162461bcd60e51b8152602060048201526014602482015273736c6963696e67206f7574206f662072616e676560601b6044820152606401610198565b60188361268284600361322a565b8151811061269257612692613457565b016020015160f81c901b6010846126aa85600261322a565b815181106126ba576126ba613457565b016020015160f81c901b6008856126d286600161322a565b815181106126e2576126e2613457565b0160200151865160f89190911c90911b9086908690811061270557612705613457565b016020015160f81c171717905092915050565b6001600160a01b038216600090815260016020526040812080548392906127409084906130f1565b9091555050604080516001600160a01b0384168152602081018390527f8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc006458791016107c0565b6001600160a01b038216600090815260026020526040812080548392906127ac9084906130f1565b9091555050604080516001600160a01b0384168152602081018390527f7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa3833491016107c0565b60606127fb82612857565b6128048361289f565b60405160200161281592919061346d565b6040516020818303038152906040529050919050565b6060612836826128ef565b61280483612973565b6000828201838110156124b957600019915050610681565b6060816000015182602001518360400151846060015185608001518660a001518760c001518860e0015189610100015160405160200161281599989796959493929190613492565b6060816101200151826101400151836101600151846101800151856101a00151866101c00151876101e0015188610200015189610220015160405160200161281599989796959493929190613516565b6060816000015182602001518360400151846060015185608001518660a001518760c0015160405160200161281597969594939291906001600160a01b03978816815295871660208701529390951660408501526001600160401b0391821660608501528116608084015260079390930b60a083015290911660c082015260e00190565b60608160e00151826101000151836101200151846101400151856101600151866101800151876101a00151604051602001612815979695949392919063ffffffff9788168152958716602087015261ffff94851660408701529290931660608501528416608084015290831660a083015290911660c082015260e00190565b634e487b7160e01b600052604160045260246000fd5b60405161024081016001600160401b0381118282101715612a2b57612a2b6129f2565b60405290565b6040516101c081016001600160401b0381118282101715612a2b57612a2b6129f2565b604051601f8201601f191681016001600160401b0381118282101715612a7c57612a7c6129f2565b604052919050565b80356bffffffffffffffffffffffff1981168114612aa157600080fd5b919050565b6001600160a01b0381168114612abb57600080fd5b50565b8035612aa181612aa6565b60006001600160401b03821115612ae257612ae26129f2565b50601f01601f191660200190565b600082601f830112612b0157600080fd5b8135612b14612b0f82612ac9565b612a54565b818152846020838601011115612b2957600080fd5b816020850160208301376000918101602001919091529392505050565b803563ffffffff81168114612aa157600080fd5b8035600781900b8114612aa157600080fd5b803561ffff81168114612aa157600080fd5b8015158114612abb57600080fd5b8035612aa181612b7e565b60006102408284031215612baa57600080fd5b612bb2612a08565b9050612bbd82612a84565b8152612bcb60208301612abe565b6020820152612bdc60408301612abe565b604082015260608201356001600160401b0380821115612bfb57600080fd5b612c0785838601612af0565b6060840152612c1860808501612abe565b608084015260a0840135915080821115612c3157600080fd5b612c3d85838601612af0565b60a084015260c084013560c084015260e084013560e08401526101009150612c66828501612abe565b8284015261012091508184013581811115612c8057600080fd5b612c8c86828701612af0565b83850152505050610140612ca1818401612b46565b90820152610160612cb3838201612b5a565b9082015261018082810135908201526101a0612cd0818401612b46565b908201526101c0612ce2838201612b46565b908201526101e0612cf4838201612b46565b90820152610200612d06838201612b6c565b90820152610220612d18838201612b8c565b9082015292915050565b600060208284031215612d3457600080fd5b81356001600160401b03811115612d4a57600080fd5b6124a084828501612b97565b600060208284031215612d6857600080fd5b5035919050565b600060208284031215612d8157600080fd5b81356124b981612aa6565b80356001600160401b0381168114612aa157600080fd5b60006101c08284031215612db657600080fd5b612dbe612a31565b9050612dc982612abe565b8152612dd760208301612abe565b6020820152612de860408301612abe565b6040820152612df960608301612d8c565b6060820152612e0a60808301612d8c565b6080820152612e1b60a08301612b5a565b60a0820152612e2c60c08301612d8c565b60c0820152612e3d60e08301612b46565b60e0820152610100612e50818401612b46565b90820152610120612e62838201612b6c565b90820152610140612e74838201612b6c565b90820152610160612e86838201612b46565b90820152610180612e98838201612b46565b908201526101a0612d18838201612b46565b60006101c08284031215612ebd57600080fd5b6124b98383612da3565b600080600080600060a08688031215612edf57600080fd5b85356001600160401b0380821115612ef657600080fd5b612f0289838a01612b97565b96506020880135915080821115612f1857600080fd5b612f2489838a01612af0565b95506040880135915080821115612f3a57600080fd5b612f4689838a01612af0565b94506060880135915080821115612f5c57600080fd5b50612f6988828901612af0565b95989497509295608001359392505050565b6000806101e08385031215612f8f57600080fd5b612f998484612da3565b91506101c08301356001600160401b03811115612fb557600080fd5b612fc185828601612af0565b9150509250929050565b600060208284031215612fdd57600080fd5b81356001600160401b03811115612ff357600080fd5b6124a084828501612af0565b6000806000806000610240868803121561301857600080fd5b6130228787612da3565b94506101c086013593506101e0860135925061020086013591506102208601356001600160401b038082111561305757600080fd5b818801915088601f83011261306b57600080fd5b813560208282111561307f5761307f6129f2565b8160051b9250613090818401612a54565b828152928401810192818101908c8511156130aa57600080fd5b948201945b848610156130c8578535825294820194908201906130af565b8096505050505050509295509295909350565b634e487b7160e01b600052601160045260246000fd5b81810381811115610681576106816130db565b6020808252600e908201526d139bdd081c9959da5cdd195c995960921b604082015260600190565b6020808252600e908201526d1499595b9d1c985b9d0818d85b1b60921b604082015260600190565b60005b8381101561316f578181015183820152602001613157565b50506000910152565b60008151808452613190816020860160208601613154565b601f01601f19169290920160200192915050565b60018060a01b03841681528260208201526060604082015260006131cb6060830184613178565b95945050505050565b6000602082840312156131e657600080fd5b81516124b981612b7e565b8082028115828204841417610681576106816130db565b60008261322557634e487b7160e01b600052601260045260246000fd5b500490565b80820180821115610681576106816130db565b6000825161324f818460208701613154565b9190910192915050565b6001600160a01b0388811682528716602082015263ffffffff861660408201526060810185905260e06080820181905260009061329890830186613178565b93151560a08301525060c0015295945050505050565b6001600160401b038181168382160190808211156132ce576132ce6130db565b5092915050565b63ffffffff8181168382160190808211156132ce576132ce6130db565b600060808201868352602086818501528560408501526080606085015281855180845260a086019150828701935060005b8181101561333f57845183529383019391830191600101613323565b50909998505050505050505050565b60006020828403121561336057600080fd5b5051919050565b600061010080835261337b8184018c613178565b90508960208401528281036040840152613395818a613178565b905087606084015282810360808401526133af8188613178565b6001600160a01b03871660a085015283810360c085015290506133d28186613178565b91505082151560e08301529998505050505050505050565b6000602082840312156133fc57600080fd5b81516001600160401b0381111561341257600080fd5b8201601f8101841361342357600080fd5b8051613431612b0f82612ac9565b81815285602083850101111561344657600080fd5b6131cb826020830160208601613154565b634e487b7160e01b600052603260045260246000fd5b6040815260006134806040830185613178565b82810360208401526131cb8185613178565b6bffffffffffffffffffffffff198a1681526001600160a01b0389811660208301528881166040830152610120606083018190526000916134d58483018b613178565b9150808916608085015283820360a08501526134f18289613178565b60c085019790975260e084019590955250509116610100909101529695505050505050565b600061012080835261352a8184018d613178565b63ffffffff9b8c16602085015260079a909a0b604084015250506060810196909652938716608086015291861660a085015290941660c083015261ffff90931660e08201529115156101009092019190915291905056fea2646970667358221220966fc0486c52546b0628687fd0489409c9532c3aec1a4ac827add3f9ec84995c64736f6c63430008110033",
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

// GetPegOutBalance is a free data retrieval call binding the contract method 0x644be07d.
//
// Solidity: function getPegOutBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetPegOutBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getPegOutBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPegOutBalance is a free data retrieval call binding the contract method 0x644be07d.
//
// Solidity: function getPegOutBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetPegOutBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetPegOutBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetPegOutBalance is a free data retrieval call binding the contract method 0x644be07d.
//
// Solidity: function getPegOutBalance(address addr) view returns(uint256)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetPegOutBalance(addr common.Address) (*big.Int, error) {
	return _LiquidityBridgeContract.Contract.GetPegOutBalance(&_LiquidityBridgeContract.CallOpts, addr)
}

// GetProcessedQuote is a free data retrieval call binding the contract method 0x65fd4a34.
//
// Solidity: function getProcessedQuote(bytes32 key) view returns(uint8)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) GetProcessedQuote(opts *bind.CallOpts, key [32]byte) (uint8, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "getProcessedQuote", key)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetProcessedQuote is a free data retrieval call binding the contract method 0x65fd4a34.
//
// Solidity: function getProcessedQuote(bytes32 key) view returns(uint8)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) GetProcessedQuote(key [32]byte) (uint8, error) {
	return _LiquidityBridgeContract.Contract.GetProcessedQuote(&_LiquidityBridgeContract.CallOpts, key)
}

// GetProcessedQuote is a free data retrieval call binding the contract method 0x65fd4a34.
//
// Solidity: function getProcessedQuote(bytes32 key) view returns(uint8)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) GetProcessedQuote(key [32]byte) (uint8, error) {
	return _LiquidityBridgeContract.Contract.GetProcessedQuote(&_LiquidityBridgeContract.CallOpts, key)
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

// HashPegoutQuote is a free data retrieval call binding the contract method 0x4e43529e.
//
// Solidity: function hashPegoutQuote((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCaller) HashPegoutQuote(opts *bind.CallOpts, quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	var out []interface{}
	err := _LiquidityBridgeContract.contract.Call(opts, &out, "hashPegoutQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegoutQuote is a free data retrieval call binding the contract method 0x4e43529e.
//
// Solidity: function hashPegoutQuote((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) HashPegoutQuote(quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashPegoutQuote(&_LiquidityBridgeContract.CallOpts, quote)
}

// HashPegoutQuote is a free data retrieval call binding the contract method 0x4e43529e.
//
// Solidity: function hashPegoutQuote((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote) view returns(bytes32)
func (_LiquidityBridgeContract *LiquidityBridgeContractCallerSession) HashPegoutQuote(quote LiquidityBridgeContractPegOutQuote) ([32]byte, error) {
	return _LiquidityBridgeContract.Contract.HashPegoutQuote(&_LiquidityBridgeContract.CallOpts, quote)
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

// RefundPegOut is a paid mutator transaction binding the contract method 0xb0e1cb75.
//
// Solidity: function refundPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RefundPegOut(opts *bind.TransactOpts, quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "refundPegOut", quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xb0e1cb75.
//
// Solidity: function refundPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RefundPegOut(quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundPegOut(&_LiquidityBridgeContract.TransactOpts, quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xb0e1cb75.
//
// Solidity: function refundPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes32 btcTxHash, bytes32 btcBlockHeaderHash, uint256 partialMerkleTree, bytes32[] merkleBranchHashes) returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RefundPegOut(quote LiquidityBridgeContractPegOutQuote, btcTxHash [32]byte, btcBlockHeaderHash [32]byte, partialMerkleTree *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RefundPegOut(&_LiquidityBridgeContract.TransactOpts, quote, btcTxHash, btcBlockHeaderHash, partialMerkleTree, merkleBranchHashes)
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

// RegisterPegOut is a paid mutator transaction binding the contract method 0x797e152d.
//
// Solidity: function registerPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactor) RegisterPegOut(opts *bind.TransactOpts, quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.contract.Transact(opts, "registerPegOut", quote, signature)
}

// RegisterPegOut is a paid mutator transaction binding the contract method 0x797e152d.
//
// Solidity: function registerPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractSession) RegisterPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
}

// RegisterPegOut is a paid mutator transaction binding the contract method 0x797e152d.
//
// Solidity: function registerPegOut((address,address,address,uint64,uint64,int64,uint64,uint32,uint32,uint16,uint16,uint32,uint32,uint32) quote, bytes signature) payable returns()
func (_LiquidityBridgeContract *LiquidityBridgeContractTransactorSession) RegisterPegOut(quote LiquidityBridgeContractPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _LiquidityBridgeContract.Contract.RegisterPegOut(&_LiquidityBridgeContract.TransactOpts, quote, signature)
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

// LiquidityBridgeContractPegOutIterator is returned from FilterPegOut and is used to iterate over the raw logs and unpacked data for PegOut events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutIterator struct {
	Event *LiquidityBridgeContractPegOut // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractPegOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOut)
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
		it.Event = new(LiquidityBridgeContractPegOut)
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
func (it *LiquidityBridgeContractPegOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOut represents a PegOut event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOut struct {
	From      common.Address
	Amount    *big.Int
	Quotehash [32]byte
	Processed *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPegOut is a free log retrieval operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOut(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOut")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutIterator{contract: _LiquidityBridgeContract.contract, event: "PegOut", logs: logs, sub: sub}, nil
}

// WatchPegOut is a free log subscription operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOut(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOut) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOut)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOut", log); err != nil {
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

// ParsePegOut is a log parse operation binding the contract event 0xed3e6789842b3369f529c844ab6575be53f29ffeabd4d8b84c04c8431535b1ec.
//
// Solidity: event PegOut(address from, uint256 amount, bytes32 quotehash, uint256 processed)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOut(log types.Log) (*LiquidityBridgeContractPegOut, error) {
	event := new(LiquidityBridgeContractPegOut)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutBalanceDecreaseIterator is returned from FilterPegOutBalanceDecrease and is used to iterate over the raw logs and unpacked data for PegOutBalanceDecrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceDecreaseIterator struct {
	Event *LiquidityBridgeContractPegOutBalanceDecrease // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutBalanceDecrease)
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
		it.Event = new(LiquidityBridgeContractPegOutBalanceDecrease)
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
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutBalanceDecrease represents a PegOutBalanceDecrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegOutBalanceDecrease is a free log retrieval operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutBalanceDecrease(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutBalanceDecreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutBalanceDecrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutBalanceDecreaseIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutBalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchPegOutBalanceDecrease is a free log subscription operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutBalanceDecrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutBalanceDecrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutBalanceDecrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutBalanceDecrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceDecrease", log); err != nil {
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

// ParsePegOutBalanceDecrease is a log parse operation binding the contract event 0x7862918831efd3b8f1079c5d7e9344c4a47c2940a11a6b449a03224c4fa38334.
//
// Solidity: event PegOutBalanceDecrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutBalanceDecrease(log types.Log) (*LiquidityBridgeContractPegOutBalanceDecrease, error) {
	event := new(LiquidityBridgeContractPegOutBalanceDecrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityBridgeContractPegOutBalanceIncreaseIterator is returned from FilterPegOutBalanceIncrease and is used to iterate over the raw logs and unpacked data for PegOutBalanceIncrease events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceIncreaseIterator struct {
	Event *LiquidityBridgeContractPegOutBalanceIncrease // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractPegOutBalanceIncrease)
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
		it.Event = new(LiquidityBridgeContractPegOutBalanceIncrease)
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
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractPegOutBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractPegOutBalanceIncrease represents a PegOutBalanceIncrease event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractPegOutBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegOutBalanceIncrease is a free log retrieval operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterPegOutBalanceIncrease(opts *bind.FilterOpts) (*LiquidityBridgeContractPegOutBalanceIncreaseIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "PegOutBalanceIncrease")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractPegOutBalanceIncreaseIterator{contract: _LiquidityBridgeContract.contract, event: "PegOutBalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchPegOutBalanceIncrease is a free log subscription operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchPegOutBalanceIncrease(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractPegOutBalanceIncrease) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "PegOutBalanceIncrease")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractPegOutBalanceIncrease)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceIncrease", log); err != nil {
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

// ParsePegOutBalanceIncrease is a log parse operation binding the contract event 0x7eb93adae7d5cfb024d663c84ccba97b3104e572bbe138ca323c854f666ae531.
//
// Solidity: event PegOutBalanceIncrease(address dest, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParsePegOutBalanceIncrease(log types.Log) (*LiquidityBridgeContractPegOutBalanceIncrease, error) {
	event := new(LiquidityBridgeContractPegOutBalanceIncrease)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "PegOutBalanceIncrease", log); err != nil {
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
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterRegister(opts *bind.FilterOpts) (*LiquidityBridgeContractRegisterIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Register")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractRegisterIterator{contract: _LiquidityBridgeContract.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
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

// ParseRegister is a log parse operation binding the contract event 0x007dc6ab80cc84c043b7b8d4fcafc802187470087f7ea7fccd2e17aecd0256a1.
//
// Solidity: event Register(address from, uint256 amount)
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

// LiquidityBridgeContractTestIterator is returned from FilterTest and is used to iterate over the raw logs and unpacked data for Test events raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractTestIterator struct {
	Event *LiquidityBridgeContractTest // Event containing the contract specifics and raw log

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
func (it *LiquidityBridgeContractTestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityBridgeContractTest)
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
		it.Event = new(LiquidityBridgeContractTest)
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
func (it *LiquidityBridgeContractTestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityBridgeContractTestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityBridgeContractTest represents a Test event raised by the LiquidityBridgeContract contract.
type LiquidityBridgeContractTest struct {
	Send common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterTest is a free log retrieval operation binding the contract event 0xaa9449f2bca09a7b28319d46fd3f3b58a1bb7d94039fc4b69b7bfe5d8535d527.
//
// Solidity: event Test(address send)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) FilterTest(opts *bind.FilterOpts) (*LiquidityBridgeContractTestIterator, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.FilterLogs(opts, "Test")
	if err != nil {
		return nil, err
	}
	return &LiquidityBridgeContractTestIterator{contract: _LiquidityBridgeContract.contract, event: "Test", logs: logs, sub: sub}, nil
}

// WatchTest is a free log subscription operation binding the contract event 0xaa9449f2bca09a7b28319d46fd3f3b58a1bb7d94039fc4b69b7bfe5d8535d527.
//
// Solidity: event Test(address send)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) WatchTest(opts *bind.WatchOpts, sink chan<- *LiquidityBridgeContractTest) (event.Subscription, error) {

	logs, sub, err := _LiquidityBridgeContract.contract.WatchLogs(opts, "Test")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityBridgeContractTest)
				if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Test", log); err != nil {
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

// ParseTest is a log parse operation binding the contract event 0xaa9449f2bca09a7b28319d46fd3f3b58a1bb7d94039fc4b69b7bfe5d8535d527.
//
// Solidity: event Test(address send)
func (_LiquidityBridgeContract *LiquidityBridgeContractFilterer) ParseTest(log types.Log) (*LiquidityBridgeContractTest, error) {
	event := new(LiquidityBridgeContractTest)
	if err := _LiquidityBridgeContract.contract.UnpackLog(event, "Test", log); err != nil {
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
