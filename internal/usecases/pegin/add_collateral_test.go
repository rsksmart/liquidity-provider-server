package pegin_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddCollateralUseCase_Run(t *testing.T) {
	lbc := new(test.LbcMock)
	lp := new(test.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("AddCollateral", value).Return(nil)
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(100), nil)
	useCase := pegin.NewAddCollateralUseCase(lbc, lp)
	result, err := useCase.Run(value)
	lp.AssertExpectations(t)
	lbc.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(1100), result)
}

func TestAddCollateralUseCase_Run_NotEnough(t *testing.T) {
	lbc := new(test.LbcMock)
	lp := new(test.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(2000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(100), nil)
	useCase := pegin.NewAddCollateralUseCase(lbc, lp)
	result, err := useCase.Run(value)
	lp.AssertExpectations(t)
	lbc.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddCollateral", mock.Anything)
	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	assert.Nil(t, result)
}

func TestAddCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	lp := new(test.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	cases := test.Table[func(lbc *test.LbcMock), error]{
		{
			Value: func(lbc *test.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(nil, assert.AnError)
			},
		},
		{
			Value: func(lbc *test.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
				lbc.On("GetCollateral", mock.Anything).Return(nil, assert.AnError)
			},
		},
		{
			Value: func(lbc *test.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(100), nil)
				lbc.On("AddCollateral", mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, c := range cases {
		lbc := new(test.LbcMock)
		c.Value(lbc)
		useCase := pegin.NewAddCollateralUseCase(lbc, lp)
		result, err := useCase.Run(entities.NewWei(100))
		lbc.AssertExpectations(t)
		assert.Nil(t, result)
		require.Error(t, err)
	}
}
