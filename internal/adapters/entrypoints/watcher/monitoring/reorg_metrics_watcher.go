package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

type ReorgMetricsWatcher struct {
	appMetrics   *Metrics
	eventBus     entities.EventBus
	closeChannel chan struct{}
}

func NewReorgMetricsWatcher(appMetrics *Metrics, eventBus entities.EventBus) *ReorgMetricsWatcher {
	closeChannel := make(chan struct{}, 1)
	return &ReorgMetricsWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: closeChannel,
	}
}

func (watcher *ReorgMetricsWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *ReorgMetricsWatcher) Start() {
	checkCh := watcher.eventBus.Subscribe(blockchain.NodeReorgCheckEventId)
	errCh := watcher.eventBus.Subscribe(blockchain.NodeReorgCheckErrorEventId)
	alertCh := watcher.eventBus.Subscribe(blockchain.NodeReorgAlertSentEventId)

metricLoop:
	for {
		select {
		case event := <-checkCh:
			if ev, ok := event.(blockchain.NodeReorgCheckEvent); ok {
				watcher.appMetrics.UpdateNodeReorgStatus(
					string(ev.NodeType),
					float64(ev.CurrentDepth),
					float64(ev.MaxAllowedDepth),
					ev.AboveThreshold,
				)
			}
		case event := <-errCh:
			if ev, ok := event.(blockchain.NodeReorgCheckErrorEvent); ok {
				watcher.appMetrics.IncrementNodeReorgCheckError(string(ev.NodeType))
			}
		case event := <-alertCh:
			if ev, ok := event.(blockchain.NodeReorgAlertSentEvent); ok {
				watcher.appMetrics.IncrementNodeReorgAlert(string(ev.NodeType))
			}
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			break metricLoop
		}
	}
}

func (watcher *ReorgMetricsWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.closeChannel <- struct{}{}
	closeChannel <- true
	log.Debug("Reorg metrics watcher shutdown completed")
}
