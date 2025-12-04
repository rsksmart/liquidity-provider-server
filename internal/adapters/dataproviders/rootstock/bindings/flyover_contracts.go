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
	CallFee                     *big.Int
	PenaltyFee                  *big.Int
	Value                       *big.Int
	ProductFeeAmount            *big.Int
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
	CallFee               *big.Int
	PenaltyFee            *big.Int
	Value                 *big.Int
	ProductFeeAmount      *big.Int
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

// QuotesV2PegOutQuote is an auto generated low-level Go binding around an user-defined struct.
type QuotesV2PegOutQuote struct {
	LbcAddress            common.Address
	LpRskAddress          common.Address
	BtcRefundAddress      []byte
	RskRefundAddress      common.Address
	LpBtcAddress          []byte
	CallFee               *big.Int
	PenaltyFee            *big.Int
	Nonce                 int64
	DeposityAddress       []byte
	Value                 *big.Int
	AgreementTimestamp    uint32
	DepositDateLimit      uint32
	DepositConfirmations  uint16
	TransferConfirmations uint16
	TransferTime          uint32
	ExpireDate            uint32
	ExpireBlock           uint32
	ProductFeeAmount      *big.Int
	GasFee                *big.Int
}

// FlyoverMetaData contains all meta data concerning the Flyover contract.
var FlyoverMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"heightOrHash\",\"type\":\"bytes32\"}],\"name\":\"EmptyBlockHeader\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"IncorrectContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"InsufficientAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"NoBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NoContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"passedAmount\",\"type\":\"uint256\"}],\"name\":\"Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"PaymentFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ProviderNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteNotFound\",\"type\":\"error\"}]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212204950f2ab918c24153aca756ac5d2f6ef8e1f2d64d15998f39f1a3c907a4966ad64736f6c63430008190033",
}

// FlyoverABI is the input ABI used to generate the binding from.
// Deprecated: Use FlyoverMetaData.ABI instead.
var FlyoverABI = FlyoverMetaData.ABI

// FlyoverBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FlyoverMetaData.Bin instead.
var FlyoverBin = FlyoverMetaData.Bin

// DeployFlyover deploys a new Ethereum contract, binding an instance of Flyover to it.
func DeployFlyover(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Flyover, error) {
	parsed, err := FlyoverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FlyoverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Flyover{FlyoverCaller: FlyoverCaller{contract: contract}, FlyoverTransactor: FlyoverTransactor{contract: contract}, FlyoverFilterer: FlyoverFilterer{contract: contract}}, nil
}

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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"key\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rskKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"mstKey\",\"type\":\"bytes\"}],\"name\":\"addFederatorPublicKeyMultikey\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"},{\"internalType\":\"int256\",\"name\":\"maxTransferValue\",\"type\":\"int256\"}],\"name\":\"addOneOffLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"txhash\",\"type\":\"bytes\"}],\"name\":\"addSignature\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"addUnlimitedLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"hash\",\"type\":\"bytes\"}],\"name\":\"commitFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveFederationCreationBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActivePowpegRedeemScript\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestBlockHeader\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainBestChainHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"depth\",\"type\":\"int256\"}],\"name\":\"getBtcBlockchainBlockHashAtDepth\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"btcBlockHeight\",\"type\":\"uint256\"}],\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBtcBlockchainInitialBlockHeight\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"btcBlockHash\",\"type\":\"bytes32\"}],\"name\":\"getBtcBlockchainParentBlockHeaderByHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"txHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"getBtcTransactionConfirmations\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"getBtcTxHashProcessedHeight\",\"outputs\":[{\"internalType\":\"int64\",\"name\":\"\",\"type\":\"int64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePerKb\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"getLockWhitelistEntryByAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockWhitelistSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLockingCap\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLockTxValue\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getPendingFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getPendingFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationBlockNumber\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationCreationTime\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationSize\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiringFederationThreshold\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"}],\"name\":\"getRetiringFederatorPublicKey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"index\",\"type\":\"int256\"},{\"internalType\":\"string\",\"name\":\"atype\",\"type\":\"string\"}],\"name\":\"getRetiringFederatorPublicKeyOfType\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForBtcReleaseClient\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStateForDebugging\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"}],\"name\":\"hasBtcBlockCoinbaseTransactionInformation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"newLockingCap\",\"type\":\"int256\"}],\"name\":\"increaseLockingCap\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"isBtcTxHashAlreadyProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"ablock\",\"type\":\"bytes\"}],\"name\":\"receiveHeader\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"blocks\",\"type\":\"bytes[]\"}],\"name\":\"receiveHeaders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"witnessMerkleRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"witnessReservedValue\",\"type\":\"bytes32\"}],\"name\":\"registerBtcCoinbaseTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"atx\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"height\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"pmt\",\"type\":\"bytes\"}],\"name\":\"registerBtcTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"btcTxSerialized\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"pmtSerialized\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bool\",\"name\":\"shouldTransferToContract\",\"type\":\"bool\"}],\"name\":\"registerFastBridgeBtcTransaction\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"aaddress\",\"type\":\"string\"}],\"name\":\"removeLockWhitelistAddress\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollbackFederation\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"disableDelay\",\"type\":\"int256\"}],\"name\":\"setLockWhitelistDisableBlockDelay\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateCollections\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"feePerKb\",\"type\":\"int256\"}],\"name\":\"voteFeePerKbChange\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"AlreadyResigned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotResigned\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NothingToWithdraw\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"resignationBlockNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"resignDelayInBlocks\",\"type\":\"uint256\"}],\"name\":\"ResignationDelayNotMet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawalFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegInCollateralAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutCollateralAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"liquidityProvider\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"collateralType\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"Penalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"Resigned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardsWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawCollateral\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addPegInCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addPegInCollateralTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addPegOutCollateral\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addPegOutCollateralTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegInCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getPegOutCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPenalties\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getResignDelayInBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getResignationBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRewardPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isCollateralSufficient\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resign\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"slashPegInCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"punisher\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"slashPegOutCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawCollateral\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// IDaoContributorMetaData contains all meta data concerning the IDaoContributor contract.
var IDaoContributorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getCurrentContribution\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IDaoContributorABI is the input ABI used to generate the binding from.
// Deprecated: Use IDaoContributorMetaData.ABI instead.
var IDaoContributorABI = IDaoContributorMetaData.ABI

// IDaoContributor is an auto generated Go binding around an Ethereum contract.
type IDaoContributor struct {
	IDaoContributorCaller     // Read-only binding to the contract
	IDaoContributorTransactor // Write-only binding to the contract
	IDaoContributorFilterer   // Log filterer for contract events
}

// IDaoContributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type IDaoContributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDaoContributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IDaoContributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDaoContributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IDaoContributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDaoContributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IDaoContributorSession struct {
	Contract     *IDaoContributor  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IDaoContributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IDaoContributorCallerSession struct {
	Contract *IDaoContributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IDaoContributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IDaoContributorTransactorSession struct {
	Contract     *IDaoContributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IDaoContributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type IDaoContributorRaw struct {
	Contract *IDaoContributor // Generic contract binding to access the raw methods on
}

// IDaoContributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IDaoContributorCallerRaw struct {
	Contract *IDaoContributorCaller // Generic read-only contract binding to access the raw methods on
}

// IDaoContributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IDaoContributorTransactorRaw struct {
	Contract *IDaoContributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIDaoContributor creates a new instance of IDaoContributor, bound to a specific deployed contract.
func NewIDaoContributor(address common.Address, backend bind.ContractBackend) (*IDaoContributor, error) {
	contract, err := bindIDaoContributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IDaoContributor{IDaoContributorCaller: IDaoContributorCaller{contract: contract}, IDaoContributorTransactor: IDaoContributorTransactor{contract: contract}, IDaoContributorFilterer: IDaoContributorFilterer{contract: contract}}, nil
}

// NewIDaoContributorCaller creates a new read-only instance of IDaoContributor, bound to a specific deployed contract.
func NewIDaoContributorCaller(address common.Address, caller bind.ContractCaller) (*IDaoContributorCaller, error) {
	contract, err := bindIDaoContributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IDaoContributorCaller{contract: contract}, nil
}

// NewIDaoContributorTransactor creates a new write-only instance of IDaoContributor, bound to a specific deployed contract.
func NewIDaoContributorTransactor(address common.Address, transactor bind.ContractTransactor) (*IDaoContributorTransactor, error) {
	contract, err := bindIDaoContributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IDaoContributorTransactor{contract: contract}, nil
}

// NewIDaoContributorFilterer creates a new log filterer instance of IDaoContributor, bound to a specific deployed contract.
func NewIDaoContributorFilterer(address common.Address, filterer bind.ContractFilterer) (*IDaoContributorFilterer, error) {
	contract, err := bindIDaoContributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IDaoContributorFilterer{contract: contract}, nil
}

// bindIDaoContributor binds a generic wrapper to an already deployed contract.
func bindIDaoContributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IDaoContributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDaoContributor *IDaoContributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDaoContributor.Contract.IDaoContributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDaoContributor *IDaoContributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDaoContributor.Contract.IDaoContributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDaoContributor *IDaoContributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDaoContributor.Contract.IDaoContributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDaoContributor *IDaoContributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDaoContributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDaoContributor *IDaoContributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDaoContributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDaoContributor *IDaoContributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDaoContributor.Contract.contract.Transact(opts, method, params...)
}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IDaoContributor *IDaoContributorCaller) GetCurrentContribution(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDaoContributor.contract.Call(opts, &out, "getCurrentContribution")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IDaoContributor *IDaoContributorSession) GetCurrentContribution() (*big.Int, error) {
	return _IDaoContributor.Contract.GetCurrentContribution(&_IDaoContributor.CallOpts)
}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IDaoContributor *IDaoContributorCallerSession) GetCurrentContribution() (*big.Int, error) {
	return _IDaoContributor.Contract.GetCurrentContribution(&_IDaoContributor.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IDaoContributor *IDaoContributorCaller) GetFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDaoContributor.contract.Call(opts, &out, "getFeeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IDaoContributor *IDaoContributorSession) GetFeeCollector() (common.Address, error) {
	return _IDaoContributor.Contract.GetFeeCollector(&_IDaoContributor.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IDaoContributor *IDaoContributorCallerSession) GetFeeCollector() (common.Address, error) {
	return _IDaoContributor.Contract.GetFeeCollector(&_IDaoContributor.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IDaoContributor *IDaoContributorCaller) GetFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDaoContributor.contract.Call(opts, &out, "getFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IDaoContributor *IDaoContributorSession) GetFeePercentage() (*big.Int, error) {
	return _IDaoContributor.Contract.GetFeePercentage(&_IDaoContributor.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IDaoContributor *IDaoContributorCallerSession) GetFeePercentage() (*big.Int, error) {
	return _IDaoContributor.Contract.GetFeePercentage(&_IDaoContributor.CallOpts)
}

// IFlyoverDiscoveryMetaData contains all meta data concerning the IFlyoverDiscovery contract.
var IFlyoverDiscoveryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"AlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InsufficientCollateral\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"InvalidProviderData\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"}],\"name\":\"InvalidProviderType\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"NotEOA\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"ProviderStatusSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"ProviderUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Register\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"}],\"name\":\"getProvider\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"internalType\":\"structFlyover.LiquidityProvider\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProviders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"providerAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"internalType\":\"structFlyover.LiquidityProvider[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isOperational\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"enumFlyover.ProviderType\",\"name\":\"providerType\",\"type\":\"uint8\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"providerId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"setProviderStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"}],\"name\":\"updateProvider\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// ILegacyLiquidityBridgeContractMetaData contains all meta data concerning the ILegacyLiquidityBridgeContract contract.
var ILegacyLiquidityBridgeContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"bytes\",\"name\":\"deposityAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"}],\"internalType\":\"structQuotesV2.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"depositPegout\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"apiBaseUrl\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"providerType\",\"type\":\"string\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ILegacyLiquidityBridgeContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ILegacyLiquidityBridgeContractMetaData.ABI instead.
var ILegacyLiquidityBridgeContractABI = ILegacyLiquidityBridgeContractMetaData.ABI

// ILegacyLiquidityBridgeContract is an auto generated Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContract struct {
	ILegacyLiquidityBridgeContractCaller     // Read-only binding to the contract
	ILegacyLiquidityBridgeContractTransactor // Write-only binding to the contract
	ILegacyLiquidityBridgeContractFilterer   // Log filterer for contract events
}

// ILegacyLiquidityBridgeContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILegacyLiquidityBridgeContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILegacyLiquidityBridgeContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ILegacyLiquidityBridgeContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ILegacyLiquidityBridgeContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ILegacyLiquidityBridgeContractSession struct {
	Contract     *ILegacyLiquidityBridgeContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ILegacyLiquidityBridgeContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ILegacyLiquidityBridgeContractCallerSession struct {
	Contract *ILegacyLiquidityBridgeContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ILegacyLiquidityBridgeContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ILegacyLiquidityBridgeContractTransactorSession struct {
	Contract     *ILegacyLiquidityBridgeContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ILegacyLiquidityBridgeContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContractRaw struct {
	Contract *ILegacyLiquidityBridgeContract // Generic contract binding to access the raw methods on
}

// ILegacyLiquidityBridgeContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContractCallerRaw struct {
	Contract *ILegacyLiquidityBridgeContractCaller // Generic read-only contract binding to access the raw methods on
}

// ILegacyLiquidityBridgeContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ILegacyLiquidityBridgeContractTransactorRaw struct {
	Contract *ILegacyLiquidityBridgeContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewILegacyLiquidityBridgeContract creates a new instance of ILegacyLiquidityBridgeContract, bound to a specific deployed contract.
func NewILegacyLiquidityBridgeContract(address common.Address, backend bind.ContractBackend) (*ILegacyLiquidityBridgeContract, error) {
	contract, err := bindILegacyLiquidityBridgeContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ILegacyLiquidityBridgeContract{ILegacyLiquidityBridgeContractCaller: ILegacyLiquidityBridgeContractCaller{contract: contract}, ILegacyLiquidityBridgeContractTransactor: ILegacyLiquidityBridgeContractTransactor{contract: contract}, ILegacyLiquidityBridgeContractFilterer: ILegacyLiquidityBridgeContractFilterer{contract: contract}}, nil
}

// NewILegacyLiquidityBridgeContractCaller creates a new read-only instance of ILegacyLiquidityBridgeContract, bound to a specific deployed contract.
func NewILegacyLiquidityBridgeContractCaller(address common.Address, caller bind.ContractCaller) (*ILegacyLiquidityBridgeContractCaller, error) {
	contract, err := bindILegacyLiquidityBridgeContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ILegacyLiquidityBridgeContractCaller{contract: contract}, nil
}

// NewILegacyLiquidityBridgeContractTransactor creates a new write-only instance of ILegacyLiquidityBridgeContract, bound to a specific deployed contract.
func NewILegacyLiquidityBridgeContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ILegacyLiquidityBridgeContractTransactor, error) {
	contract, err := bindILegacyLiquidityBridgeContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ILegacyLiquidityBridgeContractTransactor{contract: contract}, nil
}

// NewILegacyLiquidityBridgeContractFilterer creates a new log filterer instance of ILegacyLiquidityBridgeContract, bound to a specific deployed contract.
func NewILegacyLiquidityBridgeContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ILegacyLiquidityBridgeContractFilterer, error) {
	contract, err := bindILegacyLiquidityBridgeContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ILegacyLiquidityBridgeContractFilterer{contract: contract}, nil
}

// bindILegacyLiquidityBridgeContract binds a generic wrapper to an already deployed contract.
func bindILegacyLiquidityBridgeContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ILegacyLiquidityBridgeContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILegacyLiquidityBridgeContract.Contract.ILegacyLiquidityBridgeContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.ILegacyLiquidityBridgeContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.ILegacyLiquidityBridgeContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ILegacyLiquidityBridgeContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.contract.Transact(opts, method, params...)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x8beb537a.
//
// Solidity: function depositPegout((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32,uint256,uint256) quote, bytes signature) payable returns()
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactor) DepositPegout(opts *bind.TransactOpts, quote QuotesV2PegOutQuote, signature []byte) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.contract.Transact(opts, "depositPegout", quote, signature)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x8beb537a.
//
// Solidity: function depositPegout((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32,uint256,uint256) quote, bytes signature) payable returns()
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractSession) DepositPegout(quote QuotesV2PegOutQuote, signature []byte) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.DepositPegout(&_ILegacyLiquidityBridgeContract.TransactOpts, quote, signature)
}

// DepositPegout is a paid mutator transaction binding the contract method 0x8beb537a.
//
// Solidity: function depositPegout((address,address,bytes,address,bytes,uint256,uint256,int64,bytes,uint256,uint32,uint32,uint16,uint16,uint32,uint32,uint32,uint256,uint256) quote, bytes signature) payable returns()
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactorSession) DepositPegout(quote QuotesV2PegOutQuote, signature []byte) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.DepositPegout(&_ILegacyLiquidityBridgeContract.TransactOpts, quote, signature)
}

// Register is a paid mutator transaction binding the contract method 0x41705518.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, string providerType) payable returns(uint256)
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactor) Register(opts *bind.TransactOpts, name string, apiBaseUrl string, status bool, providerType string) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.contract.Transact(opts, "register", name, apiBaseUrl, status, providerType)
}

// Register is a paid mutator transaction binding the contract method 0x41705518.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, string providerType) payable returns(uint256)
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractSession) Register(name string, apiBaseUrl string, status bool, providerType string) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.Register(&_ILegacyLiquidityBridgeContract.TransactOpts, name, apiBaseUrl, status, providerType)
}

// Register is a paid mutator transaction binding the contract method 0x41705518.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, string providerType) payable returns(uint256)
func (_ILegacyLiquidityBridgeContract *ILegacyLiquidityBridgeContractTransactorSession) Register(name string, apiBaseUrl string, status bool, providerType string) (*types.Transaction, error) {
	return _ILegacyLiquidityBridgeContract.Contract.Register(&_ILegacyLiquidityBridgeContract.TransactOpts, name, apiBaseUrl, status, providerType)
}

// IPegInMetaData contains all meta data concerning the IPegIn contract.
var IPegInMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AmountUnderMinimum\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasRequired\",\"type\":\"uint256\"}],\"name\":\"InsufficientGas\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"refundAddress\",\"type\":\"bytes\"}],\"name\":\"InvalidRefundAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughConfirmations\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyProcessed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"UnexpectedBridgeError\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceDecrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BalanceIncrease\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"BridgeCapExceeded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"CallForUser\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"transferredAmount\",\"type\":\"uint256\"}],\"name\":\"PegInRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"Refund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"callForUser\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentContribution\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinPegIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"getQuoteStatus\",\"outputs\":[{\"internalType\":\"enumIPegIn.PegInStates\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegInQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRawTransaction\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"partialMerkleTree\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"registerPegIn\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"}],\"name\":\"validatePegInDepositAddress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegIn *IPegInCaller) GetCurrentContribution(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getCurrentContribution")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegIn *IPegInSession) GetCurrentContribution() (*big.Int, error) {
	return _IPegIn.Contract.GetCurrentContribution(&_IPegIn.CallOpts)
}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegIn *IPegInCallerSession) GetCurrentContribution() (*big.Int, error) {
	return _IPegIn.Contract.GetCurrentContribution(&_IPegIn.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegIn *IPegInCaller) GetFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getFeeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegIn *IPegInSession) GetFeeCollector() (common.Address, error) {
	return _IPegIn.Contract.GetFeeCollector(&_IPegIn.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegIn *IPegInCallerSession) GetFeeCollector() (common.Address, error) {
	return _IPegIn.Contract.GetFeeCollector(&_IPegIn.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegIn *IPegInCaller) GetFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegIn.contract.Call(opts, &out, "getFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegIn *IPegInSession) GetFeePercentage() (*big.Int, error) {
	return _IPegIn.Contract.GetFeePercentage(&_IPegIn.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegIn *IPegInCallerSession) GetFeePercentage() (*big.Int, error) {
	return _IPegIn.Contract.GetFeePercentage(&_IPegIn.CallOpts)
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
	ABI: "[{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"heightOrHash\",\"type\":\"bytes32\"}],\"name\":\"EmptyBlockHeader\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeCollectorUnset\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"IncorrectContract\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expectedAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"usedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"IncorrectSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"InsufficientAmount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"expected\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"actual\",\"type\":\"bytes\"}],\"name\":\"InvalidDestination\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"InvalidQuoteHash\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"actual\",\"type\":\"address\"}],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"outputScript\",\"type\":\"bytes\"}],\"name\":\"MalformedTransaction\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wanted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actual\",\"type\":\"uint256\"}],\"name\":\"NoBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"NoContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoFees\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"required\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"}],\"name\":\"NotEnoughConfirmations\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"reason\",\"type\":\"bytes\"}],\"name\":\"PaymentFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PaymentNotAllowed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ProviderNotRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyCompleted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"}],\"name\":\"QuoteExpiredByBlocks\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"}],\"name\":\"QuoteExpiredByTime\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteNotExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"QuoteNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"errorCode\",\"type\":\"int256\"}],\"name\":\"UnableToGetConfirmations\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldTime\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newTime\",\"type\":\"uint256\"}],\"name\":\"BtcBlockTimeSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAddress\",\"type\":\"address\"}],\"name\":\"CollateralManagementSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeCollector\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"ContributionsConfigured\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contributor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DaoContribution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"claimer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DaoFeesClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldThreshold\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"DustThresholdSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"EmergencyPaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"EmergencyUnpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"change\",\"type\":\"uint256\"}],\"name\":\"PegOutChangePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PegOutDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"PegOutRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"PegOutUserRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"btcBlockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimContribution\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"feeCollector\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feePercentage\",\"type\":\"uint256\"}],\"name\":\"configureContributions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"depositPegOut\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dustThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentContribution\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeePercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"hashPegOutQuote\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"bridge\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"dustThreshold_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"collateralManagement\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"mainnet\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"btcBlockTime_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"daoFeePercentage\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"daoFeeCollector\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"isQuoteCompleted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isPaused\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"since\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"btcTx\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"btcBlockHeaderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"merkleBranchPath\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\"}],\"name\":\"refundPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"}],\"name\":\"refundUserPegOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockTime\",\"type\":\"uint256\"}],\"name\":\"setBtcBlockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"collateralManagement\",\"type\":\"address\"}],\"name\":\"setCollateralManagement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setDustThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"btcTx\",\"type\":\"bytes\"}],\"name\":\"validatePegout\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IPegOut *IPegOutCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IPegOut *IPegOutSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _IPegOut.Contract.DEFAULTADMINROLE(&_IPegOut.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_IPegOut *IPegOutCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _IPegOut.Contract.DEFAULTADMINROLE(&_IPegOut.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IPegOut *IPegOutCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IPegOut *IPegOutSession) VERSION() (string, error) {
	return _IPegOut.Contract.VERSION(&_IPegOut.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (_IPegOut *IPegOutCallerSession) VERSION() (string, error) {
	return _IPegOut.Contract.VERSION(&_IPegOut.CallOpts)
}

// BtcBlockTime is a free data retrieval call binding the contract method 0x23982937.
//
// Solidity: function btcBlockTime() view returns(uint256)
func (_IPegOut *IPegOutCaller) BtcBlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "btcBlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BtcBlockTime is a free data retrieval call binding the contract method 0x23982937.
//
// Solidity: function btcBlockTime() view returns(uint256)
func (_IPegOut *IPegOutSession) BtcBlockTime() (*big.Int, error) {
	return _IPegOut.Contract.BtcBlockTime(&_IPegOut.CallOpts)
}

// BtcBlockTime is a free data retrieval call binding the contract method 0x23982937.
//
// Solidity: function btcBlockTime() view returns(uint256)
func (_IPegOut *IPegOutCallerSession) BtcBlockTime() (*big.Int, error) {
	return _IPegOut.Contract.BtcBlockTime(&_IPegOut.CallOpts)
}

// DustThreshold is a free data retrieval call binding the contract method 0xe8462e8f.
//
// Solidity: function dustThreshold() view returns(uint256)
func (_IPegOut *IPegOutCaller) DustThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "dustThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DustThreshold is a free data retrieval call binding the contract method 0xe8462e8f.
//
// Solidity: function dustThreshold() view returns(uint256)
func (_IPegOut *IPegOutSession) DustThreshold() (*big.Int, error) {
	return _IPegOut.Contract.DustThreshold(&_IPegOut.CallOpts)
}

// DustThreshold is a free data retrieval call binding the contract method 0xe8462e8f.
//
// Solidity: function dustThreshold() view returns(uint256)
func (_IPegOut *IPegOutCallerSession) DustThreshold() (*big.Int, error) {
	return _IPegOut.Contract.DustThreshold(&_IPegOut.CallOpts)
}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegOut *IPegOutCaller) GetCurrentContribution(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "getCurrentContribution")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegOut *IPegOutSession) GetCurrentContribution() (*big.Int, error) {
	return _IPegOut.Contract.GetCurrentContribution(&_IPegOut.CallOpts)
}

// GetCurrentContribution is a free data retrieval call binding the contract method 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (_IPegOut *IPegOutCallerSession) GetCurrentContribution() (*big.Int, error) {
	return _IPegOut.Contract.GetCurrentContribution(&_IPegOut.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegOut *IPegOutCaller) GetFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "getFeeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegOut *IPegOutSession) GetFeeCollector() (common.Address, error) {
	return _IPegOut.Contract.GetFeeCollector(&_IPegOut.CallOpts)
}

// GetFeeCollector is a free data retrieval call binding the contract method 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (_IPegOut *IPegOutCallerSession) GetFeeCollector() (common.Address, error) {
	return _IPegOut.Contract.GetFeeCollector(&_IPegOut.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegOut *IPegOutCaller) GetFeePercentage(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "getFeePercentage")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegOut *IPegOutSession) GetFeePercentage() (*big.Int, error) {
	return _IPegOut.Contract.GetFeePercentage(&_IPegOut.CallOpts)
}

// GetFeePercentage is a free data retrieval call binding the contract method 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (_IPegOut *IPegOutCallerSession) GetFeePercentage() (*big.Int, error) {
	return _IPegOut.Contract.GetFeePercentage(&_IPegOut.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IPegOut *IPegOutCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IPegOut *IPegOutSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _IPegOut.Contract.GetRoleAdmin(&_IPegOut.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_IPegOut *IPegOutCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _IPegOut.Contract.GetRoleAdmin(&_IPegOut.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IPegOut *IPegOutCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IPegOut *IPegOutSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _IPegOut.Contract.HasRole(&_IPegOut.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_IPegOut *IPegOutCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _IPegOut.Contract.HasRole(&_IPegOut.CallOpts, role, account)
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPegOut *IPegOutCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPegOut *IPegOutSession) Owner() (common.Address, error) {
	return _IPegOut.Contract.Owner(&_IPegOut.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IPegOut *IPegOutCallerSession) Owner() (common.Address, error) {
	return _IPegOut.Contract.Owner(&_IPegOut.CallOpts)
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

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IPegOut *IPegOutCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IPegOut *IPegOutSession) Paused() (bool, error) {
	return _IPegOut.Contract.Paused(&_IPegOut.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_IPegOut *IPegOutCallerSession) Paused() (bool, error) {
	return _IPegOut.Contract.Paused(&_IPegOut.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IPegOut *IPegOutCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _IPegOut.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IPegOut *IPegOutSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IPegOut.Contract.SupportsInterface(&_IPegOut.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_IPegOut *IPegOutCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _IPegOut.Contract.SupportsInterface(&_IPegOut.CallOpts, interfaceId)
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

// ClaimContribution is a paid mutator transaction binding the contract method 0x0114a690.
//
// Solidity: function claimContribution() returns()
func (_IPegOut *IPegOutTransactor) ClaimContribution(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "claimContribution")
}

// ClaimContribution is a paid mutator transaction binding the contract method 0x0114a690.
//
// Solidity: function claimContribution() returns()
func (_IPegOut *IPegOutSession) ClaimContribution() (*types.Transaction, error) {
	return _IPegOut.Contract.ClaimContribution(&_IPegOut.TransactOpts)
}

// ClaimContribution is a paid mutator transaction binding the contract method 0x0114a690.
//
// Solidity: function claimContribution() returns()
func (_IPegOut *IPegOutTransactorSession) ClaimContribution() (*types.Transaction, error) {
	return _IPegOut.Contract.ClaimContribution(&_IPegOut.TransactOpts)
}

// ConfigureContributions is a paid mutator transaction binding the contract method 0x10c2f1c5.
//
// Solidity: function configureContributions(address feeCollector, uint256 feePercentage) returns()
func (_IPegOut *IPegOutTransactor) ConfigureContributions(opts *bind.TransactOpts, feeCollector common.Address, feePercentage *big.Int) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "configureContributions", feeCollector, feePercentage)
}

// ConfigureContributions is a paid mutator transaction binding the contract method 0x10c2f1c5.
//
// Solidity: function configureContributions(address feeCollector, uint256 feePercentage) returns()
func (_IPegOut *IPegOutSession) ConfigureContributions(feeCollector common.Address, feePercentage *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.ConfigureContributions(&_IPegOut.TransactOpts, feeCollector, feePercentage)
}

// ConfigureContributions is a paid mutator transaction binding the contract method 0x10c2f1c5.
//
// Solidity: function configureContributions(address feeCollector, uint256 feePercentage) returns()
func (_IPegOut *IPegOutTransactorSession) ConfigureContributions(feeCollector common.Address, feePercentage *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.ConfigureContributions(&_IPegOut.TransactOpts, feeCollector, feePercentage)
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

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.GrantRole(&_IPegOut.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.GrantRole(&_IPegOut.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xd4b31024.
//
// Solidity: function initialize(address owner, address bridge, uint256 dustThreshold_, address collateralManagement, bool mainnet, uint256 btcBlockTime_, uint256 daoFeePercentage, address daoFeeCollector) returns()
func (_IPegOut *IPegOutTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, bridge common.Address, dustThreshold_ *big.Int, collateralManagement common.Address, mainnet bool, btcBlockTime_ *big.Int, daoFeePercentage *big.Int, daoFeeCollector common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "initialize", owner, bridge, dustThreshold_, collateralManagement, mainnet, btcBlockTime_, daoFeePercentage, daoFeeCollector)
}

// Initialize is a paid mutator transaction binding the contract method 0xd4b31024.
//
// Solidity: function initialize(address owner, address bridge, uint256 dustThreshold_, address collateralManagement, bool mainnet, uint256 btcBlockTime_, uint256 daoFeePercentage, address daoFeeCollector) returns()
func (_IPegOut *IPegOutSession) Initialize(owner common.Address, bridge common.Address, dustThreshold_ *big.Int, collateralManagement common.Address, mainnet bool, btcBlockTime_ *big.Int, daoFeePercentage *big.Int, daoFeeCollector common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.Initialize(&_IPegOut.TransactOpts, owner, bridge, dustThreshold_, collateralManagement, mainnet, btcBlockTime_, daoFeePercentage, daoFeeCollector)
}

// Initialize is a paid mutator transaction binding the contract method 0xd4b31024.
//
// Solidity: function initialize(address owner, address bridge, uint256 dustThreshold_, address collateralManagement, bool mainnet, uint256 btcBlockTime_, uint256 daoFeePercentage, address daoFeeCollector) returns()
func (_IPegOut *IPegOutTransactorSession) Initialize(owner common.Address, bridge common.Address, dustThreshold_ *big.Int, collateralManagement common.Address, mainnet bool, btcBlockTime_ *big.Int, daoFeePercentage *big.Int, daoFeeCollector common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.Initialize(&_IPegOut.TransactOpts, owner, bridge, dustThreshold_, collateralManagement, mainnet, btcBlockTime_, daoFeePercentage, daoFeeCollector)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string reason) returns()
func (_IPegOut *IPegOutTransactor) Pause(opts *bind.TransactOpts, reason string) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "pause", reason)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string reason) returns()
func (_IPegOut *IPegOutSession) Pause(reason string) (*types.Transaction, error) {
	return _IPegOut.Contract.Pause(&_IPegOut.TransactOpts, reason)
}

// Pause is a paid mutator transaction binding the contract method 0x6da66355.
//
// Solidity: function pause(string reason) returns()
func (_IPegOut *IPegOutTransactorSession) Pause(reason string) (*types.Transaction, error) {
	return _IPegOut.Contract.Pause(&_IPegOut.TransactOpts, reason)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPegOut *IPegOutTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPegOut *IPegOutSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPegOut.Contract.RenounceOwnership(&_IPegOut.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IPegOut *IPegOutTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IPegOut.Contract.RenounceOwnership(&_IPegOut.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IPegOut *IPegOutTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IPegOut *IPegOutSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.RenounceRole(&_IPegOut.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_IPegOut *IPegOutTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.RenounceRole(&_IPegOut.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.RevokeRole(&_IPegOut.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_IPegOut *IPegOutTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.RevokeRole(&_IPegOut.TransactOpts, role, account)
}

// SetBtcBlockTime is a paid mutator transaction binding the contract method 0x3492426f.
//
// Solidity: function setBtcBlockTime(uint256 blockTime) returns()
func (_IPegOut *IPegOutTransactor) SetBtcBlockTime(opts *bind.TransactOpts, blockTime *big.Int) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "setBtcBlockTime", blockTime)
}

// SetBtcBlockTime is a paid mutator transaction binding the contract method 0x3492426f.
//
// Solidity: function setBtcBlockTime(uint256 blockTime) returns()
func (_IPegOut *IPegOutSession) SetBtcBlockTime(blockTime *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.SetBtcBlockTime(&_IPegOut.TransactOpts, blockTime)
}

// SetBtcBlockTime is a paid mutator transaction binding the contract method 0x3492426f.
//
// Solidity: function setBtcBlockTime(uint256 blockTime) returns()
func (_IPegOut *IPegOutTransactorSession) SetBtcBlockTime(blockTime *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.SetBtcBlockTime(&_IPegOut.TransactOpts, blockTime)
}

// SetCollateralManagement is a paid mutator transaction binding the contract method 0x05dfbf62.
//
// Solidity: function setCollateralManagement(address collateralManagement) returns()
func (_IPegOut *IPegOutTransactor) SetCollateralManagement(opts *bind.TransactOpts, collateralManagement common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "setCollateralManagement", collateralManagement)
}

// SetCollateralManagement is a paid mutator transaction binding the contract method 0x05dfbf62.
//
// Solidity: function setCollateralManagement(address collateralManagement) returns()
func (_IPegOut *IPegOutSession) SetCollateralManagement(collateralManagement common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.SetCollateralManagement(&_IPegOut.TransactOpts, collateralManagement)
}

// SetCollateralManagement is a paid mutator transaction binding the contract method 0x05dfbf62.
//
// Solidity: function setCollateralManagement(address collateralManagement) returns()
func (_IPegOut *IPegOutTransactorSession) SetCollateralManagement(collateralManagement common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.SetCollateralManagement(&_IPegOut.TransactOpts, collateralManagement)
}

// SetDustThreshold is a paid mutator transaction binding the contract method 0xad7e55ba.
//
// Solidity: function setDustThreshold(uint256 threshold) returns()
func (_IPegOut *IPegOutTransactor) SetDustThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "setDustThreshold", threshold)
}

// SetDustThreshold is a paid mutator transaction binding the contract method 0xad7e55ba.
//
// Solidity: function setDustThreshold(uint256 threshold) returns()
func (_IPegOut *IPegOutSession) SetDustThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.SetDustThreshold(&_IPegOut.TransactOpts, threshold)
}

// SetDustThreshold is a paid mutator transaction binding the contract method 0xad7e55ba.
//
// Solidity: function setDustThreshold(uint256 threshold) returns()
func (_IPegOut *IPegOutTransactorSession) SetDustThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _IPegOut.Contract.SetDustThreshold(&_IPegOut.TransactOpts, threshold)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPegOut *IPegOutTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPegOut *IPegOutSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.TransferOwnership(&_IPegOut.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IPegOut *IPegOutTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IPegOut.Contract.TransferOwnership(&_IPegOut.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IPegOut *IPegOutTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IPegOut *IPegOutSession) Unpause() (*types.Transaction, error) {
	return _IPegOut.Contract.Unpause(&_IPegOut.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_IPegOut *IPegOutTransactorSession) Unpause() (*types.Transaction, error) {
	return _IPegOut.Contract.Unpause(&_IPegOut.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPegOut *IPegOutTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPegOut.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPegOut *IPegOutSession) Receive() (*types.Transaction, error) {
	return _IPegOut.Contract.Receive(&_IPegOut.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPegOut *IPegOutTransactorSession) Receive() (*types.Transaction, error) {
	return _IPegOut.Contract.Receive(&_IPegOut.TransactOpts)
}

// IPegOutBtcBlockTimeSetIterator is returned from FilterBtcBlockTimeSet and is used to iterate over the raw logs and unpacked data for BtcBlockTimeSet events raised by the IPegOut contract.
type IPegOutBtcBlockTimeSetIterator struct {
	Event *IPegOutBtcBlockTimeSet // Event containing the contract specifics and raw log

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
func (it *IPegOutBtcBlockTimeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutBtcBlockTimeSet)
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
		it.Event = new(IPegOutBtcBlockTimeSet)
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
func (it *IPegOutBtcBlockTimeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutBtcBlockTimeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutBtcBlockTimeSet represents a BtcBlockTimeSet event raised by the IPegOut contract.
type IPegOutBtcBlockTimeSet struct {
	OldTime *big.Int
	NewTime *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBtcBlockTimeSet is a free log retrieval operation binding the contract event 0x644c2a75d9785c3196efad342f8ea14846206549e770ae1eb99d24675a1bb8e0.
//
// Solidity: event BtcBlockTimeSet(uint256 indexed oldTime, uint256 indexed newTime)
func (_IPegOut *IPegOutFilterer) FilterBtcBlockTimeSet(opts *bind.FilterOpts, oldTime []*big.Int, newTime []*big.Int) (*IPegOutBtcBlockTimeSetIterator, error) {

	var oldTimeRule []interface{}
	for _, oldTimeItem := range oldTime {
		oldTimeRule = append(oldTimeRule, oldTimeItem)
	}
	var newTimeRule []interface{}
	for _, newTimeItem := range newTime {
		newTimeRule = append(newTimeRule, newTimeItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "BtcBlockTimeSet", oldTimeRule, newTimeRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutBtcBlockTimeSetIterator{contract: _IPegOut.contract, event: "BtcBlockTimeSet", logs: logs, sub: sub}, nil
}

// WatchBtcBlockTimeSet is a free log subscription operation binding the contract event 0x644c2a75d9785c3196efad342f8ea14846206549e770ae1eb99d24675a1bb8e0.
//
// Solidity: event BtcBlockTimeSet(uint256 indexed oldTime, uint256 indexed newTime)
func (_IPegOut *IPegOutFilterer) WatchBtcBlockTimeSet(opts *bind.WatchOpts, sink chan<- *IPegOutBtcBlockTimeSet, oldTime []*big.Int, newTime []*big.Int) (event.Subscription, error) {

	var oldTimeRule []interface{}
	for _, oldTimeItem := range oldTime {
		oldTimeRule = append(oldTimeRule, oldTimeItem)
	}
	var newTimeRule []interface{}
	for _, newTimeItem := range newTime {
		newTimeRule = append(newTimeRule, newTimeItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "BtcBlockTimeSet", oldTimeRule, newTimeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutBtcBlockTimeSet)
				if err := _IPegOut.contract.UnpackLog(event, "BtcBlockTimeSet", log); err != nil {
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

// ParseBtcBlockTimeSet is a log parse operation binding the contract event 0x644c2a75d9785c3196efad342f8ea14846206549e770ae1eb99d24675a1bb8e0.
//
// Solidity: event BtcBlockTimeSet(uint256 indexed oldTime, uint256 indexed newTime)
func (_IPegOut *IPegOutFilterer) ParseBtcBlockTimeSet(log types.Log) (*IPegOutBtcBlockTimeSet, error) {
	event := new(IPegOutBtcBlockTimeSet)
	if err := _IPegOut.contract.UnpackLog(event, "BtcBlockTimeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutCollateralManagementSetIterator is returned from FilterCollateralManagementSet and is used to iterate over the raw logs and unpacked data for CollateralManagementSet events raised by the IPegOut contract.
type IPegOutCollateralManagementSetIterator struct {
	Event *IPegOutCollateralManagementSet // Event containing the contract specifics and raw log

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
func (it *IPegOutCollateralManagementSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutCollateralManagementSet)
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
		it.Event = new(IPegOutCollateralManagementSet)
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
func (it *IPegOutCollateralManagementSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutCollateralManagementSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutCollateralManagementSet represents a CollateralManagementSet event raised by the IPegOut contract.
type IPegOutCollateralManagementSet struct {
	OldAddress common.Address
	NewAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCollateralManagementSet is a free log retrieval operation binding the contract event 0x95aab800aba1612fbcd4ca8b23b6fb7e7a3b24792b0132fcddfbe48891ed13d0.
//
// Solidity: event CollateralManagementSet(address indexed oldAddress, address indexed newAddress)
func (_IPegOut *IPegOutFilterer) FilterCollateralManagementSet(opts *bind.FilterOpts, oldAddress []common.Address, newAddress []common.Address) (*IPegOutCollateralManagementSetIterator, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "CollateralManagementSet", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutCollateralManagementSetIterator{contract: _IPegOut.contract, event: "CollateralManagementSet", logs: logs, sub: sub}, nil
}

// WatchCollateralManagementSet is a free log subscription operation binding the contract event 0x95aab800aba1612fbcd4ca8b23b6fb7e7a3b24792b0132fcddfbe48891ed13d0.
//
// Solidity: event CollateralManagementSet(address indexed oldAddress, address indexed newAddress)
func (_IPegOut *IPegOutFilterer) WatchCollateralManagementSet(opts *bind.WatchOpts, sink chan<- *IPegOutCollateralManagementSet, oldAddress []common.Address, newAddress []common.Address) (event.Subscription, error) {

	var oldAddressRule []interface{}
	for _, oldAddressItem := range oldAddress {
		oldAddressRule = append(oldAddressRule, oldAddressItem)
	}
	var newAddressRule []interface{}
	for _, newAddressItem := range newAddress {
		newAddressRule = append(newAddressRule, newAddressItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "CollateralManagementSet", oldAddressRule, newAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutCollateralManagementSet)
				if err := _IPegOut.contract.UnpackLog(event, "CollateralManagementSet", log); err != nil {
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

// ParseCollateralManagementSet is a log parse operation binding the contract event 0x95aab800aba1612fbcd4ca8b23b6fb7e7a3b24792b0132fcddfbe48891ed13d0.
//
// Solidity: event CollateralManagementSet(address indexed oldAddress, address indexed newAddress)
func (_IPegOut *IPegOutFilterer) ParseCollateralManagementSet(log types.Log) (*IPegOutCollateralManagementSet, error) {
	event := new(IPegOutCollateralManagementSet)
	if err := _IPegOut.contract.UnpackLog(event, "CollateralManagementSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutContributionsConfiguredIterator is returned from FilterContributionsConfigured and is used to iterate over the raw logs and unpacked data for ContributionsConfigured events raised by the IPegOut contract.
type IPegOutContributionsConfiguredIterator struct {
	Event *IPegOutContributionsConfigured // Event containing the contract specifics and raw log

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
func (it *IPegOutContributionsConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutContributionsConfigured)
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
		it.Event = new(IPegOutContributionsConfigured)
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
func (it *IPegOutContributionsConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutContributionsConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutContributionsConfigured represents a ContributionsConfigured event raised by the IPegOut contract.
type IPegOutContributionsConfigured struct {
	FeeCollector  common.Address
	FeePercentage *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterContributionsConfigured is a free log retrieval operation binding the contract event 0xf86fbded534f84710d18920a93885d142d4b5e3687e2d3e4a58b8a17c70f8a4c.
//
// Solidity: event ContributionsConfigured(address indexed feeCollector, uint256 indexed feePercentage)
func (_IPegOut *IPegOutFilterer) FilterContributionsConfigured(opts *bind.FilterOpts, feeCollector []common.Address, feePercentage []*big.Int) (*IPegOutContributionsConfiguredIterator, error) {

	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}
	var feePercentageRule []interface{}
	for _, feePercentageItem := range feePercentage {
		feePercentageRule = append(feePercentageRule, feePercentageItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "ContributionsConfigured", feeCollectorRule, feePercentageRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutContributionsConfiguredIterator{contract: _IPegOut.contract, event: "ContributionsConfigured", logs: logs, sub: sub}, nil
}

// WatchContributionsConfigured is a free log subscription operation binding the contract event 0xf86fbded534f84710d18920a93885d142d4b5e3687e2d3e4a58b8a17c70f8a4c.
//
// Solidity: event ContributionsConfigured(address indexed feeCollector, uint256 indexed feePercentage)
func (_IPegOut *IPegOutFilterer) WatchContributionsConfigured(opts *bind.WatchOpts, sink chan<- *IPegOutContributionsConfigured, feeCollector []common.Address, feePercentage []*big.Int) (event.Subscription, error) {

	var feeCollectorRule []interface{}
	for _, feeCollectorItem := range feeCollector {
		feeCollectorRule = append(feeCollectorRule, feeCollectorItem)
	}
	var feePercentageRule []interface{}
	for _, feePercentageItem := range feePercentage {
		feePercentageRule = append(feePercentageRule, feePercentageItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "ContributionsConfigured", feeCollectorRule, feePercentageRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutContributionsConfigured)
				if err := _IPegOut.contract.UnpackLog(event, "ContributionsConfigured", log); err != nil {
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

// ParseContributionsConfigured is a log parse operation binding the contract event 0xf86fbded534f84710d18920a93885d142d4b5e3687e2d3e4a58b8a17c70f8a4c.
//
// Solidity: event ContributionsConfigured(address indexed feeCollector, uint256 indexed feePercentage)
func (_IPegOut *IPegOutFilterer) ParseContributionsConfigured(log types.Log) (*IPegOutContributionsConfigured, error) {
	event := new(IPegOutContributionsConfigured)
	if err := _IPegOut.contract.UnpackLog(event, "ContributionsConfigured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutDaoContributionIterator is returned from FilterDaoContribution and is used to iterate over the raw logs and unpacked data for DaoContribution events raised by the IPegOut contract.
type IPegOutDaoContributionIterator struct {
	Event *IPegOutDaoContribution // Event containing the contract specifics and raw log

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
func (it *IPegOutDaoContributionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutDaoContribution)
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
		it.Event = new(IPegOutDaoContribution)
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
func (it *IPegOutDaoContributionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutDaoContributionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutDaoContribution represents a DaoContribution event raised by the IPegOut contract.
type IPegOutDaoContribution struct {
	Contributor common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDaoContribution is a free log retrieval operation binding the contract event 0x0bea2810a9b120f7c61cad402d6fb7224d38f3b20102000398102d9e68f32503.
//
// Solidity: event DaoContribution(address indexed contributor, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) FilterDaoContribution(opts *bind.FilterOpts, contributor []common.Address, amount []*big.Int) (*IPegOutDaoContributionIterator, error) {

	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "DaoContribution", contributorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutDaoContributionIterator{contract: _IPegOut.contract, event: "DaoContribution", logs: logs, sub: sub}, nil
}

// WatchDaoContribution is a free log subscription operation binding the contract event 0x0bea2810a9b120f7c61cad402d6fb7224d38f3b20102000398102d9e68f32503.
//
// Solidity: event DaoContribution(address indexed contributor, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) WatchDaoContribution(opts *bind.WatchOpts, sink chan<- *IPegOutDaoContribution, contributor []common.Address, amount []*big.Int) (event.Subscription, error) {

	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "DaoContribution", contributorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutDaoContribution)
				if err := _IPegOut.contract.UnpackLog(event, "DaoContribution", log); err != nil {
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

// ParseDaoContribution is a log parse operation binding the contract event 0x0bea2810a9b120f7c61cad402d6fb7224d38f3b20102000398102d9e68f32503.
//
// Solidity: event DaoContribution(address indexed contributor, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) ParseDaoContribution(log types.Log) (*IPegOutDaoContribution, error) {
	event := new(IPegOutDaoContribution)
	if err := _IPegOut.contract.UnpackLog(event, "DaoContribution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutDaoFeesClaimedIterator is returned from FilterDaoFeesClaimed and is used to iterate over the raw logs and unpacked data for DaoFeesClaimed events raised by the IPegOut contract.
type IPegOutDaoFeesClaimedIterator struct {
	Event *IPegOutDaoFeesClaimed // Event containing the contract specifics and raw log

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
func (it *IPegOutDaoFeesClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutDaoFeesClaimed)
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
		it.Event = new(IPegOutDaoFeesClaimed)
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
func (it *IPegOutDaoFeesClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutDaoFeesClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutDaoFeesClaimed represents a DaoFeesClaimed event raised by the IPegOut contract.
type IPegOutDaoFeesClaimed struct {
	Claimer  common.Address
	Receiver common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDaoFeesClaimed is a free log retrieval operation binding the contract event 0x6efd2761705c3e459833b29e5f0b7c2e9ff5ed7e85e2c03b57a1963e6ca2f2ff.
//
// Solidity: event DaoFeesClaimed(address indexed claimer, address indexed receiver, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) FilterDaoFeesClaimed(opts *bind.FilterOpts, claimer []common.Address, receiver []common.Address, amount []*big.Int) (*IPegOutDaoFeesClaimedIterator, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "DaoFeesClaimed", claimerRule, receiverRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutDaoFeesClaimedIterator{contract: _IPegOut.contract, event: "DaoFeesClaimed", logs: logs, sub: sub}, nil
}

// WatchDaoFeesClaimed is a free log subscription operation binding the contract event 0x6efd2761705c3e459833b29e5f0b7c2e9ff5ed7e85e2c03b57a1963e6ca2f2ff.
//
// Solidity: event DaoFeesClaimed(address indexed claimer, address indexed receiver, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) WatchDaoFeesClaimed(opts *bind.WatchOpts, sink chan<- *IPegOutDaoFeesClaimed, claimer []common.Address, receiver []common.Address, amount []*big.Int) (event.Subscription, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "DaoFeesClaimed", claimerRule, receiverRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutDaoFeesClaimed)
				if err := _IPegOut.contract.UnpackLog(event, "DaoFeesClaimed", log); err != nil {
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

// ParseDaoFeesClaimed is a log parse operation binding the contract event 0x6efd2761705c3e459833b29e5f0b7c2e9ff5ed7e85e2c03b57a1963e6ca2f2ff.
//
// Solidity: event DaoFeesClaimed(address indexed claimer, address indexed receiver, uint256 indexed amount)
func (_IPegOut *IPegOutFilterer) ParseDaoFeesClaimed(log types.Log) (*IPegOutDaoFeesClaimed, error) {
	event := new(IPegOutDaoFeesClaimed)
	if err := _IPegOut.contract.UnpackLog(event, "DaoFeesClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutDustThresholdSetIterator is returned from FilterDustThresholdSet and is used to iterate over the raw logs and unpacked data for DustThresholdSet events raised by the IPegOut contract.
type IPegOutDustThresholdSetIterator struct {
	Event *IPegOutDustThresholdSet // Event containing the contract specifics and raw log

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
func (it *IPegOutDustThresholdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutDustThresholdSet)
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
		it.Event = new(IPegOutDustThresholdSet)
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
func (it *IPegOutDustThresholdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutDustThresholdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutDustThresholdSet represents a DustThresholdSet event raised by the IPegOut contract.
type IPegOutDustThresholdSet struct {
	OldThreshold *big.Int
	NewThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDustThresholdSet is a free log retrieval operation binding the contract event 0x7c3840c706b556b72613551196289967782829a05b71bfd3f030ec84a97d060a.
//
// Solidity: event DustThresholdSet(uint256 indexed oldThreshold, uint256 indexed newThreshold)
func (_IPegOut *IPegOutFilterer) FilterDustThresholdSet(opts *bind.FilterOpts, oldThreshold []*big.Int, newThreshold []*big.Int) (*IPegOutDustThresholdSetIterator, error) {

	var oldThresholdRule []interface{}
	for _, oldThresholdItem := range oldThreshold {
		oldThresholdRule = append(oldThresholdRule, oldThresholdItem)
	}
	var newThresholdRule []interface{}
	for _, newThresholdItem := range newThreshold {
		newThresholdRule = append(newThresholdRule, newThresholdItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "DustThresholdSet", oldThresholdRule, newThresholdRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutDustThresholdSetIterator{contract: _IPegOut.contract, event: "DustThresholdSet", logs: logs, sub: sub}, nil
}

// WatchDustThresholdSet is a free log subscription operation binding the contract event 0x7c3840c706b556b72613551196289967782829a05b71bfd3f030ec84a97d060a.
//
// Solidity: event DustThresholdSet(uint256 indexed oldThreshold, uint256 indexed newThreshold)
func (_IPegOut *IPegOutFilterer) WatchDustThresholdSet(opts *bind.WatchOpts, sink chan<- *IPegOutDustThresholdSet, oldThreshold []*big.Int, newThreshold []*big.Int) (event.Subscription, error) {

	var oldThresholdRule []interface{}
	for _, oldThresholdItem := range oldThreshold {
		oldThresholdRule = append(oldThresholdRule, oldThresholdItem)
	}
	var newThresholdRule []interface{}
	for _, newThresholdItem := range newThreshold {
		newThresholdRule = append(newThresholdRule, newThresholdItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "DustThresholdSet", oldThresholdRule, newThresholdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutDustThresholdSet)
				if err := _IPegOut.contract.UnpackLog(event, "DustThresholdSet", log); err != nil {
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

// ParseDustThresholdSet is a log parse operation binding the contract event 0x7c3840c706b556b72613551196289967782829a05b71bfd3f030ec84a97d060a.
//
// Solidity: event DustThresholdSet(uint256 indexed oldThreshold, uint256 indexed newThreshold)
func (_IPegOut *IPegOutFilterer) ParseDustThresholdSet(log types.Log) (*IPegOutDustThresholdSet, error) {
	event := new(IPegOutDustThresholdSet)
	if err := _IPegOut.contract.UnpackLog(event, "DustThresholdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutEmergencyPausedIterator is returned from FilterEmergencyPaused and is used to iterate over the raw logs and unpacked data for EmergencyPaused events raised by the IPegOut contract.
type IPegOutEmergencyPausedIterator struct {
	Event *IPegOutEmergencyPaused // Event containing the contract specifics and raw log

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
func (it *IPegOutEmergencyPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutEmergencyPaused)
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
		it.Event = new(IPegOutEmergencyPaused)
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
func (it *IPegOutEmergencyPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutEmergencyPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutEmergencyPaused represents a EmergencyPaused event raised by the IPegOut contract.
type IPegOutEmergencyPaused struct {
	By     common.Address
	Reason string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEmergencyPaused is a free log retrieval operation binding the contract event 0x375c0abd968f4602b557f6ac9a48ffc89820233aa9becc5d7ff1176fd09eafff.
//
// Solidity: event EmergencyPaused(address indexed by, string reason)
func (_IPegOut *IPegOutFilterer) FilterEmergencyPaused(opts *bind.FilterOpts, by []common.Address) (*IPegOutEmergencyPausedIterator, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "EmergencyPaused", byRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutEmergencyPausedIterator{contract: _IPegOut.contract, event: "EmergencyPaused", logs: logs, sub: sub}, nil
}

// WatchEmergencyPaused is a free log subscription operation binding the contract event 0x375c0abd968f4602b557f6ac9a48ffc89820233aa9becc5d7ff1176fd09eafff.
//
// Solidity: event EmergencyPaused(address indexed by, string reason)
func (_IPegOut *IPegOutFilterer) WatchEmergencyPaused(opts *bind.WatchOpts, sink chan<- *IPegOutEmergencyPaused, by []common.Address) (event.Subscription, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "EmergencyPaused", byRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutEmergencyPaused)
				if err := _IPegOut.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
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

// ParseEmergencyPaused is a log parse operation binding the contract event 0x375c0abd968f4602b557f6ac9a48ffc89820233aa9becc5d7ff1176fd09eafff.
//
// Solidity: event EmergencyPaused(address indexed by, string reason)
func (_IPegOut *IPegOutFilterer) ParseEmergencyPaused(log types.Log) (*IPegOutEmergencyPaused, error) {
	event := new(IPegOutEmergencyPaused)
	if err := _IPegOut.contract.UnpackLog(event, "EmergencyPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutEmergencyUnpausedIterator is returned from FilterEmergencyUnpaused and is used to iterate over the raw logs and unpacked data for EmergencyUnpaused events raised by the IPegOut contract.
type IPegOutEmergencyUnpausedIterator struct {
	Event *IPegOutEmergencyUnpaused // Event containing the contract specifics and raw log

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
func (it *IPegOutEmergencyUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutEmergencyUnpaused)
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
		it.Event = new(IPegOutEmergencyUnpaused)
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
func (it *IPegOutEmergencyUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutEmergencyUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutEmergencyUnpaused represents a EmergencyUnpaused event raised by the IPegOut contract.
type IPegOutEmergencyUnpaused struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyUnpaused is a free log retrieval operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed by)
func (_IPegOut *IPegOutFilterer) FilterEmergencyUnpaused(opts *bind.FilterOpts, by []common.Address) (*IPegOutEmergencyUnpausedIterator, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "EmergencyUnpaused", byRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutEmergencyUnpausedIterator{contract: _IPegOut.contract, event: "EmergencyUnpaused", logs: logs, sub: sub}, nil
}

// WatchEmergencyUnpaused is a free log subscription operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed by)
func (_IPegOut *IPegOutFilterer) WatchEmergencyUnpaused(opts *bind.WatchOpts, sink chan<- *IPegOutEmergencyUnpaused, by []common.Address) (event.Subscription, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "EmergencyUnpaused", byRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutEmergencyUnpaused)
				if err := _IPegOut.contract.UnpackLog(event, "EmergencyUnpaused", log); err != nil {
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

// ParseEmergencyUnpaused is a log parse operation binding the contract event 0xf5cbf596165cc457b2cd92e8d8450827ee314968160a5696402d75766fc52caf.
//
// Solidity: event EmergencyUnpaused(address indexed by)
func (_IPegOut *IPegOutFilterer) ParseEmergencyUnpaused(log types.Log) (*IPegOutEmergencyUnpaused, error) {
	event := new(IPegOutEmergencyUnpaused)
	if err := _IPegOut.contract.UnpackLog(event, "EmergencyUnpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the IPegOut contract.
type IPegOutInitializedIterator struct {
	Event *IPegOutInitialized // Event containing the contract specifics and raw log

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
func (it *IPegOutInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutInitialized)
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
		it.Event = new(IPegOutInitialized)
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
func (it *IPegOutInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutInitialized represents a Initialized event raised by the IPegOut contract.
type IPegOutInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPegOut *IPegOutFilterer) FilterInitialized(opts *bind.FilterOpts) (*IPegOutInitializedIterator, error) {

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &IPegOutInitializedIterator{contract: _IPegOut.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPegOut *IPegOutFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *IPegOutInitialized) (event.Subscription, error) {

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutInitialized)
				if err := _IPegOut.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_IPegOut *IPegOutFilterer) ParseInitialized(log types.Log) (*IPegOutInitialized, error) {
	event := new(IPegOutInitialized)
	if err := _IPegOut.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IPegOut contract.
type IPegOutOwnershipTransferredIterator struct {
	Event *IPegOutOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IPegOutOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutOwnershipTransferred)
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
		it.Event = new(IPegOutOwnershipTransferred)
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
func (it *IPegOutOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutOwnershipTransferred represents a OwnershipTransferred event raised by the IPegOut contract.
type IPegOutOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPegOut *IPegOutFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IPegOutOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutOwnershipTransferredIterator{contract: _IPegOut.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPegOut *IPegOutFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IPegOutOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutOwnershipTransferred)
				if err := _IPegOut.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IPegOut *IPegOutFilterer) ParseOwnershipTransferred(log types.Log) (*IPegOutOwnershipTransferred, error) {
	event := new(IPegOutOwnershipTransferred)
	if err := _IPegOut.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the IPegOut contract.
type IPegOutPausedIterator struct {
	Event *IPegOutPaused // Event containing the contract specifics and raw log

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
func (it *IPegOutPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutPaused)
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
		it.Event = new(IPegOutPaused)
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
func (it *IPegOutPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutPaused represents a Paused event raised by the IPegOut contract.
type IPegOutPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IPegOut *IPegOutFilterer) FilterPaused(opts *bind.FilterOpts) (*IPegOutPausedIterator, error) {

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &IPegOutPausedIterator{contract: _IPegOut.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IPegOut *IPegOutFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *IPegOutPaused) (event.Subscription, error) {

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutPaused)
				if err := _IPegOut.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_IPegOut *IPegOutFilterer) ParsePaused(log types.Log) (*IPegOutPaused, error) {
	event := new(IPegOutPaused)
	if err := _IPegOut.contract.UnpackLog(event, "Paused", log); err != nil {
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

// IPegOutRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the IPegOut contract.
type IPegOutRoleAdminChangedIterator struct {
	Event *IPegOutRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *IPegOutRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutRoleAdminChanged)
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
		it.Event = new(IPegOutRoleAdminChanged)
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
func (it *IPegOutRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutRoleAdminChanged represents a RoleAdminChanged event raised by the IPegOut contract.
type IPegOutRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_IPegOut *IPegOutFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*IPegOutRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutRoleAdminChangedIterator{contract: _IPegOut.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_IPegOut *IPegOutFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *IPegOutRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutRoleAdminChanged)
				if err := _IPegOut.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_IPegOut *IPegOutFilterer) ParseRoleAdminChanged(log types.Log) (*IPegOutRoleAdminChanged, error) {
	event := new(IPegOutRoleAdminChanged)
	if err := _IPegOut.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the IPegOut contract.
type IPegOutRoleGrantedIterator struct {
	Event *IPegOutRoleGranted // Event containing the contract specifics and raw log

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
func (it *IPegOutRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutRoleGranted)
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
		it.Event = new(IPegOutRoleGranted)
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
func (it *IPegOutRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutRoleGranted represents a RoleGranted event raised by the IPegOut contract.
type IPegOutRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*IPegOutRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutRoleGrantedIterator{contract: _IPegOut.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *IPegOutRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutRoleGranted)
				if err := _IPegOut.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) ParseRoleGranted(log types.Log) (*IPegOutRoleGranted, error) {
	event := new(IPegOutRoleGranted)
	if err := _IPegOut.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the IPegOut contract.
type IPegOutRoleRevokedIterator struct {
	Event *IPegOutRoleRevoked // Event containing the contract specifics and raw log

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
func (it *IPegOutRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutRoleRevoked)
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
		it.Event = new(IPegOutRoleRevoked)
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
func (it *IPegOutRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutRoleRevoked represents a RoleRevoked event raised by the IPegOut contract.
type IPegOutRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*IPegOutRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &IPegOutRoleRevokedIterator{contract: _IPegOut.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *IPegOutRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutRoleRevoked)
				if err := _IPegOut.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_IPegOut *IPegOutFilterer) ParseRoleRevoked(log types.Log) (*IPegOutRoleRevoked, error) {
	event := new(IPegOutRoleRevoked)
	if err := _IPegOut.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPegOutUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the IPegOut contract.
type IPegOutUnpausedIterator struct {
	Event *IPegOutUnpaused // Event containing the contract specifics and raw log

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
func (it *IPegOutUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPegOutUnpaused)
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
		it.Event = new(IPegOutUnpaused)
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
func (it *IPegOutUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPegOutUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPegOutUnpaused represents a Unpaused event raised by the IPegOut contract.
type IPegOutUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IPegOut *IPegOutFilterer) FilterUnpaused(opts *bind.FilterOpts) (*IPegOutUnpausedIterator, error) {

	logs, sub, err := _IPegOut.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &IPegOutUnpausedIterator{contract: _IPegOut.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IPegOut *IPegOutFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *IPegOutUnpaused) (event.Subscription, error) {

	logs, sub, err := _IPegOut.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPegOutUnpaused)
				if err := _IPegOut.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_IPegOut *IPegOutFilterer) ParseUnpaused(log types.Log) (*IPegOutUnpaused, error) {
	event := new(IPegOutUnpaused)
	if err := _IPegOut.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// QuotesMetaData contains all meta data concerning the Quotes contract.
var QuotesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"target\",\"type\":\"uint256\"}],\"name\":\"AmountTooLow\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"transferredAmount\",\"type\":\"uint256\"}],\"name\":\"checkAgreedAmount\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lpRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"depositDateLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"transferTime\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireDate\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"expireBlock\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"transferConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bytes\",\"name\":\"depositAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"lpBtcAddress\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegOutQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"encodePegOutQuote\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"callFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"penaltyFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"productFeeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes20\",\"name\":\"fedBtcAddress\",\"type\":\"bytes20\"},{\"internalType\":\"address\",\"name\":\"lbcAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"rskRefundAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"nonce\",\"type\":\"int64\"},{\"internalType\":\"uint32\",\"name\":\"gasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"agreementTimestamp\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"timeForDeposit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"callTime\",\"type\":\"uint32\"},{\"internalType\":\"uint16\",\"name\":\"depositConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"bool\",\"name\":\"callOnRegister\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"btcRefundAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structQuotes.PegInQuote\",\"name\":\"quote\",\"type\":\"tuple\"}],\"name\":\"encodeQuote\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x610b2c610039600b82828239805160001a607314602c57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361061004b5760003560e01c80632770727a14610050578063e5648b9a14610065578063fb19b88b1461008e575b600080fd5b61006361005e36600461034a565b6100a1565b005b61007861007336600461038e565b61012a565b6040516100859190610410565b60405180910390f35b61007861009c36600461042a565b610175565b6000608083013560608401356100bc8535604087013561047b565b6100c6919061047b565b6100d0919061047b565b905060006100e061271083610494565b9050826100ed82846104b6565b111561012457826100fe82846104b6565b604051631cc6243f60e01b81526004810192909252602482015260440160405180910390fd5b50505050565b606061013d61013883610633565b610199565b61014e61014984610633565b6101e4565b60405160200161015f9291906107d6565b6040516020818303038152906040529050919050565b606061018861018383610804565b61023f565b61014e61019484610804565b61028a565b60608160a001518260c001518360e001518461022001518561012001518661024001518760000151886020015189610100015160405160200161015f9998979695949392919061097a565b60608161026001518261016001518361014001518460400151856101800151866101a00151876101c00151886101e001518961020001518a606001518b6080015160405160200161015f9b9a999897969594939291906109f9565b60608160a001518260c001518361022001518460e001518561024001518660000151876020015188610100015189610200015160405160200161015f99989796959493929190610a72565b60608160400151826101200151836101400151846101c00151856101e00151866101600151876101800151886101a0015189606001518a6080015160405160200161015f9a99989796959493929190998a5263ffffffff98891660208b015296881660408a015261ffff95861660608a015293909416608088015290851660a0870152841660c0860152921660e08401526101008301919091526101208201526101400190565b6000610280828403121561034457600080fd5b50919050565b6000806040838503121561035d57600080fd5b82356001600160401b0381111561037357600080fd5b61037f85828601610331565b95602094909401359450505050565b6000602082840312156103a057600080fd5b81356001600160401b038111156103b657600080fd5b6103c284828501610331565b949350505050565b6000815180845260005b818110156103f0576020818501810151868301820152016103d4565b506000602082860101526020601f19601f83011685010191505092915050565b60208152600061042360208301846103ca565b9392505050565b60006020828403121561043c57600080fd5b81356001600160401b0381111561045257600080fd5b8201610260818503121561042357600080fd5b634e487b7160e01b600052601160045260246000fd5b8082018082111561048e5761048e610465565b92915050565b6000826104b157634e487b7160e01b600052601260045260246000fd5b500490565b8181038181111561048e5761048e610465565b634e487b7160e01b600052604160045260246000fd5b60405161028081016001600160401b0381118282101715610502576105026104c9565b60405290565b60405161026081016001600160401b0381118282101715610502576105026104c9565b80356001600160601b03198116811461054357600080fd5b919050565b80356001600160a01b038116811461054357600080fd5b8035600781900b811461054357600080fd5b803563ffffffff8116811461054357600080fd5b803561ffff8116811461054357600080fd5b8035801515811461054357600080fd5b600082601f8301126105b857600080fd5b81356001600160401b03808211156105d2576105d26104c9565b604051601f8301601f19908116603f011681019082821181831017156105fa576105fa6104c9565b8160405283815286602085880101111561061357600080fd5b836020870160208301376000602085830101528094505050505092915050565b6000610280823603121561064657600080fd5b61064e6104df565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261068660a0840161052b565b60a082015261069760c08401610548565b60c08201526106a860e08401610548565b60e08201526101006106bb818501610548565b908201526101206106cd848201610548565b908201526101406106df84820161055f565b908201526101606106f1848201610571565b90820152610180610703848201610571565b908201526101a0610715848201610571565b908201526101c0610727848201610571565b908201526101e0610739848201610585565b9082015261020061074b848201610597565b90820152610220838101356001600160401b038082111561076b57600080fd5b610777368388016105a7565b8385015261024092508286013591508082111561079357600080fd5b61079f368388016105a7565b838501526102609250828601359150808211156107bb57600080fd5b506107c8368287016105a7565b918301919091525092915050565b6040815260006107e960408301856103ca565b82810360208401526107fb81856103ca565b95945050505050565b6000610260823603121561081757600080fd5b61081f610508565b823581526020830135602082015260408301356040820152606083013560608201526080830135608082015261085760a08401610548565b60a082015261086860c08401610548565b60c082015261087960e08401610548565b60e082015261010061088c81850161055f565b9082015261012061089e848201610571565b908201526101406108b0848201610571565b908201526101606108c2848201610571565b908201526101806108d4848201610571565b908201526101a06108e6848201610571565b908201526101c06108f8848201610585565b908201526101e061090a848201610585565b90820152610200838101356001600160401b038082111561092a57600080fd5b610936368388016105a7565b8385015261022092508286013591508082111561095257600080fd5b61095e368388016105a7565b838501526102409250828601359150808211156107bb57600080fd5b6001600160601b03198a1681526001600160a01b0389811660208301528881166040830152610120606083018190526000916109b88483018b6103ca565b9150808916608085015283820360a08501526109d482896103ca565b60c085019790975260e084019590955250509116610100909101529695505050505050565b6000610160808352610a0d8184018f6103ca565b63ffffffff9d8e16602085015260079c909c0b604084015250506060810198909852958916608088015293881660a08701529190961660c085015261ffff90951660e08401529315156101008301526101208201939093526101400191909152919050565b6001600160a01b038a81168252898116602083015261012060408301819052600091610aa08483018c6103ca565b908a16606085015283810360808501529050610abc81896103ca565b90508660a08401528560c08401528460070b60e0840152828103610100840152610ae681856103ca565b9c9b50505050505050505050505056fea26469706673582212206b918946ed9107ffd1f4f22d1afe315ae833a25486c0ccc8d06428588be7d2a364736f6c63430008190033",
}

// QuotesABI is the input ABI used to generate the binding from.
// Deprecated: Use QuotesMetaData.ABI instead.
var QuotesABI = QuotesMetaData.ABI

// QuotesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use QuotesMetaData.Bin instead.
var QuotesBin = QuotesMetaData.Bin

// DeployQuotes deploys a new Ethereum contract, binding an instance of Quotes to it.
func DeployQuotes(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Quotes, error) {
	parsed, err := QuotesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(QuotesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Quotes{QuotesCaller: QuotesCaller{contract: contract}, QuotesTransactor: QuotesTransactor{contract: contract}, QuotesFilterer: QuotesFilterer{contract: contract}}, nil
}

// Quotes is an auto generated Go binding around an Ethereum contract.
type Quotes struct {
	QuotesCaller     // Read-only binding to the contract
	QuotesTransactor // Write-only binding to the contract
	QuotesFilterer   // Log filterer for contract events
}

// QuotesCaller is an auto generated read-only Go binding around an Ethereum contract.
type QuotesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuotesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type QuotesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuotesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type QuotesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// QuotesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type QuotesSession struct {
	Contract     *Quotes           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuotesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type QuotesCallerSession struct {
	Contract *QuotesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// QuotesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type QuotesTransactorSession struct {
	Contract     *QuotesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// QuotesRaw is an auto generated low-level Go binding around an Ethereum contract.
type QuotesRaw struct {
	Contract *Quotes // Generic contract binding to access the raw methods on
}

// QuotesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type QuotesCallerRaw struct {
	Contract *QuotesCaller // Generic read-only contract binding to access the raw methods on
}

// QuotesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type QuotesTransactorRaw struct {
	Contract *QuotesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewQuotes creates a new instance of Quotes, bound to a specific deployed contract.
func NewQuotes(address common.Address, backend bind.ContractBackend) (*Quotes, error) {
	contract, err := bindQuotes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Quotes{QuotesCaller: QuotesCaller{contract: contract}, QuotesTransactor: QuotesTransactor{contract: contract}, QuotesFilterer: QuotesFilterer{contract: contract}}, nil
}

// NewQuotesCaller creates a new read-only instance of Quotes, bound to a specific deployed contract.
func NewQuotesCaller(address common.Address, caller bind.ContractCaller) (*QuotesCaller, error) {
	contract, err := bindQuotes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &QuotesCaller{contract: contract}, nil
}

// NewQuotesTransactor creates a new write-only instance of Quotes, bound to a specific deployed contract.
func NewQuotesTransactor(address common.Address, transactor bind.ContractTransactor) (*QuotesTransactor, error) {
	contract, err := bindQuotes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &QuotesTransactor{contract: contract}, nil
}

// NewQuotesFilterer creates a new log filterer instance of Quotes, bound to a specific deployed contract.
func NewQuotesFilterer(address common.Address, filterer bind.ContractFilterer) (*QuotesFilterer, error) {
	contract, err := bindQuotes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &QuotesFilterer{contract: contract}, nil
}

// bindQuotes binds a generic wrapper to an already deployed contract.
func bindQuotes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := QuotesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quotes *QuotesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quotes.Contract.QuotesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quotes *QuotesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quotes.Contract.QuotesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quotes *QuotesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quotes.Contract.QuotesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Quotes *QuotesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Quotes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Quotes *QuotesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Quotes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Quotes *QuotesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Quotes.Contract.contract.Transact(opts, method, params...)
}

// CheckAgreedAmount is a free data retrieval call binding the contract method 0x62e2e937.
//
// Solidity: function checkAgreedAmount((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, uint256 transferredAmount) pure returns()
func (_Quotes *QuotesCaller) CheckAgreedAmount(opts *bind.CallOpts, quote QuotesPegInQuote, transferredAmount *big.Int) error {
	var out []interface{}
	err := _Quotes.contract.Call(opts, &out, "checkAgreedAmount", quote, transferredAmount)

	if err != nil {
		return err
	}

	return err

}

// CheckAgreedAmount is a free data retrieval call binding the contract method 0x62e2e937.
//
// Solidity: function checkAgreedAmount((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, uint256 transferredAmount) pure returns()
func (_Quotes *QuotesSession) CheckAgreedAmount(quote QuotesPegInQuote, transferredAmount *big.Int) error {
	return _Quotes.Contract.CheckAgreedAmount(&_Quotes.CallOpts, quote, transferredAmount)
}

// CheckAgreedAmount is a free data retrieval call binding the contract method 0x62e2e937.
//
// Solidity: function checkAgreedAmount((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, uint256 transferredAmount) pure returns()
func (_Quotes *QuotesCallerSession) CheckAgreedAmount(quote QuotesPegInQuote, transferredAmount *big.Int) error {
	return _Quotes.Contract.CheckAgreedAmount(&_Quotes.CallOpts, quote, transferredAmount)
}

// EncodePegOutQuote is a free data retrieval call binding the contract method 0x1bc6d4a4.
//
// Solidity: function encodePegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesCaller) EncodePegOutQuote(opts *bind.CallOpts, quote QuotesPegOutQuote) ([]byte, error) {
	var out []interface{}
	err := _Quotes.contract.Call(opts, &out, "encodePegOutQuote", quote)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodePegOutQuote is a free data retrieval call binding the contract method 0x1bc6d4a4.
//
// Solidity: function encodePegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesSession) EncodePegOutQuote(quote QuotesPegOutQuote) ([]byte, error) {
	return _Quotes.Contract.EncodePegOutQuote(&_Quotes.CallOpts, quote)
}

// EncodePegOutQuote is a free data retrieval call binding the contract method 0x1bc6d4a4.
//
// Solidity: function encodePegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesCallerSession) EncodePegOutQuote(quote QuotesPegOutQuote) ([]byte, error) {
	return _Quotes.Contract.EncodePegOutQuote(&_Quotes.CallOpts, quote)
}

// EncodeQuote is a free data retrieval call binding the contract method 0xb77b6c0b.
//
// Solidity: function encodeQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesCaller) EncodeQuote(opts *bind.CallOpts, quote QuotesPegInQuote) ([]byte, error) {
	var out []interface{}
	err := _Quotes.contract.Call(opts, &out, "encodeQuote", quote)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeQuote is a free data retrieval call binding the contract method 0xb77b6c0b.
//
// Solidity: function encodeQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesSession) EncodeQuote(quote QuotesPegInQuote) ([]byte, error) {
	return _Quotes.Contract.EncodeQuote(&_Quotes.CallOpts, quote)
}

// EncodeQuote is a free data retrieval call binding the contract method 0xb77b6c0b.
//
// Solidity: function encodeQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) pure returns(bytes)
func (_Quotes *QuotesCallerSession) EncodeQuote(quote QuotesPegInQuote) ([]byte, error) {
	return _Quotes.Contract.EncodeQuote(&_Quotes.CallOpts, quote)
}

// SignatureValidatorMetaData contains all meta data concerning the SignatureValidator contract.
var SignatureValidatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"expectedAddress\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"usedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"IncorrectSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"quoteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"verify\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6102a7610039600b82828239805160001a607314602c57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106100355760003560e01c80631a86b5501461003a575b600080fd5b61004d610048366004610167565b610061565b604051901515815260200160405180910390f35b602081810151604080840151606085015182518084018452601c81527b0ca2ba3432b932bab69029b4b3b732b21026b2b9b9b0b3b29d05199960211b818701529251600095929391861a9286916100bc9184918b910161023f565b60408051601f1981840301815282825280516020918201206000845290830180835281905260ff861691830191909152606082018790526080820186905291506001600160a01b038a169060019060a0016020604051602081039080840390855afa15801561012f573d6000803e3d6000fd5b505050602060405103516001600160a01b031614955050505050509392505050565b634e487b7160e01b600052604160045260246000fd5b60008060006060848603121561017c57600080fd5b83356001600160a01b038116811461019357600080fd5b92506020840135915060408401356001600160401b03808211156101b657600080fd5b818601915086601f8301126101ca57600080fd5b8135818111156101dc576101dc610151565b604051601f8201601f19908116603f0116810190838211818310171561020457610204610151565b8160405282815289602084870101111561021d57600080fd5b8260208601602083013760006020848301015280955050505050509250925092565b6000835160005b818110156102605760208187018101518583015201610246565b50919091019182525060200191905056fea2646970667358221220aa28be7fb4c916241e4be09ce37be2e061b3ad7002abda8fe0eaee1d9430d8b564736f6c63430008190033",
}

// SignatureValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use SignatureValidatorMetaData.ABI instead.
var SignatureValidatorABI = SignatureValidatorMetaData.ABI

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
	parsed, err := SignatureValidatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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
