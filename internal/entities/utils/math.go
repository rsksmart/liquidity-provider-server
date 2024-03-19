package utils

import (
	"errors"
	"math"
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
