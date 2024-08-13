package routes_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

func TestGetPublicEndpoints(t *testing.T) {
	registryMock := &mocks.UseCaseRegistryMock{}
	registryMock.EXPECT().HealthUseCase().Return(&usecases.HealthUseCase{})
	registryMock.EXPECT().GetProvidersUseCase().Return(&liquidity_provider.GetProvidersUseCase{})
	registryMock.EXPECT().GetPeginQuoteUseCase().Return(&pegin.GetQuoteUseCase{})
	registryMock.EXPECT().GetAcceptPeginQuoteUseCase().Return(&pegin.AcceptQuoteUseCase{})
	registryMock.EXPECT().GetPegoutQuoteUseCase().Return(&pegout.GetQuoteUseCase{})
	registryMock.EXPECT().GetAcceptPegoutQuoteUseCase().Return(&pegout.AcceptQuoteUseCase{})
	registryMock.EXPECT().GetUserDepositsUseCase().Return(&pegout.GetUserDepositsUseCase{})
	registryMock.EXPECT().GetProviderDetailUseCase().Return(&liquidity_provider.GetDetailUseCase{})
	registryMock.EXPECT().GetPeginStatusUseCase().Return(&pegin.StatusUseCase{})
	registryMock.EXPECT().GetPegoutStatusUseCase().Return(&pegout.StatusUseCase{})
	registryMock.EXPECT().GetAvailableLiquidityUseCase().Return(&liquidity_provider.GetAvailableLiquidityUseCase{})

	endpoints := routes.GetPublicEndpoints(registryMock)
	specBytes := test.ReadFile(t, "OpenApi.yml")
	spec := &openApiSpecification{}

	err := yaml.Unmarshal(specBytes, spec)
	require.NoError(t, err)

	assert.Len(t, endpoints, 11)
	for _, endpoint := range endpoints {
		lowerCaseMethod := strings.ToLower(endpoint.Method)
		assert.NotNilf(t, spec.Paths[endpoint.Path][lowerCaseMethod], "Handler not found for path %s and verb %s", endpoint.Path, endpoint.Method)
	}
	registryMock.AssertExpectations(t)
}
