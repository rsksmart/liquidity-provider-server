package monitoring

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

type PeerMetricsWatcher struct {
	appMetrics   *Metrics
	eventBus     entities.EventBus
	closeChannel chan struct{}
}

func NewPeerMetricsWatcher(appMetrics *Metrics, eventBus entities.EventBus) *PeerMetricsWatcher {
	closeChannel := make(chan struct{}, 1)
	return &PeerMetricsWatcher{
		appMetrics:   appMetrics,
		eventBus:     eventBus,
		closeChannel: closeChannel,
	}
}

func (watcher *PeerMetricsWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *PeerMetricsWatcher) Start() {
	peerCheckChannel := watcher.eventBus.Subscribe(blockchain.NodePeerCheckEventId)
	peerCheckErrorChannel := watcher.eventBus.Subscribe(blockchain.NodePeerCheckErrorEventId)
	peerAlertSentChannel := watcher.eventBus.Subscribe(blockchain.NodePeerAlertSentEventId)

metricLoop:
	for {
		select {
		case event := <-peerCheckChannel:
			if peerEvent, ok := event.(blockchain.NodePeerCheckEvent); ok {
				watcher.appMetrics.UpdateNodePeerStatus(
					string(peerEvent.NodeType),
					float64(peerEvent.CurrentPeers),
					float64(peerEvent.MinPeers),
					peerEvent.BelowThreshold,
				)
			}
		case event := <-peerCheckErrorChannel:
			if peerErrorEvent, ok := event.(blockchain.NodePeerCheckErrorEvent); ok {
				watcher.appMetrics.IncrementNodePeerCheckError(string(peerErrorEvent.NodeType))
			}
		case event := <-peerAlertSentChannel:
			if peerAlertEvent, ok := event.(blockchain.NodePeerAlertSentEvent); ok {
				watcher.appMetrics.IncrementNodePeerAlert(string(peerAlertEvent.NodeType))
			}
		case <-watcher.closeChannel:
			close(watcher.closeChannel)
			break metricLoop
		}
	}
}

func (watcher *PeerMetricsWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.closeChannel <- struct{}{}
	closeChannel <- true
	log.Debug("Peer metrics watcher shutdown completed")
}
