package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	rsk "github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var dummyClient = rootstock.NewRskClient(nil)

const (
	batchPegOutIteratorString = "*bindings.RskBridgeBatchPegoutCreatedIterator"
)

func TestNewRskBridgeImpl(t *testing.T) {
	config := rootstock.RskBridgeConfig{Address: test.AnyAddress, RequiredConfirmations: 10, ErpKeys: []string{"key1", "key2", "key3"}, UseSegwitFederation: true}
	client := rootstock.NewRskClient(&mocks.RpcClientBindingMock{})
	bridge := rootstock.NewRskBridgeImpl(config, &mocks.RskBridgeAdapterMock{}, client, &chaincfg.TestNet3Params, rootstock.RetryParams{Retries: 1, Sleep: time.Duration(1)}, &mocks.TransactionSignerMock{}, time.Duration(1))
	test.AssertNonZeroValues(t, bridge)
}

func TestRskBridgeImpl_GetAddress(t *testing.T) {
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{Address: test.AnyAddress}, nil, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
	assert.Equal(t, test.AnyAddress, bridge.GetAddress())
}

func TestRskBridgeImpl_GetRequiredTxConfirmations(t *testing.T) {
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{RequiredConfirmations: 10}, nil, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
	assert.Equal(t, uint64(10), bridge.GetRequiredTxConfirmations())
}

func TestRskBridgeImpl_GetFedAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
		result, err := bridge.GetFedAddress()
		assert.Equal(t, test.AnyAddress, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetFederationAddress call fail", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetFederationAddress", mock.Anything).Return("", assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
		result, err := bridge.GetFedAddress()
		assert.Empty(t, result)
		require.Error(t, err)
	})
}

func TestRskBridgeImpl_GetMinimumLockTxValue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetMinimumLockTxValue", mock.Anything).Return(big.NewInt(5), nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
		result, err := bridge.GetMinimumLockTxValue()
		assert.IsType(t, &entities.Wei{}, result)
		assert.Equal(t, entities.NewWei(50000000000), result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinimumLockTxValue call fail", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetMinimumLockTxValue", mock.Anything).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, nil, rootstock.RetryParams{}, nil, time.Duration(1))
		result, err := bridge.GetMinimumLockTxValue()
		assert.Nil(t, result)
		require.Error(t, err)
	})
}

func TestRskBridgeImpl_GetFlyoverDerivationAddress(t *testing.T) {
	const redeemScriptString = "64522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae670350cd00b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec8054ae68"
	lbcAddress, err := hex.DecodeString("2ff74F841b95E000625b3A77fed03714874C4fEa")
	require.NoError(t, err)
	quoteHash, err := hex.DecodeString("4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0")
	require.NoError(t, err)
	userBtcAddress, err := bitcoin.DecodeAddressBase58("mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk", true)
	require.NoError(t, err)
	lpBtcAddress, err := bitcoin.DecodeAddressBase58("mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk", true)
	require.NoError(t, err)
	args := rsk.FlyoverDerivationArgs{
		FedInfo:              mocks.GetFakeFedInfo(),
		LbcAdress:            lbcAddress,
		UserBtcRefundAddress: userBtcAddress,
		LpBtcAddress:         lpBtcAddress,
		QuoteHash:            quoteHash,
	}
	args.FedInfo.FedAddress = "2NCxHG5oK8CWLDrBpTQq6pgKE8jyoB2DpTe"
	t.Run("Success", func(t *testing.T) {
		var testError error
		var redeemScriptBytes []byte
		var result rsk.FlyoverDerivation
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		redeemScriptBytes, testError = hex.DecodeString(redeemScriptString)
		require.NoError(t, testError)
		bridgeMock.On("GetActivePowpegRedeemScript", mock.Anything).Return(redeemScriptBytes, nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Equal(t, rsk.FlyoverDerivation{
			Address:      "2MxeEHVx71taCeVsXFsfQ7TKK6v943PFVEu",
			RedeemScript: "20ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c727564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae670350cd00b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec8054ae68",
		}, result)
		require.NoError(t, testError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on GetActivePowpegRedeemScript call fail", func(t *testing.T) {
		var testError error
		var result rsk.FlyoverDerivation
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetActivePowpegRedeemScript", mock.Anything).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Empty(t, result)
		require.ErrorContains(t, testError, "error retreiving fed redeem script from bridge")
	})
}

// nolint:funlen
func TestRskBridgeImpl_FetchFederationInfo(t *testing.T) {
	t.Run("Success for segwit federation", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(2), nil).Once()
		bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
		bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(1), "btc").Return([]byte{0x0a, 0x0b, 0x0c}, nil).Once()
		bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
		bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil).Once()
		bridgeMock.On("GetActiveFederationCreationBlockHeight", mock.Anything).Return(big.NewInt(500), nil).Once()

		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{ErpKeys: []string{"key1", "key2", "key3"}, UseSegwitFederation: true},
			bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
		fedInfo, err := bridge.FetchFederationInfo()
		require.NoError(t, err)
		assert.Equal(t, rsk.FederationInfo{
			FedSize:              2,
			FedThreshold:         5,
			FedAddress:           test.AnyAddress,
			PubKeys:              []string{"010203", "0a0b0c"},
			ActiveFedBlockHeight: 500,
			ErpKeys:              []string{"key1", "key2", "key3"},
			UseSegwit:            true,
		}, fedInfo)
		bridgeMock.AssertExpectations(t)
	})

	t.Run("Success for legacy federation", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(big.NewInt(2), nil).Once()
		bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
		bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(1), "btc").Return([]byte{0x0a, 0x0b, 0x0c}, nil).Once()
		bridgeMock.EXPECT().GetRetiringFederationThreshold(mock.Anything).Return(big.NewInt(5), nil).Once()
		bridgeMock.EXPECT().GetRetiringFederationAddress(mock.Anything).Return(test.AnyAddress, nil).Once()
		bridgeMock.EXPECT().GetRetiringFederationCreationBlockNumber(mock.Anything).Return(big.NewInt(500), nil).Once()

		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{ErpKeys: []string{"key1", "key2", "key3"}, UseSegwitFederation: false},
			bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
		fedInfo, err := bridge.FetchFederationInfo()
		require.NoError(t, err)
		assert.Equal(t, rsk.FederationInfo{
			FedSize:              2,
			FedThreshold:         5,
			FedAddress:           test.AnyAddress,
			PubKeys:              []string{"010203", "0a0b0c"},
			ActiveFedBlockHeight: 500,
			ErpKeys:              []string{"key1", "key2", "key3"},
			UseSegwit:            false,
		}, fedInfo)
		bridgeMock.AssertExpectations(t)
	})

	t.Run("Error handling segwit federation", func(t *testing.T) {
		for _, setUp := range fetchFedInfoErrorSetUps() {
			bridgeMock := &mocks.RskBridgeAdapterMock{}
			setUp(bridgeMock)
			bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{UseSegwitFederation: true}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
			result, err := bridge.FetchFederationInfo()
			require.Error(t, err)
			assert.Empty(t, result)
			bridgeMock.AssertExpectations(t)
		}
	})

	t.Run("Error handling legacy federation", func(t *testing.T) {
		for _, setUp := range fetchLegacyFedInfoErrorSetUps() {
			bridgeMock := &mocks.RskBridgeAdapterMock{}
			setUp(bridgeMock)
			bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{UseSegwitFederation: false}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, time.Duration(1))
			result, err := bridge.FetchFederationInfo()
			require.Error(t, err)
			assert.Empty(t, result)
			bridgeMock.AssertExpectations(t)
		}
	})
}

func fetchFedInfoErrorSetUps() []func(bridgeMock *mocks.RskBridgeAdapterMock) {
	return []func(bridgeMock *mocks.RskBridgeAdapterMock){
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.On("GetFederationAddress", mock.Anything).Return("", assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.On("GetFederationSize", mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.On("GetFederatorPublicKeyOfType", mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.On("GetFederationThreshold", mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.On("GetFederationAddress", mock.Anything).Return(test.AnyAddress, nil).Once()
			bridgeMock.On("GetActiveFederationCreationBlockHeight", mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}

func fetchLegacyFedInfoErrorSetUps() []func(bridgeMock *mocks.RskBridgeAdapterMock) {
	return []func(bridgeMock *mocks.RskBridgeAdapterMock){
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(0), "btc").Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationThreshold(mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationThreshold(mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationAddress(mock.Anything).Return("", assert.AnError).Once()
		},
		func(bridgeMock *mocks.RskBridgeAdapterMock) {
			bridgeMock.EXPECT().GetRetiringFederationSize(mock.Anything).Return(big.NewInt(1), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederatorPublicKeyOfType(mock.Anything, big.NewInt(0), "btc").Return([]byte{0x01, 0x02, 0x03}, nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationThreshold(mock.Anything).Return(big.NewInt(5), nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationAddress(mock.Anything).Return(test.AnyAddress, nil).Once()
			bridgeMock.EXPECT().GetRetiringFederationCreationBlockNumber(mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}

// nolint:funlen
func TestRskBridgeImpl_RegisterBtcCoinbaseTransaction(t *testing.T) {
	signerMock := &mocks.TransactionSignerMock{}
	signerMock.On("Address").Return(parsedAddress)
	coinbaseInfo := rsk.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      utils.MustGetRandomBytes(32),
		BlockHash:            utils.To32Bytes(utils.MustGetRandomBytes(32)),
		BlockHeight:          big.NewInt(500),
		SerializedPmt:        utils.MustGetRandomBytes(64),
		WitnessMerkleRoot:    utils.To32Bytes(utils.MustGetRandomBytes(32)),
		WitnessReservedValue: utils.To32Bytes(utils.MustGetRandomBytes(32)),
	}
	t.Run("Should handle error getting best chain height", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should return WaitingForBridgeError if block is higher than best chain height", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(300), nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when validating if tx was registered", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should skip registration if tx was already registered", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(true, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when sending tx", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(mock.Anything, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error")
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should handle error when tx reverts", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		mockClient := &mocks.RpcClientBindingMock{}
		tx := prepareTxMocks(mockClient, signerMock, false)
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(mock.Anything, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(tx, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error: transaction reverted (0xfe6fc232343284368505aa7bad1ccdd865498df6e6691b53e128f14e5e21bb74)")
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Should register coinbase transaction successfully", func(t *testing.T) {
		bridgeMock := &mocks.RskBridgeAdapterMock{}
		mockClient := &mocks.RpcClientBindingMock{}
		tx := prepareTxMocks(mockClient, signerMock, true)
		bridgeMock.EXPECT().GetBtcBlockchainBestChainHeight(mock.Anything).Return(big.NewInt(600), nil).Once()
		bridgeMock.EXPECT().HasBtcBlockCoinbaseTransactionInformation(mock.Anything, coinbaseInfo.BlockHash).Return(false, nil).Once()
		matchFunc := mock.MatchedBy(func(opts *bind.TransactOpts) bool {
			return opts.From == parsedAddress && opts.GasLimit == 100000
		})
		bridgeMock.EXPECT().RegisterBtcCoinbaseTransaction(matchFunc, coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue).Return(tx, nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, bridgeMock, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.NoError(t, err)
		bridgeMock.AssertExpectations(t)
	})
}

// nolint:funlen
func TestRskBridgeImpl_GetBatchPegOutCreatedEvent(t *testing.T) {
	batches := []bindings.RskBridgeBatchPegoutCreated{
		{
			BtcTxHash:          common.HexToHash("0x88c9bcbc80a4885335ddf8b26e45b4d3f1fbaf9ba3aebae9f2f315208a20bb88"),
			ReleaseRskTxHashes: hexutil.MustDecode("0x70d2a975d67783a11716dd4ace98730a7b25dacc80058bb8c0b6b501f8b1ddd7"),
			Raw: types.Log{
				Address: common.HexToAddress("0x0000000000000000000000000000000001000006"),
				Topics: []common.Hash{
					common.HexToHash("0x483d0191cc4e784b04a41f6c4801a0766b43b1fdd0b9e3e6bfdca74e5b05c2eb"),
					common.HexToHash("0x88c9bcbc80a4885335ddf8b26e45b4d3f1fbaf9ba3aebae9f2f315208a20bb88"),
				},
				Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000002070d2a975d67783a11716dd4ace98730a7b25dacc80058bb8c0b6b501f8b1ddd7"),
				BlockNumber: 7803285,
				TxHash:      common.HexToHash("0x3e116c9d63ad8cb15b21914be1a427ac01e28dbb6f1ce17e5f6628fc35e63aa6"),
				TxIndex:     6,
				BlockHash:   common.HexToHash("0x73aa5fa17f1999e62f4651261d3fbdf86c2992c8f8cb044562f8d2374cc4bf3d"),
				Index:       22,
				Removed:     false,
			},
		},
	}
	parsedBatches := []rsk.BatchPegOut{
		{
			TransactionHash:    "0x3e116c9d63ad8cb15b21914be1a427ac01e28dbb6f1ce17e5f6628fc35e63aa6",
			BlockHash:          "0x73aa5fa17f1999e62f4651261d3fbdf86c2992c8f8cb044562f8d2374cc4bf3d",
			BlockNumber:        7803285,
			BtcTxHash:          "88c9bcbc80a4885335ddf8b26e45b4d3f1fbaf9ba3aebae9f2f315208a20bb88",
			ReleaseRskTxHashes: []string{"0x70d2a975d67783a11716dd4ace98730a7b25dacc80058bb8c0b6b501f8b1ddd7"},
		},
	}
	bridgeMock := &mocks.RskBridgeAdapterMock{}
	iteratorMock := &mocks.EventIteratorAdapterMock[bindings.RskBridgeBatchPegoutCreated]{}
	filterMatchFunc := func(from uint64, to uint64) func(opts *bind.FilterOpts) bool {
		return func(opts *bind.FilterOpts) bool {
			return from == opts.Start && to == *opts.End && opts.Context != nil
		}
	}
	signer := &mocks.TransactionSignerMock{}
	bridge := rootstock.NewRskBridgeImpl(
		rootstock.RskBridgeConfig{},
		bridgeMock,
		dummyClient,
		&chaincfg.TestNet3Params,
		rootstock.RetryParams{},
		signer,
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		bridgeMock.EXPECT().FilterBatchPegoutCreated(mock.MatchedBy(filterMatchFunc(from, to)), [][32]byte(nil)).
			Return(&bindings.RskBridgeBatchPegoutCreatedIterator{}, nil).Once()
		bridgeMock.EXPECT().BatchPegOutCreatedIteratorAdapter(mock.AnythingOfType(batchPegOutIteratorString)).Return(iteratorMock)
		iteratorMock.EXPECT().Next().Return(true).Times(len(batches))
		iteratorMock.EXPECT().Next().Return(false).Once()
		for _, batch := range batches {
			iteratorMock.EXPECT().Event().Return(&batch).Once()
		}
		iteratorMock.EXPECT().Error().Return(nil).Once()
		iteratorMock.EXPECT().Close().Return(nil).Once()
		result, err := bridge.GetBatchPegOutCreatedEvent(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedBatches, result)
		bridgeMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
	t.Run("Error handling when failed to get iterator", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		bridgeMock.EXPECT().FilterBatchPegoutCreated(mock.MatchedBy(filterMatchFunc(from, to)), [][32]byte(nil)).
			Return(nil, assert.AnError).Once()
		bridgeMock.EXPECT().BatchPegOutCreatedIteratorAdapter(mock.AnythingOfType(batchPegOutIteratorString)).Return(nil)
		result, err := bridge.GetBatchPegOutCreatedEvent(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		bridgeMock.AssertExpectations(t)
	})
	t.Run("Error handling on iterator error", func(t *testing.T) {
		var from uint64 = 700
		var to uint64 = 1200
		bridgeMock.EXPECT().FilterBatchPegoutCreated(mock.MatchedBy(filterMatchFunc(from, to)), [][32]byte(nil)).
			Return(&bindings.RskBridgeBatchPegoutCreatedIterator{}, nil).Once()
		bridgeMock.EXPECT().BatchPegOutCreatedIteratorAdapter(mock.AnythingOfType(batchPegOutIteratorString)).Return(iteratorMock)
		iteratorMock.EXPECT().Next().Return(false).Once()
		iteratorMock.EXPECT().Error().Return(assert.AnError).Once()
		iteratorMock.EXPECT().Close().Return(nil).Once()
		result, err := bridge.GetBatchPegOutCreatedEvent(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		bridgeMock.AssertExpectations(t)
		iteratorMock.AssertExpectations(t)
	})
}
