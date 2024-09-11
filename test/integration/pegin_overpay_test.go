package integration_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func (s *IntegrationTestSuite) TestPegInOverPayFlow() {
	var quote pkg.GetPeginQuoteResponse
	var acceptedQuote pkg.AcceptPeginRespose
	var extraAmount *big.Int
	var value, callFee, gasFee *entities.Wei
	URL := s.config.Lps.Url

	s.Run("Should be able to get pegin quote", func() { getPeginQuoteTest(s, URL, &quote) })
	value = entities.NewUWei(quote.Quote.Value)
	callFee = entities.NewUWei(quote.Quote.CallFee)
	gasFee = entities.NewUWei(quote.Quote.GasFee)
	s.Run("Should be able to accept pegin quote", func() { acceptPeginQuoteTest(s, URL, quote, &acceptedQuote) })
	s.Run("Should deposit 1.5 of the expected BTC and receive the requested amount in the callForUser", func() {
		productFee := entities.NewUWei(quote.Quote.ProductFeeAmount)
		totalFees := new(entities.Wei).Add(new(entities.Wei).Add(callFee, gasFee), productFee)
		totalAmount := new(entities.Wei).Add(totalFees, value)
		extraAmount = new(big.Int).Div(totalAmount.AsBigInt(), big.NewInt(2))
		callForUserTest(s, quote, acceptedQuote, new(entities.Wei).Add(totalAmount, entities.NewBigWei(extraAmount)))
	})
	s.Run("Should call registerPegIn and pay the extra 0.5 to the user in RBTC", func() {
		balanceIncreaseChannel := make(chan *bindings.LiquidityBridgeContractBalanceIncrease)
		balanceIncreaseSubscription, err := s.lbc.WatchBalanceIncrease(nil, balanceIncreaseChannel)
		s.Require().NoError(err, "Error listening for balance increase")

		refundChannel := make(chan *bindings.LiquidityBridgeContractRefund)
		refundSubscription, err := s.lbc.WatchRefund(nil, refundChannel)
		s.Require().NoError(err, "Error listening for refund")

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		var registered, refunded bool
		for !(registered && refunded) {
			select {
			case refund := <-refundChannel:
				refundSubscription.Unsubscribe()
				s.Require().Equal(extraAmount, refund.Amount, "User wasn't refunded with the correct amount")
				s.Require().Equal(quote.Quote.RSKRefundAddr, refund.Dest.String(), "User wasn't refunded to the correct address")
				refunded = true
			case balanceIncreaseEvent := <-balanceIncreaseChannel:
				balanceIncreaseSubscription.Unsubscribe()
				refundedAmount := new(entities.Wei).Add(value, new(entities.Wei).Add(callFee, gasFee))
				s.Require().Equal(strings.ToLower(quote.Quote.LPRSKAddr), strings.ToLower(balanceIncreaseEvent.Dest.String()))
				s.Require().Equal(refundedAmount.AsBigInt().Int64(), balanceIncreaseEvent.Amount.Int64())
				registered = true
			case err = <-balanceIncreaseSubscription.Err():
				if err != nil {
					s.FailNow("Error listening for registerPegIn", err)
				}
			case err = <-refundSubscription.Err():
				if err != nil {
					s.FailNow("Error listening for refund", err)
				}
			case <-done:
				balanceIncreaseSubscription.Unsubscribe()
				refundSubscription.Unsubscribe()
				s.FailNow("Test cancelled")
			}
		}
	})
}
