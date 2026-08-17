package lps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	env                  environment.Environment
	timeouts             environment.ApplicationTimeouts
	peginAddressRegistry environment.PegInAddressRegistryWatcherConfig
	lpRegistry           *registry.LiquidityProvider
	useCaseRegistry      *registry.UseCaseRegistry
	watcherRegistry      *registry.WatcherRegistry
	rskRegistry          *registry.Rootstock
	btcRegistry          *registry.Bitcoin
	dbRegistry           *registry.Database
	messagingRegistry    *registry.Messaging
	runningServices      []entities.Closeable
	doneChannel          chan os.Signal
}

type peginAddressRegistryDeploymentVerifier interface {
	IsDeploymentBlock(ctx context.Context, blockNumber uint64) (bool, error)
}

func NewApplication(initCtx context.Context, env environment.Environment, timeouts environment.ApplicationTimeouts) *Application {
	secretLoader, err := secrets.GetSecretLoader(initCtx, env)
	if err != nil {
		log.Fatal("Error getting secret loader:", err)
	}

	rskClient, err := bootstrap.Rootstock(initCtx, env)
	if err != nil {
		log.Fatal("Error connecting to RSK node: ", err)
	}
	log.Debug("Connected to RSK node")
	walletFactory, err := wallet.NewFactory(env, wallet.FactoryCreationArgs{
		Ctx: initCtx, Env: env, SecretLoader: secretLoader, RskClient: rskClient, Timeouts: timeouts,
	})
	if err != nil {
		log.Fatal("Error creating wallet factory: ", err)
	}
	btcConnection, err := btc_bootstrap.Bitcoin(env.Btc)
	if err != nil {
		log.Fatal("Error connecting to the bitcoin node: ", err)
	}
	log.Debug("Connected to BTC node RPC server")
	dbConnection, err := bootstrap.Mongo(initCtx, env.Mongo, timeouts)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	log.Debug("Connected to MongoDB")
	externalClients, err := createExternalRpc(initCtx, env)
	if err != nil {
		log.Fatal(err)
	}

	btcRegistry := NewBitcoinRegistry(walletFactory, btcConnection, os.Exit)
	dbRegistry := registry.NewDatabaseRegistry(dbConnection)
	rootstockRegistry, err := registry.NewRootstockRegistry(env, rskClient, walletFactory, timeouts)
	if err != nil {
		log.Fatal("Error creating Rootstock registry:", err)
	}
	messagingRegistry := registry.NewMessagingRegistry(initCtx, env, rskClient, btcConnection, externalClients)
	lpRegistry, err := registry.NewLiquidityProviderRegistry(dbRegistry, rootstockRegistry, btcRegistry, messagingRegistry, walletFactory)
	if err != nil {
		log.Fatal("Error creating Liquidity Provider registry:", err)
	}
	mutexes := environment.NewApplicationMutexes()

	// Resolved once, because reporting an incomplete configuration is a side effect of reading it.
	peginAddressRegistry, err := env.Pegin.AddressRegistryWatcherConfig()
	if err != nil {
		log.Fatal("Error reading PegIn address registry watcher configuration: ", err)
	}

	useCaseRegistry := registry.NewUseCaseRegistry(env, rootstockRegistry, btcRegistry, dbRegistry, lpRegistry, messagingRegistry, mutexes)
	watcherRegistry := registry.NewWatcherRegistry(env, useCaseRegistry, rootstockRegistry, btcRegistry, lpRegistry, dbRegistry, messagingRegistry, watcher.NewApplicationTickers(), timeouts, peginAddressRegistry)
	return &Application{
		env: env, timeouts: timeouts, peginAddressRegistry: peginAddressRegistry,
		lpRegistry: lpRegistry, useCaseRegistry: useCaseRegistry,
		rskRegistry: rootstockRegistry, btcRegistry: btcRegistry,
		dbRegistry: dbRegistry, messagingRegistry: messagingRegistry,
		watcherRegistry: watcherRegistry, runningServices: make([]entities.Closeable, 0),
	}
}

func createExternalRpc(ctx context.Context, env environment.Environment) (registry.ExternalRpc, error) {
	externalRskSources, err := bootstrap.ExternalRskSources(ctx, env)
	if err != nil {
		return registry.ExternalRpc{}, fmt.Errorf("error connecting to external RSK clients: %w", err)
	} else if len(externalRskSources) == 0 {
		log.Warn("No external RSK clients configured")
	}

	externalBtcSources, err := btc_bootstrap.ExternalBitcoinSources(env)
	if err != nil {
		return registry.ExternalRpc{}, fmt.Errorf("error connecting to external BTC clients: %w", err)
	} else if len(externalBtcSources) == 0 {
		log.Warn("No external BTC sources configured")
	}
	return registry.ExternalRpc{
		RskExternalRpc: externalRskSources,
		BtcExternalRpc: externalBtcSources,
	}, nil
}

func NewBitcoinRegistry(walletFactory wallet.AbstractFactory, btcConnection *bitcoin.Connection, exitFn func(int)) *registry.Bitcoin {
	btcRegistry, err := registry.NewBitcoinRegistry(walletFactory, btcConnection)
	if errors.Is(err, bitcoin.ErrWalletScanning) {
		log.Info("Bitcoin wallet rescan in progress. The server will start once the rescan completes. Exiting cleanly.")
		exitFn(0)
		return nil
	} else if err != nil {
		log.Error("Error creating BTC registry:", err)
		exitFn(1)
		return nil
	}
	return btcRegistry
}

func (app *Application) Run(ctx context.Context, env environment.Environment, logLevel log.Level) {
	app.addRunningService(app.dbRegistry.Connection)
	app.addRunningService(app.rskRegistry.Client)
	app.addRunningService(app.btcRegistry.RpcConnection)
	app.addRunningService(app.btcRegistry.PaymentWallet)
	app.addRunningService(app.btcRegistry.MonitoringWallet)
	app.addRunningService(app.messagingRegistry.EventBus)

	registerParams := blockchain.NewProviderRegistrationParams(app.env.Provider.Name, app.env.Provider.ApiBaseUrl, true, app.env.Provider.ProviderType())
	id, err := app.useCaseRegistry.GetRegistrationUseCase().Run(ctx, registerParams)
	switch {
	case errors.Is(err, usecases.RegistrationRejectedError):
		log.Fatal("Registration rejected by admin while waiting for approval; stopping LPS. Restart to submit a new registration request.")
	case errors.Is(err, usecases.RegistrationWithdrawnError):
		log.Fatal("Registration was withdrawn by the LP owner while waiting for approval; stopping LPS. Restart to submit a new registration request.")
	case err != nil:
		log.Fatal("Error registering provider: ", err)
	default:
		log.Info("Provider registered with ID ", id)
	}

	err = app.useCaseRegistry.GenerateDefaultCredentialsUseCase().Run(ctx, os.TempDir())
	if err != nil {
		log.Fatal("Error generating default password for management interface: ", err)
	}

	err = app.useCaseRegistry.InitializeStateConfigurationUseCase().Run(ctx)
	if err != nil {
		log.Fatal("Error initializing state configuration: ", err)
	}

	if err = app.useCaseRegistry.CheckColdWalletAddressChangeUseCase().Run(ctx); err != nil {
		log.Error("Error checking cold wallet address change: ", err)
	}

	watchers, err := app.prepareWatchers(ctx)
	if err != nil {
		log.Fatal("Error initializing watchers: ", err)
	}
	for _, w := range watchers {
		go w.Start()
	}

	applicationServer, done := server.NewServer(env, app.useCaseRegistry, logLevel, app.timeouts)
	app.doneChannel = done
	app.addRunningService(applicationServer)
	go applicationServer.Start()
	<-done
}

func (app *Application) addRunningService(service entities.Closeable) {
	app.runningServices = append(app.runningServices, service)
}

func (app *Application) prepareWatchers(ctx context.Context) ([]watcher.Watcher, error) {
	var err error
	prepareCtx, cancel := context.WithTimeout(ctx, app.timeouts.WatcherPreparation.Seconds())
	defer cancel()

	if err = app.validatePegInAddressRegistryConfiguration(prepareCtx); err != nil {
		return nil, err
	}
	watchers := app.enabledWatchers()

	for _, w := range watchers {
		if err = w.Prepare(prepareCtx); err != nil {
			return nil, err
		}
		app.addRunningService(w)
	}
	return watchers, nil
}

func (app *Application) validatePegInAddressRegistryConfiguration(ctx context.Context) error {
	if !app.peginAddressRegistry.Enabled {
		return nil
	}

	contract, head, err := app.peginAddressRegistryValidationContext(ctx)
	if err != nil {
		return err
	}
	startBlock := app.peginAddressRegistry.StartBlock
	if err = validatePegInAddressRegistryDeployment(ctx, contract, startBlock); err != nil {
		return err
	}
	if head == startBlock {
		return nil
	}
	if _, err = contract.GetRegistrationRoot(ctx, head); err != nil {
		return fmt.Errorf("validate PegIn address registry at active RSK head %d: %w", head, err)
	}
	return nil
}

func (app *Application) peginAddressRegistryValidationContext(
	ctx context.Context,
) (blockchain.PegInAddressRegistryContract, uint64, error) {
	if app.rskRegistry == nil || app.rskRegistry.Contracts.PegInAddressRegistry == nil {
		return nil, 0, errors.New(
			"PEGIN_ADDRESS_REGISTRY_ADDRESS is required when PegIn address registry watchers are enabled",
		)
	}
	if app.messagingRegistry == nil || app.messagingRegistry.Rpc.Rsk == nil {
		return nil, 0, errors.New("active RSK RPC is required to validate PegIn address registry configuration")
	}

	head, err := app.messagingRegistry.Rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("get active RSK head for PegIn address registry validation: %w", err)
	}
	startBlock := app.peginAddressRegistry.StartBlock
	if startBlock > head {
		return nil, 0, fmt.Errorf(
			"PegIn address registry deployment block %d exceeds active RSK head %d",
			startBlock,
			head,
		)
	}
	return app.rskRegistry.Contracts.PegInAddressRegistry, head, nil
}

func validatePegInAddressRegistryDeployment(
	ctx context.Context,
	contract blockchain.PegInAddressRegistryContract,
	startBlock uint64,
) error {
	deploymentVerifier, ok := contract.(peginAddressRegistryDeploymentVerifier)
	if !ok {
		return errors.New("PegIn address registry adapter does not support exact deployment-block validation")
	}
	isDeploymentBlock, err := deploymentVerifier.IsDeploymentBlock(ctx, startBlock)
	if err != nil {
		return fmt.Errorf("prove PegIn address registry deployment block %d: %w", startBlock, err)
	}
	if !isDeploymentBlock {
		return fmt.Errorf("configured start block %d is not the PegIn address registry deployment block", startBlock)
	}
	deploymentRoot, err := contract.GetRegistrationRoot(ctx, startBlock)
	if err != nil {
		return fmt.Errorf("validate PegIn address registry at deployment block %d: %w", startBlock, err)
	}
	toBlock := startBlock
	deploymentEvents, err := contract.GetAddressRegisteredEvents(ctx, startBlock, &toBlock)
	if err != nil {
		return fmt.Errorf("read PegIn address registry events at deployment block %d: %w", startBlock, err)
	}
	replayedDeploymentRoot, err := replayPegInAddressRegistryDeployment(deploymentEvents, startBlock)
	if err != nil {
		return err
	}
	if deploymentRoot != replayedDeploymentRoot {
		return fmt.Errorf("PegIn address registry deployment block %d already has registry state", startBlock)
	}
	return nil
}

func replayPegInAddressRegistryDeployment(
	deploymentEvents []blockchain.AddressRegistered,
	startBlock uint64,
) ([32]byte, error) {
	sort.Slice(deploymentEvents, func(first, second int) bool {
		return deploymentEvents[first].LogIndex < deploymentEvents[second].LogIndex
	})
	replayedDeploymentRoot := [32]byte{}
	for _, event := range deploymentEvents {
		if event.BlockNumber != startBlock {
			return [32]byte{}, fmt.Errorf(
				"PegIn address registry returned block %d while validating deployment block %d",
				event.BlockNumber,
				startBlock,
			)
		}
		var err error
		replayedDeploymentRoot, err = blockchain.FoldPegInAddressRegistryRoot(replayedDeploymentRoot, event.RskAddress)
		if err != nil {
			return [32]byte{}, fmt.Errorf(
				"validate PegIn address registry event at deployment block %d: %w",
				startBlock,
				err,
			)
		}
		if event.RegistrationRoot != replayedDeploymentRoot {
			return [32]byte{}, fmt.Errorf(
				"PegIn address registry event root differs at deployment block %d",
				startBlock,
			)
		}
	}
	return replayedDeploymentRoot, nil
}

func (app *Application) enabledWatchers() []watcher.Watcher {
	watchers := []watcher.Watcher{
		app.watcherRegistry.PeginDepositAddressWatcher,
		app.watcherRegistry.PeginBridgeWatcher,
		app.watcherRegistry.PegoutRskDepositWatcher,
		app.watcherRegistry.PegoutBtcTransferWatcher,
		app.watcherRegistry.LiquidityCheckWatcher,
		app.watcherRegistry.PenalizationAlertWatcher,
		app.watcherRegistry.PegoutBridgeWatcher,
		app.watcherRegistry.BtcReleaseWatcher,
		app.watcherRegistry.BitcoinPeerWatcher,
		app.watcherRegistry.RootstockPeerWatcher,
		app.watcherRegistry.QuoteMetricsWatcher,
		app.watcherRegistry.PeerMetricsWatcher,
		app.watcherRegistry.AssetReportWatcher,
		app.watcherRegistry.TransferColdWalletWatcher,
		app.watcherRegistry.ColdWalletMetricsWatcher,
		app.watcherRegistry.BitcoinReorgWatcher,
		app.watcherRegistry.RootstockReorgWatcher,
		app.watcherRegistry.ReorgMetricsWatcher,
		app.watcherRegistry.PegInAddressRegistryMetricsWatcher,
	}

	if app.watcherRegistry.PegInClaimWatcher != nil {
		watchers = append(watchers, app.watcherRegistry.PegInClaimWatcher)
	}

	if app.env.Eclipse.Enabled {
		watchers = append(watchers, app.watcherRegistry.RskEclipseWatcher)
		watchers = append(watchers, app.watcherRegistry.BitcoinEclipseWatcher)
	}

	if app.peginAddressRegistryWatchersEnabled() {
		log.Infof(
			"PegIn address registry watcher enabled on RSK chain id %d from block %d with page size %d",
			app.env.Rsk.ChainId,
			app.peginAddressRegistry.StartBlock,
			app.peginAddressRegistry.PageSize,
		)
		watchers = append(watchers, app.watcherRegistry.PegInAddressRegistryWatcher)
	}

	return watchers
}

func (app *Application) peginAddressRegistryWatchersEnabled() bool {
	if !app.peginAddressRegistry.Enabled {
		log.Info("PegIn address registry watchers are disabled")
		return false
	}
	// The registry adapter is only built when its address is configured, so registering the
	// watcher without it would leave the loop calling a nil contract on its first tick.
	if app.rskRegistry.Contracts.PegInAddressRegistry == nil {
		log.Error("PegIn address registry watchers are disabled because PEGIN_ADDRESS_REGISTRY_ADDRESS is missing")
		return false
	}
	return true
}

func (app *Application) ShutdownServices() {
	log.Info("Starting graceful shutdown...")
	numberOfServices := len(app.runningServices)
	closeChannel := make(chan bool, numberOfServices)
	for _, service := range app.runningServices {
		service.Shutdown(closeChannel)
	}
	for i := 0; i < numberOfServices; i++ {
		<-closeChannel
	}
	log.Info("Shutdown completed")
}

func (app *Application) ForceShutdown() {
	app.doneChannel <- syscall.SIGINT
}
