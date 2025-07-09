package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createValidAddRequest() *http.Request {
	btcLockingCap := new(big.Int)
	btcLockingCap.SetString("1000000000000000000", 10)
	rbtcLockingCap := new(big.Int)
	rbtcLockingCap.SetString("2000000000000000000", 10)
	reqBody := &pkg.TrustedAccountRequest{
		Address:        "0x123",
		Name:           "Test Account",
		BtcLockingCap:  btcLockingCap,
		RbtcLockingCap: rbtcLockingCap,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}
	request := httptest.NewRequest("POST", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	return request
}

// nolint:funlen
func TestNewAddTrustedAccountHandler(t *testing.T) {
	t.Run("should return 204 on success", func(t *testing.T) {
		request := createValidAddRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("AddTrustedAccount", mock.Anything, mock.Anything).Return(nil)
		useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("should return 400 on invalid request", func(t *testing.T) {
		reqBody := pkg.TrustedAccountRequest{}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should handle malformed json", func(t *testing.T) {
		malformedJSON := []byte(`{"address": "0x123", "name": "Test Account",`)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/management/trusted-accounts", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should return 400 when account already exists", func(t *testing.T) {
		request := createValidAddRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("AddTrustedAccount", mock.Anything, mock.Anything).Return(lp.DuplicateTrustedAccountError)
		useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("should return 500 on unexpected error", func(t *testing.T) {
		request := createValidAddRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("AddTrustedAccount", mock.Anything, mock.Anything).Return(errors.New("database error"))
		useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
