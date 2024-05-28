package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestUpdatePeginDepositUseCase_Run(t *testing.T) {
	peginQuote := quote.PeginQuote{
		Value:              entities.NewWei(5000),
		CallFee:            entities.NewWei(100),
		ProductFeeAmount:   1,
		GasFee:             entities.NewWei(150),
		AgreementTimestamp: 900,
		TimeForDeposit:     200,
	}
	retainedQuote := quote.RetainedPeginQuote{
		State:          quote.PeginStateWaitingForDeposit,
		DepositAddress: test.AnyAddress,
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	block := blockchain.BitcoinBlockInformation{
		Hash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Height: big.NewInt(500),
		Time:   time.Unix(1000, 0),
	}
	tx := blockchain.BitcoinTransactionInformation{
		Hash:          test.AnyString,
		Confirmations: 5,
		Outputs: map[string][]*entities.Wei{
			test.AnyAddress: {entities.NewWei(5251)},
		},
	}
	quoteRepository.On(
		"UpdateRetainedQuote",
		test.AnyCtx,
		mock.MatchedBy(func(q quote.RetainedPeginQuote) bool {
			return q.UserBtcTxHash == test.AnyString && q.State == quote.PeginStateWaitingForDepositConfirmations
		}),
	).Return(nil)
	useCase := watcher.NewUpdatePeginDepositUseCase(quoteRepository)
	watchedPeginQuote, err := useCase.Run(context.Background(), quote.NewWatchedPeginQuote(peginQuote, retainedQuote), block, tx)
	quoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, quote.PeginStateWaitingForDepositConfirmations, watchedPeginQuote.RetainedQuote.State)
	assert.Equal(t, tx.Hash, watchedPeginQuote.RetainedQuote.UserBtcTxHash)
}

func TestUpdatePeginDepositUseCase_Run_ErrorHandling(t *testing.T) {
	const bitcoinTxErrorMsg = "invalid bitcoin transaction for quote"
	peginQuote := quote.PeginQuote{
		Value:              entities.NewWei(1),
		CallFee:            entities.NewWei(2),
		ProductFeeAmount:   3,
		GasFee:             entities.NewWei(4),
		AgreementTimestamp: 500,
		TimeForDeposit:     50,
	}
	t.Run("Fail by illegal quote state", func(t *testing.T) {
		states := []quote.PeginState{
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		block := blockchain.BitcoinBlockInformation{Time: time.Unix(510, 0)}
		tx := blockchain.BitcoinTransactionInformation{Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(10)}}}
		for _, state := range states {
			retainedQuote := quote.RetainedPeginQuote{State: state, DepositAddress: test.AnyAddress}
			useCase := watcher.NewUpdatePeginDepositUseCase(quoteRepository)
			watchedPeginQuote, err := useCase.Run(context.Background(), quote.NewWatchedPeginQuote(peginQuote, retainedQuote), block, tx)
			require.ErrorIs(t, err, usecases.IllegalQuoteStateError)
			assert.Empty(t, watchedPeginQuote)
		}
	})
	t.Run("Fail by bitcoin transaction amount", func(t *testing.T) {
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		block := blockchain.BitcoinBlockInformation{Time: time.Unix(510, 0)}
		tx := blockchain.BitcoinTransactionInformation{Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(5)}}}
		retainedQuote := quote.RetainedPeginQuote{State: quote.PeginStateWaitingForDeposit, DepositAddress: test.AnyAddress}
		useCase := watcher.NewUpdatePeginDepositUseCase(quoteRepository)
		watchedPeginQuote, err := useCase.Run(context.Background(), quote.NewWatchedPeginQuote(peginQuote, retainedQuote), block, tx)
		require.ErrorContains(t, err, bitcoinTxErrorMsg)
		assert.Empty(t, watchedPeginQuote)
	})
	t.Run("Fail by bitcoin transaction timestamp", func(t *testing.T) {
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		block := blockchain.BitcoinBlockInformation{Time: time.Unix(2000, 0)}
		tx := blockchain.BitcoinTransactionInformation{Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(10)}}}
		retainedQuote := quote.RetainedPeginQuote{State: quote.PeginStateWaitingForDeposit, DepositAddress: test.AnyAddress}
		useCase := watcher.NewUpdatePeginDepositUseCase(quoteRepository)
		watchedPeginQuote, err := useCase.Run(context.Background(), quote.NewWatchedPeginQuote(peginQuote, retainedQuote), block, tx)
		require.ErrorContains(t, err, bitcoinTxErrorMsg)
		assert.Empty(t, watchedPeginQuote)
	})
	t.Run("Fail to update retained quote", func(t *testing.T) {
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		quoteRepository.On("UpdateRetainedQuote", test.AnyCtx, mock.Anything).Return(assert.AnError)
		block := blockchain.BitcoinBlockInformation{Time: time.Unix(510, 0)}
		tx := blockchain.BitcoinTransactionInformation{Outputs: map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(10)}}}
		retainedQuote := quote.RetainedPeginQuote{State: quote.PeginStateWaitingForDeposit, DepositAddress: test.AnyAddress}
		useCase := watcher.NewUpdatePeginDepositUseCase(quoteRepository)
		watchedPeginQuote, err := useCase.Run(context.Background(), quote.NewWatchedPeginQuote(peginQuote, retainedQuote), block, tx)
		require.Error(t, err)
		assert.Empty(t, watchedPeginQuote)
	})
}
