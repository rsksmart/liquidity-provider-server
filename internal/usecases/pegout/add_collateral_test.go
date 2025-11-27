package pegout_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddCollateralUseCase_Run(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	lp := new(mocks.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	collateral.On("AddPegoutCollateral", value).Return(nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(100), nil)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := pegout.NewAddCollateralUseCase(contracts, lp)
	result, err := useCase.Run(value)
	lp.AssertExpectations(t)
	collateral.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(1100), result)
}

func TestAddCollateralUseCase_Run_NotEnough(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	lp := new(mocks.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(2000), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(100), nil)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := pegout.NewAddCollateralUseCase(contracts, lp)
	result, err := useCase.Run(value)
	lp.AssertExpectations(t)
	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddPegoutCollateral", mock.Anything)
	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	assert.Nil(t, result)
}

func TestAddCollateralUseCase_Run_ErrorHandling(t *testing.T) {
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	cases := test.Table[func(collateral *mocks.CollateralManagementContractMock), error]{
		{
			Value: func(collateral *mocks.CollateralManagementContractMock) {
				collateral.On("GetMinimumCollateral").Return(nil, assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(nil, assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(100), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(100), nil)
				collateral.On("AddPegoutCollateral", mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, c := range cases {
		collateral := new(mocks.CollateralManagementContractMock)
		c.Value(collateral)
		contracts := blockchain.RskContracts{CollateralManagement: collateral}
		useCase := pegout.NewAddCollateralUseCase(contracts, lp)
		result, err := useCase.Run(entities.NewWei(100))
		collateral.AssertExpectations(t)
		assert.Nil(t, result)
		require.Error(t, err)
	}
}
