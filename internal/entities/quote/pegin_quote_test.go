package quote_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPeginQuote_Total(t *testing.T) {
	var result *entities.Wei
	quotes := test.Table[quote.PeginQuote, *entities.Wei]{
		{
			Value: quote.PeginQuote{
				CallFee:          nil,
				Value:            entities.NewWei(400000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: entities.NewWei(700000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            nil,
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: entities.NewWei(600000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            entities.NewWei(400000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 0,
			},
			Result: entities.NewWei(800000000000000000),
		},
		{
			Value: quote.PeginQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            entities.NewWei(400000000000000000),
				ProductFeeAmount: 200000000000000000,
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
				FedBtcAddress:    test.AnyAddress,
				LbcAddress:       test.AnyAddress,
				LpRskAddress:     test.AnyAddress,
				BtcRefundAddress: test.AnyAddress,
				RskRefundAddress: test.AnyAddress,
				LpBtcAddress:     test.AnyAddress,
				CallFee:          entities.NewWei(0),
				PenaltyFee:       entities.NewWei(1),
				ContractAddress:  test.AnyAddress,
				Data:             "",
				GasLimit:         1,
				Nonce:            1,
				Value:            entities.NewWei(400000000000000000),
				// nolint:gosec
				AgreementTimestamp: uint32(now - 61),
				// nolint:gosec
				TimeForDeposit:   uint32(time.Minute.Seconds()),
				LpCallTime:       1,
				Confirmations:    1,
				CallOnRegister:   false,
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: true,
		},
		{
			Value: quote.PeginQuote{
				FedBtcAddress:    test.AnyAddress,
				LbcAddress:       test.AnyAddress,
				LpRskAddress:     test.AnyAddress,
				BtcRefundAddress: test.AnyAddress,
				RskRefundAddress: test.AnyAddress,
				LpBtcAddress:     test.AnyAddress,
				CallFee:          entities.NewWei(300000000000000000),
				PenaltyFee:       entities.NewWei(1),
				ContractAddress:  test.AnyAddress,
				Data:             "",
				GasLimit:         1,
				Nonce:            1,
				Value:            entities.NewWei(0),
				// nolint:gosec
				AgreementTimestamp: uint32(now),
				// nolint:gosec
				TimeForDeposit:   uint32(time.Minute.Seconds()),
				LpCallTime:       1,
				Confirmations:    1,
				CallOnRegister:   false,
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: false,
		},
	}
	test.RunTable(t, quotes, func(value quote.PeginQuote) bool {
		return value.IsExpired()
	})
}
