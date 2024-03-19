package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetProvidersUseCase struct {
	lbc blockchain.LiquidityBridgeContract
}

func NewGetProvidersUseCase(lbc blockchain.LiquidityBridgeContract) *GetProvidersUseCase {
	return &GetProvidersUseCase{lbc: lbc}
}

func (useCase *GetProvidersUseCase) Run() ([]liquidity_provider.RegisteredLiquidityProvider, error) {
	var err error
	var providers []liquidity_provider.RegisteredLiquidityProvider
	if providers, err = useCase.lbc.GetProviders(); err != nil {
		return providers, usecases.WrapUseCaseError(usecases.GetProvidersId, err)
	}
	return providers, nil
}
