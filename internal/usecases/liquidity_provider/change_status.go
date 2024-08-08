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
	var id uint64

	id, err = ValidateConfiguredProvider(useCase.provider, useCase.contracts.Lbc)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}

	if err = useCase.contracts.Lbc.SetProviderStatus(id, newStatus); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}
	return nil
}
