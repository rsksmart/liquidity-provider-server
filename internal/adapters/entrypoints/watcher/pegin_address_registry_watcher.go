package watcher

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type PegInWatcherUseCases struct {
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase
	catchUpUseCase    *w.CatchUpRegisteredAddressesUseCase
	importUseCase     *w.ImportRegisteredAddressesUseCase
}

func NewPegInWatcherUseCases(
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase,
	catchUpUseCase *w.CatchUpRegisteredAddressesUseCase,
	importUseCase *w.ImportRegisteredAddressesUseCase,
) *PegInWatcherUseCases {
	return &PegInWatcherUseCases{
		getWatchedUseCase: getWatchedUseCase,
		catchUpUseCase:    catchUpUseCase,
		importUseCase:     importUseCase,
	}
}

type PegInWatcher struct {
	getWatchedUseCase  *w.GetWatchedRegisteredAddressesUseCase
	catchUpUseCase     *w.CatchUpRegisteredAddressesUseCase
	importUseCase      *w.ImportRegisteredAddressesUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	watches            []rootstock.PegInWatch
	stateMutex         sync.RWMutex
	reloadPending      bool
}

func NewPegInWatcher(
	useCases *PegInWatcherUseCases,
	ticker utils.Ticker,
) *PegInWatcher {
	return &PegInWatcher{
		getWatchedUseCase:  useCases.getWatchedUseCase,
		catchUpUseCase:     useCases.catchUpUseCase,
		importUseCase:      useCases.importUseCase,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInWatcher) Prepare(ctx context.Context) error {
	watches, err := watcher.getWatchedUseCase.Run(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry watch set: %w", err)
	}

	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	watcher.watches = watches
	return nil
}

func (watcher *PegInWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			if err := watcher.scan(context.Background()); err != nil {
				log.Errorf("PegIn address registry watcher scan failed: %v", err)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegInWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("PegInWatcher shut down")
}

func (watcher *PegInWatcher) scan(ctx context.Context) error {
	result, catchUpErr := watcher.catchUpUseCase.Run(ctx)
	if result.Rebuilt || watcher.reloadPending {
		if err := watcher.reload(ctx); err != nil {
			return errors.Join(catchUpErr, err)
		}
	} else {
		for _, watch := range result.Watches {
			watcher.rememberWatch(watch)
		}
	}
	watcher.stateMutex.RLock()
	watches := append([]rootstock.PegInWatch(nil), watcher.watches...)
	watcher.stateMutex.RUnlock()
	importResult, importErr := watcher.importUseCase.Run(ctx, watches)
	for _, watch := range importResult.Watches {
		watcher.rememberWatch(watch)
	}
	return errors.Join(catchUpErr, importErr)
}

func (watcher *PegInWatcher) reload(ctx context.Context) error {
	err := watcher.Prepare(ctx)
	watcher.reloadPending = err != nil
	return err
}

func (watcher *PegInWatcher) rememberWatch(watch rootstock.PegInWatch) {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for index := range watcher.watches {
		if watcher.watches[index].RskAddress == watch.RskAddress {
			watcher.watches[index] = watch
			return
		}
	}
	watcher.watches = append(watcher.watches, watch)
}
