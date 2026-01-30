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

	// Mock existing state config
	existingTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	existingStateConfig := &entities.Signed[lpEntity.StateConfiguration]{
		Value: lpEntity.StateConfiguration{
			LastBtcToColdWalletTransfer:  &existingTime,
			LastRbtcToColdWalletTransfer: &existingTime,
		},
		Signature: "existing signature",
		Hash:      "existing hash",
	}

	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(existingStateConfig, nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)
	// Should not call UpsertStateConfiguration when config already exists
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_CreateNewStateConfig(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}

	// Mock that state config doesn't exist
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, nil).Once()

	// Capture the upserted config
	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		now := time.Now()
		if config.Value.LastBtcToColdWalletTransfer == nil || config.Value.LastRbtcToColdWalletTransfer == nil {
			return false
		}
		btcDiff := now.Sub(*config.Value.LastBtcToColdWalletTransfer)
		rbtcDiff := now.Sub(*config.Value.LastRbtcToColdWalletTransfer)

		return btcDiff >= 0 && btcDiff < 5*time.Second &&
			rbtcDiff >= 0 && rbtcDiff < 5*time.Second &&
			config.Signature == ""
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	lpRepository.AssertExpectations(t)

	assert.NotNil(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotNil(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Empty(t, capturedConfig.Signature)
}

func TestInitializeStateConfigurationUseCase_Run_GetStateConfigurationError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}

	// Mock GetStateConfiguration to return an error
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	lpRepository.AssertExpectations(t)
	// Should not attempt to upsert if get fails
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_UpsertStateConfigurationError(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}

	// Mock that state config doesn't exist
	lpRepository.On("GetStateConfiguration", test.AnyCtx).Return(nil, nil).Once()

	// Mock UpsertStateConfiguration to return an error
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(lpRepository)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	lpRepository.AssertExpectations(t)
}
