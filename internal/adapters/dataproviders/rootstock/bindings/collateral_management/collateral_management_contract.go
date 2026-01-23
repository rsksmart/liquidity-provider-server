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

// CollateralManagementContractMetaData contains all meta data concerning the CollateralManagementContract contract.
var CollateralManagementContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addPegInCollateral\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addPegInCollateralTo\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addPegOutCollateral\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addPegOutCollateralTo\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getMinCollateral\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInCollateral\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutCollateral\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPenalties\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getResignDelayInBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getResignationBlock\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRewardPercentage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRewards\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isCollateralSufficient\",\"inputs\":[{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"providerType\",\"type\":\"uint8\",\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resign\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashPegInCollateral\",\"inputs\":[{\"name\":\"punisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegInQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fedBtcAddress\",\"type\":\"bytes20\",\"internalType\":\"bytes20\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"liquidityProviderRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"gasLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"timeForDeposit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"callTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"callOnRegister\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashPegOutCollateral\",\"inputs\":[{\"name\":\"punisher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"productFeeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawRewards\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"PegInCollateralAdded\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutCollateralAdded\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Penalized\",\"inputs\":[{\"name\":\"liquidityProvider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"punisher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"collateralType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumFlyover.ProviderType\"},{\"name\":\"penalty\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reward\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Resigned\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsWithdrawn\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawCollateral\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyResigned\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotResigned\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NothingToWithdraw\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ResignationDelayNotMet\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"resignationBlockNum\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"resignDelayInBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"WithdrawalFailed\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "CollateralManagementContract",
}

// CollateralManagementContract is an auto generated Go binding around an Ethereum contract.
type CollateralManagementContract struct {
	abi abi.ABI
}

// NewCollateralManagementContract creates a new instance of CollateralManagementContract.
func NewCollateralManagementContract() *CollateralManagementContract {
	parsed, err := CollateralManagementContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &CollateralManagementContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *CollateralManagementContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde567d6d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addPegInCollateral() payable returns()
func (collateralManagementContract *CollateralManagementContract) PackAddPegInCollateral() []byte {
	enc, err := collateralManagementContract.abi.Pack("addPegInCollateral")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xde567d6d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addPegInCollateral() payable returns()
func (collateralManagementContract *CollateralManagementContract) TryPackAddPegInCollateral() ([]byte, error) {
	return collateralManagementContract.abi.Pack("addPegInCollateral")
}

// PackAddPegInCollateralTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x83fe87f9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addPegInCollateralTo(address addr) payable returns()
func (collateralManagementContract *CollateralManagementContract) PackAddPegInCollateralTo(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("addPegInCollateralTo", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddPegInCollateralTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x83fe87f9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addPegInCollateralTo(address addr) payable returns()
func (collateralManagementContract *CollateralManagementContract) TryPackAddPegInCollateralTo(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("addPegInCollateralTo", addr)
}

// PackAddPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b2318d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addPegOutCollateral() payable returns()
func (collateralManagementContract *CollateralManagementContract) PackAddPegOutCollateral() []byte {
	enc, err := collateralManagementContract.abi.Pack("addPegOutCollateral")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x52b2318d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addPegOutCollateral() payable returns()
func (collateralManagementContract *CollateralManagementContract) TryPackAddPegOutCollateral() ([]byte, error) {
	return collateralManagementContract.abi.Pack("addPegOutCollateral")
}

// PackAddPegOutCollateralTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ee5b1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addPegOutCollateralTo(address addr) payable returns()
func (collateralManagementContract *CollateralManagementContract) PackAddPegOutCollateralTo(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("addPegOutCollateralTo", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddPegOutCollateralTo is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ee5b1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addPegOutCollateralTo(address addr) payable returns()
func (collateralManagementContract *CollateralManagementContract) TryPackAddPegOutCollateralTo(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("addPegOutCollateralTo", addr)
}

// PackGetMinCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe830b690.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetMinCollateral() []byte {
	enc, err := collateralManagementContract.abi.Pack("getMinCollateral")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMinCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe830b690.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetMinCollateral() ([]byte, error) {
	return collateralManagementContract.abi.Pack("getMinCollateral")
}

// UnpackGetMinCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe830b690.
//
// Solidity: function getMinCollateral() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetMinCollateral(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getMinCollateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x003c3317.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetPegInCollateral(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("getPegInCollateral", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x003c3317.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetPegInCollateral(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("getPegInCollateral", addr)
}

// UnpackGetPegInCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x003c3317.
//
// Solidity: function getPegInCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetPegInCollateral(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getPegInCollateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82b90e93.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetPegOutCollateral(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("getPegOutCollateral", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82b90e93.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetPegOutCollateral(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("getPegOutCollateral", addr)
}

// UnpackGetPegOutCollateral is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x82b90e93.
//
// Solidity: function getPegOutCollateral(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetPegOutCollateral(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getPegOutCollateral", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPenalties is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6ef2a38.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPenalties() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetPenalties() []byte {
	enc, err := collateralManagementContract.abi.Pack("getPenalties")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPenalties is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6ef2a38.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPenalties() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetPenalties() ([]byte, error) {
	return collateralManagementContract.abi.Pack("getPenalties")
}

// UnpackGetPenalties is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe6ef2a38.
//
// Solidity: function getPenalties() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetPenalties(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getPenalties", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetResignDelayInBlocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27887ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetResignDelayInBlocks() []byte {
	enc, err := collateralManagementContract.abi.Pack("getResignDelayInBlocks")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetResignDelayInBlocks is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x27887ffc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetResignDelayInBlocks() ([]byte, error) {
	return collateralManagementContract.abi.Pack("getResignDelayInBlocks")
}

// UnpackGetResignDelayInBlocks is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x27887ffc.
//
// Solidity: function getResignDelayInBlocks() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetResignDelayInBlocks(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getResignDelayInBlocks", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetResignationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd36933d3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetResignationBlock(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("getResignationBlock", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetResignationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd36933d3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetResignationBlock(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("getResignationBlock", addr)
}

// UnpackGetResignationBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd36933d3.
//
// Solidity: function getResignationBlock(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetResignationBlock(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getResignationBlock", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRewardPercentage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7213163.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetRewardPercentage() []byte {
	enc, err := collateralManagementContract.abi.Pack("getRewardPercentage")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRewardPercentage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7213163.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetRewardPercentage() ([]byte, error) {
	return collateralManagementContract.abi.Pack("getRewardPercentage")
}

// UnpackGetRewardPercentage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc7213163.
//
// Solidity: function getRewardPercentage() view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetRewardPercentage(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getRewardPercentage", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRewards is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ee54f7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) PackGetRewards(addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("getRewards", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRewards is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x79ee54f7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) TryPackGetRewards(addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("getRewards", addr)
}

// UnpackGetRewards is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x79ee54f7.
//
// Solidity: function getRewards(address addr) view returns(uint256)
func (collateralManagementContract *CollateralManagementContract) UnpackGetRewards(data []byte) (*big.Int, error) {
	out, err := collateralManagementContract.abi.Unpack("getRewards", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackIsCollateralSufficient is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x718c5aa8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) PackIsCollateralSufficient(providerType uint8, addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("isCollateralSufficient", providerType, addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsCollateralSufficient is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x718c5aa8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) TryPackIsCollateralSufficient(providerType uint8, addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("isCollateralSufficient", providerType, addr)
}

// UnpackIsCollateralSufficient is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x718c5aa8.
//
// Solidity: function isCollateralSufficient(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) UnpackIsCollateralSufficient(data []byte) (bool, error) {
	out, err := collateralManagementContract.abi.Unpack("isCollateralSufficient", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x900daa73.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) PackIsRegistered(providerType uint8, addr common.Address) []byte {
	enc, err := collateralManagementContract.abi.Pack("isRegistered", providerType, addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x900daa73.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) TryPackIsRegistered(providerType uint8, addr common.Address) ([]byte, error) {
	return collateralManagementContract.abi.Pack("isRegistered", providerType, addr)
}

// UnpackIsRegistered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x900daa73.
//
// Solidity: function isRegistered(uint8 providerType, address addr) view returns(bool)
func (collateralManagementContract *CollateralManagementContract) UnpackIsRegistered(data []byte) (bool, error) {
	out, err := collateralManagementContract.abi.Unpack("isRegistered", data)
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
func (collateralManagementContract *CollateralManagementContract) PackPause(reason string) []byte {
	enc, err := collateralManagementContract.abi.Pack("pause", reason)
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
func (collateralManagementContract *CollateralManagementContract) TryPackPause(reason string) ([]byte, error) {
	return collateralManagementContract.abi.Pack("pause", reason)
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (collateralManagementContract *CollateralManagementContract) PackPauseStatus() []byte {
	enc, err := collateralManagementContract.abi.Pack("pauseStatus")
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
func (collateralManagementContract *CollateralManagementContract) TryPackPauseStatus() ([]byte, error) {
	return collateralManagementContract.abi.Pack("pauseStatus")
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
func (collateralManagementContract *CollateralManagementContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := collateralManagementContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackResign is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x69652fcf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function resign() returns()
func (collateralManagementContract *CollateralManagementContract) PackResign() []byte {
	enc, err := collateralManagementContract.abi.Pack("resign")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackResign is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x69652fcf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function resign() returns()
func (collateralManagementContract *CollateralManagementContract) TryPackResign() ([]byte, error) {
	return collateralManagementContract.abi.Pack("resign")
}

// PackSlashPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e4de194.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function slashPegInCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (collateralManagementContract *CollateralManagementContract) PackSlashPegInCollateral(punisher common.Address, quote QuotesPegInQuote, quoteHash [32]byte) []byte {
	enc, err := collateralManagementContract.abi.Pack("slashPegInCollateral", punisher, quote, quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSlashPegInCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3e4de194.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function slashPegInCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,bytes20,address,address,address,address,int64,uint32,uint32,uint32,uint32,uint16,bool,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (collateralManagementContract *CollateralManagementContract) TryPackSlashPegInCollateral(punisher common.Address, quote QuotesPegInQuote, quoteHash [32]byte) ([]byte, error) {
	return collateralManagementContract.abi.Pack("slashPegInCollateral", punisher, quote, quoteHash)
}

// PackSlashPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f6fcee8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function slashPegOutCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (collateralManagementContract *CollateralManagementContract) PackSlashPegOutCollateral(punisher common.Address, quote QuotesPegOutQuote, quoteHash [32]byte) []byte {
	enc, err := collateralManagementContract.abi.Pack("slashPegOutCollateral", punisher, quote, quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSlashPegOutCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f6fcee8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function slashPegOutCollateral(address punisher, (uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes32 quoteHash) returns()
func (collateralManagementContract *CollateralManagementContract) TryPackSlashPegOutCollateral(punisher common.Address, quote QuotesPegOutQuote, quoteHash [32]byte) ([]byte, error) {
	return collateralManagementContract.abi.Pack("slashPegOutCollateral", punisher, quote, quoteHash)
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unpause() returns()
func (collateralManagementContract *CollateralManagementContract) PackUnpause() []byte {
	enc, err := collateralManagementContract.abi.Pack("unpause")
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
func (collateralManagementContract *CollateralManagementContract) TryPackUnpause() ([]byte, error) {
	return collateralManagementContract.abi.Pack("unpause")
}

// PackWithdrawCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59c153be.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawCollateral() returns()
func (collateralManagementContract *CollateralManagementContract) PackWithdrawCollateral() []byte {
	enc, err := collateralManagementContract.abi.Pack("withdrawCollateral")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawCollateral is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59c153be.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawCollateral() returns()
func (collateralManagementContract *CollateralManagementContract) TryPackWithdrawCollateral() ([]byte, error) {
	return collateralManagementContract.abi.Pack("withdrawCollateral")
}

// PackWithdrawRewards is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7b8981c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdrawRewards() returns()
func (collateralManagementContract *CollateralManagementContract) PackWithdrawRewards() []byte {
	enc, err := collateralManagementContract.abi.Pack("withdrawRewards")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdrawRewards is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc7b8981c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdrawRewards() returns()
func (collateralManagementContract *CollateralManagementContract) TryPackWithdrawRewards() ([]byte, error) {
	return collateralManagementContract.abi.Pack("withdrawRewards")
}

// CollateralManagementContractPegInCollateralAdded represents a PegInCollateralAdded event raised by the CollateralManagementContract contract.
type CollateralManagementContractPegInCollateralAdded struct {
	Addr   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractPegInCollateralAddedEventName = "PegInCollateralAdded"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractPegInCollateralAdded) ContractEventName() string {
	return CollateralManagementContractPegInCollateralAddedEventName
}

// UnpackPegInCollateralAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInCollateralAdded(address indexed addr, uint256 indexed amount)
func (collateralManagementContract *CollateralManagementContract) UnpackPegInCollateralAddedEvent(log *types.Log) (*CollateralManagementContractPegInCollateralAdded, error) {
	event := "PegInCollateralAdded"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractPegInCollateralAdded)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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

// CollateralManagementContractPegOutCollateralAdded represents a PegOutCollateralAdded event raised by the CollateralManagementContract contract.
type CollateralManagementContractPegOutCollateralAdded struct {
	Addr   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractPegOutCollateralAddedEventName = "PegOutCollateralAdded"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractPegOutCollateralAdded) ContractEventName() string {
	return CollateralManagementContractPegOutCollateralAddedEventName
}

// UnpackPegOutCollateralAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutCollateralAdded(address indexed addr, uint256 indexed amount)
func (collateralManagementContract *CollateralManagementContract) UnpackPegOutCollateralAddedEvent(log *types.Log) (*CollateralManagementContractPegOutCollateralAdded, error) {
	event := "PegOutCollateralAdded"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractPegOutCollateralAdded)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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

// CollateralManagementContractPenalized represents a Penalized event raised by the CollateralManagementContract contract.
type CollateralManagementContractPenalized struct {
	LiquidityProvider common.Address
	Punisher          common.Address
	QuoteHash         [32]byte
	CollateralType    uint8
	Penalty           *big.Int
	Reward            *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractPenalizedEventName = "Penalized"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractPenalized) ContractEventName() string {
	return CollateralManagementContractPenalizedEventName
}

// UnpackPenalizedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Penalized(address indexed liquidityProvider, address indexed punisher, bytes32 indexed quoteHash, uint8 collateralType, uint256 penalty, uint256 reward)
func (collateralManagementContract *CollateralManagementContract) UnpackPenalizedEvent(log *types.Log) (*CollateralManagementContractPenalized, error) {
	event := "Penalized"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractPenalized)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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

// CollateralManagementContractResigned represents a Resigned event raised by the CollateralManagementContract contract.
type CollateralManagementContractResigned struct {
	Addr common.Address
	Raw  *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractResignedEventName = "Resigned"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractResigned) ContractEventName() string {
	return CollateralManagementContractResignedEventName
}

// UnpackResignedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Resigned(address indexed addr)
func (collateralManagementContract *CollateralManagementContract) UnpackResignedEvent(log *types.Log) (*CollateralManagementContractResigned, error) {
	event := "Resigned"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractResigned)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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

// CollateralManagementContractRewardsWithdrawn represents a RewardsWithdrawn event raised by the CollateralManagementContract contract.
type CollateralManagementContractRewardsWithdrawn struct {
	Addr   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractRewardsWithdrawnEventName = "RewardsWithdrawn"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractRewardsWithdrawn) ContractEventName() string {
	return CollateralManagementContractRewardsWithdrawnEventName
}

// UnpackRewardsWithdrawnEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RewardsWithdrawn(address indexed addr, uint256 indexed amount)
func (collateralManagementContract *CollateralManagementContract) UnpackRewardsWithdrawnEvent(log *types.Log) (*CollateralManagementContractRewardsWithdrawn, error) {
	event := "RewardsWithdrawn"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractRewardsWithdrawn)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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

// CollateralManagementContractWithdrawCollateral represents a WithdrawCollateral event raised by the CollateralManagementContract contract.
type CollateralManagementContractWithdrawCollateral struct {
	Addr   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const CollateralManagementContractWithdrawCollateralEventName = "WithdrawCollateral"

// ContractEventName returns the user-defined event name.
func (CollateralManagementContractWithdrawCollateral) ContractEventName() string {
	return CollateralManagementContractWithdrawCollateralEventName
}

// UnpackWithdrawCollateralEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WithdrawCollateral(address indexed addr, uint256 indexed amount)
func (collateralManagementContract *CollateralManagementContract) UnpackWithdrawCollateralEvent(log *types.Log) (*CollateralManagementContractWithdrawCollateral, error) {
	event := "WithdrawCollateral"
	if len(log.Topics) == 0 || log.Topics[0] != collateralManagementContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(CollateralManagementContractWithdrawCollateral)
	if len(log.Data) > 0 {
		if err := collateralManagementContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range collateralManagementContract.abi.Events[event].Inputs {
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
func (collateralManagementContract *CollateralManagementContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], collateralManagementContract.abi.Errors["AlreadyResigned"].ID.Bytes()[:4]) {
		return collateralManagementContract.UnpackAlreadyResignedError(raw[4:])
	}
	if bytes.Equal(raw[:4], collateralManagementContract.abi.Errors["NotResigned"].ID.Bytes()[:4]) {
		return collateralManagementContract.UnpackNotResignedError(raw[4:])
	}
	if bytes.Equal(raw[:4], collateralManagementContract.abi.Errors["NothingToWithdraw"].ID.Bytes()[:4]) {
		return collateralManagementContract.UnpackNothingToWithdrawError(raw[4:])
	}
	if bytes.Equal(raw[:4], collateralManagementContract.abi.Errors["ResignationDelayNotMet"].ID.Bytes()[:4]) {
		return collateralManagementContract.UnpackResignationDelayNotMetError(raw[4:])
	}
	if bytes.Equal(raw[:4], collateralManagementContract.abi.Errors["WithdrawalFailed"].ID.Bytes()[:4]) {
		return collateralManagementContract.UnpackWithdrawalFailedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// CollateralManagementContractAlreadyResigned represents a AlreadyResigned error raised by the CollateralManagementContract contract.
type CollateralManagementContractAlreadyResigned struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyResigned(address from)
func CollateralManagementContractAlreadyResignedErrorID() common.Hash {
	return common.HexToHash("0x742cab66bae35d9196c31064ad2a9361adc32e5a5c8b1b46917e3efe3e4ced66")
}

// UnpackAlreadyResignedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyResigned(address from)
func (collateralManagementContract *CollateralManagementContract) UnpackAlreadyResignedError(raw []byte) (*CollateralManagementContractAlreadyResigned, error) {
	out := new(CollateralManagementContractAlreadyResigned)
	if err := collateralManagementContract.abi.UnpackIntoInterface(out, "AlreadyResigned", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CollateralManagementContractNotResigned represents a NotResigned error raised by the CollateralManagementContract contract.
type CollateralManagementContractNotResigned struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotResigned(address from)
func CollateralManagementContractNotResignedErrorID() common.Hash {
	return common.HexToHash("0x977254571e2eb4cdd907de3701759fa3386f94c2248e8793202bc75baba41ccf")
}

// UnpackNotResignedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotResigned(address from)
func (collateralManagementContract *CollateralManagementContract) UnpackNotResignedError(raw []byte) (*CollateralManagementContractNotResigned, error) {
	out := new(CollateralManagementContractNotResigned)
	if err := collateralManagementContract.abi.UnpackIntoInterface(out, "NotResigned", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CollateralManagementContractNothingToWithdraw represents a NothingToWithdraw error raised by the CollateralManagementContract contract.
type CollateralManagementContractNothingToWithdraw struct {
	From common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NothingToWithdraw(address from)
func CollateralManagementContractNothingToWithdrawErrorID() common.Hash {
	return common.HexToHash("0xdc69dc16e4a8405b44139edda9a06ee66b303337358b6f4633533d65de15b7f4")
}

// UnpackNothingToWithdrawError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NothingToWithdraw(address from)
func (collateralManagementContract *CollateralManagementContract) UnpackNothingToWithdrawError(raw []byte) (*CollateralManagementContractNothingToWithdraw, error) {
	out := new(CollateralManagementContractNothingToWithdraw)
	if err := collateralManagementContract.abi.UnpackIntoInterface(out, "NothingToWithdraw", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CollateralManagementContractResignationDelayNotMet represents a ResignationDelayNotMet error raised by the CollateralManagementContract contract.
type CollateralManagementContractResignationDelayNotMet struct {
	From                common.Address
	ResignationBlockNum *big.Int
	ResignDelayInBlocks *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ResignationDelayNotMet(address from, uint256 resignationBlockNum, uint256 resignDelayInBlocks)
func CollateralManagementContractResignationDelayNotMetErrorID() common.Hash {
	return common.HexToHash("0xf6cf333563a32441fbc8282a0448cf3475ae2def39a6a786fb688587bff1fb72")
}

// UnpackResignationDelayNotMetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ResignationDelayNotMet(address from, uint256 resignationBlockNum, uint256 resignDelayInBlocks)
func (collateralManagementContract *CollateralManagementContract) UnpackResignationDelayNotMetError(raw []byte) (*CollateralManagementContractResignationDelayNotMet, error) {
	out := new(CollateralManagementContractResignationDelayNotMet)
	if err := collateralManagementContract.abi.UnpackIntoInterface(out, "ResignationDelayNotMet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// CollateralManagementContractWithdrawalFailed represents a WithdrawalFailed error raised by the CollateralManagementContract contract.
type CollateralManagementContractWithdrawalFailed struct {
	From   common.Address
	Amount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WithdrawalFailed(address from, uint256 amount)
func CollateralManagementContractWithdrawalFailedErrorID() common.Hash {
	return common.HexToHash("0x92873d130824b495f22ad10f7f14028200557770e5986714318e78c54f3aa83c")
}

// UnpackWithdrawalFailedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WithdrawalFailed(address from, uint256 amount)
func (collateralManagementContract *CollateralManagementContract) UnpackWithdrawalFailedError(raw []byte) (*CollateralManagementContractWithdrawalFailed, error) {
	out := new(CollateralManagementContractWithdrawalFailed)
	if err := collateralManagementContract.abi.UnpackIntoInterface(out, "WithdrawalFailed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
