package testmocks

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/mock"
)

type AddressWatcherMock struct {
	mock.Mock
}

func (a *AddressWatcherMock) OnNewConfirmation(txHash string, confirmations int64, amount btcutil.Amount) {
	a.Called(txHash, confirmations, amount)
}

func (a *AddressWatcherMock) OnExpire() {
	a.Called()
}

func (a *AddressWatcherMock) Done() <-chan struct{} {
	args := a.Called()
	return args.Get(0).(<-chan struct{})
}
