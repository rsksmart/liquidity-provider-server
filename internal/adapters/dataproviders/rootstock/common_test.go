package rootstock_test

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
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
	contractBinding := &mocks.PeginContractAdapterMock{}
	contract := rootstock.NewPeginContractImpl(
		dummyClient,
		test.AnyAddress,
		contractBinding,
		nil,
		rootstock.RetryParams{Retries: retries, Sleep: 1 * time.Second},
		time.Duration(1),
		Abis,
	)
	t.Run("Error on every attempt", func(t *testing.T) {
		contractBinding.EXPECT().HashPegInQuote(mock.Anything, mock.Anything).Return([32]byte{}, assert.AnError).Times(retries)
		start := time.Now()
		result, err := contract.HashPeginQuote(peginQuote)
		end := time.Now()
		assert.WithinRange(t, end, start, start.Add(3*time.Second).Add(500*time.Millisecond))
		require.Error(t, err)
		assert.Empty(t, result)
		contractBinding.AssertExpectations(t)
	})
	t.Run("Error first attempt", func(t *testing.T) {
		contractBinding.EXPECT().HashPegInQuote(mock.Anything, mock.Anything).Return([32]byte{}, assert.AnError).Once()
		contractBinding.EXPECT().HashPegInQuote(mock.Anything, mock.Anything).Return([32]byte{1, 2, 3}, nil).Once()
		result, err := contract.HashPeginQuote(peginQuote)
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		contractBinding.AssertExpectations(t)
	})
}

func TestAwaitTxWithCtx(t *testing.T) {
	t.Run("should return receipt if tx is successful", func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		signerMock := &mocks.TransactionSignerMock{}
		tx := prepareTxMocks(clientMock, signerMock, true)
		defer test.AssertLogContains(t, fmt.Sprintf("Transaction success tx (%s) executed successfully", tx.Hash()))()
		receipt, err := rootstock.AwaitTxWithCtx(clientMock, time.Duration(1), "success tx", context.Background(), func() (*geth.Transaction, error) {
			return tx, nil
		})
		require.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, uint64(1), receipt.Status)
	})
	t.Run("should return receipt if tx reverts", func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		signerMock := &mocks.TransactionSignerMock{}
		tx := prepareTxMocks(clientMock, signerMock, false)
		defer test.AssertLogContains(t, fmt.Sprintf("Transaction fail tx (%s) reverted", tx.Hash()))()
		receipt, err := rootstock.AwaitTxWithCtx(clientMock, time.Duration(1), "fail tx", context.Background(), func() (*geth.Transaction, error) {
			return tx, nil
		})
		require.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Zero(t, receipt.Status)
	})
	t.Run("should return error if tx to be mined", func(t *testing.T) {
		clientMock := &mocks.RpcClientBindingMock{}
		signerMock := &mocks.TransactionSignerMock{}
		tx := prepareTxMocks(clientMock, signerMock, true)
		clientMock.ExpectedCalls = []*mock.Call{}
		clientMock.On("TransactionReceipt", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer test.AssertLogContains(t, fmt.Sprintf("Error waiting for transaction Test tx (%s) to be mined", tx.Hash()))()
		receipt, err := rootstock.AwaitTxWithCtx(clientMock, time.Duration(1), "Test tx", ctx, func() (*geth.Transaction, error) {
			return tx, nil
		})
		cancel()
		require.Error(t, err)
		assert.Nil(t, receipt)
	})
}

func TestParseRevertReason(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		result, err := rootstock.ParseRevertReason(Abis.PegOut, nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
	t.Run("non-DataError returns error", func(t *testing.T) {
		result, err := rootstock.ParseRevertReason(Abis.PegOut, assert.AnError)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("DataError with non-string ErrorData returns error", func(t *testing.T) {
		e := rskRpcErrorWithIntData{message: "revert", data: 42}
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("DataError with invalid hex returns error", func(t *testing.T) {
		e := NewRskRpcError("revert", "0xZZ")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("DataError with generic Error(string) ABI revert returns error", func(t *testing.T) {
		// Error("test") ABI-encoded: selector 0x08c379a0 + offset + length + data
		e := NewRskRpcError("revert", "0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("DataError with empty data returns ErrShortRevertData", func(t *testing.T) {
		e := NewRskRpcError("execution reverted", "0x")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.ErrorIs(t, err, rootstock.ErrShortRevertData)
		assert.Nil(t, result)
	})
	t.Run("DataError with data shorter than selector returns ErrShortRevertData", func(t *testing.T) {
		e := NewRskRpcError("execution reverted", "0xaabbcc")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.ErrorIs(t, err, rootstock.ErrShortRevertData)
		assert.Nil(t, result)
	})
	t.Run("DataError with unknown 4-byte selector returns error", func(t *testing.T) {
		e := NewRskRpcError("revert", "0xdeadbeef")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("DataError with known ABI error selector returns parsed error", func(t *testing.T) {
		// NotEnoughConfirmations selector from the pegout ABI with valid ABI-encoded arguments
		e := NewRskRpcError("transaction reverted", "0xd2506f8c00000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000002")
		result, err := rootstock.ParseRevertReason(Abis.PegOut, e)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "NotEnoughConfirmations", result.Name)
	})
}

// rskRpcErrorWithIntData is a DataError whose ErrorData returns a non-string value.
type rskRpcErrorWithIntData struct {
	message string
	data    int
}

func (r rskRpcErrorWithIntData) Error() string          { return r.message }
func (r rskRpcErrorWithIntData) ErrorData() interface{} { return r.data }
