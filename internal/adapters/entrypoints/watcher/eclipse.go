package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
	"time"
)

type EclipseCheckUseCase interface {
	Run(ctx context.Context, nodeType entities.NodeType) error
}

type EclipseWatcher struct {
	eclipseCheckUseCase EclipseCheckUseCase
	ticker              utils.Ticker
	watcherStopChannel  chan struct{}
	cooldownChannel     chan struct{}
	cooldownSeconds     uint64
	coolingDown         bool
	cooldownTimer       *time.Timer
	nodeType            entities.NodeType
}

func NewEclipseWatcher(
	eclipseCheckUseCase EclipseCheckUseCase,
	nodeType entities.NodeType,
	cooldownSeconds uint64,
	ticker utils.Ticker,
) *EclipseWatcher {
	watcherStopChannel := make(chan struct{}, 1)
	cooldownChannel := make(chan struct{}, 1)
	return &EclipseWatcher{
		eclipseCheckUseCase: eclipseCheckUseCase,
		ticker:              ticker,
		watcherStopChannel:  watcherStopChannel,
		cooldownChannel:     cooldownChannel,
		coolingDown:         false,
		cooldownTimer:       nil,
		nodeType:            nodeType,
		cooldownSeconds:     cooldownSeconds,
	}
}

func (watcher *EclipseWatcher) Prepare(ctx context.Context) error {
	return nil
}

func (watcher *EclipseWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			if !watcher.coolingDown {
				watcher.runEclipseCheck()
			}
		case <-watcher.cooldownChannel:
			watcher.coolingDown = false
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			if watcher.cooldownTimer != nil {
				watcher.cooldownTimer.Stop()
			}
			close(watcher.watcherStopChannel)
			close(watcher.cooldownChannel)
			break watcherLoop
		}
	}
}

func (watcher *EclipseWatcher) runEclipseCheck() {
	err := watcher.eclipseCheckUseCase.Run(context.Background(), watcher.nodeType)
	if errors.Is(err, w.NodeEclipseDetectedError) {
		watcher.coolingDown = true
		watcher.cooldownTimer = time.AfterFunc(
			time.Duration(watcher.cooldownSeconds)*time.Second,
			func() {
				if watcher.cooldownChannel != nil {
					watcher.cooldownChannel <- struct{}{}
				}
			},
		)
	} else if err != nil {
		log.Error("Error executing eclipse check: ", err)
	}
}

func (watcher *EclipseWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("EclipseWatcher shut down")
}
