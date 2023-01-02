package http

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/connectors"
	server "github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tkanos/gonfig"
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

type config struct {
	LogFile              string
	Debug                bool
	IrisActivationHeight int
	ErpKeys              []string

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK struct {
		Endpoint                    string
		LBCAddr                     string
		BridgeAddr                  string
		RequiredBridgeConfirmations int64
	}
	BTC struct {
		Endpoint string
		Username string
		Password string
		Network  string
	}
	Provider providers.ProviderConfig
}

var (
	cfg config
)

func loadConfig() {
	err := gonfig.GetConf("config.json", &cfg)

	if err != nil {
		panic("config file is missing")
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
	assert.NotNil(suite.T(), quote.DerivationAddress)
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

func (suite *ExampleTestSuite) TestAcceptQuotePegOut() {
	suite.SetupTest()
	quote := getQuote(suite)

	h, err := suite.rsk.HashPegOutQuote(quote.Quote)

	if err != nil {
		fmt.Println(err)
		assert.Fail(suite.T(), "response error")
	}

	acceptQuote := acceptQuote(suite, h, quote.DerivationAddress)

	assert.NotEmpty(suite.T(), acceptQuote.Signature)
}

func acceptQuote(suite *ExampleTestSuite, hash string, derivationAddress string) server.AcceptResPegOut {
	jsonStr := fmt.Sprintf(`{
		"derivationAddress": "%v",
		"quoteHash": "%v"
	}`, derivationAddress, hash)

	acceptQuote := server.AcceptResPegOut{}
	err := execute(&Execution{
		Body:     jsonStr,
		URL:      "http://localhost:8080/pegout/acceptQuote",
		Method:   "POST",
		Response: &acceptQuote,
	})

	if err != nil {
		fmt.Println(err)
		assert.Fail(suite.T(), "response error")
	}
	return acceptQuote
}

var runIntegration = flag.Bool("integration", false, "Run the integration testsuite (in addition to the unit tests)")

func TestPegoutSuite(t *testing.T) {
	if !*runIntegration {
		t.Skip("skipping integration tests")
	}
	suite.Run(t, new(ExampleTestSuite))
}
