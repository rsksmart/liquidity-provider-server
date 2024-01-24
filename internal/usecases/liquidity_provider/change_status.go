package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ChangeStatusUseCase struct {
	lbc      blockchain.LiquidityBridgeContract
	provider entities.LiquidityProvider
}

func (useCase *ChangeStatusUseCase) Run(newStatus bool) error {
	var err error
	var id uint64

	id, err = ValidateConfiguredProvider(useCase.provider, useCase.lbc)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}

	if err = useCase.lbc.SetProviderStatus(id, newStatus); err != nil {
		return usecases.WrapUseCaseError(usecases.ChangeProviderStatusId, err)
	}
	return nil
}
