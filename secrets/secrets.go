package secrets

import (
	"time"
)

const (
	WriteTimeout = 10 * time.Second
	ReadTimeout
)

type SecretStorage[secretType any] interface {
	SaveJsonSecret(name string, secret *secretType) error
	SaveTextSecret(name, secret string) error
	GetJsonSecret(name string) (*secretType, error)
	GetTextSecret(name string) (string, error)
	DeleteSecret(name string) error
}
