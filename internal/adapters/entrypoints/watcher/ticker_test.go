package watcher_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
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
