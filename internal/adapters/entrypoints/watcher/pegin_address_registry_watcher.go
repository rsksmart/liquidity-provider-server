package watcher

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type PegInWatcherUseCases struct {
	getWatchedUseCase   *w.GetWatchedRegisteredAddressesUseCase
	catchUpUseCase      *w.CatchUpRegisteredAddressesUseCase
	discoverUseCase     *w.DiscoverRegisteredAddressUseCase
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase
	recordErrorUseCase  *w.RecordRegisteredAddressWatchErrorUseCase
}

func NewPegInWatcherUseCases(
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase,
	catchUpUseCase *w.CatchUpRegisteredAddressesUseCase,
	discoverUseCase *w.DiscoverRegisteredAddressUseCase,
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase,
	recordErrorUseCase *w.RecordRegisteredAddressWatchErrorUseCase,
) *PegInWatcherUseCases {
	return &PegInWatcherUseCases{
		getWatchedUseCase:   getWatchedUseCase,
		catchUpUseCase:      catchUpUseCase,
		discoverUseCase:     discoverUseCase,
		markImportedUseCase: markImportedUseCase,
		recordErrorUseCase:  recordErrorUseCase,
	}
}

type PegInWatcher struct {
	getWatchedUseCase   *w.GetWatchedRegisteredAddressesUseCase
	catchUpUseCase      *w.CatchUpRegisteredAddressesUseCase
	discoverUseCase     *w.DiscoverRegisteredAddressUseCase
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase
	recordErrorUseCase  *w.RecordRegisteredAddressWatchErrorUseCase
	btcNetwork          blockchain.BitcoinNetwork
	wallet              blockchain.BitcoinWallet
	ticker              utils.Ticker
	watcherStopChannel  chan struct{}
	watches             []rootstock.PegInWatch
	stateMutex          sync.RWMutex
	reloadPending       bool
}

func NewPegInWatcher(
	useCases *PegInWatcherUseCases,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
) *PegInWatcher {
	return &PegInWatcher{
		getWatchedUseCase:   useCases.getWatchedUseCase,
		catchUpUseCase:      useCases.catchUpUseCase,
		discoverUseCase:     useCases.discoverUseCase,
		markImportedUseCase: useCases.markImportedUseCase,
		recordErrorUseCase:  useCases.recordErrorUseCase,
		btcNetwork:          btcNetwork,
		wallet:              wallet,
		ticker:              ticker,
		watcherStopChannel:  make(chan struct{}, 1),
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
	pending, err := watcher.retryDiscoveredWatches(ctx, []rootstock.PegInWatch{})
	if err == nil {
		err = watcher.rescanPendingImports(ctx, pending)
	}
	return errors.Join(catchUpErr, err)
}

func (watcher *PegInWatcher) reload(ctx context.Context) error {
	err := watcher.Prepare(ctx)
	watcher.reloadPending = err != nil
	return err
}

func (watcher *PegInWatcher) retryDiscoveredWatches(
	ctx context.Context,
	pending []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	watcher.stateMutex.RLock()
	watches := append([]rootstock.PegInWatch(nil), watcher.watches...)
	watcher.stateMutex.RUnlock()
	var err error
	for index := range watches {
		if watches[index].State != rootstock.PegInWatchDiscovered {
			continue
		}
		pending, err = watcher.discoverEvent(ctx, blockchain.NewAddressRegisteredFromWatchEntry(watches[index]), pending)
		if err != nil {
			return pending, err
		}
	}
	return pending, nil
}

func (watcher *PegInWatcher) discoverEvent(
	ctx context.Context,
	event blockchain.AddressRegistered,
	pending []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	result, err := watcher.discoverUseCase.Run(ctx, event)
	if err != nil {
		return pending, err
	}
	if result.Watch == nil {
		return pending, nil
	}
	watcher.rememberWatch(*result.Watch)
	if !result.NeedsRescan || slices.ContainsFunc(pending, func(watch rootstock.PegInWatch) bool {
		return watch.RskAddress == result.Watch.RskAddress
	}) {
		return pending, nil
	}
	return append(pending, *result.Watch), nil
}

func (watcher *PegInWatcher) rescanPendingImports(
	ctx context.Context,
	pending []rootstock.PegInWatch,
) error {
	if len(pending) == 0 {
		return nil
	}
	tip, err := watcher.btcNetwork.GetHeight()
	if err != nil {
		return watcher.recordPendingRescanError(ctx, pending, fmt.Errorf("get BTC height for registry rescan: %w", err))
	}
	fromHeight := max(tip.Int64()-peginAddressRegistryRescanDepthBlocks, 0)
	if _, err = watcher.wallet.RescanBlockchain(fromHeight); err != nil {
		return watcher.recordPendingRescanError(ctx, pending, fmt.Errorf("rescan PegIn addresses: %w", err))
	}
	for index := range pending {
		watch := pending[index]
		if err = watcher.markImportedUseCase.Run(ctx, &watch); err != nil {
			return err
		}
		pending[index] = watch
		watcher.rememberWatch(watch)
	}
	return nil
}

func (watcher *PegInWatcher) recordPendingRescanError(
	ctx context.Context,
	pending []rootstock.PegInWatch,
	rescanErr error,
) error {
	for index := range pending {
		watch := pending[index]
		if err := watcher.recordErrorUseCase.Run(ctx, &watch, rescanErr); err != nil {
			return err
		}
		pending[index] = watch
		watcher.rememberWatch(watch)
	}
	return nil
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
