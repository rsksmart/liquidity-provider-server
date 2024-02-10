package pegout_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCollateralUseCase_Run(t *testing.T) {
	lbc := new(test.LbcMock)
	lp := new(test.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetPegoutCollateral", "rskAddress").Return(value, nil)
	useCase := pegout.NewGetCollateralUseCase(lbc, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, value, result)
}

func TestGetCollateralUseCase_Run_Error(t *testing.T) {
	lbc := new(test.LbcMock)
	lp := new(test.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetPegoutCollateral", "rskAddress").Return(entities.NewWei(0), assert.AnError)
	useCase := pegout.NewGetCollateralUseCase(lbc, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
