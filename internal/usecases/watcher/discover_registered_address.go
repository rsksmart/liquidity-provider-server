package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type DiscoverRegisteredAddressUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
	registry   blockchain.PegInAddressRegistryContract
	wallet     blockchain.BitcoinWallet
}

func NewDiscoverRegisteredAddressUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	wallet blockchain.BitcoinWallet,
) *DiscoverRegisteredAddressUseCase {
	return &DiscoverRegisteredAddressUseCase{
		repository: repository,
		registry:   registry,
		wallet:     wallet,
	}
}

func (useCase *DiscoverRegisteredAddressUseCase) Run(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (*rootstock.PegInAddressRegistryWatchEntry, bool, error) {
	entry, err := loadOrCreateWatchEntry(ctx, useCase.repository, event, usecases.DiscoverRegisteredAddressId)
	if err != nil {
		return nil, false, err
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
		return entry, false, err
	}
	return entry, needsRescan, nil
}
