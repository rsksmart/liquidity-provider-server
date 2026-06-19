package utils_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	}
	test.RunTable(t, errorCases, func(value args) error {
		_, err = utils.SafeAdd(value.A, value.B)
		return err
	})
	test.RunTable(t, successCases, func(value args) uint64 {
		result, err = utils.SafeAdd(value.A, value.B)
		require.NoError(t, err)
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

// nolint:funlen
func TestApplyPercentageIncrease(t *testing.T) {
	t.Run("20% increase on 1000", func(t *testing.T) {
		value := big.NewInt(1000)
		percentage := big.NewFloat(20.0)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 1000 * (1 + 20/100) = 1000 * 1.2 = 1200
		assert.Equal(t, "1200", result.String())
	})

	t.Run("10% increase on 5000", func(t *testing.T) {
		value := big.NewInt(5000)
		percentage := big.NewFloat(10.0)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 5000 * (1 + 10/100) = 5000 * 1.1 = 5500
		assert.Equal(t, "5500", result.String())
	})

	t.Run("0% increase on 1000", func(t *testing.T) {
		value := big.NewInt(1000)
		percentage := big.NewFloat(0.0)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 1000 * (1 + 0/100) = 1000
		assert.Equal(t, "1000", result.String())
	})

	t.Run("50% increase on 2000", func(t *testing.T) {
		value := big.NewInt(2000)
		percentage := big.NewFloat(50.0)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 2000 * (1 + 50/100) = 2000 * 1.5 = 3000
		assert.Equal(t, "3000", result.String())
	})

	t.Run("5.5% increase on 10000", func(t *testing.T) {
		value := big.NewInt(10000)
		percentage := big.NewFloat(5.5)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 10000 * (1 + 5.5/100) = 10000 * 1.055 = 10550
		assert.Equal(t, "10550", result.String())
	})

	t.Run("Large value with percentage", func(t *testing.T) {
		value := big.NewInt(1000000000000000000) // 1 ETH in wei
		percentage := big.NewFloat(15.0)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 1000000000000000000 * 1.15 = 1150000000000000000
		assert.Equal(t, "1150000000000000000", result.String())
	})

	t.Run("Fractional percentage 0.33%", func(t *testing.T) {
		value := big.NewInt(100000)
		percentage := big.NewFloat(0.33)

		result := utils.ApplyPercentageIncrease(value, percentage)

		// 100000 * (1 + 0.33/100) = 100000 * 1.0033 = 100330
		assert.Equal(t, "100330", result.String())
	})
}
