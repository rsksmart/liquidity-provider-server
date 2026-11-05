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
	root, watches, err := useCase.verifyPage(ctx, checkpoint.Root, events, toBlock)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(err)
	}

	nextCheckpoint := rootstock.PegInWatchCheckpoint{Root: root, Block: toBlock}
	if err = useCase.repository.SaveVerifiedRegistrations(ctx, watches, nextCheckpoint); err != nil {
		return CatchUpRegisteredAddressesResult{}, useCase.wrap(
			fmt.Errorf("save PegIn address registry checkpoint at block %d: %w", toBlock, err),
		)
	}
	return CatchUpRegisteredAddressesResult{Watches: watches, Checkpoint: nextCheckpoint}, nil
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

func (useCase *CatchUpRegisteredAddressesUseCase) verifyPage(
	ctx context.Context,
	root [32]byte,
	events []blockchain.AddressRegistered,
	toBlock uint64,
) ([32]byte, []rootstock.PegInWatch, error) {
	useCase.sortAddressRegisteredEvents(events)
	watches := make([]rootstock.PegInWatch, 0, len(events))
	var err error
	for _, event := range events {
		if root, err = blockchain.FoldPegInAddressRegistryRoot(useCase.hashFunction, root, event.RskAddress); err != nil {
			return [32]byte{}, nil, fmt.Errorf("fold AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
		}
		if root != event.RegistrationRoot {
			log.Warn(registryRootMismatch(event.BlockNumber, root, "event", event.RegistrationRoot))
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

	contractRoot, err := useCase.registry.GetRegistrationRoot(ctx, toBlock)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("get PegIn address registry root at block %d: %w", toBlock, err)
	}
	if root != contractRoot {
		log.Warn(registryRootMismatch(toBlock, root, "contract", contractRoot))
	}
	return root, watches, nil
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

func registryRootMismatch(block uint64, localRoot [32]byte, source string, expectedRoot [32]byte) string {
	return fmt.Sprintf(
		"PegIn address registry root mismatch at block %d: local root 0x%x, %s root 0x%x",
		block,
		localRoot,
		source,
		expectedRoot,
	)
}
