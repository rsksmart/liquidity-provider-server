package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider-server/storage"
	"github.com/rsksmart/liquidity-provider/providers"
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

var (
	cfg     config
	srv     http.Server
	cfgData http.ConfigData
)

func loadConfig() {
	err := gonfig.GetConf("config.json", &cfg)

	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
}

func initLogger() {
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if cfg.LogFile == "" {
		return
	}
	f, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		log.Error(fmt.Sprintf("cannot open file %v: ", cfg.LogFile), err)
	} else {
		log.SetOutput(f)
	}
}

func startServer(rsk *connectors.RSK, btc *connectors.BTC, db *storage.DB) {
	lpRepository := storage.NewLPRepository(db, rsk)
	lp, err := providers.NewLocalProvider(cfg.Provider, lpRepository)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	initCfgData()

	srv = http.New(rsk, btc, db, cfgData)
	log.Debug("registering local provider (this might take a while)")
	err = srv.AddProvider(lp)
	if err != nil {
		log.Fatalf("error registering local provider: %v", err)
	}
	port := cfg.Server.Port

	if port == 0 {
		port = 8080
	}
	go func() {
		err := srv.Start(port)

		if err != nil {
			log.Error("server error: ", err.Error())
		}
	}()
}

func main() {
	loadConfig()
	initLogger()
	rand.Seed(time.Now().UnixNano())

	log.Info("starting liquidity provider server")
	log.Debugf("loaded config %+v", cfg)

	db, err := storage.Connect(cfg.DB.Path)
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}

	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, cfg.RSK.BridgeAddr, cfg.RSK.RequiredBridgeConfirmations, cfg.IrisActivationHeight, cfg.ErpKeys)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(cfg.RSK.Endpoint, cfg.Provider.ChainId)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	btc, err := connectors.NewBTC(cfg.BTC.Network)
	if err != nil {
		log.Fatal("error initializing BTC connector: ", err)
	}

	err = btc.Connect(cfg.BTC.Endpoint, cfg.BTC.Username, cfg.BTC.Password)
	if err != nil {
		log.Fatal("error connecting to BTC: ", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, btc, db)

	<-done

	srv.Shutdown()
	rsk.Close()
	btc.Close()

	err = db.Close()
	if err != nil {
		log.Fatal("error closing DB connection: ", err)
	}
}

func initCfgData() {
	cfgData.MaxQuoteValue = cfg.MaxQuoteValue

	cfgData.RSK = cfg.RSK
}
