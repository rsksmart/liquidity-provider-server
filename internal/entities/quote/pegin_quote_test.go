package quote_test

import (
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
)

func TestPeginQuote_Total(t *testing.T) {
	var result *entities.Wei
	quotes := test.Table[quote.PeginQuote, *entities.Wei]{
		{
			Value: quote.PeginQuote{
				CallFee:          nil,
				Value:            entities.NewWei(400000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: entities.NewWei(200000000000000000),
			},
			Result: entities.NewWei(700000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            nil,
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: entities.NewWei(200000000000000000),
			},
			Result: entities.NewWei(600000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            entities.NewWei(400000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: entities.NewWei(0),
			},
			Result: entities.NewWei(800000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            entities.NewWei(400000000000000000),
				ProductFeeAmount: entities.NewWei(200000000000000000),
				GasFee:           nil,
			},
			Result: entities.NewWei(900000000000000000),
		},
	}
	test.RunTable(t, quotes, func(value quote.PeginQuote) *entities.Wei {
		result = value.Total()
		assert.NotNil(t, value.Value)
		assert.NotNil(t, value.CallFee)
		assert.NotNil(t, value.GasFee)
		return result
	})
}

func TestPeginQuote_IsExpired(t *testing.T) {
	now := time.Now().Unix()

	quotes := test.Table[quote.PeginQuote, bool]{
		{
			Value: quote.PeginQuote{
				FedBtcAddress:      test.AnyAddress,
				LbcAddress:         test.AnyAddress,
				LpRskAddress:       test.AnyAddress,
				BtcRefundAddress:   test.AnyAddress,
				RskRefundAddress:   test.AnyAddress,
				LpBtcAddress:       test.AnyAddress,
				CallFee:            entities.NewWei(0),
				PenaltyFee:         entities.NewWei(1),
				ContractAddress:    test.AnyAddress,
				Data:               "",
				GasLimit:           1,
				Nonce:              1,
				Value:              entities.NewWei(400000000000000000),
				AgreementTimestamp: uint32(now - 61),
				TimeForDeposit:     uint32(time.Minute.Seconds()),
				LpCallTime:         1,
				Confirmations:      1,
				CallOnRegister:     false,
				GasFee:             entities.NewWei(100000000000000000),
				ProductFeeAmount:   entities.NewWei(200000000000000000),
			},
			Result: true,
		},
		{
			Value: quote.PeginQuote{
				FedBtcAddress:      test.AnyAddress,
				LbcAddress:         test.AnyAddress,
				LpRskAddress:       test.AnyAddress,
				BtcRefundAddress:   test.AnyAddress,
				RskRefundAddress:   test.AnyAddress,
				LpBtcAddress:       test.AnyAddress,
				CallFee:            entities.NewWei(300000000000000000),
				PenaltyFee:         entities.NewWei(1),
				ContractAddress:    test.AnyAddress,
				Data:               "",
				GasLimit:           1,
				Nonce:              1,
				Value:              entities.NewWei(0),
				AgreementTimestamp: uint32(now),
				TimeForDeposit:     uint32(time.Minute.Seconds()),
				LpCallTime:         1,
				Confirmations:      1,
				CallOnRegister:     false,
				GasFee:             entities.NewWei(100000000000000000),
				ProductFeeAmount:   entities.NewWei(200000000000000000),
			},
			Result: false,
		},
	}
	test.RunTable(t, quotes, func(value quote.PeginQuote) bool {
		return value.IsExpired()
	})
}

//nolint:funlen
func TestEnsureRetainedPeginQuoteZeroValues(t *testing.T) {
	testCases := []struct {
		name     string
		input    quote.RetainedPeginQuote
		expected quote.RetainedPeginQuote
	}{
		{
			name: "should set nil gas prices to zero",
			input: quote.RetainedPeginQuote{
				QuoteHash:             "0x123",
				CallForUserGasPrice:   nil,
				RegisterPeginGasPrice: nil,
				CallForUserGasUsed:    100,
				RegisterPeginGasUsed:  200,
			},
			expected: quote.RetainedPeginQuote{
				QuoteHash:             "0x123",
				CallForUserGasPrice:   entities.NewWei(0),
				RegisterPeginGasPrice: entities.NewWei(0),
				CallForUserGasUsed:    100,
				RegisterPeginGasUsed:  200,
			},
		},
		{
			name: "should not modify existing non-nil gas prices",
			input: quote.RetainedPeginQuote{
				QuoteHash:             "0x456",
				CallForUserGasPrice:   entities.NewWei(1000),
				RegisterPeginGasPrice: entities.NewWei(2000),
				CallForUserGasUsed:    150,
			},
			expected: quote.RetainedPeginQuote{
				QuoteHash:             "0x456",
				CallForUserGasPrice:   entities.NewWei(1000),
				RegisterPeginGasPrice: entities.NewWei(2000),
				CallForUserGasUsed:    150,
			},
		},
		{
			name: "should handle mixed nil and non-nil values",
			input: quote.RetainedPeginQuote{
				QuoteHash:             "0x789",
				CallForUserGasPrice:   entities.NewWei(500),
				RegisterPeginGasPrice: nil,
				OwnerAccountAddress:   "0xowner",
			},
			expected: quote.RetainedPeginQuote{
				QuoteHash:             "0x789",
				CallForUserGasPrice:   entities.NewWei(500),
				RegisterPeginGasPrice: entities.NewWei(0),
				OwnerAccountAddress:   "0xowner",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualQuote := tc.input
			quote.EnsureRetainedPeginQuoteZeroValues(&actualQuote)

			assert.Equal(t, tc.expected.QuoteHash, actualQuote.QuoteHash)
			assert.Equal(t, tc.expected.CallForUserGasPrice, actualQuote.CallForUserGasPrice)
			assert.Equal(t, tc.expected.RegisterPeginGasPrice, actualQuote.RegisterPeginGasPrice)
			assert.Equal(t, tc.expected.CallForUserGasUsed, actualQuote.CallForUserGasUsed)
			assert.Equal(t, tc.expected.RegisterPeginGasUsed, actualQuote.RegisterPeginGasUsed)
			assert.Equal(t, tc.expected.OwnerAccountAddress, actualQuote.OwnerAccountAddress)
		})
	}

	t.Run("should be idempotent", func(t *testing.T) {
		originalQuote := quote.RetainedPeginQuote{
			QuoteHash:             "0xabc",
			CallForUserGasPrice:   nil,
			RegisterPeginGasPrice: nil,
		}

		quote.EnsureRetainedPeginQuoteZeroValues(&originalQuote)
		firstCallResult := originalQuote

		quote.EnsureRetainedPeginQuoteZeroValues(&originalQuote)
		secondCallResult := originalQuote

		assert.Equal(t, firstCallResult.CallForUserGasPrice, secondCallResult.CallForUserGasPrice)
		assert.Equal(t, firstCallResult.RegisterPeginGasPrice, secondCallResult.RegisterPeginGasPrice)
		assert.NotNil(t, secondCallResult.CallForUserGasPrice)
		assert.NotNil(t, secondCallResult.RegisterPeginGasPrice)
	})
}
