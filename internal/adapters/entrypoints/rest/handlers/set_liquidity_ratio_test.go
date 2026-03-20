package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetLiquidityRatioHandler_Success(t *testing.T) {
	useCase := new(mocks.SetLiquidityRatioUseCaseMock)
	useCase.EXPECT().Run(mock.Anything, uint64(60)).Return(nil)

	handler := handlers.NewSetLiquidityRatioHandler(useCase)
	req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{"btcPercentage": 60}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	useCase.AssertExpectations(t)
}

func TestSetLiquidityRatioHandler_ValidationErrors(t *testing.T) {
	t.Run("should return bad request if btcPercentage is below minimum", func(t *testing.T) {
		useCase := new(mocks.SetLiquidityRatioUseCaseMock)
		handler := handlers.NewSetLiquidityRatioHandler(useCase)
		req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{"btcPercentage": 5}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return bad request if btcPercentage is above maximum", func(t *testing.T) {
		useCase := new(mocks.SetLiquidityRatioUseCaseMock)
		handler := handlers.NewSetLiquidityRatioHandler(useCase)
		req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{"btcPercentage": 95}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return bad request if btcPercentage is missing", func(t *testing.T) {
		useCase := new(mocks.SetLiquidityRatioUseCaseMock)
		handler := handlers.NewSetLiquidityRatioHandler(useCase)
		req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})

	t.Run("should return bad request on invalid JSON", func(t *testing.T) {
		useCase := new(mocks.SetLiquidityRatioUseCaseMock)
		handler := handlers.NewSetLiquidityRatioHandler(useCase)
		req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{invalid`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
}

func TestSetLiquidityRatioHandler_UseCaseError(t *testing.T) {
	useCase := new(mocks.SetLiquidityRatioUseCaseMock)
	useCase.EXPECT().Run(mock.Anything, uint64(60)).Return(assert.AnError)

	handler := handlers.NewSetLiquidityRatioHandler(useCase)
	req := httptest.NewRequest(http.MethodPost, "/management/liquidity-ratio", strings.NewReader(`{"btcPercentage": 60}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "unknown error")
	useCase.AssertExpectations(t)
}
