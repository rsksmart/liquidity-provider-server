package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAcceptPegoutQuoteHandlerHappyPath(t *testing.T) {
	quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	acceptedQuote := quote.AcceptedQuote{
		Signature:      "signedHash123",
		DepositAddress: "0xLbcAddress456",
	}

	reqBody := pkg.AcceptQuoteRequest{
		QuoteHash: quoteHash,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(acceptedQuote, nil)

	handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.AcceptPegoutResponse
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
	assert.Equal(t, acceptedQuote.DepositAddress, responseBody.LbcAddress)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen, maintidx
func TestAcceptPegoutQuoteHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		malformedJSON := []byte(`{"quoteHash": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"`)
		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "Error decoding request", errorResponse["message"])
		assert.Contains(t, errorResponse, "details")
		assert.NotEmpty(t, errorResponse["details"])
	})

	t.Run("should handle request validation failure", func(t *testing.T) {
		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: "", // Empty quote hash will fail validation
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Contains(t, errorResponse["message"], "validation error")
	})

	t.Run("should return 400 on invalid quote hash length", func(t *testing.T) {
		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: "1234567890abcdef", // Invalid - wrong length (16 chars instead of 64)
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
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

	t.Run("should return 400 on quote hash with invalid hex characters", func(t *testing.T) {
		// 64 characters but contains invalid hex character 'G'
		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeG",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
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

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.QuoteNotFoundError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "invalid quote hash", errorResponse["message"])
	})

	t.Run("should return 410 when quote is expired", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.ExpiredQuoteError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
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

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.NoLiquidityError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
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

	t.Run("should return 200 with already accepted quote", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		acceptedQuote := quote.AcceptedQuote{
			Signature:      "previouslySignedHash",
			DepositAddress: "0xPreviousLbcAddress",
		}

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(acceptedQuote, nil)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var responseBody pkg.AcceptPegoutResponse
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
		assert.Equal(t, acceptedQuote.DepositAddress, responseBody.LbcAddress)
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected database error")
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, unexpectedError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
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

// nolint:funlen
func TestAcceptPegoutQuoteHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.QuoteNotFoundError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.QuoteNotFoundError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include recoverable flag in error response", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.ExpiredQuoteError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, true, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on unexpected errors", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected error")
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, unexpectedError)

		handlerFunc := handlers.NewAcceptPegoutQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, false, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})
}
