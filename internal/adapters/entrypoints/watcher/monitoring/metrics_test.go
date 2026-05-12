package monitoring_test

import (
	"math/big"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMetrics(t *testing.T) {
	t.Run("should create metrics with proper configuration and register them", func(t *testing.T) {
		metrics, registerer := createMetricsWithMock(t)

		require.NotNil(t, metrics)
		assert.NotNil(t, metrics.PeginQuotesMetric)
		assert.NotNil(t, metrics.PegoutQuotesMetric)
		assert.NotNil(t, metrics.ServerInfoMetric)
		assert.NotNil(t, metrics.AssetsMetrics)

		// Verify metric names and help text by checking descriptors
		peginDesc := getMetricDesc(metrics.PeginQuotesMetric)
		pegoutDesc := getMetricDesc(metrics.PegoutQuotesMetric)
		serverInfoDesc := getMetricDesc(metrics.ServerInfoMetric)
		assetsDesc := getMetricDesc(metrics.AssetsMetrics)

		assert.Contains(t, peginDesc, "lps_pegin_quotes")
		assert.Contains(t, peginDesc, "Pegin quotes processed")
		assert.Contains(t, peginDesc, "state")

		assert.Contains(t, pegoutDesc, "lps_pegout_quotes")
		assert.Contains(t, pegoutDesc, "Pegout quotes processed")
		assert.Contains(t, pegoutDesc, "state")

		assert.Contains(t, serverInfoDesc, "lps_server_info")
		assert.Contains(t, serverInfoDesc, "Server information")
		assert.Contains(t, serverInfoDesc, "version")
		assert.Contains(t, serverInfoDesc, "commit")

		assert.Contains(t, assetsDesc, "lps_assets_balances")
		assert.Contains(t, assetsDesc, "Liquidity provider asset balances and metrics")
		assert.Contains(t, assetsDesc, "currency")
		assert.Contains(t, assetsDesc, "type")

		registerer.AssertExpectations(t)
	})
}

// nolint:funlen
func TestMetrics_UpdateAssetsFromReport(t *testing.T) {
	t.Run("should update all asset metrics with correct values and labels", func(t *testing.T) {
		metrics, _ := createMetricsWithMock(t)
		report := createTestAssetReport()

		metrics.UpdateAssetsFromReport(report)

		// Verify RBTC metrics are set with correct values
		assert.InDelta(t, 9.5, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "total"), 0.0001)
		assert.InDelta(t, 3.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "location_rsk_wallet"), 0.0001)
		assert.InDelta(t, 1.5, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "location_lbc"), 0.0001)
		assert.InDelta(t, 5.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "location_federation"), 0.0001)
		assert.InDelta(t, 2.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_reserved_for_users"), 0.0001)
		assert.InDelta(t, 0.5, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_waiting_refund"), 0.0001)
		assert.InDelta(t, 5.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_available"), 0.0001)

		// Verify BTC metrics are set with correct values
		assert.InDelta(t, 9.2, getGaugeVecValue(metrics.AssetsMetrics, "btc", "total"), 0.0001)
		assert.InDelta(t, 2.8, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_btc_wallet"), 0.0001)
		assert.InDelta(t, 1.8, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_federation"), 0.0001)
		assert.InDelta(t, 3.6, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_rsk_wallet"), 0.0001)
		assert.InDelta(t, 0.1, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_lbc"), 0.0001)
		assert.InDelta(t, 1.8, getGaugeVecValue(metrics.AssetsMetrics, "btc", "allocation_reserved_for_users"), 0.0001)
		assert.InDelta(t, 0.1, getGaugeVecValue(metrics.AssetsMetrics, "btc", "allocation_waiting_refund"), 0.0001)
		assert.InDelta(t, 4.5, getGaugeVecValue(metrics.AssetsMetrics, "btc", "allocation_available"), 0.0001)

		// Verify labels are correct for RBTC metrics
		rbtcLabels := getGaugeVecLabels(metrics.AssetsMetrics, "rbtc", "total")
		assert.Equal(t, "rbtc", rbtcLabels["currency"])
		assert.Equal(t, "total", rbtcLabels["type"])

		// Verify labels are correct for BTC metrics
		btcLabels := getGaugeVecLabels(metrics.AssetsMetrics, "btc", "location_lbc")
		assert.Equal(t, "btc", btcLabels["currency"])
		assert.Equal(t, "location_lbc", btcLabels["type"])
	})

	t.Run("should handle edge case values correctly (zero, large, complex decimals)", func(t *testing.T) {
		metrics, _ := createMetricsWithMock(t)

		// Create report with mixed edge case values:
		// - Zero values for some metrics
		// - Large values for others
		// - Complex decimal values
		report := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: createWeiFromString("62277000000000000000"), // Total BTC
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(2500), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 2500 BTC
					Federation: createWeiFromString("56789012345678900000"),                                                                     // 56.7890123456789 BTC
					RskWallet:  entities.NewWei(0),                                                                                              // 0 BTC
					Lbc:        entities.NewWei(987654321000000000),                                                                             // 0.987654321 BTC
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: createWeiFromString("56789012345678900000"), // 56.7890123456789 BTC
					WaitingForRefund: entities.NewWei(987654321000000000),         // 0.987654321 BTC
					Available:        entities.NewWei(0),                          // 0 BTC
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(1050), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // ~1050 RBTC
				Location: reports.RbtcAssetLocation{
					RskWallet:  createWeiFromString("7890123456789000000"),                                                                      // 7.890123456789 RBTC
					Lbc:        entities.NewWei(0),                                                                                              // 0 RBTC
					Federation: entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 1000 RBTC
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: createWeiFromString("34598535894857007656"),                                                                     // 34.598535894857007656 RBTC
					WaitingForRefund: entities.NewWei(123456789000000),                                                                                // 0.000123456789 RBTC
					Available:        entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 1000 RBTC
				},
			},
		}

		metrics.UpdateAssetsFromReport(report)

		// Verify zero values
		assert.InDelta(t, 0.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "location_lbc"), 0.0001)
		assert.InDelta(t, 0.0, getGaugeVecValue(metrics.AssetsMetrics, "btc", "allocation_available"), 0.0001)

		// Verify complex decimal values
		assert.InDelta(t, 34.598535894857007656, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_reserved_for_users"), 1e-15)
		assert.InDelta(t, 0.000123456789, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_waiting_refund"), 1e-15)
		assert.InDelta(t, 7.890123456789, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "location_rsk_wallet"), 1e-12)
		assert.InDelta(t, 56.7890123456789, getGaugeVecValue(metrics.AssetsMetrics, "btc", "allocation_reserved_for_users"), 1e-12)
		assert.InDelta(t, 0.987654321, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_lbc"), 1e-9)

		// Verify large values
		assert.InDelta(t, 1000.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "allocation_available"), 0.0001)
		assert.InDelta(t, 2500.0, getGaugeVecValue(metrics.AssetsMetrics, "btc", "location_btc_wallet"), 0.0001)
	})
}

// Helper function to create test asset report with known values
func createTestAssetReport() reports.GetAssetsReportResult {
	return reports.GetAssetsReportResult{
		BtcAssetReport: reports.BtcAssetReport{
			Total: entities.NewBigWei(big.NewInt(9200000000000000000)), // 9.2 BTC total
			Location: reports.BtcAssetLocation{
				BtcWallet:  entities.NewBigWei(big.NewInt(2800000000000000000)), // 2.8 BTC
				Federation: entities.NewBigWei(big.NewInt(1800000000000000000)), // 1.8 BTC
				RskWallet:  entities.NewBigWei(big.NewInt(3600000000000000000)), // 3.6 BTC
				Lbc:        entities.NewWei(100000000000000000),                 // 0.1 BTC
			},
			Allocation: reports.BtcAssetAllocation{
				ReservedForUsers: entities.NewBigWei(big.NewInt(1800000000000000000)), // 1.8 BTC
				WaitingForRefund: entities.NewWei(100000000000000000),                 // 0.1 BTC
				Available:        entities.NewBigWei(big.NewInt(4500000000000000000)), // 4.5 BTC
			},
		},
		RbtcAssetReport: reports.RbtcAssetReport{
			Total: entities.NewBigWei(big.NewInt(0).Add(big.NewInt(9000000000000000000), big.NewInt(500000000000000000))), // 9.5 RBTC total
			Location: reports.RbtcAssetLocation{
				RskWallet:  entities.NewBigWei(big.NewInt(3000000000000000000)), // 3.0 RBTC
				Lbc:        entities.NewBigWei(big.NewInt(1500000000000000000)), // 1.5 RBTC
				Federation: entities.NewBigWei(big.NewInt(5000000000000000000)), // 5.0 RBTC
			},
			Allocation: reports.RbtcAssetAllocation{
				ReservedForUsers: entities.NewBigWei(big.NewInt(2000000000000000000)), // 2.0 RBTC
				WaitingForRefund: entities.NewWei(500000000000000000),                 // 0.5 RBTC
				Available:        entities.NewBigWei(big.NewInt(5000000000000000000)), // 5.0 RBTC
			},
		},
	}
}

// Helper function to create metrics with fresh mock registerer for each test
func createMetricsWithMock(t *testing.T) (*monitoring.Metrics, *mocks.RegistererMock) {
	registerer := mocks.NewRegistererMock(t)
	registerer.On("MustRegister",
		mock.AnythingOfType("*prometheus.CounterVec"), // PegoutQuotesMetric
		mock.AnythingOfType("*prometheus.CounterVec"), // PeginQuotesMetric
		mock.AnythingOfType("*prometheus.GaugeVec"),   // ServerInfoMetric
		mock.AnythingOfType("*prometheus.GaugeVec"),   // AssetsMetrics
	).Return()

	metrics := monitoring.NewMetrics(registerer)
	return metrics, registerer
}

// Helper function to get metric descriptor as string
func getMetricDesc(metric prometheus.Collector) string {
	descChannel := make(chan *prometheus.Desc, 1)
	metric.Describe(descChannel)
	desc := <-descChannel
	return desc.String()
}

// Helper function to create Wei values from string (for large numbers)
func createWeiFromString(weiStr string) *entities.Wei {
	val := new(big.Int)
	val.SetString(weiStr, 10)
	return entities.NewBigWei(val)
}
