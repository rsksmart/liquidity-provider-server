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

// Two quotes: A total=400 (300+50+50), B total=300 (200+30+70), combined=700.
var utxoSplitQuotes = []quote.WatchedPegoutQuote{
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "us-01", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:   entities.NewWei(300),
			CallFee: entities.NewWei(50),
			GasFee:  entities.NewWei(50),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "us-02", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:   entities.NewWei(200),
			CallFee: entities.NewWei(30),
			GasFee:  entities.NewWei(70),
		},
	},
}

const utxoSplitTotal = int64(700)

func utxoSplitConfig(bridgeMin int64) liquidity_provider.PegoutConfiguration {
	return liquidity_provider.PegoutConfiguration{BridgeTransactionMin: entities.NewWei(bridgeMin)}
}

func utxoSplitReceipt(txHash string, value int64) blockchain.TransactionReceipt {
	return blockchain.TransactionReceipt{
		TransactionHash: txHash,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(value),
		Status:          true,
	}
}

func newUtxoSplitHandler(repo quote.PegoutQuoteRepository, wallet blockchain.RootstockWallet, bridge rootstock.Bridge, mutex *mocks.MutexMock) *pegout.UtxoSplitHandler {
	return pegout.NewUtxoSplitHandler(repo, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
}

func utxoSplitGasPerTx() int64 {
	return int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
}

//nolint:funlen,maintidx
func TestUtxoSplitHandler_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(2*utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		// bridgeMin=300, total=700 → N=2, remainder=100, chunks=[400, 300]
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(400)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xtx1", 400), nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(300)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xtx2", 300), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		var updatedQuotes []quote.RetainedPegoutQuote
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				q, ok := args.Get(1).(quote.RetainedPegoutQuote)
				require.True(t, ok, "expected quote.RetainedPegoutQuote")
				updatedQuotes = append(updatedQuotes, q)
			}).Return(nil)

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(300), testQuotes)

		require.NoError(t, err)

		q1 := findUpdatedQuote(updatedQuotes, "us-01")
		require.NotNil(t, q1)
		assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xtx1", q1.BridgeRefundTxHash)
		assert.Len(t, q1.BridgeRebalances, 1)
		assert.Equal(t, "0xtx1", q1.BridgeRebalances[0].TxHash)

		q2 := findUpdatedQuote(updatedQuotes, "us-02")
		require.NotNil(t, q2)
		assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xtx2", q2.BridgeRefundTxHash)
		assert.Len(t, q2.BridgeRebalances, 1)
		assert.Equal(t, "0xtx2", q2.BridgeRebalances[0].TxHash)

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

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(5000), testQuotes)

		require.ErrorIs(t, err, usecases.TxBelowMinimumError)
		wallet.AssertNotCalled(t, "GetBalance")
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuote")
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

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(500), testQuotes)

		require.Error(t, err)
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuote")
		bridge.AssertNotCalled(t, "GetAddress")
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		// N=2 chunks need 2*gasPerTx, only provide 1
		lowBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(lowBalance, nil).Once()
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(300), testQuotes)

		require.ErrorIs(t, err, usecases.InsufficientAmountError)
		wallet.AssertNotCalled(t, "SendRbtc")
		repo.AssertNotCalled(t, "UpdateRetainedQuote")
		bridge.AssertNotCalled(t, "GetAddress")
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("second chunk fails", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(2*utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(400)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xtx1", 400), nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(300)) == 0 }),
			test.AnyAddress,
		).Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		var updatedQuotes []quote.RetainedPegoutQuote
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				q, ok := args.Get(1).(quote.RetainedPegoutQuote)
				require.True(t, ok, "expected quote.RetainedPegoutQuote")
				updatedQuotes = append(updatedQuotes, q)
			}).Return(nil)

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(300), testQuotes)

		require.NoError(t, err)
		q1 := findUpdatedQuote(updatedQuotes, "us-01")
		require.NotNil(t, q1)
		assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))

		q2 := findUpdatedQuote(updatedQuotes, "us-02")
		assert.Nil(t, q2, "us-02 should not be updated because the second chunk failed")

		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("all chunks fail", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(2*utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
			Return(blockchain.TransactionReceipt{}, assert.AnError)
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(300), testQuotes)

		require.NoError(t, err)
		repo.AssertNotCalled(t, "UpdateRetainedQuote")
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("DB update failure", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		// bridgeMin=600, total=700 → N=1, chunks=[700]
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(700)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xdbfail", 700), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Return(errors.New("db connection lost")).Once()

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(600), testQuotes)

		require.ErrorContains(t, err, "db connection lost")
		wallet.AssertExpectations(t)
		repo.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("chunk covers two quotes", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(utxoSplitTotal), entities.NewWei(utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		// bridgeMin=600, total=700 → N=1, chunks=[700]
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(700)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xspan", 700), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		var updatedQuotes []quote.RetainedPegoutQuote
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				q, ok := args.Get(1).(quote.RetainedPegoutQuote)
				require.True(t, ok, "expected quote.RetainedPegoutQuote")
				updatedQuotes = append(updatedQuotes, q)
			}).Return(nil)

		testQuotes := make([]quote.WatchedPegoutQuote, len(utxoSplitQuotes))
		copy(testQuotes, utxoSplitQuotes)
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(600), testQuotes)

		require.NoError(t, err)

		q1 := findUpdatedQuote(updatedQuotes, "us-01")
		require.NotNil(t, q1)
		assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xspan", q1.BridgeRefundTxHash)
		assert.Len(t, q1.BridgeRebalances, 1)

		q2 := findUpdatedQuote(updatedQuotes, "us-02")
		require.NotNil(t, q2)
		assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xspan", q2.BridgeRefundTxHash)
		assert.Len(t, q2.BridgeRebalances, 1)

		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("quote spans two chunks", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		walletBalance := new(entities.Wei).Add(entities.NewWei(800), entities.NewWei(2*utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		// bridgeMin=300, total=800 → N=2, remainder=200, chunks=[500, 300]
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(500)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xmulti1", 500), nil).Once()
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(300)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xmulti2", 300), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		var updatedQuotes []quote.RetainedPegoutQuote
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				q, ok := args.Get(1).(quote.RetainedPegoutQuote)
				require.True(t, ok, "expected quote.RetainedPegoutQuote")
				updatedQuotes = append(updatedQuotes, q)
			}).Return(nil)

		customQuotes := []quote.WatchedPegoutQuote{
			{
				RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "multi-01", State: quote.PegoutStateRefundPegOutSucceeded},
				PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(800), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
			},
		}
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(300), customQuotes)

		require.NoError(t, err)
		require.Len(t, updatedQuotes, 2)

		assert.Equal(t, "multi-01", updatedQuotes[0].QuoteHash)
		assert.NotEqual(t, quote.PegoutStateBridgeTxSucceeded, updatedQuotes[0].State)
		assert.Equal(t, 0, updatedQuotes[0].RemainingToRefund.Cmp(entities.NewWei(300)))
		assert.Equal(t, "0xmulti1", updatedQuotes[0].BridgeRefundTxHash)
		assert.Len(t, updatedQuotes[0].BridgeRebalances, 1)

		assert.Equal(t, "multi-01", updatedQuotes[1].QuoteHash)
		assert.Equal(t, quote.PegoutStateBridgeTxSucceeded, updatedQuotes[1].State)
		assert.Equal(t, 0, updatedQuotes[1].RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xmulti1", updatedQuotes[1].BridgeRefundTxHash)
		assert.Len(t, updatedQuotes[1].BridgeRebalances, 2)
		assert.Equal(t, "0xmulti1", updatedQuotes[1].BridgeRebalances[0].TxHash)
		assert.Equal(t, "0xmulti2", updatedQuotes[1].BridgeRebalances[1].TxHash)

		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})

	t.Run("retry with pre-existing partial refund", func(t *testing.T) {
		repo := &mocks.PegoutQuoteRepositoryMock{}
		wallet := &mocks.RskWalletMock{}
		bridge := &mocks.BridgeMock{}
		mutex := &mocks.MutexMock{}

		// retry-01: Total=500, RemainingToRefund=200, retry-02: Total=300 → adjusted=500
		walletBalance := new(entities.Wei).Add(entities.NewWei(500), entities.NewWei(utxoSplitGasPerTx()))
		wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
		// bridgeMin=400, adjusted=500 → N=1, chunks=[500]
		wallet.On("SendRbtc", mock.Anything,
			mock.MatchedBy(func(c blockchain.TransactionConfig) bool { return c.Value.Cmp(entities.NewWei(500)) == 0 }),
			test.AnyAddress,
		).Return(utxoSplitReceipt("0xretry", 500), nil).Once()
		bridge.On("GetAddress").Return(test.AnyAddress)
		mutex.On("Lock").Return().Once()
		mutex.On("Unlock").Return().Once()

		var updatedQuotes []quote.RetainedPegoutQuote
		repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				q, ok := args.Get(1).(quote.RetainedPegoutQuote)
				require.True(t, ok, "expected quote.RetainedPegoutQuote")
				updatedQuotes = append(updatedQuotes, q)
			}).Return(nil)

		customQuotes := []quote.WatchedPegoutQuote{
			{
				RetainedQuote: quote.RetainedPegoutQuote{
					QuoteHash:          "retry-01",
					State:              quote.PegoutStateRefundPegOutSucceeded,
					RemainingToRefund:  entities.NewWei(200),
					BridgeRefundTxHash: "0xprev",
					BridgeRebalances: []quote.BridgeRebalanceAllocation{
						{TxHash: "0xprev", GasUsed: 21000, GasPrice: entities.NewWei(pegout.BridgeConversionGasPrice)},
					},
				},
				PegoutQuote: quote.PegoutQuote{Value: entities.NewWei(500), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
			},
			{
				RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "retry-02", State: quote.PegoutStateRefundPegOutSucceeded},
				PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
			},
		}
		handler := newUtxoSplitHandler(repo, wallet, bridge, mutex)
		err := handler.Execute(context.Background(), utxoSplitConfig(400), customQuotes)

		require.NoError(t, err)

		q1 := findUpdatedQuote(updatedQuotes, "retry-01")
		require.NotNil(t, q1)
		assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xprev", q1.BridgeRefundTxHash)
		assert.Len(t, q1.BridgeRebalances, 2)
		assert.Equal(t, "0xprev", q1.BridgeRebalances[0].TxHash)
		assert.Equal(t, "0xretry", q1.BridgeRebalances[1].TxHash)

		q2 := findUpdatedQuote(updatedQuotes, "retry-02")
		require.NotNil(t, q2)
		assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
		assert.Equal(t, "0xretry", q2.BridgeRefundTxHash)
		assert.Len(t, q2.BridgeRebalances, 1)

		repo.AssertExpectations(t)
		wallet.AssertExpectations(t)
		mutex.AssertExpectations(t)
	})
}
