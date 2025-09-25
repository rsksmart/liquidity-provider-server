package rootstock_test

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var contractProviders = []bindings.FlyoverLiquidityProvider{
	{
		Id:              big.NewInt(1),
		Name:            "test",
		ApiBaseUrl:      "http://test.com",
		Status:          true,
		ProviderAddress: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		ProviderType:    uint8(liquidity_provider.PeginProvider),
	},
	{
		Id:              big.NewInt(2),
		Name:            "test2",
		ApiBaseUrl:      "http://test2.com",
		Status:          true,
		ProviderAddress: common.HexToAddress("0x1233557799abcdef1234567890abcdef12345678"),
		ProviderType:    uint8(liquidity_provider.PegoutProvider),
	},
}
var parsedProviders = []liquidity_provider.RegisteredLiquidityProvider{
	{
		Id:           1,
		Address:      "0x1234567890AbcdEF1234567890aBcdef12345678",
		Name:         "test",
		ApiBaseUrl:   "http://test.com",
		Status:       true,
		ProviderType: liquidity_provider.PeginProvider,
	},
	{
		Id:           2,
		Address:      "0x1233557799ABcDEf1234567890abCdef12345678",
		Name:         "test2",
		ApiBaseUrl:   "http://test2.com",
		Status:       true,
		ProviderType: liquidity_provider.PegoutProvider,
	},
}

func TestNewDiscoveryContractImpl(t *testing.T) {
	contract := rootstock.NewDiscoveryContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		&mocks.DiscoveryBindingMock{},
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestDiscoveryContractImpl_GetAddress(t *testing.T) {
	discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	assert.Equal(t, test.AnyAddress, discovery.GetAddress())
}

func TestDiscoveryContractImpl_SetProviderStatus(t *testing.T) {
	contractBinding := &mocks.DiscoveryBindingMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("SetProviderStatus", mock.Anything, big.NewInt(2), true).Return(tx, nil).Once()
		err := discovery.SetProviderStatus(2, true)
		require.NoError(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending setProviderStatus tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("SetProviderStatus", mock.Anything, big.NewInt(1), true).Return(nil, assert.AnError).Once()
		err := discovery.SetProviderStatus(1, true)
		require.Error(t, err)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (setProviderStatus tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.On("SetProviderStatus", mock.Anything, big.NewInt(1), false).Return(tx, nil).Once()
		err := discovery.SetProviderStatus(1, false)
		require.ErrorContains(t, err, "setProviderStatus transaction failed")
		contractBinding.AssertExpectations(t)
	})
}

func TestDiscoveryContractImpl_GetProvider(t *testing.T) {
	contractBinding := &mocks.DiscoveryBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.FlyoverLiquidityProvider{
			Id:              big.NewInt(5),
			ProviderAddress: parsedAddress,
			Name:            test.AnyString,
			ApiBaseUrl:      test.AnyUrl,
			Status:          true,
			ProviderType:    uint8(liquidity_provider.FullProvider),
		}, nil).Once()
		result, err := discovery.GetProvider(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.RegisteredLiquidityProvider{
			Id:           5,
			Address:      parsedAddress.String(),
			Name:         test.AnyString,
			ApiBaseUrl:   test.AnyUrl,
			Status:       true,
			ProviderType: liquidity_provider.FullProvider,
		}, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on GetProvider call fail", func(t *testing.T) {
		contractBinding.On("GetProvider", mock.Anything, parsedAddress).Return(bindings.FlyoverLiquidityProvider{}, assert.AnError).Once()
		result, err := discovery.GetProvider(parsedAddress.String())
		require.Error(t, err)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run(invalidAddressTest, func(t *testing.T) {
		result, err := discovery.GetProvider(test.AnyString)
		require.ErrorIs(t, err, blockchain.InvalidAddressError)
		assert.Empty(t, result)
	})
	t.Run("Invalid type", func(t *testing.T) {
		contractBinding.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.FlyoverLiquidityProvider{
			Id:              big.NewInt(5),
			ProviderAddress: parsedAddress,
			Name:            test.AnyString,
			ApiBaseUrl:      test.AnyUrl,
			Status:          true,
			ProviderType:    5,
		}, nil).Once()
		result, err := discovery.GetProvider(parsedAddress.String())
		require.ErrorIs(t, err, liquidity_provider.InvalidProviderTypeError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Provider not found", func(t *testing.T) {
		e := NewRskRpcError("execution reverted", "0x232cb27a000000000000000000000000dd2fd4581271e230360230f9337d5c0430bf44c0")
		contractBinding.EXPECT().GetProvider(mock.Anything, parsedAddress).Return(bindings.FlyoverLiquidityProvider{}, e).Once()
		result, err := discovery.GetProvider(parsedAddress.String())
		require.ErrorIs(t, err, liquidity_provider.ProviderNotFoundError)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
	})
}

func TestDiscoveryContractImpl_GetProviders(t *testing.T) {
	contractBinding := &mocks.DiscoveryBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
	t.Run("Success", func(t *testing.T) {
		contractBinding.On("GetProviders", mock.Anything).Return(contractProviders, nil).Once()
		result, err := discovery.GetProviders()
		require.NoError(t, err)
		assert.Equal(t, parsedProviders, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling on invalid provider type", func(t *testing.T) {
		invalidProviders := contractProviders
		invalidProviders[0].ProviderType = 5
		contractBinding.On("GetProviders", mock.Anything).Return(invalidProviders, nil).Once()
		result, err := discovery.GetProviders()
		require.ErrorIs(t, err, liquidity_provider.InvalidProviderTypeError)
		assert.Nil(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling GetProviders", func(t *testing.T) {
		contractBinding.On("GetProviders", mock.Anything).Return(nil, assert.AnError).Once()
		result, err := discovery.GetProviders()
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDiscoveryContractImpl_UpdateProvider(t *testing.T) {
	const (
		name = "test name"
		url  = "http://test.update.example.com"
	)

	contractBinding := &mocks.DiscoveryBindingMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.EXPECT().UpdateProvider(mock.Anything, name, url).Return(tx, nil).Once()
		result, err := discovery.UpdateProvider(name, url)
		require.NoError(t, err)
		assert.Equal(t, tx.Hash().String(), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling when sending updateProvider tx", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.EXPECT().UpdateProvider(mock.Anything, name, url).Return(nil, assert.AnError).Once()
		result, err := discovery.UpdateProvider(name, url)
		require.Error(t, err)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)

		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.EXPECT().UpdateProvider(mock.Anything, name, url).Return(nil, nil).Once()
		result, err = discovery.UpdateProvider(name, url)
		require.Error(t, err)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (updateProvider tx reverted)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.EXPECT().UpdateProvider(mock.Anything, name, url).Return(tx, nil).Once()
		result, err := discovery.UpdateProvider(name, url)
		require.ErrorContains(t, err, "update provider error")
		contractBinding.AssertExpectations(t)
		assert.Equal(t, tx.Hash().String(), result)
	})
}

func TestDiscoveryContractImpl_RegisterProvider(t *testing.T) {
	contractBinding := &mocks.DiscoveryBindingMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(800)}
	params := blockchain.ProviderRegistrationParams{
		Name:       "mock provider",
		ApiBaseUrl: "url.com",
		Status:     true,
		Type:       liquidity_provider.FullProvider,
	}
	t.Run("Success", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		receipt, err := mockClient.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		data, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000d529ae9e860000")
		require.NoError(t, err)
		receipt.Logs = append(receipt.Logs, &geth.Log{
			Address: common.HexToAddress("0xAa9caf1e3967600578727f975F283446a3dA6612"),
			Topics: []common.Hash{
				common.HexToHash("0xa9d44d6e13bb3fee938c3f66d1103e91f8dc6b12d4405a55eea558e8f275aa6e"),
				common.HexToHash("0x0000000000000000000000004202bac9919c3412fc7c8be4e678e26279386603"),
			},
			Data:        data,
			BlockNumber: 5778711,
			TxHash:      common.HexToHash("0x37e52bd50866063727188751052e35510b8bc7d5de72541b84168cb2cb8b9c6c"),
			TxIndex:     0,
			BlockHash:   common.HexToHash("0xdc48007bd41ed3d8027aaac9c67fe1142107453b80b1fead090490fa8cbd751a"),
			Index:       0,
			Removed:     false,
		})
		contractBinding.On("ParseRegister", *receipt.Logs[0]).Return(&bindings.IFlyoverDiscoveryRegister{
			Id:     big.NewInt(1),
			From:   parsedAddress,
			Amount: txConfig.Value.AsBigInt(),
			Raw:    *receipt.Logs[0],
		}, nil)
		mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
		contractBinding.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, uint8(params.Type)).
			Return(tx, nil).Once()
		result, err := discovery.RegisterProvider(txConfig, params)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)
		contractBinding.AssertExpectations(t)
	})
}

func TestDiscoveryContractImpl_RegisterProvider_ErrorHandling(t *testing.T) {
	const incompleteReceipt = "incomplete receipt"
	contractBinding := &mocks.DiscoveryBindingMock{}
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	discovery := rootstock.NewDiscoveryContractImpl(
		rootstock.NewRskClient(mockClient),
		test.AnyAddress,
		contractBinding,
		signerMock,
		rootstock.RetryParams{},
		time.Duration(1),
		Abis,
	)
	txConfig := blockchain.TransactionConfig{Value: entities.NewWei(800)}
	params := blockchain.ProviderRegistrationParams{Name: "mock provider", ApiBaseUrl: "url.com", Status: true, Type: liquidity_provider.FullProvider}
	t.Run("Error handling (send transaction error)", func(t *testing.T) {
		_ = prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, uint8(params.Type)).
			Return(nil, assert.AnError).Once()
		result, err := discovery.RegisterProvider(txConfig, params)
		require.Error(t, err)
		assert.Equal(t, int64(0), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (receipt without event)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		contractBinding.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, uint8(params.Type)).
			Return(tx, nil).Once()
		result, err := discovery.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, incompleteReceipt)
		assert.Equal(t, int64(0), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (transaction revert)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, false)
		contractBinding.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, uint8(params.Type)).
			Return(tx, nil).Once()
		result, err := discovery.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, incompleteReceipt)
		assert.Equal(t, int64(0), result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error handling (parsing error)", func(t *testing.T) {
		tx := prepareTxMocks(mockClient, signerMock, true)
		receipt, err := mockClient.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		receipt.Logs = append(receipt.Logs, &geth.Log{})
		contractBinding.On("ParseRegister", *receipt.Logs[0]).Return(nil, assert.AnError)
		mockClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil).Once()
		contractBinding.On("Register", mock.Anything,
			params.Name, params.ApiBaseUrl, params.Status, uint8(params.Type)).
			Return(tx, nil).Once()
		result, err := discovery.RegisterProvider(txConfig, params)
		require.ErrorContains(t, err, "error parsing register event")
		assert.Equal(t, int64(0), result)
		contractBinding.AssertExpectations(t)
	})
}

// nolint:funlen
func TestDiscoveryContractImpl_IsOperational(t *testing.T) {
	t.Run("is operational for pegin", func(t *testing.T) {
		contractBinding := &mocks.DiscoveryBindingMock{}
		discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
		t.Run("Success", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.PeginProvider), parsedAddress).Return(true, nil).Once()
			result, err := discovery.IsOperational(liquidity_provider.PeginProvider, parsedAddress.String())
			require.NoError(t, err)
			assert.True(t, result)
			contractBinding.AssertExpectations(t)
		})
		t.Run("Error handling on IsOperational call fail", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.PeginProvider), parsedAddress).Return(true, assert.AnError).Once()
			_, err := discovery.IsOperational(liquidity_provider.PeginProvider, parsedAddress.String())
			require.Error(t, err)
		})
		t.Run(invalidAddressTest, func(t *testing.T) {
			result, err := discovery.IsOperational(liquidity_provider.PeginProvider, test.AnyString)
			require.ErrorIs(t, err, blockchain.InvalidAddressError)
			assert.False(t, result)
		})
	})
	t.Run("is operational for pegout", func(t *testing.T) {
		contractBinding := &mocks.DiscoveryBindingMock{}
		discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
		t.Run("Success", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.PegoutProvider), parsedAddress).Return(true, nil).Once()
			result, err := discovery.IsOperational(liquidity_provider.PegoutProvider, parsedAddress.String())
			require.NoError(t, err)
			assert.True(t, result)
			contractBinding.AssertExpectations(t)
		})
		t.Run("Error handling on IsOperational call fail", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.PegoutProvider), parsedAddress).Return(true, assert.AnError).Once()
			_, err := discovery.IsOperational(liquidity_provider.PegoutProvider, parsedAddress.String())
			require.Error(t, err)
		})
		t.Run(invalidAddressTest, func(t *testing.T) {
			result, err := discovery.IsOperational(liquidity_provider.PegoutProvider, test.AnyString)
			require.ErrorIs(t, err, blockchain.InvalidAddressError)
			assert.False(t, result)
		})
	})
	t.Run("is operational for both", func(t *testing.T) {
		contractBinding := &mocks.DiscoveryBindingMock{}
		discovery := rootstock.NewDiscoveryContractImpl(dummyClient, test.AnyAddress, contractBinding, nil, rootstock.RetryParams{}, time.Duration(1), Abis)
		t.Run("Success", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.FullProvider), parsedAddress).Return(true, nil).Once()
			result, err := discovery.IsOperational(liquidity_provider.FullProvider, parsedAddress.String())
			require.NoError(t, err)
			assert.True(t, result)
			contractBinding.AssertExpectations(t)
		})
		t.Run("Error handling on IsOperational call fail", func(t *testing.T) {
			contractBinding.EXPECT().IsOperational(mock.Anything, uint8(liquidity_provider.FullProvider), parsedAddress).Return(true, assert.AnError).Once()
			_, err := discovery.IsOperational(liquidity_provider.FullProvider, parsedAddress.String())
			require.Error(t, err)
		})
		t.Run(invalidAddressTest, func(t *testing.T) {
			result, err := discovery.IsOperational(liquidity_provider.FullProvider, test.AnyString)
			require.ErrorIs(t, err, blockchain.InvalidAddressError)
			assert.False(t, result)
		})
	})
}
