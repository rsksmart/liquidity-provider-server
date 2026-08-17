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
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	entries := []rootstock.PegInAddressRegistryWatchEntry{
		{RskAddress: "0xa", State: rootstock.PegInAddressRegistryWatchImported},
		{RskAddress: "0xb", State: rootstock.PegInAddressRegistryWatchDiscovered},
	}
	repository.EXPECT().List(test.AnyCtx).Return(entries, nil).Once()

	useCase := watcher.NewGetWatchedRegisteredAddressesUseCase(repository)
	got, err := useCase.Run(context.Background())

	require.NoError(t, err)
	assert.Equal(t, entries, got)
}

func TestGetWatchedRegisteredAddressesUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	repository.EXPECT().List(test.AnyCtx).Return(nil, assert.AnError).Once()

	useCase := watcher.NewGetWatchedRegisteredAddressesUseCase(repository)
	_, err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.GetWatchedRegisteredAddressesId))
}
