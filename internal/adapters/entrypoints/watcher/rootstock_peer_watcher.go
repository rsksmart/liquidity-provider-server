package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher/monitoring"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type RootstockPeerWatcher struct {
	rpc                blockchain.Rpc
	peerAlertUseCase   NodePeerAlertUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	minPeers           uint64
	validationTimeout  time.Duration
	metrics            *monitoring.Metrics
	alertCooldown      time.Duration
	lastAlertTime      time.Time
}

func NewRootstockPeerWatcher(
	rpc blockchain.Rpc,
	peerAlertUseCase NodePeerAlertUseCase,
	ticker utils.Ticker,
	minPeers uint64,
	validationTimeout time.Duration,
	alertCooldown time.Duration,
	metrics *monitoring.Metrics,
) *RootstockPeerWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	return &RootstockPeerWatcher{
		rpc:                rpc,
		peerAlertUseCase:   peerAlertUseCase,
		ticker:             ticker,
		watcherStopChannel: watcherStopChannel,
		minPeers:           minPeers,
		validationTimeout:  validationTimeout,
		alertCooldown:      alertCooldown,
		metrics:            metrics,
	}
}

func (watcher *RootstockPeerWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *RootstockPeerWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.checkPeers()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *RootstockPeerWatcher) checkPeers() {
	ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
	defer cancel()
	currentPeers, err := watcher.rpc.Rsk.PeerCount(ctx)
	if err != nil {
		log.Error("RootstockPeerWatcher: error getting peer count: ", err)
		watcher.metrics.IncrementNodePeerCheckError(entities.NodeTypeRootstock)
		return
	}
	belowThreshold := currentPeers < watcher.minPeers
	watcher.metrics.UpdateNodePeerStatus(entities.NodeTypeRootstock, float64(currentPeers), float64(watcher.minPeers), belowThreshold)
	if !belowThreshold {
		return
	}
	log.Warnf("RootstockPeerWatcher: peer count %d is below minimum %d", currentPeers, watcher.minPeers)
	if time.Since(watcher.lastAlertTime) < watcher.alertCooldown {
		return
	}
	if alertErr := watcher.peerAlertUseCase.Run(ctx, entities.NodeTypeRootstock, int64(currentPeers), watcher.minPeers); alertErr != nil {
		log.Error("RootstockPeerWatcher: error sending low peer alert: ", alertErr)
	} else {
		watcher.lastAlertTime = time.Now()
		watcher.metrics.IncrementNodePeerAlert(entities.NodeTypeRootstock)
	}
}

func (watcher *RootstockPeerWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("RootstockPeerWatcher shut down")
}
