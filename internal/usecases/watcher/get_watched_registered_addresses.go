package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetWatchedRegisteredAddressesUseCase struct {
	repository rootstock.PegInWatchRepository
}

func NewGetWatchedRegisteredAddressesUseCase(
	repository rootstock.PegInWatchRepository,
) *GetWatchedRegisteredAddressesUseCase {
	return &GetWatchedRegisteredAddressesUseCase{repository: repository}
}

func (useCase *GetWatchedRegisteredAddressesUseCase) Run(
	ctx context.Context,
) ([]rootstock.PegInWatch, error) {
	watches, err := useCase.repository.List(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetWatchedRegisteredAddressesId, err)
	}
	return watches, nil
}
