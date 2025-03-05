package pegin_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPeginReportUseCase_Run(t *testing.T) {
	ctx := context.Background()

	retainedQuotes := []quote.RetainedPeginQuote{
		{QuoteHash: "hash1"},
		{QuoteHash: "hash2"},
	}

	peginQuotes := []quote.PeginQuote{
		{
			Value:   entities.NewWei(1000),
			CallFee: entities.NewWei(10),
		},
		{
			Value:   entities.NewWei(3000),
			CallFee: entities.NewWei(20),
		},
	}

	expectedMinimumValue := entities.NewWei(1000)
	expectedMaximumValue := entities.NewWei(3000)
	expectedAverageValue := entities.NewWei((1000 + 3000) / 2)
	expectedTotalFees := entities.NewWei(10 + 20)
	expectedAverageFee := entities.NewWei((10 + 20) / 2)

	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
	peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
		Return(retainedQuotes, nil).Once()
	peginQuoteRepository.On("GetQuote", ctx, "hash1").Return(&peginQuotes[0], nil).Once()
	peginQuoteRepository.On("GetQuote", ctx, "hash2").Return(&peginQuotes[1], nil).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx)

	peginQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, 2, result.NumberOfQuotes)
	assert.Equal(t, expectedMinimumValue, result.MinimumQuoteValue)
	assert.Equal(t, expectedMaximumValue, result.MaximumQuoteValue)
	assert.Equal(t, expectedAverageValue, result.AverageQuoteValue)
	assert.Equal(t, expectedTotalFees, result.TotalFeesCollected)
	assert.Equal(t, expectedAverageFee, result.AverageFeePerQuote)
}

func TestGetPeginReportUseCase_Run_EmptyQuotes(t *testing.T) {
	ctx := context.Background()

	retainedQuotes := []quote.RetainedPeginQuote{}

	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
	peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
		Return(retainedQuotes, nil).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx)

	peginQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, 0, result.NumberOfQuotes)
	assert.Equal(t, entities.NewWei(0), result.MinimumQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.MaximumQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.AverageQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.TotalFeesCollected)
	assert.Equal(t, entities.NewWei(0), result.AverageFeePerQuote)
}

func TestGetPeginReportUseCase_Run_ErrorFetchingQuotes(t *testing.T) {
	ctx := context.Background()

	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
	peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
		Return(nil, assert.AnError).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx)

	peginQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
	assert.Zero(t, result.NumberOfQuotes)
}
