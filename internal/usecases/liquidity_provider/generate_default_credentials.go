package liquidity_provider

import (
	"context"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

const defaultPasswordBytes = 64

type GenerateDefaultCredentialsUseCase struct {
	lpRepository liquidity_provider.LiquidityProviderRepository
	eventBus     entities.EventBus
}

func NewGenerateDefaultCredentialsUseCase(
	lpRepository liquidity_provider.LiquidityProviderRepository,
	eventBus entities.EventBus,
) *GenerateDefaultCredentialsUseCase {
	return &GenerateDefaultCredentialsUseCase{lpRepository: lpRepository, eventBus: eventBus}
}

func (useCase *GenerateDefaultCredentialsUseCase) Run(ctx context.Context, targetDir string) error {
	credentials, err := useCase.lpRepository.GetCredentials(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DefaultCredentialsId, err)
	}
	if credentials != nil {
		return nil
	}
	passwordBytes, err := utils.GetRandomBytes(defaultPasswordBytes)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DefaultCredentialsId, err)
	}
	stringPassword := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(passwordBytes)
	passwordFile := path.Join(targetDir, "management_password.txt")
	err = os.WriteFile(passwordFile, []byte(stringPassword), 0600)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DefaultCredentialsId, fmt.Errorf("error writing password file: %w", err))
	}

	hashedUsername, usernameSalt, err := utils.HashAndSaltArgon2(defaultUsername, credentialSaltSize)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DefaultCredentialsId, err)
	}
	hashedPassword, passwordSalt, err := utils.HashAndSaltArgon2(stringPassword, credentialSaltSize)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.DefaultCredentialsId, err)
	}

	useCase.eventBus.Publish(liquidity_provider.DefaultCredentialsSetEvent{
		Event: entities.NewBaseEvent(liquidity_provider.DefaultCredentialsSetEventId),
		Credentials: &liquidity_provider.HashedCredentials{
			HashedUsername: hex.EncodeToString(hashedUsername),
			HashedPassword: hex.EncodeToString(hashedPassword),
			UsernameSalt:   usernameSalt,
			PasswordSalt:   passwordSalt,
		},
	})
	log.Infof("There was no password detected in the database. A new password has been generated and saved in the file %s."+
		"Please keep this file safe. The first time you open the management interface, you will be asked to change this password.", passwordFile)
	return nil
}
