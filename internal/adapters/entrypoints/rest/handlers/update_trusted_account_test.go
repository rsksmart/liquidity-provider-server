package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createValidRequest() *http.Request {
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
	request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func TestNewUpdateTrustedAccountHandler(t *testing.T) { //nolint:funlen
	t.Run("should return 204 on success", func(t *testing.T) {
		request := createValidRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.Anything).Return(nil)
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
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
		request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(jsonBody))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should handle malformed json", func(t *testing.T) {
		malformedJSON := []byte(`{"address": "0x123", "name": "Test Account",`)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", "/management/trusted-accounts", bytes.NewBuffer(malformedJSON))
		request.Header.Set("Content-Type", "application/json")
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
	t.Run("should return 404 when account not found", func(t *testing.T) {
		request := createValidRequest()
		recorder := httptest.NewRecorder()
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.Anything).Return(liquidity_provider.TrustedAccountNotFoundError)
		useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}

// Helper function to create update validation handler with mocks
func createUpdateValidationHandler(expectError bool) (http.HandlerFunc, *mocks.TrustedAccountRepositoryMock, *mocks.TransactionSignerMock, *mocks.HashMock) {
	repo := &mocks.TrustedAccountRepositoryMock{}
	signer := &mocks.TransactionSignerMock{}
	hashMock := &mocks.HashMock{}

	if !expectError {
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.Anything).Return(nil)
	}

	useCase := lpuc.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
	handler := http.HandlerFunc(handlers.NewUpdateTrustedAccountHandler(useCase))
	return handler, repo, signer, hashMock
}

// Helper function to validate field error responses
func validateUpdateFieldErrorResponse(t *testing.T, body []byte, expectedSubstring string, fieldNames ...string) {
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
	assert.True(t, found, "Expected error substring '%s' not found in any of the fields %v. Details: %+v", expectedSubstring, fieldNames, errorResponse.Details)
}

// Address validation test cases
var updateAddressValidationTests = []struct {
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
		address:        "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b1234567890abcdef", // Too long
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: eth_addr",
	},
	{
		name:           "Invalid characters should fail - bug scenario",
		address:        "0xGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG", // Invalid hex characters
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
		address:      "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b", // lowercase
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
}

func TestUpdateTrustedAccountHandler_AddressValidation(t *testing.T) {
	for _, tc := range updateAddressValidationTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks
			handler, repo, signer, hashMock := createUpdateValidationHandler(tc.expectError)

			// Create request with test address
			requestBody := pkg.TrustedAccountRequest{
				Address:        tc.address,
				Name:           "Test Account",
				BtcLockingCap:  big.NewInt(1000000),
				RbtcLockingCap: big.NewInt(1000000),
			}

			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), "PUT", "/management/trusted-accounts", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			if tc.expectError {
				validateUpdateFieldErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring, "Address")
			} else {
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}

// Cap validation test cases
var updateCapValidationTests = []struct {
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
		btcCap:         big.NewInt(-1000000),
		rbtcCap:        big.NewInt(1000000),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Negative RBTC cap should fail",
		btcCap:         big.NewInt(1000000),
		rbtcCap:        big.NewInt(-1000000),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
	{
		name:           "Both negative caps should fail",
		btcCap:         big.NewInt(-1000000),
		rbtcCap:        big.NewInt(-1000000),
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "must be a positive integer",
	},
}

func TestUpdateTrustedAccountHandler_CapValidation(t *testing.T) {
	for _, tc := range updateCapValidationTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks
			handler, repo, signer, hashMock := createUpdateValidationHandler(tc.expectError)

			// Create request with test cap values
			requestBody := pkg.TrustedAccountRequest{
				Address:        "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
				Name:           "Test Account",
				BtcLockingCap:  tc.btcCap,
				RbtcLockingCap: tc.rbtcCap,
			}

			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), "PUT", "/management/trusted-accounts", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			if tc.expectError {
				validateUpdateFieldErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring, "BtcLockingCap", "RbtcLockingCap")
			} else {
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}

// Decimal value test cases
var updateDecimalValueTests = []struct {
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

func TestUpdateTrustedAccountHandler_DecimalValueRejection(t *testing.T) {
	for _, tc := range updateDecimalValueTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks (but they won't be called due to JSON parse failure)
			handler, _, _, _ := createUpdateValidationHandler(true)

			req, err := http.NewRequestWithContext(context.Background(), "PUT", "/management/trusted-accounts", strings.NewReader(tc.jsonPayload))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			// Verify error message contains expected substring
			body := rr.Body.String()
			assert.Contains(t, body, tc.errorSubstring, "Error message should contain expected substring")
		})
	}
}

// Name validation test cases
var updateNameValidationTests = []struct {
	name           string
	nameValue      string
	expectStatus   int
	expectError    bool
	errorSubstring string
}{
	{
		name:         "Valid company name should pass",
		nameValue:    "Appleton Inc.",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:         "Valid name with international characters should pass",
		nameValue:    "Juan & María Inc.",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:         "Valid name with special characters should pass",
		nameValue:    "O'Rulz & Associates",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:         "Valid name with numbers and punctuation should pass",
		nameValue:    "Company-2024 (Holdings) Ltd.",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:         "Single character name should pass",
		nameValue:    "A",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:         "Name at max length (100 chars) should pass",
		nameValue:    "Company Name Test String For Maximum Length Validation Testing Purposesxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		expectStatus: http.StatusNoContent,
		expectError:  false,
	},
	{
		name:           "Empty name should fail",
		nameValue:      "",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "is required",
	},
	{
		name:           "Whitespace-only name should fail",
		nameValue:      "    ",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "cannot be blank",
	},
	{
		name:           "Tab and newline only name should fail",
		nameValue:      "\t\n\r ",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "cannot be blank",
	},
	{
		name:           "Name exceeding max length (101 chars) should fail",
		nameValue:      "Company Name Test String For Maximum Length Validation Testing Purposesxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		expectStatus:   http.StatusBadRequest,
		expectError:    true,
		errorSubstring: "validation failed: max",
	},
}

func TestUpdateTrustedAccountHandler_NameValidation(t *testing.T) {
	for _, tc := range updateNameValidationTests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup handler with mocks
			handler, repo, signer, hashMock := createUpdateValidationHandler(tc.expectError)

			// Create request with test name value
			requestBody := pkg.TrustedAccountRequest{
				Address:        "0x7C4890A0f1D4bBf2C669Ac2d1efFa185c505359b",
				Name:           tc.nameValue,
				BtcLockingCap:  big.NewInt(1000000),
				RbtcLockingCap: big.NewInt(1000000),
			}

			bodyBytes, err := json.Marshal(requestBody)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(context.Background(), "PUT", "/management/trusted-accounts", bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectStatus, rr.Code, "Status code should match expected")

			if tc.expectError {
				validateUpdateFieldErrorResponse(t, rr.Body.Bytes(), tc.errorSubstring, "Name")
			} else {
				repo.AssertExpectations(t)
				signer.AssertExpectations(t)
				hashMock.AssertExpectations(t)
			}
		})
	}
}
