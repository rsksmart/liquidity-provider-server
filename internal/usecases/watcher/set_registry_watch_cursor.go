package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetRegistryWatchCursorUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewSetRegistryWatchCursorUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *SetRegistryWatchCursorUseCase {
	return &SetRegistryWatchCursorUseCase{repository: repository}
}

func (useCase *SetRegistryWatchCursorUseCase) Run(ctx context.Context, lastScannedBlock uint64) error {
	if err := useCase.repository.SetCursor(ctx, lastScannedBlock); err != nil {
		return usecases.WrapUseCaseError(usecases.SetRegistryWatchCursorId, err)
	}
	return nil
}
