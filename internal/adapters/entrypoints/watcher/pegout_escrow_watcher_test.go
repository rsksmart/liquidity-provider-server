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
