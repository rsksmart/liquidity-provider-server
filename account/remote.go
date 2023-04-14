package account

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rsksmart/liquidity-provider-server/secrets"
	log "github.com/sirupsen/logrus"
	"os"
)

type AccountSecretNames struct {
	KeySecretName      string
	PasswordSecretName string
}

type remoteAccountProvider struct {
	keyDir      string
	accountNum  int
	keyStorage  secrets.SecretStorage[any]
	secretNames *AccountSecretNames
}

func NewRemoteAccountProvider(keyDir string, accountNum int, secretNames *AccountSecretNames, keyStorage secrets.SecretStorage[any]) AccountProvider {
	return &remoteAccountProvider{keyDir: keyDir, accountNum: accountNum, secretNames: secretNames, keyStorage: keyStorage}
}

func (provider *remoteAccountProvider) GetAccount() (*RSKAccount, error) {
	if provider.keyDir == "" {
		provider.keyDir = "keystore"
	}
	if err := os.MkdirAll(provider.keyDir, 0700); err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(provider.keyDir, keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := provider.retrieveOrCreateAccount(ks, provider.accountNum)
	if err != nil {
		return nil, err
	}

	return &RSKAccount{
		Account:  account,
		Keystore: ks,
	}, nil
}

func (provider *remoteAccountProvider) retrieveOrCreateAccount(ks *keystore.KeyStore, accountNum int) (*accounts.Account, error) {
	if cap(ks.Accounts()) == 0 {
		log.Info("no RSK account found")
		acc, err := provider.createAccount(ks)
		return acc, err
	} else {
		if cap(ks.Accounts()) <= accountNum {
			return nil, fmt.Errorf("account number %v not found", accountNum)
		}
		acc := ks.Accounts()[accountNum]
		password, err := provider.keyStorage.GetTextSecret(provider.secretNames.PasswordSecretName)
		if err != nil {
			return nil, err
		}

		err = ks.Unlock(acc, password)
		return &acc, err
	}
}

func (provider *remoteAccountProvider) createAccount(ks *keystore.KeyStore) (*accounts.Account, error) {
	key, err := provider.keyStorage.GetTextSecret(provider.secretNames.KeySecretName)
	if err != nil {
		return nil, err
	}

	password, err := provider.keyStorage.GetTextSecret(provider.secretNames.PasswordSecretName)
	if err != nil {
		return nil, err
	}

	account, err := ks.Import([]byte(key), password, password)
	if err != nil {
		return nil, err
	}

	err = ks.Unlock(account, password)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
