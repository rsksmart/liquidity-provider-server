package pkg_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToPegoutQuoteDTO(t *testing.T) {
	pegoutQuote := quote.PegoutQuote{
		LbcAddress:            "0x12",
		LpRskAddress:          "0x34",
		BtcRefundAddress:      "btc1",
		RskRefundAddress:      "0x56",
		LpBtcAddress:          "btc2",
		CallFee:               entities.NewWei(5),
		PenaltyFee:            10,
		Nonce:                 15,
		DepositAddress:        "btc3",
		Value:                 entities.NewWei(20),
		AgreementTimestamp:    25,
		DepositDateLimit:      30,
		DepositConfirmations:  35,
		TransferConfirmations: 40,
		TransferTime:          45,
		ExpireDate:            50,
		ExpireBlock:           55,
		GasFee:                entities.NewWei(60),
		ProductFeeAmount:      65,
	}

	dto := pkg.ToPegoutQuoteDTO(pegoutQuote)

	assert.Equal(t, pegoutQuote.LbcAddress, dto.LBCAddr)
	assert.Equal(t, pegoutQuote.LpRskAddress, dto.LPRSKAddr)
	assert.Equal(t, pegoutQuote.BtcRefundAddress, dto.BtcRefundAddr)
	assert.Equal(t, pegoutQuote.RskRefundAddress, dto.RSKRefundAddr)
	assert.Equal(t, pegoutQuote.LpBtcAddress, dto.LpBTCAddr)
	assert.Equal(t, pegoutQuote.CallFee.Uint64(), dto.CallFee)
	assert.Equal(t, pegoutQuote.PenaltyFee, dto.PenaltyFee)
	assert.Equal(t, pegoutQuote.Nonce, dto.Nonce)
	assert.Equal(t, pegoutQuote.DepositAddress, dto.DepositAddr)
	assert.Equal(t, pegoutQuote.Value.Uint64(), dto.Value)
	assert.Equal(t, pegoutQuote.AgreementTimestamp, dto.AgreementTimestamp)
	assert.Equal(t, pegoutQuote.DepositDateLimit, dto.DepositDateLimit)
	assert.Equal(t, pegoutQuote.DepositConfirmations, dto.DepositConfirmations)
	assert.Equal(t, pegoutQuote.TransferConfirmations, dto.TransferConfirmations)
	assert.Equal(t, pegoutQuote.TransferTime, dto.TransferTime)
	assert.Equal(t, pegoutQuote.ExpireDate, dto.ExpireDate)
	assert.Equal(t, pegoutQuote.ExpireBlock, dto.ExpireBlock)
	assert.Equal(t, pegoutQuote.GasFee.Uint64(), dto.GasFee)
	assert.Equal(t, pegoutQuote.ProductFeeAmount, dto.ProductFeeAmount)
	const expectedFields = 19
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(pegoutQuote))
}

func TestToRetainedPegoutQuoteDTO(t *testing.T) {
	pegoutQuote := quote.RetainedPegoutQuote{
		QuoteHash:            "0x12",
		Signature:            "0x34",
		DepositAddress:       "btc1",
		RequiredLiquidity:    entities.NewWei(5),
		State:                quote.PegoutStateWaitingForDepositConfirmations,
		UserRskTxHash:        "0x56",
		LpBtcTxHash:          "btc2",
		RefundPegoutTxHash:   "0x78",
		RefundPegoutGasUsed:  entities.NewWei(800),
		RefundPegoutGasPrice: entities.NewWei(8),
		BridgeRefundTxHash:   "0x90",
		BridgePegoutGasUsed:  entities.NewWei(900),
		BridgePegoutGasPrice: entities.NewWei(9),
		LpBtcTxFee:           entities.NewWei(1000),
	}

	dto := pkg.ToRetainedPegoutQuoteDTO(pegoutQuote)

	assert.Equal(t, pegoutQuote.QuoteHash, dto.QuoteHash)
	assert.Equal(t, pegoutQuote.Signature, dto.Signature)
	assert.Equal(t, pegoutQuote.DepositAddress, dto.DepositAddress)
	assert.Equal(t, pegoutQuote.RequiredLiquidity.Uint64(), dto.RequiredLiquidity.Uint64())
	assert.Equal(t, string(pegoutQuote.State), dto.State)
	assert.Equal(t, pegoutQuote.UserRskTxHash, dto.UserRskTxHash)
	assert.Equal(t, pegoutQuote.LpBtcTxHash, dto.LpBtcTxHash)
	assert.Equal(t, pegoutQuote.LpBtcTxFee.AsBigInt(), dto.LpBtcTxFee)
	assert.Equal(t, pegoutQuote.RefundPegoutTxHash, dto.RefundPegoutTxHash)
	assert.Equal(t, pegoutQuote.RefundPegoutGasUsed.AsBigInt(), dto.RefundPegoutGasUsed)
	assert.Equal(t, pegoutQuote.RefundPegoutGasPrice.AsBigInt(), dto.RefundPegoutGasPrice)
	assert.Equal(t, pegoutQuote.BridgeRefundTxHash, dto.BridgeRefundTxHash)
	assert.Equal(t, pegoutQuote.BridgePegoutGasUsed.AsBigInt(), dto.BridgeRefundGasUsed)
	assert.Equal(t, pegoutQuote.BridgePegoutGasPrice.AsBigInt(), dto.BridgeRefundGasPrice)
	const expectedFields = 14
	assert.Equal(t, expectedFields, test.CountNonZeroValues(dto))
	assert.Equal(t, expectedFields, test.CountNonZeroValues(pegoutQuote))
}

func TestFromPegoutQuoteDTO(t *testing.T) {
	pegoutCreationData := quote.PegoutCreationData{
		FeeRate:       utils.NewBigFloat64(31251.3333333),
		FeePercentage: utils.NewBigFloat64(11.57),
		GasPrice:      entities.NewWei(10003124124),
		FixedFee:      entities.NewWei(15000000),
	}

	dto := pkg.ToPegoutCreationDataDTO(pegoutCreationData)

	feeRate, _ := pegoutCreationData.FeeRate.Native().Float64()
	feePercentage, _ := pegoutCreationData.FeePercentage.Native().Float64()
	assert.InDelta(t, feeRate, dto.FeeRate, 0.000000001)
	assert.InDelta(t, feePercentage, dto.FeePercentage, 0.000000001)
	assert.Equal(t, pegoutCreationData.GasPrice.Uint64(), dto.GasPrice)
	assert.Equal(t, pegoutCreationData.FixedFee.Uint64(), dto.FixedFee)
}
