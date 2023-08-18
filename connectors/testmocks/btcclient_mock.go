package testmocks

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/mock"
)

type BTCClientMock struct {
	mock.Mock
}

func (B *BTCClientMock) CreateRawTransaction(inputs []btcjson.TransactionInput, amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {
	return nil, nil
}

func (B *BTCClientMock) SignRawTransactionWithWallet(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	return nil, false, nil
}

func (B *BTCClientMock) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return nil, nil
}

func (B *BTCClientMock) SetTxFee(fee btcutil.Amount) error {
	return nil
}

func (B *BTCClientMock) LockUnspent(shouldUnlock bool, txToUnlock []*wire.OutPoint) error {
	args := B.Called(shouldUnlock, txToUnlock)
	return args.Error(0)
}

func (B *BTCClientMock) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	args := B.Called()
	return args.Get(0).([]btcjson.ListUnspentResult), args.Error(1)
}

func (B *BTCClientMock) ListLockUnspent() ([]*wire.OutPoint, error) {
	args := B.Called()
	return args.Get(0).([]*wire.OutPoint), args.Error(1)
}

func (B *BTCClientMock) GetTxOut(txHash *chainhash.Hash, index uint32, mempool bool) (*btcjson.GetTxOutResult, error) {
	args := B.Called(txHash, index, mempool)
	return args.Get(0).(*btcjson.GetTxOutResult), args.Error(1)
}

func (B *BTCClientMock) SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error) {
	args := B.Called(address, amount)
	return args.Get(0).(*chainhash.Hash), args.Error(1)
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

func (B *BTCClientMock) GetBalance(address string) (btcutil.Amount, error) {
	args := B.Called(address)
	return args.Get(0).(btcutil.Amount), args.Error(1)
}
