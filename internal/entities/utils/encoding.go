package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

func (bf *BigFloat) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if bf == nil {
		return bson.MarshalValue(float64(0))
	}
	value, _ := bf.Native().Float64()
	return bson.MarshalValue(value)
}

func (bf *BigFloat) UnmarshalBSONValue(bsonType bsontype.Type, bytes []byte) error {
	if bf == nil || bsonType != bson.TypeDouble || len(bytes) == 0 {
		return entities.DeserializationError
	}
	var value float64
	if err := bson.UnmarshalValue(bsonType, bytes, &value); err != nil {
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
