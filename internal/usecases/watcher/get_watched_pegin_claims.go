package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetWatchedPegInClaimsUseCase struct {
	repository rootstock.PegInClaimRepository
}

func NewGetWatchedPegInClaimsUseCase(
	repository rootstock.PegInClaimRepository,
) *GetWatchedPegInClaimsUseCase {
	return &GetWatchedPegInClaimsUseCase{repository: repository}
}

func (useCase *GetWatchedPegInClaimsUseCase) Run(
	ctx context.Context,
	states ...rootstock.PegInClaimState,
) ([]rootstock.PegInClaim, error) {
	claims, err := useCase.repository.ListByState(ctx, states...)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetWatchedPegInClaimsId, err)
	}
	return claims, nil
}
