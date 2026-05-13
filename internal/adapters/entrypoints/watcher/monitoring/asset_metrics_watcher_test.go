package monitoring_test

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAssetReportWatcher(t *testing.T) {
	t.Run("should create asset report watcher with all dependencies properly configured", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)

		require.NotNil(t, watcher)

		// Use reflection to verify internal structure (type checking only)
		watcherValue := reflect.ValueOf(watcher).Elem()
		watcherType := watcherValue.Type()

		// Verify all expected fields exist with correct types
		appMetricsField, found := watcherType.FieldByName("appMetrics")
		require.True(t, found, "appMetrics field should exist")
		assert.Equal(t, "*monitoring.Metrics", appMetricsField.Type.String())

		getAssetReportUseCaseField, found := watcherType.FieldByName("getAssetReportUseCase")
		require.True(t, found, "getAssetReportUseCase field should exist")
		assert.Equal(t, "monitoring.GetAssetReportUseCase", getAssetReportUseCaseField.Type.String())

		tickerField, found := watcherType.FieldByName("ticker")
		require.True(t, found, "ticker field should exist")
		assert.Equal(t, "utils.Ticker", tickerField.Type.String())

		watcherStopChannelField, found := watcherType.FieldByName("watcherStopChannel")
		require.True(t, found, "watcherStopChannel field should exist")
		assert.Equal(t, "chan struct {}", watcherStopChannelField.Type.String())

		// Verify watcherStopChannel is properly initialized by checking it's not zero value
		watcherStopChannelValue := watcherValue.FieldByName("watcherStopChannel")
		assert.False(t, watcherStopChannelValue.IsZero(), "watcherStopChannel should be initialized")
		assert.Equal(t, reflect.Chan, watcherStopChannelValue.Kind())
		assert.Equal(t, 1, watcherStopChannelValue.Cap(), "watcherStopChannel should have capacity of 1")
	})
}

func TestAssetReportWatcher_Prepare(t *testing.T) {
	t.Run("should always return nil without error", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)
		err := watcher.Prepare(context.Background())

		assert.NoError(t, err)
	})

	t.Run("should return nil regardless of context cancellation", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)
		err := watcher.Prepare(ctx)

		assert.NoError(t, err)
	})
}

// nolint:funlen
func TestAssetReportWatcher_Start(t *testing.T) {
	t.Run("should update metrics when ticker fires successfully", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		tickerChannel := make(chan time.Time, 1)
		ticker.On("C").Return((<-chan time.Time)(tickerChannel))
		ticker.On("Stop").Return()

		testReport := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(9200000000000000000), // 9.2 BTC total
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(2800000000000000000), // 2.8 BTC
					Federation: entities.NewWei(1800000000000000000), // 1.8 BTC
					RskWallet:  entities.NewWei(3600000000000000000), // 3.6 BTC
					Lbc:        entities.NewWei(1000000000000000000), // 1.0 BTC
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(1800000000000000000), // 1.8 BTC
					WaitingForRefund: entities.NewWei(100000000000000000),  // 0.1 BTC
					Available:        entities.NewWei(4500000000000000000), // 4.5 BTC
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewBigWei(big.NewInt(0).Add(big.NewInt(9000000000000000000), big.NewInt(500000000000000000))), // 9.5 RBTC total
				Location: reports.RbtcAssetLocation{
					RskWallet:  entities.NewWei(3000000000000000000), // 3.0 RBTC
					Lbc:        entities.NewWei(1500000000000000000), // 1.5 RBTC
					Federation: entities.NewWei(5000000000000000000), // 5.0 RBTC
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: entities.NewWei(2000000000000000000), // 2.0 RBTC
					WaitingForRefund: entities.NewWei(500000000000000000),  // 0.5 RBTC
					Available:        entities.NewWei(5000000000000000000), // 5.0 RBTC
				},
			},
		}
		assetsReportUseCase.On("Run", mock.Anything).Return(testReport, nil)

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)

		// Verify initial metric values are zero
		initialRbtcTotal := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelTotal)
		initialBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 0.0, initialRbtcTotal, 0.0001)
		assert.InDelta(t, 0.0, initialBtcAvailable, 0.0001)

		go func() {
			watcher.Start()
		}()

		// Send ticker event to trigger metric update
		tickerChannel <- time.Now()

		// Give some time for the event to be processed
		time.Sleep(50 * time.Millisecond)

		// Verify metrics were updated with expected values
		finalRbtcTotal := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelTotal)
		finalBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 9.5, finalRbtcTotal, 0.0001)
		assert.InDelta(t, 4.5, finalBtcAvailable, 0.0001)

		// Verify labels are correct
		rbtcLabels := getAssetGaugeLabels(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelTotal)
		assert.Equal(t, monitoring.MetricLabelRbtc, rbtcLabels["currency"])
		assert.Equal(t, monitoring.MetricLabelTotal, rbtcLabels["type"])

		// Properly shutdown and wait for completion
		closeChannel := make(chan bool, 1)
		watcher.Shutdown(closeChannel)
		<-closeChannel // Wait for shutdown to complete

		// Give a small delay to ensure the Start() goroutine processes the shutdown
		time.Sleep(10 * time.Millisecond)
		ticker.AssertExpectations(t)
	})

	t.Run("should handle use case error gracefully without updating metrics", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		tickerChannel := make(chan time.Time, 1)
		ticker.On("C").Return((<-chan time.Time)(tickerChannel))
		ticker.On("Stop").Return()

		// Mock use case to return error
		expectedError := errors.New("failed to get asset report")
		assetsReportUseCase.On("Run", mock.Anything).Return(reports.GetAssetsReportResult{}, expectedError)

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)

		// Verify initial metric values are zero
		initialRbtcTotal := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelTotal)
		initialBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 0.0, initialRbtcTotal, 0.0001)
		assert.InDelta(t, 0.0, initialBtcAvailable, 0.0001)

		go func() {
			watcher.Start()
		}()

		// Send ticker event to trigger metric update attempt
		tickerChannel <- time.Now()

		// Give some time for the event to be processed
		time.Sleep(50 * time.Millisecond)

		// Verify metrics remain unchanged (zero) due to error
		finalRbtcTotal := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelTotal)
		finalBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 0.0, finalRbtcTotal, 0.0001)
		assert.InDelta(t, 0.0, finalBtcAvailable, 0.0001)

		// Properly shutdown
		closeChannel := make(chan bool, 1)
		watcher.Shutdown(closeChannel)
		<-closeChannel

		time.Sleep(10 * time.Millisecond)
		ticker.AssertExpectations(t)
	})

	t.Run("should process multiple ticker events and update metrics consistently", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		tickerChannel := make(chan time.Time, 3)
		ticker.On("C").Return((<-chan time.Time)(tickerChannel))
		ticker.On("Stop").Return()

		firstReport := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(3600000000000000000), // 3.6 BTC total
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(1000000000000000000), // 1.0 BTC
					Federation: entities.NewWei(500000000000000000),  // 0.5 BTC
					RskWallet:  entities.NewWei(1500000000000000000), // 1.5 BTC
					Lbc:        entities.NewWei(600000000000000000),  // 0.6 BTC
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(500000000000000000),  // 0.5 BTC
					WaitingForRefund: entities.NewWei(100000000000000000),  // 0.1 BTC
					Available:        entities.NewWei(2000000000000000000), // 2.0 BTC
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewWei(7500000000000000000), // 7.5 RBTC total
				Location: reports.RbtcAssetLocation{
					RskWallet:  entities.NewWei(1500000000000000000), // 1.5 RBTC
					Lbc:        entities.NewWei(1000000000000000000), // 1.0 RBTC
					Federation: entities.NewWei(5000000000000000000), // 5.0 RBTC
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: entities.NewWei(2000000000000000000), // 2.0 RBTC
					WaitingForRefund: entities.NewWei(0),                   // 0 RBTC
					Available:        entities.NewWei(3000000000000000000), // 3.0 RBTC
				},
			},
		}

		secondReport := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(5800000000000000000), // 5.8 BTC total
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(1800000000000000000), // 1.8 BTC
					Federation: entities.NewWei(800000000000000000),  // 0.8 BTC
					RskWallet:  entities.NewWei(2400000000000000000), // 2.4 BTC
					Lbc:        entities.NewWei(800000000000000000),  // 0.8 BTC
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(800000000000000000),  // 0.8 BTC
					WaitingForRefund: entities.NewWei(200000000000000000),  // 0.2 BTC
					Available:        entities.NewWei(3000000000000000000), // 3.0 BTC
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 10.0 RBTC total
				Location: reports.RbtcAssetLocation{
					RskWallet:  entities.NewWei(2500000000000000000), // 2.5 RBTC
					Lbc:        entities.NewWei(2000000000000000000), // 2.0 RBTC
					Federation: entities.NewWei(5500000000000000000), // 5.5 RBTC
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: entities.NewWei(3000000000000000000), // 3.0 RBTC
					WaitingForRefund: entities.NewWei(500000000000000000),  // 0.5 RBTC
					Available:        entities.NewWei(4000000000000000000), // 4.0 RBTC
				},
			},
		}

		// Set up expectations for multiple calls
		assetsReportUseCase.On("Run", mock.Anything).Return(firstReport, nil).Once()
		assetsReportUseCase.On("Run", mock.Anything).Return(secondReport, nil).Once()

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)

		go func() {
			watcher.Start()
		}()

		// First ticker event
		tickerChannel <- time.Now()
		time.Sleep(50 * time.Millisecond)

		// Verify first update
		firstRbtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelAllocationAvailable)
		firstBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 3.0, firstRbtcAvailable, 0.0001)
		assert.InDelta(t, 2.0, firstBtcAvailable, 0.0001)

		// Second ticker event
		tickerChannel <- time.Now()
		time.Sleep(50 * time.Millisecond)

		// Verify second update (metrics should reflect new values)
		secondRbtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelRbtc, monitoring.MetricLabelAllocationAvailable)
		secondBtcAvailable := getAssetGaugeValue(appMetrics.AssetsMetrics, monitoring.MetricLabelBtc, monitoring.MetricLabelAllocationAvailable)
		assert.InDelta(t, 4.0, secondRbtcAvailable, 0.0001)
		assert.InDelta(t, 3.0, secondBtcAvailable, 0.0001)

		// Properly shutdown
		closeChannel := make(chan bool, 1)
		watcher.Shutdown(closeChannel)
		<-closeChannel

		time.Sleep(10 * time.Millisecond)
		ticker.AssertExpectations(t)
	})
}

func TestAssetReportWatcher_Shutdown(t *testing.T) {
	t.Run("should shutdown gracefully and stop ticker", func(t *testing.T) {
		appMetrics := createTestAssetMetrics()
		assetsReportUseCase := mocks.NewGetAssetsReportUseCaseMock(t)
		ticker := &mocks.TickerMock{}

		ticker.On("C").Return(make(<-chan time.Time))
		ticker.On("Stop").Return()

		watcher := monitoring.NewAssetReportWatcher(appMetrics, assetsReportUseCase, ticker)

		go func() {
			watcher.Start()
		}()

		// Give the Start method time to begin listening
		time.Sleep(50 * time.Millisecond)

		closeChannel := make(chan bool, 1)
		watcher.Shutdown(closeChannel)

		// Wait for shutdown signal
		select {
		case <-closeChannel:
			// Shutdown completed successfully
		case <-time.After(time.Second):
			t.Fatal("Shutdown timed out")
		}

		// Give a small delay to ensure the Start() goroutine processes the shutdown
		time.Sleep(10 * time.Millisecond)
		ticker.AssertExpectations(t)
	})
}

// Helper function to create test asset metrics
func createTestAssetMetrics() *monitoring.Metrics {
	return &monitoring.Metrics{
		AssetsMetrics: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "test_assets_balances"}, []string{"currency", "type"}),
	}
}

// Helper function to get metric value from Assets GaugeVec
func getAssetGaugeValue(gaugeVec *prometheus.GaugeVec, labelValues ...string) float64 {
	gauge := gaugeVec.WithLabelValues(labelValues...)
	metric := &dto.Metric{}
	// nolint:errcheck
	gauge.Write(metric)
	return metric.GetGauge().GetValue()
}

// Helper function to get metric labels from Assets GaugeVec
func getAssetGaugeLabels(gaugeVec *prometheus.GaugeVec, labelValues ...string) map[string]string {
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
