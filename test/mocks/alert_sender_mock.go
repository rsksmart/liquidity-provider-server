package mocks

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/stretchr/testify/mock"
)

type AlertSenderMock struct {
	entities.AlertSender
	mock.Mock
}

func (m *AlertSenderMock) SendAlert(ctx context.Context, subject, body string, recipients []string) error {
	args := m.Called(ctx, subject, body, recipients)
	return args.Error(0)
}
