package quote

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"time"
)

const (
	AcceptedPegoutQuoteEventId  entities.EventId = "AcceptedPegoutQuote"
	PegoutBtcSentEventId        entities.EventId = "PegoutBtcSent"
	PegoutQuoteCompletedEventId entities.EventId = "PegoutQuoteCompleted"
)

type PegoutState string

const (
	PegoutStateWaitingForDeposit              PegoutState = "WaitingForDeposit"
	PegoutStateTimeForDepositElapsed          PegoutState = "TimeForDepositElapsed"
	PegoutStateSendPegoutSucceeded            PegoutState = "SendPegoutSucceeded"
	PegoutStateSendPegoutFailed               PegoutState = "SendPegoutFailed"
	PegoutStateRefundPegOutSucceeded          PegoutState = "RefundPegOutSucceeded"
	PegoutStateRefundPegOutFailed             PegoutState = "RefundPegOutFailed"
	PegoutStateWaitingForDepositConfirmations PegoutState = "WaitingForDepositConfirmations"
	PegoutStateBridgeTxSucceeded              PegoutState = "BridgeTxSucceeded"
	PegoutStateBridgeTxFailed                 PegoutState = "BridgeTxFailed"
)

type PegoutQuoteRepository interface {
	InsertQuote(ctx context.Context, quote CreatedPegoutQuote) error
	GetPegoutCreationData(ctx context.Context, hash string) PegoutCreationData
	GetQuote(ctx context.Context, hash string) (*PegoutQuote, error)
	GetRetainedQuote(ctx context.Context, hash string) (*RetainedPegoutQuote, error)
	InsertRetainedQuote(ctx context.Context, quote RetainedPegoutQuote) error
	ListPegoutDepositsByAddress(ctx context.Context, address string) ([]PegoutDeposit, error)
	UpdateRetainedQuote(ctx context.Context, quote RetainedPegoutQuote) error
	UpdateRetainedQuotes(ctx context.Context, quotes []RetainedPegoutQuote) error
	GetRetainedQuoteByState(ctx context.Context, states ...PegoutState) ([]RetainedPegoutQuote, error)
	// DeleteQuotes deletes both regular and retained quotes
	DeleteQuotes(ctx context.Context, quotes []string) (uint, error)
	UpsertPegoutDeposit(ctx context.Context, deposit PegoutDeposit) error
	UpsertPegoutDeposits(ctx context.Context, deposits []PegoutDeposit) error
}

type CreatedPegoutQuote struct {
	Hash         string
	Quote        PegoutQuote
	CreationData PegoutCreationData
}

type PegoutCreationData struct {
	FeeRate       *utils.BigFloat `json:"feeRate" bson:"fee_rate" validate:"required"`
	FeePercentage *utils.BigFloat `json:"feePercentage" bson:"percentage_fee" validate:"required"`
	GasPrice      *entities.Wei   `json:"gasPrice" bson:"gas_price" validate:"required"`
	FixedFee      *entities.Wei   `json:"fixedFee" bson:"fixed_fee" validate:"required"`
}

func PegoutCreationDataZeroValue() PegoutCreationData {
	return PegoutCreationData{
		FeeRate:       utils.NewBigFloat64(0),
		FeePercentage: utils.NewBigFloat64(0),
		GasPrice:      entities.NewWei(0),
		FixedFee:      entities.NewWei(0),
	}
}

type PegoutQuote struct {
	LbcAddress            string        `json:"lbcAddress" bson:"lbc_address" validate:"required"`
	LpRskAddress          string        `json:"lpRskAddress" bson:"lp_rsk_address" validate:"required"`
	BtcRefundAddress      string        `json:"btcRefundAddress" bson:"btc_refund_address" validate:"required"`
	RskRefundAddress      string        `json:"rskRefundAddress" bson:"rsk_refund_address" validate:"required"`
	LpBtcAddress          string        `json:"lpBtcAddress" bson:"lp_btc_address" validate:"required"`
	CallFee               *entities.Wei `json:"callFee" bson:"call_fee" validate:"required"`
	PenaltyFee            uint64        `json:"penaltyFee" bson:"penalty_fee" validate:"required"`
	Nonce                 int64         `json:"nonce" bson:"nonce" validate:"required"`
	DepositAddress        string        `json:"depositAddress" bson:"deposit_address" validate:"required"`
	Value                 *entities.Wei `json:"value" bson:"value" validate:"required"`
	AgreementTimestamp    uint32        `json:"agreementTimestamp" bson:"agreement_timestamp" validate:"required"`
	DepositDateLimit      uint32        `json:"depositDateLimit" bson:"deposit_date_limit" validate:"required"`
	DepositConfirmations  uint16        `json:"depositConfirmations" bson:"deposit_confirmations" validate:"required"`
	TransferConfirmations uint16        `json:"transferConfirmations" bson:"transfer_confirmations" validate:"required"`
	TransferTime          uint32        `json:"transferTime" bson:"transfer_time" validate:"required"`
	ExpireDate            uint32        `json:"expireDate" bson:"expire_date" validate:"required"`
	ExpireBlock           uint32        `json:"expireBlocks" bson:"expire_blocks" validate:"required"`
	GasFee                *entities.Wei `json:"gasFee" bson:"gas_fee" validate:"required"`
	ProductFeeAmount      uint64        `json:"productFeeAmount" bson:"product_fee_amount" validate:""`
}

func (quote *PegoutQuote) ExpireTime() time.Time {
	return time.Unix(int64(quote.ExpireDate), 0)
}

func (quote *PegoutQuote) IsExpired() bool {
	return time.Now().After(quote.ExpireTime())
}

func GetCreationBlock(pegoutConfig liquidity_provider.PegoutConfiguration, pegoutQuote PegoutQuote) uint64 {
	return utils.SafeSub(uint64(pegoutQuote.ExpireBlock), pegoutConfig.ExpireBlocks)
}

func (quote *PegoutQuote) Total() *entities.Wei {
	if quote.Value == nil {
		quote.Value = entities.NewWei(0)
	}
	if quote.CallFee == nil {
		quote.CallFee = entities.NewWei(0)
	}
	if quote.GasFee == nil {
		quote.GasFee = entities.NewWei(0)
	}
	total := new(entities.Wei)
	total.Add(total, quote.Value)
	total.Add(total, quote.CallFee)
	total.Add(total, entities.NewUWei(quote.ProductFeeAmount))
	total.Add(total, quote.GasFee)
	return total
}

type RetainedPegoutQuote struct {
	QuoteHash          string        `json:"quoteHash" bson:"quote_hash" validate:"required"`
	DepositAddress     string        `json:"depositAddress" bson:"deposit_address" validate:"required"`
	Signature          string        `json:"signature" bson:"signature" validate:"required"`
	RequiredLiquidity  *entities.Wei `json:"requiredLiquidity" bson:"required_liquidity" validate:"required"`
	State              PegoutState   `json:"state" bson:"state" validate:"required"`
	UserRskTxHash      string        `json:"userRskTxHash" bson:"user_rsk_tx_hash"`
	LpBtcTxHash        string        `json:"lpBtcTxHash" bson:"lp_btc_tx_hash"`
	RefundPegoutTxHash string        `json:"refundPegoutTxHash" bson:"refund_pegout_tx_hash"`
	BridgeRefundTxHash string        `json:"BridgeRefundTxHash" bson:"bridge_refund_tx_hash"`
}

type WatchedPegoutQuote struct {
	PegoutQuote   PegoutQuote
	RetainedQuote RetainedPegoutQuote
	CreationData  PegoutCreationData
}

func NewWatchedPegoutQuote(pegoutQuote PegoutQuote, retainedQuote RetainedPegoutQuote, creationData PegoutCreationData) WatchedPegoutQuote {
	return WatchedPegoutQuote{PegoutQuote: pegoutQuote, RetainedQuote: retainedQuote, CreationData: creationData}
}

type AcceptedPegoutQuoteEvent struct {
	entities.Event
	Quote         PegoutQuote
	RetainedQuote RetainedPegoutQuote
	CreationData  PegoutCreationData
}

type PegoutDeposit struct {
	TxHash      string        `json:"txHash" bson:"tx_hash"`
	QuoteHash   string        `json:"quoteHash" bson:"quote_hash"`
	Amount      *entities.Wei `json:"amount" bson:"amount"`
	Timestamp   time.Time     `json:"timestamp" bson:"timestamp"`
	BlockNumber uint64        `json:"blockNumber" bson:"block_number"`
	From        string        `json:"from" bson:"from"`
}

func (deposit *PegoutDeposit) IsValidForQuote(quote PegoutQuote) bool {
	enoughAmount := deposit.Amount.Cmp(quote.Total()) >= 0
	nonExpiredInTime := deposit.Timestamp.Before(quote.ExpireTime())
	nonExpiredInBlocks := deposit.BlockNumber <= uint64(quote.ExpireBlock)
	return enoughAmount && nonExpiredInTime && nonExpiredInBlocks
}

type PegoutQuoteCompletedEvent struct {
	entities.Event
	RetainedQuote RetainedPegoutQuote
	Error         error
}

type PegoutBtcSentToUserEvent struct {
	entities.Event
	PegoutQuote   PegoutQuote
	RetainedQuote RetainedPegoutQuote
	CreationData  PegoutCreationData
	Error         error
}
