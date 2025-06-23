package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewTickerWrapper(t *testing.T) {
	duration := time.Second
	ticker := utils.NewTickerWrapper(duration)
	require.NotNil(t, ticker)
	require.NotNil(t, ticker.C())
	ticker.Stop()
}

func TestTickerWrapperChannelReceivesTicks(t *testing.T) {
	duration := 10 * time.Millisecond
	ticker := utils.NewTickerWrapper(duration)
	defer ticker.Stop()

	select {
	case <-ticker.C():
		// Tick received
	case <-time.After(50 * time.Millisecond):
		t.FailNow()
	}
}

func TestTickerWrapperStopPreventsTicks(t *testing.T) {
	duration := 10 * time.Millisecond
	ticker := utils.NewTickerWrapper(duration)
	ticker.Stop()

	select {
	case <-ticker.C():
		t.FailNow()
	case <-time.After(50 * time.Millisecond):
		// No tick received
	}
}
