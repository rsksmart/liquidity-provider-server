package entities_test

import (
	"database/sql/driver"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

func TestSatoshiToWei(t *testing.T) {
	type args struct {
		x uint64
	}
	tests := []struct {
		name string
		args args
		want *entities.Wei
	}{
		{
			name: "zero sat to wei",
			args: args{x: 0},
			want: entities.NewWei(0),
		},
		{
			name: "one sat to wei",
			args: args{x: 1},
			want: entities.NewWei(int64(math.Pow(10, 10))),
		},
		{
			name: "10**8 sat (1 btc) to wei",
			args: args{x: uint64(math.Pow(10, 8))},
			want: entities.NewWei(int64(math.Pow(10, 18))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.SatoshiToWei(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SatoshiToWei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBigWei(t *testing.T) {
	type args struct {
		x *big.Int
	}
	tests := []struct {
		name string
		args args
		want *entities.Wei
	}{
		{
			name: "new big wei",
			args: args{x: big.NewInt(1)},
			want: entities.NewWei(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.NewBigWei(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBigWei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWei(t *testing.T) {
	type args struct {
		x int64
	}
	tests := []struct {
		name string
		args args
		want *entities.Wei
	}{
		{
			name: "new zero wei",
			args: args{x: 0},
			want: new(entities.Wei),
		},
		{
			name: "new one wei",
			args: args{x: 1},
			want: (*entities.Wei)(big.NewInt(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entities.NewWei(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWei() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_AsBigInt(t *testing.T) {
	tests := []struct {
		name string
		w    *entities.Wei
		want entities.BigIntPtr
	}{
		{
			name: "as big.int",
			w:    entities.NewWei(1),
			want: big.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.AsBigInt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_ToRbtc(t *testing.T) {
	tests := []struct {
		name string
		w    *entities.Wei
		want *big.Float
	}{
		{
			name: "1 wei to rbtc",
			w:    entities.NewWei(1),
			want: new(big.Float).Quo(new(big.Float).SetInt64(1), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))),
		},
		{
			name: "2*(10**10) wei to rbtc",
			w:    entities.NewWei(int64(2 * math.Pow(10, 18))),
			want: big.NewFloat(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.ToRbtc(); got.Cmp(tt.want) != 0 {
				t.Errorf("ToRbtc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_ToSatoshi(t *testing.T) {
	tests := []struct {
		name string
		w    *entities.Wei
		want *big.Float
	}{
		{
			name: "zero wei to sat",
			w:    entities.NewWei(0),
			want: big.NewFloat(0),
		},
		{
			name: "1 wei to sat",
			w:    entities.NewWei(1),
			want: big.NewFloat(1),
		},
		{
			name: "72160329123080000 wei to 7216033 sat",
			w:    entities.NewWei(72160329123080000),
			want: big.NewFloat(7216033),
		},
		{
			name: "4360000000000000 wei to 436000 sat",
			w:    entities.NewWei(4360000000000000),
			want: big.NewFloat(436000),
		},
		{
			name: "1 RBTC to 100000000 sat",
			w:    entities.NewWei(1000000000000000000),
			want: big.NewFloat(100000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.ToSatoshi(); got.Cmp(tt.want) != 0 {
				t.Errorf("ToSatoshi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_Uint64(t *testing.T) {
	tests := []struct {
		name string
		w    *entities.Wei
		want uint64
	}{
		{
			name: "wei to uint64",
			w:    entities.NewWei(1),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Uint64(); got != tt.want {
				t.Errorf("Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_Value(t *testing.T) {
	tests := []struct {
		name    string
		w       *entities.Wei
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "wei value",
			w:       entities.NewWei(1),
			want:    "1",
			wantErr: false,
		},
		{
			name:    "<nil> wei value",
			w:       nil,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		w       *entities.Wei
		args    args
		wantErr bool
	}{
		{
			name:    "valid value",
			w:       new(entities.Wei),
			args:    args{src: "100"},
			wantErr: false,
		},
		{
			name:    "valid big value",
			w:       new(entities.Wei),
			args:    args{src: new(big.Int).Mul(new(big.Int).SetUint64(math.MaxUint64), big.NewInt(10)).String()}, // 10 * math.MaxUint64
			wantErr: false,
		},
		{
			name:    "<nil> value",
			w:       new(entities.Wei),
			args:    args{src: nil},
			wantErr: true,
		},
		{
			name:    "invalid value",
			w:       new(entities.Wei),
			args:    args{src: "abc"},
			wantErr: true,
		},
		{
			name:    "invalid type",
			w:       new(entities.Wei),
			args:    args{src: true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				stringArg, ok := tt.args.src.(string)
				require.True(t, ok)
				val, ok := new(big.Int).SetString(stringArg, 10)
				if !ok {
					t.Fatal("invalid arg")
				}
				if val.Cmp(tt.w.AsBigInt()) != 0 {
					t.Errorf("Scan() = %v, want %v", tt.w, val)
				}
			}
		})
	}
}

func TestWei_String(t *testing.T) {
	tests := []struct {
		name string
		w    *entities.Wei
		want string
	}{
		{
			name: "wei to string",
			w:    entities.NewWei(100),
			want: "100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_Copy(t *testing.T) {
	w := entities.NewWei(100)
	tests := []struct {
		name string
		w    *entities.Wei
		want *entities.Wei
	}{
		{
			name: "copy wei",
			w:    w,
			want: w,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Copy(); tt.w == got || got.AsBigInt().Cmp(tt.want.AsBigInt()) != 0 {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_Cmp(t *testing.T) {
	type args struct {
		y *entities.Wei
	}
	tests := []struct {
		name  string
		x     *entities.Wei
		args  args
		wantR int
	}{
		{
			name:  "eq wei",
			x:     entities.NewWei(2),
			args:  args{y: entities.NewWei(2)},
			wantR: 0,
		},
		{
			name:  "gt wei",
			x:     entities.NewWei(2),
			args:  args{y: entities.NewWei(1)},
			wantR: 1,
		},
		{
			name:  "lt wei",
			x:     entities.NewWei(1),
			args:  args{y: entities.NewWei(2)},
			wantR: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := tt.x.Cmp(tt.args.y); gotR != tt.wantR {
				t.Errorf("Cmp() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestWei_MarshalJSON(t *testing.T) {
	bigIntToBytes := func(i *big.Int) []byte {
		bytes, err := i.MarshalJSON()
		require.NoError(t, err)
		return bytes
	}
	tests := []struct {
		name    string
		w       *entities.Wei
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal wei",
			w:       entities.NewWei(100),
			want:    bigIntToBytes(big.NewInt(100)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWei_UnmarshalJSON(t *testing.T) {
	bigIntToBytes := func(i *big.Int) []byte {
		bytes, err := i.MarshalJSON()
		require.NoError(t, err)
		return bytes
	}
	type args struct {
		val   *big.Int
		bytes []byte
	}
	tests := []struct {
		name    string
		w       *entities.Wei
		args    args
		wantErr bool
	}{
		{
			name:    "unmarshal wei",
			w:       new(entities.Wei),
			args:    args{val: big.NewInt(100), bytes: bigIntToBytes(big.NewInt(100))},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.w.AsBigInt().Cmp(tt.args.val) != 0 {
				t.Errorf("tt.w = %v, want %v", tt.w, tt.args.val)
			}
		})
	}
}

func TestWei_UnmarshalBSONValue(t *testing.T) {
	dataTypeCases := test.Table[bsontype.Type, error]{
		{Value: bson.TypeDBPointer, Result: entities.DeserializationError},
		{Value: bson.TypeBinary, Result: entities.DeserializationError},
		{Value: bson.TypeDouble, Result: entities.DeserializationError},
	}

	zeroInt64 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	zeroString := []byte{0x2, 0x0, 0x0, 0x0, 0x30, 0x0}

	t.Run("should return error for unsupported bson type", func(t *testing.T) {
		var nilWei *entities.Wei
		zeroWei := entities.NewWei(0)
		require.ErrorIs(t, nilWei.UnmarshalBSONValue(bson.TypeInt64, []byte{}), entities.DeserializationError)
		require.ErrorIs(t, nilWei.UnmarshalBSONValue(bson.TypeString, []byte{}), entities.DeserializationError)
		require.NoError(t, zeroWei.UnmarshalBSONValue(bson.TypeString, zeroString))
		require.NoError(t, zeroWei.UnmarshalBSONValue(bson.TypeInt64, zeroInt64))
		test.RunTable(t, dataTypeCases, func(dataType bsontype.Type) error {
			return nilWei.UnmarshalBSONValue(dataType, zeroInt64)
		})
	})
	t.Run("should handle null values gracefully", func(t *testing.T) {
		zeroWei := entities.NewWei(0)
		// When BSON contains a null value, UnmarshalBSONValue should succeed
		// This allows the Go MongoDB driver to set the field to nil
		require.NoError(t, zeroWei.UnmarshalBSONValue(bson.TypeNull, []byte{}))
	})
	t.Run("should handle '<nil>' string values gracefully", func(t *testing.T) {
		zeroWei := entities.NewWei(0)
		// When MongoDB stores nil pointers as the string "<nil>", we should handle it gracefully
		nilStringBytes := []byte{0x06, 0x00, 0x00, 0x00, 0x3c, 0x6e, 0x69, 0x6c, 0x3e, 0x00} // BSON string "<nil>"
		require.NoError(t, zeroWei.UnmarshalBSONValue(bson.TypeString, nilStringBytes))

		// Test with nil Wei pointer should still error
		var nilWei *entities.Wei
		require.ErrorIs(t, nilWei.UnmarshalBSONValue(bson.TypeString, nilStringBytes), entities.DeserializationError)
	})
}

func TestWei_MarshalBSONValue(t *testing.T) {
	t.Run("should marshal nil Wei as BSON null", func(t *testing.T) {
		var nilWei *entities.Wei
		bsonType, bytes, err := nilWei.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bson.TypeNull, bsonType)
		assert.Empty(t, bytes)
	})

	t.Run("should marshal non-nil Wei as string", func(t *testing.T) {
		wei := entities.NewWei(12345)
		bsonType, bytes, err := wei.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bson.TypeString, bsonType)
		assert.NotEmpty(t, bytes)

		// Verify we can unmarshal it back
		var result string
		err = bson.UnmarshalValue(bsonType, bytes, &result)
		require.NoError(t, err)
		assert.Equal(t, "12345", result)
	})
}

func TestWei_UnmarshalBSONValue_Integration(t *testing.T) {
	zeroInt64 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	zeroString := []byte{0x2, 0x0, 0x0, 0x0, 0x30, 0x0}

	t.Run("should unmarshal from bson.TypeInt64", func(t *testing.T) {
		int64Cases := test.Table[[]byte, *entities.Wei]{
			{Value: zeroInt64, Result: entities.NewWei(0)},
			{Value: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: entities.NewWei(1)},
			{Value: []byte{0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: entities.NewWei(42)},
			{Value: []byte{0x63, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: entities.NewWei(99)},
			{Value: []byte{0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: entities.NewWei(100)},
			{Value: []byte{0xf4, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, Result: entities.NewWei(500)},
			{Value: []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, Result: entities.NewWei(math.MaxInt64 - 1)},
			{Value: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, Result: entities.NewWei(math.MaxInt64)},
		}
		for _, tc := range int64Cases {
			result := new(entities.Wei)
			err := result.UnmarshalBSONValue(bson.TypeInt64, tc.Value)
			require.NoError(t, err)
			assert.Equal(t, tc.Result, result)
		}
	})
	t.Run("should unmarshal from bson.TypeString", func(t *testing.T) {
		stringCases := test.Table[[]byte, *entities.Wei]{
			{Value: zeroString, Result: entities.NewWei(0)},
			{Value: []byte{0x02, 0x00, 0x00, 0x00, 0x31, 0x00}, Result: entities.NewWei(1)},
			{Value: []byte{2, 0, 0, 0, 53, 0}, Result: entities.NewWei(5)},
			{Value: []byte{0x03, 0x00, 0x00, 0x00, 0x34, 0x32, 0x00}, Result: entities.NewWei(42)},
			{Value: []byte{3, 0, 0, 0, 55, 55, 0}, Result: entities.NewWei(77)},
			{Value: []byte{0x03, 0x00, 0x00, 0x00, 0x39, 0x39, 0x00}, Result: entities.NewWei(99)},
			{Value: []byte{0x04, 0x00, 0x00, 0x00, 0x31, 0x30, 0x30, 0x00}, Result: entities.NewWei(100)},
			{Value: []byte{0x04, 0x00, 0x00, 0x00, 0x35, 0x30, 0x30, 0x00}, Result: entities.NewWei(500)},
			{Value: []byte{5, 0, 0, 0, 53, 54, 55, 56, 0}, Result: entities.NewWei(5678)},
			{Value: []byte{20, 0, 0, 0, 57, 50, 50, 51, 51, 55, 50, 48, 51, 54, 56, 53, 52, 55, 55, 53, 51, 48, 55, 0}, Result: entities.NewWei(math.MaxInt64 - 500)},
			{Value: []byte{0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37, 0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35, 0x38, 0x30, 0x36, 0x00}, Result: entities.NewWei(math.MaxInt64 - 1)},
			{Value: []byte{0x14, 0x00, 0x00, 0x00, 0x39, 0x32, 0x32, 0x33, 0x33, 0x37, 0x32, 0x30, 0x33, 0x36, 0x38, 0x35, 0x34, 0x37, 0x37, 0x35, 0x38, 0x30, 0x37, 0x00}, Result: entities.NewWei(math.MaxInt64)},
		}
		for _, tc := range stringCases {
			result := new(entities.Wei)
			err := result.UnmarshalBSONValue(bson.TypeString, tc.Value)
			require.NoError(t, err)
			assert.Equal(t, tc.Result, result)
		}
	})
}

func TestWei_Div(t *testing.T) {
	tests := []struct {
		name    string
		w       *entities.Wei
		x       *entities.Wei
		y       *entities.Wei
		want    *entities.Wei
		wantErr bool
	}{
		{
			name:    "divide two positive Wei values",
			w:       entities.NewWei(0),
			x:       entities.NewWei(10),
			y:       entities.NewWei(5),
			want:    entities.NewWei(2),
			wantErr: false,
		},
		{
			name:    "divide by zero",
			w:       entities.NewWei(0),
			x:       entities.NewWei(10),
			y:       entities.NewWei(0),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "divide zero by a number",
			w:       entities.NewWei(0),
			x:       entities.NewWei(0),
			y:       entities.NewWei(5),
			want:    entities.NewWei(0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Div(tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Div() got = %v, want %v", got, tt.want)
			}
		})
	}
}
