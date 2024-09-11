package integration_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	"os"
	"testing"
	"time"
)

type IntegrationTestSuite struct {
	suite.Suite
	btc       *rpcclient.Client
	rsk       *ethclient.Client
	lbc       *bindings.LiquidityBridgeContract
	btcParams chaincfg.Params
	config    SuiteConfig
	app       *lps.Application
}

type SuiteConfig struct {
	Network string `json:"network"`
	Lps     struct {
		Url             string `json:"url"`
		UseTestInstance bool   `json:"useTestInstance"`
	} `json:"lps"`
	Btc struct {
		RpcEndpoint    string `json:"rpcEndpoint"`
		User           string `json:"user"`
		Password       string `json:"password"`
		WalletPassword string `json:"walletPassword"`
	} `json:"btc"`
	Rsk struct {
		RpcUrl         string `json:"rpcUrl"`
		LbcAddress     string `json:"lbcAddress"`
		UserPrivateKey string `json:"userPrivateKey"`
	} `json:"rsk"`
}

func (s *IntegrationTestSuite) SetupSuite() {
	var err error
	var configBytes []byte
	var configFile *os.File

	log.Debug("Setting up integration tests...")
	if configFile, err = os.Open("./integration-test.config.json"); err != nil {
		s.FailNow("Error reading configuration file", err)
	}
	defer func(configFile *os.File) {
		if closingErr := configFile.Close(); closingErr != nil {
			s.FailNow("Error closing configuration file", err)
		}
	}(configFile)

	if configBytes, err = io.ReadAll(configFile); err != nil {
		s.FailNow("Error reading configuration", err)
	}

	if err = json.Unmarshal(configBytes, &s.config); err != nil {
		s.FailNow("Error reading configuration", err)
	}

	if err = s.setupBtc(); err != nil {
		s.FailNow("Error setting up Bitcoin client", err)
	}

	if err = s.setupRsk(); err != nil {
		s.FailNow("Error setting up RSK client", err)
	}

	if s.config.Lps.UseTestInstance {
		s.setupLps()
		time.Sleep(3 * time.Second)
	}

	log.Debug("Set up completed")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.config.Lps.UseTestInstance {
		time.Sleep(3 * time.Second) // to allow LPS to finish updating the database after blockchain calls
		s.app.ForceShutdown()
	}
}

func (s *IntegrationTestSuite) setupLps() {
	fatalHook := &FatalHook{suite: s}
	referenceChannel := make(chan *lps.Application)
	go setUpLps(referenceChannel, fatalHook)
	s.app = <-referenceChannel
}

func (s *IntegrationTestSuite) setupBtc() error {
	switch s.config.Network {
	case "mainnet":
		s.btcParams = chaincfg.MainNetParams
	case "testnet":
		s.btcParams = chaincfg.TestNet3Params
	case "regtest":
		s.btcParams = chaincfg.RegressionNetParams
	default:
		return fmt.Errorf("invalid network name: %v", s.config.Network)
	}

	config := rpcclient.ConnConfig{
		Host:         s.config.Btc.RpcEndpoint,
		User:         s.config.Btc.User,
		Pass:         s.config.Btc.Password,
		Params:       s.btcParams.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	btc, err := rpcclient.New(&config, nil)
	if err != nil {
		return err
	}
	s.btc = btc
	return nil
}

func (s *IntegrationTestSuite) setupRsk() error {
	rsk, err := ethclient.Dial(s.config.Rsk.RpcUrl)
	if err != nil {
		return err
	}
	if !common.IsHexAddress(s.config.Rsk.LbcAddress) {
		return errors.New("invalid LBC address")
	}
	lbcAddress := common.HexToAddress(s.config.Rsk.LbcAddress)
	lbc, err := bindings.NewLiquidityBridgeContract(lbcAddress, rsk)
	if err != nil {
		return err
	}
	s.rsk = rsk
	s.lbc = lbc
	return nil
}

func (s *IntegrationTestSuite) AssertFields(expectedFields []string, object map[string]any) {
	for _, field := range expectedFields {
		_, exists := object[field]
		s.Require().True(exists, fmt.Sprintf("Field %v is missing", field))
	}
}

type FatalHook struct {
	suite *IntegrationTestSuite
}

func (h *FatalHook) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
}

func (h *FatalHook) Fire(e *log.Entry) error {
	h.suite.app.ShutdownServices()
	h.suite.Fail("Unexpected server error", e.Message)
	return nil
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
