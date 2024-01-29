package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"time"
)

type PenalizationAlertWatcher struct {
	rskRpc                   blockchain.RootstockRpcServer
	penalizationAlertUseCase *liquidity_provider.PenalizationAlertUseCase
	currentBlock             uint64
	ticker                   *time.Ticker
	watcherStopChannel       chan bool
}

func NewPenalizationAlertWatcher(rskRpc blockchain.RootstockRpcServer, penalizationAlertUseCase *liquidity_provider.PenalizationAlertUseCase) *PenalizationAlertWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &PenalizationAlertWatcher{rskRpc: rskRpc, penalizationAlertUseCase: penalizationAlertUseCase, watcherStopChannel: watcherStopChannel}
}

func (watcher *PenalizationAlertWatcher) Prepare(ctx context.Context) error {
	var err error
	var height uint64
	if height, err = watcher.rskRpc.GetHeight(ctx); err != nil {
		return err
	}
	watcher.currentBlock = height
	return nil
}

func (watcher *PenalizationAlertWatcher) Start() {
	var cancel context.CancelFunc
	var ctx context.Context
	var err error
	var height uint64
	watcher.ticker = time.NewTicker(penalizationCheckInterval)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			ctx, cancel = context.WithTimeout(context.Background(), watcherValidationTimeout)
			if height, err = watcher.rskRpc.GetHeight(ctx); err != nil {
				log.Error("Error checking penalization events inside watcher: ", err)
			} else {
				if err = watcher.penalizationAlertUseCase.Run(ctx, watcher.currentBlock, height); err == nil {
					watcher.currentBlock = height - 1
				}
			}
			cancel()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PenalizationAlertWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("PenalizationAlertWatcher shut down")
}
