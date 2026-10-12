package regtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Step is what the driver mines on each round. Mine RSK blocks on every round of a phase that
// needs the bridge: powpeg relays Bitcoin headers to it only while RSK blocks are being mined.
type Step struct {
	BtcBlocks int64
	RskBlocks int
}

// Condition reports whether the awaited state is reached and what it observed, for the log and
// the failure message.
type Condition func(t *testing.T) (done bool, observed string)

// Driver advances both chains until the stack reaches a state. Neither regtest chain advances
// by itself, so every wait in a test goes through it.
type Driver struct {
	btc      *Bitcoin
	rsk      *Rootstock
	interval time.Duration
}

// NewDriver waits interval between rounds, which gives the LPS watchers time to run.
func NewDriver(btc *Bitcoin, rsk *Rootstock, interval time.Duration) *Driver {
	return &Driver{btc: btc, rsk: rsk, interval: interval}
}

// AdvanceUntil checks condition and, until it holds, mines step and waits. It fails the test
// when timeout passes, naming the phase and the last observed state.
func (driver *Driver) AdvanceUntil(t *testing.T, phase string, timeout time.Duration, step Step, condition Condition) {
	t.Helper()
	start := time.Now()
	deadline := start.Add(timeout)
	lastObserved := ""
	for {
		done, observed := condition(t)
		if observed != lastObserved {
			t.Logf("[%s] %s", phase, observed)
			lastObserved = observed
		}
		if done {
			t.Logf("[%s] done after %s", phase, time.Since(start).Round(100*time.Millisecond))
			return
		}
		if time.Now().After(deadline) {
			require.FailNowf(t, phase+" timed out", "not done after %s; last observed: %s", timeout, observed)
		}
		if step.BtcBlocks > 0 {
			driver.btc.Mine(t, step.BtcBlocks)
		}
		if step.RskBlocks > 0 {
			driver.rsk.Mine(t, step.RskBlocks)
		}
		select {
		case <-t.Context().Done():
			require.FailNowf(t, phase+" interrupted", "last observed: %s", observed)
		case <-time.After(driver.interval):
		}
	}
}
