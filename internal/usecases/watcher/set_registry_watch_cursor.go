package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type SetRegistryWatchCursorUseCase struct {
	repository rootstock.PegInWatchRepository
}

func NewSetRegistryWatchCursorUseCase(
	repository rootstock.PegInWatchRepository,
) *SetRegistryWatchCursorUseCase {
	return &SetRegistryWatchCursorUseCase{repository: repository}
}

func (useCase *SetRegistryWatchCursorUseCase) Run(ctx context.Context, lastScannedBlock uint64) error {
	if err := useCase.repository.SetCursor(ctx, lastScannedBlock); err != nil {
		return usecases.WrapUseCaseError(usecases.SetRegistryWatchCursorId, err)
	}
	return nil
}
