package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testPegoutQuote = quote.PegoutQuote{
	LbcAddress:            "0x85FaB18a0d06fb14651c8F5EE9C7f4b00D80d70c",
	LpRskAddress:          "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
	BtcRefundAddress:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
	RskRefundAddress:      "0x1234567890123456789012345678901234567890",
	LpBtcAddress:          "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2",
	CallFee:               entities.NewWei(100),
	PenaltyFee:            entities.NewWei(200),
	Nonce:                 1,
	DepositAddress:        "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
	Value:                 entities.NewWei(1000),
	AgreementTimestamp:    1640995200,
	DepositDateLimit:      1641001200,
	DepositConfirmations:  6,
	TransferConfirmations: 3,
	TransferTime:          3600,
	ExpireDate:            1641000000,
	ExpireBlock:           1000000,
	GasFee:                entities.NewWei(50),
	ProductFeeAmount:      entities.NewWei(25),
}

var testRetainedPegoutQuote = quote.RetainedPegoutQuote{
	QuoteHash:          "", // Set per test case
	DepositAddress:     "0x9D93929A9099be4355fC2389FbF253982F9dF47c",
	Signature:          "b24831aac7230910087d9818b378a31679be5e3991a7227cc160bc3add09e1645a26e9c740e3467f53953d7ec086c82bf8ef0eb03c118d0382ee6049a8f0119f1c",
	RequiredLiquidity:  entities.NewWei(1500),
	State:              quote.PegoutStateSendPegoutSucceeded,
	UserRskTxHash:      "0x2c73de184c80797c04a655217d121588e8d5c228d3e0cc26187cb249123aa7c3",
	LpBtcTxHash:        "619c4d69ccaa5f78aaa2284817cf070609ac40af3792916ca3d0ef82b14af75f",
	RefundPegoutTxHash: "0x3a0feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db89",
	BridgeRefundTxHash: "0x4b1feaef4d803468ba5bfc1db78f4d2568de1b7cf002dec5991c469e6719db90",
}

var testPegoutCreationData = quote.PegoutCreationData{
	FeeRate:       utils.NewBigFloat64(0.0001),
	FeePercentage: utils.NewBigFloat64(1.5),
	GasPrice:      entities.NewWei(55),
	FixedFee:      entities.NewWei(100000),
}

//nolint:funlen
func TestNewGetPegoutQuoteStatusHandler_SuccessfulResponse(t *testing.T) {
	testQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"

	retainedQuote := testRetainedPegoutQuote
	retainedQuote.QuoteHash = testQuoteHash

	testWatchedQuote := quote.NewWatchedPegoutQuote(testPegoutQuote, retainedQuote, testPegoutCreationData)

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	mockUseCase.On("Run", mock.Anything, testQuoteHash).Return(
		testWatchedQuote,
		nil,
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+testQuoteHash, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var response pkg.PegoutQuoteStatusDTO
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)

	// Validate detail (PegoutQuote)
	expectedDetail := pkg.ToPegoutQuoteDTO(testPegoutQuote)
	assert.Equal(t, expectedDetail, response.Detail)

	// Validate status (RetainedPegoutQuote)
	expectedStatus := pkg.ToRetainedPegoutQuoteDTO(retainedQuote)
	assert.Equal(t, expectedStatus, response.Status)

	// Validate creation data
	expectedCreationData := pkg.ToPegoutCreationDataDTO(testPegoutCreationData)
	assert.Equal(t, expectedCreationData, response.CreationData)

	// Validate specific critical fields
	assert.Equal(t, testQuoteHash, response.Status.QuoteHash)
	assert.Equal(t, string(quote.PegoutStateSendPegoutSucceeded), response.Status.State)
	assert.Equal(t, retainedQuote.Signature, response.Status.Signature)
	assert.Equal(t, retainedQuote.DepositAddress, response.Status.DepositAddress)
	assert.Equal(t, retainedQuote.RequiredLiquidity.AsBigInt(), response.Status.RequiredLiquidity)
	assert.Equal(t, retainedQuote.UserRskTxHash, response.Status.UserRskTxHash)
	assert.Equal(t, retainedQuote.LpBtcTxHash, response.Status.LpBtcTxHash)
	assert.Equal(t, retainedQuote.RefundPegoutTxHash, response.Status.RefundPegoutTxHash)

	// Validate quote detail fields
	assert.Equal(t, testPegoutQuote.BtcRefundAddress, response.Detail.BtcRefundAddr)
	assert.Equal(t, testPegoutQuote.RskRefundAddress, response.Detail.RSKRefundAddr)
	assert.Equal(t, testPegoutQuote.LpRskAddress, response.Detail.LPRSKAddr)
	assert.Equal(t, testPegoutQuote.Value.AsBigInt(), response.Detail.Value)
	assert.Equal(t, testPegoutQuote.CallFee.AsBigInt(), response.Detail.CallFee)
	assert.Equal(t, testPegoutQuote.PenaltyFee.AsBigInt(), response.Detail.PenaltyFee)

	// Validate creation data fields
	assert.Equal(t, testPegoutCreationData.GasPrice.AsBigInt(), response.CreationData.GasPrice)
	assert.Equal(t, testPegoutCreationData.FixedFee.AsBigInt(), response.CreationData.FixedFee)
	expectedFeePercentage, _ := testPegoutCreationData.FeePercentage.Native().Float64()
	assert.InDelta(t, expectedFeePercentage, response.CreationData.FeePercentage, 0.0)
	expectedFeeRate, _ := testPegoutCreationData.FeeRate.Native().Float64()
	assert.InDelta(t, expectedFeeRate, response.CreationData.FeeRate, 0.0)

	// Count fields in all three main objects and assert expected counts
	detailFieldCount := reflect.TypeOf(response.Detail).NumField()
	statusFieldCount := reflect.TypeOf(response.Status).NumField()
	creationDataFieldCount := reflect.TypeOf(response.CreationData).NumField()

	const expectedDetailFields = 19      // PegoutQuoteDTO has 19 fields
	const expectedStatusFields = 9       // RetainedPegoutQuoteDTO has 9 fields
	const expectedCreationDataFields = 4 // PegoutCreationDataDTO has 4 fields

	assert.Equal(t, expectedDetailFields, detailFieldCount, "Detail object should have exactly %d fields", expectedDetailFields)
	assert.Equal(t, expectedStatusFields, statusFieldCount, "Status object should have exactly %d fields", expectedStatusFields)
	assert.Equal(t, expectedCreationDataFields, creationDataFieldCount, "CreationData object should have exactly %d fields", expectedCreationDataFields)
}

func TestNewGetPegoutQuoteStatusHandler_MissingQuoteHashQueryParameter(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 0")
}

func TestNewGetPegoutQuoteStatusHandler_QuoteNotFound(t *testing.T) {
	testQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	mockUseCase.On("Run", mock.Anything, testQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		usecases.QuoteNotFoundError,
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+testQuoteHash, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Equal(t, "Quote not found", errorResponse.Message)
}

func TestNewGetPegoutQuoteStatusHandler_QuoteNotAccepted(t *testing.T) {
	testQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	mockUseCase.On("Run", mock.Anything, testQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		usecases.QuoteNotAcceptedError,
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+testQuoteHash, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusConflict, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Equal(t, usecases.QuoteNotAcceptedError.Error(), errorResponse.Message)
}

func TestNewGetPegoutQuoteStatusHandler_UnhandledError(t *testing.T) {
	var errorMessage = "database connection failed"
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)

	// Use a valid quote hash so validation passes
	validQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"
	mockUseCase.On("Run", mock.Anything, validQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		errors.New(errorMessage),
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+validQuoteHash, nil)
	require.NoError(t, err)

	// Assert that the error is logged
	defer test.AssertLogContains(t, errorMessage)()

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Should return 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, res.Code)

	// Should contain generic error message but not internal error details for security
	responseBody := res.Body.String()
	assert.Contains(t, responseBody, "Internal server error")
	assert.NotContains(t, responseBody, "database connection failed")

	// Should be valid JSON error response
	var errorResponse rest.ErrorResponse
	err = json.Unmarshal([]byte(responseBody), &errorResponse)
	require.NoError(t, err)
	assert.Equal(t, "Internal server error", errorResponse.Message)
	assert.False(t, errorResponse.Recoverable)
}

func TestNewGetPegoutQuoteStatusHandler_WrappedUseCaseError(t *testing.T) {
	testQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	wrappedErr := fmt.Errorf("wrapped: %w", usecases.QuoteNotFoundError)
	mockUseCase.On("Run", mock.Anything, testQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		wrappedErr,
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+testQuoteHash, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Equal(t, "Quote not found", errorResponse.Message)
}

func TestNewGetPegoutQuoteStatusHandler_EmptyQueryParameterValue(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash=", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 0")
}

func TestNewGetPegoutQuoteStatusHandler_QueryParameterWithSpaces(t *testing.T) {
	quoteHashWithSpaces := " 8d1ba2cb559a6ebe41f  19131602467e1d939682d651b2a91e55b86bc664a6819 "

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+quoteHashWithSpaces, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Should fail validation due to spaces making length != 64
	assert.Equal(t, http.StatusBadRequest, res.Code)

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
}

func TestNewGetPegoutQuoteStatusHandler_WrongParameterName(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?hash=8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Should fail validation because wrong parameter name means empty quoteHash
	assert.Equal(t, http.StatusBadRequest, res.Code)

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 0")
}

func TestNewGetPegoutQuoteStatusHandler_CaseSensitiveParameterName(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?QuoteHash=8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Should fail validation because wrong case means empty quoteHash
	assert.Equal(t, http.StatusBadRequest, res.Code)

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 0")
}

func TestNewGetPegoutQuoteStatusHandler_MultipleQueryParameters(t *testing.T) {
	testQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"

	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	mockUseCase.On("Run", mock.Anything, testQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		usecases.QuoteNotFoundError,
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	// Create request with multiple parameters
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?otherParam=value&quoteHash="+testQuoteHash, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	// Should reach use case since quoteHash is valid
	assert.NotEqual(t, http.StatusBadRequest, res.Code)
}

func TestNewGetPegoutQuoteStatusHandler_InvalidQuoteHashFormat(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	// Test with invalid hex characters (exactly 64 characters but contains 'G' which is invalid hex)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash=8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a681G", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "must be a valid hex string")
}

func TestNewGetPegoutQuoteStatusHandler_ErrorResponseFormat(t *testing.T) {
	// Test that all error responses follow the expected format
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)

	// Use valid quote hash so validation passes
	validQuoteHash := "8d1ba2cb559a6ebe41f19131602467e1d939682d651b2a91e55b86bc664a6819"
	mockUseCase.On("Run", mock.Anything, validQuoteHash).Return(
		quote.WatchedPegoutQuote{},
		errors.New("test error"),
	).Once()

	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash="+validQuoteHash, nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify response structure
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.True(t, json.Valid(rr.Body.Bytes()))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	// Verify error response has required fields
	assert.NotEmpty(t, errorResponse.Message)
	assert.NotZero(t, errorResponse.Timestamp)
	assert.NotNil(t, errorResponse.Details)
}

func TestNewGetPegoutQuoteStatusHandler_QuoteHashEmptyString(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	// Test with `""` string as quote hash value
	req, err := http.NewRequestWithContext(context.Background(), "GET", `/pegout/status?quoteHash=""`, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 2")
}

func TestNewGetPegoutQuoteStatusHandler_QuoteHashNullString(t *testing.T) {
	mockUseCase := mocks.NewPegoutStatusUseCaseMock(t)
	// No mock setup needed since validation fails before reaching use case
	handler := handlers.NewGetPegoutQuoteStatusHandler(mockUseCase)

	// Test with literal "null" as quote hash value
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/pegout/status?quoteHash=null", nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))

	var errorResponse rest.ErrorResponse
	err = json.Unmarshal(res.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.True(t, errorResponse.Recoverable)
	assert.Contains(t, errorResponse.Message, "invalid or missing parameter quoteHash")
	assert.Contains(t, errorResponse.Details["error"], "expected 64 characters, got 4")
}
