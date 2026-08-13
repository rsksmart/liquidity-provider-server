package lps

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"time"

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

type peginAddressRegistryDeploymentProof struct {
	blockchain.PegInAddressRegistryContract
	isDeploymentBlock func(context.Context, uint64) (bool, error)
}

func (proof *peginAddressRegistryDeploymentProof) IsDeploymentBlock(
	ctx context.Context,
	blockNumber uint64,
) (bool, error) {
	return proof.isDeploymentBlock(ctx, blockNumber)
}

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

func peginAddressRegistryValidationApplication(
	t *testing.T,
	config environment.PegInAddressRegistryWatcherConfig,
	contract blockchain.PegInAddressRegistryContract,
	rskRpc blockchain.RootstockRpcServer,
) *Application {
	t.Helper()
	if _, ok := contract.(peginAddressRegistryDeploymentVerifier); !ok {
		contract = &peginAddressRegistryDeploymentProof{
			PegInAddressRegistryContract: contract,
			isDeploymentBlock: func(context.Context, uint64) (bool, error) {
				return true, nil
			},
		}
	}
	return &Application{
		env:                  environment.Environment{Rsk: environment.RskEnv{ChainId: testChainId}},
		peginAddressRegistry: config,
		rskRegistry:          &registry.Rootstock{Contracts: blockchain.RskContracts{PegInAddressRegistry: contract}},
		messagingRegistry:    &registry.Messaging{Rpc: blockchain.Rpc{Rsk: rskRpc}},
	}
}

func TestApplication_prepareWatchers_appliesPreparationTimeoutToRegistryValidation(t *testing.T) {
	parentCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	var validationDeadline time.Time
	rskRpc.EXPECT().
		GetHeight(mock.Anything).
		RunAndReturn(func(ctx context.Context) (uint64, error) {
			var found bool
			validationDeadline, found = ctx.Deadline()
			require.True(t, found, "registry validation did not receive a deadline")
			return 0, assert.AnError
		}).
		Once()
	app := &Application{
		timeouts:             environment.ApplicationTimeouts{WatcherPreparation: 1},
		peginAddressRegistry: environment.PegInAddressRegistryWatcherConfig{Enabled: true},
		rskRegistry: &registry.Rootstock{Contracts: blockchain.RskContracts{
			PegInAddressRegistry: mocks.NewPegInAddressRegistryContractMock(t),
		}},
		messagingRegistry: &registry.Messaging{Rpc: blockchain.Rpc{Rsk: rskRpc}},
	}
	startedAt := time.Now()

	_, err := app.prepareWatchers(parentCtx)

	require.ErrorIs(t, err, assert.AnError)
	assert.LessOrEqual(
		t,
		validationDeadline.Sub(startedAt),
		1500*time.Millisecond,
		"registry validation did not use WATCHER_PREPARATION_TIMEOUT",
	)
}

func TestApplication_validatePegInAddressRegistryConfiguration_rejectsUnsafeConfiguration(t *testing.T) {
	config := environment.PegInAddressRegistryWatcherConfig{StartBlock: 100, PageSize: 10, Enabled: true}
	t.Run("deployment block is ahead of the active network", func(t *testing.T) {
		testRegistryDeploymentAheadOfNetwork(t, config)
	})
	t.Run("rejects a start block that is not the contract deployment block", func(t *testing.T) {
		testRegistryStartBlockIsNotDeployment(t, config)
	})
	t.Run("cross-network or incompatible address has no registry at the deployment block", func(t *testing.T) {
		testRegistryIncompatibleAtDeployment(t, config)
	})
	t.Run("rejects a deployment root that does not match its events", func(t *testing.T) {
		testRegistryDeploymentRootMismatch(t, config)
	})
	t.Run("deployment block may contain the first registry event", func(t *testing.T) {
		testRegistryDeploymentContainsFirstEvent(t, config)
	})
	t.Run("address becomes incompatible at the active head", func(t *testing.T) {
		testRegistryIncompatibleAtHead(t, config)
	})
}

func testRegistryDeploymentAheadOfNetwork(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(99), nil).Once()
	app := peginAddressRegistryValidationApplication(
		t,
		config,
		mocks.NewPegInAddressRegistryContractMock(t),
		rskRpc,
	)

	err := app.validatePegInAddressRegistryConfiguration(context.Background())

	require.ErrorContains(t, err, "deployment block 100 exceeds active RSK head 99")
}

func testRegistryStartBlockIsNotDeployment(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(200), nil).Once()
	contract := mocks.NewPegInAddressRegistryContractMock(t)
	proofCalled := false
	proof := &peginAddressRegistryDeploymentProof{
		PegInAddressRegistryContract: contract,
		isDeploymentBlock: func(_ context.Context, blockNumber uint64) (bool, error) {
			proofCalled = true
			assert.Equal(t, uint64(100), blockNumber)
			return false, nil
		},
	}
	app := peginAddressRegistryValidationApplication(t, config, proof, rskRpc)

	err := app.validatePegInAddressRegistryConfiguration(context.Background())

	require.ErrorContains(t, err, "configured start block 100 is not the PegIn address registry deployment block")
	assert.True(t, proofCalled, "deployment-block proof was not requested")
}

func testRegistryIncompatibleAtDeployment(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	incompatibleError := errors.New("no contract code or incompatible ABI")
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(200), nil).Once()
	contract := mocks.NewPegInAddressRegistryContractMock(t)
	contract.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(100)).
		Return([32]byte{}, incompatibleError).
		Once()
	app := peginAddressRegistryValidationApplication(t, config, contract, rskRpc)

	err := app.validatePegInAddressRegistryConfiguration(context.Background())

	require.ErrorContains(t, err, "validate PegIn address registry at deployment block 100")
	require.ErrorIs(t, err, incompatibleError)
}

func testRegistryDeploymentRootMismatch(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
	contract := mocks.NewPegInAddressRegistryContractMock(t)
	contract.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(100)).
		Return([32]byte{1}, nil).
		Once()
	contract.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(100), mock.Anything).
		Return(nil, nil).
		Once()
	app := peginAddressRegistryValidationApplication(t, config, contract, rskRpc)

	err := app.validatePegInAddressRegistryConfiguration(context.Background())

	require.ErrorContains(t, err, "deployment block 100 already has registry state")
}

func testRegistryDeploymentContainsFirstEvent(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	rskAddress := "0x00000000000000000000000000000000000000a1"
	rootBytes, err := hex.DecodeString("5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e")
	require.NoError(t, err)
	deploymentRoot := [32]byte(rootBytes)
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
	contract := mocks.NewPegInAddressRegistryContractMock(t)
	contract.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(100)).
		Return(deploymentRoot, nil).
		Once()
	contract.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(100), mock.Anything).
		Return([]blockchain.AddressRegistered{{
			RskAddress:       rskAddress,
			RegistrationRoot: deploymentRoot,
			BlockNumber:      100,
		}}, nil).
		Once()
	app := peginAddressRegistryValidationApplication(t, config, contract, rskRpc)

	err = app.validatePegInAddressRegistryConfiguration(context.Background())

	require.NoError(t, err)
}

func testRegistryIncompatibleAtHead(t *testing.T, config environment.PegInAddressRegistryWatcherConfig) {
	incompatibleError := errors.New("registry implementation changed")
	rskRpc := mocks.NewRootstockRpcServerMock(t)
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(200), nil).Once()
	contract := mocks.NewPegInAddressRegistryContractMock(t)
	contract.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(100)).
		Return([32]byte{}, nil).
		Once()
	contract.EXPECT().
		GetAddressRegisteredEvents(mock.Anything, uint64(100), mock.Anything).
		Return(nil, nil).
		Once()
	contract.EXPECT().
		GetRegistrationRoot(mock.Anything, uint64(200)).
		Return([32]byte{}, incompatibleError).
		Once()
	app := peginAddressRegistryValidationApplication(t, config, contract, rskRpc)

	err := app.validatePegInAddressRegistryConfiguration(context.Background())

	require.ErrorContains(t, err, "validate PegIn address registry at active RSK head 200")
	require.ErrorIs(t, err, incompatibleError)
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
