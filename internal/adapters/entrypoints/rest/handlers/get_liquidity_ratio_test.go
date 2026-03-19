package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func sampleLiquidityRatioDetail() lp.LiquidityRatioDetail {
	return lp.LiquidityRatioDetail{
		BtcPercentage:           50,
		RbtcPercentage:          50,
		MaxLiquidity:            entities.NewWei(1000),
		BtcTarget:               entities.NewWei(500),
		BtcThreshold:            entities.NewWei(600),
		RbtcTarget:              entities.NewWei(500),
		RbtcThreshold:           entities.NewWei(600),
		BtcCurrentBalance:       entities.NewWei(400),
		RbtcCurrentBalance:      entities.NewWei(550),
		BtcImpact:               lp.NetworkImpactDetail{Type: lp.NetworkImpactDeficit, Amount: entities.NewWei(100)},
		RbtcImpact:              lp.NetworkImpactDetail{Type: lp.NetworkImpactWithinTolerance, Amount: entities.NewWei(0)},
		CooldownActive:          false,
		CooldownEndTimestamp:    0,
		CooldownDurationSeconds: 10800,
		IsPreview:               false,
	}
}

func TestGetLiquidityRatioHandler_SuccessWithoutQueryParam(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/management/liquidity-ratio", nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.GetLiquidityRatioUseCaseMock)
	mockUseCase.On("Run", mock.Anything, uint64(0)).Return(sampleLiquidityRatioDetail(), nil)

	handler := handlers.NewGetLiquidityRatioHandler(mockUseCase)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var response pkg.LiquidityRatioResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, uint64(50), response.BtcPercentage)
	assert.Equal(t, uint64(50), response.RbtcPercentage)
	assert.Equal(t, "500", response.BtcTarget.String())
	assert.Equal(t, "400", response.BtcCurrentBalance.String())
	assert.Equal(t, "deficit", response.BtcImpact.Type)
	assert.Equal(t, "100", response.BtcImpact.Amount.String())
	assert.Equal(t, "withinTolerance", response.RbtcImpact.Type)
	assert.False(t, response.CooldownActive)
	assert.False(t, response.IsPreview)

	mockUseCase.AssertExpectations(t)
}

func TestGetLiquidityRatioHandler_SuccessWithQueryParam(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/management/liquidity-ratio?btcPercentage=70", nil)
	recorder := httptest.NewRecorder()

	previewDetail := sampleLiquidityRatioDetail()
	previewDetail.BtcPercentage = 70
	previewDetail.RbtcPercentage = 30
	previewDetail.IsPreview = true

	mockUseCase := new(mocks.GetLiquidityRatioUseCaseMock)
	mockUseCase.On("Run", mock.Anything, uint64(70)).Return(previewDetail, nil)

	handler := handlers.NewGetLiquidityRatioHandler(mockUseCase)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response pkg.LiquidityRatioResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, uint64(70), response.BtcPercentage)
	assert.Equal(t, uint64(30), response.RbtcPercentage)
	assert.True(t, response.IsPreview)

	mockUseCase.AssertExpectations(t)
}

func TestGetLiquidityRatioHandler_InvalidQueryParam(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/management/liquidity-ratio?btcPercentage=abc", nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.GetLiquidityRatioUseCaseMock)

	handler := handlers.NewGetLiquidityRatioHandler(mockUseCase)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockUseCase.AssertNotCalled(t, "Run", mock.Anything, mock.Anything)
}

func TestGetLiquidityRatioHandler_OutOfRangeQueryParam(t *testing.T) {
	cases := []string{"0", "9", "91", "101", "200"}
	for _, param := range cases {
		request := httptest.NewRequest(http.MethodGet, "/management/liquidity-ratio?btcPercentage="+param, nil)
		recorder := httptest.NewRecorder()

		mockUseCase := new(mocks.GetLiquidityRatioUseCaseMock)

		handler := handlers.NewGetLiquidityRatioHandler(mockUseCase)
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, "expected 400 for btcPercentage=%s", param)
		mockUseCase.AssertNotCalled(t, "Run", mock.Anything, mock.Anything)
	}
}

func TestGetLiquidityRatioHandler_UseCaseError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/management/liquidity-ratio", nil)
	recorder := httptest.NewRecorder()

	mockUseCase := new(mocks.GetLiquidityRatioUseCaseMock)
	mockUseCase.On("Run", mock.Anything, uint64(0)).Return(lp.LiquidityRatioDetail{}, assert.AnError)

	handler := handlers.NewGetLiquidityRatioHandler(mockUseCase)
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "unknown error")

	mockUseCase.AssertExpectations(t)
}
