package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

var storedCredentials = &entities.Signed[lpEntity.HashedCredentials]{
	Value: lpEntity.HashedCredentials{
		HashedUsername: "12ff30a29822669a598d9ad86afa00d48c5c25917c8c75cff1c4302051d7e9a5",
		HashedPassword: "100844abd86da73cd12ed3d7949267f24d3b6f81dd6f38f3bbb9d55d596d3e1e",
		UsernameSalt:   "f80a2b53ce41bc884f9d574f78ca58ab97d58f785b6a1d13356a64217afffb9b",
		PasswordSalt:   "eaa2bded32cc585d3f37c5319abe8890ad28a697ed66d5823f10536cc9c0fdb9",
	},
}

var oldCredentials = lpEntity.Credentials{Username: "oldUsername", Password: "oldPassword1*"}

func TestSetCredentialsUseCase_DefaultPassword(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	mockChannel := make(chan entities.Event)
	eventBus.On("Subscribe", mock.Anything).Return((<-chan entities.Event)(mockChannel)).Once()
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)
	credentials := &lpEntity.HashedCredentials{
		HashedUsername: test.AnyString,
		HashedPassword: test.AnyString,
		UsernameSalt:   test.AnyString,
		PasswordSalt:   test.AnyString,
	}
	assert.Empty(t, useCase.DefaultCredentials())
	useCase.SetDefaultCredentials(credentials)
	assert.Equal(t, credentials, useCase.DefaultCredentials())
}

func TestSetCredentialsUseCase_GetDefaultPasswordChannel(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := dataproviders.NewLocalEventBus()
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)
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

func TestSetCredentialsUseCase_Run_StoredCredentials(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil)
	lpRepository.On(
		"UpsertCredentials",
		test.AnyCtx,
		mock.MatchedBy(func(credentials entities.Signed[lpEntity.HashedCredentials]) bool {
			const expectedLength = 64
			// we cant assert the exact value because the salt is random
			return len(credentials.Value.HashedUsername) == expectedLength &&
				len(credentials.Value.HashedPassword) == expectedLength &&
				len(credentials.Value.UsernameSalt) == expectedLength &&
				len(credentials.Value.PasswordSalt) == expectedLength &&
				credentials.Signature == "04030201" &&
				credentials.Hash == "01020304"
		})).
		Return(nil).Once()
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return((<-chan entities.Event)(make(chan entities.Event))).Once()
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{0x01, 0x02, 0x03, 0x04}).Once()
	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)
	newCredentials := lpEntity.Credentials{Username: "newUsername", Password: "newPassword1*"}

	t.Run("Correct login", func(t *testing.T) {
		err := useCase.Run(context.Background(), oldCredentials, newCredentials)
		require.NoError(t, err)
		lpRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		walletMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
	t.Run("Incorrect login", func(t *testing.T) {
		lpRepository.Mock.Calls = []mock.Call{}
		incorrectCredentials := lpEntity.Credentials{Username: "oldUsernameIncorrect", Password: "oldPassword123!+"}
		err := useCase.Run(context.Background(), incorrectCredentials, newCredentials)
		require.ErrorIs(t, err, liquidity_provider.BadLoginError)
		lpRepository.AssertCalled(t, "GetCredentials", test.AnyCtx)
	})
}

func TestSetCredentialsUseCase_Run_DefaultCredentials(t *testing.T) {
	var waitGroup sync.WaitGroup
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
	lpRepository.On(
		"UpsertCredentials",
		test.AnyCtx,
		mock.MatchedBy(func(credentials entities.Signed[lpEntity.HashedCredentials]) bool {
			const expectedLength = 64
			// we cant assert the exact value because the salt is random
			return len(credentials.Value.HashedUsername) == expectedLength &&
				len(credentials.Value.HashedPassword) == expectedLength &&
				len(credentials.Value.UsernameSalt) == expectedLength &&
				len(credentials.Value.PasswordSalt) == expectedLength &&
				credentials.Signature == "04030201" &&
				credentials.Hash == "01020304"
		})).
		Return(nil).Once()
	eventBus := &mocks.EventBusMock{}
	defaultPasswordChannel := make(chan entities.Event, 1)
	eventBus.On("Subscribe", lpEntity.DefaultCredentialsSetEventId).Return((<-chan entities.Event)(defaultPasswordChannel)).Once()
	walletMock := &mocks.RskWalletMock{}
	walletMock.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil)
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{0x01, 0x02, 0x03, 0x04}).Once()

	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)
	newCredentials := lpEntity.Credentials{Username: "newUsername", Password: "newPassword1*"}
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup, eventChannel chan entities.Event) {
		defer wg.Done()
		eventChannel <- lpEntity.DefaultCredentialsSetEvent{
			Event:       entities.NewBaseEvent(lpEntity.DefaultCredentialsSetEventId),
			Credentials: hashedDefaultCredentialsMock,
		}
	}(&waitGroup, defaultPasswordChannel)
	waitGroup.Wait()

	t.Run("Correct login", func(t *testing.T) {
		err := useCase.Run(context.Background(), defaultCredentialsMock, newCredentials)
		require.NoError(t, err)
		lpRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		walletMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Incorrect login", func(t *testing.T) {
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, nil).Once()
		incorrectCredentials := lpEntity.Credentials{Username: "oldUsernameIncorrect", Password: defaultCredentialsMock.Password + "wrong"}
		err := useCase.Run(context.Background(), incorrectCredentials, newCredentials)
		require.ErrorIs(t, err, liquidity_provider.BadLoginError)
		lpRepository.AssertExpectations(t)
	})
}

func TestSetCredentialsUseCase_Run_InvalidPassword(t *testing.T) {
	passwords := []string{
		"short",
		"nouppercase!123",
		"NOLOWERCASE**456",
		"NoDigit!!#%qdasv",
		"NoSpecialChar123",
		"longeaa2bded32cc585d3f37c5319abe8890ad28a697ed66d5823f10536cc9c0fdb9eaa2bded32cc585d3f37c5319abe8890ad28a697ed66d5823f10536cc9c0fdb9",
	}

	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Times(len(passwords))
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return((<-chan entities.Event)(make(chan entities.Event))).Once()
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)

	for _, password := range passwords {
		newCredentials := lpEntity.Credentials{Username: "newUsername", Password: password}
		err := useCase.Run(context.Background(), oldCredentials, newCredentials)
		require.ErrorIs(t, err, utils.PasswordComplexityError)
	}
	hashMock.AssertNotCalled(t, "Hash")
	walletMock.AssertNotCalled(t, "SignBytes")
	lpRepository.AssertExpectations(t)
}

func TestSetCredentialsUseCase_Run_ErrorHandling(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	eventBus.On("Subscribe", mock.Anything).Return((<-chan entities.Event)(make(chan entities.Event))).Once()
	walletMock := &mocks.RskWalletMock{}
	hashMock := &mocks.HashMock{}
	hashMock.On("Hash", mock.Anything).Return([]byte{0x01, 0x02, 0x03, 0x04})
	useCase := liquidity_provider.NewSetCredentialsUseCase(lpRepository, walletMock, hashMock.Hash, eventBus)
	newCredentials := lpEntity.Credentials{Username: test.AnyString, Password: test.AnyString + "A1*"}

	t.Run("GetCredentials error", func(t *testing.T) {
		lpRepository.On("GetCredentials", test.AnyCtx).Return(nil, assert.AnError).Once()
		err := useCase.Run(context.Background(), oldCredentials, newCredentials)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		walletMock.AssertNotCalled(t, "SignBytes")
		hashMock.AssertNotCalled(t, "Hash")
	})

	t.Run("UpsertCredentials error", func(t *testing.T) {
		walletMock.On("SignBytes", mock.Anything).Return([]byte{4, 3, 2, 1}, nil).Once()
		lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
		lpRepository.On("UpsertCredentials", test.AnyCtx, mock.Anything).Return(assert.AnError).Once()
		err := useCase.Run(context.Background(), oldCredentials, newCredentials)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		walletMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("SignBytes error", func(t *testing.T) {
		lpRepository.On("GetCredentials", test.AnyCtx).Return(storedCredentials, nil).Once()
		walletMock.On("SignBytes", mock.Anything).Return(nil, assert.AnError).Once()
		err := useCase.Run(context.Background(), oldCredentials, newCredentials)
		require.Error(t, err)
		lpRepository.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		walletMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
