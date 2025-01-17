package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"slices"
	"strings"
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
	contracts                   blockchain.RskContracts
	baseUrl                     string
}

func NewGetManagementUiDataUseCase(
	liquidityProviderRepository liquidity_provider.LiquidityProviderRepository,
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	contracts blockchain.RskContracts,
	baseUrl string,
) *GetManagementUiDataUseCase {
	return &GetManagementUiDataUseCase{
		liquidityProviderRepository: liquidityProviderRepository,
		lp:                          lp,
		peginLp:                     peginLp,
		pegoutLp:                    pegoutLp,
		contracts:                   contracts,
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
	BtcAddress     string
	RskAddress     string
	ProviderData   liquidity_provider.RegisteredLiquidityProvider
	Configuration  FullConfiguration
}

func (useCase *GetManagementUiDataUseCase) Run(ctx context.Context, loggedIn bool) (*ManagementTemplate, error) {
	if !loggedIn {
		return useCase.getLoginTemplateData(ctx)
	}
	return useCase.getManagementTemplateData(ctx)
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

func (useCase *GetManagementUiDataUseCase) getManagementTemplateData(ctx context.Context) (*ManagementTemplate, error) {
	generalConfiguration := useCase.lp.GeneralConfiguration(ctx)
	peginConfiguration := useCase.peginLp.PeginConfiguration(ctx)
	pegoutConfiguration := useCase.pegoutLp.PegoutConfiguration(ctx)

	rskAddress := useCase.lp.RskAddress()
	// TODO change to getProvider in 2.1.0
	providers, err := useCase.contracts.Lbc.GetProviders()
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetManagementUiId, err)
	}

	index := slices.IndexFunc(providers, func(providerInfo liquidity_provider.RegisteredLiquidityProvider) bool {
		return strings.EqualFold(providerInfo.Address, rskAddress)
	})
	if index == -1 {
		return nil, usecases.WrapUseCaseError(usecases.GetManagementUiId, usecases.ProviderNotFoundError)
	}

	return &ManagementTemplate{
		Name: ManagementUiTemplate,
		Data: ManagementTemplateData{
			CredentialsSet: true,
			BaseUrl:        useCase.baseUrl,
			BtcAddress:     useCase.lp.BtcAddress(),
			RskAddress:     rskAddress,
			ProviderData:   providers[index],
			Configuration: FullConfiguration{
				General: generalConfiguration,
				Pegin:   peginConfiguration,
				Pegout:  pegoutConfiguration,
			},
		},
	}, nil
}
