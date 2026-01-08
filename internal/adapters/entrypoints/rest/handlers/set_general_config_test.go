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

//nolint:funlen
func TestSetGeneralConfigHandler(t *testing.T) {
	t.Run("should return success response if there are no errors", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "0", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return bad request if it can't decode the request", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": }`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if the request validation fails", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"NaN": 10}, "rskConfirmations": {"10": 20}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if btcConfirmations is empty map", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {}, "rskConfirmations": {"10": 20}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "must not be empty")
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if rskConfirmations is empty map", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "must not be empty")
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if maxLiquidity is negative", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "-1000"}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if maxLiquidity is not a number", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "nan"}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if maxLiquidity has decimal places", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		// Backend expects wei values (integers), so decimal values should fail
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000.5"}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return success with maxLiquidity having 18 digits precision", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		// 1 RBTC in wei (18 decimal places) - valid large integer
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000000000000000000", "excessToleranceFixed": "0", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return server internal error if the use case fails", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(assert.AnError)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "0", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return success with valid excessToleranceFixed", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "5000000000000000000", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return success with valid excessTolerancePercentage", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "0", "excessTolerancePercentage": 5.5}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return bad request if excessToleranceFixed is negative", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "-100", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excessToleranceFixed is not a number", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "not_a_number", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excessToleranceFixed has decimal places", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "100.5", "excessTolerancePercentage": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excessTolerancePercentage is negative", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "0", "excessTolerancePercentage": -5}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excessTolerancePercentage exceeds 100", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessToleranceFixed": "0", "excessTolerancePercentage": 150}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
}
