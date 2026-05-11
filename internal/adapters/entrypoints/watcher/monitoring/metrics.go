package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
)

type Metrics struct {
	PeginQuotesMetric             *prometheus.CounterVec
	PegoutQuotesMetric            *prometheus.CounterVec
	ServerInfoMetric              *prometheus.GaugeVec
	AssetsMetrics                 *prometheus.GaugeVec
	ColdWalletTransfersMetric     *prometheus.CounterVec
	ColdWalletLastAmountMetric    *prometheus.GaugeVec
	NodeReorgDepthMetric          *prometheus.GaugeVec
	NodeReorgMaxDepthMetric       *prometheus.GaugeVec
	NodeReorgAboveThresholdMetric *prometheus.GaugeVec
	NodeReorgCheckErrorsMetric    *prometheus.CounterVec
	NodeReorgAlertsMetric         *prometheus.CounterVec
	NodePeerCountMetric           *prometheus.GaugeVec
	NodePeerMinThresholdMetric    *prometheus.GaugeVec
	NodePeerBelowThreshold        *prometheus.GaugeVec
	NodePeerCheckErrors           *prometheus.CounterVec
	NodePeerAlerts                *prometheus.CounterVec
}

type nodeReorgMetrics struct {
	Depth          *prometheus.GaugeVec
	MaxDepth       *prometheus.GaugeVec
	AboveThreshold *prometheus.GaugeVec
	CheckErrors    *prometheus.CounterVec
	Alerts         *prometheus.CounterVec
}

func newNodeReorgMetrics() nodeReorgMetrics {
	depth := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_depth",
		Help: "Detected blockchain reorganization depth for the node",
	}, []string{"node"})
	maxDepth := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_max_depth_threshold",
		Help: "Configured maximum reorganization depth before alerting",
	}, []string{"node"})
	aboveThreshold := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_reorg_above_threshold",
		Help: "Whether reorganization depth exceeds configured threshold (1=yes, 0=no)",
	}, []string{"node"})
	checkErrors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_reorg_check_errors_total",
		Help: "Total number of reorg check RPC errors",
	}, []string{"node"})
	alerts := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_reorg_alerts_total",
		Help: "Total number of reorganization alerts sent",
	}, []string{"node"})
	return nodeReorgMetrics{
		Depth:          depth,
		MaxDepth:       maxDepth,
		AboveThreshold: aboveThreshold,
		CheckErrors:    checkErrors,
		Alerts:         alerts,
	}
}

func newNodePeerMetrics() (
	peerCount, minThreshold, belowThreshold *prometheus.GaugeVec,
	checkErrors, alerts *prometheus.CounterVec,
) {
	peerCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_peer_count",
		Help: "Current number of peers connected to the node",
	}, []string{"node"})
	minThreshold = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_peer_min_threshold",
		Help: "Configured minimum peer threshold for the node",
	}, []string{"node"})
	belowThreshold = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lps_node_peer_below_threshold",
		Help: "Whether the node peer count is below the configured threshold (1=below, 0=ok)",
	}, []string{"node"})
	checkErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_peer_check_errors_total",
		Help: "Total number of peer check RPC errors",
	}, []string{"node"})
	alerts = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "lps_node_peer_alerts_total",
		Help: "Total number of low peer count alerts sent",
	}, []string{"node"})
	return
}

// nolint:funlen
func NewMetrics(reg prometheus.Registerer) *Metrics {
	reorg := newNodeReorgMetrics()
	peerCount, minThreshold, belowThreshold, checkErrors, peerAlerts := newNodePeerMetrics()
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
		ColdWalletTransfersMetric: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "lps_cold_wallet_transfers",
				Help: "Cold wallet transfers executed by reason",
			},
			[]string{"currency", "reason"},
		),
		ColdWalletLastAmountMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "lps_cold_wallet_last_transfer_amount",
				Help: "Amount transferred in last cold wallet transfer (in BTC/RBTC units)",
			},
			[]string{"currency"},
		),
		NodeReorgDepthMetric:          reorg.Depth,
		NodeReorgMaxDepthMetric:       reorg.MaxDepth,
		NodeReorgAboveThresholdMetric: reorg.AboveThreshold,
		NodeReorgCheckErrorsMetric:    reorg.CheckErrors,
		NodeReorgAlertsMetric:         reorg.Alerts,
		NodePeerCountMetric:           peerCount,
		NodePeerMinThresholdMetric:    minThreshold,
		NodePeerBelowThreshold:        belowThreshold,
		NodePeerCheckErrors:           checkErrors,
		NodePeerAlerts:                peerAlerts,
	}

	reg.MustRegister(
		appMetrics.PegoutQuotesMetric,
		appMetrics.PeginQuotesMetric,
		appMetrics.ServerInfoMetric,
		appMetrics.AssetsMetrics,
		appMetrics.ColdWalletTransfersMetric,
		appMetrics.ColdWalletLastAmountMetric,
		reorg.Depth,
		reorg.MaxDepth,
		reorg.AboveThreshold,
		reorg.CheckErrors,
		reorg.Alerts,
		peerCount, minThreshold, belowThreshold, checkErrors, peerAlerts,
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
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelTotal).Set(report.RbtcAssetReport.Total.ToRbtcFloat64())

	// RBTC Location metrics
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelLocationRskWallet).Set(report.RbtcAssetReport.Location.RskWallet.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelLocationLbc).Set(report.RbtcAssetReport.Location.Lbc.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelLocationFederation).Set(report.RbtcAssetReport.Location.Federation.ToRbtcFloat64())

	// RBTC Allocation metrics
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelAllocationReservedForUsers).Set(report.RbtcAssetReport.Allocation.ReservedForUsers.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelAllocationWaitingRefund).Set(report.RbtcAssetReport.Allocation.WaitingForRefund.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelRbtc, MetricLabelAllocationAvailable).Set(report.RbtcAssetReport.Allocation.Available.ToRbtcFloat64())

	// BTC metrics - Total
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelTotal).Set(report.BtcAssetReport.Total.ToRbtcFloat64())

	// BTC Location metrics
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelLocationBtcWallet).Set(report.BtcAssetReport.Location.BtcWallet.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelLocationFederation).Set(report.BtcAssetReport.Location.Federation.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelLocationRskWallet).Set(report.BtcAssetReport.Location.RskWallet.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelLocationLbc).Set(report.BtcAssetReport.Location.Lbc.ToRbtcFloat64())

	// BTC Allocation metrics
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelAllocationReservedForUsers).Set(report.BtcAssetReport.Allocation.ReservedForUsers.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelAllocationWaitingRefund).Set(report.BtcAssetReport.Allocation.WaitingForRefund.ToRbtcFloat64())
	m.AssetsMetrics.WithLabelValues(MetricLabelBtc, MetricLabelAllocationAvailable).Set(report.BtcAssetReport.Allocation.Available.ToRbtcFloat64())
}

func (m *Metrics) UpdateNodePeerStatus(node string, currentPeers float64, minPeers float64, belowThreshold bool) {
	m.NodePeerCountMetric.WithLabelValues(node).Set(currentPeers)
	m.NodePeerMinThresholdMetric.WithLabelValues(node).Set(minPeers)
	if belowThreshold {
		m.NodePeerBelowThreshold.WithLabelValues(node).Set(1)
	} else {
		m.NodePeerBelowThreshold.WithLabelValues(node).Set(0)
	}
}

func (m *Metrics) IncrementNodePeerCheckError(node string) {
	m.NodePeerCheckErrors.WithLabelValues(node).Inc()
}

func (m *Metrics) IncrementNodePeerAlert(node string) {
	m.NodePeerAlerts.WithLabelValues(node).Inc()
}
