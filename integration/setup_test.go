package integration_test

import (
	"context"
	"fmt"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/rsksmart/liquidity-provider-server/account"
	"github.com/rsksmart/liquidity-provider-server/config"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider-server/secrets"
	"github.com/rsksmart/liquidity-provider-server/storage"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// TODO this file is very likely to change after LPS refactor

var (
	cfg     config.Config
	srv     http.Server
	cfgData http.ConfigData
)

func loadConfig() error {
	err := config.LoadEnv(&cfg)
	return err
}

func initLogger(hooks ...log.Hook) {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	for _, hook := range hooks {
		log.AddHook(hook)
	}
}

func startServer(rsk *connectors.RSK, btc *connectors.BTC, dbMongo *mongoDB.DB, endChannel chan<- os.Signal, readyChannel chan<- error) {
	lpRepository := storage.NewLPRepository(dbMongo, rsk, btc)

	awsConfiguration, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("error loading configuration: ", err.Error())
	}

	secretsStorage := secrets.NewSecretsManagerStorage[any](awsConfiguration)
	secretNames := &account.AccountSecretNames{
		KeySecretName:      cfg.ProviderCredentials.KeySecret,
		PasswordSecretName: cfg.ProviderCredentials.PasswordSecret,
	}
	accountProvider := account.NewRemoteAccountProvider(
		cfg.ProviderCredentials.Keydir,
		cfg.ProviderCredentials.AccountNum,
		secretNames,
		secretsStorage,
	)
	lp, err := pegin.NewLocalProvider(cfg.Provider, lpRepository, accountProvider, cfg.RSK.ChainId)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}
	lpPegOut, err := pegout.NewLocalProvider(&cfg.PegoutProvier, lpRepository, accountProvider, cfg.RSK.ChainId)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	srv = http.New(rsk, btc, dbMongo, cfgData, lpRepository, cfg.Provider, cfg.PegoutProvier, accountProvider, awsConfiguration)
	log.Debug("registering local provider (this might take a while)")
	req := types.ProviderRegisterRequest{
		Name:         cfg.ProviderName,
		ApiBaseUrl:   cfg.BaseURL,
		Status:       true,
		ProviderType: cfg.ProviderType,
	}

	err = srv.AddProvider(lp, lpPegOut, req)
	if err != nil {
		log.Fatalf("error registering local provider: %v", err)
	}
	port := cfg.Server.Port

	if port == 0 {
		port = 8080
	}
	go func() {
		readyChannel <- nil
		err := srv.Start(port)
		if err != nil {
			log.Error("server error: ", err.Error())
			endChannel <- syscall.SIGTERM
		}
	}()
}

func initCfgData() {
	cfgData.RSK = cfg.RSK
	cfgData.QuoteCacheStartBlock = cfg.QuoteCacheStartBlock
	cfgData.CaptchaSecretKey = cfg.CaptchaSecretKey
	cfgData.CaptchaThreshold = cfg.CaptchaThreshold
	cfgData.CaptchaSiteKey = cfg.CaptchaSiteKey
}

func setup(readyChannel chan<- error, doneChannel chan os.Signal, logHooks ...log.Hook) {
	initLogger(logHooks...)
	err := loadConfig()
	if err != nil {
		readyChannel <- fmt.Errorf("error loading configuration: %v", err)
		return
	}
	initCfgData()
	rand.Seed(time.Now().UnixNano())

	log.Info("starting liquidity provider server")
	log.Debugf("loaded config %+v", cfg)

	dbMongo, err := mongoDB.Connect()
	if err != nil {
		readyChannel <- fmt.Errorf("error connecting to DB: %v", err)
		return
	}

	erpKeys := strings.Split(os.Getenv("ERP_KEYS"), ",")
	log.Debug("ERP Keys: ", erpKeys)
	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, cfg.RSK.BridgeAddr, cfg.RSK.RequiredBridgeConfirmations, cfg.IrisActivationHeight, erpKeys)
	if err != nil {
		readyChannel <- fmt.Errorf("RSK error: %v", err)
		return
	}

	chainId, err := strconv.ParseInt(os.Getenv("RSK_CHAIN_ID"), 10, 64)
	if err != nil {
		readyChannel <- fmt.Errorf("Error getting the chain ID: %v", err)
		return
	}

	err = rsk.Connect(os.Getenv("RSKJ_CONNECTION_STRING"), big.NewInt(chainId))
	if err != nil {
		readyChannel <- fmt.Errorf("error connecting to RSK: %v", err)
		return
	}

	btc, err := connectors.NewBTC(os.Getenv("BTC_NETWORK"))
	if err != nil {
		readyChannel <- fmt.Errorf("error initializing BTC connector: %v", err)
		return
	}

	err = btc.Connect(cfg.BTC)
	if err != nil {
		readyChannel <- fmt.Errorf("error connecting to BTC: %v", err)
		return
	}

	signal.Notify(doneChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startServer(rsk, btc, dbMongo, doneChannel, readyChannel)
	<-doneChannel
	srv.Shutdown()
	rsk.Close()
	btc.Close()
	err = dbMongo.Close()
	if err != nil {
		log.Fatal("error closing DB connection: ", err)
	}
}
