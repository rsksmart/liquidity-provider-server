package blockchain

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

// EscrowedPegOutState is the on-chain lifecycle of one escrowed peg-out.
// NONE is the storage default so unset ids are never confused with REQUESTED.
type EscrowedPegOutState uint8

const (
	EscrowedPegOutStateNone EscrowedPegOutState = iota
	EscrowedPegOutStateRequested
	EscrowedPegOutStateClaimed
	EscrowedPegOutStateCancelled
	EscrowedPegOutStateFulfilled
	EscrowedPegOutStateRefunded
)

// PegOutRequested is the decoded PegOutRequested event of PegOutEscrow.
type PegOutRequested struct {
	RequestHash        string        `bson:"request_hash"`
	RefundAddress      string        `bson:"refund_address"`
	Amount             *entities.Wei `bson:"amount"`
	DestinationAddress []byte        `bson:"destination_address"`
	TxHash             string        `bson:"tx_hash"`
	BlockNumber        uint64        `bson:"block_number"`
}

// PegOutClaimed is the decoded PegOutClaimed event of PegOutEscrow.
type PegOutClaimed struct {
	LpAddress   string
	RequestHash string
	TxHash      string
	BlockNumber uint64
}

// PegOutCancelled is the decoded PegOutCancelled event of PegOutEscrow.
type PegOutCancelled struct {
	RequestHash string
	TxHash      string
	BlockNumber uint64
}

// PegOutEscrowContract is the LPS port over the frozen IPegOutEscrow ABI.
type PegOutEscrowContract interface {
	GetAddress() string
	GetPegOutState(requestHash string) (EscrowedPegOutState, error)
	GetPegOutQuote(requestHash string) (quote.PegoutQuote, error)
	GetMaxMinerFee(requestHash string) (*entities.Wei, error)
	RestrictedUntil(lpAddress string) (uint64, error)
	ClaimPegOut(txConfig TransactionConfig, requestHash string, signature []byte) (TransactionReceipt, error)
	// GetPegOutRequestedEvents returns PegOutRequested logs in [fromBlock, toBlock].
	// A nil toBlock reads up to the latest block.
	GetPegOutRequestedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]PegOutRequested, error)
	GetPegOutClaimedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]PegOutClaimed, error)
	GetPegOutCancelledEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]PegOutCancelled, error)
}

type PegOutEscrowWatchRepository interface {
	GetCheckpoint(ctx context.Context) (lastScannedBlock uint64, found bool, err error)
	SetCheckpoint(ctx context.Context, lastScannedBlock uint64) error
	UpsertCandidate(ctx context.Context, candidate PegOutRequested) error
	DeleteCandidate(ctx context.Context, requestHash string) error
	ListCandidates(ctx context.Context) ([]PegOutRequested, error)
}
