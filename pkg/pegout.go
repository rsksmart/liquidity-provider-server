package pkg

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"math/big"
	"time"
)

type PegoutQuoteRequest struct {
	To                   string `json:"to" required:"" validate:"required" description:"Bitcoin address that will receive the BTC amount"`
	ValueToTransfer      uint64 `json:"valueToTransfer" required:"" validate:"required" example:"10000000000000" description:"ValueToTransfer"`
	RskRefundAddress     string `json:"rskRefundAddress" required:"" validate:"required,eth_addr" example:"0x0" description:"RskRefundAddress"`
	BitcoinRefundAddress string `json:"bitcoinRefundAddress" required:"" validate:"required" example:"0x0" description:"BitcoinRefundAddress"`
}

type PegoutQuoteDTO struct {
	LBCAddr               string `json:"lbcAddress" required:"" validate:"required"`
	LPRSKAddr             string `json:"liquidityProviderRskAddress" required:"" validate:"required"`
	BtcRefundAddr         string `json:"btcRefundAddress" required:"" validate:"required"`
	RSKRefundAddr         string `json:"rskRefundAddress" required:"" validate:"required"`
	LpBTCAddr             string `json:"lpBtcAddr" required:"" validate:"required"`
	CallFee               uint64 `json:"callFee" required:"" validate:"required"`
	PenaltyFee            uint64 `json:"penaltyFee" required:"" validate:"required"`
	Nonce                 int64  `json:"nonce" required:"" validate:"required"`
	DepositAddr           string `json:"depositAddr" required:"" validate:"required"`
	Value                 uint64 `json:"value" required:"" validate:"required"`
	AgreementTimestamp    uint32 `json:"agreementTimestamp" required:"" validate:"required"`
	DepositDateLimit      uint32 `json:"depositDateLimit" required:"" validate:"required"`
	DepositConfirmations  uint16 `json:"depositConfirmations" required:"" validate:"required"`
	TransferConfirmations uint16 `json:"transferConfirmations" required:"" validate:"required"`
	TransferTime          uint32 `json:"transferTime" required:"" validate:"required"`
	ExpireDate            uint32 `json:"expireDate" required:"" validate:"required"`
	ExpireBlock           uint32 `json:"expireBlocks" required:"" validate:"required"`
	GasFee                uint64 `json:"gasFee" required:"" description:"Fee to pay for the gas of every call done during the pegout (call on behalf of the user in Bitcoin network and call to the dao fee collector in Rootstock)"`
	ProductFeeAmount      uint64 `json:"productFeeAmount" required:"" description:"The DAO fee amount"`
}

type RetainedPegoutQuoteDTO struct {
	QuoteHash          string   `json:"quoteHash" required:"" description:"32-byte long hash of the quote that acts as a unique identifier"`
	Signature          string   `json:"signature" required:"" description:"Signature of the liquidity provider expressing commitment on the quote"`
	DepositAddress     string   `json:"depositAddress" required:"" description:"Address of the smart contract where the user should execute depositPegout function"`
	RequiredLiquidity  *big.Int `json:"requiredLiquidity" required:"" description:"BTC liquidity that the LP locks to guarantee the service. It is different from the total amount that the user needs to pay."`
	State              string   `json:"state" required:"" description:"Current state of the quote. Possible values are:\n - WaitingForDeposit\n - WaitingForDepositConfirmations\n - TimeForDepositElapsed\n - SendPegoutSucceeded\n - SendPegoutFailed\n - RefundPegOutSucceeded\n - RefundPegOutFailed\n - BridgeTxSucceeded\n - BridgeTxFailed\n"`
	UserRskTxHash      string   `json:"userRskTxHash" required:"" description:"The hash of the depositPegout transaction made by the user"`
	LpBtcTxHash        string   `json:"lpBtcTxHash" required:"" description:"The hash of the BTC transaction from the LP to the user"`
	RefundPegoutTxHash string   `json:"refundPegoutTxHash" required:"" description:"The hash of the transaction from the LP to the LBC where the LP got the refund in RBTC"`
	BridgeRefundTxHash string   `json:"bridgeRefundTxHash" required:"" description:"The hash of the transaction from the LP to the bridge to convert the refunded RBTC into BTC"`
}

type PegoutQuoteStatusDTO struct {
	Detail PegoutQuoteDTO         `json:"detail" required:"" description:"Agreed specification of the quote"`
	Status RetainedPegoutQuoteDTO `json:"status" required:"" description:"Current status of the quote"`
}

func ToPegoutQuoteDTO(entity quote.PegoutQuote) PegoutQuoteDTO {
	return PegoutQuoteDTO{
		LBCAddr:               entity.LbcAddress,
		LPRSKAddr:             entity.LpRskAddress,
		BtcRefundAddr:         entity.BtcRefundAddress,
		RSKRefundAddr:         entity.RskRefundAddress,
		LpBTCAddr:             entity.LpBtcAddress,
		CallFee:               entity.CallFee.Uint64(),
		PenaltyFee:            entity.PenaltyFee,
		Nonce:                 entity.Nonce,
		DepositAddr:           entity.DepositAddress,
		Value:                 entity.Value.Uint64(),
		AgreementTimestamp:    entity.AgreementTimestamp,
		DepositDateLimit:      entity.DepositDateLimit,
		DepositConfirmations:  entity.DepositConfirmations,
		TransferConfirmations: entity.TransferConfirmations,
		TransferTime:          entity.TransferTime,
		ExpireDate:            entity.ExpireDate,
		ExpireBlock:           entity.ExpireBlock,
		GasFee:                entity.GasFee.Uint64(),
		ProductFeeAmount:      entity.ProductFeeAmount,
	}
}

func ToRetainedPegoutQuoteDTO(entity quote.RetainedPegoutQuote) RetainedPegoutQuoteDTO {
	return RetainedPegoutQuoteDTO{
		QuoteHash:          entity.QuoteHash,
		Signature:          entity.Signature,
		DepositAddress:     entity.DepositAddress,
		RequiredLiquidity:  entity.RequiredLiquidity.AsBigInt(),
		State:              string(entity.State),
		UserRskTxHash:      entity.UserRskTxHash,
		LpBtcTxHash:        entity.LpBtcTxHash,
		RefundPegoutTxHash: entity.RefundPegoutTxHash,
		BridgeRefundTxHash: entity.BridgeRefundTxHash,
	}
}

type GetPegoutQuoteResponse struct {
	Quote     PegoutQuoteDTO `json:"quote" required:"" description:"Detail of the quote"`
	QuoteHash string         `json:"quoteHash" required:"" description:"This is a 64 digit number that derives from a quote object"`
}

type AcceptPegoutResponse struct {
	Signature  string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	LbcAddress string `json:"lbcAddress" required:"" example:"0x0" description:"LBC address to execute depositPegout function"`
}

type DepositEventDTO struct {
	TxHash      string    `json:"-"`
	QuoteHash   string    `json:"quoteHash" example:"0x0" description:"QuoteHash"`
	Amount      *big.Int  `json:"amount" example:"10000" description:"Event Value"`
	Timestamp   time.Time `json:"timestamp" example:"10000" description:"Event Timestamp"`
	BlockNumber uint64    `json:"-"`
	From        string    `json:"from" example:"0x0" description:"From Address"`
}
