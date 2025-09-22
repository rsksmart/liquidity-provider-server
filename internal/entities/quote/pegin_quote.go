package quote

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

const (
	AcceptedPeginQuoteEventId     entities.EventId = "AcceptedPeginQuote"
	CallForUserCompletedEventId   entities.EventId = "CallForUserCompleted"
	RegisterPeginCompletedEventId entities.EventId = "RegisterPeginCompleted"
)

type PeginState string

const (
	PeginStateWaitingForDeposit              PeginState = "WaitingForDeposit"
	PeginStateWaitingForDepositConfirmations PeginState = "WaitingForDepositConfirmations"
	PeginStateTimeForDepositElapsed          PeginState = "TimeForDepositElapsed"
	PeginStateCallForUserSucceeded           PeginState = "CallForUserSucceeded"
	PeginStateCallForUserFailed              PeginState = "CallForUserFailed"
	PeginStateRegisterPegInSucceeded         PeginState = "RegisterPegInSucceeded"
	PeginStateRegisterPegInFailed            PeginState = "RegisterPegInFailed"
)

type PeginQuoteRepository interface {
	InsertQuote(ctx context.Context, quote CreatedPeginQuote) error
	GetQuote(ctx context.Context, hash string) (*PeginQuote, error)
	GetPeginCreationData(ctx context.Context, hash string) PeginCreationData
	GetQuotesByHashesAndDate(ctx context.Context, hashes []string, startDate, endDate time.Time) ([]PeginQuote, error)
	GetRetainedQuote(ctx context.Context, hash string) (*RetainedPeginQuote, error)
	InsertRetainedQuote(ctx context.Context, quote RetainedPeginQuote) error
	UpdateRetainedQuote(ctx context.Context, quote RetainedPeginQuote) error
	GetRetainedQuoteByState(ctx context.Context, states ...PeginState) ([]RetainedPeginQuote, error)
	// DeleteQuotes deletes both regular and retained quotes
	DeleteQuotes(ctx context.Context, quotes []string) (uint, error)
	ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time, page, perPage int) ([]PeginQuoteWithRetained, int, error)
	GetRetainedQuotesForAddress(ctx context.Context, address string, states ...PeginState) ([]RetainedPeginQuote, error)
}

type PeginQuoteWithRetained struct {
	Quote         PeginQuote
	RetainedQuote RetainedPeginQuote
}

type CreatedPeginQuote struct {
	Hash         string
	Quote        PeginQuote
	CreationData PeginCreationData
}

type PeginCreationData struct {
	GasPrice      *entities.Wei   `json:"gasPrice" bson:"gas_price" validate:"required"`
	FeePercentage *utils.BigFloat `json:"feePercentage" bson:"fee_percentage" validate:"required"`
	FixedFee      *entities.Wei   `json:"fixedFee" bson:"fixed_fee" validate:"required"`
}

func PeginCreationDataZeroValue() PeginCreationData {
	return PeginCreationData{
		GasPrice:      entities.NewWei(0),
		FeePercentage: utils.NewBigFloat64(0),
		FixedFee:      entities.NewWei(0),
	}
}

type PeginQuote struct {
	FedBtcAddress      string        `json:"fedBTCAddress" bson:"fed_address" validate:"required"`
	LbcAddress         string        `json:"lbcAddress" bson:"lbc_address" validate:"required"`
	LpRskAddress       string        `json:"lpRskAddress" bson:"lp_rsk_address"  validate:"required"`
	BtcRefundAddress   string        `json:"btcRefundAddress" bson:"btc_refund_address"  validate:"required"`
	RskRefundAddress   string        `json:"rskRefundAddress" bson:"rsk_refund_address"  validate:"required"`
	LpBtcAddress       string        `json:"lpBtcAddress" bson:"lp_btc_address"  validate:"required"`
	CallFee            *entities.Wei `json:"callFee" bson:"call_fee"  validate:"required"`
	PenaltyFee         *entities.Wei `json:"penaltyFee" bson:"penalty_fee"  validate:"required"`
	ContractAddress    string        `json:"contractAddress" bson:"contract_address"  validate:"required"`
	Data               string        `json:"data" bson:"data"  validate:""`
	GasLimit           uint32        `json:"gasLimit,omitempty" bson:"gas_limit"  validate:"required"`
	Nonce              int64         `json:"nonce" bson:"nonce"  validate:"required"`
	Value              *entities.Wei `json:"value" bson:"value"  validate:"required"`
	AgreementTimestamp uint32        `json:"agreementTimestamp" bson:"agreement_timestamp"  validate:"required"`
	TimeForDeposit     uint32        `json:"timeForDeposit" bson:"time_for_deposit"  validate:"required"`
	LpCallTime         uint32        `json:"lpCallTime" bson:"lp_call_time"  validate:"required"`
	Confirmations      uint16        `json:"confirmations" bson:"confirmations"  validate:"required"`
	CallOnRegister     bool          `json:"callOnRegister" bson:"call_on_register"`
	GasFee             *entities.Wei `json:"gasFee" bson:"gas_fee"  validate:"required"`
	ProductFeeAmount   *entities.Wei `json:"productFeeAmount" bson:"product_fee_amount"  validate:""`
}

func (quote *PeginQuote) ExpireTime() time.Time {
	return time.Unix(int64(quote.AgreementTimestamp+quote.TimeForDeposit), 0)
}

func (quote *PeginQuote) IsExpired() bool {
	return time.Now().After(quote.ExpireTime())
}

func (quote *PeginQuote) Total() *entities.Wei {
	if quote.Value == nil {
		quote.Value = entities.NewWei(0)
	}
	if quote.CallFee == nil {
		quote.CallFee = entities.NewWei(0)
	}
	if quote.GasFee == nil {
		quote.GasFee = entities.NewWei(0)
	}
	if quote.ProductFeeAmount == nil {
		quote.ProductFeeAmount = entities.NewWei(0)
	}
	total := new(entities.Wei)
	total.Add(total, quote.Value)
	total.Add(total, quote.CallFee)
	total.Add(total, quote.ProductFeeAmount)
	total.Add(total, quote.GasFee)
	return total
}

type RetainedPeginQuote struct {
	QuoteHash             string        `json:"quoteHash" bson:"quote_hash" validate:"required"`
	DepositAddress        string        `json:"depositAddress" bson:"deposit_address" validate:"required"`
	Signature             string        `json:"signature" bson:"signature" validate:"required"`
	RequiredLiquidity     *entities.Wei `json:"requiredLiquidity" bson:"required_liquidity" validate:"required"`
	State                 PeginState    `json:"state" bson:"state" validate:"required"`
	UserBtcTxHash         string        `json:"userBtcTxHash" bson:"user_btc_tx_hash"`
	CallForUserTxHash     string        `json:"callForUserTxHash" bson:"call_for_user_tx_hash"`
	RegisterPeginTxHash   string        `json:"registerPeginTxHash" bson:"register_pegin_tx_hash"`
	CallForUserGasUsed    uint64        `json:"callForUserGasUsed" bson:"call_for_user_gas_used"`
	CallForUserGasPrice   *entities.Wei `json:"callForUserGasPrice" bson:"call_for_user_gas_price"`
	RegisterPeginGasUsed  uint64        `json:"registerPeginGasUsed" bson:"register_pegin_gas_used"`
	RegisterPeginGasPrice *entities.Wei `json:"registerPeginGasPrice" bson:"register_pegin_gas_price"`
	OwnerAccountAddress   string        `json:"ownerAccountAddress" bson:"owner_account_address"`
}

type WatchedPeginQuote struct {
	PeginQuote    PeginQuote
	RetainedQuote RetainedPeginQuote
	CreationData  PeginCreationData
}

func NewWatchedPeginQuote(peginQuote PeginQuote, retainedQuote RetainedPeginQuote, creationData PeginCreationData) WatchedPeginQuote {
	return WatchedPeginQuote{PeginQuote: peginQuote, RetainedQuote: retainedQuote, CreationData: creationData}
}

type AcceptedPeginQuoteEvent struct {
	entities.Event
	Quote         PeginQuote
	RetainedQuote RetainedPeginQuote
	CreationData  PeginCreationData
}

type CallForUserCompletedEvent struct {
	entities.Event
	PeginQuote    PeginQuote
	RetainedQuote RetainedPeginQuote
	CreationData  PeginCreationData
	Error         error
}

type RegisterPeginCompletedEvent struct {
	entities.Event
	RetainedQuote RetainedPeginQuote
	Error         error
}
