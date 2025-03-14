package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPegoutReportUseCase_Run(t *testing.T) {
	ctx := context.Background()

	retainedQuotes := []quote.RetainedPegoutQuote{
		{QuoteHash: "hash1"},
		{QuoteHash: "hash2"},
		{QuoteHash: "hash3"},
		{QuoteHash: "hash4"},
		{QuoteHash: "hash5"},
		{QuoteHash: "hash6"},
		{QuoteHash: "hash7"},
		{QuoteHash: "hash8"},
		{QuoteHash: "hash9"},
		{QuoteHash: "hash10"},
	}

	quoteHashes := []string{"hash1", "hash2", "hash3", "hash4", "hash5", "hash6", "hash7", "hash8", "hash9", "hash10"}

	pegoutQuotes := []quote.PegoutQuote{
		{Value: entities.NewWei(1000), CallFee: entities.NewWei(10)},
		{Value: entities.NewWei(2000), CallFee: entities.NewWei(20)},
		{Value: entities.NewWei(3000), CallFee: entities.NewWei(30)},
		{Value: entities.NewWei(4000), CallFee: entities.NewWei(40)},
		{Value: entities.NewWei(5000), CallFee: entities.NewWei(50)},
		{Value: entities.NewWei(6000), CallFee: entities.NewWei(60)},
		{Value: entities.NewWei(7000), CallFee: entities.NewWei(70)},
		{Value: entities.NewWei(8000), CallFee: entities.NewWei(80)},
		{Value: entities.NewWei(9000), CallFee: entities.NewWei(90)},
		{Value: entities.NewWei(10000), CallFee: entities.NewWei(100)},
	}

	// Calculate expected values
	expectedMinimumValue := entities.NewWei(1000)
	expectedMaximumValue := entities.NewWei(10000)
	expectedAverageValue := entities.NewWei(5500)
	expectedTotalFees := entities.NewWei(550)
	expectedAverageFee := entities.NewWei(55)

	pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
		Return(retainedQuotes, nil).Once()

	pegoutQuoteRepository.On("GetQuotes", ctx, quoteHashes).Return(pegoutQuotes, nil).Once()

	useCase := pegout.NewGetPegoutReportUseCase(pegoutQuoteRepository)

	result, err := useCase.Run(ctx)

	pegoutQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, 10, result.NumberOfQuotes)
	assert.Equal(t, expectedMinimumValue, result.MinimumQuoteValue)
	assert.Equal(t, expectedMaximumValue, result.MaximumQuoteValue)
	assert.Equal(t, expectedAverageValue, result.AverageQuoteValue)
	assert.Equal(t, expectedTotalFees, result.TotalFeesCollected)
	assert.Equal(t, expectedAverageFee, result.AverageFeePerQuote)
}

func TestGetPegoutReportUseCase_Run_EmptyQuotes(t *testing.T) {
	ctx := context.Background()

	retainedQuotes := []quote.RetainedPegoutQuote{}

	pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
		Return(retainedQuotes, nil).Once()

	useCase := pegout.NewGetPegoutReportUseCase(pegoutQuoteRepository)

	result, err := useCase.Run(ctx)

	pegoutQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, 0, result.NumberOfQuotes)
	assert.Equal(t, entities.NewWei(0), result.MinimumQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.MaximumQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.AverageQuoteValue)
	assert.Equal(t, entities.NewWei(0), result.TotalFeesCollected)
	assert.Equal(t, entities.NewWei(0), result.AverageFeePerQuote)
}

func TestGetPegoutReportUseCase_Run_ErrorFetchingQuotes(t *testing.T) {
	ctx := context.Background()

	pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
	pegoutQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateRefundPegOutSucceeded).
		Return(nil, assert.AnError).Once()

	useCase := pegout.NewGetPegoutReportUseCase(pegoutQuoteRepository)

	result, err := useCase.Run(ctx)

	pegoutQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
	assert.Zero(t, result.NumberOfQuotes)
}
