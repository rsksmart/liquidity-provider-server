package liquidity_provider

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

type ProviderType string

const (
	PeginProvider  ProviderType = "pegin"
	PegoutProvider ProviderType = "pegout"
	FullProvider   ProviderType = "both"
)

const (
	DefaultCredentialsSetEventId entities.EventId = "CredentialsSet"
)

var InvalidProviderTypeError = errors.New("invalid liquidity provider type")

func (p ProviderType) IsValid() bool {
	switch p {
	case PegoutProvider, PeginProvider, FullProvider:
		return true
	default:
		return false
	}
}

func (p ProviderType) AcceptsPegin() bool {
	return p == PeginProvider || p == FullProvider
}

func (p ProviderType) AcceptsPegout() bool {
	return p == PegoutProvider || p == FullProvider
}

func ToProviderType(value string) (ProviderType, error) {
	providerType := ProviderType(value)
	if providerType.IsValid() {
		return providerType, nil
	} else {
		return "", InvalidProviderTypeError
	}
}

type LiquidityProvider interface {
	RskAddress() string
	BtcAddress() string
	SignQuote(quoteHash string) (string, error)
	GeneralConfiguration(ctx context.Context) GeneralConfiguration
}

type PeginLiquidityProvider interface {
	HasPeginLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error
	PeginConfiguration(ctx context.Context) PeginConfiguration
}

type PegoutLiquidityProvider interface {
	HasPegoutLiquidity(ctx context.Context, requiredLiquidity *entities.Wei) error
	PegoutConfiguration(ctx context.Context) PegoutConfiguration
}

type LiquidityProviderRepository interface {
	GetPeginConfiguration(ctx context.Context) (*entities.Signed[PeginConfiguration], error)
	UpsertPeginConfiguration(ctx context.Context, configuration entities.Signed[PeginConfiguration]) error
	GetPegoutConfiguration(ctx context.Context) (*entities.Signed[PegoutConfiguration], error)
	UpsertPegoutConfiguration(ctx context.Context, configuration entities.Signed[PegoutConfiguration]) error
	GetGeneralConfiguration(ctx context.Context) (*entities.Signed[GeneralConfiguration], error)
	UpsertGeneralConfiguration(ctx context.Context, configuration entities.Signed[GeneralConfiguration]) error
	GetCredentials(ctx context.Context) (*entities.Signed[HashedCredentials], error)
	UpsertCredentials(ctx context.Context, credentials entities.Signed[HashedCredentials]) error
}

type RegisteredLiquidityProvider struct {
	Id           uint64       `json:"id" validate:"required"`
	Address      string       `json:"address" validate:"required"`
	Name         string       `json:"name" validate:"required"`
	ApiBaseUrl   string       `json:"apiBaseUrl" validate:"required"`
	Status       bool         `json:"status" validate:"required"`
	ProviderType ProviderType `json:"providerType" validate:"required"`
}

type LiquidityProviderDetail struct {
	Fee                   *entities.Wei `json:"fee" validate:"required"`
	MinTransactionValue   *entities.Wei `json:"minTransactionValue"  validate:"required"`
	MaxTransactionValue   *entities.Wei `json:"maxTransactionValue"  validate:"required"`
	RequiredConfirmations uint16        `json:"requiredConfirmations"  validate:"required"`
}

type PunishmentEvent struct {
	LiquidityProvider string
	Penalty           *entities.Wei
	QuoteHash         string
}

type Credentials struct {
	Username string
	Password string
}

type DefaultCredentialsSetEvent struct {
	entities.Event
	Credentials *HashedCredentials
}
