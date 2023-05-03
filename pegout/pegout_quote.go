package pegout

import (
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
	GasLimit              uint32     `json:"gasLimit" db:"gas_limit" validate:"required"`
	Value                 *types.Wei `json:"value" db:"value" validate:"required"`
	AgreementTimestamp    uint32     `json:"agreementTimestamp" db:"agreement_timestamp" validate:"required"`
	DepositDateLimit      uint32     `json:"depositDateLimit" db:"deposit_date_limit" validate:"required"`
	DepositConfirmations  uint16     `json:"depositConfirmations" db:"deposit_confirmations" validate:"required"`
	TransferConfirmations uint16     `json:"transferConfirmations" db:"transfer_confirmations" validate:"required"`
	TransferTime          uint32     `json:"transferTime" db:"transfer_time" validate:"required"`
	ExpireDate            uint32     `json:"expireDate" db:"expire_date" validate:"required"`
	ExpireBlock           uint32     `json:"expireBlocks" db:"expire_blocks" validate:"required"`
}

func (q *Quote) GetExpirationTime() time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.DepositDateLimit), 0)
}

type QuoteState struct {
	StatusCode     uint8
	ReceivedAmount *big.Int
}

type DepositEvent struct {
	QuoteHash         string
	AccumulatedAmount *big.Int
	Timestamp         time.Time
	BlockNumber       uint64
}

func (event *DepositEvent) IsValidForQuote(quote *Quote) bool {
	enoughAmount := event.AccumulatedAmount.Cmp(new(types.Wei).Add(quote.Value, quote.CallFee).AsBigInt()) >= 0
	nonExpired := event.Timestamp.Before(quote.GetExpirationTime())
	return enoughAmount && nonExpired
}
