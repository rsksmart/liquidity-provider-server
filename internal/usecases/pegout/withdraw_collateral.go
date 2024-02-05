package pegout

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"strings"
)

type WithdrawCollateralUseCase struct {
	lbc blockchain.LiquidityBridgeContract
}

func NewWithdrawCollateralUseCase(lbc blockchain.LiquidityBridgeContract) *WithdrawCollateralUseCase {
	return &WithdrawCollateralUseCase{lbc: lbc}
}

func (useCase *WithdrawCollateralUseCase) Run() error {
	var err error
	err = useCase.lbc.WithdrawPegoutCollateral()
	if err != nil && (strings.Contains(err.Error(), "LBC021") || strings.Contains(err.Error(), "LBC022")) {
		return usecases.WrapUseCaseError(usecases.WithdrawCollateralId, usecases.ProviderNotResignedError)
	} else if err != nil {
		return usecases.WrapUseCaseError(usecases.WithdrawCollateralId, err)
	}
	return nil
}
