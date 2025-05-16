package mocks

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/stretchr/testify/mock"
)

// AcceptQuoteUseCaseMock is a manual mock for the AcceptQuoteUseCase
type AcceptQuoteUseCaseMock struct {
	mock.Mock
}

func (m *AcceptQuoteUseCaseMock) Run(ctx context.Context, quoteHash, signature string) (quote.AcceptedQuote, error) {
	args := m.Called(ctx, quoteHash, signature)
	return args.Get(0).(quote.AcceptedQuote), args.Error(1)
}
