package monitoring_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewQuoteMetricsWatcher(t *testing.T) {
	t.Run("should create metric watcher with all dependencies properly configured", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		require.NotNil(t, watcher)

		// Use reflection to verify internal structure (type checking only)
		watcherValue := reflect.ValueOf(watcher).Elem()
		watcherType := watcherValue.Type()

		// Verify all expected fields exist with correct types
		appMetricsField, found := watcherType.FieldByName("appMetrics")
		require.True(t, found, "appMetrics field should exist")
		assert.Equal(t, "*monitoring.Metrics", appMetricsField.Type.String())

		eventBusField, found := watcherType.FieldByName("eventBus")
		require.True(t, found, "eventBus field should exist")
		assert.Equal(t, "entities.EventBus", eventBusField.Type.String())

		serverInfoField, found := watcherType.FieldByName("serverInfo")
		require.True(t, found, "serverInfo field should exist")
		assert.Equal(t, "*liquidity_provider.ServerInfoUseCase", serverInfoField.Type.String())

		closeChannelField, found := watcherType.FieldByName("closeChannel")
		require.True(t, found, "closeChannel field should exist")
		assert.Equal(t, "chan struct {}", closeChannelField.Type.String())

		// Verify closeChannel is properly initialized by checking it's not zero value
		closeChannelValue := watcherValue.FieldByName("closeChannel")
		assert.False(t, closeChannelValue.IsZero(), "closeChannel should be initialized")
		assert.Equal(t, reflect.Chan, closeChannelValue.Kind())
		assert.Equal(t, 1, closeChannelValue.Cap(), "closeChannel should have capacity of 1")

		// Test behavioral verification: the watcher should be ready to call Prepare
		// This indirectly tests that serverInfo was properly assigned
		setServerInfoBuildVars("test", "test")
		err := watcher.Prepare(context.Background())
		assert.NoError(t, err, "Prepare should work if serverInfo is properly set")
	})
}

func TestQuoteMetricsWatcher_Prepare(t *testing.T) {
	t.Run("should set server info metrics successfully", func(t *testing.T) {
		setServerInfoBuildVars("1.0.0", "abc123")

		appMetrics := &monitoring.Metrics{
			ServerInfoMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_server"}, []string{"version", "commit"}),
		}
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)
		err := watcher.Prepare(context.Background())

		require.NoError(t, err)

		// Verify the metric was set with the correct values (success path)
		metricValue := getGaugeVecValue(appMetrics.ServerInfoMetric, "1.0.0", "abc123")
		assert.InEpsilon(t, float64(1), metricValue, 0.0)

		// Verify the labels are set correctly
		labels := getGaugeVecLabels(appMetrics.ServerInfoMetric, "1.0.0", "abc123")
		assert.Equal(t, "1.0.0", labels["version"], "Version label should be set correctly")
		assert.Equal(t, "abc123", labels["commit"], "Commit label should be set correctly")
	})

	t.Run("should handle server info error and set default values", func(t *testing.T) {
		setServerInfoBuildVars("", "")

		appMetrics := &monitoring.Metrics{
			ServerInfoMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_server"}, []string{"version", "commit"}),
		}
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)
		err := watcher.Prepare(context.Background())

		require.NoError(t, err)

		// Verify the metric was set with default values (error path)
		// This confirms that the "if err != nil" block was executed
		metricValue := getGaugeVecValue(appMetrics.ServerInfoMetric, "Not provided", "Not provided")
		assert.InEpsilon(t, float64(1), metricValue, 0.0)

		// Verify the labels are set to default values
		labels := getGaugeVecLabels(appMetrics.ServerInfoMetric, "Not provided", "Not provided")
		assert.Equal(t, "Not provided", labels["version"], "Version label should be set to default value")
		assert.Equal(t, "Not provided", labels["commit"], "Commit label should be set to default value")
	})
}

// nolint:funlen
func TestQuoteMetricsWatcher_Start_PegoutMetrics(t *testing.T) {
	t.Run("should increment accepted pegout quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		acceptedPegoutChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return((<-chan entities.Event)(acceptedPegoutChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		// Verify counter starts at 0
		initialValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "accepted")
		assert.Equal(t, 0, int(initialValue))

		// This pattern of Goroutine keeps the infinite loop of Start() method running but
		// without blocking the test execution
		go func() {
			watcher.Start()
		}()

		// Send events to trigger metric increment
		acceptedPegoutChannel <- entities.BaseEvent{}
		acceptedPegoutChannel <- entities.BaseEvent{}
		acceptedPegoutChannel <- entities.BaseEvent{}

		// Give some time for the event to be processed
		time.Sleep(10 * time.Millisecond)

		// Verify counter was incremented
		finalValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "accepted")
		assert.Equal(t, 3, int(finalValue))

		// Verify the label is correct
		labels := getCounterVecLabels(appMetrics.PegoutQuotesMetric, "accepted")
		assert.Equal(t, "accepted", labels["state"], "State label should be 'accepted'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})

	t.Run("should increment sent pegout quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		sendPegoutChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return((<-chan entities.Event)(sendPegoutChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		initialValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "paid")
		assert.Equal(t, 0, int(initialValue))

		go func() {
			watcher.Start()
		}()

		sendPegoutChannel <- entities.BaseEvent{}
		sendPegoutChannel <- entities.BaseEvent{}
		sendPegoutChannel <- entities.BaseEvent{}
		sendPegoutChannel <- entities.BaseEvent{}
		time.Sleep(10 * time.Millisecond)

		finalValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "paid")
		assert.Equal(t, 4, int(finalValue))

		labels := getCounterVecLabels(appMetrics.PegoutQuotesMetric, "paid")
		assert.Equal(t, "paid", labels["state"], "State label should be 'paid'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})

	t.Run("should increment refunded pegout quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		pegoutRefundChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.PegoutQuoteCompletedEventId).Return((<-chan entities.Event)(pegoutRefundChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		initialValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "lp_refunded")
		assert.Equal(t, 0, int(initialValue))

		go func() {
			watcher.Start()
		}()

		pegoutRefundChannel <- entities.BaseEvent{}
		pegoutRefundChannel <- entities.BaseEvent{}
		time.Sleep(10 * time.Millisecond)

		finalValue := getCounterVecValue(appMetrics.PegoutQuotesMetric, "lp_refunded")
		assert.Equal(t, 2, int(finalValue))

		labels := getCounterVecLabels(appMetrics.PegoutQuotesMetric, "lp_refunded")
		assert.Equal(t, "lp_refunded", labels["state"], "State label should be 'lp_refunded'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestQuoteMetricsWatcher_Start_PeginMetrics(t *testing.T) {
	t.Run("should increment accepted pegin quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		acceptedPeginChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.AcceptedPeginQuoteEventId).Return((<-chan entities.Event)(acceptedPeginChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		initialValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "accepted")
		assert.Equal(t, 0, int(initialValue))

		go func() {
			watcher.Start()
		}()

		acceptedPeginChannel <- entities.BaseEvent{}
		acceptedPeginChannel <- entities.BaseEvent{}
		time.Sleep(10 * time.Millisecond)

		finalValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "accepted")
		assert.Equal(t, 2, int(finalValue))

		labels := getCounterVecLabels(appMetrics.PeginQuotesMetric, "accepted")
		assert.Equal(t, "accepted", labels["state"], "State label should be 'accepted'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})

	t.Run("should increment call for user completed pegin quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		callForUserChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.CallForUserCompletedEventId).Return((<-chan entities.Event)(callForUserChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		initialValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "paid")
		assert.Equal(t, 0, int(initialValue))

		go func() {
			watcher.Start()
		}()

		callForUserChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}

		time.Sleep(10 * time.Millisecond)

		finalValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "paid")
		assert.Equal(t, 3, int(finalValue))

		labels := getCounterVecLabels(appMetrics.PeginQuotesMetric, "paid")
		assert.Equal(t, "paid", labels["state"], "State label should be 'paid'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})

	t.Run("should increment register pegin completed quotes metric", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		registerPeginChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", quote.RegisterPeginCompletedEventId).Return((<-chan entities.Event)(registerPeginChannel))
		eventBus.On("Subscribe", mock.AnythingOfType("entities.EventId")).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		initialValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "lp_refunded")
		assert.Equal(t, 0, int(initialValue))

		go func() {
			watcher.Start()
		}()

		registerPeginChannel <- entities.BaseEvent{}
		time.Sleep(10 * time.Millisecond)

		finalValue := getCounterVecValue(appMetrics.PeginQuotesMetric, "lp_refunded")
		assert.Equal(t, 1, int(finalValue))

		labels := getCounterVecLabels(appMetrics.PeginQuotesMetric, "lp_refunded")
		assert.Equal(t, "lp_refunded", labels["state"], "State label should be 'lp_refunded'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestQuoteMetricsWatcher_Start_MultipleEvents(t *testing.T) {
	t.Run("should handle multiple events simultaneously", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		// Create channels for ALL event types
		acceptedPegoutChannel := make(chan entities.Event, 5)
		sendPegoutChannel := make(chan entities.Event, 5)
		pegoutRefundChannel := make(chan entities.Event, 5)
		acceptedPeginChannel := make(chan entities.Event, 5)
		callForUserChannel := make(chan entities.Event, 5)
		registerPeginChannel := make(chan entities.Event, 5)

		// Mock all Subscribe calls
		eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return((<-chan entities.Event)(acceptedPegoutChannel))
		eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return((<-chan entities.Event)(sendPegoutChannel))
		eventBus.On("Subscribe", quote.PegoutQuoteCompletedEventId).Return((<-chan entities.Event)(pegoutRefundChannel))
		eventBus.On("Subscribe", quote.AcceptedPeginQuoteEventId).Return((<-chan entities.Event)(acceptedPeginChannel))
		eventBus.On("Subscribe", quote.CallForUserCompletedEventId).Return((<-chan entities.Event)(callForUserChannel))
		eventBus.On("Subscribe", quote.RegisterPeginCompletedEventId).Return((<-chan entities.Event)(registerPeginChannel))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		// Verify all counters start at 0
		initialPegoutAccepted := getCounterVecValue(appMetrics.PegoutQuotesMetric, "accepted")
		initialPegoutPaid := getCounterVecValue(appMetrics.PegoutQuotesMetric, "paid")
		initialPegoutRefunded := getCounterVecValue(appMetrics.PegoutQuotesMetric, "lp_refunded")
		initialPeginAccepted := getCounterVecValue(appMetrics.PeginQuotesMetric, "accepted")
		initialPeginPaid := getCounterVecValue(appMetrics.PeginQuotesMetric, "paid")
		initialPeginRefunded := getCounterVecValue(appMetrics.PeginQuotesMetric, "lp_refunded")

		assert.Equal(t, 0, int(initialPegoutAccepted))
		assert.Equal(t, 0, int(initialPegoutPaid))
		assert.Equal(t, 0, int(initialPegoutRefunded))
		assert.Equal(t, 0, int(initialPeginAccepted))
		assert.Equal(t, 0, int(initialPeginPaid))
		assert.Equal(t, 0, int(initialPeginRefunded))

		go func() {
			watcher.Start()
		}()

		// Send mixed events in different order - some channels get multiple events, others get none
		// This will test that counters accumulate properly and some remain at 0
		acceptedPegoutChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}
		acceptedPegoutChannel <- entities.BaseEvent{}
		acceptedPeginChannel <- entities.BaseEvent{}
		acceptedPegoutChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}
		acceptedPeginChannel <- entities.BaseEvent{}
		sendPegoutChannel <- entities.BaseEvent{}
		callForUserChannel <- entities.BaseEvent{}

		// Intentionally DO NOT send events to:
		// - pegoutRefundChannel (should remain at 0)
		// - registerPeginChannel (should remain at 0)

		time.Sleep(10 * time.Millisecond)

		finalPegoutAccepted := getCounterVecValue(appMetrics.PegoutQuotesMetric, "accepted")
		finalPegoutPaid := getCounterVecValue(appMetrics.PegoutQuotesMetric, "paid")
		finalPegoutRefunded := getCounterVecValue(appMetrics.PegoutQuotesMetric, "lp_refunded")
		finalPeginAccepted := getCounterVecValue(appMetrics.PeginQuotesMetric, "accepted")
		finalPeginPaid := getCounterVecValue(appMetrics.PeginQuotesMetric, "paid")
		finalPeginRefunded := getCounterVecValue(appMetrics.PeginQuotesMetric, "lp_refunded")

		// Verify counters that received events
		assert.Equal(t, 3, int(finalPegoutAccepted))
		assert.Equal(t, 1, int(finalPegoutPaid))
		assert.Equal(t, 2, int(finalPeginAccepted))
		assert.Equal(t, 4, int(finalPeginPaid))

		// Verify counters that received NO events (should remain at 0)
		assert.Equal(t, 0, int(finalPegoutRefunded))
		assert.Equal(t, 0, int(finalPeginRefunded))

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

func TestQuoteMetricsWatcher_Shutdown(t *testing.T) {
	t.Run("should shutdown gracefully", func(t *testing.T) {
		appMetrics := createTestMetrics()
		eventBus := &mocks.EventBusMock{}
		serverInfo := liquidity_provider.NewServerInfoUseCase()

		eventBus.On("Subscribe", quote.AcceptedPegoutQuoteEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", quote.PegoutBtcSentEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", quote.PegoutQuoteCompletedEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", quote.AcceptedPeginQuoteEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", quote.CallForUserCompletedEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", quote.RegisterPeginCompletedEventId).Return(make(<-chan entities.Event))

		watcher := monitoring.NewQuoteMetricsWatcher(appMetrics, eventBus, serverInfo)

		go func() {
			watcher.Start()
		}()

		// Give the Start method time to subscribe to events
		time.Sleep(10 * time.Millisecond)

		closeChannel := make(chan bool, 1)
		watcher.Shutdown(closeChannel)

		// Wait for shutdown signal
		select {
		case <-closeChannel:
			// Shutdown completed successfully
		case <-time.After(time.Second):
			t.Fatal("Shutdown timed out")
		}

		eventBus.AssertExpectations(t)
	})
}

// Helper function to create test metrics
func createTestMetrics() *monitoring.Metrics {
	return &monitoring.Metrics{
		PeginQuotesMetric:  prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_pegin_quotes"}, []string{"state"}),
		PegoutQuotesMetric: prometheus.NewCounterVec(prometheus.CounterOpts{Name: "test_pegout_quotes"}, []string{"state"}),
		ServerInfoMetric:   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_server_info"}, []string{"version", "commit"}),
	}
}

// Helper function to set ServerInfo build variables for tests that need them
func setServerInfoBuildVars(version, revision string) {
	liquidity_provider.BuildVersion = version
	liquidity_provider.BuildRevision = revision
}

// Helper function to get metric value from a GaugeVec
func getGaugeVecValue(gaugeVec *prometheus.GaugeVec, labelValues ...string) float64 {
	gauge := gaugeVec.WithLabelValues(labelValues...)
	metric := &dto.Metric{}
	// nolint:errcheck
	gauge.Write(metric)
	return metric.GetGauge().GetValue()
}

// Helper function to get metric labels from a GaugeVec
func getGaugeVecLabels(gaugeVec *prometheus.GaugeVec, labelValues ...string) map[string]string {
	gauge := gaugeVec.WithLabelValues(labelValues...)
	metric := &dto.Metric{}
	// nolint:errcheck
	gauge.Write(metric)

	labels := make(map[string]string)
	for _, label := range metric.GetLabel() {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}

// Helper function to get counter value from a CounterVec
func getCounterVecValue(counterVec *prometheus.CounterVec, labelValues ...string) float64 {
	counter := counterVec.WithLabelValues(labelValues...)
	metric := &dto.Metric{}
	// nolint:errcheck
	counter.Write(metric)
	return metric.GetCounter().GetValue()
}

// Helper function to get counter labels from a CounterVec
func getCounterVecLabels(counterVec *prometheus.CounterVec, labelValues ...string) map[string]string {
	counter := counterVec.WithLabelValues(labelValues...)
	metric := &dto.Metric{}
	// nolint:errcheck
	counter.Write(metric)

	labels := make(map[string]string)
	for _, label := range metric.GetLabel() {
		labels[label.GetName()] = label.GetValue()
	}
	return labels
}
