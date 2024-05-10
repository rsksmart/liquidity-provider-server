package integration_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"os"
)

const testLogLevel = log.DebugLevel

func setUpLps(referenceChannel chan<- *lps.Application, hooks ...log.Hook) {
	initCtx, cancel := context.WithTimeout(context.Background(), lps.BootstrapTimeout)

	env := environment.LoadEnv()

	log.SetLevel(testLogLevel)
	log.SetOutput(os.Stdout)
	for _, hook := range hooks {
		log.AddHook(hook)
	}
	log.Debugf("Environment loaded: %+v", env)

	log.Info("Initializing application...")
	app := lps.NewApplication(initCtx, *env)
	log.Info("Application initialized successfully")
	cancel()
	log.Info("Starting application...")
	referenceChannel <- app
	app.Run(*env, testLogLevel)
	app.ShutdownServices()
}
