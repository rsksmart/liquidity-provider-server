package lps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

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
	env               environment.Environment
	timeouts          environment.ApplicationTimeouts
	lpRegistry        *registry.LiquidityProvider
	useCaseRegistry   *registry.UseCaseRegistry
	watcherRegistry   *registry.WatcherRegistry
	rskRegistry       *registry.Rootstock
	btcRegistry       *registry.Bitcoin
	dbRegistry        *registry.Database
	messagingRegistry *registry.Messaging
	runningServices   []entities.Closeable
	doneChannel       chan os.Signal
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

	btcRegistry, err := registry.NewBitcoinRegistry(walletFactory, btcConnection)
	if err != nil {
		log.Fatal("Error creating BTC registry:", err)
	}

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

	useCaseRegistry := registry.NewUseCaseRegistry(env, rootstockRegistry, btcRegistry, dbRegistry, lpRegistry, messagingRegistry, mutexes)
	watcherRegistry := registry.NewWatcherRegistry(env, useCaseRegistry, rootstockRegistry, btcRegistry, lpRegistry, messagingRegistry, watcher.NewApplicationTickers(), timeouts)
	return &Application{
		env: env, timeouts: timeouts,
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

func (app *Application) Run(env environment.Environment, logLevel log.Level) {
	app.addRunningService(app.dbRegistry.Connection)
	app.addRunningService(app.rskRegistry.Client)
	app.addRunningService(app.btcRegistry.RpcConnection)
	app.addRunningService(app.btcRegistry.PaymentWallet)
	app.addRunningService(app.btcRegistry.MonitoringWallet)
	app.addRunningService(app.messagingRegistry.EventBus)

	registerParams := blockchain.NewProviderRegistrationParams(app.env.Provider.Name, app.env.Provider.ApiBaseUrl, true, app.env.Provider.ProviderType())
	id, err := app.useCaseRegistry.GetRegistrationUseCase().Run(registerParams)
	if errors.Is(err, usecases.AlreadyRegisteredError) {
		log.Info("Provider already registered")
	} else if err != nil {
		log.Fatal("Error registering provider: ", err)
	} else {
		log.Info("Provider registered with ID ", id)
	}

	err = app.useCaseRegistry.GenerateDefaultCredentialsUseCase().Run(context.Background(), os.TempDir())
	if err != nil {
		log.Fatal("Error generating default password for management interface: ", err)
	}

	err = app.useCaseRegistry.InitializeStateConfigurationUseCase().Run(context.Background())
	if err != nil {
		log.Fatal("Error initializing state configuration: ", err)
	}

	if err = app.useCaseRegistry.CheckColdWalletAddressChangeUseCase().Run(context.Background()); err != nil {
		log.Error("Error checking cold wallet address change: ", err)
	}

	watchers, err := app.prepareWatchers()
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

func (app *Application) prepareWatchers() ([]watcher.Watcher, error) {
	var err error
	watchers := []watcher.Watcher{
		app.watcherRegistry.PeginDepositAddressWatcher,
		app.watcherRegistry.PeginBridgeWatcher,
		app.watcherRegistry.PegoutRskDepositWatcher,
		app.watcherRegistry.PegoutBtcTransferWatcher,
		app.watcherRegistry.LiquidityCheckWatcher,
		app.watcherRegistry.PenalizationAlertWatcher,
		app.watcherRegistry.PegoutBridgeWatcher,
		app.watcherRegistry.BtcReleaseWatcher,
		app.watcherRegistry.QuoteMetricsWatcher,
		app.watcherRegistry.AssetReportWatcher,
		app.watcherRegistry.TransferColdWalletWatcher,
		app.watcherRegistry.ColdWalletMetricsWatcher,
	}

	if app.env.Eclipse.Enabled {
		watchers = append(watchers, app.watcherRegistry.RskEclipseWatcher)
		watchers = append(watchers, app.watcherRegistry.BitcoinEclipseWatcher)
	}

	ctx, cancel := context.WithTimeout(context.Background(), app.timeouts.WatcherPreparation.Seconds())
	defer cancel()
	for _, w := range watchers {
		if err = w.Prepare(ctx); err != nil {
			return nil, err
		}
		app.addRunningService(w)
	}
	return watchers, nil
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
