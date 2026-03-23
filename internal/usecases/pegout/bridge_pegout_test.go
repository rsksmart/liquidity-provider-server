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
				q.BridgeRefundGasPrice != nil && q.BridgeRefundGasPrice.Cmp(entities.NewWei(pegout.BridgeConversionGasPrice)) == 0 &&
				len(q.BridgeRebalances) == 1 &&
				q.BridgeRebalances[0].TxHash == test.AnyHash) {
				return false
			}
		}
		return true
	})).Return(nil).Once()
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(5000),
	}).Once()
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	mutex.AssertExpectations(t)
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
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	handler := pegout.NewAllAtOnceHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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

func TestNewRebalanceHandler(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	wallet := &mocks.RskWalletMock{}
	bridge := &mocks.BridgeMock{}
	contracts := blockchain.RskContracts{Bridge: bridge}

	t.Run("ALL_AT_ONCE returns AllAtOnceHandler", func(t *testing.T) {
		handler := pegout.NewRebalanceHandler(pegout.AllAtOnce, pegoutRepository, wallet, contracts)
		assert.IsType(t, &pegout.AllAtOnceHandler{}, handler)
	})
	t.Run("UTXO_SPLIT returns UtxoSplitHandler", func(t *testing.T) {
		handler := pegout.NewRebalanceHandler(pegout.UtxoSplit, pegoutRepository, wallet, contracts)
		assert.IsType(t, &pegout.UtxoSplitHandler{}, handler)
	})
	t.Run("unknown value defaults to AllAtOnceHandler", func(t *testing.T) {
		handler := pegout.NewRebalanceHandler("UNKNOWN", pegoutRepository, wallet, contracts)
		assert.IsType(t, &pegout.AllAtOnceHandler{}, handler)
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

func matchRetainedQuotes(count int, state quote.PegoutState, txHash string) func([]quote.RetainedPegoutQuote) bool {
	return func(quotes []quote.RetainedPegoutQuote) bool {
		if len(quotes) != count {
			return false
		}
		for _, q := range quotes {
			if q.State != state || q.BridgeRefundTxHash != txHash {
				return false
			}
			if txHash != "" && (len(q.BridgeRebalances) == 0 || q.BridgeRebalances[0].TxHash != txHash) {
				return false
			}
		}
		return true
	}
}

// total=558, BridgeMin=200 => N=2, R=158 => 1st tx: 358, 2nd tx: 200
// quotes: q1(164), q2(134), q3(260)
// chunk1 (value=358): fills q1(164) fully, q2(134) fully, q3 partial (60)
// chunk2 (value=200): fills q3 remaining (200)
func testUtxoSplitSuccess(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(2*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt1 := blockchain.TransactionReceipt{
		TransactionHash: "0xtx1", GasUsed: big.NewInt(21000),
		GasPrice: entities.NewWei(pegout.BridgeConversionGasPrice),
		Value: entities.NewWei(358),
	}
	receipt2 := blockchain.TransactionReceipt{
		TransactionHash: "0xtx2", GasUsed: big.NewInt(21000),
		GasPrice: entities.NewWei(pegout.BridgeConversionGasPrice),
		Value: entities.NewWei(200),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(358)) == 0
	}), test.AnyAddress).Return(receipt1, nil).Once()
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
	var updatedQuotes []quote.RetainedPegoutQuote
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(1).(quote.RetainedPegoutQuote)
			updatedQuotes = append(updatedQuotes, q)
		}).Return(nil)
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	require.NoError(t, err)

	// q1(164): fully filled by chunk1, state=BridgeTxSucceeded, 1 allocation
	q1 := findUpdatedQuote(updatedQuotes, "02", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q1, "q1 should be updated to BridgeTxSucceeded")
	assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xtx1", q1.BridgeRefundTxHash)
	assert.Len(t, q1.BridgeRebalances, 1)

	// q2(134): fully filled by chunk1, state=BridgeTxSucceeded, 1 allocation
	q2 := findUpdatedQuote(updatedQuotes, "03", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q2, "q2 should be updated to BridgeTxSucceeded")
	assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xtx1", q2.BridgeRefundTxHash)
	assert.Len(t, q2.BridgeRebalances, 1)

	// q3(260): partial fill by chunk1 (60), completed by chunk2 (200)
	// last update should be BridgeTxSucceeded with 2 allocations
	q3 := findUpdatedQuote(updatedQuotes, "05", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q3, "q3 should be updated to BridgeTxSucceeded")
	assert.Equal(t, 0, q3.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xtx1", q3.BridgeRefundTxHash) // deprecated field from first allocation
	assert.Len(t, q3.BridgeRebalances, 2)
	assert.Equal(t, "0xtx1", q3.BridgeRebalances[0].TxHash)
	assert.Equal(t, "0xtx2", q3.BridgeRebalances[1].TxHash)

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
		Value:           entities.NewWei(350),
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
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "n1-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(200), CallFee: entities.NewWei(50), GasFee: entities.NewWei(100)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
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
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
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
		Value:           entities.NewWei(200),
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
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "em-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(100), GasFee: entities.NewWei(200)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)
	pegoutRepository.AssertExpectations(t)
	wallet.AssertExpectations(t)
}

// total=558, BridgeMin=200 => N=2, 2nd tx fails
// quotes: q1(164), q2(134), q3(260)
// chunk1 (value=358): fills q1, q2, q3 partial (60) — all persisted
// chunk2 (200): tx fails → skipped, q3 stays partially filled
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
		Value:           entities.NewWei(358),
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
	var updatedQuotes []quote.RetainedPegoutQuote
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(1).(quote.RetainedPegoutQuote)
			updatedQuotes = append(updatedQuotes, q)
		}).Return(nil)
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	// Failed chunks are skipped (logged), not fatal
	require.NoError(t, err)

	// q1 and q2 fully filled by chunk1
	q1 := findUpdatedQuote(updatedQuotes, "02", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q1)
	q2 := findUpdatedQuote(updatedQuotes, "03", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q2)

	// q3 partially filled (60 of 260), chunk2 failed so it stays partial
	q3Partial := findUpdatedQuoteByHash(updatedQuotes, "05")
	require.NotNil(t, q3Partial)
	assert.NotEqual(t, quote.PegoutStateBridgeTxSucceeded, q3Partial.State)
	assert.Equal(t, 0, q3Partial.RemainingToRefund.Cmp(entities.NewWei(200))) // 260 - 60 = 200
	assert.Len(t, q3Partial.BridgeRebalances, 1)

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
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	testQuotes := make([]quote.WatchedPegoutQuote, len(bridgePegoutTestWatchedQuotes))
	copy(testQuotes, bridgePegoutTestWatchedQuotes)
	err := useCase.Run(context.Background(), testQuotes[1], testQuotes[2], testQuotes[4])
	require.ErrorIs(t, err, usecases.InsufficientAmountError)
	wallet.AssertNotCalled(t, "SendRbtc")
}

// rbtc returns n * 10^17 (0.1 RBTC units for convenience)
func rbtc(n int64) *big.Int {
	v := big.NewInt(n)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil)
	return v.Mul(v, exp)
}

var amountIntegrityCases = []struct {
	name      string
	total     *big.Int
	bridgeMin *big.Int
	wantN     int
}{
	{"total equals bridgeMin exactly (N=1)", rbtc(15), rbtc(15), 1},
	{"total just above bridgeMin (N=1, remainder absorbed)", rbtc(29), rbtc(15), 1},
	{"total exactly 2x bridgeMin (N=2, no remainder)", rbtc(30), rbtc(15), 2},
	{"total just below 2x bridgeMin (N=1)", rbtc(29), rbtc(15), 1},
	{"total 2x+1wei above bridgeMin boundary",
		new(big.Int).Add(new(big.Int).Mul(rbtc(15), big.NewInt(2)), big.NewInt(1)), rbtc(15), 2},
	{"5 transactions with large remainder in first chunk", rbtc(83), rbtc(15), 5},
	{"3 transactions exact multiple", rbtc(45), rbtc(15), 3},
	{"large value: 100 RBTC split by 1.5 RBTC min", rbtc(1000), rbtc(15), 66},
}

// TestUtxoSplit_AmountIntegrity verifies that the sum of all transaction amounts
// sent to the bridge equals the original totalValue exactly, with no rounding errors.
func TestUtxoSplit_AmountIntegrity(t *testing.T) {
	for _, tc := range amountIntegrityCases {
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
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(totalWei.Copy(), entities.NewWei(int64(wantN)*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: test.AnyHash, GasUsed: big.NewInt(21000),
		GasPrice: entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:    totalWei.Copy(),
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
	pegoutLp := &mocks.ProviderMock{}
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: bridgeMinWei,
	}).Once()
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "integrity-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: totalWei, CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, walletMutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)
	assertAmountIntegrity(t, sentAmounts, total, bridgeMin, wantN)
	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

func assertAmountIntegrity(t *testing.T, sentAmounts []*big.Int, total, bridgeMin *big.Int, wantN int) {
	t.Helper()
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
}

// findUpdatedQuote returns the last update for the given quoteHash with the given state, or nil.
func findUpdatedQuote(quotes []quote.RetainedPegoutQuote, quoteHash string, state quote.PegoutState) *quote.RetainedPegoutQuote {
	var found *quote.RetainedPegoutQuote
	for i := len(quotes) - 1; i >= 0; i-- {
		if quotes[i].QuoteHash == quoteHash && quotes[i].State == state {
			found = &quotes[i]
			break
		}
	}
	return found
}

// findUpdatedQuoteByHash returns the last update for the given quoteHash regardless of state.
func findUpdatedQuoteByHash(quotes []quote.RetainedPegoutQuote, quoteHash string) *quote.RetainedPegoutQuote {
	var found *quote.RetainedPegoutQuote
	for i := len(quotes) - 1; i >= 0; i-- {
		if quotes[i].QuoteHash == quoteHash {
			found = &quotes[i]
			break
		}
	}
	return found
}

func TestUtxoSplit_Distribution(t *testing.T) {
	t.Run("one chunk spans two quotes", func(t *testing.T) {
		testUtxoSplitChunkSpansTwoQuotes(t)
	})
	t.Run("quote spans multiple chunks with partial fill", func(t *testing.T) {
		testUtxoSplitQuoteSpansMultipleChunks(t)
	})
	t.Run("DB update failure during allocation", func(t *testing.T) {
		testUtxoSplitDbUpdateFailure(t)
	})
	t.Run("retry with RemainingToRefund already set", func(t *testing.T) {
		testUtxoSplitRetryWithRemaining(t)
	})
	t.Run("all chunks fail", func(t *testing.T) {
		testUtxoSplitAllChunksFail(t)
	})
}

// Scenario B: one chunk spans two quotes
// Q1 need 300, Q2 need 400 => total=700, BridgeMin=500 => N=1, chunk=700
// chunk1 (value=700): fills Q1(300) fully, Q2(400) fully
func testUtxoSplitChunkSpansTwoQuotes(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: "0xspan",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(700),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(700)) == 0
	}), test.AnyAddress).Return(receipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(500),
	}).Once()
	var updatedQuotes []quote.RetainedPegoutQuote
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(1).(quote.RetainedPegoutQuote)
			updatedQuotes = append(updatedQuotes, q)
		}).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "span-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
		},
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "span-02", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(400), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)

	// Both quotes filled by same chunk, same tx hash
	q1 := findUpdatedQuote(updatedQuotes, "span-01", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q1)
	assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xspan", q1.BridgeRefundTxHash)
	assert.Len(t, q1.BridgeRebalances, 1)
	assert.Equal(t, "0xspan", q1.BridgeRebalances[0].TxHash)

	q2 := findUpdatedQuote(updatedQuotes, "span-02", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q2)
	assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xspan", q2.BridgeRefundTxHash)
	assert.Len(t, q2.BridgeRebalances, 1)

	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

// Scenario C: quote spans multiple chunks
// Q1 need 800 => total=800, BridgeMin=500 => N=1, chunk=800
// Actually to get multiple chunks: BridgeMin=300 => N=2, R=200 => chunk0=500, chunk1=300
// chunk0 (value=500): Q1 partial (500 of 800)
// chunk1 (value=300): Q1 remaining (300)
func testUtxoSplitQuoteSpansMultipleChunks(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(2*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt1 := blockchain.TransactionReceipt{
		TransactionHash: "0xmulti1",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(500),
	}
	receipt2 := blockchain.TransactionReceipt{
		TransactionHash: "0xmulti2",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(300),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(500)) == 0
	}), test.AnyAddress).Return(receipt1, nil).Once()
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(300)) == 0
	}), test.AnyAddress).Return(receipt2, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(300),
	}).Once()
	var updatedQuotes []quote.RetainedPegoutQuote
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(1).(quote.RetainedPegoutQuote)
			updatedQuotes = append(updatedQuotes, q)
		}).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "multi-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(800), CallFee: entities.NewWei(0), GasFee: entities.NewWei(0)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)

	// Should have 2 updates for multi-01: partial then complete
	require.Len(t, updatedQuotes, 2)

	// First update: partial fill (500 of 800)
	assert.Equal(t, "multi-01", updatedQuotes[0].QuoteHash)
	assert.NotEqual(t, quote.PegoutStateBridgeTxSucceeded, updatedQuotes[0].State)
	assert.Equal(t, 0, updatedQuotes[0].RemainingToRefund.Cmp(entities.NewWei(300)))
	assert.Equal(t, "0xmulti1", updatedQuotes[0].BridgeRefundTxHash) // deprecated field from first allocation
	assert.Len(t, updatedQuotes[0].BridgeRebalances, 1)

	// Second update: completed (remaining 300)
	assert.Equal(t, "multi-01", updatedQuotes[1].QuoteHash)
	assert.Equal(t, quote.PegoutStateBridgeTxSucceeded, updatedQuotes[1].State)
	assert.Equal(t, 0, updatedQuotes[1].RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xmulti1", updatedQuotes[1].BridgeRefundTxHash) // stays from first allocation
	assert.Len(t, updatedQuotes[1].BridgeRebalances, 2)
	assert.Equal(t, "0xmulti1", updatedQuotes[1].BridgeRebalances[0].TxHash)
	assert.Equal(t, "0xmulti2", updatedQuotes[1].BridgeRebalances[1].TxHash)

	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

// Scenario E: DB update failure during allocation returns error immediately
func testUtxoSplitDbUpdateFailure(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: "0xdbfail",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(500),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(500)) == 0
	}), test.AnyAddress).Return(receipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(400),
	}).Once()
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Return(errors.New("db connection lost")).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "dbfail-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(100), GasFee: entities.NewWei(100)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.ErrorContains(t, err, "db connection lost")
	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

// Retry scenario: quotes have RemainingToRefund from a previous partial run
// Q1 has remaining=200 (originally 500), Q2 is fresh with contribution=300
// adjustedTotal = 200 + 300 = 500, BridgeMin=400 => N=1, chunk=500
func testUtxoSplitRetryWithRemaining(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	receipt := blockchain.TransactionReceipt{
		TransactionHash: "0xretry",
		GasUsed:         big.NewInt(21000),
		GasPrice:        entities.NewWei(pegout.BridgeConversionGasPrice),
		Value:           entities.NewWei(500),
	}
	wallet.On("SendRbtc", mock.Anything, mock.MatchedBy(func(config blockchain.TransactionConfig) bool {
		return config.Value.Cmp(entities.NewWei(500)) == 0
	}), test.AnyAddress).Return(receipt, nil).Once()
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(400),
	}).Once()
	var updatedQuotes []quote.RetainedPegoutQuote
	pegoutRepository.On("UpdateRetainedQuote", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			q := args.Get(1).(quote.RetainedPegoutQuote)
			updatedQuotes = append(updatedQuotes, q)
		}).Return(nil)
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:         "retry-01",
				State:             quote.PegoutStateRefundPegOutSucceeded,
				RemainingToRefund: entities.NewWei(200), // partial from previous run
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
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	require.NoError(t, err)

	// Q1: had remaining=200, now fully filled. Should keep previous allocation + new one
	q1 := findUpdatedQuote(updatedQuotes, "retry-01", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q1)
	assert.Equal(t, 0, q1.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xprev", q1.BridgeRefundTxHash) // deprecated field preserved from first allocation
	assert.Len(t, q1.BridgeRebalances, 2)
	assert.Equal(t, "0xprev", q1.BridgeRebalances[0].TxHash)
	assert.Equal(t, "0xretry", q1.BridgeRebalances[1].TxHash)

	// Q2: fresh quote, fully filled
	q2 := findUpdatedQuote(updatedQuotes, "retry-02", quote.PegoutStateBridgeTxSucceeded)
	require.NotNil(t, q2)
	assert.Equal(t, 0, q2.RemainingToRefund.Cmp(entities.NewWei(0)))
	assert.Equal(t, "0xretry", q2.BridgeRefundTxHash)
	assert.Len(t, q2.BridgeRebalances, 1)

	wallet.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
}

// All chunks fail: no quotes are updated, no error returned (failures are logged and skipped)
func testUtxoSplitAllChunksFail(t *testing.T) {
	pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutLp := &mocks.ProviderMock{}
	wallet := &mocks.RskWalletMock{}
	gasPerTx := int64(pegout.BridgeConversionGasLimit * pegout.BridgeConversionGasPrice)
	walletBalance := new(entities.Wei).Add(entities.NewWei(1000), entities.NewWei(2*gasPerTx))
	wallet.On("GetBalance", mock.Anything).Return(walletBalance, nil).Once()
	emptyReceipt := blockchain.TransactionReceipt{}
	wallet.On("SendRbtc", mock.Anything, mock.Anything, test.AnyAddress).
		Return(emptyReceipt, assert.AnError)
	mutex := &mocks.MutexMock{}
	mutex.On("Lock").Return().Once()
	mutex.On("Unlock").Return().Once()
	bridge := &mocks.BridgeMock{}
	bridge.On("GetAddress").Return(test.AnyAddress)
	pegoutLp.On("PegoutConfiguration", mock.Anything).Return(liquidity_provider.PegoutConfiguration{
		BridgeTransactionMin: entities.NewWei(200),
	}).Once()
	customQuotes := []quote.WatchedPegoutQuote{
		{
			RetainedQuote: quote.RetainedPegoutQuote{QuoteHash: "allfail-01", State: quote.PegoutStateRefundPegOutSucceeded},
			PegoutQuote:   quote.PegoutQuote{Value: entities.NewWei(300), CallFee: entities.NewWei(50), GasFee: entities.NewWei(50)},
		},
	}
	handler := pegout.NewUtxoSplitHandler(pegoutRepository, wallet, blockchain.RskContracts{Bridge: bridge})
	useCase := pegout.NewBridgePegoutUseCase(pegoutLp, mutex, handler)
	err := useCase.Run(context.Background(), customQuotes...)
	// All chunks failed but function returns nil (errors are logged)
	require.NoError(t, err)
	// No DB updates should have been made
	pegoutRepository.AssertNotCalled(t, "UpdateRetainedQuote")
	wallet.AssertExpectations(t)
}
