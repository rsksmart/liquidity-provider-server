package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var defaultCredentialsMock = lpEntity.Credentials{Username: "admin", Password: "a default password"}
var hashedDefaultCredentialsMock = &lpEntity.HashedCredentials{
	HashedUsername: "e58faef24d13f93d99cb3c68e381d05d1e131029f90cbb1469ff99aa8b2ca8c2",
	HashedPassword: "20cdf83f3e87da259cb72609bdcaa220d9a4c69135ce9d8dd39eb9dd738ee503",
	UsernameSalt:   "4948388a01e926807fd86a5f1c2426dba97030717001f5f9d7106950e03724b2",
	PasswordSalt:   "9baf3a40312f39849f46dad1040f2f039f1cffa1238c41e9db675315cfad39b6",
}

func TestReadDefaultPassword_AlreadyRead(t *testing.T) {
	credentialsProvider := &mocks.DefaultCredentialsProviderMock{}
	credentials := &lpEntity.HashedCredentials{
		HashedUsername: test.AnyString,
		HashedPassword: test.AnyString,
		UsernameSalt:   test.AnyString,
		PasswordSalt:   test.AnyString,
	}
	credentialsProvider.On("DefaultCredentials").Return(credentials)
	defaultCredentials, err := liquidity_provider.ReadDefaultCredentials(credentialsProvider)
	assert.Equal(t, *credentials, defaultCredentials)
	require.NoError(t, err)
	credentialsProvider.AssertExpectations(t)
}

func TestReadDefaultPassword_NotRead(t *testing.T) {
	passwordChannel := make(chan entities.Event, 1)
	passwordChannel <- lpEntity.DefaultCredentialsSetEvent{
		Event:       entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
		Credentials: hashedDefaultCredentialsMock,
	}
	defaultCredentialsProvider := &mocks.DefaultCredentialsProviderMock{}
	defaultCredentialsProvider.On("DefaultCredentials").Return(nil).Once()
	defaultCredentialsProvider.On("GetDefaultCredentialsChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
	defaultCredentialsProvider.On("SetDefaultCredentials", hashedDefaultCredentialsMock).Return().Once()
	defaultCredentialsProvider.On("DefaultCredentials").Return(hashedDefaultCredentialsMock).Once()
	defaultCredentials, err := liquidity_provider.ReadDefaultCredentials(defaultCredentialsProvider)
	assert.Equal(t, *hashedDefaultCredentialsMock, defaultCredentials)
	require.NoError(t, err)
	defaultCredentialsProvider.AssertExpectations(t)
}

func TestReadDefaultPassword_ErroHandling(t *testing.T) {
	t.Run("Default passsword not set", func(t *testing.T) {
		passwordChannel := make(chan entities.Event, 1)
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
		passwordProvider.On("DefaultCredentials").Return(nil).Once()
		passwordProvider.On("GetDefaultCredentialsChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
		password, err := liquidity_provider.ReadDefaultCredentials(passwordProvider)
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
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
		passwordProvider.On("DefaultCredentials").Return(nil).Once()
		passwordProvider.On("GetDefaultCredentialsChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
		password, err := liquidity_provider.ReadDefaultCredentials(passwordProvider)
		assert.Empty(t, password)
		require.ErrorContains(t, err, "wrong event error")
		passwordProvider.AssertExpectations(t)
	})
}

func TestValidateCredentials_DefaultCredentials(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	passwordProvider := &mocks.DefaultCredentialsProviderMock{}
	useDefaultPasswordSetUp(lpRepository, passwordProvider)
	err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, defaultCredentialsMock)
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
	passwordProvider := &mocks.DefaultCredentialsProviderMock{}
	useStoredCredentialsSetUp(lpRepository, passwordProvider)
	err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
	require.NoError(t, err)
	passwordProvider.AssertNotCalled(t, "DefaultCredentials")
	passwordProvider.AssertNotCalled(t, "SetDefaultCredentials")
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
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
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
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
		useStoredCredentialsSetUp(lpRepository, passwordProvider)
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, credentials)
		require.ErrorIs(t, err, liquidity_provider.BadLoginError)
		passwordProvider.AssertNotCalled(t, "DefaultCredentials")
		passwordProvider.AssertNotCalled(t, "SetDefaultCredentials")
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
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
		err := liquidity_provider.ValidateCredentials(context.Background(), passwordProvider, lpEntity.Credentials{})
		require.Error(t, err)
		passwordProvider.AssertExpectations(t)
		lpRepository.AssertExpectations(t)
	})
	t.Run("Default password not set error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		passwordProvider := &mocks.DefaultCredentialsProviderMock{}
		passwordChannel := make(chan entities.Event, 1)
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
		passwordProvider.On("DefaultCredentials").Return(nil).Once()
		passwordProvider.On("GetDefaultCredentialsChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
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
			passwordProvider := &mocks.DefaultCredentialsProviderMock{}
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
	passwordProvider *mocks.DefaultCredentialsProviderMock,
) {
	lpRepository.On("GetCredentials", context.Background()).Return(nil, nil).Once()
	passwordChannel := make(chan entities.Event, 1)
	passwordChannel <- lpEntity.DefaultCredentialsSetEvent{
		Event:       entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
		Credentials: hashedDefaultCredentialsMock,
	}
	passwordProvider.On("DefaultCredentials").Return(nil).Once()
	passwordProvider.On("GetDefaultCredentialsChannel").Return((<-chan entities.Event)(passwordChannel)).Once()
	passwordProvider.On("SetDefaultCredentials", hashedDefaultCredentialsMock).Return().Once()
	passwordProvider.On("DefaultCredentials").Return(hashedDefaultCredentialsMock).Once()
	passwordProvider.On("LiquidityProviderRepository").Return(lpRepository).Once()
}

func useStoredCredentialsSetUp(
	lpRepository *mocks.LiquidityProviderRepositoryMock,
	passwordProvider *mocks.DefaultCredentialsProviderMock,
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
