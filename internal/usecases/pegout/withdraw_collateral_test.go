package pegout_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithdrawCollateralUseCase_Run(t *testing.T) {
	lbc := new(test.LbcMock)
	lbc.On("WithdrawPegoutCollateral").Return(nil)
	useCase := pegout.NewWithdrawCollateralUseCase(lbc)
	err := useCase.Run()
	lbc.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestWithdrawCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	lbc := new(test.LbcMock)
	useCase := pegout.NewWithdrawCollateralUseCase(lbc)

	lbc.On("WithdrawPegoutCollateral").Return(errors.New("LBC021")).Once()
	err := useCase.Run()
	assert.ErrorIs(t, err, usecases.ProviderNotResignedError)

	lbc.On("WithdrawPegoutCollateral").Return(errors.New("LBC022")).Once()
	err = useCase.Run()
	assert.ErrorIs(t, err, usecases.ProviderNotResignedError)

	lbc.On("WithdrawPegoutCollateral").Return(assert.AnError).Once()
	err = useCase.Run()
	assert.NotErrorIs(t, err, usecases.ProviderNotResignedError)
	assert.NotNil(t, err)
}
