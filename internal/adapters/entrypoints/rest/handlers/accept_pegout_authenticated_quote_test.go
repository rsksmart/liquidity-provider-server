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

func TestAcceptPegoutAuthenticatedQuoteHandler_Success(t *testing.T) {
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

	request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
	mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(acceptedQuote, nil)

	handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBody pkg.AcceptPegoutResponse
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
	assert.Equal(t, acceptedQuote.DepositAddress, responseBody.LbcAddress)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen
func TestAcceptPegoutAuthenticatedQuoteHandler_RequestErrors(t *testing.T) {
	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		// Create a request with malformed JSON (missing closing brace)
		malformedJSON := []byte(`{"quoteHash": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", "signature": "validSignature123"`)
		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
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

	t.Run("should handle request validation failure - missing quote hash", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "", // Empty quote hash will fail validation
			Signature: "validSignature123",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
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

	t.Run("should handle request validation failure - missing signature", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			Signature: "", // Empty signature will fail validation
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
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
}

func TestAcceptPegoutAuthenticatedQuoteHandler_QuoteHashValidation(t *testing.T) {
	t.Run("should return 400 on invalid quote hash format - too short", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "123456789abcdef", // Too short
			Signature: "validSignature123",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
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

	t.Run("should return 400 on invalid quote hash format - non-hex characters", func(t *testing.T) {
		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: "123456789XABCDEF123456789XABCDEF123456789XABCDEF123456789XABCDEF", // Contains non-hex character 'X'
			Signature: "validSignature123",
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
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
}

// nolint:funlen
func TestAcceptPegoutAuthenticatedQuoteHandler_UseCaseErrors(t *testing.T) {
	testCases := []struct {
		name           string
		useCaseError   error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "quote not found error",
			useCaseError:   usecases.QuoteNotFoundError,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "quote not found",
		},
		{
			name:           "expired quote error",
			useCaseError:   usecases.ExpiredQuoteError,
			expectedStatus: http.StatusGone,
			expectedMsg:    "expired quote",
		},
		{
			name:           "no liquidity error",
			useCaseError:   usecases.NoLiquidityError,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "not enough liquidity",
		},
		{
			name:           "locking cap exceeded error",
			useCaseError:   usecases.LockingCapExceededError,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "locking cap exceeded",
		},
		{
			name:           "tampered trusted account error",
			useCaseError:   liquidity_provider.TamperedTrustedAccountError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "error fetching trusted account",
		},
		{
			name:           "unexpected error",
			useCaseError:   errors.New("unexpected database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "unknown error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
			signature := "validSignature123"

			reqBody := pkg.AcceptAuthenticatedQuoteRequest{
				QuoteHash: quoteHash,
				Signature: signature,
			}
			jsonBody, err := json.Marshal(reqBody)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
			mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(quote.AcceptedQuote{}, tc.useCaseError)

			handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
			handler := http.HandlerFunc(handlerFunc)

			handler.ServeHTTP(recorder, request)

			// Verify the response status and message (this indirectly tests that HandleAcceptQuoteError is working correctly)
			assert.Equal(t, tc.expectedStatus, recorder.Code)

			var errorResponse map[string]interface{}
			err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse, "message")
			assert.Equal(t, tc.expectedMsg, errorResponse["message"])

			mockUseCase.AssertExpectations(t)
		})
	}
}

// nolint:funlen
func TestAcceptPegoutAuthenticatedQuoteHandler_SignatureProcessing(t *testing.T) {
	t.Run("should strip 0x prefix from signature", func(t *testing.T) {
		quoteHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		signatureWithPrefix := "0xvalidSignature123"
		signatureWithoutPrefix := "validSignature123"
		acceptedQuote := quote.AcceptedQuote{
			Signature:      "signedHash123",
			DepositAddress: "depositAddress456",
		}

		reqBody := pkg.AcceptAuthenticatedQuoteRequest{
			QuoteHash: quoteHash,
			Signature: signatureWithPrefix,
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		// Verify that the use case receives the signature WITHOUT the "0x" prefix
		mockUseCase.On("Run", mock.Anything, quoteHash, signatureWithoutPrefix).Return(acceptedQuote, nil)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody pkg.AcceptPegoutResponse
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		require.NoError(t, err)

		assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
		assert.Equal(t, acceptedQuote.DepositAddress, responseBody.LbcAddress)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("should leave signature unchanged when no 0x prefix", func(t *testing.T) {
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

		request := httptest.NewRequest(http.MethodPost, "/pegout/acceptAuthenticatedQuote", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AcceptQuoteUseCaseMock)
		// Verify that the use case receives the signature unchanged
		mockUseCase.On("Run", mock.Anything, quoteHash, signature).Return(acceptedQuote, nil)

		handlerFunc := handlers.NewAcceptPegoutAuthenticatedQuoteHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var responseBody pkg.AcceptPegoutResponse
		err = json.NewDecoder(recorder.Body).Decode(&responseBody)
		require.NoError(t, err)

		assert.Equal(t, acceptedQuote.Signature, responseBody.Signature)
		assert.Equal(t, acceptedQuote.DepositAddress, responseBody.LbcAddress)

		mockUseCase.AssertExpectations(t)
	})
}
