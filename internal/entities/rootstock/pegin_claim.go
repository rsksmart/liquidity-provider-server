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
	PegInID     string          `json:"pegInId" bson:"peg_in_id"`
	ReservedWei *entities.Wei   `json:"reservedWei" bson:"reserved_wei"`
	CreatedAt   time.Time       `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" bson:"updated_at"`
}

type PegInClaimRepository interface {
	Insert(context.Context, PegInClaim) error
	Get(ctx context.Context, rskAddress, depositTxID string) (*PegInClaim, error)
	Update(context.Context, PegInClaim) error
	ListByStates(ctx context.Context, states ...PegInClaimState) ([]PegInClaim, error)
}
