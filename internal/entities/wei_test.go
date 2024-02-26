package entities_test

import (
	"database/sql/driver"
	entities "github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"math"
	"math/big"
	"reflect"
	"testing"
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
			want: new(big.Float).Quo(new(big.Float).SetInt64(1), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil))),
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
				val, ok := new(big.Int).SetString(tt.args.src.(string), 10)
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
		bytes, _ := i.MarshalJSON()
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
		bytes, _ := i.MarshalJSON()
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
		{Value: bson.TypeInt64, Result: entities.DeserializationError},
		{Value: bson.TypeString},
		{Value: bson.TypeDBPointer, Result: entities.DeserializationError},
		{Value: bson.TypeBinary, Result: entities.DeserializationError},
		{Value: bson.TypeDouble, Result: entities.DeserializationError},
	}

	zeroRepresentation := []byte{2, 0, 0, 0, 48, 0}
	successCases := test.Table[*entities.Wei, []byte]{
		{Value: entities.NewWei(0), Result: zeroRepresentation},
		{Value: entities.NewWei(5), Result: []byte{2, 0, 0, 0, 53, 0}},
		{Value: entities.NewWei(77), Result: []byte{3, 0, 0, 0, 55, 55, 0}},
		{Value: entities.NewWei(5678), Result: []byte{5, 0, 0, 0, 53, 54, 55, 56, 0}},
		{Value: entities.NewWei(math.MaxInt64 - 500), Result: []byte{20, 0, 0, 0, 57, 50, 50, 51, 51, 55, 50, 48, 51, 54, 56, 53, 52, 55, 55, 53, 51, 48, 55, 0}},
		{Value: entities.NewWei(math.MaxInt64), Result: []byte{20, 0, 0, 0, 57, 50, 50, 51, 51, 55, 50, 48, 51, 54, 56, 53, 52, 55, 55, 53, 56, 48, 55, 0}},
	}

	var nilWei *entities.Wei
	var bytes []byte
	var bsonTypeResult bsontype.Type
	weiValue := entities.NewWei(1)
	require.ErrorIs(t, nilWei.UnmarshalBSONValue(bson.TypeInt64, []byte{}), entities.DeserializationError)
	test.RunTable(t, dataTypeCases, func(bsonType bsontype.Type) error {
		return weiValue.UnmarshalBSONValue(bsonType, zeroRepresentation)
	})
	test.RunTable(t, successCases, func(value *entities.Wei) []byte {
		bsonTypeResult, bytes, _ = value.MarshalBSONValue()
		assert.Equal(t, bson.TypeString, bsonTypeResult)
		return bytes
	})

}
