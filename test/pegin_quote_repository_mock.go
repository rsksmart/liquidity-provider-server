package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/stretchr/testify/mock"
)

type PeginQuoteRepositoryMock struct {
	quote.PeginQuoteRepository
	mock.Mock
}

func (m *PeginQuoteRepositoryMock) GetQuote(ctx context.Context, hash string) (*quote.PeginQuote, error) {
	args := m.Called(ctx, hash)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.(*quote.PeginQuote), args.Error(1)
	}
}

func (m *PeginQuoteRepositoryMock) GetRetainedQuoteByState(ctx context.Context, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, states)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.([]quote.RetainedPeginQuote), args.Error(1)
	}
}

func (m *PeginQuoteRepositoryMock) DeleteQuotes(ctx context.Context, hashes []string) (uint, error) {
	args := m.Called(ctx, hashes)
	return args.Get(0).(uint), args.Error(1)
}
