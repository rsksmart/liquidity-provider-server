package testmocks

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/stretchr/testify/mock"
	"time"
)

type BtcMock struct {
	mock.Mock
}

func (b *BtcMock) AddAddressWatcher(address string, minAmount btcutil.Amount, interval time.Duration, exp time.Time, w connectors.AddressWatcher) error {
	b.Called(address, minAmount, interval, exp, w)
	return nil
}

func (b *BtcMock) Connect(endpoint string, username string, password string) error {
	b.Called(endpoint, username, password)
	return nil
}

func (b *BtcMock) CheckConnection() error {
	args := b.Called()
	return args.Error(0)
}

func (b *BtcMock) GetParams() chaincfg.Params {
	b.Called()
	return chaincfg.TestNet3Params
}

func (b *BtcMock) RemoveAddressWatcher(address string) {
	b.Called(address)
}

func (b *BtcMock) Close() {
	b.Called()
}

func (b *BtcMock) SerializePMT(txHash string) ([]byte, error) {
	b.Called(txHash)
	return nil, nil
}

func (b *BtcMock) SerializeTx(txHash string) ([]byte, error) {
	b.Called(txHash)
	return nil, nil
}

func (b *BtcMock) GetBlockNumberByTx(txHash string) (int64, error) {
	b.Called(txHash)
	return 0, nil
}

func (b *BtcMock) GetDerivedBitcoinAddress(fedInfo *connectors.FedInfo, userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error) {
	b.Called(fedInfo, userBtcRefundAddr, lbcAddress, lpBtcAddress, derivationArgumentsHash)
	return "", nil
}
