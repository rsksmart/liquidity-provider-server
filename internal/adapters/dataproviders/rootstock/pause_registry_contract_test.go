package rootstock_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newPauseRegistryTestContract() (boundContractMock, *bindings.PauseRegistryContract, blockchain.PauseRegistryContract) {
	contractMock := createBoundContractMock()
	registryBinding := bindings.NewPauseRegistryContract()
	registry := rootstock.NewPauseRegistryContractImpl(dummyClient, test.AnyAddress, contractMock.contract, rootstock.RetryParams{}, registryBinding, Abis)
	return contractMock, registryBinding, registry
}

func TestNewPauseRegistryContractImpl(t *testing.T) {
	boundContract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewPauseRegistryContract()
	contract := rootstock.NewPauseRegistryContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		boundContract,
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		contractBinding,
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPauseRegistryContractImpl_GetAddress(t *testing.T) {
	registry := rootstock.NewPauseRegistryContractImpl(dummyClient, test.AnyAddress, nil, rootstock.RetryParams{}, nil, Abis)
	assert.Equal(t, test.AnyAddress, registry.GetAddress())
}

func TestPauseRegistryContractImpl_PauseLevel(t *testing.T) {
	contractMock, registryBinding, registry := newPauseRegistryTestContract()
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseLevel()),
			mock.Anything,
		).Return(mustPackUint8(t, blockchain.PauseLevelHard), nil).Once()
		result, err := registry.PauseLevel()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseLevelHard, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseLevel()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.PauseLevel()
		require.Error(t, err)
		assert.Equal(t, uint8(0), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on decode fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseLevel()),
			mock.Anything,
		).Return([]byte{0x01}, nil).Once()
		result, err := registry.PauseLevel()
		require.Error(t, err)
		assert.Equal(t, uint8(0), result)
		contractMock.caller.AssertExpectations(t)
	})
}

func TestPauseRegistryContractImpl_PauseStatus(t *testing.T) {
	contractMock, registryBinding, registry := newPauseRegistryTestContract()
	status := generalPauseStatus{IsPaused: true, Reason: "maintenance", Since: 42}
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseStatus()),
			mock.Anything,
		).Return(mustPackPauseStatus(t, status), nil).Once()
		result, err := registry.PauseStatus()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseStatus{
			IsPaused: true,
			Reason:   "maintenance",
			Since:    42,
		}, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseStatus()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := registry.PauseStatus()
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on decode fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(registryBinding.PackPauseStatus()),
			mock.Anything,
		).Return([]byte{0x01}, nil).Once()
		result, err := registry.PauseStatus()
		require.Error(t, err)
		assert.Empty(t, result)
		contractMock.caller.AssertExpectations(t)
	})
}
