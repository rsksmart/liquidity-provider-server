package entities

import (
	"context"
	"errors"
)

type ProviderType string

const (
	PeginProvider  ProviderType = "pegin"
	PegoutProvider ProviderType = "pegout"
	FullProvider   ProviderType = "both"
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
	GetBitcoinConfirmationsForValue(value *Wei) uint16
	GetRootstockConfirmationsForValue(value *Wei) uint16
}

type PeginLiquidityProvider interface {
	ValidateAmountForPegin(amount *Wei) error
	HasPeginLiquidity(ctx context.Context, requiredLiquidity *Wei) error
	CallTime() uint32
	CallFeePegin() *Wei
	PenaltyFeePegin() *Wei
	TimeForDepositPegin() uint32
	MaxPegin() *Wei
	MinPegin() *Wei
	MaxPeginConfirmations() uint16
}

type PegoutLiquidityProvider interface {
	ValidateAmountForPegout(amount *Wei) error
	HasPegoutLiquidity(ctx context.Context, requiredLiquidity *Wei) error
	CallFeePegout() *Wei
	PenaltyFeePegout() *Wei
	TimeForDepositPegout() uint32
	ExpireBlocksPegout() uint64
	MaxPegout() *Wei
	MinPegout() *Wei
	MaxPegoutConfirmations() uint16
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
	Fee                   *Wei   `json:"fee" validate:"required"`
	MinTransactionValue   *Wei   `json:"minTransactionValue"  validate:"required"`
	MaxTransactionValue   *Wei   `json:"maxTransactionValue"  validate:"required"`
	RequiredConfirmations uint16 `json:"requiredConfirmations"  validate:"required"`
}

type PunishmentEvent struct {
	LiquidityProvider string
	Penalty           *Wei
	QuoteHash         string
}
