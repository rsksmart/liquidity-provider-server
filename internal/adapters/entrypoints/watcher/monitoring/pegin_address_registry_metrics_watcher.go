package monitoring

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNilMetrics  = errors.New("monitoring watcher requires application metrics")
	ErrNilEventBus = errors.New("monitoring watcher requires an event bus")
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
	return watcher.validateDependencies()
}

func (watcher *PegInAddressRegistryMetricsWatcher) validateDependencies() error {
	if watcher.appMetrics == nil {
		return ErrNilMetrics
	}
	if watcher.eventBus == nil {
		return ErrNilEventBus
	}
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
