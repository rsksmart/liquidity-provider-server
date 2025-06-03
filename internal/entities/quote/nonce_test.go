package quote_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"math"
	"testing"
)

func TestNonce_UnmarshalBSONValue(t *testing.T) {
	dataTypeCases := test.Table[bsontype.Type, error]{
		{Value: bson.TypeDBPointer, Result: entities.DeserializationError},
		{Value: bson.TypeBinary, Result: entities.DeserializationError},
		{Value: bson.TypeDouble, Result: entities.DeserializationError},
	}
	stringZeroRepresentation := []byte{0x2, 0x0, 0x0, 0x0, 0x30, 0x0}
	intZeroRepresentation := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	t.Run("should return error for unsupported bson type", func(t *testing.T) {
		var nilNonce *quote.Nonce
		zeroNonce := quote.NewNonce(0)
		require.ErrorIs(t, nilNonce.UnmarshalBSONValue(bson.TypeString, []byte{}), entities.DeserializationError)
		require.NoError(t, zeroNonce.UnmarshalBSONValue(bson.TypeString, stringZeroRepresentation))
		require.NoError(t, zeroNonce.UnmarshalBSONValue(bson.TypeInt64, intZeroRepresentation))
		test.RunTable(t, dataTypeCases, func(dataType bsontype.Type) error {
			return nilNonce.UnmarshalBSONValue(dataType, stringZeroRepresentation)
		})
	})
	t.Run("should unmarshal from bson.TypeInt64", func(t *testing.T) {
		int64Cases := test.Table[[]byte, *quote.Nonce]{
			{Value: intZeroRepresentation, Result: quote.NewNonce(0)},
			{Value: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: quote.NewNonce(1)},
			{Value: []byte{0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: quote.NewNonce(42)},
			{Value: []byte{0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: quote.NewNonce(99)},
			{Value: []byte{0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: quote.NewNonce(100)},
			{Value: []byte{0xf4, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: quote.NewNonce(500)},
			{Value: []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, Result: quote.NewNonce(math.MaxInt64 - 1)},
			{Value: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, Result: quote.NewNonce(math.MaxInt64)},
		}
		for _, tc := range int64Cases {
			result := new(quote.Nonce)
			err := result.UnmarshalBSONValue(bson.TypeInt64, tc.Value)
			require.NoError(t, err)
			assert.Equal(t, tc.Result, result)
		}
	})
	t.Run("should unmarshal from bson.TypeString", func(t *testing.T) {
		stringCases := test.Table[[]byte, *quote.Nonce]{
			{Value: stringZeroRepresentation, Result: quote.NewNonce(0)},
			{Value: []byte{0x02, 0x00, 0x00, 0x00, 0x31, 0x00}, Result: quote.NewNonce(1)},
			{Value: []byte{0x03, 0x00, 0x00, 0x00, 0x34, 0x32, 0x00}, Result: quote.NewNonce(42)},
			{Value: []byte{0x03, 0x00, 0x00, 0x00, 0x39, 0x39, 0x00}, Result: quote.NewNonce(99)},
			{Value: []byte{0x04, 0x00, 0x00, 0x00, 0x31, 0x30, 0x30, 0x00}, Result: quote.NewNonce(100)},
			{Value: []byte{0x04, 0x00, 0x00, 0x00, 0x35, 0x30, 0x30, 0x00}, Result: quote.NewNonce(500)},
			{
				Value:  []byte{0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37, 0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35, 0x38, 0x30, 0x36, 0x00},
				Result: quote.NewNonce(math.MaxInt64 - 1),
			},
			{
				Value:  []byte{0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37, 0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35, 0x38, 0x30, 0x37, 0x00},
				Result: quote.NewNonce(math.MaxInt64),
			},
		}
		for _, tc := range stringCases {
			result := new(quote.Nonce)
			err := result.UnmarshalBSONValue(bson.TypeString, tc.Value)
			require.NoError(t, err)
			assert.Equal(t, tc.Result, result)
		}
	})
}

func TestNonce_MarshalBSONValue(t *testing.T) {
	var bytes []byte
	var bsonTypeResult bsontype.Type
	var err error

	zeroRepresentation := []byte{0x2, 0x0, 0x0, 0x0, 0x30, 0x0}

	successCases := test.Table[*quote.Nonce, []byte]{
		{Value: quote.NewNonce(0), Result: zeroRepresentation},
		{Value: quote.NewNonce(1), Result: []byte{0x02, 0x00, 0x00, 0x00, 0x31, 0x00}},
		{Value: quote.NewNonce(77), Result: []byte{0x3, 0x0, 0x0, 0x0, 0x37, 0x37, 0x0}},
		{Value: quote.NewNonce(5678), Result: []byte{0x05, 0x00, 0x00, 0x00, 0x35, 0x36, 0x37, 0x38, 0x00}},
		{Value: quote.NewNonce(math.MaxInt64 - 1), Result: []byte{
			0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37,
			0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35,
			0x38, 0x30, 0x36, 0x00,
		}},
		{Value: quote.NewNonce(math.MaxInt64), Result: []byte{
			0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37,
			0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35,
			0x38, 0x30, 0x37, 0x00,
		}},
	}

	test.RunTable(t, successCases, func(value *quote.Nonce) []byte {
		bsonTypeResult, bytes, err = value.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bson.TypeString, bsonTypeResult)
		return bytes
	})
}

func TestNonce_String(t *testing.T) {
	var nilNonce *quote.Nonce
	zeroNonce := quote.NewNonce(0)
	oneNonce := quote.NewNonce(1)
	twoNonce := quote.NewNonce(2)
	maxNonce := quote.NewNonce(math.MaxInt64)
	negativeNonce := quote.NewNonce(-1)

	assert.Equal(t, "0", zeroNonce.String())
	assert.Equal(t, "1", oneNonce.String())
	assert.Equal(t, "2", twoNonce.String())
	assert.Equal(t, "-1", negativeNonce.String())
	assert.Equal(t, "9223372036854775807", maxNonce.String())
	assert.Equal(t, "0", nilNonce.String())
}

func TestNonce_Int64(t *testing.T) {
	tests := []struct {
		name     string
		nonce    *quote.Nonce
		expected int64
	}{
		{
			name:     "nil nonce returns zero",
			nonce:    nil,
			expected: 0,
		},
		{
			name:     "zero value",
			nonce:    quote.NewNonce(0),
			expected: 0,
		},
		{
			name:     "positive value",
			nonce:    quote.NewNonce(42),
			expected: 42,
		},
		{
			name:     "negative value",
			nonce:    quote.NewNonce(-123),
			expected: -123,
		},
		{
			name:     "max int64",
			nonce:    quote.NewNonce(math.MaxInt64),
			expected: math.MaxInt64,
		},
		{
			name:     "min int64",
			nonce:    quote.NewNonce(math.MinInt64),
			expected: math.MinInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nonce.Int64()
			assert.Equal(t, tt.expected, result)
		})
	}
}
