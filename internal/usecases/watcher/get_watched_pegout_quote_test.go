package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/big"
	"testing"
)

var okPegoutStates = []quote.PegoutState{quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations, quote.PegoutStateSendPegoutSucceeded}

var retainedPegoutQuotes = []quote.RetainedPegoutQuote{
	{QuoteHash: "01", State: quote.PegoutStateWaitingForDeposit},
	{QuoteHash: "04", State: quote.PegoutStateSendPegoutSucceeded},
	{QuoteHash: "03", State: quote.PegoutStateWaitingForDeposit},
	{QuoteHash: "05", State: quote.PegoutStateSendPegoutFailed},
	{QuoteHash: "06", State: quote.PegoutStateWaitingForDepositConfirmations},
}

var pegoutQuotes = []quote.PegoutQuote{
	{Nonce: 1},
	{Nonce: 2},
	{Nonce: 3},
	{Nonce: 4},
	{Nonce: 5},
	{Nonce: 6},
}

func TestGetWatchedPegoutQuoteUseCase_Run_WaitingForDeposit(t *testing.T) {
	quoteRepository := new(test.PegoutQuoteRepositoryMock)
	quoteRepository.On(
		"GetRetainedQuoteByState",
		mock.AnythingOfType("context.backgroundCtx"),
		[]quote.PegoutState{quote.PegoutStateWaitingForDeposit},
	).Return([]quote.RetainedPegoutQuote{retainedPegoutQuotes[0], retainedPegoutQuotes[2]}, nil)
	quoteRepository.On(
		"GetRetainedQuoteByState",
		mock.AnythingOfType("context.backgroundCtx"),
		[]quote.PegoutState{quote.PegoutStateWaitingForDepositConfirmations},
	).Return([]quote.RetainedPegoutQuote{retainedPegoutQuotes[4]}, nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPegoutQuotes[0].QuoteHash).Return(&pegoutQuotes[0], nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPegoutQuotes[2].QuoteHash).Return(&pegoutQuotes[2], nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPegoutQuotes[4].QuoteHash).Return(&pegoutQuotes[5], nil)
	useCase := watcher.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations)
	quoteRepository.AssertExpectations(t)
	assert.True(t, len(watchedQuotes) == 3)
	assert.Nil(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		// this is just to validate that the watched quotes are built with the correct pairs,
		// the nonce is not related to the hash in the business logic
		assert.True(t, parsedHash.Int64() == watchedQuote.PegoutQuote.Nonce)
	}
}

func TestGetWatchedPegoutQuoteUseCase_Run_CallForUserSucceed(t *testing.T) {
	quoteRepository := new(test.PegoutQuoteRepositoryMock)
	quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), []quote.PegoutState{quote.PegoutStateSendPegoutSucceeded}).
		Return([]quote.RetainedPegoutQuote{retainedPegoutQuotes[1]}, nil)
	quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), retainedPegoutQuotes[1].QuoteHash).Return(&pegoutQuotes[3], nil)
	useCase := watcher.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
	watchedQuotes, err := useCase.Run(context.Background(), quote.PegoutStateSendPegoutSucceeded)
	quoteRepository.AssertExpectations(t)
	assert.True(t, len(watchedQuotes) == 1)
	assert.Nil(t, err)
	var parsedHash big.Int
	for _, watchedQuote := range watchedQuotes {
		parsedHash.SetString(watchedQuote.RetainedQuote.QuoteHash, 16)
		// this is just to validate that the watched quotes are built with the correct pairs,
		// the nonce is not related to the hash in the business logic
		assert.True(t, parsedHash.Int64() == watchedQuote.PegoutQuote.Nonce)
	}
}

func TestGetWatchedPegoutQuoteUseCase_Run_WrongState(t *testing.T) {
	wrongStates := []quote.PegoutState{
		quote.PegoutStateTimeForDepositElapsed,
		quote.PegoutStateSendPegoutFailed,
		quote.PegoutStateRefundPegOutFailed,
		quote.PegoutStateRefundPegOutSucceeded,
	}
	quoteRepository := new(test.PegoutQuoteRepositoryMock)
	useCase := watcher.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
	for _, state := range wrongStates {
		_, err := useCase.Run(context.Background(), state)
		assert.NotNil(t, err)
	}
}

func TestGetWatchedPegoutQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	setups := []func(quoteRepository *test.PegoutQuoteRepositoryMock){
		func(quoteRepository *test.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(nil, assert.AnError)
		},
		func(quoteRepository *test.PegoutQuoteRepositoryMock) {
			quoteRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedPegoutQuotes, nil)
			quoteRepository.On("GetQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(nil, assert.AnError)
		},
	}
	for _, setup := range setups {
		quoteRepository := new(test.PegoutQuoteRepositoryMock)
		setup(quoteRepository)
		useCase := watcher.NewGetWatchedPegoutQuoteUseCase(quoteRepository)
		_, err := useCase.Run(context.Background(), quote.PegoutStateWaitingForDeposit)
		assert.NotNil(t, err)
	}
}
