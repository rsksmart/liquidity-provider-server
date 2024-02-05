package liquidity_provider

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetDetailUseCase struct {
	captchaSiteKey string
	peginProvider  entities.PeginLiquidityProvider
	pegoutProvider entities.PegoutLiquidityProvider
}

func NewGetDetailUseCase(
	captchaSiteKey string,
	peginProvider entities.PeginLiquidityProvider,
	pegoutProvider entities.PegoutLiquidityProvider,
) *GetDetailUseCase {
	return &GetDetailUseCase{
		captchaSiteKey: captchaSiteKey,
		peginProvider:  peginProvider,
		pegoutProvider: pegoutProvider,
	}
}

type FullLiquidityProvider struct {
	SiteKey string                           `json:"siteKey"`
	Pegin   entities.LiquidityProviderDetail `json:"pegin"`
	Pegout  entities.LiquidityProviderDetail `json:"pegout"`
}

func (useCase *GetDetailUseCase) Run() (FullLiquidityProvider, error) {
	var err error

	detail := FullLiquidityProvider{
		SiteKey: useCase.captchaSiteKey,
		Pegin: entities.LiquidityProviderDetail{
			Fee:                   useCase.peginProvider.CallFeePegin(),
			MinTransactionValue:   useCase.peginProvider.MinPegin(),
			MaxTransactionValue:   useCase.peginProvider.MaxPegin(),
			RequiredConfirmations: useCase.peginProvider.MaxPeginConfirmations(),
		},
		Pegout: entities.LiquidityProviderDetail{
			Fee:                   useCase.pegoutProvider.CallFeePegout(),
			MinTransactionValue:   useCase.pegoutProvider.MinPegout(),
			MaxTransactionValue:   useCase.pegoutProvider.MaxPegout(),
			RequiredConfirmations: useCase.pegoutProvider.MaxPegoutConfirmations(),
		},
	}

	if detail.SiteKey == "" {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, errors.New("missing captcha key"))
	} else if err = entities.ValidateStruct(detail.Pegin); err != nil {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, err)
	} else if err = entities.ValidateStruct(detail.Pegout); err != nil {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, err)
	}

	return detail, nil
}
