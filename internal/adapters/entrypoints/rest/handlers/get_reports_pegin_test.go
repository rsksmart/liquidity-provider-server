package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/singleflight"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// nolint:funlen
func TestNewGetReportsPeginHandler(t *testing.T) {
	type testCase struct {
		name      string
		startDate string
		endDate   string
		mockSetup func(mock *mocks.GetPeginReportUseCaseMock)
		result    int
	}

	tests := []testCase{
		{
			name:      "should return 400 if startDate is missing",
			startDate: "",
			endDate:   "2025-08-27",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is empty",
			startDate: " ",
			endDate:   "2025-08-27",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is missing",
			startDate: "2025-08-27",
			endDate:   "",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if endDate is empty",
			startDate: "2025-08-27",
			endDate:   " ",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 if startDate is after endDate",
			startDate: "2025-08-27",
			endDate:   "2025-08-26",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid startDate format",
			startDate: "Mon, 02 Jan 2024 15:04:05 MST",
			endDate:   "2025-08-26",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 400 for invalid endDate format",
			startDate: "2024-01-01",
			endDate:   "Mon, 02 Jan 2025 15:04:05 MST",
			mockSetup: func(mock *mocks.GetPeginReportUseCaseMock) {},
			result:    http.StatusBadRequest,
		},
		{
			name:      "should return 500 if use case returns an error",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetPeginReportUseCaseMock) {
				useCase.EXPECT().
					Run(mock.Anything, mock.Anything, mock.Anything).
					Return(reports.GetPeginReportResult{}, assert.AnError).
					Once()
			},
			result: http.StatusInternalServerError,
		},
		{
			name:      "should return 200 if use case succeeds",
			startDate: "2024-01-01",
			endDate:   "2025-08-27",
			mockSetup: func(useCase *mocks.GetPeginReportUseCaseMock) {
				useCase.EXPECT().
					Run(mock.Anything, mock.Anything, mock.Anything).
					Return(reports.GetPeginReportResult{
						NumberOfQuotes:     1,
						MinimumQuoteValue:  entities.NewWei(1),
						MaximumQuoteValue:  entities.NewWei(3),
						AverageQuoteValue:  entities.NewWei(2),
						TotalFeesCollected: entities.NewWei(3),
						AverageFeePerQuote: entities.NewWei(1),
					}, nil).Once()
			},
			result: http.StatusOK,
		},
	}

	const groupKey = "test-key"
	handlerSingleFlight := new(singleflight.Group)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := &mocks.GetPeginReportUseCaseMock{}
			tc.mockSetup(useCase)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/pegin", nil)
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

			handler := handlers.NewGetReportsPeginHandler(handlerSingleFlight, groupKey, useCase)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.result, rr.Code)
			useCase.AssertExpectations(t)
		})
	}
}

// nolint:funlen
func TestNewGetReportsPeginHandler_Concurrency(t *testing.T) {
	t.Run("should only call use case once for concurrent requests", func(t *testing.T) {
		useCase := &mocks.GetPeginReportUseCaseMock{}
		const groupKey = "test-key"
		handlerSingleFlight := new(singleflight.Group)
		result := reports.GetPeginReportResult{
			NumberOfQuotes:     1,
			MinimumQuoteValue:  entities.NewWei(2),
			MaximumQuoteValue:  entities.NewWei(3),
			AverageQuoteValue:  entities.NewWei(4),
			TotalFeesCollected: entities.NewWei(5),
			AverageFeePerQuote: entities.NewWei(6),
		}
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, start time.Time, end time.Time) (reports.GetPeginReportResult, error) {
				time.Sleep(1 * time.Second)
				return result, nil
			}).Once() // The call to Once is very important, as we ensure the mock is not called by more than one request

		handler := handlers.NewGetReportsPeginHandler(handlerSingleFlight, groupKey, useCase)

		const testSize = 30
		waitGroup := new(sync.WaitGroup)
		waitGroup.Add(testSize)
		results := make(map[int]pkg.GetPeginReportResponse, testSize)
		mutex := new(sync.Mutex)
		for i := 0; i < testSize; i++ {
			go func(requestId int, wg *sync.WaitGroup, resultMap map[int]pkg.GetPeginReportResponse, mapMutex *sync.Mutex) {
				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/pegin?startDate=2024-08-27&endDate=2025-08-27", nil)
				assert.NoError(t, err)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				assert.Equal(t, http.StatusOK, rr.Code)
				var body pkg.GetPeginReportResponse
				assert.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
				mapMutex.Lock()
				resultMap[requestId] = body
				mapMutex.Unlock()
				wg.Done()
			}(i, waitGroup, results, mutex)
		}
		waitGroup.Wait()

		expected := pkg.GetPeginReportResponse{
			NumberOfQuotes:     1,
			MinimumQuoteValue:  big.NewInt(2),
			MaximumQuoteValue:  big.NewInt(3),
			AverageQuoteValue:  big.NewInt(4),
			TotalFeesCollected: big.NewInt(5),
			AverageFeePerQuote: big.NewInt(6),
		}
		for i := 0; i < testSize; i++ {
			assert.Equal(t, expected, results[i], "Result of request #%d not matching. Expected %v got %v", i, expected, results[i])
		}
	})
	t.Run("should be capable of handling multiple groups", func(t *testing.T) {
		useCase := &mocks.GetPeginReportUseCaseMock{}
		const groupKey = "test-key"
		handlerSingleFlight := new(singleflight.Group)
		resultGroupOne := reports.GetPeginReportResult{
			NumberOfQuotes:     1,
			MinimumQuoteValue:  entities.NewWei(2),
			MaximumQuoteValue:  entities.NewWei(3),
			AverageQuoteValue:  entities.NewWei(4),
			TotalFeesCollected: entities.NewWei(5),
			AverageFeePerQuote: entities.NewWei(6),
		}
		resultGroupTwo := reports.GetPeginReportResult{
			NumberOfQuotes:     7,
			MinimumQuoteValue:  entities.NewWei(8),
			MaximumQuoteValue:  entities.NewWei(9),
			AverageQuoteValue:  entities.NewWei(10),
			TotalFeesCollected: entities.NewWei(11),
			AverageFeePerQuote: entities.NewWei(12),
		}
		const (
			groupOneStart = "2025-07-27"
			groupOneEnd   = "2025-08-27"
			groupTwoStart = "2025-08-28"
			groupTwoEnd   = "2025-08-29"
		)

		dateMatcher := func(target time.Time) any {
			return mock.MatchedBy(func(date time.Time) bool {
				return target.Day() == date.Day() && target.Month() == date.Month() && target.Year() == date.Year()
			})
		}
		useCase.EXPECT().Run(mock.Anything, dateMatcher(test.MustParseDate(groupOneStart)), dateMatcher(test.MustParseDate(groupOneEnd))).
			RunAndReturn(func(ctx context.Context, start time.Time, end time.Time) (reports.GetPeginReportResult, error) {
				time.Sleep(2 * time.Second)
				return resultGroupOne, nil
			}).Once()
		useCase.EXPECT().Run(mock.Anything, dateMatcher(test.MustParseDate(groupTwoStart)), dateMatcher(test.MustParseDate(groupTwoEnd))).
			RunAndReturn(func(ctx context.Context, start time.Time, end time.Time) (reports.GetPeginReportResult, error) {
				time.Sleep(1 * time.Second)
				return resultGroupTwo, nil
			}).Once()

		handler := handlers.NewGetReportsPeginHandler(handlerSingleFlight, groupKey, useCase)

		const testSize = 30
		waitGroup := new(sync.WaitGroup)
		waitGroup.Add(testSize)
		results := make(map[int]pkg.GetPeginReportResponse, testSize)
		mutex := new(sync.Mutex)
		for i := 0; i < testSize; i++ {
			go func(requestId int, wg *sync.WaitGroup, resultMap map[int]pkg.GetPeginReportResponse, mapMutex *sync.Mutex) {
				var err error
				var req *http.Request
				if i%2 == 0 {
					req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/reports/pegin?startDate=%s&endDate=%s", groupOneStart, groupOneEnd), nil)
				} else {
					req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("/reports/pegin?startDate=%s&endDate=%s", groupTwoStart, groupTwoEnd), nil)
				}
				assert.NoError(t, err)
				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)
				assert.Equal(t, http.StatusOK, rr.Code)
				var body pkg.GetPeginReportResponse
				assert.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
				mapMutex.Lock()
				resultMap[requestId] = body
				mapMutex.Unlock()
				wg.Done()
			}(i, waitGroup, results, mutex)
		}
		waitGroup.Wait()

		expectedGroupOne := pkg.GetPeginReportResponse{
			NumberOfQuotes:     1,
			MinimumQuoteValue:  big.NewInt(2),
			MaximumQuoteValue:  big.NewInt(3),
			AverageQuoteValue:  big.NewInt(4),
			TotalFeesCollected: big.NewInt(5),
			AverageFeePerQuote: big.NewInt(6),
		}
		expectedGroupTwo := pkg.GetPeginReportResponse{
			NumberOfQuotes:     7,
			MinimumQuoteValue:  big.NewInt(8),
			MaximumQuoteValue:  big.NewInt(9),
			AverageQuoteValue:  big.NewInt(10),
			TotalFeesCollected: big.NewInt(11),
			AverageFeePerQuote: big.NewInt(12),
		}
		useCase.AssertExpectations(t)
		for i := 0; i < testSize; i++ {
			if i%2 == 0 {
				assert.Equal(t, expectedGroupOne, results[i], "Result of request #%d not matching. Expected %v got %v", i, expectedGroupOne, results[i])
			} else {
				assert.Equal(t, expectedGroupTwo, results[i], "Result of request #%d not matching. Expected %v got %v", i, expectedGroupTwo, results[i])
			}
		}
	})
}
