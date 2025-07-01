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

func TestMetrics_UpdateAssetsFromReport(t *testing.T) {
	t.Run("should update all asset metrics with correct values and labels", func(t *testing.T) {
		metrics, _ := createMetricsWithMock(t)
		report := createTestAssetReport()

		metrics.UpdateAssetsFromReport(report)

		// Verify RBTC metrics are set with correct values
		assert.InDelta(t, 1.5, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "locked_lbc"), 0.0001)
		assert.InDelta(t, 2.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "locked_for_users"), 0.0001)
		assert.InDelta(t, 0.5, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "waiting_refund"), 0.0001)
		assert.InDelta(t, 5.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "liquidity"), 0.0001)
		assert.InDelta(t, 3.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "wallet_balance"), 0.0001)

		// Verify BTC metrics are set with correct values
		assert.InDelta(t, 1.8, getGaugeVecValue(metrics.AssetsMetrics, "btc", "locked_for_users"), 0.0001)
		assert.InDelta(t, 4.5, getGaugeVecValue(metrics.AssetsMetrics, "btc", "liquidity"), 0.0001)
		assert.InDelta(t, 2.8, getGaugeVecValue(metrics.AssetsMetrics, "btc", "wallet_balance"), 0.0001)
		assert.InDelta(t, 0.1, getGaugeVecValue(metrics.AssetsMetrics, "btc", "rebalancing"), 0.0001)

		// Verify labels are correct for RBTC metrics
		rbtcLabels := getGaugeVecLabels(metrics.AssetsMetrics, "rbtc", "locked_lbc")
		assert.Equal(t, "rbtc", rbtcLabels["currency"])
		assert.Equal(t, "locked_lbc", rbtcLabels["type"])

		// Verify labels are correct for BTC metrics
		btcLabels := getGaugeVecLabels(metrics.AssetsMetrics, "btc", "rebalancing")
		assert.Equal(t, "btc", btcLabels["currency"])
		assert.Equal(t, "rebalancing", btcLabels["type"])
	})

	t.Run("should handle edge case values correctly (zero, large, complex decimals)", func(t *testing.T) {
		metrics, _ := createMetricsWithMock(t)

		// Create report with mixed edge case values:
		// - Zero values for some metrics
		// - Large values for others
		// - Complex decimal values
		report := reports.GetAssetReportResult{
			RbtcLockedLbc:      entities.NewWei(0),                                                                                              // 0 RBTC
			RbtcLockedForUsers: createWeiFromString("34598535894857007656"),                                                                     // 34.598535894857007656 RBTC
			RbtcWaitingRefund:  entities.NewWei(123456789000000),                                                                                // 0.000123456789 RBTC
			RbtcLiquidity:      entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 1000 RBTC
			RbtcWalletBalance:  createWeiFromString("7890123456789000000"),                                                                      // 7.890123456789 RBTC
			BtcLockedForUsers:  createWeiFromString("56789012345678900000"),                                                                     // 56.7890123456789 BTC
			BtcLiquidity:       entities.NewWei(0),                                                                                              // 0 BTC
			BtcWalletBalance:   entities.NewBigWei(big.NewInt(0).Mul(big.NewInt(2500), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))), // 2500 BTC
			BtcRebalancing:     entities.NewWei(987654321000000000),                                                                             // 0.987654321 BTC
		}

		metrics.UpdateAssetsFromReport(report)

		// Verify zero values
		assert.InDelta(t, 0.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "locked_lbc"), 0.0001)
		assert.InDelta(t, 0.0, getGaugeVecValue(metrics.AssetsMetrics, "btc", "liquidity"), 0.0001)

		// Verify complex decimal values
		assert.InDelta(t, 34.598535894857007656, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "locked_for_users"), 1e-15)
		assert.InDelta(t, 0.000123456789, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "waiting_refund"), 1e-15)
		assert.InDelta(t, 7.890123456789, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "wallet_balance"), 1e-12)
		assert.InDelta(t, 56.7890123456789, getGaugeVecValue(metrics.AssetsMetrics, "btc", "locked_for_users"), 1e-12)
		assert.InDelta(t, 0.987654321, getGaugeVecValue(metrics.AssetsMetrics, "btc", "rebalancing"), 1e-9)

		// Verify large values
		assert.InDelta(t, 1000.0, getGaugeVecValue(metrics.AssetsMetrics, "rbtc", "liquidity"), 0.0001)
		assert.InDelta(t, 2500.0, getGaugeVecValue(metrics.AssetsMetrics, "btc", "wallet_balance"), 0.0001)
	})
}

// Helper function to create test asset report with known values
func createTestAssetReport() reports.GetAssetReportResult {
	return reports.GetAssetReportResult{
		RbtcLockedLbc:      entities.NewBigWei(big.NewInt(1500000000000000000)), // 1.5 RBTC
		RbtcLockedForUsers: entities.NewBigWei(big.NewInt(2000000000000000000)), // 2.0 RBTC
		RbtcWaitingRefund:  entities.NewWei(500000000000000000),                 // 0.5 RBTC
		RbtcLiquidity:      entities.NewBigWei(big.NewInt(5000000000000000000)), // 5.0 RBTC
		RbtcWalletBalance:  entities.NewBigWei(big.NewInt(3000000000000000000)), // 3.0 RBTC
		BtcLockedForUsers:  entities.NewBigWei(big.NewInt(1800000000000000000)), // 1.8 BTC equivalent
		BtcLiquidity:       entities.NewBigWei(big.NewInt(4500000000000000000)), // 4.5 BTC equivalent
		BtcWalletBalance:   entities.NewBigWei(big.NewInt(2800000000000000000)), // 2.8 BTC equivalent
		BtcRebalancing:     entities.NewWei(100000000000000000),                 // 0.1 BTC equivalent
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
