package liquidity_provider

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetDetailUseCase struct {
	captchaSiteKey   string
	captchaDisabled  bool
	segwitFederation bool
	provider         liquidity_provider.LiquidityProvider
	peginProvider    liquidity_provider.PeginLiquidityProvider
	pegoutProvider   liquidity_provider.PegoutLiquidityProvider
}

func NewGetDetailUseCase(
	captchaSiteKey string,
	captchaDisabled bool,
	segwitFederation bool,
	provider liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
) *GetDetailUseCase {
	return &GetDetailUseCase{
		captchaSiteKey:   captchaSiteKey,
		captchaDisabled:  captchaDisabled,
		provider:         provider,
		peginProvider:    peginProvider,
		pegoutProvider:   pegoutProvider,
		segwitFederation: segwitFederation,
	}
}

type FullLiquidityProvider struct {
	SiteKey               string                                     `json:"siteKey"`
	LiquidityCheckEnabled bool                                       `json:"liquidityCheckEnabled"`
	UsingSegwitFederation bool                                       `json:"usingSegwitFederation"`
	Pegin                 liquidity_provider.LiquidityProviderDetail `json:"pegin"`
	Pegout                liquidity_provider.LiquidityProviderDetail `json:"pegout"`
}

func (useCase *GetDetailUseCase) Run(ctx context.Context) (FullLiquidityProvider, error) {
	var err error
	generalConfiguration := useCase.provider.GeneralConfiguration(ctx)
	peginConfig := useCase.peginProvider.PeginConfiguration(ctx)
	pegoutConfig := useCase.pegoutProvider.PegoutConfiguration(ctx)
	detail := FullLiquidityProvider{
		SiteKey:               useCase.captchaSiteKey,
		LiquidityCheckEnabled: generalConfiguration.PublicLiquidityCheck,
		UsingSegwitFederation: useCase.segwitFederation,
		Pegin: liquidity_provider.LiquidityProviderDetail{
			FixedFee:              peginConfig.FixedFee,
			FeePercentage:         peginConfig.FeePercentage,
			MinTransactionValue:   peginConfig.MinValue,
			MaxTransactionValue:   peginConfig.MaxValue,
			RequiredConfirmations: generalConfiguration.BtcConfirmations.Max(),
		},
		Pegout: liquidity_provider.LiquidityProviderDetail{
			FixedFee:              pegoutConfig.FixedFee,
			FeePercentage:         pegoutConfig.FeePercentage,
			MinTransactionValue:   pegoutConfig.MinValue,
			MaxTransactionValue:   pegoutConfig.MaxValue,
			RequiredConfirmations: generalConfiguration.RskConfirmations.Max(),
		},
	}

	if detail.SiteKey == "" && !useCase.captchaDisabled {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, errors.New("missing captcha key"))
	} else if err = entities.ValidateStruct(detail.Pegin); err != nil {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, err)
	} else if err = entities.ValidateStruct(detail.Pegout); err != nil {
		return FullLiquidityProvider{}, usecases.WrapUseCaseError(usecases.ProviderDetailId, err)
	}

	return detail, nil
}
