package monitoring

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

type Metrics struct {
	PeginQuotesMetric  *prometheus.CounterVec
	PegoutQuotesMetric *prometheus.CounterVec
	ServerInfoMetric   *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) Metrics {
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
	}
	reg.MustRegister(
		appMetrics.PegoutQuotesMetric,
		appMetrics.PeginQuotesMetric,
		appMetrics.ServerInfoMetric,
	)
	return appMetrics
}

type MetricWatcher struct {
	appMetrics   Metrics
	serverInfo   *liquidity_provider.ServerInfoUseCase
	eventBus     entities.EventBus
	closeChannel chan struct{}
}

func NewMetricWatcher(
	appMetrics Metrics,
	eventBus entities.EventBus,
	serverInfo *liquidity_provider.ServerInfoUseCase,
) *MetricWatcher {
	closeChannel := make(chan struct{}, 1)
	return &MetricWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: closeChannel,
		serverInfo:   serverInfo,
	}
}

func (watcher *MetricWatcher) Prepare(ctx context.Context) error {
	var info lp.ServerInfo
	var err error
	info, err = watcher.serverInfo.Run()
	if err != nil {
		info = lp.ServerInfo{
			Version:  "Not provided",
			Revision: "Not provided",
		}
	}
	watcher.appMetrics.ServerInfoMetric.WithLabelValues(info.Version, info.Revision).Set(1)
	return nil
}

func (watcher *MetricWatcher) Start() {
	acceptedPegoutChannel := watcher.eventBus.Subscribe(quote.AcceptedPegoutQuoteEventId)
	sendPegoutChannel := watcher.eventBus.Subscribe(quote.PegoutBtcSentEventId)
	pegoutRefundChannel := watcher.eventBus.Subscribe(quote.PegoutQuoteCompletedEventId)
	acceptedPeginChannel := watcher.eventBus.Subscribe(quote.AcceptedPeginQuoteEventId)
	callForUserChannel := watcher.eventBus.Subscribe(quote.CallForUserCompletedEventId)
	registerPeginChannel := watcher.eventBus.Subscribe(quote.RegisterPeginCompletedEventId)

metricLoop:
	for {
		select {
		case <-acceptedPegoutChannel:
			watcher.appMetrics.PegoutQuotesMetric.WithLabelValues("accepted").Inc()
		case <-sendPegoutChannel:
			watcher.appMetrics.PegoutQuotesMetric.WithLabelValues("paid").Inc()
		case <-pegoutRefundChannel:
			watcher.appMetrics.PegoutQuotesMetric.WithLabelValues("lp_refunded").Inc()
		case <-acceptedPeginChannel:
			watcher.appMetrics.PeginQuotesMetric.WithLabelValues("accepted").Inc()
		case <-callForUserChannel:
			watcher.appMetrics.PeginQuotesMetric.WithLabelValues("paid").Inc()
		case <-registerPeginChannel:
			watcher.appMetrics.PeginQuotesMetric.WithLabelValues("lp_refunded").Inc()
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			break metricLoop
		}
	}
}

func (watcher *MetricWatcher) Shutdown(closeChannel chan<- bool) {
	<-watcher.closeChannel
	closeChannel <- true
	log.Debug("Metrics watcher shutdown completed")
}
