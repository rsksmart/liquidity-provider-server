package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSummariesUseCase struct {
	mock.Mock
}

func (m *MockSummariesUseCase) Run(ctx context.Context, startDate, endDate time.Time) (liquidity_provider.SummariesResponse, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return liquidity_provider.SummariesResponse{}, args.Error(1)
	}
	response, ok := args.Get(0).(liquidity_provider.SummariesResponse)
	if !ok {
		return liquidity_provider.SummariesResponse{}, errors.New("invalid response type")
	}
	return response, args.Error(1)
}

func getReportSummariesHandlerForTest(useCase *MockSummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		startDate, endDate, valid := rest.ValidateDateRange(w, req, liquidity_provider.DateFormat)
		if !valid {
			return
		}
		response, err := useCase.Run(req.Context(), startDate, endDate)
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("An error occurred while processing your request", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}

func TestGetReportSummariesHandler(t *testing.T) { //nolint:funlen
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
					TotalAcceptedQuotesCount:  10,
					ConfirmedQuotesCount:      8,
					TotalQuotedAmount:         "1000",
					TotalAcceptedQuotedAmount: "800",
					TotalFeesCollected:        "50",
					RefundedQuotesCount:       2,
					TotalPenaltyAmount:        "20",
					LpEarnings:                "30",
				},
				PegoutSummary: liquidity_provider.SummaryData{
					TotalAcceptedQuotesCount:  5,
					ConfirmedQuotesCount:      4,
					TotalQuotedAmount:         "500",
					TotalAcceptedQuotedAmount: "400",
					TotalFeesCollected:        "40",
					RefundedQuotesCount:       1,
					TotalPenaltyAmount:        "0",
					LpEarnings:                "20",
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
			name:           "EndDate before StartDate",
			url:            "/report/summaries?startDate=2023-02-01&endDate=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        nil,
		},
		{
			name:           "Error in use case",
			url:            "/report/summaries?startDate=2023-01-01&endDate=2023-01-31",
			expectedStatus: http.StatusInternalServerError,
			mockResponse:   liquidity_provider.SummariesResponse{},
			mockErr:        errors.New("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(MockSummariesUseCase)
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusInternalServerError {
				startDate, err := time.Parse(liquidity_provider.DateFormat, "2023-01-01")
				require.NoError(t, err)
				endDate, err := time.Parse(liquidity_provider.DateFormat, "2023-01-31")
				require.NoError(t, err)
				endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
				mockUseCase.On("Run", mock.Anything, startDate, endDate).Return(tt.mockResponse, tt.mockErr)
			}
			handler := getReportSummariesHandlerForTest(mockUseCase)
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.url, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var response liquidity_provider.SummariesResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tt.mockResponse, response)
			}
			mockUseCase.AssertExpectations(t)
		})
	}
}
