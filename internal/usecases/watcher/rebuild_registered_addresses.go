package watcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	log "github.com/sirupsen/logrus"
)

var RegistryOldBlockReadError = errors.New("cannot read the PegIn address registry root at old block")

func (useCase *CatchUpRegisteredAddressesUseCase) rebuild(
	ctx context.Context,
	checkpoint rootstock.PegInWatchCheckpoint,
	finalizedHead uint64,
) (CatchUpRegisteredAddressesResult, error) {
	contractRoot, err := useCase.contractRootAt(ctx, checkpoint.Block)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, err
	}
	if contractRoot == checkpoint.Root {
		return useCase.rebuildFrom(ctx, checkpoint.Block+1, checkpoint.Root, finalizedHead)
	}

	fromBlock, err := useCase.firstMismatchingBlock(ctx, useCase.startBlock, checkpoint.Block)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, fmt.Errorf("find the first mismatching PegIn address registry block: %w", err)
	}
	return useCase.rebuildFromStoredRoot(ctx, fromBlock, finalizedHead)
}

func (useCase *CatchUpRegisteredAddressesUseCase) rebuildFromStoredRoot(
	ctx context.Context,
	fromBlock uint64,
	finalizedHead uint64,
) (CatchUpRegisteredAddressesResult, error) {
	root, err := useCase.storedRootBefore(ctx, fromBlock)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, err
	}
	return useCase.rebuildFrom(ctx, fromBlock, root, finalizedHead)
}

func (useCase *CatchUpRegisteredAddressesUseCase) rebuildFrom(
	ctx context.Context,
	fromBlock uint64,
	root [32]byte,
	finalizedHead uint64,
) (CatchUpRegisteredAddressesResult, error) {
	rebuilt, err := useCase.readVerifiedRange(ctx, fromBlock, finalizedHead, root)
	if err != nil {
		return CatchUpRegisteredAddressesResult{}, fmt.Errorf("rebuild the PegIn watch list from block %d: %w", fromBlock, err)
	}
	if err = useCase.repository.ReplaceVerifiedRegistrations(ctx, fromBlock, rebuilt.Watches, rebuilt.Checkpoint); err != nil {
		return CatchUpRegisteredAddressesResult{}, fmt.Errorf("save the PegIn watch list rebuilt from block %d: %w", fromBlock, err)
	}
	log.Info(LogRegistryWatchListRebuilt(fromBlock, finalizedHead, len(rebuilt.Watches)))
	rebuilt.Rebuilt = true
	return rebuilt, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) firstMismatchingBlock(ctx context.Context, low, high uint64) (uint64, error) {
	for low < high {
		middle := low + (high-low)/2
		matches, err := useCase.storedRootMatchesAt(ctx, middle)
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

func (useCase *CatchUpRegisteredAddressesUseCase) storedRootMatchesAt(ctx context.Context, block uint64) (bool, error) {
	storedRoot, err := useCase.storedRootAt(ctx, block)
	if err != nil {
		return false, err
	}
	contractRoot, err := useCase.contractRootAt(ctx, block)
	if err != nil {
		return false, err
	}
	return storedRoot == contractRoot, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) storedRootBefore(ctx context.Context, block uint64) ([32]byte, error) {
	if block == 0 {
		return [32]byte{}, nil
	}
	return useCase.storedRootAt(ctx, block-1)
}

func (useCase *CatchUpRegisteredAddressesUseCase) storedRootAt(ctx context.Context, block uint64) ([32]byte, error) {
	root, err := useCase.repository.GetRootAt(ctx, block)
	if err != nil {
		return [32]byte{}, fmt.Errorf("get stored PegIn watch root at block %d: %w", block, err)
	}
	return root, nil
}

func (useCase *CatchUpRegisteredAddressesUseCase) contractRootAt(ctx context.Context, block uint64) ([32]byte, error) {
	root, err := useCase.registry.GetRegistrationRoot(ctx, block)
	if err != nil {
		return [32]byte{}, fmt.Errorf("%w %d, check that the RSK node can answer calls at old blocks: %w", RegistryOldBlockReadError, block, err)
	}
	return root, nil
}
