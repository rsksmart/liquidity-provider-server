package pkg_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToPeginQuoteDTO(t *testing.T) {
	peginQuote := quote.PeginQuote{
		FedBtcAddress:      "0x12",
		LbcAddress:         "0x34",
		LpRskAddress:       "0x56",
		BtcRefundAddress:   "btc1",
		RskRefundAddress:   "0x90",
		LpBtcAddress:       "btc2",
		CallFee:            entities.NewWei(5),
		PenaltyFee:         entities.NewWei(10),
		ContractAddress:    "0xab",
		Data:               "cd",
		GasLimit:           15,
		Nonce:              20,
		Value:              entities.NewWei(25),
		AgreementTimestamp: 25,
		TimeForDeposit:     30,
		LpCallTime:         35,
		Confirmations:      40,
		CallOnRegister:     true,
		GasFee:             entities.NewWei(45),
		ProductFeeAmount:   50,
	}
	dto := pkg.ToPeginQuoteDTO(peginQuote)

	assert.Equal(t, peginQuote.FedBtcAddress, dto.FedBTCAddr)
	assert.Equal(t, peginQuote.LbcAddress, dto.LBCAddr)
	assert.Equal(t, peginQuote.LpRskAddress, dto.LPRSKAddr)
	assert.Equal(t, peginQuote.BtcRefundAddress, dto.BTCRefundAddr)
	assert.Equal(t, peginQuote.RskRefundAddress, dto.RSKRefundAddr)
	assert.Equal(t, peginQuote.LpBtcAddress, dto.LPBTCAddr)
	assert.Equal(t, peginQuote.CallFee.Uint64(), dto.CallFee)
	assert.Equal(t, peginQuote.PenaltyFee.Uint64(), dto.PenaltyFee)
	assert.Equal(t, peginQuote.ContractAddress, dto.ContractAddr)
	assert.Equal(t, peginQuote.Data, dto.Data)
	assert.Equal(t, peginQuote.GasLimit, dto.GasLimit)
	assert.Equal(t, peginQuote.Nonce, dto.Nonce)
	assert.Equal(t, peginQuote.Value.Uint64(), dto.Value)
	assert.Equal(t, peginQuote.AgreementTimestamp, dto.AgreementTimestamp)
	assert.Equal(t, peginQuote.TimeForDeposit, dto.TimeForDeposit)
	assert.Equal(t, peginQuote.LpCallTime, dto.LpCallTime)
	assert.Equal(t, peginQuote.Confirmations, dto.Confirmations)
	assert.Equal(t, peginQuote.CallOnRegister, dto.CallOnRegister)
	assert.Equal(t, peginQuote.GasFee.Uint64(), dto.GasFee)
	assert.Equal(t, peginQuote.ProductFeeAmount, dto.ProductFeeAmount)
	const expectedFields = 20
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(peginQuote))
}

func TestToRetainedPeginQuoteDTO(t *testing.T) {
	peginQuote := quote.RetainedPeginQuote{
		QuoteHash:           "0x12",
		Signature:           "0x34",
		DepositAddress:      "0x56",
		RequiredLiquidity:   entities.NewWei(5),
		State:               quote.PeginStateWaitingForDeposit,
		UserBtcTxHash:       "0x78",
		CallForUserTxHash:   "0x90",
		RegisterPeginTxHash: "0xab",
	}
	dto := pkg.ToRetainedPeginQuoteDTO(peginQuote)

	assert.Equal(t, peginQuote.QuoteHash, dto.QuoteHash)
	assert.Equal(t, peginQuote.Signature, dto.Signature)
	assert.Equal(t, peginQuote.DepositAddress, dto.DepositAddress)
	assert.Equal(t, peginQuote.RequiredLiquidity.AsBigInt(), dto.RequiredLiquidity)
	assert.Equal(t, string(peginQuote.State), dto.State)
	assert.Equal(t, peginQuote.UserBtcTxHash, dto.UserBtcTxHash)
	assert.Equal(t, peginQuote.CallForUserTxHash, dto.CallForUserTxHash)
	assert.Equal(t, peginQuote.RegisterPeginTxHash, dto.RegisterPeginTxHash)
	const expectedFields = 8
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(peginQuote))
}
