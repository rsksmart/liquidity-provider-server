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

// EmergencyPauseContractMetaData contains all meta data concerning the EmergencyPauseContract contract.
var EmergencyPauseContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"pauseRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPauseRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]}]",
	ID:  "EmergencyPauseContract",
}

// EmergencyPauseContract is an auto generated Go binding around an Ethereum contract.
type EmergencyPauseContract struct {
	abi abi.ABI
}

// NewEmergencyPauseContract creates a new instance of EmergencyPauseContract.
func NewEmergencyPauseContract() *EmergencyPauseContract {
	parsed, err := EmergencyPauseContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &EmergencyPauseContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *EmergencyPauseContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackPauseRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c82732d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseRegistry() view returns(address)
func (emergencyPauseContract *EmergencyPauseContract) PackPauseRegistry() []byte {
	enc, err := emergencyPauseContract.abi.Pack("pauseRegistry")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseRegistry is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1c82732d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseRegistry() view returns(address)
func (emergencyPauseContract *EmergencyPauseContract) TryPackPauseRegistry() ([]byte, error) {
	return emergencyPauseContract.abi.Pack("pauseRegistry")
}

// UnpackPauseRegistry is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1c82732d.
//
// Solidity: function pauseRegistry() view returns(address)
func (emergencyPauseContract *EmergencyPauseContract) UnpackPauseRegistry(data []byte) (common.Address, error) {
	out, err := emergencyPauseContract.abi.Unpack("pauseRegistry", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (emergencyPauseContract *EmergencyPauseContract) PackPauseStatus() []byte {
	enc, err := emergencyPauseContract.abi.Pack("pauseStatus")
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
func (emergencyPauseContract *EmergencyPauseContract) TryPackPauseStatus() ([]byte, error) {
	return emergencyPauseContract.abi.Pack("pauseStatus")
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
func (emergencyPauseContract *EmergencyPauseContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := emergencyPauseContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// EmergencyPauseContractInitialized represents a Initialized event raised by the EmergencyPauseContract contract.
type EmergencyPauseContractInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const EmergencyPauseContractInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (EmergencyPauseContractInitialized) ContractEventName() string {
	return EmergencyPauseContractInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (emergencyPauseContract *EmergencyPauseContract) UnpackInitializedEvent(log *types.Log) (*EmergencyPauseContractInitialized, error) {
	event := "Initialized"
	if len(log.Topics) == 0 || log.Topics[0] != emergencyPauseContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(EmergencyPauseContractInitialized)
	if len(log.Data) > 0 {
		if err := emergencyPauseContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range emergencyPauseContract.abi.Events[event].Inputs {
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
func (emergencyPauseContract *EmergencyPauseContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], emergencyPauseContract.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return emergencyPauseContract.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], emergencyPauseContract.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return emergencyPauseContract.UnpackNotInitializingError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// EmergencyPauseContractInvalidInitialization represents a InvalidInitialization error raised by the EmergencyPauseContract contract.
type EmergencyPauseContractInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func EmergencyPauseContractInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (emergencyPauseContract *EmergencyPauseContract) UnpackInvalidInitializationError(raw []byte) (*EmergencyPauseContractInvalidInitialization, error) {
	out := new(EmergencyPauseContractInvalidInitialization)
	if err := emergencyPauseContract.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// EmergencyPauseContractNotInitializing represents a NotInitializing error raised by the EmergencyPauseContract contract.
type EmergencyPauseContractNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func EmergencyPauseContractNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (emergencyPauseContract *EmergencyPauseContract) UnpackNotInitializingError(raw []byte) (*EmergencyPauseContractNotInitializing, error) {
	out := new(EmergencyPauseContractNotInitializing)
	if err := emergencyPauseContract.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}
