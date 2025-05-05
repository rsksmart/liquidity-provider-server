package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewGetTrustedAccountsHandler(t *testing.T) {
	t.Run("should return 200 with accounts on success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/management/trusted-accounts", nil)
		mockAccounts := []liquidity_provider.TrustedAccountDetails{
			{
				Address:          "0x123",
				Name:             "Test Account 1",
				Btc_locking_cap:  entities.NewWei(100),
				Rbtc_locking_cap: entities.NewWei(200),
			},
			{
				Address:          "0x456",
				Name:             "Test Account 2",
				Btc_locking_cap:  entities.NewWei(300),
				Rbtc_locking_cap: entities.NewWei(400),
			},
		}
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(mockAccounts, nil)
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo)
		handler := http.HandlerFunc(handlers.NewGetTrustedAccountsHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusOK, recorder.Code)
		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, 2, len(response.Accounts))
		assert.Equal(t, "0x123", response.Accounts[0].Address)
		assert.Equal(t, "Test Account 1", response.Accounts[0].Name)
		assert.Equal(t, "0x456", response.Accounts[1].Address)
		assert.Equal(t, "Test Account 2", response.Accounts[1].Name)
		repo.AssertExpectations(t)
	})

	t.Run("should return 500 on error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/management/trusted-accounts", nil)
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(nil, errors.New("database error"))
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo)
		handler := http.HandlerFunc(handlers.NewGetTrustedAccountsHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		repo.AssertExpectations(t)
	})
}
