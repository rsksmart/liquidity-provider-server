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

// QuotesPegOutQuote is an auto generated low-level Go binding around an user-defined struct.
type QuotesPegOutQuote struct {
	ChainId               *big.Int
	CallFee               *big.Int
	PenaltyFee            *big.Int
	Value                 *big.Int
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

// PegOutEscrowContractMetaData contains all meta data concerning the PegOutEscrowContract contract.
var PegOutEscrowContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"cancelPegOut\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimFailCount\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimPegOut\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getMaxMinerFee\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutQuote\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPegOutState\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIPegOutEscrow.EscrowedPegOutState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onClaimFail\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onSettlement\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"finalState\",\"type\":\"uint8\",\"internalType\":\"enumIPegOutEscrow.EscrowedPegOutState\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refundOnNoClaim\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestIdAt\",\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestPegOut\",\"inputs\":[{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"refundAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"restrictedUntil\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalRequests\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unrevoke\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GlobalSlashSkipped\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutCancelled\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutClaimed\",\"inputs\":[{\"name\":\"lpAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"requestHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutRefundedOnNoClaim\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"refundAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutRequested\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"refundAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ClaimWindowClosed\",\"inputs\":[{\"name\":\"depositDateLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ClaimWindowOpen\",\"inputs\":[{\"name\":\"depositDateLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"CollateralManagementNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDestination\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidState\",\"inputs\":[{\"name\":\"requestHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expected\",\"type\":\"uint8\",\"internalType\":\"enumIPegOutEscrow.EscrowedPegOutState\"},{\"name\":\"actual\",\"type\":\"uint8\",\"internalType\":\"enumIPegOutEscrow.EscrowedPegOutState\"}]},{\"type\":\"error\",\"name\":\"LpRestricted\",\"inputs\":[{\"name\":\"lp\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"restrictedUntil\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotServiceable\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OnlyPegOutContract\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PegOutContractNotSet\",\"inputs\":[]}]",
	ID:  "PegOutEscrowContract",
}

// PegOutEscrowContract is an auto generated Go binding around an Ethereum contract.
type PegOutEscrowContract struct {
	abi abi.ABI
}

// NewPegOutEscrowContract creates a new instance of PegOutEscrowContract.
func NewPegOutEscrowContract() *PegOutEscrowContract {
	parsed, err := PegOutEscrowContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PegOutEscrowContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PegOutEscrowContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackCancelPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x288d45f3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function cancelPegOut(bytes32 requestHash) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackCancelPegOut(requestHash [32]byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("cancelPegOut", requestHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCancelPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x288d45f3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function cancelPegOut(bytes32 requestHash) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackCancelPegOut(requestHash [32]byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("cancelPegOut", requestHash)
}

// PackClaimFailCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66e352c6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimFailCount(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) PackClaimFailCount(lp common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("claimFailCount", lp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimFailCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66e352c6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimFailCount(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackClaimFailCount(lp common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("claimFailCount", lp)
}

// UnpackClaimFailCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x66e352c6.
//
// Solidity: function claimFailCount(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackClaimFailCount(data []byte) (*big.Int, error) {
	out, err := pegOutEscrowContract.abi.Unpack("claimFailCount", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackClaimPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4a0e9dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claimPegOut(bytes32 requestHash, bytes signature) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackClaimPegOut(requestHash [32]byte, signature []byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("claimPegOut", requestHash, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaimPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd4a0e9dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claimPegOut(bytes32 requestHash, bytes signature) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackClaimPegOut(requestHash [32]byte, signature []byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("claimPegOut", requestHash, signature)
}

// PackGetMaxMinerFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb469343d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMaxMinerFee(bytes32 requestHash) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) PackGetMaxMinerFee(requestHash [32]byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("getMaxMinerFee", requestHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMaxMinerFee is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb469343d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMaxMinerFee(bytes32 requestHash) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackGetMaxMinerFee(requestHash [32]byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("getMaxMinerFee", requestHash)
}

// UnpackGetMaxMinerFee is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb469343d.
//
// Solidity: function getMaxMinerFee(bytes32 requestHash) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackGetMaxMinerFee(data []byte) (*big.Int, error) {
	out, err := pegOutEscrowContract.abi.Unpack("getMaxMinerFee", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPegOutQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3998c1c3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutQuote(bytes32 requestHash) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes))
func (pegOutEscrowContract *PegOutEscrowContract) PackGetPegOutQuote(requestHash [32]byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("getPegOutQuote", requestHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3998c1c3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutQuote(bytes32 requestHash) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes))
func (pegOutEscrowContract *PegOutEscrowContract) TryPackGetPegOutQuote(requestHash [32]byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("getPegOutQuote", requestHash)
}

// UnpackGetPegOutQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3998c1c3.
//
// Solidity: function getPegOutQuote(bytes32 requestHash) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes))
func (pegOutEscrowContract *PegOutEscrowContract) UnpackGetPegOutQuote(data []byte) (QuotesPegOutQuote, error) {
	out, err := pegOutEscrowContract.abi.Unpack("getPegOutQuote", data)
	if err != nil {
		return *new(QuotesPegOutQuote), err
	}
	out0 := *abi.ConvertType(out[0], new(QuotesPegOutQuote)).(*QuotesPegOutQuote)
	return out0, nil
}

// PackGetPegOutState is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa03a0f25.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPegOutState(bytes32 requestHash) view returns(uint8)
func (pegOutEscrowContract *PegOutEscrowContract) PackGetPegOutState(requestHash [32]byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("getPegOutState", requestHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPegOutState is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa03a0f25.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPegOutState(bytes32 requestHash) view returns(uint8)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackGetPegOutState(requestHash [32]byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("getPegOutState", requestHash)
}

// UnpackGetPegOutState is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa03a0f25.
//
// Solidity: function getPegOutState(bytes32 requestHash) view returns(uint8)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackGetPegOutState(data []byte) (uint8, error) {
	out, err := pegOutEscrowContract.abi.Unpack("getPegOutState", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackOnClaimFail is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56311623.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onClaimFail(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackOnClaimFail(lp common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("onClaimFail", lp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnClaimFail is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56311623.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onClaimFail(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackOnClaimFail(lp common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("onClaimFail", lp)
}

// PackOnSettlement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe706b5eb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function onSettlement(bytes32 requestHash, uint8 finalState) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackOnSettlement(requestHash [32]byte, finalState uint8) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("onSettlement", requestHash, finalState)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOnSettlement is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe706b5eb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function onSettlement(bytes32 requestHash, uint8 finalState) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackOnSettlement(requestHash [32]byte, finalState uint8) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("onSettlement", requestHash, finalState)
}

// PackRefundOnNoClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xebda26bc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refundOnNoClaim(bytes32 requestHash) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackRefundOnNoClaim(requestHash [32]byte) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("refundOnNoClaim", requestHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefundOnNoClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xebda26bc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refundOnNoClaim(bytes32 requestHash) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackRefundOnNoClaim(requestHash [32]byte) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("refundOnNoClaim", requestHash)
}

// PackRequestIdAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8a33b1c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestIdAt(uint256 nonce) view returns(bytes32)
func (pegOutEscrowContract *PegOutEscrowContract) PackRequestIdAt(nonce *big.Int) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("requestIdAt", nonce)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestIdAt is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd8a33b1c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestIdAt(uint256 nonce) view returns(bytes32)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackRequestIdAt(nonce *big.Int) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("requestIdAt", nonce)
}

// UnpackRequestIdAt is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd8a33b1c.
//
// Solidity: function requestIdAt(uint256 nonce) view returns(bytes32)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackRequestIdAt(data []byte) ([32]byte, error) {
	out, err := pegOutEscrowContract.abi.Unpack("requestIdAt", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRequestPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9ea29d7e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestPegOut(bytes destinationAddress, address refundAddress) payable returns(bytes32 requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) PackRequestPegOut(destinationAddress []byte, refundAddress common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("requestPegOut", destinationAddress, refundAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9ea29d7e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestPegOut(bytes destinationAddress, address refundAddress) payable returns(bytes32 requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackRequestPegOut(destinationAddress []byte, refundAddress common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("requestPegOut", destinationAddress, refundAddress)
}

// UnpackRequestPegOut is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x9ea29d7e.
//
// Solidity: function requestPegOut(bytes destinationAddress, address refundAddress) payable returns(bytes32 requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackRequestPegOut(data []byte) ([32]byte, error) {
	out, err := pegOutEscrowContract.abi.Unpack("requestPegOut", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRestrictedUntil is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xed622994.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function restrictedUntil(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) PackRestrictedUntil(lp common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("restrictedUntil", lp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRestrictedUntil is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xed622994.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function restrictedUntil(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackRestrictedUntil(lp common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("restrictedUntil", lp)
}

// UnpackRestrictedUntil is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xed622994.
//
// Solidity: function restrictedUntil(address lp) view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackRestrictedUntil(data []byte) (*big.Int, error) {
	out, err := pegOutEscrowContract.abi.Unpack("restrictedUntil", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74a8f103.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revoke(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackRevoke(lp common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("revoke", lp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x74a8f103.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revoke(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackRevoke(lp common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("revoke", lp)
}

// PackTotalRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8aea61dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalRequests() view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) PackTotalRequests() []byte {
	enc, err := pegOutEscrowContract.abi.Pack("totalRequests")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalRequests is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8aea61dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalRequests() view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) TryPackTotalRequests() ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("totalRequests")
}

// UnpackTotalRequests is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8aea61dc.
//
// Solidity: function totalRequests() view returns(uint256)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackTotalRequests(data []byte) (*big.Int, error) {
	out, err := pegOutEscrowContract.abi.Unpack("totalRequests", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUnrevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe06145a4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function unrevoke(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) PackUnrevoke(lp common.Address) []byte {
	enc, err := pegOutEscrowContract.abi.Pack("unrevoke", lp)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUnrevoke is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe06145a4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function unrevoke(address lp) returns()
func (pegOutEscrowContract *PegOutEscrowContract) TryPackUnrevoke(lp common.Address) ([]byte, error) {
	return pegOutEscrowContract.abi.Pack("unrevoke", lp)
}

// PegOutEscrowContractGlobalSlashSkipped represents a GlobalSlashSkipped event raised by the PegOutEscrowContract contract.
type PegOutEscrowContractGlobalSlashSkipped struct {
	RequestHash [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const PegOutEscrowContractGlobalSlashSkippedEventName = "GlobalSlashSkipped"

// ContractEventName returns the user-defined event name.
func (PegOutEscrowContractGlobalSlashSkipped) ContractEventName() string {
	return PegOutEscrowContractGlobalSlashSkippedEventName
}

// UnpackGlobalSlashSkippedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event GlobalSlashSkipped(bytes32 indexed requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackGlobalSlashSkippedEvent(log *types.Log) (*PegOutEscrowContractGlobalSlashSkipped, error) {
	event := "GlobalSlashSkipped"
	if len(log.Topics) == 0 || log.Topics[0] != pegOutEscrowContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegOutEscrowContractGlobalSlashSkipped)
	if len(log.Data) > 0 {
		if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegOutEscrowContract.abi.Events[event].Inputs {
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

// PegOutEscrowContractPegOutCancelled represents a PegOutCancelled event raised by the PegOutEscrowContract contract.
type PegOutEscrowContractPegOutCancelled struct {
	RequestHash [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const PegOutEscrowContractPegOutCancelledEventName = "PegOutCancelled"

// ContractEventName returns the user-defined event name.
func (PegOutEscrowContractPegOutCancelled) ContractEventName() string {
	return PegOutEscrowContractPegOutCancelledEventName
}

// UnpackPegOutCancelledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutCancelled(bytes32 indexed requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackPegOutCancelledEvent(log *types.Log) (*PegOutEscrowContractPegOutCancelled, error) {
	event := "PegOutCancelled"
	if len(log.Topics) == 0 || log.Topics[0] != pegOutEscrowContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegOutEscrowContractPegOutCancelled)
	if len(log.Data) > 0 {
		if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegOutEscrowContract.abi.Events[event].Inputs {
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

// PegOutEscrowContractPegOutClaimed represents a PegOutClaimed event raised by the PegOutEscrowContract contract.
type PegOutEscrowContractPegOutClaimed struct {
	LpAddress   common.Address
	RequestHash [32]byte
	Raw         *types.Log // Blockchain specific contextual infos
}

const PegOutEscrowContractPegOutClaimedEventName = "PegOutClaimed"

// ContractEventName returns the user-defined event name.
func (PegOutEscrowContractPegOutClaimed) ContractEventName() string {
	return PegOutEscrowContractPegOutClaimedEventName
}

// UnpackPegOutClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutClaimed(address indexed lpAddress, bytes32 indexed requestHash)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackPegOutClaimedEvent(log *types.Log) (*PegOutEscrowContractPegOutClaimed, error) {
	event := "PegOutClaimed"
	if len(log.Topics) == 0 || log.Topics[0] != pegOutEscrowContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegOutEscrowContractPegOutClaimed)
	if len(log.Data) > 0 {
		if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegOutEscrowContract.abi.Events[event].Inputs {
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

// PegOutEscrowContractPegOutRefundedOnNoClaim represents a PegOutRefundedOnNoClaim event raised by the PegOutEscrowContract contract.
type PegOutEscrowContractPegOutRefundedOnNoClaim struct {
	RequestHash   [32]byte
	RefundAddress common.Address
	Amount        *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const PegOutEscrowContractPegOutRefundedOnNoClaimEventName = "PegOutRefundedOnNoClaim"

// ContractEventName returns the user-defined event name.
func (PegOutEscrowContractPegOutRefundedOnNoClaim) ContractEventName() string {
	return PegOutEscrowContractPegOutRefundedOnNoClaimEventName
}

// UnpackPegOutRefundedOnNoClaimEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutRefundedOnNoClaim(bytes32 indexed requestHash, address indexed refundAddress, uint256 amount)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackPegOutRefundedOnNoClaimEvent(log *types.Log) (*PegOutEscrowContractPegOutRefundedOnNoClaim, error) {
	event := "PegOutRefundedOnNoClaim"
	if len(log.Topics) == 0 || log.Topics[0] != pegOutEscrowContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegOutEscrowContractPegOutRefundedOnNoClaim)
	if len(log.Data) > 0 {
		if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegOutEscrowContract.abi.Events[event].Inputs {
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

// PegOutEscrowContractPegOutRequested represents a PegOutRequested event raised by the PegOutEscrowContract contract.
type PegOutEscrowContractPegOutRequested struct {
	RequestHash        [32]byte
	RefundAddress      common.Address
	Amount             *big.Int
	DestinationAddress []byte
	Raw                *types.Log // Blockchain specific contextual infos
}

const PegOutEscrowContractPegOutRequestedEventName = "PegOutRequested"

// ContractEventName returns the user-defined event name.
func (PegOutEscrowContractPegOutRequested) ContractEventName() string {
	return PegOutEscrowContractPegOutRequestedEventName
}

// UnpackPegOutRequestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutRequested(bytes32 indexed requestHash, address indexed refundAddress, uint256 indexed amount, bytes destinationAddress)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackPegOutRequestedEvent(log *types.Log) (*PegOutEscrowContractPegOutRequested, error) {
	event := "PegOutRequested"
	if len(log.Topics) == 0 || log.Topics[0] != pegOutEscrowContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegOutEscrowContractPegOutRequested)
	if len(log.Data) > 0 {
		if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegOutEscrowContract.abi.Events[event].Inputs {
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
func (pegOutEscrowContract *PegOutEscrowContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["ClaimWindowClosed"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackClaimWindowClosedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["ClaimWindowOpen"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackClaimWindowOpenError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["CollateralManagementNotSet"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackCollateralManagementNotSetError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["InvalidDestination"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackInvalidDestinationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["InvalidState"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackInvalidStateError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["LpRestricted"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackLpRestrictedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["NotServiceable"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackNotServiceableError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["OnlyPegOutContract"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackOnlyPegOutContractError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegOutEscrowContract.abi.Errors["PegOutContractNotSet"].ID.Bytes()[:4]) {
		return pegOutEscrowContract.UnpackPegOutContractNotSetError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PegOutEscrowContractClaimWindowClosed represents a ClaimWindowClosed error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractClaimWindowClosed struct {
	DepositDateLimit *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ClaimWindowClosed(uint256 depositDateLimit)
func PegOutEscrowContractClaimWindowClosedErrorID() common.Hash {
	return common.HexToHash("0x8fb6d6875cba88232b860f532fb9b3d4a0dccd84f2f275a3db894cf8a0a05f2c")
}

// UnpackClaimWindowClosedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ClaimWindowClosed(uint256 depositDateLimit)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackClaimWindowClosedError(raw []byte) (*PegOutEscrowContractClaimWindowClosed, error) {
	out := new(PegOutEscrowContractClaimWindowClosed)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "ClaimWindowClosed", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractClaimWindowOpen represents a ClaimWindowOpen error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractClaimWindowOpen struct {
	DepositDateLimit *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ClaimWindowOpen(uint256 depositDateLimit)
func PegOutEscrowContractClaimWindowOpenErrorID() common.Hash {
	return common.HexToHash("0xffdab61d336b4e791cb18d74737af581b8f13961dc539ce462d04974bf0bac44")
}

// UnpackClaimWindowOpenError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ClaimWindowOpen(uint256 depositDateLimit)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackClaimWindowOpenError(raw []byte) (*PegOutEscrowContractClaimWindowOpen, error) {
	out := new(PegOutEscrowContractClaimWindowOpen)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "ClaimWindowOpen", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractCollateralManagementNotSet represents a CollateralManagementNotSet error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractCollateralManagementNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollateralManagementNotSet()
func PegOutEscrowContractCollateralManagementNotSetErrorID() common.Hash {
	return common.HexToHash("0x94f06056204868752f999c36d8da7b399c7c84dfd1c8e04d1c29d43029d2964b")
}

// UnpackCollateralManagementNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollateralManagementNotSet()
func (pegOutEscrowContract *PegOutEscrowContract) UnpackCollateralManagementNotSetError(raw []byte) (*PegOutEscrowContractCollateralManagementNotSet, error) {
	out := new(PegOutEscrowContractCollateralManagementNotSet)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "CollateralManagementNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractInvalidDestination represents a InvalidDestination error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractInvalidDestination struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidDestination()
func PegOutEscrowContractInvalidDestinationErrorID() common.Hash {
	return common.HexToHash("0xac6b05f542ae1b66b4b43710bf4773616692aad79ee0fe8cb9c5db0f38dd0764")
}

// UnpackInvalidDestinationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidDestination()
func (pegOutEscrowContract *PegOutEscrowContract) UnpackInvalidDestinationError(raw []byte) (*PegOutEscrowContractInvalidDestination, error) {
	out := new(PegOutEscrowContractInvalidDestination)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "InvalidDestination", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractInvalidState represents a InvalidState error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractInvalidState struct {
	RequestHash [32]byte
	Expected    uint8
	Actual      uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidState(bytes32 requestHash, uint8 expected, uint8 actual)
func PegOutEscrowContractInvalidStateErrorID() common.Hash {
	return common.HexToHash("0x4e4d80a5a1bafc929c906c8d7c3bf009923e60165b6c4e3b2e2163aba9634adc")
}

// UnpackInvalidStateError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidState(bytes32 requestHash, uint8 expected, uint8 actual)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackInvalidStateError(raw []byte) (*PegOutEscrowContractInvalidState, error) {
	out := new(PegOutEscrowContractInvalidState)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "InvalidState", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractLpRestricted represents a LpRestricted error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractLpRestricted struct {
	Lp              common.Address
	RestrictedUntil *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LpRestricted(address lp, uint256 restrictedUntil)
func PegOutEscrowContractLpRestrictedErrorID() common.Hash {
	return common.HexToHash("0x35251d8f6d53bbcd56d82926b81b245bddb129948d3b059c41849f2e6aade265")
}

// UnpackLpRestrictedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LpRestricted(address lp, uint256 restrictedUntil)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackLpRestrictedError(raw []byte) (*PegOutEscrowContractLpRestricted, error) {
	out := new(PegOutEscrowContractLpRestricted)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "LpRestricted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractNotServiceable represents a NotServiceable error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractNotServiceable struct {
	Amount    *big.Int
	MinAmount *big.Int
	MaxAmount *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotServiceable(uint256 amount, uint256 minAmount, uint256 maxAmount)
func PegOutEscrowContractNotServiceableErrorID() common.Hash {
	return common.HexToHash("0xa3d25c1cd616834ed6cd725091b2f8e0b37567fd88299b4402f0cbb67d76e60c")
}

// UnpackNotServiceableError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotServiceable(uint256 amount, uint256 minAmount, uint256 maxAmount)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackNotServiceableError(raw []byte) (*PegOutEscrowContractNotServiceable, error) {
	out := new(PegOutEscrowContractNotServiceable)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "NotServiceable", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractOnlyPegOutContract represents a OnlyPegOutContract error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractOnlyPegOutContract struct {
	Caller common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OnlyPegOutContract(address caller)
func PegOutEscrowContractOnlyPegOutContractErrorID() common.Hash {
	return common.HexToHash("0x96122122e327409abf975214cde141ef185c2504614c526886b84cc4a870d7f0")
}

// UnpackOnlyPegOutContractError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OnlyPegOutContract(address caller)
func (pegOutEscrowContract *PegOutEscrowContract) UnpackOnlyPegOutContractError(raw []byte) (*PegOutEscrowContractOnlyPegOutContract, error) {
	out := new(PegOutEscrowContractOnlyPegOutContract)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "OnlyPegOutContract", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegOutEscrowContractPegOutContractNotSet represents a PegOutContractNotSet error raised by the PegOutEscrowContract contract.
type PegOutEscrowContractPegOutContractNotSet struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PegOutContractNotSet()
func PegOutEscrowContractPegOutContractNotSetErrorID() common.Hash {
	return common.HexToHash("0x58bfb6e7f457e27a5f924ffdf2f6366a1ccbe3f6ebe0b6f3a6d9874ff16e0f8a")
}

// UnpackPegOutContractNotSetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PegOutContractNotSet()
func (pegOutEscrowContract *PegOutEscrowContract) UnpackPegOutContractNotSetError(raw []byte) (*PegOutEscrowContractPegOutContractNotSet, error) {
	out := new(PegOutEscrowContractPegOutContractNotSet)
	if err := pegOutEscrowContract.abi.UnpackIntoInterface(out, "PegOutContractNotSet", raw); err != nil {
		return nil, err
	}
	return out, nil
}
