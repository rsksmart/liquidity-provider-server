package quote

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
	"math/big"
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

type PegConfiguration interface {
	GetFixedFee() *entities.Wei
	GetFeePercentage() *utils.BigFloat
	ValidateAmount(amount *entities.Wei) error
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
		return errors.New("invalid quote hash format: must be a valid hex string")
	}

	return nil
}

func CalculateCallFee(amount *entities.Wei, config PegConfiguration) *entities.Wei {
	result := new(entities.Wei)

	percentageFee := calculatePercentageFee(amount, config.GetFeePercentage())
	result.Add(percentageFee, config.GetFixedFee())

	log.Debugf("Percentage fee: %v%% of %v = %v", config.GetFeePercentage(), amount, percentageFee)
	log.Debugf("Fixed fee: %v", config.GetFixedFee())
	log.Debugf("Call fee: %v + %v = %v", percentageFee, config.GetFixedFee(), result)
	return result
}

func calculatePercentageFee(amount *entities.Wei, percentage *utils.BigFloat) *entities.Wei {
	const scale = 1000 // the scale needs to have at least as many zeros as the amount of decimals we want to support in the percentage
	amountAsRat := new(big.Rat).SetInt(amount.AsBigInt())
	floatPercentage, _ := percentage.Native().Float64()

	percentageAsFraction := new(big.Rat).SetFrac(
		big.NewInt(int64(floatPercentage*scale)), // Scale to avoid precision loss
		big.NewInt(100*scale),
	)
	percentageFee := new(big.Rat).Mul(amountAsRat, percentageAsFraction)

	remainder := new(big.Int)
	result, _ := new(big.Int).QuoRem(
		percentageFee.Num(),
		percentageFee.Denom(),
		remainder,
	)

	// if remainder is more than half denominator round up
	if new(big.Int).Mul(remainder, big.NewInt(2)).Cmp(percentageFee.Denom()) >= 0 {
		result.Add(result, big.NewInt(1))
	}
	return entities.NewBigWei(result)
}
