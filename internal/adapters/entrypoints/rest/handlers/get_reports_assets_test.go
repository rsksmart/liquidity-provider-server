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

//nolint:funlen
func TestNewGetReportsAssetsHandler_Success(t *testing.T) {
	t.Run("should return 200 with correct structure on success", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)

		expectedResult := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(67500000), // 0.675 BTC
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(50000000), // 0.5 BTC
					Federation: entities.NewWei(5000000),  // 0.05 BTC
					RskWallet:  entities.NewWei(6500000),  // 0.065 BTC
					Lbc:        entities.NewWei(6000000),  // 0.06 BTC
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(4500000),  // 0.045 BTC
					WaitingForRefund: entities.NewWei(11500000), // 0.115 BTC
					Available:        entities.NewWei(51500000), // 0.515 BTC
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewBigWei(new(big.Int).SetUint64(17000000000000000000)), // 17 RBTC
				Location: reports.RbtcAssetLocation{
					RskWallet:  entities.NewBigWei(new(big.Int).SetUint64(10000000000000000000)), // 10 RBTC
					Lbc:        entities.NewBigWei(new(big.Int).SetUint64(5000000000000000000)),  // 5 RBTC
					Federation: entities.NewBigWei(new(big.Int).SetUint64(2000000000000000000)),  // 2 RBTC
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: entities.NewBigWei(new(big.Int).SetUint64(3000000000000000000)),  // 3 RBTC
					WaitingForRefund: entities.NewBigWei(new(big.Int).SetUint64(2000000000000000000)),  // 2 RBTC
					Available:        entities.NewBigWei(new(big.Int).SetUint64(12000000000000000000)), // 12 RBTC
				},
			},
		}

		useCase.EXPECT().Run(mock.Anything).Return(expectedResult, nil).Once()

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/assets", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetAssetsReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, expectedResult.BtcAssetReport.Total.AsBigInt(), response.BtcAssetReport.Total)
		assert.Equal(t, expectedResult.BtcAssetReport.Location.BtcWallet.AsBigInt(), response.BtcAssetReport.Location.BtcWallet)
		assert.Equal(t, expectedResult.BtcAssetReport.Location.Federation.AsBigInt(), response.BtcAssetReport.Location.Federation)
		assert.Equal(t, expectedResult.BtcAssetReport.Location.RskWallet.AsBigInt(), response.BtcAssetReport.Location.RskWallet)
		assert.Equal(t, expectedResult.BtcAssetReport.Location.Lbc.AsBigInt(), response.BtcAssetReport.Location.Lbc)
		assert.Equal(t, expectedResult.BtcAssetReport.Allocation.ReservedForUsers.AsBigInt(), response.BtcAssetReport.Allocation.ReservedForUsers)
		assert.Equal(t, expectedResult.BtcAssetReport.Allocation.WaitingForRefund.AsBigInt(), response.BtcAssetReport.Allocation.WaitingForRefund)
		assert.Equal(t, expectedResult.BtcAssetReport.Allocation.Available.AsBigInt(), response.BtcAssetReport.Allocation.Available)

		assert.Equal(t, expectedResult.RbtcAssetReport.Total.AsBigInt(), response.RbtcAssetReport.Total)
		assert.Equal(t, expectedResult.RbtcAssetReport.Location.RskWallet.AsBigInt(), response.RbtcAssetReport.Location.RskWallet)
		assert.Equal(t, expectedResult.RbtcAssetReport.Location.Lbc.AsBigInt(), response.RbtcAssetReport.Location.Lbc)
		assert.Equal(t, expectedResult.RbtcAssetReport.Location.Federation.AsBigInt(), response.RbtcAssetReport.Location.Federation)
		assert.Equal(t, expectedResult.RbtcAssetReport.Allocation.ReservedForUsers.AsBigInt(), response.RbtcAssetReport.Allocation.ReservedForUsers)
		assert.Equal(t, expectedResult.RbtcAssetReport.Allocation.WaitingForRefund.AsBigInt(), response.RbtcAssetReport.Allocation.WaitingForRefund)
		assert.Equal(t, expectedResult.RbtcAssetReport.Allocation.Available.AsBigInt(), response.RbtcAssetReport.Allocation.Available)

		useCase.AssertExpectations(t)
	})

	t.Run("should handle zero balances correctly", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)

		expectedResult := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(0),
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(0),
					Federation: entities.NewWei(0),
					RskWallet:  entities.NewWei(0),
					Lbc:        entities.NewWei(0),
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(0),
					WaitingForRefund: entities.NewWei(0),
					Available:        entities.NewWei(0),
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: entities.NewWei(0),
				Location: reports.RbtcAssetLocation{
					RskWallet:  entities.NewWei(0),
					Lbc:        entities.NewWei(0),
					Federation: entities.NewWei(0),
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: entities.NewWei(0),
					WaitingForRefund: entities.NewWei(0),
					Available:        entities.NewWei(0),
				},
			},
		}

		useCase.EXPECT().Run(mock.Anything).Return(expectedResult, nil).Once()

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/assets", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetAssetsReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, int64(0), response.BtcAssetReport.Total.Int64())
		assert.Equal(t, int64(0), response.RbtcAssetReport.Total.Int64())

		useCase.AssertExpectations(t)
	})

	t.Run("should handle large balances correctly", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)

		expectedResult := reports.GetAssetsReportResult{
			BtcAssetReport: reports.BtcAssetReport{
				Total: entities.NewWei(2100000000000000), // 21 million BTC in satoshis
				Location: reports.BtcAssetLocation{
					BtcWallet:  entities.NewWei(1000000000000000),
					Federation: entities.NewWei(500000000000000),
					RskWallet:  entities.NewWei(300000000000000),
					Lbc:        entities.NewWei(300000000000000),
				},
				Allocation: reports.BtcAssetAllocation{
					ReservedForUsers: entities.NewWei(100000000000000),
					WaitingForRefund: entities.NewWei(800000000000000),
					Available:        entities.NewWei(1200000000000000),
				},
			},
			RbtcAssetReport: reports.RbtcAssetReport{
				Total: mustParseBigWei("21000000000000000000000000"), // 21 million RBTC in wei
				Location: reports.RbtcAssetLocation{
					RskWallet:  mustParseBigWei("10000000000000000000000000"),
					Lbc:        mustParseBigWei("8000000000000000000000000"),
					Federation: mustParseBigWei("3000000000000000000000000"),
				},
				Allocation: reports.RbtcAssetAllocation{
					ReservedForUsers: mustParseBigWei("2000000000000000000000000"),
					WaitingForRefund: mustParseBigWei("3000000000000000000000000"),
					Available:        mustParseBigWei("16000000000000000000000000"),
				},
			},
		}

		useCase.EXPECT().Run(mock.Anything).Return(expectedResult, nil).Once()

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/assets", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response pkg.GetAssetsReportResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, expectedResult.BtcAssetReport.Total.AsBigInt(), response.BtcAssetReport.Total)
		assert.Equal(t, expectedResult.RbtcAssetReport.Total.AsBigInt(), response.RbtcAssetReport.Total)

		useCase.AssertExpectations(t)
	})
}

func TestNewGetReportsAssetsHandler_Error(t *testing.T) {
	t.Run("should return 500 when use case returns an error", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)

		// Mock the use case to return an error
		useCase.EXPECT().Run(mock.Anything).Return(reports.GetAssetsReportResult{}, assert.AnError).Once()

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/assets", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		// Verify error status code
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Verify error response contains error message
		var errorResponse map[string]interface{}
		err = json.NewDecoder(rr.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse, "details")
		assert.Equal(t, "unknown error", errorResponse["message"])

		useCase.AssertExpectations(t)
	})

	t.Run("should return proper error structure", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)

		useCase.EXPECT().Run(mock.Anything).Return(reports.GetAssetsReportResult{}, assert.AnError).Once()

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/reports/assets", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		useCase.AssertExpectations(t)
	})
}

// mustParseBigWei creates a *entities.Wei from a string representation of a large number
func mustParseBigWei(s string) *entities.Wei {
	val := new(big.Int)
	val, ok := val.SetString(s, 10)
	if !ok {
		panic("failed to parse big wei: " + s)
	}
	return entities.NewBigWei(val)
}
