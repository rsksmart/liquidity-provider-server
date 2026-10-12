package watcher_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestNewApplicationTickers(t *testing.T) {
	tickers := watcher.NewApplicationTickers(watcher.PegInWatcherIntervals{
		AddressRegistry: time.Hour,
		Claim:           time.Hour,
		Resolve:         time.Hour,
	})
	require.NotNil(t, tickers)
	value := reflect.ValueOf(tickers).Elem()
	for i := 0; i < value.Type().NumField(); i++ {
		if value.Field(i).IsNil() {
			t.Errorf("Field %s of application tickers is nil", value.Type().Field(i).Name)
		}
	}
}

func TestNewApplicationTickers_PegInIntervals(t *testing.T) {
	const fast, slow = 10 * time.Millisecond, time.Hour
	cases := []struct {
		name      string
		intervals watcher.PegInWatcherIntervals
		ticker    func(*watcher.ApplicationTickers) utils.Ticker
	}{
		{
			name:      "address registry",
			intervals: watcher.PegInWatcherIntervals{AddressRegistry: fast, Claim: slow, Resolve: slow},
			ticker: func(tickers *watcher.ApplicationTickers) utils.Ticker {
				return tickers.PegInAddressRegistryWatcherTicker
			},
		},
		{
			name:      "claim",
			intervals: watcher.PegInWatcherIntervals{AddressRegistry: slow, Claim: fast, Resolve: slow},
			ticker:    func(tickers *watcher.ApplicationTickers) utils.Ticker { return tickers.PegInClaimWatcherTicker },
		},
		{
			name:      "resolve",
			intervals: watcher.PegInWatcherIntervals{AddressRegistry: slow, Claim: slow, Resolve: fast},
			ticker:    func(tickers *watcher.ApplicationTickers) utils.Ticker { return tickers.PegInResolveWatcherTicker },
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name+" ticker uses its own interval", func(t *testing.T) {
			tickers := watcher.NewApplicationTickers(testCase.intervals)
			t.Cleanup(func() {
				tickers.PegInAddressRegistryWatcherTicker.Stop()
				tickers.PegInClaimWatcherTicker.Stop()
				tickers.PegInResolveWatcherTicker.Stop()
			})
			select {
			case <-testCase.ticker(tickers).C():
			case <-time.After(time.Second):
				assert.Fail(t, "the ticker did not tick with its configured interval")
			}
		})
	}
}
