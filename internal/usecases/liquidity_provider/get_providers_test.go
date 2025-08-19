package liquidity_provider_test

import (
	"errors"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProvidersUseCase_Run(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}

	provider := lpEntity.RegisteredLiquidityProvider{
		Id:           1,
		Address:      "0x01",
		Name:         "one",
		ApiBaseUrl:   "api1.com",
		Status:       true,
		ProviderType: "both",
	}
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{provider}, nil).Once()

	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewGetProvidersUseCase(contracts)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{provider}, result)
}

func TestGetProvidersUseCase_Run_Fail(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}

	lbc.On("GetProviders").Return(
		[]lpEntity.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()

	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewGetProvidersUseCase(contracts)
	result, err := useCase.Run()

	lbc.AssertExpectations(t)
	require.Error(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{}, result)
}
