package watcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
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
	entry, err := useCase.loadOrCreate(ctx, event)
	if err != nil {
		return nil, false, err
	}
	if entry.State == rootstock.PegInAddressRegistryWatchImported ||
		entry.State == rootstock.PegInAddressRegistryWatchUnsupportedEncoding {
		return entry, false, nil
	}

	pegInAddress, err := useCase.registry.GetPegInAddress(entry.RskAddress)
	if err != nil {
		return entry, false, useCase.RecordError(
			ctx,
			entry,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	entry.Encoding = uint8(pegInAddress.Encoding)
	entry.UpdatedAt = time.Now().UTC()
	if pegInAddress.Encoding != blockchain.PegInAddressRegistryEncodingBase58 {
		entry.State = rootstock.PegInAddressRegistryWatchUnsupportedEncoding
		entry.LastError = ""
		if err = useCase.repository.Update(ctx, *entry); err != nil {
			return nil, false, usecases.WrapUseCaseError(
				usecases.DiscoverRegisteredAddressId,
				fmt.Errorf("persist unsupported encoding for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
			)
		}
		log.Errorf(
			"PegIn address registry event %s/%d uses unsupported encoding %d",
			entry.TxHash,
			entry.LogIndex,
			pegInAddress.Encoding,
		)
		return entry, false, nil
	}
	if entry.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return entry, false, useCase.RecordError(
			ctx,
			entry,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	if err = useCase.wallet.ImportAddress(entry.BtcAddress); err != nil && !isAlreadyImportedError(err) {
		return entry, false, useCase.RecordError(
			ctx,
			entry,
			fmt.Errorf("import PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	return entry, true, nil
}

func (useCase *DiscoverRegisteredAddressUseCase) MarkImported(
	ctx context.Context,
	entry *rootstock.PegInAddressRegistryWatchEntry,
) error {
	entry.State = rootstock.PegInAddressRegistryWatchImported
	entry.LastError = ""
	entry.UpdatedAt = time.Now().UTC()
	if err := useCase.repository.Update(ctx, *entry); err != nil {
		return usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("persist imported state for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	return nil
}

func (useCase *DiscoverRegisteredAddressUseCase) RecordError(
	ctx context.Context,
	entry *rootstock.PegInAddressRegistryWatchEntry,
	entryErr error,
) error {
	log.Error(entryErr)
	if entry.LastError == entryErr.Error() {
		return nil
	}
	entry.LastError = entryErr.Error()
	entry.UpdatedAt = time.Now().UTC()
	if err := useCase.repository.Update(ctx, *entry); err != nil {
		return usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("persist PegIn address registry entry error for %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	return nil
}

func (useCase *DiscoverRegisteredAddressUseCase) loadOrCreate(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (*rootstock.PegInAddressRegistryWatchEntry, error) {
	entry, err := useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
		)
	}
	if entry != nil {
		return entry, nil
	}

	now := time.Now().UTC()
	created := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:           event.TxHash,
		LogIndex:         event.LogIndex,
		BlockNumber:      event.BlockNumber,
		RskAddress:       event.RskAddress,
		Registrant:       event.Registrant,
		RegistrationRoot: event.RegistrationRoot,
		State:            rootstock.PegInAddressRegistryWatchDiscovered,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err = useCase.repository.Upsert(ctx, created); err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("persist AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err),
		)
	}
	entry, err = useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
		)
	}
	if entry == nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: entry not found after upsert", event.RskAddress),
		)
	}
	return entry, nil
}

func isAlreadyImportedError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already imported")
}
