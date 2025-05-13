package handlers_test

import (
	"encoding/hex"
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
		mockHashBytes := []byte("mockhash12345678")
		mockHashHex := hex.EncodeToString(mockHashBytes)
		account1 := liquidity_provider.TrustedAccountDetails{
			Address:          "0x123",
			Name:             "Test Account 1",
			Btc_locking_cap:  entities.NewWei(100),
			Rbtc_locking_cap: entities.NewWei(200),
		}
		account2 := liquidity_provider.TrustedAccountDetails{
			Address:          "0x456",
			Name:             "Test Account 2",
			Btc_locking_cap:  entities.NewWei(300),
			Rbtc_locking_cap: entities.NewWei(400),
		}
		mockSignedAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
			{
				Value:     account1,
				Hash:      mockHashHex,
				Signature: "sig1",
			},
			{
				Value:     account2,
				Hash:      mockHashHex,
				Signature: "sig2",
			},
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(mockSignedAccounts, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(mockHashBytes)
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewGetTrustedAccountsHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusOK, recorder.Code)
		var response pkg.TrustedAccountsResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Len(t, response.Accounts, 2)
		assert.Equal(t, "0x123", response.Accounts[0].Address)
		assert.Equal(t, "Test Account 1", response.Accounts[0].Name)
		assert.Equal(t, "0x456", response.Accounts[1].Address)
		assert.Equal(t, "Test Account 2", response.Accounts[1].Name)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("should return 500 on error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/management/trusted-accounts", nil)
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(nil, errors.New("database error"))
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo, hashMock.Hash)
		handler := http.HandlerFunc(handlers.NewGetTrustedAccountsHandler(useCase))
		handler.ServeHTTP(recorder, request)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		repo.AssertExpectations(t)
		hashMock.AssertNotCalled(t, "Hash")
	})
}
