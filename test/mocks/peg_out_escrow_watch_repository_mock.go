package mocks

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/mock"
)

// PegOutEscrowWatchRepositoryMock is a testify mock for blockchain.PegOutEscrowWatchRepository.
type PegOutEscrowWatchRepositoryMock struct {
	mock.Mock
}

func (m *PegOutEscrowWatchRepositoryMock) GetCheckpoint(ctx context.Context) (uint64, bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Bool(1), args.Error(2)
}

func (m *PegOutEscrowWatchRepositoryMock) SetCheckpoint(ctx context.Context, lastScannedBlock uint64) error {
	args := m.Called(ctx, lastScannedBlock)
	return args.Error(0)
}

func (m *PegOutEscrowWatchRepositoryMock) UpsertCandidate(ctx context.Context, candidate blockchain.PegOutRequested) error {
	args := m.Called(ctx, candidate)
	return args.Error(0)
}

func (m *PegOutEscrowWatchRepositoryMock) DeleteCandidate(ctx context.Context, requestHash string) error {
	args := m.Called(ctx, requestHash)
	return args.Error(0)
}

func (m *PegOutEscrowWatchRepositoryMock) ListCandidates(ctx context.Context) ([]blockchain.PegOutRequested, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]blockchain.PegOutRequested), args.Error(1)
}
