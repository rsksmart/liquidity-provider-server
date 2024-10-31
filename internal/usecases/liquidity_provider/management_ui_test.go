package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen
func TestGetManagementUiDataUseCase_Run(t *testing.T) {
	const testUrl = "http://localhost:8080"
	t.Run("Return correct data when not logged in and credentials not set", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lbcMock := &mocks.LbcMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Lbc: lbcMock}, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementLoginTemplate, result.Name)
		assert.False(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Empty(t, result.Data.Configuration)
		assert.Empty(t, result.Data.ProviderData)
		assert.Empty(t, result.Data.BtcAddress)
		assert.Empty(t, result.Data.RskAddress)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "GetProviders")
	})
	t.Run("Return correct data when not logged in and credentials set", func(t *testing.T) {
		lpMock := &mocks.ProviderMock{}
		lbcMock := &mocks.LbcMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Lbc: lbcMock}, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementLoginTemplate, result.Name)
		assert.True(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Empty(t, result.Data.Configuration)
		assert.Empty(t, result.Data.ProviderData)
		assert.Empty(t, result.Data.BtcAddress)
		assert.Empty(t, result.Data.RskAddress)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		lbcMock.AssertNotCalled(t, "GetProviders")
	})
	t.Run("Return correct data when logged in", func(t *testing.T) {
		const (
			btcAddress = test.AnyAddress
			rskAddress = test.AnyHash
		)
		lbcMock := &mocks.LbcMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		fullConfig := liquidity_provider.FullConfiguration{
			General: lp.DefaultGeneralConfiguration(),
			Pegin:   lp.DefaultPeginConfiguration(),
			Pegout:  lp.DefaultPegoutConfiguration(),
		}
		providersInfo := []lp.RegisteredLiquidityProvider{
			{
				Id:           1,
				Address:      "otherAddress",
				Name:         "otherName",
				ApiBaseUrl:   "otherUrl",
				Status:       true,
				ProviderType: lp.PeginProvider,
			},
			{
				Id:           2,
				Address:      rskAddress,
				Name:         test.AnyString,
				ApiBaseUrl:   test.AnyUrl,
				Status:       true,
				ProviderType: lp.FullProvider,
			},
		}
		lbcMock.On("GetProviders").Return(providersInfo, nil).Once()
		lpMock.On("GeneralConfiguration", test.AnyCtx).Return(fullConfig.General).Once()
		lpMock.On("PeginConfiguration", test.AnyCtx).Return(fullConfig.Pegin).Once()
		lpMock.On("PegoutConfiguration", test.AnyCtx).Return(fullConfig.Pegout).Once()
		lpMock.On("BtcAddress").Return(btcAddress).Once()
		lpMock.On("RskAddress").Return(rskAddress).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Lbc: lbcMock}, testUrl)
		result, err := useCase.Run(context.Background(), true)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementUiTemplate, result.Name)
		assert.True(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Equal(t, fullConfig, result.Data.Configuration)
		assert.Equal(t, providersInfo[1], result.Data.ProviderData)
		assert.Equal(t, btcAddress, result.Data.BtcAddress)
		assert.Equal(t, rskAddress, result.Data.RskAddress)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		lbcMock.AssertExpectations(t)
	})
	t.Run("Return error when repository fails", func(t *testing.T) {
		lbcMock := &mocks.LbcMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Lbc: lbcMock}, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Return error when provider doesn't exists", func(t *testing.T) {
		const rskAddress = test.AnyHash

		lbcMock := &mocks.LbcMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		fullConfig := liquidity_provider.FullConfiguration{
			General: lp.DefaultGeneralConfiguration(),
			Pegin:   lp.DefaultPeginConfiguration(),
			Pegout:  lp.DefaultPegoutConfiguration(),
		}
		providersInfo := []lp.RegisteredLiquidityProvider{
			{
				Id:           1,
				Address:      "otherAddress",
				Name:         "otherName",
				ApiBaseUrl:   "otherUrl",
				Status:       true,
				ProviderType: lp.PeginProvider,
			},
			{
				Id:           2,
				Address:      rskAddress,
				Name:         test.AnyString,
				ApiBaseUrl:   test.AnyUrl,
				Status:       true,
				ProviderType: lp.FullProvider,
			},
		}
		lbcMock.On("GetProviders").Return(providersInfo, nil).Once()
		lpMock.On("GeneralConfiguration", test.AnyCtx).Return(fullConfig.General).Once()
		lpMock.On("PeginConfiguration", test.AnyCtx).Return(fullConfig.Pegin).Once()
		lpMock.On("PegoutConfiguration", test.AnyCtx).Return(fullConfig.Pegout).Once()
		lpMock.On("RskAddress").Return("nonExistingAddress").Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Lbc: lbcMock}, testUrl)
		result, err := useCase.Run(context.Background(), true)
		require.ErrorIs(t, err, usecases.ProviderNotFoundError)
		assert.Empty(t, result)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		lbcMock.AssertExpectations(t)
	})
}
