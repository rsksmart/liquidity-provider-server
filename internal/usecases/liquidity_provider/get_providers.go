package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetProvidersUseCase struct {
	lbc blockchain.LiquidityBridgeContract
}

func NewGetProvidersUseCase(lbc blockchain.LiquidityBridgeContract) *GetProvidersUseCase {
	return &GetProvidersUseCase{lbc: lbc}
}

func (useCase *GetProvidersUseCase) Run() ([]entities.RegisteredLiquidityProvider, error) {
	var err error
	var providers []entities.RegisteredLiquidityProvider
	if providers, err = useCase.lbc.GetProviders(); err != nil {
		return providers, usecases.WrapUseCaseError(usecases.GetProvidersId, err)
	}
	return providers, nil
}
