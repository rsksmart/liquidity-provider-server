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
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddPegoutCollateralHandlerHappyPath(t *testing.T) {
	reqBody := pkg.AddCollateralRequest{
		Amount: big.NewInt(1000000000000000000),
	}
	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	expectedBalance := entities.NewWei(2000000000000000000)

	mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
	mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).Return(expectedBalance, nil)

	handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody pkg.AddCollateralResponse
	err = json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	assert.Equal(t, expectedBalance.AsBigInt(), responseBody.NewCollateralBalance)

	mockUseCase.AssertExpectations(t)
}

// nolint:funlen
func TestAddPegoutCollateralHandlerErrorCases(t *testing.T) {

	t.Run("should handle malformed JSON in request body", func(t *testing.T) {
		malformedJSON := []byte(`{"amount": "not-a-number"`)
		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		mockUseCase.AssertNotCalled(t, "Run")

		var errorResponse map[string]interface{}
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "message")
		assert.Equal(t, "Error decoding request", errorResponse["message"])
	})

	t.Run("should handle request validation failure - missing required fields", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: nil, // Missing required field
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
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

	t.Run("should return 409 on InsufficientAmountError", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(100), // Too small amount
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), usecases.InsufficientAmountError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusConflict, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "not enough for minimum collateral", errorResponse["message"])
	})

	t.Run("should return 500 on unexpected errors", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(1000000000000000000),
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		unexpectedError := errors.New("unexpected blockchain error")
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), unexpectedError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)

		mockUseCase.AssertExpectations(t)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Equal(t, "unknown error", errorResponse["message"])
	})
}

// nolint:funlen
func TestAddPegoutCollateralHandlerErrorResponseFormat(t *testing.T) {

	t.Run("should set correct content type header on error", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(100),
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), usecases.InsufficientAmountError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should include timestamp in error response", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(100),
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), usecases.InsufficientAmountError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "timestamp")
		assert.NotZero(t, errorResponse["timestamp"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false for InsufficientAmountError", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(100),
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), usecases.InsufficientAmountError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
		handler := http.HandlerFunc(handlerFunc)

		handler.ServeHTTP(recorder, request)

		var errorResponse map[string]interface{}
		err = json.NewDecoder(recorder.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "recoverable")
		assert.Equal(t, false, errorResponse["recoverable"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should set recoverable to false on unexpected errors", func(t *testing.T) {
		reqBody := pkg.AddCollateralRequest{
			Amount: big.NewInt(1000000000000000000),
		}
		jsonBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		request := httptest.NewRequest(http.MethodPost, "/pegout/addCollateral", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.AddPegoutCollateralUseCaseMock)
		unexpectedError := errors.New("unexpected error")
		mockUseCase.On("Run", mock.AnythingOfType("*entities.Wei")).
			Return((*entities.Wei)(nil), unexpectedError)

		handlerFunc := handlers.NewAddPegoutCollateralHandler(mockUseCase)
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
