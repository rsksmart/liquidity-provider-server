package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

const (
	estimationAddress              = "0x462d7082F3671a3be160638Be3F8c23cA354f48a"
	estimationBaseGas       uint64 = 57000
	estimationNewAccountGas uint64 = 25000
	txHash                         = "0x0e5a74de4d3f7eceff661d953f75270041c82ba0b0b787ec8daf7d566a53baa5"
	blockHash                      = "0x010203"
)

var (
	estimationData, _ = hex.DecodeString("5a68669900000000000000000000000000000000000000000000000002dda2a7ea1e40000000000000000000000000000000000000000000000000000000000066223d930000000000000000000000009d4b2c05818a0086e641437fcb64ab6098c7bbec")
	estimationValue   = entities.NewWei(300)
)

func TestRskjRpcServer_EstimateGas_NewAccount(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})

	toAddress := common.HexToAddress(estimationAddress)
	client.On("NonceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(uint64(0), nil).Once()
	client.On("CodeAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(nil, nil).Once()
	client.On("BalanceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(big.NewInt(0), nil).Once()
	client.On("EstimateGas", test.AnyCtx, ethereum.CallMsg{
		To:    &toAddress,
		Data:  estimationData,
		Value: estimationValue.AsBigInt(),
	}).Return(estimationBaseGas, nil).Once()
	result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
	require.NoError(t, err)
	assert.Equal(t, entities.NewUWei(estimationBaseGas+estimationNewAccountGas), result)
	client.AssertExpectations(t)
}

func TestRskjRpcServer_EstimateGas_ExistingAccount(t *testing.T) {
	toAddress := common.HexToAddress(estimationAddress)
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	client.On("EstimateGas", test.AnyCtx, ethereum.CallMsg{
		To:    &toAddress,
		Data:  estimationData,
		Value: estimationValue.AsBigInt(),
	}).Return(estimationBaseGas, nil).Times(3)
	t.Run("Existing nonce", func(t *testing.T) {
		client.On("NonceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(uint64(1), nil).Once()
		client.On("CodeAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(nil, nil).Once()
		client.On("BalanceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(big.NewInt(0), nil).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.NoError(t, err)
		assert.Equal(t, entities.NewUWei(estimationBaseGas), result)
	})
	t.Run("Existing code", func(t *testing.T) {
		client.On("NonceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(uint64(0), nil).Once()
		client.On("CodeAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return([]byte{1, 2, 3}, nil).Once()
		client.On("BalanceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(big.NewInt(0), nil).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.NoError(t, err)
		assert.Equal(t, entities.NewUWei(estimationBaseGas), result)
	})
	t.Run("Existing balance", func(t *testing.T) {
		client.On("NonceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(uint64(0), nil).Once()
		client.On("CodeAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(nil, nil).Once()
		client.On("BalanceAt", test.AnyCtx, toAddress, (*big.Int)(nil)).Return(big.NewInt(1), nil).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.NoError(t, err)
		assert.Equal(t, entities.NewUWei(estimationBaseGas), result)
	})
	client.AssertExpectations(t)
}

func TestRskjRpcServer_EstimateGas_ErrorHandling(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Invalid address", func(t *testing.T) {
		result, err := rpc.EstimateGas(context.Background(), test.AnyString, estimationValue, estimationData)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error calling CodeAt", func(t *testing.T) {
		client.On("CodeAt", test.AnyCtx, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error calling BalanceAt", func(t *testing.T) {
		client.On("CodeAt", test.AnyCtx, mock.Anything, mock.Anything).Return(nil, nil).Once()
		client.On("BalanceAt", test.AnyCtx, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error calling NonceAt", func(t *testing.T) {
		client.On("CodeAt", test.AnyCtx, mock.Anything, mock.Anything).Return(nil, nil).Once()
		client.On("BalanceAt", test.AnyCtx, mock.Anything, mock.Anything).Return(big.NewInt(0), nil).Once()
		client.On("NonceAt", test.AnyCtx, mock.Anything, mock.Anything).Return(uint64(0), assert.AnError).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error calling EstimateGas", func(t *testing.T) {
		client.On("CodeAt", test.AnyCtx, mock.Anything, mock.Anything).Return(nil, nil).Once()
		client.On("BalanceAt", test.AnyCtx, mock.Anything, mock.Anything).Return(big.NewInt(0), nil).Once()
		client.On("NonceAt", test.AnyCtx, mock.Anything, mock.Anything).Return(uint64(0), nil).Once()
		client.On("EstimateGas", test.AnyCtx, mock.Anything).Return(uint64(0), assert.AnError).Once()
		result, err := rpc.EstimateGas(context.Background(), estimationAddress, estimationValue, estimationData)
		require.Error(t, err)
		assert.Nil(t, result)
		client.AssertExpectations(t)
	})
}

func TestRskjRpcServer_GasPrice(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Success", func(t *testing.T) {
		client.On("SuggestGasPrice", test.AnyCtx).Return(big.NewInt(200), nil).Once()
		gasPrice, err := rpc.GasPrice(context.Background())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(200), gasPrice)
	})
	t.Run("Error calling SuggestGasPrice", func(t *testing.T) {
		client.On("SuggestGasPrice", test.AnyCtx).Return(nil, assert.AnError).Once()
		gasPrice, err := rpc.GasPrice(context.Background())
		require.Error(t, err)
		assert.Nil(t, gasPrice)
	})
}

func TestRskjRpcServer_GetBalance(t *testing.T) {
	var blockNumber *big.Int = nil
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Success", func(t *testing.T) {
		client.On("BalanceAt", test.AnyCtx, parsedAddress, blockNumber).Return(big.NewInt(100), nil).Once()
		balance, err := rpc.GetBalance(context.Background(), parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(100), balance)
	})
	t.Run("Error calling BalanceAt", func(t *testing.T) {
		client.On("BalanceAt", test.AnyCtx, parsedAddress, blockNumber).Return(nil, assert.AnError).Once()
		balance, err := rpc.GetBalance(context.Background(), parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, balance)
	})
	t.Run("Invalid address", func(t *testing.T) {
		balance, err := rpc.GetBalance(context.Background(), test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Nil(t, balance)
	})
}

func TestRskjRpcServer_GetHeight(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Success", func(t *testing.T) {
		client.On("BlockNumber", test.AnyCtx).Return(uint64(50), nil).Once()
		blockNumber, err := rpc.GetHeight(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(50), blockNumber)
	})
	t.Run("Error calling BlockNumber", func(t *testing.T) {
		client.On("BlockNumber", test.AnyCtx).Return(uint64(0), assert.AnError).Once()
		blockNumber, err := rpc.GetHeight(context.Background())
		require.Error(t, err)
		assert.Zero(t, blockNumber)
	})
}

func TestRskjRpcServer_GetTransactionReceipt(t *testing.T) {
	const (
		v int64 = 0x62
		r       = "73e409ecab98206d4f2afbf6953739ed30002bda88760e2a211e23334766b467"
		s       = "3a020211dfe07777d3d6373771fc848e0a777b2647ee8c4df5c1e44b22e13b39"
	)
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Success", func(t *testing.T) {
		client.On("TransactionReceipt", test.AnyCtx, common.HexToHash(txHash)).Return(&types.Receipt{
			GasUsed:           456,
			CumulativeGasUsed: 123,
			TxHash:            common.HexToHash(txHash),
			BlockHash:         common.HexToHash(blockHash),
			BlockNumber:       big.NewInt(500),
		}, nil).Once()
		parsedToAddress := common.HexToAddress("0x462d7082F3671a3be160638Be3F8c23cA354f48a")
		rAsBigInt := new(big.Int)
		rAsBigInt.SetString(r, 16)
		sAsBigInt := new(big.Int)
		sAsBigInt.SetString(s, 16)
		data, _ := hex.DecodeString("5a68669900000000000000000000000000000000000000000000000002dda2a7ea1e40000000000000000000000000000000000000000000000000000000000066223d930000000000000000000000009d4b2c05818a0086e641437fcb64ab6098c7bbec")
		client.On("TransactionByHash", test.AnyCtx, common.HexToHash(txHash)).
			Return(types.NewTx(&types.LegacyTx{
				Nonce:    741514,
				GasPrice: big.NewInt(65826581),
				Gas:      200000,
				To:       &parsedToAddress,
				Value:    big.NewInt(0),
				Data:     data,
				V:        big.NewInt(v),
				R:        rAsBigInt,
				S:        sAsBigInt,
			}), false, nil).Once()
		receipt, err := rpc.GetTransactionReceipt(context.Background(), txHash)
		require.NoError(t, err)
		assert.Equal(t, blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         "0x0000000000000000000000000000000000000000000000000000000000010203",
			BlockNumber:       500,
			From:              "0xC67D9EE30d2119A384E02de568BE80fe785074Ba",
			To:                parsedToAddress.String(),
			CumulativeGasUsed: big.NewInt(123),
			GasUsed:           big.NewInt(456),
			Value:             entities.NewWei(0),
		}, receipt)
	})
}

func TestRskjRpcServer_GetTransactionReceipt_ErrorHandling(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Error error getting receipt", func(t *testing.T) {
		client.On("TransactionReceipt", test.AnyCtx, common.HexToHash(txHash)).
			Return(nil, assert.AnError).Once()
		receipt, err := rpc.GetTransactionReceipt(context.Background(), txHash)
		require.Error(t, err)
		assert.Empty(t, receipt)
	})
	t.Run("Error error getting transaction", func(t *testing.T) {
		client.On("TransactionReceipt", test.AnyCtx, common.HexToHash(txHash)).Return(&types.Receipt{
			GasUsed:           456,
			CumulativeGasUsed: 123,
			TxHash:            common.HexToHash(txHash),
			BlockHash:         common.HexToHash(blockHash),
			BlockNumber:       big.NewInt(500),
		}, nil).Once()
		client.On("TransactionByHash", test.AnyCtx, common.HexToHash(txHash)).
			Return(nil, false, assert.AnError).Once()
		receipt, err := rpc.GetTransactionReceipt(context.Background(), txHash)
		require.Error(t, err)
		assert.Empty(t, receipt)
	})
	t.Run("Invalid tx hash", func(t *testing.T) {
		receipt, err := rpc.GetTransactionReceipt(context.Background(), test.AnyString)
		require.Error(t, err)
		assert.Empty(t, receipt)
	})
}

func TestRskjRpcServer_GetBlockByHash(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	var now int64 = 1714471719922
	client.On("BlockByHash", test.AnyCtx, common.HexToHash(blockHash)).Return(types.NewBlock(
		&types.Header{
			Number: big.NewInt(123),
			Time:   uint64(now),
			Nonce:  [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		}, nil, nil, nil), nil).Once()
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	block, err := rpc.GetBlockByHash(context.Background(), blockHash)
	require.NoError(t, err)
	assert.Equal(t, blockchain.BlockInfo{
		Hash:      "0xde378ac47c11cdc8182c05f10edd90899fced079aa2b141f4f548b354deac5d8",
		Number:    123,
		Timestamp: time.Unix(now, 0),
		Nonce:     72623859790382856,
	}, block)
	client.AssertExpectations(t)
}

func TestRskjRpcServer_GetBlockByHash_ErrorHandling(t *testing.T) {
	client := &mocks.RpcClientBindingMock{}
	rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
	t.Run("Error error getting block", func(t *testing.T) {
		client.On("BlockByHash", test.AnyCtx, common.HexToHash(blockHash)).Return(nil, assert.AnError).Once()
		block, err := rpc.GetBlockByHash(context.Background(), blockHash)
		require.Error(t, err)
		assert.Empty(t, block)
	})
	t.Run("Invalid tx hash", func(t *testing.T) {
		block, err := rpc.GetBlockByHash(context.Background(), test.AnyString)
		require.Error(t, err)
		assert.Empty(t, block)
	})
}
