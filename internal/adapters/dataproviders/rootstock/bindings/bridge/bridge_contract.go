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

// RskBridgeMetaData contains all meta data concerning the RskBridge contract.
var RskBridgeMetaData = bind.MetaData{
	ABI: "[{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addFederatorPublicKey\",\"inputs\":[{\"name\":\"key\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addFederatorPublicKeyMultikey\",\"inputs\":[{\"name\":\"btcKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"rskKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mstKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addLockWhitelistAddress\",\"inputs\":[{\"name\":\"aaddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"maxTransferValue\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addOneOffLockWhitelistAddress\",\"inputs\":[{\"name\":\"aaddress\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"maxTransferValue\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addSignature\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"txhash\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addUnlimitedLockWhitelistAddress\",\"inputs\":[{\"name\":\"aaddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitFederation\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createFederation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getActiveFederationCreationBlockHeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActivePowpegRedeemScript\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainBestBlockHeader\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainBestChainHeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainBlockHashAtDepth\",\"inputs\":[{\"name\":\"depth\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainBlockHeaderByHash\",\"inputs\":[{\"name\":\"btcBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainBlockHeaderByHeight\",\"inputs\":[{\"name\":\"btcBlockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainInitialBlockHeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcBlockchainParentBlockHeaderByHash\",\"inputs\":[{\"name\":\"btcBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcTransactionConfirmations\",\"inputs\":[{\"name\":\"txHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleBranchPath\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBtcTxHashProcessedHeight\",\"inputs\":[{\"name\":\"hash\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int64\",\"internalType\":\"int64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederationAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederationCreationBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederationCreationTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederationSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederationThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederatorPublicKey\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFederatorPublicKeyOfType\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"atype\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeePerKb\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockWhitelistAddress\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockWhitelistEntryByAddress\",\"inputs\":[{\"name\":\"aaddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockWhitelistSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLockingCap\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinimumLockTxValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingFederationHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingFederationSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingFederatorPublicKey\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingFederatorPublicKeyOfType\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"atype\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederationAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederationCreationBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederationCreationTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederationSize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederationThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederatorPublicKey\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRetiringFederatorPublicKeyOfType\",\"inputs\":[{\"name\":\"index\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"atype\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStateForBtcReleaseClient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStateForDebugging\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasBtcBlockCoinbaseTransactionInformation\",\"inputs\":[{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseLockingCap\",\"inputs\":[{\"name\":\"newLockingCap\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isBtcTxHashAlreadyProcessed\",\"inputs\":[{\"name\":\"hash\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"receiveHeader\",\"inputs\":[{\"name\":\"ablock\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"receiveHeaders\",\"inputs\":[{\"name\":\"blocks\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerBtcCoinbaseTransaction\",\"inputs\":[{\"name\":\"btcTxSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pmtSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"witnessMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"witnessReservedValue\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerBtcTransaction\",\"inputs\":[{\"name\":\"atx\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"pmt\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerFastBridgeBtcTransaction\",\"inputs\":[{\"name\":\"btcTxSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pmtSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"derivationArgumentsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"userRefundBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"liquidityBridgeContractAddress\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"liquidityProviderBtcAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"shouldTransferToContract\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeLockWhitelistAddress\",\"inputs\":[{\"name\":\"aaddress\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rollbackFederation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLockWhitelistDisableBlockDelay\",\"inputs\":[{\"name\":\"disableDelay\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCollections\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"voteFeePerKbChange\",\"inputs\":[{\"name\":\"feePerKb\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"batch_pegout_created\",\"inputs\":[{\"name\":\"btcTx\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"releaseRskTxHashes\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	ID:  "RskBridge",
}

// RskBridge is an auto generated Go binding around an Ethereum contract.
type RskBridge struct {
	abi abi.ABI
}

// NewRskBridge creates a new instance of RskBridge.
func NewRskBridge() *RskBridge {
	parsed, err := RskBridgeMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &RskBridge{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *RskBridge) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecefd339.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (rskBridge *RskBridge) PackAddFederatorPublicKey(key []byte) []byte {
	enc, err := rskBridge.abi.Pack("addFederatorPublicKey", key)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xecefd339.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (rskBridge *RskBridge) TryPackAddFederatorPublicKey(key []byte) ([]byte, error) {
	return rskBridge.abi.Pack("addFederatorPublicKey", key)
}

// UnpackAddFederatorPublicKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xecefd339.
//
// Solidity: function addFederatorPublicKey(bytes key) returns(int256)
func (rskBridge *RskBridge) UnpackAddFederatorPublicKey(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("addFederatorPublicKey", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAddFederatorPublicKeyMultikey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x444ff9da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (rskBridge *RskBridge) PackAddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) []byte {
	enc, err := rskBridge.abi.Pack("addFederatorPublicKeyMultikey", btcKey, rskKey, mstKey)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddFederatorPublicKeyMultikey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x444ff9da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (rskBridge *RskBridge) TryPackAddFederatorPublicKeyMultikey(btcKey []byte, rskKey []byte, mstKey []byte) ([]byte, error) {
	return rskBridge.abi.Pack("addFederatorPublicKeyMultikey", btcKey, rskKey, mstKey)
}

// UnpackAddFederatorPublicKeyMultikey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x444ff9da.
//
// Solidity: function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns(int256)
func (rskBridge *RskBridge) UnpackAddFederatorPublicKeyMultikey(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("addFederatorPublicKeyMultikey", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAddLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x502bbbce.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) PackAddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("addLockWhitelistAddress", aaddress, maxTransferValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x502bbbce.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) TryPackAddLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("addLockWhitelistAddress", aaddress, maxTransferValue)
}

// UnpackAddLockWhitelistAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x502bbbce.
//
// Solidity: function addLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) UnpackAddLockWhitelistAddress(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("addLockWhitelistAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAddOneOffLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x848206d9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) PackAddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("addOneOffLockWhitelistAddress", aaddress, maxTransferValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddOneOffLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x848206d9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) TryPackAddOneOffLockWhitelistAddress(aaddress string, maxTransferValue *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("addOneOffLockWhitelistAddress", aaddress, maxTransferValue)
}

// UnpackAddOneOffLockWhitelistAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x848206d9.
//
// Solidity: function addOneOffLockWhitelistAddress(string aaddress, int256 maxTransferValue) returns(int256)
func (rskBridge *RskBridge) UnpackAddOneOffLockWhitelistAddress(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("addOneOffLockWhitelistAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackAddSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf10b9c59.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (rskBridge *RskBridge) PackAddSignature(pubkey []byte, signatures [][]byte, txhash []byte) []byte {
	enc, err := rskBridge.abi.Pack("addSignature", pubkey, signatures, txhash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddSignature is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf10b9c59.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addSignature(bytes pubkey, bytes[] signatures, bytes txhash) returns()
func (rskBridge *RskBridge) TryPackAddSignature(pubkey []byte, signatures [][]byte, txhash []byte) ([]byte, error) {
	return rskBridge.abi.Pack("addSignature", pubkey, signatures, txhash)
}

// PackAddUnlimitedLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb906c938.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) PackAddUnlimitedLockWhitelistAddress(aaddress string) []byte {
	enc, err := rskBridge.abi.Pack("addUnlimitedLockWhitelistAddress", aaddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddUnlimitedLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb906c938.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) TryPackAddUnlimitedLockWhitelistAddress(aaddress string) ([]byte, error) {
	return rskBridge.abi.Pack("addUnlimitedLockWhitelistAddress", aaddress)
}

// UnpackAddUnlimitedLockWhitelistAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb906c938.
//
// Solidity: function addUnlimitedLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) UnpackAddUnlimitedLockWhitelistAddress(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("addUnlimitedLockWhitelistAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCommitFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1533330f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (rskBridge *RskBridge) PackCommitFederation(hash []byte) []byte {
	enc, err := rskBridge.abi.Pack("commitFederation", hash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCommitFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1533330f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (rskBridge *RskBridge) TryPackCommitFederation(hash []byte) ([]byte, error) {
	return rskBridge.abi.Pack("commitFederation", hash)
}

// UnpackCommitFederation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1533330f.
//
// Solidity: function commitFederation(bytes hash) returns(int256)
func (rskBridge *RskBridge) UnpackCommitFederation(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("commitFederation", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackCreateFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1183d5d1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function createFederation() returns(int256)
func (rskBridge *RskBridge) PackCreateFederation() []byte {
	enc, err := rskBridge.abi.Pack("createFederation")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCreateFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1183d5d1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function createFederation() returns(int256)
func (rskBridge *RskBridge) TryPackCreateFederation() ([]byte, error) {
	return rskBridge.abi.Pack("createFederation")
}

// UnpackCreateFederation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1183d5d1.
//
// Solidity: function createFederation() returns(int256)
func (rskBridge *RskBridge) UnpackCreateFederation(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("createFederation", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetActiveFederationCreationBlockHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x177d6e18.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (rskBridge *RskBridge) PackGetActiveFederationCreationBlockHeight() []byte {
	enc, err := rskBridge.abi.Pack("getActiveFederationCreationBlockHeight")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetActiveFederationCreationBlockHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x177d6e18.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (rskBridge *RskBridge) TryPackGetActiveFederationCreationBlockHeight() ([]byte, error) {
	return rskBridge.abi.Pack("getActiveFederationCreationBlockHeight")
}

// UnpackGetActiveFederationCreationBlockHeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x177d6e18.
//
// Solidity: function getActiveFederationCreationBlockHeight() view returns(uint256)
func (rskBridge *RskBridge) UnpackGetActiveFederationCreationBlockHeight(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getActiveFederationCreationBlockHeight", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetActivePowpegRedeemScript is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d73d5dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (rskBridge *RskBridge) PackGetActivePowpegRedeemScript() []byte {
	enc, err := rskBridge.abi.Pack("getActivePowpegRedeemScript")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetActivePowpegRedeemScript is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1d73d5dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (rskBridge *RskBridge) TryPackGetActivePowpegRedeemScript() ([]byte, error) {
	return rskBridge.abi.Pack("getActivePowpegRedeemScript")
}

// UnpackGetActivePowpegRedeemScript is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1d73d5dd.
//
// Solidity: function getActivePowpegRedeemScript() view returns(bytes)
func (rskBridge *RskBridge) UnpackGetActivePowpegRedeemScript(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getActivePowpegRedeemScript", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcBlockchainBestBlockHeader is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0b2424b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (rskBridge *RskBridge) PackGetBtcBlockchainBestBlockHeader() []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainBestBlockHeader")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainBestBlockHeader is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf0b2424b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainBestBlockHeader() ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainBestBlockHeader")
}

// UnpackGetBtcBlockchainBestBlockHeader is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf0b2424b.
//
// Solidity: function getBtcBlockchainBestBlockHeader() view returns(bytes)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainBestBlockHeader(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainBestBlockHeader", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcBlockchainBestChainHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14c89c01.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (rskBridge *RskBridge) PackGetBtcBlockchainBestChainHeight() []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainBestChainHeight")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainBestChainHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x14c89c01.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainBestChainHeight() ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainBestChainHeight")
}

// UnpackGetBtcBlockchainBestChainHeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x14c89c01.
//
// Solidity: function getBtcBlockchainBestChainHeight() view returns(int256)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainBestChainHeight(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainBestChainHeight", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetBtcBlockchainBlockHashAtDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefd44418.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (rskBridge *RskBridge) PackGetBtcBlockchainBlockHashAtDepth(depth *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainBlockHashAtDepth", depth)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainBlockHashAtDepth is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xefd44418.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainBlockHashAtDepth(depth *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainBlockHashAtDepth", depth)
}

// UnpackGetBtcBlockchainBlockHashAtDepth is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xefd44418.
//
// Solidity: function getBtcBlockchainBlockHashAtDepth(int256 depth) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainBlockHashAtDepth(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainBlockHashAtDepth", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcBlockchainBlockHeaderByHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x739e364a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) PackGetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainBlockHeaderByHash", btcBlockHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainBlockHeaderByHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x739e364a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainBlockHeaderByHash", btcBlockHash)
}

// UnpackGetBtcBlockchainBlockHeaderByHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x739e364a.
//
// Solidity: function getBtcBlockchainBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainBlockHeaderByHash(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainBlockHeaderByHash", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcBlockchainBlockHeaderByHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd0c1fff.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (rskBridge *RskBridge) PackGetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainBlockHeaderByHeight", btcBlockHeight)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainBlockHeaderByHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbd0c1fff.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainBlockHeaderByHeight(btcBlockHeight *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainBlockHeaderByHeight", btcBlockHeight)
}

// UnpackGetBtcBlockchainBlockHeaderByHeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbd0c1fff.
//
// Solidity: function getBtcBlockchainBlockHeaderByHeight(uint256 btcBlockHeight) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainBlockHeaderByHeight(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainBlockHeaderByHeight", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcBlockchainInitialBlockHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4897193f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (rskBridge *RskBridge) PackGetBtcBlockchainInitialBlockHeight() []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainInitialBlockHeight")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainInitialBlockHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4897193f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainInitialBlockHeight() ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainInitialBlockHeight")
}

// UnpackGetBtcBlockchainInitialBlockHeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4897193f.
//
// Solidity: function getBtcBlockchainInitialBlockHeight() view returns(int256)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainInitialBlockHeight(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainInitialBlockHeight", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetBtcBlockchainParentBlockHeaderByHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0236724.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) PackGetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) []byte {
	enc, err := rskBridge.abi.Pack("getBtcBlockchainParentBlockHeaderByHash", btcBlockHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcBlockchainParentBlockHeaderByHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe0236724.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetBtcBlockchainParentBlockHeaderByHash(btcBlockHash [32]byte) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcBlockchainParentBlockHeaderByHash", btcBlockHash)
}

// UnpackGetBtcBlockchainParentBlockHeaderByHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe0236724.
//
// Solidity: function getBtcBlockchainParentBlockHeaderByHash(bytes32 btcBlockHash) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetBtcBlockchainParentBlockHeaderByHash(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getBtcBlockchainParentBlockHeaderByHash", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetBtcTransactionConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5b644587.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (rskBridge *RskBridge) PackGetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) []byte {
	enc, err := rskBridge.abi.Pack("getBtcTransactionConfirmations", txHash, blockHash, merkleBranchPath, merkleBranchHashes)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcTransactionConfirmations is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5b644587.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (rskBridge *RskBridge) TryPackGetBtcTransactionConfirmations(txHash [32]byte, blockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcTransactionConfirmations", txHash, blockHash, merkleBranchPath, merkleBranchHashes)
}

// UnpackGetBtcTransactionConfirmations is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5b644587.
//
// Solidity: function getBtcTransactionConfirmations(bytes32 txHash, bytes32 blockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) view returns(int256)
func (rskBridge *RskBridge) UnpackGetBtcTransactionConfirmations(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getBtcTransactionConfirmations", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetBtcTxHashProcessedHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97fcca7d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (rskBridge *RskBridge) PackGetBtcTxHashProcessedHeight(hash string) []byte {
	enc, err := rskBridge.abi.Pack("getBtcTxHashProcessedHeight", hash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetBtcTxHashProcessedHeight is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x97fcca7d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (rskBridge *RskBridge) TryPackGetBtcTxHashProcessedHeight(hash string) ([]byte, error) {
	return rskBridge.abi.Pack("getBtcTxHashProcessedHeight", hash)
}

// UnpackGetBtcTxHashProcessedHeight is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x97fcca7d.
//
// Solidity: function getBtcTxHashProcessedHeight(string hash) view returns(int64)
func (rskBridge *RskBridge) UnpackGetBtcTxHashProcessedHeight(data []byte) (int64, error) {
	out, err := rskBridge.abi.Unpack("getBtcTxHashProcessedHeight", data)
	if err != nil {
		return *new(int64), err
	}
	out0 := *abi.ConvertType(out[0], new(int64)).(*int64)
	return out0, nil
}

// PackGetFederationAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6923fa85.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederationAddress() view returns(string)
func (rskBridge *RskBridge) PackGetFederationAddress() []byte {
	enc, err := rskBridge.abi.Pack("getFederationAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederationAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6923fa85.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederationAddress() view returns(string)
func (rskBridge *RskBridge) TryPackGetFederationAddress() ([]byte, error) {
	return rskBridge.abi.Pack("getFederationAddress")
}

// UnpackGetFederationAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6923fa85.
//
// Solidity: function getFederationAddress() view returns(string)
func (rskBridge *RskBridge) UnpackGetFederationAddress(data []byte) (string, error) {
	out, err := rskBridge.abi.Unpack("getFederationAddress", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackGetFederationCreationBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b2045ee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) PackGetFederationCreationBlockNumber() []byte {
	enc, err := rskBridge.abi.Pack("getFederationCreationBlockNumber")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederationCreationBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1b2045ee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) TryPackGetFederationCreationBlockNumber() ([]byte, error) {
	return rskBridge.abi.Pack("getFederationCreationBlockNumber")
}

// UnpackGetFederationCreationBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1b2045ee.
//
// Solidity: function getFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) UnpackGetFederationCreationBlockNumber(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getFederationCreationBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetFederationCreationTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e2db9d4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) PackGetFederationCreationTime() []byte {
	enc, err := rskBridge.abi.Pack("getFederationCreationTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederationCreationTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5e2db9d4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) TryPackGetFederationCreationTime() ([]byte, error) {
	return rskBridge.abi.Pack("getFederationCreationTime")
}

// UnpackGetFederationCreationTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5e2db9d4.
//
// Solidity: function getFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) UnpackGetFederationCreationTime(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getFederationCreationTime", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x802ad4b6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederationSize() view returns(int256)
func (rskBridge *RskBridge) PackGetFederationSize() []byte {
	enc, err := rskBridge.abi.Pack("getFederationSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x802ad4b6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederationSize() view returns(int256)
func (rskBridge *RskBridge) TryPackGetFederationSize() ([]byte, error) {
	return rskBridge.abi.Pack("getFederationSize")
}

// UnpackGetFederationSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x802ad4b6.
//
// Solidity: function getFederationSize() view returns(int256)
func (rskBridge *RskBridge) UnpackGetFederationSize(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getFederationSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetFederationThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fd47456.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) PackGetFederationThreshold() []byte {
	enc, err := rskBridge.abi.Pack("getFederationThreshold")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederationThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0fd47456.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) TryPackGetFederationThreshold() ([]byte, error) {
	return rskBridge.abi.Pack("getFederationThreshold")
}

// UnpackGetFederationThreshold is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0fd47456.
//
// Solidity: function getFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) UnpackGetFederationThreshold(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getFederationThreshold", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b89a1af.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) PackGetFederatorPublicKey(index *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getFederatorPublicKey", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6b89a1af.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetFederatorPublicKey(index *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getFederatorPublicKey", index)
}

// UnpackGetFederatorPublicKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6b89a1af.
//
// Solidity: function getFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetFederatorPublicKey(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getFederatorPublicKey", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x549cfd1c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) PackGetFederatorPublicKeyOfType(index *big.Int, atype string) []byte {
	enc, err := rskBridge.abi.Pack("getFederatorPublicKeyOfType", index, atype)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x549cfd1c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return rskBridge.abi.Pack("getFederatorPublicKeyOfType", index, atype)
}

// UnpackGetFederatorPublicKeyOfType is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x549cfd1c.
//
// Solidity: function getFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetFederatorPublicKeyOfType(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getFederatorPublicKeyOfType", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetFeePerKb is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x724ec886.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getFeePerKb() view returns(int256)
func (rskBridge *RskBridge) PackGetFeePerKb() []byte {
	enc, err := rskBridge.abi.Pack("getFeePerKb")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetFeePerKb is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x724ec886.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getFeePerKb() view returns(int256)
func (rskBridge *RskBridge) TryPackGetFeePerKb() ([]byte, error) {
	return rskBridge.abi.Pack("getFeePerKb")
}

// UnpackGetFeePerKb is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x724ec886.
//
// Solidity: function getFeePerKb() view returns(int256)
func (rskBridge *RskBridge) UnpackGetFeePerKb(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getFeePerKb", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93988b76.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (rskBridge *RskBridge) PackGetLockWhitelistAddress(index *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getLockWhitelistAddress", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x93988b76.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (rskBridge *RskBridge) TryPackGetLockWhitelistAddress(index *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getLockWhitelistAddress", index)
}

// UnpackGetLockWhitelistAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x93988b76.
//
// Solidity: function getLockWhitelistAddress(int256 index) view returns(string)
func (rskBridge *RskBridge) UnpackGetLockWhitelistAddress(data []byte) (string, error) {
	out, err := rskBridge.abi.Unpack("getLockWhitelistAddress", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackGetLockWhitelistEntryByAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x251c5f7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (rskBridge *RskBridge) PackGetLockWhitelistEntryByAddress(aaddress string) []byte {
	enc, err := rskBridge.abi.Pack("getLockWhitelistEntryByAddress", aaddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetLockWhitelistEntryByAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x251c5f7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (rskBridge *RskBridge) TryPackGetLockWhitelistEntryByAddress(aaddress string) ([]byte, error) {
	return rskBridge.abi.Pack("getLockWhitelistEntryByAddress", aaddress)
}

// UnpackGetLockWhitelistEntryByAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x251c5f7b.
//
// Solidity: function getLockWhitelistEntryByAddress(string aaddress) view returns(int256)
func (rskBridge *RskBridge) UnpackGetLockWhitelistEntryByAddress(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getLockWhitelistEntryByAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetLockWhitelistSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9e658dc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (rskBridge *RskBridge) PackGetLockWhitelistSize() []byte {
	enc, err := rskBridge.abi.Pack("getLockWhitelistSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetLockWhitelistSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe9e658dc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (rskBridge *RskBridge) TryPackGetLockWhitelistSize() ([]byte, error) {
	return rskBridge.abi.Pack("getLockWhitelistSize")
}

// UnpackGetLockWhitelistSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe9e658dc.
//
// Solidity: function getLockWhitelistSize() view returns(int256)
func (rskBridge *RskBridge) UnpackGetLockWhitelistSize(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getLockWhitelistSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetLockingCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f9db977.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getLockingCap() view returns(int256)
func (rskBridge *RskBridge) PackGetLockingCap() []byte {
	enc, err := rskBridge.abi.Pack("getLockingCap")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetLockingCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f9db977.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getLockingCap() view returns(int256)
func (rskBridge *RskBridge) TryPackGetLockingCap() ([]byte, error) {
	return rskBridge.abi.Pack("getLockingCap")
}

// UnpackGetLockingCap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f9db977.
//
// Solidity: function getLockingCap() view returns(int256)
func (rskBridge *RskBridge) UnpackGetLockingCap(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getLockingCap", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetMinimumLockTxValue is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f8d158f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (rskBridge *RskBridge) PackGetMinimumLockTxValue() []byte {
	enc, err := rskBridge.abi.Pack("getMinimumLockTxValue")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetMinimumLockTxValue is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f8d158f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (rskBridge *RskBridge) TryPackGetMinimumLockTxValue() ([]byte, error) {
	return rskBridge.abi.Pack("getMinimumLockTxValue")
}

// UnpackGetMinimumLockTxValue is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2f8d158f.
//
// Solidity: function getMinimumLockTxValue() view returns(int256)
func (rskBridge *RskBridge) UnpackGetMinimumLockTxValue(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getMinimumLockTxValue", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPendingFederationHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6ce0ed5a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (rskBridge *RskBridge) PackGetPendingFederationHash() []byte {
	enc, err := rskBridge.abi.Pack("getPendingFederationHash")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingFederationHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6ce0ed5a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (rskBridge *RskBridge) TryPackGetPendingFederationHash() ([]byte, error) {
	return rskBridge.abi.Pack("getPendingFederationHash")
}

// UnpackGetPendingFederationHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6ce0ed5a.
//
// Solidity: function getPendingFederationHash() view returns(bytes)
func (rskBridge *RskBridge) UnpackGetPendingFederationHash(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getPendingFederationHash", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetPendingFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ac72b38.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (rskBridge *RskBridge) PackGetPendingFederationSize() []byte {
	enc, err := rskBridge.abi.Pack("getPendingFederationSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3ac72b38.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (rskBridge *RskBridge) TryPackGetPendingFederationSize() ([]byte, error) {
	return rskBridge.abi.Pack("getPendingFederationSize")
}

// UnpackGetPendingFederationSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3ac72b38.
//
// Solidity: function getPendingFederationSize() view returns(int256)
func (rskBridge *RskBridge) UnpackGetPendingFederationSize(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getPendingFederationSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetPendingFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x492f7c44.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) PackGetPendingFederatorPublicKey(index *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getPendingFederatorPublicKey", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x492f7c44.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetPendingFederatorPublicKey(index *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getPendingFederatorPublicKey", index)
}

// UnpackGetPendingFederatorPublicKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x492f7c44.
//
// Solidity: function getPendingFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetPendingFederatorPublicKey(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getPendingFederatorPublicKey", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetPendingFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc61295d9.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) PackGetPendingFederatorPublicKeyOfType(index *big.Int, atype string) []byte {
	enc, err := rskBridge.abi.Pack("getPendingFederatorPublicKeyOfType", index, atype)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetPendingFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc61295d9.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetPendingFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return rskBridge.abi.Pack("getPendingFederatorPublicKeyOfType", index, atype)
}

// UnpackGetPendingFederatorPublicKeyOfType is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc61295d9.
//
// Solidity: function getPendingFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetPendingFederatorPublicKeyOfType(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getPendingFederatorPublicKeyOfType", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetRetiringFederationAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47227286.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (rskBridge *RskBridge) PackGetRetiringFederationAddress() []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederationAddress")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederationAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x47227286.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (rskBridge *RskBridge) TryPackGetRetiringFederationAddress() ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederationAddress")
}

// UnpackGetRetiringFederationAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x47227286.
//
// Solidity: function getRetiringFederationAddress() view returns(string)
func (rskBridge *RskBridge) UnpackGetRetiringFederationAddress(data []byte) (string, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederationAddress", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackGetRetiringFederationCreationBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905153f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) PackGetRetiringFederationCreationBlockNumber() []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederationCreationBlockNumber")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederationCreationBlockNumber is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd905153f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) TryPackGetRetiringFederationCreationBlockNumber() ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederationCreationBlockNumber")
}

// UnpackGetRetiringFederationCreationBlockNumber is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd905153f.
//
// Solidity: function getRetiringFederationCreationBlockNumber() view returns(int256)
func (rskBridge *RskBridge) UnpackGetRetiringFederationCreationBlockNumber(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederationCreationBlockNumber", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRetiringFederationCreationTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f0ce9b1.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) PackGetRetiringFederationCreationTime() []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederationCreationTime")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederationCreationTime is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f0ce9b1.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) TryPackGetRetiringFederationCreationTime() ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederationCreationTime")
}

// UnpackGetRetiringFederationCreationTime is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x3f0ce9b1.
//
// Solidity: function getRetiringFederationCreationTime() view returns(int256)
func (rskBridge *RskBridge) UnpackGetRetiringFederationCreationTime(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederationCreationTime", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRetiringFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd970b0fd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (rskBridge *RskBridge) PackGetRetiringFederationSize() []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederationSize")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederationSize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd970b0fd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (rskBridge *RskBridge) TryPackGetRetiringFederationSize() ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederationSize")
}

// UnpackGetRetiringFederationSize is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd970b0fd.
//
// Solidity: function getRetiringFederationSize() view returns(int256)
func (rskBridge *RskBridge) UnpackGetRetiringFederationSize(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederationSize", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRetiringFederationThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07bbdfc4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) PackGetRetiringFederationThreshold() []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederationThreshold")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederationThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x07bbdfc4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) TryPackGetRetiringFederationThreshold() ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederationThreshold")
}

// UnpackGetRetiringFederationThreshold is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x07bbdfc4.
//
// Solidity: function getRetiringFederationThreshold() view returns(int256)
func (rskBridge *RskBridge) UnpackGetRetiringFederationThreshold(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederationThreshold", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackGetRetiringFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4675d6de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) PackGetRetiringFederatorPublicKey(index *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederatorPublicKey", index)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederatorPublicKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4675d6de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetRetiringFederatorPublicKey(index *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederatorPublicKey", index)
}

// UnpackGetRetiringFederatorPublicKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4675d6de.
//
// Solidity: function getRetiringFederatorPublicKey(int256 index) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetRetiringFederatorPublicKey(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederatorPublicKey", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetRetiringFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68bc2b2b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) PackGetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) []byte {
	enc, err := rskBridge.abi.Pack("getRetiringFederatorPublicKeyOfType", index, atype)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRetiringFederatorPublicKeyOfType is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x68bc2b2b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) TryPackGetRetiringFederatorPublicKeyOfType(index *big.Int, atype string) ([]byte, error) {
	return rskBridge.abi.Pack("getRetiringFederatorPublicKeyOfType", index, atype)
}

// UnpackGetRetiringFederatorPublicKeyOfType is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x68bc2b2b.
//
// Solidity: function getRetiringFederatorPublicKeyOfType(int256 index, string atype) view returns(bytes)
func (rskBridge *RskBridge) UnpackGetRetiringFederatorPublicKeyOfType(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getRetiringFederatorPublicKeyOfType", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetStateForBtcReleaseClient is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4fbca20.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (rskBridge *RskBridge) PackGetStateForBtcReleaseClient() []byte {
	enc, err := rskBridge.abi.Pack("getStateForBtcReleaseClient")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetStateForBtcReleaseClient is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc4fbca20.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (rskBridge *RskBridge) TryPackGetStateForBtcReleaseClient() ([]byte, error) {
	return rskBridge.abi.Pack("getStateForBtcReleaseClient")
}

// UnpackGetStateForBtcReleaseClient is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc4fbca20.
//
// Solidity: function getStateForBtcReleaseClient() view returns(bytes)
func (rskBridge *RskBridge) UnpackGetStateForBtcReleaseClient(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getStateForBtcReleaseClient", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackGetStateForDebugging is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d0cee93.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (rskBridge *RskBridge) PackGetStateForDebugging() []byte {
	enc, err := rskBridge.abi.Pack("getStateForDebugging")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetStateForDebugging is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0d0cee93.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (rskBridge *RskBridge) TryPackGetStateForDebugging() ([]byte, error) {
	return rskBridge.abi.Pack("getStateForDebugging")
}

// UnpackGetStateForDebugging is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0d0cee93.
//
// Solidity: function getStateForDebugging() view returns(bytes)
func (rskBridge *RskBridge) UnpackGetStateForDebugging(data []byte) ([]byte, error) {
	out, err := rskBridge.abi.Unpack("getStateForDebugging", data)
	if err != nil {
		return *new([]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	return out0, nil
}

// PackHasBtcBlockCoinbaseTransactionInformation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x253b944b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (rskBridge *RskBridge) PackHasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) []byte {
	enc, err := rskBridge.abi.Pack("hasBtcBlockCoinbaseTransactionInformation", blockHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasBtcBlockCoinbaseTransactionInformation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x253b944b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (rskBridge *RskBridge) TryPackHasBtcBlockCoinbaseTransactionInformation(blockHash [32]byte) ([]byte, error) {
	return rskBridge.abi.Pack("hasBtcBlockCoinbaseTransactionInformation", blockHash)
}

// UnpackHasBtcBlockCoinbaseTransactionInformation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x253b944b.
//
// Solidity: function hasBtcBlockCoinbaseTransactionInformation(bytes32 blockHash) view returns(bool)
func (rskBridge *RskBridge) UnpackHasBtcBlockCoinbaseTransactionInformation(data []byte) (bool, error) {
	out, err := rskBridge.abi.Unpack("hasBtcBlockCoinbaseTransactionInformation", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIncreaseLockingCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2910aeb2.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (rskBridge *RskBridge) PackIncreaseLockingCap(newLockingCap *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("increaseLockingCap", newLockingCap)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIncreaseLockingCap is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2910aeb2.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (rskBridge *RskBridge) TryPackIncreaseLockingCap(newLockingCap *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("increaseLockingCap", newLockingCap)
}

// UnpackIncreaseLockingCap is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2910aeb2.
//
// Solidity: function increaseLockingCap(int256 newLockingCap) returns(bool)
func (rskBridge *RskBridge) UnpackIncreaseLockingCap(data []byte) (bool, error) {
	out, err := rskBridge.abi.Unpack("increaseLockingCap", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsBtcTxHashAlreadyProcessed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a982d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (rskBridge *RskBridge) PackIsBtcTxHashAlreadyProcessed(hash string) []byte {
	enc, err := rskBridge.abi.Pack("isBtcTxHashAlreadyProcessed", hash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsBtcTxHashAlreadyProcessed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a982d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (rskBridge *RskBridge) TryPackIsBtcTxHashAlreadyProcessed(hash string) ([]byte, error) {
	return rskBridge.abi.Pack("isBtcTxHashAlreadyProcessed", hash)
}

// UnpackIsBtcTxHashAlreadyProcessed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a982d.
//
// Solidity: function isBtcTxHashAlreadyProcessed(string hash) view returns(bool)
func (rskBridge *RskBridge) UnpackIsBtcTxHashAlreadyProcessed(data []byte) (bool, error) {
	out, err := rskBridge.abi.Unpack("isBtcTxHashAlreadyProcessed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackReceiveHeader is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x884bdd86.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (rskBridge *RskBridge) PackReceiveHeader(ablock []byte) []byte {
	enc, err := rskBridge.abi.Pack("receiveHeader", ablock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReceiveHeader is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x884bdd86.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (rskBridge *RskBridge) TryPackReceiveHeader(ablock []byte) ([]byte, error) {
	return rskBridge.abi.Pack("receiveHeader", ablock)
}

// UnpackReceiveHeader is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x884bdd86.
//
// Solidity: function receiveHeader(bytes ablock) returns(int256)
func (rskBridge *RskBridge) UnpackReceiveHeader(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("receiveHeader", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackReceiveHeaders is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe5400e7b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (rskBridge *RskBridge) PackReceiveHeaders(blocks [][]byte) []byte {
	enc, err := rskBridge.abi.Pack("receiveHeaders", blocks)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackReceiveHeaders is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe5400e7b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function receiveHeaders(bytes[] blocks) returns()
func (rskBridge *RskBridge) TryPackReceiveHeaders(blocks [][]byte) ([]byte, error) {
	return rskBridge.abi.Pack("receiveHeaders", blocks)
}

// PackRegisterBtcCoinbaseTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccf417ae.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (rskBridge *RskBridge) PackRegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) []byte {
	enc, err := rskBridge.abi.Pack("registerBtcCoinbaseTransaction", btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterBtcCoinbaseTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccf417ae.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerBtcCoinbaseTransaction(bytes btcTxSerialized, bytes32 blockHash, bytes pmtSerialized, bytes32 witnessMerkleRoot, bytes32 witnessReservedValue) returns()
func (rskBridge *RskBridge) TryPackRegisterBtcCoinbaseTransaction(btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) ([]byte, error) {
	return rskBridge.abi.Pack("registerBtcCoinbaseTransaction", btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
}

// PackRegisterBtcTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43dc0656.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (rskBridge *RskBridge) PackRegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) []byte {
	enc, err := rskBridge.abi.Pack("registerBtcTransaction", atx, height, pmt)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterBtcTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x43dc0656.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerBtcTransaction(bytes atx, int256 height, bytes pmt) returns()
func (rskBridge *RskBridge) TryPackRegisterBtcTransaction(atx []byte, height *big.Int, pmt []byte) ([]byte, error) {
	return rskBridge.abi.Pack("registerBtcTransaction", atx, height, pmt)
}

// PackRegisterFastBridgeBtcTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6adc0133.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (rskBridge *RskBridge) PackRegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) []byte {
	enc, err := rskBridge.abi.Pack("registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterFastBridgeBtcTransaction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6adc0133.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (rskBridge *RskBridge) TryPackRegisterFastBridgeBtcTransaction(btcTxSerialized []byte, height *big.Int, pmtSerialized []byte, derivationArgumentsHash [32]byte, userRefundBtcAddress []byte, liquidityBridgeContractAddress common.Address, liquidityProviderBtcAddress []byte, shouldTransferToContract bool) ([]byte, error) {
	return rskBridge.abi.Pack("registerFastBridgeBtcTransaction", btcTxSerialized, height, pmtSerialized, derivationArgumentsHash, userRefundBtcAddress, liquidityBridgeContractAddress, liquidityProviderBtcAddress, shouldTransferToContract)
}

// UnpackRegisterFastBridgeBtcTransaction is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x6adc0133.
//
// Solidity: function registerFastBridgeBtcTransaction(bytes btcTxSerialized, uint256 height, bytes pmtSerialized, bytes32 derivationArgumentsHash, bytes userRefundBtcAddress, address liquidityBridgeContractAddress, bytes liquidityProviderBtcAddress, bool shouldTransferToContract) returns(int256)
func (rskBridge *RskBridge) UnpackRegisterFastBridgeBtcTransaction(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("registerFastBridgeBtcTransaction", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRemoveLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfcdeb46f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) PackRemoveLockWhitelistAddress(aaddress string) []byte {
	enc, err := rskBridge.abi.Pack("removeLockWhitelistAddress", aaddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRemoveLockWhitelistAddress is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfcdeb46f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) TryPackRemoveLockWhitelistAddress(aaddress string) ([]byte, error) {
	return rskBridge.abi.Pack("removeLockWhitelistAddress", aaddress)
}

// UnpackRemoveLockWhitelistAddress is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfcdeb46f.
//
// Solidity: function removeLockWhitelistAddress(string aaddress) returns(int256)
func (rskBridge *RskBridge) UnpackRemoveLockWhitelistAddress(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("removeLockWhitelistAddress", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackRollbackFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8dec3d32.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function rollbackFederation() returns(int256)
func (rskBridge *RskBridge) PackRollbackFederation() []byte {
	enc, err := rskBridge.abi.Pack("rollbackFederation")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRollbackFederation is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8dec3d32.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function rollbackFederation() returns(int256)
func (rskBridge *RskBridge) TryPackRollbackFederation() ([]byte, error) {
	return rskBridge.abi.Pack("rollbackFederation")
}

// UnpackRollbackFederation is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8dec3d32.
//
// Solidity: function rollbackFederation() returns(int256)
func (rskBridge *RskBridge) UnpackRollbackFederation(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("rollbackFederation", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackSetLockWhitelistDisableBlockDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1cc54f5.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (rskBridge *RskBridge) PackSetLockWhitelistDisableBlockDelay(disableDelay *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("setLockWhitelistDisableBlockDelay", disableDelay)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetLockWhitelistDisableBlockDelay is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc1cc54f5.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (rskBridge *RskBridge) TryPackSetLockWhitelistDisableBlockDelay(disableDelay *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("setLockWhitelistDisableBlockDelay", disableDelay)
}

// UnpackSetLockWhitelistDisableBlockDelay is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc1cc54f5.
//
// Solidity: function setLockWhitelistDisableBlockDelay(int256 disableDelay) returns(int256)
func (rskBridge *RskBridge) UnpackSetLockWhitelistDisableBlockDelay(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("setLockWhitelistDisableBlockDelay", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackUpdateCollections is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c5a9990.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function updateCollections() returns()
func (rskBridge *RskBridge) PackUpdateCollections() []byte {
	enc, err := rskBridge.abi.Pack("updateCollections")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackUpdateCollections is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0c5a9990.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function updateCollections() returns()
func (rskBridge *RskBridge) TryPackUpdateCollections() ([]byte, error) {
	return rskBridge.abi.Pack("updateCollections")
}

// PackVoteFeePerKbChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0461313e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (rskBridge *RskBridge) PackVoteFeePerKbChange(feePerKb *big.Int) []byte {
	enc, err := rskBridge.abi.Pack("voteFeePerKbChange", feePerKb)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVoteFeePerKbChange is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x0461313e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (rskBridge *RskBridge) TryPackVoteFeePerKbChange(feePerKb *big.Int) ([]byte, error) {
	return rskBridge.abi.Pack("voteFeePerKbChange", feePerKb)
}

// UnpackVoteFeePerKbChange is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x0461313e.
//
// Solidity: function voteFeePerKbChange(int256 feePerKb) returns(int256)
func (rskBridge *RskBridge) UnpackVoteFeePerKbChange(data []byte) (*big.Int, error) {
	out, err := rskBridge.abi.Unpack("voteFeePerKbChange", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// RskBridgeBatchPegoutCreated represents a batch_pegout_created event raised by the RskBridge contract.
type RskBridgeBatchPegoutCreated struct {
	BtcTx              [32]byte
	ReleaseRskTxHashes []byte
	Raw                *types.Log // Blockchain specific contextual infos
}

const RskBridgeBatchPegoutCreatedEventName = "batch_pegout_created"

// ContractEventName returns the user-defined event name.
func (RskBridgeBatchPegoutCreated) ContractEventName() string {
	return RskBridgeBatchPegoutCreatedEventName
}

// UnpackBatchPegoutCreatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event batch_pegout_created(bytes32 indexed btcTx, bytes releaseRskTxHashes)
func (rskBridge *RskBridge) UnpackBatchPegoutCreatedEvent(log *types.Log) (*RskBridgeBatchPegoutCreated, error) {
	event := "batch_pegout_created"
	if len(log.Topics) == 0 || log.Topics[0] != rskBridge.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(RskBridgeBatchPegoutCreated)
	if len(log.Data) > 0 {
		if err := rskBridge.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range rskBridge.abi.Events[event].Inputs {
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
