package rootstock_test

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var dummyClient = rootstock.NewRskClient(nil)

func TestRskBridgeImpl_GetAddress(t *testing.T) {
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{Address: test.AnyAddress}, nil, dummyClient, nil, rootstock.RetryParams{}, nil)
	assert.Equal(t, test.AnyAddress, bridge.GetAddress())
}

func TestRskBridgeImpl_GetRequiredTxConfirmations(t *testing.T) {
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{RequiredConfirmations: 10}, nil, dummyClient, nil, rootstock.RetryParams{}, nil)
	assert.Equal(t, uint64(10), bridge.GetRequiredTxConfirmations())
}

func TestRskBridgeImpl_GetFedAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil)
		result, err := bridge.GetFedAddress()
		assert.Equal(t, test.AnyAddress, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetFederationAddress call fail", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetFederationAddress", mock.Anything).Return("", assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil)
		result, err := bridge.GetFedAddress()
		assert.Empty(t, result)
		require.Error(t, err)
	})
}

func TestRskBridgeImpl_GetMinimumLockTxValue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetMinimumLockTxValue", mock.Anything).Return(big.NewInt(5), nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil)
		result, err := bridge.GetMinimumLockTxValue()
		assert.IsType(t, &entities.Wei{}, result)
		assert.Equal(t, entities.NewWei(50000000000), result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinimumLockTxValue call fail", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetMinimumLockTxValue", mock.Anything).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil)
		result, err := bridge.GetMinimumLockTxValue()
		assert.Nil(t, result)
		require.Error(t, err)
	})
}

func TestRskBridgeImpl_GetFlyoverDerivationAddress(t *testing.T) {
	const redeemScriptString = "522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae"
	lbcAddress, err := hex.DecodeString("2ff74F841b95E000625b3A77fed03714874C4fEa")
	require.NoError(t, err)
	quoteHash, err := hex.DecodeString("4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0")
	require.NoError(t, err)
	userBtcAddress, err := bitcoin.DecodeAddressBase58("mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk", true)
	require.NoError(t, err)
	lpBtcAddress, err := bitcoin.DecodeAddressBase58("mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk", true)
	require.NoError(t, err)
	args := blockchain.FlyoverDerivationArgs{
		FedInfo:              mocks.GetFakeFedInfo(),
		LbcAdress:            lbcAddress,
		UserBtcRefundAddress: userBtcAddress,
		LpBtcAddress:         lpBtcAddress,
		QuoteHash:            quoteHash,
	}
	args.FedInfo.FedAddress = "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p"
	t.Run("Success", func(t *testing.T) {
		var testError error
		var redeemScriptBytes []byte
		var result blockchain.FlyoverDerivation
		bridgeMock := &mocks.RskBridgeBindingMock{}
		redeemScriptBytes, testError = hex.DecodeString(redeemScriptString)
		require.NoError(t, testError)
		bridgeMock.On("GetActivePowpegRedeemScript", mock.Anything).Return(redeemScriptBytes, nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil)
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Equal(t, blockchain.FlyoverDerivation{
			Address:      "2Mx7jaPHtsgJTbqGnjU5UqBpkekHgfigXay",
			RedeemScript: "20ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c7275522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae",
		}, result)
		require.NoError(t, testError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetActivePowpegRedeemScript call fail", func(t *testing.T) {
		var testError error
		var result blockchain.FlyoverDerivation
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetActivePowpegRedeemScript", mock.Anything).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil)
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Empty(t, result)
		require.ErrorContains(t, testError, "error retreiving fed redeem script from bridge")
	})
}

func TestRskBridgeImpl_FetchFederationInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(2), nil).Once()
		bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
		bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(1), "btc").Return([]byte{0x0a, 0x0b, 0x0c}, nil).Once()
		bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
		bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil).Once()
		bridgeMock.On("GetActiveFederationCreationBlockHeight", mock.Anything).Return(big.NewInt(500), nil).Once()

		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100, ErpKeys: []string{"key1", "key2", "key3"}},
			bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil)
		fedInfo, err := bridge.FetchFederationInfo()
		require.NoError(t, err)
		assert.Equal(t, blockchain.FederationInfo{
			FedSize:              2,
			FedThreshold:         5,
			FedAddress:           test.AnyAddress,
			PubKeys:              []string{"010203", "0a0b0c"},
			ActiveFedBlockHeight: 500,
			IrisActivationHeight: 100,
			ErpKeys:              []string{"key1", "key2", "key3"},
		}, fedInfo)
		bridgeMock.AssertExpectations(t)
	})

	t.Run("Error handling", func(t *testing.T) {
		for _, setUp := range fetchFedInfoErrorSetUps() {
			bridgeMock := &mocks.RskBridgeBindingMock{}
			setUp(bridgeMock)
			bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{IrisActivationHeight: 100}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil)
			result, err := bridge.FetchFederationInfo()
			require.Error(t, err)
			assert.Empty(t, result)
			bridgeMock.AssertExpectations(t)
		}
	})
}

func fetchFedInfoErrorSetUps() []func(bridgeMock *mocks.RskBridgeBindingMock) {
	return []func(bridgeMock *mocks.RskBridgeBindingMock){
		func(bridgeMock *mocks.RskBridgeBindingMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeBindingMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeBindingMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeBindingMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.On("GetFederationAddress", mock.Anything).Return("", assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeBindingMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil).Once()
			bridgeMock.On("GetActiveFederationCreationBlockHeight", mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}

// nolint:funlen
func TestRskBridgeImpl_RegisterBtcCoinbaseTransaction(t *testing.T) {
	signerMock := &mocks.TransactionSignerMock{}
	signerMock.On("Address").Return(parsedAddress)
	coinbaseInfo := blockchain.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      utils.MustGetRandomBytes(32),
		BlockHash:            utils.To32Bytes(utils.MustGetRandomBytes(32)),
		BlockHeight:          big.NewInt(500),
		SerializedPmt:        utils.MustGetRandomBytes(64),
		WitnessMerkleRoot:    utils.To32Bytes(utils.MustGetRandomBytes(32)),
		WitnessReservedValue: utils.To32Bytes(utils.MustGetRandomBytes(32)),
	}
	t.Run("Should handle error getting best chain height", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should return WaitingForBridgeError if block is higher than best chain height", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(300), nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when validating if tx was registered", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should skip registration if tx was already registered", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(true, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when sending tx", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(mock.Anything, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error")
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when tx reverts", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		mockClient := &mocks.RpcClientBindingMock{}
		tx := prepareTxMocks(mockClient, signerMock, false)
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(mock.Anything, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(tx, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error: transaction reverted (0xfe6fc232343284368505aa7bad1ccdd865498df6e6691b53e128f14e5e21bb74)")
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should register coinbase transaction successfully", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeBindingMock{}
		mockClient := &mocks.RpcClientBindingMock{}
		tx := prepareTxMocks(mockClient, signerMock, true)
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		matchFunc := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
			return opts.From == parsedAddress && opts.GasLimit == 100000
		})
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(matchFunc, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(tx, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock)
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
}
