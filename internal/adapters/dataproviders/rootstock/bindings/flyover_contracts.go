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

// FlyoverLiquidityProvider is an auto generated low-level Go binding around an user-defined struct.
type FlyoverLiquidityProvider struct {
	Id              *big.Int
	ProviderAddress common.Address
	Status          bool
	ProviderType    uint8
	Name            string
	ApiBaseUrl      string
}

// QuotesPegInQuote is an auto generated low-level Go binding around an user-defined struct.
type QuotesPegInQuote struct {
	ChainId                     *big.Int
	CallFee                     *big.Int
	PenaltyFee                  *big.Int
	Value                       *big.Int
	GasFee                      *big.Int
	FedBtcAddress               [20]byte
	LbcAddress                  common.Address
	LiquidityProviderRskAddress common.Address
	ContractAddress             common.Address
	RskRefundAddress            common.Address
	Nonce                       int64
	GasLimit                    uint32
	AgreementTimestamp          uint32
	TimeForDeposit              uint32
	CallTime                    uint32
	DepositConfirmations        uint16
	CallOnRegister              bool
	BtcRefundAddress            []byte
	LiquidityProviderBtcAddress []byte
	Data                        []byte
}

// QuotesPegOutQuote is an auto generated low-level Go binding around an user-defined struct.
type QuotesPegOutQuote struct {
	ChainId               *big.Int
	CallFee               *big.Int
	PenaltyFee            *big.Int
	Value                 *big.Int
	GasFee                *big.Int
	LbcAddress            common.Address
	LpRskAddress          common.Address
	RskRefundAddress      common.Address
	Nonce                 int64
	AgreementTimestamp    uint32
	DepositDateLimit      uint32
	TransferTime          uint32
	ExpireDate            uint32
	ExpireBlock           uint32
	DepositConfirmations  uint16
	TransferConfirmations uint16
	DepositAddress        []byte
	BtcRefundAddress      []byte
	LpBtcAddress          []byte
}

// FlyoverMetaData contains all meta data concerning the Flyover contract.
var FlyoverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"heightOrHash\",\"type\":\"bytes32\"}],\"name\":\"EmptyBlockHeader\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"IncorrectContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"InsufficientAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"InvalidChainId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"NoBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NoContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"passedAmount\",\"type\":\"uint256\"}],\"name\":\"Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"PaymentFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ProviderNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteNotFound\",\"type\":\"error\"}]",
}

// FlyoverABI is the input ABI used to generate the binding from.
// Deprecated: Use FlyoverMetaData.ABI instead.
var FlyoverABI = FlyoverMetaData.ABI

// Flyover is an auto generated Go binding around an Ethereum contract.
type Flyover struct {
	FlyoverCaller     // Read-only binding to the contract
	FlyoverTransactor // Write-only binding to the contract
	FlyoverFilterer   // Log filterer for contract events
}

// FlyoverCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlyoverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlyoverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlyoverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlyoverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlyoverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlyoverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlyoverSession struct {
	Contract     *Flyover          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlyoverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlyoverCallerSession struct {
	Contract *FlyoverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// FlyoverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlyoverTransactorSession struct {
	Contract     *FlyoverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// FlyoverRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlyoverRaw struct {
	Contract *Flyover // Generic contract binding to access the raw methods on
}

// FlyoverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlyoverCallerRaw struct {
	Contract *FlyoverCaller // Generic read-only contract binding to access the raw methods on
}

// FlyoverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlyoverTransactorRaw struct {
	Contract *FlyoverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlyover creates a new instance of Flyover, bound to a specific deployed contract.
func NewFlyover(address common.Address, backend bind.ContractBackend) (*Flyover, error) {
	contract, err := bindFlyover(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Flyover{FlyoverCaller: FlyoverCaller{contract: contract}, FlyoverTransactor: FlyoverTransactor{contract: contract}, FlyoverFilterer: FlyoverFilterer{contract: contract}}, nil
}

// NewFlyoverCaller creates a new read-only instance of Flyover, bound to a specific deployed contract.
func NewFlyoverCaller(address common.Address, caller bind.ContractCaller) (*FlyoverCaller, error) {
	contract, err := bindFlyover(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlyoverCaller{contract: contract}, nil
}

// NewFlyoverTransactor creates a new write-only instance of Flyover, bound to a specific deployed contract.
func NewFlyoverTransactor(address common.Address, transactor bind.ContractTransactor) (*FlyoverTransactor, error) {
	contract, err := bindFlyover(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlyoverTransactor{contract: contract}, nil
}

// NewFlyoverFilterer creates a new log filterer instance of Flyover, bound to a specific deployed contract.
func NewFlyoverFilterer(address common.Address, filterer bind.ContractFilterer) (*FlyoverFilterer, error) {
	contract, err := bindFlyover(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlyoverFilterer{contract: contract}, nil
}

// bindFlyover binds a generic wrapper to an already deployed contract.
func bindFlyover(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FlyoverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flyover *FlyoverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flyover.Contract.FlyoverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flyover *FlyoverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flyover.Contract.FlyoverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flyover *FlyoverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flyover.Contract.FlyoverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flyover *FlyoverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flyover.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flyover *FlyoverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flyover.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flyover *FlyoverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flyover.Contract.contract.Transact(opts, method, params...)
}

// IBridgeMetaData contains all meta data concerning the IBridge contract.
var IBridgeMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rskKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"mstKey\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKeyMultikey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addOneOffLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"txhash\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"addUnlimitedLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"hash\",\"type\":\"bytes\"}],\"name\":\"commitFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveFederationCreationBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActivePowpegRedeemScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestBlockHeader\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestChainHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"depth\",\"type\":\"int256\"}],\"name\":\"getBtcBlockchainBlockHashAtDepth\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"btcBlockHeight\",\"type\":\"uint256\"}],\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainInitialBlockHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainParentBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"getBtcTransactionConfirmations\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"getBtcTxHashProcessedHeight\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePerKb\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"getLockWhitelistEntryByAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockWhitelistSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockingCap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLockTxValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getPendingFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getPendingFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getRetiringFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getRetiringFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForBtcReleaseClient\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForDebugging\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"hasBtcBlockCoinbaseTransactionInformation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"newLockingCap\",\"type\":\"int256\"}],\"name\":\"increaseLockingCap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"isBtcTxHashAlreadyProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ablock\",\"type\":\"bytes\"}],\"name\":\"receiveHeader\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"blocks\",\"type\":\"bytes[]\"}],\"name\":\"receiveHeaders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"witnessMerkleRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"witnessReservedValue\",\"type\":\"bytes32\"}],\"name\":\"registerBtcCoinbaseTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"atx\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"height\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"pmt\",\"type\":\"bytes\"}],\"name\":\"registerBtcTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"shouldTransferToContract\",\"type\":\"bool\"}],\"name\":\"registerFastBridgeBtcTransaction\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"removeLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollbackFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"disableDelay\",\"type\":\"int256\"}],\"name\":\"setLockWhitelistDisableBlockDelay\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateCollections\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"feePerKb\",\"type\":\"int256\"}],\"name\":\"voteFeePerKbChange\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// IBridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use IBridgeMetaData.ABI instead.
var IBridgeABI = IBridgeMetaData.ABI

// IBridge is an auto generated Go binding around an Ethereum contract.
type IBridge struct {
	IBridgeCaller     // Read-only binding to the contract
	IBridgeTransactor // Write-only binding to the contract
	IBridgeFilterer   // Log filterer for contract events
}

// IBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBridgeSession struct {
	Contract     *IBridge          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBridgeCallerSession struct {
	Contract *IBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBridgeTransactorSession struct {
	Contract     *IBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBridgeRaw struct {
	Contract *IBridge // Generic contract binding to access the raw methods on
}

// IBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBridgeCallerRaw struct {
	Contract *IBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// IBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBridgeTransactorRaw struct {
	Contract *IBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBridge creates a new instance of IBridge, bound to a specific deployed contract.
func NewIBridge(address common.Address, backend bind.ContractBackend) (*IBridge, error) {
	contract, err := bindIBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBridge{IBridgeCaller: IBridgeCaller{contract: contract}, IBridgeTransactor: IBridgeTransactor{contract: contract}, IBridgeFilterer: IBridgeFilterer{contract: contract}}, nil
}

// NewIBridgeCaller creates a new read-only instance of IBridge, bound to a specific deployed contract.
func NewIBridgeCaller(address common.Address, caller bind.ContractCaller) (*IBridgeCaller, error) {
	contract, err := bindIBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeCaller{contract: contract}, nil
}

// NewIBridgeTransactor creates a new write-only instance of IBridge, bound to a specific deployed contract.
func NewIBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*IBridgeTransactor, error) {
	contract, err := bindIBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBridgeTransactor{contract: contract}, nil
}

// NewIBridgeFilterer creates a new log filterer instance of IBridge, bound to a specific deployed contract.
func NewIBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*IBridgeFilterer, error) {
	contract, err := bindIBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBridgeFilterer{contract: contract}, nil
}

// bindIBridge binds a generic wrapper to an already deployed contract.
func bindIBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridge *IBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridge.Contract.IBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridge *IBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.Contract.IBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridge *IBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridge.Contract.IBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBridge *IBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBridge *IBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBridge *IBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBridge.Contract.contract.Transact(opts, method, params...)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_IBridge *IBridgeCaller) GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getActiveFederationCreationBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_IBridge *IBridgeSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _IBridge.Contract.GetActiveFederationCreationBlockHeight(&_IBridge.CallOpts)
}

// GetActiveFederationCreationBlockHeight is a free data retrieval call binding the contract method 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (_IBridge *IBridgeCallerSession) GetActiveFederationCreationBlockHeight() (*big.Int, error) {
	return _IBridge.Contract.GetActiveFederationCreationBlockHeight(&_IBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (_IBridge *IBridgeCaller) GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getActivePowpegRedeemScript")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (_IBridge *IBridgeSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _IBridge.Contract.GetActivePowpegRedeemScript(&_IBridge.CallOpts)
}

// GetActivePowpegRedeemScript is a free data retrieval call binding the contract method 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetActivePowpegRedeemScript() ([]byte, error) {
	return _IBridge.Contract.GetActivePowpegRedeemScript(&_IBridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_IBridge *IBridgeCaller) GetBtcBlockchainBestBlockHeader(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainBestBlockHeader")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_IBridge *IBridgeSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBestBlockHeader(&_IBridge.CallOpts)
}

// GetBtcBlockchainBestBlockHeader is a free data retrieval call binding the contract method 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBestBlockHeader(&_IBridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_IBridge *IBridgeCaller) GetBtcBlockchainBestChainHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainBestChainHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_IBridge *IBridgeSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _IBridge.Contract.GetBtcBlockchainBestChainHeight(&_IBridge.CallOpts)
}

// GetBtcBlockchainBestChainHeight is a free data retrieval call binding the contract method 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainBestChainHeight() (*big.Int, error) {
	return _IBridge.Contract.GetBtcBlockchainBestChainHeight(&_IBridge.CallOpts)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_IBridge *IBridgeCaller) GetBtcBlockchainBlockHashAtDepth(opts *bind.CallOpts, depth *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHashAtDepth", depth)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_IBridge *IBridgeSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_IBridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHashAtDepth is a free data retrieval call binding the contract method 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHashAtDepth(&_IBridge.CallOpts, depth)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeCaller) GetBtcBlockchainBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_IBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHash is a free data retrieval call binding the contract method 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHeaderByHash(&_IBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_IBridge *IBridgeCaller) GetBtcBlockchainBlockHeaderByHeight(opts *bind.CallOpts, btcBlockHeight *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainBlockHeaderByHeight", btcBlockHeight)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_IBridge *IBridgeSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_IBridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainBlockHeaderByHeight is a free data retrieval call binding the contract method 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainBlockHeaderByHeight(&_IBridge.CallOpts, btcBlockHeight)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_IBridge *IBridgeCaller) GetBtcBlockchainInitialBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainInitialBlockHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_IBridge *IBridgeSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _IBridge.Contract.GetBtcBlockchainInitialBlockHeight(&_IBridge.CallOpts)
}

// GetBtcBlockchainInitialBlockHeight is a free data retrieval call binding the contract method 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainInitialBlockHeight() (*big.Int, error) {
	return _IBridge.Contract.GetBtcBlockchainInitialBlockHeight(&_IBridge.CallOpts)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeCaller) GetBtcBlockchainParentBlockHeaderByHash(opts *bind.CallOpts, btcBlockHash [32]byte) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcBlockchainParentBlockHeaderByHash", btcBlockHash)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_IBridge.CallOpts, btcBlockHash)
}

// GetBtcBlockchainParentBlockHeaderByHash is a free data retrieval call binding the contract method 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return _IBridge.Contract.GetBtcBlockchainParentBlockHeaderByHash(&_IBridge.CallOpts, btcBlockHash)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_IBridge *IBridgeCaller) GetBtcTransactionConfirmations(opts *bind.CallOpts, txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcTransactionConfirmations", txHash, blockHash, merkleBranchPath, merkleBranchHashes)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_IBridge *IBridgeSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _IBridge.Contract.GetBtcTransactionConfirmations(&_IBridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTransactionConfirmations is a free data retrieval call binding the contract method 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (_IBridge *IBridgeCallerSession) GetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*big.Int, error) {
	return _IBridge.Contract.GetBtcTransactionConfirmations(&_IBridge.CallOpts, txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_IBridge *IBridgeCaller) GetBtcTxHashProcessedHeight(opts *bind.CallOpts, hash string) (int64, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getBtcTxHashProcessedHeight", hash)

	if err != nil {
		return *new(int64), err
	}

	out0 := *abi.ConvertType(out[0], new(int64)).(*int64)

	return out0, err

}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_IBridge *IBridgeSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _IBridge.Contract.GetBtcTxHashProcessedHeight(&_IBridge.CallOpts, hash)
}

// GetBtcTxHashProcessedHeight is a free data retrieval call binding the contract method 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (_IBridge *IBridgeCallerSession) GetBtcTxHashProcessedHeight(hash string) (int64, error) {
	return _IBridge.Contract.GetBtcTxHashProcessedHeight(&_IBridge.CallOpts, hash)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_IBridge *IBridgeCaller) GetFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_IBridge *IBridgeSession) GetFederationAddress() (string, error) {
	return _IBridge.Contract.GetFederationAddress(&_IBridge.CallOpts)
}

// GetFederationAddress is a free data retrieval call binding the contract method 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (_IBridge *IBridgeCallerSession) GetFederationAddress() (string, error) {
	return _IBridge.Contract.GetFederationAddress(&_IBridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeCaller) GetFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _IBridge.Contract.GetFederationCreationBlockNumber(&_IBridge.CallOpts)
}

// GetFederationCreationBlockNumber is a free data retrieval call binding the contract method 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetFederationCreationBlockNumber() (*big.Int, error) {
	return _IBridge.Contract.GetFederationCreationBlockNumber(&_IBridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeCaller) GetFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeSession) GetFederationCreationTime() (*big.Int, error) {
	return _IBridge.Contract.GetFederationCreationTime(&_IBridge.CallOpts)
}

// GetFederationCreationTime is a free data retrieval call binding the contract method 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetFederationCreationTime() (*big.Int, error) {
	return _IBridge.Contract.GetFederationCreationTime(&_IBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_IBridge *IBridgeCaller) GetFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_IBridge *IBridgeSession) GetFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetFederationSize(&_IBridge.CallOpts)
}

// GetFederationSize is a free data retrieval call binding the contract method 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetFederationSize(&_IBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_IBridge *IBridgeCaller) GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_IBridge *IBridgeSession) GetFederationThreshold() (*big.Int, error) {
	return _IBridge.Contract.GetFederationThreshold(&_IBridge.CallOpts)
}

// GetFederationThreshold is a free data retrieval call binding the contract method 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetFederationThreshold() (*big.Int, error) {
	return _IBridge.Contract.GetFederationThreshold(&_IBridge.CallOpts)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCaller) GetFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetFederatorPublicKey is a free data retrieval call binding the contract method 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCaller) GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_IBridge *IBridgeCaller) GetFeePerKb(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getFeePerKb")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_IBridge *IBridgeSession) GetFeePerKb() (*big.Int, error) {
	return _IBridge.Contract.GetFeePerKb(&_IBridge.CallOpts)
}

// GetFeePerKb is a free data retrieval call binding the contract method 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetFeePerKb() (*big.Int, error) {
	return _IBridge.Contract.GetFeePerKb(&_IBridge.CallOpts)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_IBridge *IBridgeCaller) GetLockWhitelistAddress(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getLockWhitelistAddress", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_IBridge *IBridgeSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _IBridge.Contract.GetLockWhitelistAddress(&_IBridge.CallOpts, index)
}

// GetLockWhitelistAddress is a free data retrieval call binding the contract method 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (_IBridge *IBridgeCallerSession) GetLockWhitelistAddress(index *big.Int) (string, error) {
	return _IBridge.Contract.GetLockWhitelistAddress(&_IBridge.CallOpts, index)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_IBridge *IBridgeCaller) GetLockWhitelistEntryByAddress(opts *bind.CallOpts, aaddress string) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getLockWhitelistEntryByAddress", aaddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_IBridge *IBridgeSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _IBridge.Contract.GetLockWhitelistEntryByAddress(&_IBridge.CallOpts, aaddress)
}

// GetLockWhitelistEntryByAddress is a free data retrieval call binding the contract method 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (_IBridge *IBridgeCallerSession) GetLockWhitelistEntryByAddress(aaddress string) (*big.Int, error) {
	return _IBridge.Contract.GetLockWhitelistEntryByAddress(&_IBridge.CallOpts, aaddress)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_IBridge *IBridgeCaller) GetLockWhitelistSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getLockWhitelistSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_IBridge *IBridgeSession) GetLockWhitelistSize() (*big.Int, error) {
	return _IBridge.Contract.GetLockWhitelistSize(&_IBridge.CallOpts)
}

// GetLockWhitelistSize is a free data retrieval call binding the contract method 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetLockWhitelistSize() (*big.Int, error) {
	return _IBridge.Contract.GetLockWhitelistSize(&_IBridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_IBridge *IBridgeCaller) GetLockingCap(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getLockingCap")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_IBridge *IBridgeSession) GetLockingCap() (*big.Int, error) {
	return _IBridge.Contract.GetLockingCap(&_IBridge.CallOpts)
}

// GetLockingCap is a free data retrieval call binding the contract method 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetLockingCap() (*big.Int, error) {
	return _IBridge.Contract.GetLockingCap(&_IBridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_IBridge *IBridgeCaller) GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getMinimumLockTxValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_IBridge *IBridgeSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _IBridge.Contract.GetMinimumLockTxValue(&_IBridge.CallOpts)
}

// GetMinimumLockTxValue is a free data retrieval call binding the contract method 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetMinimumLockTxValue() (*big.Int, error) {
	return _IBridge.Contract.GetMinimumLockTxValue(&_IBridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_IBridge *IBridgeCaller) GetPendingFederationHash(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getPendingFederationHash")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_IBridge *IBridgeSession) GetPendingFederationHash() ([]byte, error) {
	return _IBridge.Contract.GetPendingFederationHash(&_IBridge.CallOpts)
}

// GetPendingFederationHash is a free data retrieval call binding the contract method 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetPendingFederationHash() ([]byte, error) {
	return _IBridge.Contract.GetPendingFederationHash(&_IBridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_IBridge *IBridgeCaller) GetPendingFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getPendingFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_IBridge *IBridgeSession) GetPendingFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetPendingFederationSize(&_IBridge.CallOpts)
}

// GetPendingFederationSize is a free data retrieval call binding the contract method 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetPendingFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetPendingFederationSize(&_IBridge.CallOpts)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCaller) GetPendingFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getPendingFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetPendingFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetPendingFederatorPublicKey is a free data retrieval call binding the contract method 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetPendingFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCaller) GetPendingFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getPendingFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetPendingFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetPendingFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetPendingFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_IBridge *IBridgeCaller) GetRetiringFederationAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederationAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_IBridge *IBridgeSession) GetRetiringFederationAddress() (string, error) {
	return _IBridge.Contract.GetRetiringFederationAddress(&_IBridge.CallOpts)
}

// GetRetiringFederationAddress is a free data retrieval call binding the contract method 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (_IBridge *IBridgeCallerSession) GetRetiringFederationAddress() (string, error) {
	return _IBridge.Contract.GetRetiringFederationAddress(&_IBridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeCaller) GetRetiringFederationCreationBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederationCreationBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationCreationBlockNumber(&_IBridge.CallOpts)
}

// GetRetiringFederationCreationBlockNumber is a free data retrieval call binding the contract method 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetRetiringFederationCreationBlockNumber() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationCreationBlockNumber(&_IBridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeCaller) GetRetiringFederationCreationTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederationCreationTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationCreationTime(&_IBridge.CallOpts)
}

// GetRetiringFederationCreationTime is a free data retrieval call binding the contract method 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetRetiringFederationCreationTime() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationCreationTime(&_IBridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_IBridge *IBridgeCaller) GetRetiringFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederationSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_IBridge *IBridgeSession) GetRetiringFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationSize(&_IBridge.CallOpts)
}

// GetRetiringFederationSize is a free data retrieval call binding the contract method 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetRetiringFederationSize() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationSize(&_IBridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_IBridge *IBridgeCaller) GetRetiringFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_IBridge *IBridgeSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationThreshold(&_IBridge.CallOpts)
}

// GetRetiringFederationThreshold is a free data retrieval call binding the contract method 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (_IBridge *IBridgeCallerSession) GetRetiringFederationThreshold() (*big.Int, error) {
	return _IBridge.Contract.GetRetiringFederationThreshold(&_IBridge.CallOpts)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCaller) GetRetiringFederatorPublicKey(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederatorPublicKey", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetRetiringFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetRetiringFederatorPublicKey is a free data retrieval call binding the contract method 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return _IBridge.Contract.GetRetiringFederatorPublicKey(&_IBridge.CallOpts, index)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCaller) GetRetiringFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getRetiringFederatorPublicKeyOfType", index, atype)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetRetiringFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetRetiringFederatorPublicKeyOfType is a free data retrieval call binding the contract method 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return _IBridge.Contract.GetRetiringFederatorPublicKeyOfType(&_IBridge.CallOpts, index, atype)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_IBridge *IBridgeCaller) GetStateForBtcReleaseClient(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getStateForBtcReleaseClient")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_IBridge *IBridgeSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _IBridge.Contract.GetStateForBtcReleaseClient(&_IBridge.CallOpts)
}

// GetStateForBtcReleaseClient is a free data retrieval call binding the contract method 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetStateForBtcReleaseClient() ([]byte, error) {
	return _IBridge.Contract.GetStateForBtcReleaseClient(&_IBridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_IBridge *IBridgeCaller) GetStateForDebugging(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "getStateForDebugging")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_IBridge *IBridgeSession) GetStateForDebugging() ([]byte, error) {
	return _IBridge.Contract.GetStateForDebugging(&_IBridge.CallOpts)
}

// GetStateForDebugging is a free data retrieval call binding the contract method 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (_IBridge *IBridgeCallerSession) GetStateForDebugging() ([]byte, error) {
	return _IBridge.Contract.GetStateForDebugging(&_IBridge.CallOpts)
}

// HasBtcBlockCoinbaseTransactionInformation is a free data retrieval call binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (_IBridge *IBridgeCaller) HasBtcBlockCoinbaseTransactionInformation(opts *bind.CallOpts, blockHash [32]byte) (bool, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "hasBtcBlockCoinbaseTransactionInformation", blockHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasBtcBlockCoinbaseTransactionInformation is a free data retrieval call binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (_IBridge *IBridgeSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (bool, error) {
	return _IBridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_IBridge.CallOpts, blockHash)
}

// HasBtcBlockCoinbaseTransactionInformation is a free data retrieval call binding the contract method 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (_IBridge *IBridgeCallerSession) HasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) (bool, error) {
	return _IBridge.Contract.HasBtcBlockCoinbaseTransactionInformation(&_IBridge.CallOpts, blockHash)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_IBridge *IBridgeCaller) IsBtcTxHashAlreadyProcessed(opts *bind.CallOpts, hash string) (bool, error) {
	var out []interface{}
	err := _IBridge.contract.Call(opts, &out, "isBtcTxHashAlreadyProcessed", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_IBridge *IBridgeSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _IBridge.Contract.IsBtcTxHashAlreadyProcessed(&_IBridge.CallOpts, hash)
}

// IsBtcTxHashAlreadyProcessed is a free data retrieval call binding the contract method 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (_IBridge *IBridgeCallerSession) IsBtcTxHashAlreadyProcessed(hash string) (bool, error) {
	return _IBridge.Contract.IsBtcTxHashAlreadyProcessed(&_IBridge.CallOpts, hash)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_IBridge *IBridgeTransactor) AddFederatorPublicKey(opts *bind.TransactOpts, key []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addFederatorPublicKey", key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_IBridge *IBridgeSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddFederatorPublicKey(&_IBridge.TransactOpts, key)
}

// AddFederatorPublicKey is a paid mutator transaction binding the contract method 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (_IBridge *IBridgeTransactorSession) AddFederatorPublicKey(key []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddFederatorPublicKey(&_IBridge.TransactOpts, key)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_IBridge *IBridgeTransactor) AddFederatorPublicKeyMultikey(opts *bind.TransactOpts, btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addFederatorPublicKeyMultikey", btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_IBridge *IBridgeSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddFederatorPublicKeyMultikey(&_IBridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddFederatorPublicKeyMultikey is a paid mutator transaction binding the contract method 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (_IBridge *IBridgeTransactorSession) AddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddFederatorPublicKeyMultikey(&_IBridge.TransactOpts, btcKey, rskKey, mstKey)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeTransactor) AddLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.AddLockWhitelistAddress(&_IBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddLockWhitelistAddress is a paid mutator transaction binding the contract method 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeTransactorSession) AddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.AddLockWhitelistAddress(&_IBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeTransactor) AddOneOffLockWhitelistAddress(opts *bind.TransactOpts, aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addOneOffLockWhitelistAddress", aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.AddOneOffLockWhitelistAddress(&_IBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddOneOffLockWhitelistAddress is a paid mutator transaction binding the contract method 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (_IBridge *IBridgeTransactorSession) AddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.AddOneOffLockWhitelistAddress(&_IBridge.TransactOpts, aaddress, maxTransferValue)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_IBridge *IBridgeTransactor) AddSignature(opts *bind.TransactOpts, pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addSignature", pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_IBridge *IBridgeSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddSignature(&_IBridge.TransactOpts, pubkey, signatures, txhash)
}

// AddSignature is a paid mutator transaction binding the contract method 0xf10b9c59.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (_IBridge *IBridgeTransactorSession) AddSignature(pubkey []byte, signatures [][]byte, txhash []byte) (*types.Transaction, error) {
	return _IBridge.Contract.AddSignature(&_IBridge.TransactOpts, pubkey, signatures, txhash)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeTransactor) AddUnlimitedLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "addUnlimitedLockWhitelistAddress", aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _IBridge.Contract.AddUnlimitedLockWhitelistAddress(&_IBridge.TransactOpts, aaddress)
}

// AddUnlimitedLockWhitelistAddress is a paid mutator transaction binding the contract method 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeTransactorSession) AddUnlimitedLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _IBridge.Contract.AddUnlimitedLockWhitelistAddress(&_IBridge.TransactOpts, aaddress)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_IBridge *IBridgeTransactor) CommitFederation(opts *bind.TransactOpts, hash []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "commitFederation", hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_IBridge *IBridgeSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _IBridge.Contract.CommitFederation(&_IBridge.TransactOpts, hash)
}

// CommitFederation is a paid mutator transaction binding the contract method 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (_IBridge *IBridgeTransactorSession) CommitFederation(hash []byte) (*types.Transaction, error) {
	return _IBridge.Contract.CommitFederation(&_IBridge.TransactOpts, hash)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_IBridge *IBridgeTransactor) CreateFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "createFederation")
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_IBridge *IBridgeSession) CreateFederation() (*types.Transaction, error) {
	return _IBridge.Contract.CreateFederation(&_IBridge.TransactOpts)
}

// CreateFederation is a paid mutator transaction binding the contract method 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (_IBridge *IBridgeTransactorSession) CreateFederation() (*types.Transaction, error) {
	return _IBridge.Contract.CreateFederation(&_IBridge.TransactOpts)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_IBridge *IBridgeTransactor) IncreaseLockingCap(opts *bind.TransactOpts, newLockingCap *big.Int) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "increaseLockingCap", newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_IBridge *IBridgeSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.IncreaseLockingCap(&_IBridge.TransactOpts, newLockingCap)
}

// IncreaseLockingCap is a paid mutator transaction binding the contract method 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (_IBridge *IBridgeTransactorSession) IncreaseLockingCap(newLockingCap *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.IncreaseLockingCap(&_IBridge.TransactOpts, newLockingCap)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_IBridge *IBridgeTransactor) ReceiveHeader(opts *bind.TransactOpts, ablock []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "receiveHeader", ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_IBridge *IBridgeSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _IBridge.Contract.ReceiveHeader(&_IBridge.TransactOpts, ablock)
}

// ReceiveHeader is a paid mutator transaction binding the contract method 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (_IBridge *IBridgeTransactorSession) ReceiveHeader(ablock []byte) (*types.Transaction, error) {
	return _IBridge.Contract.ReceiveHeader(&_IBridge.TransactOpts, ablock)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_IBridge *IBridgeTransactor) ReceiveHeaders(opts *bind.TransactOpts, blocks [][]byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "receiveHeaders", blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_IBridge *IBridgeSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _IBridge.Contract.ReceiveHeaders(&_IBridge.TransactOpts, blocks)
}

// ReceiveHeaders is a paid mutator transaction binding the contract method 0xe5400e7b.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (_IBridge *IBridgeTransactorSession) ReceiveHeaders(blocks [][]byte) (*types.Transaction, error) {
	return _IBridge.Contract.ReceiveHeaders(&_IBridge.TransactOpts, blocks)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_IBridge *IBridgeTransactor) RegisterBtcCoinbaseTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "registerBtcCoinbaseTransaction", btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_IBridge *IBridgeSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterBtcCoinbaseTransaction(&_IBridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcCoinbaseTransaction is a paid mutator transaction binding the contract method 0xccf417ae.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (_IBridge *IBridgeTransactorSession) RegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterBtcCoinbaseTransaction(&_IBridge.TransactOpts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_IBridge *IBridgeTransactor) RegisterBtcTransaction(opts *bind.TransactOpts, atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "registerBtcTransaction", atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_IBridge *IBridgeSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterBtcTransaction(&_IBridge.TransactOpts, atx, height, pmt)
}

// RegisterBtcTransaction is a paid mutator transaction binding the contract method 0x43dc0656.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (_IBridge *IBridgeTransactorSession) RegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterBtcTransaction(&_IBridge.TransactOpts, atx, height, pmt)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_IBridge *IBridgeTransactor) RegisterFastBridgeBtcTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_IBridge *IBridgeSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterFastBridgeBtcTransaction(&_IBridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RegisterFastBridgeBtcTransaction is a paid mutator transaction binding the contract method 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (_IBridge *IBridgeTransactorSession) RegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) (*types.Transaction, error) {
	return _IBridge.Contract.RegisterFastBridgeBtcTransaction(&_IBridge.TransactOpts, btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeTransactor) RemoveLockWhitelistAddress(opts *bind.TransactOpts, aaddress string) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "removeLockWhitelistAddress", aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _IBridge.Contract.RemoveLockWhitelistAddress(&_IBridge.TransactOpts, aaddress)
}

// RemoveLockWhitelistAddress is a paid mutator transaction binding the contract method 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (_IBridge *IBridgeTransactorSession) RemoveLockWhitelistAddress(aaddress string) (*types.Transaction, error) {
	return _IBridge.Contract.RemoveLockWhitelistAddress(&_IBridge.TransactOpts, aaddress)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_IBridge *IBridgeTransactor) RollbackFederation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "rollbackFederation")
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_IBridge *IBridgeSession) RollbackFederation() (*types.Transaction, error) {
	return _IBridge.Contract.RollbackFederation(&_IBridge.TransactOpts)
}

// RollbackFederation is a paid mutator transaction binding the contract method 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (_IBridge *IBridgeTransactorSession) RollbackFederation() (*types.Transaction, error) {
	return _IBridge.Contract.RollbackFederation(&_IBridge.TransactOpts)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_IBridge *IBridgeTransactor) SetLockWhitelistDisableBlockDelay(opts *bind.TransactOpts, disableDelay *big.Int) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "setLockWhitelistDisableBlockDelay", disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_IBridge *IBridgeSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.SetLockWhitelistDisableBlockDelay(&_IBridge.TransactOpts, disableDelay)
}

// SetLockWhitelistDisableBlockDelay is a paid mutator transaction binding the contract method 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (_IBridge *IBridgeTransactorSession) SetLockWhitelistDisableBlockDelay(disableDelay *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.SetLockWhitelistDisableBlockDelay(&_IBridge.TransactOpts, disableDelay)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_IBridge *IBridgeTransactor) UpdateCollections(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "updateCollections")
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_IBridge *IBridgeSession) UpdateCollections() (*types.Transaction, error) {
	return _IBridge.Contract.UpdateCollections(&_IBridge.TransactOpts)
}

// UpdateCollections is a paid mutator transaction binding the contract method 0x0c5a9990.
//
// Solidity: function updateCollections() returns()
func (_IBridge *IBridgeTransactorSession) UpdateCollections() (*types.Transaction, error) {
	return _IBridge.Contract.UpdateCollections(&_IBridge.TransactOpts)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_IBridge *IBridgeTransactor) VoteFeePerKbChange(opts *bind.TransactOpts, feePerKb *big.Int) (*types.Transaction, error) {
	return _IBridge.contract.Transact(opts, "voteFeePerKbChange", feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_IBridge *IBridgeSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.VoteFeePerKbChange(&_IBridge.TransactOpts, feePerKb)
}

// VoteFeePerKbChange is a paid mutator transaction binding the contract method 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (_IBridge *IBridgeTransactorSession) VoteFeePerKbChange(feePerKb *big.Int) (*types.Transaction, error) {
	return _IBridge.Contract.VoteFeePerKbChange(&_IBridge.TransactOpts, feePerKb)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IBridge *IBridgeTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBridge.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IBridge *IBridgeSession) Receive() (*types.Transaction, error) {
	return _IBridge.Contract.Receive(&_IBridge.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IBridge *IBridgeTransactorSession) Receive() (*types.Transaction, error) {
	return _IBridge.Contract.Receive(&_IBridge.TransactOpts)
}

// ICollateralManagementMetaData contains all meta data concerning the ICollateralManagement contract.
var ICollateralManagementMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"addPegInCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addPegInCollateralTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addPegOutCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addPegOutCollateralTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegInCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegOutCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPenalties\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayInBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getResignationBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isCollateralSufficient\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isPaused\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"since\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"slashPegInCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"slashPegOutCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegInCollateralAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutCollateralAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"collateralType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"AlreadyResigned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotResigned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NothingToWithdraw\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"resignationBlockNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"resignDelayInBlocks\",\"type\":\"uint256\"}],\"name\":\"ResignationDelayNotMet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawalFailed\",\"type\":\"error\"}]",
}

// ICollateralManagementABI is the input ABI used to generate the binding from.
// Deprecated: Use ICollateralManagementMetaData.ABI instead.
var ICollateralManagementABI = ICollateralManagementMetaData.ABI

// ICollateralManagement is an auto generated Go binding around an Ethereum contract.
type ICollateralManagement struct {
	ICollateralManagementCaller     // Read-only binding to the contract
	ICollateralManagementTransactor // Write-only binding to the contract
	ICollateralManagementFilterer   // Log filterer for contract events
}

// ICollateralManagementCaller is an auto generated read-only Go binding around an Ethereum contract.
type ICollateralManagementCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICollateralManagementTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ICollateralManagementTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICollateralManagementFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ICollateralManagementFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICollateralManagementSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ICollateralManagementSession struct {
	Contract     *ICollateralManagement // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ICollateralManagementCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ICollateralManagementCallerSession struct {
	Contract *ICollateralManagementCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ICollateralManagementTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ICollateralManagementTransactorSession struct {
	Contract     *ICollateralManagementTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ICollateralManagementRaw is an auto generated low-level Go binding around an Ethereum contract.
type ICollateralManagementRaw struct {
	Contract *ICollateralManagement // Generic contract binding to access the raw methods on
}

// ICollateralManagementCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ICollateralManagementCallerRaw struct {
	Contract *ICollateralManagementCaller // Generic read-only contract binding to access the raw methods on
}

// ICollateralManagementTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ICollateralManagementTransactorRaw struct {
	Contract *ICollateralManagementTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICollateralManagement creates a new instance of ICollateralManagement, bound to a specific deployed contract.
func NewICollateralManagement(address common.Address, backend bind.ContractBackend) (*ICollateralManagement, error) {
	contract, err := bindICollateralManagement(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagement{ICollateralManagementCaller: ICollateralManagementCaller{contract: contract}, ICollateralManagementTransactor: ICollateralManagementTransactor{contract: contract}, ICollateralManagementFilterer: ICollateralManagementFilterer{contract: contract}}, nil
}

// NewICollateralManagementCaller creates a new read-only instance of ICollateralManagement, bound to a specific deployed contract.
func NewICollateralManagementCaller(address common.Address, caller bind.ContractCaller) (*ICollateralManagementCaller, error) {
	contract, err := bindICollateralManagement(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementCaller{contract: contract}, nil
}

// NewICollateralManagementTransactor creates a new write-only instance of ICollateralManagement, bound to a specific deployed contract.
func NewICollateralManagementTransactor(address common.Address, transactor bind.ContractTransactor) (*ICollateralManagementTransactor, error) {
	contract, err := bindICollateralManagement(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementTransactor{contract: contract}, nil
}

// NewICollateralManagementFilterer creates a new log filterer instance of ICollateralManagement, bound to a specific deployed contract.
func NewICollateralManagementFilterer(address common.Address, filterer bind.ContractFilterer) (*ICollateralManagementFilterer, error) {
	contract, err := bindICollateralManagement(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementFilterer{contract: contract}, nil
}

// bindICollateralManagement binds a generic wrapper to an already deployed contract.
func bindICollateralManagement(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICollateralManagementMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICollateralManagement *ICollateralManagementRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICollateralManagement.Contract.ICollateralManagementCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICollateralManagement *ICollateralManagementRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.ICollateralManagementTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICollateralManagement *ICollateralManagementRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.ICollateralManagementTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICollateralManagement *ICollateralManagementCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICollateralManagement.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICollateralManagement *ICollateralManagementTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICollateralManagement *ICollateralManagementTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.contract.Transact(opts, method, params...)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetMinCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getMinCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetMinCollateral() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetMinCollateral(&_ICollateralManagement.CallOpts)
}

// GetMinCollateral is a free data retrieval call binding the contract method 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetMinCollateral() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetMinCollateral(&_ICollateralManagement.CallOpts)
}

// GetPegInCollateral is a free data retrieval call binding the contract method 0x003c3317.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetPegInCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getPegInCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPegInCollateral is a free data retrieval call binding the contract method 0x003c3317.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetPegInCollateral(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPegInCollateral(&_ICollateralManagement.CallOpts, addr)
}

// GetPegInCollateral is a free data retrieval call binding the contract method 0x003c3317.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetPegInCollateral(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPegInCollateral(&_ICollateralManagement.CallOpts, addr)
}

// GetPegOutCollateral is a free data retrieval call binding the contract method 0x82b90e93.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetPegOutCollateral(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getPegOutCollateral", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPegOutCollateral is a free data retrieval call binding the contract method 0x82b90e93.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetPegOutCollateral(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPegOutCollateral(&_ICollateralManagement.CallOpts, addr)
}

// GetPegOutCollateral is a free data retrieval call binding the contract method 0x82b90e93.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetPegOutCollateral(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPegOutCollateral(&_ICollateralManagement.CallOpts, addr)
}

// GetPenalties is a free data retrieval call binding the contract method 0xe6ef2a38.
//
// Solidity: function getPenalties() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetPenalties(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getPenalties")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPenalties is a free data retrieval call binding the contract method 0xe6ef2a38.
//
// Solidity: function getPenalties() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetPenalties() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPenalties(&_ICollateralManagement.CallOpts)
}

// GetPenalties is a free data retrieval call binding the contract method 0xe6ef2a38.
//
// Solidity: function getPenalties() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetPenalties() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetPenalties(&_ICollateralManagement.CallOpts)
}

// GetResignDelayInBlocks is a free data retrieval call binding the contract method 0x27887ffc.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetResignDelayInBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getResignDelayInBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetResignDelayInBlocks is a free data retrieval call binding the contract method 0x27887ffc.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetResignDelayInBlocks() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetResignDelayInBlocks(&_ICollateralManagement.CallOpts)
}

// GetResignDelayInBlocks is a free data retrieval call binding the contract method 0x27887ffc.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetResignDelayInBlocks() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetResignDelayInBlocks(&_ICollateralManagement.CallOpts)
}

// GetResignationBlock is a free data retrieval call binding the contract method 0xd36933d3.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetResignationBlock(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getResignationBlock", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetResignationBlock is a free data retrieval call binding the contract method 0xd36933d3.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetResignationBlock(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetResignationBlock(&_ICollateralManagement.CallOpts, addr)
}

// GetResignationBlock is a free data retrieval call binding the contract method 0xd36933d3.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetResignationBlock(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetResignationBlock(&_ICollateralManagement.CallOpts, addr)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetRewardPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getRewardPercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetRewardPercentage() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetRewardPercentage(&_ICollateralManagement.CallOpts)
}

// GetRewardPercentage is a free data retrieval call binding the contract method 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetRewardPercentage() (*big.Int, error) {
	return _ICollateralManagement.Contract.GetRewardPercentage(&_ICollateralManagement.CallOpts)
}

// GetRewards is a free data retrieval call binding the contract method 0x79ee54f7.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCaller) GetRewards(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "getRewards", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRewards is a free data retrieval call binding the contract method 0x79ee54f7.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementSession) GetRewards(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetRewards(&_ICollateralManagement.CallOpts, addr)
}

// GetRewards is a free data retrieval call binding the contract method 0x79ee54f7.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (_ICollateralManagement *ICollateralManagementCallerSession) GetRewards(addr common.Address) (*big.Int, error) {
	return _ICollateralManagement.Contract.GetRewards(&_ICollateralManagement.CallOpts, addr)
}

// IsCollateralSufficient is a free data retrieval call binding the contract method 0x718c5aa8.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementCaller) IsCollateralSufficient(opts *bind.CallOpts, providerType uint8, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "isCollateralSufficient", providerType, addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsCollateralSufficient is a free data retrieval call binding the contract method 0x718c5aa8.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementSession) IsCollateralSufficient(providerType uint8, addr common.Address) (bool, error) {
	return _ICollateralManagement.Contract.IsCollateralSufficient(&_ICollateralManagement.CallOpts, providerType, addr)
}

// IsCollateralSufficient is a free data retrieval call binding the contract method 0x718c5aa8.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementCallerSession) IsCollateralSufficient(providerType uint8, addr common.Address) (bool, error) {
	return _ICollateralManagement.Contract.IsCollateralSufficient(&_ICollateralManagement.CallOpts, providerType, addr)
}

// IsRegistered is a free data retrieval call binding the contract method 0x900daa73.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementCaller) IsRegistered(opts *bind.CallOpts, providerType uint8, addr common.Address) (bool, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "isRegistered", providerType, addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0x900daa73.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementSession) IsRegistered(providerType uint8, addr common.Address) (bool, error) {
	return _ICollateralManagement.Contract.IsRegistered(&_ICollateralManagement.CallOpts, providerType, addr)
}

// IsRegistered is a free data retrieval call binding the contract method 0x900daa73.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (_ICollateralManagement *ICollateralManagementCallerSession) IsRegistered(providerType uint8, addr common.Address) (bool, error) {
	return _ICollateralManagement.Contract.IsRegistered(&_ICollateralManagement.CallOpts, providerType, addr)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_ICollateralManagement *ICollateralManagementCaller) PauseStatus(opts *bind.CallOpts) (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	var out []interface{}
	err := _ICollateralManagement.contract.Call(opts, &out, "pauseStatus")

	outstruct := new(struct {
		IsPaused bool
		Reason   string
		Since    uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_ICollateralManagement *ICollateralManagementSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _ICollateralManagement.Contract.PauseStatus(&_ICollateralManagement.CallOpts)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_ICollateralManagement *ICollateralManagementCallerSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _ICollateralManagement.Contract.PauseStatus(&_ICollateralManagement.CallOpts)
}

// AddPegInCollateral is a paid mutator transaction binding the contract method 0xde567d6d.
//
// Solidity: function addPegInCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementTransactor) AddPegInCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "addPegInCollateral")
}

// AddPegInCollateral is a paid mutator transaction binding the contract method 0xde567d6d.
//
// Solidity: function addPegInCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementSession) AddPegInCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegInCollateral(&_ICollateralManagement.TransactOpts)
}

// AddPegInCollateral is a paid mutator transaction binding the contract method 0xde567d6d.
//
// Solidity: function addPegInCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) AddPegInCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegInCollateral(&_ICollateralManagement.TransactOpts)
}

// AddPegInCollateralTo is a paid mutator transaction binding the contract method 0x83fe87f9.
//
// Solidity: function addPegInCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementTransactor) AddPegInCollateralTo(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "addPegInCollateralTo", addr)
}

// AddPegInCollateralTo is a paid mutator transaction binding the contract method 0x83fe87f9.
//
// Solidity: function addPegInCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementSession) AddPegInCollateralTo(addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegInCollateralTo(&_ICollateralManagement.TransactOpts, addr)
}

// AddPegInCollateralTo is a paid mutator transaction binding the contract method 0x83fe87f9.
//
// Solidity: function addPegInCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) AddPegInCollateralTo(addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegInCollateralTo(&_ICollateralManagement.TransactOpts, addr)
}

// AddPegOutCollateral is a paid mutator transaction binding the contract method 0x52b2318d.
//
// Solidity: function addPegOutCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementTransactor) AddPegOutCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "addPegOutCollateral")
}

// AddPegOutCollateral is a paid mutator transaction binding the contract method 0x52b2318d.
//
// Solidity: function addPegOutCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementSession) AddPegOutCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegOutCollateral(&_ICollateralManagement.TransactOpts)
}

// AddPegOutCollateral is a paid mutator transaction binding the contract method 0x52b2318d.
//
// Solidity: function addPegOutCollateral() payable returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) AddPegOutCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegOutCollateral(&_ICollateralManagement.TransactOpts)
}

// AddPegOutCollateralTo is a paid mutator transaction binding the contract method 0x313ee5b1.
//
// Solidity: function addPegOutCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementTransactor) AddPegOutCollateralTo(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "addPegOutCollateralTo", addr)
}

// AddPegOutCollateralTo is a paid mutator transaction binding the contract method 0x313ee5b1.
//
// Solidity: function addPegOutCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementSession) AddPegOutCollateralTo(addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegOutCollateralTo(&_ICollateralManagement.TransactOpts, addr)
}

// AddPegOutCollateralTo is a paid mutator transaction binding the contract method 0x313ee5b1.
//
// Solidity: function addPegOutCollateralTo(address addr) payable returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) AddPegOutCollateralTo(addr common.Address) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.AddPegOutCollateralTo(&_ICollateralManagement.TransactOpts, addr)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_ICollateralManagement *ICollateralManagementTransactor) Resign(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "resign")
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_ICollateralManagement *ICollateralManagementSession) Resign() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.Resign(&_ICollateralManagement.TransactOpts)
}

// Resign is a paid mutator transaction binding the contract method 0x69652fcf.
//
// Solidity: function resign() returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) Resign() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.Resign(&_ICollateralManagement.TransactOpts)
}

// SlashPegInCollateral is a paid mutator transaction binding the contract method 0x3e4de194.
//
// Solidity: function slashPegInCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementTransactor) SlashPegInCollateral(opts *bind.TransactOpts, punisher common.Address, quote QuotesPegInQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "slashPegInCollateral", punisher, quote, quoteHash)
}

// SlashPegInCollateral is a paid mutator transaction binding the contract method 0x3e4de194.
//
// Solidity: function slashPegInCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementSession) SlashPegInCollateral(punisher common.Address, quote QuotesPegInQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.SlashPegInCollateral(&_ICollateralManagement.TransactOpts, punisher, quote, quoteHash)
}

// SlashPegInCollateral is a paid mutator transaction binding the contract method 0x3e4de194.
//
// Solidity: function slashPegInCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) SlashPegInCollateral(punisher common.Address, quote QuotesPegInQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.SlashPegInCollateral(&_ICollateralManagement.TransactOpts, punisher, quote, quoteHash)
}

// SlashPegOutCollateral is a paid mutator transaction binding the contract method 0x2f6fcee8.
//
// Solidity: function slashPegOutCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementTransactor) SlashPegOutCollateral(opts *bind.TransactOpts, punisher common.Address, quote QuotesPegOutQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "slashPegOutCollateral", punisher, quote, quoteHash)
}

// SlashPegOutCollateral is a paid mutator transaction binding the contract method 0x2f6fcee8.
//
// Solidity: function slashPegOutCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementSession) SlashPegOutCollateral(punisher common.Address, quote QuotesPegOutQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.SlashPegOutCollateral(&_ICollateralManagement.TransactOpts, punisher, quote, quoteHash)
}

// SlashPegOutCollateral is a paid mutator transaction binding the contract method 0x2f6fcee8.
//
// Solidity: function slashPegOutCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) SlashPegOutCollateral(punisher common.Address, quote QuotesPegOutQuote, quoteHash [32]byte) (*types.Transaction, error) {
	return _ICollateralManagement.Contract.SlashPegOutCollateral(&_ICollateralManagement.TransactOpts, punisher, quote, quoteHash)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_ICollateralManagement *ICollateralManagementTransactor) WithdrawCollateral(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "withdrawCollateral")
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_ICollateralManagement *ICollateralManagementSession) WithdrawCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.WithdrawCollateral(&_ICollateralManagement.TransactOpts)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59c153be.
//
// Solidity: function withdrawCollateral() returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) WithdrawCollateral() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.WithdrawCollateral(&_ICollateralManagement.TransactOpts)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ICollateralManagement *ICollateralManagementTransactor) WithdrawRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICollateralManagement.contract.Transact(opts, "withdrawRewards")
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ICollateralManagement *ICollateralManagementSession) WithdrawRewards() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.WithdrawRewards(&_ICollateralManagement.TransactOpts)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0xc7b8981c.
//
// Solidity: function withdrawRewards() returns()
func (_ICollateralManagement *ICollateralManagementTransactorSession) WithdrawRewards() (*types.Transaction, error) {
	return _ICollateralManagement.Contract.WithdrawRewards(&_ICollateralManagement.TransactOpts)
}

// ICollateralManagementPegInCollateralAddedIterator is returned from FilterPegInCollateralAdded and is used to iterate over the raw logs and unpacked data for PegInCollateralAdded events raised by the ICollateralManagement contract.
type ICollateralManagementPegInCollateralAddedIterator struct {
	Event *ICollateralManagementPegInCollateralAdded // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementPegInCollateralAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementPegInCollateralAdded)
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
		it.Event = new(ICollateralManagementPegInCollateralAdded)
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
func (it *ICollateralManagementPegInCollateralAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementPegInCollateralAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementPegInCollateralAdded represents a PegInCollateralAdded event raised by the ICollateralManagement contract.
type ICollateralManagementPegInCollateralAdded struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegInCollateralAdded is a free log retrieval operation binding the contract event 0x6367fadf75ad09195dbdec01f12bc496bae0594c2428c82e680c37058cfba11d.
//
// Solidity: event PegInCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterPegInCollateralAdded(opts *bind.FilterOpts, addr []common.Address, amount []*big.Int) (*ICollateralManagementPegInCollateralAddedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "PegInCollateralAdded", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementPegInCollateralAddedIterator{contract: _ICollateralManagement.contract, event: "PegInCollateralAdded", logs: logs, sub: sub}, nil
}

// WatchPegInCollateralAdded is a free log subscription operation binding the contract event 0x6367fadf75ad09195dbdec01f12bc496bae0594c2428c82e680c37058cfba11d.
//
// Solidity: event PegInCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchPegInCollateralAdded(opts *bind.WatchOpts, sink chan<- *ICollateralManagementPegInCollateralAdded, addr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "PegInCollateralAdded", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementPegInCollateralAdded)
				if err := _ICollateralManagement.contract.UnpackLog(event, "PegInCollateralAdded", log); err != nil {
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

// ParsePegInCollateralAdded is a log parse operation binding the contract event 0x6367fadf75ad09195dbdec01f12bc496bae0594c2428c82e680c37058cfba11d.
//
// Solidity: event PegInCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) ParsePegInCollateralAdded(log types.Log) (*ICollateralManagementPegInCollateralAdded, error) {
	event := new(ICollateralManagementPegInCollateralAdded)
	if err := _ICollateralManagement.contract.UnpackLog(event, "PegInCollateralAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICollateralManagementPegOutCollateralAddedIterator is returned from FilterPegOutCollateralAdded and is used to iterate over the raw logs and unpacked data for PegOutCollateralAdded events raised by the ICollateralManagement contract.
type ICollateralManagementPegOutCollateralAddedIterator struct {
	Event *ICollateralManagementPegOutCollateralAdded // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementPegOutCollateralAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementPegOutCollateralAdded)
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
		it.Event = new(ICollateralManagementPegOutCollateralAdded)
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
func (it *ICollateralManagementPegOutCollateralAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementPegOutCollateralAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementPegOutCollateralAdded represents a PegOutCollateralAdded event raised by the ICollateralManagement contract.
type ICollateralManagementPegOutCollateralAdded struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPegOutCollateralAdded is a free log retrieval operation binding the contract event 0xdc9c5f97f3f996ff1e10bfd5c3986e9f65d4ad9b630dc9db357d014124f60512.
//
// Solidity: event PegOutCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterPegOutCollateralAdded(opts *bind.FilterOpts, addr []common.Address, amount []*big.Int) (*ICollateralManagementPegOutCollateralAddedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "PegOutCollateralAdded", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementPegOutCollateralAddedIterator{contract: _ICollateralManagement.contract, event: "PegOutCollateralAdded", logs: logs, sub: sub}, nil
}

// WatchPegOutCollateralAdded is a free log subscription operation binding the contract event 0xdc9c5f97f3f996ff1e10bfd5c3986e9f65d4ad9b630dc9db357d014124f60512.
//
// Solidity: event PegOutCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchPegOutCollateralAdded(opts *bind.WatchOpts, sink chan<- *ICollateralManagementPegOutCollateralAdded, addr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "PegOutCollateralAdded", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementPegOutCollateralAdded)
				if err := _ICollateralManagement.contract.UnpackLog(event, "PegOutCollateralAdded", log); err != nil {
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

// ParsePegOutCollateralAdded is a log parse operation binding the contract event 0xdc9c5f97f3f996ff1e10bfd5c3986e9f65d4ad9b630dc9db357d014124f60512.
//
// Solidity: event PegOutCollateralAdded(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) ParsePegOutCollateralAdded(log types.Log) (*ICollateralManagementPegOutCollateralAdded, error) {
	event := new(ICollateralManagementPegOutCollateralAdded)
	if err := _ICollateralManagement.contract.UnpackLog(event, "PegOutCollateralAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICollateralManagementPenalizedIterator is returned from FilterPenalized and is used to iterate over the raw logs and unpacked data for Penalized events raised by the ICollateralManagement contract.
type ICollateralManagementPenalizedIterator struct {
	Event *ICollateralManagementPenalized // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementPenalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementPenalized)
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
		it.Event = new(ICollateralManagementPenalized)
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
func (it *ICollateralManagementPenalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementPenalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementPenalized represents a Penalized event raised by the ICollateralManagement contract.
type ICollateralManagementPenalized struct {
	LiquidityProvider common.Address
	Punisher          common.Address
	QuoteHash         [32]byte
	CollateralType    uint8
	Penalty           *big.Int
	Reward            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPenalized is a free log retrieval operation binding the contract event 0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f.
//
// Solidity: event Penalized(address indexed liquidityProvider, address indexed punisher, bytes32 indexed quoteHash, uint8 collateralType, uint256 penalty, uint256 reward)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterPenalized(opts *bind.FilterOpts, liquidityProvider []common.Address, punisher []common.Address, quoteHash [][32]byte) (*ICollateralManagementPenalizedIterator, error) {

	var liquidityProviderRule []interface{}
	for _, liquidityProviderItem := range liquidityProvider {
		liquidityProviderRule = append(liquidityProviderRule, liquidityProviderItem)
	}
	var punisherRule []interface{}
	for _, punisherItem := range punisher {
		punisherRule = append(punisherRule, punisherItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "Penalized", liquidityProviderRule, punisherRule, quoteHashRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementPenalizedIterator{contract: _ICollateralManagement.contract, event: "Penalized", logs: logs, sub: sub}, nil
}

// WatchPenalized is a free log subscription operation binding the contract event 0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f.
//
// Solidity: event Penalized(address indexed liquidityProvider, address indexed punisher, bytes32 indexed quoteHash, uint8 collateralType, uint256 penalty, uint256 reward)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchPenalized(opts *bind.WatchOpts, sink chan<- *ICollateralManagementPenalized, liquidityProvider []common.Address, punisher []common.Address, quoteHash [][32]byte) (event.Subscription, error) {

	var liquidityProviderRule []interface{}
	for _, liquidityProviderItem := range liquidityProvider {
		liquidityProviderRule = append(liquidityProviderRule, liquidityProviderItem)
	}
	var punisherRule []interface{}
	for _, punisherItem := range punisher {
		punisherRule = append(punisherRule, punisherItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "Penalized", liquidityProviderRule, punisherRule, quoteHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementPenalized)
				if err := _ICollateralManagement.contract.UnpackLog(event, "Penalized", log); err != nil {
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

// ParsePenalized is a log parse operation binding the contract event 0x32d8dcdc3bd4d5d6dd9053c2e1d421c681715c97c6232e33a8658b7ae0bef13f.
//
// Solidity: event Penalized(address indexed liquidityProvider, address indexed punisher, bytes32 indexed quoteHash, uint8 collateralType, uint256 penalty, uint256 reward)
func (_ICollateralManagement *ICollateralManagementFilterer) ParsePenalized(log types.Log) (*ICollateralManagementPenalized, error) {
	event := new(ICollateralManagementPenalized)
	if err := _ICollateralManagement.contract.UnpackLog(event, "Penalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICollateralManagementResignedIterator is returned from FilterResigned and is used to iterate over the raw logs and unpacked data for Resigned events raised by the ICollateralManagement contract.
type ICollateralManagementResignedIterator struct {
	Event *ICollateralManagementResigned // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementResignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementResigned)
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
		it.Event = new(ICollateralManagementResigned)
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
func (it *ICollateralManagementResignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementResignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementResigned represents a Resigned event raised by the ICollateralManagement contract.
type ICollateralManagementResigned struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterResigned is a free log retrieval operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address indexed addr)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterResigned(opts *bind.FilterOpts, addr []common.Address) (*ICollateralManagementResignedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "Resigned", addrRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementResignedIterator{contract: _ICollateralManagement.contract, event: "Resigned", logs: logs, sub: sub}, nil
}

// WatchResigned is a free log subscription operation binding the contract event 0xa6674aa33cd1b7435474751667707bf05fde99e537d67043ec5f907782577d86.
//
// Solidity: event Resigned(address indexed addr)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchResigned(opts *bind.WatchOpts, sink chan<- *ICollateralManagementResigned, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "Resigned", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementResigned)
				if err := _ICollateralManagement.contract.UnpackLog(event, "Resigned", log); err != nil {
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
// Solidity: event Resigned(address indexed addr)
func (_ICollateralManagement *ICollateralManagementFilterer) ParseResigned(log types.Log) (*ICollateralManagementResigned, error) {
	event := new(ICollateralManagementResigned)
	if err := _ICollateralManagement.contract.UnpackLog(event, "Resigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICollateralManagementRewardsWithdrawnIterator is returned from FilterRewardsWithdrawn and is used to iterate over the raw logs and unpacked data for RewardsWithdrawn events raised by the ICollateralManagement contract.
type ICollateralManagementRewardsWithdrawnIterator struct {
	Event *ICollateralManagementRewardsWithdrawn // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementRewardsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementRewardsWithdrawn)
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
		it.Event = new(ICollateralManagementRewardsWithdrawn)
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
func (it *ICollateralManagementRewardsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementRewardsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementRewardsWithdrawn represents a RewardsWithdrawn event raised by the ICollateralManagement contract.
type ICollateralManagementRewardsWithdrawn struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardsWithdrawn is a free log retrieval operation binding the contract event 0x8a43c4352486ec339f487f64af78ca5cbf06cd47833f073d3baf3a193e503161.
//
// Solidity: event RewardsWithdrawn(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterRewardsWithdrawn(opts *bind.FilterOpts, addr []common.Address, amount []*big.Int) (*ICollateralManagementRewardsWithdrawnIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "RewardsWithdrawn", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementRewardsWithdrawnIterator{contract: _ICollateralManagement.contract, event: "RewardsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchRewardsWithdrawn is a free log subscription operation binding the contract event 0x8a43c4352486ec339f487f64af78ca5cbf06cd47833f073d3baf3a193e503161.
//
// Solidity: event RewardsWithdrawn(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchRewardsWithdrawn(opts *bind.WatchOpts, sink chan<- *ICollateralManagementRewardsWithdrawn, addr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "RewardsWithdrawn", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementRewardsWithdrawn)
				if err := _ICollateralManagement.contract.UnpackLog(event, "RewardsWithdrawn", log); err != nil {
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

// ParseRewardsWithdrawn is a log parse operation binding the contract event 0x8a43c4352486ec339f487f64af78ca5cbf06cd47833f073d3baf3a193e503161.
//
// Solidity: event RewardsWithdrawn(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) ParseRewardsWithdrawn(log types.Log) (*ICollateralManagementRewardsWithdrawn, error) {
	event := new(ICollateralManagementRewardsWithdrawn)
	if err := _ICollateralManagement.contract.UnpackLog(event, "RewardsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICollateralManagementWithdrawCollateralIterator is returned from FilterWithdrawCollateral and is used to iterate over the raw logs and unpacked data for WithdrawCollateral events raised by the ICollateralManagement contract.
type ICollateralManagementWithdrawCollateralIterator struct {
	Event *ICollateralManagementWithdrawCollateral // Event containing the contract specifics and raw log

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
func (it *ICollateralManagementWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICollateralManagementWithdrawCollateral)
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
		it.Event = new(ICollateralManagementWithdrawCollateral)
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
func (it *ICollateralManagementWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICollateralManagementWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICollateralManagementWithdrawCollateral represents a WithdrawCollateral event raised by the ICollateralManagement contract.
type ICollateralManagementWithdrawCollateral struct {
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawCollateral is a free log retrieval operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) FilterWithdrawCollateral(opts *bind.FilterOpts, addr []common.Address, amount []*big.Int) (*ICollateralManagementWithdrawCollateralIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.FilterLogs(opts, "WithdrawCollateral", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &ICollateralManagementWithdrawCollateralIterator{contract: _ICollateralManagement.contract, event: "WithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchWithdrawCollateral is a free log subscription operation binding the contract event 0xa8e76b822fc682be77f3b1c822ea81f6bda5aed92ba82e6873bfd889f328d1d2.
//
// Solidity: event WithdrawCollateral(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) WatchWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *ICollateralManagementWithdrawCollateral, addr []common.Address, amount []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _ICollateralManagement.contract.WatchLogs(opts, "WithdrawCollateral", addrRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICollateralManagementWithdrawCollateral)
				if err := _ICollateralManagement.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
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
// Solidity: event WithdrawCollateral(address indexed addr, uint256 indexed amount)
func (_ICollateralManagement *ICollateralManagementFilterer) ParseWithdrawCollateral(log types.Log) (*ICollateralManagementWithdrawCollateral, error) {
	event := new(ICollateralManagementWithdrawCollateral)
	if err := _ICollateralManagement.contract.UnpackLog(event, "WithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFlyoverDiscoveryMetaData contains all meta data concerning the IFlyoverDiscovery contract.
var IFlyoverDiscoveryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"}],\"name\":\"getProvider\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"internalType\":\"structFlyover.LiquidityProvider\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProviders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"internalType\":\"structFlyover.LiquidityProvider[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProvidersId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isPaused\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"since\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providerId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setProviderStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"updateProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"ProviderStatusSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"ProviderUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"AlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InsufficientCollateral\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"InvalidProviderData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"}],\"name\":\"InvalidProviderType\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotEOA\",\"type\":\"error\"}]",
}

// IFlyoverDiscoveryABI is the input ABI used to generate the binding from.
// Deprecated: Use IFlyoverDiscoveryMetaData.ABI instead.
var IFlyoverDiscoveryABI = IFlyoverDiscoveryMetaData.ABI

// IFlyoverDiscovery is an auto generated Go binding around an Ethereum contract.
type IFlyoverDiscovery struct {
	IFlyoverDiscoveryCaller     // Read-only binding to the contract
	IFlyoverDiscoveryTransactor // Write-only binding to the contract
	IFlyoverDiscoveryFilterer   // Log filterer for contract events
}

// IFlyoverDiscoveryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IFlyoverDiscoveryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFlyoverDiscoveryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IFlyoverDiscoveryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFlyoverDiscoveryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IFlyoverDiscoveryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IFlyoverDiscoverySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IFlyoverDiscoverySession struct {
	Contract     *IFlyoverDiscovery // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IFlyoverDiscoveryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IFlyoverDiscoveryCallerSession struct {
	Contract *IFlyoverDiscoveryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IFlyoverDiscoveryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IFlyoverDiscoveryTransactorSession struct {
	Contract     *IFlyoverDiscoveryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IFlyoverDiscoveryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IFlyoverDiscoveryRaw struct {
	Contract *IFlyoverDiscovery // Generic contract binding to access the raw methods on
}

// IFlyoverDiscoveryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IFlyoverDiscoveryCallerRaw struct {
	Contract *IFlyoverDiscoveryCaller // Generic read-only contract binding to access the raw methods on
}

// IFlyoverDiscoveryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IFlyoverDiscoveryTransactorRaw struct {
	Contract *IFlyoverDiscoveryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIFlyoverDiscovery creates a new instance of IFlyoverDiscovery, bound to a specific deployed contract.
func NewIFlyoverDiscovery(address common.Address, backend bind.ContractBackend) (*IFlyoverDiscovery, error) {
	contract, err := bindIFlyoverDiscovery(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscovery{IFlyoverDiscoveryCaller: IFlyoverDiscoveryCaller{contract: contract}, IFlyoverDiscoveryTransactor: IFlyoverDiscoveryTransactor{contract: contract}, IFlyoverDiscoveryFilterer: IFlyoverDiscoveryFilterer{contract: contract}}, nil
}

// NewIFlyoverDiscoveryCaller creates a new read-only instance of IFlyoverDiscovery, bound to a specific deployed contract.
func NewIFlyoverDiscoveryCaller(address common.Address, caller bind.ContractCaller) (*IFlyoverDiscoveryCaller, error) {
	contract, err := bindIFlyoverDiscovery(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryCaller{contract: contract}, nil
}

// NewIFlyoverDiscoveryTransactor creates a new write-only instance of IFlyoverDiscovery, bound to a specific deployed contract.
func NewIFlyoverDiscoveryTransactor(address common.Address, transactor bind.ContractTransactor) (*IFlyoverDiscoveryTransactor, error) {
	contract, err := bindIFlyoverDiscovery(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryTransactor{contract: contract}, nil
}

// NewIFlyoverDiscoveryFilterer creates a new log filterer instance of IFlyoverDiscovery, bound to a specific deployed contract.
func NewIFlyoverDiscoveryFilterer(address common.Address, filterer bind.ContractFilterer) (*IFlyoverDiscoveryFilterer, error) {
	contract, err := bindIFlyoverDiscovery(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryFilterer{contract: contract}, nil
}

// bindIFlyoverDiscovery binds a generic wrapper to an already deployed contract.
func bindIFlyoverDiscovery(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IFlyoverDiscoveryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFlyoverDiscovery *IFlyoverDiscoveryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IFlyoverDiscovery.Contract.IFlyoverDiscoveryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFlyoverDiscovery *IFlyoverDiscoveryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.IFlyoverDiscoveryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFlyoverDiscovery *IFlyoverDiscoveryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.IFlyoverDiscoveryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IFlyoverDiscovery.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.contract.Transact(opts, method, params...)
}

// GetProvider is a free data retrieval call binding the contract method 0x55f21eb7.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (_IFlyoverDiscovery *IFlyoverDiscoveryCaller) GetProvider(opts *bind.CallOpts, providerAddress common.Address) (FlyoverLiquidityProvider, error) {
	var out []interface{}
	err := _IFlyoverDiscovery.contract.Call(opts, &out, "getProvider", providerAddress)

	if err != nil {
		return *new(FlyoverLiquidityProvider), err
	}

	out0 := *abi.ConvertType(out[0], new(FlyoverLiquidityProvider)).(*FlyoverLiquidityProvider)

	return out0, err

}

// GetProvider is a free data retrieval call binding the contract method 0x55f21eb7.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) GetProvider(providerAddress common.Address) (FlyoverLiquidityProvider, error) {
	return _IFlyoverDiscovery.Contract.GetProvider(&_IFlyoverDiscovery.CallOpts, providerAddress)
}

// GetProvider is a free data retrieval call binding the contract method 0x55f21eb7.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerSession) GetProvider(providerAddress common.Address) (FlyoverLiquidityProvider, error) {
	return _IFlyoverDiscovery.Contract.GetProvider(&_IFlyoverDiscovery.CallOpts, providerAddress)
}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (_IFlyoverDiscovery *IFlyoverDiscoveryCaller) GetProviders(opts *bind.CallOpts) ([]FlyoverLiquidityProvider, error) {
	var out []interface{}
	err := _IFlyoverDiscovery.contract.Call(opts, &out, "getProviders")

	if err != nil {
		return *new([]FlyoverLiquidityProvider), err
	}

	out0 := *abi.ConvertType(out[0], new([]FlyoverLiquidityProvider)).(*[]FlyoverLiquidityProvider)

	return out0, err

}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) GetProviders() ([]FlyoverLiquidityProvider, error) {
	return _IFlyoverDiscovery.Contract.GetProviders(&_IFlyoverDiscovery.CallOpts)
}

// GetProviders is a free data retrieval call binding the contract method 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerSession) GetProviders() ([]FlyoverLiquidityProvider, error) {
	return _IFlyoverDiscovery.Contract.GetProviders(&_IFlyoverDiscovery.CallOpts)
}

// GetProvidersId is a free data retrieval call binding the contract method 0x122dab09.
//
// Solidity: function getProvidersId() view returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCaller) GetProvidersId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IFlyoverDiscovery.contract.Call(opts, &out, "getProvidersId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProvidersId is a free data retrieval call binding the contract method 0x122dab09.
//
// Solidity: function getProvidersId() view returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) GetProvidersId() (*big.Int, error) {
	return _IFlyoverDiscovery.Contract.GetProvidersId(&_IFlyoverDiscovery.CallOpts)
}

// GetProvidersId is a free data retrieval call binding the contract method 0x122dab09.
//
// Solidity: function getProvidersId() view returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerSession) GetProvidersId() (*big.Int, error) {
	return _IFlyoverDiscovery.Contract.GetProvidersId(&_IFlyoverDiscovery.CallOpts)
}

// IsOperational is a free data retrieval call binding the contract method 0xbf50daf0.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCaller) IsOperational(opts *bind.CallOpts, providerType uint8, addr common.Address) (bool, error) {
	var out []interface{}
	err := _IFlyoverDiscovery.contract.Call(opts, &out, "isOperational", providerType, addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperational is a free data retrieval call binding the contract method 0xbf50daf0.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) IsOperational(providerType uint8, addr common.Address) (bool, error) {
	return _IFlyoverDiscovery.Contract.IsOperational(&_IFlyoverDiscovery.CallOpts, providerType, addr)
}

// IsOperational is a free data retrieval call binding the contract method 0xbf50daf0.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerSession) IsOperational(providerType uint8, addr common.Address) (bool, error) {
	return _IFlyoverDiscovery.Contract.IsOperational(&_IFlyoverDiscovery.CallOpts, providerType, addr)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCaller) PauseStatus(opts *bind.CallOpts) (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	var out []interface{}
	err := _IFlyoverDiscovery.contract.Call(opts, &out, "pauseStatus")

	outstruct := new(struct {
		IsPaused bool
		Reason   string
		Since    uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IFlyoverDiscovery.Contract.PauseStatus(&_IFlyoverDiscovery.CallOpts)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IFlyoverDiscovery *IFlyoverDiscoveryCallerSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IFlyoverDiscovery.Contract.PauseStatus(&_IFlyoverDiscovery.CallOpts)
}

// Register is a paid mutator transaction binding the contract method 0x4487ce11.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactor) Register(opts *bind.TransactOpts, name string, apiBaseUrl string, status bool, providerType uint8) (*types.Transaction, error) {
	return _IFlyoverDiscovery.contract.Transact(opts, "register", name, apiBaseUrl, status, providerType)
}

// Register is a paid mutator transaction binding the contract method 0x4487ce11.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) Register(name string, apiBaseUrl string, status bool, providerType uint8) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.Register(&_IFlyoverDiscovery.TransactOpts, name, apiBaseUrl, status, providerType)
}

// Register is a paid mutator transaction binding the contract method 0x4487ce11.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactorSession) Register(name string, apiBaseUrl string, status bool, providerType uint8) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.Register(&_IFlyoverDiscovery.TransactOpts, name, apiBaseUrl, status, providerType)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 providerId, bool status) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactor) SetProviderStatus(opts *bind.TransactOpts, providerId *big.Int, status bool) (*types.Transaction, error) {
	return _IFlyoverDiscovery.contract.Transact(opts, "setProviderStatus", providerId, status)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 providerId, bool status) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) SetProviderStatus(providerId *big.Int, status bool) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.SetProviderStatus(&_IFlyoverDiscovery.TransactOpts, providerId, status)
}

// SetProviderStatus is a paid mutator transaction binding the contract method 0x72cbf4e8.
//
// Solidity: function setProviderStatus(uint256 providerId, bool status) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactorSession) SetProviderStatus(providerId *big.Int, status bool) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.SetProviderStatus(&_IFlyoverDiscovery.TransactOpts, providerId, status)
}

// UpdateProvider is a paid mutator transaction binding the contract method 0x0220f41d.
//
// Solidity: function updateProvider(string name, string apiBaseUrl) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactor) UpdateProvider(opts *bind.TransactOpts, name string, apiBaseUrl string) (*types.Transaction, error) {
	return _IFlyoverDiscovery.contract.Transact(opts, "updateProvider", name, apiBaseUrl)
}

// UpdateProvider is a paid mutator transaction binding the contract method 0x0220f41d.
//
// Solidity: function updateProvider(string name, string apiBaseUrl) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoverySession) UpdateProvider(name string, apiBaseUrl string) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.UpdateProvider(&_IFlyoverDiscovery.TransactOpts, name, apiBaseUrl)
}

// UpdateProvider is a paid mutator transaction binding the contract method 0x0220f41d.
//
// Solidity: function updateProvider(string name, string apiBaseUrl) returns()
func (_IFlyoverDiscovery *IFlyoverDiscoveryTransactorSession) UpdateProvider(name string, apiBaseUrl string) (*types.Transaction, error) {
	return _IFlyoverDiscovery.Contract.UpdateProvider(&_IFlyoverDiscovery.TransactOpts, name, apiBaseUrl)
}

// IFlyoverDiscoveryProviderStatusSetIterator is returned from FilterProviderStatusSet and is used to iterate over the raw logs and unpacked data for ProviderStatusSet events raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryProviderStatusSetIterator struct {
	Event *IFlyoverDiscoveryProviderStatusSet // Event containing the contract specifics and raw log

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
func (it *IFlyoverDiscoveryProviderStatusSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFlyoverDiscoveryProviderStatusSet)
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
		it.Event = new(IFlyoverDiscoveryProviderStatusSet)
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
func (it *IFlyoverDiscoveryProviderStatusSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFlyoverDiscoveryProviderStatusSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFlyoverDiscoveryProviderStatusSet represents a ProviderStatusSet event raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryProviderStatusSet struct {
	Id     *big.Int
	Status bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProviderStatusSet is a free log retrieval operation binding the contract event 0x833990204fe208883ab0b3d6185f6c17f549a01a19c388ec488206ac1dbbc65d.
//
// Solidity: event ProviderStatusSet(uint256 indexed id, bool indexed status)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) FilterProviderStatusSet(opts *bind.FilterOpts, id []*big.Int, status []bool) (*IFlyoverDiscoveryProviderStatusSetIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.FilterLogs(opts, "ProviderStatusSet", idRule, statusRule)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryProviderStatusSetIterator{contract: _IFlyoverDiscovery.contract, event: "ProviderStatusSet", logs: logs, sub: sub}, nil
}

// WatchProviderStatusSet is a free log subscription operation binding the contract event 0x833990204fe208883ab0b3d6185f6c17f549a01a19c388ec488206ac1dbbc65d.
//
// Solidity: event ProviderStatusSet(uint256 indexed id, bool indexed status)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) WatchProviderStatusSet(opts *bind.WatchOpts, sink chan<- *IFlyoverDiscoveryProviderStatusSet, id []*big.Int, status []bool) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.WatchLogs(opts, "ProviderStatusSet", idRule, statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFlyoverDiscoveryProviderStatusSet)
				if err := _IFlyoverDiscovery.contract.UnpackLog(event, "ProviderStatusSet", log); err != nil {
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

// ParseProviderStatusSet is a log parse operation binding the contract event 0x833990204fe208883ab0b3d6185f6c17f549a01a19c388ec488206ac1dbbc65d.
//
// Solidity: event ProviderStatusSet(uint256 indexed id, bool indexed status)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) ParseProviderStatusSet(log types.Log) (*IFlyoverDiscoveryProviderStatusSet, error) {
	event := new(IFlyoverDiscoveryProviderStatusSet)
	if err := _IFlyoverDiscovery.contract.UnpackLog(event, "ProviderStatusSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFlyoverDiscoveryProviderUpdateIterator is returned from FilterProviderUpdate and is used to iterate over the raw logs and unpacked data for ProviderUpdate events raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryProviderUpdateIterator struct {
	Event *IFlyoverDiscoveryProviderUpdate // Event containing the contract specifics and raw log

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
func (it *IFlyoverDiscoveryProviderUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFlyoverDiscoveryProviderUpdate)
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
		it.Event = new(IFlyoverDiscoveryProviderUpdate)
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
func (it *IFlyoverDiscoveryProviderUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFlyoverDiscoveryProviderUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFlyoverDiscoveryProviderUpdate represents a ProviderUpdate event raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryProviderUpdate struct {
	From       common.Address
	Name       string
	ApiBaseUrl string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProviderUpdate is a free log retrieval operation binding the contract event 0xc15f90eb34a098bb02f2641dff62935246fb005d8f06e13d5cc6be0bddcce8e3.
//
// Solidity: event ProviderUpdate(address indexed from, string name, string apiBaseUrl)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) FilterProviderUpdate(opts *bind.FilterOpts, from []common.Address) (*IFlyoverDiscoveryProviderUpdateIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.FilterLogs(opts, "ProviderUpdate", fromRule)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryProviderUpdateIterator{contract: _IFlyoverDiscovery.contract, event: "ProviderUpdate", logs: logs, sub: sub}, nil
}

// WatchProviderUpdate is a free log subscription operation binding the contract event 0xc15f90eb34a098bb02f2641dff62935246fb005d8f06e13d5cc6be0bddcce8e3.
//
// Solidity: event ProviderUpdate(address indexed from, string name, string apiBaseUrl)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) WatchProviderUpdate(opts *bind.WatchOpts, sink chan<- *IFlyoverDiscoveryProviderUpdate, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.WatchLogs(opts, "ProviderUpdate", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFlyoverDiscoveryProviderUpdate)
				if err := _IFlyoverDiscovery.contract.UnpackLog(event, "ProviderUpdate", log); err != nil {
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

// ParseProviderUpdate is a log parse operation binding the contract event 0xc15f90eb34a098bb02f2641dff62935246fb005d8f06e13d5cc6be0bddcce8e3.
//
// Solidity: event ProviderUpdate(address indexed from, string name, string apiBaseUrl)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) ParseProviderUpdate(log types.Log) (*IFlyoverDiscoveryProviderUpdate, error) {
	event := new(IFlyoverDiscoveryProviderUpdate)
	if err := _IFlyoverDiscovery.contract.UnpackLog(event, "ProviderUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IFlyoverDiscoveryRegisterIterator is returned from FilterRegister and is used to iterate over the raw logs and unpacked data for Register events raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryRegisterIterator struct {
	Event *IFlyoverDiscoveryRegister // Event containing the contract specifics and raw log

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
func (it *IFlyoverDiscoveryRegisterIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IFlyoverDiscoveryRegister)
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
		it.Event = new(IFlyoverDiscoveryRegister)
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
func (it *IFlyoverDiscoveryRegisterIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IFlyoverDiscoveryRegisterIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IFlyoverDiscoveryRegister represents a Register event raised by the IFlyoverDiscovery contract.
type IFlyoverDiscoveryRegister struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegister is a free log retrieval operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 indexed amount)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) FilterRegister(opts *bind.FilterOpts, id []*big.Int, from []common.Address, amount []*big.Int) (*IFlyoverDiscoveryRegisterIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.FilterLogs(opts, "Register", idRule, fromRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IFlyoverDiscoveryRegisterIterator{contract: _IFlyoverDiscovery.contract, event: "Register", logs: logs, sub: sub}, nil
}

// WatchRegister is a free log subscription operation binding the contract event 0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 indexed amount)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) WatchRegister(opts *bind.WatchOpts, sink chan<- *IFlyoverDiscoveryRegister, id []*big.Int, from []common.Address, amount []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IFlyoverDiscovery.contract.WatchLogs(opts, "Register", idRule, fromRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IFlyoverDiscoveryRegister)
				if err := _IFlyoverDiscovery.contract.UnpackLog(event, "Register", log); err != nil {
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
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 indexed amount)
func (_IFlyoverDiscovery *IFlyoverDiscoveryFilterer) ParseRegister(log types.Log) (*IFlyoverDiscoveryRegister, error) {
	event := new(IFlyoverDiscoveryRegister)
	if err := _IFlyoverDiscovery.contract.UnpackLog(event, "Register", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInMetaData contains all meta data concerning the IPegIn contract.
var IPegInMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinPegIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"getQuoteStatus\",\"outputs\":[{\"internalType\":\"enumIPegIn.PegInStates\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegInQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegInQuoteEIP712\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isPaused\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"since\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"}],\"name\":\"validatePegInDepositAddress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"transferredAmount\",\"type\":\"uint256\"}],\"name\":\"PegInRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AmountUnderMinimum\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasRequired\",\"type\":\"uint256\"}],\"name\":\"InsufficientGas\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"refundAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRefundAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughConfirmations\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyProcessed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"UnexpectedBridgeError\",\"type\":\"error\"}]",
}

// IPegInABI is the input ABI used to generate the binding from.
// Deprecated: Use IPegInMetaData.ABI instead.
var IPegInABI = IPegInMetaData.ABI

// IPegIn is an auto generated Go binding around an Ethereum contract.
type IPegIn struct {
	IPegInCaller     // Read-only binding to the contract
	IPegInTransactor // Write-only binding to the contract
	IPegInFilterer   // Log filterer for contract events
}

// IPegInCaller is an auto generated read-only Go binding around an Ethereum contract.
type IPegInCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegInTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IPegInTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegInFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPegInFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegInSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPegInSession struct {
	Contract     *IPegIn           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPegInCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPegInCallerSession struct {
	Contract *IPegInCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IPegInTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPegInTransactorSession struct {
	Contract     *IPegInTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPegInRaw is an auto generated low-level Go binding around an Ethereum contract.
type IPegInRaw struct {
	Contract *IPegIn // Generic contract binding to access the raw methods on
}

// IPegInCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPegInCallerRaw struct {
	Contract *IPegInCaller // Generic read-only contract binding to access the raw methods on
}

// IPegInTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPegInTransactorRaw struct {
	Contract *IPegInTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPegIn creates a new instance of IPegIn, bound to a specific deployed contract.
func NewIPegIn(address common.Address, backend bind.ContractBackend) (*IPegIn, error) {
	contract, err := bindIPegIn(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPegIn{IPegInCaller: IPegInCaller{contract: contract}, IPegInTransactor: IPegInTransactor{contract: contract}, IPegInFilterer: IPegInFilterer{contract: contract}}, nil
}

// NewIPegInCaller creates a new read-only instance of IPegIn, bound to a specific deployed contract.
func NewIPegInCaller(address common.Address, caller bind.ContractCaller) (*IPegInCaller, error) {
	contract, err := bindIPegIn(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPegInCaller{contract: contract}, nil
}

// NewIPegInTransactor creates a new write-only instance of IPegIn, bound to a specific deployed contract.
func NewIPegInTransactor(address common.Address, transactor bind.ContractTransactor) (*IPegInTransactor, error) {
	contract, err := bindIPegIn(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPegInTransactor{contract: contract}, nil
}

// NewIPegInFilterer creates a new log filterer instance of IPegIn, bound to a specific deployed contract.
func NewIPegInFilterer(address common.Address, filterer bind.ContractFilterer) (*IPegInFilterer, error) {
	contract, err := bindIPegIn(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPegInFilterer{contract: contract}, nil
}

// bindIPegIn binds a generic wrapper to an already deployed contract.
func bindIPegIn(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPegInMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPegIn *IPegInRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPegIn.Contract.IPegInCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPegIn *IPegInRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegIn.Contract.IPegInTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPegIn *IPegInRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPegIn.Contract.IPegInTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPegIn *IPegInCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPegIn.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPegIn *IPegInTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegIn.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPegIn *IPegInTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPegIn.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegIn *IPegInCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegIn *IPegInSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IPegIn.Contract.Eip712Domain(&_IPegIn.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegIn *IPegInCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IPegIn.Contract.Eip712Domain(&_IPegIn.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_IPegIn *IPegInCaller) GetBalance(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getBalance", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_IPegIn *IPegInSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _IPegIn.Contract.GetBalance(&_IPegIn.CallOpts, addr)
}

// GetBalance is a free data retrieval call binding the contract method 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (_IPegIn *IPegInCallerSession) GetBalance(addr common.Address) (*big.Int, error) {
	return _IPegIn.Contract.GetBalance(&_IPegIn.CallOpts, addr)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_IPegIn *IPegInCaller) GetMinPegIn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getMinPegIn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_IPegIn *IPegInSession) GetMinPegIn() (*big.Int, error) {
	return _IPegIn.Contract.GetMinPegIn(&_IPegIn.CallOpts)
}

// GetMinPegIn is a free data retrieval call binding the contract method 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (_IPegIn *IPegInCallerSession) GetMinPegIn() (*big.Int, error) {
	return _IPegIn.Contract.GetMinPegIn(&_IPegIn.CallOpts)
}

// GetQuoteStatus is a free data retrieval call binding the contract method 0xf93c8ec2.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (_IPegIn *IPegInCaller) GetQuoteStatus(opts *bind.CallOpts, quoteHash [32]byte) (uint8, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getQuoteStatus", quoteHash)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetQuoteStatus is a free data retrieval call binding the contract method 0xf93c8ec2.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (_IPegIn *IPegInSession) GetQuoteStatus(quoteHash [32]byte) (uint8, error) {
	return _IPegIn.Contract.GetQuoteStatus(&_IPegIn.CallOpts, quoteHash)
}

// GetQuoteStatus is a free data retrieval call binding the contract method 0xf93c8ec2.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (_IPegIn *IPegInCallerSession) GetQuoteStatus(quoteHash [32]byte) (uint8, error) {
	return _IPegIn.Contract.GetQuoteStatus(&_IPegIn.CallOpts, quoteHash)
}

// HashPegInQuote is a free data retrieval call binding the contract method 0xf218a7d8.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInCaller) HashPegInQuote(opts *bind.CallOpts, quote QuotesPegInQuote) ([32]byte, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "hashPegInQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegInQuote is a free data retrieval call binding the contract method 0xf218a7d8.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInSession) HashPegInQuote(quote QuotesPegInQuote) ([32]byte, error) {
	return _IPegIn.Contract.HashPegInQuote(&_IPegIn.CallOpts, quote)
}

// HashPegInQuote is a free data retrieval call binding the contract method 0xf218a7d8.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInCallerSession) HashPegInQuote(quote QuotesPegInQuote) ([32]byte, error) {
	return _IPegIn.Contract.HashPegInQuote(&_IPegIn.CallOpts, quote)
}

// HashPegInQuoteEIP712 is a free data retrieval call binding the contract method 0x928f4598.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInCaller) HashPegInQuoteEIP712(opts *bind.CallOpts, quote QuotesPegInQuote) ([32]byte, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "hashPegInQuoteEIP712", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegInQuoteEIP712 is a free data retrieval call binding the contract method 0x928f4598.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInSession) HashPegInQuoteEIP712(quote QuotesPegInQuote) ([32]byte, error) {
	return _IPegIn.Contract.HashPegInQuoteEIP712(&_IPegIn.CallOpts, quote)
}

// HashPegInQuoteEIP712 is a free data retrieval call binding the contract method 0x928f4598.
//
// Solidity: function hashPegInQuoteEIP712((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegIn *IPegInCallerSession) HashPegInQuoteEIP712(quote QuotesPegInQuote) ([32]byte, error) {
	return _IPegIn.Contract.HashPegInQuoteEIP712(&_IPegIn.CallOpts, quote)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegIn *IPegInCaller) PauseStatus(opts *bind.CallOpts) (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "pauseStatus")

	outstruct := new(struct {
		IsPaused bool
		Reason   string
		Since    uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegIn *IPegInSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IPegIn.Contract.PauseStatus(&_IPegIn.CallOpts)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegIn *IPegInCallerSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IPegIn.Contract.PauseStatus(&_IPegIn.CallOpts)
}

// ValidatePegInDepositAddress is a free data retrieval call binding the contract method 0xe9accea2.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (_IPegIn *IPegInCaller) ValidatePegInDepositAddress(opts *bind.CallOpts, quote QuotesPegInQuote, depositAddress []byte) (bool, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "validatePegInDepositAddress", quote, depositAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidatePegInDepositAddress is a free data retrieval call binding the contract method 0xe9accea2.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (_IPegIn *IPegInSession) ValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) (bool, error) {
	return _IPegIn.Contract.ValidatePegInDepositAddress(&_IPegIn.CallOpts, quote, depositAddress)
}

// ValidatePegInDepositAddress is a free data retrieval call binding the contract method 0xe9accea2.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (_IPegIn *IPegInCallerSession) ValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) (bool, error) {
	return _IPegIn.Contract.ValidatePegInDepositAddress(&_IPegIn.CallOpts, quote, depositAddress)
}

// CallForUser is a paid mutator transaction binding the contract method 0xc7a3dc3c.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (_IPegIn *IPegInTransactor) CallForUser(opts *bind.TransactOpts, quote QuotesPegInQuote) (*types.Transaction, error) {
	return _IPegIn.contract.Transact(opts, "callForUser", quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xc7a3dc3c.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (_IPegIn *IPegInSession) CallForUser(quote QuotesPegInQuote) (*types.Transaction, error) {
	return _IPegIn.Contract.CallForUser(&_IPegIn.TransactOpts, quote)
}

// CallForUser is a paid mutator transaction binding the contract method 0xc7a3dc3c.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (_IPegIn *IPegInTransactorSession) CallForUser(quote QuotesPegInQuote) (*types.Transaction, error) {
	return _IPegIn.Contract.CallForUser(&_IPegIn.TransactOpts, quote)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IPegIn *IPegInTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegIn.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IPegIn *IPegInSession) Deposit() (*types.Transaction, error) {
	return _IPegIn.Contract.Deposit(&_IPegIn.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IPegIn *IPegInTransactorSession) Deposit() (*types.Transaction, error) {
	return _IPegIn.Contract.Deposit(&_IPegIn.TransactOpts)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x3823c753.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_IPegIn *IPegInTransactor) RegisterPegIn(opts *bind.TransactOpts, quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _IPegIn.contract.Transact(opts, "registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x3823c753.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_IPegIn *IPegInSession) RegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _IPegIn.Contract.RegisterPegIn(&_IPegIn.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// RegisterPegIn is a paid mutator transaction binding the contract method 0x3823c753.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (_IPegIn *IPegInTransactorSession) RegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) (*types.Transaction, error) {
	return _IPegIn.Contract.RegisterPegIn(&_IPegIn.TransactOpts, quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IPegIn *IPegInTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _IPegIn.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IPegIn *IPegInSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _IPegIn.Contract.Withdraw(&_IPegIn.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_IPegIn *IPegInTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _IPegIn.Contract.Withdraw(&_IPegIn.TransactOpts, amount)
}

// IPegInBalanceDecreaseIterator is returned from FilterBalanceDecrease and is used to iterate over the raw logs and unpacked data for BalanceDecrease events raised by the IPegIn contract.
type IPegInBalanceDecreaseIterator struct {
	Event *IPegInBalanceDecrease // Event containing the contract specifics and raw log

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
func (it *IPegInBalanceDecreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInBalanceDecrease)
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
		it.Event = new(IPegInBalanceDecrease)
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
func (it *IPegInBalanceDecreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInBalanceDecreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInBalanceDecrease represents a BalanceDecrease event raised by the IPegIn contract.
type IPegInBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceDecrease is a free log retrieval operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) FilterBalanceDecrease(opts *bind.FilterOpts, dest []common.Address, amount []*big.Int) (*IPegInBalanceDecreaseIterator, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "BalanceDecrease", destRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegInBalanceDecreaseIterator{contract: _IPegIn.contract, event: "BalanceDecrease", logs: logs, sub: sub}, nil
}

// WatchBalanceDecrease is a free log subscription operation binding the contract event 0x8e51a4493a6f66c76e13fd9e3b754eafbfe21343c04508deb61be8ccc0064587.
//
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) WatchBalanceDecrease(opts *bind.WatchOpts, sink chan<- *IPegInBalanceDecrease, dest []common.Address, amount []*big.Int) (event.Subscription, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "BalanceDecrease", destRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInBalanceDecrease)
				if err := _IPegIn.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
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
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) ParseBalanceDecrease(log types.Log) (*IPegInBalanceDecrease, error) {
	event := new(IPegInBalanceDecrease)
	if err := _IPegIn.contract.UnpackLog(event, "BalanceDecrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInBalanceIncreaseIterator is returned from FilterBalanceIncrease and is used to iterate over the raw logs and unpacked data for BalanceIncrease events raised by the IPegIn contract.
type IPegInBalanceIncreaseIterator struct {
	Event *IPegInBalanceIncrease // Event containing the contract specifics and raw log

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
func (it *IPegInBalanceIncreaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInBalanceIncrease)
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
		it.Event = new(IPegInBalanceIncrease)
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
func (it *IPegInBalanceIncreaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInBalanceIncreaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInBalanceIncrease represents a BalanceIncrease event raised by the IPegIn contract.
type IPegInBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBalanceIncrease is a free log retrieval operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) FilterBalanceIncrease(opts *bind.FilterOpts, dest []common.Address, amount []*big.Int) (*IPegInBalanceIncreaseIterator, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "BalanceIncrease", destRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegInBalanceIncreaseIterator{contract: _IPegIn.contract, event: "BalanceIncrease", logs: logs, sub: sub}, nil
}

// WatchBalanceIncrease is a free log subscription operation binding the contract event 0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53.
//
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) WatchBalanceIncrease(opts *bind.WatchOpts, sink chan<- *IPegInBalanceIncrease, dest []common.Address, amount []*big.Int) (event.Subscription, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "BalanceIncrease", destRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInBalanceIncrease)
				if err := _IPegIn.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
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
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) ParseBalanceIncrease(log types.Log) (*IPegInBalanceIncrease, error) {
	event := new(IPegInBalanceIncrease)
	if err := _IPegIn.contract.UnpackLog(event, "BalanceIncrease", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInBridgeCapExceededIterator is returned from FilterBridgeCapExceeded and is used to iterate over the raw logs and unpacked data for BridgeCapExceeded events raised by the IPegIn contract.
type IPegInBridgeCapExceededIterator struct {
	Event *IPegInBridgeCapExceeded // Event containing the contract specifics and raw log

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
func (it *IPegInBridgeCapExceededIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInBridgeCapExceeded)
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
		it.Event = new(IPegInBridgeCapExceeded)
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
func (it *IPegInBridgeCapExceededIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInBridgeCapExceededIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInBridgeCapExceeded represents a BridgeCapExceeded event raised by the IPegIn contract.
type IPegInBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBridgeCapExceeded is a free log retrieval operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 indexed quoteHash, int256 indexed errorCode)
func (_IPegIn *IPegInFilterer) FilterBridgeCapExceeded(opts *bind.FilterOpts, quoteHash [][32]byte, errorCode []*big.Int) (*IPegInBridgeCapExceededIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var errorCodeRule []interface{}
	for _, errorCodeItem := range errorCode {
		errorCodeRule = append(errorCodeRule, errorCodeItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "BridgeCapExceeded", quoteHashRule, errorCodeRule)
	if err != nil {
		return nil, err
	}
	return &IPegInBridgeCapExceededIterator{contract: _IPegIn.contract, event: "BridgeCapExceeded", logs: logs, sub: sub}, nil
}

// WatchBridgeCapExceeded is a free log subscription operation binding the contract event 0xfb209329d5ab5b7bcb2e92f45f4534814b6e68fa5ad1f171dabc1d17d26f0ebe.
//
// Solidity: event BridgeCapExceeded(bytes32 indexed quoteHash, int256 indexed errorCode)
func (_IPegIn *IPegInFilterer) WatchBridgeCapExceeded(opts *bind.WatchOpts, sink chan<- *IPegInBridgeCapExceeded, quoteHash [][32]byte, errorCode []*big.Int) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var errorCodeRule []interface{}
	for _, errorCodeItem := range errorCode {
		errorCodeRule = append(errorCodeRule, errorCodeItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "BridgeCapExceeded", quoteHashRule, errorCodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInBridgeCapExceeded)
				if err := _IPegIn.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
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
// Solidity: event BridgeCapExceeded(bytes32 indexed quoteHash, int256 indexed errorCode)
func (_IPegIn *IPegInFilterer) ParseBridgeCapExceeded(log types.Log) (*IPegInBridgeCapExceeded, error) {
	event := new(IPegInBridgeCapExceeded)
	if err := _IPegIn.contract.UnpackLog(event, "BridgeCapExceeded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInCallForUserIterator is returned from FilterCallForUser and is used to iterate over the raw logs and unpacked data for CallForUser events raised by the IPegIn contract.
type IPegInCallForUserIterator struct {
	Event *IPegInCallForUser // Event containing the contract specifics and raw log

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
func (it *IPegInCallForUserIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInCallForUser)
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
		it.Event = new(IPegInCallForUser)
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
func (it *IPegInCallForUserIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInCallForUserIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInCallForUser represents a CallForUser event raised by the IPegIn contract.
type IPegInCallForUser struct {
	From      common.Address
	Dest      common.Address
	QuoteHash [32]byte
	GasLimit  *big.Int
	Value     *big.Int
	Data      []byte
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCallForUser is a free log retrieval operation binding the contract event 0x29a638a7bf9fc6a3c0bdf6ad339d1bba4555c740f5f80ddd3747cfe8dae172d9.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, bytes32 indexed quoteHash, uint256 gasLimit, uint256 value, bytes data, bool success)
func (_IPegIn *IPegInFilterer) FilterCallForUser(opts *bind.FilterOpts, from []common.Address, dest []common.Address, quoteHash [][32]byte) (*IPegInCallForUserIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "CallForUser", fromRule, destRule, quoteHashRule)
	if err != nil {
		return nil, err
	}
	return &IPegInCallForUserIterator{contract: _IPegIn.contract, event: "CallForUser", logs: logs, sub: sub}, nil
}

// WatchCallForUser is a free log subscription operation binding the contract event 0x29a638a7bf9fc6a3c0bdf6ad339d1bba4555c740f5f80ddd3747cfe8dae172d9.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, bytes32 indexed quoteHash, uint256 gasLimit, uint256 value, bytes data, bool success)
func (_IPegIn *IPegInFilterer) WatchCallForUser(opts *bind.WatchOpts, sink chan<- *IPegInCallForUser, from []common.Address, dest []common.Address, quoteHash [][32]byte) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "CallForUser", fromRule, destRule, quoteHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInCallForUser)
				if err := _IPegIn.contract.UnpackLog(event, "CallForUser", log); err != nil {
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

// ParseCallForUser is a log parse operation binding the contract event 0x29a638a7bf9fc6a3c0bdf6ad339d1bba4555c740f5f80ddd3747cfe8dae172d9.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, bytes32 indexed quoteHash, uint256 gasLimit, uint256 value, bytes data, bool success)
func (_IPegIn *IPegInFilterer) ParseCallForUser(log types.Log) (*IPegInCallForUser, error) {
	event := new(IPegInCallForUser)
	if err := _IPegIn.contract.UnpackLog(event, "CallForUser", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the IPegIn contract.
type IPegInEIP712DomainChangedIterator struct {
	Event *IPegInEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *IPegInEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInEIP712DomainChanged)
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
		it.Event = new(IPegInEIP712DomainChanged)
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
func (it *IPegInEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInEIP712DomainChanged represents a EIP712DomainChanged event raised by the IPegIn contract.
type IPegInEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegIn *IPegInFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*IPegInEIP712DomainChangedIterator, error) {

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &IPegInEIP712DomainChangedIterator{contract: _IPegIn.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegIn *IPegInFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *IPegInEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInEIP712DomainChanged)
				if err := _IPegIn.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegIn *IPegInFilterer) ParseEIP712DomainChanged(log types.Log) (*IPegInEIP712DomainChanged, error) {
	event := new(IPegInEIP712DomainChanged)
	if err := _IPegIn.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInPegInRegisteredIterator is returned from FilterPegInRegistered and is used to iterate over the raw logs and unpacked data for PegInRegistered events raised by the IPegIn contract.
type IPegInPegInRegisteredIterator struct {
	Event *IPegInPegInRegistered // Event containing the contract specifics and raw log

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
func (it *IPegInPegInRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInPegInRegistered)
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
		it.Event = new(IPegInPegInRegistered)
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
func (it *IPegInPegInRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInPegInRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInPegInRegistered represents a PegInRegistered event raised by the IPegIn contract.
type IPegInPegInRegistered struct {
	QuoteHash         [32]byte
	TransferredAmount *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPegInRegistered is a free log retrieval operation binding the contract event 0x0405e68a1f0887bc595391af9c93e3d8ac89077d862794c6bc78e061ead0f170.
//
// Solidity: event PegInRegistered(bytes32 indexed quoteHash, uint256 indexed transferredAmount)
func (_IPegIn *IPegInFilterer) FilterPegInRegistered(opts *bind.FilterOpts, quoteHash [][32]byte, transferredAmount []*big.Int) (*IPegInPegInRegisteredIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var transferredAmountRule []interface{}
	for _, transferredAmountItem := range transferredAmount {
		transferredAmountRule = append(transferredAmountRule, transferredAmountItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "PegInRegistered", quoteHashRule, transferredAmountRule)
	if err != nil {
		return nil, err
	}
	return &IPegInPegInRegisteredIterator{contract: _IPegIn.contract, event: "PegInRegistered", logs: logs, sub: sub}, nil
}

// WatchPegInRegistered is a free log subscription operation binding the contract event 0x0405e68a1f0887bc595391af9c93e3d8ac89077d862794c6bc78e061ead0f170.
//
// Solidity: event PegInRegistered(bytes32 indexed quoteHash, uint256 indexed transferredAmount)
func (_IPegIn *IPegInFilterer) WatchPegInRegistered(opts *bind.WatchOpts, sink chan<- *IPegInPegInRegistered, quoteHash [][32]byte, transferredAmount []*big.Int) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var transferredAmountRule []interface{}
	for _, transferredAmountItem := range transferredAmount {
		transferredAmountRule = append(transferredAmountRule, transferredAmountItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "PegInRegistered", quoteHashRule, transferredAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInPegInRegistered)
				if err := _IPegIn.contract.UnpackLog(event, "PegInRegistered", log); err != nil {
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

// ParsePegInRegistered is a log parse operation binding the contract event 0x0405e68a1f0887bc595391af9c93e3d8ac89077d862794c6bc78e061ead0f170.
//
// Solidity: event PegInRegistered(bytes32 indexed quoteHash, uint256 indexed transferredAmount)
func (_IPegIn *IPegInFilterer) ParsePegInRegistered(log types.Log) (*IPegInPegInRegistered, error) {
	event := new(IPegInPegInRegistered)
	if err := _IPegIn.contract.UnpackLog(event, "PegInRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInRefundIterator is returned from FilterRefund and is used to iterate over the raw logs and unpacked data for Refund events raised by the IPegIn contract.
type IPegInRefundIterator struct {
	Event *IPegInRefund // Event containing the contract specifics and raw log

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
func (it *IPegInRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInRefund)
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
		it.Event = new(IPegInRefund)
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
func (it *IPegInRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInRefund represents a Refund event raised by the IPegIn contract.
type IPegInRefund struct {
	Dest      common.Address
	QuoteHash [32]byte
	Amount    *big.Int
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRefund is a free log retrieval operation binding the contract event 0xc1724559d229f11ed3d78e5fec91a4f596a160bf689f7a40ff3a2a8230cb9515.
//
// Solidity: event Refund(address indexed dest, bytes32 indexed quoteHash, uint256 indexed amount, bool success)
func (_IPegIn *IPegInFilterer) FilterRefund(opts *bind.FilterOpts, dest []common.Address, quoteHash [][32]byte, amount []*big.Int) (*IPegInRefundIterator, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "Refund", destRule, quoteHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegInRefundIterator{contract: _IPegIn.contract, event: "Refund", logs: logs, sub: sub}, nil
}

// WatchRefund is a free log subscription operation binding the contract event 0xc1724559d229f11ed3d78e5fec91a4f596a160bf689f7a40ff3a2a8230cb9515.
//
// Solidity: event Refund(address indexed dest, bytes32 indexed quoteHash, uint256 indexed amount, bool success)
func (_IPegIn *IPegInFilterer) WatchRefund(opts *bind.WatchOpts, sink chan<- *IPegInRefund, dest []common.Address, quoteHash [][32]byte, amount []*big.Int) (event.Subscription, error) {

	var destRule []interface{}
	for _, destItem := range dest {
		destRule = append(destRule, destItem)
	}
	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "Refund", destRule, quoteHashRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInRefund)
				if err := _IPegIn.contract.UnpackLog(event, "Refund", log); err != nil {
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

// ParseRefund is a log parse operation binding the contract event 0xc1724559d229f11ed3d78e5fec91a4f596a160bf689f7a40ff3a2a8230cb9515.
//
// Solidity: event Refund(address indexed dest, bytes32 indexed quoteHash, uint256 indexed amount, bool success)
func (_IPegIn *IPegInFilterer) ParseRefund(log types.Log) (*IPegInRefund, error) {
	event := new(IPegInRefund)
	if err := _IPegIn.contract.UnpackLog(event, "Refund", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegInWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the IPegIn contract.
type IPegInWithdrawalIterator struct {
	Event *IPegInWithdrawal // Event containing the contract specifics and raw log

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
func (it *IPegInWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegInWithdrawal)
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
		it.Event = new(IPegInWithdrawal)
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
func (it *IPegInWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegInWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegInWithdrawal represents a Withdrawal event raised by the IPegIn contract.
type IPegInWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed from, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) FilterWithdrawal(opts *bind.FilterOpts, from []common.Address, amount []*big.Int) (*IPegInWithdrawalIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.FilterLogs(opts, "Withdrawal", fromRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegInWithdrawalIterator{contract: _IPegIn.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed from, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *IPegInWithdrawal, from []common.Address, amount []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegIn.contract.WatchLogs(opts, "Withdrawal", fromRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegInWithdrawal)
				if err := _IPegIn.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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
// Solidity: event Withdrawal(address indexed from, uint256 indexed amount)
func (_IPegIn *IPegInFilterer) ParseWithdrawal(log types.Log) (*IPegInWithdrawal, error) {
	event := new(IPegInWithdrawal)
	if err := _IPegIn.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutMetaData contains all meta data concerning the IPegOut contract.
var IPegOutMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"depositPegOut\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegOutQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegOutQuoteEIP712\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"isQuoteCompleted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isPaused\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"since\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"btcTx\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"btcBlockHeaderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"refundPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"refundUserPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"btcTx\",\"type\":\"bytes\"}],\"name\":\"validatePegout\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"change\",\"type\":\"uint256\"}],\"name\":\"PegOutChangePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"PegOutRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"PegOutUserRefunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expected\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"actual\",\"type\":\"bytes\"}],\"name\":\"InvalidDestination\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"InvalidQuoteHash\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"outputScript\",\"type\":\"bytes\"}],\"name\":\"MalformedTransaction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"required\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"}],\"name\":\"NotEnoughConfirmations\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyCompleted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"name\":\"QuoteExpiredByBlocks\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"}],\"name\":\"QuoteExpiredByTime\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteNotExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"UnableToGetConfirmations\",\"type\":\"error\"}]",
}

// IPegOutABI is the input ABI used to generate the binding from.
// Deprecated: Use IPegOutMetaData.ABI instead.
var IPegOutABI = IPegOutMetaData.ABI

// IPegOut is an auto generated Go binding around an Ethereum contract.
type IPegOut struct {
	IPegOutCaller     // Read-only binding to the contract
	IPegOutTransactor // Write-only binding to the contract
	IPegOutFilterer   // Log filterer for contract events
}

// IPegOutCaller is an auto generated read-only Go binding around an Ethereum contract.
type IPegOutCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegOutTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IPegOutTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegOutFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IPegOutFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPegOutSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IPegOutSession struct {
	Contract     *IPegOut          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IPegOutCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IPegOutCallerSession struct {
	Contract *IPegOutCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IPegOutTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IPegOutTransactorSession struct {
	Contract     *IPegOutTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IPegOutRaw is an auto generated low-level Go binding around an Ethereum contract.
type IPegOutRaw struct {
	Contract *IPegOut // Generic contract binding to access the raw methods on
}

// IPegOutCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IPegOutCallerRaw struct {
	Contract *IPegOutCaller // Generic read-only contract binding to access the raw methods on
}

// IPegOutTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IPegOutTransactorRaw struct {
	Contract *IPegOutTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPegOut creates a new instance of IPegOut, bound to a specific deployed contract.
func NewIPegOut(address common.Address, backend bind.ContractBackend) (*IPegOut, error) {
	contract, err := bindIPegOut(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPegOut{IPegOutCaller: IPegOutCaller{contract: contract}, IPegOutTransactor: IPegOutTransactor{contract: contract}, IPegOutFilterer: IPegOutFilterer{contract: contract}}, nil
}

// NewIPegOutCaller creates a new read-only instance of IPegOut, bound to a specific deployed contract.
func NewIPegOutCaller(address common.Address, caller bind.ContractCaller) (*IPegOutCaller, error) {
	contract, err := bindIPegOut(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPegOutCaller{contract: contract}, nil
}

// NewIPegOutTransactor creates a new write-only instance of IPegOut, bound to a specific deployed contract.
func NewIPegOutTransactor(address common.Address, transactor bind.ContractTransactor) (*IPegOutTransactor, error) {
	contract, err := bindIPegOut(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPegOutTransactor{contract: contract}, nil
}

// NewIPegOutFilterer creates a new log filterer instance of IPegOut, bound to a specific deployed contract.
func NewIPegOutFilterer(address common.Address, filterer bind.ContractFilterer) (*IPegOutFilterer, error) {
	contract, err := bindIPegOut(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPegOutFilterer{contract: contract}, nil
}

// bindIPegOut binds a generic wrapper to an already deployed contract.
func bindIPegOut(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPegOutMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPegOut *IPegOutRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPegOut.Contract.IPegOutCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPegOut *IPegOutRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.Contract.IPegOutTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPegOut *IPegOutRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPegOut.Contract.IPegOutTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPegOut *IPegOutCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPegOut.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPegOut *IPegOutTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPegOut *IPegOutTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPegOut.Contract.contract.Transact(opts, method, params...)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegOut *IPegOutCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegOut *IPegOutSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IPegOut.Contract.Eip712Domain(&_IPegOut.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_IPegOut *IPegOutCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _IPegOut.Contract.Eip712Domain(&_IPegOut.CallOpts)
}

// HashPegOutQuote is a free data retrieval call binding the contract method 0x6408f6fe.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutCaller) HashPegOutQuote(opts *bind.CallOpts, quote QuotesPegOutQuote) ([32]byte, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "hashPegOutQuote", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegOutQuote is a free data retrieval call binding the contract method 0x6408f6fe.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutSession) HashPegOutQuote(quote QuotesPegOutQuote) ([32]byte, error) {
	return _IPegOut.Contract.HashPegOutQuote(&_IPegOut.CallOpts, quote)
}

// HashPegOutQuote is a free data retrieval call binding the contract method 0x6408f6fe.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutCallerSession) HashPegOutQuote(quote QuotesPegOutQuote) ([32]byte, error) {
	return _IPegOut.Contract.HashPegOutQuote(&_IPegOut.CallOpts, quote)
}

// HashPegOutQuoteEIP712 is a free data retrieval call binding the contract method 0x5966252a.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutCaller) HashPegOutQuoteEIP712(opts *bind.CallOpts, quote QuotesPegOutQuote) ([32]byte, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "hashPegOutQuoteEIP712", quote)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashPegOutQuoteEIP712 is a free data retrieval call binding the contract method 0x5966252a.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutSession) HashPegOutQuoteEIP712(quote QuotesPegOutQuote) ([32]byte, error) {
	return _IPegOut.Contract.HashPegOutQuoteEIP712(&_IPegOut.CallOpts, quote)
}

// HashPegOutQuoteEIP712 is a free data retrieval call binding the contract method 0x5966252a.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (_IPegOut *IPegOutCallerSession) HashPegOutQuoteEIP712(quote QuotesPegOutQuote) ([32]byte, error) {
	return _IPegOut.Contract.HashPegOutQuoteEIP712(&_IPegOut.CallOpts, quote)
}

// IsQuoteCompleted is a free data retrieval call binding the contract method 0x35bf61f1.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (_IPegOut *IPegOutCaller) IsQuoteCompleted(opts *bind.CallOpts, quoteHash [32]byte) (bool, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "isQuoteCompleted", quoteHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsQuoteCompleted is a free data retrieval call binding the contract method 0x35bf61f1.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (_IPegOut *IPegOutSession) IsQuoteCompleted(quoteHash [32]byte) (bool, error) {
	return _IPegOut.Contract.IsQuoteCompleted(&_IPegOut.CallOpts, quoteHash)
}

// IsQuoteCompleted is a free data retrieval call binding the contract method 0x35bf61f1.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (_IPegOut *IPegOutCallerSession) IsQuoteCompleted(quoteHash [32]byte) (bool, error) {
	return _IPegOut.Contract.IsQuoteCompleted(&_IPegOut.CallOpts, quoteHash)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegOut *IPegOutCaller) PauseStatus(opts *bind.CallOpts) (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "pauseStatus")

	outstruct := new(struct {
		IsPaused bool
		Reason   string
		Since    uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegOut *IPegOutSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IPegOut.Contract.PauseStatus(&_IPegOut.CallOpts)
}

// PauseStatus is a free data retrieval call binding the contract method 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (_IPegOut *IPegOutCallerSession) PauseStatus() (struct {
	IsPaused bool
	Reason   string
	Since    uint64
}, error) {
	return _IPegOut.Contract.PauseStatus(&_IPegOut.CallOpts)
}

// ValidatePegout is a free data retrieval call binding the contract method 0x7846150c.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (_IPegOut *IPegOutCaller) ValidatePegout(opts *bind.CallOpts, quoteHash [32]byte, btcTx []byte) (QuotesPegOutQuote, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "validatePegout", quoteHash, btcTx)

	if err != nil {
		return *new(QuotesPegOutQuote), err
	}

	out0 := *abi.ConvertType(out[0], new(QuotesPegOutQuote)).(*QuotesPegOutQuote)

	return out0, err

}

// ValidatePegout is a free data retrieval call binding the contract method 0x7846150c.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (_IPegOut *IPegOutSession) ValidatePegout(quoteHash [32]byte, btcTx []byte) (QuotesPegOutQuote, error) {
	return _IPegOut.Contract.ValidatePegout(&_IPegOut.CallOpts, quoteHash, btcTx)
}

// ValidatePegout is a free data retrieval call binding the contract method 0x7846150c.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (_IPegOut *IPegOutCallerSession) ValidatePegout(quoteHash [32]byte, btcTx []byte) (QuotesPegOutQuote, error) {
	return _IPegOut.Contract.ValidatePegout(&_IPegOut.CallOpts, quoteHash, btcTx)
}

// DepositPegOut is a paid mutator transaction binding the contract method 0x083bc4b2.
//
// Solidity: function depositPegOut((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes signature) payable returns()
func (_IPegOut *IPegOutTransactor) DepositPegOut(opts *bind.TransactOpts, quote QuotesPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "depositPegOut", quote, signature)
}

// DepositPegOut is a paid mutator transaction binding the contract method 0x083bc4b2.
//
// Solidity: function depositPegOut((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes signature) payable returns()
func (_IPegOut *IPegOutSession) DepositPegOut(quote QuotesPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _IPegOut.Contract.DepositPegOut(&_IPegOut.TransactOpts, quote, signature)
}

// DepositPegOut is a paid mutator transaction binding the contract method 0x083bc4b2.
//
// Solidity: function depositPegOut((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes signature) payable returns()
func (_IPegOut *IPegOutTransactorSession) DepositPegOut(quote QuotesPegOutQuote, signature []byte) (*types.Transaction, error) {
	return _IPegOut.Contract.DepositPegOut(&_IPegOut.TransactOpts, quote, signature)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xd6c70de8.
//
// Solidity: function refundPegOut(bytes32 quoteHash, bytes btcTx, bytes32 btcBlockHeaderHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (_IPegOut *IPegOutTransactor) RefundPegOut(opts *bind.TransactOpts, quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "refundPegOut", quoteHash, btcTx, btcBlockHeaderHash, merkleBranchPath, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xd6c70de8.
//
// Solidity: function refundPegOut(bytes32 quoteHash, bytes btcTx, bytes32 btcBlockHeaderHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (_IPegOut *IPegOutSession) RefundPegOut(quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _IPegOut.Contract.RefundPegOut(&_IPegOut.TransactOpts, quoteHash, btcTx, btcBlockHeaderHash, merkleBranchPath, merkleBranchHashes)
}

// RefundPegOut is a paid mutator transaction binding the contract method 0xd6c70de8.
//
// Solidity: function refundPegOut(bytes32 quoteHash, bytes btcTx, bytes32 btcBlockHeaderHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (_IPegOut *IPegOutTransactorSession) RefundPegOut(quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) (*types.Transaction, error) {
	return _IPegOut.Contract.RefundPegOut(&_IPegOut.TransactOpts, quoteHash, btcTx, btcBlockHeaderHash, merkleBranchPath, merkleBranchHashes)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0x8f91797d.
//
// Solidity: function refundUserPegOut(bytes32 quoteHash) returns()
func (_IPegOut *IPegOutTransactor) RefundUserPegOut(opts *bind.TransactOpts, quoteHash [32]byte) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "refundUserPegOut", quoteHash)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0x8f91797d.
//
// Solidity: function refundUserPegOut(bytes32 quoteHash) returns()
func (_IPegOut *IPegOutSession) RefundUserPegOut(quoteHash [32]byte) (*types.Transaction, error) {
	return _IPegOut.Contract.RefundUserPegOut(&_IPegOut.TransactOpts, quoteHash)
}

// RefundUserPegOut is a paid mutator transaction binding the contract method 0x8f91797d.
//
// Solidity: function refundUserPegOut(bytes32 quoteHash) returns()
func (_IPegOut *IPegOutTransactorSession) RefundUserPegOut(quoteHash [32]byte) (*types.Transaction, error) {
	return _IPegOut.Contract.RefundUserPegOut(&_IPegOut.TransactOpts, quoteHash)
}

// IPegOutEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the IPegOut contract.
type IPegOutEIP712DomainChangedIterator struct {
	Event *IPegOutEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *IPegOutEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutEIP712DomainChanged)
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
		it.Event = new(IPegOutEIP712DomainChanged)
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
func (it *IPegOutEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutEIP712DomainChanged represents a EIP712DomainChanged event raised by the IPegOut contract.
type IPegOutEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegOut *IPegOutFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*IPegOutEIP712DomainChangedIterator, error) {

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &IPegOutEIP712DomainChangedIterator{contract: _IPegOut.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegOut *IPegOutFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *IPegOutEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutEIP712DomainChanged)
				if err := _IPegOut.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_IPegOut *IPegOutFilterer) ParseEIP712DomainChanged(log types.Log) (*IPegOutEIP712DomainChanged, error) {
	event := new(IPegOutEIP712DomainChanged)
	if err := _IPegOut.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutPegOutChangePaidIterator is returned from FilterPegOutChangePaid and is used to iterate over the raw logs and unpacked data for PegOutChangePaid events raised by the IPegOut contract.
type IPegOutPegOutChangePaidIterator struct {
	Event *IPegOutPegOutChangePaid // Event containing the contract specifics and raw log

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
func (it *IPegOutPegOutChangePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutPegOutChangePaid)
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
		it.Event = new(IPegOutPegOutChangePaid)
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
func (it *IPegOutPegOutChangePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutPegOutChangePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutPegOutChangePaid represents a PegOutChangePaid event raised by the IPegOut contract.
type IPegOutPegOutChangePaid struct {
	QuoteHash   [32]byte
	UserAddress common.Address
	Change      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPegOutChangePaid is a free log retrieval operation binding the contract event 0xc5de476c050573583f2f9b92de1a4e269559a6da0d5935920753c74ee4edbca0.
//
// Solidity: event PegOutChangePaid(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed change)
func (_IPegOut *IPegOutFilterer) FilterPegOutChangePaid(opts *bind.FilterOpts, quoteHash [][32]byte, userAddress []common.Address, change []*big.Int) (*IPegOutPegOutChangePaidIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var userAddressRule []interface{}
	for _, userAddressItem := range userAddress {
		userAddressRule = append(userAddressRule, userAddressItem)
	}
	var changeRule []interface{}
	for _, changeItem := range change {
		changeRule = append(changeRule, changeItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "PegOutChangePaid", quoteHashRule, userAddressRule, changeRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutPegOutChangePaidIterator{contract: _IPegOut.contract, event: "PegOutChangePaid", logs: logs, sub: sub}, nil
}

// WatchPegOutChangePaid is a free log subscription operation binding the contract event 0xc5de476c050573583f2f9b92de1a4e269559a6da0d5935920753c74ee4edbca0.
//
// Solidity: event PegOutChangePaid(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed change)
func (_IPegOut *IPegOutFilterer) WatchPegOutChangePaid(opts *bind.WatchOpts, sink chan<- *IPegOutPegOutChangePaid, quoteHash [][32]byte, userAddress []common.Address, change []*big.Int) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var userAddressRule []interface{}
	for _, userAddressItem := range userAddress {
		userAddressRule = append(userAddressRule, userAddressItem)
	}
	var changeRule []interface{}
	for _, changeItem := range change {
		changeRule = append(changeRule, changeItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "PegOutChangePaid", quoteHashRule, userAddressRule, changeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutPegOutChangePaid)
				if err := _IPegOut.contract.UnpackLog(event, "PegOutChangePaid", log); err != nil {
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

// ParsePegOutChangePaid is a log parse operation binding the contract event 0xc5de476c050573583f2f9b92de1a4e269559a6da0d5935920753c74ee4edbca0.
//
// Solidity: event PegOutChangePaid(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed change)
func (_IPegOut *IPegOutFilterer) ParsePegOutChangePaid(log types.Log) (*IPegOutPegOutChangePaid, error) {
	event := new(IPegOutPegOutChangePaid)
	if err := _IPegOut.contract.UnpackLog(event, "PegOutChangePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutPegOutDepositIterator is returned from FilterPegOutDeposit and is used to iterate over the raw logs and unpacked data for PegOutDeposit events raised by the IPegOut contract.
type IPegOutPegOutDepositIterator struct {
	Event *IPegOutPegOutDeposit // Event containing the contract specifics and raw log

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
func (it *IPegOutPegOutDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutPegOutDeposit)
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
		it.Event = new(IPegOutPegOutDeposit)
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
func (it *IPegOutPegOutDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutPegOutDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutPegOutDeposit represents a PegOutDeposit event raised by the IPegOut contract.
type IPegOutPegOutDeposit struct {
	QuoteHash [32]byte
	Sender    common.Address
	Timestamp *big.Int
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPegOutDeposit is a free log retrieval operation binding the contract event 0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f.
//
// Solidity: event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 indexed timestamp, uint256 amount)
func (_IPegOut *IPegOutFilterer) FilterPegOutDeposit(opts *bind.FilterOpts, quoteHash [][32]byte, sender []common.Address, timestamp []*big.Int) (*IPegOutPegOutDepositIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "PegOutDeposit", quoteHashRule, senderRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutPegOutDepositIterator{contract: _IPegOut.contract, event: "PegOutDeposit", logs: logs, sub: sub}, nil
}

// WatchPegOutDeposit is a free log subscription operation binding the contract event 0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f.
//
// Solidity: event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 indexed timestamp, uint256 amount)
func (_IPegOut *IPegOutFilterer) WatchPegOutDeposit(opts *bind.WatchOpts, sink chan<- *IPegOutPegOutDeposit, quoteHash [][32]byte, sender []common.Address, timestamp []*big.Int) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var timestampRule []interface{}
	for _, timestampItem := range timestamp {
		timestampRule = append(timestampRule, timestampItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "PegOutDeposit", quoteHashRule, senderRule, timestampRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutPegOutDeposit)
				if err := _IPegOut.contract.UnpackLog(event, "PegOutDeposit", log); err != nil {
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

// ParsePegOutDeposit is a log parse operation binding the contract event 0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f.
//
// Solidity: event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 indexed timestamp, uint256 amount)
func (_IPegOut *IPegOutFilterer) ParsePegOutDeposit(log types.Log) (*IPegOutPegOutDeposit, error) {
	event := new(IPegOutPegOutDeposit)
	if err := _IPegOut.contract.UnpackLog(event, "PegOutDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutPegOutRefundedIterator is returned from FilterPegOutRefunded and is used to iterate over the raw logs and unpacked data for PegOutRefunded events raised by the IPegOut contract.
type IPegOutPegOutRefundedIterator struct {
	Event *IPegOutPegOutRefunded // Event containing the contract specifics and raw log

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
func (it *IPegOutPegOutRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutPegOutRefunded)
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
		it.Event = new(IPegOutPegOutRefunded)
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
func (it *IPegOutPegOutRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutPegOutRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutPegOutRefunded represents a PegOutRefunded event raised by the IPegOut contract.
type IPegOutPegOutRefunded struct {
	QuoteHash [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPegOutRefunded is a free log retrieval operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 indexed quoteHash)
func (_IPegOut *IPegOutFilterer) FilterPegOutRefunded(opts *bind.FilterOpts, quoteHash [][32]byte) (*IPegOutPegOutRefundedIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "PegOutRefunded", quoteHashRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutPegOutRefundedIterator{contract: _IPegOut.contract, event: "PegOutRefunded", logs: logs, sub: sub}, nil
}

// WatchPegOutRefunded is a free log subscription operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 indexed quoteHash)
func (_IPegOut *IPegOutFilterer) WatchPegOutRefunded(opts *bind.WatchOpts, sink chan<- *IPegOutPegOutRefunded, quoteHash [][32]byte) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "PegOutRefunded", quoteHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutPegOutRefunded)
				if err := _IPegOut.contract.UnpackLog(event, "PegOutRefunded", log); err != nil {
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

// ParsePegOutRefunded is a log parse operation binding the contract event 0xb781856ec73fd0dc39351043d1634ea22cd3277b0866ab93e7ec1801766bb384.
//
// Solidity: event PegOutRefunded(bytes32 indexed quoteHash)
func (_IPegOut *IPegOutFilterer) ParsePegOutRefunded(log types.Log) (*IPegOutPegOutRefunded, error) {
	event := new(IPegOutPegOutRefunded)
	if err := _IPegOut.contract.UnpackLog(event, "PegOutRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutPegOutUserRefundedIterator is returned from FilterPegOutUserRefunded and is used to iterate over the raw logs and unpacked data for PegOutUserRefunded events raised by the IPegOut contract.
type IPegOutPegOutUserRefundedIterator struct {
	Event *IPegOutPegOutUserRefunded // Event containing the contract specifics and raw log

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
func (it *IPegOutPegOutUserRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutPegOutUserRefunded)
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
		it.Event = new(IPegOutPegOutUserRefunded)
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
func (it *IPegOutPegOutUserRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutPegOutUserRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutPegOutUserRefunded represents a PegOutUserRefunded event raised by the IPegOut contract.
type IPegOutPegOutUserRefunded struct {
	QuoteHash   [32]byte
	UserAddress common.Address
	Value       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPegOutUserRefunded is a free log retrieval operation binding the contract event 0x7cafa4b6e67efda68d661be7bc7a1514d0f1864db803311ec6f5f5778fb2c92a.
//
// Solidity: event PegOutUserRefunded(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed value)
func (_IPegOut *IPegOutFilterer) FilterPegOutUserRefunded(opts *bind.FilterOpts, quoteHash [][32]byte, userAddress []common.Address, value []*big.Int) (*IPegOutPegOutUserRefundedIterator, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var userAddressRule []interface{}
	for _, userAddressItem := range userAddress {
		userAddressRule = append(userAddressRule, userAddressItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "PegOutUserRefunded", quoteHashRule, userAddressRule, valueRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutPegOutUserRefundedIterator{contract: _IPegOut.contract, event: "PegOutUserRefunded", logs: logs, sub: sub}, nil
}

// WatchPegOutUserRefunded is a free log subscription operation binding the contract event 0x7cafa4b6e67efda68d661be7bc7a1514d0f1864db803311ec6f5f5778fb2c92a.
//
// Solidity: event PegOutUserRefunded(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed value)
func (_IPegOut *IPegOutFilterer) WatchPegOutUserRefunded(opts *bind.WatchOpts, sink chan<- *IPegOutPegOutUserRefunded, quoteHash [][32]byte, userAddress []common.Address, value []*big.Int) (event.Subscription, error) {

	var quoteHashRule []interface{}
	for _, quoteHashItem := range quoteHash {
		quoteHashRule = append(quoteHashRule, quoteHashItem)
	}
	var userAddressRule []interface{}
	for _, userAddressItem := range userAddress {
		userAddressRule = append(userAddressRule, userAddressItem)
	}
	var valueRule []interface{}
	for _, valueItem := range value {
		valueRule = append(valueRule, valueItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "PegOutUserRefunded", quoteHashRule, userAddressRule, valueRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutPegOutUserRefunded)
				if err := _IPegOut.contract.UnpackLog(event, "PegOutUserRefunded", log); err != nil {
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

// ParsePegOutUserRefunded is a log parse operation binding the contract event 0x7cafa4b6e67efda68d661be7bc7a1514d0f1864db803311ec6f5f5778fb2c92a.
//
// Solidity: event PegOutUserRefunded(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed value)
func (_IPegOut *IPegOutFilterer) ParsePegOutUserRefunded(log types.Log) (*IPegOutPegOutUserRefunded, error) {
	event := new(IPegOutPegOutUserRefunded)
	if err := _IPegOut.contract.UnpackLog(event, "PegOutUserRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
