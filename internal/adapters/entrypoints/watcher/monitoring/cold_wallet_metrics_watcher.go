package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	log "github.com/sirupsen/logrus"
)

type ColdWalletMetricsWatcher struct {
	appMetrics   *Metrics
	eventBus     entities.EventBus
	closeChannel chan struct{}
}

func NewColdWalletMetricsWatcher(
	appMetrics *Metrics,
	eventBus entities.EventBus,
) *ColdWalletMetricsWatcher {
	closeChannel := make(chan struct{}, 1)
	return &ColdWalletMetricsWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: closeChannel,
	}
}

func (watcher *ColdWalletMetricsWatcher) Prepare(ctx context.Context) error {
	return nil
}

// nolint: cyclop
func (watcher *ColdWalletMetricsWatcher) Start() {
	rbtcThresholdChannel := watcher.eventBus.Subscribe(cold_wallet.RbtcTransferredDueToThresholdEventId)
	btcThresholdChannel := watcher.eventBus.Subscribe(cold_wallet.BtcTransferredDueToThresholdEventId)
	rbtcTimeForcingChannel := watcher.eventBus.Subscribe(cold_wallet.RbtcTransferredDueToTimeForcingEventId)
	btcTimeForcingChannel := watcher.eventBus.Subscribe(cold_wallet.BtcTransferredDueToTimeForcingEventId)

metricLoop:
	for {
		select {
		case event := <-rbtcThresholdChannel:
			if typedEvent, ok := event.(cold_wallet.RbtcTransferredDueToThresholdEvent); ok {
				watcher.appMetrics.ColdWalletTransfersMetric.WithLabelValues(MetricLabelRbtc, MetricLabelThreshold).Inc()
				watcher.appMetrics.ColdWalletLastAmountMetric.WithLabelValues(MetricLabelRbtc).Set(typedEvent.Amount.ToRbtcFloat64())
			}
		case event := <-btcThresholdChannel:
			if typedEvent, ok := event.(cold_wallet.BtcTransferredDueToThresholdEvent); ok {
				watcher.appMetrics.ColdWalletTransfersMetric.WithLabelValues(MetricLabelBtc, MetricLabelThreshold).Inc()
				watcher.appMetrics.ColdWalletLastAmountMetric.WithLabelValues(MetricLabelBtc).Set(typedEvent.Amount.ToRbtcFloat64())
			}
		case event := <-rbtcTimeForcingChannel:
			if typedEvent, ok := event.(cold_wallet.RbtcTransferredDueToTimeForcingEvent); ok {
				watcher.appMetrics.ColdWalletTransfersMetric.WithLabelValues(MetricLabelRbtc, MetricLabelTimeForcing).Inc()
				watcher.appMetrics.ColdWalletLastAmountMetric.WithLabelValues(MetricLabelRbtc).Set(typedEvent.Amount.ToRbtcFloat64())
			}
		case event := <-btcTimeForcingChannel:
			if typedEvent, ok := event.(cold_wallet.BtcTransferredDueToTimeForcingEvent); ok {
				watcher.appMetrics.ColdWalletTransfersMetric.WithLabelValues(MetricLabelBtc, MetricLabelTimeForcing).Inc()
				watcher.appMetrics.ColdWalletLastAmountMetric.WithLabelValues(MetricLabelBtc).Set(typedEvent.Amount.ToRbtcFloat64())
			}
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			break metricLoop
		}
	}
}

func (watcher *ColdWalletMetricsWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.closeChannel <- struct{}{}
	closeChannel <- true
	log.Debug("Cold wallet metrics watcher shutdown completed")
}
