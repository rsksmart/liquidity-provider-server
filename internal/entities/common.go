package entities

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	DeserializationError = errors.New("error during value deserialization")
	SerializationError   = errors.New("error during value serialization")
	IntegrityError       = errors.New("error during value integrity check, stored hash doesn't match actual hash")
	validate             = validator.New(validator.WithRequiredStructEnabled())
	DivideByZeroError    = errors.New("divide by zero error")
)

type NodeType = string

const (
	NodeTypeRootstock NodeType = "rootstock"
	NodeTypeBitcoin   NodeType = "bitcoin"
)

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

type Closeable interface {
	Shutdown(closeChannel chan<- bool)
}

type Service interface {
	CheckConnection(ctx context.Context) bool
}

type HashFunction func(...[]byte) []byte

type Signer interface {
	SignBytes(msg []byte) ([]byte, error)
	Validate(signature, hash string) bool
}

type Signed[T any] struct {
	Value     T      `bson:",inline"`
	Signature string `json:"signature" bson:"signature"`
	Hash      string `json:"hash" bson:"hash"`
}

func (signedValue Signed[T]) CheckIntegrity(hashFunction HashFunction) error {
	valueBytes, err := json.Marshal(signedValue.Value)
	if err != nil {
		return err
	}
	hash := hashFunction(valueBytes)
	storedHash, err := hex.DecodeString(signedValue.Hash)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, storedHash) {
		return IntegrityError
	}
	return nil
}
