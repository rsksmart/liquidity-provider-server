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
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rsksmart/liquidity-provider-server/account"
	"github.com/rsksmart/liquidity-provider-server/secrets"

	"github.com/sethvargo/go-envconfig"

	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider/types"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/http"

	"github.com/rsksmart/liquidity-provider-server/storage"
	log "github.com/sirupsen/logrus"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
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

	awsConfiguration, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("error loading configuration: ", err.Error())
	}

	peginSecretsStorage := secrets.NewSecretsManagerStorage[any](awsConfiguration)
	peginSecretNames := &account.AccountSecretNames{KeySecretName: cfg.Provider.KeySecret, PasswordSecretName: cfg.Provider.PasswordSecret}
	peginAccountProvider := account.NewRemoteAccountProvider(cfg.Provider.Keydir, cfg.Provider.AccountNum, peginSecretNames, peginSecretsStorage)
	lp, err := pegin.NewLocalProvider(*cfg.Provider, lpRepository, peginAccountProvider)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	pegoutSecretsStorage := secrets.NewSecretsManagerStorage[any](awsConfiguration)
	pegoutSecretNames := &account.AccountSecretNames{KeySecretName: cfg.PegoutProvier.KeySecret, PasswordSecretName: cfg.PegoutProvier.PasswordSecret}
	pegoutAccountProvider := account.NewRemoteAccountProvider(cfg.Provider.Keydir, cfg.Provider.AccountNum, pegoutSecretNames, pegoutSecretsStorage)
	lpPegOut, err := pegout.NewLocalProvider(cfg.PegoutProvier, lpRepository, pegoutAccountProvider)
	if err != nil {
		log.Fatal("cannot create local provider: ", err)
	}

	key, err := pegoutSecretsStorage.GetTextSecret(os.Getenv("ENCRYPT_APP_KEY"))
	if err != nil {
		key = generateRandomKey(32)
		pegoutSecretsStorage.SaveTextSecret(os.Getenv("ENCRYPT_APP_KEY"), key)
	}

	cfgData.EncryptKey = key

	srv = http.New(rsk, btc, dbMongo, cfgData, lpRepository, *cfg.Provider, peginAccountProvider)
	log.Debug("registering local provider (this might take a while)")
	req := types.ProviderRegisterRequest{
		Name:                cfg.PeginProviderName,
		Fee:                 cfg.PeginFee,
		QuoteExpiration:     cfg.PeginQuoteExp,
		MinTransactionValue: cfg.PeginMinTransactValue,
		MaxTransactionValue: cfg.PeginMaxTransactValue,
		ApiBaseUrl:          cfg.BaseURL,
		Status:              true,
	}

	err = srv.AddProvider(lp, req)
	if err != nil {
		log.Fatalf("error registering local provider: %v", err)
	}
	req2 := types.ProviderRegisterRequest{
		Name:                cfg.PegoutProviderName,
		Fee:                 cfg.PegoutFee,
		QuoteExpiration:     cfg.PegoutQuoteExp,
		MinTransactionValue: cfg.PegoutMinTransactValue,
		MaxTransactionValue: cfg.PegoutMaxTransactValue,
		ApiBaseUrl:          cfg.BaseURL,
		Status:              true,
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

	err = btc.Connect(cfg.BTC)
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
	cfgData.RSK = cfg.RSK
	cfgData.QuoteCacheStartBlock = cfg.QuoteCacheStartBlock
}

func generateRandomKey(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,!#@&")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
