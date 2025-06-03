package quote

import (
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"math/big"
)

type Nonce big.Int

func NewNonce(value int64) *Nonce {
	n := new(Nonce)
	n.AsBigInt().SetInt64(value)
	return n
}

func (n *Nonce) String() string {
	if n == nil {
		return "0"
	}
	return (*big.Int)(n).String()
}

func (n *Nonce) AsBigInt() *big.Int {
	return (*big.Int)(n)
}

func (n *Nonce) Int64() int64 {
	if n == nil || n.AsBigInt() == nil {
		return 0
	}
	return n.AsBigInt().Int64()
}

func (n *Nonce) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(n.AsBigInt().String())
}

func (n *Nonce) UnmarshalBSONValue(bsonType bsontype.Type, bytes []byte) error {
	supportedType := bsonType == bson.TypeInt64 || bsonType == bson.TypeString
	if n == nil || !supportedType || len(bytes) == 0 {
		return entities.DeserializationError
	}

	// we are supporting int64 to be able to handle legacy quotes as well
	if bsonType == bson.TypeInt64 {
		var value int64
		if err := bson.UnmarshalValue(bsonType, bytes, &value); err != nil {
			return errors.Join(entities.DeserializationError, err)
		}
		n.AsBigInt().SetInt64(value)
		return nil
	}

	var value string
	if err := bson.UnmarshalValue(bsonType, bytes, &value); err != nil {
		return errors.Join(entities.DeserializationError, err)
	}
	_, ok := n.AsBigInt().SetString(value, 10)

	if !ok {
		return fmt.Errorf("%w: cannot unmarshal value %s to Nonce", entities.DeserializationError, value)
	}
	return nil
}
