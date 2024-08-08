package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestLoginUseCase_DefaultPassword(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	mockChannel := make(chan entities.Event)
	eventBus.On("Subscribe", mock.Anything).Return((<-chan entities.Event)(mockChannel)).Once()
	useCase := liquidity_provider.NewLoginUseCase(lpRepository, eventBus)
	var credentials = &lpEntity.HashedCredentials{
		HashedUsername: test.AnyString,
		HashedPassword: test.AnyString,
		UsernameSalt:   test.AnyString,
		PasswordSalt:   test.AnyString,
	}

	assert.Nil(t, useCase.DefaultCredentials())
	useCase.SetDefaultCredentials(credentials)
	assert.Equal(t, credentials, useCase.DefaultCredentials())
}

func TestLoginUseCase_GetDefaultPasswordChannel(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := dataproviders.NewLocalEventBus()
	useCase := liquidity_provider.NewLoginUseCase(lpRepository, eventBus)
	eventBus.Publish(lpEntity.DefaultCredentialsSetEvent{
		Event:       entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
		Credentials: hashedDefaultCredentialsMock,
	})
	assert.NotNil(t, useCase.GetDefaultCredentialsChannel())
	select {
	case content := <-useCase.GetDefaultCredentialsChannel():
		assert.Equal(t, hashedDefaultCredentialsMock, content.(lpEntity.DefaultCredentialsSetEvent).Credentials)
	default:
		assert.Fail(t, "expected to receive an event")
	}
}

func TestLoginUseCase_Run_UseDefaultPassword(t *testing.T) {
	var waitGroup sync.WaitGroup
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Twice()
	eventBus := dataproviders.NewLocalEventBus()
	useCase := liquidity_provider.NewLoginUseCase(lpRepository, eventBus)
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup, bus entities.EventBus) {
		defer wg.Done()
		bus.Publish(lpEntity.DefaultCredentialsSetEvent{
			Event:       entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
			Credentials: hashedDefaultCredentialsMock,
		})
	}(&waitGroup, eventBus)
	waitGroup.Wait()
	err := useCase.Run(context.Background(), defaultCredentialsMock)
	require.NoError(t, err)
	err = useCase.Run(context.Background(), lpEntity.Credentials{
		Username: defaultCredentialsMock.Username,
		Password: defaultCredentialsMock.Password + "1234",
	})
	require.ErrorIs(t, err, liquidity_provider.BadLoginError)
	lpRepository.AssertExpectations(t)
}

func TestLoginUseCase_Run_UseStoredPassword(t *testing.T) {
	storedCredentials := &entities.Signed[lpEntity.HashedCredentials]{
		Value: lpEntity.HashedCredentials{
			HashedUsername: "d5e7cb7636083de780d8d32a7267b1aca58d27105c28462352e75dbc9b4aa938",
			HashedPassword: "b59ce56c879d1980ce8136c11c57b5a26a7d96cb30a3ba831805affdac142dcb",
			UsernameSalt:   "c009436ca9dbfc146dc3b5c47cb1937f95a28a5c55962721ca0851fff1dd7e17",
			PasswordSalt:   "b646012160b9b3dfdc35e2d0e65c741e49ca58d1f728f67923ecf6b5ecafbe08",
		},
	}
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Twice()
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).
		Return((<-chan entities.Event)(make(chan entities.Event))).Once()
	useCase := liquidity_provider.NewLoginUseCase(lpRepository, eventBus)

	err := useCase.Run(context.Background(), lpEntity.Credentials{
		Username: "fakeUser",
		Password: "MyFakeCredential1!",
	})
	require.NoError(t, err)
	err = useCase.Run(context.Background(), lpEntity.Credentials{
		Username: "otherFakeUser",
		Password: "wrong password",
	})
	require.ErrorIs(t, err, liquidity_provider.BadLoginError)
	eventBus.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}
