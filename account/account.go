package account

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type RSKAccount struct {
	Account  *accounts.Account
	Keystore *keystore.KeyStore
}

type AccountProvider interface {
	GetAccount() (*RSKAccount, error)
}
