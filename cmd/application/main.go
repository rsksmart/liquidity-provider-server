package main

import (
	"context"
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

// @Version 1.2.1
// @Title Liquidity Provider Server
// @Server https://lps.testnet.flyover.rif.technology Testnet
// @Server https://lps.flyover.rif.technology Mainnet

var (
	BuildVersion string
	BuildTime    string
)

func main() {
	memguard.CatchInterrupt()
	defer memguard.Purge()

	env := environment.LoadEnv()
	timeouts, err := environment.TimeoutsFromEnv(env.Timeouts)
	if err != nil {
		log.Fatal("Error parsing timeouts: ", err)
	}
	initCtx, cancel := context.WithTimeout(context.Background(), timeouts.Bootstrap.Seconds())

	logLevel := setUpLogger(*env)
	logBuildInfo()
	log.Debugf("Environment loaded: %+v", env)

	log.Info("Initializing application...")
	log.Debugf("Using following timeouts (in seconds): %+v", timeouts)
	app := lps.NewApplication(initCtx, *env, timeouts)
	log.Info("Application initialized successfully")
	cancel()
	log.Info("Starting application...")
	app.Run(*env, logLevel)
	app.ShutdownServices()
}

func setUpLogger(env environment.Environment) log.Level {
	var file *os.File
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
	return logLevel
}

func logBuildInfo() {
	if BuildVersion == "" {
		BuildVersion = "No version provided during build"
	}
	if BuildTime == "" {
		BuildTime = "No time provided during build"
	}
	log.Info("Build version: ", BuildVersion)
	log.Info("Build time: ", BuildTime)
}
