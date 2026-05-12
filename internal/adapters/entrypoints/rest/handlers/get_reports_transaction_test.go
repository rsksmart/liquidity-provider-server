package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetReportsTransactionHandler_SuccessfulPeginRequest(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	mockQuotes := []quote.PeginQuoteWithRetained{
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
	}

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 1, 10).
		Return(mockQuotes, 1, nil)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01&endDate=2023-01-31&page=1&perPage=10", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pkg.GetTransactionsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Data, 1)
	assert.Equal(t, "0x123", response.Data[0].QuoteHash)
	assert.Equal(t, 1, response.Pagination.Total)
	assert.Equal(t, 10, response.Pagination.PerPage)
}

func TestGetReportsTransactionHandler_SuccessfulPegoutRequest(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	mockQuotes := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				Value:   entities.NewWei(2000000000000000000),
				CallFee: entities.NewWei(100000000000000000),
				GasFee:  entities.NewWei(20000000000000000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "0x456",
				State:     quote.PegoutStateBridgeTxSucceeded,
			},
		},
	}

	pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 2, 5).
		Return(mockQuotes, 25, nil)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegout&startDate=2023-01-01&endDate=2023-01-31&page=2&perPage=5", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pkg.GetTransactionsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Data, 1)
	assert.Equal(t, "0x456", response.Data[0].QuoteHash)
	assert.Equal(t, 25, response.Pagination.Total)
	assert.Equal(t, 5, response.Pagination.PerPage)
}

func TestGetReportsTransactionHandler_MissingTypeParameter(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?startDate=2023-01-01&endDate=2023-01-31", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetReportsTransactionHandler_InvalidTypeParameter(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=invalid&startDate=2023-01-01&endDate=2023-01-31", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetReportsTransactionHandler_InvalidPageParameter(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01&endDate=2023-01-31&page=invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetReportsTransactionHandler_InvalidPerPageParameter(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01&endDate=2023-01-31&perPage=invalid", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetReportsTransactionHandler_InvalidDateFormat(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=invalid-date&endDate=2023-01-31", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetReportsTransactionHandler_UseCaseError(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 1, 10).
		Return([]quote.PeginQuoteWithRetained{}, 0, assert.AnError)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01&endDate=2023-01-31", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetReportsTransactionHandler_DefaultPaginationParameters(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 1, 10).
		Return([]quote.PeginQuoteWithRetained{}, 0, nil)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01&endDate=2023-01-31", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pkg.GetTransactionsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, response.Pagination.Total)
	assert.Equal(t, 10, response.Pagination.PerPage)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 1, response.Pagination.TotalPages)
}

func TestGetReportsTransactionHandler_ISO8601DateFormat(t *testing.T) {
	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	startDate := time.Date(2023, 1, 1, 9, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 1, 17, 0, 0, 0, time.UTC)

	peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 1, 10).
		Return([]quote.PeginQuoteWithRetained{}, 0, nil)

	useCase := reports.NewGetTransactionsUseCase(peginRepo, pegoutRepo)
	handler := handlers.NewGetReportsTransactionHandler(useCase)

	req := httptest.NewRequest(http.MethodGet, "/reports/transactions?type=pegin&startDate=2023-01-01T09:00:00Z&endDate=2023-01-01T17:00:00Z", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response pkg.GetTransactionsResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, response.Pagination.Total)
	assert.Equal(t, 10, response.Pagination.PerPage)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 1, response.Pagination.TotalPages)
}
