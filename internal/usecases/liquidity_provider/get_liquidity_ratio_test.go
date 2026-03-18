package liquidity_provider

import (
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateTarget(t *testing.T) {
	tests := []struct {
		name         string
		maxLiquidity *entities.Wei
		percentage   uint64
		expected     *entities.Wei
	}{
		{
			name:         "50% of 1000",
			maxLiquidity: entities.NewWei(1000),
			percentage:   50,
			expected:     entities.NewWei(500),
		},
		{
			name:         "10% of 1000",
			maxLiquidity: entities.NewWei(1000),
			percentage:   10,
			expected:     entities.NewWei(100),
		},
		{
			name:         "90% of 1000",
			maxLiquidity: entities.NewWei(1000),
			percentage:   90,
			expected:     entities.NewWei(900),
		},
		{
			name:         "50% of 0",
			maxLiquidity: entities.NewWei(0),
			percentage:   50,
			expected:     entities.NewWei(0),
		},
		{
			name:         "0% of 1000",
			maxLiquidity: entities.NewWei(1000),
			percentage:   0,
			expected:     entities.NewWei(0),
		},
		{
			name:         "100% of 1000",
			maxLiquidity: entities.NewWei(1000),
			percentage:   100,
			expected:     entities.NewWei(1000),
		},
		{
			name:         "30% of large value (10 ether)",
			maxLiquidity: entities.NewBigWei(new(big.Int).Mul(big.NewInt(10), big.NewInt(1_000_000_000_000_000_000))),
			percentage:   30,
			expected:     entities.NewBigWei(new(big.Int).Mul(big.NewInt(3), big.NewInt(1_000_000_000_000_000_000))),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculateTarget(tc.maxLiquidity, tc.percentage)
			require.NoError(t, err)
			assert.Equal(t, tc.expected.AsBigInt().String(), result.AsBigInt().String())
		})
	}
}

func TestCalculateThreshold(t *testing.T) {
	t.Run("fixed tolerance", func(t *testing.T) {
		target := entities.NewWei(500)
		tolerance := lpEntity.ExcessTolerance{
			IsFixed:         true,
			FixedValue:      entities.NewWei(100),
			PercentageValue: utils.NewBigFloat64(0),
		}

		result := calculateThreshold(target, tolerance)

		assert.Equal(t, "600", result.AsBigInt().String())
	})

	t.Run("percentage tolerance 20%", func(t *testing.T) {
		target := entities.NewWei(500)
		tolerance := lpEntity.ExcessTolerance{
			IsFixed:         false,
			FixedValue:      entities.NewWei(0),
			PercentageValue: utils.NewBigFloat64(20),
		}

		result := calculateThreshold(target, tolerance)

		assert.Equal(t, "600", result.AsBigInt().String())
	})

	t.Run("percentage tolerance 0%", func(t *testing.T) {
		target := entities.NewWei(500)
		tolerance := lpEntity.ExcessTolerance{
			IsFixed:         false,
			FixedValue:      entities.NewWei(0),
			PercentageValue: utils.NewBigFloat64(0),
		}

		result := calculateThreshold(target, tolerance)

		assert.Equal(t, "500", result.AsBigInt().String())
	})

	t.Run("fixed tolerance zero", func(t *testing.T) {
		target := entities.NewWei(500)
		tolerance := lpEntity.ExcessTolerance{
			IsFixed:         true,
			FixedValue:      entities.NewWei(0),
			PercentageValue: utils.NewBigFloat64(0),
		}

		result := calculateThreshold(target, tolerance)

		assert.Equal(t, "500", result.AsBigInt().String())
	})
}

//nolint:funlen
func TestCalculateImpact(t *testing.T) {
	tests := []struct {
		name           string
		balance        *entities.Wei
		target         *entities.Wei
		threshold      *entities.Wei
		expectedType   NetworkImpactType
		expectedAmount *entities.Wei
	}{
		{
			name:           "deficit when balance below target",
			balance:        entities.NewWei(300),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactDeficit,
			expectedAmount: entities.NewWei(200),
		},
		{
			name:           "excess when balance above threshold",
			balance:        entities.NewWei(700),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactExcess,
			expectedAmount: entities.NewWei(200),
		},
		{
			name:           "within tolerance when balance equals target",
			balance:        entities.NewWei(500),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactWithinTolerance,
			expectedAmount: entities.NewWei(0),
		},
		{
			name:           "within tolerance when balance between target and threshold",
			balance:        entities.NewWei(550),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactWithinTolerance,
			expectedAmount: entities.NewWei(0),
		},
		{
			name:           "within tolerance when balance equals threshold",
			balance:        entities.NewWei(600),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactWithinTolerance,
			expectedAmount: entities.NewWei(0),
		},
		{
			name:           "excess when balance just above threshold",
			balance:        entities.NewWei(601),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactExcess,
			expectedAmount: entities.NewWei(101),
		},
		{
			name:           "deficit when balance just below target",
			balance:        entities.NewWei(499),
			target:         entities.NewWei(500),
			threshold:      entities.NewWei(600),
			expectedType:   NetworkImpactDeficit,
			expectedAmount: entities.NewWei(1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateImpact(tc.balance, tc.target, tc.threshold)
			assert.Equal(t, tc.expectedType, result.Type)
			assert.Equal(t, tc.expectedAmount.AsBigInt().String(), result.Amount.AsBigInt().String())
		})
	}
}
