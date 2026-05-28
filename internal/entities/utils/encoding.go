package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BigFloat big.Float

func (bf *BigFloat) Native() *big.Float {
	return (*big.Float)(bf)
}

func NewBigFloat64(value float64) *BigFloat {
	return (*BigFloat)(big.NewFloat(value))
}

func NewBigFloat(value *big.Float) *BigFloat {
	return (*BigFloat)(value)
}

func (bf *BigFloat) MarshalJSON() ([]byte, error) {
	if bf == nil {
		return []byte("null"), nil
	}
	value, _ := bf.Native().Float64()
	return json.Marshal(value)
}

func (bf *BigFloat) UnmarshalJSON(b []byte) error {
	var value float64
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	bf.Native().SetFloat64(value)
	return nil
}

func (bf *BigFloat) MarshalBSONValue() (byte, []byte, error) {
	if bf == nil {
		t, data, err := bson.MarshalValue(float64(0))
		return byte(t), data, err
	}
	value, _ := bf.Native().Float64()
	t, data, err := bson.MarshalValue(value)
	return byte(t), data, err
}

func (bf *BigFloat) UnmarshalBSONValue(bsonType byte, bytes []byte) error {
	typ := bson.Type(bsonType)
	if bf == nil || typ != bson.TypeDouble || len(bytes) == 0 {
		return entities.DeserializationError
	}
	var value float64
	if err := bson.UnmarshalValue(typ, bytes, &value); err != nil {
		return errors.Join(entities.DeserializationError, err)
	}
	result := big.NewFloat(value)
	bf.Native().Set(result)
	return nil
}

func (bf *BigFloat) String() string {
	return bf.Native().Text('f', -1)
}

func DecodeKey(key string, expectedBytes int) ([]byte, error) {
	var err error
	var bytes []byte
	if bytes, err = hex.DecodeString(key); err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}
	if len(bytes) != expectedBytes {
		return nil, fmt.Errorf("key length is not %d bytes, key is %d bytes long", expectedBytes, len(bytes))
	}
	return bytes, nil
}

// To32Bytes utility to convert a byte slice to a fixed size byte array, if input has
// more than 32 bytes they won't be copied.
func To32Bytes(value []byte) [32]byte {
	var bytes [32]byte
	copy(bytes[:], value)
	return bytes
}

func CompareIgnore0x(a, b string) bool {
	const prefix = "0x"
	return strings.EqualFold(strings.TrimPrefix(a, prefix), strings.TrimPrefix(b, prefix))
}

func Prepend0x(s string) string {
	const prefix = "0x"
	if strings.HasPrefix(strings.ToLower(s), prefix) {
		return s
	}
	return prefix + s
}
