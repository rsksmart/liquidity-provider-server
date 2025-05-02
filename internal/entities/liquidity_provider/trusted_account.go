package liquidity_provider

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

var (
	ErrTrustedAccountNotFound = errors.New("trusted account not found")
	ErrDuplicateAddress       = errors.New("address already exists")
)

type TrustedAccountDetails struct {
	Address          string        `json:"address" bson:"address" validate:"required"`
	Name             string        `json:"name" bson:"name" validate:"required"`
	Btc_locking_cap  *entities.Wei `json:"btc_locking_cap" bson:"btc_locking_cap" validate:"required"`
	Rbtc_locking_cap *entities.Wei `json:"rbtc_locking_cap" bson:"rbtc_locking_cap" validate:"required"`
	Signature        string        `json:"signature" bson:"signature"`
	Hash             string        `json:"hash" bson:"hash"`
}

type TrustedAccountRepository interface {
	GetTrustedAccount(ctx context.Context, address string) (*TrustedAccountDetails, error)
	GetAllTrustedAccounts(ctx context.Context) ([]TrustedAccountDetails, error)
	AddTrustedAccount(ctx context.Context, account TrustedAccountDetails) error
	UpdateTrustedAccount(ctx context.Context, account TrustedAccountDetails) error
	DeleteTrustedAccount(ctx context.Context, address string) error
}
