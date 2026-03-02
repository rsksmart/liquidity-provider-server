package utils

import (
	"errors"
	"math"
	"math/big"
)

var OverFlowError = errors.New("uint overflow")

func SafeSub(a, b uint64) uint64 {
	if a < b {
		return 0
	}
	return a - b
}

func SafeAdd(a, b uint64) (uint64, error) {
	var highest, lowest uint64
	if a > b {
		highest, lowest = a, b
	} else {
		highest, lowest = b, a
	}
	if lowest > math.MaxUint64-highest {
		return 0, OverFlowError
	}
	return a + b, nil
}

func RoundToNDecimals(value float64, decimals uint) float64 {
	ratio := math.Pow(10, float64(decimals))
	return math.Round(value*ratio) / ratio
}

// ApplyPercentageIncrease calculates value * (1 + percentage/100) using integer arithmetic
// to avoid floating-point precision issues.
//
// The calculation uses a scale factor (Scale constant = 10_000) for precision:
// result = value * (Scale + basisPoints) / Scale
// where basisPoints = percentage * 100
func ApplyPercentageIncrease(value *big.Int, percentage *big.Float) *big.Int {
	// Convert percentage to basis points (percentage * 100)
	hundred := big.NewFloat(100)
	basisPointsFloat := new(big.Float).Mul(percentage, hundred)
	basisPointsInt, _ := basisPointsFloat.Int(nil)

	// Calculate: value * (Scale + basisPoints) / Scale
	scale := big.NewInt(Scale)
	multiplier := new(big.Int).Add(scale, basisPointsInt)
	numerator := new(big.Int).Mul(value, multiplier)
	result := new(big.Int).Div(numerator, scale)

	return result
}
