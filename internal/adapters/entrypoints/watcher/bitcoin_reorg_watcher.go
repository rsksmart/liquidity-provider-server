package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type NodeReorgCheckUseCase interface {
	Run(ctx context.Context, nodeType entities.NodeType) error
}

type BitcoinReorgWatcher struct {
	reorgCheckUseCase  NodeReorgCheckUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	validationTimeout  time.Duration
}

func NewBitcoinReorgWatcher(
	reorgCheckUseCase NodeReorgCheckUseCase,
	ticker utils.Ticker,
	validationTimeout time.Duration,
) *BitcoinReorgWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	return &BitcoinReorgWatcher{
		reorgCheckUseCase:  reorgCheckUseCase,
		ticker:             ticker,
		watcherStopChannel: watcherStopChannel,
		validationTimeout:  validationTimeout,
	}
}

func (watcher *BitcoinReorgWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *BitcoinReorgWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
			err := watcher.reorgCheckUseCase.Run(ctx, entities.NodeTypeBitcoin)
			cancel()
			if err != nil {
				log.Error("BitcoinReorgWatcher: error running reorg check: ", err)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *BitcoinReorgWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("BitcoinReorgWatcher shut down")
}
