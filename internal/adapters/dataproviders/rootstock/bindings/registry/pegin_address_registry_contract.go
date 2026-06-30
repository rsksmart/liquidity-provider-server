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

// PegInAddressRegistryMetaData contains all meta data concerning the PegInAddressRegistry contract.
var PegInAddressRegistryMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"ADDRESS_ENCODING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPegInAddressRegistry.Encoding\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DERIVATION_DOMAIN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"beginDefaultAdminTransfer\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"cancelDefaultAdminTransfer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"changeDefaultAdminDelay\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"defaultAdminDelayIncreaseWait\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBridge\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBridge\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPegInAddressRegistry.Encoding\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegInAddresses\",\"inputs\":[{\"name\":\"addrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"derivationAddresses\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"encoding\",\"type\":\"uint8\",\"internalType\":\"enumIPegInAddressRegistry.Encoding\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrationBlock\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrationCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistrationRoot\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"bridge\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mainnet\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"p2shAddressFromScript\",\"inputs\":[{\"name\":\"segwitScript\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mainnet\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingDefaultAdminDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTx\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackDefaultAdminDelay\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AddressRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"registrationRoot\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminDelayChangeScheduled\",\"inputs\":[{\"name\":\"newDelay\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"effectSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferCanceled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultAdminTransferScheduled\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"acceptSchedule\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminDelay\",\"inputs\":[{\"name\":\"schedule\",\"type\":\"uint48\",\"internalType\":\"uint48\"}]},{\"type\":\"error\",\"name\":\"AccessControlEnforcedDefaultAdminRules\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlInvalidDefaultAdmin\",\"inputs\":[{\"name\":\"defaultAdmin\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidDepositProof\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bridgeResult\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoContract\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Overflow\",\"inputs\":[{\"name\":\"passedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PaymentNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
	ID:  "PegInAddressRegistry",
}

// PegInAddressRegistry is an auto generated Go binding around an Ethereum contract.
type PegInAddressRegistry struct {
	abi abi.ABI
}

// NewPegInAddressRegistry creates a new instance of PegInAddressRegistry.
func NewPegInAddressRegistry() *PegInAddressRegistry {
	parsed, err := PegInAddressRegistryMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PegInAddressRegistry{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PegInAddressRegistry) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackADDRESSENCODING is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd5a5c5f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ADDRESS_ENCODING() view returns(uint8)
func (pegInAddressRegistry *PegInAddressRegistry) PackADDRESSENCODING() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("ADDRESS_ENCODING")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackADDRESSENCODING is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd5a5c5f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ADDRESS_ENCODING() view returns(uint8)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackADDRESSENCODING() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("ADDRESS_ENCODING")
}

// UnpackADDRESSENCODING is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd5a5c5f.
//
// Solidity: function ADDRESS_ENCODING() view returns(uint8)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackADDRESSENCODING(data []byte) (uint8, error) {
	out, err := pegInAddressRegistry.abi.Unpack("ADDRESS_ENCODING", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) PackDEFAULTADMINROLE() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := pegInAddressRegistry.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDERIVATIONDOMAIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x75f8be01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DERIVATION_DOMAIN() view returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) PackDERIVATIONDOMAIN() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("DERIVATION_DOMAIN")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDERIVATIONDOMAIN is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x75f8be01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DERIVATION_DOMAIN() view returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackDERIVATIONDOMAIN() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("DERIVATION_DOMAIN")
}

// UnpackDERIVATIONDOMAIN is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x75f8be01.
//
// Solidity: function DERIVATION_DOMAIN() view returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDERIVATIONDOMAIN(data []byte) ([]byte, error) {
	out, err := pegInAddressRegistry.abi.Unpack("DERIVATION_DOMAIN", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa1ad74.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function VERSION() view returns(string)
func (pegInAddressRegistry *PegInAddressRegistry) PackVERSION() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("VERSION")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVERSION is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffa1ad74.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function VERSION() view returns(string)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackVERSION() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("VERSION")
}

// UnpackVERSION is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffa1ad74.
//
// Solidity: function VERSION() view returns(string)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackVERSION(data []byte) (string, error) {
	out, err := pegInAddressRegistry.abi.Unpack("VERSION", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackAcceptDefaultAdminTransfer() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("acceptDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAcceptDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcefc1429.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function acceptDefaultAdminTransfer() returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackAcceptDefaultAdminTransfer() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("acceptDefaultAdminTransfer")
}

// PackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackBeginDefaultAdminTransfer(newAdmin common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("beginDefaultAdminTransfer", newAdmin)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBeginDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x634e93da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function beginDefaultAdminTransfer(address newAdmin) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackBeginDefaultAdminTransfer(newAdmin common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("beginDefaultAdminTransfer", newAdmin)
}

// PackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackCancelDefaultAdminTransfer() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("cancelDefaultAdminTransfer")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelDefaultAdminTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd602b9fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelDefaultAdminTransfer() returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackCancelDefaultAdminTransfer() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("cancelDefaultAdminTransfer")
}

// PackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackChangeDefaultAdminDelay(newDelay *big.Int) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("changeDefaultAdminDelay", newDelay)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackChangeDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x649a5ec7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function changeDefaultAdminDelay(uint48 newDelay) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackChangeDefaultAdminDelay(newDelay *big.Int) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("changeDefaultAdminDelay", newDelay)
}

// PackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdmin() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) PackDefaultAdmin() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("defaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84ef8ffc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdmin() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackDefaultAdmin() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("defaultAdmin")
}

// UnpackDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x84ef8ffc.
//
// Solidity: function defaultAdmin() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdmin(data []byte) (common.Address, error) {
	out, err := pegInAddressRegistry.abi.Unpack("defaultAdmin", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) PackDefaultAdminDelay() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("defaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcc8463c8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackDefaultAdminDelay() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("defaultAdminDelay")
}

// UnpackDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcc8463c8.
//
// Solidity: function defaultAdminDelay() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminDelay(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistry.abi.Unpack("defaultAdminDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) PackDefaultAdminDelayIncreaseWait() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("defaultAdminDelayIncreaseWait")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDefaultAdminDelayIncreaseWait is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x022d63fb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackDefaultAdminDelayIncreaseWait() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("defaultAdminDelayIncreaseWait")
}

// UnpackDefaultAdminDelayIncreaseWait is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x022d63fb.
//
// Solidity: function defaultAdminDelayIncreaseWait() view returns(uint48)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminDelayIncreaseWait(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistry.abi.Unpack("defaultAdminDelayIncreaseWait", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetBridge is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fffbaf3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBridge() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetBridge() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getBridge")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBridge is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fffbaf3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBridge() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetBridge() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getBridge")
}

// UnpackGetBridge is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0fffbaf3.
//
// Solidity: function getBridge() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetBridge(data []byte) (common.Address, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getBridge", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackGetPegInAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x831adb16.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInAddress(address addr) view returns(bytes, uint8)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetPegInAddress(addr common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getPegInAddress", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x831adb16.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInAddress(address addr) view returns(bytes, uint8)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetPegInAddress(addr common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getPegInAddress", addr)
}

// GetPegInAddressOutput serves as a container for the return parameters of contract
// method GetPegInAddress.
type GetPegInAddressOutput struct {
	Arg0 []byte
	Arg1 uint8
}

// UnpackGetPegInAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x831adb16.
//
// Solidity: function getPegInAddress(address addr) view returns(bytes, uint8)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetPegInAddress(data []byte) (GetPegInAddressOutput, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getPegInAddress", data)
	outstruct := new(GetPegInAddressOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.Arg0 = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Arg1 = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackGetPegInAddresses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd4a25a4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegInAddresses(address[] addrs) view returns(bytes[] derivationAddresses, uint8 encoding)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetPegInAddresses(addrs []common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getPegInAddresses", addrs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegInAddresses is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd4a25a4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegInAddresses(address[] addrs) view returns(bytes[] derivationAddresses, uint8 encoding)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetPegInAddresses(addrs []common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getPegInAddresses", addrs)
}

// GetPegInAddressesOutput serves as a container for the return parameters of contract
// method GetPegInAddresses.
type GetPegInAddressesOutput struct {
	DerivationAddresses [][]byte
	Encoding            uint8
}

// UnpackGetPegInAddresses is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd4a25a4.
//
// Solidity: function getPegInAddresses(address[] addrs) view returns(bytes[] derivationAddresses, uint8 encoding)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetPegInAddresses(data []byte) (GetPegInAddressesOutput, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getPegInAddresses", data)
	outstruct := new(GetPegInAddressesOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.DerivationAddresses = *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	outstruct.Encoding = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	return *outstruct, nil
}

// PackGetRegistrationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5017bc63.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRegistrationBlock(address addr) view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetRegistrationBlock(addr common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getRegistrationBlock", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrationBlock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5017bc63.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrationBlock(address addr) view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetRegistrationBlock(addr common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getRegistrationBlock", addr)
}

// UnpackGetRegistrationBlock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5017bc63.
//
// Solidity: function getRegistrationBlock(address addr) view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetRegistrationBlock(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getRegistrationBlock", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRegistrationCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59bfaee1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRegistrationCount() view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetRegistrationCount() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getRegistrationCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrationCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x59bfaee1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrationCount() view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetRegistrationCount() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getRegistrationCount")
}

// UnpackGetRegistrationCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x59bfaee1.
//
// Solidity: function getRegistrationCount() view returns(uint256)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetRegistrationCount(data []byte) (*big.Int, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getRegistrationCount", data)
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
// Solidity: function getRegistrationRoot() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetRegistrationRoot() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getRegistrationRoot")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRegistrationRoot is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfe625000.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRegistrationRoot() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetRegistrationRoot() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getRegistrationRoot")
}

// UnpackGetRegistrationRoot is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfe625000.
//
// Solidity: function getRegistrationRoot() view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetRegistrationRoot(data []byte) ([32]byte, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getRegistrationRoot", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := pegInAddressRegistry.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackHasRole(data []byte) (bool, error) {
	out, err := pegInAddressRegistry.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb014bdb8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function initialize(address defaultAdmin, uint48 initialDelay, address bridge, bool mainnet) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackInitialize(defaultAdmin common.Address, initialDelay *big.Int, bridge common.Address, mainnet bool) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("initialize", defaultAdmin, initialDelay, bridge, mainnet)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackInitialize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb014bdb8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function initialize(address defaultAdmin, uint48 initialDelay, address bridge, bool mainnet) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackInitialize(defaultAdmin common.Address, initialDelay *big.Int, bridge common.Address, mainnet bool) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("initialize", defaultAdmin, initialDelay, bridge, mainnet)
}

// PackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3c5a547.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) PackIsRegistered(addr common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("isRegistered", addr)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsRegistered is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc3c5a547.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackIsRegistered(addr common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("isRegistered", addr)
}

// UnpackIsRegistered is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc3c5a547.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackIsRegistered(data []byte) (bool, error) {
	out, err := pegInAddressRegistry.abi.Unpack("isRegistered", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) PackOwner() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackOwner() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackOwner(data []byte) (common.Address, error) {
	out, err := pegInAddressRegistry.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackP2shAddressFromScript is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedbef7bd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function p2shAddressFromScript(bytes segwitScript, bool mainnet) pure returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) PackP2shAddressFromScript(segwitScript []byte, mainnet bool) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("p2shAddressFromScript", segwitScript, mainnet)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackP2shAddressFromScript is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xedbef7bd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function p2shAddressFromScript(bytes segwitScript, bool mainnet) pure returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackP2shAddressFromScript(segwitScript []byte, mainnet bool) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("p2shAddressFromScript", segwitScript, mainnet)
}

// UnpackP2shAddressFromScript is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xedbef7bd.
//
// Solidity: function p2shAddressFromScript(bytes segwitScript, bool mainnet) pure returns(bytes)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackP2shAddressFromScript(data []byte) ([]byte, error) {
	out, err := pegInAddressRegistry.abi.Unpack("p2shAddressFromScript", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) PackPendingDefaultAdmin() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("pendingDefaultAdmin")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcf6eefb7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackPendingDefaultAdmin() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("pendingDefaultAdmin")
}

// PendingDefaultAdminOutput serves as a container for the return parameters of contract
// method PendingDefaultAdmin.
type PendingDefaultAdminOutput struct {
	NewAdmin common.Address
	Schedule *big.Int
}

// UnpackPendingDefaultAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcf6eefb7.
//
// Solidity: function pendingDefaultAdmin() view returns(address newAdmin, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackPendingDefaultAdmin(data []byte) (PendingDefaultAdminOutput, error) {
	out, err := pegInAddressRegistry.abi.Unpack("pendingDefaultAdmin", data)
	outstruct := new(PendingDefaultAdminOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewAdmin = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) PackPendingDefaultAdminDelay() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("pendingDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackPendingDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa1eda53c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackPendingDefaultAdminDelay() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("pendingDefaultAdminDelay")
}

// PendingDefaultAdminDelayOutput serves as a container for the return parameters of contract
// method PendingDefaultAdminDelay.
type PendingDefaultAdminDelayOutput struct {
	NewDelay *big.Int
	Schedule *big.Int
}

// UnpackPendingDefaultAdminDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa1eda53c.
//
// Solidity: function pendingDefaultAdminDelay() view returns(uint48 newDelay, uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackPendingDefaultAdminDelay(data []byte) (PendingDefaultAdminDelayOutput, error) {
	out, err := pegInAddressRegistry.abi.Unpack("pendingDefaultAdminDelay", data)
	outstruct := new(PendingDefaultAdminDelayOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.NewDelay = abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	outstruct.Schedule = abi.ConvertType(out[1], new(big.Int)).(*big.Int)
	return *outstruct, nil
}

// PackRegisterAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x19e9d8b0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerAddress(address addr, bytes btcTx, uint256 blockHeight, bytes merkleProof) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackRegisterAddress(addr common.Address, btcTx []byte, blockHeight *big.Int, merkleProof []byte) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("registerAddress", addr, btcTx, blockHeight, merkleProof)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x19e9d8b0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerAddress(address addr, bytes btcTx, uint256 blockHeight, bytes merkleProof) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackRegisterAddress(addr common.Address, btcTx []byte, blockHeight *big.Int, merkleProof []byte) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("registerAddress", addr, btcTx, blockHeight, merkleProof)
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackRenounceRole(role [32]byte, account common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("renounceRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackRenounceRole(role [32]byte, account common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("renounceRole", role, account)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("revokeRole", role, account)
}

// PackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (pegInAddressRegistry *PegInAddressRegistry) PackRollbackDefaultAdminDelay() []byte {
	enc, err := pegInAddressRegistry.abi.Pack("rollbackDefaultAdminDelay")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRollbackDefaultAdminDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0aa6220b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rollbackDefaultAdminDelay() returns()
func (pegInAddressRegistry *PegInAddressRegistry) TryPackRollbackDefaultAdminDelay() ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("rollbackDefaultAdminDelay")
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := pegInAddressRegistry.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return pegInAddressRegistry.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := pegInAddressRegistry.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PegInAddressRegistryAddressRegistered represents a AddressRegistered event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAddressRegistered struct {
	Addr             common.Address
	RegistrationRoot [32]byte
	Raw              *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryAddressRegisteredEventName = "AddressRegistered"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryAddressRegistered) ContractEventName() string {
	return PegInAddressRegistryAddressRegisteredEventName
}

// UnpackAddressRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AddressRegistered(address indexed addr, bytes32 indexed registrationRoot)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAddressRegisteredEvent(log *types.Log) (*PegInAddressRegistryAddressRegistered, error) {
	event := "AddressRegistered"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryAddressRegistered)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryDefaultAdminDelayChangeCanceled represents a DefaultAdminDelayChangeCanceled event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryDefaultAdminDelayChangeCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryDefaultAdminDelayChangeCanceledEventName = "DefaultAdminDelayChangeCanceled"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryDefaultAdminDelayChangeCanceled) ContractEventName() string {
	return PegInAddressRegistryDefaultAdminDelayChangeCanceledEventName
}

// UnpackDefaultAdminDelayChangeCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeCanceled()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminDelayChangeCanceledEvent(log *types.Log) (*PegInAddressRegistryDefaultAdminDelayChangeCanceled, error) {
	event := "DefaultAdminDelayChangeCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryDefaultAdminDelayChangeCanceled)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryDefaultAdminDelayChangeScheduled represents a DefaultAdminDelayChangeScheduled event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryDefaultAdminDelayChangeScheduled struct {
	NewDelay       *big.Int
	EffectSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryDefaultAdminDelayChangeScheduledEventName = "DefaultAdminDelayChangeScheduled"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryDefaultAdminDelayChangeScheduled) ContractEventName() string {
	return PegInAddressRegistryDefaultAdminDelayChangeScheduledEventName
}

// UnpackDefaultAdminDelayChangeScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminDelayChangeScheduled(uint48 newDelay, uint48 effectSchedule)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminDelayChangeScheduledEvent(log *types.Log) (*PegInAddressRegistryDefaultAdminDelayChangeScheduled, error) {
	event := "DefaultAdminDelayChangeScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryDefaultAdminDelayChangeScheduled)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryDefaultAdminTransferCanceled represents a DefaultAdminTransferCanceled event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryDefaultAdminTransferCanceled struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryDefaultAdminTransferCanceledEventName = "DefaultAdminTransferCanceled"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryDefaultAdminTransferCanceled) ContractEventName() string {
	return PegInAddressRegistryDefaultAdminTransferCanceledEventName
}

// UnpackDefaultAdminTransferCanceledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferCanceled()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminTransferCanceledEvent(log *types.Log) (*PegInAddressRegistryDefaultAdminTransferCanceled, error) {
	event := "DefaultAdminTransferCanceled"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryDefaultAdminTransferCanceled)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryDefaultAdminTransferScheduled represents a DefaultAdminTransferScheduled event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryDefaultAdminTransferScheduled struct {
	NewAdmin       common.Address
	AcceptSchedule *big.Int
	Raw            *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryDefaultAdminTransferScheduledEventName = "DefaultAdminTransferScheduled"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryDefaultAdminTransferScheduled) ContractEventName() string {
	return PegInAddressRegistryDefaultAdminTransferScheduledEventName
}

// UnpackDefaultAdminTransferScheduledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DefaultAdminTransferScheduled(address indexed newAdmin, uint48 acceptSchedule)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackDefaultAdminTransferScheduledEvent(log *types.Log) (*PegInAddressRegistryDefaultAdminTransferScheduled, error) {
	event := "DefaultAdminTransferScheduled"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryDefaultAdminTransferScheduled)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryInitialized represents a Initialized event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryInitialized struct {
	Version uint64
	Raw     *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryInitializedEventName = "Initialized"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryInitialized) ContractEventName() string {
	return PegInAddressRegistryInitializedEventName
}

// UnpackInitializedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Initialized(uint64 version)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackInitializedEvent(log *types.Log) (*PegInAddressRegistryInitialized, error) {
	event := "Initialized"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryInitialized)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryRoleAdminChanged) ContractEventName() string {
	return PegInAddressRegistryRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackRoleAdminChangedEvent(log *types.Log) (*PegInAddressRegistryRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryRoleGranted represents a RoleGranted event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryRoleGranted) ContractEventName() string {
	return PegInAddressRegistryRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackRoleGrantedEvent(log *types.Log) (*PegInAddressRegistryRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryRoleGranted)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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

// PegInAddressRegistryRoleRevoked represents a RoleRevoked event raised by the PegInAddressRegistry contract.
type PegInAddressRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const PegInAddressRegistryRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (PegInAddressRegistryRoleRevoked) ContractEventName() string {
	return PegInAddressRegistryRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackRoleRevokedEvent(log *types.Log) (*PegInAddressRegistryRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 || log.Topics[0] != pegInAddressRegistry.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegInAddressRegistryRoleRevoked)
	if len(log.Data) > 0 {
		if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegInAddressRegistry.abi.Events[event].Inputs {
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
func (pegInAddressRegistry *PegInAddressRegistry) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AccessControlEnforcedDefaultAdminDelay"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAccessControlEnforcedDefaultAdminDelayError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AccessControlEnforcedDefaultAdminRules"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAccessControlEnforcedDefaultAdminRulesError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AccessControlInvalidDefaultAdmin"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAccessControlInvalidDefaultAdminError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["AlreadyRegistered"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["InvalidDepositProof"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackInvalidDepositProofError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["InvalidInitialization"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackInvalidInitializationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["NoContract"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackNoContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["NotInitializing"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackNotInitializingError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["Overflow"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackOverflowError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["PaymentNotAllowed"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackPaymentNotAllowedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegInAddressRegistry.abi.Errors["SafeCastOverflowedUintDowncast"].ID.Bytes()[:4]) {
		return pegInAddressRegistry.UnpackSafeCastOverflowedUintDowncastError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PegInAddressRegistryAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func PegInAddressRegistryAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAccessControlBadConfirmationError(raw []byte) (*PegInAddressRegistryAccessControlBadConfirmation, error) {
	out := new(PegInAddressRegistryAccessControlBadConfirmation)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryAccessControlEnforcedDefaultAdminDelay represents a AccessControlEnforcedDefaultAdminDelay error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAccessControlEnforcedDefaultAdminDelay struct {
	Schedule *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func PegInAddressRegistryAccessControlEnforcedDefaultAdminDelayErrorID() common.Hash {
	return common.HexToHash("0x19ca5ebb8fb33f00e502c9392eddab1501674629178bf69b853cf037aaf4bb5d")
}

// UnpackAccessControlEnforcedDefaultAdminDelayError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminDelay(uint48 schedule)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAccessControlEnforcedDefaultAdminDelayError(raw []byte) (*PegInAddressRegistryAccessControlEnforcedDefaultAdminDelay, error) {
	out := new(PegInAddressRegistryAccessControlEnforcedDefaultAdminDelay)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminDelay", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryAccessControlEnforcedDefaultAdminRules represents a AccessControlEnforcedDefaultAdminRules error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAccessControlEnforcedDefaultAdminRules struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func PegInAddressRegistryAccessControlEnforcedDefaultAdminRulesErrorID() common.Hash {
	return common.HexToHash("0x3fc3c27ae3db78c81b8f6e685172134623efa268ee8cd8d54be38ad2a74fc13b")
}

// UnpackAccessControlEnforcedDefaultAdminRulesError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlEnforcedDefaultAdminRules()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAccessControlEnforcedDefaultAdminRulesError(raw []byte) (*PegInAddressRegistryAccessControlEnforcedDefaultAdminRules, error) {
	out := new(PegInAddressRegistryAccessControlEnforcedDefaultAdminRules)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AccessControlEnforcedDefaultAdminRules", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryAccessControlInvalidDefaultAdmin represents a AccessControlInvalidDefaultAdmin error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAccessControlInvalidDefaultAdmin struct {
	DefaultAdmin common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func PegInAddressRegistryAccessControlInvalidDefaultAdminErrorID() common.Hash {
	return common.HexToHash("0xc22c8022f2a840d6b6a9f113407715f5bbd4e88c1b0dd9434dc00700ba609ed4")
}

// UnpackAccessControlInvalidDefaultAdminError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlInvalidDefaultAdmin(address defaultAdmin)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAccessControlInvalidDefaultAdminError(raw []byte) (*PegInAddressRegistryAccessControlInvalidDefaultAdmin, error) {
	out := new(PegInAddressRegistryAccessControlInvalidDefaultAdmin)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AccessControlInvalidDefaultAdmin", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func PegInAddressRegistryAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*PegInAddressRegistryAccessControlUnauthorizedAccount, error) {
	out := new(PegInAddressRegistryAccessControlUnauthorizedAccount)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryAlreadyRegistered represents a AlreadyRegistered error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryAlreadyRegistered struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AlreadyRegistered(address addr)
func PegInAddressRegistryAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x45ed80e9399c87887ea54f615514a1e3dde31e9b6c027ddceb4ffd503b70e428")
}

// UnpackAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AlreadyRegistered(address addr)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackAlreadyRegisteredError(raw []byte) (*PegInAddressRegistryAlreadyRegistered, error) {
	out := new(PegInAddressRegistryAlreadyRegistered)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "AlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryInvalidDepositProof represents a InvalidDepositProof error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryInvalidDepositProof struct {
	Addr         common.Address
	BridgeResult *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidDepositProof(address addr, int256 bridgeResult)
func PegInAddressRegistryInvalidDepositProofErrorID() common.Hash {
	return common.HexToHash("0xe752092735389923328dfbd794b89402fdc5ea0a92bd311b815dce62336801f5")
}

// UnpackInvalidDepositProofError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidDepositProof(address addr, int256 bridgeResult)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackInvalidDepositProofError(raw []byte) (*PegInAddressRegistryInvalidDepositProof, error) {
	out := new(PegInAddressRegistryInvalidDepositProof)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "InvalidDepositProof", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryInvalidInitialization represents a InvalidInitialization error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryInvalidInitialization struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidInitialization()
func PegInAddressRegistryInvalidInitializationErrorID() common.Hash {
	return common.HexToHash("0xf92ee8a957075833165f68c320933b1a1294aafc84ee6e0dd3fb178008f9aaf5")
}

// UnpackInvalidInitializationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidInitialization()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackInvalidInitializationError(raw []byte) (*PegInAddressRegistryInvalidInitialization, error) {
	out := new(PegInAddressRegistryInvalidInitialization)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "InvalidInitialization", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryNoContract represents a NoContract error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryNoContract struct {
	Addr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NoContract(address addr)
func PegInAddressRegistryNoContractErrorID() common.Hash {
	return common.HexToHash("0x5f15d672b6235f8600ffc72925d8d2f9dcea14be067296327891153847185a5c")
}

// UnpackNoContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NoContract(address addr)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackNoContractError(raw []byte) (*PegInAddressRegistryNoContract, error) {
	out := new(PegInAddressRegistryNoContract)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "NoContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryNotInitializing represents a NotInitializing error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryNotInitializing struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotInitializing()
func PegInAddressRegistryNotInitializingErrorID() common.Hash {
	return common.HexToHash("0xd7e6bcf8597daa127dc9f0048d2f08d5ef140a2cb659feabd700beff1f7a8302")
}

// UnpackNotInitializingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotInitializing()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackNotInitializingError(raw []byte) (*PegInAddressRegistryNotInitializing, error) {
	out := new(PegInAddressRegistryNotInitializing)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "NotInitializing", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryOverflow represents a Overflow error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryOverflow struct {
	PassedAmount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Overflow(uint256 passedAmount)
func PegInAddressRegistryOverflowErrorID() common.Hash {
	return common.HexToHash("0xe0fb6a7ce291b396fa814871fbb6fcc26c1a1454a6e18a2e7c911a8763b928dc")
}

// UnpackOverflowError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Overflow(uint256 passedAmount)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackOverflowError(raw []byte) (*PegInAddressRegistryOverflow, error) {
	out := new(PegInAddressRegistryOverflow)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "Overflow", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryPaymentNotAllowed represents a PaymentNotAllowed error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryPaymentNotAllowed struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PaymentNotAllowed()
func PegInAddressRegistryPaymentNotAllowedErrorID() common.Hash {
	return common.HexToHash("0x8619bd43ab22b4b01742bd29d231dff1e50413ee3a444878bed65970c80c97df")
}

// UnpackPaymentNotAllowedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PaymentNotAllowed()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackPaymentNotAllowedError(raw []byte) (*PegInAddressRegistryPaymentNotAllowed, error) {
	out := new(PegInAddressRegistryPaymentNotAllowed)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "PaymentNotAllowed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistryReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the PegInAddressRegistry contract.
type PegInAddressRegistryReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func PegInAddressRegistryReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (pegInAddressRegistry *PegInAddressRegistry) UnpackReentrancyGuardReentrantCallError(raw []byte) (*PegInAddressRegistryReentrancyGuardReentrantCall, error) {
	out := new(PegInAddressRegistryReentrancyGuardReentrantCall)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegInAddressRegistrySafeCastOverflowedUintDowncast represents a SafeCastOverflowedUintDowncast error raised by the PegInAddressRegistry contract.
type PegInAddressRegistrySafeCastOverflowedUintDowncast struct {
	Bits  uint8
	Value *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func PegInAddressRegistrySafeCastOverflowedUintDowncastErrorID() common.Hash {
	return common.HexToHash("0x6dfcc6503a32754ce7a89698e18201fc5294fd4aad43edefee786f88423b1a12")
}

// UnpackSafeCastOverflowedUintDowncastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeCastOverflowedUintDowncast(uint8 bits, uint256 value)
func (pegInAddressRegistry *PegInAddressRegistry) UnpackSafeCastOverflowedUintDowncastError(raw []byte) (*PegInAddressRegistrySafeCastOverflowedUintDowncast, error) {
	out := new(PegInAddressRegistrySafeCastOverflowedUintDowncast)
	if err := pegInAddressRegistry.abi.UnpackIntoInterface(out, "SafeCastOverflowedUintDowncast", raw); err != nil {
		return nil, err
	}
	return out, nil
}
