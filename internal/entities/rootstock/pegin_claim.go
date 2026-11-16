package rootstock

import (
	"context"
	"errors"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

var (
	PegInClaimExistsError   = errors.New("pegin claim already exists")
	PegInClaimNotFoundError = errors.New("pegin claim not found")
)

type PegInClaimState string

const (
	PegInClaimSubmitting       PegInClaimState = "submitting"
	PegInClaimClaimed          PegInClaimState = "claimed"
	PegInClaimRaceLost         PegInClaimState = "race_lost"
	PegInClaimRetryableFailure PegInClaimState = "retryable_failure"
)

type PegInClaim struct {
	RskAddress  string          `json:"rskAddress" bson:"rsk_address"`
	DepositTxID string          `json:"depositTxId" bson:"deposit_txid"`
	BtcAddress  string          `json:"btcAddress" bson:"btc_address"`
	PegInID     string          `json:"pegInId" bson:"peg_in_id"`
	State       PegInClaimState `json:"state" bson:"state"`
	TxHash      string          `json:"txHash" bson:"tx_hash"`
	Amount      *entities.Wei   `json:"amount" bson:"amount"`
	Fee         *entities.Wei   `json:"fee" bson:"fee"`
	GasFee      *entities.Wei   `json:"gasFee" bson:"gas_fee"`
	Value       *entities.Wei   `json:"value" bson:"value"`
	CallSuccess bool            `json:"callSuccess" bson:"call_success"`
	CreatedAt   time.Time       `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" bson:"updated_at"`
}

type PegInClaimRepository interface {
	Add(ctx context.Context, claim PegInClaim) error
	Get(ctx context.Context, rskAddress, depositTxID string) (*PegInClaim, error)
	Update(ctx context.Context, claim PegInClaim) error
}
