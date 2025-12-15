package pegout_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCollateralUseCase_Run(t *testing.T) {
	lbc := new(mocks.LiquidityBridgeContractMock)
	lp := new(mocks.ProviderMock)
	value := entities.NewWei(1000)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetPegoutCollateral", "rskAddress").Return(value, nil)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewGetCollateralUseCase(contracts, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestGetCollateralUseCase_Run_Error(t *testing.T) {
	lbc := new(mocks.LiquidityBridgeContractMock)
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return("rskAddress")
	lbc.On("GetPegoutCollateral", "rskAddress").Return(entities.NewWei(0), assert.AnError)
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := pegout.NewGetCollateralUseCase(contracts, lp)
	result, err := useCase.Run()
	lbc.AssertExpectations(t)
	require.Error(t, err)
	assert.Nil(t, result)
}
