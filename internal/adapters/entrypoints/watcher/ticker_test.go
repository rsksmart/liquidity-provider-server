package watcher_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestNewApplicationTickers(t *testing.T) {
	tickers := watcher.NewApplicationTickers()
	require.NotNil(t, tickers)
	value := reflect.ValueOf(tickers).Elem()
	for i := 0; i < value.Type().NumField(); i++ {
		if value.Field(i).IsNil() {
			t.Errorf("Field %s of application tickers is nil", value.Type().Field(i).Name)
		}
	}
}

// The registry scan loop and the deposit reconciliation loop must not share a ticker: a stopped
// or drained ticker would otherwise take both loops down together.
func TestNewApplicationTickers_PegInAddressRegistryLoopsAreIndependent(t *testing.T) {
	tickers := watcher.NewApplicationTickers()
	assert.NotSame(t, tickers.PegInAddressRegistryWatcherTicker, tickers.PegInAddressRegistryDepositWatcherTicker)
	assert.NotEqual(t, tickers.PegInAddressRegistryWatcherTicker.C(), tickers.PegInAddressRegistryDepositWatcherTicker.C())
}
