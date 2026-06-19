package watcher

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type RootstockReorgWatcher struct {
	reorgCheckUseCase  NodeReorgCheckUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	validationTimeout  time.Duration
}

func NewRootstockReorgWatcher(
	reorgCheckUseCase NodeReorgCheckUseCase,
	ticker utils.Ticker,
	validationTimeout time.Duration,
) *RootstockReorgWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	return &RootstockReorgWatcher{
		reorgCheckUseCase:  reorgCheckUseCase,
		ticker:             ticker,
		watcherStopChannel: watcherStopChannel,
		validationTimeout:  validationTimeout,
	}
}

func (watcher *RootstockReorgWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *RootstockReorgWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
			err := watcher.reorgCheckUseCase.Run(ctx, entities.NodeTypeRootstock)
			cancel()
			if err != nil {
				log.Error("RootstockReorgWatcher: error running reorg check: ", err)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *RootstockReorgWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("RootstockReorgWatcher shut down")
}
