package watcher

import (
	"context"
	"slices"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetPegInWatchesUseCase struct {
	repository rootstock.PegInWatchRepository
}

func NewGetPegInWatchesUseCase(repository rootstock.PegInWatchRepository) *GetPegInWatchesUseCase {
	return &GetPegInWatchesUseCase{repository: repository}
}

func (useCase *GetPegInWatchesUseCase) Run(
	ctx context.Context,
	states ...rootstock.PegInWatchState,
) ([]rootstock.PegInWatch, error) {
	if len(states) == 0 {
		return []rootstock.PegInWatch{}, nil
	}
	watches, err := useCase.repository.List(ctx)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetPegInWatchesId, err)
	}
	result := make([]rootstock.PegInWatch, 0)
	for _, watch := range watches {
		if slices.Contains(states, watch.State) {
			result = append(result, watch)
		}
	}
	return result, nil
}
