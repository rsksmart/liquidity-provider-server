package watcher_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRegistryWatchCursorUseCase_Run(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	repository.EXPECT().GetCursor(test.AnyCtx).Return(uint64(42), true, nil).Once()

	useCase := watcher.NewGetRegistryWatchCursorUseCase(repository)
	block, found, err := useCase.Run(context.Background())

	require.NoError(t, err)
	assert.Equal(t, uint64(42), block)
	assert.True(t, found)
}

func TestGetRegistryWatchCursorUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	repository.EXPECT().GetCursor(test.AnyCtx).Return(uint64(0), false, assert.AnError).Once()

	useCase := watcher.NewGetRegistryWatchCursorUseCase(repository)
	_, _, err := useCase.Run(context.Background())

	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.GetRegistryWatchCursorId))
}
