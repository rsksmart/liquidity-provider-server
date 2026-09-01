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
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type ReplayRegisteredAddressesUseCase struct {
	repository    rootstock.PegInWatchRepositorySet
	registry      blockchain.PegInAddressRegistryContract
	rskRpc        blockchain.RootstockRpcServer
	eventBus      entities.EventBus
	wallet        blockchain.BitcoinWallet
	startBlock    uint64
	pageSize      uint64
	finalityDepth uint64
	replayMutex   sync.Mutex
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

func NewReplayRegisteredAddressesUseCase(
	repository rootstock.PegInWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	eventBus entities.EventBus,
	wallet blockchain.BitcoinWallet,
	finalityDepth uint64,
) *ReplayRegisteredAddressesUseCase {
	return &ReplayRegisteredAddressesUseCase{
		repository:    repository,
		registry:      registry,
		rskRpc:        rskRpc,
		eventBus:      eventBus,
		wallet:        wallet,
		finalityDepth: finalityDepth,
	}
}

func (useCase *ReplayRegisteredAddressesUseCase) Run(
	ctx context.Context,
	discardCheckpoint bool,
	startBlock uint64,
	pageSize uint64,
) ([]*rootstock.PegInWatch, error) {
	useCase.replayMutex.Lock()
	defer useCase.replayMutex.Unlock()

	useCase.startBlock = startBlock
	useCase.pageSize = pageSize

	pending, err := useCase.runReplay(ctx, discardCheckpoint)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ReplayRegisteredAddressesId, err)
	}
	return pending, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) runReplay(
	ctx context.Context,
	discardCheckpoint bool,
) ([]*rootstock.PegInWatch, error) {
	replay, err := useCase.prepareReplay(ctx, discardCheckpoint)
	if err != nil {
		return nil, err
	}
	if replay == nil {
		return nil, nil
	}
	pending := []*rootstock.PegInWatch{}
	pending, err = useCase.retryDiscoveredEntries(ctx, pending)
	if err != nil {
		return nil, err
	}
	if replay.recoveryReason != "" {
		useCase.reportResyncStarted(replay.recoveryReason)
	}
	checkpoint, pending, err := useCase.replayAndVerify(ctx, replay.trustedCheckpoint, replay.finalizedHead, replay.head, pending)
	if err == nil {
		if err = useCase.publishReplayCheckpoint(ctx, checkpoint, replay.trustedCheckpoint); err != nil {
			return nil, err
		}
		return pending, nil
	}
	return useCase.recoverFromRootMismatch(ctx, err, replay.finalizedHead, replay.head)
}

type pegInAddressRegistryReplay struct {
	head              uint64
	finalizedHead     uint64
	trustedCheckpoint *rootstock.PegInWatchCheckpoint
	recoveryReason    string
}

func (useCase *ReplayRegisteredAddressesUseCase) prepareReplay(
	ctx context.Context,
	discardCheckpoint bool,
) (*pegInAddressRegistryReplay, error) {
	recoveryReason := ""
	if discardCheckpoint {
		if err := useCase.discardCheckpoint(ctx); err != nil {
			return nil, err
		}
		recoveryReason = "chain_reorganization"
	}
	head, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("get RSK height: %w", err)
	}
	if head < useCase.finalityDepth {
		return nil, nil
	}
	finalizedHead := head - useCase.finalityDepth
	if finalizedHead < useCase.startBlock {
		return nil, nil
	}
	trustedCheckpoint, err := useCase.loadTrustedCheckpoint(ctx)
	if err != nil {
		return nil, err
	}
	return &pegInAddressRegistryReplay{
		head:              head,
		finalizedHead:     finalizedHead,
		trustedCheckpoint: trustedCheckpoint,
		recoveryReason:    recoveryReason,
	}, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) loadTrustedCheckpoint(
	ctx context.Context,
) (*rootstock.PegInWatchCheckpoint, error) {
	checkpoint, found, err := useCase.repository.GetCheckpoint(ctx)
	if err != nil {
		return nil, fmt.Errorf("load PegIn address registry checkpoint: %w", err)
	}
	if !found {
		return nil, nil
	}
	return &checkpoint, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) replayAndVerify(
	ctx context.Context,
	trustedCheckpoint *rootstock.PegInWatchCheckpoint,
	finalizedHead uint64,
	head uint64,
	pending []*rootstock.PegInWatch,
) (rootstock.PegInWatchCheckpoint, []*rootstock.PegInWatch, error) {
	checkpoint, pending, err := useCase.replay(ctx, trustedCheckpoint, finalizedHead, pending)
	if err != nil {
		return checkpoint, pending, err
	}
	err = useCase.verifyReplayRoot(ctx, checkpoint.LocalRoot, checkpoint.LastProcessedBlock, head)
	return checkpoint, pending, err
}

func (useCase *ReplayRegisteredAddressesUseCase) publishReplayCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
	trustedCheckpoint *rootstock.PegInWatchCheckpoint,
) error {
	if trustedCheckpoint != nil && checkpoint == *trustedCheckpoint {
		return nil
	}
	return useCase.publishCheckpoint(ctx, checkpoint)
}

func (useCase *ReplayRegisteredAddressesUseCase) recoverFromRootMismatch(
	ctx context.Context,
	replayErr error,
	finalizedHead uint64,
	head uint64,
) ([]*rootstock.PegInWatch, error) {
	if !errors.Is(replayErr, errPegInAddressRegistryRootMismatch) {
		return nil, replayErr
	}
	useCase.reportRootMismatch(replayErr)
	if discardErr := useCase.discardCheckpoint(ctx); discardErr != nil {
		return nil, errors.Join(replayErr, discardErr)
	}
	useCase.reportResyncStarted("root_mismatch")
	recovered := []*rootstock.PegInWatch{}
	checkpoint, recovered, err := useCase.replayAndVerify(ctx, nil, finalizedHead, head, recovered)
	if err == nil {
		if err = useCase.publishCheckpoint(ctx, checkpoint); err != nil {
			return nil, err
		}
		return recovered, nil
	}
	return nil, errors.Join(err, useCase.discardCheckpoint(ctx))
}

func (useCase *ReplayRegisteredAddressesUseCase) reportRootMismatch(err error) {
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
	if useCase.eventBus != nil {
		useCase.eventBus.Publish(event)
	}
}

func (useCase *ReplayRegisteredAddressesUseCase) reportResyncStarted(reason string) {
	if useCase.eventBus == nil {
		return
	}
	useCase.eventBus.Publish(blockchain.PegInAddressRegistryResyncStartedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryResyncStartedEventId),
		Reason:    reason,
	})
}

func (useCase *ReplayRegisteredAddressesUseCase) discardCheckpoint(ctx context.Context) error {
	if err := useCase.repository.DeleteCheckpoint(ctx); err != nil {
		return fmt.Errorf("delete PegIn address registry checkpoint: %w", err)
	}
	return nil
}

func (useCase *ReplayRegisteredAddressesUseCase) replay(
	ctx context.Context,
	trustedCheckpoint *rootstock.PegInWatchCheckpoint,
	finalizedHead uint64,
	pending []*rootstock.PegInWatch,
) (rootstock.PegInWatchCheckpoint, []*rootstock.PegInWatch, error) {
	fromBlock := useCase.startBlock
	localRoot := [32]byte{}
	checkpoint := rootstock.PegInWatchCheckpoint{
		LocalRoot:          localRoot,
		LastProcessedBlock: finalizedHead,
	}
	if trustedCheckpoint != nil {
		checkpoint = *trustedCheckpoint
		if trustedCheckpoint.LastProcessedBlock >= finalizedHead {
			return *trustedCheckpoint, pending, nil
		}
		fromBlock = trustedCheckpoint.LastProcessedBlock + 1
		localRoot = trustedCheckpoint.LocalRoot
	}
	for fromBlock <= finalizedHead {
		toBlock := finalizedHead
		if useCase.pageSize-1 <= finalizedHead-fromBlock {
			toBlock = fromBlock + useCase.pageSize - 1
		}
		events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
		if err != nil {
			return checkpoint, pending, fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err)
		}
		events = eventsWithinRange(events, fromBlock, toBlock)
		localRoot, pending, err = useCase.processEvents(ctx, events, localRoot, pending)
		if err != nil {
			return checkpoint, pending, err
		}
		checkpoint = rootstock.PegInWatchCheckpoint{
			LocalRoot:          localRoot,
			LastProcessedBlock: toBlock,
		}
		if toBlock == finalizedHead {
			break
		}
		fromBlock = toBlock + 1
	}
	return checkpoint, pending, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) publishCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
) error {
	if err := useCase.repository.SetCheckpoint(ctx, checkpoint); err != nil {
		return fmt.Errorf("persist PegIn address registry checkpoint: %w", err)
	}
	return nil
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

func (useCase *ReplayRegisteredAddressesUseCase) verifyReplayRoot(
	ctx context.Context,
	localRoot [32]byte,
	finalizedHead uint64,
	head uint64,
) error {
	verifiedRoot := localRoot
	if finalizedHead < head {
		fromBlock := finalizedHead + 1
		unconfirmedEvents, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &head)
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
	chainRoot, err := useCase.registry.GetRegistrationRoot(ctx, head)
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

func (useCase *ReplayRegisteredAddressesUseCase) processEvents(
	ctx context.Context,
	events []blockchain.AddressRegistered,
	localRoot [32]byte,
	pending []*rootstock.PegInWatch,
) ([32]byte, []*rootstock.PegInWatch, error) {
	events = orderedUniqueEvents(events)
	for _, event := range events {
		var err error
		localRoot, err = blockchain.FoldPegInAddressRegistryRoot(localRoot, event.RskAddress)
		if err != nil {
			return [32]byte{}, pending, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if event.RegistrationRoot != localRoot {
			return [32]byte{}, pending, &pegInAddressRegistryRootMismatchError{
				blockNumber: event.BlockNumber,
				localRoot:   localRoot,
				chainRoot:   event.RegistrationRoot,
				source:      fmt.Sprintf("event_%s_%d", event.TxHash, event.LogIndex),
			}
		}
		pending, err = useCase.discoverEvent(ctx, event, pending)
		if err != nil {
			return [32]byte{}, pending, err
		}
	}
	return localRoot, pending, nil
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

func (useCase *ReplayRegisteredAddressesUseCase) retryDiscoveredEntries(
	ctx context.Context,
	pending []*rootstock.PegInWatch,
) ([]*rootstock.PegInWatch, error) {
	entries, err := useCase.repository.List(ctx)
	if err != nil {
		return pending, err
	}
	for index := range entries {
		if entries[index].State != rootstock.PegInWatchDiscovered {
			continue
		}
		pending, err = useCase.discoverEvent(ctx, blockchain.NewAddressRegisteredFromWatchEntry(entries[index]), pending)
		if err != nil {
			return pending, err
		}
	}
	return pending, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) discoverEvent(
	ctx context.Context,
	event blockchain.AddressRegistered,
	pending []*rootstock.PegInWatch,
) ([]*rootstock.PegInWatch, error) {
	if rootstock.PegInWatches(pending).Contains(event.RskAddress) {
		return pending, nil
	}
	entry, err := loadOrCreateWatchEntry(ctx, useCase.repository, event, usecases.ReplayRegisteredAddressesId)
	if err != nil {
		return pending, err
	}
	needsRescan, err := resolveAndImportWatchEntry(
		ctx,
		useCase.repository,
		useCase.registry,
		useCase.wallet,
		entry,
		usecases.ReplayRegisteredAddressesId,
	)
	if err != nil {
		return pending, err
	}
	if entry == nil || !needsRescan {
		return pending, nil
	}
	return append(pending, entry), nil
}
