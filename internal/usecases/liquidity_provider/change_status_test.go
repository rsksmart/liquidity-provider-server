package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeStatusUseCase_Run(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{
			Id:      1,
			Address: "0x01",
		},
		{
			Id:      2,
			Address: "0x02",
		},
		{
			Id:      3,
			Address: "0x03",
		},
	}, nil).Once()
	lbc.On("SetProviderStatus", uint64(2), false).Return(nil).Once()

	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("0x02")

	contracts := blockchain.RskContracts{Lbc: lbc}
	err := liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)

	lbc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestChangeStatusUseCase_Run_Fail(t *testing.T) {
	lbc := &mocks.LbcMock{}
	provider := &mocks.ProviderMock{}

	lbc.On("GetProviders").Return(
		[]lp.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()
	contracts := blockchain.RskContracts{Lbc: lbc}
	err := liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)
	lbc.AssertExpectations(t)
	require.Error(t, err)

	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{Id: 1, Address: "0x01"},
	}, nil).Once()
	provider.On("RskAddress").Return("0x01")
	lbc.On("SetProviderStatus", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
	err = liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)
	lbc.AssertExpectations(t)
	require.Error(t, err)
}
