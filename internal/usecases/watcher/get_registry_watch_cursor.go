package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetRegistryWatchCursorUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewGetRegistryWatchCursorUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *GetRegistryWatchCursorUseCase {
	return &GetRegistryWatchCursorUseCase{repository: repository}
}

func (useCase *GetRegistryWatchCursorUseCase) Run(
	ctx context.Context,
) (lastScannedBlock uint64, found bool, err error) {
	lastScannedBlock, found, err = useCase.repository.GetCursor(ctx)
	if err != nil {
		return 0, false, usecases.WrapUseCaseError(usecases.GetRegistryWatchCursorId, err)
	}
	return lastScannedBlock, found, nil
}
