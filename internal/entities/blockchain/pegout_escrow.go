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
	RequestHash        string
	RefundAddress      string
	Amount             *entities.Wei
	DestinationAddress []byte
	TxHash             string
	BlockNumber        uint64
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
