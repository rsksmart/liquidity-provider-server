// @Version 0.5
// @Title Liquidity Provider Server
// @Server https://flyover-lps.testnet.rsk.co Testnet
// @Server https://flyover-lps.mainnet.rifcomputing.net Mainnet
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token
package main

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/account"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sethvargo/go-envconfig"

	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"

	"github.com/rsksmart/liquidity-provider-server/storage"
	log "github.com/sirupsen/logrus"
)

var (
	cfg     config
	srv     http.Server
	cfgData http.ConfigData
)

func loadConfig() {
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
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

func startServer(rsk *connectors.RSK, btc *connectors.BTC, dbMongo *mongoDB.DB, endChannel chan<- os.Signal) {
	lpRepository := storage.NewLPRepository(dbMongo, rsk, btc)
	accountProvider := account.NewLocalAccountProvider(cfg.Provider.Keydir, cfg.Provider.PwdFile, cfg.Provider.AccountNum)
	lp, err := pegin.NewLocalProvider(*cfg.Provider, lpRepository, accountProvider)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	lpPegOut, err := pegout.NewLocalProvider(cfg.PegoutProvier, lpRepository)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	srv = http.New(rsk, btc, dbMongo, cfgData, lpRepository, *cfg.Provider)
	log.Debug("registering local provider (this might take a while)")
	req := types.ProviderRegisterRequest{
		Name:                    cfg.PeginProviderName,
		Fee:                     cfg.PeginFee,
		QuoteExpiration:         cfg.PeginQuoteExp,
		AcceptedQuoteExpiration: cfg.PeginAcceptedQuoteExp,
		MinTransactionValue:     cfg.PeginMinTransactValue,
		MaxTransactionValue:     cfg.PeginMaxTransactValue,
		ApiBaseUrl:              cfg.BaseURL,
		Status:                  true,
	}
	err = srv.AddProvider(lp, req)
	if err != nil {
		log.Fatalf("error registering local provider: %v", err)
	}
	req2 := types.ProviderRegisterRequest{
		Name:                    cfg.PegoutProviderName,
		Fee:                     cfg.PegoutFee,
		QuoteExpiration:         cfg.PegoutQuoteExp,
		AcceptedQuoteExpiration: cfg.PegoutAcceptedQuoteExp,
		MinTransactionValue:     cfg.PegoutMinTransactValue,
		MaxTransactionValue:     cfg.PegoutMaxTransactValue,
		ApiBaseUrl:              cfg.BaseURL,
		Status:                  true,
	}
	err = srv.AddPegOutProvider(lpPegOut, req2)

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
			endChannel <- syscall.SIGTERM
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

	chainId, err := strconv.ParseInt(os.Getenv("RSK_CHAIN_ID"), 10, 64)

	if err != nil {
		log.Fatal("Error getting the chain ID: ", err)
	}

	err = rsk.Connect(os.Getenv("RSKJ_CONNECTION_STRING"), big.NewInt(chainId))
	if err != nil {
		log.Fatal("error connecting to RSK: ", err)
	}

	btc, err := connectors.NewBTC(os.Getenv("BTC_NETWORK"))
	if err != nil {
		log.Fatal("error initializing BTC connector: ", err)
	}

	err = btc.Connect(os.Getenv("BTC_ENDPOINT"), os.Getenv("BTC_USERNAME"), os.Getenv("BTC_PASSWORD"))
	if err != nil {
		log.Fatal("error connecting to BTC: ", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, btc, dbMongo, done)

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
