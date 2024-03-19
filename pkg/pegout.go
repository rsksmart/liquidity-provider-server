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
