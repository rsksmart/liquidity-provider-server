package watcher_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNodePeerCheckUseCase_Run(t *testing.T) {
	t.Run("btc: should publish check event and not alert above threshold", testNodePeerCheckNoAlertAboveThreshold)
	t.Run("btc: should send alert and publish alert event below threshold", testNodePeerCheckSendAlertBelowThreshold)
	t.Run("btc: should publish error event when rpc fails", testNodePeerCheckPublishErrorOnRPCFailure)
	t.Run("btc: should suppress alert during cooldown", testNodePeerCheckCooldownSuppressesAlert)
	t.Run("btc: should return error when send alert fails", testNodePeerCheckSendAlertFailure)
	t.Run("rsk: should publish check event and not alert above threshold", testNodePeerCheckRskNoAlertAboveThreshold)
	t.Run("rsk: should send alert and publish alert event below threshold", testNodePeerCheckRskSendAlertBelowThreshold)
	t.Run("rsk: should publish error event when rpc fails", testNodePeerCheckRskPublishErrorOnRPCFailure)
	t.Run("rsk: should suppress alert during cooldown", testNodePeerCheckRskCooldownSuppressesAlert)
	t.Run("rsk: should return error when send alert fails", testNodePeerCheckRskSendAlertFailure)
}

func testNodePeerCheckNoAlertAboveThreshold(t *testing.T) {
	const recipient = "alert@example.com"
	btcRpc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	btcRpc.On("GetConnectionCount").Return(int64(5), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin && !ev.BelowThreshold && ev.CurrentPeers == 5 && ev.MinPeers == 3
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Btc: btcRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	btcRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func testNodePeerCheckSendAlertBelowThreshold(t *testing.T) {
	const recipient = "alert@example.com"
	btcRpc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	btcRpc.On("GetConnectionCount").Return(int64(1), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin && ev.BelowThreshold && ev.CurrentPeers == 1 && ev.MinPeers == 3
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeBitcoin, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Btc: btcRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	btcRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertExpectations(t)
}

func testNodePeerCheckPublishErrorOnRPCFailure(t *testing.T) {
	const recipient = "alert@example.com"
	btcRpc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	btcRpc.On("GetConnectionCount").Return(int64(0), assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckErrorEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Btc: btcRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	eventBus.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func testNodePeerCheckCooldownSuppressesAlert(t *testing.T) {
	const recipient = "alert@example.com"
	btcRpc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	btcRpc.On("GetConnectionCount").Return(int64(1), nil).Twice()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin
	})).Return().Twice()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeBitcoin, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Btc: btcRpc}, alertSender, recipient, eventBus, 3, 3, time.Hour)
	err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	err = useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	alertSender.AssertNumberOfCalls(t, "SendAlert", 1)
}

func testNodePeerCheckSendAlertFailure(t *testing.T) {
	const recipient = "alert@example.com"
	btcRpc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	btcRpc.On("GetConnectionCount").Return(int64(1), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin && ev.BelowThreshold && ev.CurrentPeers == 1 && ev.MinPeers == 3
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeBitcoin, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(assert.AnError).Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Btc: btcRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeBitcoin)
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	btcRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok
	}))
}

func testNodePeerCheckRskSendAlertFailure(t *testing.T) {
	const recipient = "alert@example.com"
	rskRpc := &mocks.RootstockRpcServerMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	rskRpc.On("PeerCount", mock.Anything).Return(uint64(1), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock && ev.BelowThreshold && ev.CurrentPeers == 1 && ev.MinPeers == 3
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeRootstock, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(assert.AnError).Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Rsk: rskRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	rskRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	eventBus.AssertNotCalled(t, "Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok
	}))
}

func testNodePeerCheckRskNoAlertAboveThreshold(t *testing.T) {
	const recipient = "alert@example.com"
	rskRpc := &mocks.RootstockRpcServerMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	rskRpc.On("PeerCount", mock.Anything).Return(uint64(5), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock && !ev.BelowThreshold && ev.CurrentPeers == 5 && ev.MinPeers == 3
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Rsk: rskRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	rskRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func testNodePeerCheckRskSendAlertBelowThreshold(t *testing.T) {
	const recipient = "alert@example.com"
	rskRpc := &mocks.RootstockRpcServerMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	rskRpc.On("PeerCount", mock.Anything).Return(uint64(1), nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock && ev.BelowThreshold && ev.CurrentPeers == 1 && ev.MinPeers == 3
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeRootstock, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Rsk: rskRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	rskRpc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	alertSender.AssertExpectations(t)
}

func testNodePeerCheckRskPublishErrorOnRPCFailure(t *testing.T) {
	const recipient = "alert@example.com"
	rskRpc := &mocks.RootstockRpcServerMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	rskRpc.On("PeerCount", mock.Anything).Return(uint64(0), assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckErrorEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock
	})).Return().Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Rsk: rskRpc}, alertSender, recipient, eventBus, 3, 3, 30*time.Minute)
	err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
	eventBus.AssertExpectations(t)
	alertSender.AssertNotCalled(t, "SendAlert", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func testNodePeerCheckRskCooldownSuppressesAlert(t *testing.T) {
	const recipient = "alert@example.com"
	rskRpc := &mocks.RootstockRpcServerMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	rskRpc.On("PeerCount", mock.Anything).Return(uint64(1), nil).Twice()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock
	})).Return().Twice()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodePeerAlertSentEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock
	})).Return().Once()
	expectedBody := fmt.Sprintf(watcher.LowPeersAlertBodyTemplate, entities.NodeTypeRootstock, 1, 3)
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectLowPeers, expectedBody, []string{recipient}).Return(nil).Once()
	useCase := watcher.NewNodePeerCheckUseCase(blockchain.Rpc{Rsk: rskRpc}, alertSender, recipient, eventBus, 3, 3, time.Hour)
	err := useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	err = useCase.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	alertSender.AssertNumberOfCalls(t, "SendAlert", 1)
}
