package bitcoin_test

import (
	"cmp"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const (
	mockPassword      = "pwd"
	mockFeeRate       = 0.0001
	mockAddress       = "mx5ySMGiiDd9rjkfwcZkSCo3ATQ16PEiJM"
	testnetAddress    = "mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5"
	mainnetAddress    = "141dsd6YZxdKcmTZckG4Q9qGzJbR1Jc9kv"
	expiredTime       = 1711098457 // 2024-03-22
	unexpiredTime     = 1900400857 // 2030-03-22
	paymentScriptMock = "a payment script"
)

func TestBitcoindWallet_Unlock(t *testing.T) {
	expiredLockUntil := expiredTime
	nonExpiredLockUntil := unexpiredTime
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
		UnlockedUntil: &expiredLockUntil,
	}, nil).Once()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
		UnlockedUntil: &nonExpiredLockUntil,
	}, nil).Once()
	client.On("WalletPassphrase", mockPassword, int64(bitcoin.WalletUnlockSeconds)).Return(nil).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(nil, client), mockAddress, mockFeeRate, true, mockPassword)
	err := rpc.Unlock()
	require.NoError(t, err)
	err = rpc.Unlock()
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_Unlock_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(nil, client), mockAddress, mockFeeRate, true, mockPassword)
	err := rpc.Unlock()
	require.Error(t, err)
}

func TestBitcoindWallet_ImportAddress(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("ImportAddressRescan", testnetAddress, "", false).Return(nil).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	err := rpc.ImportAddress(testnetAddress)
	require.NoError(t, err)
	client.AssertExpectations(t)

	client = &mocks.ClientAdapterMock{}
	client.On("ImportAddressRescan", mainnetAddress, "", false).Return(nil).Once()
	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.MainNetParams, client), mockAddress, mockFeeRate, true, mockPassword)
	err = rpc.ImportAddress(mainnetAddress)
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_ImportAddress_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.MainNetParams, client), mockAddress, mockFeeRate, true, mockPassword)
	err := rpc.ImportAddress(testnetAddress)
	require.Error(t, err)

	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	err = rpc.ImportAddress(mainnetAddress)
	require.Error(t, err)
}

func TestBitcoindWallet_EstimateTxFees(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	amount := entities.NewWei(5000000000000000)
	var changePosition int64 = 2
	var input []btcjson.PsbtInput
	var lockTime *uint32
	var bip32Derivs *bool
	feeRate := mockFeeRate
	floatAmount, _ := amount.ToRbtc().Float64()
	client.On("WalletCreateFundedPsbt",
		input,
		[]btcjson.PsbtOutput{
			{testnetAddress: floatAmount},
			{"data": "0000000000000000000000000000000000000000000000000000000000000000"},
		},
		lockTime,
		&btcjson.WalletCreateFundedPsbtOpts{
			ChangePosition: &changePosition,
			FeeRate:        &feeRate,
		},
		bip32Derivs,
	).Return(&btcjson.WalletCreateFundedPsbtResult{
		Fee: 0.0006,
	}, nil).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	fee, err := rpc.EstimateTxFees(testnetAddress, amount)
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(600000000000000), fee)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_EstimateTxFees_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	fee, err := rpc.EstimateTxFees(mainnetAddress, entities.NewWei(1))
	require.Error(t, err)
	assert.Nil(t, fee)

	client.On("WalletCreateFundedPsbt",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, assert.AnError).Once()
	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	fee, err = rpc.EstimateTxFees(testnetAddress, entities.NewWei(1))
	require.Error(t, err)
	assert.Nil(t, fee)

	client.On("WalletCreateFundedPsbt",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(&btcjson.WalletCreateFundedPsbtResult{Fee: math.NaN()}, nil).Once()
	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	fee, err = rpc.EstimateTxFees(testnetAddress, entities.NewWei(1))
	require.Error(t, err)
	assert.Nil(t, fee)
}

func TestBitcoindWallet_GetBalance(t *testing.T) {
	absolutePath, err := filepath.Abs("../../../../test/mocks/listUnspent.json")
	require.NoError(t, err)
	rpcResponse, err := os.ReadFile(absolutePath)
	require.NoError(t, err)
	var result []btcjson.ListUnspentResult
	err = json.Unmarshal(rpcResponse, &result)
	require.NoError(t, err)
	client := &mocks.ClientAdapterMock{}
	client.On("ListUnspent").Return(result, nil).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	balance, err := rpc.GetBalance()
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(57962080000000000), balance)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_GetBalance_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("ListUnspent").Return(nil, assert.AnError).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	balance, err := rpc.GetBalance()
	require.Error(t, err)
	assert.Nil(t, balance)

	client.On("ListUnspent").Return([]btcjson.ListUnspentResult{{Amount: math.NaN(), Spendable: true}}, nil).Once()
	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	balance, err = rpc.GetBalance()
	require.Error(t, err)
	assert.Nil(t, balance)
}

func setupSendWithOpReturnTest(t *testing.T, client *mocks.ClientAdapterMock, encrypted bool) {
	var input []btcjson.TransactionInput
	var lockTime *int64
	satoshis := 50000000

	address, err := btcutil.DecodeAddress(testnetAddress, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	client.On("CreateRawTransaction",
		input,
		mock.MatchedBy(func(outputs map[btcutil.Address]btcutil.Amount) bool {
			for k, v := range outputs {
				require.Equal(t, address, k)
				require.Equal(t, btcutil.Amount(satoshis), v)
			}
			return len(outputs) == 1
		}),
		lockTime,
	).Return(&wire.MsgTx{
		Version:  0,
		TxIn:     nil,
		TxOut:    []*wire.TxOut{{Value: int64(satoshis), PkScript: []byte(paymentScriptMock)}},
		LockTime: 0,
	}, nil).Once()

	if encrypted {
		nonExpiredLockUntil := unexpiredTime
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
			UnlockedUntil: &nonExpiredLockUntil,
		}, nil).Once()
	}

	tx := &wire.MsgTx{
		Version: 0,
		TxIn:    nil,
		TxOut: []*wire.TxOut{
			{Value: int64(satoshis), PkScript: []byte(paymentScriptMock)},
			{Value: int64(0), PkScript: []byte{0x6a, 0x08, 0x02, 0x01, 0x00, 0x07, 0x02, 0x00, 0x00, 0x00}},
		},
		LockTime: 0,
	}

	changePos := 2
	feeRate := mockFeeRate
	var isWitness *bool
	client.On("FundRawTransaction", tx, btcjson.FundRawTransactionOpts{
		ChangePosition: &changePos,
		FeeRate:        &feeRate,
	}, isWitness).Return(&btcjson.FundRawTransactionResult{Transaction: tx, Fee: 0, ChangePosition: 2}, nil).Once()
	client.On("SignRawTransactionWithWallet", tx).Return(tx, true, nil).Once()
	client.On("SendRawTransaction", tx, false).Return(chainhash.NewHashFromStr(testnetTestTxHash)).Once()
}

func TestBitcoindWallet_SendWithOpReturn(t *testing.T) {
	data := []byte{2, 1, 0, 7, 2, 0, 0, 0}
	params := &chaincfg.TestNet3Params

	client := &mocks.ClientAdapterMock{}
	setupSendWithOpReturnTest(t, client, true)
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(params, client), mockAddress, mockFeeRate, true, mockPassword)
	txHash, err := rpc.SendWithOpReturn(testnetAddress, entities.NewWei(500000000000000000), data)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)
	assert.Equal(t, testnetTestTxHash, txHash)
	client.AssertExpectations(t)

	client = &mocks.ClientAdapterMock{}
	setupSendWithOpReturnTest(t, client, false)
	rpc = bitcoin.NewBitcoindWallet(bitcoin.NewConnection(params, client), mockAddress, mockFeeRate, false, mockPassword)
	txHash, err = rpc.SendWithOpReturn(testnetAddress, entities.NewWei(500000000000000000), data)
	require.NoError(t, err)
	assert.NotEmpty(t, txHash)
	assert.Equal(t, testnetTestTxHash, txHash)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_SendWithOpReturn_InvalidAddress(t *testing.T) {
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, &mocks.ClientAdapterMock{}), mockAddress, mockFeeRate, true, mockPassword)
	txHash, err := rpc.SendWithOpReturn(test.AnyString, entities.NewWei(500000000000000000), []byte{})
	require.Error(t, err)
	assert.Empty(t, txHash)
}

func TestBitcoindWallet_Address(t *testing.T) {
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, &mocks.ClientAdapterMock{}), mockAddress, mockFeeRate, true, mockPassword)
	assert.Equal(t, mockAddress, rpc.Address())
}

func TestBitcoindWallet_SendWithOpReturn_ErrorHandling(t *testing.T) {
	setups := sendWithOpReturnErrorSetups()
	for _, setup := range setups {
		client := &mocks.ClientAdapterMock{}
		data := []byte{2, 1, 0, 7, 2, 0, 0, 0}
		setup(client, &data)
		rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
		txHash, err := rpc.SendWithOpReturn(testnetAddress, entities.NewWei(500000000000000000), data)
		require.Error(t, err)
		assert.Empty(t, txHash)
		client.AssertExpectations(t)
	}
}

// nolint:funlen
func sendWithOpReturnErrorSetups() []func(client *mocks.ClientAdapterMock, data *[]byte) {
	return []func(client *mocks.ClientAdapterMock, data *[]byte){
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&wire.MsgTx{
				Version:  0,
				TxIn:     nil,
				TxOut:    []*wire.TxOut{{Value: int64(1), PkScript: []byte(paymentScriptMock)}},
				LockTime: 0,
			}, nil).Once()
			client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
		},
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&wire.MsgTx{
				Version:  0,
				TxIn:     nil,
				TxOut:    []*wire.TxOut{{Value: int64(1), PkScript: []byte(paymentScriptMock)}},
				LockTime: 0,
			}, nil).Once()
			for i := 0; i < txscript.MaxDataCarrierSize; i++ {
				*data = append(*data, byte(i))
			}
		},
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&wire.MsgTx{
				Version:  0,
				TxIn:     nil,
				TxOut:    []*wire.TxOut{{Value: int64(1), PkScript: []byte(paymentScriptMock)}},
				LockTime: 0,
			}, nil).Once()
			nonExpiredLockUntil := unexpiredTime
			client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
				UnlockedUntil: &nonExpiredLockUntil,
			}, nil).Once()
			var isWitness *bool
			client.On("FundRawTransaction", mock.Anything, mock.Anything, isWitness).Return(nil, assert.AnError).Once()
		},
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&wire.MsgTx{
				Version:  0,
				TxIn:     nil,
				TxOut:    []*wire.TxOut{{Value: int64(1), PkScript: []byte(paymentScriptMock)}},
				LockTime: 0,
			}, nil).Once()
			nonExpiredLockUntil := unexpiredTime
			client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
				UnlockedUntil: &nonExpiredLockUntil,
			}, nil).Once()
			var isWitness *bool
			client.On("FundRawTransaction", mock.Anything, mock.Anything, isWitness).Return(&btcjson.FundRawTransactionResult{}, nil).Once()
			client.On("SignRawTransactionWithWallet", mock.Anything).Return(nil, false, assert.AnError).Once()
		},
		func(client *mocks.ClientAdapterMock, data *[]byte) {
			client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&wire.MsgTx{
				Version:  0,
				TxIn:     nil,
				TxOut:    []*wire.TxOut{{Value: int64(1), PkScript: []byte(paymentScriptMock)}},
				LockTime: 0,
			}, nil).Once()
			nonExpiredLockUntil := unexpiredTime
			client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
				UnlockedUntil: &nonExpiredLockUntil,
			}, nil).Once()
			var isWitness *bool
			client.On("FundRawTransaction", mock.Anything, mock.Anything, isWitness).Return(&btcjson.FundRawTransactionResult{}, nil).Once()
			client.On("SignRawTransactionWithWallet", mock.Anything).Return(&wire.MsgTx{}, false, nil).Once()
			client.On("SendRawTransaction", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}

func TestBitcoindWallet_GetTransactions(t *testing.T) {
	absolutePath, err := filepath.Abs("../../../../test/mocks/listUnspentByAddress.json")
	require.NoError(t, err)
	rpcResponse, err := os.ReadFile(absolutePath)
	require.NoError(t, err)
	var result []btcjson.ListUnspentResult
	err = json.Unmarshal(rpcResponse, &result)
	require.NoError(t, err)
	client := &mocks.ClientAdapterMock{}
	parsedAddress, err := btcutil.DecodeAddress(testnetAddress, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	client.On("ListUnspentMinMaxAddresses", 0, 9999999, []btcutil.Address{parsedAddress}).Return(result, nil).Once()
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	transactions, err := rpc.GetTransactions(testnetAddress)
	require.NoError(t, err)
	slices.SortFunc(transactions, func(i, j blockchain.BitcoinTransactionInformation) int {
		return cmp.Compare(i.Hash, j.Hash)
	})
	assert.Equal(t, []blockchain.BitcoinTransactionInformation{
		{
			Hash:          "2ba6da53badd14349c5d6379e88c345e88193598aad714815d4b57c691a9fbdf",
			Confirmations: 2439,
			Outputs: map[string][]*entities.Wei{
				"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(2531000000000000)},
			},
		},
		{
			Hash:          "586c51dc94452aed9a373b0f52936c3e343c0db90f1155e985fd60e3c2e5c2b2",
			Confirmations: 6,
			Outputs: map[string][]*entities.Wei{
				"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(2000000000000000)},
			},
		},
		{
			Hash:          "da28401c76d618e8c3b1c3e15dfe1c10d4b24875f23768f30bcc26c99b9c82d4",
			Confirmations: 2,
			Outputs: map[string][]*entities.Wei{
				"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(200000000000000), entities.NewWei(1000000000000000), entities.NewWei(1000000000000000)},
			},
		},
		{
			Hash:          "fda421ccdff7324a382067d1746f6a387132435de6af336a0ebbf3f720eaae4d",
			Confirmations: 6,
			Outputs: map[string][]*entities.Wei{
				"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(20000000000000000)},
			},
		},
	}, transactions)
	client.AssertExpectations(t)
}

func TestBitcoindWallet_GetTransactions_ErrorHandling(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	rpc := bitcoin.NewBitcoindWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), mockAddress, mockFeeRate, true, mockPassword)
	transactions, err := rpc.GetTransactions("invalidAddress")
	require.Error(t, err)
	assert.Nil(t, transactions)

	client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	transactions, err = rpc.GetTransactions(testnetAddress)
	require.Error(t, err)
	assert.Nil(t, transactions)

	client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything).Return([]btcjson.ListUnspentResult{{Amount: math.NaN()}}, nil).Once()
	transactions, err = rpc.GetTransactions(testnetAddress)
	require.Error(t, err)
	assert.Nil(t, transactions)
}
