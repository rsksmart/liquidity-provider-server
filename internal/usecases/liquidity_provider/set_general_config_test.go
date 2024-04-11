package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetGeneralConfigUseCase_Run(t *testing.T) {
	config := entities.Signed[lp.GeneralConfiguration]{
		Value: lp.GeneralConfiguration{
			RskConfirmations: map[int]uint16{5: 10},
			BtcConfirmations: map[int]uint16{10: 20},
		},
		Signature: "010203",
		Hash:      "040506",
	}

	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("UpsertGeneralConfiguration", test.AnyCtx, config).Return(nil)
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	useCase := liquidity_provider.NewSetGeneralConfigUseCase(lpRepository, walletMock, hashMock.Hash)

	err := useCase.Run(context.Background(), config.Value)
	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
}

func TestSetGeneralConfigUseCase_Run_ErrorHandling(t *testing.T) {
	config := entities.Signed[lp.GeneralConfiguration]{
		Value: lp.GeneralConfiguration{
			RskConfirmations: map[int]uint16{5: 10},
			BtcConfirmations: map[int]uint16{10: 20},
		},
		Signature: "010203",
		Hash:      "040506",
	}

	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	errorSetups := []func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock){
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)
		},
		func(lpRepository *mocks.LiquidityProviderRepositoryMock, walletMock *mocks.RskWalletMock) {
			walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
			lpRepository.On("UpsertGeneralConfiguration", test.AnyCtx, config).Return(assert.AnError)
		},
	}

	for _, errorSetup := range errorSetups {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		errorSetup(lpRepository, walletMock)
		useCase := liquidity_provider.NewSetGeneralConfigUseCase(lpRepository, walletMock, hashMock.Hash)
		err := useCase.Run(context.Background(), config.Value)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		walletMock.AssertExpectations(t)
	}
}
