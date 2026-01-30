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

// FlyoverMetaData contains all meta data concerning the Flyover contract.
var FlyoverMetaData = bind.MetaData{
	ABI: "[{\"type\":\"error\",\"name\":\"EmptyBlockHeader\",\"inputs\":[{\"name\":\"heightOrHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IncorrectContract\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSender\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoBalance\",\"inputs\":[{\"name\":\"wanted\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NoContract\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Overflow\",\"inputs\":[{\"name\":\"passedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PaymentFailed\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"PaymentNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProviderNotRegistered\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"QuoteNotFound\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "Flyover",
}

// Flyover is an auto generated Go binding around an Ethereum contract.
type Flyover struct {
	abi abi.ABI
}

// NewFlyover creates a new instance of Flyover.
func NewFlyover() *Flyover {
	parsed, err := FlyoverMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Flyover{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Flyover) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (flyover *Flyover) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], flyover.abi.Errors["EmptyBlockHeader"].ID.Bytes()[:4]) {
		return flyover.UnpackEmptyBlockHeaderError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["IncorrectContract"].ID.Bytes()[:4]) {
		return flyover.UnpackIncorrectContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["InsufficientAmount"].ID.Bytes()[:4]) {
		return flyover.UnpackInsufficientAmountError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["InvalidAddress"].ID.Bytes()[:4]) {
		return flyover.UnpackInvalidAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["InvalidSender"].ID.Bytes()[:4]) {
		return flyover.UnpackInvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["NoBalance"].ID.Bytes()[:4]) {
		return flyover.UnpackNoBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["NoContract"].ID.Bytes()[:4]) {
		return flyover.UnpackNoContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["Overflow"].ID.Bytes()[:4]) {
		return flyover.UnpackOverflowError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["PaymentFailed"].ID.Bytes()[:4]) {
		return flyover.UnpackPaymentFailedError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["PaymentNotAllowed"].ID.Bytes()[:4]) {
		return flyover.UnpackPaymentNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["ProviderNotRegistered"].ID.Bytes()[:4]) {
		return flyover.UnpackProviderNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyover.abi.Errors["QuoteNotFound"].ID.Bytes()[:4]) {
		return flyover.UnpackQuoteNotFoundError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// FlyoverEmptyBlockHeader represents a EmptyBlockHeader error raised by the Flyover contract.
type FlyoverEmptyBlockHeader struct {
	HeightOrHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error EmptyBlockHeader(bytes32 heightOrHash)
func FlyoverEmptyBlockHeaderErrorID() common.Hash {
	return common.HexToHash("0xc1a923b4e595599b5ebca706a34bfaa111ec5aad01c417609e91334f899d99e4")
}

// UnpackEmptyBlockHeaderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error EmptyBlockHeader(bytes32 heightOrHash)
func (flyover *Flyover) UnpackEmptyBlockHeaderError(raw []byte) (*FlyoverEmptyBlockHeader, error) {
	out := new(FlyoverEmptyBlockHeader)
	if err := flyover.abi.UnpackIntoInterface(out, "EmptyBlockHeader", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverIncorrectContract represents a IncorrectContract error raised by the Flyover contract.
type FlyoverIncorrectContract struct {
	Expected common.Address
	Actual   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectContract(address expected, address actual)
func FlyoverIncorrectContractErrorID() common.Hash {
	return common.HexToHash("0x367b77278f2bd6dd9afab1117babaedb89c7f420646aa9a343c7f6bd654b7740")
}

// UnpackIncorrectContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectContract(address expected, address actual)
func (flyover *Flyover) UnpackIncorrectContractError(raw []byte) (*FlyoverIncorrectContract, error) {
	out := new(FlyoverIncorrectContract)
	if err := flyover.abi.UnpackIntoInterface(out, "IncorrectContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverInsufficientAmount represents a InsufficientAmount error raised by the Flyover contract.
type FlyoverInsufficientAmount struct {
	Amount *big.Int
	Target *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientAmount(uint256 amount, uint256 target)
func FlyoverInsufficientAmountErrorID() common.Hash {
	return common.HexToHash("0x6d400e382e49fdfa5a03b18b5b2c938638a3fb351ac4810276f70a093eb3f20f")
}

// UnpackInsufficientAmountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientAmount(uint256 amount, uint256 target)
func (flyover *Flyover) UnpackInsufficientAmountError(raw []byte) (*FlyoverInsufficientAmount, error) {
	out := new(FlyoverInsufficientAmount)
	if err := flyover.abi.UnpackIntoInterface(out, "InsufficientAmount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverInvalidAddress represents a InvalidAddress error raised by the Flyover contract.
type FlyoverInvalidAddress struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidAddress(address addr)
func FlyoverInvalidAddressErrorID() common.Hash {
	return common.HexToHash("0x8e4c8aa64faa2ab4276a2c6878068416c0eae6ead59e25aa3c28d438ef21fae3")
}

// UnpackInvalidAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidAddress(address addr)
func (flyover *Flyover) UnpackInvalidAddressError(raw []byte) (*FlyoverInvalidAddress, error) {
	out := new(FlyoverInvalidAddress)
	if err := flyover.abi.UnpackIntoInterface(out, "InvalidAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverInvalidSender represents a InvalidSender error raised by the Flyover contract.
type FlyoverInvalidSender struct {
	Expected common.Address
	Actual   common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSender(address expected, address actual)
func FlyoverInvalidSenderErrorID() common.Hash {
	return common.HexToHash("0xe1130dbad6e77228912cd79cc3b53cd156f090ef6a73d9fdb2720c4f9d40d9d3")
}

// UnpackInvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSender(address expected, address actual)
func (flyover *Flyover) UnpackInvalidSenderError(raw []byte) (*FlyoverInvalidSender, error) {
	out := new(FlyoverInvalidSender)
	if err := flyover.abi.UnpackIntoInterface(out, "InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverNoBalance represents a NoBalance error raised by the Flyover contract.
type FlyoverNoBalance struct {
	Wanted *big.Int
	Actual *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoBalance(uint256 wanted, uint256 actual)
func FlyoverNoBalanceErrorID() common.Hash {
	return common.HexToHash("0x292266533ee1631c0f0faf752ebfa5783238c0e3e2fcdef002a1685294062289")
}

// UnpackNoBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoBalance(uint256 wanted, uint256 actual)
func (flyover *Flyover) UnpackNoBalanceError(raw []byte) (*FlyoverNoBalance, error) {
	out := new(FlyoverNoBalance)
	if err := flyover.abi.UnpackIntoInterface(out, "NoBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverNoContract represents a NoContract error raised by the Flyover contract.
type FlyoverNoContract struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoContract(address addr)
func FlyoverNoContractErrorID() common.Hash {
	return common.HexToHash("0x5f15d672b6235f8600ffc72925d8d2f9dcea14be067296327891153847185a5c")
}

// UnpackNoContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoContract(address addr)
func (flyover *Flyover) UnpackNoContractError(raw []byte) (*FlyoverNoContract, error) {
	out := new(FlyoverNoContract)
	if err := flyover.abi.UnpackIntoInterface(out, "NoContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverOverflow represents a Overflow error raised by the Flyover contract.
type FlyoverOverflow struct {
	PassedAmount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Overflow(uint256 passedAmount)
func FlyoverOverflowErrorID() common.Hash {
	return common.HexToHash("0xe0fb6a7ce291b396fa814871fbb6fcc26c1a1454a6e18a2e7c911a8763b928dc")
}

// UnpackOverflowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Overflow(uint256 passedAmount)
func (flyover *Flyover) UnpackOverflowError(raw []byte) (*FlyoverOverflow, error) {
	out := new(FlyoverOverflow)
	if err := flyover.abi.UnpackIntoInterface(out, "Overflow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverPaymentFailed represents a PaymentFailed error raised by the Flyover contract.
type FlyoverPaymentFailed struct {
	Addr   common.Address
	Amount *big.Int
	Reason []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentFailed(address addr, uint256 amount, bytes reason)
func FlyoverPaymentFailedErrorID() common.Hash {
	return common.HexToHash("0xadca8d516d2aaa483a86cefb25d722eccb15750e54cc37c21033c80cc79b13e3")
}

// UnpackPaymentFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentFailed(address addr, uint256 amount, bytes reason)
func (flyover *Flyover) UnpackPaymentFailedError(raw []byte) (*FlyoverPaymentFailed, error) {
	out := new(FlyoverPaymentFailed)
	if err := flyover.abi.UnpackIntoInterface(out, "PaymentFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverPaymentNotAllowed represents a PaymentNotAllowed error raised by the Flyover contract.
type FlyoverPaymentNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentNotAllowed()
func FlyoverPaymentNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x8619bd43ab22b4b01742bd29d231dff1e50413ee3a444878bed65970c80c97df")
}

// UnpackPaymentNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentNotAllowed()
func (flyover *Flyover) UnpackPaymentNotAllowedError(raw []byte) (*FlyoverPaymentNotAllowed, error) {
	out := new(FlyoverPaymentNotAllowed)
	if err := flyover.abi.UnpackIntoInterface(out, "PaymentNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverProviderNotRegistered represents a ProviderNotRegistered error raised by the Flyover contract.
type FlyoverProviderNotRegistered struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ProviderNotRegistered(address from)
func FlyoverProviderNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x232cb27a4b9e96657e43917628f0b0ddd34885ba8495a2108b78da7512210fb9")
}

// UnpackProviderNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ProviderNotRegistered(address from)
func (flyover *Flyover) UnpackProviderNotRegisteredError(raw []byte) (*FlyoverProviderNotRegistered, error) {
	out := new(FlyoverProviderNotRegistered)
	if err := flyover.abi.UnpackIntoInterface(out, "ProviderNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverQuoteNotFound represents a QuoteNotFound error raised by the Flyover contract.
type FlyoverQuoteNotFound struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteNotFound(bytes32 quoteHash)
func FlyoverQuoteNotFoundErrorID() common.Hash {
	return common.HexToHash("0xa871f056dfd91173197c96810e0eba4f526ef57ede4af25e166ed0471a99b2d7")
}

// UnpackQuoteNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteNotFound(bytes32 quoteHash)
func (flyover *Flyover) UnpackQuoteNotFoundError(raw []byte) (*FlyoverQuoteNotFound, error) {
	out := new(FlyoverQuoteNotFound)
	if err := flyover.abi.UnpackIntoInterface(out, "QuoteNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}
