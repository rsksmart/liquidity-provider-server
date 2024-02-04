package entities

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"math/big"
	"slices"
)

type Wei big.Int

type BigIntPtr = *big.Int

var bTen = big.NewInt(10)
var bEighteen = big.NewInt(18)
var bTenPowTen = new(big.Int).Exp(bTen, bTen, nil)           // 10**10
var bTenPowEighteen = new(big.Int).Exp(bTen, bEighteen, nil) // 10**18

func NewWei(x int64) *Wei {
	w := new(Wei)
	w.AsBigInt().SetInt64(x)
	return w
}

func NewUWei(x uint64) *Wei {
	w := new(Wei)
	w.AsBigInt().SetUint64(x)
	return w
}

func NewBigWei(x *big.Int) *Wei {
	w := new(Wei)
	w.AsBigInt().Set(x)
	return w
}

func SatoshiToWei(x uint64) *Wei {
	sat := new(big.Int).SetUint64(x)
	w := new(Wei)
	w.AsBigInt().Mul(sat, bTenPowTen)
	return w
}

func (w *Wei) Copy() *Wei {
	return NewBigWei(w.AsBigInt())
}

func (w *Wei) Cmp(y *Wei) int {
	return w.AsBigInt().Cmp(y.AsBigInt())
}

func (w *Wei) AsBigInt() BigIntPtr {
	return BigIntPtr(w)
}

func (w *Wei) Uint64() uint64 {
	return w.AsBigInt().Uint64()
}

func (w *Wei) ToRbtc() *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(w.AsBigInt()), new(big.Float).SetInt(bTenPowEighteen))
}

func (w *Wei) ToSatoshi() *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(w.AsBigInt()), new(big.Float).SetInt(bTenPowTen))
}

func (w *Wei) String() string {
	return w.AsBigInt().String()
}

func (w *Wei) Value() (driver.Value, error) {
	if w == nil {
		return "", errors.New("cannot retrieve value from <nil>")
	}
	return w.AsBigInt().String(), nil
}

func (w *Wei) Scan(src interface{}) error {
	switch src := src.(type) {
	case string:
		_, ok := w.AsBigInt().SetString(src, 10)
		if !ok {
			return errors.New("cannot scan invalid value")
		}
		return nil
	case nil:
		return errors.New("cannot scan <nil> value")
	default:
		return errors.New("cannot scan invalid type of value")
	}
}

func (w *Wei) MarshalJSON() ([]byte, error) {
	return w.AsBigInt().MarshalJSON()
}

func (w *Wei) UnmarshalJSON(bytes []byte) error {
	return w.AsBigInt().UnmarshalJSON(bytes)
}

func (w *Wei) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if w == nil {
		return bson.TypeInt64, make([]byte, 0), SerializationError
	}
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, w.Uint64())
	return bson.TypeInt64, value, nil
}

func (w *Wei) UnmarshalBSONValue(bsonType bsontype.Type, bytes []byte) error {
	if w == nil || bsonType != bson.TypeInt64 {
		return DeserializationError
	}
	slices.Reverse(bytes)
	w.AsBigInt().SetBytes(bytes)
	return nil
}

func (w *Wei) Add(x, y *Wei) *Wei {
	w.AsBigInt().Add(x.AsBigInt(), y.AsBigInt())
	return w
}

func (w *Wei) Sub(x, y *Wei) *Wei {
	w.AsBigInt().Sub(x.AsBigInt(), y.AsBigInt())
	return w
}

func (w *Wei) Mul(x, y *Wei) *Wei {
	w.AsBigInt().Mul(x.AsBigInt(), y.AsBigInt())
	return w
}
