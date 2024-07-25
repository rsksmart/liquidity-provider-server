package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ManagementTemplateId string

const (
	ManagementLoginTemplate ManagementTemplateId = "login.html"
	ManagementUiTemplate    ManagementTemplateId = "management.html"
	ManagementErrorTemplate ManagementTemplateId = "error.html"
)

type GetManagementUiDataUseCase struct {
	liquidityProviderRepository liquidity_provider.LiquidityProviderRepository
	lp                          liquidity_provider.LiquidityProvider
	peginLp                     liquidity_provider.PeginLiquidityProvider
	pegoutLp                    liquidity_provider.PegoutLiquidityProvider
	baseUrl                     string
}

func NewGetManagementUiDataUseCase(
	liquidityProviderRepository liquidity_provider.LiquidityProviderRepository,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	baseUrl string,
) *GetManagementUiDataUseCase {
	return &GetManagementUiDataUseCase{
		liquidityProviderRepository: liquidityProviderRepository,
		lp:                          lp,
		peginLp:                     peginLp,
		pegoutLp:                    pegoutLp,
		baseUrl:                     baseUrl,
	}
}

type ManagementTemplate struct {
	Name ManagementTemplateId
	Data ManagementTemplateData
}

type ManagementTemplateData struct {
	CredentialsSet bool
	BaseUrl        string
	Configuration  FullConfiguration
}

func (useCase *GetManagementUiDataUseCase) Run(ctx context.Context, loggedIn bool) (*ManagementTemplate, error) {
	if !loggedIn {
		return useCase.getLoginTemplateData(ctx)
	}
	return useCase.getManagementTemplateData(ctx), nil
}

func (useCase *GetManagementUiDataUseCase) getLoginTemplateData(ctx context.Context) (*ManagementTemplate, error) {
	credentials, err := useCase.liquidityProviderRepository.GetCredentials(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetManagementUiId, err)
	}

	return &ManagementTemplate{
		Name: ManagementLoginTemplate,
		Data: ManagementTemplateData{
			CredentialsSet: credentials != nil,
			BaseUrl:        useCase.baseUrl,
		},
	}, nil
}

func (useCase *GetManagementUiDataUseCase) getManagementTemplateData(ctx context.Context) *ManagementTemplate {
	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	peginConfiguration := useCase.peginLp.PeginConfiguration(ctx)
	pegoutConfiguration := useCase.pegoutLp.PegoutConfiguration(ctx)

	return &ManagementTemplate{
		Name: ManagementUiTemplate,
		Data: ManagementTemplateData{
			CredentialsSet: true,
			BaseUrl:        useCase.baseUrl,
			Configuration: FullConfiguration{
				General: generalConfiguration,
				Pegin:   peginConfiguration,
				Pegout:  pegoutConfiguration,
			},
		},
	}
}
