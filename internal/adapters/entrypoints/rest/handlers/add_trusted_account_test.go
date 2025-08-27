package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
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
		Address:        "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
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

// Helper function to create test handler for address validation
func createAddressValidationHandler(expectError bool) (http.HandlerFunc, *mocks.TrustedAccountRepositoryMock, *mocks.TransactionSignerMock, *mocks.HashMock) {
	repo := &mocks.TrustedAccountRepositoryMock{}
	signer := &mocks.TransactionSignerMock{}
	hashMock := &mocks.HashMock{}

	if !expectError {
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("AddTrustedAccount", mock.Anything, mock.Anything).Return(nil)
	}

	useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
	handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))
	return handler, repo, signer, hashMock
}

// Helper function to validate error response
func validateErrorResponse(t *testing.T, body []byte, expectedSubstring string) {
	var errorResponse rest.ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	require.NoError(t, err, "Should be able to unmarshal error response")
	assert.Equal(t, "validation error", errorResponse.Message, "Main error message should be 'validation error'")

	if addressError, exists := errorResponse.Details["Address"]; exists {
		addressErrorStr, ok := addressError.(string)
		require.True(t, ok, "Address error should be a string")
		assert.Contains(t, addressErrorStr, expectedSubstring, "Address validation error should contain expected substring")
	} else {
		t.Errorf("Expected 'Address' field in error details, but got: %+v", errorResponse.Details)
	}
}

// Address validation test cases
var addressValidationTests = []struct {
	name           string
	address        string
	expectStatus   int
	expectError    bool
	errorSubstring string
}{
	{
		name:         "Valid RSK address should pass",
		address:      "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:           "A2",
		address:        "0x45400c53ebd0853CD246bC44", // Exact address from bug report
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:           "Long address should fail - bug scenario",
		address:        "0x45400c53ebd0853CD246bC44Extra12345678901234567890",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:           "Invalid characters should fail - bug scenario",
		address:        "0x45400c53ebd0853CD246bC44XYZ123456789012345",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:           "Missing 0x prefix should fail",
		address:        "7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:           "Empty address should fail",
		address:        "",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "is required",
	},
	{
		name:           "Only 0x should fail",
		address:        "0x",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:         "Wrong case but valid hex should pass",
		address:      "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
}

func TestAddTrustedAccountHandler_AddressValidation(t *testing.T) {
	for _, tc := range addressValidationTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks
			handler, repo, signer, hashMock := createAddressValidationHandler(tc.expectError)

			// Create request with test address
			requestBody := pkg.TrustedAccountRequest{
				Address:        tc.address,
				Name:           "Test Account",
				BtcLockingCap:  big.NewInt(0),
				RbtcLockingCap: big.NewInt(0),
			}

			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), "POST", "/management/trusted-accounts", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Validate response
			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			if tc.expectError {
				validateErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring)
			} else {
				// Verify mock expectations for successful cases
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}
