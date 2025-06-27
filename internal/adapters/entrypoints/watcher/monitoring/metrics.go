package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
)

type Metrics struct {
	PeginQuotesMetric  *prometheus.CounterVec
	PegoutQuotesMetric *prometheus.CounterVec
	ServerInfoMetric   *prometheus.GaugeVec
	AssetsMetrics      *prometheus.GaugeVec
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
	}

	reg.MustRegister(
		appMetrics.PegoutQuotesMetric,
		appMetrics.PeginQuotesMetric,
		appMetrics.ServerInfoMetric,
		appMetrics.AssetsMetrics,
	)
	return &appMetrics
}

func (m *Metrics) UpdateAssetsFromReport(report reports.GetAssetReportResult) {
	// RBTC metrics
	m.AssetsMetrics.WithLabelValues("rbtc", "locked_lbc").Set(weiToBtcFloat64(report.RbtcLockedLbc))
	m.AssetsMetrics.WithLabelValues("rbtc", "locked_for_users").Set(weiToBtcFloat64(report.RbtcLockedForUsers))
	m.AssetsMetrics.WithLabelValues("rbtc", "waiting_refund").Set(weiToBtcFloat64(report.RbtcWaitingRefund))
	m.AssetsMetrics.WithLabelValues("rbtc", "liquidity").Set(weiToBtcFloat64(report.RbtcLiquidity))
	m.AssetsMetrics.WithLabelValues("rbtc", "wallet_balance").Set(weiToBtcFloat64(report.RbtcWalletBalance))

	// BTC metrics
	m.AssetsMetrics.WithLabelValues("btc", "locked_for_users").Set(weiToBtcFloat64(report.BtcLockedForUsers))
	m.AssetsMetrics.WithLabelValues("btc", "liquidity").Set(weiToBtcFloat64(report.BtcLiquidity))
	m.AssetsMetrics.WithLabelValues("btc", "wallet_balance").Set(weiToBtcFloat64(report.BtcWalletBalance))
	m.AssetsMetrics.WithLabelValues("btc", "rebalancing").Set(weiToBtcFloat64(report.BtcRebalancing))
}

func weiToBtcFloat64(weiValue *entities.Wei) float64 {
	asRbtc := weiValue.ToRbtc()
	asFloat, _ := asRbtc.Float64()
	return asFloat
}
