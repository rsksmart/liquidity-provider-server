package liquidity_provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetLiquidityRatioUseCase_Run_HappyPath(t *testing.T) {
	ctx := context.Background()
	provider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	walletMock := new(mocks.RskWalletMock)

	stateConfig := lpEntity.StateConfiguration{
		BtcLiquidityTargetPercentage: 50,
		RatioCooldownEndTimestamp:    0,
	}
	provider.On("StateConfiguration", ctx).Return(stateConfig, nil)
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)

	now := time.Now().Unix()
	lpRepository.On("UpsertStateConfiguration", ctx, mock.MatchedBy(func(config entities.Signed[lpEntity.StateConfiguration]) bool {
		return config.Value.BtcLiquidityTargetPercentage == 70 &&
			config.Value.RatioCooldownEndTimestamp >= now+liquidity_provider.CooldownAfterRatioChange-1 &&
			config.Value.RatioCooldownEndTimestamp <= now+liquidity_provider.CooldownAfterRatioChange+1
	})).Return(nil)

	useCase := liquidity_provider.NewSetLiquidityRatioUseCase(provider, lpRepository, walletMock, crypto.Keccak256)
	err := useCase.Run(ctx, 70)

	require.NoError(t, err)
	provider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
	walletMock.AssertExpectations(t)
}

func TestSetLiquidityRatioUseCase_Run_NoOpSamePercentage(t *testing.T) {
	ctx := context.Background()
	provider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	walletMock := new(mocks.RskWalletMock)

	stateConfig := lpEntity.StateConfiguration{
		BtcLiquidityTargetPercentage: 50,
	}
	provider.On("StateConfiguration", ctx).Return(stateConfig, nil)

	useCase := liquidity_provider.NewSetLiquidityRatioUseCase(provider, lpRepository, walletMock, crypto.Keccak256)
	err := useCase.Run(ctx, 50)

	require.NoError(t, err)
	provider.AssertExpectations(t)
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
	walletMock.AssertNotCalled(t, "SignBytes", mock.Anything)
}

func TestSetLiquidityRatioUseCase_Run_StateConfigurationError(t *testing.T) {
	ctx := context.Background()
	provider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	walletMock := new(mocks.RskWalletMock)

	provider.On("StateConfiguration", ctx).Return(lpEntity.StateConfiguration{}, assert.AnError)

	useCase := liquidity_provider.NewSetLiquidityRatioUseCase(provider, lpRepository, walletMock, crypto.Keccak256)
	err := useCase.Run(ctx, 70)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SetLiquidityRatio")
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestSetLiquidityRatioUseCase_Run_SigningError(t *testing.T) {
	ctx := context.Background()
	provider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	walletMock := new(mocks.RskWalletMock)

	stateConfig := lpEntity.StateConfiguration{
		BtcLiquidityTargetPercentage: 50,
	}
	provider.On("StateConfiguration", ctx).Return(stateConfig, nil)
	walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError)

	useCase := liquidity_provider.NewSetLiquidityRatioUseCase(provider, lpRepository, walletMock, crypto.Keccak256)
	err := useCase.Run(ctx, 70)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SetLiquidityRatio")
	lpRepository.AssertNotCalled(t, "UpsertStateConfiguration", mock.Anything, mock.Anything)
}

func TestSetLiquidityRatioUseCase_Run_UpsertError(t *testing.T) {
	ctx := context.Background()
	provider := new(mocks.ProviderMock)
	lpRepository := new(mocks.LiquidityProviderRepositoryMock)
	walletMock := new(mocks.RskWalletMock)

	stateConfig := lpEntity.StateConfiguration{
		BtcLiquidityTargetPercentage: 50,
	}
	provider.On("StateConfiguration", ctx).Return(stateConfig, nil)
	walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
	lpRepository.On("UpsertStateConfiguration", ctx, mock.Anything).Return(assert.AnError)

	useCase := liquidity_provider.NewSetLiquidityRatioUseCase(provider, lpRepository, walletMock, crypto.Keccak256)
	err := useCase.Run(ctx, 70)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SetLiquidityRatio")
}
