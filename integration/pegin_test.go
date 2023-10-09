package integration_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	lps "github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func (s *IntegrationTestSuite) TestSuccessfulPegInFlow() {

	var quote lps.QuoteReturn
	var acceptedQuote lps.AcceptRes
	URL := s.config.Lps.Url

	s.Run("Should be able to get pegin quote", func() {
		body := lps.QuoteRequest{
			CallEoaOrContractAddress: "0x79568c2989232dCa1840087D73d403602364c0D4",
			CallContractArguments:    "",
			ValueToTransfer:          600000000000000000,
			RskRefundAddress:         "0x79568c2989232dCa1840087D73d403602364c0D4",
			BitcoinRefundAddress:     "mxEp7KGqyjFiLnWJoXU6MNXpop8BYe9Gv1",
		}

		result, err := execute[[]lps.QuoteReturn](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegin/getQuote",
			Body:   body,
		})

		if err != nil {
			assert.Fail(s.T(), "Unexpected error: ", err)
		}

		expectedFields := []string{
			"callCost",
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
		if err != nil {
			assert.Fail(s.T(), "Response does not have required format")
		}
		assert.Equal(s.T(), http.StatusOK, result.StatusCode)
		assert.NotEmpty(s.T(), rawResponse[0]["quoteHash"])
		assert.NotEmpty(s.T(), rawResponse[0]["quote"])
		quoteFields, ok := rawResponse[0]["quote"].(map[string]any)
		if !ok {
			assert.Fail(s.T(), "Quote is not an object")
		}
		s.AssertFields(expectedFields, quoteFields)
		quote = result.Response[0]
	})

	s.Run("Should be able to accept pegin quote", func() {
		body := lps.AcceptReq{QuoteHash: quote.QuoteHash}
		result, err := execute[lps.AcceptRes](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegin/acceptQuote",
			Body:   body,
		})
		if err != nil {
			assert.Fail(s.T(), "Unexpected error: ", err)
		}

		expectedFields := []string{
			"signature",
			"bitcoinDepositAddressHash",
		}

		assert.Equal(s.T(), http.StatusOK, result.StatusCode)
		var rawResponse map[string]any
		err = json.Unmarshal(result.RawResponse, &rawResponse)
		if err != nil {
			assert.Fail(s.T(), "Response does not have required format")
		}
		s.AssertFields(expectedFields, rawResponse)
		acceptedQuote = result.Response
	})

	s.Run("Should process bitcoin deposit and callForUser", func() {
		address, err := btcutil.DecodeAddress(acceptedQuote.BitcoinDepositAddressHash, &s.btcParams)
		if err != nil {
			assert.Fail(s.T(), "Invalid derivation address")
		}
		value := types.NewWei(int64(quote.Quote.Value))
		callFee := types.NewWei(int64(quote.Quote.CallFee))
		callCost := types.NewWei(int64(quote.Quote.CallCost))
		totalFees := new(types.Wei).Add(callFee, callCost)
		totalAmount := new(types.Wei).Add(totalFees, value)
		floatAmount, _ := totalAmount.ToRbtc().Float64()
		btcAmount, err := btcutil.NewAmount(floatAmount)
		if err != nil {
			assert.Fail(s.T(), err.Error())
		}
		amount, _ := btcutil.NewAmount(0.00025)
		_ = s.btc.SetTxFee(amount)
		_, err = s.btc.SendToAddress(address, btcAmount)
		if err != nil {
			assert.FailNow(s.T(), "Error sending btc transaction")
		}

		eventChannel := make(chan *bindings.LiquidityBridgeContractCallForUser)
		lpAddress := common.HexToAddress(quote.Quote.LPRSKAddr)
		toAddress := common.HexToAddress(quote.Quote.ContractAddr)
		subscription, err := s.lbc.WatchCallForUser(
			nil,
			eventChannel,
			[]common.Address{lpAddress},
			[]common.Address{toAddress},
		)
		if err != nil {
			assert.FailNow(s.T(), "Error listening for callForUser")
		}

		select {
		case callForUser := <-eventChannel:
			subscription.Unsubscribe()
			assert.True(s.T(), callForUser.Success, "Call for user failed")
		case err = <-subscription.Err():
			assert.FailNow(s.T(), "Error listening for callForUser")
		}
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
			assert.FailNow(s.T(), "Error listening for callForUser")
		}

		select {
		case registerPegIn := <-eventChannel:
			subscription.Unsubscribe()
			assert.Positive(s.T(), registerPegIn.TransferredAmount.Int64(), "Register PegIn failed")
		case err = <-subscription.Err():
			assert.FailNow(s.T(), "Error listening for callForUser")
		}
	})
}
