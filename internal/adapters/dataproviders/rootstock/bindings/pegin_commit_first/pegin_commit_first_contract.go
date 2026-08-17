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

// PeginCommitFirstContractMetaData contains all meta data concerning the PeginCommitFirstContract contract.
var PeginCommitFirstContractMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"requestPegIn\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxSerialized\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"opReturn\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"btcBlockHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleBranchPath\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"merkleBranchHashes\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"resolvePegIn\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"btcRawTransaction\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"partialMerkleTree\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"height\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"registerResult\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"PegInRequested\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"rskAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"netToUser\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"callSuccess\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PegInResolved\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"registrant\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"released\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"claimerPayout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"registrantFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"userPayout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressNotRegistered\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DepositOutputNotFound\",\"inputs\":[{\"name\":\"rskAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"btcTxHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"IncorrectFronting\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientConfirmations\",\"inputs\":[{\"name\":\"have\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"PegInAlreadyProcessed\",\"inputs\":[{\"name\":\"pegInId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	ID:  "PeginCommitFirstContract",
}

// PeginCommitFirstContract is an auto generated Go binding around an Ethereum contract.
type PeginCommitFirstContract struct {
	abi abi.ABI
}

// NewPeginCommitFirstContract creates a new instance of PeginCommitFirstContract.
func NewPeginCommitFirstContract() *PeginCommitFirstContract {
	parsed, err := PeginCommitFirstContractMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &PeginCommitFirstContract{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *PeginCommitFirstContract) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackRequestPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa355e935.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requestPegIn(address rskAddr, bytes btcTxSerialized, bytes opReturn, bytes32 btcBlockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) payable returns(bytes32 pegInId)
func (peginCommitFirstContract *PeginCommitFirstContract) PackRequestPegIn(rskAddr common.Address, btcTxSerialized []byte, opReturn []byte, btcBlockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) []byte {
	enc, err := peginCommitFirstContract.abi.Pack("requestPegIn", rskAddr, btcTxSerialized, opReturn, btcBlockHash, merkleBranchPath, merkleBranchHashes)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequestPegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa355e935.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requestPegIn(address rskAddr, bytes btcTxSerialized, bytes opReturn, bytes32 btcBlockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) payable returns(bytes32 pegInId)
func (peginCommitFirstContract *PeginCommitFirstContract) TryPackRequestPegIn(rskAddr common.Address, btcTxSerialized []byte, opReturn []byte, btcBlockHash [32]byte, merkleBranchPath *big.Int, merkleBranchHashes [][32]byte) ([]byte, error) {
	return peginCommitFirstContract.abi.Pack("requestPegIn", rskAddr, btcTxSerialized, opReturn, btcBlockHash, merkleBranchPath, merkleBranchHashes)
}

// UnpackRequestPegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa355e935.
//
// Solidity: function requestPegIn(address rskAddr, bytes btcTxSerialized, bytes opReturn, bytes32 btcBlockHash, uint256 merkleBranchPath, bytes32[] merkleBranchHashes) payable returns(bytes32 pegInId)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackRequestPegIn(data []byte) ([32]byte, error) {
	out, err := peginCommitFirstContract.abi.Unpack("requestPegIn", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackResolvePegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x958c8fb0.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256 registerResult)
func (peginCommitFirstContract *PeginCommitFirstContract) PackResolvePegIn(rskAddr common.Address, btcTxHash [32]byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) []byte {
	enc, err := peginCommitFirstContract.abi.Pack("resolvePegIn", rskAddr, btcTxHash, btcRawTransaction, partialMerkleTree, height)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackResolvePegIn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x958c8fb0.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256 registerResult)
func (peginCommitFirstContract *PeginCommitFirstContract) TryPackResolvePegIn(rskAddr common.Address, btcTxHash [32]byte, btcRawTransaction []byte, partialMerkleTree []byte, height *big.Int) ([]byte, error) {
	return peginCommitFirstContract.abi.Pack("resolvePegIn", rskAddr, btcTxHash, btcRawTransaction, partialMerkleTree, height)
}

// UnpackResolvePegIn is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x958c8fb0.
//
// Solidity: function resolvePegIn(address rskAddr, bytes32 btcTxHash, bytes btcRawTransaction, bytes partialMerkleTree, uint256 height) returns(int256 registerResult)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackResolvePegIn(data []byte) (*big.Int, error) {
	out, err := peginCommitFirstContract.abi.Unpack("resolvePegIn", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PeginCommitFirstContractPegInRequested represents a PegInRequested event raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractPegInRequested struct {
	PegInId     [32]byte
	Claimer     common.Address
	RskAddr     common.Address
	Amount      *big.Int
	NetToUser   *big.Int
	CallSuccess bool
	Raw         *types.Log // Blockchain specific contextual infos
}

const PeginCommitFirstContractPegInRequestedEventName = "PegInRequested"

// ContractEventName returns the user-defined event name.
func (PeginCommitFirstContractPegInRequested) ContractEventName() string {
	return PeginCommitFirstContractPegInRequestedEventName
}

// UnpackPegInRequestedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInRequested(bytes32 indexed pegInId, address indexed claimer, address indexed rskAddr, uint256 amount, uint256 netToUser, bool callSuccess)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackPegInRequestedEvent(log *types.Log) (*PeginCommitFirstContractPegInRequested, error) {
	event := "PegInRequested"
	if len(log.Topics) == 0 || log.Topics[0] != peginCommitFirstContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginCommitFirstContractPegInRequested)
	if len(log.Data) > 0 {
		if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginCommitFirstContract.abi.Events[event].Inputs {
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

// PeginCommitFirstContractPegInResolved represents a PegInResolved event raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractPegInResolved struct {
	PegInId       [32]byte
	Claimer       common.Address
	Registrant    common.Address
	Released      *big.Int
	ClaimerPayout *big.Int
	RegistrantFee *big.Int
	UserPayout    *big.Int
	Raw           *types.Log // Blockchain specific contextual infos
}

const PeginCommitFirstContractPegInResolvedEventName = "PegInResolved"

// ContractEventName returns the user-defined event name.
func (PeginCommitFirstContractPegInResolved) ContractEventName() string {
	return PeginCommitFirstContractPegInResolvedEventName
}

// UnpackPegInResolvedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PegInResolved(bytes32 indexed pegInId, address indexed claimer, address indexed registrant, uint256 released, uint256 claimerPayout, uint256 registrantFee, uint256 userPayout)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackPegInResolvedEvent(log *types.Log) (*PeginCommitFirstContractPegInResolved, error) {
	event := "PegInResolved"
	if len(log.Topics) == 0 || log.Topics[0] != peginCommitFirstContract.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(PeginCommitFirstContractPegInResolved)
	if len(log.Data) > 0 {
		if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range peginCommitFirstContract.abi.Events[event].Inputs {
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
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], peginCommitFirstContract.abi.Errors["AddressNotRegistered"].ID.Bytes()[:4]) {
		return peginCommitFirstContract.UnpackAddressNotRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginCommitFirstContract.abi.Errors["DepositOutputNotFound"].ID.Bytes()[:4]) {
		return peginCommitFirstContract.UnpackDepositOutputNotFoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginCommitFirstContract.abi.Errors["IncorrectFronting"].ID.Bytes()[:4]) {
		return peginCommitFirstContract.UnpackIncorrectFrontingError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginCommitFirstContract.abi.Errors["InsufficientConfirmations"].ID.Bytes()[:4]) {
		return peginCommitFirstContract.UnpackInsufficientConfirmationsError(raw[4:])
	}
	if bytes.Equal(raw[:4], peginCommitFirstContract.abi.Errors["PegInAlreadyProcessed"].ID.Bytes()[:4]) {
		return peginCommitFirstContract.UnpackPegInAlreadyProcessedError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// PeginCommitFirstContractAddressNotRegistered represents a AddressNotRegistered error raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractAddressNotRegistered struct {
	RskAddr common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressNotRegistered(address rskAddr)
func PeginCommitFirstContractAddressNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0xf3a8c4eea98f33ac9ed01788a298f73b7cc71b3532c7507ff436e3d37bd70a2a")
}

// UnpackAddressNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressNotRegistered(address rskAddr)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackAddressNotRegisteredError(raw []byte) (*PeginCommitFirstContractAddressNotRegistered, error) {
	out := new(PeginCommitFirstContractAddressNotRegistered)
	if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, "AddressNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginCommitFirstContractDepositOutputNotFound represents a DepositOutputNotFound error raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractDepositOutputNotFound struct {
	RskAddr   common.Address
	BtcTxHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error DepositOutputNotFound(address rskAddr, bytes32 btcTxHash)
func PeginCommitFirstContractDepositOutputNotFoundErrorID() common.Hash {
	return common.HexToHash("0xf32c452cbd044c3a0dae986f7dc3d620f8532983c066582f2ed2c4593c320ee7")
}

// UnpackDepositOutputNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error DepositOutputNotFound(address rskAddr, bytes32 btcTxHash)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackDepositOutputNotFoundError(raw []byte) (*PeginCommitFirstContractDepositOutputNotFound, error) {
	out := new(PeginCommitFirstContractDepositOutputNotFound)
	if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, "DepositOutputNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginCommitFirstContractIncorrectFronting represents a IncorrectFronting error raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractIncorrectFronting struct {
	Expected *big.Int
	Actual   *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error IncorrectFronting(uint256 expected, uint256 actual)
func PeginCommitFirstContractIncorrectFrontingErrorID() common.Hash {
	return common.HexToHash("0x4c0c9bd9730b66d0ffcd547ddbe6f891acd58a39d5bbd9bd3984d4e760e964fe")
}

// UnpackIncorrectFrontingError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error IncorrectFronting(uint256 expected, uint256 actual)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackIncorrectFrontingError(raw []byte) (*PeginCommitFirstContractIncorrectFronting, error) {
	out := new(PeginCommitFirstContractIncorrectFronting)
	if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, "IncorrectFronting", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginCommitFirstContractInsufficientConfirmations represents a InsufficientConfirmations error raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractInsufficientConfirmations struct {
	Have     *big.Int
	Required *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientConfirmations(uint256 have, uint256 required)
func PeginCommitFirstContractInsufficientConfirmationsErrorID() common.Hash {
	return common.HexToHash("0x22c53f1f9da074136409181cd7a7bfae7388312a843098ca3ab118f5de3d1890")
}

// UnpackInsufficientConfirmationsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientConfirmations(uint256 have, uint256 required)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackInsufficientConfirmationsError(raw []byte) (*PeginCommitFirstContractInsufficientConfirmations, error) {
	out := new(PeginCommitFirstContractInsufficientConfirmations)
	if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, "InsufficientConfirmations", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// PeginCommitFirstContractPegInAlreadyProcessed represents a PegInAlreadyProcessed error raised by the PeginCommitFirstContract contract.
type PeginCommitFirstContractPegInAlreadyProcessed struct {
	PegInId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error PegInAlreadyProcessed(bytes32 pegInId)
func PeginCommitFirstContractPegInAlreadyProcessedErrorID() common.Hash {
	return common.HexToHash("0x61ef28332480fb1a452c181f028122bc7068d7c84387e53bbca7f60f4f3cb6f9")
}

// UnpackPegInAlreadyProcessedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error PegInAlreadyProcessed(bytes32 pegInId)
func (peginCommitFirstContract *PeginCommitFirstContract) UnpackPegInAlreadyProcessedError(raw []byte) (*PeginCommitFirstContractPegInAlreadyProcessed, error) {
	out := new(PeginCommitFirstContractPegInAlreadyProcessed)
	if err := peginCommitFirstContract.abi.UnpackIntoInterface(out, "PegInAlreadyProcessed", raw); err != nil {
		return nil, err
	}
	return out, nil
}
