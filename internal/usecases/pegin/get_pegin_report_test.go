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
	"time"
)

// nolint:funlen
func TestGetPeginReportUseCase_Run(t *testing.T) {
	ctx := context.Background()

	retainedQuotes := []quote.RetainedPeginQuote{
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

	peginQuotes := []quote.PeginQuote{
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

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 365 * 10)

	filters := []quote.QueryFilter{
		{
			Field:    "agreement_timestamp",
			Operator: "$gte",
			Value:    startDate.Unix(),
		},
		{
			Field:    "agreement_timestamp",
			Operator: "$lte",
			Value:    endDate.Unix(),
		},
	}

	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}

	peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
		Return(retainedQuotes, nil).Once()

	peginQuoteRepository.On("GetQuotes", ctx, filters, quoteHashes).Return(peginQuotes, nil).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx, startDate, endDate)

	peginQuoteRepository.AssertExpectations(t)
	require.NoError(t, err)
	assert.Equal(t, 10, result.NumberOfQuotes)
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

	result, err := useCase.Run(ctx, time.Now(), time.Now().Add(time.Hour*24*365*10))

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

	result, err := useCase.Run(ctx, time.Now(), time.Now().Add(time.Hour*24*365*10))

	peginQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
	assert.Zero(t, result.NumberOfQuotes)
}
