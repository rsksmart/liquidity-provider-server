package pegin_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	"github.com/rsksmart/liquidity-provider-server/test/integration/contracts"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"
)

type PegInSuite struct {
	suite.Suite
	QuoteResponse       pkg.GetPeginQuoteResponse
	AcceptQuoteResponse pkg.AcceptPeginRespose

	// --- initialized during setup ---
	serverUrl              string
	args                   TestArguments
	rpcWalletPassword      string
	rpcWalletUnlockSeconds int64
	maxCallForUserWait     time.Duration
	maxRegisterPeginWait   time.Duration

	btcParams   *chaincfg.Params
	rskClient   *ethclient.Client
	btcClient   *rpcclient.Client
	lbcExecutor contracts.LiquidityBridgeContractExecutor
	// ---------------------------------
}

type TestArguments struct {
	CallEoaOrContractAddress string   `json:"callEoaOrContractAddress"`
	CallContractArguments    string   `json:"callContractArguments"`
	ValueToTransfer          *big.Int `json:"valueToTransfer"`
	RskRefundAddress         string   `json:"rskRefundAddress"`
}

func TestPegInSuite(t *testing.T) {
	suite.Run(t, new(PegInSuite))
}

func (s *PegInSuite) SetupSuite() {
	config := integration.ReadTestConfig(s.T())
	s.serverUrl = config.Lps.Url
	s.args = TestArguments{
		CallEoaOrContractAddress: config.Tests.Pegin.DestinationAddress,
		CallContractArguments:    config.Tests.Pegin.Data,
		ValueToTransfer:          config.Tests.Pegin.Value,
		RskRefundAddress:         config.Tests.Pegin.RefundAddress,
	}
	s.rpcWalletPassword = config.Btc.WalletPassword
	s.rpcWalletUnlockSeconds = config.Btc.WalletUnlockSeconds
	s.maxCallForUserWait = time.Second * time.Duration(config.Tests.Pegin.CallForUserWait)
	s.maxRegisterPeginWait = time.Minute * time.Duration(config.Tests.Pegin.RegisterPeginWait)

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
	s.btcParams = btcParams
	rsk, err := ethclient.Dial(config.Rsk.RpcUrl)
	if err != nil {
		panic(err)
	}
	s.rskClient = rsk
	btcConfig := rpcclient.ConnConfig{
		Host:         config.Btc.RpcEndpoint,
		User:         config.Btc.User,
		Pass:         config.Btc.Password,
		Params:       s.btcParams.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	btc, err := rpcclient.New(&btcConfig, nil)
	if err != nil {
		panic(err)
	}
	s.btcClient = btc
	lbcExecutor, err := contracts.NewSplityLbcExecutor(contracts.SplitAddresses{
		Discovery:            config.Rsk.DiscoveryContract,
		Pegout:               config.Rsk.PegoutContract,
		Pegin:                config.Rsk.PeginContract,
		CollateralManagement: config.Rsk.CollateralManagementContract,
	}, rsk)
	if err != nil {
		panic(err)
	}
	s.lbcExecutor = lbcExecutor
}

func (s *PegInSuite) Raw() *suite.Suite {
	return &s.Suite
}

func (s *PegInSuite) RskClient() *ethclient.Client {
	return s.rskClient
}

func (s *PegInSuite) TestSuccessfulPegInFlow() {

	s.Run("Should get a pegin quote", func() {
		s.executeGetQuote()
	})

	s.Run("Should be able to accept pegin quote", func() {
		s.executeAcceptQuote()
	})

	s.Run("Should process bitcoin deposit and callForUser", func() {
		s.sendBtcTransaction()
		s.verifyCallForUser()
	})

	s.Run("Should call registerPegIn after proper confirmations", func() {
		s.verifyRegisterPegin()
	})
}

func (s *PegInSuite) executeGetQuote() {
	body := pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: s.args.CallEoaOrContractAddress,
		CallContractArguments:    s.args.CallContractArguments,
		ValueToTransfer:          s.args.ValueToTransfer,
		RskRefundAddress:         s.args.RskRefundAddress,
	}

	result, err := integration.ExecuteHttpRequest[[]pkg.GetPeginQuoteResponse](integration.Execution{
		Method: http.MethodPost,
		URL:    s.serverUrl + "/pegin/getQuote",
		Body:   body,
	})
	s.Raw().Require().NoError(err)

	expectedFields := []string{
		"gasFee",
		"callOnRegister",
		"callFee",
		"value",
		"gasLimit",
		"confirmations",
		"btcRefundAddr",
		"data",
		"lpRSKAddr",
		"fedBTCAddr",
		"lpBTCAddr",
		"contractAddr",
		"penaltyFee",
		"rskRefundAddr",
		"nonce",
		"timeForDeposit",
		"lpCallTime",
		"agreementTimestamp",
		"lbcAddr",
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

func (s *PegInSuite) executeAcceptQuote() {
	body := pkg.AcceptQuoteRequest{QuoteHash: s.QuoteResponse.QuoteHash}
	result, err := integration.ExecuteHttpRequest[pkg.AcceptPeginRespose](integration.Execution{
		Method: http.MethodPost,
		URL:    s.serverUrl + "/pegin/acceptQuote",
		Body:   body,
	})
	s.Require().NoError(err)

	expectedFields := []string{
		"signature",
		"bitcoinDepositAddressHash",
	}

	s.Equal(http.StatusOK, result.StatusCode)
	var rawResponse map[string]any
	err = json.Unmarshal(result.RawResponse, &rawResponse)
	s.Require().NoError(err, "Response does not have required format")
	integration.AssertFields(&s.Suite, expectedFields, rawResponse)
	s.AcceptQuoteResponse = result.Response
}

func (s *PegInSuite) sendBtcTransaction() {
	err := s.btcClient.WalletPassphrase(s.rpcWalletPassword, s.rpcWalletUnlockSeconds)
	s.Require().NoError(err)

	total := integration.SumAll(
		s.QuoteResponse.Quote.Value,
		s.QuoteResponse.Quote.GasFee,
		s.QuoteResponse.Quote.CallFee,
		s.QuoteResponse.Quote.ProductFeeAmount,
	)
	rbtcAmount, _ := entities.NewBigWei(total).ToRbtc().Float64()
	btcAmount, err := btcutil.NewAmount(rbtcAmount)
	s.Require().NoError(err)

	btcAddress, err := btcutil.DecodeAddress(s.AcceptQuoteResponse.BitcoinDepositAddressHash, s.btcParams)
	s.Require().NoError(err)

	txHash, err := s.btcClient.SendToAddress(btcAddress, btcAmount)
	s.Require().NoError(err)
	log.Infof("[Integration test] Pegin payment transaction hash: %s", txHash)
}

func (s *PegInSuite) verifyCallForUser() {
	callForUser := s.lbcExecutor.GetCallForUserEvent(
		s,
		s.maxCallForUserWait,
		s.args.CallEoaOrContractAddress,
		s.QuoteResponse.Quote.LPRSKAddr,
	)
	s.Equal(s.args.CallContractArguments, hex.EncodeToString(callForUser.Data))
	s.Equal(strings.ToLower(s.args.CallEoaOrContractAddress), strings.ToLower(callForUser.To))
	s.Equal(strings.ToLower(s.QuoteResponse.Quote.LPRSKAddr), strings.ToLower(callForUser.From))
	s.Equal(uint64(s.QuoteResponse.Quote.GasLimit), callForUser.GasLimit)
	s.Equal(s.QuoteResponse.QuoteHash, callForUser.QuoteHash)
	s.True(callForUser.Success)
	s.Equal(s.args.ValueToTransfer, callForUser.Value)
}

func (s *PegInSuite) verifyRegisterPegin() {
	registerPegin := s.lbcExecutor.GetPeginRegisteredEvent(
		s,
		s.maxRegisterPeginWait,
		s.QuoteResponse.QuoteHash,
	)
	total := integration.SumAll(
		s.QuoteResponse.Quote.Value,
		s.QuoteResponse.Quote.GasFee,
		s.QuoteResponse.Quote.CallFee,
		s.QuoteResponse.Quote.ProductFeeAmount,
	)
	var bigIntSats big.Int
	entities.NewBigWei(total).ToSatoshi().Int(&bigIntSats)
	s.Equal(bigIntSats.String(), registerPegin.Amount.ToSatoshi().String())
	s.Equal(s.QuoteResponse.QuoteHash, registerPegin.QuoteHash)
}
