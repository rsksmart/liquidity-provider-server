package account

import (
	"encoding/hex"
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"unsafe"
)

var NoDerivationError = fmt.Errorf("btc derivation wasn't enabled for this account")

type RskAccount struct {
	Account  *accounts.Account
	Keystore *keystore.KeyStore
	btc      *btcDerivationInfo
}

type btcDerivationInfo struct {
	pubKey       *btcec.PublicKey
	address      btcutil.Address
	protectedWif *memguard.Enclave
}

type CreationArgs struct {
	KeyDir        string
	AccountNum    int
	EncryptedJson string
	Password      string
}

type CreationWithDerivationArgs struct {
	CreationArgs
	BtcParams *chaincfg.Params
}

func GetRskAccount(args CreationArgs) (*RskAccount, error) {
	if err := os.MkdirAll(args.KeyDir, 0700); err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore(args.KeyDir, keystore.StandardScryptN, keystore.StandardScryptP)
	if account, err := retrieveOrCreateAccount(ks, args.AccountNum, args.EncryptedJson, args.Password); err != nil {
		return nil, err
	} else {
		return &RskAccount{
			Account:  account,
			Keystore: ks,
		}, nil
	}
}

// GetRskAccountWithDerivation returns an RSK account with the corresponding BTC derivative information
func GetRskAccountWithDerivation(args CreationWithDerivationArgs) (*RskAccount, error) {
	account, err := GetRskAccount(args.CreationArgs)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey([]byte(args.EncryptedJson), args.Password)
	if err != nil {
		return nil, err
	}

	privateKey, pubKey := btcec.PrivKeyFromBytes(key.PrivateKey.D.Bytes())
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), args.BtcParams)
	if err != nil {
		return nil, err
	}

	protectedWifBuffer, protectedWif := utils.GetSecurePointer[btcutil.WIF]()
	unprotectedWif, err := btcutil.NewWIF(privateKey, args.BtcParams, true)
	if err != nil {
		return nil, err
	}

	// this line is to write the content of the unprotectedWif to the protected memory address inside the locked buffer, the protectedWif
	// variable is just to allow us to write inside buffer's memory address, then we set unprotectedWif to its zero value
	*protectedWif = *unprotectedWif
	*unprotectedWif = btcutil.WIF{}
	account.btc = &btcDerivationInfo{pubKey: pubKey, address: address, protectedWif: protectedWifBuffer.Seal()}
	return account, nil
}

func (account *RskAccount) BtcPubKey() (string, error) {
	if account.btc == nil || account.btc.pubKey == nil {
		return "", NoDerivationError
	}
	pubKeyBytes := account.btc.pubKey.SerializeCompressed()
	return hex.EncodeToString(pubKeyBytes), nil
}

func (account *RskAccount) BtcAddress() (btcutil.Address, error) {
	if account.btc == nil || account.btc.address == nil {
		return nil, NoDerivationError
	}
	return account.btc.address, nil
}

func (account *RskAccount) UsePrivateKeyWif(usageFunc func(wif *btcutil.WIF) error) error {
	if account.btc == nil {
		return NoDerivationError
	}
	buffer, err := account.btc.protectedWif.Open()
	defer func(b *memguard.LockedBuffer) {
		if b != nil {
			b.Destroy()
		}
	}(buffer)
	if err != nil {
		return err
	}
	wif := (*btcutil.WIF)(unsafe.Pointer(&buffer.Bytes()[0]))
	return usageFunc(wif)
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
