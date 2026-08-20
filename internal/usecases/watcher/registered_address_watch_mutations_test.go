package watcher_test

import (
	"context"
	"fmt"
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
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	watch := &rootstock.PegInAddressRegistryWatch{
		TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk",
		State:     rootstock.PegInAddressRegistryWatchDiscovered,
		LastError: "previous",
	}
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatch) bool {
		return updated.TxHash == watch.TxHash &&
			updated.State == rootstock.PegInAddressRegistryWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()

	useCase := watcher.NewMarkRegisteredAddressImportedUseCase(repository)
	require.NoError(t, useCase.Run(context.Background(), watch))
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, watch.State)
}

func TestMarkRegisteredAddressImportedUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	watch := &rootstock.PegInAddressRegistryWatch{TxHash: "0xreg", LogIndex: 1}
	repository.On("Update", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := watcher.NewMarkRegisteredAddressImportedUseCase(repository)
	err := useCase.Run(context.Background(), watch)
	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.MarkRegisteredAddressImportedId))
}

func TestRecordRegisteredAddressWatchErrorUseCase_Run_SuppressesIdenticalWrites(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	watch := &rootstock.PegInAddressRegistryWatch{
		TxHash: "stuck", LogIndex: 1, RskAddress: "stuck-rsk",
		State:     rootstock.PegInAddressRegistryWatchDiscovered,
		LastError: fmt.Sprintf("import PegIn address for event stuck/1: %v", assert.AnError),
	}
	useCase := watcher.NewRecordRegisteredAddressWatchErrorUseCase(repository)
	require.NoError(t, useCase.Run(
		context.Background(),
		watch,
		fmt.Errorf("import PegIn address for event stuck/1: %w", assert.AnError),
	))
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestRecordRegisteredAddressWatchErrorUseCase_Run_PersistsNewError(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	watch := &rootstock.PegInAddressRegistryWatch{
		TxHash: "stuck", LogIndex: 1,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatch) bool {
		return updated.TxHash == watch.TxHash && updated.LastError != ""
	})).Return(nil).Once()

	useCase := watcher.NewRecordRegisteredAddressWatchErrorUseCase(repository)
	require.NoError(t, useCase.Run(context.Background(), watch, assert.AnError))
	assert.Equal(t, assert.AnError.Error(), watch.LastError)
}

func TestRecordRegisteredAddressWatchErrorUseCase_Run_WrapsError(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	watch := &rootstock.PegInAddressRegistryWatch{TxHash: "stuck", LogIndex: 1}
	repository.On("Update", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	useCase := watcher.NewRecordRegisteredAddressWatchErrorUseCase(repository)
	err := useCase.Run(context.Background(), watch, assert.AnError)
	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.RecordRegisteredAddressWatchErrorId))
}
