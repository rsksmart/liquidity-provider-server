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
	cfg config
	srv http.Server
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
	pwdFile, err := os.Open(cfg.Provider.PwdFilePath)
	providerCfg := providers.ProviderConfig{
		PwdFile:    pwdFile,
		Keydir:     cfg.Provider.Keystore,
		BtcAddr:    cfg.Provider.BtcAddress,
		AccountNum: cfg.Provider.RskAccountNum,
	}
	lp, err := providers.NewLocalProvider(providerCfg)

	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	srv = http.New(rsk, btc, db, cfg.IrisActivationHeight, cfg.ErpKeys, cfg.RSK.LBCAddr)
	srv.AddProvider(lp)
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

	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, cfg.RSK.BridgeAddr)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(cfg.RSK.Endpoint)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	btc := connectors.NewBTC()
	err = btc.Connect(cfg.BTC.Endpoint, cfg.BTC.Username, cfg.BTC.Password, cfg.BTC.Network)
	if err != nil {
		log.Fatal("error connecting to BTC: ", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, btc, db)

	<-done

	srv.Shutdown()
	db.Close()
	rsk.Close()
	btc.Close()
}
