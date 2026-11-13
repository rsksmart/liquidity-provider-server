package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type PegInAddressRegistryMetricsWatcher struct {
	appMetrics     *Metrics
	eventBus       entities.EventBus
	mismatchEvents <-chan entities.Event
	resyncEvents   <-chan entities.Event
	closeChannel   chan struct{}
}

func NewPegInAddressRegistryMetricsWatcher(appMetrics *Metrics, eventBus entities.EventBus) *PegInAddressRegistryMetricsWatcher {
	return &PegInAddressRegistryMetricsWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInAddressRegistryMetricsWatcher) Prepare(context.Context) error {
	watcher.mismatchEvents = watcher.eventBus.Subscribe(blockchain.PegInAddressRegistryRootMismatchEventId)
	watcher.resyncEvents = watcher.eventBus.Subscribe(blockchain.PegInAddressRegistryResyncEventId)
	return nil
}

func (watcher *PegInAddressRegistryMetricsWatcher) Start() {
	for {
		select {
		case event := <-watcher.mismatchEvents:
			if _, isMismatch := event.(blockchain.PegInAddressRegistryRootMismatchEvent); isMismatch {
				watcher.appMetrics.IncrementPegInAddressRegistryRootMismatch()
			}
		case event := <-watcher.resyncEvents:
			if _, isResync := event.(blockchain.PegInAddressRegistryResyncEvent); isResync {
				watcher.appMetrics.IncrementPegInAddressRegistryResync()
			}
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			return
		}
	}
}

func (watcher *PegInAddressRegistryMetricsWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.closeChannel <- struct{}{}
	closeChannel <- true
}
