package watcher

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type CatchUpRegisteredAddressesUseCase struct {
	repository    rootstock.PegInWatchRepository
	registry      blockchain.PegInAddressRegistryContract
	rskRpc        blockchain.RootstockRpcServer
	startBlock    uint64
	pageSize      uint64
	finalityDepth uint64
}

func NewCatchUpRegisteredAddressesUseCase(
	repository rootstock.PegInWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
	startBlock uint64,
	pageSize uint64,
	finalityDepth uint64,
) *CatchUpRegisteredAddressesUseCase {
	return &CatchUpRegisteredAddressesUseCase{
		repository:    repository,
		registry:      registry,
		rskRpc:        rskRpc,
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
	head, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(fmt.Errorf("get RSK height: %w", err))
	}
	unchanged := CatchUpRegisteredAddressesResult{Checkpoint: checkpoint}
	if head < useCase.finalityDepth {
		return unchanged, nil
	}
	fromBlock, toBlock := useCase.nextRange(checkpoint.Block, found, head-useCase.finalityDepth)
	if fromBlock > toBlock {
		return unchanged, nil
	}

	events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err),
		)
	}
	SortAddressRegisteredEvents(events)

	next := rootstock.PegInWatchCheckpoint{Root: checkpoint.Root, Block: toBlock}
	watches := make([]rootstock.PegInWatch, 0, len(events))
	for _, event := range events {
		next.Root = event.RegistrationRoot
		watches = append(watches, rootstock.NewPegInWatch(
			event.TxHash,
			event.LogIndex,
			event.BlockNumber,
			event.RskAddress,
			event.Registrant,
			event.RegistrationRoot,
		))
	}
	if err = useCase.repository.SaveVerifiedRegistrations(ctx, watches, next); err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("save PegIn address registry checkpoint at block %d: %w", toBlock, err),
		)
	}
	return CatchUpRegisteredAddressesResult{Watches: watches, Checkpoint: next}, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) nextRange(lastBlock uint64, found bool, finalizedHead uint64) (uint64, uint64) {
	if !found {
		if useCase.startBlock > finalizedHead || useCase.pageSize-1 > finalizedHead-useCase.startBlock {
			return useCase.startBlock, finalizedHead
		}
		return useCase.startBlock, useCase.startBlock + useCase.pageSize - 1
	}

	fromBlock := useCase.startBlock
	if lastBlock+1 > useCase.finalityDepth {
		fromBlock = max(useCase.startBlock, lastBlock+1-useCase.finalityDepth)
	}
	if fromBlock > finalizedHead {
		return fromBlock, finalizedHead
	}
	if lastBlock >= finalizedHead || useCase.pageSize > finalizedHead-lastBlock {
		return fromBlock, finalizedHead
	}
	return fromBlock, lastBlock + useCase.pageSize
}

func (useCase *CatchUpRegisteredAddressesUseCase) wrap(err error) error {
	return usecases.WrapUseCaseError(usecases.CatchUpRegisteredAddressesId, err)
}

func SortAddressRegisteredEvents(events []blockchain.AddressRegistered) {
	slices.SortStableFunc(events, func(first, second blockchain.AddressRegistered) int {
		if byBlock := cmp.Compare(first.BlockNumber, second.BlockNumber); byBlock != 0 {
			return byBlock
		}
		return cmp.Compare(first.LogIndex, second.LogIndex)
	})
}
