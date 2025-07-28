package liquidity_provider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateTrustedAccountUseCase_Run(t *testing.T) { //nolint:funlen
	t.Run("Success case", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:        "0x123456",
			Name:           "Test Account",
			BtcLockingCap:  entities.NewWei(1000),
			RbtcLockingCap: entities.NewWei(1000),
		}
		expectedSignedAccount := entities.Signed[liquidity_provider.TrustedAccountDetails]{
			Value:     account,
			Signature: "04030201",
			Hash:      "01020304",
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.MatchedBy(func(a entities.Signed[liquidity_provider.TrustedAccountDetails]) bool {
			return a.Value.Address == account.Address &&
				a.Value.Name == account.Name &&
				a.Signature == expectedSignedAccount.Signature &&
				a.Hash == expectedSignedAccount.Hash
		})).Return(nil)
		useCase := lp.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.NoError(t, err)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("SignConfiguration error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:        "0x123456",
			Name:           "Test Account",
			BtcLockingCap:  entities.NewWei(1000),
			RbtcLockingCap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return(nil, errors.New("signing error"))
		useCase := lp.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.Error(t, err)
		repo.AssertNotCalled(t, "UpdateTrustedAccount")
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("UpdateTrustedAccount error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:        "0x123456",
			Name:           "Test Account",
			BtcLockingCap:  entities.NewWei(1000),
			RbtcLockingCap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.MatchedBy(func(a entities.Signed[liquidity_provider.TrustedAccountDetails]) bool {
			return a.Value.Address == account.Address && a.Value.Name == account.Name
		})).Return(assert.AnError)
		useCase := lp.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.Error(t, err)
		require.Contains(t, err.Error(), usecases.UpdateTrustedAccountId)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("Account not found error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:        "0x123456",
			Name:           "Test Account",
			BtcLockingCap:  entities.NewWei(1000),
			RbtcLockingCap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, mock.MatchedBy(func(a entities.Signed[liquidity_provider.TrustedAccountDetails]) bool {
			return a.Value.Address == account.Address && a.Value.Name == account.Name
		})).Return(liquidity_provider.TrustedAccountNotFoundError)
		useCase := lp.NewUpdateTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.Error(t, err)
		require.Contains(t, err.Error(), liquidity_provider.TrustedAccountNotFoundError.Error())
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
