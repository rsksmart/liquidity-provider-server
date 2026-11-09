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

func mustPackConfigurationOutput(t *testing.T, method string, value any) []byte {
	t.Helper()
	parsedAbi, err := bindings.FlyoverConfigurationsContractMetaData.ParseABI()
	require.NoError(t, err)
	out, err := parsedAbi.Methods[method].Outputs.Pack(value)
	require.NoError(t, err)
	return out
}

func TestFlyoverConfigurationsContractImpl_GetPegInConfiguration(t *testing.T) {
	contractMock, configurationsBinding, configurations := newFlyoverConfigurationsTestContract()
	t.Run("maps every field of the peg-in configuration", func(t *testing.T) {
		raw := bindings.IFlyoverConfigurationsPegConfiguration{
			FixedFee:      big.NewInt(100000000000000),
			PercentageFee: big.NewInt(25),
			MinAmount:     big.NewInt(5000000000000000),
			MaxAmount:     new(big.Int).SetUint64(10000000000000000000),
			RegistrantFee: big.NewInt(100000000000000),
			ClaimGasFee:   big.NewInt(20000000000000),
			ConfirmationTiers: []bindings.IFlyoverConfigurationsConfirmationTier{
				{MaxAmount: big.NewInt(1000000000000000000), Confirmations: big.NewInt(2)},
				{MaxAmount: new(big.Int).SetUint64(10000000000000000000), Confirmations: big.NewInt(6)},
			},
		}
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackGetPegInConfiguration()),
			mock.Anything,
		).Return(mustPackConfigurationOutput(t, "getPegInConfiguration", raw), nil).Once()
		result, err := configurations.GetPegInConfiguration()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PegInConfiguration{
			FixedFee:      entities.NewWei(100000000000000),
			PercentageFee: 25,
			MinAmount:     entities.NewWei(5000000000000000),
			MaxAmount:     entities.NewUWei(10000000000000000000),
			RegistrantFee: entities.NewWei(100000000000000),
			ClaimGasFee:   entities.NewWei(20000000000000),
			ConfirmationTiers: []blockchain.ConfirmationTier{
				{MaxAmount: entities.NewWei(1000000000000000000), Confirmations: 2},
				{MaxAmount: entities.NewUWei(10000000000000000000), Confirmations: 6},
			},
		}, result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on call fail", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(configurationsBinding.PackGetPegInConfiguration()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := configurations.GetPegInConfiguration()
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, blockchain.PegInConfiguration{}, result)
		contractMock.caller.AssertExpectations(t)
	})
}
