package pegout_test

import (
	"context"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	"github.com/rsksmart/liquidity-provider-server/test/integration/contracts"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type PegOutSuite struct {
	suite.Suite
	QuoteResponse       pkg.GetPegoutQuoteResponse
	AcceptQuoteResponse pkg.AcceptPegoutResponse

	// --- initialized during setup ---
	serverUrl      string
	args           TestArguments
	userAccountKey string
	maxPegoutWait  time.Duration
	maxRefundWait  time.Duration

	rskClient   *ethclient.Client
	btcClient   *rpcclient.Client
	lbcExecutor contracts.LiquidityBridgeContractExecutor
	// --------------------------------
}

type TestArguments struct {
	To               string   `json:"to"`
	RskRefundAddress string   `json:"rskRefundAddress"`
	ValueToTransfer  *big.Int `json:"valueToTransfer"`
}

func TestPegOutSuite(t *testing.T) {
	suite.Run(t, new(PegOutSuite))
}

func (s *PegOutSuite) SetupSuite() {
	config := integration.ReadTestConfig(s.T())
	s.serverUrl = config.Lps.Url
	s.args = TestArguments{
		To:               config.Tests.Pegout.DestinationAddress,
		ValueToTransfer:  config.Tests.Pegout.Value,
		RskRefundAddress: config.Tests.Pegout.RefundAddress,
	}
	s.userAccountKey = config.Rsk.UserPrivateKey
	s.maxPegoutWait = time.Second * time.Duration(config.Tests.Pegout.PegoutPaymentWait)
	s.maxRefundWait = time.Second * time.Duration(config.Tests.Pegout.PegoutRefundWait)

	var btcParams *chaincfg.Params
	switch config.Network {
	case "testnet":
		btcParams = &chaincfg.TestNet3Params
	case "mainnet":
		btcParams = &chaincfg.MainNetParams
	case "regtest":
		btcParams = &chaincfg.RegressionNetParams
	default:
		panic("invalid network")
	}
	rsk, err := ethclient.Dial(config.Rsk.RpcUrl)
	if err != nil {
		panic(err)
	}
	s.rskClient = rsk
	btcConfig := rpcclient.ConnConfig{
		Host:         config.Btc.RpcEndpoint,
		User:         config.Btc.User,
		Pass:         config.Btc.Password,
		Params:       btcParams.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	btc, err := rpcclient.New(&btcConfig, nil)
	if err != nil {
		panic(err)
	}
	s.btcClient = btc
	lbcExecutor, err := contracts.NewLegacyLbcExecutor(config.Rsk.LbcAddress, rsk)
	if err != nil {
		panic(err)
	}
	s.lbcExecutor = lbcExecutor
}

func (s *PegOutSuite) Raw() *suite.Suite {
	return &s.Suite
}

func (s *PegOutSuite) RskClient() *ethclient.Client {
	return s.rskClient
}

func (s *PegOutSuite) TestSuccessfulPegOutFlow() {
	s.Run("Should get a pegout quote", func() {
		s.executeGetQuote()
	})

	s.Run("Should be able to accept pegout quote", func() {
		s.executeAcceptQuote()
	})

	s.Run("Should transfer bitcoin to the user after they deposit the peg out", func() {
		s.depositPegout()
		s.receiveBtc()
	})

	s.Run("Should refund pegout to liquidity provider", func() {
		s.verifyRefundPegout()
	})
}

func (s *PegOutSuite) executeGetQuote() {
	body := pkg.PegoutQuoteRequest{
		To:               s.args.To,
		ValueToTransfer:  s.args.ValueToTransfer,
		RskRefundAddress: s.args.RskRefundAddress,
	}

	result, err := integration.ExecuteHttpRequest[[]pkg.GetPegoutQuoteResponse](integration.Execution{
		Method: http.MethodPost,
		URL:    s.serverUrl + "/pegout/getQuotes",
		Body:   body,
	})
	s.Raw().Require().NoError(err)

	expectedFields := []string{
		"lbcAddress",
		"liquidityProviderRskAddress",
		"btcRefundAddress",
		"rskRefundAddress",
		"lpBtcAddr",
		"callFee",
		"penaltyFee",
		"nonce",
		"depositAddr",
		"value",
		"agreementTimestamp",
		"depositDateLimit",
		"depositConfirmations",
		"transferConfirmations",
		"transferTime",
		"expireDate",
		"expireBlocks",
		"gasFee",
		"productFeeAmount",
	}

	var rawResponse []map[string]any
	err = json.Unmarshal(result.RawResponse, &rawResponse)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, result.StatusCode)
	s.Len(rawResponse, 1)
	s.NotEmpty(rawResponse[0]["quoteHash"])
	s.NotEmpty(rawResponse[0]["quote"])
	quoteFields, ok := rawResponse[0]["quote"].(map[string]any)
	s.True(ok, "Quote is not an object")

	integration.AssertFields(&s.Suite, expectedFields, quoteFields)
	s.QuoteResponse = result.Response[0]
}

func (s *PegOutSuite) executeAcceptQuote() {
	body := pkg.AcceptQuoteRequest{QuoteHash: s.QuoteResponse.QuoteHash}
	result, err := integration.ExecuteHttpRequest[pkg.AcceptPegoutResponse](integration.Execution{
		Method: http.MethodPost,
		URL:    s.serverUrl + "/pegout/acceptQuote",
		Body:   body,
	})
	s.Require().NoError(err)

	expectedFields := []string{
		"signature",
		"lbcAddress",
	}

	s.Equal(http.StatusOK, result.StatusCode)
	var rawResponse map[string]any
	err = json.Unmarshal(result.RawResponse, &rawResponse)
	s.Require().NoError(err, "Response does not have required format")
	integration.AssertFields(&s.Suite, expectedFields, rawResponse)
	s.AcceptQuoteResponse = result.Response
}

func (s *PegOutSuite) depositPegout() {
	var err error
	ctx := context.Background()
	privateKey, err := crypto.HexToECDSA(s.userAccountKey)
	s.Require().NoError(err)
	chainId, err := s.rskClient.ChainID(ctx)
	s.Require().NoError(err)
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	s.Require().NoError(err)
	receipt, _ := s.lbcExecutor.DepositPegout(s, opts, s.QuoteResponse.Quote, s.AcceptQuoteResponse.Signature)
	s.Require().Equal(uint64(1), receipt.Status)
}

func (s *PegOutSuite) receiveBtc() {
	var height, newHeight int64
	var pegoutTx string
	var err error

	height, err = s.btcClient.GetBlockCount()
	s.Require().NoError(err)
	if pegoutTx = s.findPegoutTransaction(height); pegoutTx != "" {
		log.Infof("[Integration test] Pegout received in %s", pegoutTx)
		return
	}

	testTolerance := time.NewTimer(s.maxPegoutWait)
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer ticker.Stop()
	defer testTolerance.Stop()

waitLoop:
	for {
		select {
		case <-ticker.C:
			if newHeight, err = s.btcClient.GetBlockCount(); err != nil {
				s.T().Fatal(err)
			} else if newHeight == height {
				return
			}
			height++
			if pegoutTx = s.findPegoutTransaction(height); pegoutTx != "" {
				log.Infof("[Integration test] Pegout received in %s", pegoutTx)
				done <- syscall.SIGTERM
			}
		case <-testTolerance.C:
			s.T().Fatalf("timeout waiting for pegout transaction")
		case <-done:
			break waitLoop
		}
	}
}

func (s *PegOutSuite) verifyRefundPegout() {
	refund := s.lbcExecutor.GetRefundPegoutEvent(s, s.maxPegoutWait, s.QuoteResponse.QuoteHash)
	s.Equal(s.QuoteResponse.QuoteHash, refund.QuoteHash)
}

func (s *PegOutSuite) findPegoutTransaction(height int64) string {
	const (
		QuoteHashOutputIndex  = 1
		QuoteHashOutputPrefix = "6a20"
	)
	var btcAmount btcutil.Amount
	blockHash, err := s.btcClient.GetBlockHash(height)
	s.Require().NoError(err)
	block, err := s.btcClient.GetBlockVerboseTx(blockHash)
	s.Require().NoError(err)

	for _, tx := range block.Tx {
		for _, output := range tx.Vout {
			if output.ScriptPubKey.Address == s.args.To {
				btcAmount, err = btcutil.NewAmount(output.Value)
				s.Require().NoError(err)
				weiAmount := entities.SatoshiToWei(uint64(btcAmount.ToUnit(btcutil.AmountSatoshi)))
				s.Equal(s.args.ValueToTransfer, weiAmount.AsBigInt())
				s.Equal(QuoteHashOutputPrefix+s.QuoteResponse.QuoteHash, tx.Vout[QuoteHashOutputIndex].ScriptPubKey.Hex)
				return tx.Txid
			}
		}
	}
	return ""
}
