package watcher_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWatchedPegInClaimsUseCase_Run(t *testing.T) {
	repository := mocks.NewPegInClaimRepositoryMock(t)
	claims := []rootstock.PegInClaim{
		{RskAddress: "0xa", DepositTxID: "deposit-a", State: rootstock.PegInClaimClaimed},
		{RskAddress: "0xb", DepositTxID: "deposit-b", State: rootstock.PegInClaimResolved},
	}
	repository.EXPECT().ListByState(test.AnyCtx, rootstock.PegInClaimClaimed, rootstock.PegInClaimResolved).
		Return(claims, nil).Once()

	useCase := watcher.NewGetWatchedPegInClaimsUseCase(repository)
	got, err := useCase.Run(context.Background(), rootstock.PegInClaimClaimed, rootstock.PegInClaimResolved)

	require.NoError(t, err)
	assert.Equal(t, claims, got)
}

func TestGetWatchedPegInClaimsUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInClaimRepositoryMock(t)
	repository.EXPECT().ListByState(test.AnyCtx, rootstock.PegInClaimClaimed).Return(nil, assert.AnError).Once()

	useCase := watcher.NewGetWatchedPegInClaimsUseCase(repository)
	got, err := useCase.Run(context.Background(), rootstock.PegInClaimClaimed)

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.GetWatchedPegInClaimsId))
	assert.Nil(t, got)
}
