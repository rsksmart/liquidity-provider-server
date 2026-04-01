package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type NodePeerCheckUseCase interface {
	Run(ctx context.Context, nodeType entities.NodeType) error
}

type BitcoinPeerWatcher struct {
	peerCheckUseCase   NodePeerCheckUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	validationTimeout  time.Duration
}

func NewBitcoinPeerWatcher(
	peerCheckUseCase NodePeerCheckUseCase,
	ticker utils.Ticker,
	validationTimeout time.Duration,
) *BitcoinPeerWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	return &BitcoinPeerWatcher{
		peerCheckUseCase:   peerCheckUseCase,
		ticker:             ticker,
		watcherStopChannel: watcherStopChannel,
		validationTimeout:  validationTimeout,
	}
}

func (watcher *BitcoinPeerWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *BitcoinPeerWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
			err := watcher.peerCheckUseCase.Run(ctx, entities.NodeTypeBitcoin)
			cancel()
			if err != nil {
				log.Error("BitcoinPeerWatcher: error running peer check: ", err)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *BitcoinPeerWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("BitcoinPeerWatcher shut down")
}
