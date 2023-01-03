package testmocks

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/mock"
)

type BTCClientMock struct {
	mock.Mock
}

func (B *BTCClientMock) SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error) {
	panic("implement me")
}

func (B *BTCClientMock) GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error) {
	args := B.Called()
	return args.Get(0).(*btcjson.GetNetworkInfoResult), args.Error(1)
}

func (B *BTCClientMock) ImportAddressRescan(address string, account string, rescan bool) error {
	args := B.Called(address, account, rescan)
	return args.Error(0)
}

func (B *BTCClientMock) GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error) {
	args := B.Called(txHash)
	return args.Get(0).(*btcjson.GetTransactionResult), args.Error(1)
}

func (B *BTCClientMock) GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	B.Called(blockHash)
	return new(btcjson.GetBlockVerboseResult), nil
}

func (B *BTCClientMock) ListUnspentMinMaxAddresses(minConf, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	args := B.Called(minConf, maxConf, addrs)
	return args.Get(0).([]btcjson.ListUnspentResult), args.Error(1)
}

func (B *BTCClientMock) GetBlock(blockHash *chainhash.Hash) (*wire.MsgBlock, error) {
	B.Called(blockHash)
	return new(wire.MsgBlock), nil
}

func (B *BTCClientMock) GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error) {
	B.Called(txHash)
	return new(btcutil.Tx), nil
}

func (B *BTCClientMock) Disconnect() {
	B.Called()
}
