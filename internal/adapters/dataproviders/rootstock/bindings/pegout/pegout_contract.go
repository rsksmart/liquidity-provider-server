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

// PegoutContractMetaData contains all meta data concerning the PegoutContract contract.
var PegoutContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"depositPegOut\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBalance\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashPegOutQuote\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashPegOutQuoteEIP712\",\"inputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isQuoteCompleted\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"isPaused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"since\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"refundPegOut\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"btcTx\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcBlockHeaderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleBranchPath\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refundUserPegOut\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatePegout\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"btcTx\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"quote\",\"type\":\"tuple\",\"internalType\":\"structQuotes.PegOutQuote\",\"components\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"callFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lbcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpRskAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rskRefundAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"int64\",\"internalType\":\"int64\"},{\"name\":\"agreementTimestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"transferTime\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"depositConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"transferConfirmations\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"depositAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcRefundAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"lpBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BalanceDecrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BalanceIncrease\",\"inputs\":[{\"name\":\"dest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutChangePaid\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"userAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"change\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutDeposit\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutRefunded\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegOutUserRefunded\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"userAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidDestination\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"actual\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidQuoteHash\",\"inputs\":[{\"name\":\"expected\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actual\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"MalformedTransaction\",\"inputs\":[{\"name\":\"outputScript\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NotEnoughConfirmations\",\"inputs\":[{\"name\":\"required\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"current\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"type\":\"error\",\"name\":\"QuoteAlreadyCompleted\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"QuoteAlreadyRegistered\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"QuoteExpiredByBlocks\",\"inputs\":[{\"name\":\"expireBlock\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"QuoteExpiredByTime\",\"inputs\":[{\"name\":\"depositDateLimit\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"expireDate\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"type\":\"error\",\"name\":\"QuoteNotExpired\",\"inputs\":[{\"name\":\"quoteHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnableToGetConfirmations\",\"inputs\":[{\"name\":\"errorCode\",\"type\":\"int256\",\"internalType\":\"int256\"}]}]",
	ID:  "PegoutContract",
}

// PegoutContract is an auto generated Go binding around an Ethereum contract.
type PegoutContract struct {
	abi abi.ABI
}

// NewPegoutContract creates a new instance of PegoutContract.
func NewPegoutContract() *PegoutContract {
	parsed, err := PegoutContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PegoutContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PegoutContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackDepositPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x083bc4b2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function depositPegOut((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes signature) payable returns()
func (pegoutContract *PegoutContract) PackDepositPegOut(quote QuotesPegOutQuote, signature []byte) []byte {
	enc, err := pegoutContract.abi.Pack("depositPegOut", quote, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDepositPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x083bc4b2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function depositPegOut((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote, bytes signature) payable returns()
func (pegoutContract *PegoutContract) TryPackDepositPegOut(quote QuotesPegOutQuote, signature []byte) ([]byte, error) {
	return pegoutContract.abi.Pack("depositPegOut", quote, signature)
}

// PackEip712Domain is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84b0196e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (pegoutContract *PegoutContract) PackEip712Domain() []byte {
	enc, err := pegoutContract.abi.Pack("eip712Domain")
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
func (pegoutContract *PegoutContract) TryPackEip712Domain() ([]byte, error) {
	return pegoutContract.abi.Pack("eip712Domain")
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
func (pegoutContract *PegoutContract) UnpackEip712Domain(data []byte) (Eip712DomainOutput, error) {
	out, err := pegoutContract.abi.Unpack("eip712Domain", data)
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
func (pegoutContract *PegoutContract) PackGetBalance(addr common.Address) []byte {
	enc, err := pegoutContract.abi.Pack("getBalance", addr)
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
func (pegoutContract *PegoutContract) TryPackGetBalance(addr common.Address) ([]byte, error) {
	return pegoutContract.abi.Pack("getBalance", addr)
}

// UnpackGetBalance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf8b2cb4f.
//
// Solidity: function getBalance(address addr) view returns(uint256)
func (pegoutContract *PegoutContract) UnpackGetBalance(data []byte) (*big.Int, error) {
	out, err := pegoutContract.abi.Unpack("getBalance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackHashPegOutQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6408f6fe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) PackHashPegOutQuote(quote QuotesPegOutQuote) []byte {
	enc, err := pegoutContract.abi.Pack("hashPegOutQuote", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashPegOutQuote is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6408f6fe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) TryPackHashPegOutQuote(quote QuotesPegOutQuote) ([]byte, error) {
	return pegoutContract.abi.Pack("hashPegOutQuote", quote)
}

// UnpackHashPegOutQuote is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6408f6fe.
//
// Solidity: function hashPegOutQuote((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) UnpackHashPegOutQuote(data []byte) ([32]byte, error) {
	out, err := pegoutContract.abi.Unpack("hashPegOutQuote", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackHashPegOutQuoteEIP712 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5966252a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) PackHashPegOutQuoteEIP712(quote QuotesPegOutQuote) []byte {
	enc, err := pegoutContract.abi.Pack("hashPegOutQuoteEIP712", quote)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashPegOutQuoteEIP712 is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5966252a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) TryPackHashPegOutQuoteEIP712(quote QuotesPegOutQuote) ([]byte, error) {
	return pegoutContract.abi.Pack("hashPegOutQuoteEIP712", quote)
}

// UnpackHashPegOutQuoteEIP712 is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5966252a.
//
// Solidity: function hashPegOutQuoteEIP712((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote) view returns(bytes32)
func (pegoutContract *PegoutContract) UnpackHashPegOutQuoteEIP712(data []byte) ([32]byte, error) {
	out, err := pegoutContract.abi.Unpack("hashPegOutQuoteEIP712", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackIsQuoteCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35bf61f1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (pegoutContract *PegoutContract) PackIsQuoteCompleted(quoteHash [32]byte) []byte {
	enc, err := pegoutContract.abi.Pack("isQuoteCompleted", quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsQuoteCompleted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x35bf61f1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (pegoutContract *PegoutContract) TryPackIsQuoteCompleted(quoteHash [32]byte) ([]byte, error) {
	return pegoutContract.abi.Pack("isQuoteCompleted", quoteHash)
}

// UnpackIsQuoteCompleted is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x35bf61f1.
//
// Solidity: function isQuoteCompleted(bytes32 quoteHash) view returns(bool)
func (pegoutContract *PegoutContract) UnpackIsQuoteCompleted(data []byte) (bool, error) {
	out, err := pegoutContract.abi.Unpack("isQuoteCompleted", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackPauseStatus is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x466916ca.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function pauseStatus() view returns(bool isPaused, string reason, uint64 since)
func (pegoutContract *PegoutContract) PackPauseStatus() []byte {
	enc, err := pegoutContract.abi.Pack("pauseStatus")
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
func (pegoutContract *PegoutContract) TryPackPauseStatus() ([]byte, error) {
	return pegoutContract.abi.Pack("pauseStatus")
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
func (pegoutContract *PegoutContract) UnpackPauseStatus(data []byte) (PauseStatusOutput, error) {
	out, err := pegoutContract.abi.Unpack("pauseStatus", data)
	outstruct := new(PauseStatusOutput)
	if err != nil {
		return *outstruct, err
	}
	outstruct.IsPaused = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Since = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	return *outstruct, nil
}

// PackRefundPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6c70de8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refundPegOut(bytes32 quoteHash, bytes btcTx, bytes32 btcBlockHeaderHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (pegoutContract *PegoutContract) PackRefundPegOut(quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) []byte {
	enc, err := pegoutContract.abi.Pack("refundPegOut", quoteHash, btcTx, btcBlockHeaderHash, merkleBranchPath, merkleBranchHashes)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefundPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6c70de8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refundPegOut(bytes32 quoteHash, bytes btcTx, bytes32 btcBlockHeaderHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) returns()
func (pegoutContract *PegoutContract) TryPackRefundPegOut(quoteHash [32]byte, btcTx []byte, btcBlockHeaderHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) ([]byte, error) {
	return pegoutContract.abi.Pack("refundPegOut", quoteHash, btcTx, btcBlockHeaderHash, merkleBranchPath, merkleBranchHashes)
}

// PackRefundUserPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f91797d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refundUserPegOut(bytes32 quoteHash) returns()
func (pegoutContract *PegoutContract) PackRefundUserPegOut(quoteHash [32]byte) []byte {
	enc, err := pegoutContract.abi.Pack("refundUserPegOut", quoteHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefundUserPegOut is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8f91797d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refundUserPegOut(bytes32 quoteHash) returns()
func (pegoutContract *PegoutContract) TryPackRefundUserPegOut(quoteHash [32]byte) ([]byte, error) {
	return pegoutContract.abi.Pack("refundUserPegOut", quoteHash)
}

// PackValidatePegout is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7846150c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (pegoutContract *PegoutContract) PackValidatePegout(quoteHash [32]byte, btcTx []byte) []byte {
	enc, err := pegoutContract.abi.Pack("validatePegout", quoteHash, btcTx)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackValidatePegout is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7846150c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (pegoutContract *PegoutContract) TryPackValidatePegout(quoteHash [32]byte, btcTx []byte) ([]byte, error) {
	return pegoutContract.abi.Pack("validatePegout", quoteHash, btcTx)
}

// UnpackValidatePegout is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7846150c.
//
// Solidity: function validatePegout(bytes32 quoteHash, bytes btcTx) view returns((uint256,uint256,uint256,uint256,uint256,address,address,address,int64,uint32,uint32,uint32,uint32,uint32,uint16,uint16,bytes,bytes,bytes) quote)
func (pegoutContract *PegoutContract) UnpackValidatePegout(data []byte) (QuotesPegOutQuote, error) {
	out, err := pegoutContract.abi.Unpack("validatePegout", data)
	if err != nil {
		return *new(QuotesPegOutQuote), err
	}
	out0 := *abi.ConvertType(out[0], new(QuotesPegOutQuote)).(*QuotesPegOutQuote)
	return out0, nil
}

// PackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3fef3a3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function withdraw(address addr, uint256 amount) returns()
func (pegoutContract *PegoutContract) PackWithdraw(addr common.Address, amount *big.Int) []byte {
	enc, err := pegoutContract.abi.Pack("withdraw", addr, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackWithdraw is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf3fef3a3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function withdraw(address addr, uint256 amount) returns()
func (pegoutContract *PegoutContract) TryPackWithdraw(addr common.Address, amount *big.Int) ([]byte, error) {
	return pegoutContract.abi.Pack("withdraw", addr, amount)
}

// PegoutContractBalanceDecrease represents a BalanceDecrease event raised by the PegoutContract contract.
type PegoutContractBalanceDecrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PegoutContractBalanceDecreaseEventName = "BalanceDecrease"

// ContractEventName returns the user-defined event name.
func (PegoutContractBalanceDecrease) ContractEventName() string {
	return PegoutContractBalanceDecreaseEventName
}

// UnpackBalanceDecreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceDecrease(address indexed dest, uint256 indexed amount)
func (pegoutContract *PegoutContract) UnpackBalanceDecreaseEvent(log *types.Log) (*PegoutContractBalanceDecrease, error) {
	event := "BalanceDecrease"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractBalanceDecrease)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractBalanceIncrease represents a BalanceIncrease event raised by the PegoutContract contract.
type PegoutContractBalanceIncrease struct {
	Dest   common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PegoutContractBalanceIncreaseEventName = "BalanceIncrease"

// ContractEventName returns the user-defined event name.
func (PegoutContractBalanceIncrease) ContractEventName() string {
	return PegoutContractBalanceIncreaseEventName
}

// UnpackBalanceIncreaseEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event BalanceIncrease(address indexed dest, uint256 indexed amount)
func (pegoutContract *PegoutContract) UnpackBalanceIncreaseEvent(log *types.Log) (*PegoutContractBalanceIncrease, error) {
	event := "BalanceIncrease"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractBalanceIncrease)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractEIP712DomainChanged represents a EIP712DomainChanged event raised by the PegoutContract contract.
type PegoutContractEIP712DomainChanged struct {
	Raw *types.Log // Blockchain specific contextual infos
}

const PegoutContractEIP712DomainChangedEventName = "EIP712DomainChanged"

// ContractEventName returns the user-defined event name.
func (PegoutContractEIP712DomainChanged) ContractEventName() string {
	return PegoutContractEIP712DomainChangedEventName
}

// UnpackEIP712DomainChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event EIP712DomainChanged()
func (pegoutContract *PegoutContract) UnpackEIP712DomainChangedEvent(log *types.Log) (*PegoutContractEIP712DomainChanged, error) {
	event := "EIP712DomainChanged"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractEIP712DomainChanged)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractPegOutChangePaid represents a PegOutChangePaid event raised by the PegoutContract contract.
type PegoutContractPegOutChangePaid struct {
	QuoteHash   [32]byte
	UserAddress common.Address
	Change      *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const PegoutContractPegOutChangePaidEventName = "PegOutChangePaid"

// ContractEventName returns the user-defined event name.
func (PegoutContractPegOutChangePaid) ContractEventName() string {
	return PegoutContractPegOutChangePaidEventName
}

// UnpackPegOutChangePaidEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutChangePaid(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed change)
func (pegoutContract *PegoutContract) UnpackPegOutChangePaidEvent(log *types.Log) (*PegoutContractPegOutChangePaid, error) {
	event := "PegOutChangePaid"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractPegOutChangePaid)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractPegOutDeposit represents a PegOutDeposit event raised by the PegoutContract contract.
type PegoutContractPegOutDeposit struct {
	QuoteHash [32]byte
	Sender    common.Address
	Timestamp *big.Int
	Amount    *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const PegoutContractPegOutDepositEventName = "PegOutDeposit"

// ContractEventName returns the user-defined event name.
func (PegoutContractPegOutDeposit) ContractEventName() string {
	return PegoutContractPegOutDepositEventName
}

// UnpackPegOutDepositEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 indexed timestamp, uint256 amount)
func (pegoutContract *PegoutContract) UnpackPegOutDepositEvent(log *types.Log) (*PegoutContractPegOutDeposit, error) {
	event := "PegOutDeposit"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractPegOutDeposit)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractPegOutRefunded represents a PegOutRefunded event raised by the PegoutContract contract.
type PegoutContractPegOutRefunded struct {
	QuoteHash [32]byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const PegoutContractPegOutRefundedEventName = "PegOutRefunded"

// ContractEventName returns the user-defined event name.
func (PegoutContractPegOutRefunded) ContractEventName() string {
	return PegoutContractPegOutRefundedEventName
}

// UnpackPegOutRefundedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutRefunded(bytes32 indexed quoteHash)
func (pegoutContract *PegoutContract) UnpackPegOutRefundedEvent(log *types.Log) (*PegoutContractPegOutRefunded, error) {
	event := "PegOutRefunded"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractPegOutRefunded)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractPegOutUserRefunded represents a PegOutUserRefunded event raised by the PegoutContract contract.
type PegoutContractPegOutUserRefunded struct {
	QuoteHash   [32]byte
	UserAddress common.Address
	Value       *big.Int
	Raw         *types.Log // Blockchain specific contextual infos
}

const PegoutContractPegOutUserRefundedEventName = "PegOutUserRefunded"

// ContractEventName returns the user-defined event name.
func (PegoutContractPegOutUserRefunded) ContractEventName() string {
	return PegoutContractPegOutUserRefundedEventName
}

// UnpackPegOutUserRefundedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegOutUserRefunded(bytes32 indexed quoteHash, address indexed userAddress, uint256 indexed value)
func (pegoutContract *PegoutContract) UnpackPegOutUserRefundedEvent(log *types.Log) (*PegoutContractPegOutUserRefunded, error) {
	event := "PegOutUserRefunded"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractPegOutUserRefunded)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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

// PegoutContractWithdrawal represents a Withdrawal event raised by the PegoutContract contract.
type PegoutContractWithdrawal struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const PegoutContractWithdrawalEventName = "Withdrawal"

// ContractEventName returns the user-defined event name.
func (PegoutContractWithdrawal) ContractEventName() string {
	return PegoutContractWithdrawalEventName
}

// UnpackWithdrawalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Withdrawal(address indexed from, address indexed to, uint256 indexed amount)
func (pegoutContract *PegoutContract) UnpackWithdrawalEvent(log *types.Log) (*PegoutContractWithdrawal, error) {
	event := "Withdrawal"
	if len(log.Topics) == 0 || log.Topics[0] != pegoutContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PegoutContractWithdrawal)
	if len(log.Data) > 0 {
		if err := pegoutContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range pegoutContract.abi.Events[event].Inputs {
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
func (pegoutContract *PegoutContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["InvalidDestination"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackInvalidDestinationError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["InvalidQuoteHash"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackInvalidQuoteHashError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["MalformedTransaction"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackMalformedTransactionError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["NotEnoughConfirmations"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackNotEnoughConfirmationsError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["QuoteAlreadyCompleted"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackQuoteAlreadyCompletedError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["QuoteAlreadyRegistered"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackQuoteAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["QuoteExpiredByBlocks"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackQuoteExpiredByBlocksError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["QuoteExpiredByTime"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackQuoteExpiredByTimeError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["QuoteNotExpired"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackQuoteNotExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], pegoutContract.abi.Errors["UnableToGetConfirmations"].ID.Bytes()[:4]) {
		return pegoutContract.UnpackUnableToGetConfirmationsError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PegoutContractInvalidDestination represents a InvalidDestination error raised by the PegoutContract contract.
type PegoutContractInvalidDestination struct {
	Expected []byte
	Actual   []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidDestination(bytes expected, bytes actual)
func PegoutContractInvalidDestinationErrorID() common.Hash {
	return common.HexToHash("0x7c5722b2339c1bac07667c7c774049367da0fdd6712c59c7f4d37cb4c4ed9bc9")
}

// UnpackInvalidDestinationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidDestination(bytes expected, bytes actual)
func (pegoutContract *PegoutContract) UnpackInvalidDestinationError(raw []byte) (*PegoutContractInvalidDestination, error) {
	out := new(PegoutContractInvalidDestination)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "InvalidDestination", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractInvalidQuoteHash represents a InvalidQuoteHash error raised by the PegoutContract contract.
type PegoutContractInvalidQuoteHash struct {
	Expected [32]byte
	Actual   [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidQuoteHash(bytes32 expected, bytes32 actual)
func PegoutContractInvalidQuoteHashErrorID() common.Hash {
	return common.HexToHash("0x7826d5fac496a7d2b822d6be19026f7c9d348dbcb9b7a093d3e68c466ec2dafd")
}

// UnpackInvalidQuoteHashError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidQuoteHash(bytes32 expected, bytes32 actual)
func (pegoutContract *PegoutContract) UnpackInvalidQuoteHashError(raw []byte) (*PegoutContractInvalidQuoteHash, error) {
	out := new(PegoutContractInvalidQuoteHash)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "InvalidQuoteHash", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractMalformedTransaction represents a MalformedTransaction error raised by the PegoutContract contract.
type PegoutContractMalformedTransaction struct {
	OutputScript []byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error MalformedTransaction(bytes outputScript)
func PegoutContractMalformedTransactionErrorID() common.Hash {
	return common.HexToHash("0x7201f86d13a7581b1c224e80b1bab6463f79bfc02165d6aa3ef3fe920ec19562")
}

// UnpackMalformedTransactionError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error MalformedTransaction(bytes outputScript)
func (pegoutContract *PegoutContract) UnpackMalformedTransactionError(raw []byte) (*PegoutContractMalformedTransaction, error) {
	out := new(PegoutContractMalformedTransaction)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "MalformedTransaction", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractNotEnoughConfirmations represents a NotEnoughConfirmations error raised by the PegoutContract contract.
type PegoutContractNotEnoughConfirmations struct {
	Required *big.Int
	Current  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotEnoughConfirmations(int256 required, int256 current)
func PegoutContractNotEnoughConfirmationsErrorID() common.Hash {
	return common.HexToHash("0xd2506f8c9e9ad186c7e4f7debe041c835ed9b36457957ead5160c46f1a622788")
}

// UnpackNotEnoughConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotEnoughConfirmations(int256 required, int256 current)
func (pegoutContract *PegoutContract) UnpackNotEnoughConfirmationsError(raw []byte) (*PegoutContractNotEnoughConfirmations, error) {
	out := new(PegoutContractNotEnoughConfirmations)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "NotEnoughConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractQuoteAlreadyCompleted represents a QuoteAlreadyCompleted error raised by the PegoutContract contract.
type PegoutContractQuoteAlreadyCompleted struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteAlreadyCompleted(bytes32 quoteHash)
func PegoutContractQuoteAlreadyCompletedErrorID() common.Hash {
	return common.HexToHash("0x86ef90ba285627a7f8550084d34a1c1c7b16b3bce815576ae0f444b9c09dea34")
}

// UnpackQuoteAlreadyCompletedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteAlreadyCompleted(bytes32 quoteHash)
func (pegoutContract *PegoutContract) UnpackQuoteAlreadyCompletedError(raw []byte) (*PegoutContractQuoteAlreadyCompleted, error) {
	out := new(PegoutContractQuoteAlreadyCompleted)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "QuoteAlreadyCompleted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractQuoteAlreadyRegistered represents a QuoteAlreadyRegistered error raised by the PegoutContract contract.
type PegoutContractQuoteAlreadyRegistered struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteAlreadyRegistered(bytes32 quoteHash)
func PegoutContractQuoteAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x5bf1b7fb72ebc6fd6814b3e552d529c3af61c4bca358a54bd4765ab4cd6d6f93")
}

// UnpackQuoteAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteAlreadyRegistered(bytes32 quoteHash)
func (pegoutContract *PegoutContract) UnpackQuoteAlreadyRegisteredError(raw []byte) (*PegoutContractQuoteAlreadyRegistered, error) {
	out := new(PegoutContractQuoteAlreadyRegistered)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "QuoteAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractQuoteExpiredByBlocks represents a QuoteExpiredByBlocks error raised by the PegoutContract contract.
type PegoutContractQuoteExpiredByBlocks struct {
	ExpireBlock uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteExpiredByBlocks(uint32 expireBlock)
func PegoutContractQuoteExpiredByBlocksErrorID() common.Hash {
	return common.HexToHash("0x1dc8f1505efa9c351cbbee012d8ca856e0aaffe5f28c0bd5116c8ea1db449d8f")
}

// UnpackQuoteExpiredByBlocksError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteExpiredByBlocks(uint32 expireBlock)
func (pegoutContract *PegoutContract) UnpackQuoteExpiredByBlocksError(raw []byte) (*PegoutContractQuoteExpiredByBlocks, error) {
	out := new(PegoutContractQuoteExpiredByBlocks)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "QuoteExpiredByBlocks", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractQuoteExpiredByTime represents a QuoteExpiredByTime error raised by the PegoutContract contract.
type PegoutContractQuoteExpiredByTime struct {
	DepositDateLimit uint32
	ExpireDate       uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteExpiredByTime(uint32 depositDateLimit, uint32 expireDate)
func PegoutContractQuoteExpiredByTimeErrorID() common.Hash {
	return common.HexToHash("0x289e847e8f4c12d2213053147fe76ff2247863c2ea21a6dff757458f6c315c5f")
}

// UnpackQuoteExpiredByTimeError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteExpiredByTime(uint32 depositDateLimit, uint32 expireDate)
func (pegoutContract *PegoutContract) UnpackQuoteExpiredByTimeError(raw []byte) (*PegoutContractQuoteExpiredByTime, error) {
	out := new(PegoutContractQuoteExpiredByTime)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "QuoteExpiredByTime", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractQuoteNotExpired represents a QuoteNotExpired error raised by the PegoutContract contract.
type PegoutContractQuoteNotExpired struct {
	QuoteHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error QuoteNotExpired(bytes32 quoteHash)
func PegoutContractQuoteNotExpiredErrorID() common.Hash {
	return common.HexToHash("0xe7af165007a17453578995b3bdbe29fc875ac731917d0190d1cdca6b3a68951b")
}

// UnpackQuoteNotExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error QuoteNotExpired(bytes32 quoteHash)
func (pegoutContract *PegoutContract) UnpackQuoteNotExpiredError(raw []byte) (*PegoutContractQuoteNotExpired, error) {
	out := new(PegoutContractQuoteNotExpired)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "QuoteNotExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PegoutContractUnableToGetConfirmations represents a UnableToGetConfirmations error raised by the PegoutContract contract.
type PegoutContractUnableToGetConfirmations struct {
	ErrorCode *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnableToGetConfirmations(int256 errorCode)
func PegoutContractUnableToGetConfirmationsErrorID() common.Hash {
	return common.HexToHash("0xd06e366aa562b88cc28592d73c0fea0b2c6f7f7544f39af59c0d16c892dd5ad4")
}

// UnpackUnableToGetConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnableToGetConfirmations(int256 errorCode)
func (pegoutContract *PegoutContract) UnpackUnableToGetConfirmationsError(raw []byte) (*PegoutContractUnableToGetConfirmations, error) {
	out := new(PegoutContractUnableToGetConfirmations)
	if err := pegoutContract.abi.UnpackIntoInterface(out, "UnableToGetConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}
