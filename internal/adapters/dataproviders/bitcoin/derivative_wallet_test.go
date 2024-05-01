package bitcoin_test

import (
	"cmp"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

const (
	btcAddress     = "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6"
	pubKey         = "0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22"
	changePosition = 2
	feeRate        = 0.0001
)

var (
	walletNotFoundErr     = errors.New("wallet not found")
	rawExistingAddress    = []byte("{\n  \"address\": \"n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6\",\n  \"scriptPubKey\": \"76a914ddb677f36498f7a4901a74e882df68fd00cf473588ac\",\n  \"ismine\": false,\n  \"solvable\": true,\n  \"desc\": \"pkh([ddb677f3]0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22)#ts3jjdae\",\n  \"iswatchonly\": true,\n  \"isscript\": false,\n  \"iswitness\": false,\n  \"pubkey\": \"0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22\",\n  \"iscompressed\": true,\n  \"ischange\": false,\n  \"timestamp\": 1,\n  \"labels\": [\n    \"\"\n  ]\n}")
	rawNonExistingAddress = []byte("{\n  \"address\": \"mjaGtyj74LYn7gApr17prZxDPDnfuUnRa5\",\n  \"scriptPubKey\": \"76a9142c81478132b5dda64ffc484a0d225096c4b22ad588ac\",\n  \"ismine\": false,\n  \"solvable\": false,\n  \"iswatchonly\": false,\n  \"isscript\": false,\n  \"iswitness\": false,\n  \"ischange\": false,\n  \"labels\": [\n  ]\n}")
)

func TestNewDerivativeWallet(t *testing.T) {
	existingAddressInfo := new(btcjson.GetAddressInfoResult)
	nonExistingAddressInfo := new(btcjson.GetAddressInfoResult)
	e := existingAddressInfo.UnmarshalJSON(rawExistingAddress)
	require.NoError(t, e)
	e = nonExistingAddressInfo.UnmarshalJSON(rawNonExistingAddress)
	require.NoError(t, e)

	rskAccount := test.OpenDerivativeWalletForTest(t, "derivative-wallet-creation")
	t.Run("Fail if doesn't have RSK wallet id", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewConnection(&chaincfg.TestNet3Params, client), &account.RskAccount{})
		require.ErrorContains(t, err, "derivative wallet can only be created with wallet id rsk-wallet")
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Fail if RSK account doesn't have derivation enabled", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), &account.RskAccount{})
		require.ErrorContains(t, err, "derivative wallet can only be used if RSK account has derivation enabled")
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Load wallet if not loaded", func(t *testing.T) { testLoadWallet(t, rskAccount, existingAddressInfo) })
	t.Run("Creates watch-only wallet if not exists", func(t *testing.T) { testCreateWatchOnlyWallet(t, rskAccount, existingAddressInfo) })
	t.Run("Imports pubkey if not imported and starts rescan", func(t *testing.T) { testImportPubKeyAndRescan(t, rskAccount, nonExistingAddressInfo) })
	t.Run("Returns error if wallet is scanning", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{
			WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: btcjson.ScanProgress{Duration: 5, Progress: 50}},
		}, nil).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.ErrorContains(t, err, "wallet is still scanning")
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Starts normally if wallet is created and key is imported", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
		client.On("GetAddressInfo", btcAddress).Return(existingAddressInfo, nil).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.NoError(t, err)
		assert.NotNil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Error handling", func(t *testing.T) { derivativeWalleCreationtErrorHandlingTests(t, rskAccount, nonExistingAddressInfo) })
}

func derivativeWalleCreationtErrorHandlingTests(t *testing.T, rskAccount *account.RskAccount, nonExistingAddressInfo *btcjson.GetAddressInfoResult) {
	t.Run("Error creating wallet", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
		client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, walletNotFoundErr).Once()
		client.On("CreateReadonlyWallet", mock.Anything).Return(assert.AnError).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.ErrorContains(t, err, "error while creating rsk-wallet wallet")
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Error getting address info", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
		client.On("GetAddressInfo", btcAddress).Return(nil, assert.AnError).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.Error(t, err)
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Error importing pubkey", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
		client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, walletNotFoundErr).Once()
		client.On("CreateReadonlyWallet", mock.Anything).Return(nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
		client.On("GetAddressInfo", btcAddress).Return(nonExistingAddressInfo, nil).Once()
		client.On("ImportPubKey", pubKey).Return(assert.AnError).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.Error(t, err)
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
	t.Run("Error starting rescan", func(t *testing.T) {
		client := &mocks.ClientAdapterMock{}
		client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
		client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, walletNotFoundErr).Once()
		client.On("CreateReadonlyWallet", mock.Anything).Return(nil).Once()
		client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
		client.On("GetAddressInfo", btcAddress).Return(nonExistingAddressInfo, nil).Once()
		client.On("ImportPubKey", pubKey).Return(nil).Once()
		client.On("ImportAddressRescan", btcAddress, "", true).Return(assert.AnError).Once()
		wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
		require.Error(t, err)
		assert.Nil(t, wallet)
		client.AssertExpectations(t)
	})
}

func testLoadWallet(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
	client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, nil).Once()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	assert.NotNil(t, wallet)
	client.AssertExpectations(t)
}

func testCreateWatchOnlyWallet(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
	client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, walletNotFoundErr).Once()
	client.On("CreateReadonlyWallet", btcclient.ReadonlyWalletRequest{
		WalletName: bitcoin.DerivativeWalletId, DisablePrivateKeys: true, Blank: true, AvoidReuse: false, Descriptors: false,
	}).Return(nil).Once()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	assert.NotNil(t, wallet)
	client.AssertExpectations(t)
}

func testImportPubKeyAndRescan(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(nil, walletNotFoundErr).Once()
	client.On("LoadWallet", bitcoin.DerivativeWalletId).Return(nil, walletNotFoundErr).Once()
	client.On("CreateReadonlyWallet", btcclient.ReadonlyWalletRequest{
		WalletName: bitcoin.DerivativeWalletId, DisablePrivateKeys: true, Blank: true, AvoidReuse: false, Descriptors: false,
	}).Return(nil).Once()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	client.On("ImportPubKey", pubKey).Return(nil).Once()
	client.On("ImportAddressRescan", btcAddress, "", true).Return(nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.ErrorContains(t, err, "public key imported, rescan started")
	assert.Nil(t, wallet)
	client.AssertExpectations(t)
}

func TestDerivativeWallet(t *testing.T) {
	existingAddressInfo := new(btcjson.GetAddressInfoResult)
	nonExistingAddressInfo := new(btcjson.GetAddressInfoResult)
	e := existingAddressInfo.UnmarshalJSON(rawExistingAddress)
	require.NoError(t, e)
	e = nonExistingAddressInfo.UnmarshalJSON(rawNonExistingAddress)
	require.NoError(t, e)
	rskAccount := test.OpenDerivativeWalletForTest(t, "derivative-wallet")

	t.Run("Address", func(t *testing.T) { testAddress(t, rskAccount, existingAddressInfo) })
	t.Run("Unlock", func(t *testing.T) { testUnlock(t, rskAccount, existingAddressInfo) })
	t.Run("ImportAddress", func(t *testing.T) { testImportAddress(t, rskAccount, existingAddressInfo) })

	t.Run("GetTransactions", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) { testGetTransactions(t, rskAccount, existingAddressInfo) })
	})

	t.Run("GetBalance", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) { testGetBalance(t, rskAccount, existingAddressInfo) })
	})

	t.Run("EstimateTxFees", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) { testEstimateFees(t, rskAccount, existingAddressInfo) })
		t.Run("Add extra value if estimation blocks is higher than the target", func(t *testing.T) { testEstimateFeesExtra(t, rskAccount, existingAddressInfo) })
		t.Run("Error handling", func(t *testing.T) {
			cases := derivativeWalletEstimateTxFeesErrorSetups(rskAccount)
			for _, testCase := range cases {
				client := &mocks.ClientAdapterMock{}
				client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
				client.On("GetAddressInfo", btcAddress).Return(existingAddressInfo, nil).Once()
				t.Run(testCase.description, func(t *testing.T) {
					testCase.setup(t, client)
				})
			}
		})
	})

	t.Run("SendWithOpReturn", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) { testSendWithOpReturn(t, rskAccount, existingAddressInfo) })
		t.Run("Error handling", func(t *testing.T) {
			cases := derivativeWalletSendWithOpReturnErrorSetups(rskAccount)
			for _, testCase := range cases {
				client := &mocks.ClientAdapterMock{}
				client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
				client.On("GetAddressInfo", btcAddress).Return(existingAddressInfo, nil).Once()
				t.Run(testCase.description, func(t *testing.T) {
					testCase.setup(t, client)
				})
			}
		})
	})
}

func testUnlock(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	err = wallet.Unlock()
	require.ErrorContains(t, err, "derivative wallet does not support unlocking as it is a watch-only wallet")
}

func testImportAddress(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	err = wallet.ImportAddress("n12ja1bZfZhpkxy8KHkQvj6rZM74kbhUWs")
	require.ErrorContains(t, err, "address importing is not supported in this type of wallet")
}

func testAddress(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	assert.Equal(t, btcAddress, wallet.Address())
}

func testGetBalance(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	parsedAddress, err := btcutil.DecodeAddress(btcAddress, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	client.On("ListUnspentMinMaxAddresses", bitcoin.MinConfirmationsForUtxos, bitcoin.MaxConfirmationsForUtxos, mock.MatchedBy(func(addresses []btcutil.Address) bool {
		return len(addresses) == 1 && addresses[0].EncodeAddress() == parsedAddress.EncodeAddress()
	})).Return([]btcjson.ListUnspentResult{
		{Amount: 50000000, Confirmations: 1},
		{Amount: 80000000, Confirmations: 3},
		{Amount: 50000000, Confirmations: 0},
		{Amount: 30000000, Confirmations: 2},
	}, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	balance, err := wallet.GetBalance()
	require.NoError(t, err)
	expected := new(big.Int)
	expected.SetString("160000000000000000000000000", 10)
	require.Equal(t, entities.NewBigWei(expected), balance)
	client.AssertExpectations(t)
}

func testGetTransactions(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	absolutePath, err := filepath.Abs("../../../../test/mocks/listUnspentByAddress.json")
	require.NoError(t, err)
	rpcResponse, err := os.ReadFile(absolutePath)
	require.NoError(t, err)
	var result []btcjson.ListUnspentResult
	err = json.Unmarshal(rpcResponse, &result)
	require.NoError(t, err)
	client := &mocks.ClientAdapterMock{}
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	parsedAddress, err := btcutil.DecodeAddress(testnetAddress, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	client.On("ListUnspentMinMaxAddresses", 0, 9999999, []btcutil.Address{parsedAddress}).Return(result, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	transactions, err := wallet.GetTransactions(testnetAddress)
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

func testEstimateFees(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	amount := entities.NewWei(5000000000000000)
	floatAmount, _ := amount.ToRbtc().Float64()
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	client.On("EstimateSmartFee", int64(1), &btcjson.EstimateModeConservative).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
	client.On("WalletCreateFundedPsbt",
		([]btcjson.PsbtInput)(nil),
		[]btcjson.PsbtOutput{
			{testnetAddress: floatAmount},
			{"data": "0000000000000000000000000000000000000000000000000000000000000000"},
		},
		(*uint32)(nil),
		&btcjson.WalletCreateFundedPsbtOpts{
			ChangeAddress:   btcjson.String(btcAddress),
			ChangePosition:  btcjson.Int64(changePosition),
			IncludeWatching: btcjson.Bool(true),
			FeeRate:         btcjson.Float64(feeRate),
		},
		(*bool)(nil),
	).Return(&btcjson.WalletCreateFundedPsbtResult{Fee: 0.0006}, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	fee, err := wallet.EstimateTxFees(testnetAddress, amount)
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(600000000000000), fee)
	client.AssertExpectations(t)
}

func testEstimateFeesExtra(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	amount := entities.NewWei(5000000000000000)
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	client.On("EstimateSmartFee", int64(1), &btcjson.EstimateModeConservative).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 2}, nil).Once()
	client.On("WalletCreateFundedPsbt",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&btcjson.WalletCreateFundedPsbtOpts{
			ChangeAddress:   btcjson.String(btcAddress),
			ChangePosition:  btcjson.Int64(changePosition),
			IncludeWatching: btcjson.Bool(true),
			FeeRate:         btcjson.Float64(0.00011),
		},
		mock.Anything,
	).Return(&btcjson.WalletCreateFundedPsbtResult{Fee: 0.0006}, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	fee, err := wallet.EstimateTxFees(testnetAddress, amount)
	require.NoError(t, err)
	assert.Equal(t, entities.NewWei(600000000000000), fee)
	client.AssertExpectations(t)
}

func testSendWithOpReturn(t *testing.T, rskAccount *account.RskAccount, addressInfo *btcjson.GetAddressInfoResult) {
	client := &mocks.ClientAdapterMock{}
	value := entities.NewWei(600000000000000000)
	satoshis, _ := value.ToSatoshi().Float64()
	address, err := btcutil.DecodeAddress(testnetAddress, &chaincfg.TestNet3Params)
	require.NoError(t, err)
	client.On("GetWalletInfo").Return(&btcjson.GetWalletInfoResult{WalletName: bitcoin.DerivativeWalletId, Scanning: btcjson.ScanningOrFalse{Value: false}}, nil).Once()
	client.On("GetAddressInfo", btcAddress).Return(addressInfo, nil).Once()
	client.On("CreateRawTransaction",
		([]btcjson.TransactionInput)(nil),
		mock.MatchedBy(func(outputs map[btcutil.Address]btcutil.Amount) bool {
			for k, v := range outputs {
				require.Equal(t, address, k)
				require.Equal(t, btcutil.Amount(satoshis), v)
			}
			return len(outputs) == 1
		}),
		(*int64)(nil),
	).Return(&wire.MsgTx{
		Version:  0,
		TxIn:     nil,
		TxOut:    []*wire.TxOut{{Value: int64(satoshis), PkScript: []byte(paymentScriptMock)}},
		LockTime: 0,
	}, nil).Once()
	tx := &wire.MsgTx{
		Version: 0,
		TxIn:    nil,
		TxOut: []*wire.TxOut{
			{Value: int64(satoshis), PkScript: []byte(paymentScriptMock)},
			{Value: int64(0), PkScript: []byte{0x6a, 0x05, 0xf1, 0xf2, 0xf3, 0xf4, 0x00}},
		},
		LockTime: 0,
	}
	client.On("EstimateSmartFee", int64(1), &btcjson.EstimateModeConservative).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
	client.On("FundRawTransaction", tx, btcjson.FundRawTransactionOpts{
		ChangeAddress:   btcjson.String(btcAddress),
		ChangePosition:  btcjson.Int(changePosition),
		IncludeWatching: btcjson.Bool(true),
		LockUnspents:    btcjson.Bool(true),
		FeeRate:         btcjson.Float64(feeRate),
		Replaceable:     btcjson.Bool(true),
	}, (*bool)(nil)).Return(&btcjson.FundRawTransactionResult{Transaction: tx, Fee: 50, ChangePosition: 2}, nil).Once()
	client.On("SendRawTransaction", tx, false).Return(chainhash.NewHashFromStr(testnetTestTxHash)).Once()
	client.On("SignRawTransactionWithKey", tx, mock.MatchedBy(func(pks []string) bool {
		return len(pks) == 1 && pks[0] != ""
	})).Return(tx, true, nil).Once()
	wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
	require.NoError(t, err)
	result, err := wallet.SendWithOpReturn(testnetAddress, value, []byte{0xf1, 0xf2, 0xf3, 0xf4, 0x00})
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	client.AssertExpectations(t)
}

// nolint:funlen
func derivativeWalletSendWithOpReturnErrorSetups(rskAccount *account.RskAccount) []struct {
	description string
	setup       func(t *testing.T, client *mocks.ClientAdapterMock)
} {
	rawTx := &wire.MsgTx{TxOut: []*wire.TxOut{{Value: int64(50000000), PkScript: []byte(paymentScriptMock)}}}
	return []struct {
		description string
		setup       func(t *testing.T, client *mocks.ClientAdapterMock)
	}{
		{
			description: "error parsing address",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(test.AnyString, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error creating raw tx",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error estimating fees",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error estimating fees (RPC error)",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{
					Errors: []string{assert.AnError.Error()},
				}, nil).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error funding raw tx",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
				client.On("FundRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error funding raw tx",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
				client.On("FundRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error signing tx",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
				client.On("FundRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&btcjson.FundRawTransactionResult{Transaction: rawTx, Fee: 50, ChangePosition: 2}, nil).Once()
				client.On("SignRawTransactionWithKey", mock.Anything, mock.Anything).Return(nil, false, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error sending tx",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
				client.On("FundRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&btcjson.FundRawTransactionResult{Transaction: rawTx, Fee: 50, ChangePosition: 2}, nil).Once()
				client.On("SignRawTransactionWithKey", mock.Anything, mock.Anything).Return(rawTx, true, nil).Once()
				client.On("SendRawTransaction", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error sending tx (incomplete signatures)",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("CreateRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(rawTx, nil).Once()
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(feeRate), Blocks: 1}, nil).Once()
				client.On("FundRawTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&btcjson.FundRawTransactionResult{Transaction: rawTx, Fee: 50, ChangePosition: 2}, nil).Once()
				client.On("SignRawTransactionWithKey", mock.Anything, mock.Anything).Return(rawTx, false, nil).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.SendWithOpReturn(testnetAddress, entities.NewWei(1), []byte{0xf1})
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
	}
}

func derivativeWalletEstimateTxFeesErrorSetups(rskAccount *account.RskAccount) []struct {
	description string
	setup       func(t *testing.T, client *mocks.ClientAdapterMock)
} {
	return []struct {
		description string
		setup       func(t *testing.T, client *mocks.ClientAdapterMock)
	}{
		{
			description: "estimate for invalid address",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.EstimateTxFees(test.AnyString, entities.NewWei(1))
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error when getting estimation",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.EstimateTxFees(testnetAddress, entities.NewWei(1))
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error when getting estimation (RPC error)",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{
					Errors: []string{assert.AnError.Error()},
				}, nil).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.EstimateTxFees(testnetAddress, entities.NewWei(1))
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
		{
			description: "error when funding psbt",
			setup: func(t *testing.T, client *mocks.ClientAdapterMock) {
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{FeeRate: btcjson.Float64(0.001), Blocks: 1}, nil).Once()
				client.On("WalletCreateFundedPsbt", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
				wallet, err := bitcoin.NewDerivativeWallet(bitcoin.NewWalletConnection(&chaincfg.TestNet3Params, client, bitcoin.DerivativeWalletId), rskAccount)
				require.NoError(t, err)
				result, err := wallet.EstimateTxFees(testnetAddress, entities.NewWei(1))
				require.Error(t, err)
				assert.Empty(t, result)
				client.AssertExpectations(t)
			},
		},
	}
}
