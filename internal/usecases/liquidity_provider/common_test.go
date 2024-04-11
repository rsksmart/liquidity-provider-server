package liquidity_provider_test

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// not a password, only a random hash
// nolint:gosec
const mockDefaultPassword = "2071bae7f92c6272f614e40d57272518480c48fb4c7a4c39525fdc14e6c97c1d"

func TestValidateConfiguredProvider(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{
		{
			Id:           1,
			Address:      "0x01",
			Name:         "one",
			ApiBaseUrl:   "api1.com",
			Status:       true,
			ProviderType: "both",
		},
		{
			Id:           2,
			Address:      "0x02",
			Name:         "two",
			ApiBaseUrl:   "api2.com",
			Status:       true,
			ProviderType: "pegin",
		},
		{
			Id:           3,
			Address:      "0x03",
			Name:         "three",
			ApiBaseUrl:   "api3.com",
			Status:       true,
			ProviderType: "pegout",
		},
	}, nil)

	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("0x02")

	id, err := liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(2), id)
	require.NoError(t, err)
}

func TestValidateConfiguredProvider_Fail(t *testing.T) {
	lbc := &mocks.LbcMock{}
	var provider *mocks.ProviderMock = nil
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{}, errors.New("some error")).Once()
	id, err := liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(0), id)
	require.Error(t, err)

	provider = &mocks.ProviderMock{}
	provider.On("RskAddress").Return("0x02")
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{
		{
			Id:           3,
			Address:      "0x03",
			Name:         "three",
			ApiBaseUrl:   "api3.com",
			Status:       true,
			ProviderType: "pegout",
		},
	}, nil).Once()
	id, err = liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(0), id)
	require.ErrorIs(t, err, usecases.ProviderConfigurationError)
}

func TestReadDefaultPassword_AlreadyRead(t *testing.T) {
	passwordProvider := &mocks.DefaultPasswordProviderMock{}
	passwordProvider.On("DefaultPassword").Return(test.AnyString)
	password, err := liquidity_provider.ReadDefaultPassword(passwordProvider)
	assert.Equal(t, test.AnyString, password)
	require.NoError(t, err)
	passwordProvider.AssertExpectations(t)
}

func TestReadDefaultPassword_NotRead(t *testing.T) {
	passwordChannel := make(chan entities.Event, 1)
	passwordChannel <- lpEntity.DefaultCredentialsSetEvent{
		Event:    entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
		Password: mockDefaultPassword,
	}
	passwordProvider := &mocks.DefaultPasswordProviderMock{}
	passwordProvider.On("DefaultPassword").Return("").Once()
	passwordProvider.On("GetDefaultPasswordChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
	passwordProvider.On("SetDefaultPassword", mockDefaultPassword).Return().Once()
	passwordProvider.On("DefaultPassword").Return(mockDefaultPassword).Once()
	password, err := liquidity_provider.ReadDefaultPassword(passwordProvider)
	assert.Equal(t, mockDefaultPassword, password)
	require.NoError(t, err)
	passwordProvider.AssertExpectations(t)
}

func TestReadDefaultPassword_ErroHandling(t *testing.T) {
	t.Run("Default passsword not set", func(t *testing.T) {
		passwordChannel := make(chan entities.Event, 1)
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		passwordProvider.On("DefaultPassword").Return("").Once()
		passwordProvider.On("GetDefaultPasswordChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
		password, err := liquidity_provider.ReadDefaultPassword(passwordProvider)
		assert.Empty(t, password)
		require.ErrorContains(t, err, "default password not set")
		passwordProvider.AssertExpectations(t)
	})
	t.Run("Wrong event", func(t *testing.T) {
		passwordChannel := make(chan entities.Event, 1)
		passwordChannel <- quote.AcceptedPeginQuoteEvent{
			Event:         entities.NewBaseEvent(quote.AcceptedPeginQuoteEventId),
			Quote:         quote.PeginQuote{},
			RetainedQuote: quote.RetainedPeginQuote{},
		}
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		passwordProvider.On("DefaultPassword").Return("").Once()
		passwordProvider.On("GetDefaultPasswordChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
		password, err := liquidity_provider.ReadDefaultPassword(passwordProvider)
		assert.Empty(t, password)
		require.ErrorContains(t, err, "wrong event error")
		passwordProvider.AssertExpectations(t)
	})
}

func TestValidateCredentials_DefaultCredentials(t *testing.T) {
	credentials := lpEntity.Credentials{
		Username: "admin",
		Password: mockDefaultPassword,
	}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	passwordProvider := &mocks.DefaultPasswordProviderMock{}
	useDefaultPasswordSetUp(lpRepository, passwordProvider)
	err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
	require.NoError(t, err)
	passwordProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestValidateCredentials_DefaultCredentials_StoredCredentials(t *testing.T) {
	credentials := lpEntity.Credentials{
		Username: "fakeUser",
		Password: "MyFakeCredential1!",
	}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	passwordProvider := &mocks.DefaultPasswordProviderMock{}
	useStoredCredentialsSetUp(lpRepository, passwordProvider)
	err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
	require.NoError(t, err)
	passwordProvider.AssertNotCalled(t, "DefaultPassword")
	passwordProvider.AssertNotCalled(t, "SetDefaultPassword")
	passwordProvider.AssertNotCalled(t, "DefaultPasswordChannel")
	passwordProvider.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestValidateCredentials_Badlogin(t *testing.T) {
	t.Run("Default password bad login", func(t *testing.T) {
		credentials := lpEntity.Credentials{
			Username: "admin",
			Password: "wrong password",
		}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		useDefaultPasswordSetUp(lpRepository, passwordProvider)
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
		require.ErrorIs(t, err, liquidity_provider.BadLoginError)
		passwordProvider.AssertExpectations(t)
		lpRepository.AssertExpectations(t)
	})
	t.Run("Stored password bad login", func(t *testing.T) {
		credentials := lpEntity.Credentials{
			Username: "any user",
			Password: "other wrong password",
		}
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		useStoredCredentialsSetUp(lpRepository, passwordProvider)
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
		require.ErrorIs(t, err, liquidity_provider.BadLoginError)
		passwordProvider.AssertNotCalled(t, "DefaultPassword")
		passwordProvider.AssertNotCalled(t, "SetDefaultPassword")
		passwordProvider.AssertNotCalled(t, "DefaultPasswordChannel")
		passwordProvider.AssertExpectations(t)
		lpRepository.AssertExpectations(t)
	})
}

func TestValidateCredentials_ErrorHandling(t *testing.T) {
	const (
		noHexString = "no hex"
		hexString   = "adefaa4dfb9723b813f5137850140bf2f23b8289f75ed532c75fc2a2ae51c082"
	)
	credentials := lpEntity.Credentials{Username: test.AnyString, Password: test.AnyString}
	t.Run("GetCredentials error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, lpEntity.Credentials{})
		require.Error(t, err)
		passwordProvider.AssertExpectations(t)
		lpRepository.AssertExpectations(t)
	})
	t.Run("Default password not set error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		passwordProvider := &mocks.DefaultPasswordProviderMock{}
		passwordChannel := make(chan entities.Event, 1)
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
		passwordProvider.On("DefaultPassword").Return("").Once()
		passwordProvider.On("GetDefaultPasswordChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, lpEntity.Credentials{})
		require.Error(t, err)
		passwordProvider.AssertExpectations(t)
		lpRepository.AssertExpectations(t)
	})
	errorCredentials := []lpEntity.HashedCredentials{
		{HashedUsername: noHexString, HashedPassword: hexString, UsernameSalt: hexString, PasswordSalt: hexString},
		{HashedUsername: hexString, HashedPassword: noHexString, UsernameSalt: hexString, PasswordSalt: hexString},
		{HashedUsername: hexString, HashedPassword: hexString, UsernameSalt: noHexString, PasswordSalt: hexString},
		{HashedUsername: hexString, HashedPassword: hexString, UsernameSalt: hexString, PasswordSalt: noHexString},
	}
	for _, errorCredential := range errorCredentials {
		t.Run("HashArgon2 error", func(t *testing.T) {
			storedCredentials := &entities.Signed[lpEntity.HashedCredentials]{Value: errorCredential}
			lpRepository := &mocks.LiquidityProviderRepositoryMock{}
			passwordProvider := &mocks.DefaultPasswordProviderMock{}
			lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
			passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
			err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
			require.Error(t, err)
			passwordProvider.AssertExpectations(t)
			lpRepository.AssertExpectations(t)
		})
	}
}

func useDefaultPasswordSetUp(
	lpRepository *mocks.LiquidityProviderRepositoryMock,
	passwordProvider *mocks.DefaultPasswordProviderMock,
) {
	lpRepository.On("GetCredentials", context.Background()).Return(nil, nil).Once()
	passwordChannel := make(chan entities.Event, 1)
	passwordChannel <- lpEntity.DefaultCredentialsSetEvent{
		Event:    entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
		Password: mockDefaultPassword,
	}
	passwordProvider.On("DefaultPassword").Return("").Once()
	passwordProvider.On("GetDefaultPasswordChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
	passwordProvider.On("SetDefaultPassword", mockDefaultPassword).Return().Once()
	passwordProvider.On("DefaultPassword").Return(mockDefaultPassword).Once()
	passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
}

func useStoredCredentialsSetUp(
	lpRepository *mocks.LiquidityProviderRepositoryMock,
	passwordProvider *mocks.DefaultPasswordProviderMock,
) {
	storedCredentials := &entities.Signed[lpEntity.HashedCredentials]{
		Value: lpEntity.HashedCredentials{
			HashedUsername: "d5e7cb7636083de780d8d32a7267b1aca58d27105c28462352e75dbc9b4aa938",
			HashedPassword: "b59ce56c879d1980ce8136c11c57b5a26a7d96cb30a3ba831805affdac142dcb",
			UsernameSalt:   "c009436ca9dbfc146dc3b5c47cb1937f95a28a5c55962721ca0851fff1dd7e17",
			PasswordSalt:   "b646012160b9b3dfdc35e2d0e65c741e49ca58d1f728f67923ecf6b5ecafbe08",
		},
	}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
	passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
}
