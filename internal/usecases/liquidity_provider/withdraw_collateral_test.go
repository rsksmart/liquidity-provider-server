package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithdrawCollateralUseCase_Run(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	collateral.On("WithdrawCollateral").Return(nil)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := liquidity_provider.NewWithdrawCollateralUseCase(contracts)
	err := useCase.Run()
	collateral.AssertExpectations(t)
	require.NoError(t, err)
}

func TestWithdrawCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := liquidity_provider.NewWithdrawCollateralUseCase(contracts)

	collateral.On("WithdrawCollateral").Return(lpEntity.ProviderNotResignedError).Once()
	err := useCase.Run()
	require.ErrorIs(t, err, lpEntity.ProviderNotResignedError)

	collateral.On("WithdrawCollateral").Return(errors.New("some error")).Once()
	err = useCase.Run()
	require.ErrorContains(t, err, "some error")

	collateral.On("WithdrawCollateral").Return(assert.AnError).Once()
	err = useCase.Run()
	require.NotErrorIs(t, err, lpEntity.ProviderNotResignedError)
	require.Error(t, err)
}
