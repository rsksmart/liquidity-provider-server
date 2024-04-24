package account

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetRskAccount(keyDir string, accountNum int, encryptedJson, password string) (*rootstock.RskAccount, error) {
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)
	if account, err := retrieveOrCreateAccount(ks, accountNum, encryptedJson, password); err != nil {
		return nil, err
	} else {
		return &rootstock.RskAccount{
			Account:  account,
			Keystore: ks,
		}, nil
	}
}

func createAccount(ks *keystore.KeyStore, encryptedJson, password string) (*accounts.Account, error) {
	account, err := ks.Import([]byte(encryptedJson), password, password)
	if err != nil {
		return nil, err
	}

	err = ks.Unlock(account, password)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func retrieveOrCreateAccount(ks *keystore.KeyStore, accountNum int, encryptedJson, password string) (*accounts.Account, error) {
	if cap(ks.Accounts()) == 0 {
		log.Debug("No RSK account found")
		acc, err := createAccount(ks, encryptedJson, password)
		return acc, err
	} else {
		if cap(ks.Accounts()) <= accountNum {
			return nil, fmt.Errorf("account number %v not found", accountNum)
		}
		acc := ks.Accounts()[accountNum]
		err := ks.Unlock(acc, password)
		return &acc, err
	}
}
