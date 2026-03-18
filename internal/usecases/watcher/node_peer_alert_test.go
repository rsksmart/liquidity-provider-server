package watcher_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNodePeerAlertUseCase_Run(t *testing.T) {
	const recipient = "alert@example.com"
	t.Run("should send alert successfully", func(t *testing.T) {
		alertSender := &mocks.AlertSenderMock{}
		expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeBitcoin, 1, 3)
		alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
		useCase := watcher.NewNodePeerAlertUseCase(alertSender, recipient)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin, 1, 3)
		require.NoError(t, err)
		alertSender.AssertExpectations(t)
	})
	t.Run("should propagate sender error", func(t *testing.T) {
		alertSender := &mocks.AlertSenderMock{}
		alertSender.On("SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
		useCase := watcher.NewNodePeerAlertUseCase(alertSender, recipient)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin, 1, 3)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		alertSender.AssertExpectations(t)
	})
	t.Run("should include node type and peer counts in alert body", func(t *testing.T) {
		alertSender := &mocks.AlertSenderMock{}
		expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeBitcoin, 2, 5)
		alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
		useCase := watcher.NewNodePeerAlertUseCase(alertSender, recipient)
		err := useCase.Run(context.Background(), entities.NodeTypeBitcoin, 2, 5)
		require.NoError(t, err)
		assert.Contains(t, expectedBody, "bitcoin")
		assert.Contains(t, expectedBody, "2 peers")
		assert.Contains(t, expectedBody, "5")
		alertSender.AssertExpectations(t)
	})
}
