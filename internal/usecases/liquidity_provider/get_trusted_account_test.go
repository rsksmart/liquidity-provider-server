package liquidity_provider_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lpuc "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestGetTrustedAccountUseCase_Run(t *testing.T) {
	t.Run("should return account and nil when account is found and integrity check passes", func(t *testing.T) {
		mockHashBytes := []byte("mockhash12345678")
		mockHashHex := hex.EncodeToString(mockHashBytes)
		signedAccount := &entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value: liquidity_provider.TrustedAccountDetails{
				Address:          "0x123",
				Name:             "Test Account",
				Btc_locking_cap:  entities.NewWei(100),
				Rbtc_locking_cap: entities.NewWei(200),
			},
			Hash:      mockHashHex,
			Signature: "valid-signature",
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, "0x123").Return(signedAccount, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(mockHashBytes)
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash)
		result, err := useCase.Run(context.Background(), "0x123")
		require.NoError(t, err)
		assert.Equal(t, signedAccount, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("should return error when account is not found", func(t *testing.T) {
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, "0x123").Return(nil, liquidity_provider.ErrTrustedAccountNotFound)
		hashMock := &mocks.HashMock{}
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash)
		result, err := useCase.Run(context.Background(), "0x123")
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTrustedAccountNotFound, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertNotCalled(t, "Hash")
	})
	t.Run("should return error when integrity check fails", func(t *testing.T) {
		storedHashHex := hex.EncodeToString([]byte("stored-hash-value"))
		actualHashBytes := []byte("actual-hash-value")
		signedAccount := &entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value: liquidity_provider.TrustedAccountDetails{
				Address:          "0x123",
				Name:             "Test Account",
				Btc_locking_cap:  entities.NewWei(100),
				Rbtc_locking_cap: entities.NewWei(200),
			},
			Hash:      storedHashHex,
			Signature: "signature",
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, "0x123").Return(signedAccount, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(actualHashBytes)
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash)
		result, err := useCase.Run(context.Background(), "0x123")
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTamperedTrustedAccount, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
