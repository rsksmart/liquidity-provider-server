// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

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

// PeginContractMetaData contains all meta data concerning the PeginContract contract.
var PeginContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"callForUser\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentContribution\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeeCollector\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeePercentage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinPegIn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getQuoteStatus\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPegIn.PegInStates\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashPegInQuote\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerPegIn\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRawTransaction\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"partialMerkleTree\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePegInDepositAddress\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BalanceDecrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BalanceIncrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BridgeCapExceeded\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"errorCode\",\"type\":\"int256\",\"indexed\":true,\"internalType\":\"int256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CallForUser\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"gasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInRegistered\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"transferredAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refund\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"success\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AmountUnderMinimum\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGas\",\"inputs\":[{\"name\":\"gasLeft\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasRequired\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidRefundAddress\",\"inputs\":[{\"name\":\"refundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NotEnoughConfirmations\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuoteAlreadyProcessed\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnexpectedBridgeError\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"int256\",\"internalType\":\"int256\"}]}]",
	ID:  "PeginContract",
}

// PeginContract is an auto generated Go binding around an Ethereum contract.
type PeginContract struct {
	abi abi.ABI
}

// NewPeginContract creates a new instance of PeginContract.
func NewPeginContract() *PeginContract {
	parsed, err := PeginContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PeginContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PeginContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackCallForUser is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7a3dc3c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) PackCallForUser(quote QuotesPegInQuote) []byte {
	enc, err := peginContract.abi.Pack("callForUser", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCallForUser is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7a3dc3c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) TryPackCallForUser(quote QuotesPegInQuote) ([]byte, error) {
	return peginContract.abi.Pack("callForUser", quote)
}

// UnpackCallForUser is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7a3dc3c.
//
// Solidity: function callForUser((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) payable returns(bool)
func (peginContract *PeginContract) UnpackCallForUser(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("callForUser", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit() payable returns()
func (peginContract *PeginContract) PackDeposit() []byte {
	enc, err := peginContract.abi.Pack("deposit")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd0e30db0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit() payable returns()
func (peginContract *PeginContract) TryPackDeposit() ([]byte, error) {
	return peginContract.abi.Pack("deposit")
}

// PackGetBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8b2cb4f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) PackGetBalance(addr common.Address) []byte {
	enc, err := peginContract.abi.Pack("getBalance", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBalance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8b2cb4f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) TryPackGetBalance(addr common.Address) ([]byte, error) {
	return peginContract.abi.Pack("getBalance", addr)
}

// UnpackGetBalance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (peginContract *PeginContract) UnpackGetBalance(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getBalance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetCurrentContribution is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8623d53.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (peginContract *PeginContract) PackGetCurrentContribution() []byte {
	enc, err := peginContract.abi.Pack("getCurrentContribution")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetCurrentContribution is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb8623d53.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (peginContract *PeginContract) TryPackGetCurrentContribution() ([]byte, error) {
	return peginContract.abi.Pack("getCurrentContribution")
}

// UnpackGetCurrentContribution is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb8623d53.
//
// Solidity: function getCurrentContribution() view returns(uint256)
func (peginContract *PeginContract) UnpackGetCurrentContribution(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getCurrentContribution", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x12fde4b7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFeeCollector() view returns(address)
func (peginContract *PeginContract) PackGetFeeCollector() []byte {
	enc, err := peginContract.abi.Pack("getFeeCollector")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFeeCollector is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x12fde4b7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFeeCollector() view returns(address)
func (peginContract *PeginContract) TryPackGetFeeCollector() ([]byte, error) {
	return peginContract.abi.Pack("getFeeCollector")
}

// UnpackGetFeeCollector is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x12fde4b7.
//
// Solidity: function getFeeCollector() view returns(address)
func (peginContract *PeginContract) UnpackGetFeeCollector(data []byte) (common.Address, error) {
	out, err := peginContract.abi.Unpack("getFeeCollector", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetFeePercentage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11efbf61.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (peginContract *PeginContract) PackGetFeePercentage() []byte {
	enc, err := peginContract.abi.Pack("getFeePercentage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFeePercentage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x11efbf61.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (peginContract *PeginContract) TryPackGetFeePercentage() ([]byte, error) {
	return peginContract.abi.Pack("getFeePercentage")
}

// UnpackGetFeePercentage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x11efbf61.
//
// Solidity: function getFeePercentage() view returns(uint256)
func (peginContract *PeginContract) UnpackGetFeePercentage(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getFeePercentage", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa88dcde.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) PackGetMinPegIn() []byte {
	enc, err := peginContract.abi.Pack("getMinPegIn")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMinPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfa88dcde.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) TryPackGetMinPegIn() ([]byte, error) {
	return peginContract.abi.Pack("getMinPegIn")
}

// UnpackGetMinPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfa88dcde.
//
// Solidity: function getMinPegIn() view returns(uint256)
func (peginContract *PeginContract) UnpackGetMinPegIn(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("getMinPegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetQuoteStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf93c8ec2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) PackGetQuoteStatus(quoteHash [32]byte) []byte {
	enc, err := peginContract.abi.Pack("getQuoteStatus", quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetQuoteStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf93c8ec2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) TryPackGetQuoteStatus(quoteHash [32]byte) ([]byte, error) {
	return peginContract.abi.Pack("getQuoteStatus", quoteHash)
}

// UnpackGetQuoteStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf93c8ec2.
//
// Solidity: function getQuoteStatus(bytes32 quoteHash) view returns(uint8)
func (peginContract *PeginContract) UnpackGetQuoteStatus(data []byte) (uint8, error) {
	out, err := peginContract.abi.Unpack("getQuoteStatus", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackHashPegInQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf218a7d8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) PackHashPegInQuote(quote QuotesPegInQuote) []byte {
	enc, err := peginContract.abi.Pack("hashPegInQuote", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashPegInQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf218a7d8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) TryPackHashPegInQuote(quote QuotesPegInQuote) ([]byte, error) {
	return peginContract.abi.Pack("hashPegInQuote", quote)
}

// UnpackHashPegInQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf218a7d8.
//
// Solidity: function hashPegInQuote((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote) view returns(bytes32)
func (peginContract *PeginContract) UnpackHashPegInQuote(data []byte) ([32]byte, error) {
	out, err := peginContract.abi.Unpack("hashPegInQuote", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6da66355.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pause(string reason) returns()
func (peginContract *PeginContract) PackPause(reason string) []byte {
	enc, err := peginContract.abi.Pack("pause", reason)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6da66355.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pause(string reason) returns()
func (peginContract *PeginContract) TryPackPause(reason string) ([]byte, error) {
	return peginContract.abi.Pack("pause", reason)
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) PackPauseStatus() []byte {
	enc, err := peginContract.abi.Pack("pauseStatus")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) TryPackPauseStatus() ([]byte, error) {
	return peginContract.abi.Pack("pauseStatus")
}

// PauseStatusOutput serves as a container for the return parameters of contract
// method PauseStatus.
type PauseStatusOutput struct {
	IsPaused bool
	Reason   string
	Since    uint64
}

// UnpackPauseStatus is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x466916ca.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (peginContract *PeginContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := peginContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackRegisterPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3823c753.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) PackRegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) []byte {
	enc, err := peginContract.abi.Pack("registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3823c753.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) TryPackRegisterPegIn(quote QuotesPegInQuote, signature []byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("registerPegIn", quote, signature, btcRawTransaction, partialMerkleTree, height)
}

// UnpackRegisterPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3823c753.
//
// Solidity: function registerPegIn((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes signature, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256)
func (peginContract *PeginContract) UnpackRegisterPegIn(data []byte) (*big.Int, error) {
	out, err := peginContract.abi.Unpack("registerPegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpause() returns()
func (peginContract *PeginContract) PackUnpause() []byte {
	enc, err := peginContract.abi.Pack("unpause")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unpause() returns()
func (peginContract *PeginContract) TryPackUnpause() ([]byte, error) {
	return peginContract.abi.Pack("unpause")
}

// PackValidatePegInDepositAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9accea2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) PackValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) []byte {
	enc, err := peginContract.abi.Pack("validatePegInDepositAddress", quote, depositAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidatePegInDepositAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9accea2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) TryPackValidatePegInDepositAddress(quote QuotesPegInQuote, depositAddress []byte) ([]byte, error) {
	return peginContract.abi.Pack("validatePegInDepositAddress", quote, depositAddress)
}

// UnpackValidatePegInDepositAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe9accea2.
//
// Solidity: function validatePegInDepositAddress((uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes depositAddress) view returns(bool)
func (peginContract *PeginContract) UnpackValidatePegInDepositAddress(data []byte) (bool, error) {
	out, err := peginContract.abi.Unpack("validatePegInDepositAddress", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(uint256 amount) returns()
func (peginContract *PeginContract) PackWithdraw(amount *big.Int) []byte {
	enc, err := peginContract.abi.Pack("withdraw", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2e1a7d4d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(uint256 amount) returns()
func (peginContract *PeginContract) TryPackWithdraw(amount *big.Int) ([]byte, error) {
	return peginContract.abi.Pack("withdraw", amount)
}

// PeginContractBalanceDecrease represents a BalanceDecrease event raised by the PeginContract contract.
type PeginContractBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractBalanceDecreaseEventName = "BalanceDecrease"

// ContractEventName returns the user-defined event name.
func (PeginContractBalanceDecrease) ContractEventName() string {
	return PeginContractBalanceDecreaseEventName
}

// UnpackBalanceDecreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (peginContract *PeginContract) UnpackBalanceDecreaseEvent(log *types.Log) (*PeginContractBalanceDecrease, error) {
	event := "BalanceDecrease"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBalanceDecrease)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractBalanceIncrease represents a BalanceIncrease event raised by the PeginContract contract.
type PeginContractBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractBalanceIncreaseEventName = "BalanceIncrease"

// ContractEventName returns the user-defined event name.
func (PeginContractBalanceIncrease) ContractEventName() string {
	return PeginContractBalanceIncreaseEventName
}

// UnpackBalanceIncreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (peginContract *PeginContract) UnpackBalanceIncreaseEvent(log *types.Log) (*PeginContractBalanceIncrease, error) {
	event := "BalanceIncrease"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBalanceIncrease)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractBridgeCapExceeded represents a BridgeCapExceeded event raised by the PeginContract contract.
type PeginContractBridgeCapExceeded struct {
	QuoteHash [32]byte
	ErrorCode *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractBridgeCapExceededEventName = "BridgeCapExceeded"

// ContractEventName returns the user-defined event name.
func (PeginContractBridgeCapExceeded) ContractEventName() string {
	return PeginContractBridgeCapExceededEventName
}

// UnpackBridgeCapExceededEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BridgeCapExceeded(bytes32 indexed quoteHash, int256 indexed errorCode)
func (peginContract *PeginContract) UnpackBridgeCapExceededEvent(log *types.Log) (*PeginContractBridgeCapExceeded, error) {
	event := "BridgeCapExceeded"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractBridgeCapExceeded)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractCallForUser represents a CallForUser event raised by the PeginContract contract.
type PeginContractCallForUser struct {
	From      common.Address
	Dest      common.Address
	QuoteHash [32]byte
	GasLimit  *big.Int
	Value     *big.Int
	Data      []byte
	Success   bool
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractCallForUserEventName = "CallForUser"

// ContractEventName returns the user-defined event name.
func (PeginContractCallForUser) ContractEventName() string {
	return PeginContractCallForUserEventName
}

// UnpackCallForUserEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CallForUser(address indexed from, address indexed dest, bytes32 indexed quoteHash, uint256 gasLimit, uint256 value, bytes data, bool success)
func (peginContract *PeginContract) UnpackCallForUserEvent(log *types.Log) (*PeginContractCallForUser, error) {
	event := "CallForUser"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractCallForUser)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractPegInRegistered represents a PegInRegistered event raised by the PeginContract contract.
type PeginContractPegInRegistered struct {
	QuoteHash         [32]byte
	TransferredAmount *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const PeginContractPegInRegisteredEventName = "PegInRegistered"

// ContractEventName returns the user-defined event name.
func (PeginContractPegInRegistered) ContractEventName() string {
	return PeginContractPegInRegisteredEventName
}

// UnpackPegInRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInRegistered(bytes32 indexed quoteHash, uint256 indexed transferredAmount)
func (peginContract *PeginContract) UnpackPegInRegisteredEvent(log *types.Log) (*PeginContractPegInRegistered, error) {
	event := "PegInRegistered"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractPegInRegistered)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractRefund represents a Refund event raised by the PeginContract contract.
type PeginContractRefund struct {
	Dest      common.Address
	QuoteHash [32]byte
	Amount    *big.Int
	Success   bool
	Raw       *types.Log // Blockchain specific contextual infos
}

const PeginContractRefundEventName = "Refund"

// ContractEventName returns the user-defined event name.
func (PeginContractRefund) ContractEventName() string {
	return PeginContractRefundEventName
}

// UnpackRefundEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Refund(address indexed dest, bytes32 indexed quoteHash, uint256 indexed amount, bool success)
func (peginContract *PeginContract) UnpackRefundEvent(log *types.Log) (*PeginContractRefund, error) {
	event := "Refund"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractRefund)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// PeginContractWithdrawal represents a Withdrawal event raised by the PeginContract contract.
type PeginContractWithdrawal struct {
	From   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PeginContractWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (PeginContractWithdrawal) ContractEventName() string {
	return PeginContractWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(address indexed from, uint256 indexed amount)
func (peginContract *PeginContract) UnpackWithdrawalEvent(log *types.Log) (*PeginContractWithdrawal, error) {
	event := "Withdrawal"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractWithdrawal)
	if len(log.Data) > 0 {
		if err := peginContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginContract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (peginContract *PeginContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], peginContract.abi.Errors["AmountUnderMinimum"].ID.Bytes()[:4]) {
		return peginContract.UnpackAmountUnderMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InsufficientGas"].ID.Bytes()[:4]) {
		return peginContract.UnpackInsufficientGasError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["InvalidRefundAddress"].ID.Bytes()[:4]) {
		return peginContract.UnpackInvalidRefundAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["NotEnoughConfirmations"].ID.Bytes()[:4]) {
		return peginContract.UnpackNotEnoughConfirmationsError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["QuoteAlreadyProcessed"].ID.Bytes()[:4]) {
		return peginContract.UnpackQuoteAlreadyProcessedError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginContract.abi.Errors["UnexpectedBridgeError"].ID.Bytes()[:4]) {
		return peginContract.UnpackUnexpectedBridgeErrorError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PeginContractAmountUnderMinimum represents a AmountUnderMinimum error raised by the PeginContract contract.
type PeginContractAmountUnderMinimum struct {
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AmountUnderMinimum(uint256 amount)
func PeginContractAmountUnderMinimumErrorID() common.Hash {
	return common.HexToHash("0x12b21ac5391b7c8c532d3e1a87c2e8173cdc084bbd0c28c2cabc2c2c296d2f1a")
}

// UnpackAmountUnderMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AmountUnderMinimum(uint256 amount)
func (peginContract *PeginContract) UnpackAmountUnderMinimumError(raw []byte) (*PeginContractAmountUnderMinimum, error) {
	out := new(PeginContractAmountUnderMinimum)
	if err := peginContract.abi.UnpackIntoInterface(out, "AmountUnderMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInsufficientGas represents a InsufficientGas error raised by the PeginContract contract.
type PeginContractInsufficientGas struct {
	GasLeft     *big.Int
	GasRequired *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientGas(uint256 gasLeft, uint256 gasRequired)
func PeginContractInsufficientGasErrorID() common.Hash {
	return common.HexToHash("0x23e228cb30ba888d86279f51c27a6690ec0fe016ea074bbd6332de5b2571f9b3")
}

// UnpackInsufficientGasError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientGas(uint256 gasLeft, uint256 gasRequired)
func (peginContract *PeginContract) UnpackInsufficientGasError(raw []byte) (*PeginContractInsufficientGas, error) {
	out := new(PeginContractInsufficientGas)
	if err := peginContract.abi.UnpackIntoInterface(out, "InsufficientGas", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractInvalidRefundAddress represents a InvalidRefundAddress error raised by the PeginContract contract.
type PeginContractInvalidRefundAddress struct {
	RefundAddress []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidRefundAddress(bytes refundAddress)
func PeginContractInvalidRefundAddressErrorID() common.Hash {
	return common.HexToHash("0x17b97374cf0657ace771eafcd95943e733cce4db0d5b4e147411067d764c8f99")
}

// UnpackInvalidRefundAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidRefundAddress(bytes refundAddress)
func (peginContract *PeginContract) UnpackInvalidRefundAddressError(raw []byte) (*PeginContractInvalidRefundAddress, error) {
	out := new(PeginContractInvalidRefundAddress)
	if err := peginContract.abi.UnpackIntoInterface(out, "InvalidRefundAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractNotEnoughConfirmations represents a NotEnoughConfirmations error raised by the PeginContract contract.
type PeginContractNotEnoughConfirmations struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEnoughConfirmations()
func PeginContractNotEnoughConfirmationsErrorID() common.Hash {
	return common.HexToHash("0xb9310b562727f0fb376475537f3a4e5f39f5fed59dbda43f984b828c8ef037d0")
}

// UnpackNotEnoughConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEnoughConfirmations()
func (peginContract *PeginContract) UnpackNotEnoughConfirmationsError(raw []byte) (*PeginContractNotEnoughConfirmations, error) {
	out := new(PeginContractNotEnoughConfirmations)
	if err := peginContract.abi.UnpackIntoInterface(out, "NotEnoughConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractQuoteAlreadyProcessed represents a QuoteAlreadyProcessed error raised by the PeginContract contract.
type PeginContractQuoteAlreadyProcessed struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteAlreadyProcessed(bytes32 quoteHash)
func PeginContractQuoteAlreadyProcessedErrorID() common.Hash {
	return common.HexToHash("0xda4bb665b5917c51a08a9d79f0cf72c95ee90f52ca97526abf4b52ae2d737c77")
}

// UnpackQuoteAlreadyProcessedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteAlreadyProcessed(bytes32 quoteHash)
func (peginContract *PeginContract) UnpackQuoteAlreadyProcessedError(raw []byte) (*PeginContractQuoteAlreadyProcessed, error) {
	out := new(PeginContractQuoteAlreadyProcessed)
	if err := peginContract.abi.UnpackIntoInterface(out, "QuoteAlreadyProcessed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginContractUnexpectedBridgeError represents a UnexpectedBridgeError error raised by the PeginContract contract.
type PeginContractUnexpectedBridgeError struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnexpectedBridgeError(int256 errorCode)
func PeginContractUnexpectedBridgeErrorErrorID() common.Hash {
	return common.HexToHash("0xab2e19b07f35862bb35f42ec9b8ce9397c135ed8ed749f2117205465c4166a09")
}

// UnpackUnexpectedBridgeErrorError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnexpectedBridgeError(int256 errorCode)
func (peginContract *PeginContract) UnpackUnexpectedBridgeErrorError(raw []byte) (*PeginContractUnexpectedBridgeError, error) {
	out := new(PeginContractUnexpectedBridgeError)
	if err := peginContract.abi.UnpackIntoInterface(out, "UnexpectedBridgeError", raw); err != nil {
		return nil, err
	}
	return out, nil
}
