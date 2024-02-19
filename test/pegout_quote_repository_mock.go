package test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/stretchr/testify/mock"
)

type PegoutQuoteRepositoryMock struct {
	quote.PegoutQuoteRepository
	mock.Mock
}

func (m *PegoutQuoteRepositoryMock) InsertQuote(ctx context.Context, hash string, quote quote.PegoutQuote) error {
	args := m.Called(ctx, hash, quote)
	return args.Error(0)
}

func (m *PegoutQuoteRepositoryMock) InsertRetainedQuote(ctx context.Context, quote quote.RetainedPegoutQuote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *PegoutQuoteRepositoryMock) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	args := m.Called(ctx, hash)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.(*quote.PegoutQuote), args.Error(1)
	}
}

func (m *PegoutQuoteRepositoryMock) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, hash)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.(*quote.RetainedPegoutQuote), args.Error(1)
	}
}

func (m *PegoutQuoteRepositoryMock) UpsertPegoutDeposits(ctx context.Context, deposits []quote.PegoutDeposit) error {
	args := m.Called(ctx, deposits)
	return args.Error(0)
}

func (m *PegoutQuoteRepositoryMock) UpdateRetainedQuote(ctx context.Context, quote quote.RetainedPegoutQuote) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

func (m *PegoutQuoteRepositoryMock) ListPegoutDepositsByAddress(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	args := m.Called(ctx, address)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.([]quote.PegoutDeposit), args.Error(1)
	}
}

func (m *PegoutQuoteRepositoryMock) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, states)
	arg := args.Get(0)
	if arg == nil {
		return nil, args.Error(1)
	} else {
		return arg.([]quote.RetainedPegoutQuote), args.Error(1)
	}
}

func (m *PegoutQuoteRepositoryMock) DeleteQuotes(ctx context.Context, hashes []string) (uint, error) {
	args := m.Called(ctx, hashes)
	return args.Get(0).(uint), args.Error(1)
}

func (m *PegoutQuoteRepositoryMock) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	args := m.Called(ctx, deposit)
	return args.Error(0)
}
