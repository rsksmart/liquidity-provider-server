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

func TestGetPeginReportUseCase_Run(t *testing.T) {
	ctx := context.Background()

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

	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
	filter := quote.GetPeginQuotesByStateFilter{
		States:    []quote.PeginState{quote.PeginStateRegisterPegInSucceeded},
		StartDate: uint32(time.Now().Unix()),
		EndDate:   uint32(time.Now().Add(time.Hour * 24 * 365 * 10).Unix()),
	}

	peginQuoteRepository.On("GetQuotesByState", ctx, filter).Return(peginQuotes, nil).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx, time.Now(), time.Now().Add(time.Hour*24*365*10))

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

	quotes := make([]quote.PeginQuote, 0)
	peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
	filter := quote.GetPeginQuotesByStateFilter{
		States:    []quote.PeginState{quote.PeginStateRegisterPegInSucceeded},
		StartDate: uint32(time.Now().Unix()),
		EndDate:   uint32(time.Now().Add(time.Hour * 24 * 365 * 10).Unix()),
	}
	peginQuoteRepository.On("GetQuotesByState", ctx, filter).
		Return(quotes, nil).Once()

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
	filter := quote.GetPeginQuotesByStateFilter{
		States:    []quote.PeginState{quote.PeginStateRegisterPegInSucceeded},
		StartDate: uint32(time.Now().Unix()),
		EndDate:   uint32(time.Now().Add(time.Hour * 24 * 365 * 10).Unix()),
	}
	peginQuoteRepository.On("GetQuotesByState", ctx, filter).
		Return(nil, assert.AnError).Once()

	useCase := pegin.NewGetPeginReportUseCase(peginQuoteRepository)

	result, err := useCase.Run(ctx, time.Now(), time.Now().Add(time.Hour*24*365*10))

	peginQuoteRepository.AssertExpectations(t)
	require.Error(t, err)
	assert.Zero(t, result.NumberOfQuotes)
}
