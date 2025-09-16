package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetReportSummariesHandler(t *testing.T) { //nolint:funlen
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		mockResponse   reports.SummaryResult
		mockErr        error
		setupMocks     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock)
	}{
		{
			name:           "Success with valid date range",
			url:            "/reports/summaries?startDate=2023-01-01&endDate=2023-01-31",
			expectedStatus: http.StatusOK,
			mockResponse: reports.SummaryResult{
				PeginSummary: reports.SummaryData{
					TotalQuotesCount:    10,
					AcceptedQuotesCount: 8,
					PaidQuotesCount:     6,
					PaidQuotesAmount:    entities.NewWei(1000),
					TotalFeesCollected:  entities.NewWei(50),
					RefundedQuotesCount: 2,
					TotalPenaltyAmount:  entities.NewWei(20),
					LpEarnings:          entities.NewWei(30),
				},
				PegoutSummary: reports.SummaryData{
					TotalQuotesCount:    5,
					AcceptedQuotesCount: 4,
					PaidQuotesCount:     3,
					PaidQuotesAmount:    entities.NewWei(500),
					TotalFeesCollected:  entities.NewWei(40),
					RefundedQuotesCount: 1,
					TotalPenaltyAmount:  entities.NewWei(0),
					LpEarnings:          entities.NewWei(40),
				},
			},
			mockErr: nil,
			setupMocks: func(t *testing.T, peginRepo *mocks.PeginQuoteRepositoryMock, pegoutRepo *mocks.PegoutQuoteRepositoryMock) {
				startDate, err := time.Parse(reports.DateFormat, "2023-01-01")
				require.NoError(t, err)
				endDate, err := time.Parse(reports.DateFormat, "2023-01-31")
				require.NoError(t, err)
				endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
				peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 0, 0).
					Return([]quote.PeginQuoteWithRetained{}, 0, nil)
				pegoutRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 0, 0).
					Return([]quote.PegoutQuoteWithRetained{}, 0, nil)
			},
		},
		{
			name:           "Missing startDate parameter",
			url:            "/reports/summaries?endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   reports.SummaryResult{},
			mockErr:        nil,
			setupMocks:     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock) {},
		},
		{
			name:           "Missing endDate parameter",
			url:            "/reports/summaries?startDate=2023-01-01",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   reports.SummaryResult{},
			mockErr:        nil,
			setupMocks:     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock) {},
		},
		{
			name:           "Invalid startDate format",
			url:            "/reports/summaries?startDate=01/01/2023&endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   reports.SummaryResult{},
			mockErr:        nil,
			setupMocks:     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock) {},
		},
		{
			name:           "Invalid endDate format",
			url:            "/reports/summaries?startDate=2023-01-01&endDate=31/01/2023",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   reports.SummaryResult{},
			mockErr:        nil,
			setupMocks:     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock) {},
		},
		{
			name:           "EndDate before StartDate",
			url:            "/reports/summaries?startDate=2023-02-01&endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   reports.SummaryResult{},
			mockErr:        nil,
			setupMocks:     func(*testing.T, *mocks.PeginQuoteRepositoryMock, *mocks.PegoutQuoteRepositoryMock) {},
		},
		{
			name:           "Error in use case",
			url:            "/reports/summaries?startDate=2023-01-01&endDate=2023-01-31",
			expectedStatus: http.StatusInternalServerError,
			mockResponse:   reports.SummaryResult{},
			mockErr:        errors.New("test error"),
			setupMocks: func(t *testing.T, peginRepo *mocks.PeginQuoteRepositoryMock, pegoutRepo *mocks.PegoutQuoteRepositoryMock) {
				startDate, err := time.Parse(reports.DateFormat, "2023-01-01")
				require.NoError(t, err)
				endDate, err := time.Parse(reports.DateFormat, "2023-01-31")
				require.NoError(t, err)
				endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
				peginRepo.On("ListQuotesByDateRange", mock.Anything, startDate, endDate, 0, 0).
					Return([]quote.PeginQuoteWithRetained{}, 0, errors.New("test error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peginRepoMock := mocks.NewPeginQuoteRepositoryMock(t)
			pegoutRepoMock := mocks.NewPegoutQuoteRepositoryMock(t)
			tt.setupMocks(t, peginRepoMock, pegoutRepoMock)
			useCase := reports.NewSummariesUseCase(peginRepoMock, pegoutRepoMock, nil)
			handler := handlers.NewGetReportSummariesHandler(useCase)
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var response reports.SummaryResult
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.NotNil(t, response)
			}
			peginRepoMock.AssertExpectations(t)
			pegoutRepoMock.AssertExpectations(t)
		})
	}
}
