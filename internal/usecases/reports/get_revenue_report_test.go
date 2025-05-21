package reports_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// nolint:funlen
func TestGetRevenueReportUseCase_Run(t *testing.T) {
	retainedPeginQuotes := []quote.RetainedPeginQuote{
		{QuoteHash: "hash1"},
		{QuoteHash: "hash2"},
		{QuoteHash: "hash3"},
		{QuoteHash: "hash4"},
		{QuoteHash: "hash5"},
	}

	retainedPegoutQuotes := []quote.RetainedPegoutQuote{
		{QuoteHash: "hash6"},
		{QuoteHash: "hash7"},
		{QuoteHash: "hash8"},
		{QuoteHash: "hash9"},
		{QuoteHash: "hash10"},
	}

	peginQuoteHashes := []string{"hash1", "hash2", "hash3", "hash4", "hash5"}
	pegoutQuoteHashes := []string{"hash6", "hash7", "hash8", "hash9", "hash10"}
	allQuoteHashes := append(peginQuoteHashes, pegoutQuoteHashes...)

	peginQuotes := []quote.PeginQuote{
		{Value: entities.NewWei(1000), CallFee: entities.NewWei(100)},
		{Value: entities.NewWei(2000), CallFee: entities.NewWei(200)},
		{Value: entities.NewWei(3000), CallFee: entities.NewWei(300)},
		{Value: entities.NewWei(4000), CallFee: entities.NewWei(400)},
		{Value: entities.NewWei(5000), CallFee: entities.NewWei(500)},
	}
	pegoutQuotes := []quote.PegoutQuote{
		{Value: entities.NewWei(6000), CallFee: entities.NewWei(600)},
		{Value: entities.NewWei(7000), CallFee: entities.NewWei(700)},
		{Value: entities.NewWei(8000), CallFee: entities.NewWei(800)},
		{Value: entities.NewWei(9000), CallFee: entities.NewWei(900)},
		{Value: entities.NewWei(10000), CallFee: entities.NewWei(1000)},
	}
	penalizedEvents := []penalization.PenalizedEvent{
		{QuoteHash: "hash1", Penalty: entities.NewWei(100)},
		{QuoteHash: "hash7", Penalty: entities.NewWei(200)},
		{QuoteHash: "hash9", Penalty: entities.NewWei(300)},
	}

	t.Run("Should get the report correctly", func(t *testing.T) {
		ctx := context.Background()

		startDate := time.Now()
		endDate := time.Now().Add(time.Hour * 24 * 365 * 10)

		peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
		penalizedRepository := &mocks.PenalizedEventRepositoryMock{}

		peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
			Return(retainedPeginQuotes, nil).Once()
		peginQuoteRepository.On("GetQuotesByHashesAndDate", ctx, peginQuoteHashes, startDate, endDate).Return(peginQuotes, nil).Once()
		pegoutQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded).
			Return(retainedPegoutQuotes, nil).Once()
		pegoutQuoteRepository.On("GetQuotesByHashesAndDate", ctx, pegoutQuoteHashes, startDate, endDate).Return(pegoutQuotes, nil).Once()
		penalizedRepository.On("GetPenalizationsByQuoteHashes", ctx, allQuoteHashes).Return(penalizedEvents, nil).Once()

		useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepository, pegoutQuoteRepository, penalizedRepository)

		result, err := useCase.Run(ctx, startDate, endDate)

		peginQuoteRepository.AssertExpectations(t)
		pegoutQuoteRepository.AssertExpectations(t)
		penalizedRepository.AssertExpectations(t)
		require.NoError(t, err)
		assert.Equal(t, 5500, int(result.TotalQuoteCallFees.Uint64()))
		assert.Equal(t, 600, int(result.TotalPenalizations.Uint64()))
		assert.Equal(t, 4900, int(result.TotalProfit.Uint64()))
	})
	t.Run("Should return an error when fetching pegin retained quotes fails", func(t *testing.T) {
		ctx := context.Background()

		startDate := time.Now()
		endDate := time.Now().Add(time.Hour * 24 * 365 * 10)

		peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
		penalizedRepository := &mocks.PenalizedEventRepositoryMock{}

		peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
			Return(nil, assert.AnError).Once()

		useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepository, pegoutQuoteRepository, penalizedRepository)

		_, err := useCase.Run(ctx, startDate, endDate)

		peginQuoteRepository.AssertExpectations(t)
		pegoutQuoteRepository.AssertExpectations(t)
		penalizedRepository.AssertExpectations(t)
		require.Error(t, err)
	})
	t.Run("Should return an error when fetching pegin quotes fails", func(t *testing.T) {
		ctx := context.Background()

		startDate := time.Now()
		endDate := time.Now().Add(time.Hour * 24 * 365 * 10)

		peginQuoteRepository := &mocks.PeginQuoteRepositoryMock{}
		pegoutQuoteRepository := &mocks.PegoutQuoteRepositoryMock{}
		penalizedRepository := &mocks.PenalizedEventRepositoryMock{}

		peginQuoteRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateRegisterPegInSucceeded).
			Return(retainedPeginQuotes, nil).Once()
		peginQuoteRepository.On("GetQuotesByHashesAndDate", ctx, peginQuoteHashes, startDate, endDate).Return(nil, assert.AnError).Once()

		useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepository, pegoutQuoteRepository, penalizedRepository)

		_, err := useCase.Run(ctx, startDate, endDate)

		peginQuoteRepository.AssertExpectations(t)
		pegoutQuoteRepository.AssertExpectations(t)
		penalizedRepository.AssertExpectations(t)
		require.Error(t, err)
	})
}
