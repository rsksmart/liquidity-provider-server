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

func TestSetRegistryWatchCursorUseCase_Run(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	repository.EXPECT().SetCursor(test.AnyCtx, uint64(99)).Return(nil).Once()

	useCase := watcher.NewSetRegistryWatchCursorUseCase(repository)
	require.NoError(t, useCase.Run(context.Background(), 99))
}

func TestSetRegistryWatchCursorUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	repository.EXPECT().SetCursor(test.AnyCtx, uint64(99)).Return(assert.AnError).Once()

	useCase := watcher.NewSetRegistryWatchCursorUseCase(repository)
	err := useCase.Run(context.Background(), 99)

	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.SetRegistryWatchCursorId))
}
