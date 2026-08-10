package rootstock_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover_configurations"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFlyoverConfigurationsContractImpl(t *testing.T) {
	boundContract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewFlyoverConfigurationsContract()
	contract := rootstock.NewFlyoverConfigurationsContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		boundContract,
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		contractBinding,
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestFlyoverConfigurationsContractImpl_GetAddress(t *testing.T) {
	configurations := rootstock.NewFlyoverConfigurationsContractImpl(dummyClient, test.AnyAddress, nil, rootstock.RetryParams{}, nil, Abis)
	assert.Equal(t, test.AnyAddress, configurations.GetAddress())
}

// newFlyoverConfigurationsTestContract builds a configurations adapter wired to a fresh bound-contract mock.
func newFlyoverConfigurationsTestContract() (boundContractMock, *bindings.FlyoverConfigurationsContract, blockchain.FlyoverConfigurationsContract) {
	contractMock := createBoundContractMock()
	configurationsBinding := bindings.NewFlyoverConfigurationsContract()
	configurations := rootstock.NewFlyoverConfigurationsContractImpl(dummyClient, test.AnyAddress, contractMock.contract, rootstock.RetryParams{}, configurationsBinding, Abis)
	return contractMock, configurationsBinding, configurations
}

func TestFlyoverConfigurationsContractImpl_CalculatePegInFee(t *testing.T) {
	contractMock, configurationsBinding, configurations := newFlyoverConfigurationsTestContract()
	amount := entities.NewWei(1000)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackCalculatePegInFee(amount.AsBigInt())),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(25)), nil).Once()
		result, err := configurations.CalculatePegInFee(amount)
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(25), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackCalculatePegInFee(amount.AsBigInt())),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := configurations.CalculatePegInFee(amount)
		require.Error(t, err)
		assert.Nil(t, result)
		contractMock.caller.AssertExpectations(t)
	})
}

func TestFlyoverConfigurationsContractImpl_GetRequiredPegInBtcConfirmations(t *testing.T) {
	contractMock, configurationsBinding, configurations := newFlyoverConfigurationsTestContract()
	amount := entities.NewWei(2000)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackGetRequiredPegInBtcConfirmations(amount.AsBigInt())),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(10)), nil).Once()
		result, err := configurations.GetRequiredPegInBtcConfirmations(amount)
		require.NoError(t, err)
		assert.Equal(t, uint64(10), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackGetRequiredPegInBtcConfirmations(amount.AsBigInt())),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := configurations.GetRequiredPegInBtcConfirmations(amount)
		require.Error(t, err)
		assert.Equal(t, uint64(0), result)
		contractMock.caller.AssertExpectations(t)
	})
}
