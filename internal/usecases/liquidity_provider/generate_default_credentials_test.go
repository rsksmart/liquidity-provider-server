package liquidity_provider_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestGenerateDefaultCredentialsUseCase_Run(t *testing.T) {
	var emittedEvent lpEntity.DefaultCredentialsSetEvent
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	lpRepository.On("GetCredentials", context.Background()).Return(nil, nil)
	eventBus.On("Publish", mock.MatchedBy(func(input entities.Event) bool {
		event, ok := input.(lpEntity.DefaultCredentialsSetEvent)
		require.True(t, ok)
		emittedEvent = event
		return event.Id() == lpEntity.DefaultCredentialsSetEventId && event.Password != ""
	})).Once()
	useCase := liquidity_provider.NewGenerateDefaultCredentialsUseCase(lpRepository, eventBus)
	dir := os.TempDir()
	buff := new(bytes.Buffer)
	log.SetOutput(buff)
	err := useCase.Run(context.Background(), dir)
	require.NoError(t, err)
	passwordFile := path.Join(dir, "management_password.txt")
	writtenPassword, err := os.ReadFile(passwordFile)
	require.NoError(t, err)
	assert.Equal(t, emittedEvent.Password, string(writtenPassword))
	expectedLog := fmt.Sprintf("There was no password detected in the database. A new password has "+
		"been generated and saved in the file %s.Please keep this file safe. The first time you open the "+
		"management interface, you will be asked to change this password.", passwordFile)
	assert.Contains(t, buff.String(), expectedLog)
	assert.Positive(t, len(writtenPassword))
	eventBus.AssertExpectations(t)
	lpRepository.AssertExpectations(t)
}

func TestGenerateDefaultCredentialsUseCase_Run_StoredCredentials(t *testing.T) {
	lpRepository := &mocks.LiquidityProviderRepositoryMock{}
	eventBus := &mocks.EventBusMock{}
	lpRepository.On("GetCredentials", context.Background()).Return(&entities.Signed[lpEntity.HashedCredentials]{
		Value: lpEntity.HashedCredentials{
			HashedUsername: test.AnyString,
			HashedPassword: test.AnyString,
			UsernameSalt:   test.AnyString,
			PasswordSalt:   test.AnyString,
		},
	}, nil)
	useCase := liquidity_provider.NewGenerateDefaultCredentialsUseCase(lpRepository, eventBus)
	dir := os.TempDir()
	err := useCase.Run(context.Background(), dir)
	require.NoError(t, err)
	eventBus.AssertNotCalled(t, "Publish")
	lpRepository.AssertExpectations(t)
}

func TestGenerateDefaultCredentialsUseCase_Run_ErrorHandling(t *testing.T) {
	t.Run("GetCredentials error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		lpRepository.On("GetCredentials", context.Background()).Return(nil, assert.AnError)
		useCase := liquidity_provider.NewGenerateDefaultCredentialsUseCase(lpRepository, eventBus)
		dir := os.TempDir()
		err := useCase.Run(context.Background(), dir)
		require.Error(t, err)
		eventBus.AssertNotCalled(t, "Publish")
		lpRepository.AssertExpectations(t)
	})
	t.Run("Write file error", func(t *testing.T) {
		lpRepository := &mocks.LiquidityProviderRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		lpRepository.On("GetCredentials", context.Background()).Return(nil, nil)
		useCase := liquidity_provider.NewGenerateDefaultCredentialsUseCase(lpRepository, eventBus)
		err := useCase.Run(context.Background(), "not a dir")
		require.ErrorContains(t, err, "error writing password file")
		eventBus.AssertNotCalled(t, "Publish")
		lpRepository.AssertExpectations(t)
	})
}
