package liquidity_provider_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckLiquidityUseCase_Run(t *testing.T) {
	bridge := &mocks.BridgeMock{}
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	contracts := blockchain.RskContracts{Bridge: bridge}
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, contracts, alertSender, "recipient")
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	provider.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert")
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoPeginLiquidity(t *testing.T) {
	bridge := &mocks.BridgeMock{}
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	recipient := "recipient@test.com"
	alertSender.On("SendAlert",
		test.AnyCtx,
		entities.AlertSubjectPeginOutOfLiquidity,
		"You are out of liquidity to perform a PegIn. Please, do a deposit",
		[]string{recipient},
	).Return(nil).Once()
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	contracts := blockchain.RskContracts{Bridge: bridge}
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, contracts, alertSender, recipient)
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	provider.AssertExpectations(t)
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoPegoutLiquidity(t *testing.T) {
	bridge := &mocks.BridgeMock{}
	provider := &mocks.ProviderMock{}
	alertSender := &mocks.AlertSenderMock{}
	recipient := "recipient@test.com"
	alertSender.On("SendAlert",
		test.AnyCtx,
		entities.AlertSubjectPegoutOutOfLiquidity,
		"You are out of liquidity to perform a PegOut. Please, do a deposit",
		[]string{recipient},
	).Return(nil).Once()
	provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
	provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
	contracts := blockchain.RskContracts{Bridge: bridge}
	useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, contracts, alertSender, recipient)
	err := useCase.Run(context.Background())
	bridge.AssertExpectations(t)
	provider.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	require.NoError(t, err)
}

func TestCheckLiquidityUseCase_Run_NoRecoverableErrorHandling(t *testing.T) {
	recipient := "anything"
	cases := test.Table[func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock), error]{
		{
			Value: func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(0), assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}
	for _, testCase := range cases {
		bridge := &mocks.BridgeMock{}
		provider := &mocks.ProviderMock{}
		sender := &mocks.AlertSenderMock{}
		testCase.Value(bridge, provider, sender)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, contracts, sender, recipient)
		err := useCase.Run(context.Background())
		bridge.AssertExpectations(t)
		provider.AssertExpectations(t)

		sender.AssertExpectations(t)
		require.Error(t, err)
	}
}

func TestCheckLiquidityUseCase_Run_OnlyLogSendErrors(t *testing.T) {
	recipient := "anything"
	cases := test.Table[func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock), error]{
		{
			Value: func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				sender.On("SendAlert", test.AnyCtx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			Value: func(bridge *mocks.BridgeMock, provider *mocks.ProviderMock, sender *mocks.AlertSenderMock) {
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
				provider.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil).Once()
				provider.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError).Once()
				sender.On("SendAlert", test.AnyCtx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}

	for _, testCase := range cases {
		bridge := &mocks.BridgeMock{}
		provider := &mocks.ProviderMock{}
		sender := &mocks.AlertSenderMock{}
		buff := new(bytes.Buffer)
		testCase.Value(bridge, provider, sender)
		log.SetOutput(buff)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := liquidity_provider.NewCheckLiquidityUseCase(provider, provider, contracts, sender, recipient)
		err := useCase.Run(context.Background())
		assert.Positive(t, buff.Bytes())
		bridge.AssertExpectations(t)
		provider.AssertExpectations(t)
		sender.AssertExpectations(t)
		require.NoError(t, err)
	}
}
