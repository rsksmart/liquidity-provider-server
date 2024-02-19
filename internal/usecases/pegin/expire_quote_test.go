package pegin_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpiredPeginQuoteUseCase_Run(t *testing.T) {
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         "0x1234",
		DepositAddress:    "0xa1b2c3",
		Signature:         "0x4321",
		RequiredLiquidity: entities.NewWei(1),
		State:             quote.PeginStateWaitingForDeposit,
	}

	expectedRetainedQuote := retainedQuote
	expectedRetainedQuote.State = quote.PeginStateTimeForDepositElapsed
	peginQuoteRepository := new(test.PeginQuoteRepositoryMock)
	peginQuoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), expectedRetainedQuote).Return(nil)
	useCase := pegin.NewExpiredPeginQuoteUseCase(peginQuoteRepository)
	err := useCase.Run(context.Background(), retainedQuote)
	peginQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
}

func TestExpiredPeginQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	retainedQuote := quote.RetainedPeginQuote{
		QuoteHash:         "0x1234",
		DepositAddress:    "0xa1b2c3",
		Signature:         "0x4321",
		RequiredLiquidity: entities.NewWei(1),
		State:             quote.PeginStateWaitingForDeposit,
	}
	peginQuoteRepository := new(test.PeginQuoteRepositoryMock)
	peginQuoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError)
	useCase := pegin.NewExpiredPeginQuoteUseCase(peginQuoteRepository)
	err := useCase.Run(context.Background(), retainedQuote)
	peginQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
}
