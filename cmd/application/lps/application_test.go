package lps_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/watcher"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPegInWatcherIntervals(t *testing.T) {
	t.Run("uses the defaults when the intervals are not set", func(t *testing.T) {
		intervals := lps.NewPegInWatcherIntervals(environment.PeginEnv{})
		assert.Equal(t, watcher.PegInWatcherIntervals{
			AddressRegistry: time.Minute,
			Claim:           time.Minute,
			Resolve:         3 * time.Minute,
		}, intervals)
	})
	t.Run("converts the configured seconds", func(t *testing.T) {
		intervals := lps.NewPegInWatcherIntervals(environment.PeginEnv{
			AddressRegistryWatcherIntervalSeconds: 2,
			ClaimWatcherIntervalSeconds:           3,
			ResolveWatcherIntervalSeconds:         4,
		})
		assert.Equal(t, watcher.PegInWatcherIntervals{
			AddressRegistry: 2 * time.Second,
			Claim:           3 * time.Second,
			Resolve:         4 * time.Second,
		}, intervals)
	})
}

func TestNewBitcoinRegistry_Success(t *testing.T) {
	paymentWallet := mocks.NewBitcoinWalletMock(t)
	monitoringWallet := mocks.NewBitcoinWalletMock(t)
	factory := mocks.NewAbstractFactoryMock(t)
	factory.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).Return(paymentWallet, nil)
	factory.On("BitcoinMonitoringWallet", bitcoin.PeginWalletId).Return(monitoringWallet, nil)

	exitCalled := false
	mockExit := func(int) { exitCalled = true }

	result := lps.NewBitcoinRegistry(factory, nil, mockExit)

	require.NotNil(t, result)
	assert.False(t, exitCalled)
}

func TestNewBitcoinRegistry_WalletScanning(t *testing.T) {
	factory := mocks.NewAbstractFactoryMock(t)
	factory.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).
		Return(nil, fmt.Errorf("rescan triggered: %w", bitcoin.ErrWalletScanning))

	var capturedCode int
	mockExit := func(code int) { capturedCode = code }

	result := lps.NewBitcoinRegistry(factory, nil, mockExit)

	assert.Nil(t, result)
	assert.Equal(t, 0, capturedCode)
}

func TestNewBitcoinRegistry_Error(t *testing.T) {
	factory := mocks.NewAbstractFactoryMock(t)
	factory.On("BitcoinPaymentWallet", bitcoin.DerivativeWalletId).
		Return(nil, errors.New("unexpected rpc error"))

	var capturedCode int
	mockExit := func(code int) { capturedCode = code }

	result := lps.NewBitcoinRegistry(factory, nil, mockExit)

	assert.Nil(t, result)
	assert.Equal(t, 1, capturedCode)
}
