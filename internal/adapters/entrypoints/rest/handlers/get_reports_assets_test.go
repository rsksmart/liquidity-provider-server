package handlers

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"math/big"
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
		ctx := context.Background()
		useCase := mocks.NewAssetsReportUseCaseMock(t)

		useCase.On("Run", ctx).Return(successReturn, nil)

		handler := NewGetReportsAssetsHandler(useCase)
		assert.HTTPSuccess(t, handler, verb, path, nil)
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcBalance":%d`, successReturn.BtcBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcBalance":%d`, successReturn.RbtcBalance))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLocked":%d`, successReturn.BtcLocked))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLocked":%d`, successReturn.RbtcLocked))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"btcLiquidity":%d`, successReturn.BtcLiquidity))
		assert.HTTPBodyContains(t, handler, verb, path, nil, fmt.Sprintf(`"rbtcLiquidity":%d`, successReturn.RbtcLiquidity))
	})
	t.Run("Should return an error when use case fail", func(t *testing.T) {
		ctx := context.Background()
		useCase := mocks.NewAssetsReportUseCaseMock(t)

		useCase.On("Run", ctx).Return(failReturn, assert.AnError)
		handler := NewGetReportsAssetsHandler(useCase)
		assert.HTTPError(t, handler, verb, path, nil)
		assert.HTTPBodyContains(t, handler, verb, path, nil, `{"error":"assert.AnError general error for testing"}`)
	})
}
