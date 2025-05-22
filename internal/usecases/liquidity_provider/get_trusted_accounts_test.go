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
func TestGetTrustedAccountsUseCase_Run(t *testing.T) {
	t.Run("should return accounts and nil when accounts are found and integrity check passes", func(t *testing.T) {
		mockHashBytes := []byte("mockhash12345678")
		mockHashHex := hex.EncodeToString(mockHashBytes)
		signedAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
			{
				Value: liquidity_provider.TrustedAccountDetails{
					Address:        "0x123",
					Name:           "Test Account 1",
					BtcLockingCap:  entities.NewWei(100),
					RbtcLockingCap: entities.NewWei(200),
				},
				Hash:      mockHashHex,
				Signature: "valid-signature-1",
			},
			{
				Value: liquidity_provider.TrustedAccountDetails{
					Address:        "0x456",
					Name:           "Test Account 2",
					BtcLockingCap:  entities.NewWei(300),
					RbtcLockingCap: entities.NewWei(400),
				},
				Hash:      mockHashHex,
				Signature: "valid-signature-2",
			},
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(signedAccounts, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(mockHashBytes)
		signerMock := &mocks.SignerMock{}
		signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background())
		require.NoError(t, err)
		assert.Equal(t, signedAccounts, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})
	t.Run("should return error when repository returns error", func(t *testing.T) {
		repoErr := assert.AnError
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(nil, repoErr)
		hashMock := &mocks.HashMock{}
		signerMock := &mocks.SignerMock{}
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background())
		require.Error(t, err)
		assert.Equal(t, repoErr, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertNotCalled(t, "Hash")
		signerMock.AssertNotCalled(t, "Validate")
	})
	t.Run("should return error when integrity check fails for any account", func(t *testing.T) {
		validHashBytes := []byte("valid-hash-bytes")
		validHashHex := hex.EncodeToString(validHashBytes)
		tamperedHashHex := hex.EncodeToString([]byte("tampered-hash-bytes"))
		signedAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
			{
				Value: liquidity_provider.TrustedAccountDetails{
					Address:        "0x123",
					Name:           "Test Account 1",
					BtcLockingCap:  entities.NewWei(100),
					RbtcLockingCap: entities.NewWei(200),
				},
				Hash:      validHashHex,
				Signature: "valid-signature",
			},
			{
				Value: liquidity_provider.TrustedAccountDetails{
					Address:        "0x456",
					Name:           "Test Account 2",
					BtcLockingCap:  entities.NewWei(300),
					RbtcLockingCap: entities.NewWei(400),
				},
				Hash:      tamperedHashHex,
				Signature: "signature",
			},
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetAllTrustedAccounts", mock.Anything).Return(signedAccounts, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(validHashBytes).Once()
		hashMock.On("Hash", mock.Anything).Return([]byte("different-hash-bytes")).Once()
		signerMock := &mocks.SignerMock{}
		signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)
		useCase := lpuc.NewGetTrustedAccountsUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background())
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.ErrTamperedTrustedAccount, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})
}
