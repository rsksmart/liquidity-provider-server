package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetPegInClaimsUseCase struct {
	repository rootstock.PegInClaimRepository
}

func NewGetPegInClaimsUseCase(repository rootstock.PegInClaimRepository) *GetPegInClaimsUseCase {
	return &GetPegInClaimsUseCase{repository: repository}
}

func (useCase *GetPegInClaimsUseCase) Run(
	ctx context.Context,
	states ...rootstock.PegInClaimState,
) ([]rootstock.PegInClaim, error) {
	if len(states) == 0 {
		return []rootstock.PegInClaim{}, nil
	}
	claims, err := useCase.repository.ListByStates(ctx, states...)
	if err != nil {
		return nil, usecases.WrapUseCaseError(usecases.GetPegInClaimsId, err)
	}
	return claims, nil
}
