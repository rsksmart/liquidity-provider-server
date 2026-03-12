package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
)

type Metrics struct {
	PeginQuotesMetric          *prometheus.CounterVec
	PegoutQuotesMetric         *prometheus.CounterVec
	ServerInfoMetric           *prometheus.GaugeVec
	AssetsMetrics              *prometheus.GaugeVec
	ColdWalletTransfersMetric  *prometheus.CounterVec
	ColdWalletLastAmountMetric *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
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
	}

	reg.MustRegister(
		appMetrics.PegoutQuotesMetric,
		appMetrics.PeginQuotesMetric,
		appMetrics.ServerInfoMetric,
		appMetrics.AssetsMetrics,
		appMetrics.ColdWalletTransfersMetric,
		appMetrics.ColdWalletLastAmountMetric,
	)
	return &appMetrics
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
