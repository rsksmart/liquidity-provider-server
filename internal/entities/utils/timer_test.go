package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewTimerWrapper(t *testing.T) {
	duration := time.Second
	timer := utils.NewTimerWrapper(duration)
	require.NotNil(t, timer)
	require.NotNil(t, timer.C())
	timer.Stop()
}

func TestTimerWrapperChannelReceivesSignal(t *testing.T) {
	duration := 10 * time.Millisecond
	timer := utils.NewTimerWrapper(duration)
	defer timer.Stop()

	select {
	case <-timer.C():
		// Signal received
	case <-time.After(50 * time.Millisecond):
		t.FailNow()
	}
}

func TestTimerWrapperStopPreventsSignal(t *testing.T) {
	duration := 10 * time.Millisecond
	timer := utils.NewTimerWrapper(duration)
	timer.Stop()

	select {
	case <-timer.C():
		t.FailNow()
	case <-time.After(50 * time.Millisecond):
		// No signal received
	}
}
