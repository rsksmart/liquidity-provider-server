package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithdrawCollateralUseCase_Run(t *testing.T) {
	lbc := new(mocks.LbcMock)
	lbc.On("WithdrawCollateral").Return(nil)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewWithdrawCollateralUseCase(contracts)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestWithdrawCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	lbc := new(mocks.LbcMock)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewWithdrawCollateralUseCase(contracts)

	lbc.On("WithdrawCollateral").Return(errors.New("LBC021")).Once()
	err := useCase.Run()
	require.ErrorIs(t, err, usecases.ProviderNotResignedError)

	lbc.On("WithdrawCollateral").Return(errors.New("LBC022")).Once()
	err = useCase.Run()
	require.ErrorIs(t, err, usecases.ProviderNotResignedError)

	lbc.On("WithdrawCollateral").Return(assert.AnError).Once()
	err = useCase.Run()
	require.NotErrorIs(t, err, usecases.ProviderNotResignedError)
	require.Error(t, err)
}
