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
	repository   rootstock.PegInWatchRepository
	checkpoints  rootstock.PegInWatchCheckpointRepository
	registry     blockchain.PegInAddressRegistryContract
	rskRpc       blockchain.RootstockRpcServer
	eventBus     entities.EventBus
	wallet       blockchain.BitcoinWallet
	hashFunction entities.HashFunction
	startBlock   uint64
	pageSize     uint64
	replayMutex  sync.Mutex
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
	watches rootstock.PegInWatchRepository,
	checkpoints rootstock.PegInWatchCheckpointRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	eventBus entities.EventBus,
	wallet blockchain.BitcoinWallet,
	hashFunction entities.HashFunction,
) *ReplayRegisteredAddressesUseCase {
	return &ReplayRegisteredAddressesUseCase{
		repository:   watches,
		checkpoints:  checkpoints,
		registry:     registry,
		rskRpc:       rskRpc,
		eventBus:     eventBus,
		wallet:       wallet,
		hashFunction: hashFunction,
	}
}

func (useCase *ReplayRegisteredAddressesUseCase) Run(
	ctx context.Context,
	startBlock uint64,
	pageSize uint64,
) ([]*rootstock.PegInWatch, error) {
	useCase.replayMutex.Lock()
	defer useCase.replayMutex.Unlock()

	useCase.startBlock = startBlock
	useCase.pageSize = pageSize

	if pageSize == 0 {
		return nil, usecases.WrapUseCaseError(
			usecases.ReplayRegisteredAddressesId,
			errors.New("replay page size must be greater than zero"),
		)
	}
	pending, err := useCase.runReplay(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.ReplayRegisteredAddressesId, err)
	}
	return pending, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) runReplay(
	ctx context.Context,
) ([]*rootstock.PegInWatch, error) {
	head, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("get RSK height: %w", err)
	}
	if head < useCase.startBlock {
		return []*rootstock.PegInWatch{}, nil
	}
	state, err := useCase.loadInitialState(ctx, head)
	if err != nil {
		return nil, err
	}
	if state.localRootAtHead == state.chainRootAtHead {
		return useCase.completeCurrentState(ctx, state)
	}

	plan, err := useCase.planReplay(ctx, state)
	if plan.resyncReason == "root_mismatch" {
		useCase.reportRootMismatch(&pegInAddressRegistryRootMismatchError{
			blockNumber: state.head,
			localRoot:   state.localRootAtHead,
			chainRoot:   state.chainRootAtHead,
			source:      "captured_head",
		})
	}
	if err != nil {
		return nil, err
	}
	return useCase.rebuildFromPlan(ctx, state, plan)
}

type initialReconciliation struct {
	head            uint64
	checkpoint      *rootstock.PegInWatchCheckpoint
	timeline        localRootTimeline
	localRootAtHead [32]byte
	chainRootAtHead [32]byte
}

func (useCase *ReplayRegisteredAddressesUseCase) loadInitialState(
	ctx context.Context,
	head uint64,
) (initialReconciliation, error) {
	checkpoint, err := useCase.checkpoints.GetCheckpoint(ctx)
	if err != nil {
		return initialReconciliation{}, fmt.Errorf("load PegIn address registry checkpoint: %w", err)
	}
	entries, err := useCase.repository.List(ctx)
	if err != nil {
		return initialReconciliation{}, fmt.Errorf("list PegIn watches: %w", err)
	}
	timeline, err := buildLocalRootTimeline(entries, useCase.startBlock, head, useCase.hashFunction)
	if err != nil {
		return initialReconciliation{}, err
	}
	if timeline.hasBeforeStart {
		return initialReconciliation{}, fmt.Errorf(
			"PegIn watch exists below configured start block %d",
			useCase.startBlock,
		)
	}
	chainRootAtHead, err := useCase.registry.GetRegistrationRoot(ctx, head)
	if err != nil {
		return initialReconciliation{}, fmt.Errorf(
			"get PegIn address registry root at block %d: %w",
			head,
			err,
		)
	}
	return initialReconciliation{
		head:            head,
		checkpoint:      checkpoint,
		timeline:        timeline,
		localRootAtHead: timeline.RootAt(head),
		chainRootAtHead: chainRootAtHead,
	}, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) completeCurrentState(
	ctx context.Context,
	state initialReconciliation,
) ([]*rootstock.PegInWatch, error) {
	if state.timeline.hasAfterHead {
		return nil, fmt.Errorf("PegIn watch exists above current head %d", state.head)
	}
	finalCheckpoint := rootstock.PegInWatchCheckpoint{
		LocalRoot:          state.chainRootAtHead,
		LastProcessedBlock: state.head,
	}
	if state.checkpoint == nil || *state.checkpoint != finalCheckpoint {
		if err := useCase.publishCheckpoint(ctx, finalCheckpoint); err != nil {
			return nil, err
		}
	}
	return useCase.retryDiscoveredEntries(ctx, []*rootstock.PegInWatch{})
}

type replayPlan struct {
	fromBlock    uint64
	seedRoot     [32]byte
	resyncReason string
}

func (useCase *ReplayRegisteredAddressesUseCase) planReplay(
	ctx context.Context,
	state initialReconciliation,
) (replayPlan, error) {
	trusted, err := useCase.checkpointSurvivesVerification(ctx, state)
	if err != nil {
		return replayPlan{}, err
	}
	if trusted {
		return replayPlan{
			fromBlock:    state.checkpoint.LastProcessedBlock + 1,
			seedRoot:     state.checkpoint.LocalRoot,
			resyncReason: "catch_up",
		}, nil
	}
	plan := replayPlan{resyncReason: "root_mismatch"}
	fromBlock, err := useCase.findFirstMismatchingBlock(
		ctx,
		state.timeline,
		useCase.startBlock,
		state.head,
	)
	if err != nil {
		return plan, err
	}
	plan.fromBlock = fromBlock
	if fromBlock > useCase.startBlock {
		plan.seedRoot = state.timeline.RootAt(fromBlock - 1)
	}
	return plan, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) checkpointSurvivesVerification(
	ctx context.Context,
	state initialReconciliation,
) (bool, error) {
	if state.checkpoint == nil {
		return false, nil
	}
	checkpointBlock := state.checkpoint.LastProcessedBlock
	if checkpointBlock < useCase.startBlock || checkpointBlock >= state.head {
		return false, nil
	}
	if state.timeline.RootAt(checkpointBlock) != state.checkpoint.LocalRoot {
		return false, nil
	}
	chainRoot, err := useCase.registry.GetRegistrationRoot(ctx, checkpointBlock)
	if err != nil {
		return false, fmt.Errorf(
			"get PegIn address registry root at checkpoint block %d: %w",
			checkpointBlock,
			err,
		)
	}
	return chainRoot == state.checkpoint.LocalRoot, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) rebuildFromPlan(
	ctx context.Context,
	state initialReconciliation,
	plan replayPlan,
) ([]*rootstock.PegInWatch, error) {
	if plan.fromBlock > state.head {
		return nil, fmt.Errorf(
			"replay start block %d is above captured head %d",
			plan.fromBlock,
			state.head,
		)
	}
	useCase.reportResyncStarted(plan.resyncReason)
	if err := useCase.repository.DeleteFromBlock(ctx, plan.fromBlock); err != nil {
		return nil, fmt.Errorf("delete PegIn watches from block %d: %w", plan.fromBlock, err)
	}
	pending, err := useCase.replayRange(ctx, plan.fromBlock, state.head, plan.seedRoot)
	if err != nil {
		return nil, err
	}
	if err = useCase.verifyPersistedState(ctx, state); err != nil {
		return nil, err
	}
	if err = useCase.publishCheckpoint(ctx, rootstock.PegInWatchCheckpoint{
		LocalRoot:          state.chainRootAtHead,
		LastProcessedBlock: state.head,
	}); err != nil {
		return nil, err
	}
	return useCase.retryDiscoveredEntries(ctx, pending)
}

func (useCase *ReplayRegisteredAddressesUseCase) verifyPersistedState(
	ctx context.Context,
	state initialReconciliation,
) error {
	rebuiltEntries, err := useCase.repository.List(ctx)
	if err != nil {
		return fmt.Errorf("list rebuilt PegIn watches: %w", err)
	}
	rebuiltTimeline, err := buildLocalRootTimeline(
		rebuiltEntries,
		useCase.startBlock,
		state.head,
		useCase.hashFunction,
	)
	if err != nil {
		return err
	}
	rebuiltRoot := rebuiltTimeline.RootAt(state.head)
	if !rebuiltTimeline.hasBeforeStart &&
		!rebuiltTimeline.hasAfterHead &&
		rebuiltRoot == state.chainRootAtHead {
		return nil
	}
	mismatchErr := &pegInAddressRegistryRootMismatchError{
		blockNumber: state.head,
		localRoot:   rebuiltRoot,
		chainRoot:   state.chainRootAtHead,
		source:      "persisted_replay",
	}
	useCase.reportRootMismatch(mismatchErr)
	return mismatchErr
}

// RegistrationRoot is intentionally absent. A stored row cannot prove its own integrity.
type localRegistration struct {
	blockNumber uint64
	logIndex    uint
	txHash      string
	rskAddress  string
}

type localRootPoint struct {
	blockNumber uint64
	root        [32]byte
}

type localRootTimeline struct {
	points         []localRootPoint
	hasRows        bool
	hasBeforeStart bool
	hasAfterHead   bool
}

func collectLocalRegistrations(
	entries []rootstock.PegInWatch,
	startBlock uint64,
	head uint64,
) ([]localRegistration, localRootTimeline) {
	registrations := make([]localRegistration, 0, len(entries))
	timeline := localRootTimeline{hasRows: len(entries) != 0}
	for _, entry := range entries {
		if entry.BlockNumber < startBlock {
			timeline.hasBeforeStart = true
		}
		if entry.BlockNumber > head {
			timeline.hasAfterHead = true
		}
		registrations = append(registrations, localRegistration{
			blockNumber: entry.BlockNumber,
			logIndex:    entry.LogIndex,
			txHash:      entry.TxHash,
			rskAddress:  entry.RskAddress,
		})
	}
	return registrations, timeline
}

func localRegistrationBefore(first localRegistration, second localRegistration) bool {
	if first.blockNumber != second.blockNumber {
		return first.blockNumber < second.blockNumber
	}
	if first.logIndex != second.logIndex {
		return first.logIndex < second.logIndex
	}
	if first.txHash != second.txHash {
		return first.txHash < second.txHash
	}
	return first.rskAddress < second.rskAddress
}

func buildLocalRootTimeline(
	entries []rootstock.PegInWatch,
	startBlock uint64,
	head uint64,
	hashFunction entities.HashFunction,
) (localRootTimeline, error) {
	registrations, timeline := collectLocalRegistrations(entries, startBlock, head)
	sort.Slice(registrations, func(firstIndex, secondIndex int) bool {
		return localRegistrationBefore(registrations[firstIndex], registrations[secondIndex])
	})

	root := [32]byte{}
	for _, registration := range registrations {
		if registration.blockNumber < startBlock || registration.blockNumber > head {
			continue
		}
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(hashFunction, root, registration.rskAddress)
		if err != nil {
			return localRootTimeline{}, fmt.Errorf(
				"fold stored AddressRegistered event %s/%d: %w",
				registration.txHash,
				registration.logIndex,
				err,
			)
		}
		point := localRootPoint{blockNumber: registration.blockNumber, root: root}
		lastPoint := len(timeline.points) - 1
		if lastPoint >= 0 && timeline.points[lastPoint].blockNumber == registration.blockNumber {
			timeline.points[lastPoint] = point
		} else {
			timeline.points = append(timeline.points, point)
		}
	}
	return timeline, nil
}

func (timeline localRootTimeline) RootAt(height uint64) [32]byte {
	index := sort.Search(len(timeline.points), func(index int) bool {
		return timeline.points[index].blockNumber > height
	})
	if index == 0 {
		return [32]byte{}
	}
	return timeline.points[index-1].root
}

func (useCase *ReplayRegisteredAddressesUseCase) localMatchesChainAt(
	ctx context.Context,
	timeline localRootTimeline,
	height uint64,
) (bool, error) {
	chainRoot, err := useCase.registry.GetRegistrationRoot(ctx, height)
	if err != nil {
		return false, fmt.Errorf("get PegIn address registry root at block %d: %w", height, err)
	}
	return timeline.RootAt(height) == chainRoot, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) findFirstMismatchingBlock(
	ctx context.Context,
	timeline localRootTimeline,
	startBlock uint64,
	head uint64,
) (uint64, error) {
	low, high := startBlock, head
	for low < high {
		middle := low + (high-low)/2
		matches, err := useCase.localMatchesChainAt(ctx, timeline, middle)
		if err != nil {
			return 0, err
		}
		if matches {
			low = middle + 1
		} else {
			high = middle
		}
	}
	return low, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) replayRange(
	ctx context.Context,
	fromBlock uint64,
	head uint64,
	seedRoot [32]byte,
) ([]*rootstock.PegInWatch, error) {
	pending := []*rootstock.PegInWatch{}
	if fromBlock > head {
		return pending, fmt.Errorf(
			"replay start block %d is above captured head %d",
			fromBlock,
			head,
		)
	}
	localRoot := seedRoot
	for fromBlock <= head {
		toBlock := head
		if useCase.pageSize-1 <= head-fromBlock {
			toBlock = fromBlock + useCase.pageSize - 1
		}
		events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
		if err != nil {
			return pending, fmt.Errorf(
				"get AddressRegistered events for blocks %d-%d: %w",
				fromBlock,
				toBlock,
				err,
			)
		}
		localRoot, pending, err = useCase.processEvents(ctx, events, localRoot, pending)
		if err != nil {
			return pending, err
		}
		if toBlock == head {
			break
		}
		fromBlock = toBlock + 1
	}
	return pending, nil
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

func (useCase *ReplayRegisteredAddressesUseCase) publishCheckpoint(
	ctx context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
) error {
	if err := useCase.checkpoints.SetCheckpoint(ctx, checkpoint); err != nil {
		return fmt.Errorf("persist PegIn address registry checkpoint: %w", err)
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
		localRoot, err = blockchain.FoldPegInAddressRegistryRoot(useCase.hashFunction, localRoot, event.RskAddress)
		if err != nil {
			return [32]byte{}, pending, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if event.RegistrationRoot != localRoot {
			mismatchErr := &pegInAddressRegistryRootMismatchError{
				blockNumber: event.BlockNumber,
				localRoot:   localRoot,
				chainRoot:   event.RegistrationRoot,
				source:      fmt.Sprintf("event_%s_%d", event.TxHash, event.LogIndex),
			}
			useCase.reportRootMismatch(mismatchErr)
			return [32]byte{}, pending, mismatchErr
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
