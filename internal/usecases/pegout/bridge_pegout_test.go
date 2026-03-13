package pegout_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var bridgePegoutTestWatchedQuotes = []quote.WatchedPegoutQuote{
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "01", State: quote.PegoutStateSendPegoutFailed},
		PegoutQuote: quote.PegoutQuote{
			Value:      entities.NewWei(100),
			CallFee:    entities.NewWei(10),
			PenaltyFee: entities.NewWei(5),
			GasFee:     entities.NewWei(30),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "02", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:      entities.NewWei(77),
			CallFee:    entities.NewWei(32),
			PenaltyFee: entities.NewWei(5),
			GasFee:     entities.NewWei(55),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "03", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:      entities.NewWei(123),
			CallFee:    entities.NewWei(8),
			PenaltyFee: entities.NewWei(1),
			GasFee:     entities.NewWei(3),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "04", State: quote.PegoutStateWaitingForDeposit},
		PegoutQuote: quote.PegoutQuote{
			Value:      entities.NewWei(1000),
			CallFee:    entities.NewWei(11),
			PenaltyFee: entities.NewWei(7),
			GasFee:     entities.NewWei(210),
		},
	},
	{
		RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "05", State: quote.PegoutStateRefundPegOutSucceeded},
		PegoutQuote: quote.PegoutQuote{
			Value:      entities.NewWei(200),
			CallFee:    entities.NewWei(20),
			PenaltyFee: entities.NewWei(15),
			GasFee:     entities.NewWei(40),
		},
	},
}

func TestBridgePegoutUseCase_Run(t *testing.T) {
	t.Run("make bridge pegout successfully", func(t *testing.T) {
		testBridgePegoutUseCaseSuccess(t)
	})
	t.Run("when the total value to pegout is below the minimum", func(t *testing.T) {
		testBridgePegoutUseCaseValueBelowMinimum(t)
	})
	t.Run("when some of the quotes have not been refunded", func(t *testing.T) {
		testBridgePegoutUseCaseQuotesNotRefunded(t)
	})
	t.Run("error getting wallet balance", func(t *testing.T) {
		testBridgePegoutUseCaseWalletBalanceError(t)
	})
	t.Run("wallet doesn't have enough balance", func(t *testing.T) {
		testBridgePegoutUseCaseWalletWithoutBalance(t)
	})
	t.Run("bridge tx fails", func(t *testing.T) {
		testBridgePegoutUseCaseTxFails(t)
	})
	t.Run("quotes update fails", func(t *testing.T) {
		testBridgePegoutUseCaseUpdateFails(t)
	})
}

func testBridgePegoutUseCaseSuccess(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	sendRbtcReceipt := blockchain.TransactionReceipt{
		TransactionHash:   test.AnyHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              "0x123",
		To:                test.AnyAddress,
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(558),
		GasPrice:          entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(558)) == 0 &&
			*config.GasLimit == pegout.BridgeConversionGasLimit &&
			config.GasPrice.Cmp(entities.NewWei(pegout.BridgeConversionGasPrice)) == 0
	}), test.AnyAddress).Return(sendRbtcReceipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress).Once()
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(550),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if !(q.State == quote.PegoutStateBridgeTxSucceeded &&
				q.BridgeRefundTxHash == test.AnyHash &&
				q.BridgeRefundGasUsed == uint64(21000) &&
				q.BridgeRefundGasPrice != nil && q.BridgeRefundGasPrice.Cmp(entities.NewWei(pegout.BridgeConversionGasPrice)) == 0) {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.NoError(t, err)
	pegoutRepository.AssertExpectations(t)
	pegoutLp.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func testBridgePegoutUseCaseValueBelowMinimum(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	mutex := &mocks.MutexMock{}
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(5000),
	}).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.ErrorIs(t, err, usecases.TxBelowMinimumError)
	pegoutRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	pegoutLp.AssertExpectations(t)
	wallet.AssertNotCalled(t, "GetBalance")
	wallet.AssertNotCalled(t, "SendRbtc")
	mutex.AssertNotCalled(t, "Unlock")
	mutex.AssertNotCalled(t, "Lock")
	bridge.AssertNotCalled(t, "GetAddress")
}

func testBridgePegoutUseCaseQuotesNotRefunded(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	mutex := &mocks.MutexMock{}
	bridge := &mocks.BridgeMock{}
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	err := useCase.Run(context.Background(), bridgePegoutTestWatchedQuotes...)
	require.ErrorContains(t, err, "not all quotes were refunded successfully")
	pegoutRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	pegoutLp.AssertNotCalled(t, "PegoutConfiguration")
	wallet.AssertNotCalled(t, "GetBalance")
	wallet.AssertNotCalled(t, "SendRbtc")
	mutex.AssertNotCalled(t, "Unlock")
	mutex.AssertNotCalled(t, "Lock")
	bridge.AssertNotCalled(t, "GetAddress")
}

func testBridgePegoutUseCaseWalletBalanceError(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	wallet.On("GetBalance", mock.Anything).Return((*entities.Wei)(nil), assert.AnError).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(550),
	}).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.Error(t, err)
	pegoutRepository.AssertNotCalled(t, "UpdateRetainedQuotes")
	pegoutLp.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertNotCalled(t, "GetAddress")
}

func testBridgePegoutUseCaseWalletWithoutBalance(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	walletBalance := new(entities.Wei).Add(entities.NewWei(500), entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(550),
	}).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	pegoutRepository.AssertNotCalled(t, "UpdateRetainedQuotes")
	pegoutLp.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertNotCalled(t, "GetAddress")
}

func testBridgePegoutUseCaseTxFails(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	emptyReceipt := blockchain.TransactionReceipt{}
	wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).Return(emptyReceipt, assert.AnError).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress).Once()
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(550),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if !(q.State == quote.PegoutStateBridgeTxFailed &&
				q.BridgeRefundTxHash == "" &&
				q.BridgeRefundGasUsed == uint64(0) &&
				q.BridgeRefundGasPrice == nil) {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.Error(t, err)
	pegoutRepository.AssertExpectations(t)
	pegoutLp.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func testBridgePegoutUseCaseUpdateFails(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(pegout.BridgeConversionGasLimit*pegout.BridgeConversionGasPrice))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	successReceipt := blockchain.TransactionReceipt{
		TransactionHash:   test.AnyHash,
		BlockHash:         "0xblock123",
		BlockNumber:       uint64(1000),
		From:              "0x123",
		To:                test.AnyAddress,
		CumulativeGasUsed: big.NewInt(21000),
		GasUsed:           big.NewInt(21000),
		Value:             entities.NewWei(0),
		GasPrice:          entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).Return(successReceipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress).Once()
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(550),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.Anything).Return(errors.New("update error")).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.AllAtOnce)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(
		context.Background(),
		testQuotes[1],
		testQuotes[2],
		testQuotes[4],
	)
	require.ErrorContains(t, err, "update error")
	pegoutRepository.AssertExpectations(t)
	pegoutLp.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
	bridge.AssertExpectations(t)
}

func TestParseRebalanceStrategy(t *testing.T) {
	t.Run("parse UTXO_SPLIT", func(t *testing.T) {
		strategy, err := pegout.ParseRebalanceStrategy("UTXO_SPLIT")
		require.NoError(t, err)
		assert.Equal(t, pegout.UtxoSplit, strategy)
	})
	t.Run("parse ALL_AT_ONCE", func(t *testing.T) {
		strategy, err := pegout.ParseRebalanceStrategy("ALL_AT_ONCE")
		require.NoError(t, err)
		assert.Equal(t, pegout.AllAtOnce, strategy)
	})
	t.Run("empty string returns error", func(t *testing.T) {
		_, err := pegout.ParseRebalanceStrategy("")
		require.Error(t, err)
	})
	t.Run("unknown value returns error", func(t *testing.T) {
		_, err := pegout.ParseRebalanceStrategy("UNKNOWN")
		require.Error(t, err)
	})
}

func TestBridgePegoutUseCase_UtxoSplit(t *testing.T) {
	t.Run("split into correct number of txs", func(t *testing.T) {
		testUtxoSplitSuccess(t)
	})
	t.Run("no split when N=1", func(t *testing.T) {
		testUtxoSplitNoSplitWhenN1(t)
	})
	t.Run("below minimum", func(t *testing.T) {
		testUtxoSplitBelowMinimum(t)
	})
	t.Run("exact multiple", func(t *testing.T) {
		testUtxoSplitExactMultiple(t)
	})
	t.Run("fail mid-split", func(t *testing.T) {
		testUtxoSplitFailMidSplit(t)
	})
	t.Run("balance check with multi-tx gas", func(t *testing.T) {
		testUtxoSplitInsufficientGas(t)
	})
}

// total=558, BridgeMin=200 => N=2, R=158 => 1st tx: 358, 2nd tx: 200
func testUtxoSplitSuccess(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(2*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt1 := blockchain.TransactionReceipt{
		TransactionHash: "0xtx1",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	receipt2 := blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	// First tx: 358
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(358)) == 0
	}), test.AnyAddress).Return(receipt1, nil).Once()
	// Second tx: 200
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(200)) == 0
	}), test.AnyAddress).Return(receipt2, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if q.State != quote.PegoutStateBridgeTxSucceeded || q.BridgeRefundTxHash != test.AnyHash {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	require.NoError(t, err)
	pegoutRepository.AssertExpectations(t)
	wallet.AssertExpectations(t)
	mutex.AssertExpectations(t)
}

// total=350, BridgeMin=200 => N=1, single tx of 350
func testUtxoSplitNoSplitWhenN1(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(500), entities.NewWei(gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(350)) == 0
	}), test.AnyAddress).Return(receipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if q.State != quote.PegoutStateBridgeTxSucceeded {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "n1-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(200), CallFee: entities.NewWei(50), GasFee: entities.NewWei(100)},
		},
	}
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)
	pegoutRepository.AssertExpectations(t)
	wallet.AssertExpectations(t)
}

// total=150, BridgeMin=200 => below minimum
func testUtxoSplitBelowMinimum(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	mutex := &mocks.MutexMock{}
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "bm-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(100), CallFee: entities.NewWei(30), GasFee: entities.NewWei(20)},
		},
	}
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	err := useCase.Run(context.Background(), customQuotes...)
	require.ErrorIs(t, err, usecases.TxBelowMinimumError)
	wallet.AssertNotCalled(t, "SendRbtc")
}

// total=600, BridgeMin=200 => N=3, R=0 => three txs of 200 each
func testUtxoSplitExactMultiple(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(3*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	// All three txs should be 200
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(200)) == 0
	}), test.AnyAddress).Return(receipt, nil).Times(3)
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if q.State != quote.PegoutStateBridgeTxSucceeded {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "em-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(100), GasFee: entities.NewWei(200)},
		},
	}
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)
	pegoutRepository.AssertExpectations(t)
	wallet.AssertExpectations(t)
}

// total=558, BridgeMin=200 => N=2, 2nd tx fails => all quotes BridgeTxFailed
func testUtxoSplitFailMidSplit(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(2*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt1 := blockchain.TransactionReceipt{
		TransactionHash: "0xtx1",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	emptyReceipt := blockchain.TransactionReceipt{}
	// First tx: 358 succeeds
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(358)) == 0
	}), test.AnyAddress).Return(receipt1, nil).Once()
	// Second tx: 200 fails
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(200)) == 0
	}), test.AnyAddress).Return(emptyReceipt, assert.AnError).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
		for _, q := range quotes {
			if q.State != quote.PegoutStateBridgeTxFailed {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	require.Error(t, err)
	pegoutRepository.AssertExpectations(t)
	wallet.AssertExpectations(t)
}

// total=558, BridgeMin=200 => N=2, needs 2*gasPerTx but balance only has 1*gasPerTx
func testUtxoSplitInsufficientGas(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	// Balance covers totalValue + 1 gas, but not 2 gas
	walletBalance := new(entities.Wei).Add(entities.NewWei(558), entities.NewWei(gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex, pegout.UtxoSplit)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	wallet.AssertNotCalled(t, "SendRbtc")
}

// TestUtxoSplit_AmountIntegrity verifies that the sum of all transaction amounts
// sent to the bridge equals the original totalValue exactly, with no rounding errors.
// Each sub-test captures every SendRbtc call's value and asserts sum == total.
func TestUtxoSplit_AmountIntegrity(t *testing.T) {
	rbtc := func(n int64) *big.Int {
		// n * 10^17 (0.1 RBTC units for convenience)
		v := big.NewInt(n)
		exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil)
		return v.Mul(v, exp)
	}

	cases := []struct {
		name      string
		total     *big.Int
		bridgeMin *big.Int
		wantN     int // expected number of SendRbtc calls
	}{
		{
			name:      "total equals bridgeMin exactly (N=1)",
			total:     rbtc(15), // 1.5 RBTC
			bridgeMin: rbtc(15),
			wantN:     1,
		},
		{
			name:      "total just above bridgeMin (N=1, remainder absorbed)",
			total:     rbtc(29), // 2.9 RBTC, min=1.5 => N=1
			bridgeMin: rbtc(15),
			wantN:     1,
		},
		{
			name:      "total exactly 2x bridgeMin (N=2, no remainder)",
			total:     rbtc(30), // 3.0 RBTC, min=1.5
			bridgeMin: rbtc(15),
			wantN:     2,
		},
		{
			name:      "total just below 2x bridgeMin (N=1)",
			total:     rbtc(29), // 2.9 RBTC, min=1.5 => N=1 (int div 29/15=1)
			bridgeMin: rbtc(15),
			wantN:     1,
		},
		{
			name:      "total 2x+1wei above bridgeMin boundary",
			total:     new(big.Int).Add(new(big.Int).Mul(rbtc(15), big.NewInt(2)), big.NewInt(1)), // 3.0 RBTC + 1 wei
			bridgeMin: rbtc(15),
			wantN:     2, // N=2, first=1.5 RBTC + 1 wei, second=1.5 RBTC
		},
		{
			name:      "5 transactions with large remainder in first chunk",
			total:     rbtc(83), // 8.3 RBTC, min=1.5 => N=5, first=2.3 RBTC (1.5+0.8), rest=1.5 RBTC each
			bridgeMin: rbtc(15),
			wantN:     5,
		},
		{
			name:      "3 transactions exact multiple",
			total:     rbtc(45), // 4.5 RBTC, min=1.5 => N=3, no remainder
			bridgeMin: rbtc(15),
			wantN:     3,
		},
		{
			name:      "large value: 100 RBTC split by 1.5 RBTC min",
			total:     rbtc(1000), // 100 RBTC, min=1.5 => N=66, remainder=1.0 RBTC
			bridgeMin: rbtc(15),
			wantN:     66,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAmountIntegritySubtest(t, tc.total, tc.bridgeMin, tc.wantN)
		})
	}
}

func runAmountIntegritySubtest(t *testing.T, total, bridgeMin *big.Int, wantN int) {
	t.Helper()
	totalWei := entities.NewBigWei(total)
	bridgeMinWei := entities.NewBigWei(bridgeMin)
	var sentAmounts []*big.Int
	var mu sync.Mutex
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(totalWei.Copy(), entities.NewWei(int64(wantN)*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash,
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
	}
	wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
		Run(func(args mock.Arguments) {
			config, ok := args.Get(1).(blockchain.TransactionConfig)
			require.True(t, ok)
			mu.Lock()
			sentAmounts = append(sentAmounts, config.Value.AsBigInt())
			mu.Unlock()
		}).Return(receipt, nil)
	walletMutex := &mocks.MutexMock{}
	walletMutex.On("Lock").Return()
	walletMutex.On("Unlock").Return()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: bridgeMinWei,
	}).Once()
	pegoutRepository.On("UpdateRetainedQuotes", mock.Anything, mock.Anything).Return(nil).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "integrity-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: totalWei, CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
		},
	}
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, walletMutex, pegout.UtxoSplit)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)
	require.Len(t, sentAmounts, wantN, "unexpected number of SendRbtc calls")
	sum := new(big.Int)
	for _, amount := range sentAmounts {
		sum.Add(sum, amount)
	}
	assert.Equal(t, 0, sum.Cmp(total),
		"sum of sent amounts (%s) != total (%s), diff = %s",
		sum.String(), total.String(), new(big.Int).Sub(sum, total).String())
	if wantN > 1 {
		for i, amount := range sentAmounts {
			assert.GreaterOrEqual(t, amount.Cmp(bridgeMin), 0,
				"tx %d amount %s is below bridgeMin %s", i, amount.String(), bridgeMin.String())
		}
	}
	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}
