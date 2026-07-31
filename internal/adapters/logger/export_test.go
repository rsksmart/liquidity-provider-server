package logger

import "time"

// This file is only included in the tests compilation, the purpose of these functions
// is for testing only.

// WithClock returns a copy of the config with a deterministic clock.
func (c Config) WithClock(clock func() time.Time) Config {
	c.clock = clock
	return c
}

// WithExit returns a copy of the config with a custom process-exit function so
// tests can assert Fatal behaviour without terminating the test binary.
func (c Config) WithExit(exit func(int)) Config {
	c.exit = exit
	return c
}
