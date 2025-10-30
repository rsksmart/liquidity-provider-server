package handlers_test

import (
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
func TestSetPegoutConfigHandler(t *testing.T) {
	t.Run("should return success response if there are no errors", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("UpsertPegoutConfiguration", mock.Anything, mock.Anything).Return(nil)
		walletMock := &mocks.RskWalletMock{}
		walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(100), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}

		useCase := uc_lp.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPegoutConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "expireTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000", "expireBlocks": 100, "bridgeTransactionMin": "5000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegout/configuration", strings.NewReader(reqBody))
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
		useCase := uc_lp.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPegoutConfigHandler(useCase)
		reqBody := `{"configuration": }`
		req := httptest.NewRequest(http.MethodPost, "/pegout/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return bad request if bridgeTransactionMin is below or equal to bridge minimum", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		walletMock := &mocks.RskWalletMock{}
		hashMock := &mocks.HashMock{}
		bridge := &mocks.BridgeMock{}
		// Set bridge minimum to 1000, so bridgeTransactionMin of 100 will fail
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPegoutConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "expireTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "500", "expireBlocks": 100, "bridgeTransactionMin": "100"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegout/configuration", strings.NewReader(reqBody))
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
		useCase := uc_lp.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPegoutConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "expireTime": 300, "penaltyFee": "1000", "fixedFee": "-500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000", "expireBlocks": 100, "bridgeTransactionMin": "5000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegout/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Validation error")
	})

	t.Run("should return server internal error if use case fails with non-validation error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("UpsertPegoutConfiguration", mock.Anything, mock.Anything).Return(assert.AnError)
		walletMock := &mocks.RskWalletMock{}
		walletMock.On("SignBytes", mock.Anything).Return([]byte{1, 2, 3}, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return([]byte{4, 5, 6})
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(100), nil)
		contracts := blockchain.RskContracts{Bridge: bridge}
		useCase := uc_lp.NewSetPegoutConfigUseCase(lpRepository, walletMock, hashMock.Hash, contracts)
		handler := handlers.NewSetPegoutConfigHandler(useCase)
		reqBody := `{"configuration": {"timeForDeposit": 600, "expireTime": 300, "penaltyFee": "1000", "fixedFee": "500", "feePercentage": 1.5, "maxValue": "10000000", "minValue": "1000", "expireBlocks": 100, "bridgeTransactionMin": "5000"}}`
		req := httptest.NewRequest(http.MethodPost, "/pegout/configuration", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "unknown error")
	})
}
