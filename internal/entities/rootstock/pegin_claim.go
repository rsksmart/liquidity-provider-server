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
	// ErrIncorrectFronting is the same sentinel as blockchain.ErrIncorrectFronting.
	ErrIncorrectFronting = errors.New("incorrect fronting")
)

type PegInClaimState string

const (
	PegInClaimCandidate        PegInClaimState = "candidate"
	PegInClaimSubmitting       PegInClaimState = "submitting"
	PegInClaimClaimed          PegInClaimState = "claimed"
	PegInClaimRaceLost         PegInClaimState = "race_lost"
	PegInClaimRetryableFailure PegInClaimState = "retryable_failure"
)

type PegInClaim struct {
	RskAddress  string          `json:"rskAddress" bson:"rsk_address"`
	DepositTxID string          `json:"depositTxId" bson:"deposit_txid"`
	BtcAddress  string          `json:"btcAddress" bson:"btc_address"`
	State       PegInClaimState `json:"state" bson:"state"`
	TxHash      string          `json:"txHash" bson:"tx_hash"`
	// PegInID records the contract event identifier. It is empty when settlement
	// uses a successful receipt without a decodable event.
	PegInID     string        `json:"pegInId" bson:"peg_in_id"`
	ReservedWei *entities.Wei `json:"reservedWei" bson:"reserved_wei"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updated_at"`
}

type PegInClaimRepository interface {
	Insert(context.Context, PegInClaim) error
	Get(ctx context.Context, rskAddress, depositTxID string) (*PegInClaim, error)
	Update(context.Context, PegInClaim) error
	ListByStates(ctx context.Context, states ...PegInClaimState) ([]PegInClaim, error)
}

// CalculatePegInClaimPayableValue returns amount - fee without mutating either input.
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
	return claim.State == PegInClaimClaimed || claim.State == PegInClaimRaceLost
}

func NewCandidatePegInClaim(entry PegInWatch, depositTxID string, reserved *entities.Wei, existing *PegInClaim) PegInClaim {
	now := time.Now().UTC()
	copied := entities.NewWei(0)
	if reserved != nil {
		copied = reserved.Copy()
	}
	claim := PegInClaim{
		RskAddress:  entry.RskAddress,
		DepositTxID: depositTxID,
		BtcAddress:  entry.BtcAddress,
		State:       PegInClaimCandidate,
		ReservedWei: copied,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if existing != nil {
		claim.CreatedAt = existing.CreatedAt
	}
	return claim
}
