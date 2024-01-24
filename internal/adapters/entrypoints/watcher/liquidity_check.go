package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"time"
)

type LiquidityCheckWatcher struct {
	checkLiquidityUseCase *liquidity_provider.CheckLiquidityUseCase
	watcherStopChannel    chan bool
	ticker                *time.Ticker
}

func NewLiquidityCheckWatcher(checkLiquidityUseCase *liquidity_provider.CheckLiquidityUseCase) *LiquidityCheckWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &LiquidityCheckWatcher{checkLiquidityUseCase: checkLiquidityUseCase, watcherStopChannel: watcherStopChannel}
}

func (watcher *LiquidityCheckWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("PeginBridgeWatcher shut down")
}

func (watcher *LiquidityCheckWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *LiquidityCheckWatcher) Start() {
	var ctx context.Context
	var cancel context.CancelFunc
	watcher.ticker = time.NewTicker(liquidityCheckInterval)

watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
			if err := watcher.checkLiquidityUseCase.Run(ctx); err != nil {
				log.Error("Error checking liquidity inside watcher: ", err)
			}
			cancel()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}
