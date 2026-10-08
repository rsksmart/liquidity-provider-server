package rootstock

import (
	"context"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

var (
	ErrPegInClaimAlreadyExists = errors.New("pegin claim already exists")
	ErrPegInClaimNotFound      = errors.New("pegin claim not found")
	// Owned here so payable-value and the adapter share ErrorIs. blockchain aliases this var.
	ErrIncorrectFronting = errors.New("incorrect fronting")
)

const PegInClaimCompletedEventId entities.EventId = "PegInClaimCompleted"

type PegInClaimState string

const (
	PegInClaimCandidate        PegInClaimState = "candidate"
	PegInClaimSubmitting       PegInClaimState = "submitting"
	PegInClaimClaimed          PegInClaimState = "claimed"
	PegInClaimRaceLost         PegInClaimState = "race_lost"
	PegInClaimRetryableFailure PegInClaimState = "retryable_failure"
	PegInClaimResolved         PegInClaimState = "resolved"
	PegInClaimResolveFailed    PegInClaimState = "resolve_failed"
)

type PegInClaim struct {
	RskAddress    string          `json:"rskAddress" bson:"rsk_address"`
	DepositTxID   string          `json:"depositTxId" bson:"deposit_txid"`
	BtcAddress    string          `json:"btcAddress" bson:"btc_address"`
	State         PegInClaimState `json:"state" bson:"state"`
	RequestTxHash string          `json:"requestTxHash" bson:"request_tx_hash"`
	// Empty when a successful receipt has no PegInRequested event to unpack.
	PegInID         string        `json:"pegInId" bson:"peg_in_id"`
	ResolveTxHash   string        `json:"resolveTxHash" bson:"resolve_tx_hash"`
	ResolveGasUsed  uint64        `json:"resolveGasUsed" bson:"resolve_gas_used"`
	ResolveGasPrice *entities.Wei `json:"resolveGasPrice" bson:"resolve_gas_price"`
	ClaimerPayout   *entities.Wei `json:"claimerPayout" bson:"claimer_payout"`
	RegistrantFee   *entities.Wei `json:"registrantFee" bson:"registrant_fee"`
	CreatedAt       time.Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updatedAt" bson:"updated_at"`
}

type PegInClaimRepository interface {
	Insert(context.Context, PegInClaim) error
	Get(ctx context.Context, rskAddress, depositTxID string) (*PegInClaim, error)
	Update(context.Context, PegInClaim) error
	ListByStates(ctx context.Context, states ...PegInClaimState) ([]PegInClaim, error)
}

func CalculatePegInClaimPayableValue(amount, fee *entities.Wei) (*entities.Wei, error) {
	if amount == nil || fee == nil {
		return nil, ErrIncorrectFronting
	}
	if amount.Cmp(fee) < 0 {
		return nil, ErrIncorrectFronting
	}
	return new(entities.Wei).Sub(amount, fee), nil
}

func (claim *PegInClaim) IsTerminal() bool {
	if claim == nil {
		return false
	}
	return claim.State == PegInClaimRaceLost ||
		claim.State == PegInClaimResolved ||
		claim.State == PegInClaimResolveFailed
}

func (claim *PegInClaim) IsClaimed() bool {
	return claim != nil && claim.State == PegInClaimClaimed
}

func (claim *PegInClaim) IsResolvable() bool {
	return claim.IsClaimed() && claim.PegInID != ""
}

func NewCandidatePegInClaim(entry PegInWatch, depositTxID string, existing *PegInClaim) PegInClaim {
	now := time.Now().UTC()
	claim := PegInClaim{
		RskAddress:  entry.RskAddress,
		DepositTxID: depositTxID,
		BtcAddress:  entry.BtcAddress,
		State:       PegInClaimCandidate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if existing != nil {
		claim.CreatedAt = existing.CreatedAt
	}
	return claim
}

type PegInClaimCompletedEvent struct {
	entities.Event
	Claim PegInClaim
}

func NewPegInClaimCompletedEvent(claim PegInClaim) PegInClaimCompletedEvent {
	return PegInClaimCompletedEvent{
		Event: entities.NewBaseEvent(PegInClaimCompletedEventId),
		Claim: claim,
	}
}
