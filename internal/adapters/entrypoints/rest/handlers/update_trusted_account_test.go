package handlers_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createValidRequest() *http.Request {
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
	request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestNewUpdateTrustedAccountHandler(t *testing.T) { //nolint:funlen
	t.Run("should return 204 on success", func(t *testing.T) {
		request := createValidRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.Anything).Return(nil)
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
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
		request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should handle malformed json", func(t *testing.T) {
		malformedJSON := []byte(`{"address": "0x123", "name": "Test Account",`)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should return 404 when account not found", func(t *testing.T) {
		request := createValidRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.Anything).Return(liquidity_provider.ErrTrustedAccountNotFound)
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
