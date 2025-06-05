package pkg_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"math/big"
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
		Nonce:              quote.NewNonce(20),
		Value:              entities.NewWei(25),
		AgreementTimestamp: 25,
		TimeForDeposit:     30,
		LpCallTime:         35,
		Confirmations:      40,
		CallOnRegister:     true,
		GasFee:             entities.NewWei(45),
		ProductFeeAmount:   entities.NewWei(50),
	}
	dto := pkg.ToPeginQuoteDTO(peginQuote)

	assert.Equal(t, peginQuote.FedBtcAddress, dto.FedBTCAddr)
	assert.Equal(t, peginQuote.LbcAddress, dto.LBCAddr)
	assert.Equal(t, peginQuote.LpRskAddress, dto.LPRSKAddr)
	assert.Equal(t, peginQuote.BtcRefundAddress, dto.BTCRefundAddr)
	assert.Equal(t, peginQuote.RskRefundAddress, dto.RSKRefundAddr)
	assert.Equal(t, peginQuote.LpBtcAddress, dto.LPBTCAddr)
	assert.Equal(t, peginQuote.CallFee.AsBigInt(), dto.CallFee)
	assert.Equal(t, peginQuote.PenaltyFee.AsBigInt(), dto.PenaltyFee)
	assert.Equal(t, peginQuote.ContractAddress, dto.ContractAddr)
	assert.Equal(t, peginQuote.Data, dto.Data)
	assert.Equal(t, peginQuote.GasLimit, dto.GasLimit)
	assert.Equal(t, peginQuote.Nonce.AsBigInt(), dto.Nonce)
	assert.Equal(t, peginQuote.Value.AsBigInt(), dto.Value)
	assert.Equal(t, peginQuote.AgreementTimestamp, dto.AgreementTimestamp)
	assert.Equal(t, peginQuote.TimeForDeposit, dto.TimeForDeposit)
	assert.Equal(t, peginQuote.LpCallTime, dto.LpCallTime)
	assert.Equal(t, peginQuote.Confirmations, dto.Confirmations)
	assert.Equal(t, peginQuote.CallOnRegister, dto.CallOnRegister)
	assert.Equal(t, peginQuote.GasFee.AsBigInt(), dto.GasFee)
	assert.Equal(t, peginQuote.ProductFeeAmount.AsBigInt(), dto.ProductFeeAmount)
	const expectedFields = 20
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(peginQuote))
}

func TestFromPeginQuoteDTO(t *testing.T) {
	dto := pkg.PeginQuoteDTO{
		FedBTCAddr:         "0x12",
		LBCAddr:            "0x34",
		LPRSKAddr:          "0x56",
		BTCRefundAddr:      "btc1",
		RSKRefundAddr:      "0x90",
		LPBTCAddr:          "btc2",
		CallFee:            big.NewInt(5),
		PenaltyFee:         big.NewInt(10),
		ContractAddr:       "0xab",
		Data:               "cd",
		GasLimit:           15,
		Nonce:              big.NewInt(20),
		Value:              big.NewInt(25),
		AgreementTimestamp: 25,
		TimeForDeposit:     30,
		LpCallTime:         35,
		Confirmations:      40,
		CallOnRegister:     true,
		GasFee:             big.NewInt(45),
		ProductFeeAmount:   big.NewInt(50),
	}
	peginQuote := pkg.FromPeginQuoteDTO(dto)

	assert.Equal(t, dto.FedBTCAddr, peginQuote.FedBtcAddress)
	assert.Equal(t, dto.LBCAddr, peginQuote.LbcAddress)
	assert.Equal(t, dto.LPRSKAddr, peginQuote.LpRskAddress)
	assert.Equal(t, dto.BTCRefundAddr, peginQuote.BtcRefundAddress)
	assert.Equal(t, dto.RSKRefundAddr, peginQuote.RskRefundAddress)
	assert.Equal(t, dto.LPBTCAddr, peginQuote.LpBtcAddress)
	assert.Equal(t, dto.CallFee.String(), peginQuote.CallFee.String())
	assert.Equal(t, dto.PenaltyFee.String(), peginQuote.PenaltyFee.String())
	assert.Equal(t, dto.ContractAddr, peginQuote.ContractAddress)
	assert.Equal(t, dto.Data, peginQuote.Data)
	assert.Equal(t, dto.GasLimit, peginQuote.GasLimit)
	assert.Equal(t, dto.Nonce, peginQuote.Nonce.AsBigInt())
	assert.Equal(t, dto.Value.String(), peginQuote.Value.String())
	assert.Equal(t, dto.AgreementTimestamp, peginQuote.AgreementTimestamp)
	assert.Equal(t, dto.TimeForDeposit, peginQuote.TimeForDeposit)
	assert.Equal(t, dto.LpCallTime, peginQuote.LpCallTime)
	assert.Equal(t, dto.Confirmations, peginQuote.Confirmations)
	assert.Equal(t, dto.CallOnRegister, peginQuote.CallOnRegister)
	assert.Equal(t, dto.GasFee.String(), peginQuote.GasFee.String())
	assert.Equal(t, dto.ProductFeeAmount.String(), peginQuote.ProductFeeAmount.String())
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

func TestToPeginCreationDataDTO(t *testing.T) {
	peginCreationData := quote.PeginCreationData{
		GasPrice:      entities.NewWei(5),
		FeePercentage: utils.NewBigFloat64(10.54),
		FixedFee:      entities.NewWei(15000000),
	}
	dto := pkg.ToPeginCreationDataDTO(peginCreationData)

	feePercentage, _ := peginCreationData.FeePercentage.Native().Float64()
	assert.Equal(t, peginCreationData.GasPrice.String(), dto.GasPrice.String())
	assert.InDelta(t, feePercentage, dto.FeePercentage, 0.000000001)
	assert.Equal(t, peginCreationData.FixedFee.String(), dto.FixedFee.String())

	const expectedFields = 3
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(peginCreationData))
}
