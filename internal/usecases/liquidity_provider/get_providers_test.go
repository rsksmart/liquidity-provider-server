package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetProvidersUseCase_Run(t *testing.T) {
	discovery := &mocks.DiscoveryContractMock{}

	provider := lpEntity.RegisteredLiquidityProvider{
		Id:           1,
		Address:      "0x01",
		Name:         "one",
		ApiBaseUrl:   "api1.com",
		Status:       true,
		ProviderType: lpEntity.FullProvider,
	}
	discovery.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{provider}, nil).Once()

	contracts := blockchain.RskContracts{Discovery: discovery}
	useCase := liquidity_provider.NewGetProvidersUseCase(contracts)
	result, err := useCase.Run()

	discovery.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{provider}, result)
}

func TestGetProvidersUseCase_Run_Fail(t *testing.T) {
	discovery := &mocks.DiscoveryContractMock{}

	discovery.On("GetProviders").Return(
		[]lpEntity.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()

	contracts := blockchain.RskContracts{Discovery: discovery}
	useCase := liquidity_provider.NewGetProvidersUseCase(contracts)
	result, err := useCase.Run()

	discovery.AssertExpectations(t)
	require.Error(t, err)
	assert.Equal(t, []lpEntity.RegisteredLiquidityProvider{}, result)
}
