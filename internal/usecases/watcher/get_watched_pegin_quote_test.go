package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var retainedQuotes = []quote.RetainedPeginQuote{
	{QuoteHash: "01", State: quote.PeginStateWaitingForDeposit},
	{QuoteHash: "02", State: quote.PeginStateCallForUserSucceeded},
	{QuoteHash: "04", State: quote.PeginStateCallForUserSucceeded},
	{QuoteHash: "03", State: quote.PeginStateWaitingForDeposit},
	{QuoteHash: "05", State: quote.PeginStateCallForUserFailed},
	{QuoteHash: "07", State: quote.PeginStateWaitingForDepositConfirmations},
}

var peginQuotes = []quote.PeginQuote{
	{Nonce: quote.NewNonce(1)},
	{Nonce: quote.NewNonce(2)},
	{Nonce: quote.NewNonce(3)},
	{Nonce: quote.NewNonce(4)},
	{Nonce: quote.NewNonce(5)},
	{Nonce: quote.NewNonce(7)},
	{Nonce: quote.NewNonce(8)},
}

var peginsCreationData = []quote.PeginCreationData{
	{GasPrice: entities.NewWei(1)},
	{GasPrice: entities.NewWei(2)},
	{GasPrice: entities.NewWei(3)},
	{GasPrice: entities.NewWei(4)},
	{GasPrice: entities.NewWei(5)},
	{GasPrice: entities.NewWei(6)},
	{GasPrice: entities.NewWei(7)},
}

func TestGetWatchedPeginQuoteUseCase_Run_WaitingForDeposit(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDeposit).
		Return([]quote.RetainedPeginQuote{retainedQuotes[0], retainedQuotes[3]}, nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[0].QuoteHash).Return(&peginQuotes[0], nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[3].QuoteHash).Return(&peginQuotes[2], nil)
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[0].QuoteHash).Return(peginsCreationData[0])
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[3].QuoteHash).Return(peginsCreationData[3])
	useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PeginStateWaitingForDeposit)
	quoteRepository.AssertExpectations(t)
	assert.Len(t, watchedQuotes, 2)
	require.NoError(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		// this is just to validate that the watched quotes are built with the correct pairs,
		// the nonce is not related to the hash in the business logic
		assert.Equal(t, quote.NewNonce(parsedHash.Int64()), watchedQuote.PeginQuote.Nonce)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_WaitingForDepositConfirmations(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPeginQuote{retainedQuotes[1], retainedQuotes[5]}, nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[1].QuoteHash).Return(&peginQuotes[1], nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[5].QuoteHash).Return(&peginQuotes[5], nil)
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[1].QuoteHash).Return(peginsCreationData[1])
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[5].QuoteHash).Return(peginsCreationData[5])
	useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PeginStateWaitingForDepositConfirmations)
	quoteRepository.AssertExpectations(t)
	assert.Len(t, watchedQuotes, 2)
	require.NoError(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		// this is just to validate that the watched quotes are built with the correct pairs,
		// the nonce is not related to the hash in the business logic
		assert.Equal(t, quote.NewNonce(parsedHash.Int64()), watchedQuote.PeginQuote.Nonce)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_MoreThanOneState(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDeposit).
		Return([]quote.RetainedPeginQuote{retainedQuotes[0], retainedQuotes[3]}, nil)
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPeginQuote{retainedQuotes[1], retainedQuotes[5]}, nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[1].QuoteHash).Return(&peginQuotes[1], nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[5].QuoteHash).Return(&peginQuotes[5], nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[0].QuoteHash).Return(&peginQuotes[0], nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[3].QuoteHash).Return(&peginQuotes[2], nil)
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[1].QuoteHash).Return(peginsCreationData[1])
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[5].QuoteHash).Return(peginsCreationData[5])
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[0].QuoteHash).Return(peginsCreationData[0])
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[3].QuoteHash).Return(peginsCreationData[3])
	useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)
	quoteRepository.AssertExpectations(t)
	assert.Len(t, watchedQuotes, 4)
	require.NoError(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		assert.Equal(t, quote.NewNonce(parsedHash.Int64()), watchedQuote.PeginQuote.Nonce)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_CallForUserSucceed(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, quote.PeginStateCallForUserSucceeded).
		Return([]quote.RetainedPeginQuote{retainedQuotes[2]}, nil)
	quoteRepository.On("GetQuote", test.AnyCtx, retainedQuotes[2].QuoteHash).Return(&peginQuotes[3], nil)
	quoteRepository.EXPECT().GetPeginCreationData(test.AnyCtx, retainedQuotes[2].QuoteHash).Return(peginsCreationData[2])
	useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PeginStateCallForUserSucceeded)
	quoteRepository.AssertExpectations(t)
	assert.Len(t, watchedQuotes, 1)
	require.NoError(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		// this is just to validate that the watched quotes are built with the correct pairs,
		// the nonce is not related to the hash in the business logic
		assert.Equal(t, quote.NewNonce(parsedHash.Int64()), watchedQuote.PeginQuote.Nonce)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_WrongState(t *testing.T) {
	wrongStates := []quote.PeginState{
		quote.PeginStateTimeForDepositElapsed,
		quote.PeginStateCallForUserFailed,
		quote.PeginStateRegisterPegInFailed,
		quote.PeginStateRegisterPegInSucceeded,
	}
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
	for _, state := range wrongStates {
		_, err := useCase.Run(context.Background(), state)
		require.Error(t, err)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	setups := []func(quoteRepository *mocks.PeginQuoteRepositoryMock){
		func(quoteRepository *mocks.PeginQuoteRepositoryMock) {
			quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, mock.Anything).
				Return(nil, assert.AnError)
		},
		func(quoteRepository *mocks.PeginQuoteRepositoryMock) {
			quoteRepository.On("GetRetainedQuoteByState", test.AnyCtx, mock.Anything).
				Return(retainedQuotes, nil)
			quoteRepository.On("GetQuote", test.AnyCtx, mock.Anything).
				Return(nil, assert.AnError)
		},
	}
	for _, setup := range setups {
		quoteRepository := new(mocks.PeginQuoteRepositoryMock)
		setup(quoteRepository)
		useCase := watcher.NewGetWatchedPeginQuoteUseCase(quoteRepository)
		_, err := useCase.Run(context.Background(), quote.PeginStateWaitingForDeposit)
		require.Error(t, err)
	}
}
