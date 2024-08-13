package liquidity_provider

import (
	"errors"
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

	_, err = useCase.contracts.Lbc.GetProvider(useCase.provider.RskAddress())
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, errors.Join(err, usecases.ProviderConfigurationError))
	}

	if err = useCase.contracts.Lbc.ProviderResign(); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, err)
	}
	return nil
}
