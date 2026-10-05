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

func pegInWatchesFixture() []rootstock.PegInWatch {
	return []rootstock.PegInWatch{
		{RskAddress: "0xa", State: rootstock.PegInWatchImported},
		{RskAddress: "0xb", State: rootstock.PegInWatchDiscovered},
		{RskAddress: "0xc", State: rootstock.PegInWatchUnsupportedEncoding},
		{RskAddress: "0xd", State: rootstock.PegInWatchDiscovered},
	}
}

func TestGetPegInWatchesUseCase_Run_FiltersByState(t *testing.T) {
	watches := pegInWatchesFixture()
	cases := []struct {
		name   string
		states []rootstock.PegInWatchState
		want   []rootstock.PegInWatch
	}{
		{name: "discovered", states: []rootstock.PegInWatchState{rootstock.PegInWatchDiscovered}, want: []rootstock.PegInWatch{watches[1], watches[3]}},
		{name: "imported", states: []rootstock.PegInWatchState{rootstock.PegInWatchImported}, want: []rootstock.PegInWatch{watches[0]}},
		{
			name:   "several states",
			states: []rootstock.PegInWatchState{rootstock.PegInWatchImported, rootstock.PegInWatchUnsupportedEncoding},
			want:   []rootstock.PegInWatch{watches[0], watches[2]},
		},
		{name: "no states", states: nil, want: []rootstock.PegInWatch{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := mocks.NewPegInWatchRepositoryMock(t)
			repository.EXPECT().List(test.AnyCtx).Return(watches, nil).Once()

			got, err := watcher.NewGetPegInWatchesUseCase(repository).Run(context.Background(), tc.states...)

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetPegInWatchesUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	repository.EXPECT().List(test.AnyCtx).Return(nil, assert.AnError).Once()

	_, err := watcher.NewGetPegInWatchesUseCase(repository).Run(context.Background(), rootstock.PegInWatchImported)

	require.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, string(usecases.GetPegInWatchesId))
}
