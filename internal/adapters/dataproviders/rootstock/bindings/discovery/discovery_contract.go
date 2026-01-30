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

// FlyoverLiquidityProvider is an auto generated low-level Go binding around an user-defined struct.
type FlyoverLiquidityProvider struct {
	Id              *big.Int
	ProviderAddress common.Address
	Status          bool
	ProviderType    uint8
	Name            string
	ApiBaseUrl      string
}

// FlyoverDiscoveryMetaData contains all meta data concerning the FlyoverDiscovery contract.
var FlyoverDiscoveryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getProvider\",\"inputs\":[{\"name\":\"providerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structFlyover.LiquidityProvider\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"providerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProviders\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structFlyover.LiquidityProvider[]\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"providerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProvidersId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperational\",\"inputs\":[{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setProviderStatus\",\"inputs\":[{\"name\":\"providerId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProvider\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ProviderStatusSet\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":true,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderUpdate\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Register\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientCollateral\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidProviderData\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"apiBaseUrl\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"InvalidProviderType\",\"inputs\":[{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"}]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotEOA\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "FlyoverDiscovery",
}

// FlyoverDiscovery is an auto generated Go binding around an Ethereum contract.
type FlyoverDiscovery struct {
	abi abi.ABI
}

// NewFlyoverDiscovery creates a new instance of FlyoverDiscovery.
func NewFlyoverDiscovery() *FlyoverDiscovery {
	parsed, err := FlyoverDiscoveryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &FlyoverDiscovery{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *FlyoverDiscovery) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackGetProvider is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55f21eb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (flyoverDiscovery *FlyoverDiscovery) PackGetProvider(providerAddress common.Address) []byte {
	enc, err := flyoverDiscovery.abi.Pack("getProvider", providerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetProvider is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x55f21eb7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (flyoverDiscovery *FlyoverDiscovery) TryPackGetProvider(providerAddress common.Address) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("getProvider", providerAddress)
}

// UnpackGetProvider is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x55f21eb7.
//
// Solidity: function getProvider(address providerAddress) view returns((uint256,address,bool,uint8,string,string))
func (flyoverDiscovery *FlyoverDiscovery) UnpackGetProvider(data []byte) (FlyoverLiquidityProvider, error) {
	out, err := flyoverDiscovery.abi.Unpack("getProvider", data)
	if err != nil {
		return *new(FlyoverLiquidityProvider), err
	}
	out0 := *abi.ConvertType(out[0], new(FlyoverLiquidityProvider)).(*FlyoverLiquidityProvider)
	return out0, nil
}

// PackGetProviders is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedc922a9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (flyoverDiscovery *FlyoverDiscovery) PackGetProviders() []byte {
	enc, err := flyoverDiscovery.abi.Pack("getProviders")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetProviders is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedc922a9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (flyoverDiscovery *FlyoverDiscovery) TryPackGetProviders() ([]byte, error) {
	return flyoverDiscovery.abi.Pack("getProviders")
}

// UnpackGetProviders is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xedc922a9.
//
// Solidity: function getProviders() view returns((uint256,address,bool,uint8,string,string)[])
func (flyoverDiscovery *FlyoverDiscovery) UnpackGetProviders(data []byte) ([]FlyoverLiquidityProvider, error) {
	out, err := flyoverDiscovery.abi.Unpack("getProviders", data)
	if err != nil {
		return *new([]FlyoverLiquidityProvider), err
	}
	out0 := *abi.ConvertType(out[0], new([]FlyoverLiquidityProvider)).(*[]FlyoverLiquidityProvider)
	return out0, nil
}

// PackGetProvidersId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x122dab09.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getProvidersId() view returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) PackGetProvidersId() []byte {
	enc, err := flyoverDiscovery.abi.Pack("getProvidersId")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetProvidersId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x122dab09.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getProvidersId() view returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) TryPackGetProvidersId() ([]byte, error) {
	return flyoverDiscovery.abi.Pack("getProvidersId")
}

// UnpackGetProvidersId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x122dab09.
//
// Solidity: function getProvidersId() view returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) UnpackGetProvidersId(data []byte) (*big.Int, error) {
	out, err := flyoverDiscovery.abi.Unpack("getProvidersId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackIsOperational is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf50daf0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (flyoverDiscovery *FlyoverDiscovery) PackIsOperational(providerType uint8, addr common.Address) []byte {
	enc, err := flyoverDiscovery.abi.Pack("isOperational", providerType, addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsOperational is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbf50daf0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (flyoverDiscovery *FlyoverDiscovery) TryPackIsOperational(providerType uint8, addr common.Address) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("isOperational", providerType, addr)
}

// UnpackIsOperational is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbf50daf0.
//
// Solidity: function isOperational(uint8 providerType, address addr) view returns(bool)
func (flyoverDiscovery *FlyoverDiscovery) UnpackIsOperational(data []byte) (bool, error) {
	out, err := flyoverDiscovery.abi.Unpack("isOperational", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6da66355.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pause(string reason) returns()
func (flyoverDiscovery *FlyoverDiscovery) PackPause(reason string) []byte {
	enc, err := flyoverDiscovery.abi.Pack("pause", reason)
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
func (flyoverDiscovery *FlyoverDiscovery) TryPackPause(reason string) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("pause", reason)
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (flyoverDiscovery *FlyoverDiscovery) PackPauseStatus() []byte {
	enc, err := flyoverDiscovery.abi.Pack("pauseStatus")
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
func (flyoverDiscovery *FlyoverDiscovery) TryPackPauseStatus() ([]byte, error) {
	return flyoverDiscovery.abi.Pack("pauseStatus")
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
func (flyoverDiscovery *FlyoverDiscovery) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := flyoverDiscovery.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackRegister is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4487ce11.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) PackRegister(name string, apiBaseUrl string, status bool, providerType uint8) []byte {
	enc, err := flyoverDiscovery.abi.Pack("register", name, apiBaseUrl, status, providerType)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegister is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4487ce11.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) TryPackRegister(name string, apiBaseUrl string, status bool, providerType uint8) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("register", name, apiBaseUrl, status, providerType)
}

// UnpackRegister is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4487ce11.
//
// Solidity: function register(string name, string apiBaseUrl, bool status, uint8 providerType) payable returns(uint256)
func (flyoverDiscovery *FlyoverDiscovery) UnpackRegister(data []byte) (*big.Int, error) {
	out, err := flyoverDiscovery.abi.Unpack("register", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetProviderStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72cbf4e8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setProviderStatus(uint256 providerId, bool status) returns()
func (flyoverDiscovery *FlyoverDiscovery) PackSetProviderStatus(providerId *big.Int, status bool) []byte {
	enc, err := flyoverDiscovery.abi.Pack("setProviderStatus", providerId, status)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetProviderStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x72cbf4e8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setProviderStatus(uint256 providerId, bool status) returns()
func (flyoverDiscovery *FlyoverDiscovery) TryPackSetProviderStatus(providerId *big.Int, status bool) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("setProviderStatus", providerId, status)
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpause() returns()
func (flyoverDiscovery *FlyoverDiscovery) PackUnpause() []byte {
	enc, err := flyoverDiscovery.abi.Pack("unpause")
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
func (flyoverDiscovery *FlyoverDiscovery) TryPackUnpause() ([]byte, error) {
	return flyoverDiscovery.abi.Pack("unpause")
}

// PackUpdateProvider is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0220f41d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateProvider(string name, string apiBaseUrl) returns()
func (flyoverDiscovery *FlyoverDiscovery) PackUpdateProvider(name string, apiBaseUrl string) []byte {
	enc, err := flyoverDiscovery.abi.Pack("updateProvider", name, apiBaseUrl)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateProvider is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0220f41d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateProvider(string name, string apiBaseUrl) returns()
func (flyoverDiscovery *FlyoverDiscovery) TryPackUpdateProvider(name string, apiBaseUrl string) ([]byte, error) {
	return flyoverDiscovery.abi.Pack("updateProvider", name, apiBaseUrl)
}

// FlyoverDiscoveryProviderStatusSet represents a ProviderStatusSet event raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryProviderStatusSet struct {
	Id     *big.Int
	Status bool
	Raw    *types.Log // Blockchain specific contextual infos
}

const FlyoverDiscoveryProviderStatusSetEventName = "ProviderStatusSet"

// ContractEventName returns the user-defined event name.
func (FlyoverDiscoveryProviderStatusSet) ContractEventName() string {
	return FlyoverDiscoveryProviderStatusSetEventName
}

// UnpackProviderStatusSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ProviderStatusSet(uint256 indexed id, bool indexed status)
func (flyoverDiscovery *FlyoverDiscovery) UnpackProviderStatusSetEvent(log *types.Log) (*FlyoverDiscoveryProviderStatusSet, error) {
	event := "ProviderStatusSet"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverDiscovery.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverDiscoveryProviderStatusSet)
	if len(log.Data) > 0 {
		if err := flyoverDiscovery.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverDiscovery.abi.Events[event].Inputs {
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

// FlyoverDiscoveryProviderUpdate represents a ProviderUpdate event raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryProviderUpdate struct {
	From       common.Address
	Name       string
	ApiBaseUrl string
	Raw        *types.Log // Blockchain specific contextual infos
}

const FlyoverDiscoveryProviderUpdateEventName = "ProviderUpdate"

// ContractEventName returns the user-defined event name.
func (FlyoverDiscoveryProviderUpdate) ContractEventName() string {
	return FlyoverDiscoveryProviderUpdateEventName
}

// UnpackProviderUpdateEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ProviderUpdate(address indexed from, string name, string apiBaseUrl)
func (flyoverDiscovery *FlyoverDiscovery) UnpackProviderUpdateEvent(log *types.Log) (*FlyoverDiscoveryProviderUpdate, error) {
	event := "ProviderUpdate"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverDiscovery.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverDiscoveryProviderUpdate)
	if len(log.Data) > 0 {
		if err := flyoverDiscovery.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverDiscovery.abi.Events[event].Inputs {
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

// FlyoverDiscoveryRegister represents a Register event raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryRegister struct {
	Id     *big.Int
	From   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const FlyoverDiscoveryRegisterEventName = "Register"

// ContractEventName returns the user-defined event name.
func (FlyoverDiscoveryRegister) ContractEventName() string {
	return FlyoverDiscoveryRegisterEventName
}

// UnpackRegisterEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Register(uint256 indexed id, address indexed from, uint256 indexed amount)
func (flyoverDiscovery *FlyoverDiscovery) UnpackRegisterEvent(log *types.Log) (*FlyoverDiscoveryRegister, error) {
	event := "Register"
	if len(log.Topics) == 0 || log.Topics[0] != flyoverDiscovery.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(FlyoverDiscoveryRegister)
	if len(log.Data) > 0 {
		if err := flyoverDiscovery.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range flyoverDiscovery.abi.Events[event].Inputs {
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
func (flyoverDiscovery *FlyoverDiscovery) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["AlreadyRegistered"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["InsufficientCollateral"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackInsufficientCollateralError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["InvalidProviderData"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackInvalidProviderDataError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["InvalidProviderType"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackInvalidProviderTypeError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["NotAuthorized"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackNotAuthorizedError(raw[4:])
	}
	if bytes.Equal(raw[:4], flyoverDiscovery.abi.Errors["NotEOA"].ID.Bytes()[:4]) {
		return flyoverDiscovery.UnpackNotEOAError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// FlyoverDiscoveryAlreadyRegistered represents a AlreadyRegistered error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryAlreadyRegistered struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyRegistered(address from)
func FlyoverDiscoveryAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x45ed80e9399c87887ea54f615514a1e3dde31e9b6c027ddceb4ffd503b70e428")
}

// UnpackAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyRegistered(address from)
func (flyoverDiscovery *FlyoverDiscovery) UnpackAlreadyRegisteredError(raw []byte) (*FlyoverDiscoveryAlreadyRegistered, error) {
	out := new(FlyoverDiscoveryAlreadyRegistered)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "AlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverDiscoveryInsufficientCollateral represents a InsufficientCollateral error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryInsufficientCollateral struct {
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientCollateral(uint256 amount)
func FlyoverDiscoveryInsufficientCollateralErrorID() common.Hash {
	return common.HexToHash("0x2b3bc98557cfd9172d79165d096e2f1f67c00ee208ef0ac5df06614762a558fb")
}

// UnpackInsufficientCollateralError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientCollateral(uint256 amount)
func (flyoverDiscovery *FlyoverDiscovery) UnpackInsufficientCollateralError(raw []byte) (*FlyoverDiscoveryInsufficientCollateral, error) {
	out := new(FlyoverDiscoveryInsufficientCollateral)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "InsufficientCollateral", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverDiscoveryInvalidProviderData represents a InvalidProviderData error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryInvalidProviderData struct {
	Name       string
	ApiBaseUrl string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProviderData(string name, string apiBaseUrl)
func FlyoverDiscoveryInvalidProviderDataErrorID() common.Hash {
	return common.HexToHash("0x98c8014f7ee8c66171df31ea3333fea87d47597b799c865862100f29e9a069d8")
}

// UnpackInvalidProviderDataError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProviderData(string name, string apiBaseUrl)
func (flyoverDiscovery *FlyoverDiscovery) UnpackInvalidProviderDataError(raw []byte) (*FlyoverDiscoveryInvalidProviderData, error) {
	out := new(FlyoverDiscoveryInvalidProviderData)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "InvalidProviderData", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverDiscoveryInvalidProviderType represents a InvalidProviderType error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryInvalidProviderType struct {
	ProviderType uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidProviderType(uint8 providerType)
func FlyoverDiscoveryInvalidProviderTypeErrorID() common.Hash {
	return common.HexToHash("0x36ecd6f0622f170b4f2f53cb59b015a30d326ea6dcce1e3cf6a480eb7bd0bff0")
}

// UnpackInvalidProviderTypeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidProviderType(uint8 providerType)
func (flyoverDiscovery *FlyoverDiscovery) UnpackInvalidProviderTypeError(raw []byte) (*FlyoverDiscoveryInvalidProviderType, error) {
	out := new(FlyoverDiscoveryInvalidProviderType)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "InvalidProviderType", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverDiscoveryNotAuthorized represents a NotAuthorized error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryNotAuthorized struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAuthorized(address from)
func FlyoverDiscoveryNotAuthorizedErrorID() common.Hash {
	return common.HexToHash("0x4a0bfec1fa3ea832f47a765e92fe13fe78ddadee056368ea6e89037b3cf70498")
}

// UnpackNotAuthorizedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAuthorized(address from)
func (flyoverDiscovery *FlyoverDiscovery) UnpackNotAuthorizedError(raw []byte) (*FlyoverDiscoveryNotAuthorized, error) {
	out := new(FlyoverDiscoveryNotAuthorized)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "NotAuthorized", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FlyoverDiscoveryNotEOA represents a NotEOA error raised by the FlyoverDiscovery contract.
type FlyoverDiscoveryNotEOA struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEOA(address from)
func FlyoverDiscoveryNotEOAErrorID() common.Hash {
	return common.HexToHash("0x77817ac35fb75c204a372975a1886a5b3d89550bb4560d9ba4dc2d5aec75cca1")
}

// UnpackNotEOAError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEOA(address from)
func (flyoverDiscovery *FlyoverDiscovery) UnpackNotEOAError(raw []byte) (*FlyoverDiscoveryNotEOA, error) {
	out := new(FlyoverDiscoveryNotEOA)
	if err := flyoverDiscovery.abi.UnpackIntoInterface(out, "NotEOA", raw); err != nil {
		return nil, err
	}
	return out, nil
}
