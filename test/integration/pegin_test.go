package integration_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (s *IntegrationTestSuite) TestSuccessfulPegInFlow() {
	var quote pkg.GetPeginQuoteResponse
	var acceptedQuote pkg.AcceptPeginRespose
	URL := s.config.Lps.Url

	s.Run("Should be able to get pegin quote", func() {
		getPeginQuoteTest(s, URL, &quote)
	})

	s.Run("Should be able to accept pegin quote", func() {
		acceptPeginQuoteTest(s, URL, quote, &acceptedQuote)
	})

	s.Run("Should process bitcoin deposit and callForUser", func() {
		value := entities.NewUWei(quote.Quote.Value)
		callFee := entities.NewUWei(quote.Quote.CallFee)
		gasFee := entities.NewUWei(quote.Quote.GasFee)
		productFee := entities.NewUWei(quote.Quote.ProductFeeAmount)
		totalFees := new(entities.Wei).Add(new(entities.Wei).Add(callFee, gasFee), productFee)
		totalAmount := new(entities.Wei).Add(totalFees, value)
		callForUserTest(s, quote, acceptedQuote, totalAmount)
	})

	s.Run("Should call registerPegIn after proper confirmations", func() {
		eventChannel := make(chan *bindings.LiquidityBridgeContractPegInRegistered)
		var quoteHash [32]byte
		hashBytes, _ := hex.DecodeString(quote.QuoteHash)
		copy(quoteHash[:], hashBytes)
		subscription, err := s.lbc.WatchPegInRegistered(
			nil,
			eventChannel,
			[][32]byte{quoteHash},
		)
		if err != nil {
			s.FailNow("Error listening for registerPegIn")
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case registerPegIn := <-eventChannel:
			subscription.Unsubscribe()
			s.Positive(registerPegIn.TransferredAmount.Int64(), "Register PegIn failed")
		case err = <-subscription.Err():
			if err != nil {
				s.FailNow("Error listening for registerPegIn", err)
			}
		case <-done:
			subscription.Unsubscribe()
			s.FailNow("Test cancelled")
		}
	})
}

func getPeginQuoteTest(s *IntegrationTestSuite, url string, quoteResponse *pkg.GetPeginQuoteResponse) {
	body := pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: "0x79568c2989232dCa1840087D73d403602364c0D4",
		CallContractArguments:    "",
		ValueToTransfer:          600000000000000000,
		RskRefundAddress:         "0x79568c2989232dCa1840087D73d403602364c0D4",
	}

	result, err := execute[[]pkg.GetPeginQuoteResponse](Execution{
		Method: http.MethodPost,
		URL:    url + "/pegin/getQuote",
		Body:   body,
	})
	s.NoError(err)

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
	}

	var rawResponse []map[string]any
	err = json.Unmarshal(result.RawResponse, &rawResponse)
	s.NoError(err)
	s.Equal(http.StatusOK, result.StatusCode)
	s.NotEmpty(rawResponse[0]["quoteHash"])
	s.NotEmpty(rawResponse[0]["quote"])
	quoteFields, ok := rawResponse[0]["quote"].(map[string]any)
	if !ok {
		s.FailNow("Quote is not an object")
	}
	s.AssertFields(expectedFields, quoteFields)
	*quoteResponse = result.Response[0]
}

func acceptPeginQuoteTest(s *IntegrationTestSuite, url string, quote pkg.GetPeginQuoteResponse, acceptedQuote *pkg.AcceptPeginRespose) {
	body := pkg.AcceptQuoteRequest{QuoteHash: quote.QuoteHash}
	result, err := execute[pkg.AcceptPeginRespose](Execution{
		Method: http.MethodPost,
		URL:    url + "/pegin/acceptQuote",
		Body:   body,
	})
	s.NoError(err)

	expectedFields := []string{"signature", "bitcoinDepositAddressHash"}

	s.Require().Equal(http.StatusOK, result.StatusCode)
	var rawResponse map[string]any
	err = json.Unmarshal(result.RawResponse, &rawResponse)
	s.NoError(err)
	s.AssertFields(expectedFields, rawResponse)
	*acceptedQuote = result.Response
}

func callForUserTest(s *IntegrationTestSuite, quote pkg.GetPeginQuoteResponse, acceptedQuote pkg.AcceptPeginRespose, amountToPay *entities.Wei) {
	address, err := btcutil.DecodeAddress(acceptedQuote.BitcoinDepositAddressHash, &s.btcParams)
	s.NoError(err)
	floatAmount, _ := amountToPay.ToRbtc().Float64()
	btcAmount, err := btcutil.NewAmount(floatAmount)
	s.NoError(err)
	amount, _ := btcutil.NewAmount(0.00025)
	err = s.btc.WalletPassphrase(s.config.Btc.WalletPassword, 60)
	s.Require().NoError(err)
	err = s.btc.SetTxFee(amount)
	s.Require().NoError(err)
	_, err = s.btc.SendToAddress(address, btcAmount)
	s.Require().NoError(err)

	eventChannel := make(chan *bindings.LiquidityBridgeContractCallForUser)
	lpAddress := common.HexToAddress(quote.Quote.LPRSKAddr)
	toAddress := common.HexToAddress(quote.Quote.ContractAddr)
	subscription, err := s.lbc.WatchCallForUser(
		nil,
		eventChannel,
		[]common.Address{lpAddress},
		[]common.Address{toAddress},
	)
	s.NoError(err)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case callForUser := <-eventChannel:
		subscription.Unsubscribe()
		s.Require().Equal(quote.Quote.Value, callForUser.Value.Uint64(), "Call for user value does not match")
		s.Require().True(callForUser.Success, "Call for user failed")
	case err = <-subscription.Err():
		if err != nil {
			s.FailNow("Error listening for callForUser", err)
		}
	case <-done:
		subscription.Unsubscribe()
		s.FailNow("Test cancelled")
	}
}
