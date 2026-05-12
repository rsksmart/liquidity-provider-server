package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PenalizationAlertWatcher struct {
	rpc                      blockchain.Rpc
	penalizationAlertUseCase *liquidity_provider.PenalizationAlertUseCase
	currentBlock             uint64
	currentBlockMutex        sync.RWMutex
	ticker                   utils.Ticker
	watcherStopChannel       chan bool
	validationTimeout        time.Duration
}

func NewPenalizationAlertWatcher(
	rpc blockchain.Rpc,
	penalizationAlertUseCase *liquidity_provider.PenalizationAlertUseCase,
	ticker utils.Ticker,
	validationTimeout time.Duration,
) *PenalizationAlertWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &PenalizationAlertWatcher{
		rpc:                      rpc,
		penalizationAlertUseCase: penalizationAlertUseCase,
		watcherStopChannel:       watcherStopChannel,
		ticker:                   ticker,
		currentBlockMutex:        sync.RWMutex{},
		validationTimeout:        validationTimeout,
	}
}

func (watcher *PenalizationAlertWatcher) Prepare(ctx context.Context) error {
	var err error
	var height uint64
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	if height, err = watcher.rpc.Rsk.GetHeight(ctx); err != nil {
		return err
	}
	watcher.currentBlock = height
	return nil
}

func (watcher *PenalizationAlertWatcher) Start() {
	var err error
	var height uint64
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.currentBlockMutex.Lock()
			ctx, cancel := context.WithTimeout(context.Background(), watcher.validationTimeout)
			if height, err = watcher.rpc.Rsk.GetHeight(ctx); err != nil {
				log.Error("Error checking penalization events inside watcher: ", err)
			} else {
				if err = watcher.penalizationAlertUseCase.Run(ctx, watcher.currentBlock, height); err == nil {
					watcher.currentBlock = height - 1
				}
			}
			cancel()
			watcher.currentBlockMutex.Unlock()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PenalizationAlertWatcher) GetCurrentBlock() uint64 {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}

func (watcher *PenalizationAlertWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("PenalizationAlertWatcher shut down")
}
