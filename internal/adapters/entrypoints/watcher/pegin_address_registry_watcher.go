package watcher

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type PegInWatcherUseCases struct {
	getWatchedUseCase   *w.GetWatchedRegisteredAddressesUseCase
	getCursorUseCase    *w.GetRegistryWatchCursorUseCase
	setCursorUseCase    *w.SetRegistryWatchCursorUseCase
	discoverUseCase     *w.DiscoverRegisteredAddressUseCase
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase
	recordErrorUseCase  *w.RecordRegisteredAddressWatchErrorUseCase
}

func NewPegInWatcherUseCases(
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase,
	getCursorUseCase *w.GetRegistryWatchCursorUseCase,
	setCursorUseCase *w.SetRegistryWatchCursorUseCase,
	discoverUseCase *w.DiscoverRegisteredAddressUseCase,
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase,
	recordErrorUseCase *w.RecordRegisteredAddressWatchErrorUseCase,
) *PegInWatcherUseCases {
	return &PegInWatcherUseCases{
		getWatchedUseCase:   getWatchedUseCase,
		getCursorUseCase:    getCursorUseCase,
		setCursorUseCase:    setCursorUseCase,
		discoverUseCase:     discoverUseCase,
		markImportedUseCase: markImportedUseCase,
		recordErrorUseCase:  recordErrorUseCase,
	}
}

type PegInWatcher struct {
	getWatchedUseCase   *w.GetWatchedRegisteredAddressesUseCase
	getCursorUseCase    *w.GetRegistryWatchCursorUseCase
	setCursorUseCase    *w.SetRegistryWatchCursorUseCase
	discoverUseCase     *w.DiscoverRegisteredAddressUseCase
	markImportedUseCase *w.MarkRegisteredAddressImportedUseCase
	recordErrorUseCase  *w.RecordRegisteredAddressWatchErrorUseCase
	registry            blockchain.PegInAddressRegistryContract
	rskRpc              blockchain.RootstockRpcServer
	btcNetwork          blockchain.BitcoinNetwork
	wallet              blockchain.BitcoinWallet
	ticker              utils.Ticker
	watcherStopChannel  chan struct{}
	startBlock          uint64
	pageSize            uint64
	finalityDepth       uint64
	lastScannedBlock    uint64
	hasCursor           bool
	watches             []rootstock.PegInWatch
	stateMutex          sync.RWMutex
}

func NewPegInWatcher(
	useCases *PegInWatcherUseCases,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *PegInWatcher {
	return &PegInWatcher{
		getWatchedUseCase:   useCases.getWatchedUseCase,
		getCursorUseCase:    useCases.getCursorUseCase,
		setCursorUseCase:    useCases.setCursorUseCase,
		discoverUseCase:     useCases.discoverUseCase,
		markImportedUseCase: useCases.markImportedUseCase,
		recordErrorUseCase:  useCases.recordErrorUseCase,
		registry:            registry,
		rskRpc:              rskRpc,
		btcNetwork:          btcNetwork,
		wallet:              wallet,
		ticker:              ticker,
		watcherStopChannel:  make(chan struct{}, 1),
		startBlock:          startBlock,
		pageSize:            pageSize,
		finalityDepth:       finalityDepth,
	}
}

func (watcher *PegInWatcher) Prepare(ctx context.Context) error {
	watches, err := watcher.getWatchedUseCase.Run(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry watch set: %w", err)
	}
	lastScannedBlock, found, err := watcher.getCursorUseCase.Run(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry scan cursor: %w", err)
	}

	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	watcher.watches = watches
	watcher.lastScannedBlock = lastScannedBlock
	watcher.hasCursor = found
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
	pending := []rootstock.PegInWatch{}
	pending, err := watcher.retryDiscoveredWatches(ctx, pending)
	if err != nil {
		return err
	}
	head, err := watcher.rskRpc.GetHeight(ctx)
	if err != nil {
		return fmt.Errorf("get RSK height: %w", err)
	}
	if head < watcher.finalityDepth {
		return watcher.rescanPendingImports(ctx, pending)
	}

	watcher.stateMutex.RLock()
	fromBlock, toBlock := watcher.nextRange(head - watcher.finalityDepth)
	watcher.stateMutex.RUnlock()
	if fromBlock > toBlock {
		return watcher.rescanPendingImports(ctx, pending)
	}

	events, err := watcher.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err)
	}
	pending, err = watcher.processEvents(ctx, events, pending)
	if err != nil {
		return err
	}
	if err = watcher.rescanPendingImports(ctx, pending); err != nil {
		return err
	}
	if err = watcher.setCursorUseCase.Run(ctx, toBlock); err != nil {
		return fmt.Errorf("persist PegIn address registry scan cursor: %w", err)
	}

	watcher.stateMutex.Lock()
	watcher.lastScannedBlock = toBlock
	watcher.hasCursor = true
	watcher.stateMutex.Unlock()
	return nil
}

func (watcher *PegInWatcher) processEvents(
	ctx context.Context,
	events []blockchain.AddressRegistered,
	pending []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	sort.Slice(events, func(firstIndex, secondIndex int) bool {
		if events[firstIndex].BlockNumber == events[secondIndex].BlockNumber {
			return events[firstIndex].LogIndex < events[secondIndex].LogIndex
		}
		return events[firstIndex].BlockNumber < events[secondIndex].BlockNumber
	})

	var err error
	for _, event := range events {
		pending, err = watcher.discoverEvent(ctx, event, pending)
		if err != nil {
			return pending, err
		}
	}
	return pending, nil
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

func (watcher *PegInWatcher) nextRange(finalizedHead uint64) (uint64, uint64) {
	// This is the first scan
	if !watcher.hasCursor {
		if watcher.startBlock > finalizedHead ||
			watcher.pageSize-1 > finalizedHead-watcher.startBlock {
			return watcher.startBlock, finalizedHead
		}
		return watcher.startBlock, watcher.startBlock + watcher.pageSize - 1
	}

	fromBlock := watcher.startBlock
	if watcher.lastScannedBlock+1 > watcher.finalityDepth {
		fromBlock = max(watcher.startBlock, watcher.lastScannedBlock+1-watcher.finalityDepth)
	}
	if fromBlock > finalizedHead {
		return fromBlock, finalizedHead
	}
	if watcher.lastScannedBlock >= finalizedHead || watcher.pageSize > finalizedHead-watcher.lastScannedBlock {
		return fromBlock, finalizedHead
	}
	return fromBlock, watcher.lastScannedBlock + watcher.pageSize
}
