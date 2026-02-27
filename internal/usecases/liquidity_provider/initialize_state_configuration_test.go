package liquidity_provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitializeStateConfigurationUseCase_Run_StateConfigAlreadyExists(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	existingConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  existingUnix,
		LastRbtcToColdWalletTransfer: existingUnix,
	}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(existingConfig, nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_CreateNewStateConfig(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		nowUnix := time.Now().Unix()
		btcDiff := nowUnix - config.Value.LastBtcToColdWalletTransfer
		rbtcDiff := nowUnix - config.Value.LastRbtcToColdWalletTransfer

		return btcDiff >= 0 && btcDiff < 5 &&
			rbtcDiff >= 0 && rbtcDiff < 5 &&
			config.Signature != ""
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)

	assert.NotZero(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotZero(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Equal(t, "010203", capturedConfig.Signature)
	assert.Equal(t, "040506", capturedConfig.Hash)
}

func TestInitializeStateConfigurationUseCase_Run_UpsertStateConfigurationError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestInitializeStateConfigurationUseCase_Run_PartialInitialization(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	partialConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  existingUnix,
		LastRbtcToColdWalletTransfer: 0,
	}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(partialConfig, nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{7, 8, 9}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{10, 11, 12})

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		nowUnix := time.Now().Unix()

		if config.Value.LastBtcToColdWalletTransfer != existingUnix {
			return false
		}

		rbtcDiff := nowUnix - config.Value.LastRbtcToColdWalletTransfer
		if rbtcDiff < 0 || rbtcDiff >= 5 {
			return false
		}

		return config.Signature != "" && config.Signature != "old signature"
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)

	assert.NotZero(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotZero(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Equal(t, existingUnix, capturedConfig.Value.LastBtcToColdWalletTransfer, "BTC timestamp should remain unchanged")
	assert.Equal(t, "070809", capturedConfig.Signature)
	assert.Equal(t, "0a0b0c", capturedConfig.Hash)
}

func TestInitializeStateConfigurationUseCase_Run_ProviderError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	providerMock.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_SigningError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
	walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	providerMock.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}
