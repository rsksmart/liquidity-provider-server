package watcher_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCleanExpiredQuotesUseCase_Run(t *testing.T) {
	peginExpiredQuotes := []quote.RetainedPeginQuote{
		{QuoteHash: "peginHash1", State: quote.PeginStateTimeForDepositElapsed},
		{QuoteHash: "peginHash2", State: quote.PeginStateTimeForDepositElapsed},
		{QuoteHash: "peginHash3", State: quote.PeginStateCallForUserSucceeded},
		{QuoteHash: "peginHash4", State: quote.PeginStateTimeForDepositElapsed},
		{QuoteHash: "peginHash5", State: quote.PeginStateRegisterPegInFailed},
		{QuoteHash: "peginHash6", State: quote.PeginStateRegisterPegInSucceeded},
	}

	pegoutExpiredQuotes := []quote.RetainedPegoutQuote{
		{QuoteHash: "pegoutHash1", State: quote.PegoutStateRefundPegOutSucceeded},
		{QuoteHash: "pegoutHash2", State: quote.PegoutStateSendPegoutSucceeded},
		{QuoteHash: "pegoutHash3", State: quote.PegoutStateTimeForDepositElapsed},
		{QuoteHash: "pegoutHash4", State: quote.PegoutStateSendPegoutFailed},
		{QuoteHash: "pegoutHash5", State: quote.PegoutStateTimeForDepositElapsed},
		{QuoteHash: "pegoutHash6", State: quote.PegoutStateTimeForDepositElapsed},
	}

	peginRepository := new(test.PeginQuoteRepositoryMock)
	peginRepository.On(
		"GetRetainedQuoteByState",
		mock.AnythingOfType("context.backgroundCtx"),
		[]quote.PeginState{quote.PeginStateTimeForDepositElapsed},
	).Return([]quote.RetainedPeginQuote{peginExpiredQuotes[0], peginExpiredQuotes[1], peginExpiredQuotes[3]}, nil)
	peginRepository.On(
		"DeleteQuotes",
		mock.AnythingOfType("context.backgroundCtx"),
		[]string{"peginHash1", "peginHash2", "peginHash4"},
	).Return(uint(3), nil)

	pegoutRepository := new(test.PegoutQuoteRepositoryMock)
	pegoutRepository.On(
		"GetRetainedQuoteByState",
		mock.AnythingOfType("context.backgroundCtx"),
		[]quote.PegoutState{quote.PegoutStateTimeForDepositElapsed},
	).Return([]quote.RetainedPegoutQuote{pegoutExpiredQuotes[2], pegoutExpiredQuotes[4], pegoutExpiredQuotes[5]}, nil)
	pegoutRepository.On(
		"DeleteQuotes",
		mock.AnythingOfType("context.backgroundCtx"),
		[]string{"pegoutHash3", "pegoutHash5", "pegoutHash6"},
	).Return(uint(3), nil)

	useCase := watcher.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
	hashes, err := useCase.Run(context.Background())

	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	assert.Nil(t, err)
	assert.True(t, len(hashes) == 6)
	assert.Equal(t, []string{"peginHash1", "peginHash2", "peginHash4", "pegoutHash3", "pegoutHash5", "pegoutHash6"}, hashes)
}

func TestCleanExpiredQuotesUseCase_Run_ErrorHandling(t *testing.T) {
	setups := []func(peginRepository *test.PeginQuoteRepositoryMock, pegoutRepository *test.PegoutQuoteRepositoryMock){
		func(peginRepository *test.PeginQuoteRepositoryMock, pegoutRepository *test.PegoutQuoteRepositoryMock) {
			peginRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(nil, assert.AnError)
		},
		func(peginRepository *test.PeginQuoteRepositoryMock, pegoutRepository *test.PegoutQuoteRepositoryMock) {
			peginRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedQuotes, nil)
			pegoutRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(nil, assert.AnError)
		},
		func(peginRepository *test.PeginQuoteRepositoryMock, pegoutRepository *test.PegoutQuoteRepositoryMock) {
			peginRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedQuotes, nil)
			pegoutRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedPegoutQuotes, nil)
			peginRepository.On("DeleteQuotes", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(uint(0), assert.AnError)
		},
		func(peginRepository *test.PeginQuoteRepositoryMock, pegoutRepository *test.PegoutQuoteRepositoryMock) {
			peginRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedQuotes, nil)
			pegoutRepository.On("GetRetainedQuoteByState", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(retainedPegoutQuotes, nil)
			peginRepository.On("DeleteQuotes", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(uint(5), nil)
			pegoutRepository.On("DeleteQuotes", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).
				Return(uint(0), assert.AnError)
		},
	}

	for _, setup := range setups {
		peginRepository := new(test.PeginQuoteRepositoryMock)
		pegoutRepository := new(test.PegoutQuoteRepositoryMock)
		setup(peginRepository, pegoutRepository)
		useCase := watcher.NewCleanExpiredQuotesUseCase(peginRepository, pegoutRepository)
		_, err := useCase.Run(context.Background())
		assert.NotNil(t, err)
		peginRepository.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
	}
}
