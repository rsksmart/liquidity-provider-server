package liquidity_provider_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSummariesUseCase_Run(t *testing.T) { //nolint:funlen
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)
	t.Run("Success with full set of data", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginQuotes := []quote.PeginQuote{
			{
				Value:            entities.NewWei(100),
				CallFee:          entities.NewWei(5),
				GasFee:           entities.NewWei(2),
				PenaltyFee:       entities.NewWei(1),
				ProductFeeAmount: 3,
			},
			{
				Value:            entities.NewWei(200),
				CallFee:          entities.NewWei(10),
				GasFee:           entities.NewWei(4),
				PenaltyFee:       entities.NewWei(2),
				ProductFeeAmount: 6,
			},
		}
		retainedPeginQuotes := []quote.RetainedPeginQuote{
			{
				QuoteHash:         "hash1",
				Signature:         "sig1",
				DepositAddress:    "addr1",
				State:             quote.PeginStateCallForUserSucceeded,
				UserBtcTxHash:     "user_tx1",
				CallForUserTxHash: "call_tx1",
			},
			{
				QuoteHash:         "hash2",
				Signature:         "sig2",
				DepositAddress:    "addr2",
				State:             quote.PeginStateCallForUserFailed,
				UserBtcTxHash:     "user_tx2",
				CallForUserTxHash: "",
			},
		}
		pegoutQuotes := []quote.PegoutQuote{
			{
				Value:            entities.NewWei(300),
				CallFee:          entities.NewWei(15),
				GasFee:           entities.NewWei(6),
				PenaltyFee:       10,
				ProductFeeAmount: 9,
			},
			{
				Value:            entities.NewWei(400),
				CallFee:          entities.NewWei(20),
				GasFee:           entities.NewWei(8),
				PenaltyFee:       15,
				ProductFeeAmount: 12,
			},
		}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{
			{
				QuoteHash: "hash3",
				Signature: "sig3",
				State:     quote.PegoutStateSendPegoutSucceeded,
			},
			{
				QuoteHash: "hash4",
				Signature: "sig4",
				State:     quote.PegoutStateSendPegoutFailed,
			},
		}
		peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
			{
				Quote:         peginQuotes[0],
				RetainedQuote: retainedPeginQuotes[0],
			},
			{
				Quote:         peginQuotes[1],
				RetainedQuote: retainedPeginQuotes[1],
			},
		}
		pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
			{
				Quote:         pegoutQuotes[0],
				RetainedQuote: retainedPegoutQuotes[0],
			},
			{
				Quote:         pegoutQuotes[1],
				RetainedQuote: retainedPegoutQuotes[1],
			},
		}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotesWithRetained, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotesWithRetained, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo, nil)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotesWithRetained)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(100)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(7)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(5)))
		assert.Equal(t, int64(len(pegoutQuotesWithRetained)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(300)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(21)))
		assert.Equal(t, 0, result.PegoutSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.LpEarnings.Cmp(entities.NewWei(15)))
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Success with only regular quotes (no retained quotes)", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginQuotes := []quote.PeginQuote{
			{
				Value:            entities.NewWei(100),
				CallFee:          entities.NewWei(5),
				GasFee:           entities.NewWei(2),
				PenaltyFee:       entities.NewWei(1),
				ProductFeeAmount: 3,
			},
			{
				Value:            entities.NewWei(200),
				CallFee:          entities.NewWei(10),
				GasFee:           entities.NewWei(4),
				PenaltyFee:       entities.NewWei(2),
				ProductFeeAmount: 6,
			},
		}
		pegoutQuotes := []quote.PegoutQuote{
			{
				Value:            entities.NewWei(300),
				CallFee:          entities.NewWei(15),
				GasFee:           entities.NewWei(6),
				PenaltyFee:       10,
				ProductFeeAmount: 9,
			},
		}
		peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
			{
				Quote:         peginQuotes[0],
				RetainedQuote: quote.RetainedPeginQuote{},
			},
			{
				Quote:         peginQuotes[1],
				RetainedQuote: quote.RetainedPeginQuote{},
			},
		}
		pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
			{
				Quote:         pegoutQuotes[0],
				RetainedQuote: quote.RetainedPegoutQuote{},
			},
		}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotesWithRetained, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotesWithRetained, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo, nil)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotesWithRetained)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(0)))
		assert.Equal(t, int64(len(pegoutQuotesWithRetained)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.LpEarnings.Cmp(entities.NewWei(0)))
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Success with only retained quotes (no regular quotes)", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginQuote := quote.PeginQuote{
			Value:            entities.NewWei(100),
			CallFee:          entities.NewWei(5),
			GasFee:           entities.NewWei(2),
			PenaltyFee:       entities.NewWei(1),
			ProductFeeAmount: 3,
		}
		retainedPeginQuote := quote.RetainedPeginQuote{
			QuoteHash:         "hash1",
			Signature:         "sig1",
			DepositAddress:    "addr1",
			State:             quote.PeginStateCallForUserSucceeded,
			UserBtcTxHash:     "user_tx1",
			CallForUserTxHash: "call_tx1",
		}
		pegoutQuote := quote.PegoutQuote{
			Value:            entities.NewWei(300),
			CallFee:          entities.NewWei(15),
			GasFee:           entities.NewWei(6),
			PenaltyFee:       10,
			ProductFeeAmount: 9,
		}
		retainedPegoutQuote := quote.RetainedPegoutQuote{
			QuoteHash: "hash3",
			Signature: "sig3",
			State:     quote.PegoutStateSendPegoutSucceeded,
		}
		peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
			{
				Quote:         peginQuote,
				RetainedQuote: retainedPeginQuote,
			},
		}
		pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
			{
				Quote:         pegoutQuote,
				RetainedQuote: retainedPegoutQuote,
			},
		}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotesWithRetained, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotesWithRetained, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo, nil)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotesWithRetained)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(100)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(7)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(5)))
		assert.Equal(t, int64(len(pegoutQuotesWithRetained)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(300)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(21)))
		assert.Equal(t, 0, result.PegoutSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.LpEarnings.Cmp(entities.NewWei(15)))
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Error getting pegin quotes", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PeginQuoteWithRetained{}, errors.New("db error"))
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo, nil)
		_, err := useCase.Run(context.Background(), startDate, endDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Error getting pegout quotes", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PeginQuoteWithRetained{}, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PegoutQuoteWithRetained{}, errors.New("db error"))
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo, nil)
		_, err := useCase.Run(context.Background(), startDate, endDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
}
