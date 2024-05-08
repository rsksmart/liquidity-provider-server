package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"strings"
)

type WithdrawCollateralUseCase struct {
	contracts blockchain.RskContracts
}

func NewWithdrawCollateralUseCase(contracts blockchain.RskContracts) *WithdrawCollateralUseCase {
	return &WithdrawCollateralUseCase{contracts: contracts}
}

func (useCase *WithdrawCollateralUseCase) Run() error {
	err := useCase.contracts.Lbc.WithdrawPegoutCollateral()
	if err != nil && (strings.Contains(err.Error(), "LBC021") || strings.Contains(err.Error(), "LBC022")) {
		return usecases.WrapUseCaseError(usecases.WithdrawCollateralId, usecases.ProviderNotResignedError)
	} else if err != nil {
		return usecases.WrapUseCaseError(usecases.WithdrawCollateralId, err)
	}
	return nil
}