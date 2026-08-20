package watcher

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type MarkRegisteredAddressImportedUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewMarkRegisteredAddressImportedUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *MarkRegisteredAddressImportedUseCase {
	return &MarkRegisteredAddressImportedUseCase{repository: repository}
}

func (useCase *MarkRegisteredAddressImportedUseCase) Run(
	ctx context.Context,
	watch *rootstock.PegInAddressRegistryWatch,
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
