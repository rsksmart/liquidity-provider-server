package watcher_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/stretchr/testify/require"
)

func TestNewApplicationTickers(t *testing.T) {
	tickers := watcher.NewApplicationTickers(time.Minute, time.Minute)
	require.NotNil(t, tickers)
	value := reflect.ValueOf(tickers).Elem()
	for i := 0; i < value.Type().NumField(); i++ {
		if value.Field(i).IsNil() {
			t.Errorf("Field %s of application tickers is nil", value.Type().Field(i).Name)
		}
	}
}

func TestNewApplicationTickers_UsesProvidedIntervals(t *testing.T) {
	tickers := watcher.NewApplicationTickers(2*time.Second, 2*time.Second)
	require.NotNil(t, tickers.PegInAddressRegistryWatcherTicker)
	require.NotNil(t, tickers.PegInClaimWatcherTicker)
	tickers.PegInAddressRegistryWatcherTicker.Stop()
	tickers.PegInClaimWatcherTicker.Stop()
}
