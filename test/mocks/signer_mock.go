package mocks

import (
	"github.com/stretchr/testify/mock"
)

type SignerMock struct {
	mock.Mock
}

func (m *SignerMock) SignBytes(msg []byte) ([]byte, error) {
	args := m.Called(msg)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *SignerMock) Validate(signature, hash string) bool {
	args := m.Called(signature, hash)
	return args.Bool(0)
}
