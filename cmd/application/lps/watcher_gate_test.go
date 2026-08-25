package lps

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func peginAddressRegistryApplication(scanner *watcher.PegInWatcher) *Application {
	return &Application{
		env: environment.Environment{},
		watcherRegistry: &registry.WatcherRegistry{
			PegInAddressRegistryWatcher: scanner,
			PeginDepositAddressWatcher:  &watcher.PeginDepositAddressWatcher{},
		},
	}
}

func TestApplication_enabledWatchers_registersPegInWatcherWhenEnabled(t *testing.T) {
	scanner := &watcher.PegInWatcher{}
	app := peginAddressRegistryApplication(scanner)

	watchers := app.enabledWatchers()

	require.Contains(t, watchers, scanner)
	assert.Equal(t, app.watcherRegistry.PeginDepositAddressWatcher, watchers[0])
}

func TestApplication_enabledWatchers_skipsPegInWatcherWhenDisabled(t *testing.T) {
	app := peginAddressRegistryApplication(nil)

	for _, registered := range app.enabledWatchers() {
		_, isPegInWatcher := registered.(*watcher.PegInWatcher)
		require.False(t, isPegInWatcher)
	}
}

func TestRequirePegInAddressRegistry(t *testing.T) {
	require.NoError(t, requirePegInAddressRegistry(mocks.NewPegInAddressRegistryContractMock(t)))
	require.ErrorContains(t, requirePegInAddressRegistry(nil), "PEGIN_ADDRESS_REGISTRY_ADDRESS")
}

func TestApplication_prepareWatchers_propagatesFirstWatcherPrepareError(t *testing.T) {
	peginRepository := &mocks.PeginQuoteRepositoryMock{}
	peginRepository.EXPECT().GetRetainedQuoteByState(mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	getWatchedUseCase := w.NewGetWatchedPeginQuoteUseCase(peginRepository)
	useCases := watcher.NewPeginDepositAddressWatcherUseCases(nil, getWatchedUseCase, nil, nil)
	addressWatcher := watcher.NewPeginDepositAddressWatcher(useCases, nil, blockchain.Rpc{}, nil, nil)

	app := &Application{
		timeouts: environment.ApplicationTimeouts{WatcherPreparation: 15},
		watcherRegistry: &registry.WatcherRegistry{
			PeginDepositAddressWatcher: addressWatcher,
		},
	}

	watchers, err := app.prepareWatchers(context.Background())

	require.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, watchers)
	assert.Empty(t, app.runningServices)
	peginRepository.AssertExpectations(t)
}
