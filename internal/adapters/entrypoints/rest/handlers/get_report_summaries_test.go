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

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
)

type MockPeginQuoteRepository struct {
	mock.Mock
}

func (m *MockPeginQuoteRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PeginQuote, []quote.RetainedPeginQuote, error) {
	args := m.Called(ctx, startDate, endDate)
	err := args.Error(2)
	if err != nil {
		return nil, nil, err
	}
	
	quotes, ok := args.Get(0).([]quote.PeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, nil, errors.New("invalid pegin quotes type")
	}
	
	retainedQuotes, ok := args.Get(1).([]quote.RetainedPeginQuote)
	if !ok && args.Get(1) != nil {
		return nil, nil, errors.New("invalid retained pegin quotes type")
	}
	
	return quotes, retainedQuotes, nil
}

func (m *MockPeginQuoteRepository) GetQuote(ctx context.Context, quoteHash string) (*quote.PeginQuote, error) {
	args := m.Called(ctx, quoteHash)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	
	quote, ok := args.Get(0).(*quote.PeginQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid pegin quote type")
	}
	
	return quote, nil
}

type MockPegoutQuoteRepository struct {
	mock.Mock
}

func (m *MockPegoutQuoteRepository) ListQuotesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]quote.PegoutQuote, []quote.RetainedPegoutQuote, error) {
	args := m.Called(ctx, startDate, endDate)
	err := args.Error(2)
	if err != nil {
		return nil, nil, err
	}
	
	quotes, ok := args.Get(0).([]quote.PegoutQuote)
	if !ok && args.Get(0) != nil {
		return nil, nil, errors.New("invalid pegout quotes type")
	}
	
	retainedQuotes, ok := args.Get(1).([]quote.RetainedPegoutQuote)
	if !ok && args.Get(1) != nil {
		return nil, nil, errors.New("invalid retained pegout quotes type")
	}
	
	return quotes, retainedQuotes, nil
}

func (m *MockPegoutQuoteRepository) GetQuote(ctx context.Context, quoteHash string) (*quote.PegoutQuote, error) {
	args := m.Called(ctx, quoteHash)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	
	quote, ok := args.Get(0).(*quote.PegoutQuote)
	if !ok && args.Get(0) != nil {
		return nil, errors.New("invalid pegout quote type")
	}
	
	return quote, nil
}

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
		return liquidity_provider.SummariesResponse{}, errors.New("invalid summaries response type")
	}
	
	return response, args.Error(1)
}

func validateDateParametersForTest(w http.ResponseWriter, req *http.Request) (startDate time.Time, endDate time.Time, valid bool) {
	start := req.URL.Query().Get("startDate")
	end := req.URL.Query().Get("endDate")
	
	if start == "" || end == "" {
		missing := []string{}
		if start == "" {
			missing = append(missing, "startDate")
		}
		if end == "" {
			missing = append(missing, "endDate")
		}
		jsonErr := rest.NewErrorResponseWithDetails("missing required parameters", map[string]any{
			"missing": missing,
		}, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	var err error
	startDate, err = time.Parse("2006-01-02", start)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	endDate, err = time.Parse("2006-01-02", end)
	if err != nil {
		jsonErr := rest.NewErrorResponseWithDetails("invalid date format", rest.DetailsFromError(err), true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	if endDate.Before(startDate) {
		details := map[string]any{
			"startDate": startDate.Format("2006-01-02"),
			"endDate":   endDate.Format("2006-01-02"),
		}
		jsonErr := rest.NewErrorResponseWithDetails("invalid date range", details, true)
		rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
		return time.Time{}, time.Time{}, false
	}

	return startDate, endDate, true
}

func getReportSummariesHandlerForTest(useCase *MockSummariesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		startDate, endDate, valid := validateDateParametersForTest(w, req)
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
			req, err := http.NewRequestWithContext(context.Background(), "GET", tt.url, nil)
			require.NoError(t, err)
			
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusInternalServerError {
				q := req.URL.Query()
				startDateStr := q.Get("startDate")
				endDateStr := q.Get("endDate")
				
				startDate, err := time.Parse("2006-01-02", startDateStr)
				if err == nil {
					endDate, err := time.Parse("2006-01-02", endDateStr)
					if err == nil {
						endDateWithTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
						mockUseCase.On("Run", mock.Anything, startDate, endDateWithTime).Return(tt.mockResponse, tt.mockErr)
					}
				}
			}

			rr := httptest.NewRecorder()			
			handler := getReportSummariesHandlerForTest(mockUseCase)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code, "Expected status %d but got %d", tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var response liquidity_provider.SummariesResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.mockResponse, response)
			} else if tt.expectedStatus == http.StatusBadRequest || tt.expectedStatus == http.StatusInternalServerError {
				var errorResponse rest.ErrorResponse
				err := json.NewDecoder(rr.Body).Decode(&errorResponse)
				require.NoError(t, err)
				assert.NotEmpty(t, errorResponse.Message)
			}
			
			mockUseCase.AssertExpectations(t)
		})
	}
}
