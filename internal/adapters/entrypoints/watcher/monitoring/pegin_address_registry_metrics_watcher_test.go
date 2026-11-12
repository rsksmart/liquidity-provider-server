package monitoring_test

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type registryMetricsFixture struct {
	metrics        *monitoring.Metrics
	eventBus       *mocks.EventBusMock
	mismatchEvents chan entities.Event
	resyncEvents   chan entities.Event
	watcher        *monitoring.PegInAddressRegistryMetricsWatcher
}

func newRegistryMetricsFixture(t *testing.T) *registryMetricsFixture {
	t.Helper()
	fixture := &registryMetricsFixture{
		metrics:        monitoring.NewMetrics(prometheus.NewRegistry()),
		eventBus:       &mocks.EventBusMock{},
		mismatchEvents: make(chan entities.Event),
		resyncEvents:   make(chan entities.Event),
	}
	fixture.eventBus.Test(t)
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressRegistryRootMismatchEventId).
		Return((<-chan entities.Event)(fixture.mismatchEvents)).Once()
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressRegistryResyncEventId).
		Return((<-chan entities.Event)(fixture.resyncEvents)).Once()
	fixture.watcher = monitoring.NewPegInAddressRegistryMetricsWatcher(fixture.metrics, fixture.eventBus)
	require.NoError(t, fixture.watcher.Prepare(context.Background()))
	return fixture
}

func startRegistryMetricsWatcher(t *testing.T, metricsWatcher *monitoring.PegInAddressRegistryMetricsWatcher) (stop func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		metricsWatcher.Start()
		close(done)
	}()
	return func() {
		closeChannel := make(chan bool, 1)
		metricsWatcher.Shutdown(closeChannel)
		assert.True(t, <-closeChannel)
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("the metrics watcher did not stop")
		}
	}
}

func registryCounters(metrics *monitoring.Metrics) (mismatches, resyncs float64) {
	return getCounterValue(metrics.PegInAddressRegistryRootMismatchMetric), getCounterValue(metrics.PegInAddressRegistryResyncMetric)
}

func rootMismatchEvent() blockchain.PegInAddressRegistryRootMismatchEvent {
	return blockchain.PegInAddressRegistryRootMismatchEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRootMismatchEventId),
	}
}

func resyncEvent() blockchain.PegInAddressRegistryResyncEvent {
	return blockchain.PegInAddressRegistryResyncEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryResyncEventId),
	}
}

func TestPegInAddressRegistryMetricsWatcher_Prepare(t *testing.T) {
	fixture := newRegistryMetricsFixture(t)

	fixture.eventBus.AssertExpectations(t)
}

func TestPegInAddressRegistryMetricsWatcher_CountsARecoveryPublishedBeforeStart(t *testing.T) {
	eventBus := dataproviders.NewLocalEventBus()
	t.Cleanup(func() {
		closeChannel := make(chan bool, 1)
		eventBus.Shutdown(closeChannel)
		<-closeChannel
	})
	metrics := monitoring.NewMetrics(prometheus.NewRegistry())
	metricsWatcher := monitoring.NewPegInAddressRegistryMetricsWatcher(metrics, eventBus)
	require.NoError(t, metricsWatcher.Prepare(context.Background()))

	eventBus.Publish(rootMismatchEvent())
	eventBus.Publish(resyncEvent())
	stop := startRegistryMetricsWatcher(t, metricsWatcher)

	assert.Eventually(t, func() bool {
		mismatches, resyncs := registryCounters(metrics)
		return mismatches == 1 && resyncs == 1
	}, time.Second, time.Millisecond)
	stop()
}

func TestPegInAddressRegistryMetricsWatcher_CountsEachEvent(t *testing.T) {
	cases := []struct {
		name             string
		mismatchEvents   int
		resyncEvents     int
		expectedMismatch float64
		expectedResyncs  float64
	}{
		{name: "a recovery after a mismatch counts one mismatch and one resync", mismatchEvents: 1, resyncEvents: 1, expectedMismatch: 1, expectedResyncs: 1},
		{name: "a rebuild with no stored point counts only a resync", resyncEvents: 1, expectedResyncs: 1},
		{name: "two recoveries count twice", mismatchEvents: 2, resyncEvents: 2, expectedMismatch: 2, expectedResyncs: 2},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newRegistryMetricsFixture(t)
			stop := startRegistryMetricsWatcher(t, fixture.watcher)

			for range testCase.mismatchEvents {
				fixture.mismatchEvents <- rootMismatchEvent()
			}
			for range testCase.resyncEvents {
				fixture.resyncEvents <- resyncEvent()
			}
			stop()

			mismatches, resyncs := registryCounters(fixture.metrics)
			assert.InDelta(t, testCase.expectedMismatch, mismatches, 0.001)
			assert.InDelta(t, testCase.expectedResyncs, resyncs, 0.001)
		})
	}
}

func TestPegInAddressRegistryMetricsWatcher_CountsNothingFromAClosedSubscription(t *testing.T) {
	cases := []struct {
		name              string
		closeSubscription func(fixture *registryMetricsFixture)
		publishOther      func(fixture *registryMetricsFixture)
		expectedMismatch  float64
		expectedResyncs   float64
	}{
		{
			name:              "a closed mismatch subscription",
			closeSubscription: func(fixture *registryMetricsFixture) { close(fixture.mismatchEvents) },
			publishOther:      func(fixture *registryMetricsFixture) { fixture.resyncEvents <- resyncEvent() },
			expectedResyncs:   1,
		},
		{
			name:              "a closed resync subscription",
			closeSubscription: func(fixture *registryMetricsFixture) { close(fixture.resyncEvents) },
			publishOther:      func(fixture *registryMetricsFixture) { fixture.mismatchEvents <- rootMismatchEvent() },
			expectedMismatch:  1,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newRegistryMetricsFixture(t)
			stop := startRegistryMetricsWatcher(t, fixture.watcher)

			testCase.closeSubscription(fixture)
			testCase.publishOther(fixture)
			stop()

			mismatches, resyncs := registryCounters(fixture.metrics)
			assert.InDelta(t, testCase.expectedMismatch, mismatches, 0.001)
			assert.InDelta(t, testCase.expectedResyncs, resyncs, 0.001)
		})
	}
}
