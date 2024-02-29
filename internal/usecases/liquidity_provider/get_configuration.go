package liquidity_provider

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type GetConfigUseCase struct {
	lp       liquidity_provider.LiquidityProvider
	peginLp  liquidity_provider.PeginLiquidityProvider
	pegoutLp liquidity_provider.PegoutLiquidityProvider
}

func NewGetConfigUseCase(
	lp liquidity_provider.LiquidityProvider,
	peginLp liquidity_provider.PeginLiquidityProvider,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
) *GetConfigUseCase {
	return &GetConfigUseCase{lp: lp, peginLp: peginLp, pegoutLp: pegoutLp}
}

type FullConfiguration struct {
	General liquidity_provider.GeneralConfiguration `json:"general"`
	Pegin   liquidity_provider.PeginConfiguration   `json:"pegin"`
	Pegout  liquidity_provider.PegoutConfiguration  `json:"pegout"`
}

func (useCase *GetConfigUseCase) Run(ctx context.Context) FullConfiguration {
	general := useCase.lp.GeneralConfiguration(ctx)
	pegin := useCase.peginLp.PeginConfiguration(ctx)
	pegout := useCase.pegoutLp.PegoutConfiguration(ctx)

	return FullConfiguration{
		General: general,
		Pegin:   pegin,
		Pegout:  pegout,
	}
}
