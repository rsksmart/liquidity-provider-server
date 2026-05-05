package liquidity_provider_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	lp "github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteTrustedAccountUseCase_Run(t *testing.T) {
	const addr = "0x1234567890123456789012345678901234567890"
	t.Run("Success case", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("DeleteTrustedAccount", mock.Anything, addr).Return(nil)
		useCase := lp.NewDeleteTrustedAccountUseCase(repo)
		err := useCase.Run(context.Background(), addr)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Not found error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		repo.On("DeleteTrustedAccount", mock.Anything, addr).Return(liquidity_provider.TrustedAccountNotFoundError)
		useCase := lp.NewDeleteTrustedAccountUseCase(repo)
		err := useCase.Run(context.Background(), addr)
		require.Error(t, err)
		require.Contains(t, err.Error(), liquidity_provider.TrustedAccountNotFoundError.Error())
		repo.AssertExpectations(t)
	})

	t.Run("invalid address", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		useCase := lp.NewDeleteTrustedAccountUseCase(repo)
		err := useCase.Run(context.Background(), "bad-address")
		require.Error(t, err)
		require.ErrorIs(t, err, liquidity_provider.InvalidTrustedAccountAddressError)
		repo.AssertNotCalled(t, "DeleteTrustedAccount")
	})
}
