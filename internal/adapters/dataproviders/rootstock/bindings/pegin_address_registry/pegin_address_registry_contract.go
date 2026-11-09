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

// PegInAddressRegistryContractMetaData contains all meta data concerning the PegInAddressRegistryContract contract.
var PegInAddressRegistryContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getMinDepositSats\",\"inputs\":[],\"outputs\":[{\"name\":\"minDepositSats\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInAddress\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"payload\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"encoding\",\"type\":\"uint8\",\"internalType\":\"enumIPegInAddressRegistry.Encoding\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInAddresses\",\"inputs\":[{\"name\":\"rskAddrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"payloads\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"encoding\",\"type\":\"uint8\",\"internalType\":\"enumIPegInAddressRegistry.Encoding\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrant\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"registrant\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrationBlock\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"registrationBlock\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrationRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"registrationRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"registered\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAddress\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleBranchPath\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AddressRegistered\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"registrant\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"registrationRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyRegistered\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BatchTooLarge\",\"inputs\":[{\"name\":\"requested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfigMinBelowBridge\",\"inputs\":[{\"name\":\"protocolMinSats\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bridgeMinSats\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ConfigurationsNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositBelowMinimum\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DepositNotConfirmed\",\"inputs\":[{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"DepositOutputNotFound\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PegInContractNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnableToGetConfirmations\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"int256\",\"internalType\":\"int256\"}]}]",
	ID:  "PegInAddressRegistryContract",
}

// PegInAddressRegistryContract is an auto generated Go binding around an Ethereum contract.
type PegInAddressRegistryContract struct {
	abi abi.ABI
}

// NewPegInAddressRegistryContract creates a new instance of PegInAddressRegistryContract.
func NewPegInAddressRegistryContract() *PegInAddressRegistryContract {
	parsed, err := PegInAddressRegistryContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PegInAddressRegistryContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PegInAddressRegistryContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetMinDepositSats is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a2b4622.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMinDepositSats() view returns(uint256 minDepositSats)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetMinDepositSats() []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getMinDepositSats")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMinDepositSats is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0a2b4622.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMinDepositSats() view returns(uint256 minDepositSats)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetMinDepositSats() ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getMinDepositSats")
}

// UnpackGetMinDepositSats is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0a2b4622.
//
// Solidity: function getMinDepositSats() view returns(uint256 minDepositSats)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetMinDepositSats(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getMinDepositSats", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegInAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x831adb16.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInAddress(address rskAddr) view returns(bytes payload, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetPegInAddress(rskAddr common.Address) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getPegInAddress", rskAddr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x831adb16.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInAddress(address rskAddr) view returns(bytes payload, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetPegInAddress(rskAddr common.Address) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getPegInAddress", rskAddr)
}

// GetPegInAddressOutput serves as a container for the return parameters of contract
// method GetPegInAddress.
type GetPegInAddressOutput struct {
	Payload  []byte
	Encoding uint8
}

// UnpackGetPegInAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x831adb16.
//
// Solidity: function getPegInAddress(address rskAddr) view returns(bytes payload, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetPegInAddress(data []byte) (GetPegInAddressOutput, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getPegInAddress", data)
	outstruct := new(GetPegInAddressOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Payload = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Encoding = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackGetPegInAddresses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd4a25a4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInAddresses(address[] rskAddrs) view returns(bytes[] payloads, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetPegInAddresses(rskAddrs []common.Address) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getPegInAddresses", rskAddrs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInAddresses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd4a25a4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInAddresses(address[] rskAddrs) view returns(bytes[] payloads, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetPegInAddresses(rskAddrs []common.Address) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getPegInAddresses", rskAddrs)
}

// GetPegInAddressesOutput serves as a container for the return parameters of contract
// method GetPegInAddresses.
type GetPegInAddressesOutput struct {
	Payloads [][]byte
	Encoding uint8
}

// UnpackGetPegInAddresses is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd4a25a4.
//
// Solidity: function getPegInAddresses(address[] rskAddrs) view returns(bytes[] payloads, uint8 encoding)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetPegInAddresses(data []byte) (GetPegInAddressesOutput, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getPegInAddresses", data)
	outstruct := new(GetPegInAddressesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Payloads = *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	outstruct.Encoding = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackGetRegistrant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b7dfb40.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRegistrant(address rskAddr) view returns(address registrant)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetRegistrant(rskAddr common.Address) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getRegistrant", rskAddr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrant is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b7dfb40.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrant(address rskAddr) view returns(address registrant)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetRegistrant(rskAddr common.Address) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getRegistrant", rskAddr)
}

// UnpackGetRegistrant is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6b7dfb40.
//
// Solidity: function getRegistrant(address rskAddr) view returns(address registrant)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetRegistrant(data []byte) (common.Address, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getRegistrant", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetRegistrationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5017bc63.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRegistrationBlock(address rskAddr) view returns(uint96 registrationBlock)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetRegistrationBlock(rskAddr common.Address) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getRegistrationBlock", rskAddr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5017bc63.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrationBlock(address rskAddr) view returns(uint96 registrationBlock)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetRegistrationBlock(rskAddr common.Address) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getRegistrationBlock", rskAddr)
}

// UnpackGetRegistrationBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5017bc63.
//
// Solidity: function getRegistrationBlock(address rskAddr) view returns(uint96 registrationBlock)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetRegistrationBlock(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getRegistrationBlock", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRegistrationRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe625000.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRegistrationRoot() view returns(bytes32 registrationRoot)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackGetRegistrationRoot() []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("getRegistrationRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrationRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe625000.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrationRoot() view returns(bytes32 registrationRoot)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackGetRegistrationRoot() ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("getRegistrationRoot")
}

// UnpackGetRegistrationRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe625000.
//
// Solidity: function getRegistrationRoot() view returns(bytes32 registrationRoot)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackGetRegistrationRoot(data []byte) ([32]byte, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("getRegistrationRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3c5a547.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isRegistered(address rskAddr) view returns(bool registered)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackIsRegistered(rskAddr common.Address) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("isRegistered", rskAddr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3c5a547.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isRegistered(address rskAddr) view returns(bool registered)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackIsRegistered(rskAddr common.Address) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("isRegistered", rskAddr)
}

// UnpackIsRegistered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3c5a547.
//
// Solidity: function isRegistered(address rskAddr) view returns(bool registered)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackIsRegistered(data []byte) (bool, error) {
	out, err := pegInAddressRegistryContract.abi.Unpack("isRegistered", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackRegisterAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7aefc26e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerAddress(address rskAddr, bytes btcTxSerialized, bytes32 btcBlockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (pegInAddressRegistryContract *PegInAddressRegistryContract) PackRegisterAddress(rskAddr common.Address, btcTxSerialized []byte, btcBlockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) []byte {
	enc, err := pegInAddressRegistryContract.abi.Pack("registerAddress", rskAddr, btcTxSerialized, btcBlockHash, merkleBranchPath, merkleBranchHashes)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7aefc26e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerAddress(address rskAddr, bytes btcTxSerialized, bytes32 btcBlockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (pegInAddressRegistryContract *PegInAddressRegistryContract) TryPackRegisterAddress(rskAddr common.Address, btcTxSerialized []byte, btcBlockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) ([]byte, error) {
	return pegInAddressRegistryContract.abi.Pack("registerAddress", rskAddr, btcTxSerialized, btcBlockHash, merkleBranchPath, merkleBranchHashes)
}

// PegInAddressRegistryContractAddressRegistered represents a AddressRegistered event raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractAddressRegistered struct {
	RskAddr          common.Address
	Registrant       common.Address
	RegistrationRoot [32]byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryContractAddressRegisteredEventName = "AddressRegistered"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryContractAddressRegistered) ContractEventName() string {
	return PegInAddressRegistryContractAddressRegisteredEventName
}

// UnpackAddressRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AddressRegistered(address indexed rskAddr, address indexed registrant, bytes32 registrationRoot)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackAddressRegisteredEvent(log *types.Log) (*PegInAddressRegistryContractAddressRegistered, error) {
	event := "AddressRegistered"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistryContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryContractAddressRegistered)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistryContract.abi.Events[event].Inputs {
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
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["AddressAlreadyRegistered"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackAddressAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["BatchTooLarge"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackBatchTooLargeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["ConfigMinBelowBridge"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackConfigMinBelowBridgeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["ConfigurationsNotSet"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackConfigurationsNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["DepositBelowMinimum"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackDepositBelowMinimumError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["DepositNotConfirmed"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackDepositNotConfirmedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["DepositOutputNotFound"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackDepositOutputNotFoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["PegInContractNotSet"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackPegInContractNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistryContract.abi.Errors["UnableToGetConfirmations"].ID.Bytes()[:4]) {
		return pegInAddressRegistryContract.UnpackUnableToGetConfirmationsError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PegInAddressRegistryContractAddressAlreadyRegistered represents a AddressAlreadyRegistered error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractAddressAlreadyRegistered struct {
	RskAddr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressAlreadyRegistered(address rskAddr)
func PegInAddressRegistryContractAddressAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0xef3c800bde1cd5a483e6e96380c892083c77f073dac6b6e59d9d7bfceae57cf9")
}

// UnpackAddressAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressAlreadyRegistered(address rskAddr)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackAddressAlreadyRegisteredError(raw []byte) (*PegInAddressRegistryContractAddressAlreadyRegistered, error) {
	out := new(PegInAddressRegistryContractAddressAlreadyRegistered)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "AddressAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractBatchTooLarge represents a BatchTooLarge error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractBatchTooLarge struct {
	Requested *big.Int
	Max       *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error BatchTooLarge(uint256 requested, uint256 max)
func PegInAddressRegistryContractBatchTooLargeErrorID() common.Hash {
	return common.HexToHash("0xbb1cb70bcc2ddbaabd8d251943ebb6b48ae10b9f78ce35a6f00b43eac7ae64f1")
}

// UnpackBatchTooLargeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error BatchTooLarge(uint256 requested, uint256 max)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackBatchTooLargeError(raw []byte) (*PegInAddressRegistryContractBatchTooLarge, error) {
	out := new(PegInAddressRegistryContractBatchTooLarge)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "BatchTooLarge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractConfigMinBelowBridge represents a ConfigMinBelowBridge error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractConfigMinBelowBridge struct {
	ProtocolMinSats *big.Int
	BridgeMinSats   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ConfigMinBelowBridge(uint256 protocolMinSats, uint256 bridgeMinSats)
func PegInAddressRegistryContractConfigMinBelowBridgeErrorID() common.Hash {
	return common.HexToHash("0x83a110e5c260a077179361e44fee8a2bc38d94bf3197f8090786a827ddd5175e")
}

// UnpackConfigMinBelowBridgeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ConfigMinBelowBridge(uint256 protocolMinSats, uint256 bridgeMinSats)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackConfigMinBelowBridgeError(raw []byte) (*PegInAddressRegistryContractConfigMinBelowBridge, error) {
	out := new(PegInAddressRegistryContractConfigMinBelowBridge)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "ConfigMinBelowBridge", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractConfigurationsNotSet represents a ConfigurationsNotSet error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractConfigurationsNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ConfigurationsNotSet()
func PegInAddressRegistryContractConfigurationsNotSetErrorID() common.Hash {
	return common.HexToHash("0xb886c87724e229e4303d6c62dd737b4b34c49708459241ceb77de770ce9030fd")
}

// UnpackConfigurationsNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ConfigurationsNotSet()
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackConfigurationsNotSetError(raw []byte) (*PegInAddressRegistryContractConfigurationsNotSet, error) {
	out := new(PegInAddressRegistryContractConfigurationsNotSet)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "ConfigurationsNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractDepositBelowMinimum represents a DepositBelowMinimum error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractDepositBelowMinimum struct {
	Value   *big.Int
	Minimum *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepositBelowMinimum(uint256 value, uint256 minimum)
func PegInAddressRegistryContractDepositBelowMinimumErrorID() common.Hash {
	return common.HexToHash("0xed28d08990578501140aeed4e899fa6c18673086b6de379fd15ef2a9bc2704a6")
}

// UnpackDepositBelowMinimumError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepositBelowMinimum(uint256 value, uint256 minimum)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackDepositBelowMinimumError(raw []byte) (*PegInAddressRegistryContractDepositBelowMinimum, error) {
	out := new(PegInAddressRegistryContractDepositBelowMinimum)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "DepositBelowMinimum", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractDepositNotConfirmed represents a DepositNotConfirmed error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractDepositNotConfirmed struct {
	BtcTxHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepositNotConfirmed(bytes32 btcTxHash)
func PegInAddressRegistryContractDepositNotConfirmedErrorID() common.Hash {
	return common.HexToHash("0x82891d8289b8b9bfe30d65ab3848d8bd2a3deb0a2c04a0cf7fd777a4a23ba310")
}

// UnpackDepositNotConfirmedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepositNotConfirmed(bytes32 btcTxHash)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackDepositNotConfirmedError(raw []byte) (*PegInAddressRegistryContractDepositNotConfirmed, error) {
	out := new(PegInAddressRegistryContractDepositNotConfirmed)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "DepositNotConfirmed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractDepositOutputNotFound represents a DepositOutputNotFound error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractDepositOutputNotFound struct {
	RskAddr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepositOutputNotFound(address rskAddr)
func PegInAddressRegistryContractDepositOutputNotFoundErrorID() common.Hash {
	return common.HexToHash("0xfb8505d7bfa9f884edfcc011ecf7eb6687d546f26b9d3c2a552473c1946b8836")
}

// UnpackDepositOutputNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepositOutputNotFound(address rskAddr)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackDepositOutputNotFoundError(raw []byte) (*PegInAddressRegistryContractDepositOutputNotFound, error) {
	out := new(PegInAddressRegistryContractDepositOutputNotFound)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "DepositOutputNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractPegInContractNotSet represents a PegInContractNotSet error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractPegInContractNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PegInContractNotSet()
func PegInAddressRegistryContractPegInContractNotSetErrorID() common.Hash {
	return common.HexToHash("0x86342334cd81a2843094724d6be926dc543a0726d4f1f220532563338633d274")
}

// UnpackPegInContractNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PegInContractNotSet()
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackPegInContractNotSetError(raw []byte) (*PegInAddressRegistryContractPegInContractNotSet, error) {
	out := new(PegInAddressRegistryContractPegInContractNotSet)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "PegInContractNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryContractUnableToGetConfirmations represents a UnableToGetConfirmations error raised by the PegInAddressRegistryContract contract.
type PegInAddressRegistryContractUnableToGetConfirmations struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnableToGetConfirmations(int256 errorCode)
func PegInAddressRegistryContractUnableToGetConfirmationsErrorID() common.Hash {
	return common.HexToHash("0xd06e366aa562b88cc28592d73c0fea0b2c6f7f7544f39af59c0d16c892dd5ad4")
}

// UnpackUnableToGetConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnableToGetConfirmations(int256 errorCode)
func (pegInAddressRegistryContract *PegInAddressRegistryContract) UnpackUnableToGetConfirmationsError(raw []byte) (*PegInAddressRegistryContractUnableToGetConfirmations, error) {
	out := new(PegInAddressRegistryContractUnableToGetConfirmations)
	if err := pegInAddressRegistryContract.abi.UnpackIntoInterface(out, "UnableToGetConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}
