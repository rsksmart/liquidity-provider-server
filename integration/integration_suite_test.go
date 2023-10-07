package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

// TODO this file is very likely to change after LPS refactor

type IntegrationTestSuite struct {
	suite.Suite
	SetUpCompletedChannel chan error
	ServerDoneChannel     chan os.Signal
	btc                   *rpcclient.Client
	rsk                   *ethclient.Client
	lbc                   *bindings.LiquidityBridgeContract
	btcParams             chaincfg.Params
	config                SuiteConfig
}

type SuiteConfig struct {
	Lps struct {
		Url             string `json:"url"`
		UseTestInstance bool   `json:"useTestInstance"`
	} `json:"lps"`
	Btc struct {
		RpcEndpoint string `json:"rpcEndpoint"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Network     string `json:"network"`
	} `json:"btc"`
	Rsk struct {
		RpcUrl         string `json:"rpcUrl"`
		LbcAddress     string `json:"lbcAddress"`
		UserPrivateKey string `json:"userPrivateKey"`
	} `json:"rsk"`
}

func (s *IntegrationTestSuite) SetupSuite() {
	log.Debug("Setting up integration tests...")
	configFile, err := os.Open("integration-test.config.json")
	defer configFile.Close()
	if err != nil {
		s.FailNow("Error reading configuration file", err)
	}
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		s.FailNow("Error reading configuration", err)
	}
	err = json.Unmarshal(configBytes, &s.config)
	if err != nil {
		s.FailNow("Error reading configuration", err)
	}

	if s.config.Lps.UseTestInstance {
		if err = s.setupLps(); err != nil {
			s.FailNow("Error setting up LPS", err)
		}
		time.Sleep(1 * time.Second)
	}

	if err = s.setupBtc(); err != nil {
		s.FailNow("Error setting up Bitcoin client", err)
	}

	if err = s.setupRsk(); err != nil {
		s.FailNow("Error setting up RSK client", err)
	}

	log.Debug("Set up completed")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if s.config.Lps.UseTestInstance {
		time.Sleep(3 * time.Second) // to allow LPS to finish updating the database after blockchain calls
		s.ServerDoneChannel <- syscall.SIGINT
	}
}

func (s *IntegrationTestSuite) setupLps() error {
	s.SetUpCompletedChannel = make(chan error, 1)
	s.ServerDoneChannel = make(chan os.Signal, 1)
	fatalHook := &FatalHook{suite: s}
	go setup(s.SetUpCompletedChannel, s.ServerDoneChannel, fatalHook)
	err := <-s.SetUpCompletedChannel
	return err
}

func (s *IntegrationTestSuite) setupBtc() error {
	switch s.config.Btc.Network {
	case "mainnet":
		s.btcParams = chaincfg.MainNetParams
	case "testnet":
		s.btcParams = chaincfg.TestNet3Params
	case "regtest":
		s.btcParams = chaincfg.RegressionNetParams
	default:
		return fmt.Errorf("invalid network name: %v", s.config.Btc.Network)
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
	s.rsk = rsk
	s.lbc = lbc
	return nil
}

func (s *IntegrationTestSuite) AssertFields(expectedFields []string, object map[string]any) {
	for _, field := range expectedFields {
		_, exists := object[field]
		assert.True(s.T(), exists, fmt.Sprintf("Field %v is missing", field))
	}
}

type FatalHook struct {
	suite *IntegrationTestSuite
}

func (h *FatalHook) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
}

func (h *FatalHook) Fire(e *log.Entry) error {
	h.suite.SetUpCompletedChannel <- errors.New(e.Message)
	h.suite.Fail("Unexpected server error", e.Message)
	return nil
}

type Execution struct {
	Body   any
	Method string
	URL    string
}

type Result[responseType any] struct {
	Response    responseType
	RawResponse []byte
	StatusCode  int
}

func execute[responseType any](execution Execution) (Result[responseType], error) {
	payload, err := json.Marshal(execution.Body)
	req, err := http.NewRequest(execution.Method, execution.URL, bytes.NewBuffer(payload))
	if err != nil {
		return Result[responseType]{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return Result[responseType]{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Result[responseType]{}, err
	}

	var response responseType
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return Result[responseType]{}, err
	}

	result := Result[responseType]{
		Response:    response,
		StatusCode:  res.StatusCode,
		RawResponse: bodyBytes,
	}
	return result, nil
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
