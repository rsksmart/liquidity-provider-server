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

type NodePeerAlertUseCase interface {
	Run(ctx context.Context, nodeType entities.NodeType, currentPeers int64, minPeers uint64) error
}

type BitcoinPeerWatcher struct {
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

func NewBitcoinPeerWatcher(
	rpc blockchain.Rpc,
	peerAlertUseCase NodePeerAlertUseCase,
	ticker utils.Ticker,
	minPeers uint64,
	validationTimeout time.Duration,
	alertCooldown time.Duration,
	metrics *monitoring.Metrics,
) *BitcoinPeerWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	return &BitcoinPeerWatcher{
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

func (watcher *BitcoinPeerWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *BitcoinPeerWatcher) Start() {
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

func (watcher *BitcoinPeerWatcher) checkPeers() {
	ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
	defer cancel()
	currentPeers, err := watcher.rpc.Btc.GetConnectionCount()
	if err != nil {
		log.Error("BitcoinPeerWatcher: error getting connection count: ", err)
		watcher.metrics.IncrementNodePeerCheckError(entities.NodeTypeBitcoin)
		return
	}
	belowThreshold := uint64(currentPeers) < watcher.minPeers
	watcher.metrics.UpdateNodePeerStatus(entities.NodeTypeBitcoin, float64(currentPeers), float64(watcher.minPeers), belowThreshold)
	if !belowThreshold {
		return
	}
	log.Warnf("BitcoinPeerWatcher: peer count %d is below minimum %d", currentPeers, watcher.minPeers)
	if time.Since(watcher.lastAlertTime) < watcher.alertCooldown {
		return
	}
	if alertErr := watcher.peerAlertUseCase.Run(ctx, entities.NodeTypeBitcoin, currentPeers, watcher.minPeers); alertErr != nil {
		log.Error("BitcoinPeerWatcher: error sending low peer alert: ", alertErr)
	} else {
		watcher.lastAlertTime = time.Now()
		watcher.metrics.IncrementNodePeerAlert(entities.NodeTypeBitcoin)
	}
}

func (watcher *BitcoinPeerWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("BitcoinPeerWatcher shut down")
}
