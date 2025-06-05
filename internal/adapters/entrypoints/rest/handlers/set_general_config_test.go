package handlers_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetGeneralConfigHandler(t *testing.T) {
	t.Run("should return success response if there are no errors", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(nil)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}}}`
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
	t.Run("should return server internal error if the request validation fails", func(t *testing.T) {
		useCase := new(mocks.SetGeneralConfigUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything).Return(assert.AnError)

		handler := handlers.NewSetGeneralConfigHandler(useCase)
		reqBody := `{"configuration": {"btcConfirmations": {"5": 10}, "rskConfirmations": {"10": 20}}}`
		req := httptest.NewRequest(http.MethodPost, "/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		useCase.AssertExpectations(t)
	})
}
