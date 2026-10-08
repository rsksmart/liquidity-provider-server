package watcher_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type resolveWatcherHarness struct {
	t              *testing.T
	resolveUseCase *mocks.PegInClaimResolverMock
	btc            *mocks.BitcoinNetworkMock
	ticks          chan time.Time
	events         chan entities.Event
	watcher        *watcher.PegInResolveWatcher
	done           chan struct{}
}

// newResolveWatcherHarness prepares the watcher with the given claimed rows and starts it
func newResolveWatcherHarness(t *testing.T, claims ...rootstock.PegInClaim) *resolveWatcherHarness {
	t.Helper()
	h := &resolveWatcherHarness{
		t:              t,
		resolveUseCase: mocks.NewPegInClaimResolverMock(t),
		btc:            mocks.NewBitcoinNetworkMock(t),
		ticks:          make(chan time.Time),
		events:         make(chan entities.Event),
		done:           make(chan struct{}),
	}
	getClaimsUseCase := mocks.NewPegInClaimsGetterMock(t)
	getClaimsUseCase.EXPECT().Run(mock.Anything, rootstock.PegInClaimClaimed).Return(claims, nil).Once()
	eventBus := mocks.NewEventBusMock(t)
	eventBus.EXPECT().Subscribe(rootstock.PegInClaimCompletedEventId).Return(h.events).Once()
	ticker := mocks.NewTickerMock(t)
	ticker.EXPECT().C().Return(h.ticks)
	ticker.EXPECT().Stop().Return().Once()

	h.watcher = watcher.NewPegInResolveWatcher(h.resolveUseCase, getClaimsUseCase, h.btc, eventBus, ticker)
	require.NoError(t, h.watcher.Prepare(context.Background()))
	go func() {
		// A failed mock expectation ends this goroutine, so done must still be closed
		defer close(h.done)
		h.watcher.Start()
	}()
	return h
}

func (h *resolveWatcherHarness) tick(height int64) {
	h.btc.EXPECT().GetHeight().Return(big.NewInt(height), nil).Once()
	h.sendTick()
}

func (h *resolveWatcherHarness) sendTick() {
	select {
	case h.ticks <- time.Now():
	case <-h.done:
		h.t.Fatal("watcher stopped before the tick")
	}
}

func (h *resolveWatcherHarness) sendEvent(event entities.Event) {
	select {
	case h.events <- event:
	case <-h.done:
		h.t.Fatal("watcher stopped before the event")
	}
}

// stop returns once the watcher loop has exited, so every pass is finished
func (h *resolveWatcherHarness) stop() {
	closeChannel := make(chan bool, 1)
	h.watcher.Shutdown(closeChannel)
	<-closeChannel
	<-h.done
}

func claimedPegIn(pegInID string) rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  test.AnyRskAddress,
		DepositTxID: "deposit-" + pegInID,
		State:       rootstock.PegInClaimClaimed,
		PegInID:     pegInID,
	}
}

func TestPegInResolveWatcher_Prepare(t *testing.T) {
	t.Run("loads claimed rows and runs no pass", func(t *testing.T) {
		h := newResolveWatcherHarness(t, claimedPegIn("a"), claimedPegIn("b"))
		h.stop()
	})
	t.Run("resolves the loaded rows on the next block", func(t *testing.T) {
		first, second := claimedPegIn("a"), claimedPegIn("b")
		h := newResolveWatcherHarness(t, first, second)
		h.resolveUseCase.EXPECT().Run(mock.Anything, first).Return(nil).Once()
		h.resolveUseCase.EXPECT().Run(mock.Anything, second).Return(nil).Once()
		h.tick(1)
		h.stop()
	})
	t.Run("returns the list error", func(t *testing.T) {
		getClaimsUseCase := mocks.NewPegInClaimsGetterMock(t)
		getClaimsUseCase.EXPECT().Run(mock.Anything, rootstock.PegInClaimClaimed).Return(nil, assert.AnError).Once()
		resolveWatcher := watcher.NewPegInResolveWatcher(nil, getClaimsUseCase, nil, nil, nil)
		require.ErrorIs(t, resolveWatcher.Prepare(context.Background()), assert.AnError)
	})
}

func TestPegInResolveWatcher_ClaimCompletedEvent(t *testing.T) {
	t.Run("adds a claimed claim and resolves it on the next block", func(t *testing.T) {
		claim := claimedPegIn("a")
		h := newResolveWatcherHarness(t)
		h.sendEvent(rootstock.NewPegInClaimCompletedEvent(claim))
		h.resolveUseCase.EXPECT().Run(mock.Anything, claim).Return(nil).Once()
		h.tick(1)
		h.stop()
	})
	t.Run("ignores a race_lost claim", func(t *testing.T) {
		claim := claimedPegIn("a")
		claim.State = rootstock.PegInClaimRaceLost
		h := newResolveWatcherHarness(t)
		h.sendEvent(rootstock.NewPegInClaimCompletedEvent(claim))
		h.tick(1)
		h.stop()
	})
	t.Run("doesn't add a claim twice", func(t *testing.T) {
		claim := claimedPegIn("a")
		h := newResolveWatcherHarness(t, claim)
		h.sendEvent(rootstock.NewPegInClaimCompletedEvent(claim))
		h.resolveUseCase.EXPECT().Run(mock.Anything, claim).Return(nil).Once()
		h.tick(1)
		h.stop()
	})
	t.Run("logs a wrong event", func(t *testing.T) {
		checkLog := test.AssertLogContains(t, watcher.LogPegInResolveWrongEvent)
		h := newResolveWatcherHarness(t)
		h.sendEvent(quote.PegoutQuoteCompletedEvent{Event: entities.NewBaseEvent(quote.PegoutQuoteCompletedEventId)})
		h.stop()
		assert.True(t, checkLog())
	})
}

func TestPegInResolveWatcher_HeightGate(t *testing.T) {
	t.Run("runs one pass per new block", func(t *testing.T) {
		claim := claimedPegIn("a")
		h := newResolveWatcherHarness(t, claim)
		h.resolveUseCase.EXPECT().Run(mock.Anything, claim).Return(usecases.NoEnoughConfirmationsError).Times(2)
		h.tick(1)
		h.tick(1)
		h.tick(2)
		h.stop()
	})
	t.Run("skips the pass on a height error", func(t *testing.T) {
		checkLog := test.AssertLogContains(t, "PegInResolveWatcher: error getting Bitcoin chain height")
		h := newResolveWatcherHarness(t, claimedPegIn("a"))
		h.btc.EXPECT().GetHeight().Return(nil, assert.AnError).Once()
		h.sendTick()
		h.stop()
		assert.True(t, checkLog())
	})
	t.Run("skips the pass on a nil height", func(t *testing.T) {
		h := newResolveWatcherHarness(t, claimedPegIn("a"))
		h.btc.EXPECT().GetHeight().Return(nil, nil).Once()
		h.sendTick()
		h.stop()
	})
}

func TestPegInResolveWatcher_Results(t *testing.T) {
	cases := []struct {
		name  string
		err   error
		times int
	}{
		{name: "resolved is not run again", err: nil, times: 1},
		{name: "non recoverable is not run again", err: errors.Join(assert.AnError, usecases.NonRecoverableError), times: 1},
		{name: "not enough confirmations runs again", err: usecases.NoEnoughConfirmationsError, times: 2},
		{name: "waiting for bridge runs again", err: fmt.Errorf("resolve: %w", blockchain.WaitingForBridgeError), times: 2},
		{name: "recoverable error runs again", err: assert.AnError, times: 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			claim := claimedPegIn("a")
			h := newResolveWatcherHarness(t, claim)
			h.resolveUseCase.EXPECT().Run(mock.Anything, claim).Return(c.err).Times(c.times)
			h.tick(1)
			h.tick(2)
			h.stop()
		})
	}
	t.Run("an error on one claim doesn't stop the others", func(t *testing.T) {
		failing, resolved := claimedPegIn("a"), claimedPegIn("b")
		h := newResolveWatcherHarness(t, failing, resolved)
		h.resolveUseCase.EXPECT().Run(mock.Anything, failing).Return(assert.AnError).Once()
		h.resolveUseCase.EXPECT().Run(mock.Anything, resolved).Return(nil).Once()
		h.tick(1)
		h.stop()
	})
}

func TestPegInResolveWatcher_StopsPass(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{name: "infrastructure unavailable", err: usecases.JoinInfrastructureUnavailable(assert.AnError)},
		{name: "contract paused", err: fmt.Errorf("%w. Pause level %d", blockchain.ContractPausedError, blockchain.PauseLevelHard)},
	}
	// Every claim returns the stop error, so each pass runs exactly one claim, whatever the map order
	for _, c := range cases {
		t.Run(c.name+" stops the pass", func(t *testing.T) {
			h := newResolveWatcherHarness(t, claimedPegIn("a"), claimedPegIn("b"))
			h.resolveUseCase.EXPECT().Run(mock.Anything, mock.Anything).Return(c.err).Once()
			h.tick(1)
			h.stop()
		})
		t.Run(c.name+" retries the same block", func(t *testing.T) {
			h := newResolveWatcherHarness(t, claimedPegIn("a"), claimedPegIn("b"))
			h.resolveUseCase.EXPECT().Run(mock.Anything, mock.Anything).Return(c.err).Times(2)
			h.tick(1)
			h.tick(1)
			h.stop()
		})
	}
	t.Run("logs the aborted pass", func(t *testing.T) {
		checkLog := test.AssertLogContains(t, "PegInResolveWatcher: aborted resolve pass")
		h := newResolveWatcherHarness(t, claimedPegIn("a"))
		h.resolveUseCase.EXPECT().Run(mock.Anything, mock.Anything).Return(usecases.JoinInfrastructureUnavailable(assert.AnError)).Once()
		h.tick(1)
		h.stop()
		assert.True(t, checkLog())
	})
}

func TestPegInResolveWatcher_EmptyPegInIdSharesOneEntry(t *testing.T) {
	first, second := claimedPegIn(""), claimedPegIn("")
	second.DepositTxID = "other-deposit"
	h := newResolveWatcherHarness(t, first, second)
	h.resolveUseCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil).Once()
	h.tick(1)
	h.stop()
}

func TestPegInResolveWatcher_Shutdown(t *testing.T) {
	eventBus := mocks.NewEventBusMock(t)
	eventBus.EXPECT().Subscribe(rootstock.PegInClaimCompletedEventId).Return(make(<-chan entities.Event)).Once()
	createWatcherShutdownTest(t, func(ticker utils.Ticker) watcher.Watcher {
		return watcher.NewPegInResolveWatcher(nil, nil, nil, eventBus, ticker)
	})
}
