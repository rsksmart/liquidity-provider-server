package pegout

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
	if err = usecases.CheckPauseState(useCase.contracts.CollateralManagement); err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, err)
	}
	minCollateral, err := useCase.contracts.CollateralManagement.GetMinimumCollateral()
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, err)
	}
	collateral, err := useCase.contracts.CollateralManagement.GetPegoutCollateral(useCase.lp.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, err)
	}
	result := new(entities.Wei)
	result.Add(collateral, amount)
	if minCollateral.Cmp(result) > 0 {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, usecases.InsufficientAmountError)
	}
	err = useCase.contracts.CollateralManagement.AddPegoutCollateral(amount)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, err)
	}
	return result, nil
}
