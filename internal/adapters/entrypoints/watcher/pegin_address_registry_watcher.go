package watcher

import (
	"context"
	"fmt"
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
	pending := make([]*rootstock.PegInWatch, 0)
	if err := watcher.retryDiscoveredWatches(ctx, &pending); err != nil {
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
	if err = watcher.processEvents(ctx, events, &pending); err != nil {
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
	pending *[]*rootstock.PegInWatch,
) error {
	sort.Slice(events, func(firstIndex, secondIndex int) bool {
		if events[firstIndex].BlockNumber == events[secondIndex].BlockNumber {
			return events[firstIndex].LogIndex < events[secondIndex].LogIndex
		}
		return events[firstIndex].BlockNumber < events[secondIndex].BlockNumber
	})

	for _, event := range events {
		if err := watcher.discoverEvent(ctx, event, pending); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegInWatcher) retryDiscoveredWatches(
	ctx context.Context,
	pending *[]*rootstock.PegInWatch,
) error {
	watcher.stateMutex.RLock()
	watches := append([]rootstock.PegInWatch(nil), watcher.watches...)
	watcher.stateMutex.RUnlock()
	for index := range watches {
		if watches[index].State != rootstock.PegInWatchDiscovered {
			continue
		}
		if err := watcher.discoverEvent(ctx, watchToEvent(watches[index]), pending); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegInWatcher) discoverEvent(
	ctx context.Context,
	event blockchain.AddressRegistered,
	pending *[]*rootstock.PegInWatch,
) error {
	watch, needsRescan, err := watcher.discoverUseCase.Run(ctx, event)
	if err != nil {
		return err
	}
	if watch == nil {
		return nil
	}
	watcher.rememberWatch(*watch)
	if !needsRescan || pendingContains(*pending, watch.TxHash, watch.LogIndex) {
		return nil
	}
	*pending = append(*pending, watch)
	return nil
}

func (watcher *PegInWatcher) rescanPendingImports(
	ctx context.Context,
	pending []*rootstock.PegInWatch,
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
	for _, watch := range pending {
		if err = watcher.markImportedUseCase.Run(ctx, watch); err != nil {
			return err
		}
		watcher.rememberWatch(*watch)
	}
	return nil
}

func (watcher *PegInWatcher) recordPendingRescanError(
	ctx context.Context,
	pending []*rootstock.PegInWatch,
	rescanErr error,
) error {
	for _, watch := range pending {
		if err := watcher.recordErrorUseCase.Run(ctx, watch, rescanErr); err != nil {
			return err
		}
		watcher.rememberWatch(*watch)
	}
	return nil
}

func (watcher *PegInWatcher) rememberWatch(watch rootstock.PegInWatch) {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for index := range watcher.watches {
		if watcher.watches[index].TxHash == watch.TxHash && watcher.watches[index].LogIndex == watch.LogIndex {
			watcher.watches[index] = watch
			return
		}
	}
	watcher.watches = append(watcher.watches, watch)
}

func pendingContains(pending []*rootstock.PegInWatch, txHash string, logIndex uint) bool {
	for _, watch := range pending {
		if watch.TxHash == txHash && watch.LogIndex == logIndex {
			return true
		}
	}
	return false
}

func watchToEvent(watch rootstock.PegInWatch) blockchain.AddressRegistered {
	return blockchain.AddressRegistered{
		RskAddress:       watch.RskAddress,
		Registrant:       watch.Registrant,
		RegistrationRoot: watch.RegistrationRoot,
		TxHash:           watch.TxHash,
		BlockNumber:      watch.BlockNumber,
		LogIndex:         watch.LogIndex,
	}
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
