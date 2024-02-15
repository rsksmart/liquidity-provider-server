package entities

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	DeserializationError = errors.New("error during value deserialization")
	SerializationError   = errors.New("error during value serialization")
	validate             = validator.New(validator.WithRequiredStructEnabled())
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
