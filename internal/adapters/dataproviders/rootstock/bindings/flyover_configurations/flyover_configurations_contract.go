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

// IFlyoverConfigurationsConfirmationTier is an auto generated low-level Go binding around an user-defined struct.
type IFlyoverConfigurationsConfirmationTier struct {
	MaxAmount     *big.Int
	Confirmations *big.Int
}

// IFlyoverConfigurationsPegConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IFlyoverConfigurationsPegConfiguration struct {
	FixedFee          *big.Int
	PercentageFee     *big.Int
	MinAmount         *big.Int
	MaxAmount         *big.Int
	ConfirmationTiers []IFlyoverConfigurationsConfirmationTier
}

// IFlyoverConfigurationsPegOutConfiguration is an auto generated low-level Go binding around an user-defined struct.
type IFlyoverConfigurationsPegOutConfiguration struct {
	FixedFee          *big.Int
	PercentageFee     *big.Int
	MinAmount         *big.Int
	MaxAmount         *big.Int
	ConfirmationTiers []IFlyoverConfigurationsConfirmationTier
	PenaltyFee        *big.Int
	ClaimWindow       *big.Int
	ClaimWindowBlocks *big.Int
	CallTime          *big.Int
	ExpireTime        *big.Int
	ExpireBlocks      *big.Int
	MaxMinerFee       *big.Int
}

// FlyoverConfigurationsContractMetaData contains all meta data concerning the FlyoverConfigurationsContract contract.
var FlyoverConfigurationsContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"applyChange\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"applyPegOutChange\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"calculatePegInFee\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculatePegOutFee\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInConfiguration\",\"inputs\":[],\"outputs\":[{\"name\":\"configuration\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutConfiguration\",\"inputs\":[],\"outputs\":[{\"name\":\"configuration\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegOutConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"claimWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"claimWindowBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxMinerFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredPegInBtcConfirmations\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequiredPegOutBtcConfirmations\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"queueChange\",\"inputs\":[{\"name\":\"newConfiguration\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"queuePegOutChange\",\"inputs\":[{\"name\":\"newConfiguration\",\"type\":\"tuple\",\"internalType\":\"structIFlyoverConfigurations.PegOutConfiguration\",\"components\":[{\"name\":\"fixedFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"percentageFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmationTiers\",\"type\":\"tuple[]\",\"internalType\":\"structIFlyoverConfigurations.ConfirmationTier[]\",\"components\":[{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"confirmations\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"claimWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"claimWindowBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expireBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxMinerFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	ID:  "FlyoverConfigurationsContract",
}

// FlyoverConfigurationsContract is an auto generated Go binding around an Ethereum contract.
type FlyoverConfigurationsContract struct {
	abi abi.ABI
}

// NewFlyoverConfigurationsContract creates a new instance of FlyoverConfigurationsContract.
func NewFlyoverConfigurationsContract() *FlyoverConfigurationsContract {
	parsed, err := FlyoverConfigurationsContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &FlyoverConfigurationsContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *FlyoverConfigurationsContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackApplyChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e65037e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyChange() returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackApplyChange() []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("applyChange")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4e65037e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyChange() returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackApplyChange() ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("applyChange")
}

// PackApplyPegOutChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7bab8805.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function applyPegOutChange() returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackApplyPegOutChange() []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("applyPegOutChange")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApplyPegOutChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7bab8805.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function applyPegOutChange() returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackApplyPegOutChange() ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("applyPegOutChange")
}

// PackCalculatePegInFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715a177c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackCalculatePegInFee(amount *big.Int) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("calculatePegInFee", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCalculatePegInFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715a177c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackCalculatePegInFee(amount *big.Int) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("calculatePegInFee", amount)
}

// UnpackCalculatePegInFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x715a177c.
//
// Solidity: function calculatePegInFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackCalculatePegInFee(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("calculatePegInFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCalculatePegOutFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x516ba0b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackCalculatePegOutFee(amount *big.Int) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("calculatePegOutFee", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCalculatePegOutFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x516ba0b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackCalculatePegOutFee(amount *big.Int) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("calculatePegOutFee", amount)
}

// UnpackCalculatePegOutFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x516ba0b3.
//
// Solidity: function calculatePegOutFee(uint256 amount) view returns(uint256 fee)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackCalculatePegOutFee(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("calculatePegOutFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegInConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cd5733c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[]) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackGetPegInConfiguration() []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("getPegInConfiguration")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7cd5733c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[]) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackGetPegInConfiguration() ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("getPegInConfiguration")
}

// UnpackGetPegInConfiguration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7cd5733c.
//
// Solidity: function getPegInConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[]) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackGetPegInConfiguration(data []byte) (IFlyoverConfigurationsPegConfiguration, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("getPegInConfiguration", data)
	if err != nil {
		return *new(IFlyoverConfigurationsPegConfiguration), err
	}
	out0 := *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegConfiguration)).(*IFlyoverConfigurationsPegConfiguration)
	return out0, nil
}

// PackGetPegOutConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a8a2124.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256,uint256) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackGetPegOutConfiguration() []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("getPegOutConfiguration")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutConfiguration is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8a8a2124.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256,uint256) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackGetPegOutConfiguration() ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("getPegOutConfiguration")
}

// UnpackGetPegOutConfiguration is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8a8a2124.
//
// Solidity: function getPegOutConfiguration() view returns((uint256,uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256,uint256) configuration)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackGetPegOutConfiguration(data []byte) (IFlyoverConfigurationsPegOutConfiguration, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("getPegOutConfiguration", data)
	if err != nil {
		return *new(IFlyoverConfigurationsPegOutConfiguration), err
	}
	out0 := *abi.ConvertType(out[0], new(IFlyoverConfigurationsPegOutConfiguration)).(*IFlyoverConfigurationsPegOutConfiguration)
	return out0, nil
}

// PackGetRequiredPegInBtcConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7bb29dbe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRequiredPegInBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackGetRequiredPegInBtcConfirmations(amount *big.Int) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("getRequiredPegInBtcConfirmations", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRequiredPegInBtcConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7bb29dbe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRequiredPegInBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackGetRequiredPegInBtcConfirmations(amount *big.Int) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("getRequiredPegInBtcConfirmations", amount)
}

// UnpackGetRequiredPegInBtcConfirmations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7bb29dbe.
//
// Solidity: function getRequiredPegInBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackGetRequiredPegInBtcConfirmations(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("getRequiredPegInBtcConfirmations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRequiredPegOutBtcConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00565cfb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRequiredPegOutBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackGetRequiredPegOutBtcConfirmations(amount *big.Int) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("getRequiredPegOutBtcConfirmations", amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRequiredPegOutBtcConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x00565cfb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRequiredPegOutBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackGetRequiredPegOutBtcConfirmations(amount *big.Int) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("getRequiredPegOutBtcConfirmations", amount)
}

// UnpackGetRequiredPegOutBtcConfirmations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x00565cfb.
//
// Solidity: function getRequiredPegOutBtcConfirmations(uint256 amount) view returns(uint256 confirmations)
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) UnpackGetRequiredPegOutBtcConfirmations(data []byte) (*big.Int, error) {
	out, err := flyoverConfigurationsContract.abi.Unpack("getRequiredPegOutBtcConfirmations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackQueueChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x99173258.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function queueChange((uint256,uint256,uint256,uint256,(uint256,uint256)[]) newConfiguration) returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackQueueChange(newConfiguration IFlyoverConfigurationsPegConfiguration) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("queueChange", newConfiguration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQueueChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x99173258.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function queueChange((uint256,uint256,uint256,uint256,(uint256,uint256)[]) newConfiguration) returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackQueueChange(newConfiguration IFlyoverConfigurationsPegConfiguration) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("queueChange", newConfiguration)
}

// PackQueuePegOutChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x29ae8609.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function queuePegOutChange((uint256,uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256,uint256) newConfiguration) returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) PackQueuePegOutChange(newConfiguration IFlyoverConfigurationsPegOutConfiguration) []byte {
	enc, err := flyoverConfigurationsContract.abi.Pack("queuePegOutChange", newConfiguration)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackQueuePegOutChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x29ae8609.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function queuePegOutChange((uint256,uint256,uint256,uint256,(uint256,uint256)[],uint256,uint256,uint256,uint256,uint256,uint256,uint256) newConfiguration) returns()
func (flyoverConfigurationsContract *FlyoverConfigurationsContract) TryPackQueuePegOutChange(newConfiguration IFlyoverConfigurationsPegOutConfiguration) ([]byte, error) {
	return flyoverConfigurationsContract.abi.Pack("queuePegOutChange", newConfiguration)
}
