package liquidity_provider

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

var (
	TrustedAccountNotFoundError = errors.New("trusted account not found")
	DuplicateAddressError       = errors.New("address already exists")
)

type TrustedAccount struct {
	Address        string        `json:"address" bson:"address" validate:"required"`
	Name           string        `json:"name" bson:"name" validate:"required"`
	BtcLockingCap  *entities.Wei `json:"btcLockingCap" bson:"btcLockingCap" validate:"required"`
	RbtcLockingCap *entities.Wei `json:"rbtcLockingCap" bson:"rbtcLockingCap" validate:"required"`
}

type TrustedAccountRepository interface {
	GetTrustedAccount(ctx context.Context, address string) (*TrustedAccount, error)
	GetAllTrustedAccounts(ctx context.Context) ([]TrustedAccount, error)
	AddTrustedAccount(ctx context.Context, account TrustedAccount) error
	UpdateTrustedAccount(ctx context.Context, account TrustedAccount) error
	DeleteTrustedAccount(ctx context.Context, address string) error
}
