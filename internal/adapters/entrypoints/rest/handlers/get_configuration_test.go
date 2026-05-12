package handlers_test

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	uc_lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestGetConfigurationHandlerHappyPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/configuration", nil)
	recorder := httptest.NewRecorder()

	maxValueBigInt := new(big.Int)
	maxValueBigInt.SetString("10000000000000000000", 10)

	expectedConfig := uc_lp.FullConfiguration{
		General: liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"1000000000000000000": 10,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"500000000000000000": 2,
			},
			PublicLiquidityCheck: true,
		},
		Pegin: liquidity_provider.PeginConfiguration{
			TimeForDeposit: 3600,
			CallTime:       7200,
			PenaltyFee:     entities.NewWei(1000000000000000),
			FixedFee:       entities.NewWei(500000000000000),
			FeePercentage:  utils.NewBigFloat64(0.5),
			MaxValue:       entities.NewBigWei(maxValueBigInt),
			MinValue:       entities.NewWei(100000000000000000),
		},
		Pegout: liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       1800,
			ExpireTime:           3600,
			PenaltyFee:           entities.NewWei(2000000000000000),
			FixedFee:             entities.NewWei(600000000000000),
			FeePercentage:        utils.NewBigFloat64(0.75),
			MaxValue:             entities.NewWei(5000000000000000000),
			MinValue:             entities.NewWei(50000000000000000),
			ExpireBlocks:         100,
			BridgeTransactionMin: entities.NewWei(10000000000000000),
		},
	}

	mockUseCase := new(mocks.GetConfigUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedConfig)

	handlerFunc := handlers.NewGetConfigurationHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var responseBody struct {
		General liquidity_provider.GeneralConfiguration `json:"general"`
		Pegin   pkg.PeginConfigurationDTO               `json:"pegin"`
		Pegout  pkg.PegoutConfigurationDTO              `json:"pegout"`
	}
	err := json.NewDecoder(recorder.Body).Decode(&responseBody)
	require.NoError(t, err)

	// Verify General configuration
	assert.True(t, responseBody.General.PublicLiquidityCheck)
	assert.Len(t, responseBody.General.RskConfirmations, 1)
	assert.Len(t, responseBody.General.BtcConfirmations, 1)

	// Verify Pegin configuration
	assert.Equal(t, uint32(3600), responseBody.Pegin.TimeForDeposit)
	assert.Equal(t, uint32(7200), responseBody.Pegin.CallTime)
	assert.Equal(t, "1000000000000000", responseBody.Pegin.PenaltyFee)
	assert.Equal(t, "500000000000000", responseBody.Pegin.FixedFee)
	assert.InDelta(t, 0.5, responseBody.Pegin.FeePercentage, 0.001)
	assert.Equal(t, "10000000000000000000", responseBody.Pegin.MaxValue)
	assert.Equal(t, "100000000000000000", responseBody.Pegin.MinValue)

	// Verify Pegout configuration
	assert.Equal(t, uint32(1800), responseBody.Pegout.TimeForDeposit)
	assert.Equal(t, uint32(3600), responseBody.Pegout.ExpireTime)
	assert.Equal(t, "2000000000000000", responseBody.Pegout.PenaltyFee)
	assert.Equal(t, "600000000000000", responseBody.Pegout.FixedFee)
	assert.InDelta(t, 0.75, responseBody.Pegout.FeePercentage, 0.001)
	assert.Equal(t, "5000000000000000000", responseBody.Pegout.MaxValue)
	assert.Equal(t, "50000000000000000", responseBody.Pegout.MinValue)
	assert.Equal(t, uint64(100), responseBody.Pegout.ExpireBlocks)
	assert.Equal(t, "10000000000000000", responseBody.Pegout.BridgeTransactionMin)

	mockUseCase.AssertExpectations(t)
}

func TestGetConfigurationHandlerResponseStructure(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/configuration", nil)
	recorder := httptest.NewRecorder()

	expectedConfig := uc_lp.FullConfiguration{
		General: liquidity_provider.GeneralConfiguration{
			RskConfirmations:     liquidity_provider.ConfirmationsPerAmount{},
			BtcConfirmations:     liquidity_provider.ConfirmationsPerAmount{},
			PublicLiquidityCheck: true,
		},
		Pegin: liquidity_provider.PeginConfiguration{
			TimeForDeposit: 100,
			CallTime:       200,
			PenaltyFee:     entities.NewWei(1),
			FixedFee:       entities.NewWei(2),
			FeePercentage:  utils.NewBigFloat64(1.5),
			MaxValue:       entities.NewWei(1000),
			MinValue:       entities.NewWei(10),
		},
		Pegout: liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       300,
			ExpireTime:           400,
			PenaltyFee:           entities.NewWei(3),
			FixedFee:             entities.NewWei(4),
			FeePercentage:        utils.NewBigFloat64(2.5),
			MaxValue:             entities.NewWei(2000),
			MinValue:             entities.NewWei(20),
			ExpireBlocks:         50,
			BridgeTransactionMin: entities.NewWei(5),
		},
	}

	mockUseCase := new(mocks.GetConfigUseCaseMock)
	mockUseCase.On("Run", mock.Anything).Return(expectedConfig)

	handlerFunc := handlers.NewGetConfigurationHandler(mockUseCase)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(recorder, request)

	// Verify the response contains exactly 3 top-level keys
	var rawResponse map[string]interface{}
	err := json.NewDecoder(recorder.Body).Decode(&rawResponse)
	require.NoError(t, err)

	assert.Len(t, rawResponse, 3)
	assert.Contains(t, rawResponse, "general")
	assert.Contains(t, rawResponse, "pegin")
	assert.Contains(t, rawResponse, "pegout")

	mockUseCase.AssertExpectations(t)
}
