package lps

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testChainId = 31

// The gate is exercised through enabledWatchers because the watcher set is decided before any
// watcher is prepared, so the decision can be observed without live RSK, BTC or Mongo services.
func peginAddressRegistryGateApplication(
	contract blockchain.PegInAddressRegistryContract,
) (*Application, *watcher.PegInAddressRegistryWatcher) {
	scanner := &watcher.PegInAddressRegistryWatcher{}
	app := &Application{
		env:         environment.Environment{Rsk: environment.RskEnv{ChainId: testChainId}},
		rskRegistry: &registry.Rootstock{Contracts: blockchain.RskContracts{PegInAddressRegistry: contract}},
		watcherRegistry: &registry.WatcherRegistry{
			PegInAddressRegistryWatcher: scanner,
		},
	}
	return app, scanner
}

func TestApplication_enabledWatchers_peginAddressRegistryGate(t *testing.T) {
	enabledConfig := environment.PegInAddressRegistryWatcherConfig{StartBlock: 4500000, PageSize: 100, Enabled: true}

	t.Run("gate set registers the scanner and nothing else", func(t *testing.T) {
		app, scanner := peginAddressRegistryGateApplication(mocks.NewPegInAddressRegistryContractMock(t))
		baseline := app.enabledWatchers()
		require.NotContains(t, baseline, scanner)
		assertLog := test.AssertLogContains(t, "chain id 31 from block 4500000 with page size 100")

		app.peginAddressRegistry = enabledConfig
		gated := app.enabledWatchers()

		assert.Equal(t, append(baseline, scanner), gated)
		assert.True(t, assertLog(), "registration is not reported with the active chain id and start block")
	})
	t.Run("gate unset leaves the watcher set at the baseline", func(t *testing.T) {
		app, scanner := peginAddressRegistryGateApplication(mocks.NewPegInAddressRegistryContractMock(t))

		watchers := app.enabledWatchers()

		assert.NotContains(t, watchers, scanner)
	})
	t.Run("gate set without a registry address registers neither loop", func(t *testing.T) {
		app, scanner := peginAddressRegistryGateApplication(nil)
		app.peginAddressRegistry = enabledConfig
		assertLog := test.AssertLogContains(t, "PEGIN_ADDRESS_REGISTRY_ADDRESS is missing")

		watchers := app.enabledWatchers()

		assert.NotContains(t, watchers, scanner)
		assert.True(t, assertLog(), "the missing registry address is not reported")
	})
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
