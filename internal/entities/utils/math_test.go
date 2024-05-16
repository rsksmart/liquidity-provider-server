package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"math"
	"testing"
)

func TestSafeSub(t *testing.T) {
	type args struct{ A, B uint64 }
	cases := test.Table[args, uint64]{
		{Value: args{A: 7, B: 2}, Result: 5},
		{Value: args{A: 2, B: 7}, Result: 0},
		{Value: args{A: 7, B: 7}, Result: 0},
		{Value: args{A: 0, B: 0}, Result: 0},
		{Value: args{A: 5, B: math.MaxUint64}, Result: 0},
		{Value: args{A: math.MaxUint64, B: 1<<64 - 2}, Result: 1},
	}
	test.RunTable(t, cases, func(value args) uint64 {
		return utils.SafeSub(value.A, value.B)
	})
}

func TestSafeAdd(t *testing.T) {
	var err error
	var result uint64
	type args struct{ A, B uint64 }
	errorCases := test.Table[args, error]{
		{Value: args{A: 7, B: 2}, Result: nil},
		{Value: args{A: 2, B: 7}, Result: nil},
		{Value: args{A: 0, B: 0}, Result: nil},
		{Value: args{A: 1, B: math.MaxUint64}, Result: utils.OverFlowError},
		{Value: args{A: math.MaxUint64, B: 1}, Result: utils.OverFlowError},
		{Value: args{A: math.MaxUint64, B: math.MaxUint64}, Result: utils.OverFlowError},
	}

	successCases := test.Table[args, uint64]{
		{Value: args{A: 1<<64 - 2, B: 1}, Result: math.MaxUint64},
		{Value: args{A: 7, B: 2}, Result: 9},
		{Value: args{A: 2, B: 7}, Result: 9},
		{Value: args{A: 0, B: 0}, Result: 0},
		{Value: args{A: 1, B: math.MaxUint64}, Result: 0},
		{Value: args{A: math.MaxUint64, B: 1}, Result: 0},
		{Value: args{A: math.MaxUint64, B: math.MaxUint64}, Result: 0},
	}
	test.RunTable(t, errorCases, func(value args) error {
		_, err = utils.SafeAdd(value.A, value.B)
		return err
	})
	test.RunTable(t, successCases, func(value args) uint64 {
		result, _ = utils.SafeAdd(value.A, value.B)
		return result
	})
}

func TestRoundToNDecimals(t *testing.T) {
	type args struct {
		Value    float64
		Decimals uint
	}
	cases := test.Table[args, float64]{
		{Value: args{Value: 1.123456789, Decimals: 2}, Result: 1.12},
		{Value: args{Value: 1.123456789, Decimals: 3}, Result: 1.123},
		{Value: args{Value: 1.123456789, Decimals: 4}, Result: 1.1235},
		{Value: args{Value: 1.123456789, Decimals: 5}, Result: 1.12346},
		{Value: args{Value: 0.011998954000000001, Decimals: 10}, Result: 0.0119989540},
		{Value: args{Value: 5, Decimals: 10}, Result: 5},
		{Value: args{Value: -1.123456789, Decimals: 4}, Result: -1.1235},
		{Value: args{Value: -1.123456789, Decimals: 5}, Result: -1.12346},
	}
	test.RunTable(t, cases, func(value args) float64 {
		return utils.RoundToNDecimals(value.Value, value.Decimals)
	})
}
