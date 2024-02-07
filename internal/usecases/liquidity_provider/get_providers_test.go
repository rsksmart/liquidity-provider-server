package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProvidersUseCase_Run(t *testing.T) {
	lbc := &test.LbcMock{}

	provider := entities.RegisteredLiquidityProvider{
		Id:           1,
		Address:      "0x01",
		Name:         "one",
		ApiBaseUrl:   "api1.com",
		Status:       true,
		ProviderType: "both",
	}
	lbc.On("GetProviders").Return([]entities.RegisteredLiquidityProvider{provider}, nil).Once()

	useCase := liquidity_provider.NewGetProvidersUseCase(lbc)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, []entities.RegisteredLiquidityProvider{provider}, result)
}

func TestGetProvidersUseCase_Run_Fail(t *testing.T) {
	lbc := &test.LbcMock{}

	lbc.On("GetProviders").Return(
		[]entities.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()

	useCase := liquidity_provider.NewGetProvidersUseCase(lbc)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, []entities.RegisteredLiquidityProvider{}, result)
}
