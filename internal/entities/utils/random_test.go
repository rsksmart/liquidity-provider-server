package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func TestGetRandomInt(t *testing.T) {
	var numbers []int64
	var number int64
	var err error
	for i := 0; i < 100; i++ {
		number, err = utils.GetRandomInt()
		assert.Positive(t, number)
		assert.False(t, slices.Contains(numbers, number))
		require.NoError(t, err)
		numbers = append(numbers, number)
	}
}

func TestMustGetRandomInt(t *testing.T) {
	var numbers []int64
	var number int64
	for i := 0; i < 100; i++ {
		number = utils.MustGetRandomInt()
		assert.Positive(t, number)
		assert.False(t, slices.Contains(numbers, number))
		numbers = append(numbers, number)
	}
}

func TestGetRandomBytes_Size(t *testing.T) {
	sizes := []int64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}
	for _, size := range sizes {
		bytes, err := utils.GetRandomBytes(size)
		require.NoError(t, err)
		require.Len(t, bytes, int(size))
	}
}

func TestGetRandomBytes_Random(t *testing.T) {
	const size = 32
	var generatedBytes [][]byte
	for i := 0; i < 100; i++ {
		bytes, err := utils.GetRandomBytes(size)
		require.NoError(t, err)
		require.Len(t, bytes, size)
		for _, generated := range generatedBytes {
			assert.NotEqual(t, generated, bytes)
		}
		generatedBytes = append(generatedBytes, bytes)
	}
}

func TestMustGetRandomBytes(t *testing.T) {
	const size = 32
	var generatedBytes [][]byte
	for i := 0; i < 100; i++ {
		bytes := utils.MustGetRandomBytes(size)
		require.Len(t, bytes, size)
		for _, generated := range generatedBytes {
			assert.NotEqual(t, generated, bytes)
		}
		generatedBytes = append(generatedBytes, bytes)
	}
}
