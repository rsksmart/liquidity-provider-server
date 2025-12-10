package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestGetTrustedAccountsHandlerHappyPath(t *testing.T) {
	t.Run("should return 200 with multiple accounts", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		account1 := liquidity_provider.TrustedAccountDetails{
			Address:        "0x1234567890abcdef1234567890abcdef12345678",
			Name:           "Test Account 1",
			BtcLockingCap:  entities.NewWei(1000000000000000000),
			RbtcLockingCap: entities.NewWei(2000000000000000000),
		}
		account2 := liquidity_provider.TrustedAccountDetails{
			Address:        "0xabcdef1234567890abcdef1234567890abcdef12",
			Name:           "Test Account 2",
			BtcLockingCap:  entities.NewWei(3000000000000000000),
			RbtcLockingCap: entities.NewWei(4000000000000000000),
		}
		mockAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
			{Value: account1, Hash: "hash1", Signature: "sig1"},
			{Value: account2, Hash: "hash2", Signature: "sig2"},
		}

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(mockAccounts, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify response structure
		require.Len(t, response.Accounts, 2)

		// Verify first account
		assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", response.Accounts[0].Address)
		assert.Equal(t, "Test Account 1", response.Accounts[0].Name)
		assert.Equal(t, "1000000000000000000", response.Accounts[0].BtcLockingCap.String())
		assert.Equal(t, "2000000000000000000", response.Accounts[0].RbtcLockingCap.String())

		// Verify second account
		assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef12", response.Accounts[1].Address)
		assert.Equal(t, "Test Account 2", response.Accounts[1].Name)
		assert.Equal(t, "3000000000000000000", response.Accounts[1].BtcLockingCap.String())
		assert.Equal(t, "4000000000000000000", response.Accounts[1].RbtcLockingCap.String())

		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return 200 with single account", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		account := liquidity_provider.TrustedAccountDetails{
			Address:        "0x1234567890abcdef1234567890abcdef12345678",
			Name:           "Single Account",
			BtcLockingCap:  entities.NewWei(5000000000000000000),
			RbtcLockingCap: entities.NewWei(6000000000000000000),
		}
		mockAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
			{Value: account, Hash: "hash1", Signature: "sig1"},
		}

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(mockAccounts, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Len(t, response.Accounts, 1)
		assert.Equal(t, "Single Account", response.Accounts[0].Name)

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetTrustedAccountsHandlerEdgeCases(t *testing.T) {
	t.Run("should return 200 with empty accounts list", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{}

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(mockAccounts, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Empty(t, response.Accounts)
		assert.Empty(t, response.Accounts)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return 200 with nil accounts (converted to empty)", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(nil, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Empty(t, response.Accounts)

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetTrustedAccountsHandlerErrorCases(t *testing.T) {
	t.Run("should return 500 on use case error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(nil, errors.New("database connection error"))

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Equal(t, handlers.UnknownErrorMessage, errorResponse.Message)
		assert.False(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetTrustedAccountsHandlerResponseFormat(t *testing.T) {
	t.Run("response should have correct Content-Type", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{}

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(mockAccounts, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("error response should have timestamp", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(nil, errors.New("error"))

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse rest.ErrorResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.NotZero(t, errorResponse.Timestamp)
	})

	t.Run("response should have accounts key even when empty", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/management/trusted-accounts", nil)
		recorder := httptest.NewRecorder()

		mockAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{}

		mockUseCase := new(mocks.GetTrustedAccountsUseCaseMock)
		mockUseCase.On("Run", mock.Anything).Return(mockAccounts, nil)

		handlerFunc := handlers.NewGetTrustedAccountsHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// Verify JSON structure has "accounts" key
		var rawResponse map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &rawResponse)
		require.NoError(t, err)
		assert.Contains(t, rawResponse, "accounts")
	})
}
