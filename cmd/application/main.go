package main

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const bootstrapTimeout = 10 * time.Second

func main() {
	var err error
	var file *os.File
	initCtx, cancel := context.WithTimeout(context.Background(), bootstrapTimeout)

	env := environment.LoadEnv()
	secrets := environment.LoadSecrets(initCtx, *env)

	logLevel, err := log.ParseLevel(env.LogLevel)
	if err != nil {
		log.Fatal("Error parsing log level:", err)
	}
	log.SetLevel(logLevel)

	if env.LogFile != "" {
		if err = os.MkdirAll(path.Dir(env.LogFile), 0700); err != nil {
			log.Fatal(fmt.Sprintf("cannot create log file path (%v): ", env.LogFile), err)
		}
		if file, err = os.OpenFile(env.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
			log.Fatal(fmt.Sprintf("cannot open log file %v: ", env.LogFile), err)
		} else {
			log.SetOutput(file)
		}
	}

	app := lps.NewApplication(initCtx, *env, *secrets)
	cancel()
	app.Run(*env, logLevel)
	app.Shutdown()
}
