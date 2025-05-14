package liquidity_provider

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

var (
	ErrTrustedAccountNotFound  = errors.New("trusted account not found")
	ErrDuplicateTrustedAccount = errors.New("trusted account already exists")
	ErrTamperedTrustedAccount  = errors.New("trusted account signature verification failed")
)

type TrustedAccountDetails struct {
	Address        string        `json:"address" bson:"address" validate:"required"`
	Name           string        `json:"name" bson:"name" validate:"required"`
	BtcLockingCap  *entities.Wei `json:"btcLockingCap" bson:"btc_locking_cap" validate:"required"`
	RbtcLockingCap *entities.Wei `json:"rbtcLockingCap" bson:"rbtc_locking_cap" validate:"required"`
}

type TrustedAccountRepository interface {
	GetTrustedAccount(ctx context.Context, address string) (*entities.Signed[TrustedAccountDetails], error)
	GetAllTrustedAccounts(ctx context.Context) ([]entities.Signed[TrustedAccountDetails], error)
	AddTrustedAccount(ctx context.Context, account entities.Signed[TrustedAccountDetails]) error
	UpdateTrustedAccount(ctx context.Context, account entities.Signed[TrustedAccountDetails]) error
	DeleteTrustedAccount(ctx context.Context, address string) error
}
