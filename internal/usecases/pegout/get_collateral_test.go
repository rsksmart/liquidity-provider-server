package pegout_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCollateralUseCase_Run(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	lp := new(mocks.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	collateral.On("GetPegoutCollateral", "rskAddress").Return(value, nil)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := pegout.NewGetCollateralUseCase(contracts, lp)
	result, err := useCase.Run()
	collateral.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetCollateralUseCase_Run_Error(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	collateral.On("GetPegoutCollateral", "rskAddress").Return(entities.NewWei(0), assert.AnError)
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := pegout.NewGetCollateralUseCase(contracts, lp)
	result, err := useCase.Run()
	collateral.AssertExpectations(t)
	require.Error(t, err)
	assert.Nil(t, result)
}
