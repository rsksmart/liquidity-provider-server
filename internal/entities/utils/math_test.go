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
