package pegin

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type GetCollateralUseCase struct {
	contracts     blockchain.RskContracts
	peginProvider liquidity_provider.LiquidityProvider
}

func NewGetCollateralUseCase(contracts blockchain.RskContracts, peginProvider liquidity_provider.LiquidityProvider) *GetCollateralUseCase {
	return &GetCollateralUseCase{contracts: contracts, peginProvider: peginProvider}
}

func (useCase *GetCollateralUseCase) Run() (*entities.Wei, error) {
	rskAddress := useCase.peginProvider.RskAddress()
	collateral, err := useCase.contracts.CollateralManagement.GetCollateral(rskAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"vertical":    "pegin",
			"rsk_address": rskAddress,
		}).WithError(err).Error("GetCollateral: read failed")
		return nil, usecases.WrapUseCaseError(usecases.GetCollateralId, err)
	}
	return collateral, nil
}
