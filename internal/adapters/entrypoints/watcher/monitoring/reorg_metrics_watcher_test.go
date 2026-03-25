package monitoring_test

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
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
		d := prometheusMetricGaugeValue(t, reg, "lps_node_reorg_depth", string(entities.NodeTypeBitcoin))
		return d == 3
	}, time.Second, 10*time.Millisecond)
	assert.InDelta(t, 2.0, prometheusMetricGaugeValue(t, reg, "lps_node_reorg_max_depth_threshold", string(entities.NodeTypeBitcoin)), 0.001)
	assert.InDelta(t, 1.0, prometheusMetricGaugeValue(t, reg, "lps_node_reorg_above_threshold", string(entities.NodeTypeBitcoin)), 0.001)
	assert.InDelta(t, 1.0, prometheusMetricCounterValue(t, reg, "lps_node_reorg_check_errors_total", string(entities.NodeTypeRootstock)), 0.001)
	assert.InDelta(t, 1.0, prometheusMetricCounterValue(t, reg, "lps_node_reorg_alerts_total", string(entities.NodeTypeBitcoin)), 0.001)
	closeDone := make(chan bool, 1)
	go w.Shutdown(closeDone)
	<-closeDone
	bus.AssertExpectations(t)
}

func prometheusMetricGaugeValue(t *testing.T, g prometheus.Gatherer, name, nodeLabel string) float64 {
	t.Helper()
	mfs, err := g.Gather()
	require.NoError(t, err)
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, metric := range mf.GetMetric() {
			if !labelHasNode(metric.GetLabel(), nodeLabel) || metric.GetGauge() == nil {
				continue
			}
			return metric.GetGauge().GetValue()
		}
	}
	t.Fatalf("gauge %s{node=%q} not found", name, nodeLabel)
	return 0
}

func prometheusMetricCounterValue(t *testing.T, g prometheus.Gatherer, name, nodeLabel string) float64 {
	t.Helper()
	mfs, err := g.Gather()
	require.NoError(t, err)
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, metric := range mf.GetMetric() {
			if !labelHasNode(metric.GetLabel(), nodeLabel) || metric.GetCounter() == nil {
				continue
			}
			return metric.GetCounter().GetValue()
		}
	}
	t.Fatalf("counter %s{node=%q} not found", name, nodeLabel)
	return 0
}

func labelHasNode(labels []*dto.LabelPair, node string) bool {
	for _, lp := range labels {
		if lp.GetName() == "node" && lp.GetValue() == node {
			return true
		}
	}
	return false
}
