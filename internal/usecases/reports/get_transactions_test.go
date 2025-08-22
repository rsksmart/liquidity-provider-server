package reports_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTransactionsUseCase_Run_PeginTransactions(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)

	ctx := context.Background()
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	mockQuotePairs := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				Value:   entities.NewWei(1000000000000000000),
				CallFee: entities.NewWei(50000000000000000),
				GasFee:  entities.NewWei(10000000000000000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "0x123",
				State:     quote.PeginStateRegisterPegInSucceeded,
			},
		},
		{
			Quote: quote.PeginQuote{
				Value:   entities.NewWei(2000000000000000000),
				CallFee: entities.NewWei(100000000000000000),
				GasFee:  entities.NewWei(20000000000000000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "0x456",
				State:     quote.PeginStateWaitingForDepositConfirmations,
			},
		},
	}

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startTime, endTime, 1, 10).
		Return(mockQuotePairs, 50, nil)

	result, err := useCase.Run(ctx, "pegin", startTime, endTime, 1, 10)

	require.NoError(t, err)
	assert.Len(t, result.Data, 2)

	assert.Equal(t, "0x123", result.Data[0].QuoteHash)
	assert.Equal(t, "1000000000000000000", result.Data[0].Amount.String())
	assert.Equal(t, "50000000000000000", result.Data[0].CallFee.String())
	assert.Equal(t, "10000000000000000", result.Data[0].GasFee.String())
	assert.Equal(t, string(quote.PeginStateRegisterPegInSucceeded), result.Data[0].Status)

	assert.Equal(t, "0x456", result.Data[1].QuoteHash)
	assert.Equal(t, "2000000000000000000", result.Data[1].Amount.String())
	assert.Equal(t, string(quote.PeginStateWaitingForDepositConfirmations), result.Data[1].Status)

	assert.Equal(t, 50, result.Pagination.Total)
	assert.Equal(t, 10, result.Pagination.PerPage)
	assert.Equal(t, 5, result.Pagination.TotalPages)
	assert.Equal(t, 1, result.Pagination.Page)

	peginRepo.AssertExpectations(t)
}

func TestGetTransactionsUseCase_Run_PegoutTransactions(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)

	ctx := context.Background()
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	mockQuotePairs := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				Value:   entities.NewWei(500000000000000000),
				CallFee: entities.NewWei(25000000000000000),
				GasFee:  entities.NewWei(5000000000000000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "0x789",
				State:     quote.PegoutStateBridgeTxSucceeded,
			},
		},
	}

	pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startTime, endTime, 2, 5).
		Return(mockQuotePairs, 25, nil)

	result, err := useCase.Run(ctx, "pegout", startTime, endTime, 2, 5)

	require.NoError(t, err)
	assert.Len(t, result.Data, 1)

	assert.Equal(t, "0x789", result.Data[0].QuoteHash)
	assert.Equal(t, "500000000000000000", result.Data[0].Amount.String())
	assert.Equal(t, string(quote.PegoutStateBridgeTxSucceeded), result.Data[0].Status)

	assert.Equal(t, 25, result.Pagination.Total)
	assert.Equal(t, 5, result.Pagination.PerPage)
	assert.Equal(t, 5, result.Pagination.TotalPages)
	assert.Equal(t, 2, result.Pagination.Page)

	pegoutRepo.AssertExpectations(t)
}

func TestGetTransactionsUseCase_Run_InvalidType(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)

	ctx := context.Background()
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	result, err := useCase.Run(ctx, "invalid", startTime, endTime, 1, 10)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transaction type")
	assert.Empty(t, result.Data)
}

func TestGetTransactionsUseCase_Run_EmptyResults(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)

	ctx := context.Background()
	startTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startTime, endTime, 1, 10).
		Return([]quote.PeginQuoteWithRetained{}, 0, nil)

	result, err := useCase.Run(ctx, "pegin", startTime, endTime, 1, 10)

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, 0, result.Pagination.Total)
	assert.Equal(t, 1, result.Pagination.TotalPages) // Minimum 1 page

	peginRepo.AssertExpectations(t)
}

func TestCalculatePaginationMetadata(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		perPage    int
		totalCount int
		expected   reports.PaginationMetadata
	}{
		{
			name:       "First page with full results",
			page:       1,
			perPage:    10,
			totalCount: 100,
			expected: reports.PaginationMetadata{
				Total:      100,
				PerPage:    10,
				TotalPages: 10,
				Page:       1,
			},
		},
		{
			name:       "Middle page",
			page:       5,
			perPage:    20,
			totalCount: 150,
			expected: reports.PaginationMetadata{
				Total:      150,
				PerPage:    20,
				TotalPages: 8, // Ceiling of 150/20
				Page:       5,
			},
		},
		{
			name:       "No results",
			page:       1,
			perPage:    10,
			totalCount: 0,
			expected: reports.PaginationMetadata{
				Total:      0,
				PerPage:    10,
				TotalPages: 1, // Minimum 1 page
				Page:       1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reports.CalculatePaginationMetadata(tt.page, tt.perPage, tt.totalCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}
