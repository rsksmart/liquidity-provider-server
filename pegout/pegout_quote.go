package pegout

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider/types"
	"math/big"
	"time"
)

type Quote struct {
	LBCAddr               string     `json:"lbcAddress" db:"lbc_addr" validate:"required"`
	LPRSKAddr             string     `json:"liquidityProviderRskAddress" db:"lp_rsk_addr" validate:"required"`
	BtcRefundAddr         string     `json:"btcRefundAddress" db:"btc_refund_addr" validate:"required"`
	RSKRefundAddr         string     `json:"rskRefundAddress" db:"rsk_refund_addr" validate:"required"`
	LpBTCAddr             string     `json:"lpBtcAddr" db:"lp_btc_addr" validate:"required"`
	CallFee               *types.Wei `json:"callFee" db:"callFee" validate:"required"`
	PenaltyFee            uint64     `json:"penaltyFee" db:"penalty_fee" validate:"required"`
	Nonce                 int64      `json:"nonce" db:"nonce" validate:"required"`
	DepositAddr           string     `json:"depositAddr" db:"deposit_addr" validate:"required"`
	Value                 *types.Wei `json:"value" db:"value" validate:"required"`
	AgreementTimestamp    uint32     `json:"agreementTimestamp" db:"agreement_timestamp" validate:"required"`
	DepositDateLimit      uint32     `json:"depositDateLimit" db:"deposit_date_limit" validate:"required"`
	DepositConfirmations  uint16     `json:"depositConfirmations" db:"deposit_confirmations" validate:"required"`
	TransferConfirmations uint16     `json:"transferConfirmations" db:"transfer_confirmations" validate:"required"`
	TransferTime          uint32     `json:"transferTime" db:"transfer_time" validate:"required"`
	ExpireDate            uint32     `json:"expireDate" db:"expire_date" validate:"required"`
	ExpireBlock           uint32     `json:"expireBlocks" db:"expire_blocks" validate:"required"`
	CallCost              *types.Wei `json:"callCost" db:"callCost" validate:"required"`
}

func (q *Quote) GetExpirationTime() time.Time {
	return time.Unix(int64(q.ExpireDate), 0)
}

type QuoteState struct {
	StatusCode     uint8
	ReceivedAmount *big.Int
}

type DepositEvent struct {
	TxHash      common.Hash    `json:"-"`
	QuoteHash   string         `json:"quoteHash" example:"0x0" description:"QuoteHash"`
	Amount      *big.Int       `json:"amount" example:"10000" description:"Event Value"`
	Timestamp   time.Time      `json:"timestamp" example:"10000" description:"Event Timestamp"`
	BlockNumber uint64         `json:"-"`
	From        common.Address `json:"from" example:"0x0" description:"From Address"`
}

func (event *DepositEvent) IsValidForQuote(quote *Quote) bool {
	enoughAmount := event.Amount.Cmp(new(types.Wei).Add(quote.Value, quote.CallFee).AsBigInt()) >= 0
	nonExpired := event.Timestamp.Before(quote.GetExpirationTime())
	return enoughAmount && nonExpired
}
