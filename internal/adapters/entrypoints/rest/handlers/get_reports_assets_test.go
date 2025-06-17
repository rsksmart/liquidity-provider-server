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

	successReturn := reports.GetAssetsReportResult{
		RbtcLockedLbc:      big.NewInt(100000),
		RbtcLockedForUsers: big.NewInt(100000),
		RbtcWaitingRefund:  big.NewInt(15000),
		RbtcLiquidity:      big.NewInt(17500),
		RbtcWalletBalance:  big.NewInt(85000),
		BtcLockedForUsers:  big.NewInt(67500),
		BtcLiquidity:       big.NewInt(67500),
		BtcWalletBalance:   big.NewInt(67500),
		BtcRebalancing:     big.NewInt(67500),
	}

	failReturn := reports.GetAssetsReportResult{
		RbtcLockedLbc:      big.NewInt(0),
		RbtcLockedForUsers: big.NewInt(0),
		RbtcWaitingRefund:  big.NewInt(0),
		RbtcLiquidity:      big.NewInt(0),
		RbtcWalletBalance:  big.NewInt(0),
		BtcLockedForUsers:  big.NewInt(0),
		BtcLiquidity:       big.NewInt(0),
		BtcWalletBalance:   big.NewInt(0),
		BtcRebalancing:     big.NewInt(0),
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
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLockedLbc":%d`, successReturn.RbtcLockedLbc))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLockedForUsers":%d`, successReturn.RbtcLockedForUsers))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcWaitingRefund":%d`, successReturn.RbtcWaitingRefund))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLiquidity":%d`, successReturn.RbtcLiquidity))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcWalletBalance":%d`, successReturn.RbtcWalletBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLockedForUsers":%d`, successReturn.BtcLockedForUsers))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLiquidity":%d`, successReturn.BtcLiquidity))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcWalletBalance":%d`, successReturn.BtcWalletBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcRebalancing":%d`, successReturn.BtcRebalancing))
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
