package monitoring_test

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeerMetricsWatcher_Prepare(t *testing.T) {
	watcher := monitoring.NewPeerMetricsWatcher(monitoring.NewMetrics(prometheus.NewRegistry()), &mocks.EventBusMock{})
	require.NoError(t, watcher.Prepare(context.Background()))
}

func TestPeerMetricsWatcher_Start(t *testing.T) {
	appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
	eventBus := &mocks.EventBusMock{}

	checkChannel := make(chan entities.Event, 1)
	errorChannel := make(chan entities.Event, 1)
	alertChannel := make(chan entities.Event, 1)
	eventBus.On("Subscribe", blockchain.NodePeerCheckEventId).Return((<-chan entities.Event)(checkChannel))
	eventBus.On("Subscribe", blockchain.NodePeerCheckErrorEventId).Return((<-chan entities.Event)(errorChannel))
	eventBus.On("Subscribe", blockchain.NodePeerAlertSentEventId).Return((<-chan entities.Event)(alertChannel))

	watcher := monitoring.NewPeerMetricsWatcher(appMetrics, eventBus)
	go watcher.Start()

	checkChannel <- blockchain.NodePeerCheckEvent{
		BaseEvent:      entities.NewBaseEvent(blockchain.NodePeerCheckEventId),
		NodeType:       entities.NodeTypeBitcoin,
		CurrentPeers:   2,
		MinPeers:       3,
		BelowThreshold: true,
	}
	errorChannel <- blockchain.NodePeerCheckErrorEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.NodePeerCheckErrorEventId),
		NodeType:  entities.NodeTypeBitcoin,
	}
	alertChannel <- blockchain.NodePeerAlertSentEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.NodePeerAlertSentEventId),
		NodeType:  entities.NodeTypeBitcoin,
	}

	assert.Eventually(t, func() bool {
		return getGaugeVecValue(appMetrics.NodePeerCountMetric, "bitcoin") == 2 &&
			getGaugeVecValue(appMetrics.NodePeerMinThresholdMetric, "bitcoin") == 3 &&
			getGaugeVecValue(appMetrics.NodePeerBelowThreshold, "bitcoin") == 1 &&
			getCounterVecValue(appMetrics.NodePeerCheckErrors, "bitcoin") == 1 &&
			getCounterVecValue(appMetrics.NodePeerAlerts, "bitcoin") == 1
	}, time.Second, 20*time.Millisecond)

	shutdown := make(chan bool, 1)
	watcher.Shutdown(shutdown)
	<-shutdown
	eventBus.AssertExpectations(t)
}

func TestPeerMetricsWatcher_Shutdown(t *testing.T) {
	appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", blockchain.NodePeerCheckEventId).Return(make(<-chan entities.Event))
	eventBus.On("Subscribe", blockchain.NodePeerCheckErrorEventId).Return(make(<-chan entities.Event))
	eventBus.On("Subscribe", blockchain.NodePeerAlertSentEventId).Return(make(<-chan entities.Event))

	watcher := monitoring.NewPeerMetricsWatcher(appMetrics, eventBus)
	go watcher.Start()
	time.Sleep(10 * time.Millisecond)

	closeChannel := make(chan bool, 1)
	watcher.Shutdown(closeChannel)
	select {
	case <-closeChannel:
	case <-time.After(time.Second):
		t.Fatal("Shutdown timed out")
	}
	eventBus.AssertExpectations(t)
}
