package watcher

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type PegInAddressRegistryWatcher struct {
	repository         rootstock.PegInAddressRegistryWatchRepository
	registry           blockchain.PegInAddressRegistryContract
	rskRpc             blockchain.RootstockRpcServer
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
	repository rootstock.PegInAddressRegistryWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	wallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *PegInAddressRegistryWatcher {
	return &PegInAddressRegistryWatcher{
		repository:         repository,
		registry:           registry,
		rskRpc:             rskRpc,
		wallet:             wallet,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
		startBlock:         startBlock,
		pageSize:           pageSize,
		finalityDepth:      finalityDepth,
	}
}

func (watcher *PegInAddressRegistryWatcher) Prepare(ctx context.Context) error {
	entries, err := watcher.repository.List(ctx)
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
	if err := watcher.retryDiscoveredEntries(ctx); err != nil {
		return err
	}
	head, err := watcher.rskRpc.GetHeight(ctx)
	if err != nil {
		return fmt.Errorf("get RSK height: %w", err)
	}
	if head < watcher.finalityDepth {
		return nil
	}

	watcher.stateMutex.RLock()
	fromBlock, toBlock := watcher.nextRange(head - watcher.finalityDepth)
	watcher.stateMutex.RUnlock()
	if fromBlock > toBlock {
		return nil
	}

	events, err := watcher.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err)
	}
	if err = watcher.processEvents(ctx, events); err != nil {
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
) error {
	var err error
	sort.Slice(events, func(firstIndex, secondIndex int) bool {
		if events[firstIndex].BlockNumber == events[secondIndex].BlockNumber {
			return events[firstIndex].LogIndex < events[secondIndex].LogIndex
		}
		return events[firstIndex].BlockNumber < events[secondIndex].BlockNumber
	})

	for _, event := range events {
		now := time.Now().UTC()
		entry := rootstock.PegInAddressRegistryWatchEntry{
			TxHash:           event.TxHash,
			LogIndex:         event.LogIndex,
			BlockNumber:      event.BlockNumber,
			RskAddress:       event.RskAddress,
			Registrant:       event.Registrant,
			RegistrationRoot: event.RegistrationRoot,
			State:            rootstock.PegInAddressRegistryWatchDiscovered,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err = watcher.repository.Upsert(ctx, entry); err != nil {
			return fmt.Errorf("persist AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		persistedEntry, getErr := watcher.repository.Get(ctx, event.TxHash, event.LogIndex)
		if getErr != nil {
			return fmt.Errorf("load AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, getErr)
		}
		if persistedEntry == nil {
			return fmt.Errorf("load AddressRegistered event %s/%d: entry not found after upsert", event.TxHash, event.LogIndex)
		}
		if persistedEntry.State != rootstock.PegInAddressRegistryWatchDiscovered {
			watcher.rememberEntry(*persistedEntry)
			continue
		}
		if err = watcher.processDiscoveredEntry(ctx, persistedEntry); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) retryDiscoveredEntries(ctx context.Context) error {
	watcher.stateMutex.RLock()
	entries := append([]rootstock.PegInAddressRegistryWatchEntry(nil), watcher.entries...)
	watcher.stateMutex.RUnlock()
	for index := range entries {
		if entries[index].State != rootstock.PegInAddressRegistryWatchDiscovered {
			continue
		}
		if err := watcher.processDiscoveredEntry(ctx, &entries[index]); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) processDiscoveredEntry(
	ctx context.Context,
	entry *rootstock.PegInAddressRegistryWatchEntry,
) error {
	pegInAddress, err := watcher.registry.GetPegInAddress(entry.RskAddress)
	if err != nil {
		return watcher.recordEntryError(
			ctx,
			entry,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	encoding := pegInAddress.Encoding
	entry.Encoding = uint8(encoding)
	entry.UpdatedAt = time.Now().UTC()
	// LastError is cleared only on the paths that reach a persisted state, so recordEntryError can
	// still compare against the previously persisted error and suppress repeated identical writes.
	if encoding != blockchain.PegInAddressRegistryEncodingBase58 {
		entry.State = rootstock.PegInAddressRegistryWatchUnsupportedEncoding
		entry.LastError = ""
		if err = watcher.repository.Update(ctx, *entry); err != nil {
			return fmt.Errorf("persist unsupported encoding for event %s/%d: %w", entry.TxHash, entry.LogIndex, err)
		}
		log.Errorf(
			"PegIn address registry event %s/%d uses unsupported encoding %d",
			entry.TxHash,
			entry.LogIndex,
			encoding,
		)
		watcher.rememberEntry(*entry)
		return nil
	}
	// The registry returns the base58check payload, not the address, and leaves the encoding to the
	// caller, so this step is what makes the value importable by a Bitcoin node.
	if entry.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return watcher.recordEntryError(
			ctx,
			entry,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	entry.DepositTxID, err = watcher.registry.GetRegisteredBtcTransactionHash(ctx, entry.TxHash, entry.RskAddress)
	if err != nil {
		return watcher.recordEntryError(
			ctx,
			entry,
			fmt.Errorf("resolve registered BTC transaction for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	if err = watcher.wallet.ImportAddressWithRescan(entry.BtcAddress); err != nil && !isAlreadyImportedError(err) {
		return watcher.recordEntryError(
			ctx,
			entry,
			fmt.Errorf("import PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	entry.State = rootstock.PegInAddressRegistryWatchImported
	entry.LastError = ""
	if err = watcher.repository.Update(ctx, *entry); err != nil {
		return fmt.Errorf("persist imported state for event %s/%d: %w", entry.TxHash, entry.LogIndex, err)
	}
	watcher.rememberEntry(*entry)
	return nil
}

func (watcher *PegInAddressRegistryWatcher) recordEntryError(
	ctx context.Context,
	entry *rootstock.PegInAddressRegistryWatchEntry,
	entryErr error,
) error {
	log.Error(entryErr)
	if entry.LastError == entryErr.Error() {
		return nil
	}
	entry.LastError = entryErr.Error()
	entry.UpdatedAt = time.Now().UTC()
	if err := watcher.repository.Update(ctx, *entry); err != nil {
		return fmt.Errorf("persist PegIn address registry entry error for %s/%d: %w", entry.TxHash, entry.LogIndex, err)
	}
	watcher.rememberEntry(*entry)
	return nil
}

func (watcher *PegInAddressRegistryWatcher) rememberEntry(entry rootstock.PegInAddressRegistryWatchEntry) {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for index := range watcher.entries {
		if watcher.entries[index].TxHash == entry.TxHash && watcher.entries[index].LogIndex == entry.LogIndex {
			watcher.entries[index] = entry
			return
		}
	}
	watcher.entries = append(watcher.entries, entry)
}

func isAlreadyImportedError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already imported")
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
