package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
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
	// nolint:errcheck
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

// nolint:funlen
func TestRskjRpcServer_GetTransactionReceipt(t *testing.T) {
	const (
		v int64 = 0x62
		r       = "73e409ecab98206d4f2afbf6953739ed30002bda88760e2a211e23334766b467"
		s       = "3a020211dfe07777d3d6373771fc848e0a777b2647ee8c4df5c1e44b22e13b39"
	)
	t.Run("should get receipt without logs", func(t *testing.T) {
		client := &mocks.RpcClientBindingMock{}
		rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
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
		data, err := hex.DecodeString("5a68669900000000000000000000000000000000000000000000000002dda2a7ea1e40000000000000000000000000000000000000000000000000000000000066223d930000000000000000000000009d4b2c05818a0086e641437fcb64ab6098c7bbec")
		require.NoError(t, err)
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
			Logs:              make([]blockchain.TransactionLog, 0),
		}, receipt)
	})
	t.Run("should get receipt with logs", func(t *testing.T) {
		const (
			transaction = "0x2c78749f71d5f1e870397eaf401c3683ca752577a7a2fcbf64e3eb2091c817dd"
			block       = "0x9ca0a12f5854335bd78c145c105543b724043277a71dc63fb7b8c5ed9760af02"
			blockNumber = 0x755fa7
			contract    = "0xAA9cAf1e3967600578727F975F283446A3Da6612"
		)
		client := &mocks.RpcClientBindingMock{}
		rpc := rootstock.NewRskjRpcServer(rootstock.NewRskClient(client), rootstock.RetryParams{})
		rawReceipt := &types.Receipt{
			Type:              0,
			Status:            1,
			CumulativeGasUsed: 0x2314b,
			Bloom:             types.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000080000000002000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000040000000000000000001000000000000000000000000000000000000100000040000000400000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000400000000")),
			Logs: []*types.Log{
				{
					BlockNumber: blockNumber,
					BlockHash:   common.HexToHash(block),
					TxHash:      common.HexToHash(transaction),
					Index:       0,
					Address:     common.HexToAddress(contract),
					Data:        hexutil.MustDecode("0x00000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d325785200000000000000000000000000000000000000000000000014c337856ab59600"),
					Topics: []common.Hash{
						common.HexToHash("0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53"),
					},
					TxIndex: 0,
					Removed: false,
				},
				{
					BlockNumber: blockNumber,
					BlockHash:   common.HexToHash(block),
					TxHash:      common.HexToHash(transaction),
					Index:       1,
					Address:     common.HexToAddress(contract),
					Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000014c3378607043c00"),
					Topics: []common.Hash{
						common.HexToHash("0x0629ae9d1dc61501b0ca90670a9a9b88daaf7504b54537b53e1219de794c63d2"),
						common.HexToHash("0xc9dffd38b14139b2fb49ed67bc1a221d3f84e119a638bb424511ba4cd43e4c5e"),
					},
					TxIndex: 0,
					Removed: false,
				},
			},
			TxHash:            common.HexToHash("0x2c78749f71d5f1e870397eaf401c3683ca752577a7a2fcbf64e3eb2091c817dd"),
			GasUsed:           0x2314b,
			EffectiveGasPrice: big.NewInt(0x18dbac0),
			BlockHash:         common.HexToHash("0x9ca0a12f5854335bd78c145c105543b724043277a71dc63fb7b8c5ed9760af02"),
			BlockNumber:       big.NewInt(blockNumber),
			TransactionIndex:  0,
		}
		to := common.HexToAddress(contract)
		txData := types.LegacyTx{
			Nonce:    0x1d0,
			GasPrice: big.NewInt(0x18dbac0),
			Gas:      0x2625a0,
			To:       &to,
			Value:    big.NewInt(0),
			Data:     hexutil.MustDecode("0x2d67952c00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000003c00000000000000000000000000000000000000000000000000000000000000440000000000000000000000000000000000000000000000000000000000000058000000000000000000000000000000000000000000000000000000000000dc3870ae2bb1b3bb7aed09b42aad17f4ebd54fca6f3ca000000000000000000000000000000000000000000000000aa9caf1e3967600578727f975f283446a3da661200000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d32578520000000000000000000000000000000000000000000000000000000000000280000000000000000000000000ac31a4beedd7ec916b7a48a612230cb85c1aaf5600000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000005af3107a4000000000000000000000000000000000000000000000000000000009184e72a000000000000000000000000000ac31a4beedd7ec916b7a48a612230cb85c1aaf560000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000520800000000000000000000000000000000000000000000000074546abbfd23123b00000000000000000000000000000000000000000000000014c2dc12e7fc000000000000000000000000000000000000000000000000000000000000685527e500000000000000000000000000000000000000000000000000000000000015180000000000000000000000000000000000000000000000000000000000001c200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007f723f560000000000000000000000000000000000000000000000000000000000000000150000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001500840098213fec4001cdc4a77cc3340f5bb83d9ed5000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041dc942ceb0e3bdd9815f53e536e70aaac2217fc422c9370a0f07fde7cdab39eed4729a933f1e90e6c8cbd7bbfaff15273de6f82d0f34516f99783f18bf558c7641c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011601000000053d0ab910f144862c6c9d7049f3cf70a4fa8e3eb2185726f99173caaa8b0f00a10100000000ffffffffc4655c058fabcf3b0ed7a5186d81e04fc97f1ea4b5ffcca19694da67c839a2420000000000ffffffffd870ecc8e2a3cecef46d7244264678277840aa2f6d82a2dd973acf888e0ac8dd0100000000ffffffffcacaee423f27ad6cfc76fe2cfbd4bd20171919f6615e07cb27f2fcca21bb0ed80100000000ffffffff28721107c40671b831a30fbaf2e6e8776a0947229e4bf8fb41c9d47894d59aa30200000000ffffffff0247deea080000000017a9145be63a8b8f3f848981d31fa163a0c89f354a1c50878b7d000000000000160014f011a71128cb7e04fe511e868543b9f6cf81426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000169070200000bcd88f6365a9cd71dbde59464d2c4b16b2080b631e61d62383a790dc37acfff6e6dfc3733f939bdb9a53b643139f7eddde0927c236b8b99e93fed7a1d27626abf880ca69c03569d61efa7ec9f1039d748212a73a7deeb8efa9654c79c002441867a140bd11767a41bf7c71ace188719e1d2c233fbdb64b4fb223b761daf311982ffffffc6fc50ff2cd0636b80d2d9ac06db90534cb03bfa347c631ecb40849fbade60d5e556b13f08131f047066d5d9040a1e0a3820d12bfae8a4cadd6c4d0400b9fa236eb57fd4cd2614e9c1f38f84e6d35f40800bb8076197596bc9be8db148428dbed1254a3369a99492c6c02bc99b84ca97dcf623cd88d2aa56dfc6f63fb3f3756234efce77e21d9c7d704b243b4ca842ffa1b7aede2d01edbc50717abad92f8b8b68234c8d6cb9084a2c0819f21e237781d341905cbe0d3ea47be8f4ce44d25a9f4860bf57d610f8638dd85202a12beaee7664d922e48184550501cd66c8036f2f000000000000000000000000000000000000000000000000"),
			V:        big.NewInt(0x60),
			R:        new(big.Int).SetBytes(hexutil.MustDecode("0xda4204b9a42bcdc0442d36fd87fee053dc121674e7aaf7639dcdb67631f755ac")),
			S:        new(big.Int).SetBytes(hexutil.MustDecode("0x192902edaeb78e096b555b5884b6c0a62d10809f00a28ad8883e0074a4755f7f")),
		}
		rawTx := types.NewTx(&txData)
		client.On("TransactionReceipt", mock.Anything, common.HexToHash(transaction)).Return(rawReceipt, nil).Once()
		client.On("TransactionByHash", mock.Anything, common.HexToHash(transaction)).Return(rawTx, false, nil).Once()
		receipt, err := rpc.GetTransactionReceipt(context.Background(), transaction)
		require.NoError(t, err)
		assert.Equal(t, blockchain.TransactionReceipt{
			TransactionHash:   transaction,
			BlockHash:         block,
			BlockNumber:       blockNumber,
			From:              "0x82a06eBDB97776a2da4041dF8f2b2ea8D3257852",
			To:                contract,
			CumulativeGasUsed: big.NewInt(143691),
			GasUsed:           big.NewInt(143691),
			Value:             entities.NewWei(0),
			Logs: []blockchain.TransactionLog{
				{
					Address: contract,
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53")),
					},
					Data:        hexutil.MustDecode("0x00000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d325785200000000000000000000000000000000000000000000000014c337856ab59600"),
					BlockNumber: blockNumber,
					TxHash:      transaction,
					TxIndex:     0,
					BlockHash:   block,
					Index:       0,
					Removed:     false,
				},
				{
					Address: contract,
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0x0629ae9d1dc61501b0ca90670a9a9b88daaf7504b54537b53e1219de794c63d2")),
						utils.To32Bytes(hexutil.MustDecode("0xc9dffd38b14139b2fb49ed67bc1a221d3f84e119a638bb424511ba4cd43e4c5e")),
					},
					Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000014c3378607043c00"),
					BlockNumber: blockNumber,
					TxHash:      transaction,
					TxIndex:     0,
					BlockHash:   block,
					Index:       1,
					Removed:     false,
				},
			},
		}, receipt)
		client.AssertExpectations(t)
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
