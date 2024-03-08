package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ResignUseCase struct {
	contracts blockchain.RskContracts
	provider  liquidity_provider.LiquidityProvider
}

func NewResignUseCase(contracts blockchain.RskContracts, provider liquidity_provider.LiquidityProvider) *ResignUseCase {
	return &ResignUseCase{contracts: contracts, provider: provider}
}

func (useCase *ResignUseCase) Run() error {
	var err error

	_, err = ValidateConfiguredProvider(useCase.provider, useCase.contracts.Lbc)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, err)
	}

	if err = useCase.contracts.Lbc.ProviderResign(); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, err)
	}
	return nil
}
