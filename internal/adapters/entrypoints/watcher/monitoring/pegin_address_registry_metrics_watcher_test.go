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
	metrics             *monitoring.Metrics
	eventBus            *mocks.EventBusMock
	mismatchEvents      chan entities.Event
	rebuildEvents       chan entities.Event
	rebuildFailedEvents chan entities.Event
	watcher             *monitoring.PegInAddressRegistryMetricsWatcher
}

type registryCounterValues struct {
	mismatches    float64
	rebuilds      float64
	rebuildErrors map[blockchain.PegInAddressRegistryRebuildReason]float64
}

var registryRebuildReasons = []blockchain.PegInAddressRegistryRebuildReason{
	blockchain.PegInAddressRegistryRebuildRootMismatch,
	blockchain.PegInAddressRegistryRebuildFirstStart,
	blockchain.PegInAddressRegistryRebuildStartBlockRaised,
}

func newRegistryMetricsFixture(t *testing.T) *registryMetricsFixture {
	t.Helper()
	fixture := &registryMetricsFixture{
		metrics:             monitoring.NewMetrics(prometheus.NewRegistry()),
		eventBus:            &mocks.EventBusMock{},
		mismatchEvents:      make(chan entities.Event),
		rebuildEvents:       make(chan entities.Event),
		rebuildFailedEvents: make(chan entities.Event),
	}
	fixture.eventBus.Test(t)
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressRegistryRootMismatchEventId).
		Return((<-chan entities.Event)(fixture.mismatchEvents)).Once()
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressRegistryRebuildEventId).
		Return((<-chan entities.Event)(fixture.rebuildEvents)).Once()
	fixture.eventBus.On("Subscribe", blockchain.PegInAddressRegistryRebuildFailedEventId).
		Return((<-chan entities.Event)(fixture.rebuildFailedEvents)).Once()
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

func readRegistryCounters(metrics *monitoring.Metrics) registryCounterValues {
	values := registryCounterValues{
		mismatches:    getCounterValue(metrics.PegInAddressRegistryRootMismatchMetric),
		rebuilds:      getCounterValue(metrics.PegInAddressRegistryRebuildMetric),
		rebuildErrors: map[blockchain.PegInAddressRegistryRebuildReason]float64{},
	}
	for _, reason := range registryRebuildReasons {
		values.rebuildErrors[reason] = getCounterVecValue(metrics.PegInAddressRegistryRebuildErrorsMetric, string(reason))
	}
	return values
}

func assertRegistryCounters(t *testing.T, expected registryCounterValues, metrics *monitoring.Metrics) {
	t.Helper()
	actual := readRegistryCounters(metrics)
	assert.InDelta(t, expected.mismatches, actual.mismatches, 0.001)
	assert.InDelta(t, expected.rebuilds, actual.rebuilds, 0.001)
	for _, reason := range registryRebuildReasons {
		assert.InDelta(t, expected.rebuildErrors[reason], actual.rebuildErrors[reason], 0.001, string(reason))
	}
}

func rootMismatchEvent() blockchain.PegInAddressRegistryRootMismatchEvent {
	return blockchain.PegInAddressRegistryRootMismatchEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRootMismatchEventId),
	}
}

func rebuildEvent() blockchain.PegInAddressRegistryRebuildEvent {
	return blockchain.PegInAddressRegistryRebuildEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRebuildEventId),
	}
}

func rebuildFailedEvent(reason blockchain.PegInAddressRegistryRebuildReason) blockchain.PegInAddressRegistryRebuildFailedEvent {
	return blockchain.PegInAddressRegistryRebuildFailedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRebuildFailedEventId),
		Reason:    reason,
	}
}

func TestPegInAddressRegistryMetricsWatcher_Prepare(t *testing.T) {
	fixture := newRegistryMetricsFixture(t)

	fixture.eventBus.AssertExpectations(t)
}

func TestPegInAddressRegistryMetricsWatcher_CountsEventsPublishedBeforeStart(t *testing.T) {
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
	eventBus.Publish(rebuildEvent())
	eventBus.Publish(rebuildFailedEvent(blockchain.PegInAddressRegistryRebuildFirstStart))
	stop := startRegistryMetricsWatcher(t, metricsWatcher)

	assert.Eventually(t, func() bool {
		counters := readRegistryCounters(metrics)
		return counters.mismatches == 1 && counters.rebuilds == 1 &&
			counters.rebuildErrors[blockchain.PegInAddressRegistryRebuildFirstStart] == 1
	}, time.Second, time.Millisecond)
	stop()
}

func TestPegInAddressRegistryMetricsWatcher_CountsEachEvent(t *testing.T) {
	cases := []struct {
		name           string
		mismatchEvents int
		rebuildEvents  int
		expected       registryCounterValues
	}{
		{
			name:           "a recovery after a mismatch counts one mismatch and one rebuild",
			mismatchEvents: 1,
			rebuildEvents:  1,
			expected:       registryCounterValues{mismatches: 1, rebuilds: 1},
		},
		{
			name:          "a rebuild with no stored checkpoint counts only a rebuild",
			rebuildEvents: 1,
			expected:      registryCounterValues{rebuilds: 1},
		},
		{
			name:           "two recoveries count twice",
			mismatchEvents: 2,
			rebuildEvents:  2,
			expected:       registryCounterValues{mismatches: 2, rebuilds: 2},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newRegistryMetricsFixture(t)
			stop := startRegistryMetricsWatcher(t, fixture.watcher)

			for range testCase.mismatchEvents {
				fixture.mismatchEvents <- rootMismatchEvent()
			}
			for range testCase.rebuildEvents {
				fixture.rebuildEvents <- rebuildEvent()
			}
			stop()

			assertRegistryCounters(t, testCase.expected, fixture.metrics)
		})
	}
}

func TestPegInAddressRegistryMetricsWatcher_CountsFailedRebuildsPerReason(t *testing.T) {
	fixture := newRegistryMetricsFixture(t)
	stop := startRegistryMetricsWatcher(t, fixture.watcher)

	for _, reason := range []blockchain.PegInAddressRegistryRebuildReason{
		blockchain.PegInAddressRegistryRebuildRootMismatch,
		blockchain.PegInAddressRegistryRebuildFirstStart,
		blockchain.PegInAddressRegistryRebuildStartBlockRaised,
		blockchain.PegInAddressRegistryRebuildRootMismatch,
	} {
		fixture.rebuildFailedEvents <- rebuildFailedEvent(reason)
	}
	stop()

	assertRegistryCounters(t, registryCounterValues{rebuildErrors: map[blockchain.PegInAddressRegistryRebuildReason]float64{
		blockchain.PegInAddressRegistryRebuildRootMismatch:     2,
		blockchain.PegInAddressRegistryRebuildFirstStart:       1,
		blockchain.PegInAddressRegistryRebuildStartBlockRaised: 1,
	}}, fixture.metrics)
}

func TestPegInAddressRegistryMetricsWatcher_CountsNothingFromAClosedSubscription(t *testing.T) {
	cases := []struct {
		name              string
		closeSubscription func(fixture *registryMetricsFixture)
		publishOther      func(fixture *registryMetricsFixture)
		expected          registryCounterValues
	}{
		{
			name:              "a closed mismatch subscription",
			closeSubscription: func(fixture *registryMetricsFixture) { close(fixture.mismatchEvents) },
			publishOther:      func(fixture *registryMetricsFixture) { fixture.rebuildEvents <- rebuildEvent() },
			expected:          registryCounterValues{rebuilds: 1},
		},
		{
			name:              "a closed rebuild subscription",
			closeSubscription: func(fixture *registryMetricsFixture) { close(fixture.rebuildEvents) },
			publishOther: func(fixture *registryMetricsFixture) {
				fixture.rebuildFailedEvents <- rebuildFailedEvent(blockchain.PegInAddressRegistryRebuildRootMismatch)
			},
			expected: registryCounterValues{rebuildErrors: map[blockchain.PegInAddressRegistryRebuildReason]float64{
				blockchain.PegInAddressRegistryRebuildRootMismatch: 1,
			}},
		},
		{
			name:              "a closed failed rebuild subscription",
			closeSubscription: func(fixture *registryMetricsFixture) { close(fixture.rebuildFailedEvents) },
			publishOther:      func(fixture *registryMetricsFixture) { fixture.mismatchEvents <- rootMismatchEvent() },
			expected:          registryCounterValues{mismatches: 1},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newRegistryMetricsFixture(t)
			stop := startRegistryMetricsWatcher(t, fixture.watcher)

			testCase.closeSubscription(fixture)
			testCase.publishOther(fixture)
			stop()

			assertRegistryCounters(t, testCase.expected, fixture.metrics)
		})
	}
}
