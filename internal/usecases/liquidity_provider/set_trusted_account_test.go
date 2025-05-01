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

func TestSetTrustedAccountUseCase_Run(t *testing.T) { //nolint:funlen
	t.Run("Success case", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:          "0x123456",
			Name:             "Test Account",
			Btc_locking_cap:  entities.NewWei(1000),
			Rbtc_locking_cap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, account).Return(nil)
		useCase := lp.NewSetTrustedAccountUseCase(repo, signer, hashMock.Hash)
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
			Address:          "0x123456",
			Name:             "Test Account",
			Btc_locking_cap:  entities.NewWei(1000),
			Rbtc_locking_cap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return(nil, errors.New("signing error"))
		useCase := lp.NewSetTrustedAccountUseCase(repo, signer, hashMock.Hash)
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
			Address:          "0x123456",
			Name:             "Test Account",
			Btc_locking_cap:  entities.NewWei(1000),
			Rbtc_locking_cap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, account).Return(assert.AnError)
		useCase := lp.NewSetTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.Error(t, err)
		require.Contains(t, err.Error(), usecases.SetTrustedAccountId)
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Account not found error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		signer := &mocks.TransactionSignerMock{}
		hashMock := &mocks.HashMock{}
		account := liquidity_provider.TrustedAccountDetails{
			Address:          "0x123456",
			Name:             "Test Account",
			Btc_locking_cap:  entities.NewWei(1000),
			Rbtc_locking_cap: entities.NewWei(1000),
		}
		hashMock.On("Hash", mock.Anything).Return([]byte{1, 2, 3, 4})
		signer.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
		repo.On("UpdateTrustedAccount", mock.Anything, account).Return(liquidity_provider.ErrTrustedAccountNotFound)
		useCase := lp.NewSetTrustedAccountUseCase(repo, signer, hashMock.Hash)
		err := useCase.Run(context.Background(), account)
		require.Error(t, err)
		require.Contains(t, err.Error(), liquidity_provider.ErrTrustedAccountNotFound.Error())
		repo.AssertExpectations(t)
		signer.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
