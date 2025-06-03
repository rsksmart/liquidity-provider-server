package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewErrorResponseWithDetails(t *testing.T) {
	details := map[string]any{
		test.AnyHash:   1,
		test.AnyString: test.AnyHash,
	}
	response := rest.NewErrorResponseWithDetails(test.AnyString, details, true)
	assert.Equal(t, test.AnyString, response.Message)
	assert.Equal(t, details, response.Details)
	assert.True(t, response.Recoverable)
}

func TestNewErrorResponse(t *testing.T) {
	response := rest.NewErrorResponse(test.AnyString, true)
	assert.Equal(t, test.AnyString, response.Message)
	assert.Empty(t, response.Details)
	assert.True(t, response.Recoverable)
}

func TestDetailsFromError(t *testing.T) {
	err := errors.New(test.AnyString)
	details := rest.DetailsFromError(err)
	assert.Len(t, details, 1)
	assert.Equal(t, err.Error(), details["error"])
}

func TestJsonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	rest.JsonResponse(w, http.StatusAccepted)
	assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestJsonErrorResponse(t *testing.T) {
	var body rest.ErrorResponse
	w := httptest.NewRecorder()
	response := rest.NewErrorResponse(test.AnyString, true)
	rest.JsonErrorResponse(w, http.StatusBadRequest, response)
	assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Equal(t, *response, body)
}

func TestJsonResponseWithBody(t *testing.T) {
	t.Run("response with nil body", func(t *testing.T) {
		w := httptest.NewRecorder()
		rest.JsonResponseWithBody[map[string]any](w, http.StatusAccepted, nil)
		assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
		assert.Equal(t, http.StatusAccepted, w.Code)
	})
	t.Run("response with body", func(t *testing.T) {
		var expectedBody, body map[string]string
		expectedBody = map[string]string{
			test.AnyHash:   test.AnyString,
			test.AnyString: test.AnyHash,
		}
		w := httptest.NewRecorder()
		rest.JsonResponseWithBody(w, http.StatusAccepted, &expectedBody)
		assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
		assert.Equal(t, http.StatusAccepted, w.Code)
		require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
		assert.Equal(t, expectedBody, body)
	})
	t.Run("encoding error", func(t *testing.T) {
		w := httptest.NewRecorder()
		circular := map[string]any{}
		circular["circular"] = circular
		rest.JsonResponseWithBody(w, http.StatusAccepted, &circular)
		var response rest.ErrorResponse
		require.NoError(t, json.NewDecoder(w.Body).Decode(&response))
		assert.Equal(t, "Unable to build response", response.Message)
	})
}

func TestDecodeRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New(test.AnyString)
	rest.DecodeRequestError(w, err)
	assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body rest.ErrorResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Contains(t, body.Message, test.AnyString)
	assert.True(t, body.Recoverable)
}

func TestValidateRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New(test.AnyString)
	rest.ValidateRequestError(w, err)
	assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body rest.ErrorResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Contains(t, body.Message, test.AnyString)
	assert.True(t, body.Recoverable)
}

func TestDecodeRequest(t *testing.T) {
	t.Run("decode request successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"key":"value"}`)))
		var body map[string]string
		err := rest.DecodeRequest(w, req, &body)
		require.NoError(t, err)
		assert.Equal(t, "value", body["key"])
	})
	t.Run("decode request with error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{not json}")))
		var body map[string]string
		err := rest.DecodeRequest(w, req, &body)
		require.Error(t, err)
		assert.Empty(t, body)
		assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRequiredQueryParam(t *testing.T) {
	require.ErrorContains(t, rest.RequiredQueryParam(test.AnyString), "required query parameter any value is missing")
}

func TestValidateRequest(t *testing.T) {
	t.Run("validate request successfully", func(t *testing.T) {
		req := pkg.AcceptQuoteRequest{QuoteHash: test.AnyHash}
		w := httptest.NewRecorder()
		err := rest.ValidateRequest(w, &req)
		require.NoError(t, err)
	})
	t.Run("handle not-validation error", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := map[string]string{"key": ""}
		var response rest.ErrorResponse
		err := rest.ValidateRequest(w, &body)
		require.Error(t, err)
		assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
		assert.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.NewDecoder(w.Body).Decode(&response))
		assert.Contains(t, response.Message, "Error validating request")
	})
	t.Run("handle validation error", func(t *testing.T) {
		req := pkg.PeginQuoteRequest{
			CallEoaOrContractAddress: test.AnyHash,
			ValueToTransfer:          1,
		}
		var response rest.ErrorResponse
		w := httptest.NewRecorder()
		err := rest.ValidateRequest(w, &req)
		require.Error(t, err)
		assert.Equal(t, rest.ContentTypeJson, w.Header().Get(rest.HeaderContentType))
		assert.Equal(t, http.StatusBadRequest, w.Code)
		require.NoError(t, json.NewDecoder(w.Body).Decode(&response))
		assert.Contains(t, response.Message, "validation error")
		assert.Len(t, response.Details, 2)
		for key := range response.Details {
			assert.NotEmpty(t, response.Details[key])
		}
		assert.True(t, response.Recoverable)
	})
}

func TestMaxDecimalPlacesValidation(t *testing.T) {
	type testStruct struct {
		Number float64 `validate:"max_decimal_places=4"`
	}

	testCases := []struct {
		value       float64
		expectError bool
		description string
	}{
		{value: 1.2345, expectError: false, description: "exactly 4 decimal places"},
		{value: 1.23456, expectError: true, description: "exceeds 4 decimal places"},
		{value: 1.0, expectError: false, description: "integer value"},
		{value: 1e-4, expectError: false, description: "scientific notation within limit"},
		{value: 1e-5, expectError: true, description: "scientific notation exceeds limit"},
		{value: 1.123456789, expectError: true, description: "many decimal places"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ts := testStruct{Number: tc.value}
			err := rest.RequestValidator.Struct(ts)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
