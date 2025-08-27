package utils_test

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/require"
)

func TestDecodeKey(t *testing.T) {
	type params struct {
		key   string
		bytes int
	}
	cases := test.Table[params, []byte]{
		{
			Value:  params{key: "1234567890abcdef", bytes: 8},
			Result: []byte{18, 52, 86, 120, 144, 171, 205, 239},
		},
		{
			Value:  params{key: "a2fbac02d66202e8468d2a4f1deba4fa5c2491f592e0e22e32fe1e6acac25923", bytes: 32},
			Result: []byte{0xa2, 0xfb, 0xac, 0x2, 0xd6, 0x62, 0x2, 0xe8, 0x46, 0x8d, 0x2a, 0x4f, 0x1d, 0xeb, 0xa4, 0xfa, 0x5c, 0x24, 0x91, 0xf5, 0x92, 0xe0, 0xe2, 0x2e, 0x32, 0xfe, 0x1e, 0x6a, 0xca, 0xc2, 0x59, 0x23},
		},
		{
			Value:  params{key: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", bytes: 32},
			Result: []byte{0x9f, 0x86, 0xd0, 0x81, 0x88, 0x4c, 0x7d, 0x65, 0x9a, 0x2f, 0xea, 0xa0, 0xc5, 0x5a, 0xd0, 0x15, 0xa3, 0xbf, 0x4f, 0x1b, 0x2b, 0xb, 0x82, 0x2c, 0xd1, 0x5d, 0x6c, 0x15, 0xb0, 0xf0, 0xa, 0x8},
		},
		{
			Value:  params{key: "c5ff177a86e82441f93e3772da700d5f6838157fa1bfdc0bb689d7f7e55e7aba", bytes: 32},
			Result: []byte{0xc5, 0xff, 0x17, 0x7a, 0x86, 0xe8, 0x24, 0x41, 0xf9, 0x3e, 0x37, 0x72, 0xda, 0x70, 0xd, 0x5f, 0x68, 0x38, 0x15, 0x7f, 0xa1, 0xbf, 0xdc, 0xb, 0xb6, 0x89, 0xd7, 0xf7, 0xe5, 0x5e, 0x7a, 0xba},
		},
		{
			Value:  params{key: "ab5c2d1f", bytes: 4},
			Result: []byte{0xab, 0x5c, 0x2d, 0x1f},
		},
	}
	test.RunTable(t, cases, func(p params) []byte {
		result, err := utils.DecodeKey(p.key, p.bytes)
		require.NoError(t, err)
		return result
	})
}

func TestDecodeKey_SizeError(t *testing.T) {
	sizes := []int{4, 8, 12, 24, 32}
	key := "1122abcdff1122abcdff"
	for _, size := range sizes {
		result, err := utils.DecodeKey(key, size)
		require.Error(t, err)
		require.Nil(t, result)
		require.Errorf(t, err, "key length is not %d bytes, %s is %d bytes long", size, key, 10)
	}
}

func TestDecodeKey_DecodingError(t *testing.T) {
	cases := []string{
		"no hex",
		"abcde",
		"17",
		"g8ab11",
	}
	for _, key := range cases {
		result, err := utils.DecodeKey(key, 8)
		require.Error(t, err)
		require.Nil(t, result)
	}
}

func TestTo32Bytes(t *testing.T) {
	shortSlice := []byte{0x1, 0x2, 0x3, 0x4}
	expectedShortResult := [32]byte{0x1, 0x2, 0x3, 0x4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	longSlice := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22}
	expectedLongResult := [32]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

	exactSlice := []byte{0x20, 0x1f, 0x1e, 0x1d, 0x1c, 0x1b, 0x1a, 0x19, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0xf, 0xe, 0xd, 0xc, 0xb, 0xa, 0x9, 0x8, 0x7, 0x6, 0x5, 0x4, 0x3, 0x2, 0x1}
	expectedExactResult := [32]byte{0x20, 0x1f, 0x1e, 0x1d, 0x1c, 0x1b, 0x1a, 0x19, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0xf, 0xe, 0xd, 0xc, 0xb, 0xa, 0x9, 0x8, 0x7, 0x6, 0x5, 0x4, 0x3, 0x2, 0x1}

	shortResult := utils.To32Bytes(shortSlice)
	longResult := utils.To32Bytes(longSlice)
	exactResult := utils.To32Bytes(exactSlice)

	require.Equal(t, expectedShortResult, shortResult)
	require.Equal(t, expectedLongResult, longResult)
	require.Equal(t, expectedExactResult, exactResult)
}

func TestDecodeKey_LengthErrorDoesNotExposeKey(t *testing.T) {
	// A valid hex string but wrong length
	sensitiveKey := "1234567890abcdef1234567890abcdef" // 16 bytes
	expectedBytes := 32                                // Expecting 32 bytes

	result, err := utils.DecodeKey(sensitiveKey, expectedBytes)

	require.Error(t, err)
	require.Nil(t, result)

	// Check that error message contains expected information but not the key
	require.Contains(t, err.Error(), "key length is not 32 bytes")
	require.Contains(t, err.Error(), "16 bytes long")
	require.NotContains(t, err.Error(), sensitiveKey)
}

func TestNewBigFloat(t *testing.T) {
	type args struct {
		x *big.Float
	}
	tests := []struct {
		name string
		args args
		want *utils.BigFloat
	}{
		{
			name: "new BigFloat",
			args: args{x: big.NewFloat(1.1234554321)},
			want: utils.NewBigFloat64(1.1234554321),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.NewBigFloat(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigFloat64(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want *utils.BigFloat
	}{
		{
			name: "new zero BigFloat",
			args: args{x: 0},
			want: utils.NewBigFloat(new(big.Float).SetPrec(53)),
		},
		{
			name: "new BigFloat",
			args: args{x: 1.55553333111},
			want: (*utils.BigFloat)(big.NewFloat(1.55553333111)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.NewBigFloat64(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigFloat_StringFormatsFixedPoint(t *testing.T) {
	bf := utils.NewBigFloat64(0.000012)
	expected := "0.000012"
	require.Equal(t, expected, bf.String())
}

func TestStructFormatting_UsesFixedPoint(t *testing.T) {
	type sample struct {
		FeeRate *utils.BigFloat
	}
	payload := sample{FeeRate: utils.NewBigFloat64(0.000012)}
	formatted := fmt.Sprintf("%+v", payload)
	require.Contains(t, formatted, "FeeRate:0.000012")
}
func TestBigFloat_Native(t *testing.T) {
	tests := []struct {
		name string
		w    *utils.BigFloat
		want *big.Float
	}{
		{
			name: "as big.Float",
			w:    utils.NewBigFloat64(123.45567889),
			want: big.NewFloat(123.45567889),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Native(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigFloat_UnmarshalBSONValue(t *testing.T) {
	dataTypeCases := test.Table[bsontype.Type, error]{
		{Value: bson.TypeInt64, Result: entities.DeserializationError},
		{Value: bson.TypeString, Result: entities.DeserializationError},
		{Value: bson.TypeDBPointer, Result: entities.DeserializationError},
		{Value: bson.TypeBinary, Result: entities.DeserializationError},
		{Value: bson.TypeDouble},
	}

	zeroRepresentation := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	successCases := test.Table[*utils.BigFloat, []byte]{
		{Value: utils.NewBigFloat64(0), Result: zeroRepresentation},
		{Value: utils.NewBigFloat64(5.3333), Result: []byte{0xf7, 0x6, 0x5f, 0x98, 0x4c, 0x55, 0x15, 0x40}},
		{Value: utils.NewBigFloat64(77), Result: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x53, 0x40}},
		{Value: utils.NewBigFloat64(5678.51251), Result: []byte{0xdf, 0xf8, 0xda, 0x33, 0x83, 0x2e, 0xb6, 0x40}},
		{Value: utils.NewBigFloat64(math.MaxFloat64 - 500.1235), Result: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}},
		{Value: utils.NewBigFloat64(math.MaxFloat64), Result: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xef, 0x7f}},
	}

	var nilBigFloat *utils.BigFloat
	var bytes []byte
	var bsonTypeResult bsontype.Type
	var err error
	bigFloatValue := utils.NewBigFloat64(1.12351251)
	require.ErrorIs(t, nilBigFloat.UnmarshalBSONValue(bson.TypeString, []byte{}), entities.DeserializationError)
	test.RunTable(t, dataTypeCases, func(bsonType bsontype.Type) error {
		return bigFloatValue.UnmarshalBSONValue(bsonType, zeroRepresentation)
	})
	test.RunTable(t, successCases, func(value *utils.BigFloat) []byte {
		bsonTypeResult, bytes, err = value.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bson.TypeDouble, bsonTypeResult)
		return bytes
	})
}

func TestBigFloat_MarshalJSON_Nil(t *testing.T) {
	var bf *utils.BigFloat
	data, err := bf.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte("null"), data)
}

var bigFloatMarshalJSONTests = []struct {
	name      string
	bf        *utils.BigFloat
	expected  string
	expectErr bool
}{
	{
		name:      "Zero BigFloat",
		bf:        utils.NewBigFloat64(0),
		expected:  "0",
		expectErr: false,
	},
	{
		name:      "Positive BigFloat",
		bf:        utils.NewBigFloat64(123.456),
		expected:  "123.456",
		expectErr: false,
	},
	{
		name:      "Negative BigFloat",
		bf:        utils.NewBigFloat64(-123.456),
		expected:  "-123.456",
		expectErr: false,
	},
	{
		name:      "Large BigFloat",
		bf:        utils.NewBigFloat(big.NewFloat(1e+20)),
		expected:  "1e+20",
		expectErr: false,
	},
	{
		name:      "Small BigFloat",
		bf:        utils.NewBigFloat(big.NewFloat(1e-20)),
		expected:  "1e-20",
		expectErr: false,
	},
	{
		name:      "Scientific notation",
		bf:        utils.NewBigFloat(big.NewFloat(1.2345e+10)),
		expected:  "1.2345e+10",
		expectErr: false,
	},
	{
		name:      "Negative scientific notation",
		bf:        utils.NewBigFloat(big.NewFloat(-1.2345e+10)),
		expected:  "-1.2345e+10",
		expectErr: false,
	},
	{
		name:      "Maximum float64",
		bf:        utils.NewBigFloat(big.NewFloat(math.MaxFloat64)),
		expected:  "1.7976931348623157e+308",
		expectErr: false,
	},
	{
		name:      "Minimum positive float64",
		bf:        utils.NewBigFloat(big.NewFloat(math.SmallestNonzeroFloat64)),
		expected:  "5e-324",
		expectErr: false,
	},
}

func TestBigFloat_MarshalJSON(t *testing.T) {
	for _, tt := range bigFloatMarshalJSONTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.bf.MarshalJSON()
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.JSONEq(t, tt.expected, string(data))
			}
		})
	}
}

var bigFloatUnmarshalJSONValidInputs = []struct {
	name     string
	input    string
	expected *big.Float
}{
	{
		name:     "Valid positive number",
		input:    "123.456",
		expected: big.NewFloat(123.456),
	},
	{
		name:     "Valid negative number",
		input:    "-123.456",
		expected: big.NewFloat(-123.456),
	},
	{
		name:     "Zero",
		input:    "0",
		expected: big.NewFloat(0),
	},
	{
		name:     "Scientific notation",
		input:    "1.2345e+10",
		expected: big.NewFloat(1.2345e+10),
	},
	{
		name:     "Negative scientific notation",
		input:    "-1.2345e+10",
		expected: big.NewFloat(-1.2345e+10),
	},
	{
		name:     "Whitespace input",
		input:    "   123.456   ",
		expected: big.NewFloat(123.456),
	},
	{
		name:     "Maximum float64",
		input:    "1.7976931348623157e+308",
		expected: big.NewFloat(math.MaxFloat64),
	},
	{
		name:     "Minimum positive float64",
		input:    "5e-324",
		expected: big.NewFloat(math.SmallestNonzeroFloat64),
	},
	{
		name:     "Underflow number",
		input:    "1e-400",
		expected: big.NewFloat(0.0),
	},
	{
		name:     "Negative underflow number",
		input:    "-1e-400",
		expected: big.NewFloat(math.Copysign(0, -1)),
	},
}

func TestBigFloat_UnmarshalJSON_ValidInputs(t *testing.T) {
	for _, tt := range bigFloatUnmarshalJSONValidInputs {
		t.Run(tt.name, func(t *testing.T) {
			var bf utils.BigFloat
			err := bf.UnmarshalJSON([]byte(tt.input))
			require.NoError(t, err)

			expectedFloat64, _ := tt.expected.Float64()
			actualFloat64, _ := bf.Native().Float64()
			require.InDelta(t, expectedFloat64, actualFloat64, 1e-10)
		})
	}
}

func TestBigFloat_UnmarshalJSON_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid string",
			input: "\"invalid\"",
		},
		{
			name:  "Empty input",
			input: "",
		},
		{
			name:  "Invalid JSON",
			input: "{",
		},
		{
			name:  "NaN input",
			input: "\"NaN\"",
		},
		{
			name:  "Infinity input",
			input: "\"Infinity\"",
		},
		{
			name:  "Array input",
			input: "[123.456]",
		},
		{
			name:  "Object input",
			input: "{\"value\":123.456}",
		},
		{
			name:  "Boolean true input",
			input: "true",
		},
		{
			name:  "Boolean false input",
			input: "false",
		},
		{
			name:  "Overflow number",
			input: "1e+400",
		},
		{
			name:  "Negative overflow number",
			input: "-1e+400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf utils.BigFloat
			err := bf.UnmarshalJSON([]byte(tt.input))
			require.Error(t, err)
		})
	}
}

func TestBigFloat_MarshalJSON_Basic(t *testing.T) {
	tests := []struct {
		name     string
		bf       *utils.BigFloat
		expected string
	}{
		{
			name:     "Zero BigFloat",
			bf:       utils.NewBigFloat64(0),
			expected: "0",
		},
		{
			name:     "Positive BigFloat",
			bf:       utils.NewBigFloat64(123.456),
			expected: "123.456",
		},
		{
			name:     "Negative BigFloat",
			bf:       utils.NewBigFloat64(-123.456),
			expected: "-123.456",
		},
		{
			name:     "Small BigFloat",
			bf:       utils.NewBigFloat(big.NewFloat(1e-20)),
			expected: "1e-20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.bf.MarshalJSON()
			require.NoError(t, err)
			require.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestBigFloat_MarshalJSON_Extremes(t *testing.T) {
	tests := []struct {
		name     string
		bf       *utils.BigFloat
		expected string
	}{
		{
			name:     "Large BigFloat",
			bf:       utils.NewBigFloat(big.NewFloat(1e+20)),
			expected: "1e+20",
		},
		{
			name:     "Scientific notation",
			bf:       utils.NewBigFloat(big.NewFloat(1.2345e+10)),
			expected: "1.2345e+10",
		},
		{
			name:     "Negative scientific notation",
			bf:       utils.NewBigFloat(big.NewFloat(-1.2345e+10)),
			expected: "-1.2345e+10",
		},
		{
			name:     "Maximum float64",
			bf:       utils.NewBigFloat(big.NewFloat(math.MaxFloat64)),
			expected: "1.7976931348623157e+308",
		},
		{
			name:     "Minimum positive float64",
			bf:       utils.NewBigFloat(big.NewFloat(math.SmallestNonzeroFloat64)),
			expected: "5e-324",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.bf.MarshalJSON()
			require.NoError(t, err)
			require.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestBigFloat_UnmarshalJSON_ValidInputs_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *big.Float
	}{
		{
			name:     "Valid positive number",
			input:    "123.456",
			expected: big.NewFloat(123.456),
		},
		{
			name:     "Valid negative number",
			input:    "-123.456",
			expected: big.NewFloat(-123.456),
		},
		{
			name:     "Zero",
			input:    "0",
			expected: big.NewFloat(0),
		},
		{
			name:     "Whitespace input",
			input:    "   123.456   ",
			expected: big.NewFloat(123.456),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf utils.BigFloat
			err := bf.UnmarshalJSON([]byte(tt.input))
			require.NoError(t, err)

			expectedFloat64, _ := tt.expected.Float64()
			actualFloat64, _ := bf.Native().Float64()
			require.InDelta(t, expectedFloat64, actualFloat64, 1e-10)
		})
	}
}

func TestBigFloat_UnmarshalJSON_ValidInputs_Extremes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *big.Float
	}{
		{
			name:     "Scientific notation",
			input:    "1.2345e+10",
			expected: big.NewFloat(1.2345e+10),
		},
		{
			name:     "Negative scientific notation",
			input:    "-1.2345e+10",
			expected: big.NewFloat(-1.2345e+10),
		},
		{
			name:     "Maximum float64",
			input:    "1.7976931348623157e+308",
			expected: big.NewFloat(math.MaxFloat64),
		},
		{
			name:     "Minimum positive float64",
			input:    "5e-324",
			expected: big.NewFloat(math.SmallestNonzeroFloat64),
		},
		{
			name:     "Underflow number",
			input:    "1e-400",
			expected: big.NewFloat(0.0),
		},
		{
			name:     "Negative underflow number",
			input:    "-1e-400",
			expected: big.NewFloat(math.Copysign(0, -1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf utils.BigFloat
			err := bf.UnmarshalJSON([]byte(tt.input))
			require.NoError(t, err)

			expectedFloat64, _ := tt.expected.Float64()
			actualFloat64, _ := bf.Native().Float64()
			require.InDelta(t, expectedFloat64, actualFloat64, 1e-10)
		})
	}
}

func TestBigFloat_UnmarshalJSON_NullInput(t *testing.T) {
	var bf utils.BigFloat
	err := bf.UnmarshalJSON([]byte("null"))
	require.NoError(t, err)

	expected := big.NewFloat(0)
	expectedFloat64, _ := expected.Float64()
	actualFloat64, _ := bf.Native().Float64()
	require.InDelta(t, expectedFloat64, actualFloat64, 1e-10)
}

func TestCompareIgnore0x(t *testing.T) {
	cases := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "Equal with 0x prefix and without",
			a:        "0x123abc",
			b:        "123abc",
			expected: true,
		},
		{
			name:     "Different values",
			a:        "0x123abc",
			b:        "0x456def",
			expected: false,
		},
		{
			name:     "Both empty strings",
			a:        "",
			b:        "",
			expected: true,
		},
		{
			name:     "Case insensitive match",
			a:        "0xABCDEF",
			b:        "abcdef",
			expected: true,
		},
		{
			name:     "No prefix, equal",
			a:        "123abc",
			b:        "123abc",
			expected: true,
		},
		{
			name:     "Completely different",
			a:        "0x123",
			b:        "0x456",
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CompareIgnore0x(tt.a, tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}
