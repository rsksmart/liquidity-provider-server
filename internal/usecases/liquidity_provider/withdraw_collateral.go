package liquidity_provider

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type WithdrawCollateralUseCase struct {
	contracts blockchain.RskContracts
}

func NewWithdrawCollateralUseCase(contracts blockchain.RskContracts) *WithdrawCollateralUseCase {
	return &WithdrawCollateralUseCase{contracts: contracts}
}

func (useCase *WithdrawCollateralUseCase) Run() error {
	err := useCase.contracts.CollateralManagement.WithdrawCollateral()
	if err != nil {
		return usecases.WrapUseCaseError(usecases.WithdrawCollateralId, err)
	}
	return nil
}
