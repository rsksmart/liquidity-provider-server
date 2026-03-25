package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
)

type Metrics struct {
	PeginQuotesMetric             *prometheus.CounterVec
	PegoutQuotesMetric            *prometheus.CounterVec
	ServerInfoMetric              *prometheus.GaugeVec
	AssetsMetrics                 *prometheus.GaugeVec
	NodeReorgDepthMetric          *prometheus.GaugeVec
	NodeReorgMaxDepthMetric       *prometheus.GaugeVec
	NodeReorgAboveThresholdMetric *prometheus.GaugeVec
	NodeReorgCheckErrorsMetric    *prometheus.CounterVec
	NodeReorgAlertsMetric         *prometheus.CounterVec
}

func newNodeReorgMetrics() (
	depth, maxDepth, aboveThreshold *prometheus.GaugeVec,
	checkErrors, alerts *prometheus.CounterVec,
) {
	depth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_depth",
		Help: "Detected blockchain reorganization depth for the node",
	}, []string{"node"})
	maxDepth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_max_depth_threshold",
		Help: "Configured maximum reorganization depth before alerting",
	}, []string{"node"})
	aboveThreshold = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_above_threshold",
		Help: "Whether reorganization depth exceeds configured threshold (1=yes, 0=no)",
	}, []string{"node"})
	checkErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_reorg_check_errors_total",
		Help: "Total number of reorg check RPC errors",
	}, []string{"node"})
	alerts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_reorg_alerts_total",
		Help: "Total number of reorganization alerts sent",
	}, []string{"node"})
	return
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	reorgDepth, reorgMax, reorgAbove, reorgErrs, reorgAlerts := newNodeReorgMetrics()
	appMetrics := Metrics{
		PeginQuotesMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lps_pegin_quotes",
				Help: "Pegin quotes processed",
			},
			[]string{"state"},
		),
		PegoutQuotesMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lps_pegout_quotes",
				Help: "Pegout quotes processed",
			},
			[]string{"state"},
		),
		ServerInfoMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "lps_server_info",
				Help: "Server information",
			},
			[]string{"version", "commit"},
		),
		AssetsMetrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "lps_assets_balances",
				Help: "Liquidity provider asset balances and metrics (in BTC/RBTC units)",
			},
			[]string{"currency", "type"},
		),
		NodeReorgDepthMetric:          reorgDepth,
		NodeReorgMaxDepthMetric:       reorgMax,
		NodeReorgAboveThresholdMetric: reorgAbove,
		NodeReorgCheckErrorsMetric:    reorgErrs,
		NodeReorgAlertsMetric:         reorgAlerts,
	}

	reg.MustRegister(
		appMetrics.PegoutQuotesMetric,
		appMetrics.PeginQuotesMetric,
		appMetrics.ServerInfoMetric,
		appMetrics.AssetsMetrics,
		reorgDepth,
		reorgMax,
		reorgAbove,
		reorgErrs,
		reorgAlerts,
	)
	return &appMetrics
}

func (m *Metrics) UpdateNodeReorgStatus(node string, currentDepth float64, maxDepth float64, aboveThreshold bool) {
	m.NodeReorgDepthMetric.WithLabelValues(node).Set(currentDepth)
	m.NodeReorgMaxDepthMetric.WithLabelValues(node).Set(maxDepth)
	if aboveThreshold {
		m.NodeReorgAboveThresholdMetric.WithLabelValues(node).Set(1)
	} else {
		m.NodeReorgAboveThresholdMetric.WithLabelValues(node).Set(0)
	}
}

func (m *Metrics) IncrementNodeReorgCheckError(node string) {
	m.NodeReorgCheckErrorsMetric.WithLabelValues(node).Inc()
}

func (m *Metrics) IncrementNodeReorgAlert(node string) {
	m.NodeReorgAlertsMetric.WithLabelValues(node).Inc()
}

func (m *Metrics) UpdateAssetsFromReport(report reports.GetAssetsReportResult) {
	// RBTC metrics - Total
	m.AssetsMetrics.WithLabelValues("rbtc", "total").Set(weiToBtcFloat64(report.RbtcAssetReport.Total))

	// RBTC Location metrics
	m.AssetsMetrics.WithLabelValues("rbtc", "location_rsk_wallet").Set(weiToBtcFloat64(report.RbtcAssetReport.Location.RskWallet))
	m.AssetsMetrics.WithLabelValues("rbtc", "location_lbc").Set(weiToBtcFloat64(report.RbtcAssetReport.Location.Lbc))
	m.AssetsMetrics.WithLabelValues("rbtc", "location_federation").Set(weiToBtcFloat64(report.RbtcAssetReport.Location.Federation))

	// RBTC Allocation metrics
	m.AssetsMetrics.WithLabelValues("rbtc", "allocation_reserved_for_users").Set(weiToBtcFloat64(report.RbtcAssetReport.Allocation.ReservedForUsers))
	m.AssetsMetrics.WithLabelValues("rbtc", "allocation_waiting_refund").Set(weiToBtcFloat64(report.RbtcAssetReport.Allocation.WaitingForRefund))
	m.AssetsMetrics.WithLabelValues("rbtc", "allocation_available").Set(weiToBtcFloat64(report.RbtcAssetReport.Allocation.Available))

	// BTC metrics - Total
	m.AssetsMetrics.WithLabelValues("btc", "total").Set(weiToBtcFloat64(report.BtcAssetReport.Total))

	// BTC Location metrics
	m.AssetsMetrics.WithLabelValues("btc", "location_btc_wallet").Set(weiToBtcFloat64(report.BtcAssetReport.Location.BtcWallet))
	m.AssetsMetrics.WithLabelValues("btc", "location_federation").Set(weiToBtcFloat64(report.BtcAssetReport.Location.Federation))
	m.AssetsMetrics.WithLabelValues("btc", "location_rsk_wallet").Set(weiToBtcFloat64(report.BtcAssetReport.Location.RskWallet))
	m.AssetsMetrics.WithLabelValues("btc", "location_lbc").Set(weiToBtcFloat64(report.BtcAssetReport.Location.Lbc))

	// BTC Allocation metrics
	m.AssetsMetrics.WithLabelValues("btc", "allocation_reserved_for_users").Set(weiToBtcFloat64(report.BtcAssetReport.Allocation.ReservedForUsers))
	m.AssetsMetrics.WithLabelValues("btc", "allocation_waiting_refund").Set(weiToBtcFloat64(report.BtcAssetReport.Allocation.WaitingForRefund))
	m.AssetsMetrics.WithLabelValues("btc", "allocation_available").Set(weiToBtcFloat64(report.BtcAssetReport.Allocation.Available))
}

func weiToBtcFloat64(weiValue *entities.Wei) float64 {
	asRbtc := weiValue.ToRbtc()
	asFloat, _ := asRbtc.Float64()
	return asFloat
}
