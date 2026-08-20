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

func TestMarkRegisteredAddressImportedUseCase_Run(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	watch := &rootstock.PegInWatch{
		TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk",
		State:     rootstock.PegInWatchDiscovered,
		LastError: "previous",
	}
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInWatch) bool {
		return updated.TxHash == watch.TxHash &&
			updated.State == rootstock.PegInWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()

	useCase := watcher.NewMarkRegisteredAddressImportedUseCase(repository)
	require.NoError(t, useCase.Run(context.Background(), watch))
	assert.Equal(t, rootstock.PegInWatchImported, watch.State)
}

func TestMarkRegisteredAddressImportedUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInWatchRepositoryMock(t)
	watch := &rootstock.PegInWatch{TxHash: "0xreg", LogIndex: 1}
	repository.On("Update", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := watcher.NewMarkRegisteredAddressImportedUseCase(repository)
	err := useCase.Run(context.Background(), watch)
	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.MarkRegisteredAddressImportedId))
}
