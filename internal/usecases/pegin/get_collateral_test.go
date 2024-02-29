package pegin_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCollateralUseCase_Run(t *testing.T) {
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetCollateral", "rskAddress").Return(value, nil)
	useCase := pegin.NewGetCollateralUseCase(lbc, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetCollateralUseCase_Run_Error(t *testing.T) {
	lbc := new(mocks.LbcMock)
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetCollateral", "rskAddress").Return(entities.NewWei(0), assert.AnError)
	useCase := pegin.NewGetCollateralUseCase(lbc, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	require.Error(t, err)
	assert.Nil(t, result)
}
