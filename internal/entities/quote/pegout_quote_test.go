package quote_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type LpMock struct {
	mock.Mock
	liquidity_provider.PegoutLiquidityProvider
}

func (l *LpMock) ExpireBlocksPegout() uint64 {
	return 40
}

func TestPegoutQuote_Total(t *testing.T) {
	var result *entities.Wei
	quotes := test.Table[quote.PegoutQuote, *entities.Wei]{
		{
			Value: quote.PegoutQuote{
				Value:            entities.NewWei(400000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: entities.NewWei(700000000000000000),
		},
		{
			Value: quote.PegoutQuote{
				CallFee:          entities.NewWei(300000000000000000),
				GasFee:           entities.NewWei(100000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: entities.NewWei(600000000000000000),
		},
		{
			Value: quote.PegoutQuote{
				CallFee: entities.NewWei(300000000000000000),
				Value:   entities.NewWei(400000000000000000),
				GasFee:  entities.NewWei(100000000000000000),
			},
			Result: entities.NewWei(800000000000000000),
		},
		{
			Value: quote.PegoutQuote{
				CallFee:          entities.NewWei(300000000000000000),
				Value:            entities.NewWei(400000000000000000),
				ProductFeeAmount: 200000000000000000,
			},
			Result: entities.NewWei(900000000000000000),
		},
	}
	test.RunTable(t, quotes, func(value quote.PegoutQuote) *entities.Wei {
		result = value.Total()
		assert.NotNil(t, value.Value)
		assert.NotNil(t, value.CallFee)
		assert.NotNil(t, value.GasFee)
		return result
	})
}

func TestPegoutQuote_IsExpired(t *testing.T) {
	now := time.Now().Unix()

	quotes := test.Table[quote.PegoutQuote, bool]{
		{
			Value: quote.PegoutQuote{
				LbcAddress:            test.AnyAddress,
				LpRskAddress:          test.AnyAddress,
				BtcRefundAddress:      test.AnyAddress,
				RskRefundAddress:      test.AnyAddress,
				LpBtcAddress:          test.AnyAddress,
				CallFee:               entities.NewWei(300000000000000000),
				PenaltyFee:            1,
				Nonce:                 1,
				DepositAddress:        test.AnyAddress,
				Value:                 entities.NewWei(0),
				AgreementTimestamp:    uint32(now),
				DepositDateLimit:      1,
				DepositConfirmations:  1,
				TransferTime:          1,
				TransferConfirmations: 1,
				ExpireDate:            uint32(now - 61),
				ExpireBlock:           1,
				GasFee:                entities.NewWei(100000000000000000),
				ProductFeeAmount:      200000000000000000,
			},
			Result: true,
		},
		{
			Value: quote.PegoutQuote{
				LbcAddress:            test.AnyAddress,
				LpRskAddress:          test.AnyAddress,
				BtcRefundAddress:      test.AnyAddress,
				RskRefundAddress:      test.AnyAddress,
				LpBtcAddress:          test.AnyAddress,
				CallFee:               entities.NewWei(300000000000000000),
				PenaltyFee:            1,
				Nonce:                 1,
				DepositAddress:        test.AnyAddress,
				Value:                 entities.NewWei(0),
				AgreementTimestamp:    uint32(now),
				DepositDateLimit:      1,
				DepositConfirmations:  1,
				TransferTime:          1,
				TransferConfirmations: 1,
				ExpireDate:            uint32(now + 60),
				ExpireBlock:           1,
				GasFee:                entities.NewWei(100000000000000000),
				ProductFeeAmount:      200000000000000000,
			},
			Result: false,
		},
	}
	test.RunTable(t, quotes, func(value quote.PegoutQuote) bool {
		return value.IsExpired()
	})
}

func TestGetCreationBlock(t *testing.T) {
	pegoutConfig := liquidity_provider.PegoutConfiguration{
		ExpireBlocks: 40,
	}
	quotes := test.Table[quote.PegoutQuote, uint64]{
		{
			Value: quote.PegoutQuote{
				LbcAddress:            test.AnyAddress,
				LpRskAddress:          test.AnyAddress,
				BtcRefundAddress:      test.AnyAddress,
				RskRefundAddress:      test.AnyAddress,
				LpBtcAddress:          test.AnyAddress,
				CallFee:               entities.NewWei(300000000000000000),
				PenaltyFee:            1,
				Nonce:                 1,
				DepositAddress:        test.AnyAddress,
				Value:                 entities.NewWei(0),
				AgreementTimestamp:    1,
				DepositDateLimit:      1,
				DepositConfirmations:  1,
				TransferTime:          1,
				TransferConfirmations: 1,
				ExpireDate:            1,
				ExpireBlock:           40,
				GasFee:                entities.NewWei(100000000000000000),
				ProductFeeAmount:      200000000000000000,
			},
			Result: 0,
		},
		{
			Value: quote.PegoutQuote{
				LbcAddress:            test.AnyAddress,
				LpRskAddress:          test.AnyAddress,
				BtcRefundAddress:      test.AnyAddress,
				RskRefundAddress:      test.AnyAddress,
				LpBtcAddress:          test.AnyAddress,
				CallFee:               entities.NewWei(300000000000000000),
				PenaltyFee:            1,
				Nonce:                 1,
				DepositAddress:        test.AnyAddress,
				Value:                 entities.NewWei(0),
				AgreementTimestamp:    1,
				DepositDateLimit:      1,
				DepositConfirmations:  1,
				TransferTime:          1,
				TransferConfirmations: 1,
				ExpireDate:            1,
				ExpireBlock:           380,
				GasFee:                entities.NewWei(100000000000000000),
				ProductFeeAmount:      200000000000000000,
			},
			Result: 340,
		},
	}
	test.RunTable(t, quotes, func(value quote.PegoutQuote) uint64 {
		return quote.GetCreationBlock(pegoutConfig, value)
	})
}

func TestPegoutDeposit_IsValidForQuote(t *testing.T) {
	now := time.Now()
	pegoutQuote := quote.PegoutQuote{
		ExpireBlock:      500,
		ExpireDate:       uint32(now.Unix()) + 60,
		Value:            entities.NewWei(200000000000000000),
		CallFee:          entities.NewWei(100000000000000000),
		GasFee:           entities.NewWei(100000000000000000),
		ProductFeeAmount: 100000000000000000,
	}
	cases := test.Table[quote.PegoutDeposit, bool]{
		{
			Value: quote.PegoutDeposit{
				TxHash:      test.AnyHash,
				QuoteHash:   test.AnyHash,
				Amount:      entities.NewWei(490000000000000000),
				Timestamp:   now,
				BlockNumber: 499,
				From:        test.AnyAddress,
			},
			Result: false,
		},
		{
			Value: quote.PegoutDeposit{
				TxHash:      test.AnyHash,
				QuoteHash:   test.AnyHash,
				Amount:      entities.NewWei(5100000000000000000),
				Timestamp:   time.Unix(now.Unix()+61, 0),
				BlockNumber: 499,
				From:        test.AnyAddress,
			},
			Result: false,
		},
		{
			Value: quote.PegoutDeposit{
				TxHash:      test.AnyHash,
				QuoteHash:   test.AnyHash,
				Amount:      entities.NewWei(5100000000000000000),
				Timestamp:   now,
				BlockNumber: 501,
				From:        test.AnyAddress,
			},
			Result: false,
		},
		{
			Value: quote.PegoutDeposit{
				TxHash:      test.AnyHash,
				QuoteHash:   test.AnyHash,
				Amount:      entities.NewWei(5100000000000000000),
				Timestamp:   now,
				BlockNumber: 499,
				From:        test.AnyAddress,
			},
			Result: true,
		},
	}
	test.RunTable(t, cases, func(value quote.PegoutDeposit) bool {
		return value.IsValidForQuote(pegoutQuote)
	})
}
