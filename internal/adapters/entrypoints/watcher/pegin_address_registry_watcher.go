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

type PegInAddressRegistryWatcherUseCases struct {
	discoverUseCase   *w.DiscoverRegisteredAddressUseCase
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase
}

func NewPegInAddressRegistryWatcherUseCases(
	discoverUseCase *w.DiscoverRegisteredAddressUseCase,
	getWatchedUseCase *w.GetWatchedRegisteredAddressesUseCase,
) *PegInAddressRegistryWatcherUseCases {
	return &PegInAddressRegistryWatcherUseCases{
		discoverUseCase:   discoverUseCase,
		getWatchedUseCase: getWatchedUseCase,
	}
}

type PegInAddressRegistryWatcher struct {
	discoverUseCase    *w.DiscoverRegisteredAddressUseCase
	getWatchedUseCase  *w.GetWatchedRegisteredAddressesUseCase
	repository         rootstock.PegInAddressRegistryWatchRepository
	registry           blockchain.PegInAddressRegistryContract
	rskRpc             blockchain.RootstockRpcServer
	btcNetwork         blockchain.BitcoinNetwork
	wallet             blockchain.BitcoinWallet
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	startBlock         uint64
	pageSize           uint64
	finalityDepth      uint64
	lastScannedBlock   uint64
	hasCursor          bool
	entries            []rootstock.PegInAddressRegistryWatchEntry
	stateMutex         sync.RWMutex
}

func NewPegInAddressRegistryWatcher(
	useCases *PegInAddressRegistryWatcherUseCases,
	repository rootstock.PegInAddressRegistryWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *PegInAddressRegistryWatcher {
	return &PegInAddressRegistryWatcher{
		discoverUseCase:    useCases.discoverUseCase,
		getWatchedUseCase:  useCases.getWatchedUseCase,
		repository:         repository,
		registry:           registry,
		rskRpc:             rskRpc,
		btcNetwork:         btcNetwork,
		wallet:             wallet,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
		startBlock:         startBlock,
		pageSize:           pageSize,
		finalityDepth:      finalityDepth,
	}
}

func (watcher *PegInAddressRegistryWatcher) Prepare(ctx context.Context) error {
	entries, err := watcher.getWatchedUseCase.Run(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry watch set: %w", err)
	}
	lastScannedBlock, found, err := watcher.repository.GetCursor(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry scan cursor: %w", err)
	}

	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	watcher.entries = entries
	watcher.lastScannedBlock = lastScannedBlock
	watcher.hasCursor = found
	return nil
}

func (watcher *PegInAddressRegistryWatcher) Start() {
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

func (watcher *PegInAddressRegistryWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("PegInAddressRegistryWatcher shut down")
}

func (watcher *PegInAddressRegistryWatcher) scan(ctx context.Context) error {
	pending := make([]*rootstock.PegInAddressRegistryWatchEntry, 0)
	if err := watcher.retryDiscoveredEntries(ctx, &pending); err != nil {
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
	if err = watcher.repository.SetCursor(ctx, toBlock); err != nil {
		return fmt.Errorf("persist PegIn address registry scan cursor: %w", err)
	}

	watcher.stateMutex.Lock()
	watcher.lastScannedBlock = toBlock
	watcher.hasCursor = true
	watcher.stateMutex.Unlock()
	return nil
}

func (watcher *PegInAddressRegistryWatcher) processEvents(
	ctx context.Context,
	events []blockchain.AddressRegistered,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
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

func (watcher *PegInAddressRegistryWatcher) retryDiscoveredEntries(
	ctx context.Context,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) error {
	watcher.stateMutex.RLock()
	entries := append([]rootstock.PegInAddressRegistryWatchEntry(nil), watcher.entries...)
	watcher.stateMutex.RUnlock()
	for index := range entries {
		if entries[index].State != rootstock.PegInAddressRegistryWatchDiscovered {
			continue
		}
		if err := watcher.discoverEvent(ctx, watchEntryToEvent(entries[index]), pending); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) discoverEvent(
	ctx context.Context,
	event blockchain.AddressRegistered,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) error {
	entry, needsRescan, err := watcher.discoverUseCase.Run(ctx, event)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}
	watcher.rememberEntry(*entry)
	if !needsRescan || pendingContains(*pending, entry.RskAddress) {
		return nil
	}
	*pending = append(*pending, entry)
	return nil
}

func (watcher *PegInAddressRegistryWatcher) rescanPendingImports(
	ctx context.Context,
	pending []*rootstock.PegInAddressRegistryWatchEntry,
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
	for _, entry := range pending {
		if err = watcher.discoverUseCase.MarkImported(ctx, entry); err != nil {
			return err
		}
		watcher.rememberEntry(*entry)
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) recordPendingRescanError(
	ctx context.Context,
	pending []*rootstock.PegInAddressRegistryWatchEntry,
	rescanErr error,
) error {
	for _, entry := range pending {
		if err := watcher.discoverUseCase.RecordError(ctx, entry, rescanErr); err != nil {
			return err
		}
		watcher.rememberEntry(*entry)
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) rememberEntry(entry rootstock.PegInAddressRegistryWatchEntry) {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for index := range watcher.entries {
		if watcher.entries[index].RskAddress == entry.RskAddress {
			watcher.entries[index] = entry
			return
		}
	}
	watcher.entries = append(watcher.entries, entry)
}

func pendingContains(pending []*rootstock.PegInAddressRegistryWatchEntry, rskAddress string) bool {
	for _, entry := range pending {
		if entry.RskAddress == rskAddress {
			return true
		}
	}
	return false
}

func watchEntryToEvent(entry rootstock.PegInAddressRegistryWatchEntry) blockchain.AddressRegistered {
	return blockchain.AddressRegistered{
		RskAddress:       entry.RskAddress,
		Registrant:       entry.Registrant,
		RegistrationRoot: entry.RegistrationRoot,
		TxHash:           entry.TxHash,
		BlockNumber:      entry.BlockNumber,
		LogIndex:         entry.LogIndex,
	}
}

func (watcher *PegInAddressRegistryWatcher) nextRange(finalizedHead uint64) (uint64, uint64) {
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
