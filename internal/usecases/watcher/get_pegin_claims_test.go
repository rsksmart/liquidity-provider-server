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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPegInClaimsUseCase_Run_ListsRequestedStates(t *testing.T) {
	repository := mocks.NewPegInClaimRepositoryMock(t)
	claims := []rootstock.PegInClaim{
		{RskAddress: "0xa", DepositTxID: "aa", State: rootstock.PegInClaimSubmitting},
		{RskAddress: "0xb", DepositTxID: "bb", State: rootstock.PegInClaimCandidate},
	}
	repository.EXPECT().ListByStates(test.AnyCtx, rootstock.PegInClaimSubmitting, rootstock.PegInClaimCandidate).
		Return(claims, nil).Once()

	got, err := watcher.NewGetPegInClaimsUseCase(repository).
		Run(context.Background(), rootstock.PegInClaimSubmitting, rootstock.PegInClaimCandidate)

	require.NoError(t, err)
	assert.Equal(t, claims, got)
}

func TestGetPegInClaimsUseCase_Run_NoStatesReturnsEmpty(t *testing.T) {
	repository := mocks.NewPegInClaimRepositoryMock(t)

	got, err := watcher.NewGetPegInClaimsUseCase(repository).Run(context.Background())

	require.NoError(t, err)
	assert.Empty(t, got)
	repository.AssertNotCalled(t, "ListByStates", mock.Anything)
}

func TestGetPegInClaimsUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInClaimRepositoryMock(t)
	repository.EXPECT().ListByStates(test.AnyCtx, rootstock.PegInClaimSubmitting).Return(nil, assert.AnError).Once()

	_, err := watcher.NewGetPegInClaimsUseCase(repository).Run(context.Background(), rootstock.PegInClaimSubmitting)

	require.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, string(usecases.GetPegInClaimsId))
}
