package liquidity_provider_test

import (
	"errors"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetProvidersUseCase_Run(t *testing.T) {
	lbc := &test.LbcMock{}

	provider := lpEntity.RegisteredLiquidityProvider{
		Id:           1,
		Address:      "0x01",
		Name:         "one",
		ApiBaseUrl:   "api1.com",
		Status:       true,
		ProviderType: "both",
	}
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{provider}, nil).Once()

	useCase := liquidity_provider.NewGetProvidersUseCase(lbc)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{provider}, result)
}

func TestGetProvidersUseCase_Run_Fail(t *testing.T) {
	lbc := &test.LbcMock{}

	lbc.On("GetProviders").Return(
		[]lpEntity.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()

	useCase := liquidity_provider.NewGetProvidersUseCase(lbc)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	require.Error(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{}, result)
}
