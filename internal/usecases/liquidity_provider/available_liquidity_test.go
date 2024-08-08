package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAvailableLiquidityUseCase_Run(t *testing.T) {
	t.Run("Return error when feature disabled", func(t *testing.T) {
		lp := new(mocks.ProviderMock)
		lp.On("GeneralConfiguration", mock.Anything).Return(lpEntity.GeneralConfiguration{PublicLiquidityCheck: false})
		useCase := liquidity_provider.NewGetAvailableLiquidityUseCase(lp, lp, lp)
		result, err := useCase.Run(context.Background())
		require.ErrorIs(t, err, liquidity_provider.LiquidityCheckNotEnabledError)
		assert.Empty(t, result)
	})
	t.Run("Return pegin & pegout liquidity when feature enabled", func(t *testing.T) {
		lp := new(mocks.ProviderMock)
		lp.On("GeneralConfiguration", mock.Anything).Return(lpEntity.GeneralConfiguration{PublicLiquidityCheck: true})
		lp.On("AvailablePeginLiquidity", mock.Anything).Return(entities.NewWei(100), nil)
		lp.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.NewWei(200), nil)
		useCase := liquidity_provider.NewGetAvailableLiquidityUseCase(lp, lp, lp)
		result, err := useCase.Run(context.Background())
		require.NoError(t, err)
		assert.Equal(t, lpEntity.AvailableLiquidity{
			PeginLiquidity:  entities.NewWei(100),
			PegoutLiquidity: entities.NewWei(200),
		}, result)
	})
	t.Run("Handle error when getting pegin liquidity", func(t *testing.T) {
		lp := new(mocks.ProviderMock)
		lp.On("GeneralConfiguration", mock.Anything).Return(lpEntity.GeneralConfiguration{PublicLiquidityCheck: true})
		lp.On("AvailablePeginLiquidity", mock.Anything).Return((*entities.Wei)(nil), assert.AnError)
		useCase := liquidity_provider.NewGetAvailableLiquidityUseCase(lp, lp, lp)
		result, err := useCase.Run(context.Background())
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Handle error when getting pegout liquidity", func(t *testing.T) {
		lp := new(mocks.ProviderMock)
		lp.On("GeneralConfiguration", mock.Anything).Return(lpEntity.GeneralConfiguration{PublicLiquidityCheck: true})
		lp.On("AvailablePeginLiquidity", mock.Anything).Return(entities.NewWei(100), nil)
		lp.On("AvailablePegoutLiquidity", mock.Anything).Return((*entities.Wei)(nil), assert.AnError)
		useCase := liquidity_provider.NewGetAvailableLiquidityUseCase(lp, lp, lp)
		result, err := useCase.Run(context.Background())
		require.Error(t, err)
		assert.Empty(t, result)
	})
}
