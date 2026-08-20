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

func TestPegInAddressRegistryMetricsWatcher_UpdatesIntegrityCounters(t *testing.T) {
	appMetrics := monitoring.NewMetrics(prometheus.NewRegistry())
	mismatchEvents := make(chan entities.Event, 1)
	resyncEvents := make(chan entities.Event, 1)
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", blockchain.PegInAddressRegistryRootMismatchEventId).
		Return((<-chan entities.Event)(mismatchEvents)).
		Once()
	eventBus.On("Subscribe", blockchain.PegInAddressRegistryResyncStartedEventId).
		Return((<-chan entities.Event)(resyncEvents)).
		Once()
	metricsWatcher := monitoring.NewPegInAddressRegistryMetricsWatcher(appMetrics, eventBus)
	require.NoError(t, metricsWatcher.Prepare(context.Background()))
	go metricsWatcher.Start()

	mismatchEvents <- blockchain.PegInAddressRegistryRootMismatchEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRootMismatchEventId),
	}
	resyncEvents <- blockchain.PegInAddressRegistryResyncStartedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryResyncStartedEventId),
	}

	assert.Eventually(t, func() bool {
		return counterValue(appMetrics.PegInAddressRegistryRootMismatchMetric) == 1 &&
			counterValue(appMetrics.PegInAddressRegistryResyncMetric) == 1
	}, time.Second, time.Millisecond)
	closeDone := make(chan bool, 1)
	metricsWatcher.Shutdown(closeDone)
	<-closeDone
	eventBus.AssertExpectations(t)
}

func counterValue(counter prometheus.Counter) float64 {
	metric := &dto.Metric{}
	if err := counter.Write(metric); err != nil {
		return 0
	}
	return metric.GetCounter().GetValue()
}
