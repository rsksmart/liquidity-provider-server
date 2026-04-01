package pegout_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// allAtOnceWatchedQuotes has two quotes with a combined total of 700 Wei
// (quote A: 300+50+50=400, quote B: 200+30+70=300).
var allAtOnceWatchedQuotes = []quote.WatchedPegoutQuote{
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "aao-01", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:   entities.NewWei(300),
			CallFee: entities.NewWei(50),
			GasFee:  entities.NewWei(50),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "aao-02", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:   entities.NewWei(200),
			CallFee: entities.NewWei(30),
			GasFee:  entities.NewWei(70),
		},
	},
}

const allAtOnceTotal = int64(700)

func allAtOnceConfig(bridgeMin int64) liquidity_provider.PegoutConfiguration {
	return liquidity_provider.PegoutConfiguration{BridgeTransactionMin: entities.NewWei(bridgeMin)}
}

func allAtOnceReceipt() blockchain.TransactionReceipt {
	return blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		Status:          true,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
}

func allAtOnceRevertedReceipt() blockchain.TransactionReceipt {
	return blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		Status:          false,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
}

func matchAllAtOnceSuccessQuotes(quotes []quote.RetainedPegoutQuote) bool {
	for _, q := range quotes {
		checks := []func() bool{
			func() bool { return q.State == quote.PegoutStateBridgeTxSucceeded },
			func() bool { return q.RemainingToRefund != nil },
			func() bool { return q.RemainingToRefund.Cmp(entities.NewWei(0)) == 0 },
			func() bool { return q.BridgeRefundTxHash == test.AnyHash },
			func() bool { return q.BridgeRefundGasUsed == uint64(21000) },
			func() bool { return q.BridgeRefundGasPrice != nil },
			func() bool { return q.BridgeRefundGasPrice.Cmp(entities.NewWei(pegout.BridgeConversionGasPrice)) == 0 },
			func() bool { return len(q.BridgeRebalances) == 1 },
			func() bool { return q.BridgeRebalances[0].TxHash == test.AnyHash },
		}

		for _, check := range checks {
			if !check() {
				return false
			}
		}
	}
	return true
}

func matchAllAtOnceFailedQuotes(quotes []quote.RetainedPegoutQuote) bool {
	for _, q := range quotes {
		if q.State != quote.PegoutStateBridgeTxFailed ||
			q.BridgeRefundTxHash != "" ||
			q.BridgeRefundGasUsed != uint64(0) ||
			q.BridgeRefundGasPrice != nil ||
			len(q.BridgeRebalances) != 0 {
			return false
		}
	}
	return true
}

func matchAllAtOnceRevertedQuotes(quotes []quote.RetainedPegoutQuote) bool {
	for _, q := range quotes {
		if q.State != quote.PegoutStateBridgeTxFailed ||
			q.BridgeRefundTxHash != test.AnyHash ||
			len(q.BridgeRebalances) != 1 ||
			q.BridgeRebalances[0].TxHash != test.AnyHash {
			return false
		}
	}
	return true
}

func newAllAtOnceHandler(repo quote.PegoutQuoteRepository, wallet blockchain.RootstockWallet, bridge rootstock.Bridge, mutex *mocks.MutexMock) *pegout.AllAtOnceHandler {
	return pegout.NewAllAtOnceHandler(repo, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
}

//nolint:funlen
func TestAllAtOnceHandler_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(
			entities.NewWei(allAtOnceTotal),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(matchBridgeTxConfig(entities.NewWei(allAtOnceTotal))),
			test.AnyAddress,
		).Return(allAtOnceReceipt(), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		repo.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(matchAllAtOnceSuccessQuotes)).
			Return(nil).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.NoError(t, err)
		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		bridge.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("total below minimum", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(5000), testQuotes)

		require.ErrorIs(t, err, usecases.TxBelowMinimumError)
		wallet.AssertNotCalled(t, "GetBalance")
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuotes")
		bridge.AssertNotCalled(t, "GetAddress")
		mutex.AssertExpectations(t)
	})

	t.Run("balance check error", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		wallet.On("GetBalance", mock.Anything).Return((*entities.Wei)(nil), assert.AnError).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorIs(t, err, assert.AnError)
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuotes")
		bridge.AssertNotCalled(t, "GetAddress")
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		// balance covers gas but not the transfer amount
		lowBalance := new(entities.Wei).Add(
			entities.NewWei(100),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(lowBalance, nil).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorIs(t, err, usecases.InsufficientAmountError)
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuotes")
		bridge.AssertNotCalled(t, "GetAddress")
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("bridge tx fails", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(
			entities.NewWei(allAtOnceTotal),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
			Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		bridge.On("GetAddress").Return(test.AnyAddress).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		repo.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(matchAllAtOnceFailedQuotes)).
			Return(nil).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorIs(t, err, assert.AnError)
		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		bridge.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("bridge tx reverted", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(
			entities.NewWei(allAtOnceTotal),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(matchBridgeTxConfig(entities.NewWei(allAtOnceTotal))),
			test.AnyAddress,
		).Return(allAtOnceRevertedReceipt(), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		repo.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(matchAllAtOnceRevertedQuotes)).
			Return(nil).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorContains(t, err, "transaction reverted")
		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		bridge.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("quotes update fails", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(
			entities.NewWei(allAtOnceTotal),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
			Return(allAtOnceReceipt(), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		updateErr := errors.New("update error")
		repo.On("UpdateRetainedQuotes", mock.Anything, mock.Anything).Return(updateErr).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorContains(t, err, "update error")
		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		bridge.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("bridge tx fails and quotes update fails", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(
			entities.NewWei(allAtOnceTotal),
			entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice),
		)
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		txErr := errors.New("tx error")
		wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
			Return(blockchain.TransactionReceipt{}, txErr).Once()
		bridge.On("GetAddress").Return(test.AnyAddress).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		updateErr := errors.New("update error")
		repo.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(matchAllAtOnceFailedQuotes)).
			Return(updateErr).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(allAtOnceWatchedQuotes))
		copy(testQuotes, allAtOnceWatchedQuotes)
		handler := newAllAtOnceHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), allAtOnceConfig(500), testQuotes)

		require.ErrorContains(t, err, "tx error")
		require.ErrorContains(t, err, "update error")
		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		bridge.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})
}
