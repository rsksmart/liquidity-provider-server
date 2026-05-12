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
	t.Run("Success case", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		address := "0x123456"
		repo.On("DeleteTrustedAccount", mock.Anything, address).Return(nil)
		useCase := lp.NewDeleteTrustedAccountUseCase(repo)
		err := useCase.Run(context.Background(), address)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Not found error", func(t *testing.T) {
		repo := &mocks.TrustedAccountRepositoryMock{}
		address := "0x123456"
		repo.On("DeleteTrustedAccount", mock.Anything, address).Return(liquidity_provider.TrustedAccountNotFoundError)
		useCase := lp.NewDeleteTrustedAccountUseCase(repo)
		err := useCase.Run(context.Background(), address)
		require.Error(t, err)
		require.Contains(t, err.Error(), liquidity_provider.TrustedAccountNotFoundError.Error())
		repo.AssertExpectations(t)
	})
}
