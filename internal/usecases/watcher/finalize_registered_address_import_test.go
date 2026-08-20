package watcher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFinalizeRegisteredAddressImportUseCase_Run_PersistsRescanErrorAndReturnsNil(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	entry := &rootstock.PegInAddressRegistryWatchEntry{
		TxHash: "0x1", LogIndex: 1, RskAddress: "0xabc",
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	rescanErr := errors.New("rescan PegIn addresses")
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatchEntry) bool {
		return updated.RskAddress == entry.RskAddress &&
			updated.LastError == rescanErr.Error() &&
			updated.State == rootstock.PegInAddressRegistryWatchDiscovered
	})).Return(nil).Once()

	err := watcher.NewFinalizeRegisteredAddressImportUseCase(repository).
		Run(context.Background(), []*rootstock.PegInAddressRegistryWatchEntry{entry}, rescanErr)

	require.NoError(t, err)
}

func TestFinalizeRegisteredAddressImportUseCase_Run_MarksImportedWhenRescanSucceeds(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	entry := &rootstock.PegInAddressRegistryWatchEntry{
		TxHash: "0x1", LogIndex: 1, RskAddress: "0xabc",
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatchEntry) bool {
		return updated.RskAddress == entry.RskAddress &&
			updated.State == rootstock.PegInAddressRegistryWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()

	err := watcher.NewFinalizeRegisteredAddressImportUseCase(repository).
		Run(context.Background(), []*rootstock.PegInAddressRegistryWatchEntry{entry}, nil)

	require.NoError(t, err)
}

func TestFinalizeRegisteredAddressImportUseCase_Run_EmptyPendingIsNoop(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	err := watcher.NewFinalizeRegisteredAddressImportUseCase(repository).
		Run(context.Background(), nil, errors.New("ignored"))
	require.NoError(t, err)
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
