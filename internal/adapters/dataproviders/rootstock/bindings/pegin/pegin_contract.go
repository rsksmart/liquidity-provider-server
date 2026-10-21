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

// PeginContractMetaData contains all meta data concerning the PeginContract contract.
var PeginContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinPegIn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BalanceDecrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BalanceIncrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
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

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) PackEip712Domain() []byte {
	enc, err := peginContract.abi.Pack("eip712Domain")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) TryPackEip712Domain() ([]byte, error) {
	return peginContract.abi.Pack("eip712Domain")
}

// Eip712DomainOutput serves as a container for the return parameters of contract
// method Eip712Domain.
type Eip712DomainOutput struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}

// UnpackEip712Domain is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (peginContract *PeginContract) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := peginContract.abi.Unpack("eip712Domain", data)
	outstruct := new(Eip712DomainOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = abi.ConvertType(out[3], new(big.Int)).(*big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)
	return *outstruct, nil
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

// PeginContractEIP712DomainChanged represents a EIP712DomainChanged event raised by the PeginContract contract.
type PeginContractEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PeginContractEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (PeginContractEIP712DomainChanged) ContractEventName() string {
	return PeginContractEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (peginContract *PeginContract) UnpackEIP712DomainChangedEvent(log *types.Log) (*PeginContractEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 || log.Topics[0] != peginContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginContractEIP712DomainChanged)
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
