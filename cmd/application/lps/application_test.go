package lps_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/cmd/application/lps"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
