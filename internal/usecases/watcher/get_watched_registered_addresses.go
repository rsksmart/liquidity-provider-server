package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetWatchedRegisteredAddressesUseCase struct {
	repository rootstock.PegInAddressRegistryWatchRepository
}

func NewGetWatchedRegisteredAddressesUseCase(
	repository rootstock.PegInAddressRegistryWatchRepository,
) *GetWatchedRegisteredAddressesUseCase {
	return &GetWatchedRegisteredAddressesUseCase{repository: repository}
}

func (useCase *GetWatchedRegisteredAddressesUseCase) Run(
	ctx context.Context,
) ([]rootstock.PegInAddressRegistryWatchEntry, error) {
	entries, err := useCase.repository.List(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetWatchedRegisteredAddressesId, err)
	}
	return entries, nil
}
