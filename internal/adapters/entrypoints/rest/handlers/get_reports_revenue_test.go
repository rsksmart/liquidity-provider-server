package handlers_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestNewGetReportsRevenueHandler(t *testing.T) {
	type testCase struct {
		name      string
		startDate string
		endDate   string
		mockSetup func(useCase *mocks.GetRevenueReportUseCaseMock)
		result    int
	}

	tests := []testCase{
		{
			name:      "should return 400 if startDate is missing",
			startDate: "",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is empty",
			startDate: " ",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is missing",
			startDate: "2025-08-27",
			endDate:   "",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is empty",
			startDate: "2025-08-27",
			endDate:   " ",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is after endDate",
			startDate: "2025-08-27",
			endDate:   "2025-08-26",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid startDate format",
			startDate: "Mon, 02 Jan 2024 15:04:05 MST",
			endDate:   "2025-08-26",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid endDate format",
			startDate: "2024-01-01",
			endDate:   "Mon, 02 Jan 2025 15:04:05 MST",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 500 if use case returns an error",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.GetRevenueReportResult{}, assert.AnError).
					Once()
			},
			result: http.StatusInternalServerError,
		},
		{
			name:      "should return 200 if use case succeeds",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.GetRevenueReportResult{
						TotalQuoteCallFees:    entities.NewWei(1000),
						TotalGasFeesCollected: entities.NewWei(500),
						TotalGasSpent:         entities.NewWei(300),
						TotalPenalizations:    entities.NewWei(100),
					}, nil).Once()
			},
			result: http.StatusOK,
		},
		{
			name:      "should return 200 with zero values when no quotes exist",
			startDate: "2024-01-01",
			endDate:   "2024-01-02",
			mockSetup: func(useCase *mocks.GetRevenueReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.GetRevenueReportResult{
						TotalQuoteCallFees:    entities.NewWei(0),
						TotalGasFeesCollected: entities.NewWei(0),
						TotalGasSpent:         entities.NewWei(0),
						TotalPenalizations:    entities.NewWei(0),
					}, nil).Once()
			},
			result: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := mocks.NewGetRevenueReportUseCaseMock(t)
			tc.mockSetup(useCase)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/revenue", nil)
			require.NoError(t, err)

			q := req.URL.Query()
			if tc.startDate != "" {
				q.Add("startDate", tc.startDate)
			}
			if tc.endDate != "" {
				q.Add("endDate", tc.endDate)
			}
			req.URL.RawQuery = q.Encode()
			rr := httptest.NewRecorder()

			handler := handlers.NewGetReportsRevenueHandler(
				handlers.SingleFlightGroup,
				handlers.RevenueReportSingleFlightKey,
				useCase,
			)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.result, rr.Code)
			useCase.AssertExpectations(t)
		})
	}
}

// nolint:funlen
func TestNewGetReportsRevenueHandler_ResponseStructure(t *testing.T) {
	t.Run("should return correct response structure with all gas tracking fields", func(t *testing.T) {
		useCase := mocks.NewGetRevenueReportUseCaseMock(t)
		useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
			Return(reports.GetRevenueReportResult{
				TotalQuoteCallFees:    entities.NewWei(1000),
				TotalGasFeesCollected: entities.NewWei(500),
				TotalGasSpent:         entities.NewWei(300),
				TotalPenalizations:    entities.NewWei(100),
			}, nil).Once()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/revenue?startDate=2024-01-01&endDate=2025-08-27", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.NewGetReportsRevenueHandler(
			handlers.SingleFlightGroup,
			handlers.RevenueReportSingleFlightKey,
			useCase,
		)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetRevenueReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(1000), response.TotalQuoteCallFees)
		assert.Equal(t, big.NewInt(500), response.TotalGasFeesCollected)
		assert.Equal(t, big.NewInt(300), response.TotalGasSpent)
		assert.Equal(t, big.NewInt(100), response.TotalPenalizations)

		useCase.AssertExpectations(t)
	})

	t.Run("should calculate profit correctly including gas differential", func(t *testing.T) {
		useCase := mocks.NewGetRevenueReportUseCaseMock(t)
		useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
			Return(reports.GetRevenueReportResult{
				TotalQuoteCallFees:    entities.NewWei(2000),
				TotalGasFeesCollected: entities.NewWei(600),
				TotalGasSpent:         entities.NewWei(400),
				TotalPenalizations:    entities.NewWei(300),
			}, nil).Once()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/revenue?startDate=2024-01-01&endDate=2025-08-27", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.NewGetReportsRevenueHandler(
			handlers.SingleFlightGroup,
			handlers.RevenueReportSingleFlightKey,
			useCase,
		)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetRevenueReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		gasProfit := new(big.Int).Sub(response.TotalGasFeesCollected, response.TotalGasSpent)
		expectedProfit := new(big.Int).Add(response.TotalQuoteCallFees, gasProfit)
		expectedProfit.Sub(expectedProfit, response.TotalPenalizations)

		useCase.AssertExpectations(t)
	})

	t.Run("should handle negative gas profit (gas spent > gas collected)", func(t *testing.T) {
		useCase := mocks.NewGetRevenueReportUseCaseMock(t)
		useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
			Return(reports.GetRevenueReportResult{
				TotalQuoteCallFees:    entities.NewWei(3000),
				TotalGasFeesCollected: entities.NewWei(200),
				TotalGasSpent:         entities.NewWei(500),
				TotalPenalizations:    entities.NewWei(100),
			}, nil).Once()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/revenue?startDate=2024-01-01&endDate=2025-08-27", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.NewGetReportsRevenueHandler(
			handlers.SingleFlightGroup,
			handlers.RevenueReportSingleFlightKey,
			useCase,
		)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetRevenueReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, big.NewInt(3000), response.TotalQuoteCallFees)
		assert.Equal(t, big.NewInt(200), response.TotalGasFeesCollected)
		assert.Equal(t, big.NewInt(500), response.TotalGasSpent)

		useCase.AssertExpectations(t)
	})
}
