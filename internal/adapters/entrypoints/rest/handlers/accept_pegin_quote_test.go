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

func TestAcceptPeginQuoteHandlerHappyPath(t *testing.T) {
	quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	acceptedQuote := quote.AcceptedQuote{
		Signature:      "signedHash123",
		DepositAddress: "depositAddress456",
	}

	reqBody := pkg.AcceptQuoteRequest{
		QuoteHash: quoteHash,
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(acceptedQuote, nil)

	handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

// nolint:funlen
func TestAcceptPeginQuoteHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		// Create a request with malformed JSON (missing closing brace)
		malformedJSON := []byte(`{"quoteHash": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"`)
		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

	t.Run("should handle request validation failure", func(t *testing.T) {
		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: "", // Empty quote hash will fail validation
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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
		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: "123456789XABCDEF", // Invalid - contains non-hex character 'X'
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.QuoteNotFoundError)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.ExpiredQuoteError)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, usecases.NoLiquidityError)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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

	t.Run("should return 200 with already retained quote", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		retainedQuote := quote.AcceptedQuote{
			Signature:      "previouslySignedHash",
			DepositAddress: "previouslyRetainedAddress",
		}

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(retainedQuote, nil)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var responseBody pkg.AcceptPeginRespose
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		require.NoError(t, err)
		assert.Equal(t, retainedQuote.Signature, responseBody.Signature)
		assert.Equal(t, retainedQuote.DepositAddress, responseBody.BitcoinDepositAddressHash)
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		reqBody := pkg.AcceptQuoteRequest{
			QuoteHash: quoteHash,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegin/acceptQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		unexpectedError := errors.New("unexpected database error")
		mockUseCase.On("Run", mock.Anything, quoteHash, "").Return(quote.AcceptedQuote{}, unexpectedError)

		handlerFunc := handlers.NewAcceptPeginQuoteHandler(mockUseCase)
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
