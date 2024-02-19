package pegin_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithdrawCollateralUseCase_Run(t *testing.T) {
	lbc := new(test.LbcMock)
	lbc.On("WithdrawCollateral").Return(nil)
	useCase := pegin.NewWithdrawCollateralUseCase(lbc)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestWithdrawCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	lbc := new(test.LbcMock)
	useCase := pegin.NewWithdrawCollateralUseCase(lbc)

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
