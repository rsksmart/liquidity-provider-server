package liquidity_provider_test

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
	const testAddr = "0x0000000000000000000000000000000000000123"
	t.Run("should return account and nil when account is found and integrity check passes", func(t *testing.T) {
		mockHashBytes := []byte("mockhash12345678")
		mockHashHex := hex.EncodeToString(mockHashBytes)
		signedAccount := &entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value: liquidity_provider.TrustedAccountDetails{
				Address:        testAddr,
				Name:           "Test Account",
				BtcLockingCap:  entities.NewWei(100),
				RbtcLockingCap: entities.NewWei(200),
			},
			Hash:      mockHashHex,
			Signature: "valid-signature",
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, testAddr).Return(signedAccount, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(mockHashBytes)
		signerMock := &mocks.SignerMock{}
		signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background(), testAddr)
		require.NoError(t, err)
		assert.Equal(t, signedAccount, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
		signerMock.AssertExpectations(t)
	})
	t.Run("should return error when account is not found", func(t *testing.T) {
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, testAddr).Return(nil, liquidity_provider.TrustedAccountNotFoundError)
		hashMock := &mocks.HashMock{}
		signerMock := &mocks.SignerMock{}
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background(), testAddr)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.TrustedAccountNotFoundError, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertNotCalled(t, "Hash")
		signerMock.AssertNotCalled(t, "Validate")
	})
	t.Run("should return error when integrity check fails", func(t *testing.T) {
		storedHashHex := hex.EncodeToString([]byte("stored-hash-value"))
		actualHashBytes := []byte("actual-hash-value")
		signedAccount := &entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value: liquidity_provider.TrustedAccountDetails{
				Address:        testAddr,
				Name:           "Test Account",
				BtcLockingCap:  entities.NewWei(100),
				RbtcLockingCap: entities.NewWei(200),
			},
			Hash:      storedHashHex,
			Signature: "signature",
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, testAddr).Return(signedAccount, nil)
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(actualHashBytes)
		signerMock := &mocks.SignerMock{}
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background(), testAddr)
		require.Error(t, err)
		assert.Equal(t, liquidity_provider.TamperedTrustedAccountError, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
		signerMock.AssertNotCalled(t, "Validate")
	})

	t.Run("lookup uses normalized address for mixed-case input", func(t *testing.T) {
		base := common.HexToAddress("0xabcdef00112233445566778899aabbccddeef00")
		mixedLookup := base.Hex()
		normalizedLookup := strings.ToLower(base.Hex())
		mockHashBytes := []byte("mockhash12345678")
		mockHashHex := hex.EncodeToString(mockHashBytes)
		signedAccount := &entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value: liquidity_provider.TrustedAccountDetails{
				Address:        normalizedLookup,
				Name:           "Test Account",
				BtcLockingCap:  entities.NewWei(100),
				RbtcLockingCap: entities.NewWei(200),
			},
			Hash:      mockHashHex,
			Signature: "valid-signature",
		}
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		repo.On("GetTrustedAccount", mock.Anything, normalizedLookup).Return(signedAccount, nil).Once()
		hashMock := &mocks.HashMock{}
		hashMock.On("Hash", mock.Anything).Return(mockHashBytes)
		signerMock := &mocks.SignerMock{}
		signerMock.On("Validate", mock.Anything, mock.Anything).Return(true)
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background(), mixedLookup)
		require.NoError(t, err)
		assert.Equal(t, signedAccount, result)
	})

	t.Run("invalid address", func(t *testing.T) {
		repo := mocks.NewTrustedAccountRepositoryMock(t)
		hashMock := &mocks.HashMock{}
		signerMock := &mocks.SignerMock{}
		useCase := lpuc.NewGetTrustedAccountUseCase(repo, hashMock.Hash, signerMock)
		result, err := useCase.Run(context.Background(), "invalid")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, liquidity_provider.InvalidTrustedAccountAddressError))
		repo.AssertNotCalled(t, "GetTrustedAccount")
	})
}
