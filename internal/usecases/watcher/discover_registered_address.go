package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type DiscoverRegisteredAddressUseCase struct {
	repository rootstock.PegInWatchRepository
	registry   blockchain.PegInAddressRegistryContract
	wallet     blockchain.BitcoinWallet
}

func NewDiscoverRegisteredAddressUseCase(
	repository rootstock.PegInWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	wallet blockchain.BitcoinWallet,
) *DiscoverRegisteredAddressUseCase {
	return &DiscoverRegisteredAddressUseCase{
		repository: repository,
		registry:   registry,
		wallet:     wallet,
	}
}

type DiscoverRegisteredAddressResult struct {
	Watch       *rootstock.PegInWatch
	NeedsRescan bool
}

func (useCase *DiscoverRegisteredAddressUseCase) Run(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (DiscoverRegisteredAddressResult, error) {
	entry, err := loadOrCreateWatchEntry(ctx, useCase.repository, event, usecases.DiscoverRegisteredAddressId)
	if err != nil {
		return DiscoverRegisteredAddressResult{}, err
	}
	needsRescan, err := resolveAndImportWatchEntry(
		ctx,
		useCase.repository,
		useCase.registry,
		useCase.wallet,
		entry,
		usecases.DiscoverRegisteredAddressId,
	)
	if err != nil {
		return DiscoverRegisteredAddressResult{Watch: entry}, err
	}
	return DiscoverRegisteredAddressResult{Watch: entry, NeedsRescan: needsRescan}, nil
}
