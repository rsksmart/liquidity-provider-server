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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
	useCase := pegout.NewBridgePegoutUseCase(pegoutRepository, pegoutLp, wallet, blockchain.RskContracts{Bridge: bridge}, mutex)
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
