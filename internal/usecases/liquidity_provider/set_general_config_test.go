package liquidity_provider_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetGeneralConfigUseCase_Run(t *testing.T) {
	config := entities.Signed[lp.GeneralConfiguration]{
		Value: lp.GeneralConfiguration{
			RskConfirmations: map[string]uint16{"5": 10},
			BtcConfirmations: map[string]uint16{"10": 20},
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
			RskConfirmations: map[string]uint16{"5": 10},
			BtcConfirmations: map[string]uint16{"10": 20},
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

func TestSetGeneralConfigUseCase_Run_ValidateConfirmations(t *testing.T) {
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	invalidConfigs := []lp.GeneralConfiguration{
		// Empty RSK confirmations
		{
			RskConfirmations: map[string]uint16{},
			BtcConfirmations: map[string]uint16{"10": 20},
		},
		// Empty BTC confirmations
		{
			RskConfirmations: map[string]uint16{"5": 10},
			BtcConfirmations: map[string]uint16{},
		},
		// Negative key in RSK confirmations
		{
			RskConfirmations: map[string]uint16{"-1": 10},
			BtcConfirmations: map[string]uint16{"10": 20},
		},
		// Zero key in BTC confirmations
		{
			RskConfirmations: map[string]uint16{"5": 10},
			BtcConfirmations: map[string]uint16{"0": 20},
		},
	}

	for _, cfg := range invalidConfigs {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		useCase := liquidity_provider.NewSetGeneralConfigUseCase(lpRepository, walletMock, hashMock.Hash)
		err := useCase.Run(context.Background(), cfg)
		require.Error(t, err)
		if len(cfg.RskConfirmations) == 0 || len(cfg.BtcConfirmations) == 0 {
			require.ErrorIs(t, err, usecases.EmptyConfirmationsMapError)
		} else {
			require.ErrorIs(t, err, usecases.NonPositiveConfirmationKeyError)
		}
	}
}
