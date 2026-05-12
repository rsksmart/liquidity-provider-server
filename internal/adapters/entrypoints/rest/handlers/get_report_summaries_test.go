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
func TestNewGetReportSummariesHandler(t *testing.T) {
	type testCase struct {
		name      string
		startDate string
		endDate   string
		mockSetup func(useCase *mocks.GetSummariesReportUseCaseMock)
		result    int
	}

	tests := []testCase{
		{
			name:      "should return 400 if startDate is missing",
			startDate: "",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is empty",
			startDate: " ",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is missing",
			startDate: "2025-08-27",
			endDate:   "",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is empty",
			startDate: "2025-08-27",
			endDate:   " ",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is after endDate",
			startDate: "2025-08-27",
			endDate:   "2025-08-26",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid startDate format",
			startDate: "Mon, 02 Jan 2024 15:04:05 MST",
			endDate:   "2025-08-26",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid endDate format",
			startDate: "2024-01-01",
			endDate:   "Mon, 02 Jan 2025 15:04:05 MST",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 500 if use case returns an error",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.SummaryResult{}, assert.AnError).
					Once()
			},
			result: http.StatusInternalServerError,
		},
		{
			name:      "should return 200 if use case succeeds with comprehensive data",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.SummaryResult{
						PeginSummary: reports.SummaryData{
							TotalQuotesCount:          100,
							AcceptedQuotesCount:       80,
							TotalAcceptedQuotesAmount: entities.NewWei(50000),
							PaidQuotesCount:           70,
							PaidQuotesAmount:          entities.NewWei(45000),
							RefundedQuotesCount:       10,
							TotalRefundedQuotesAmount: entities.NewWei(5000),
							PenalizationsCount:        5,
							TotalPenalizationsAmount:  entities.NewWei(500),
						},
						PegoutSummary: reports.SummaryData{
							TotalQuotesCount:          90,
							AcceptedQuotesCount:       75,
							TotalAcceptedQuotesAmount: entities.NewWei(45000),
							PaidQuotesCount:           65,
							PaidQuotesAmount:          entities.NewWei(40000),
							RefundedQuotesCount:       8,
							TotalRefundedQuotesAmount: entities.NewWei(4000),
							PenalizationsCount:        3,
							TotalPenalizationsAmount:  entities.NewWei(300),
						},
					}, nil).Once()
			},
			result: http.StatusOK,
		},
		{
			name:      "should return 200 with zero values when no quotes exist",
			startDate: "2024-01-01",
			endDate:   "2024-01-02",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.SummaryResult{
						PeginSummary: reports.SummaryData{
							TotalQuotesCount:          0,
							AcceptedQuotesCount:       0,
							TotalAcceptedQuotesAmount: entities.NewWei(0),
							PaidQuotesCount:           0,
							PaidQuotesAmount:          entities.NewWei(0),
							RefundedQuotesCount:       0,
							TotalRefundedQuotesAmount: entities.NewWei(0),
							PenalizationsCount:        0,
							TotalPenalizationsAmount:  entities.NewWei(0),
						},
						PegoutSummary: reports.SummaryData{
							TotalQuotesCount:          0,
							AcceptedQuotesCount:       0,
							TotalAcceptedQuotesAmount: entities.NewWei(0),
							PaidQuotesCount:           0,
							PaidQuotesAmount:          entities.NewWei(0),
							RefundedQuotesCount:       0,
							TotalRefundedQuotesAmount: entities.NewWei(0),
							PenalizationsCount:        0,
							TotalPenalizationsAmount:  entities.NewWei(0),
						},
					}, nil).Once()
			},
			result: http.StatusOK,
		},
		{
			name:      "should support ISO 8601 date format",
			startDate: "2024-01-01T00:00:00Z",
			endDate:   "2024-01-31T23:59:59Z",
			mockSetup: func(useCase *mocks.GetSummariesReportUseCaseMock) {
				useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
					Return(reports.SummaryResult{
						PeginSummary: reports.SummaryData{
							TotalQuotesCount:          10,
							AcceptedQuotesCount:       8,
							TotalAcceptedQuotesAmount: entities.NewWei(5000),
							PaidQuotesCount:           7,
							PaidQuotesAmount:          entities.NewWei(4500),
							RefundedQuotesCount:       1,
							TotalRefundedQuotesAmount: entities.NewWei(500),
							PenalizationsCount:        1,
							TotalPenalizationsAmount:  entities.NewWei(50),
						},
						PegoutSummary: reports.SummaryData{
							TotalQuotesCount:          5,
							AcceptedQuotesCount:       4,
							TotalAcceptedQuotesAmount: entities.NewWei(2500),
							PaidQuotesCount:           3,
							PaidQuotesAmount:          entities.NewWei(2000),
							RefundedQuotesCount:       1,
							TotalRefundedQuotesAmount: entities.NewWei(500),
							PenalizationsCount:        0,
							TotalPenalizationsAmount:  entities.NewWei(0),
						},
					}, nil).Once()
			},
			result: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := mocks.NewGetSummariesReportUseCaseMock(t)
			tc.mockSetup(useCase)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/summaries", nil)
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

			handler := handlers.NewGetReportSummariesHandler(
				handlers.SingleFlightGroup,
				handlers.SummariesReportSingleFlightKey,
				useCase,
			)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.result, rr.Code)
			useCase.AssertExpectations(t)
		})
	}
}

// nolint:funlen
func TestNewGetReportSummariesHandler_ResponseStructure(t *testing.T) {
	t.Run("should return correct response structure with all summary fields", func(t *testing.T) {
		useCase := mocks.NewGetSummariesReportUseCaseMock(t)
		useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
			Return(reports.SummaryResult{
				PeginSummary: reports.SummaryData{
					TotalQuotesCount:          100,
					AcceptedQuotesCount:       80,
					TotalAcceptedQuotesAmount: entities.NewWei(50000),
					PaidQuotesCount:           70,
					PaidQuotesAmount:          entities.NewWei(45000),
					RefundedQuotesCount:       10,
					TotalRefundedQuotesAmount: entities.NewWei(5000),
					PenalizationsCount:        5,
					TotalPenalizationsAmount:  entities.NewWei(500),
				},
				PegoutSummary: reports.SummaryData{
					TotalQuotesCount:          90,
					AcceptedQuotesCount:       75,
					TotalAcceptedQuotesAmount: entities.NewWei(45000),
					PaidQuotesCount:           65,
					PaidQuotesAmount:          entities.NewWei(40000),
					RefundedQuotesCount:       8,
					TotalRefundedQuotesAmount: entities.NewWei(4000),
					PenalizationsCount:        3,
					TotalPenalizationsAmount:  entities.NewWei(300),
				},
			}, nil).Once()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/summaries?startDate=2024-01-01&endDate=2025-08-27", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.NewGetReportSummariesHandler(
			handlers.SingleFlightGroup,
			handlers.SummariesReportSingleFlightKey,
			useCase,
		)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.SummaryResultDTO
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		// Assert Pegin Summary
		assert.Equal(t, int64(100), response.PeginSummary.TotalQuotesCount)
		assert.Equal(t, int64(80), response.PeginSummary.AcceptedQuotesCount)
		assert.Equal(t, big.NewInt(50000), response.PeginSummary.TotalAcceptedQuotesAmount)
		assert.Equal(t, int64(70), response.PeginSummary.PaidQuotesCount)
		assert.Equal(t, big.NewInt(45000), response.PeginSummary.PaidQuotesAmount)
		assert.Equal(t, int64(10), response.PeginSummary.RefundedQuotesCount)
		assert.Equal(t, big.NewInt(5000), response.PeginSummary.TotalRefundedQuotesAmount)
		assert.Equal(t, int64(5), response.PeginSummary.PenalizationsCount)
		assert.Equal(t, big.NewInt(500), response.PeginSummary.TotalPenalizationsAmount)

		// Assert Pegout Summary
		assert.Equal(t, int64(90), response.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, int64(75), response.PegoutSummary.AcceptedQuotesCount)
		assert.Equal(t, big.NewInt(45000), response.PegoutSummary.TotalAcceptedQuotesAmount)
		assert.Equal(t, int64(65), response.PegoutSummary.PaidQuotesCount)
		assert.Equal(t, big.NewInt(40000), response.PegoutSummary.PaidQuotesAmount)
		assert.Equal(t, int64(8), response.PegoutSummary.RefundedQuotesCount)
		assert.Equal(t, big.NewInt(4000), response.PegoutSummary.TotalRefundedQuotesAmount)
		assert.Equal(t, int64(3), response.PegoutSummary.PenalizationsCount)
		assert.Equal(t, big.NewInt(300), response.PegoutSummary.TotalPenalizationsAmount)

		useCase.AssertExpectations(t)
	})

	t.Run("should handle scenario with only pegin quotes", func(t *testing.T) {
		useCase := mocks.NewGetSummariesReportUseCaseMock(t)
		useCase.On("Run", mock.Anything, mock.Anything, mock.Anything).
			Return(reports.SummaryResult{
				PeginSummary: reports.SummaryData{
					TotalQuotesCount:          50,
					AcceptedQuotesCount:       40,
					TotalAcceptedQuotesAmount: entities.NewWei(25000),
					PaidQuotesCount:           35,
					PaidQuotesAmount:          entities.NewWei(22000),
					RefundedQuotesCount:       5,
					TotalRefundedQuotesAmount: entities.NewWei(3000),
					PenalizationsCount:        2,
					TotalPenalizationsAmount:  entities.NewWei(200),
				},
				PegoutSummary: reports.SummaryData{
					TotalQuotesCount:          0,
					AcceptedQuotesCount:       0,
					TotalAcceptedQuotesAmount: entities.NewWei(0),
					PaidQuotesCount:           0,
					PaidQuotesAmount:          entities.NewWei(0),
					RefundedQuotesCount:       0,
					TotalRefundedQuotesAmount: entities.NewWei(0),
					PenalizationsCount:        0,
					TotalPenalizationsAmount:  entities.NewWei(0),
				},
			}, nil).Once()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/summaries?startDate=2024-01-01&endDate=2025-08-27", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.NewGetReportSummariesHandler(
			handlers.SingleFlightGroup,
			handlers.SummariesReportSingleFlightKey,
			useCase,
		)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.SummaryResultDTO
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		// Pegin should have data
		assert.Equal(t, int64(50), response.PeginSummary.TotalQuotesCount)
		assert.Equal(t, big.NewInt(25000), response.PeginSummary.TotalAcceptedQuotesAmount)

		// Pegout should be zero
		assert.Equal(t, int64(0), response.PegoutSummary.TotalQuotesCount)
		assert.Equal(t, big.NewInt(0), response.PegoutSummary.TotalAcceptedQuotesAmount)

		useCase.AssertExpectations(t)
	})
}
