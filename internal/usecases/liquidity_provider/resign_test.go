package liquidity_provider_test

import (
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResignUseCase_Run(t *testing.T) {
	lbc := &test.LbcMock{}
	provider := &test.ProviderMock{}
	provider.On("RskAddress").Return("0x01")
	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{
			Id:      1,
			Address: "0x01",
		},
	}, nil)
	lbc.On("ProviderResign").Return(nil).Once()
	useCase := liquidity_provider.NewResignUseCase(lbc, provider)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestResignUseCase_Run_NotRegistered(t *testing.T) {
	lbc := &test.LbcMock{}
	provider := &test.ProviderMock{}
	provider.On("RskAddress").Return("0x01")
	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{
			Id:      2,
			Address: "0x02",
		},
	}, nil)
	useCase := liquidity_provider.NewResignUseCase(lbc, provider)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	require.ErrorIs(t, err, usecases.ProviderConfigurationError)
}

func TestResignUseCase_Run_Error(t *testing.T) {
	lbc := &test.LbcMock{}
	provider := &test.ProviderMock{}
	provider.On("RskAddress").Return("0x01")
	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{
			Id:      1,
			Address: "0x01",
		},
	}, nil)
	lbc.On("ProviderResign").Return(assert.AnError).Once()
	useCase := liquidity_provider.NewResignUseCase(lbc, provider)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	require.Error(t, err)
}
