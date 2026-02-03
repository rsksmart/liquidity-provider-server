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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 100, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {},  "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "-1000", "reimbursementWindowBlocks": 100, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "nan", "reimbursementWindowBlocks": 100, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if reimbursementWindowBlocks is zero", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 0}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if reimbursementWindowBlocks is missing", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excess tolerance is not provided", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "reimbursementWindowBlocks": 10, "maxLiquidity": "1000"}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excess tolerance fixed value is not provided", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true, "percentageValue":0}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excess tolerance percentage value is not provided", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"100"}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excess tolerance percentage value is negative", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue": -1}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
	t.Run("should return bad request if excess tolerance fixed value is negative", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"-100", "percentageValue":1}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000.5", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000000000000000000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
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
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 100, "excessTolerance":{"isFixed":true,"fixedValue":"100", "percentageValue":0}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return success with valid excessTolerance in fixed mode", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":true,"fixedValue":"5000000000000000000", "percentageValue":0}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return success with valid excessTolerance in percentage mode", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":false,"fixedValue":"0", "percentageValue":5.5}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		useCase.AssertExpectations(t)
	})
	t.Run("should return bad request if excessTolerance percentageValue exceeds 100", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}, "publicLiquidityCheck": true, "maxLiquidity": "1000", "reimbursementWindowBlocks": 10, "excessTolerance":{"isFixed":false,"fixedValue":"0", "percentageValue":150}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		useCase.AssertNotCalled(t, "Run")
	})
}
