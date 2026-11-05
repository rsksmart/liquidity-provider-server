package watcher

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type CatchUpRegisteredAddressesUseCase struct {
	repository    rootstock.PegInWatchRepository
	registry      blockchain.PegInAddressRegistryContract
	rskRpc        blockchain.RootstockRpcServer
	hashFunction  entities.HashFunction
	startBlock    uint64
	pageSize      uint64
	finalityDepth uint64
}

func NewCatchUpRegisteredAddressesUseCase(
	repository rootstock.PegInWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	hashFunction entities.HashFunction,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *CatchUpRegisteredAddressesUseCase {
	return &CatchUpRegisteredAddressesUseCase{
		repository:    repository,
		registry:      registry,
		rskRpc:        rskRpc,
		hashFunction:  hashFunction,
		startBlock:    startBlock,
		pageSize:      pageSize,
		finalityDepth: finalityDepth,
	}
}

type CatchUpRegisteredAddressesResult struct {
	Watches    []rootstock.PegInWatch
	Checkpoint rootstock.PegInWatchCheckpoint
}

func (useCase *CatchUpRegisteredAddressesUseCase) Run(ctx context.Context) (CatchUpRegisteredAddressesResult, error) {
	checkpoint, found, err := useCase.repository.GetCheckpoint(ctx)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(fmt.Errorf("load PegIn address registry checkpoint: %w", err))
	}
	fromBlock := useCase.startBlock
	if found {
		fromBlock = checkpoint.Block + 1
	}
	toBlock, hasPage, err := useCase.nextPage(ctx, fromBlock)
	if err != nil || !hasPage {
		return CatchUpRegisteredAddressesResult{Checkpoint: checkpoint}, err
	}

	events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err),
		)
	}
	root, watches, err := useCase.foldPage(checkpoint.Root, events)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(err)
	}

	next := rootstock.PegInWatchCheckpoint{Root: root, Block: toBlock}
	if err = useCase.repository.SaveVerifiedRegistrations(ctx, watches, next); err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("save PegIn address registry checkpoint at block %d: %w", toBlock, err),
		)
	}
	return CatchUpRegisteredAddressesResult{Watches: watches, Checkpoint: next}, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) nextPage(ctx context.Context, fromBlock uint64) (uint64, bool, error) {
	head, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return 0, false, useCase.wrap(fmt.Errorf("get RSK height: %w", err))
	}
	if head < useCase.finalityDepth {
		return 0, false, nil
	}
	toBlock, hasPage := NextRegistryPage(fromBlock, head-useCase.finalityDepth, useCase.pageSize)
	return toBlock, hasPage, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) foldPage(
	root [32]byte,
	events []blockchain.AddressRegistered,
) ([32]byte, []rootstock.PegInWatch, error) {
	useCase.sortAddressRegisteredEvents(events)
	watches := make([]rootstock.PegInWatch, 0, len(events))
	var err error
	for _, event := range events {
		if root, err = blockchain.FoldPegInAddressRegistryRoot(useCase.hashFunction, root, event.RskAddress); err != nil {
			return [32]byte{}, nil, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		watches = append(watches, rootstock.NewPegInWatch(
			event.TxHash,
			event.LogIndex,
			event.BlockNumber,
			event.RskAddress,
			event.Registrant,
			event.RegistrationRoot,
		))
	}
	return root, watches, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) wrap(err error) error {
	return usecases.WrapUseCaseError(usecases.CatchUpRegisteredAddressesId, err)
}

func NextRegistryPage(fromBlock, finalizedHead, pageSize uint64) (uint64, bool) {
	if fromBlock > finalizedHead {
		return 0, false
	}
	if pageSize == 0 || finalizedHead-fromBlock < pageSize {
		return finalizedHead, true
	}
	return fromBlock + pageSize - 1, true
}

func (useCase *CatchUpRegisteredAddressesUseCase) sortAddressRegisteredEvents(events []blockchain.AddressRegistered) {
	slices.SortStableFunc(events, func(first, second blockchain.AddressRegistered) int {
		if byBlock := cmp.Compare(first.BlockNumber, second.BlockNumber); byBlock != 0 {
			return byBlock
		}
		return cmp.Compare(first.LogIndex, second.LogIndex)
	})
}
