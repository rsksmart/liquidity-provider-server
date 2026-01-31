package monitoring_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewColdWalletMetricsWatcher(t *testing.T) {
	t.Run("should create cold wallet metrics watcher with all dependencies properly configured", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

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

		closeChannelField, found := watcherType.FieldByName("closeChannel")
		require.True(t, found, "closeChannel field should exist")
		assert.Equal(t, "chan struct {}", closeChannelField.Type.String())

		// Verify closeChannel is properly initialized by checking it's not zero value
		closeChannelValue := watcherValue.FieldByName("closeChannel")
		assert.False(t, closeChannelValue.IsZero(), "closeChannel should be initialized")
		assert.Equal(t, reflect.Chan, closeChannelValue.Kind())
		assert.Equal(t, 1, closeChannelValue.Cap(), "closeChannel should have capacity of 1")
	})
}

func TestColdWalletMetricsWatcher_Prepare(t *testing.T) {
	t.Run("should always return nil without error", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)
		err := watcher.Prepare(context.Background())

		assert.NoError(t, err)
	})

	t.Run("should return nil regardless of context cancellation", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)
		err := watcher.Prepare(ctx)

		assert.NoError(t, err)
	})
}

// nolint:funlen
func TestColdWalletMetricsWatcher_Start_RbtcThresholdEvents(t *testing.T) {
	t.Run("should increment RBTC threshold transfer counter and set amount", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		rbtcThresholdChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return((<-chan entities.Event)(rbtcThresholdChannel))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

		// Verify counter starts at 0
		initialCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "threshold")
		assert.Equal(t, 0, int(initialCounter))

		// Verify gauge starts at 0
		initialGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "rbtc")
		assert.InDelta(t, 0.0, initialGauge, 0.0001)

		go func() {
			watcher.Start()
		}()

		// Send events with specific amounts
		rbtcThresholdChannel <- cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(1000000000000000000), // 1 RBTC
			TxHash: "0xabc123",
			Fee:    entities.NewWei(21000),
		}
		time.Sleep(10 * time.Millisecond)

		rbtcThresholdChannel <- cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(2500000000000000000), // 2.5 RBTC
			TxHash: "0xdef456",
			Fee:    entities.NewWei(21000),
		}
		time.Sleep(10 * time.Millisecond)

		// Verify counter was incremented twice
		finalCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "threshold")
		assert.Equal(t, 2, int(finalCounter))

		// Verify gauge has the last amount (2.5 RBTC)
		finalGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "rbtc")
		assert.InDelta(t, 2.5, finalGauge, 0.0001)

		// Verify the labels are correct
		counterLabels := getCounterVecLabels(appMetrics.ColdWalletTransfersMetric, "rbtc", "threshold")
		assert.Equal(t, "rbtc", counterLabels["currency"], "Currency label should be 'rbtc'")
		assert.Equal(t, "threshold", counterLabels["reason"], "Reason label should be 'threshold'")

		gaugeLabels := getGaugeVecLabels(appMetrics.ColdWalletLastAmountMetric, "rbtc")
		assert.Equal(t, "rbtc", gaugeLabels["currency"], "Currency label should be 'rbtc'")

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestColdWalletMetricsWatcher_Start_BtcThresholdEvents(t *testing.T) {
	t.Run("should increment BTC threshold transfer counter and set amount", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		btcThresholdChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return((<-chan entities.Event)(btcThresholdChannel))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

		initialCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "threshold")
		assert.Equal(t, 0, int(initialCounter))

		initialGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "btc")
		assert.InDelta(t, 0.0, initialGauge, 0.0001)

		go func() {
			watcher.Start()
		}()

		btcThresholdChannel <- cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(500000000000000000), // 0.5 BTC
			TxHash: "btc_tx_123",
			Fee:    entities.NewWei(5000),
		}
		time.Sleep(10 * time.Millisecond)

		btcThresholdChannel <- cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(1750000000000000000), // 1.75 BTC
			TxHash: "btc_tx_456",
			Fee:    entities.NewWei(5000),
		}
		time.Sleep(10 * time.Millisecond)

		btcThresholdChannel <- cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(3000000000000000000), // 3.0 BTC
			TxHash: "btc_tx_789",
			Fee:    entities.NewWei(5000),
		}
		time.Sleep(10 * time.Millisecond)

		// Verify counter was incremented 3 times
		finalCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "threshold")
		assert.Equal(t, 3, int(finalCounter))

		// Verify gauge has the last amount (3.0 BTC)
		finalGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "btc")
		assert.InDelta(t, 3.0, finalGauge, 0.0001)

		counterLabels := getCounterVecLabels(appMetrics.ColdWalletTransfersMetric, "btc", "threshold")
		assert.Equal(t, "btc", counterLabels["currency"])
		assert.Equal(t, "threshold", counterLabels["reason"])

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestColdWalletMetricsWatcher_Start_RbtcTimeForcingEvents(t *testing.T) {
	t.Run("should increment RBTC time forcing transfer counter and set amount", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		rbtcTimeForcingChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return((<-chan entities.Event)(rbtcTimeForcingChannel))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

		initialCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "time_forcing")
		assert.Equal(t, 0, int(initialCounter))

		go func() {
			watcher.Start()
		}()

		rbtcTimeForcingChannel <- cold_wallet.RbtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToTimeForcingEventId),
			Amount: entities.NewWei(750000000000000000), // 0.75 RBTC
			TxHash: "0xtimeforce1",
			Fee:    entities.NewWei(21000),
		}
		time.Sleep(10 * time.Millisecond)

		finalCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "time_forcing")
		assert.Equal(t, 1, int(finalCounter))

		finalGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "rbtc")
		assert.InDelta(t, 0.75, finalGauge, 0.0001)

		counterLabels := getCounterVecLabels(appMetrics.ColdWalletTransfersMetric, "rbtc", "time_forcing")
		assert.Equal(t, "rbtc", counterLabels["currency"])
		assert.Equal(t, "time_forcing", counterLabels["reason"])

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestColdWalletMetricsWatcher_Start_BtcTimeForcingEvents(t *testing.T) {
	t.Run("should increment BTC time forcing transfer counter and set amount", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		btcTimeForcingChannel := make(chan entities.Event, 1)
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return((<-chan entities.Event)(btcTimeForcingChannel))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

		initialCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "time_forcing")
		assert.Equal(t, 0, int(initialCounter))

		go func() {
			watcher.Start()
		}()

		btcTimeForcingChannel <- cold_wallet.BtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToTimeForcingEventId),
			Amount: entities.NewWei(1250000000000000000), // 1.25 BTC
			TxHash: "btc_timeforce_1",
			Fee:    entities.NewWei(5000),
		}
		time.Sleep(10 * time.Millisecond)

		btcTimeForcingChannel <- cold_wallet.BtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToTimeForcingEventId),
			Amount: entities.NewWei(2000000000000000000), // 2.0 BTC
			TxHash: "btc_timeforce_2",
			Fee:    entities.NewWei(5000),
		}
		time.Sleep(10 * time.Millisecond)

		finalCounter := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "time_forcing")
		assert.Equal(t, 2, int(finalCounter))

		finalGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "btc")
		assert.InDelta(t, 2.0, finalGauge, 0.0001)

		counterLabels := getCounterVecLabels(appMetrics.ColdWalletTransfersMetric, "btc", "time_forcing")
		assert.Equal(t, "btc", counterLabels["currency"])
		assert.Equal(t, "time_forcing", counterLabels["reason"])

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

// nolint:funlen
func TestColdWalletMetricsWatcher_Start_MultipleEvents(t *testing.T) {
	t.Run("should handle multiple events simultaneously from all channels", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		// Create channels for ALL event types
		rbtcThresholdChannel := make(chan entities.Event, 5)
		btcThresholdChannel := make(chan entities.Event, 5)
		rbtcTimeForcingChannel := make(chan entities.Event, 5)
		btcTimeForcingChannel := make(chan entities.Event, 5)

		// Mock all Subscribe calls
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return((<-chan entities.Event)(rbtcThresholdChannel))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return((<-chan entities.Event)(btcThresholdChannel))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return((<-chan entities.Event)(rbtcTimeForcingChannel))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return((<-chan entities.Event)(btcTimeForcingChannel))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

		// Verify all counters start at 0
		initialRbtcThreshold := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "threshold")
		initialBtcThreshold := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "threshold")
		initialRbtcTimeForcing := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "time_forcing")
		initialBtcTimeForcing := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "time_forcing")

		assert.Equal(t, 0, int(initialRbtcThreshold))
		assert.Equal(t, 0, int(initialBtcThreshold))
		assert.Equal(t, 0, int(initialRbtcTimeForcing))
		assert.Equal(t, 0, int(initialBtcTimeForcing))

		go func() {
			watcher.Start()
		}()

		// Send mixed events in different order
		rbtcThresholdChannel <- cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(1000000000000000000),
			TxHash: "0x1",
			Fee:    entities.NewWei(21000),
		}
		btcThresholdChannel <- cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(500000000000000000),
			TxHash: "btc1",
			Fee:    entities.NewWei(5000),
		}
		rbtcThresholdChannel <- cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(2000000000000000000),
			TxHash: "0x2",
			Fee:    entities.NewWei(21000),
		}
		rbtcTimeForcingChannel <- cold_wallet.RbtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToTimeForcingEventId),
			Amount: entities.NewWei(1500000000000000000),
			TxHash: "0x3",
			Fee:    entities.NewWei(21000),
		}
		btcTimeForcingChannel <- cold_wallet.BtcTransferredDueToTimeForcingEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToTimeForcingEventId),
			Amount: entities.NewWei(3000000000000000000),
			TxHash: "btc2",
			Fee:    entities.NewWei(5000),
		}
		btcThresholdChannel <- cold_wallet.BtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.BtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(2500000000000000000),
			TxHash: "btc3",
			Fee:    entities.NewWei(5000),
		}
		rbtcThresholdChannel <- cold_wallet.RbtcTransferredDueToThresholdEvent{
			Event:  entities.NewBaseEvent(cold_wallet.RbtcTransferredDueToThresholdEventId),
			Amount: entities.NewWei(500000000000000000),
			TxHash: "0x4",
			Fee:    entities.NewWei(21000),
		}

		// Intentionally DO NOT send more events to time_forcing channels to test varied counts

		time.Sleep(20 * time.Millisecond)

		finalRbtcThreshold := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "threshold")
		finalBtcThreshold := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "threshold")
		finalRbtcTimeForcing := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "rbtc", "time_forcing")
		finalBtcTimeForcing := getCounterVecValue(appMetrics.ColdWalletTransfersMetric, "btc", "time_forcing")

		// Verify counters that received events
		assert.Equal(t, 3, int(finalRbtcThreshold))
		assert.Equal(t, 2, int(finalBtcThreshold))
		assert.Equal(t, 1, int(finalRbtcTimeForcing))
		assert.Equal(t, 1, int(finalBtcTimeForcing))

		// Verify gauges have been updated (they track the last amount for each currency)
		// Due to the async nature of event processing, we just verify they're non-zero
		// and contain one of the expected values (either threshold or time_forcing amounts)
		finalRbtcGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "rbtc")
		finalBtcGauge := getGaugeVecValue(appMetrics.ColdWalletLastAmountMetric, "btc")

		// RBTC gauge should be either 0.5 (last threshold) or 1.5 (time_forcing) depending on processing order
		assert.True(t, finalRbtcGauge == 0.5 || finalRbtcGauge == 1.5, 
			"RBTC gauge should be 0.5 or 1.5, got %.2f", finalRbtcGauge)
		
		// BTC gauge should be either 2.5 (last threshold) or 3.0 (time_forcing) depending on processing order
		assert.True(t, finalBtcGauge == 2.5 || finalBtcGauge == 3.0, 
			"BTC gauge should be 2.5 or 3.0, got %.2f", finalBtcGauge)

		watcher.Shutdown(make(chan bool, 1))
		eventBus.AssertExpectations(t)
	})
}

func TestColdWalletMetricsWatcher_Shutdown(t *testing.T) {
	t.Run("should shutdown gracefully", func(t *testing.T) {
		appMetrics := createTestColdWalletMetrics()
		eventBus := &mocks.EventBusMock{}

		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToThresholdEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.RbtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))
		eventBus.On("Subscribe", cold_wallet.BtcTransferredDueToTimeForcingEventId).Return(make(<-chan entities.Event))

		watcher := monitoring.NewColdWalletMetricsWatcher(appMetrics, eventBus)

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

// Helper function to create test cold wallet metrics
func createTestColdWalletMetrics() *monitoring.Metrics {
	return &monitoring.Metrics{
		ColdWalletTransfersMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "test_cold_wallet_transfers"},
			[]string{"currency", "reason"},
		),
		ColdWalletLastAmountMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: "test_cold_wallet_last_amount"},
			[]string{"currency"},
		),
	}
}
