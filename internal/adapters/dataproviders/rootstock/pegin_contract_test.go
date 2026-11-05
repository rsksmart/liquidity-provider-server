package rootstock_test

import (
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
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

func TestNewPeginContractImpl(t *testing.T) {
	contract := rootstock.NewPeginContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		test.AnyAddress,
		createBoundContractMock().contract,
		&mocks.TransactionSignerMock{},
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		time.Duration(1),
		bindings.NewPeginContract(),
		Abis,
	)
	test.AssertNonZeroValues(t, contract)
}

func TestPeginContractImpl_GetBalance(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackGetBalance(parsedAddress)),
			mock.Anything,
		).Return(mustPackUint256(t, big.NewInt(600)), nil).Once()
		result, err := peginContract.GetBalance(parsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(600), result)
		contractMock.caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetBalance call error", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackGetBalance(parsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := peginContract.GetBalance(parsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting balance", func(t *testing.T) {
		result, err := peginContract.GetBalance(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

// nolint:funlen
func TestPeginContractImpl_GetAddress(t *testing.T) {
	peginContract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, nil, nil, rootstock.RetryParams{}, time.Duration(1), nil, Abis)
	assert.Equal(t, test.AnyAddress, peginContract.GetAddress())
}

// nolint:funlen
func TestPeginContractImpl_Withdraw(t *testing.T) {
	contractMock := createBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	t.Run("Success", func(t *testing.T) {
		txData := peginBinding.PackWithdraw(withdrawAmount.AsBigInt())
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(txData),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(
			mock.Anything,
			matchTransaction(contractMock.transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), txData),
		).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := peginContract.Withdraw(withdrawAmount)
		require.NoError(t, err)
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginContractImpl_Withdraw_ErrorHandling(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	signerMock.On("Address").Return(parsedAddress)
	t.Run("Error handling (dry-run revert NoBalance)", func(t *testing.T) {
		e := NewRskRpcError("transaction reverted", "0x29226653")
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, e).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw reverted with: NoBalance")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (dry-run parse error)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "error parsing withdraw result")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (transaction send error)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(assert.AnError).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, true)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
	t.Run("Error handling (transaction reverted)", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(nil).Once()
		prepareTxMocks(&contractMock, mockClient, signerMock, false)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error: transaction failed")
		contractMock.caller.AssertExpectations(t)
		contractMock.transactor.AssertExpectations(t)
	})
}

func TestPeginContractImpl_PausedStatus(t *testing.T) {
	contractMock := createBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	contract := rootstock.NewPeginContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, Abis)
	t.Run("should return pause status result", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(mustPackPauseStatus(t, generalPauseStatus{IsPaused: true, Reason: "test", Since: 123}), nil).Once()
		result, err := contract.PausedStatus()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseStatus{IsPaused: true, Reason: "test", Since: 123}, result)
	})
	t.Run("should handle error checking pause status", func(t *testing.T) {
		contractMock.caller.EXPECT().CallContract(
			mock.Anything,
			matchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := contract.PausedStatus()
		require.Error(t, err)
		assert.Empty(t, result)
	})
	contractMock.caller.AssertExpectations(t)
}
