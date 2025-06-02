package routes_test

import (
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetManagementEndpoints(t *testing.T) {
	registryMock := &mocks.UseCaseRegistryMock{}
	registryMock.EXPECT().GetPeginCollateralUseCase().Return(&pegin.GetCollateralUseCase{})
	registryMock.EXPECT().AddPeginCollateralUseCase().Return(&pegin.AddCollateralUseCase{})
	registryMock.EXPECT().GetPegoutCollateralUseCase().Return(&pegout.GetCollateralUseCase{})
	registryMock.EXPECT().AddPegoutCollateralUseCase().Return(&pegout.AddCollateralUseCase{})
	registryMock.EXPECT().ChangeStatusUseCase().Return(&liquidity_provider.ChangeStatusUseCase{})
	registryMock.EXPECT().ResignationUseCase().Return(&liquidity_provider.ResignUseCase{})
	registryMock.EXPECT().WithdrawCollateralUseCase().Return(&liquidity_provider.WithdrawCollateralUseCase{})
	registryMock.EXPECT().SummariesUseCase().Return(&reports.SummariesUseCase{})
	registryMock.EXPECT().GetConfigurationUseCase().Return(&liquidity_provider.GetConfigUseCase{})
	registryMock.EXPECT().SetGeneralConfigUseCase().Return(&liquidity_provider.SetGeneralConfigUseCase{})
	registryMock.EXPECT().SetPeginConfigUseCase().Return(&liquidity_provider.SetPeginConfigUseCase{})
	registryMock.EXPECT().SetPegoutConfigUseCase().Return(&liquidity_provider.SetPegoutConfigUseCase{})
	registryMock.EXPECT().SetCredentialsUseCase().Return(&liquidity_provider.SetCredentialsUseCase{})
	registryMock.EXPECT().LoginUseCase().Return(&liquidity_provider.LoginUseCase{})
	registryMock.EXPECT().GetManagementUiDataUseCase().Return(&liquidity_provider.GetManagementUiDataUseCase{})
	registryMock.EXPECT().GetPeginReportUseCase().Return(&reports.GetPeginReportUseCase{})
	registryMock.EXPECT().GetPegoutReportUseCase().Return(&reports.GetPegoutReportUseCase{})
	registryMock.EXPECT().GetRevenueReportUseCase().Return(&reports.GetRevenueReportUseCase{})

	endpoints := routes.GetManagementEndpoints(environment.Environment{}, registryMock, &mocks.StoreMock{})
	specBytes := test.ReadFile(t, "OpenApi.yml")
	spec := &openApiSpecification{}

	err := yaml.Unmarshal(specBytes, spec)
	require.NoError(t, err)

	assert.Len(t, endpoints, 21)
	for _, endpoint := range endpoints {
		if endpoint.Path != routes.IconPath && endpoint.Path != routes.StaticPath {
			lowerCaseMethod := strings.ToLower(endpoint.Method)
			assert.NotNilf(t, spec.Paths[endpoint.Path][lowerCaseMethod], "Handler not found for path %s and verb %s", endpoint.Path, endpoint.Method)
		}
	}
	registryMock.AssertExpectations(t)
}
