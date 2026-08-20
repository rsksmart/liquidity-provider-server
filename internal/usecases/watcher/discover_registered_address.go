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

func (useCase *DiscoverRegisteredAddressUseCase) Run(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (watch *rootstock.PegInWatch, needsRescan bool, err error) {
	watch, err = useCase.loadOrCreate(ctx, event)
	if err != nil {
		return nil, false, err
	}
	if watch.State == rootstock.PegInWatchImported ||
		watch.State == rootstock.PegInWatchUnsupportedEncoding {
		return watch, false, nil
	}

	pegInAddress, err := useCase.registry.GetPegInAddress(watch.RskAddress)
	if err != nil {
		return watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	watch.Encoding = uint8(pegInAddress.Encoding)
	if pegInAddress.Encoding != blockchain.PegInAddressRegistryEncodingBase58 {
		watch.MarkUnsupportedEncoding(uint8(pegInAddress.Encoding))
		if err = useCase.repository.Update(ctx, *watch); err != nil {
			return nil, false, usecases.WrapUseCaseError(
				usecases.DiscoverRegisteredAddressId,
				fmt.Errorf("persist unsupported encoding for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
			)
		}
		log.Errorf(
			"PegIn address registry event %s/%d uses unsupported encoding %d",
			watch.TxHash,
			watch.LogIndex,
			pegInAddress.Encoding,
		)
		return watch, false, nil
	}
	if watch.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	if err = useCase.wallet.ImportAddress(watch.BtcAddress); err != nil && !btcclient.IsAddressAlreadyImported(err) {
		return watch, false, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("import PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return watch, true, nil
}

func (useCase *DiscoverRegisteredAddressUseCase) loadOrCreate(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (*rootstock.PegInWatch, error) {
	watch, err := useCase.repository.Get(ctx, event.TxHash, event.LogIndex)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err),
		)
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
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("persist AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err),
		)
	}
	watch, err = useCase.repository.Get(ctx, event.TxHash, event.LogIndex)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s/%d: %w", event.TxHash, event.LogIndex, err),
		)
	}
	if watch == nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s/%d: watch not found after upsert", event.TxHash, event.LogIndex),
		)
	}
	return watch, nil
}

func (useCase *DiscoverRegisteredAddressUseCase) persistError(
	ctx context.Context,
	watch *rootstock.PegInWatch,
	watchErr error,
) error {
	log.Error(watchErr)
	if !watch.RecordError(watchErr) {
		return nil
	}
	if err := useCase.repository.Update(ctx, *watch); err != nil {
		return usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("persist PegIn address registry watch error for %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return nil
}
