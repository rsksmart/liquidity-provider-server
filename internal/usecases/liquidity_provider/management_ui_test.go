package liquidity_provider_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestGetManagementUiDataUseCase_Run(t *testing.T) {
	const testUrl = "http://localhost:8080"
	t.Run("Return correct data when not logged in and credentials not set", func(t *testing.T) {
		discovery := &mocks.DiscoveryContractMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		coldWallet := &mocks.ColdWalletMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Discovery: discovery}, coldWallet, testUrl)
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
		discovery.AssertNotCalled(t, "GetProvider")
	})
	t.Run("Return correct data when not logged in and credentials set", func(t *testing.T) {
		discovery := &mocks.DiscoveryContractMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		coldWallet := &mocks.ColdWalletMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Discovery: discovery}, coldWallet, testUrl)
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
		discovery.AssertNotCalled(t, "GetProvider")
	})
	t.Run("Return correct data when logged in", func(t *testing.T) {
		const (
			btcAddress = test.AnyAddress
			rskAddress = test.AnyHash
		)
		discovery := &mocks.DiscoveryContractMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		coldWallet := &mocks.ColdWalletMock{}
		fullConfig := liquidity_provider.FullConfiguration{
			General: lp.DefaultGeneralConfiguration(),
			Pegin:   lp.DefaultPeginConfiguration(),
			Pegout:  lp.DefaultPegoutConfiguration(),
		}
		lpInfo := lp.RegisteredLiquidityProvider{
			Id:           1,
			Address:      rskAddress,
			Name:         test.AnyString,
			ApiBaseUrl:   test.AnyUrl,
			Status:       true,
			ProviderType: lp.FullProvider,
		}
		discovery.On("GetProvider", rskAddress).Return(lpInfo, nil)
		lpMock.On("GeneralConfiguration", test.AnyCtx).Return(fullConfig.General).Once()
		lpMock.On("PeginConfiguration", test.AnyCtx).Return(fullConfig.Pegin).Once()
		lpMock.On("PegoutConfiguration", test.AnyCtx).Return(fullConfig.Pegout).Once()
		lpMock.On("BtcAddress").Return(btcAddress).Once()
		lpMock.On("RskAddress").Return(rskAddress).Once()
		coldWallet.EXPECT().GetRskAddress().Return(test.AnyRskAddress)
		coldWallet.EXPECT().GetBtcAddress().Return(test.AnyBtcAddress)
		coldWallet.EXPECT().GetLabel().Return("Address")
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Discovery: discovery}, coldWallet, testUrl)
		result, err := useCase.Run(context.Background(), true)
		require.NoError(t, err)
		assert.Equal(t, liquidity_provider.ManagementUiTemplate, result.Name)
		assert.True(t, result.Data.CredentialsSet)
		assert.Equal(t, testUrl, result.Data.BaseUrl)
		assert.Equal(t, fullConfig, result.Data.Configuration)
		assert.Equal(t, lpInfo, result.Data.ProviderData)
		assert.Equal(t, btcAddress, result.Data.BtcAddress)
		assert.Equal(t, rskAddress, result.Data.RskAddress)
		assert.Equal(t, test.AnyRskAddress, result.Data.ColdWallet.RskAddress)
		assert.Equal(t, test.AnyBtcAddress, result.Data.ColdWallet.BtcAddress)
		assert.Equal(t, "Address", result.Data.ColdWallet.Label)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		discovery.AssertExpectations(t)
		coldWallet.AssertExpectations(t)
	})
	t.Run("Return error when repository fails", func(t *testing.T) {
		discovery := &mocks.DiscoveryContractMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		coldWallet := &mocks.ColdWalletMock{}
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Discovery: discovery}, coldWallet, testUrl)
		result, err := useCase.Run(context.Background(), false)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("Return error when provider doesn't exists", func(t *testing.T) {
		discovery := &mocks.DiscoveryContractMock{}
		coldWallet := &mocks.ColdWalletMock{}
		lpMock := &mocks.ProviderMock{}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		fullConfig := liquidity_provider.FullConfiguration{
			General: lp.DefaultGeneralConfiguration(),
			Pegin:   lp.DefaultPeginConfiguration(),
			Pegout:  lp.DefaultPegoutConfiguration(),
		}
		discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{}, assert.AnError).Once()
		lpMock.On("GeneralConfiguration", test.AnyCtx).Return(fullConfig.General).Once()
		lpMock.On("PeginConfiguration", test.AnyCtx).Return(fullConfig.Pegin).Once()
		lpMock.On("PegoutConfiguration", test.AnyCtx).Return(fullConfig.Pegout).Once()
		lpMock.On("RskAddress").Return("nonExistingAddress").Once()
		useCase := liquidity_provider.NewGetManagementUiDataUseCase(lpRepository, lpMock, lpMock, lpMock, blockchain.RskContracts{Discovery: discovery}, coldWallet, testUrl)
		result, err := useCase.Run(context.Background(), true)
		require.Error(t, err)
		assert.Empty(t, result)
		lpRepository.AssertExpectations(t)
		lpMock.AssertExpectations(t)
		discovery.AssertExpectations(t)
	})
}
