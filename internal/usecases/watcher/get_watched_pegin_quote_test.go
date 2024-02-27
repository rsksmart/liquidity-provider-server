package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var retainedQuotes = []quote.RetainedPeginQuote{
	{QuoteHash: "01", State: quote.PeginStateWaitingForDeposit},
	{QuoteHash: "04", State: quote.PeginStateCallForUserSucceeded},
	{QuoteHash: "03", State: quote.PeginStateWaitingForDeposit},
	{QuoteHash: "05", State: quote.PeginStateCallForUserFailed},
}

var peginQuotes = []quote.PeginQuote{
	{Nonce: 1},
	{Nonce: 2},
	{Nonce: 3},
	{Nonce: 4},
	{Nonce: 5},
}

func TestGetWatchedPeginQuoteUseCase_Run_WaitingForDeposit(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), []quote.PeginState{quote.PeginStateWaitingForDeposit}).
		Return([]quote.RetainedPeginQuote{retainedQuotes[0], retainedQuotes[2]}, nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuotes[0].QuoteHash).Return(&peginQuotes[0], nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuotes[2].QuoteHash).Return(&peginQuotes[2], nil)
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
		assert.Equal(t, parsedHash.Int64(), watchedQuote.PeginQuote.Nonce)
	}
}

func TestGetWatchedPeginQuoteUseCase_Run_CallForUserSucceed(t *testing.T) {
	quoteRepository := new(mocks.PeginQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), []quote.PeginState{quote.PeginStateCallForUserSucceeded}).
		Return([]quote.RetainedPeginQuote{retainedQuotes[1]}, nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedQuotes[1].QuoteHash).Return(&peginQuotes[3], nil)
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
		assert.Equal(t, parsedHash.Int64(), watchedQuote.PeginQuote.Nonce)
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
			quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(nil, assert.AnError)
		},
		func(quoteRepository *mocks.PeginQuoteRepositoryMock) {
			quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedQuotes, nil)
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
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
