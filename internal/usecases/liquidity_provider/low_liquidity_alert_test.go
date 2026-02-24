package liquidity_provider_test

import (
	"context"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	warningThreshold  uint64 = 3
	criticalThreshold uint64 = 1

	aboveWarningCoins uint64 = 5
	atWarningCoins    uint64 = 3
	belowWarningCoins uint64 = 2
	zeroCoins         int64  = 0

	alertRecipientEmail = "recipient@test.com"
)

func TestLowLiquidityAlertUseCase_Run_NoAlertAboveWarning(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_NoAlertAtExactWarning(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(atWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(atWarningCoins), nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_WarningAlertBtc(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(belowWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityWarning,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: BTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_WarningAlertRbtc(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(belowWarningCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityWarning,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: RBTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_CriticalAlertBtc(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityCritical,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: BTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_CriticalAlertRbtc(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityCritical,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: RBTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_BothNetworksBelowCritical(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityCritical,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: BTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityCritical,
		mock.MatchedBy(func(body string) bool {
			return strings.Contains(body, "Network: RBTC")
		}),
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_OnlyCriticalWhenBelowCritical(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	alertSender.On("SendAlert",
		test.AnyCtx,
		alerts.AlertSubjectHotWalletLowLiquidityCritical,
		mock.Anything,
		[]string{alertRecipientEmail},
	).Return(nil).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert",
		mock.Anything,
		alerts.AlertSubjectHotWalletLowLiquidityWarning,
		mock.Anything,
		mock.Anything,
	)
	require.NoError(t, err)
}

func TestLowLiquidityAlertUseCase_Run_ErrorFromPegoutProvider(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), assert.AnError).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.Error(t, err)
}

func TestLowLiquidityAlertUseCase_Run_ErrorFromPeginProvider(t *testing.T) {
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("AvailablePegoutLiquidity", mock.Anything).Return(entities.CoinToWei(aboveWarningCoins), nil).Once()
	provider.On("AvailablePeginLiquidity", mock.Anything).Return(entities.NewWei(zeroCoins), assert.AnError).Once()
	useCase := liquidity_provider.NewLowLiquidityAlertUseCase(provider, provider, alertSender, alertRecipientEmail, warningThreshold, criticalThreshold)
	err := useCase.Run(context.Background())
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.Error(t, err)
}
