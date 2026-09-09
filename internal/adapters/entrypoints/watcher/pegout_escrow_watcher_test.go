package watcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const escrowRequestHash = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"

func matchUin64Ptr(expected uint64) interface{} {
	return mock.MatchedBy(func(v *uint64) bool {
		return v != nil && *v == expected
	})
}

func newEscrowWatcherFixtures(t *testing.T) (
	*mocks.PegOutEscrowContractMock,
	*mocks.PegOutEscrowWatchRepositoryMock,
	*mocks.RootstockRpcServerMock,
	*mocks.TickerMock,
	chan time.Time,
) {
	t.Helper()
	escrow := &mocks.PegOutEscrowContractMock{}
	repo := &mocks.PegOutEscrowWatchRepositoryMock{}
	rskRpc := &mocks.RootstockRpcServerMock{}
	ticker := &mocks.TickerMock{}
	tickerChannel := make(chan time.Time)
	ticker.EXPECT().C().Return(tickerChannel)
	ticker.EXPECT().Stop().Return().Maybe()
	return escrow, repo, rskRpc, ticker, tickerChannel
}

func foreignRequested() blockchain.PegOutRequested {
	return blockchain.PegOutRequested{
		RequestHash:        escrowRequestHash,
		RefundAddress:      "0x1111111111111111111111111111111111111111",
		Amount:             entities.NewWei(1_000_000),
		DestinationAddress: []byte{0x01, 0x02},
		TxHash:             "0xreqtx",
		BlockNumber:        50,
	}
}

func TestPegoutEscrowWatcher_Step5_DiscoversForeignRequest(t *testing.T) {
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowWatcherFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{}, nil).Once()

	requested := foreignRequested()
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutRequested{requested}, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutClaimed{}, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutCancelled{}, nil).Once()
	repo.On("UpsertCandidate", mock.Anything, requested).Return(nil).Once()
	repo.On("SetCheckpoint", mock.Anything, uint64(20)).Return(nil).Once()

	w := watcher.NewPegoutEscrowWatcher(contracts, rpc, repo, nil, ticker, 0, 2000, time.Second)
	require.NoError(t, w.Prepare(context.Background()))
	go w.Start()
	tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		candidates := w.GetCandidates()
		assert.Len(c, candidates, 1)
		if len(candidates) == 1 {
			assert.Equal(c, escrowRequestHash, candidates[0].RequestHash)
			assert.Equal(c, requested.Amount, candidates[0].Amount)
		}
		assert.Equal(c, uint64(20), w.LastScannedBlock())
	}, time.Second, 10*time.Millisecond)

	closeCh := make(chan bool)
	go w.Shutdown(closeCh)
	<-closeCh
	repo.AssertExpectations(t)
	escrow.AssertExpectations(t)
	rskRpc.AssertExpectations(t)
}

func TestPegoutEscrowWatcher_Step5_DropsOnClaimedEvent(t *testing.T) {
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowWatcherFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	requested := foreignRequested()

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{requested}, nil).Once()
	escrow.EXPECT().GetPegOutState(escrowRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()

	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutRequested{}, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutClaimed{{
		LpAddress:   "0xlp",
		RequestHash: escrowRequestHash,
		TxHash:      "0xclaim",
		BlockNumber: 15,
	}}, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutCancelled{}, nil).Once()
	repo.On("DeleteCandidate", mock.Anything, escrowRequestHash).Return(nil).Once()
	repo.On("SetCheckpoint", mock.Anything, uint64(20)).Return(nil).Once()

	w := watcher.NewPegoutEscrowWatcher(contracts, rpc, repo, nil, ticker, 0, 2000, time.Second)
	require.NoError(t, w.Prepare(context.Background()))
	require.Len(t, w.GetCandidates(), 1)

	go w.Start()
	tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Empty(c, w.GetCandidates())
		assert.Equal(c, uint64(20), w.LastScannedBlock())
	}, time.Second, 10*time.Millisecond)

	closeCh := make(chan bool)
	go w.Shutdown(closeCh)
	<-closeCh
	repo.AssertExpectations(t)
	escrow.AssertExpectations(t)
}

func TestPegoutEscrowWatcher_Step5_DropsOnCancelledEvent(t *testing.T) {
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowWatcherFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	requested := foreignRequested()

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{requested}, nil).Once()
	escrow.EXPECT().GetPegOutState(escrowRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()

	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutRequested{}, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutClaimed{}, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutCancelled{{
		RequestHash: escrowRequestHash,
		TxHash:      "0xcancel",
		BlockNumber: 16,
	}}, nil).Once()
	repo.On("DeleteCandidate", mock.Anything, escrowRequestHash).Return(nil).Once()
	repo.On("SetCheckpoint", mock.Anything, uint64(20)).Return(nil).Once()

	w := watcher.NewPegoutEscrowWatcher(contracts, rpc, repo, nil, ticker, 0, 2000, time.Second)
	require.NoError(t, w.Prepare(context.Background()))
	require.Len(t, w.GetCandidates(), 1)

	go w.Start()
	tickerChannel <- time.Now()

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Empty(c, w.GetCandidates())
	}, time.Second, 10*time.Millisecond)

	closeCh := make(chan bool)
	go w.Shutdown(closeCh)
	<-closeCh
	repo.AssertExpectations(t)
	escrow.AssertExpectations(t)
}

func TestPegoutEscrowWatcher_Step5_ReconcilesStaleCandidateOnRestart(t *testing.T) {
	escrow, repo, rskRpc, ticker, _ := newEscrowWatcherFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}

	stale := foreignRequested()
	stale.RequestHash = "stalehash000000000000000000000000000000000000000000000000000000000"
	alive := foreignRequested()

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(100), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{stale, alive}, nil).Once()
	escrow.EXPECT().GetPegOutState(stale.RequestHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()
	escrow.EXPECT().GetPegOutState(alive.RequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	repo.On("DeleteCandidate", mock.Anything, stale.RequestHash).Return(nil).Once()

	w := watcher.NewPegoutEscrowWatcher(contracts, rpc, repo, nil, ticker, 0, 2000, time.Second)
	require.NoError(t, w.Prepare(context.Background()))

	candidates := w.GetCandidates()
	require.Len(t, candidates, 1)
	assert.Equal(t, alive.RequestHash, candidates[0].RequestHash)
	assert.Equal(t, uint64(100), w.LastScannedBlock())
	repo.AssertExpectations(t)
	escrow.AssertExpectations(t)
	_ = rskRpc
}

func TestPegoutEscrowWatcher_Step5_AdvancesCheckpointWithoutReprocess(t *testing.T) {
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowWatcherFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{}, nil).Once()

	// First tick: scan 11..15
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(15), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(15)).Return(nil, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(15)).Return(nil, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(15)).Return(nil, nil).Once()
	repo.On("SetCheckpoint", mock.Anything, uint64(15)).Return(nil).Once()

	// Second tick: must start at 16, not re-fetch 11..15
	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(18), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(16), matchUin64Ptr(18)).Return(nil, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(16), matchUin64Ptr(18)).Return(nil, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(16), matchUin64Ptr(18)).Return(nil, nil).Once()
	repo.On("SetCheckpoint", mock.Anything, uint64(18)).Return(nil).Once()

	w := watcher.NewPegoutEscrowWatcher(contracts, rpc, repo, nil, ticker, 0, 2000, time.Second)
	require.NoError(t, w.Prepare(context.Background()))
	go w.Start()

	tickerChannel <- time.Now()
	assert.Eventually(t, func() bool { return w.LastScannedBlock() == 15 }, time.Second, 10*time.Millisecond)

	tickerChannel <- time.Now()
	assert.Eventually(t, func() bool { return w.LastScannedBlock() == 18 }, time.Second, 10*time.Millisecond)

	closeCh := make(chan bool)
	go w.Shutdown(closeCh)
	<-closeCh
	repo.AssertExpectations(t)
	escrow.AssertExpectations(t)
	rskRpc.AssertExpectations(t)
}
