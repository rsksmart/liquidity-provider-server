package rootstock_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
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

func TestParseAddress(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		addresses := []string{
			"0x79568c2989232dCa1840087D73d403602364c0D4",
			"0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa8",
			"0x892813507Bf3aBF2890759d2135Ec34f4909Fea5",
		}
		var address common.Address
		for _, addr := range addresses {
			err := rootstock.ParseAddress(&address, addr)
			require.NoError(t, err)
			assert.Equal(t, addr, address.Hex())
		}
	})
	t.Run("Error", func(t *testing.T) {
		addresses := []string{
			"0x79568c2989232dCa1840087D73d403602364c0D",
			"0xd5f00ABfbEA7A0B193836CAc6833c2Ad9D06cEa",
			"0x892813507Bf3aBF2890759d2135Ec34f4909Fea",
			"no hex",
			"0x892813507Bf3aBF2890759d2135Ec34f4909Fea51",
			"0x892813507Bf3aBF2890759d2135Ec34f4909Fea",
		}
		var address common.Address
		for _, addr := range addresses {
			err := rootstock.ParseAddress(&address, addr)
			require.ErrorIs(t, err, blockchain.InvalidAddressError)
		}
	})
}

func TestRskClient_CheckConnection(t *testing.T) {
	clientBindingMock := &mocks.RpcClientBindingMock{}
	t.Run("Success", func(t *testing.T) {
		clientBindingMock.On("ChainID", test.AnyCtx).Return(big.NewInt(31), nil).Once()
		client := rootstock.NewRskClient(clientBindingMock)
		ok := client.CheckConnection(context.Background())
		require.True(t, ok)
	})
	t.Run("Error", func(t *testing.T) {
		clientBindingMock.On("ChainID", test.AnyCtx).Return(nil, assert.AnError).Once()
		client := rootstock.NewRskClient(clientBindingMock)
		ok := client.CheckConnection(context.Background())
		require.False(t, ok)
	})
}

func TestRskClient_Rpc(t *testing.T) {
	clientBindingMock := &mocks.RpcClientBindingMock{}
	client := rootstock.NewRskClient(clientBindingMock)
	assert.Equal(t, clientBindingMock, client.Rpc())
}

func TestRskClient_Shutdown(t *testing.T) {
	clientBindingMock := &mocks.RpcClientBindingMock{}
	clientBindingMock.On("Close").Once()
	endChannel := make(chan bool, 1)
	client := rootstock.NewRskClient(clientBindingMock)
	client.Shutdown(endChannel)
	assert.True(t, <-endChannel)
	clientBindingMock.AssertCalled(t, "Close")
}

// Since the function is private, it will be tested through HashPeginQuote
func TestRskRetry(t *testing.T) {
	const retries = 3
	lbcMock := &mocks.LbcAdapterMock{}
	lbc := rootstock.NewLiquidityBridgeContractImpl(
		dummyClient,
		test.AnyAddress,
		lbcMock,
		nil,
		rootstock.RetryParams{Retries: retries, Sleep: 1 * time.Second},
	)
	t.Run("Error on every attempt", func(t *testing.T) {
		lbcMock.On("HashQuote", mock.Anything, mock.Anything).Return(nil, assert.AnError).Times(retries)
		start := time.Now()
		result, err := lbc.HashPeginQuote(peginQuote)
		end := time.Now()
		assert.WithinRange(t, end, start, start.Add(3*time.Second).Add(500*time.Millisecond))
		require.Error(t, err)
		assert.Empty(t, result)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Error first attempt", func(t *testing.T) {
		lbcMock.On("HashQuote", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		lbcMock.On("HashQuote", mock.Anything, mock.Anything).Return([32]byte{1, 2, 3}, nil).Once()
		result, err := lbc.HashPeginQuote(peginQuote)
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		lbcMock.AssertExpectations(t)
	})
}
