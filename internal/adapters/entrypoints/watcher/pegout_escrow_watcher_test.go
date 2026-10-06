package watcher_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	escrowWatcherRequestHash = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"
	escrowWatcherLpAddress   = "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b"
)

type escrowWatchRepositoryFake struct {
	candidates []blockchain.PegOutRequested
	deleted    []string
}

func (r *escrowWatchRepositoryFake) GetCheckpoint(context.Context) (uint64, bool, error) {
	return 100, true, nil
}

func (r *escrowWatchRepositoryFake) SetCheckpoint(context.Context, uint64) error {
	return nil
}

func (r *escrowWatchRepositoryFake) UpsertCandidate(context.Context, blockchain.PegOutRequested) error {
	return nil
}

func (r *escrowWatchRepositoryFake) DeleteCandidate(_ context.Context, requestHash string) error {
	r.deleted = append(r.deleted, requestHash)
	return nil
}

func (r *escrowWatchRepositoryFake) ListCandidates(context.Context) ([]blockchain.PegOutRequested, error) {
	return r.candidates, nil
}

type escrowWatcherFixtures struct {
	escrow     *mocks.PegOutEscrowContractMock
	pegout     *mocks.PegoutContractMock
	quoteRepo  *mocks.PegoutQuoteRepositoryMock
	lp         *mocks.ProviderMock
	rskRpc     *mocks.RootstockRpcServerMock
	repository *escrowWatchRepositoryFake
	watcher    *watcher.PegoutEscrowWatcher
}

func newEscrowWatcherFixtures() *escrowWatcherFixtures {
	f := &escrowWatcherFixtures{
		escrow:    &mocks.PegOutEscrowContractMock{},
		pegout:    &mocks.PegoutContractMock{},
		quoteRepo: &mocks.PegoutQuoteRepositoryMock{},
		lp:        &mocks.ProviderMock{},
		rskRpc:    &mocks.RootstockRpcServerMock{},
		repository: &escrowWatchRepositoryFake{candidates: []blockchain.PegOutRequested{
			{RequestHash: escrowWatcherRequestHash, Amount: entities.NewWei(1)},
		}},
	}
	contracts := blockchain.RskContracts{PegOut: f.pegout, PegOutEscrow: f.escrow}
	rpc := blockchain.Rpc{Rsk: f.rskRpc, Btc: &mocks.BtcRpcMock{}}
	claimUseCase := pegout.NewClaimPegOutUseCase(
		contracts, rpc, &mocks.BitcoinWalletMock{}, f.lp, f.quoteRepo, &mocks.EventBusMock{}, &sync.Mutex{},
	)
	f.watcher = watcher.NewPegoutEscrowWatcher(contracts, rpc, f.repository, claimUseCase, nil, 0, 0, time.Minute)
	f.escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.quoteRepo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimPending).Return([]quote.RetainedPegoutQuote{}, nil).Once()
	return f
}

func (f *escrowWatcherFixtures) expectNotPaused() {
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Twice()
	f.quoteRepo.On("GetRetainedQuote", mock.Anything, escrowWatcherRequestHash).Return(nil, nil).Once()
}

func TestPegoutEscrowWatcher_SkipsClaimsQuietlyWhilePaused(t *testing.T) {
	f := newEscrowWatcherFixtures()
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Reason: "soft"}, nil).Once()
	f.pegout.EXPECT().GetAddress().Return("0xpegout").Once()

	require.NoError(t, f.watcher.Prepare(context.Background()))
	f.pegout.AssertNumberOfCalls(t, "PausedStatus", 1)
	f.quoteRepo.AssertCalled(t, "GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimPending)
	f.quoteRepo.AssertNotCalled(t, "GetRetainedQuote", mock.Anything, escrowWatcherRequestHash)
	assert.Len(t, f.watcher.GetCandidates(), 1)
}

func TestPegoutEscrowWatcher_SkipsClaimsWhenPauseCheckFails(t *testing.T) {
	f := newEscrowWatcherFixtures()
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{}, assert.AnError).Once()

	require.NoError(t, f.watcher.Prepare(context.Background()))
	f.pegout.AssertNumberOfCalls(t, "PausedStatus", 1)
	f.quoteRepo.AssertNotCalled(t, "GetRetainedQuote", mock.Anything, escrowWatcherRequestHash)
	assert.Len(t, f.watcher.GetCandidates(), 1)
}

func TestPegoutEscrowWatcher_DropsClosedRequest(t *testing.T) {
	f := newEscrowWatcherFixtures()
	f.expectNotPaused()
	f.escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateNone, nil).Once()

	require.NoError(t, f.watcher.Prepare(context.Background()))
	assert.Equal(t, []string{escrowWatcherRequestHash}, f.repository.deleted)
	assert.Empty(t, f.watcher.GetCandidates())
}

func TestPegoutEscrowWatcher_KeepsSkippedRequest(t *testing.T) {
	f := newEscrowWatcherFixtures()
	f.expectNotPaused()
	f.escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.lp.On("RskAddress").Return(escrowWatcherLpAddress)
	f.escrow.EXPECT().RestrictedUntil(escrowWatcherLpAddress).Return(uint64(1_000), nil).Once()
	f.rskRpc.EXPECT().GetBlockByNumber(mock.Anything, mock.Anything).Return(blockchain.BlockInfo{Timestamp: time.Unix(10, 0)}, nil).Once()

	require.NoError(t, f.watcher.Prepare(context.Background()))
	assert.Empty(t, f.repository.deleted)
	assert.Len(t, f.watcher.GetCandidates(), 1)
}

func TestPegoutEscrowWatcher_KeepsRequestOnClaimError(t *testing.T) {
	f := newEscrowWatcherFixtures()
	f.expectNotPaused()
	f.escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateNone, assert.AnError).Once()

	require.NoError(t, f.watcher.Prepare(context.Background()))
	assert.Empty(t, f.repository.deleted)
	assert.Len(t, f.watcher.GetCandidates(), 1)
}

func matchUin64Ptr(expected uint64) interface{} {
	return mock.MatchedBy(func(v *uint64) bool {
		return v != nil && *v == expected
	})
}

func newEscrowScanFixtures(t *testing.T) (
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
		RequestHash:        escrowWatcherRequestHash,
		RefundAddress:      "0x1111111111111111111111111111111111111111",
		Amount:             entities.NewWei(1_000_000),
		DestinationAddress: []byte{0x01, 0x02},
		TxHash:             "0xreqtx",
		BlockNumber:        50,
	}
}

func TestPegoutEscrowWatcher_Step5_DiscoversForeignRequest(t *testing.T) {
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowScanFixtures(t)
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
			assert.Equal(c, escrowWatcherRequestHash, candidates[0].RequestHash)
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
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowScanFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	requested := foreignRequested()

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{requested}, nil).Once()
	escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()

	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutRequested{}, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutClaimed{{
		LpAddress:   "0xlp",
		RequestHash: escrowWatcherRequestHash,
		TxHash:      "0xclaim",
		BlockNumber: 15,
	}}, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutCancelled{}, nil).Once()
	repo.On("DeleteCandidate", mock.Anything, escrowWatcherRequestHash).Return(nil).Once()
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
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowScanFixtures(t)
	contracts := blockchain.RskContracts{PegOutEscrow: escrow}
	rpc := blockchain.Rpc{Rsk: rskRpc}
	requested := foreignRequested()

	repo.On("GetCheckpoint", mock.Anything).Return(uint64(10), true, nil).Once()
	repo.On("ListCandidates", mock.Anything).Return([]blockchain.PegOutRequested{requested}, nil).Once()
	escrow.EXPECT().GetPegOutState(escrowWatcherRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()

	rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(20), nil).Once()
	escrow.EXPECT().GetPegOutRequestedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutRequested{}, nil).Once()
	escrow.EXPECT().GetPegOutClaimedEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutClaimed{}, nil).Once()
	escrow.EXPECT().GetPegOutCancelledEvents(mock.Anything, uint64(11), matchUin64Ptr(20)).Return([]blockchain.PegOutCancelled{{
		RequestHash: escrowWatcherRequestHash,
		TxHash:      "0xcancel",
		BlockNumber: 16,
	}}, nil).Once()
	repo.On("DeleteCandidate", mock.Anything, escrowWatcherRequestHash).Return(nil).Once()
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
	escrow, repo, rskRpc, ticker, _ := newEscrowScanFixtures(t)
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
	escrow, repo, rskRpc, ticker, tickerChannel := newEscrowScanFixtures(t)
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
