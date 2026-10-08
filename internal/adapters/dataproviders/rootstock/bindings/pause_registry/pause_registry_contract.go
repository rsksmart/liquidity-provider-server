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

// PauseRegistryContractMetaData contains all meta data concerning the PauseRegistryContract contract.
var PauseRegistryContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"computePauseOverlap\",\"inputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computePauseOverlapBlocks\",\"inputs\":[{\"name\":\"startBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hardPauses\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"startTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"startBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hardPausesCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseLevel\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPauseRegistry.PauseLevel\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setPauseLevel\",\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\",\"internalType\":\"enumIPauseRegistry.PauseLevel\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	ID:  "PauseRegistryContract",
}

// PauseRegistryContract is an auto generated Go binding around an Ethereum contract.
type PauseRegistryContract struct {
	abi abi.ABI
}

// NewPauseRegistryContract creates a new instance of PauseRegistryContract.
func NewPauseRegistryContract() *PauseRegistryContract {
	parsed, err := PauseRegistryContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PauseRegistryContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PauseRegistryContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackComputePauseOverlap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x69c62b50.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function computePauseOverlap(uint256 startTimestamp, uint256 endTimestamp) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) PackComputePauseOverlap(startTimestamp *big.Int, endTimestamp *big.Int) []byte {
	enc, err := pauseRegistryContract.abi.Pack("computePauseOverlap", startTimestamp, endTimestamp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackComputePauseOverlap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x69c62b50.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function computePauseOverlap(uint256 startTimestamp, uint256 endTimestamp) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) TryPackComputePauseOverlap(startTimestamp *big.Int, endTimestamp *big.Int) ([]byte, error) {
	return pauseRegistryContract.abi.Pack("computePauseOverlap", startTimestamp, endTimestamp)
}

// UnpackComputePauseOverlap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x69c62b50.
//
// Solidity: function computePauseOverlap(uint256 startTimestamp, uint256 endTimestamp) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) UnpackComputePauseOverlap(data []byte) (*big.Int, error) {
	out, err := pauseRegistryContract.abi.Unpack("computePauseOverlap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackComputePauseOverlapBlocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5b19463d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function computePauseOverlapBlocks(uint256 startBlock, uint256 endBlock) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) PackComputePauseOverlapBlocks(startBlock *big.Int, endBlock *big.Int) []byte {
	enc, err := pauseRegistryContract.abi.Pack("computePauseOverlapBlocks", startBlock, endBlock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackComputePauseOverlapBlocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5b19463d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function computePauseOverlapBlocks(uint256 startBlock, uint256 endBlock) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) TryPackComputePauseOverlapBlocks(startBlock *big.Int, endBlock *big.Int) ([]byte, error) {
	return pauseRegistryContract.abi.Pack("computePauseOverlapBlocks", startBlock, endBlock)
}

// UnpackComputePauseOverlapBlocks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5b19463d.
//
// Solidity: function computePauseOverlapBlocks(uint256 startBlock, uint256 endBlock) view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) UnpackComputePauseOverlapBlocks(data []byte) (*big.Int, error) {
	out, err := pauseRegistryContract.abi.Unpack("computePauseOverlapBlocks", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHardPauses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x388d4f01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hardPauses(uint256 index) view returns(uint64 startTimestamp, uint64 endTimestamp, uint64 startBlock, uint64 endBlock)
func (pauseRegistryContract *PauseRegistryContract) PackHardPauses(index *big.Int) []byte {
	enc, err := pauseRegistryContract.abi.Pack("hardPauses", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHardPauses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x388d4f01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hardPauses(uint256 index) view returns(uint64 startTimestamp, uint64 endTimestamp, uint64 startBlock, uint64 endBlock)
func (pauseRegistryContract *PauseRegistryContract) TryPackHardPauses(index *big.Int) ([]byte, error) {
	return pauseRegistryContract.abi.Pack("hardPauses", index)
}

// HardPausesOutput serves as a container for the return parameters of contract
// method HardPauses.
type HardPausesOutput struct {
	StartTimestamp uint64
	EndTimestamp   uint64
	StartBlock     uint64
	EndBlock       uint64
}

// UnpackHardPauses is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x388d4f01.
//
// Solidity: function hardPauses(uint256 index) view returns(uint64 startTimestamp, uint64 endTimestamp, uint64 startBlock, uint64 endBlock)
func (pauseRegistryContract *PauseRegistryContract) UnpackHardPauses(data []byte) (HardPausesOutput, error) {
	out, err := pauseRegistryContract.abi.Unpack("hardPauses", data)
	outstruct := new(HardPausesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.StartTimestamp = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.EndTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.StartBlock = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.EndBlock = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackHardPausesCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37a8ce64.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hardPausesCount() view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) PackHardPausesCount() []byte {
	enc, err := pauseRegistryContract.abi.Pack("hardPausesCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHardPausesCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x37a8ce64.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hardPausesCount() view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) TryPackHardPausesCount() ([]byte, error) {
	return pauseRegistryContract.abi.Pack("hardPausesCount")
}

// UnpackHardPausesCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x37a8ce64.
//
// Solidity: function hardPausesCount() view returns(uint256)
func (pauseRegistryContract *PauseRegistryContract) UnpackHardPausesCount(data []byte) (*big.Int, error) {
	out, err := pauseRegistryContract.abi.Unpack("hardPausesCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackPauseLevel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe61399b8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseLevel() view returns(uint8)
func (pauseRegistryContract *PauseRegistryContract) PackPauseLevel() []byte {
	enc, err := pauseRegistryContract.abi.Pack("pauseLevel")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPauseLevel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe61399b8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pauseLevel() view returns(uint8)
func (pauseRegistryContract *PauseRegistryContract) TryPackPauseLevel() ([]byte, error) {
	return pauseRegistryContract.abi.Pack("pauseLevel")
}

// UnpackPauseLevel is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe61399b8.
//
// Solidity: function pauseLevel() view returns(uint8)
func (pauseRegistryContract *PauseRegistryContract) UnpackPauseLevel(data []byte) (uint8, error) {
	out, err := pauseRegistryContract.abi.Unpack("pauseLevel", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (pauseRegistryContract *PauseRegistryContract) PackPauseStatus() []byte {
	enc, err := pauseRegistryContract.abi.Pack("pauseStatus")
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
func (pauseRegistryContract *PauseRegistryContract) TryPackPauseStatus() ([]byte, error) {
	return pauseRegistryContract.abi.Pack("pauseStatus")
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
func (pauseRegistryContract *PauseRegistryContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := pauseRegistryContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function paused() view returns(bool)
func (pauseRegistryContract *PauseRegistryContract) PackPaused() []byte {
	enc, err := pauseRegistryContract.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function paused() view returns(bool)
func (pauseRegistryContract *PauseRegistryContract) TryPackPaused() ([]byte, error) {
	return pauseRegistryContract.abi.Pack("paused")
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (pauseRegistryContract *PauseRegistryContract) UnpackPaused(data []byte) (bool, error) {
	out, err := pauseRegistryContract.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSetPauseLevel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c5d487.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setPauseLevel(uint8 level, string reason) returns()
func (pauseRegistryContract *PauseRegistryContract) PackSetPauseLevel(level uint8, reason string) []byte {
	enc, err := pauseRegistryContract.abi.Pack("setPauseLevel", level, reason)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetPauseLevel is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf8c5d487.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setPauseLevel(uint8 level, string reason) returns()
func (pauseRegistryContract *PauseRegistryContract) TryPackSetPauseLevel(level uint8, reason string) ([]byte, error) {
	return pauseRegistryContract.abi.Pack("setPauseLevel", level, reason)
}
