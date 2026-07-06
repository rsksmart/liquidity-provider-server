package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type GetCollateralUseCase struct {
	contracts      blockchain.RskContracts
	pegoutProvider liquidity_provider.LiquidityProvider
}

func NewGetCollateralUseCase(contracts blockchain.RskContracts, pegoutProvider liquidity_provider.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{contracts: contracts, pegoutProvider: pegoutProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	rskAddress := useCase.pegoutProvider.RskAddress()
	collateral, err := useCase.contracts.CollateralManagement.GetPegoutCollateral(rskAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"vertical":    "pegout",
			"rsk_address": rskAddress,
		}).WithError(err).Error("GetCollateral: read failed")
		return nil, usecases.WrapUseCaseError(usecases.GetPegoutCollateralId, err)
	}
	return collateral, nil
}
