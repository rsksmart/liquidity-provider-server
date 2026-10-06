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

func TestGetPendingRegisteredAddressImportsUseCase_Run_FiltersDiscoveredWatches(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	watches := []rootstock.PegInWatch{
		{RskAddress: "0xa", State: rootstock.PegInWatchImported},
		{RskAddress: "0xb", State: rootstock.PegInWatchDiscovered},
		{RskAddress: "0xc", State: rootstock.PegInWatchUnsupportedEncoding},
		{RskAddress: "0xd", State: rootstock.PegInWatchDiscovered},
	}
	repository.EXPECT().List(test.AnyCtx).Return(watches, nil).Once()

	useCase := watcher.NewGetPendingRegisteredAddressImportsUseCase(repository)
	got, err := useCase.Run(context.Background())

	require.NoError(t, err)
	assert.Equal(t, []rootstock.PegInWatch{watches[1], watches[3]}, got)
}

func TestGetPendingRegisteredAddressImportsUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	repository.EXPECT().List(test.AnyCtx).Return(nil, assert.AnError).Once()

	useCase := watcher.NewGetPendingRegisteredAddressImportsUseCase(repository)
	_, err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.GetPendingRegisteredAddressImportsId))
}
