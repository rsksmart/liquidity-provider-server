package http

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/config"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	server "github.com/rsksmart/liquidity-provider-server/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type ExampleTestSuite struct {
	suite.Suite
	rsk connectors.RSKConnector
}

type Execution struct {
	Body     string
	Method   string
	URL      string
	Response interface{}
}

var (
	cfg config.Config
)

func loadConfig() {
	if err := config.LoadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("error loading config file: %v", err))
	}
}

func execute(execution *Execution) error {

	req, _ := http.NewRequest(execution.Method, execution.URL, strings.NewReader(execution.Body))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&execution.Response)

	if err != nil {
		return err
	}
	return nil
}

func (suite *ExampleTestSuite) SetupTest() {

	loadConfig()

	rsk, err := connectors.NewRSK(cfg.RSK.LBCAddr, cfg.RSK.BridgeAddr, cfg.RSK.RequiredBridgeConfirmations, cfg.IrisActivationHeight, cfg.ErpKeys)

	if err != nil {
		assert.Fail(suite.T(), "rsk did not start")
	}

	err = rsk.Connect(cfg.RSK.Endpoint, cfg.Provider.ChainId)

	if err != nil {
		assert.Fail(suite.T(), "rsk did not start")
	}

	suite.rsk = rsk
}

func (suite *ExampleTestSuite) TestGetQuotePegOut() {
	quote := getQuote(suite)
	assert.Equal(suite.T(), strconv.FormatUint(quote.Quote.Value, 10), "600000000000000000")
	assert.Equal(suite.T(), quote.Quote.RSKRefundAddr, "0xa554d96413FF72E93437C4072438302C38350EE3")
}

func getQuote(suite *ExampleTestSuite) server.QuotePegOutResponse {
	jsonStr := `{
		"from": "1NwGDBiQzGFcyH9aQqeia9XEmaftsgBS4k",
		"valueToTransfer": 600000000000000000,
		"rskRefundAddress": "0xa554d96413FF72E93437C4072438302C38350EE3",
		"bitcoinRefundAddress": "1NwGDBiQzGFcyH9aQqeia9XEmaftsgBS4k"
	}`

	quotes := []server.QuotePegOutResponse{}

	err := execute(&Execution{
		Body:     jsonStr,
		URL:      "http://localhost:8080/pegout/getQuotes",
		Method:   "POST",
		Response: &quotes,
	})

	if err != nil {
		fmt.Println(err)
		assert.Fail(suite.T(), "response error")
	}

	quote := quotes[0]
	return quote
}

var runIntegration = flag.Bool("integration", false, "Run the integration testsuite (in addition to the unit tests)")

func TestPegoutSuite(t *testing.T) {
	if !*runIntegration {
		t.Skip("skipping integration tests")
	}
	suite.Run(t, new(ExampleTestSuite))
}
