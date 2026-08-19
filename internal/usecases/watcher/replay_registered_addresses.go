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
	discoverUseCase *DiscoverRegisteredAddressUseCase
	getWatched      *GetWatchedRegisteredAddressesUseCase
	repository      rootstock.PegInAddressRegistryWatchRepositorySet
	registry        blockchain.PegInAddressRegistryContract
	rskRpc          blockchain.RootstockRpcServer
	eventBus        entities.EventBus
	startBlock      uint64
	pageSize        uint64
	finalityDepth   uint64
	replayMutex     sync.Mutex
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

type peginAddressRegistryDeploymentVerifier interface {
	IsDeploymentBlock(ctx context.Context, blockNumber uint64) (bool, error)
}

func NewReplayRegisteredAddressesUseCase(
	discoverUseCase *DiscoverRegisteredAddressUseCase,
	getWatched *GetWatchedRegisteredAddressesUseCase,
	repository rootstock.PegInAddressRegistryWatchRepositorySet,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	eventBus entities.EventBus,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *ReplayRegisteredAddressesUseCase {
	return &ReplayRegisteredAddressesUseCase{
		discoverUseCase: discoverUseCase,
		getWatched:      getWatched,
		repository:      repository,
		registry:        registry,
		rskRpc:          rskRpc,
		eventBus:        eventBus,
		startBlock:      startBlock,
		pageSize:        pageSize,
		finalityDepth:   finalityDepth,
	}
}

func (useCase *ReplayRegisteredAddressesUseCase) Run(
	ctx context.Context,
	discardCheckpoint bool,
) ([]*rootstock.PegInAddressRegistryWatchEntry, error) {
	return useCase.RunThen(ctx, discardCheckpoint, nil)
}

func (useCase *ReplayRegisteredAddressesUseCase) RunThen(
	ctx context.Context,
	discardCheckpoint bool,
	finalize func(context.Context, []*rootstock.PegInAddressRegistryWatchEntry) error,
) ([]*rootstock.PegInAddressRegistryWatchEntry, error) {
	useCase.replayMutex.Lock()
	defer useCase.replayMutex.Unlock()

	pending, err := useCase.runReplay(ctx, discardCheckpoint, finalize)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ReplayRegisteredAddressesId, err)
	}
	return pending, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) ValidateDeployment(ctx context.Context, startBlock uint64) error {
	err := useCase.validateDeployment(ctx, startBlock)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.ReplayRegisteredAddressesId, err)
	}
	return nil
}

func (useCase *ReplayRegisteredAddressesUseCase) runReplay(
	ctx context.Context,
	discardCheckpoint bool,
	finalize func(context.Context, []*rootstock.PegInAddressRegistryWatchEntry) error,
) ([]*rootstock.PegInAddressRegistryWatchEntry, error) {
	replay, err := useCase.prepareReplay(ctx, discardCheckpoint)
	if err != nil {
		return nil, err
	}
	if replay == nil {
		return nil, nil
	}
	pending := make([]*rootstock.PegInAddressRegistryWatchEntry, 0)
	if err = useCase.retryDiscoveredEntries(ctx, &pending); err != nil {
		return nil, err
	}
	if replay.recoveryReason != "" {
		useCase.reportResyncStarted(replay.recoveryReason)
	}
	checkpoint, err := useCase.replayAndVerify(ctx, replay.trustedCheckpoint, replay.finalizedHead, replay.head, &pending)
	if err == nil {
		if err = applyReplayFinalize(ctx, finalize, pending); err != nil {
			return nil, err
		}
		if err = useCase.publishReplayCheckpoint(ctx, checkpoint, replay.trustedCheckpoint); err != nil {
			return nil, err
		}
		return pending, nil
	}
	return useCase.recoverFromRootMismatch(ctx, err, replay.finalizedHead, replay.head, finalize)
}

type pegInAddressRegistryReplay struct {
	head              uint64
	finalizedHead     uint64
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint
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
) (*rootstock.PegInAddressRegistryWatchCheckpoint, error) {
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
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
	finalizedHead uint64,
	head uint64,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) (rootstock.PegInAddressRegistryWatchCheckpoint, error) {
	checkpoint, err := useCase.replay(ctx, trustedCheckpoint, finalizedHead, pending)
	if err != nil {
		return checkpoint, err
	}
	err = useCase.verifyReplayRoot(ctx, checkpoint.LocalRoot, checkpoint.LastProcessedBlock, head)
	return checkpoint, err
}

func (useCase *ReplayRegisteredAddressesUseCase) publishReplayCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
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
	finalize func(context.Context, []*rootstock.PegInAddressRegistryWatchEntry) error,
) ([]*rootstock.PegInAddressRegistryWatchEntry, error) {
	if !errors.Is(replayErr, errPegInAddressRegistryRootMismatch) {
		return nil, replayErr
	}
	useCase.reportRootMismatch(replayErr)
	if discardErr := useCase.discardCheckpoint(ctx); discardErr != nil {
		return nil, errors.Join(replayErr, discardErr)
	}
	useCase.reportResyncStarted("root_mismatch")
	recovered := make([]*rootstock.PegInAddressRegistryWatchEntry, 0)
	checkpoint, err := useCase.replayAndVerify(ctx, nil, finalizedHead, head, &recovered)
	if err == nil {
		if err = applyReplayFinalize(ctx, finalize, recovered); err != nil {
			return nil, err
		}
		if err = useCase.publishCheckpoint(ctx, checkpoint); err != nil {
			return nil, err
		}
		return recovered, nil
	}
	return nil, errors.Join(err, useCase.discardCheckpoint(ctx))
}

func applyReplayFinalize(
	ctx context.Context,
	finalize func(context.Context, []*rootstock.PegInAddressRegistryWatchEntry) error,
	pending []*rootstock.PegInAddressRegistryWatchEntry,
) error {
	if finalize == nil {
		return nil
	}
	return finalize(ctx, pending)
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
	trustedCheckpoint *rootstock.PegInAddressRegistryWatchCheckpoint,
	finalizedHead uint64,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) (rootstock.PegInAddressRegistryWatchCheckpoint, error) {
	fromBlock := useCase.startBlock
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
		if useCase.pageSize-1 <= finalizedHead-fromBlock {
			toBlock = fromBlock + useCase.pageSize - 1
		}
		events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
		if err != nil {
			return checkpoint, fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err)
		}
		events = eventsWithinRange(events, fromBlock, toBlock)
		localRoot, err = useCase.processEvents(ctx, events, localRoot, pending)
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

func (useCase *ReplayRegisteredAddressesUseCase) publishCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInAddressRegistryWatchCheckpoint,
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
		if err = useCase.discoverEvent(ctx, event, pending); err != nil {
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

func (useCase *ReplayRegisteredAddressesUseCase) retryDiscoveredEntries(
	ctx context.Context,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) error {
	entries, err := useCase.getWatched.Run(ctx)
	if err != nil {
		return err
	}
	for index := range entries {
		if entries[index].State != rootstock.PegInAddressRegistryWatchDiscovered {
			continue
		}
		if err = useCase.discoverEvent(ctx, watchEntryToEvent(entries[index]), pending); err != nil {
			return err
		}
	}
	return nil
}

func (useCase *ReplayRegisteredAddressesUseCase) discoverEvent(
	ctx context.Context,
	event blockchain.AddressRegistered,
	pending *[]*rootstock.PegInAddressRegistryWatchEntry,
) error {
	if pendingContains(*pending, event.RskAddress) {
		return nil
	}
	entry, needsRescan, err := useCase.discoverUseCase.Run(ctx, event)
	if err != nil {
		return err
	}
	if entry == nil || !needsRescan {
		return nil
	}
	*pending = append(*pending, entry)
	return nil
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

func (useCase *ReplayRegisteredAddressesUseCase) validateDeployment(ctx context.Context, startBlock uint64) error {
	deploymentVerifier, ok := useCase.registry.(peginAddressRegistryDeploymentVerifier)
	if !ok {
		return errors.New("PegIn address registry adapter does not support exact deployment-block validation")
	}
	isDeploymentBlock, err := deploymentVerifier.IsDeploymentBlock(ctx, startBlock)
	if err != nil {
		return fmt.Errorf("prove PegIn address registry deployment block %d: %w", startBlock, err)
	}
	if !isDeploymentBlock {
		return fmt.Errorf("configured start block %d is not the PegIn address registry deployment block", startBlock)
	}
	deploymentRoot, err := useCase.registry.GetRegistrationRoot(ctx, startBlock)
	if err != nil {
		return fmt.Errorf("validate PegIn address registry at deployment block %d: %w", startBlock, err)
	}
	toBlock := startBlock
	deploymentEvents, err := useCase.registry.GetAddressRegisteredEvents(ctx, startBlock, &toBlock)
	if err != nil {
		return fmt.Errorf("read PegIn address registry events at deployment block %d: %w", startBlock, err)
	}
	replayedDeploymentRoot, err := replayPegInAddressRegistryDeployment(deploymentEvents, startBlock)
	if err != nil {
		return err
	}
	if deploymentRoot != replayedDeploymentRoot {
		return fmt.Errorf("PegIn address registry deployment block %d already has registry state", startBlock)
	}
	return nil
}

func replayPegInAddressRegistryDeployment(
	deploymentEvents []blockchain.AddressRegistered,
	startBlock uint64,
) ([32]byte, error) {
	sort.Slice(deploymentEvents, func(first, second int) bool {
		return deploymentEvents[first].LogIndex < deploymentEvents[second].LogIndex
	})
	replayedDeploymentRoot := [32]byte{}
	for _, event := range deploymentEvents {
		if event.BlockNumber != startBlock {
			return [32]byte{}, fmt.Errorf(
				"PegIn address registry returned block %d while validating deployment block %d",
				event.BlockNumber,
				startBlock,
			)
		}
		var foldErr error
		replayedDeploymentRoot, foldErr = blockchain.FoldPegInAddressRegistryRoot(replayedDeploymentRoot, event.RskAddress)
		if foldErr != nil {
			return [32]byte{}, fmt.Errorf(
				"validate PegIn address registry event at deployment block %d: %w",
				startBlock,
				foldErr,
			)
		}
		if event.RegistrationRoot != replayedDeploymentRoot {
			return [32]byte{}, fmt.Errorf(
				"PegIn address registry event root differs at deployment block %d",
				startBlock,
			)
		}
	}
	return replayedDeploymentRoot, nil
}
