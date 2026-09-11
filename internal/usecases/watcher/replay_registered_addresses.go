package watcher

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type ReplayRegisteredAddressesUseCase struct {
	repository   rootstock.PegInWatchRepository
	registry     blockchain.PegInAddressRegistryContract
	rskRpc       blockchain.RootstockRpcServer
	hashFunction entities.HashFunction
	startBlock   uint64
	pageSize     uint64
}

type ReplayRootMismatch struct {
	BlockNumber uint64
	LocalRoot   [32]byte
	ChainRoot   [32]byte
	Source      string
}

type ReplayResult struct {
	Reason     blockchain.PegInAddressRegistryRecoveryReason
	Checkpoint *rootstock.PegInWatchCheckpoint
	Mismatch   *ReplayRootMismatch
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
		"%v at %s block %d: local=%x chain=%x",
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
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	hashFunction entities.HashFunction,
	startBlock uint64,
	pageSize uint64,
) (*ReplayRegisteredAddressesUseCase, error) {
	if pageSize == 0 {
		return nil, errors.New("replay page size must be greater than zero")
	}
	return &ReplayRegisteredAddressesUseCase{
		repository:   watches,
		registry:     registry,
		rskRpc:       rskRpc,
		hashFunction: hashFunction,
		startBlock:   startBlock,
		pageSize:     pageSize,
	}, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) Run(
	ctx context.Context,
	verifiedCheckpoint *rootstock.PegInWatchCheckpoint,
) (ReplayResult, error) {
	result, err := useCase.runReplay(ctx, verifiedCheckpoint)
	if err != nil {
		return ReplayResult{}, usecases.WrapUseCaseError(usecases.ReplayRegisteredAddressesId, err)
	}
	return result, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) runReplay(
	ctx context.Context,
	verifiedCheckpoint *rootstock.PegInWatchCheckpoint,
) (ReplayResult, error) {
	head, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return ReplayResult{}, fmt.Errorf("get RSK height: %w", err)
	}
	if head < useCase.startBlock {
		return ReplayResult{}, nil
	}
	state, err := useCase.loadInitialState(ctx, head, verifiedCheckpoint)
	if err != nil {
		return ReplayResult{}, err
	}
	if state.localRootAtHead == state.chainRootAtHead {
		return ReplayResult{
			Checkpoint: &rootstock.PegInWatchCheckpoint{
				LocalRoot:            state.chainRootAtHead,
				VerifiedThroughBlock: state.head,
			},
		}, nil
	}

	plan, err := useCase.planReplay(ctx, state)
	if err != nil {
		return ReplayResult{}, err
	}
	result, err := useCase.rebuildFromPlan(ctx, state, plan)
	if err != nil {
		return ReplayResult{}, err
	}
	result.Reason = plan.reason
	if plan.reason == blockchain.PegInAddressRegistryRecoveryRootMismatch {
		result.Mismatch = &ReplayRootMismatch{
			BlockNumber: state.head,
			LocalRoot:   state.localRootAtHead,
			ChainRoot:   state.chainRootAtHead,
			Source:      "captured_head",
		}
	}
	return result, nil
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
	verifiedCheckpoint *rootstock.PegInWatchCheckpoint,
) (initialReconciliation, error) {
	entries, err := useCase.repository.List(ctx)
	if err != nil {
		return initialReconciliation{}, fmt.Errorf("list PegIn watches: %w", err)
	}
	chainRootAtHead, err := useCase.registry.GetRegistrationRoot(ctx, head)
	if err != nil {
		return initialReconciliation{}, fmt.Errorf(
			"get PegIn address registry root at block %d: %w",
			head,
			err,
		)
	}
	timeline, err := NewLocalRootTimeline(entries, useCase.startBlock, head, useCase.hashFunction)
	if err != nil {
		return initialReconciliation{}, err
	}
	return initialReconciliation{
		head:            head,
		checkpoint:      verifiedCheckpoint,
		timeline:        timeline,
		localRootAtHead: timeline.RootAt(head),
		chainRootAtHead: chainRootAtHead,
	}, nil
}

type replayPlan struct {
	fromBlock uint64
	seedRoot  [32]byte
	reason    blockchain.PegInAddressRegistryRecoveryReason
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
			fromBlock: state.checkpoint.VerifiedThroughBlock + 1,
			seedRoot:  state.checkpoint.LocalRoot,
			reason:    blockchain.PegInAddressRegistryRecoveryCatchUp,
		}, nil
	}
	plan := replayPlan{reason: blockchain.PegInAddressRegistryRecoveryRootMismatch}
	fromBlock, err := findFirstMismatchingBlock(
		ctx,
		func(probeCtx context.Context, height uint64) (bool, error) {
			return useCase.localMatchesChainAt(probeCtx, state.timeline, height)
		},
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
	checkpointBlock := state.checkpoint.VerifiedThroughBlock
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
) (ReplayResult, error) {
	if plan.fromBlock > state.head {
		return ReplayResult{}, fmt.Errorf(
			"replay start block %d is above captured head %d",
			plan.fromBlock,
			state.head,
		)
	}
	events, err := useCase.fetchAndValidateReplayEvents(
		ctx,
		plan.fromBlock,
		state.head,
		plan.seedRoot,
	)
	if err != nil {
		return ReplayResult{}, err
	}
	watches := useCase.replayWatches(events)
	if err = useCase.repository.ReplaceFromBlock(ctx, plan.fromBlock, watches); err != nil {
		return ReplayResult{}, fmt.Errorf("replace PegIn watches from block %d: %w", plan.fromBlock, err)
	}
	if err = useCase.verifyPersistedState(ctx, state); err != nil {
		return ReplayResult{}, err
	}
	checkpoint := rootstock.PegInWatchCheckpoint{
		LocalRoot:            state.chainRootAtHead,
		VerifiedThroughBlock: state.head,
	}
	return ReplayResult{Checkpoint: &checkpoint}, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) verifyPersistedState(
	ctx context.Context,
	state initialReconciliation,
) error {
	rebuiltEntries, err := useCase.repository.List(ctx)
	if err != nil {
		return fmt.Errorf("list rebuilt PegIn watches: %w", err)
	}
	rebuiltTimeline, err := NewLocalRootTimeline(
		rebuiltEntries,
		useCase.startBlock,
		state.head,
		useCase.hashFunction,
	)
	if err != nil {
		return err
	}
	rebuiltRoot := rebuiltTimeline.RootAt(state.head)
	if rebuiltRoot == state.chainRootAtHead {
		return nil
	}
	mismatchErr := &pegInAddressRegistryRootMismatchError{
		blockNumber: state.head,
		localRoot:   rebuiltRoot,
		chainRoot:   state.chainRootAtHead,
		source:      "persisted_replay",
	}
	return mismatchErr
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

func (useCase *ReplayRegisteredAddressesUseCase) fetchAndValidateReplayEvents(
	ctx context.Context,
	fromBlock uint64,
	head uint64,
	seedRoot [32]byte,
) ([]blockchain.AddressRegistered, error) {
	events := make([]blockchain.AddressRegistered, 0)
	if fromBlock > head {
		return nil, fmt.Errorf(
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
		pageEvents, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
		if err != nil {
			return nil, fmt.Errorf(
				"get AddressRegistered events for blocks %d-%d: %w",
				fromBlock,
				toBlock,
				err,
			)
		}
		events = append(events, pageEvents...)
		if toBlock == head {
			break
		}
		fromBlock = toBlock + 1
	}
	events = orderedUniqueEvents(events)
	for _, event := range events {
		var err error
		localRoot, err = blockchain.FoldPegInAddressRegistryRoot(useCase.hashFunction, localRoot, event.RskAddress)
		if err != nil {
			return nil, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if event.RegistrationRoot != localRoot {
			return nil, &pegInAddressRegistryRootMismatchError{
				blockNumber: event.BlockNumber,
				localRoot:   localRoot,
				chainRoot:   event.RegistrationRoot,
				source:      fmt.Sprintf("event_%s_%d", event.TxHash, event.LogIndex),
			}
		}
	}
	return events, nil
}

func (useCase *ReplayRegisteredAddressesUseCase) replayWatches(
	events []blockchain.AddressRegistered,
) []rootstock.PegInWatch {
	watches := make([]rootstock.PegInWatch, 0, len(events))
	for _, event := range events {
		now := time.Now().UTC()
		watches = append(watches, rootstock.PegInWatch{
			TxHash:           event.TxHash,
			LogIndex:         event.LogIndex,
			BlockNumber:      event.BlockNumber,
			RskAddress:       event.RskAddress,
			Registrant:       event.Registrant,
			RegistrationRoot: event.RegistrationRoot,
			State:            rootstock.PegInWatchDiscovered,
			CreatedAt:        now,
			UpdatedAt:        now,
		})
	}
	return watches
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
