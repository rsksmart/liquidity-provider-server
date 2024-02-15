package pegin

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetCollateralUseCase struct {
	lbc           blockchain.LiquidityBridgeContract
	peginProvider entities.LiquidityProvider
}

func NewGetCollateralUseCase(lbc blockchain.LiquidityBridgeContract, peginProvider entities.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{lbc: lbc, peginProvider: peginProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	collateral, err := useCase.lbc.GetCollateral(useCase.peginProvider.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetCollateralId, err)
	}
	return collateral, nil
}
