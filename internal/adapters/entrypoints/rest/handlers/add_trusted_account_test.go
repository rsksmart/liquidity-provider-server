package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
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

// Helper function to validate error response for specific fields
func validateFieldErrorResponse(t *testing.T, body []byte, expectedSubstring string, fieldNames ...string) {
	var errorResponse rest.ErrorResponse
	err := json.Unmarshal(body, &errorResponse)
	require.NoError(t, err, "Should be able to unmarshal error response")
	assert.Equal(t, "validation error", errorResponse.Message, "Main error message should be 'validation error'")

	// Check if any of the specified fields has the expected error
	found := false
	for fieldName, fieldError := range errorResponse.Details {
		for _, expectedField := range fieldNames {
			if fieldName == expectedField {
				errorStr, ok := fieldError.(string)
				require.True(t, ok, "%s error should be a string", fieldName)
				if strings.Contains(errorStr, expectedSubstring) {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Errorf("Expected error containing '%s' in fields %v, but got: %+v", expectedSubstring, fieldNames, errorResponse.Details)
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
				BtcLockingCap:  big.NewInt(1000000),
				RbtcLockingCap: big.NewInt(1000000),
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
				validateFieldErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring, "Address")
			} else {
				// Verify mock expectations for successful cases
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}

// Cap validation test cases
var capValidationTests = []struct {
	name           string
	btcCap         *big.Int
	rbtcCap        *big.Int
	expectStatus   int
	expectError    bool
	errorSubstring string
}{
	{
		name:         "Valid positive caps should pass",
		btcCap:       big.NewInt(1000000000000000000),
		rbtcCap:      big.NewInt(2000000000000000000),
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:           "Zero BTC cap should fail",
		btcCap:         big.NewInt(0),
		rbtcCap:        big.NewInt(1),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Zero RBTC cap should fail",
		btcCap:         big.NewInt(1),
		rbtcCap:        big.NewInt(0),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Both zero caps should fail",
		btcCap:         big.NewInt(0),
		rbtcCap:        big.NewInt(0),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Negative BTC cap should fail",
		btcCap:         big.NewInt(-1),
		rbtcCap:        big.NewInt(1),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Negative RBTC cap should fail",
		btcCap:         big.NewInt(1),
		rbtcCap:        big.NewInt(-100),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Both negative caps should fail",
		btcCap:         big.NewInt(-1),
		rbtcCap:        big.NewInt(-1),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
}

func TestAddTrustedAccountHandler_CapValidation(t *testing.T) {
	for _, tc := range capValidationTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks
			handler, repo, signer, hashMock := createAddressValidationHandler(tc.expectError)

			// Create request with test cap values
			requestBody := pkg.TrustedAccountRequest{
				Address:        "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
				Name:           "Test Account",
				BtcLockingCap:  tc.btcCap,
				RbtcLockingCap: tc.rbtcCap,
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
				validateFieldErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring, "BtcLockingCap", "RbtcLockingCap")
			} else {
				// Verify mock expectations for successful cases
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}

// Decimal value test cases (testing JSON parsing behavior)
var decimalValueTests = []struct {
	name           string
	jsonPayload    string
	expectStatus   int
	errorSubstring string
}{
	{
		name:           "Decimal with .1 should be rejected during JSON parsing",
		jsonPayload:    `{"name": "Test Account", "address": "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b", "btcLockingCap": 1000000.1, "rbtcLockingCap": 2000000.5}`,
		expectStatus:   http.StatusBadRequest,
		errorSubstring: "cannot unmarshal",
	},
	{
		name:           "String values should be rejected during JSON parsing",
		jsonPayload:    `{"name": "Test Account", "address": "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b", "btcLockingCap": "1000000", "rbtcLockingCap": "2000000"}`,
		expectStatus:   http.StatusBadRequest,
		errorSubstring: "cannot unmarshal",
	},
	{
		name:           "Scientific notation should be rejected during JSON parsing",
		jsonPayload:    `{"name": "Test Account", "address": "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b", "btcLockingCap": 1e6, "rbtcLockingCap": 2e6}`,
		expectStatus:   http.StatusBadRequest,
		errorSubstring: "cannot unmarshal",
	},
}

func TestAddTrustedAccountHandler_DecimalValueRejection(t *testing.T) {
	for _, tc := range decimalValueTests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.TrustedAccountRepositoryMock{}
			signer := &mocks.TransactionSignerMock{}
			hashMock := &mocks.HashMock{}

			useCase := lpuc.NewAddTrustedAccountUseCase(repo, signer, hashMock.Hash)
			handler := http.HandlerFunc(handlers.NewAddTrustedAccountHandler(useCase))

			req, err := http.NewRequestWithContext(context.Background(), "POST", "/management/trusted-accounts", strings.NewReader(tc.jsonPayload))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Should fail during JSON parsing before reaching validation
			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			var errorResponse rest.ErrorResponse
			err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
			require.NoError(t, err, "Should be able to unmarshal error response")
			assert.Contains(t, errorResponse.Message, tc.errorSubstring, "Error message should contain expected substring")

			repo.AssertNotCalled(t, "AddTrustedAccount")
			signer.AssertNotCalled(t, "SignBytes")
			hashMock.AssertNotCalled(t, "Hash")
		})
	}
}
