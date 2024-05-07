package liquidity_provider_test

import (
	"context"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetManagementUiDataUseCase_Run(t *testing.T) {
	const testUrl = "http://localhost:8080"
	t.Run("Return correct data when not logged in and credentials not set", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementLoginTemplate, result.Name)
		assert.False(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Empty(t, result.Data.Configuration)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
	})
	t.Run("Return correct data when not logged in and credentials set", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementLoginTemplate, result.Name)
		assert.True(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Empty(t, result.Data.Configuration)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
	})
	t.Run("Return correct data when logged in", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		fullConfig := liquidity_provider.FullConfiguration{
			General: lp.DefaultGeneralConfiguration(),
			Pegin:   lp.DefaultPeginConfiguration(),
			Pegout:  lp.DefaultPegoutConfiguration(),
		}
		lpMock.On("GeneralConfiguration", test.AnyCtx).Return(fullConfig.General).Once()
		lpMock.On("PeginConfiguration", test.AnyCtx).Return(fullConfig.Pegin).Once()
		lpMock.On("PegoutConfiguration", test.AnyCtx).Return(fullConfig.Pegout).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, testUrl)
		result, err := useCase.Run(context.Background(), true)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementUiTemplate, result.Name)
		assert.True(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Equal(t, fullConfig, result.Data.Configuration)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
	})
	t.Run("Return error when repository fails", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}
