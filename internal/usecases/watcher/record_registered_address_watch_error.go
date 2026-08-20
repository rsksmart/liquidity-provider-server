package watcher

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type RecordRegisteredAddressWatchErrorUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewRecordRegisteredAddressWatchErrorUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *RecordRegisteredAddressWatchErrorUseCase {
	return &RecordRegisteredAddressWatchErrorUseCase{repository: repository}
}

func (useCase *RecordRegisteredAddressWatchErrorUseCase) Run(
	ctx context.Context,
	watch *rootstock.PegInAddressRegistryWatch,
	watchErr error,
) error {
	log.Error(watchErr)
	if !watch.RecordError(watchErr) {
		return nil
	}
	if err := useCase.repository.Update(ctx, *watch); err != nil {
		return usecases.WrapUseCaseError(
			usecases.RecordRegisteredAddressWatchErrorId,
			fmt.Errorf("persist PegIn address registry watch error for %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return nil
}
