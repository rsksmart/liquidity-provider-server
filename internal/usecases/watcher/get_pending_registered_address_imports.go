package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetPendingRegisteredAddressImportsUseCase struct {
	repository rootstock.PegInWatchRepository
}

func NewGetPendingRegisteredAddressImportsUseCase(
	repository rootstock.PegInWatchRepository,
) *GetPendingRegisteredAddressImportsUseCase {
	return &GetPendingRegisteredAddressImportsUseCase{repository: repository}
}

func (useCase *GetPendingRegisteredAddressImportsUseCase) Run(
	ctx context.Context,
) ([]rootstock.PegInWatch, error) {
	watches, err := useCase.repository.List(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetPendingRegisteredAddressImportsId, err)
	}
	pendingImports := make([]rootstock.PegInWatch, 0)
	for _, watch := range watches {
		if watch.State == rootstock.PegInWatchDiscovered {
			pendingImports = append(pendingImports, watch)
		}
	}
	return pendingImports, nil
}
