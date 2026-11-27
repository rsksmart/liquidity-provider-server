package watcher_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	resolveWatcherFirstClaim = rootstock.PegInClaim{
		RskAddress:  "0x79568c2989232dca1840087d73d403602364c0d4",
		DepositTxID: "2ba6da53badd14349c5d6379e88c345e88193598aad714815d4b57c691a9fbdf",
		State:       rootstock.PegInClaimClaimed,
	}
	resolveWatcherSecondClaim = rootstock.PegInClaim{
		RskAddress:  "0x8dccd82443b80ddde3690af86746bfd9d766f8d2",
		DepositTxID: "fda421ccdff7324a382067d1746f6a387132435de6af336a0ebbf3f720eaae4d",
		State:       rootstock.PegInClaimResolving,
	}
)

type resolveWatcherFixture struct {
	repository *mocks.PegInClaimRepositoryMock
	resolve    *mocks.ResolvePegInUseCaseMock
	btcRpc     *mocks.BitcoinNetworkMock
	ticker     *sessionTicker
}

func newResolveWatcherFixture(t *testing.T) *resolveWatcherFixture {
	t.Helper()
	return &resolveWatcherFixture{
		repository: mocks.NewPegInClaimRepositoryMock(t),
		resolve:    mocks.NewResolvePegInUseCaseMock(t),
		btcRpc:     mocks.NewBitcoinNetworkMock(t),
		ticker:     newSessionTicker(),
	}
}

func (fixture *resolveWatcherFixture) start(t *testing.T) *watcherSession {
	t.Helper()
	resolveWatcher := watcher.NewPegInResolveWatcher(
		w.NewGetWatchedPegInClaimsUseCase(fixture.repository),
		fixture.resolve,
		blockchain.Rpc{Btc: fixture.btcRpc},
		fixture.ticker,
	)
	return startWatcherSession(t, resolveWatcher, fixture.ticker)
}

func (fixture *resolveWatcherFixture) expectHeight(height int64) {
	fixture.btcRpc.EXPECT().GetHeight().Return(big.NewInt(height), nil).Once()
}

func (fixture *resolveWatcherFixture) expectUnresolved(claims ...rootstock.PegInClaim) {
	fixture.repository.EXPECT().ListByState(mock.Anything, rootstock.PegInClaimClaimed, rootstock.PegInClaimResolving).Return(claims, nil).Once()
}

func (fixture *resolveWatcherFixture) expectResolve(claim rootstock.PegInClaim, err error) {
	fixture.resolve.EXPECT().Run(mock.Anything, claim).Return(err).Once()
}

func TestPegInResolveWatcher_Shutdown(t *testing.T) {
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegInResolveWatcher(nil, nil, blockchain.Rpc{}, ticker)
	})
}

func TestPegInResolveWatcher_ResolvesEveryClaimedOrResolvingPegInOnANewBlock(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)
	fixture.expectHeight(100)
	fixture.expectUnresolved(resolveWatcherFirstClaim, resolveWatcherSecondClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, nil)
	fixture.expectResolve(resolveWatcherSecondClaim, nil)

	session.poll()
	session.stop()
}

func TestPegInResolveWatcher_ChecksOnlyWhenANewBlockArrives(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)

	fixture.expectHeight(100)
	fixture.expectUnresolved(resolveWatcherFirstClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, nil)
	session.poll()

	fixture.expectHeight(100)
	session.poll()

	fixture.expectHeight(99)
	session.poll()

	fixture.expectHeight(101)
	fixture.expectUnresolved(resolveWatcherFirstClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, nil)
	session.poll()

	session.stop()
}

func TestPegInResolveWatcher_SkipsTheTickWhenTheHeightCannotBeRead(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)
	assertLog := test.AssertLogContains(t, fmt.Sprintf(watcher.LogPegInResolveBtcHeightError, assert.AnError))
	fixture.btcRpc.EXPECT().GetHeight().Return(nil, assert.AnError).Once()

	session.poll()

	assert.True(t, assertLog())
	session.stop()
}

func TestPegInResolveWatcher_ListsTheSameBlockAgainWhenTheListingFails(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)
	listErr := usecases.WrapUseCaseError(usecases.GetWatchedPegInClaimsId, assert.AnError)
	assertLog := test.AssertLogContains(t, fmt.Sprintf(watcher.LogPegInResolveListError, listErr))
	fixture.expectHeight(100)
	fixture.repository.EXPECT().ListByState(mock.Anything, rootstock.PegInClaimClaimed, rootstock.PegInClaimResolving).Return(nil, assert.AnError).Once()
	session.poll()
	assert.True(t, assertLog())

	fixture.expectHeight(100)
	fixture.expectUnresolved(resolveWatcherFirstClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, nil)
	session.poll()

	session.stop()
}

func TestPegInResolveWatcher_KeepsResolvingTheOtherClaimsWhenOneFails(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)
	assertLog := test.AssertLogContains(t, watcher.LogPegInResolveError(
		resolveWatcherFirstClaim.RskAddress,
		resolveWatcherFirstClaim.DepositTxID,
		assert.AnError,
	))
	fixture.expectHeight(100)
	fixture.expectUnresolved(resolveWatcherFirstClaim, resolveWatcherSecondClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, assert.AnError)
	fixture.expectResolve(resolveWatcherSecondClaim, nil)

	session.poll()

	assert.True(t, assertLog())
	session.stop()
}

func TestPegInResolveWatcher_RetriesAFailedResolveOnTheNextBlock(t *testing.T) {
	fixture := newResolveWatcherFixture(t)
	session := fixture.start(t)
	fixture.expectHeight(100)
	fixture.expectUnresolved(resolveWatcherFirstClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, assert.AnError)
	session.poll()

	fixture.expectHeight(100)
	session.poll()

	fixture.expectHeight(101)
	fixture.expectUnresolved(resolveWatcherFirstClaim)
	fixture.expectResolve(resolveWatcherFirstClaim, nil)
	session.poll()

	session.stop()
}
