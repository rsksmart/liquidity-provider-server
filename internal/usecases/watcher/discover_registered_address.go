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

type DiscoverRegisteredAddressResult struct {
	Watch       *rootstock.PegInWatch
	NeedsRescan bool
}

func (useCase *DiscoverRegisteredAddressUseCase) Run(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (DiscoverRegisteredAddressResult, error) {
	watch, err := useCase.loadOrCreate(ctx, event)
	if err != nil {
		return DiscoverRegisteredAddressResult{}, err
	}
	if watch.State == rootstock.PegInWatchImported ||
		watch.State == rootstock.PegInWatchUnsupportedEncoding {
		return DiscoverRegisteredAddressResult{Watch: watch}, nil
	}

	pegInAddress, err := useCase.registry.GetPegInAddress(watch.RskAddress)
	if err != nil {
		return DiscoverRegisteredAddressResult{Watch: watch}, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("resolve PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	watch.SetEncoding(pegInAddress.Encoding)
	if watch.State == rootstock.PegInWatchUnsupportedEncoding {
		if err = useCase.repository.Update(ctx, *watch); err != nil {
			return DiscoverRegisteredAddressResult{}, usecases.WrapUseCaseError(
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
		return DiscoverRegisteredAddressResult{Watch: watch}, nil
	}
	if watch.BtcAddress, err = bitcoin.EncodeAddressBase58(pegInAddress.Payload); err != nil {
		return DiscoverRegisteredAddressResult{Watch: watch}, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("encode PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	if err = useCase.wallet.ImportAddress(watch.BtcAddress); err != nil && !btcclient.IsAddressAlreadyImported(err) {
		return DiscoverRegisteredAddressResult{Watch: watch}, useCase.persistError(
			ctx,
			watch,
			fmt.Errorf("import PegIn address for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return DiscoverRegisteredAddressResult{Watch: watch, NeedsRescan: true}, nil
}

func (useCase *DiscoverRegisteredAddressUseCase) loadOrCreate(
	ctx context.Context,
	event blockchain.AddressRegistered,
) (*rootstock.PegInWatch, error) {
	watch, err := useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
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
	watch, err = useCase.repository.Get(ctx, event.RskAddress)
	if err != nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: %w", event.RskAddress, err),
		)
	}
	if watch == nil {
		return nil, usecases.WrapUseCaseError(
			usecases.DiscoverRegisteredAddressId,
			fmt.Errorf("load AddressRegistered event %s: watch not found after upsert", event.RskAddress),
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
