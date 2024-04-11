package liquidity_provider

import (
	"cmp"
	"context"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"slices"
)

const (
	defaultUsername = "admin"
)

var BadLoginError = errors.New("incorrect username or credentials")

func ValidateConfiguredProvider(
	provider liquidity_provider.LiquidityProvider,
	lbc blockchain.LiquidityBridgeContract,
) (uint64, error) {
	var err error
	var providers []liquidity_provider.RegisteredLiquidityProvider

	if providers, err = lbc.GetProviders(); err != nil {
		return 0, err
	}

	index, found := slices.BinarySearchFunc(
		providers,
		liquidity_provider.RegisteredLiquidityProvider{Address: provider.RskAddress()},
		func(a, b liquidity_provider.RegisteredLiquidityProvider) int {
			return cmp.Compare(a.Address, b.Address)
		},
	)
	if !found {
		return 0, usecases.ProviderConfigurationError
	}
	return providers[index].Id, nil
}

// DefaultPasswordProvider this is an interface to be implemented by those use case that require to use ValidateCredentials,
// since that function requires a way to access to the default password set by the application
type DefaultPasswordProvider interface {
	LiquidityProviderRepository() liquidity_provider.LiquidityProviderRepository
	GetDefaultPasswordChannel() <-chan entities.Event
	SetDefaultPassword(password string)
	DefaultPassword() string
}

func ValidateCredentials(
	ctx context.Context,
	useCase DefaultPasswordProvider,
	credentials liquidity_provider.Credentials,
) error {
	storedCredentials, err := useCase.LiquidityProviderRepository().GetCredentials(ctx)
	if err != nil {
		return err
	}
	if storedCredentials == nil {
		return validateDefaultCredentials(useCase, credentials)
	}

	usernameHash, err := utils.HashArgon2(credentials.Username, storedCredentials.Value.UsernameSalt)
	if err != nil {
		return err
	}
	passwordHash, err := utils.HashArgon2(credentials.Password, storedCredentials.Value.PasswordSalt)
	if err != nil {
		return err
	}
	err = compareCredentials(usernameHash, passwordHash, storedCredentials.Value)
	if err != nil {
		return err
	}
	return nil
}

func ReadDefaultPassword(useCase DefaultPasswordProvider) (string, error) {
	if useCase.DefaultPassword() != "" {
		return useCase.DefaultPassword(), nil
	}
	select {
	case event := <-useCase.GetDefaultPasswordChannel():
		parsedEvent, ok := event.(liquidity_provider.DefaultCredentialsSetEvent)
		if !ok {
			return "", errors.New("wrong event error")
		}
		useCase.SetDefaultPassword(parsedEvent.Password)
		return useCase.DefaultPassword(), nil
	default:
		return "", errors.New("default password not set")
	}
}

func validateDefaultCredentials(
	useCase DefaultPasswordProvider,
	credentials liquidity_provider.Credentials,
) error {
	var defaultPassword string
	var err error
	if defaultPassword, err = ReadDefaultPassword(useCase); err != nil {
		return err
	}
	if credentials.Username == defaultUsername && credentials.Password == defaultPassword {
		return nil
	} else {
		return BadLoginError
	}
}

func compareCredentials(
	usernameHash, passwordHash []byte,
	credentials liquidity_provider.HashedCredentials,
) error {
	usernameBytes, err := hex.DecodeString(credentials.HashedUsername)
	if err != nil {
		return err
	}
	passwordBytes, err := hex.DecodeString(credentials.HashedPassword)
	if err != nil {
		return err
	}
	usernameMatch := subtle.ConstantTimeCompare(usernameBytes, usernameHash) == 1
	passwordMatch := subtle.ConstantTimeCompare(passwordBytes, passwordHash) == 1
	if usernameMatch && passwordMatch {
		return nil
	} else {
		return BadLoginError
	}
}
