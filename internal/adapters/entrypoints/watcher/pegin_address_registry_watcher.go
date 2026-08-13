package watcher

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
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
	repository         rootstock.PegInAddressRegistryWatchRepositorySet
	registry           blockchain.PegInAddressRegistryContract
	rskRpc             blockchain.RootstockRpcServer
	btcNetwork         blockchain.BitcoinNetwork
	wallet             blockchain.BitcoinWallet
	eventBus           entities.EventBus
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	startBlock         uint64
	pageSize           uint64
	finalityDepth      uint64
	localRoot          [32]byte
	lastProcessedBlock uint64
	hasCheckpoint      bool
	entries            []rootstock.PegInAddressRegistryWatchEntry
	stateMutex         sync.RWMutex
	replayMutex        sync.Mutex
}

var errPegInAddressRegistryRootMismatch = errors.New("PegIn address registry roots differ")

type pegInAddressRegistryRootMismatchError struct {
	blockNumber uint64
	localRoot   [32]byte
	chainRoot   [32]byte
	source      string
}

func (mismatch *pegInAddressRegistryRootMismatchError) Error() string {
	return fmt.Sprintf(
		"%s at %s block %d: local=%x chain=%x",
		errPegInAddressRegistryRootMismatch,
		mismatch.source,
		mismatch.blockNumber,
		mismatch.localRoot,
		mismatch.chainRoot,
	)
}

func (mismatch *pegInAddressRegistryRootMismatchError) Unwrap() error {
	return errPegInAddressRegistryRootMismatch
}

func NewPegInAddressRegistryWatcher(
	useCases *PegInAddressRegistryWatcherUseCases,
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	eventBus entities.EventBus,
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
		eventBus:           eventBus,
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
	checkpoint, found, err := watcher.repository.GetCheckpoint(ctx)
	if err != nil {
		return fmt.Errorf("load PegIn address registry checkpoint: %w", err)
	}

	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	watcher.entries = entries
	watcher.localRoot = checkpoint.LocalRoot
	watcher.lastProcessedBlock = checkpoint.LastProcessedBlock
	watcher.hasCheckpoint = found
	return nil
}

func (watcher *PegInAddressRegistryWatcher) Start() {
	var reorgEvents <-chan entities.Event
	if watcher.eventBus != nil {
		reorgEvents = watcher.eventBus.Subscribe(blockchain.NodeReorgCheckEventId)
	}
	for {
		select {
		case <-watcher.ticker.C():
			watcher.scanAndLog()
		case event, open := <-reorgEvents:
			if !open {
				reorgEvents = nil
				continue
			}
			watcher.resyncAfterReorg(event)
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			return
		}
	}
}

func (watcher *PegInAddressRegistryWatcher) scanAndLog() {
	if err := watcher.scan(context.Background()); err != nil {
		log.Errorf("PegIn address registry watcher scan failed: %v", err)
	}
}

func (watcher *PegInAddressRegistryWatcher) resyncAfterReorg(event entities.Event) {
	reorg, ok := event.(blockchain.NodeReorgCheckEvent)
	if !ok || reorg.NodeType != entities.NodeTypeRootstock || reorg.CurrentDepth == 0 {
		return
	}
	if err := watcher.resync(context.Background()); err != nil {
		log.Errorf("PegIn address registry watcher scan failed: %v", err)
	}
}

func (watcher *PegInAddressRegistryWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("PegInAddressRegistryWatcher shut down")
}

func (watcher *PegInAddressRegistryWatcher) scan(ctx context.Context) error {
	return watcher.runReplay(ctx, false)
}

func (watcher *PegInAddressRegistryWatcher) resync(ctx context.Context) error {
	return watcher.runReplay(ctx, true)
}

func (watcher *PegInAddressRegistryWatcher) runReplay(ctx context.Context, discardCheckpoint bool) error {
	watcher.replayMutex.Lock()
	defer watcher.replayMutex.Unlock()

	replay, err := watcher.prepareReplay(ctx, discardCheckpoint)
	if err != nil {
		return err
	}
	if replay == nil {
		return nil
	}
	pending := make([]*rootstock.PegInAddressRegistryWatchEntry, 0)
	if err = watcher.retryDiscoveredEntries(ctx, &pending); err != nil {
		return err
	}
	if replay.recoveryReason != "" {
		watcher.reportResyncStarted(replay.recoveryReason)
	}
	checkpoint, err := watcher.replayAndVerify(ctx, replay.trustedCheckpoint, replay.finalizedHead, replay.head, &pending)
	if err == nil {
		if rescanErr := watcher.rescanPendingImports(ctx, pending); rescanErr != nil {
			return rescanErr
		}
		return watcher.publishReplayCheckpoint(ctx, checkpoint, replay.trustedCheckpoint)
	}
	return watcher.recoverFromRootMismatch(ctx, err, replay.finalizedHead, replay.head)
}

type pegInAddressRegistryReplay struct {
	head              uint64
	finalizedHead     uint64
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint
	recoveryReason    string
}

func (watcher *PegInAddressRegistryWatcher) prepareReplay(
	ctx context.Context,
	discardCheckpoint bool,
) (*pegInAddressRegistryReplay, error) {
	recoveryReason := ""
	if discardCheckpoint {
		if err := watcher.discardCheckpoint(ctx); err != nil {
			return nil, err
		}
		recoveryReason = "chain_reorganization"
	}
	head, err := watcher.rskRpc.GetHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("get RSK height: %w", err)
	}
	if head < watcher.finalityDepth {
		return nil, nil
	}
	finalizedHead := head - watcher.finalityDepth
	if finalizedHead < watcher.startBlock {
		return nil, nil
	}
	return &pegInAddressRegistryReplay{
		head:              head,
		finalizedHead:     finalizedHead,
		trustedCheckpoint: watcher.trustedCheckpoint(),
		recoveryReason:    recoveryReason,
	}, nil
}

func (watcher *PegInAddressRegistryWatcher) trustedCheckpoint() *rootstock.PegInAddressRegistryWatchCheckpoint {
	watcher.stateMutex.RLock()
	defer watcher.stateMutex.RUnlock()
	if !watcher.hasCheckpoint {
		return nil
	}
	return &rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          watcher.localRoot,
		LastProcessedBlock: watcher.lastProcessedBlock,
	}
}

func (watcher *PegInAddressRegistryWatcher) replayAndVerify(
	ctx context.Context,
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
	finalizedHead uint64,
	head uint64,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) (rootstock.PegInAddressRegistryWatchCheckpoint, error) {
	checkpoint, err := watcher.replay(ctx, trustedCheckpoint, finalizedHead, pending)
	if err != nil {
		return checkpoint, err
	}
	err = watcher.verifyReplayRoot(ctx, checkpoint.LocalRoot, checkpoint.LastProcessedBlock, head)
	return checkpoint, err
}

func (watcher *PegInAddressRegistryWatcher) publishReplayCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	if trustedCheckpoint != nil && checkpoint == *trustedCheckpoint {
		return nil
	}
	return watcher.publishCheckpoint(ctx, checkpoint)
}

func (watcher *PegInAddressRegistryWatcher) recoverFromRootMismatch(
	ctx context.Context,
	replayErr error,
	finalizedHead uint64,
	head uint64,
) error {
	if !errors.Is(replayErr, errPegInAddressRegistryRootMismatch) {
		return replayErr
	}
	watcher.reportRootMismatch(replayErr)
	if discardErr := watcher.discardCheckpoint(ctx); discardErr != nil {
		return errors.Join(replayErr, discardErr)
	}
	watcher.reportResyncStarted("root_mismatch")
	pending := make([]*rootstock.PegInAddressRegistryWatchEntry, 0)
	checkpoint, err := watcher.replayAndVerify(ctx, nil, finalizedHead, head, &pending)
	if err == nil {
		if rescanErr := watcher.rescanPendingImports(ctx, pending); rescanErr != nil {
			return rescanErr
		}
		return watcher.publishCheckpoint(ctx, checkpoint)
	}
	return errors.Join(err, watcher.discardCheckpoint(ctx))
}

func (watcher *PegInAddressRegistryWatcher) reportRootMismatch(err error) {
	fields := log.Fields{"error": err.Error()}
	event := blockchain.PegInAddressRegistryRootMismatchEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryRootMismatchEventId),
	}
	var mismatch *pegInAddressRegistryRootMismatchError
	if errors.As(err, &mismatch) {
		fields["block_number"] = mismatch.blockNumber
		fields["local_root"] = fmt.Sprintf("0x%x", mismatch.localRoot)
		fields["chain_root"] = fmt.Sprintf("0x%x", mismatch.chainRoot)
		fields["source"] = mismatch.source
		event.BlockNumber = mismatch.blockNumber
		event.LocalRoot = mismatch.localRoot
		event.ChainRoot = mismatch.chainRoot
	}
	log.WithFields(fields).Error("PegIn address registry root mismatch")
	if watcher.eventBus != nil {
		watcher.eventBus.Publish(event)
	}
}

func (watcher *PegInAddressRegistryWatcher) reportResyncStarted(reason string) {
	if watcher.eventBus == nil {
		return
	}
	watcher.eventBus.Publish(blockchain.PegInAddressRegistryResyncStartedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryResyncStartedEventId),
		Reason:    reason,
	})
}

func (watcher *PegInAddressRegistryWatcher) discardCheckpoint(ctx context.Context) error {
	if err := watcher.repository.DeleteCheckpoint(ctx); err != nil {
		return fmt.Errorf("delete PegIn address registry checkpoint: %w", err)
	}
	watcher.stateMutex.Lock()
	watcher.localRoot = [32]byte{}
	watcher.lastProcessedBlock = 0
	watcher.hasCheckpoint = false
	watcher.stateMutex.Unlock()
	return nil
}

func (watcher *PegInAddressRegistryWatcher) replay(
	ctx context.Context,
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
	finalizedHead uint64,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) (rootstock.PegInAddressRegistryWatchCheckpoint, error) {
	fromBlock := watcher.startBlock
	localRoot := [32]byte{}
	checkpoint := rootstock.PegInAddressRegistryWatchCheckpoint{
		LocalRoot:          localRoot,
		LastProcessedBlock: finalizedHead,
	}
	if trustedCheckpoint != nil {
		checkpoint = *trustedCheckpoint
		if trustedCheckpoint.LastProcessedBlock >= finalizedHead {
			return *trustedCheckpoint, nil
		}
		fromBlock = trustedCheckpoint.LastProcessedBlock + 1
		localRoot = trustedCheckpoint.LocalRoot
	}
	for fromBlock <= finalizedHead {
		toBlock := finalizedHead
		if watcher.pageSize-1 <= finalizedHead-fromBlock {
			toBlock = fromBlock + watcher.pageSize - 1
		}
		events, err := watcher.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
		if err != nil {
			return checkpoint, fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err)
		}
		events = eventsWithinRange(events, fromBlock, toBlock)
		localRoot, err = watcher.processEvents(ctx, events, localRoot, pending)
		if err != nil {
			return checkpoint, err
		}
		checkpoint = rootstock.PegInAddressRegistryWatchCheckpoint{
			LocalRoot:          localRoot,
			LastProcessedBlock: toBlock,
		}
		if toBlock == finalizedHead {
			break
		}
		fromBlock = toBlock + 1
	}
	return checkpoint, nil
}

func (watcher *PegInAddressRegistryWatcher) publishCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) error {
	if err := watcher.repository.SetCheckpoint(ctx, checkpoint); err != nil {
		return fmt.Errorf("persist PegIn address registry checkpoint: %w", err)
	}
	watcher.rememberCheckpoint(checkpoint)
	return nil
}

func (watcher *PegInAddressRegistryWatcher) rememberCheckpoint(
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
) {
	watcher.stateMutex.Lock()
	watcher.localRoot = checkpoint.LocalRoot
	watcher.lastProcessedBlock = checkpoint.LastProcessedBlock
	watcher.hasCheckpoint = true
	watcher.stateMutex.Unlock()
}

func eventsWithinRange(
	events []blockchain.AddressRegistered,
	fromBlock uint64,
	toBlock uint64,
) []blockchain.AddressRegistered {
	inRange := make([]blockchain.AddressRegistered, 0, len(events))
	for _, event := range events {
		if event.BlockNumber >= fromBlock && event.BlockNumber <= toBlock {
			inRange = append(inRange, event)
		}
	}
	return inRange
}

func (watcher *PegInAddressRegistryWatcher) verifyReplayRoot(
	ctx context.Context,
	localRoot [32]byte,
	finalizedHead uint64,
	head uint64,
) error {
	verifiedRoot := localRoot
	if finalizedHead < head {
		fromBlock := finalizedHead + 1
		unconfirmedEvents, err := watcher.registry.GetAddressRegisteredEvents(ctx, fromBlock, &head)
		if err != nil {
			return fmt.Errorf("get unconfirmed AddressRegistered events for blocks %d-%d: %w", fromBlock, head, err)
		}
		unconfirmedEvents = eventsWithinRange(unconfirmedEvents, fromBlock, head)
		if len(unconfirmedEvents) != 0 {
			verifiedRoot, err = foldRegistrationRoots(localRoot, unconfirmedEvents)
			if err != nil {
				return err
			}
		}
	}
	chainRoot, err := watcher.registry.GetRegistrationRoot(ctx, head)
	if err != nil {
		return fmt.Errorf("get PegIn address registry root: %w", err)
	}
	if verifiedRoot != chainRoot {
		return &pegInAddressRegistryRootMismatchError{
			blockNumber: head,
			localRoot:   verifiedRoot,
			chainRoot:   chainRoot,
			source:      "captured_head",
		}
	}
	return nil
}

func (watcher *PegInAddressRegistryWatcher) processEvents(
	ctx context.Context,
	events []blockchain.AddressRegistered,
	localRoot [32]byte,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) ([32]byte, error) {
	events = orderedUniqueEvents(events)
	for _, event := range events {
		var err error
		localRoot, err = blockchain.FoldPegInAddressRegistryRoot(localRoot, event.RskAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if event.RegistrationRoot != localRoot {
			return [32]byte{}, &pegInAddressRegistryRootMismatchError{
				blockNumber: event.BlockNumber,
				localRoot:   localRoot,
				chainRoot:   event.RegistrationRoot,
				source:      fmt.Sprintf("event_%s_%d", event.TxHash, event.LogIndex),
			}
		}
		if err = watcher.discoverEvent(ctx, event, pending); err != nil {
			return [32]byte{}, err
		}
	}
	return localRoot, nil
}

type registryEventIdentity struct {
	txHash   string
	logIndex uint
}

func orderedUniqueEvents(events []blockchain.AddressRegistered) []blockchain.AddressRegistered {
	ordered := append([]blockchain.AddressRegistered(nil), events...)
	sort.Slice(ordered, func(firstIndex, secondIndex int) bool {
		if ordered[firstIndex].BlockNumber == ordered[secondIndex].BlockNumber {
			return ordered[firstIndex].LogIndex < ordered[secondIndex].LogIndex
		}
		return ordered[firstIndex].BlockNumber < ordered[secondIndex].BlockNumber
	})
	unique := make([]blockchain.AddressRegistered, 0, len(ordered))
	seen := make(map[registryEventIdentity]struct{}, len(ordered))
	for _, event := range ordered {
		identity := registryEventIdentity{txHash: event.TxHash, logIndex: event.LogIndex}
		if _, duplicate := seen[identity]; duplicate {
			continue
		}
		seen[identity] = struct{}{}
		unique = append(unique, event)
	}
	return unique
}

func foldRegistrationRoots(
	localRoot [32]byte,
	events []blockchain.AddressRegistered,
) ([32]byte, error) {
	for _, event := range orderedUniqueEvents(events) {
		var err error
		localRoot, err = blockchain.FoldPegInAddressRegistryRoot(localRoot, event.RskAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if event.RegistrationRoot != localRoot {
			return [32]byte{}, &pegInAddressRegistryRootMismatchError{
				blockNumber: event.BlockNumber,
				localRoot:   localRoot,
				chainRoot:   event.RegistrationRoot,
				source:      fmt.Sprintf("event_%s_%d", event.TxHash, event.LogIndex),
			}
		}
	}
	return localRoot, nil
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
	if pendingContains(*pending, event.RskAddress) {
		return nil
	}
	entry, needsRescan, err := watcher.discoverUseCase.Run(ctx, event)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}
	watcher.rememberEntry(*entry)
	if !needsRescan {
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
