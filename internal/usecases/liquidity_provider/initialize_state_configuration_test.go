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
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock existing state config with all fields initialized
	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	existingStateConfig := &entities.Signed[lpEntity.StateConfiguration]{
		Value: lpEntity.StateConfiguration{
			LastBtcToColdWalletTransfer:  &existingUnix,
			LastRbtcToColdWalletTransfer: &existingUnix,
		},
		Signature: "existing signature",
		Hash:      "existing hash",
	}

	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(existingStateConfig, nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	// Should not call UpsertStateConfiguration when all fields are already initialized
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_CreateNewStateConfig(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock that state config doesn't exist
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, nil).Once()

	// Mock signing operations
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	// Capture the upserted config
	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		nowUnix := time.Now().Unix()
		if config.Value.LastBtcToColdWalletTransfer == nil || config.Value.LastRbtcToColdWalletTransfer == nil {
			return false
		}
		btcDiff := nowUnix - *config.Value.LastBtcToColdWalletTransfer
		rbtcDiff := nowUnix - *config.Value.LastRbtcToColdWalletTransfer

		return btcDiff >= 0 && btcDiff < 5 &&
			rbtcDiff >= 0 && rbtcDiff < 5 &&
			config.Signature != ""
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)

	assert.NotNil(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotNil(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.NotEmpty(t, capturedConfig.Signature)
	assert.Equal(t, "010203", capturedConfig.Signature)
	assert.Equal(t, "040506", capturedConfig.Hash)
}

func TestInitializeStateConfigurationUseCase_Run_GetStateConfigurationError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock GetStateConfiguration to return an error
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	lpRepository.AssertExpectations(t)
	// Should not attempt to upsert if get fails
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_UpsertStateConfigurationError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock that state config doesn't exist
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, nil).Once()

	// Mock signing operations
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	// Mock UpsertStateConfiguration to return an error
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	lpRepository.AssertExpectations(t)
}

func TestInitializeStateConfigurationUseCase_Run_PartialInitialization(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock existing state config with only BTC field initialized
	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	existingStateConfig := &entities.Signed[lpEntity.StateConfiguration]{
		Value: lpEntity.StateConfiguration{
			LastBtcToColdWalletTransfer:  &existingUnix,
			LastRbtcToColdWalletTransfer: nil,
		},
		Signature: "old signature",
		Hash:      "old hash",
	}

	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(existingStateConfig, nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{7, 8, 9}, nil)
	hashMock.On("Hash", mock.Anything).Return([]byte{10, 11, 12})

	// Capture the upserted config
	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		nowUnix := time.Now().Unix()

		// BTC field should remain unchanged
		if config.Value.LastBtcToColdWalletTransfer == nil {
			return false
		}
		if *config.Value.LastBtcToColdWalletTransfer != existingUnix {
			return false
		}

		// RBTC field should be newly initialized
		if config.Value.LastRbtcToColdWalletTransfer == nil {
			return false
		}
		rbtcDiff := nowUnix - *config.Value.LastRbtcToColdWalletTransfer
		if rbtcDiff < 0 || rbtcDiff >= 5 {
			return false
		}

		// New signature should be generated
		return config.Signature != "" && config.Signature != "old signature"
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)

	// Verify captured config
	assert.NotNil(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotNil(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Equal(t, existingUnix, *capturedConfig.Value.LastBtcToColdWalletTransfer, "BTC timestamp should remain unchanged")
	assert.NotEqual(t, "old signature", capturedConfig.Signature, "New signature should be generated")
	assert.Equal(t, "070809", capturedConfig.Signature)
	assert.Equal(t, "0a0b0c", capturedConfig.Hash)
}

func TestInitializeStateConfigurationUseCase_Run_SigningError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	// Mock that state config doesn't exist
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, nil).Once()

	hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})

	walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository, walletMock, hashMock.Hash)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	hashMock.AssertExpectations(t)
	// Should not attempt to upsert if signing fails
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}
