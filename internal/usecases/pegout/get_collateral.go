package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetCollateralUseCase struct {
	contracts      blockchain.RskContracts
	pegoutProvider liquidity_provider.LiquidityProvider
}

func NewGetCollateralUseCase(contracts blockchain.RskContracts, pegoutProvider liquidity_provider.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{contracts: contracts, pegoutProvider: pegoutProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	collateral, err := useCase.contracts.Lbc.GetPegoutCollateral(useCase.pegoutProvider.RskAddress())
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetPegoutCollateralId, err)
	}
	return collateral, nil
}
