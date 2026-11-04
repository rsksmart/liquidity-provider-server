package watcher

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type MarkRegisteredAddressImportedUseCase struct {
	repository rootstock.PegInWatchRepository
}

func NewMarkRegisteredAddressImportedUseCase(
	repository rootstock.PegInWatchRepository,
) *MarkRegisteredAddressImportedUseCase {
	return &MarkRegisteredAddressImportedUseCase{repository: repository}
}

func (useCase *MarkRegisteredAddressImportedUseCase) Run(
	ctx context.Context,
	watch *rootstock.PegInWatch,
) error {
	watch.MarkImported()
	if err := useCase.repository.Update(ctx, *watch); err != nil {
		return usecases.WrapUseCaseError(
			usecases.MarkRegisteredAddressImportedId,
			fmt.Errorf("persist imported state for event %s/%d: %w", watch.TxHash, watch.LogIndex, err),
		)
	}
	return nil
}
