package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	uc_lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//nolint:funlen
func TestSetPeginConfigHandler(t *testing.T) {
	t.Run("should return success response if there are no errors", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("UpsertPeginConfiguration", mock.Anything, mock.Anything).Return(nil)
		walletMock := &mocks.RskWalletMock{}
		walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(100), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}

		useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPeginConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return bad request if it can't decode the request", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPeginConfigHandler(useCase)
		reqBody := `{"configuration": }`
		req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return bad request if minValue is below or equal to bridge minimum", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		// Set bridge minimum to 1000, so minValue of 100 will fail
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPeginConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "100"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Validation error")
		assert.Contains(t, w.Body.String(), "requested amount below bridge")
	})

	t.Run("should return bad request for negative wei value", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(100), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPeginConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "-500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Validation error")
	})

	t.Run("should return server internal error if use case fails with non-validation error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("UpsertPeginConfiguration", mock.Anything, mock.Anything).Return(assert.AnError)
		walletMock := &mocks.RskWalletMock{}
		walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(100), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPeginConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "unknown error")
	})
}

//nolint:funlen
func TestSetPeginConfigHandler_ValidationErrorScenarios(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}

	testCases := []struct {
		name            string
		bridgeMin       int64
		minValue        string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "minValue equals bridge minimum",
			bridgeMin:       1000,
			minValue:        "1000",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Validation error",
		},
		{
			name:            "minValue less than bridge minimum",
			bridgeMin:       1000,
			minValue:        "500",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Validation error",
		},
		{
			name:            "minValue greater than bridge minimum",
			bridgeMin:       100,
			minValue:        "1000",
			expectedStatus:  http.StatusNoContent,
			expectedMessage: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bridge := &mocks.BridgeMock{}
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(tc.bridgeMin), nil)
			contracts := blockchain.RskContracts{Bridge: bridge}

			if tc.expectedStatus == http.StatusNoContent {
				lpRepository.On("UpsertPeginConfiguration", mock.Anything, mock.Anything).Return(nil).Once()
				walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil).Once()
				hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6}).Once()
			}

			useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
			handler := handlers.NewSetPeginConfigHandler(useCase)
			reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "` + tc.minValue + `"}}`
			req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.Background())
			w := httptest.NewRecorder()
			handler(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedMessage != "" {
				assert.Contains(t, w.Body.String(), tc.expectedMessage)
			}
		})
	}
}

func TestSetPeginConfigHandler_EnsureValidationErrorsReturn400(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	bridge := &mocks.BridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil)
	contracts := blockchain.RskContracts{Bridge: bridge}

	useCase := uc_lp.NewSetPeginConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
	handler := handlers.NewSetPeginConfigHandler(useCase)

	// Test with minValue = bridgeMin (should fail)
	reqBody := `{"configuration": {"timeForDeposit": 600, "callTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000"}}`
	req := httptest.NewRequest(http.MethodPost, "/pegin/configuration", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)

	// Must return 400, not 500
	assert.Equal(t, http.StatusBadRequest, w.Code, "Validation errors MUST return 400 Bad Request")
	assert.Contains(t, w.Body.String(), "Validation error")
	assert.NotContains(t, w.Body.String(), "unknown error")
}
