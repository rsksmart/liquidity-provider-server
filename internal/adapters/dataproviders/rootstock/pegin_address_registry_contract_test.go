package rootstock_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mustPackWithEncoding ABI-encodes a payload of the given Solidity type followed by the uint8
// encoding tag, which is the return shape of both registry address reads.
func mustPackWithEncoding(t *testing.T, solidityType string, payload any, encoding uint8) []byte {
	t.Helper()
	payloadType, err := abi.NewType(solidityType, "", nil)
	require.NoError(t, err)
	uint8Type, err := abi.NewType("uint8", "", nil)
	require.NoError(t, err)
	args := abi.Arguments{{Type: payloadType}, {Type: uint8Type}}
	out, err := args.Pack(payload, encoding)
	require.NoError(t, err)
	return out
}

// newPegInAddressRegistryTestContract builds a registry adapter wired to a fresh bound-contract mock.
func newPegInAddressRegistryTestContract() (boundContractMock, *bindings.PegInAddressRegistryContract, blockchain.PegInAddressRegistryContract) {
	contractMock := createBoundContractMock()
	registryBinding := bindings.NewPegInAddressRegistryContract()
	registry := rootstock.NewPegInAddressRegistryContractImpl(dummyClient, test.AnyAddress, contractMock.contract, rootstock.RetryParams{}, registryBinding, Abis)
	return contractMock, registryBinding, registry
}

func mustPackRegistration(t *testing.T, registrant common.Address, registrationBlock *big.Int) []byte {
	t.Helper()
	regType, err := abi.NewType("tuple", "structIPegInAddressRegistry.Registration", []abi.ArgumentMarshaling{
		{Name: "registrant", Type: "address"},
		{Name: "registrationBlock", Type: "uint96"},
	})
	require.NoError(t, err)
	args := abi.Arguments{{Type: regType}}
	out, err := args.Pack(bindings.IPegInAddressRegistryRegistration{
		Registrant:        registrant,
		RegistrationBlock: registrationBlock,
	})
	require.NoError(t, err)
	return out
}

func TestNewPegInAddressRegistryContractImpl(t *testing.T) {
	boundContract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewPegInAddressRegistryContract()
	contract := rootstock.NewPegInAddressRegistryContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		boundContract,
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		contractBinding,
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPegInAddressRegistryContractImpl_GetAddress(t *testing.T) {
	registry := rootstock.NewPegInAddressRegistryContractImpl(dummyClient, test.AnyAddress, nil, rootstock.RetryParams{}, nil, Abis)
	assert.Equal(t, test.AnyAddress, registry.GetAddress())
}

func TestPegInAddressRegistryContractImpl_GetPegInAddress(t *testing.T) {
	contractMock, registryBinding, registry := newPegInAddressRegistryTestContract()
	payload := []byte{0x01, 0x02, 0x03}
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetPegInAddress(parsedAddress)),
			mock.Anything,
		).Return(mustPackWithEncoding(t, "bytes", payload, 1), nil).Once()
		result, err := registry.GetPegInAddress(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, blockchain.PegInAddress{
			Payload:  payload,
			Encoding: blockchain.PegInAddressRegistryEncodingBech32,
		}, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetPegInAddress(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.GetPegInAddress(parsedAddress.String())
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Invalid address", func(t *testing.T) {
		result, err := registry.GetPegInAddress(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Empty(t, result)
	})
}

func TestPegInAddressRegistryContractImpl_GetPegInAddresses(t *testing.T) {
	contractMock, registryBinding, registry := newPegInAddressRegistryTestContract()
	otherAddress := common.HexToAddress("0x00000000000000000000000000000000000abc")
	payloads := [][]byte{{0x01, 0x02}, {0x03, 0x04}}
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetPegInAddresses([]common.Address{parsedAddress, otherAddress})),
			mock.Anything,
		).Return(mustPackWithEncoding(t, "bytes[]", payloads, 0), nil).Once()
		result, err := registry.GetPegInAddresses([]string{parsedAddress.String(), otherAddress.String()})
		require.NoError(t, err)
		assert.Equal(t, blockchain.PegInAddressBatch{
			Payloads: payloads,
			Encoding: blockchain.PegInAddressRegistryEncodingBase58,
		}, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetPegInAddresses([]common.Address{parsedAddress, otherAddress})),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.GetPegInAddresses([]string{parsedAddress.String(), otherAddress.String()})
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Invalid address", func(t *testing.T) {
		result, err := registry.GetPegInAddresses([]string{test.AnyString})
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Empty(t, result)
	})
}

func TestPegInAddressRegistryContractImpl_IsRegistered(t *testing.T) {
	contractMock, registryBinding, registry := newPegInAddressRegistryTestContract()
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackIsRegistered(parsedAddress)),
			mock.Anything,
		).Return(mustPackBool(t, true), nil).Once()
		result, err := registry.IsRegistered(parsedAddress.String())
		require.NoError(t, err)
		assert.True(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackIsRegistered(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.IsRegistered(parsedAddress.String())
		require.Error(t, err)
		assert.False(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Invalid address", func(t *testing.T) {
		result, err := registry.IsRegistered(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.False(t, result)
	})
}

func TestPegInAddressRegistryContractImpl_GetRegistration(t *testing.T) {
	contractMock, registryBinding, registry := newPegInAddressRegistryTestContract()
	registrant := common.HexToAddress("0x00000000000000000000000000000000000def")
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetRegistration(parsedAddress)),
			mock.Anything,
		).Return(mustPackRegistration(t, registrant, big.NewInt(555)), nil).Once()
		result, err := registry.GetRegistration(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, blockchain.PegInRegistration{
			Registrant:        registrant.String(),
			RegistrationBlock: 555,
		}, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetRegistration(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.GetRegistration(parsedAddress.String())
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Invalid address", func(t *testing.T) {
		result, err := registry.GetRegistration(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Empty(t, result)
	})
}

func TestPegInAddressRegistryContractImpl_GetRegistrationRoot(t *testing.T) {
	contractMock, registryBinding, registry := newPegInAddressRegistryTestContract()
	var root [32]byte
	root[31] = 0x2b
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetRegistrationRoot()),
			mock.Anything,
		).Return(mustPackBytes32(t, root), nil).Once()
		result, err := registry.GetRegistrationRoot()
		require.NoError(t, err)
		assert.Equal(t, root, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackGetRegistrationRoot()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.GetRegistrationRoot()
		require.Error(t, err)
		assert.Equal(t, [32]byte{}, result)
		contractMock.caller.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPegInAddressRegistryContractImpl_GetAddressRegisteredEvents(t *testing.T) {
	contractMock, _, registry := newPegInAddressRegistryTestContract()

	registryAbi, err := bindings.PegInAddressRegistryContractMetaData.ParseABI()
	require.NoError(t, err)
	eventID := registryAbi.Events["AddressRegistered"].ID

	rskAddr := common.HexToAddress("0x0100000000000000000000000000000000000000")
	registrant := common.HexToAddress("0x0200000000000000000000000000000000000000")
	var root [32]byte
	root[31] = 0x2b

	addressRegisteredLogs := []geth.Log{
		{
			TxHash:      common.Hash{7},
			BlockNumber: 10,
			Index:       3,
			Topics: []common.Hash{
				eventID,
				common.BytesToHash(rskAddr.Bytes()),
				common.BytesToHash(registrant.Bytes()),
			},
			Data: root[:],
		},
		{
			TxHash:      common.Hash{8},
			BlockNumber: 10,
			Index:       9,
			Topics: []common.Hash{
				eventID,
				common.BytesToHash(rskAddr.Bytes()),
				common.BytesToHash(registrant.Bytes()),
			},
			Data: root[:],
		},
	}

	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(addressRegisteredLogs, nil).Once()
		result, err := registry.GetAddressRegisteredEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, []blockchain.AddressRegistered{
			{
				RskAddress:       rskAddr.String(),
				Registrant:       registrant.String(),
				RegistrationRoot: root,
				TxHash:           common.Hash{7}.String(),
				BlockNumber:      10,
				LogIndex:         3,
			},
			{
				RskAddress:       rskAddr.String(),
				Registrant:       registrant.String(),
				RegistrationRoot: root,
				TxHash:           common.Hash{8}.String(),
				BlockNumber:      10,
				LogIndex:         9,
			},
		}, result)
		contractMock.filterer.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(nil, assert.AnError).Once()
		result, err := registry.GetAddressRegisteredEvents(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractMock.filterer.AssertExpectations(t)
	})
	t.Run("No events found", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return([]geth.Log{}, nil).Once()
		result, err := registry.GetAddressRegisteredEvents(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Empty(t, result)
		contractMock.filterer.AssertExpectations(t)
	})
}
