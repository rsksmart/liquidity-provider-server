package liquidity_provider_test

import (
	"bytes"
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckLiquidityUseCase_Run(t *testing.T) {
	bridge := &test.BridgeMock{}
	provider := &test.ProviderMock{}
	alertSender := &test.AlertSenderMock{}
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, bridge, alertSender, "recipient")
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoPeginLiquidity(t *testing.T) {
	bridge := &test.BridgeMock{}
	provider := &test.ProviderMock{}
	alertSender := &test.AlertSenderMock{}
	recipient := "recipient@test.com"
	alertSender.On("SendAlert",
		mock.AnythingOfType("context.backgroundCtx"),
		"PegIn: Out of liquidity",
		"You are out of liquidity to perform a PegIn. Please, do a deposit",
		[]string{recipient},
	).Return(nil).Once()
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, bridge, alertSender, recipient)
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	provider.AssertExpectations(t)
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoPegoutLiquidity(t *testing.T) {
	bridge := &test.BridgeMock{}
	provider := &test.ProviderMock{}
	alertSender := &test.AlertSenderMock{}
	recipient := "recipient@test.com"
	alertSender.On("SendAlert",
		mock.AnythingOfType("context.backgroundCtx"),
		"PegOut: Out of liquidity",
		"You are out of liquidity to perform a PegOut. Please, do a deposit",
		[]string{recipient},
	).Return(nil).Once()
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, bridge, alertSender, recipient)
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoRecoverableErrorHandling(t *testing.T) {
	recipient := "anything"
	cases := test.Table[func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock), error]{
		{
			Value: func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(0), assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}
	for _, testCase := range cases {
		bridge := &test.BridgeMock{}
		provider := &test.ProviderMock{}
		sender := &test.AlertSenderMock{}
		testCase.Value(bridge, provider, sender)
		useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, bridge, sender, recipient)
		err := useCase.Run(context.Background())
		bridge.AssertExpectations(t)
		provider.AssertExpectations(t)

		sender.AssertExpectations(t)
		require.Error(t, err)
	}
}

func TestCheckLiquidityUseCase_Run_OnlyLogSendErrors(t *testing.T) {
	recipient := "anything"
	cases := test.Table[func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock), error]{
		{
			Value: func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				sender.On("SendAlert", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *test.BridgeMock, provider *test.ProviderMock, sender *test.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
				sender.On("SendAlert", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}

	for _, testCase := range cases {
		bridge := &test.BridgeMock{}
		provider := &test.ProviderMock{}
		sender := &test.AlertSenderMock{}
		buff := new(bytes.Buffer)
		testCase.Value(bridge, provider, sender)
		log.SetOutput(buff)
		useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, bridge, sender, recipient)
		err := useCase.Run(context.Background())
		assert.Positive(t, buff.Bytes())
		bridge.AssertExpectations(t)
		provider.AssertExpectations(t)
		sender.AssertExpectations(t)
		require.NoError(t, err)
	}
}
