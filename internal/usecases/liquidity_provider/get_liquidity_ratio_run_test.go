package liquidity_provider_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultGeneralConfig() lpEntity.GeneralConfiguration {
	return lpEntity.GeneralConfiguration{
		MaxLiquidity: entities.NewWei(1000),
		ExcessTolerance: lpEntity.ExcessTolerance{
			IsFixed:         false,
			PercentageValue: utils.NewBigFloat64(20),
			FixedValue:      entities.NewWei(0),
		},
	}
}

func defaultStateConfig() lpEntity.StateConfiguration {
	return lpEntity.StateConfiguration{
		BtcLiquidityTargetPercentage: 50,
		RatioCooldownEndTimestamp:    0,
	}
}

func TestGetLiquidityRatioUseCase_Run_HappyPathCurrentRatio(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(400), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(600), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 0)

	require.NoError(t, err)
	assert.Equal(t, uint64(50), result.BtcPercentage)
	assert.Equal(t, uint64(50), result.RbtcPercentage)
	assert.Equal(t, "500", result.BtcTarget.AsBigInt().String())
	assert.Equal(t, "500", result.RbtcTarget.AsBigInt().String())
	assert.Equal(t, "600", result.BtcThreshold.AsBigInt().String())
	assert.Equal(t, "600", result.RbtcThreshold.AsBigInt().String())
	assert.Equal(t, "400", result.BtcCurrentBalance.AsBigInt().String())
	assert.Equal(t, "600", result.RbtcCurrentBalance.AsBigInt().String())
	assert.Equal(t, liquidity_provider.NetworkImpactDeficit, result.BtcImpact.Type)
	assert.Equal(t, "100", result.BtcImpact.Amount.AsBigInt().String())
	assert.Equal(t, liquidity_provider.NetworkImpactWithinTolerance, result.RbtcImpact.Type)
	assert.False(t, result.CooldownActive)
	assert.False(t, result.IsPreview)

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_PreviewMode(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(700), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(300), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 70)

	require.NoError(t, err)
	assert.Equal(t, uint64(70), result.BtcPercentage)
	assert.Equal(t, uint64(30), result.RbtcPercentage)
	assert.Equal(t, "700", result.BtcTarget.AsBigInt().String())
	assert.Equal(t, "300", result.RbtcTarget.AsBigInt().String())
	assert.True(t, result.IsPreview)

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_BtcExcessRbtcDeficit(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(700), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(200), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 0)

	require.NoError(t, err)
	assert.Equal(t, liquidity_provider.NetworkImpactExcess, result.BtcImpact.Type)
	assert.Equal(t, "200", result.BtcImpact.Amount.AsBigInt().String())
	assert.Equal(t, liquidity_provider.NetworkImpactDeficit, result.RbtcImpact.Type)
	assert.Equal(t, "300", result.RbtcImpact.Amount.AsBigInt().String())

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_BothWithinTolerance(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(550), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(550), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 0)

	require.NoError(t, err)
	assert.Equal(t, liquidity_provider.NetworkImpactWithinTolerance, result.BtcImpact.Type)
	assert.Equal(t, "0", result.BtcImpact.Amount.AsBigInt().String())
	assert.Equal(t, liquidity_provider.NetworkImpactWithinTolerance, result.RbtcImpact.Type)
	assert.Equal(t, "0", result.RbtcImpact.Amount.AsBigInt().String())

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_CooldownActive(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	stateConfig := defaultStateConfig()
	stateConfig.RatioCooldownEndTimestamp = time.Now().Unix() + 3600

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig, nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(500), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(500), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 0)

	require.NoError(t, err)
	assert.True(t, result.CooldownActive)
	assert.Equal(t, stateConfig.RatioCooldownEndTimestamp, result.CooldownEndTimestamp)
	assert.Equal(t, int64(liquidity_provider.CoolDownAfterRatioChange), result.CooldownDurationSeconds)

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_CooldownExpired(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	stateConfig := defaultStateConfig()
	stateConfig.RatioCooldownEndTimestamp = time.Now().Unix() - 1

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(stateConfig, nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(500), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(500), nil)

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	result, err := useCase.Run(ctx, 0)

	require.NoError(t, err)
	assert.False(t, result.CooldownActive)

	generalProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_StateConfigurationError(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(lpEntity.StateConfiguration{}, errors.New("db error"))

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	_, err := useCase.Run(ctx, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	generalProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_PegoutLiquidityError(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return((*entities.Wei)(nil), errors.New("btc rpc error"))

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	_, err := useCase.Run(ctx, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "btc rpc error")

	generalProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
}

func TestGetLiquidityRatioUseCase_Run_PeginLiquidityError(t *testing.T) {
	ctx := context.Background()
	generalProvider := new(mocks.ProviderMock)
	peginProvider := new(mocks.ProviderMock)
	pegoutProvider := new(mocks.ProviderMock)

	generalProvider.On("GeneralConfiguration", ctx).Return(defaultGeneralConfig())
	generalProvider.On("StateConfiguration", ctx).Return(defaultStateConfig(), nil)
	pegoutProvider.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(500), nil)
	peginProvider.On("AvailablePeginLiquidity", ctx).Return((*entities.Wei)(nil), errors.New("rsk rpc error"))

	useCase := liquidity_provider.NewGetLiquidityRatioUseCase(generalProvider, peginProvider, pegoutProvider)
	_, err := useCase.Run(ctx, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rsk rpc error")

	generalProvider.AssertExpectations(t)
	pegoutProvider.AssertExpectations(t)
	peginProvider.AssertExpectations(t)
}
