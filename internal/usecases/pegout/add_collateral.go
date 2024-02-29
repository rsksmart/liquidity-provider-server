package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type AddCollateralUseCase struct {
	lbc blockchain.LiquidityBridgeContract
	lp  liquidity_provider.LiquidityProvider
}

func NewAddCollateralUseCase(lbc blockchain.LiquidityBridgeContract, lp liquidity_provider.LiquidityProvider) *AddCollateralUseCase {
	return &AddCollateralUseCase{lbc: lbc, lp: lp}
}

func (useCase *AddCollateralUseCase) Run(amount *entities.Wei) (*entities.Wei, error) {
	var err error
	minCollateral, err := useCase.lbc.GetMinimumCollateral()
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, err)
	}
	collateral, err := useCase.lbc.GetPegoutCollateral(useCase.lp.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, err)
	}
	result := new(entities.Wei)
	result.Add(collateral, amount)
	if minCollateral.Cmp(result) > 0 {
		return nil, usecases.WrapUseCaseError(usecases.AddCollateralId, usecases.InsufficientAmountError)
	}
	err = useCase.lbc.AddPegoutCollateral(amount)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.AddPegoutCollateralId, err)
	}
	return result, nil
}
