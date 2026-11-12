package watcher

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

const registeredAddressRescanDepthBlocks int64 = 100

type ImportRegisteredAddressesUseCase struct {
	repository rootstock.PegInWatchRepository
	registry   blockchain.PegInAddressRegistryContract
	btcNetwork blockchain.BitcoinNetwork
	wallet     blockchain.BitcoinWallet
}

func NewImportRegisteredAddressesUseCase(
	repository rootstock.PegInWatchRepository,
	registry blockchain.PegInAddressRegistryContract,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
) *ImportRegisteredAddressesUseCase {
	return &ImportRegisteredAddressesUseCase{
		repository: repository,
		registry:   registry,
		btcNetwork: btcNetwork,
		wallet:     wallet,
	}
}

type ImportRegisteredAddressesResult struct {
	Watches []rootstock.PegInWatch
}

func (useCase *ImportRegisteredAddressesUseCase) Run(
	ctx context.Context,
	watches []rootstock.PegInWatch,
) (ImportRegisteredAddressesResult, error) {
	updated, err := useCase.importDiscovered(ctx, watches)
	result := ImportRegisteredAddressesResult{Watches: updated}
	if err != nil {
		return result, usecases.WrapUseCaseError(usecases.ImportRegisteredAddressesId, err)
	}
	return result, nil
}

func (useCase *ImportRegisteredAddressesUseCase) importDiscovered(
	ctx context.Context,
	watches []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	updated := make([]rootstock.PegInWatch, 0)
	pending := make([]rootstock.PegInWatch, 0)
	for index := range watches {
		if watches[index].State == rootstock.PegInWatchDiscovered {
			watch, needsRescan, err := useCase.importAddress(ctx, blockchain.NewAddressRegisteredFromWatchEntry(watches[index]))
			if err != nil {
				return append(updated, pending...), err
			}
			if needsRescan {
				pending = append(pending, watch)
			} else {
				updated = append(updated, watch)
			}
		}
	}
	rescanned, err := useCase.rescan(ctx, pending)
	return append(updated, rescanned...), err
}

func (useCase *ImportRegisteredAddressesUseCase) importAddress(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (updated rootstock.PegInWatch, needsRescan bool, err error) {
	watch, err := useCase.loadOrCreate(ctx, event)
	if err != nil {
		return rootstock.PegInWatch{}, false, err
	}
	if watch.State == rootstock.PegInWatchImported ||
		watch.State == rootstock.PegInWatchUnsupportedEncoding {
		return *watch, false, nil
	}

	pegInAddress, err := useCase.registry.GetPegInAddress(watch.RskAddress)
	if err != nil {
		return *watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	watch.SetEncoding(pegInAddress.Encoding)
	if watch.State == rootstock.PegInWatchUnsupportedEncoding {
		if err = useCase.repository.Update(ctx, *watch); err != nil {
			return rootstock.PegInWatch{}, false, fmt.Errorf("persist unsupported encoding for event %s/%d: %w", watch.TxHash, watch.LogIndex, err)
		}
		log.Errorf(
			"PegIn address registry event %s/%d uses unsupported encoding %d",
			watch.TxHash,
			watch.LogIndex,
			pegInAddress.Encoding,
		)
		return *watch, false, nil
	}
	if watch.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return *watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	if err = useCase.wallet.ImportAddress(watch.BtcAddress); err != nil && !btcclient.IsAddressAlreadyImported(err) {
		return *watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("import PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return *watch, true, nil
}

func (useCase *ImportRegisteredAddressesUseCase) loadOrCreate(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (*rootstock.PegInWatch, error) {
	watch, err := useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err)
	}
	if watch != nil {
		return watch, nil
	}

	created := rootstock.NewPegInWatch(
		event.TxHash,
		event.LogIndex,
		event.BlockNumber,
		event.RskAddress,
		event.Registrant,
		event.RegistrationRoot,
	)
	if err = useCase.repository.Upsert(ctx, created); err != nil {
		return nil, fmt.Errorf("persist AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err)
	}
	watch, err = useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err)
	}
	if watch == nil {
		return nil, fmt.Errorf("load AddressRegistered event %s: watch not found after upsert", event.RskAddress)
	}
	return watch, nil
}

func (useCase *ImportRegisteredAddressesUseCase) rescan(
	ctx context.Context,
	pending []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	if len(pending) == 0 {
		return pending, nil
	}
	tip, err := useCase.btcNetwork.GetHeight()
	if err != nil {
		return useCase.recordRescanError(ctx, pending, fmt.Errorf("get BTC height for registry rescan: %w", err))
	}
	fromHeight := max(tip.Int64()-registeredAddressRescanDepthBlocks, 0)
	if _, err = useCase.wallet.RescanBlockchain(fromHeight); err != nil {
		return useCase.recordRescanError(ctx, pending, fmt.Errorf("rescan PegIn addresses: %w", err))
	}
	return useCase.markImported(ctx, pending)
}

func (useCase *ImportRegisteredAddressesUseCase) markImported(
	ctx context.Context,
	pending []rootstock.PegInWatch,
) ([]rootstock.PegInWatch, error) {
	for index := range pending {
		watch := pending[index]
		watch.MarkImported()
		if err := useCase.repository.Update(ctx, watch); err != nil {
			return pending, fmt.Errorf("persist imported state for event %s/%d: %w", watch.TxHash, watch.LogIndex, err)
		}
		pending[index] = watch
	}
	return pending, nil
}

func (useCase *ImportRegisteredAddressesUseCase) recordRescanError(
	ctx context.Context,
	pending []rootstock.PegInWatch,
	rescanErr error,
) ([]rootstock.PegInWatch, error) {
	for index := range pending {
		watch := pending[index]
		if err := useCase.persistError(ctx, &watch, rescanErr); err != nil {
			return pending, err
		}
		pending[index] = watch
	}
	return pending, nil
}

func (useCase *ImportRegisteredAddressesUseCase) persistError(
	ctx context.Context,
	watch *rootstock.PegInWatch,
	watchErr error,
) error {
	log.Error(watchErr)
	if !watch.RecordError(watchErr) {
		return nil
	}
	if err := useCase.repository.Update(ctx, *watch); err != nil {
		return fmt.Errorf("persist PegIn address registry watch error for %s/%d: %w", watch.TxHash, watch.LogIndex, err)
	}
	return nil
}
