package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/bridge"
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

func TestNewRskBridgeImpl(t *testing.T) {
	config := rootstock.RskBridgeConfig{Address: test.AnyAddress, RequiredConfirmations: 10, ErpKeys: []string{"key1", "key2", "key3"}, UseSegwitFederation: true}
	client := rootstock.NewRskClient(&mocks.RpcClientBindingMock{})
	contract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewRskBridge()
	bridge := rootstock.NewRskBridgeImpl(config, contract, client, &chaincfg.TestNet3Params, rootstock.RetryParams{Retries: 1, Sleep: time.Duration(1)}, &mocks.TransactionSignerMock{}, contractBinding, time.Duration(1))
	test.AssertNonZeroValues(t, bridge)
}

func TestRskBridgeImpl_GetAddress(t *testing.T) {
	contract := bind.NewBoundContract(common.Address{}, abi.ABI{}, nil, nil, nil)
	contractBinding := bindings.NewRskBridge()
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{Address: test.AnyAddress}, contract, dummyClient, nil, rootstock.RetryParams{}, nil, contractBinding, time.Duration(1))
	assert.Equal(t, test.AnyAddress, bridge.GetAddress())
}

func TestRskBridgeImpl_GetRequiredTxConfirmations(t *testing.T) {
	bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{RequiredConfirmations: 10}, nil, dummyClient, nil, rootstock.RetryParams{}, nil, nil, time.Duration(1))
	assert.Equal(t, uint64(10), bridge.GetRequiredTxConfirmations())
}

func TestRskBridgeImpl_GetFedAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederationAddress()),
			mock.Anything,
		).Return(mustPackString(t, test.AnyAddress), nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, nil, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
		result, err := bridge.GetFedAddress()
		assert.Equal(t, test.AnyAddress, result)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetFederationAddress call fail", func(t *testing.T) {
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederationAddress()),
			mock.Anything,
		).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, nil, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
		result, err := bridge.GetFedAddress()
		assert.Empty(t, result)
		require.Error(t, err)
		contractMock.caller.AssertExpectations(t)
	})
}

func TestRskBridgeImpl_GetMinimumLockTxValue(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetMinimumLockTxValue()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(5)), nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, nil, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
		result, err := bridge.GetMinimumLockTxValue()
		assert.IsType(t, &entities.Wei{}, result)
		assert.Equal(t, entities.NewWei(50000000000), result)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetMinimumLockTxValue call fail", func(t *testing.T) {
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetMinimumLockTxValue()),
			mock.Anything,
		).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, nil, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
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
		var result rsk.FlyoverDerivation
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetActivePowpegRedeemScript()),
			mock.Anything,
		).Return(mustPackBytes(t, redeemScriptString), nil)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Equal(t, rsk.FlyoverDerivation{
			Address:      "2MxeEHVx71taCeVsXFsfQ7TKK6v943PFVEu",
			RedeemScript: "20ff883edd54f8cb22464a8181ed62652fcdb0028e0ada18f9828afd76e0df2c727564522102cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1210362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a1242103c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db53ae670350cd00b27553210257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d42103c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f92103cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b32102370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec8054ae68",
		}, result)
		require.NoError(t, testError)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetActivePowpegRedeemScript call fail", func(t *testing.T) {
		var testError error
		var result rsk.FlyoverDerivation
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetActivePowpegRedeemScript()),
			mock.Anything,
		).Return(nil, assert.AnError)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
		result, testError = bridge.GetFlyoverDerivationAddress(args)
		assert.Empty(t, result)
		require.ErrorContains(t, testError, "error retreiving fed redeem script from bridge")
	})
}

func TestRskBridgeImpl_FetchFederationInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		contractMock := createBoundContractMock()
		bridgeBinding := bindings.NewRskBridge()
		contractMock.caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(mustPackUint256(t, big.NewInt(2)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(0), "btc")),
			mock.Anything,
		).Return(mustPackBytes(t, "010203"), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(1), "btc")),
			mock.Anything,
		).Return(mustPackBytes(t, "0a0b0c"), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederationThreshold()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(5)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetFederationAddress()),
			mock.Anything,
		).Return(mustPackString(t, test.AnyAddress), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetActiveFederationCreationBlockHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(500)), nil).Once()

		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{ErpKeys: []string{"key1", "key2", "key3"}, UseSegwitFederation: true},
			contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
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
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling", func(t *testing.T) {
		for _, setUp := range fetchFedInfoErrorSetUps() {
			contractMock := createBoundContractMock()
			bridgeBinding := bindings.NewRskBridge()
			setUp(t, bridgeBinding, contractMock.caller)
			bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, nil, bridgeBinding, time.Duration(1))
			result, err := bridge.FetchFederationInfo()
			require.Error(t, err)
			assert.Empty(t, result)
			contractMock.caller.AssertExpectations(t)
		}
	})
}

func fetchFedInfoErrorSetUps() []func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
	return []func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock){
		func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(mustPackUint256(t, big.NewInt(2)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(0), "btc")), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(mustPackUint256(t, big.NewInt(2)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(0), "btc")), mock.Anything).Return(mustPackBytes(t, "010203"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(1), "btc")), mock.Anything).Return(mustPackBytes(t, "0a0b0c"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationThreshold()), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(mustPackUint256(t, big.NewInt(2)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(0), "btc")), mock.Anything).Return(mustPackBytes(t, "010203"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(1), "btc")), mock.Anything).Return(mustPackBytes(t, "0a0b0c"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationThreshold()), mock.Anything).Return(mustPackUint256(t, big.NewInt(5)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationAddress()), mock.Anything).Return(nil, assert.AnError).Once()
		},
		func(t *testing.T, bridgeBinding *bindings.RskBridge, caller *mocks.ContractCallerMock) {
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationSize()), mock.Anything).Return(mustPackUint256(t, big.NewInt(2)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(0), "btc")), mock.Anything).Return(mustPackBytes(t, "010203"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederatorPublicKeyOfType(big.NewInt(1), "btc")), mock.Anything).Return(mustPackBytes(t, "0a0b0c"), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationThreshold()), mock.Anything).Return(mustPackUint256(t, big.NewInt(5)), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetFederationAddress()), mock.Anything).Return(mustPackString(t, test.AnyAddress), nil).Once()
			caller.EXPECT().CallContract(mock.Anything, matchCallData(bridgeBinding.PackGetActiveFederationCreationBlockHeight()), mock.Anything).Return(nil, assert.AnError).Once()
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
	bridgeBinding := bindings.NewRskBridge()
	t.Run("Should handle error getting best chain height", func(t *testing.T) {
		contractMock := createBoundContractMock()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		require.NotErrorIs(t, err, blockchain.WaitingForBridgeError)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Should return WaitingForBridgeError if block is higher than best chain height", func(t *testing.T) {
		contractMock := createBoundContractMock()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(300)), nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Should handle error when validating if tx was registered", func(t *testing.T) {
		contractMock := createBoundContractMock()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackHasBtcBlockCoinbaseTransactionInformation(coinbaseInfo.BlockHash)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.Error(t, err)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Should skip registration if tx was already registered", func(t *testing.T) {
		contractMock := createBoundContractMock()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackHasBtcBlockCoinbaseTransactionInformation(coinbaseInfo.BlockHash)),
			mock.Anything,
		).Return(mustPackBool(t, true), nil).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Should handle error when sending tx", func(t *testing.T) {
		contractMock := createBoundContractMock()
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return tx, nil
		})
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackHasBtcBlockCoinbaseTransactionInformation(coinbaseInfo.BlockHash)),
			mock.Anything,
		).Return(mustPackBool(t, false), nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 100000, big.NewInt(0), bridgeBinding.PackRegisterBtcCoinbaseTransaction(coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue)),
		).Return(assert.AnError).Once()
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, dummyClient, &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.Empty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error")
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Should handle error when tx reverts", func(t *testing.T) {
		contractMock := createBoundContractMock()
		mockClient := &mocks.RpcClientBindingMock{}
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackHasBtcBlockCoinbaseTransactionInformation(coinbaseInfo.BlockHash)),
			mock.Anything,
		).Return(mustPackBool(t, false), nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 100000, big.NewInt(0), bridgeBinding.PackRegisterBtcCoinbaseTransaction(coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue)),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.ErrorContains(t, err, "register coinbase transaction error: transaction reverted (0xd8f5d705f146230553a8aec9a290a19bf4311187fa0489d41207d7215b0b65cb)")
		contractMock.transactor.AssertExpectations(t)
	})

	t.Run("Should register coinbase transaction successfully", func(t *testing.T) {
		contractMock := createBoundContractMock()
		mockClient := &mocks.RpcClientBindingMock{}
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackGetBtcBlockchainBestChainHeight()),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(bridgeBinding.PackHasBtcBlockCoinbaseTransactionInformation(coinbaseInfo.BlockHash)),
			mock.Anything,
		).Return(mustPackBool(t, false), nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 100000, big.NewInt(0), bridgeBinding.PackRegisterBtcCoinbaseTransaction(coinbaseInfo.BtcTxSerialized, coinbaseInfo.BlockHash, coinbaseInfo.SerializedPmt, coinbaseInfo.WitnessMerkleRoot, coinbaseInfo.WitnessReservedValue)),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		bridge := rootstock.NewRskBridgeImpl(rootstock.RskBridgeConfig{}, contractMock.contract, rootstock.NewRskClient(mockClient), &chaincfg.TestNet3Params, rootstock.RetryParams{}, signerMock, bridgeBinding, time.Duration(1))
		result, err := bridge.RegisterBtcCoinbaseTransaction(coinbaseInfo)
		assert.NotEmpty(t, result)
		require.NoError(t, err)
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestRskBridgeImpl_GetBatchPegOutCreatedEvent(t *testing.T) {
	batches := []types.Log{
		{
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
	bridgeBinding := bindings.NewRskBridge()
	contractMock := createBoundContractMock()
	signer := &mocks.TransactionSignerMock{}
	bridge := rootstock.NewRskBridgeImpl(
		rootstock.RskBridgeConfig{},
		contractMock.contract,
		dummyClient,
		&chaincfg.TestNet3Params,
		rootstock.RetryParams{},
		signer,
		bridgeBinding,
		time.Duration(1),
	)
	t.Run("Success", func(t *testing.T) {
		var from uint64 = 500
		var to uint64 = 1000
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(batches, nil).Once()
		result, err := bridge.GetBatchPegOutCreatedEvent(context.Background(), from, &to)
		require.NoError(t, err)
		assert.Equal(t, parsedBatches, result)
		contractMock.filterer.AssertExpectations(t)
	})
	t.Run("Error handling when filtering logs", func(t *testing.T) {
		var from uint64 = 600
		var to uint64 = 1100
		contractMock.filterer.EXPECT().FilterLogs(mock.Anything, mock.MatchedBy(filterMatchFunc(from, to))).Return(nil, assert.AnError).Once()
		result, err := bridge.GetBatchPegOutCreatedEvent(context.Background(), from, &to)
		require.Error(t, err)
		assert.Nil(t, result)
		contractMock.filterer.AssertExpectations(t)
	})
}
