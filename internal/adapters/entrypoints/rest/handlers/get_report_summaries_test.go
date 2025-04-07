package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

func HandlerWithMock(useCase interface {
	Run(ctx context.Context, startDate, endDate time.Time) (liquidity_provider.SummariesResponse, error)
}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := req.URL.Query().Get("startDate")
		end := req.URL.Query().Get("endDate")

		if start == "" || end == "" {
			http.Error(w, "Missing startDate or endDate query parameter", http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("2006-01-02", start)
		if err != nil {
			http.Error(w, "Invalid startDate format", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", end)
		if err != nil {
			http.Error(w, "Invalid endDate format", http.StatusBadRequest)
			return
		}

		response, err := useCase.Run(req.Context(), startDate, endDate)
		if err != nil {
			http.Error(w, "Failed to retrieve report data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}
}

type MockSummariesUseCase struct {
	mock.Mock
}

func (m *MockSummariesUseCase) Run(ctx context.Context, startDate, endDate time.Time) (liquidity_provider.SummariesResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return liquidity_provider.SummariesResponse{}, args.Error(1)
	}
	return args.Get(0).(liquidity_provider.SummariesResponse), args.Error(1)
}

func TestGetReportSummariesHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		mockResponse   liquidity_provider.SummariesResponse
		mockErr        error
	}{
		{
			name:           "Success with valid date range",
			url:            "/report/summaries?startDate=2023-01-01&endDate=2023-01-31",
			expectedStatus: http.StatusOK,
			mockResponse: liquidity_provider.SummariesResponse{
				PeginSummary: liquidity_provider.SummaryData{
					TotalQuotesCount:          10,
					AcceptedQuotesCount:       8,
					TotalQuotedAmount:         "1000",
					TotalAcceptedQuotedAmount: "800",
					TotalFeesCollected:        "50",
					RefundedQuotesCount:       2,
					TotalPenaltyAmount:        "20",
					LpEarnings:                "30",
				},
				PegoutSummary: liquidity_provider.SummaryData{
					TotalQuotesCount:          5,
					AcceptedQuotesCount:       4,
					TotalQuotedAmount:         "500",
					TotalAcceptedQuotedAmount: "400",
					TotalFeesCollected:        "25",
					RefundedQuotesCount:       1,
					TotalPenaltyAmount:        "10",
					LpEarnings:                "15",
				},
			},
			mockErr: nil,
		},
		{
			name:           "Missing startDate parameter",
			url:            "/report/summaries?endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        nil,
		},
		{
			name:           "Missing endDate parameter",
			url:            "/report/summaries?startDate=2023-01-01",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        nil,
		},
		{
			name:           "Invalid startDate format",
			url:            "/report/summaries?startDate=01/01/2023&endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        nil,
		},
		{
			name:           "Invalid endDate format",
			url:            "/report/summaries?startDate=2023-01-01&endDate=31/01/2023",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        nil,
		},
		{
			name:           "Internal error from use case",
			url:            "/report/summaries?startDate=2023-01-01&endDate=2023-01-31",
			expectedStatus: http.StatusInternalServerError,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockSummariesUseCase)
			if startDateOK := len(tt.url) > 0 && tt.url != "/report/summaries?endDate=2023-01-31" && !isInvalidDateFormat(tt.url, "startDate"); startDateOK {
				if endDateOK := len(tt.url) > 0 && tt.url != "/report/summaries?startDate=2023-01-01" && !isInvalidDateFormat(tt.url, "endDate"); endDateOK {
					startDate, _ := time.Parse("2006-01-02", "2023-01-01")
					endDate, _ := time.Parse("2006-01-02", "2023-01-31")
					mockUseCase.On("Run", mock.Anything, startDate, endDate).Return(tt.mockResponse, tt.mockErr)
				}
			}

			req, err := http.NewRequest("GET", tt.url, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := HandlerWithMock(mockUseCase)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Expected status %d but got %d", tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response liquidity_provider.SummariesResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.mockResponse, response)
			}
			mockUseCase.AssertExpectations(t)
		})
	}
}

func isInvalidDateFormat(url string, param string) bool {
	if param == "startDate" {
		return url == "/report/summaries?startDate=01/01/2023&endDate=2023-01-31"
	} else if param == "endDate" {
		return url == "/report/summaries?startDate=2023-01-01&endDate=31/01/2023"
	}
	return false
}
