package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider-server/storage"
	providers "github.com/rsksmart/liquidity-provider/providers"
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

type config struct {
	LogFile      string
	FedAddr      string
	FedPubKey    string
	IsTestNet    bool
	Debug        bool
	RedeemScript string

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK struct {
		Endpoint string
		LBCAddr  string
		LBCABI   string
	}
	BTC struct {
		Endpoint string
	}
	Provider struct {
		Keystore    string
		AccountNum  uint
		PwdFilePath string
	}
}

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

func startServer(rsk *connectors.RSK, db *storage.DB) {
	pwdFile, err := os.Open(cfg.Provider.PwdFilePath)
	lp, err := providers.NewLocalProvider(cfg.Provider.Keystore, int(cfg.Provider.AccountNum), pwdFile)

	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}
	srv = http.New(rsk, db, cfg.IsTestNet)
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

	script, err := hex.DecodeString(cfg.RedeemScript)
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	db, err := storage.Connect(cfg.DB.Path)
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}

	abiFile, err := os.Open(cfg.RSK.LBCABI)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, abiFile, script)
	if err != nil {
		log.Fatal("RSK error: ", err)
	}

	err = rsk.Connect(cfg.RSK.Endpoint)
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, db)

	<-done

	srv.Shutdown()
	db.Close()
	rsk.Close()
}
