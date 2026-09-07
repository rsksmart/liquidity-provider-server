package watcher_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	watcherAdapter "github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sessionTicker is the clock a watcherSession drives.
type sessionTicker struct {
	ticks chan time.Time
	// selects counts how often the loop has entered its select statement. Go re-evaluates the
	// channel expression on every pass, so a call to C() is the signal that the previous poll
	// returned and the loop is idle again.
	selects atomic.Int64
	stops   atomic.Int64
}

func newSessionTicker() *sessionTicker {
	return &sessionTicker{ticks: make(chan time.Time)}
}

func (ticker *sessionTicker) C() <-chan time.Time {
	ticker.selects.Add(1)
	return ticker.ticks
}

func (ticker *sessionTicker) Stop() {
	ticker.stops.Add(1)
}

// watcherSession is one process lifetime of a started watcher. Starting a second session over the
// same persisted state is what a restart looks like.
type watcherSession struct {
	t            *testing.T
	watcher      watcherAdapter.Watcher
	ticker       *sessionTicker
	closeChannel chan bool
}

func startWatcherSession(t *testing.T, watcher watcherAdapter.Watcher, ticker *sessionTicker) *watcherSession {
	t.Helper()
	session := &watcherSession{
		t:            t,
		watcher:      watcher,
		ticker:       ticker,
		closeChannel: make(chan bool),
	}
	require.NoError(t, watcher.Prepare(context.Background()))
	go watcher.Start()
	session.waitUntilIdle(0)
	return session
}

// poll delivers one tick and returns once the watcher has finished acting on it.
func (session *watcherSession) poll() {
	session.t.Helper()
	idleBefore := session.ticker.selects.Load()
	select {
	case session.ticker.ticks <- time.Now():
	case <-time.After(time.Second):
		require.FailNow(session.t, "watcher did not accept the tick")
	}
	session.waitUntilIdle(idleBefore)
}

// waitUntilIdle blocks until the loop has entered its select more often than it had at idleBefore,
// which is the point at which the poll it was running has returned.
func (session *watcherSession) waitUntilIdle(idleBefore int64) {
	session.t.Helper()
	require.Eventually(
		session.t,
		func() bool { return session.ticker.selects.Load() > idleBefore },
		time.Second,
		time.Millisecond,
		"watcher poll did not return",
	)
}

func (session *watcherSession) stop() {
	session.t.Helper()
	go session.watcher.Shutdown(session.closeChannel)
	select {
	case <-session.closeChannel:
	case <-time.After(time.Second):
		require.FailNow(session.t, "watcher shutdown did not complete")
	}
	// Shutdown answers on closeChannel independently of the loop, so the ticker is released a moment
	// later.
	assert.Eventually(
		session.t,
		func() bool { return session.ticker.stops.Load() == 1 },
		time.Second,
		time.Millisecond,
		"shutdown must stop the watcher's ticker exactly once",
	)
}
