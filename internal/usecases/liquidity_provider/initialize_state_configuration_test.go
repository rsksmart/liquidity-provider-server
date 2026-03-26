package liquidity_provider_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func configuredColdWalletMock() *mocks.ColdWalletMock {
	coldWallet := new(mocks.ColdWalletMock)
	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)
	return coldWallet
}

func initHashAddress(address string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(address)))
}

func TestInitializeStateConfigurationUseCase_Run_StateConfigAlreadyExists(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	existingConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:   existingUnix,
		LastRbtcToColdWalletTransfer:  existingUnix,
		BtcColdWalletAddressHash:      initHashAddress(test.AnyBtcAddress),
		RskColdWalletAddressHash:      initHashAddress(test.AnyRskAddress),
		BtcLiquidityTargetPercentage:  50,
	}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(existingConfig, nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_CreateNewStateConfig(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	expectedBtcHash := initHashAddress(test.AnyBtcAddress)
	expectedRskHash := initHashAddress(test.AnyRskAddress)

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		nowUnix := time.Now().Unix()
		btcDiff := nowUnix - config.Value.LastBtcToColdWalletTransfer
		rbtcDiff := nowUnix - config.Value.LastRbtcToColdWalletTransfer

		return btcDiff >= 0 && btcDiff < 5 &&
			rbtcDiff >= 0 && rbtcDiff < 5 &&
			config.Value.BtcColdWalletAddressHash == expectedBtcHash &&
			config.Value.RskColdWalletAddressHash == expectedRskHash &&
			config.Signature != ""
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)

	assert.NotZero(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotZero(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Equal(t, expectedBtcHash, capturedConfig.Value.BtcColdWalletAddressHash)
	assert.Equal(t, expectedRskHash, capturedConfig.Value.RskColdWalletAddressHash)
	assert.Equal(t, "010203", capturedConfig.Signature)
}

func TestInitializeStateConfigurationUseCase_Run_UpsertStateConfigurationError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestInitializeStateConfigurationUseCase_Run_PartialInitialization(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	existingUnix := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC).Unix()
	existingBtcHash := initHashAddress(test.AnyBtcAddress)
	partialConfig := lpEntity.StateConfiguration{
		LastBtcToColdWalletTransfer:  existingUnix,
		LastRbtcToColdWalletTransfer: 0,
		BtcColdWalletAddressHash:     existingBtcHash,
		RskColdWalletAddressHash:     "",
	}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(partialConfig, nil).Once()

	walletMock.On("SignBytes", mock.Anything).Return([]byte{7, 8, 9}, nil)

	expectedRskHash := initHashAddress(test.AnyRskAddress)

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

		return config.Value.BtcColdWalletAddressHash == existingBtcHash &&
			config.Value.RskColdWalletAddressHash == expectedRskHash &&
			config.Signature != "" && config.Signature != "old signature"
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)

	assert.NotZero(t, capturedConfig.Value.LastBtcToColdWalletTransfer)
	assert.NotZero(t, capturedConfig.Value.LastRbtcToColdWalletTransfer)
	assert.Equal(t, existingUnix, capturedConfig.Value.LastBtcToColdWalletTransfer, "BTC timestamp should remain unchanged")
	assert.Equal(t, existingBtcHash, capturedConfig.Value.BtcColdWalletAddressHash, "BTC hash should remain unchanged")
	assert.Equal(t, expectedRskHash, capturedConfig.Value.RskColdWalletAddressHash, "RSK hash should be initialized")
	assert.Equal(t, "070809", capturedConfig.Signature)
}

func TestInitializeStateConfigurationUseCase_Run_ProviderError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, assert.AnError).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
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
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()

	walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "InitializeStateConfiguration")
	providerMock.AssertExpectations(t)
	walletMock.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_BtcAddressEmpty_ReturnsError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	walletMock := &mocks.RskWalletMock{}

	coldWallet.On("GetBtcAddress").Return("")
	coldWallet.On("GetRskAddress").Return(test.AnyRskAddress)

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet BTC address not configured")
	coldWallet.AssertExpectations(t)
	providerMock.AssertNotCalled(t, "StateConfiguration", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestInitializeStateConfigurationUseCase_Run_FirstRunInitializesPercentage(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := configuredColdWalletMock()
	walletMock := &mocks.RskWalletMock{}

	providerMock.On("StateConfiguration", test.AnyCtx).Return(lpEntity.StateConfiguration{}, lpEntity.ConfigurationNotFoundError).Once()
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	var capturedConfig entities.Signed[lpEntity.StateConfiguration]
	lpRepository.On("UpsertStateConfiguration", test.AnyCtx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		capturedConfig = config
		return true
	})).Return(nil).Once()

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.NoError(t, err)
	providerMock.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)

	assert.Equal(t, uint64(50), capturedConfig.Value.BtcLiquidityTargetPercentage)
	assert.Equal(t, int64(0), capturedConfig.Value.RatioCooldownEndTimestamp, "first-run should not activate cooldown")
}

func TestInitializeStateConfigurationUseCase_Run_RskAddressEmpty_ReturnsError(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	coldWallet := new(mocks.ColdWalletMock)
	walletMock := &mocks.RskWalletMock{}

	coldWallet.On("GetBtcAddress").Return(test.AnyBtcAddress)
	coldWallet.On("GetRskAddress").Return("")

	useCase := liquidity_provider.NewInitializeStateConfigurationUseCase(providerMock, lpRepository, coldWallet, walletMock, crypto.Keccak256)
	err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cold wallet RSK address not configured")
	coldWallet.AssertExpectations(t)
	providerMock.AssertNotCalled(t, "StateConfiguration", mock.Anything)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}
