package handlers_test

import (
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewGetReportsAssetsHandler(t *testing.T) {
	const (
		path = "/reports/assets"
		verb = "GET"
	)

	successReturn := reports.GetAssetsReportResponse{
		BtcBalance:    big.NewInt(100000),
		RbtcBalance:   big.NewInt(100000),
		BtcLocked:     big.NewInt(15000),
		RbtcLocked:    big.NewInt(17500),
		BtcLiquidity:  big.NewInt(85000),
		RbtcLiquidity: big.NewInt(67500),
	}

	failReturn := reports.GetAssetsReportResponse{
		BtcBalance:    big.NewInt(0),
		RbtcBalance:   big.NewInt(0),
		BtcLocked:     big.NewInt(0),
		RbtcLocked:    big.NewInt(0),
		BtcLiquidity:  big.NewInt(0),
		RbtcLiquidity: big.NewInt(0),
	}

	t.Run("should return 200 on success", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)
		useCase.EXPECT().Run(mock.Anything).Return(successReturn, nil)

		handler := handlers.NewGetReportsAssetsHandler(useCase)
		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/reports/assets", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcBalance":%d`, successReturn.BtcBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcBalance":%d`, successReturn.RbtcBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLocked":%d`, successReturn.BtcLocked))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLocked":%d`, successReturn.RbtcLocked))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLiquidity":%d`, successReturn.BtcLiquidity))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLiquidity":%d`, successReturn.RbtcLiquidity))
		useCase.AssertExpectations(t)
	})
	t.Run("Should return an error when use case fail", func(t *testing.T) {
		useCase := new(mocks.GetAssetsReportUseCaseMock)
		useCase.EXPECT().Run(mock.Anything).Return(failReturn, assert.AnError)
		handler := handlers.NewGetReportsAssetsHandler(useCase)
		assert.HTTPError(t, handler, verb, path, nil)
		assert.HTTPBodyContains(t, handler, verb, path, nil, `{"error":"assert.AnError general error for testing"}`)
	})
}
