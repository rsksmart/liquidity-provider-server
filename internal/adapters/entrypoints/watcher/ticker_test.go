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
	// A 10 ms ticker fires well inside one second; a ticker that ignored the argument would not.
	tickers := watcher.NewApplicationTickers(10*time.Millisecond, 10*time.Millisecond)
	defer tickers.PegInAddressRegistryWatcherTicker.Stop()
	defer tickers.PegInClaimWatcherTicker.Stop()
	for name, ticker := range map[string]<-chan time.Time{
		"address registry": tickers.PegInAddressRegistryWatcherTicker.C(),
		"claim":            tickers.PegInClaimWatcherTicker.C(),
	} {
		select {
		case <-ticker:
		case <-time.After(time.Second):
			t.Fatalf("%s ticker did not use the provided interval", name)
		}
	}
}
