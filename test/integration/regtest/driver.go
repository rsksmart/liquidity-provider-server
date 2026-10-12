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
	btc       *Bitcoin
	rsk       *Rootstock
	interval  time.Duration
	eachRound func(t *testing.T)
}

// NewDriver waits interval between rounds, which gives the LPS watchers time to run.
func NewDriver(btc *Bitcoin, rsk *Rootstock, interval time.Duration) *Driver {
	return &Driver{btc: btc, rsk: rsk, interval: interval}
}

// WithEachRound returns a driver that runs action at the start of every round, before the
// condition is checked, for something that has to keep happening while the stack advances.
func (driver *Driver) WithEachRound(action func(t *testing.T)) *Driver {
	copied := *driver
	copied.eachRound = action
	return &copied
}

// AdvanceUntil checks condition and, until it holds, mines step and waits. It logs the state it
// reached once, unless that state already held, and fails the test when timeout passes, naming
// the phase and the last observed state.
func (driver *Driver) AdvanceUntil(t *testing.T, phase string, timeout time.Duration, step Step, condition Condition) {
	t.Helper()
	start := time.Now()
	deadline := start.Add(timeout)
	for round := 0; ; round++ {
		if driver.eachRound != nil {
			driver.eachRound(t)
		}
		done, observed := condition(t)
		if done {
			if round > 0 {
				t.Logf("[%s] %s (after %s)", phase, observed, time.Since(start).Round(100*time.Millisecond))
			}
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
