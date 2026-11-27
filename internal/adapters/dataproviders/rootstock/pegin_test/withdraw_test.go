package pegin_test

import (
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

// nolint:funlen
func TestPeginContractImpl_Withdraw(t *testing.T) {
	contractMock := rootstocktest.NewBoundContractMock()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.Contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, commitFirstBindings.NewPeginCommitFirstContract(), rootstocktest.Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	t.Run("Success", func(t *testing.T) {
		txData := peginBinding.PackWithdraw(withdrawAmount.AsBigInt())
		contractMock.Caller.EXPECT().CallContract(
			mock.Anything,
			rootstocktest.MatchCallData(txData),
			mock.Anything,
		).Return(nil, nil).Once()
		contractMock.Transactor.EXPECT().SendTransaction(
			mock.Anything,
			rootstocktest.MatchTransaction(contractMock.Transactor, common.HexToAddress(test.AnyRskAddress), 0, big.NewInt(0), txData),
		).Return(nil).Once()
		rootstocktest.PrepareTxMocks(&contractMock, mockClient, signerMock, true)
		err := peginContract.Withdraw(withdrawAmount)
		require.NoError(t, err)
		contractMock.Caller.AssertExpectations(t)
		contractMock.Transactor.AssertExpectations(t)
	})
}

// nolint:funlen
func TestPeginContractImpl_Withdraw_ErrorHandling(t *testing.T) {
	contractMock := rootstocktest.NewBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	signerMock := &mocks.TransactionSignerMock{}
	mockClient := &mocks.RpcClientBindingMock{}
	peginContract := rootstock.NewPeginContractImpl(rootstock.NewRskClient(mockClient), test.AnyAddress, contractMock.Contract, signerMock, rootstock.RetryParams{}, time.Duration(1), peginBinding, commitFirstBindings.NewPeginCommitFirstContract(), rootstocktest.Abis)
	withdrawAmount := entities.NewWei(5000000000000000000)
	signerMock.On("Address").Return(rootstocktest.ParsedAddress)
	t.Run("Error handling (dry-run revert NoBalance)", func(t *testing.T) {
		e := rootstocktest.NewRskRpcError("transaction reverted", "0x29226653")
		contractMock.Caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, e).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw reverted with: NoBalance")
		contractMock.Caller.AssertExpectations(t)
		contractMock.Transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (dry-run parse error)", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "error parsing withdraw result")
		contractMock.Caller.AssertExpectations(t)
		contractMock.Transactor.AssertNotCalled(t, "Transact")
	})
	t.Run("Error handling (transaction send error)", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.Transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(assert.AnError).Once()
		rootstocktest.PrepareTxMocks(&contractMock, mockClient, signerMock, true)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.Transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.Transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.Transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error")
		contractMock.Caller.AssertExpectations(t)
		contractMock.Transactor.AssertExpectations(t)
	})
	t.Run("Error handling (transaction reverted)", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()
		contractMock.Transactor.EXPECT().SendTransaction(mock.Anything, mock.Anything).Return(nil).Once()
		rootstocktest.PrepareTxMocks(&contractMock, mockClient, signerMock, false)
		signerMock.EXPECT().Sign(mock.Anything, mock.Anything).RunAndReturn(func(addr common.Address, tx *geth.Transaction) (*geth.Transaction, error) {
			return tx, nil
		})
		contractMock.Transactor.EXPECT().HeaderByNumber(mock.Anything, mock.Anything).Return(&geth.Header{}, nil).Once()
		contractMock.Transactor.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1), nil).Once()
		contractMock.Transactor.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(1, nil).Once()
		err := peginContract.Withdraw(withdrawAmount)
		require.ErrorContains(t, err, "withdraw error: transaction failed")
		contractMock.Caller.AssertExpectations(t)
		contractMock.Transactor.AssertExpectations(t)
	})
}
