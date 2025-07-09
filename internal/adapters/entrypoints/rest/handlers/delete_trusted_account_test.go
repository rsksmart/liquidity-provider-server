package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDeleteTrustedAccountHandler(t *testing.T) {
	t.Run("should return 204 on success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/management/trusted-accounts?address=0x123", nil)
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("DeleteTrustedAccount", mock.Anything, "0x123").Return(nil)
		useCase := lpuc.NewDeleteTrustedAccountUseCase(repo)
		handler := http.HandlerFunc(handlers.NewDeleteTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		repo.AssertExpectations(t)
	})
	t.Run("should return 400 when address is missing", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/management/trusted-accounts", nil)
		repo := &mocks.TrustedAccountRepositoryMock{}
		useCase := lpuc.NewDeleteTrustedAccountUseCase(repo)
		handler := http.HandlerFunc(handlers.NewDeleteTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should return 404 when account not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/management/trusted-accounts?address=0x123", nil)
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("DeleteTrustedAccount", mock.Anything, "0x123").Return(lp.TrustedAccountNotFoundError)
		useCase := lpuc.NewDeleteTrustedAccountUseCase(repo)
		handler := http.HandlerFunc(handlers.NewDeleteTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		repo.AssertExpectations(t)
	})
	t.Run("should return 500 on unexpected error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", "/management/trusted-accounts?address=0x123", nil)
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("DeleteTrustedAccount", mock.Anything, "0x123").Return(errors.New("database error"))
		useCase := lpuc.NewDeleteTrustedAccountUseCase(repo)
		handler := http.HandlerFunc(handlers.NewDeleteTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		repo.AssertExpectations(t)
	})
}
