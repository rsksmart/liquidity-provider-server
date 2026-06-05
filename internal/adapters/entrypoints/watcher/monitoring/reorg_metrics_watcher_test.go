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

func TestReorgMetricsWatcher_Prepare(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := monitoring.NewMetrics(reg)
	bus := &mocks.EventBusMock{}
	w := monitoring.NewReorgMetricsWatcher(m, bus)
	require.NoError(t, w.Prepare(context.Background()))
}

func TestReorgMetricsWatcher_StartUpdatesMetricsFromEvents(t *testing.T) {
	reg := prometheus.NewRegistry()
	m := monitoring.NewMetrics(reg)
	checkCh := make(chan entities.Event, 4)
	errCh := make(chan entities.Event, 4)
	alertCh := make(chan entities.Event, 4)
	bus := &mocks.EventBusMock{}
	bus.On("Subscribe", blockchain.NodeReorgCheckEventId).Return((<-chan entities.Event)(checkCh))
	bus.On("Subscribe", blockchain.NodeReorgCheckErrorEventId).Return((<-chan entities.Event)(errCh))
	bus.On("Subscribe", blockchain.NodeReorgAlertSentEventId).Return((<-chan entities.Event)(alertCh))
	w := monitoring.NewReorgMetricsWatcher(m, bus)
	go w.Start()
	checkCh <- blockchain.NodeReorgCheckEvent{
		BaseEvent:       entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    3,
		MaxAllowedDepth: 2,
		AboveThreshold:  true,
	}
	errCh <- blockchain.NodeReorgCheckErrorEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.NodeReorgCheckErrorEventId),
		NodeType:  entities.NodeTypeRootstock,
	}
	alertCh <- blockchain.NodeReorgAlertSentEvent{
		BaseEvent:     entities.NewBaseEvent(blockchain.NodeReorgAlertSentEventId),
		NodeType:      entities.NodeTypeBitcoin,
		DetectedDepth: 4,
	}
	assert.Eventually(t, func() bool {
		return getGaugeVecValue(m.NodeReorgDepthMetric, string(entities.NodeTypeBitcoin)) == 3
	}, time.Second, 10*time.Millisecond)
	assert.InDelta(t, 2.0, getGaugeVecValue(m.NodeReorgMaxDepthMetric, string(entities.NodeTypeBitcoin)), 0.001)
	assert.InDelta(t, 1.0, getGaugeVecValue(m.NodeReorgAboveThresholdMetric, string(entities.NodeTypeBitcoin)), 0.001)
	assert.InDelta(t, 1.0, getCounterVecValue(m.NodeReorgCheckErrorsMetric, string(entities.NodeTypeRootstock)), 0.001)
	assert.InDelta(t, 1.0, getCounterVecValue(m.NodeReorgAlertsMetric, string(entities.NodeTypeBitcoin)), 0.001)
	closeDone := make(chan bool, 1)
	go w.Shutdown(closeDone)
	<-closeDone
	bus.AssertExpectations(t)
}
