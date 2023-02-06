package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegout"

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

func startServer(rsk *connectors.RSK, btc *connectors.BTC, dbMongo *mongoDB.DB) {
	lpRepository := storage.NewLPRepository(dbMongo, rsk)
	lp, err := providers.NewLocalProvider(cfg.Provider, lpRepository)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	lpPegOut, err := pegout.NewLocalProvider(cfg.PegoutProvier, lpRepository)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	srv = http.New(rsk, btc, dbMongo, cfgData)
	log.Debug("registering local provider (this might take a while)")
	err = srv.AddProvider(lp)
	if err != nil {
		log.Fatalf("error registering local provider: %v", err)
	}

	err = srv.AddPegOutProvider(lpPegOut)

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
	initCfgData()
	initLogger()
	rand.Seed(time.Now().UnixNano())

	log.Info("starting liquidity provider server")
	log.Debugf("loaded config %+v", cfg)

	dbMongo, err := mongoDB.Connect()
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}

	erpKeys := strings.Split(os.Getenv("ERP_KEYS"), ",")

	log.Debug("ERP Keys: ", erpKeys)

	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, cfg.RSK.BridgeAddr, cfg.RSK.RequiredBridgeConfirmations, cfg.IrisActivationHeight, erpKeys)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(os.Getenv("RSKJ_CONNECTION_STRING"), cfg.Provider.ChainId)
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

	startServer(rsk, btc, dbMongo)

	<-done

	srv.Shutdown()
	rsk.Close()
	btc.Close()

	err = dbMongo.Close()
	if err != nil {
		log.Fatal("error closing DB connection: ", err)
	}
}

func initCfgData() {
	cfgData.MaxQuoteValue = cfg.MaxQuoteValue
	cfgData.RSK = cfg.RSK
}
