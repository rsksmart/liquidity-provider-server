package lps

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"os"
	"syscall"
	"time"
)

const BootstrapTimeout = 3 * time.Minute // In case LP needs to register
const watcherPreparationTimeout = 3 * time.Second

type Application struct {
	env               environment.Environment
	secrets           environment.ApplicationSecrets
	liquidityProvider *dataproviders.LocalLiquidityProvider
	useCaseRegistry   *registry.UseCaseRegistry
	watcherRegistry   *registry.WatcherRegistry
	rskRegistry       *registry.Rootstock
	btcRegistry       *registry.Bitcoin
	dbRegistry        *registry.Database
	messagingRegistry *registry.Messaging
	runningServices   []entities.Closeable
	doneChannel       chan os.Signal
}

func NewApplication(initCtx context.Context, env environment.Environment, secrets environment.ApplicationSecrets) *Application {
	account, err := bootstrap.RootstockAccount(env.Rsk, env.Btc, secrets)
	if err != nil {
		log.Fatal("Error connecting to RSK account:", err)
	}
	log.Debug("Connected to RSK account")

	monitoringConnection, err := bootstrap.BitcoinWallet(env.Btc, bitcoin.PeginWalletId)
	if err != nil {
		log.Fatal("Error creating BTC monitoring connection:", err)
	}
	log.Debug("Connected to BTC node for monitoring")

	paymentConnection, err := bootstrap.BitcoinWallet(env.Btc, bitcoin.DerivativeWalletId)
	if err != nil {
		log.Fatal("Error creating BTC payment connection:", err)
	}
	log.Debug("Connected to BTC node for payments")

	btcRegistry, err := registry.NewBitcoinRegistry(monitoringConnection, paymentConnection, account)
	if err != nil {
		log.Fatal("Error creating BTC registry:", err)
	}

	connection, err := bootstrap.Mongo(initCtx, env.Mongo)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	dbRegistry := registry.NewDatabaseRegistry(connection)
	log.Debug("Connected to MongoDB")

	rskClient, err := bootstrap.Rootstock(initCtx, env.Rsk)
	if err != nil {
		log.Fatal("Error connecting to RSK node:", err)
	}
	rootstockRegistry, err := registry.NewRootstockRegistry(env.Rsk, rskClient, account, monitoringConnection)
	if err != nil {
		log.Fatal("Error creating Rootstock registry:", err)
	}
	log.Debug("Connected to RSK node")

	messagingRegistry := registry.NewMessagingRegistry(initCtx, env, rskClient, monitoringConnection)
	liquidityProvider := registry.NewLiquidityProvider(dbRegistry, rootstockRegistry, btcRegistry, messagingRegistry)
	mutexes := environment.NewApplicationMutexes()

	useCaseRegistry := registry.NewUseCaseRegistry(env, rootstockRegistry, btcRegistry, dbRegistry, liquidityProvider, messagingRegistry, mutexes)
	watcherRegistry := registry.NewWatcherRegistry(env, useCaseRegistry, rootstockRegistry, btcRegistry, liquidityProvider, messagingRegistry)

	return &Application{
		env:               env,
		secrets:           secrets,
		liquidityProvider: liquidityProvider,
		useCaseRegistry:   useCaseRegistry,
		rskRegistry:       rootstockRegistry,
		btcRegistry:       btcRegistry,
		dbRegistry:        dbRegistry,
		messagingRegistry: messagingRegistry,
		watcherRegistry:   watcherRegistry,
		runningServices:   make([]entities.Closeable, 0),
	}
}

func (app *Application) Run(env environment.Environment, logLevel log.Level) {
	app.addRunningService(app.dbRegistry.Connection)
	app.addRunningService(app.rskRegistry.Client)
	app.addRunningService(app.btcRegistry.MonitoringWalletConnection)
	app.addRunningService(app.btcRegistry.PaymentWalletConnection)
	app.addRunningService(app.messagingRegistry.EventBus)

	registerParams := blockchain.NewProviderRegistrationParams(app.env.Provider.Name, app.env.Provider.ApiBaseUrl, true, app.env.Provider.ProviderType)
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

	watchers, err := app.prepareWatchers()
	if err != nil {
		log.Fatal("Error initializing watchers: ", err)
	}
	for _, w := range watchers {
		go w.Start()
	}

	applicationServer, done := server.NewServer(env, app.useCaseRegistry, logLevel)
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
		app.watcherRegistry.QuoteCleanerWatcher,
		app.watcherRegistry.PegoutRskDepositWatcher,
		app.watcherRegistry.PegoutBtcTransferWatcher,
		app.watcherRegistry.LiquidityCheckWatcher,
		app.watcherRegistry.PenalizationAlertWatcher,
	}

	ctx, cancel := context.WithTimeout(context.Background(), watcherPreparationTimeout)
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
