package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestHandleAcceptQuoteError(t *testing.T) {
	t.Run("should return 404 when quote not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.QuoteNotFoundError)

		assert.Equal(t, 404, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "quote not found", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, usecases.QuoteNotFoundError.Error(), errorResponse.Details["error"])
	})

	t.Run("should return 410 when quote is expired", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.ExpiredQuoteError)

		assert.Equal(t, 410, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "expired quote", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, usecases.ExpiredQuoteError.Error(), errorResponse.Details["error"])
	})

	t.Run("should return 409 when not enough liquidity", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.NoLiquidityError)

		assert.Equal(t, 409, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "not enough liquidity", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, usecases.NoLiquidityError.Error(), errorResponse.Details["error"])
	})

	t.Run("should return 409 when locking cap exceeded", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.LockingCapExceededError)

		assert.Equal(t, 409, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "locking cap exceeded", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, usecases.LockingCapExceededError.Error(), errorResponse.Details["error"])
	})

	t.Run("should return 500 when tampered trusted account error", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, liquidity_provider.TamperedTrustedAccountError)

		assert.Equal(t, 500, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "error fetching trusted account", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, liquidity_provider.TamperedTrustedAccountError.Error(), errorResponse.Details["error"])
	})

	t.Run("should return 500 with unknown error message for unexpected errors", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		unexpectedError := errors.New("unexpected database connection error")

		handlers.HandleAcceptQuoteError(recorder, unexpectedError)

		assert.Equal(t, 500, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, handlers.UnknownErrorMessage, errorResponse.Message)
		assert.False(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
		assert.Equal(t, unexpectedError.Error(), errorResponse.Details["error"])
	})

	t.Run("should handle wrapped errors correctly", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		wrappedError := usecases.WrapUseCaseError(usecases.AcceptPeginQuoteId, usecases.ExpiredQuoteError)

		handlers.HandleAcceptQuoteError(recorder, wrappedError)

		assert.Equal(t, 410, recorder.Code)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "expired quote", errorResponse.Message)
		assert.True(t, errorResponse.Recoverable)
		assert.Contains(t, errorResponse.Details, "error")
	})

	t.Run("should set correct content type header", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.QuoteNotFoundError)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		handlers.HandleAcceptQuoteError(recorder, usecases.QuoteNotFoundError)

		var errorResponse rest.ErrorResponse
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.NotZero(t, errorResponse.Timestamp)
	})
}
