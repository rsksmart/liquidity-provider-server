package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSingleFlightConstants(t *testing.T) {
	t.Run("PegInReportSingleFlightKey should have correct value", func(t *testing.T) {
		assert.Equal(t, "pegin-report-singleflight", handlers.PegInReportSingleFlightKey)
	})

	t.Run("PegOutReportSingleFlightKey should have correct value", func(t *testing.T) {
		assert.Equal(t, "pegout-report-singleflight", handlers.PegOutReportSingleFlightKey)
	})

	t.Run("RevenueReportSingleFlightKey should have correct value", func(t *testing.T) {
		assert.Equal(t, "revenue-report-singleflight", handlers.RevenueReportSingleFlightKey)
	})

	t.Run("SummariesReportSingleFlightKey should have correct value", func(t *testing.T) {
		assert.Equal(t, "summaries-report-singleflight", handlers.SummariesReportSingleFlightKey)
	})
}

func TestSingleFlightGroup(t *testing.T) {
	t.Run("SingleFlightGroup should not be nil", func(t *testing.T) {
		assert.NotNil(t, handlers.SingleFlightGroup)
	})
}

func TestCalculateSingleFlightKey(t *testing.T) {
	t.Run("should combine base key with request URL", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reports/pegin?from=2024-01-01&to=2024-12-31", nil)

		result := handlers.CalculateSingleFlightKey(handlers.PegInReportSingleFlightKey, req)

		assert.Equal(t, "pegin-report-singleflight_/reports/pegin?from=2024-01-01&to=2024-12-31", result)
	})

	t.Run("should handle URL without query parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reports/pegout", nil)

		result := handlers.CalculateSingleFlightKey(handlers.PegOutReportSingleFlightKey, req)

		assert.Equal(t, "pegout-report-singleflight_/reports/pegout", result)
	})

	t.Run("should generate different keys for different URLs", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodGet, "/reports/pegin?from=2024-01-01", nil)
		req2 := httptest.NewRequest(http.MethodGet, "/reports/pegin?from=2024-06-01", nil)

		key1 := handlers.CalculateSingleFlightKey(handlers.PegInReportSingleFlightKey, req1)
		key2 := handlers.CalculateSingleFlightKey(handlers.PegInReportSingleFlightKey, req2)

		assert.NotEqual(t, key1, key2)
	})

	t.Run("should generate same key for identical requests", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodGet, "/reports/revenue?from=2024-01-01&to=2024-12-31", nil)
		req2 := httptest.NewRequest(http.MethodGet, "/reports/revenue?from=2024-01-01&to=2024-12-31", nil)

		key1 := handlers.CalculateSingleFlightKey(handlers.RevenueReportSingleFlightKey, req1)
		key2 := handlers.CalculateSingleFlightKey(handlers.RevenueReportSingleFlightKey, req2)

		assert.Equal(t, key1, key2)
	})
}
