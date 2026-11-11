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
	log "github.com/sirupsen/logrus"
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
	finalizedHead, hasHead, err := useCase.finalizedHead(ctx)
	if err != nil {
		return CatchUpRegisteredAddressesResult{Checkpoint: checkpoint}, useCase.wrap(err)
	}
	toBlock, hasPage := useCase.nextPageEnd(fromBlock, finalizedHead)
	if !hasHead || !hasPage {
		return CatchUpRegisteredAddressesResult{Checkpoint: checkpoint}, nil
	}

	events, err := useCase.registry.GetAddressRegisteredEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("get AddressRegistered events for blocks %d-%d: %w", fromBlock, toBlock, err),
		)
	}
	page, err := useCase.foldEvents(CatchUpRegisteredAddressesResult{
		Watches:    []rootstock.PegInWatch{},
		Checkpoint: rootstock.PegInWatchCheckpoint{Root: checkpoint.Root, Block: toBlock},
	}, events)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(err)
	}
	if err = useCase.checkContractRoot(ctx, page.Checkpoint); err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(err)
	}

	if err = useCase.repository.SaveVerifiedRegistrations(ctx, page.Watches, page.Checkpoint); err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("save PegIn address registry checkpoint at block %d: %w", toBlock, err),
		)
	}
	return page, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) finalizedHead(ctx context.Context) (head uint64, hasHead bool, err error) {
	height, err := useCase.rskRpc.GetHeight(ctx)
	if err != nil {
		return 0, false, fmt.Errorf("get RSK height: %w", err)
	}
	if height < useCase.finalityDepth {
		return 0, false, nil
	}
	return height - useCase.finalityDepth, true, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) foldEvents(
	folded CatchUpRegisteredAddressesResult,
	events []blockchain.AddressRegistered,
) (CatchUpRegisteredAddressesResult, error) {
	useCase.sortAddressRegisteredEvents(events)
	var err error
	for _, event := range events {
		if folded.Checkpoint.Root, err = blockchain.FoldPegInAddressRegistryRoot(useCase.hashFunction, folded.Checkpoint.Root, event.RskAddress); err != nil {
			return CatchUpRegisteredAddressesResult{}, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if folded.Checkpoint.Root != event.RegistrationRoot {
			log.Warn(LogRegistryEventRootMismatch(event.BlockNumber, folded.Checkpoint.Root, event.RegistrationRoot))
		}
		folded.Watches = append(folded.Watches, rootstock.NewPegInWatch(
			event.TxHash,
			event.LogIndex,
			event.BlockNumber,
			event.RskAddress,
			event.Registrant,
			event.RegistrationRoot,
		))
	}
	return folded, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) checkContractRoot(ctx context.Context, checkpoint rootstock.PegInWatchCheckpoint) error {
	contractRoot, err := useCase.registry.GetRegistrationRoot(ctx, checkpoint.Block)
	if err != nil {
		return fmt.Errorf("get PegIn address registry root at block %d: %w", checkpoint.Block, err)
	}
	if checkpoint.Root != contractRoot {
		log.Warn(LogRegistryContractRootMismatch(checkpoint.Block, checkpoint.Root, contractRoot))
	}
	return nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) wrap(err error) error {
	return usecases.WrapUseCaseError(usecases.CatchUpRegisteredAddressesId, err)
}

func (useCase *CatchUpRegisteredAddressesUseCase) nextPageEnd(fromBlock, lastBlock uint64) (toBlock uint64, hasPage bool) {
	if fromBlock > lastBlock {
		return 0, false
	}
	if useCase.pageSize == 0 || lastBlock-fromBlock < useCase.pageSize {
		return lastBlock, true
	}
	return fromBlock + useCase.pageSize - 1, true
}

func (useCase *CatchUpRegisteredAddressesUseCase) sortAddressRegisteredEvents(events []blockchain.AddressRegistered) {
	slices.SortStableFunc(events, func(first, second blockchain.AddressRegistered) int {
		if byBlock := cmp.Compare(first.BlockNumber, second.BlockNumber); byBlock != 0 {
			return byBlock
		}
		return cmp.Compare(first.LogIndex, second.LogIndex)
	})
}
