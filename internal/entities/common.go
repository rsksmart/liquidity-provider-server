package entities

import (
	"context"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
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
