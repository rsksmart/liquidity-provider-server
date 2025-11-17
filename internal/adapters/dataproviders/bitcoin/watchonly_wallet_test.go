package bitcoin_test

import (
	"cmp"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var getTransactionsExpectedResult = []blockchain.BitcoinTransactionInformation{
	{
		Hash:          "2ba6da53badd14349c5d6379e88c345e88193598aad714815d4b57c691a9fbdf",
		Confirmations: 288434,
		Outputs: map[string][]*entities.Wei{
			"mqbKtarYKnoEdPheFFDGRjksvEpb2vJGNh": {entities.NewWei(6000000000000000)},
			"":                                   {entities.NewWei(0)},
			"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(2531000000000000)},
		},
		HasWitness: false,
	},
	{
		Hash:          "586c51dc94452aed9a373b0f52936c3e343c0db90f1155e985fd60e3c2e5c2b2",
		Confirmations: 286001,
		Outputs: map[string][]*entities.Wei{
			"mxcLm8hdhfJ1cutzeq6zdwcUohKVfmRhPu": {entities.NewWei(992000000000000)},
			"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(2000000000000000)},
		},
		HasWitness: false,
	},
	{
		Hash:          "da28401c76d618e8c3b1c3e15dfe1c10d4b24875f23768f30bcc26c99b9c82d4",
		Confirmations: 285997,
		Outputs: map[string][]*entities.Wei{
			"mocAPSv6trAJoZRoqcn18kvXEjcxvXc9m5": {entities.NewWei(93000000000000)},
			"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(200000000000000), entities.NewWei(1000000000000000), entities.NewWei(1000000000000000)},
		},
		HasWitness: false,
	},
	{
		Hash:          "fda421ccdff7324a382067d1746f6a387132435de6af336a0ebbf3f720eaae4d",
		Confirmations: 286001,
		Outputs: map[string][]*entities.Wei{
			"n1sSgnWcHU8AeHTVFez9RQ8HxMdAVHJXui": {entities.NewWei(1873000000000000)},
			"n3HJbF1Ps5c9ZE3UvLyjGFDvyAfjzDEBkS": {entities.NewWei(20000000000000000)},
		},
		HasWitness: false,
	},
}

func TestNewWatchOnlyWallet(t *testing.T) {
	t.Run("wallet already created and loaded", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		require.NotNil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("load created wallet", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
		client.On("LoadWallet", bitcoin.PeginWalletId).Return(nil, nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		require.NotNil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("create new wallet", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
		client.On("LoadWallet", bitcoin.PeginWalletId).Return(nil, assert.AnError).Once()
		params := btcclient.ReadonlyWalletRequest{WalletName: bitcoin.PeginWalletId, DisablePrivateKeys: true, Blank: true, AvoidReuse: true, Descriptors: false}
		client.On("CreateReadonlyWallet", params).Return(nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		require.NotNil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("wallet is not watch only", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: true}, nil).Once()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.ErrorContains(t, err, "wallet is not watch-only")
		require.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("handle RPC errors", func(t *testing.T) {
		t.Run("on wallet create", func(t *testing.T) {
			client := &mocks.ClientAdapterMock{}
			client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
			client.On("LoadWallet", bitcoin.PeginWalletId).Return(nil, assert.AnError).Once()
			client.On("CreateReadonlyWallet", mock.Anything).Return(errors.New("creation error")).Once()
			wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
			require.ErrorContains(t, err, "error creating watch-only wallet: creation error")
			require.Nil(t, wallet)
			client.AssertExpectations(t)
		})
		t.Run("on get wallet info", func(t *testing.T) {
			client := &mocks.ClientAdapterMock{}
			client.On("GetWalletInfo").Return(nil, assert.AnError).Once()
			client.On("LoadWallet", bitcoin.PeginWalletId).Return(nil, nil).Once()
			client.On("GetWalletInfo", mock.Anything).Return(nil, errors.New("info error")).Once()
			wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
			require.ErrorContains(t, err, "error creating watch-only wallet: info error")
			require.Nil(t, wallet)
			client.AssertExpectations(t)
		})
	})
}

func TestWatchOnlyWallet_GetBalance(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	result, err := wallet.GetBalance()
	require.ErrorContains(t, err, "cannot get balance of a watch-only wallet")
	require.Nil(t, result)
}

func TestWatchOnlyWallet_Shutdown(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("Disconnect").Return().Once()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	shutdownChannel := make(chan bool, 1)
	wallet.Shutdown(shutdownChannel)
	<-shutdownChannel
	client.AssertExpectations(t)
}

func TestWatchOnlyWallet_Address(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	result := wallet.Address()
	require.Empty(t, result)
}

func TestWatchOnlyWallet_EstimateTxFees(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	result, err := wallet.EstimateTxFees("address", nil)
	require.ErrorContains(t, err, "cannot estimate from a watch-only wallet")
	require.Empty(t, result)
}

// TestWatchOnlyWallet_GetTransactions This test are reused from the bitcoind wallet tests suite since they share behavior
// nolint:funlen
func TestWatchOnlyWallet_GetTransactions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		absolutePathListUnspent, err := filepath.Abs("../../../../test/mocks/listUnspentByAddress.json")
		require.NoError(t, err)
		listUnspentRpcResponse, err := os.ReadFile(absolutePathListUnspent)
		require.NoError(t, err)
		var result []btcjson.ListUnspentResult
		err = json.Unmarshal(listUnspentRpcResponse, &result)
		require.NoError(t, err)

		var absolutePathTx string
		var txRpcResponse []byte
		var txId *chainhash.Hash
		client := &mocks.ClientAdapterMock{}
		mockedTxs := make(map[chainhash.Hash]*btcjson.TxRawResult)
		for _, utxo := range result {
			txResult := new(btcjson.TxRawResult)
			absolutePathTx, err = filepath.Abs("../../../../test/mocks/rawTxVerbose-" + utxo.TxID + ".json")
			require.NoError(t, err)
			txRpcResponse, err = os.ReadFile(absolutePathTx)
			require.NoError(t, err)
			err = json.Unmarshal(txRpcResponse, txResult)
			require.NoError(t, err)
			txId, err = chainhash.NewHashFromStr(utxo.TxID)
			require.NoError(t, err)
			mockedTxs[*txId] = txResult
		}
		client.EXPECT().GetRawTransactionVerbose(mock.Anything).RunAndReturn(func(hash *chainhash.Hash) (*btcjson.TxRawResult, error) {
			return mockedTxs[*hash], nil
		})

		parsedAddress, err := btcutil.DecodeAddress(testnetAddress, &chaincfg.TestNet3Params)
		require.NoError(t, err)
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Twice()
		client.On("ListUnspentMinMaxAddresses", 0, 9999999, []btcutil.Address{parsedAddress}).Return(result, nil).Once()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		transactions, err := wallet.GetTransactions(testnetAddress)
		require.NoError(t, err)
		slices.SortFunc(transactions, func(i, j blockchain.BitcoinTransactionInformation) int {
			return cmp.Compare(i.Hash, j.Hash)
		})
		assert.Equal(t, getTransactionsExpectedResult, transactions)
		client.AssertExpectations(t)
	})
	t.Run("Error on RPC call", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Times(5)
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		transactions, err := wallet.GetTransactions("invalidAddress")
		require.Error(t, err)
		assert.Nil(t, transactions)

		client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		transactions, err = wallet.GetTransactions(testnetAddress)
		require.Error(t, err)
		assert.Nil(t, transactions)

		client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything).Return([]btcjson.ListUnspentResult{
			{TxID: hex.EncodeToString(utils.MustGetRandomBytes(32))},
		}, nil).Once()
		client.On("GetRawTransactionVerbose", mock.Anything).Return(nil, assert.AnError).Once()
		transactions, err = wallet.GetTransactions(testnetAddress)
		require.Error(t, err)
		assert.Nil(t, transactions)

		client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything).Return([]btcjson.ListUnspentResult{
			{TxID: hex.EncodeToString(utils.MustGetRandomBytes(32))},
		}, nil).Once()
		client.On("GetRawTransactionVerbose", mock.Anything).Return(&btcjson.TxRawResult{
			Vout: []btcjson.Vout{{Value: math.NaN()}},
		}, nil).Once()
		transactions, err = wallet.GetTransactions(testnetAddress)
		require.Error(t, err)
		assert.Nil(t, transactions)
	})
}

// TestWatchOnlyWallet_ImportAddress This test are reused from the bitcoind wallet tests suite since they share behavior
func TestWatchOnlyWallet_ImportAddress(t *testing.T) {
	t.Run("valid address", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("ImportAddressRescan", testnetAddress, "", false).Return(nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Twice()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		err = wallet.ImportAddress(testnetAddress)
		require.NoError(t, err)
		client.AssertExpectations(t)

		client = &mocks.ClientAdapterMock{}
		client.On("ImportAddressRescan", mainnetAddress, "", false).Return(nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Twice()
		wallet, err = bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.MainNetParams, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		err = wallet.ImportAddress(mainnetAddress)
		require.NoError(t, err)
		client.AssertExpectations(t)
	})
	t.Run("invalid address", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Twice()
		wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.MainNetParams, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		err = wallet.ImportAddress(testnetAddress)
		require.Error(t, err)

		wallet, err = bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
		require.NoError(t, err)
		err = wallet.ImportAddress(mainnetAddress)
		require.Error(t, err)
	})
}

func TestWatchOnlyWallet_Unlock(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	err = wallet.Unlock()
	require.ErrorContains(t, err, "watch-only wallet does not support unlocking as it only has monitoring purposes")
}

func TestWatchOnlyWallet_SendWithOpReturn(t *testing.T) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{PrivateKeysEnabled: false}, nil).Once()
	wallet, err := bitcoin.NewWatchOnlyWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.PeginWalletId))
	require.NoError(t, err)
	result, err := wallet.SendWithOpReturn("address", nil, nil)
	require.ErrorContains(t, err, "cannot send from a watch-only wallet")
	require.Empty(t, result.Hash)
	require.Nil(t, result.Fee)
}
