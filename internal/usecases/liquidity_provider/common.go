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
	"strings"
)

const (
	defaultUsername    = "admin"
	credentialSaltSize = 32
)

var (
	BadLoginError                 = errors.New("incorrect username or credentials")
	LiquidityCheckNotEnabledError = errors.New("public liquidity check is not enabled")
)

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
			return cmp.Compare(strings.ToLower(a.Address), strings.ToLower(b.Address))
		},
	)
	if !found {
		return 0, usecases.ProviderConfigurationError
	}
	return providers[index].Id, nil
}

// DefaultCredentialsProvider this is an interface to be implemented by those use case that require to use ValidateCredentials,
// since that function requires a way to access to the default password set by the application
type DefaultCredentialsProvider interface {
	LiquidityProviderRepository() liquidity_provider.LiquidityProviderRepository
	GetDefaultCredentialsChannel() <-chan entities.Event
	SetDefaultCredentials(password *liquidity_provider.HashedCredentials)
	DefaultCredentials() *liquidity_provider.HashedCredentials
}

func ValidateCredentials(
	ctx context.Context,
	useCase DefaultCredentialsProvider,
	credentials liquidity_provider.Credentials,
) error {
	var credentialsToCompare liquidity_provider.HashedCredentials
	var err error

	storedCredentials, err := useCase.LiquidityProviderRepository().GetCredentials(ctx)
	if err != nil {
		return err
	}
	if storedCredentials == nil {
		if credentialsToCompare, err = ReadDefaultCredentials(useCase); err != nil {
			return err
		}
	} else {
		credentialsToCompare = storedCredentials.Value
	}

	usernameHash, err := utils.HashArgon2(credentials.Username, credentialsToCompare.UsernameSalt)
	if err != nil {
		return err
	}
	passwordHash, err := utils.HashArgon2(credentials.Password, credentialsToCompare.PasswordSalt)
	if err != nil {
		return err
	}
	err = compareCredentials(usernameHash, passwordHash, credentialsToCompare)
	if err != nil {
		return err
	}
	return nil
}

func ReadDefaultCredentials(useCase DefaultCredentialsProvider) (liquidity_provider.HashedCredentials, error) {
	if useCase.DefaultCredentials() != nil {
		return *useCase.DefaultCredentials(), nil
	}
	select {
	case event := <-useCase.GetDefaultCredentialsChannel():
		parsedEvent, ok := event.(liquidity_provider.DefaultCredentialsSetEvent)
		if !ok {
			return liquidity_provider.HashedCredentials{}, errors.New("wrong event error")
		}
		useCase.SetDefaultCredentials(parsedEvent.Credentials)
		return *useCase.DefaultCredentials(), nil
	default:
		return liquidity_provider.HashedCredentials{}, errors.New("default password not set")
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
