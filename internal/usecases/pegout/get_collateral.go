package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetCollateralUseCase struct {
	lbc            blockchain.LiquidityBridgeContract
	pegoutProvider entities.LiquidityProvider
}

func NewGetCollateralUseCase(lbc blockchain.LiquidityBridgeContract, pegoutProvider entities.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{lbc: lbc, pegoutProvider: pegoutProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	var err error
	collateral := new(entities.Wei)
	collateral, err = useCase.lbc.GetPegoutCollateral(useCase.pegoutProvider.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetPegoutCollateralId, err)
	}
	return collateral, nil
}
