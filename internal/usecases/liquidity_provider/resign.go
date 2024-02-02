package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ResignUseCase struct {
	lbc      blockchain.LiquidityBridgeContract
	provider entities.LiquidityProvider
}

func NewResignUseCase(lbc blockchain.LiquidityBridgeContract, provider entities.LiquidityProvider) *ResignUseCase {
	return &ResignUseCase{lbc: lbc, provider: provider}
}

func (useCase *ResignUseCase) Run() error {
	var err error

	_, err = ValidateConfiguredProvider(useCase.provider, useCase.lbc)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, err)
	}

	if err = useCase.lbc.ProviderResign(); err != nil {
		return usecases.WrapUseCaseError(usecases.ProviderResignId, err)
	}
	return nil
}