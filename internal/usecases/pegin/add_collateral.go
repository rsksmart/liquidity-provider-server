package pegin

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type AddCollateralUseCase struct {
	contracts blockchain.RskContracts
	lp        liquidity_provider.LiquidityProvider
}

func NewAddCollateralUseCase(contracts blockchain.RskContracts, lp liquidity_provider.LiquidityProvider) *AddCollateralUseCase {
	return &AddCollateralUseCase{contracts: contracts, lp: lp}
}

func (useCase *AddCollateralUseCase) Run(amount *entities.Wei) (*entities.Wei, error) {
	var err error
	minCollateral, err := useCase.contracts.Lbc.GetMinimumCollateral()
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, err)
	}
	collateral, err := useCase.contracts.Lbc.GetCollateral(useCase.lp.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, err)
	}
	result := new(entities.Wei)
	result.Add(collateral, amount)
	if minCollateral.Cmp(result) > 0 {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, usecases.InsufficientAmountError)
	}
	err = useCase.contracts.Lbc.AddCollateral(amount)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, err)
	}
	return result, nil
}
