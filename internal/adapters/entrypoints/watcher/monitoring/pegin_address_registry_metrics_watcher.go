package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

type PegInAddressRegistryMetricsWatcher struct {
	appMetrics   *Metrics
	eventBus     entities.EventBus
	closeChannel chan struct{}
}

func NewPegInAddressRegistryMetricsWatcher(
	appMetrics *Metrics,
	eventBus entities.EventBus,
) *PegInAddressRegistryMetricsWatcher {
	return &PegInAddressRegistryMetricsWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInAddressRegistryMetricsWatcher) Prepare(context.Context) error {
	return nil
}

func (watcher *PegInAddressRegistryMetricsWatcher) Start() {
	mismatchEvents := watcher.eventBus.Subscribe(blockchain.PegInAddressRegistryRootMismatchEventId)
	resyncEvents := watcher.eventBus.Subscribe(blockchain.PegInAddressRegistryResyncStartedEventId)

	for {
		select {
		case <-mismatchEvents:
			watcher.appMetrics.IncrementPegInAddressRegistryRootMismatch()
		case <-resyncEvents:
			watcher.appMetrics.IncrementPegInAddressRegistryResync()
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			log.Debug("PegIn address registry metrics watcher shutdown completed")
			return
		}
	}
}

func (watcher *PegInAddressRegistryMetricsWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.closeChannel <- struct{}{}
	closeChannel <- true
}
