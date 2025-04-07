package liquidity_provider_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPeginQuoteRepo struct {
	mock.Mock
}

func (m *mockPeginQuoteRepo) InsertQuote(ctx context.Context, q quote.CreatedPeginQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPeginQuoteRepo) InsertRetainedQuote(ctx context.Context, q quote.RetainedPeginQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPeginQuoteRepo) UpdateRetainedQuote(ctx context.Context, q quote.RetainedPeginQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPeginQuoteRepo) UpdateRetainedQuotes(ctx context.Context, quotes []quote.RetainedPeginQuote) error {
	args := m.Called(ctx, quotes)
	return args.Error(0)
}

func (m *mockPeginQuoteRepo) GetQuote(ctx context.Context, hash string) (*quote.PeginQuote, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	quote, ok := args.Get(0).(*quote.PeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid pegin quote type")
	}
	
	return quote, args.Error(1)
}

func (m *mockPeginQuoteRepo) GetPeginCreationData(ctx context.Context, hash string) quote.PeginCreationData {
	args := m.Called(ctx, hash)
	
	data, ok := args.Get(0).(quote.PeginCreationData)
	if !ok && args.Get(0) != nil {
		panic("invalid pegin creation data type")
	}
	
	return data
}

func (m *mockPeginQuoteRepo) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	quote, ok := args.Get(0).(*quote.RetainedPeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid retained pegin quote type")
	}
	
	return quote, args.Error(1)
}

func (m *mockPeginQuoteRepo) GetRetainedQuoteByDepositAddress(ctx context.Context, address string) (*quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	quote, ok := args.Get(0).(*quote.RetainedPeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid retained pegin quote type")
	}
	
	return quote, args.Error(1)
}

func (m *mockPeginQuoteRepo) GetRetainedQuoteByState(ctx context.Context, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, states)
	
	quotes, ok := args.Get(0).([]quote.RetainedPeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid retained pegin quotes type")
	}
	
	return quotes, args.Error(1)
}

func (m *mockPeginQuoteRepo) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PeginQuote, []quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, startDate, endDate)
	
	quotes, ok1 := args.Get(0).([]quote.PeginQuote)
	if !ok1 && args.Get(0) != nil {
		return nil, nil, errors.New("invalid pegin quotes type")
	}
	
	retainedQuotes, ok2 := args.Get(1).([]quote.RetainedPeginQuote)
	if !ok2 && args.Get(1) != nil {
		return nil, nil, errors.New("invalid retained pegin quotes type")
	}
	
	return quotes, retainedQuotes, args.Error(2)
}

func (m *mockPeginQuoteRepo) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	args := m.Called(ctx, quotes)
	
	count, ok := args.Get(0).(uint)
	if !ok && args.Get(0) != nil {
		return 0, errors.New("invalid count type")
	}
	
	return count, args.Error(1)
}

type mockPegoutQuoteRepo struct {
	mock.Mock
}

func (m *mockPegoutQuoteRepo) InsertQuote(ctx context.Context, q quote.CreatedPegoutQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPegoutQuoteRepo) InsertRetainedQuote(ctx context.Context, q quote.RetainedPegoutQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPegoutQuoteRepo) UpdateRetainedQuote(ctx context.Context, q quote.RetainedPegoutQuote) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockPegoutQuoteRepo) UpdateRetainedQuotes(ctx context.Context, quotes []quote.RetainedPegoutQuote) error {
	args := m.Called(ctx, quotes)
	return args.Error(0)
}

func (m *mockPegoutQuoteRepo) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	quote, ok := args.Get(0).(*quote.PegoutQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid pegout quote type")
	}
	
	return quote, args.Error(1)
}

func (m *mockPegoutQuoteRepo) GetPegoutCreationData(ctx context.Context, hash string) quote.PegoutCreationData {
	args := m.Called(ctx, hash)
	
	data, ok := args.Get(0).(quote.PegoutCreationData)
	if !ok && args.Get(0) != nil {
		panic("invalid pegout creation data type")
	}
	
	return data
}

func (m *mockPegoutQuoteRepo) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	quote, ok := args.Get(0).(*quote.RetainedPegoutQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid retained pegout quote type")
	}
	
	return quote, args.Error(1)
}

func (m *mockPegoutQuoteRepo) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, states)
	
	quotes, ok := args.Get(0).([]quote.RetainedPegoutQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid retained pegout quotes type")
	}
	
	return quotes, args.Error(1)
}

func (m *mockPegoutQuoteRepo) ListPegoutDepositsByAddress(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	args := m.Called(ctx, address)
	
	deposits, ok := args.Get(0).([]quote.PegoutDeposit)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid pegout deposits type")
	}
	
	return deposits, args.Error(1)
}

func (m *mockPegoutQuoteRepo) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PegoutQuote, []quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, startDate, endDate)
	
	quotes, ok1 := args.Get(0).([]quote.PegoutQuote)
	if !ok1 && args.Get(0) != nil {
		return nil, nil, errors.New("invalid pegout quotes type")
	}
	
	retainedQuotes, ok2 := args.Get(1).([]quote.RetainedPegoutQuote)
	if !ok2 && args.Get(1) != nil {
		return nil, nil, errors.New("invalid retained pegout quotes type")
	}
	
	return quotes, retainedQuotes, args.Error(2)
}

func (m *mockPegoutQuoteRepo) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	args := m.Called(ctx, quotes)
	
	count, ok := args.Get(0).(uint)
	if !ok && args.Get(0) != nil {
		return 0, errors.New("invalid count type")
	}
	
	return count, args.Error(1)
}

func (m *mockPegoutQuoteRepo) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	args := m.Called(ctx, deposit)
	return args.Error(0)
}

func (m *mockPegoutQuoteRepo) UpsertPegoutDeposits(ctx context.Context, deposits []quote.PegoutDeposit) error {
	args := m.Called(ctx, deposits)
	return args.Error(0)
}

func TestSummariesUseCase_Run(t *testing.T) { //nolint:funlen,maintidx
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("Success with both pegin and pegout data", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
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
				State:     quote.PegoutStateBridgeTxSucceeded,
			},
			{
				QuoteHash: "hash4",
				Signature: "sig4",
				State:     quote.PegoutStateBridgeTxFailed,
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

		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, "660", result.PeginSummary.TotalQuotedAmount)
		assert.Equal(t, "110", result.PeginSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "10", result.PeginSummary.TotalFeesCollected)
		assert.Equal(t, "2", result.PeginSummary.TotalPenaltyAmount)
		assert.Equal(t, "3", result.PeginSummary.LpEarnings)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, "1540", result.PegoutSummary.TotalQuotedAmount)
		assert.Equal(t, "330", result.PegoutSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "30", result.PegoutSummary.TotalFeesCollected)
		assert.Equal(t, "15", result.PegoutSummary.TotalPenaltyAmount)
		assert.Equal(t, "0", result.PegoutSummary.LpEarnings)
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Success with only regular quotes (no retained quotes)", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
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
		assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, "330", result.PeginSummary.TotalQuotedAmount)
		assert.Equal(t, "0", result.PeginSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "0", result.PeginSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PeginSummary.TotalPenaltyAmount)
		assert.Equal(t, "0", result.PeginSummary.LpEarnings)
		assert.Equal(t, int64(len(pegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, "330", result.PegoutSummary.TotalQuotedAmount)
		assert.Equal(t, "0", result.PegoutSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "0", result.PegoutSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PegoutSummary.TotalPenaltyAmount)
		assert.Equal(t, "0", result.PegoutSummary.LpEarnings)

		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Success with only retained quotes (no regular quotes)", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
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
				State:     quote.PegoutStateBridgeTxSucceeded,
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

		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, "110", result.PeginSummary.TotalQuotedAmount)
		assert.Equal(t, "110", result.PeginSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "10", result.PeginSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PeginSummary.TotalPenaltyAmount)
		assert.Equal(t, "5", result.PeginSummary.LpEarnings)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, "330", result.PegoutSummary.TotalQuotedAmount)
		assert.Equal(t, "330", result.PegoutSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "30", result.PegoutSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PegoutSummary.TotalPenaltyAmount)
		assert.Equal(t, "15", result.PegoutSummary.LpEarnings)

		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Success with duplicate quotes (same hash in both regular and retained)", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)

		peginQuote := quote.PeginQuote{
			Value:            entities.NewWei(100),
			CallFee:          entities.NewWei(5),
			GasFee:           entities.NewWei(2),
			PenaltyFee:       entities.NewWei(1),
			ProductFeeAmount: 3,
		}

		peginQuotes := []quote.PeginQuote{peginQuote}
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

		pegoutQuote := quote.PegoutQuote{
			Value:            entities.NewWei(300),
			CallFee:          entities.NewWei(15),
			GasFee:           entities.NewWei(6),
			PenaltyFee:       10,
			ProductFeeAmount: 9,
		}

		pegoutQuotes := []quote.PegoutQuote{pegoutQuote}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{
			{
				QuoteHash: "hash3",
				Signature: "sig3",
				State:     quote.PegoutStateBridgeTxSucceeded,
			},
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

		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, "220", result.PeginSummary.TotalQuotedAmount)
		assert.Equal(t, "110", result.PeginSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "10", result.PeginSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PeginSummary.TotalPenaltyAmount)
		assert.Equal(t, "5", result.PeginSummary.LpEarnings)
		assert.Equal(t, int64(len(retainedPegoutQuotes)), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, "660", result.PegoutSummary.TotalQuotedAmount)
		assert.Equal(t, "330", result.PegoutSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "30", result.PegoutSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PegoutSummary.TotalPenaltyAmount)
		assert.Equal(t, "15", result.PegoutSummary.LpEarnings)

		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Error retrieving pegin quotes", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
		expectedError := errors.New("database error")
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PeginQuote{}, []quote.RetainedPeginQuote{}, expectedError)

		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, liquidity_provider.SummariesResponse{}, result)
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Error retrieving pegout quotes", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
		peginQuotes := []quote.PeginQuote{}
		retainedPeginQuotes := []quote.RetainedPeginQuote{}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		expectedError := errors.New("database error")
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return([]quote.PegoutQuote{}, []quote.RetainedPegoutQuote{}, expectedError)

		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, liquidity_provider.SummariesResponse{}, result)
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("Missing quote in pegin data", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
		peginQuotes := []quote.PeginQuote{
			{
				Value:            entities.NewWei(100),
				CallFee:          entities.NewWei(5),
				GasFee:           entities.NewWei(2),
				PenaltyFee:       entities.NewWei(1),
				ProductFeeAmount: 3,
			},
		}

		retainedPeginQuotes := []quote.RetainedPeginQuote{
			{
				QuoteHash:      "hash1",
				Signature:      "sig1",
				DepositAddress: "addr1",
				State:          quote.PeginStateCallForUserSucceeded,
			},
			{
				QuoteHash:      "hash_missing",
				Signature:      "sig2",
				DepositAddress: "addr2",
				State:          quote.PeginStateCallForUserFailed,
			},
		}
		pegoutQuotes := []quote.PegoutQuote{}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{}

		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		peginRepo.On("GetQuote", mock.Anything, "hash1").
			Return(&peginQuotes[0], nil)
		peginRepo.On("GetQuote", mock.Anything, "hash_missing").
			Return(nil, errors.New("quote not found"))
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotes, retainedPegoutQuotes, nil)

		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)

		require.NoError(t, err)
		assert.Equal(t, int64(len(retainedPeginQuotes)), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(1), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})

	t.Run("No quotes found", func(t *testing.T) {
		peginRepo := new(mockPeginQuoteRepo)
		pegoutRepo := new(mockPegoutQuoteRepo)
		peginQuotes := []quote.PeginQuote{}
		retainedPeginQuotes := []quote.RetainedPeginQuote{}
		pegoutQuotes := []quote.PegoutQuote{}
		retainedPegoutQuotes := []quote.RetainedPegoutQuote{}
		peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(peginQuotes, retainedPeginQuotes, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate).
			Return(pegoutQuotes, retainedPegoutQuotes, nil)
		useCase := liquidity_provider.NewSummariesUseCase(peginRepo, pegoutRepo)
		result, err := useCase.Run(context.Background(), startDate, endDate)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, "0", result.PeginSummary.TotalQuotedAmount)
		assert.Equal(t, "0", result.PeginSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "0", result.PeginSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PeginSummary.TotalPenaltyAmount)
		assert.Equal(t, "0", result.PeginSummary.LpEarnings)
		assert.Equal(t, int64(0), result.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, "0", result.PegoutSummary.TotalQuotedAmount)
		assert.Equal(t, "0", result.PegoutSummary.TotalAcceptedQuotedAmount)
		assert.Equal(t, "0", result.PegoutSummary.TotalFeesCollected)
		assert.Equal(t, "0", result.PegoutSummary.TotalPenaltyAmount)
		assert.Equal(t, "0", result.PegoutSummary.LpEarnings)
		peginRepo.AssertExpectations(t)
		pegoutRepo.AssertExpectations(t)
	})
}
