package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAcceptPeginAuthenticatedQuoteHandler(t *testing.T) {
	quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	signature := "validSignature123"
	acceptedQuote := quote.AcceptedQuote{
		Signature:      "signedHash123",
		DepositAddress: "depositAddress456",
	}

	reqBody := pkg.AcceptAuthenticatedQuoteRequest{
		QuoteHash: quoteHash,
		Signature: signature,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(acceptedQuote, nil)

	handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody pkg.AcceptPeginRespose
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
	assert.Equal(t, acceptedQuote.DepositAddress, responseBody.BitcoinDepositAddressHash)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen,maintidx
func TestAcceptPeginAuthenticatedQuoteHandlerErrorCases(t *testing.T) {
	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		// Create a request with malformed JSON (missing closing brace)
		malformedJSON := []byte(`{"quoteHash": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", "signature": "validSignature123"`)
		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// Request decoding should fail with 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse["message"], "Error decoding request")
	})

	t.Run("should handle request validation failure - missing signature", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			Signature: "", // Empty signature will fail validation
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		// Validation should fail with 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse["message"], "validation error")
	})

	t.Run("should return 400 on invalid quote hash format", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "123456789XABCDEF", // Invalid - contains non-hex character 'X'
			Signature: "validSignature123",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse["message"], "invalid quote hash")
	})

	t.Run("should return 404 when quote not found", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, usecases.QuoteNotFoundError)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "quote not found", errorResponse["message"])
	})

	t.Run("should return 410 when quote is expired", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, usecases.ExpiredQuoteError)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusGone, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "expired quote", errorResponse["message"])
	})

	t.Run("should return 409 when not enough liquidity", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, usecases.NoLiquidityError)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "not enough liquidity", errorResponse["message"])
	})

	t.Run("should return 409 when locking cap exceeded", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, usecases.LockingCapExceededError)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "locking cap exceeded", errorResponse["message"])
	})

	t.Run("should return 500 when trusted account is tampered", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, liquidity_provider.ErrTamperedTrustedAccount)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "error fetching trusted account", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signature := "validSignature123"

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signature,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected database error")
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, unexpectedError)

		handlerFunc := handlers.NewAcceptPeginAuthenticatedQuoteHandlerWithInterface(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "unknown error", errorResponse["message"])
	})
}
