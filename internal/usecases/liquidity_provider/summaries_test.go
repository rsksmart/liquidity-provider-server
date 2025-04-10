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

func TestSummariesUseCase_Run(t *testing.T) { //nolint:funlen,maintidx
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
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		peginRepo.On("GetQuote", mock.Anything, "hash1").
			Return(&peginQuotes[0], nil)
		peginRepo.On("GetQuote", mock.Anything, "hash2").
			Return(&peginQuotes[1], nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotes, retainedPegoutQuotes, nil)
		pegoutRepo.On("GetQuote", mock.Anything, "hash3").
			Return(&pegoutQuotes[0], nil)
		pegoutRepo.On("GetQuote", mock.Anything, "hash4").
			Return(&pegoutQuotes[1], nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(2), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(660)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(10)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(5)))
		assert.Equal(t, int64(len(pegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(2), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(1540)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(30)))
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
		retainedPeginQuotes := []quote.RetainedPeginQuote{}
		pegoutQuotes := []quote.PegoutQuote{
			{
				Value:            entities.NewWei(300),
				CallFee:          entities.NewWei(15),
				GasFee:           entities.NewWei(6),
				PenaltyFee:       10,
				ProductFeeAmount: 9,
			},
		}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotes, retainedPegoutQuotes, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(330)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(0)))
		assert.Equal(t, int64(len(pegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(330)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.LpEarnings.Cmp(entities.NewWei(0)))
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Success with only retained quotes (no regular quotes)", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginQuotes := []quote.PeginQuote{}
		retainedPeginQuotes := []quote.RetainedPeginQuote{
			{
				QuoteHash:         "hash1",
				Signature:         "sig1",
				DepositAddress:    "addr1",
				State:             quote.PeginStateCallForUserSucceeded,
				UserBtcTxHash:     "user_tx1",
				CallForUserTxHash: "call_tx1",
			},
		}
		pegoutQuotes := []quote.PegoutQuote{}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{
			{
				QuoteHash: "hash3",
				Signature: "sig3",
				State:     quote.PegoutStateSendPegoutSucceeded,
			},
		}
		peginQuote := quote.PeginQuote{
			Value:            entities.NewWei(100),
			CallFee:          entities.NewWei(5),
			GasFee:           entities.NewWei(2),
			PenaltyFee:       entities.NewWei(1),
			ProductFeeAmount: 3,
		}
		pegoutQuote := quote.PegoutQuote{
			Value:            entities.NewWei(300),
			CallFee:          entities.NewWei(15),
			GasFee:           entities.NewWei(6),
			PenaltyFee:       10,
			ProductFeeAmount: 9,
		}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		peginRepo.On("GetQuote", mock.Anything, "hash1").
			Return(&peginQuote, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotes, retainedPegoutQuotes, nil)
		pegoutRepo.On("GetQuote", mock.Anything, "hash3").
			Return(&pegoutQuote, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(len(peginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PeginSummary.PaidQuotesAmount.Cmp(entities.NewWei(110)))
		assert.Equal(t, 0, result.PeginSummary.TotalFeesCollected.Cmp(entities.NewWei(10)))
		assert.Equal(t, 0, result.PeginSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PeginSummary.LpEarnings.Cmp(entities.NewWei(5)))
		assert.Equal(t, int64(len(pegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, 0, result.PegoutSummary.PaidQuotesAmount.Cmp(entities.NewWei(330)))
		assert.Equal(t, 0, result.PegoutSummary.TotalFeesCollected.Cmp(entities.NewWei(30)))
		assert.Equal(t, 0, result.PegoutSummary.TotalPenaltyAmount.Cmp(entities.NewWei(0)))
		assert.Equal(t, 0, result.PegoutSummary.LpEarnings.Cmp(entities.NewWei(15)))
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
	t.Run("Error getting pegin quotes", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PeginQuote{}, []quote.RetainedPeginQuote{}, errors.New("db error"))
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
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
			Return([]quote.PeginQuote{}, []quote.RetainedPeginQuote{}, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PegoutQuote{}, []quote.RetainedPegoutQuote{}, errors.New("db error"))
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		_, err := useCase.Run(context.Background(), startDate, endDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
}
