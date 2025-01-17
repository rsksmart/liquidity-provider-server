package quote

import (
	"encoding/hex"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

type AcceptedQuote struct {
	Signature      string `json:"signature"`
	DepositAddress string `json:"depositAddress"`
}

type Fees struct {
	CallFee          *entities.Wei
	GasFee           *entities.Wei
	PenaltyFee       *entities.Wei
	ProductFeeAmount uint64
}

// ValidateQuoteHash checks if a given string is a valid 32-byte quote hash
// Returns nil if valid, error otherwise
func ValidateQuoteHash(hash string) error {

	// Check length
	if len(hash) != 64 {
		return fmt.Errorf("invalid quote hash length: expected 64 characters, got %d", len(hash))
	}

	// Check if it's a valid hex string
	if _, err := hex.DecodeString(hash); err != nil {
		return fmt.Errorf("invalid quote hash format: must be a valid hex string")
	}

	return nil
}
