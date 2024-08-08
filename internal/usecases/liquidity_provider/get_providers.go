package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetProvidersUseCase struct {
	contracts blockchain.RskContracts
}

func NewGetProvidersUseCase(contracts blockchain.RskContracts) *GetProvidersUseCase {
	return &GetProvidersUseCase{contracts: contracts}
}

func (useCase *GetProvidersUseCase) Run() ([]liquidity_provider.RegisteredLiquidityProvider, error) {
	var err error
	var providers []liquidity_provider.RegisteredLiquidityProvider
	if providers, err = useCase.contracts.Lbc.GetProviders(); err != nil {
		return providers, usecases.WrapUseCaseError(usecases.GetProvidersId, err)
	}
	return providers, nil
}
