package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpiredPegoutQuoteUseCase_Run(t *testing.T) {
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         "0x1234",
		DepositAddress:    "0xa1b2c3",
		Signature:         "0x4321",
		RequiredLiquidity: entities.NewWei(1),
		State:             quote.PegoutStateWaitingForDeposit,
	}

	expectedRetainedQuote := retainedQuote
	expectedRetainedQuote.State = quote.PegoutStateTimeForDepositElapsed
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), expectedRetainedQuote).Return(nil)
	useCase := pegout.NewExpiredPegoutQuoteUseCase(pegoutQuoteRepository)
	err := useCase.Run(context.Background(), retainedQuote)
	pegoutQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
}

func TestExpiredPegoutQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	retainedQuote := quote.RetainedPegoutQuote{
		QuoteHash:         "0x1234",
		DepositAddress:    "0xa1b2c3",
		Signature:         "0x4321",
		RequiredLiquidity: entities.NewWei(1),
		State:             quote.PegoutStateWaitingForDeposit,
	}
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.On("UpdateRetainedQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything).Return(assert.AnError)
	useCase := pegout.NewExpiredPegoutQuoteUseCase(pegoutQuoteRepository)
	err := useCase.Run(context.Background(), retainedQuote)
	pegoutQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
}
