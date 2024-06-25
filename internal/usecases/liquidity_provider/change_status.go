package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ChangeStatusUseCase struct {
	contracts blockchain.RskContracts
	provider  liquidity_provider.LiquidityProvider
}

func NewChangeStatusUseCase(contracts blockchain.RskContracts, provider liquidity_provider.LiquidityProvider) *ChangeStatusUseCase {
	return &ChangeStatusUseCase{contracts: contracts, provider: provider}
}

func (useCase *ChangeStatusUseCase) Run(newStatus bool) error {
	var err error
	var registeredProvider liquidity_provider.RegisteredLiquidityProvider

	registeredProvider, err = useCase.contracts.Lbc.GetProvider(useCase.provider.RskAddress())
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}

	if err = useCase.contracts.Lbc.SetProviderStatus(registeredProvider.Id, newStatus); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}
	return nil
}
