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

func TestGetWatchedRegisteredAddressesUseCase_Run(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	watches := []rootstock.PegInWatch{
		{RskAddress: "0xa", State: rootstock.PegInWatchImported},
		{RskAddress: "0xb", State: rootstock.PegInWatchDiscovered},
	}
	repository.EXPECT().List(test.AnyCtx).Return(watches, nil).Once()

	useCase := watcher.NewGetWatchedRegisteredAddressesUseCase(repository)
	got, err := useCase.Run(context.Background())

	require.NoError(t, err)
	assert.Equal(t, watches, got)
}

func TestGetWatchedRegisteredAddressesUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	repository.EXPECT().List(test.AnyCtx).Return(nil, assert.AnError).Once()

	useCase := watcher.NewGetWatchedRegisteredAddressesUseCase(repository)
	_, err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.GetWatchedRegisteredAddressesId))
}
