package liquidity_provider_test

import (
	"errors"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeStatusUseCase_Run(t *testing.T) {
	lbc := &test.LbcMock{}
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

	provider := &test.ProviderMock{}
	provider.On("RskAddress").Return("0x02")

	err := liquidity_provider.NewChangeStatusUseCase(lbc, provider).Run(false)

	lbc.AssertExpectations(t)
	require.NoError(t, err)
}

func TestChangeStatusUseCase_Run_Fail(t *testing.T) {
	lbc := &test.LbcMock{}
	provider := &test.ProviderMock{}

	lbc.On("GetProviders").Return(
		[]lp.RegisteredLiquidityProvider{},
		errors.New("some error"),
	).Once()
	err := liquidity_provider.NewChangeStatusUseCase(lbc, provider).Run(false)
	lbc.AssertExpectations(t)
	require.Error(t, err)

	lbc.On("GetProviders").Return([]lp.RegisteredLiquidityProvider{
		{Id: 1, Address: "0x01"},
	}, nil).Once()
	provider.On("RskAddress").Return("0x01")
	lbc.On("SetProviderStatus", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
	err = liquidity_provider.NewChangeStatusUseCase(lbc, provider).Run(false)
	lbc.AssertExpectations(t)
	require.Error(t, err)
}
