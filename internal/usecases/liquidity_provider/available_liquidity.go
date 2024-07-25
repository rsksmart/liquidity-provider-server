package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetAvailableLiquidityUseCase struct {
	peginProvider       liquidity_provider.PeginLiquidityProvider
	pegoutProvider      liquidity_provider.PegoutLiquidityProvider
	generalProviderInfo liquidity_provider.LiquidityProvider
}

func NewGetAvailableLiquidityUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	generalProviderInfo liquidity_provider.LiquidityProvider,
) *GetAvailableLiquidityUseCase {
	return &GetAvailableLiquidityUseCase{
		peginProvider:       peginProvider,
		pegoutProvider:      pegoutProvider,
		generalProviderInfo: generalProviderInfo,
	}
}

func (useCase *GetAvailableLiquidityUseCase) Run(ctx context.Context) (liquidity_provider.AvailableLiquidity, error) {
	generalConfig := useCase.generalProviderInfo.GeneralConfiguration(ctx)
	if !generalConfig.PublicLiquidityCheck {
		return liquidity_provider.AvailableLiquidity{}, usecases.WrapUseCaseError(usecases.GetAvailableLiquidityId, LiquidityCheckNotEnabledError)
	}
	peginLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return liquidity_provider.AvailableLiquidity{}, usecases.WrapUseCaseError(usecases.GetAvailableLiquidityId, err)
	}
	pegoutLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return liquidity_provider.AvailableLiquidity{}, usecases.WrapUseCaseError(usecases.GetAvailableLiquidityId, err)
	}
	return liquidity_provider.AvailableLiquidity{
		PeginLiquidity:  peginLiquidity,
		PegoutLiquidity: pegoutLiquidity,
	}, nil
}
