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

func loadOrCreateWatchEntry(
	ctx context.Context,
	repository rootstock.PegInWatchRepository,
	event blockchain.AddressRegistered,
	useCaseId usecases.UseCaseId,
) (*rootstock.PegInWatch, error) {
	entry, err := repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
		)
	}
	if entry != nil {
		return entry, nil
	}

	now := time.Now().UTC()
	created := rootstock.PegInWatch{
		TxHash:           event.TxHash,
		LogIndex:         event.LogIndex,
		BlockNumber:      event.BlockNumber,
		RskAddress:       event.RskAddress,
		Registrant:       event.Registrant,
		RegistrationRoot: event.RegistrationRoot,
		State:            rootstock.PegInWatchDiscovered,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err = repository.Upsert(ctx, created); err != nil {
		return nil, usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("persist AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err),
		)
	}
	entry, err = repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
		)
	}
	if entry == nil {
		return nil, usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("load AddressRegistered event %s: entry not found after upsert", event.RskAddress),
		)
	}
	return entry, nil
}

func resolveAndImportWatchEntry(
	ctx context.Context,
	repository rootstock.PegInWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	wallet blockchain.BitcoinWallet,
	entry *rootstock.PegInWatch,
	useCaseId usecases.UseCaseId,
) (bool, error) {
	if entry.State == rootstock.PegInWatchImported ||
		entry.State == rootstock.PegInWatchUnsupportedEncoding {
		return false, nil
	}

	pegInAddress, err := registry.GetPegInAddress(entry.RskAddress)
	if err != nil {
		return false, persistWatchError(
			ctx,
			repository,
			entry,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
			useCaseId,
		)
	}
	entry.SetEncoding(pegInAddress.Encoding)
	if entry.State == rootstock.PegInWatchUnsupportedEncoding {
		if err = repository.Update(ctx, *entry); err != nil {
			return false, usecases.WrapUseCaseError(
				useCaseId,
				fmt.Errorf("persist unsupported encoding for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
			)
		}
		log.Errorf(
			"PegIn address registry event %s/%d uses unsupported encoding %d",
			entry.TxHash,
			entry.LogIndex,
			pegInAddress.Encoding,
		)
		return false, nil
	}
	if entry.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return false, persistWatchError(
			ctx,
			repository,
			entry,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
			useCaseId,
		)
	}
	if err = wallet.ImportAddress(entry.BtcAddress); err != nil && !isAlreadyImportedError(err) {
		return false, persistWatchError(
			ctx,
			repository,
			entry,
			fmt.Errorf("import PegIn address for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
			useCaseId,
		)
	}
	return true, nil
}

func persistWatchError(
	ctx context.Context,
	repository rootstock.PegInWatchRepository,
	entry *rootstock.PegInWatch,
	entryErr error,
	useCaseId usecases.UseCaseId,
) error {
	log.Error(entryErr)
	if entry.LastError == entryErr.Error() {
		return nil
	}
	entry.LastError = entryErr.Error()
	entry.UpdatedAt = time.Now().UTC()
	if err := repository.Update(ctx, *entry); err != nil {
		return usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("persist PegIn address registry entry error for %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	return nil
}

func persistImported(
	ctx context.Context,
	repository rootstock.PegInWatchRepository,
	entry *rootstock.PegInWatch,
	useCaseId usecases.UseCaseId,
) error {
	entry.State = rootstock.PegInWatchImported
	entry.LastError = ""
	entry.UpdatedAt = time.Now().UTC()
	if err := repository.Update(ctx, *entry); err != nil {
		return usecases.WrapUseCaseError(
			useCaseId,
			fmt.Errorf("persist imported state for event %s/%d: %w", entry.TxHash, entry.LogIndex, err),
		)
	}
	return nil
}

func isAlreadyImportedError(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already imported")
}
