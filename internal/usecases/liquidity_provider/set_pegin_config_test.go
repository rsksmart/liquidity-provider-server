package liquidity_provider_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var peginConfigMock = entities.Signed[lp.PeginConfiguration]{
	Value: lp.PeginConfiguration{
		TimeForDeposit: 1,
		CallTime:       2,
		PenaltyFee:     entities.NewWei(3),
		FixedFee:       entities.NewWei(4),
		PercentageFee:  big.NewFloat(1.25),
		MaxValue:       entities.NewWei(5),
		MinValue:       entities.NewWei(1),
	},
	Signature: "010203",
	Hash:      "040506",
}

func TestSetPeginConfigUseCase_Run(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("UpsertPeginConfiguration", test.AnyCtx, peginConfigMock).Return(nil)
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash)

	err := useCase.Run(context.Background(), peginConfigMock.Value)
	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
}

func TestSetPeginConfigUseCase_Run_ErrorHandling(t *testing.T) {
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	errorSetups := []func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock){
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)
		},
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
			lpRepository.On("UpsertPeginConfiguration", test.AnyCtx, peginConfigMock).Return(assert.AnError)
		},
	}

	for _, errorSetup := range errorSetups {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		errorSetup(lpRepository, walletMock)
		useCase := liquidity_provider.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash)
		err := useCase.Run(context.Background(), peginConfigMock.Value)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		walletMock.AssertExpectations(t)
	}
}
