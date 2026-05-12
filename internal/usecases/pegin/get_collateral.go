package pegin

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetCollateralUseCase struct {
	contracts     blockchain.RskContracts
	peginProvider liquidity_provider.LiquidityProvider
}

func NewGetCollateralUseCase(contracts blockchain.RskContracts, peginProvider liquidity_provider.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{contracts: contracts, peginProvider: peginProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	collateral, err := useCase.contracts.CollateralManagement.GetCollateral(useCase.peginProvider.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetCollateralId, err)
	}
	return collateral, nil
}
