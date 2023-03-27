package account

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

type localAccountProvider struct {
	keyDir       string
	passwordFile string
	accountNum   int
}

func NewLocalAccountProvider(keyDir string, passwordFile string, accountNum int) AccountProvider {
	return &localAccountProvider{keyDir: keyDir, passwordFile: passwordFile, accountNum: accountNum}
}

func (provider *localAccountProvider) GetAccount() (*RSKAccount, error) {
	if provider.keyDir == "" {
		provider.keyDir = "keystore"
	}
	if err := os.MkdirAll(provider.keyDir, 0700); err != nil {
		return nil, err
	}
	var f *os.File
	if provider.passwordFile != "" {
		var err error
		f, err = os.Open(provider.passwordFile)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", provider.passwordFile)
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
	}

	ks := keystore.NewKeyStore(provider.keyDir, keystore.StandardScryptN, keystore.StandardScryptP)
	if account, err := retrieveOrCreateAccount(ks, provider.accountNum, f); err != nil {
		return nil, err
	} else {
		return &RSKAccount{
			Account:  account,
			Keystore: ks,
		}, nil
	}
}

func retrieveOrCreateAccount(ks *keystore.KeyStore, accountNum int, in *os.File) (*accounts.Account, error) {
	if cap(ks.Accounts()) == 0 {
		log.Info("no RSK account found")
		acc, err := createAccount(ks, in)
		return acc, err
	} else {
		if cap(ks.Accounts()) <= accountNum {
			return nil, fmt.Errorf("account number %v not found", accountNum)
		}
		acc := ks.Accounts()[accountNum]
		passwd, err := enterPasswd(in)

		if err != nil {
			return nil, err
		}
		err = ks.Unlock(acc, passwd)
		return &acc, err
	}
}

func createAccount(ks *keystore.KeyStore, in *os.File) (*accounts.Account, error) {
	passwd, err := createPasswd(in)

	if err != nil {
		return nil, err
	}
	acc, err := ks.NewAccount(passwd)

	if err != nil {
		return &acc, err
	}
	err = ks.Unlock(acc, passwd)

	if err != nil {
		return &acc, err
	}
	log.Info("new account created: ", acc.Address)
	return &acc, err
}

func enterPasswd(in *os.File) (string, error) {
	fmt.Println("enter password for RSK account")
	fmt.Print("password: ")
	var pwd string
	var err error
	if in == nil {
		pwd, err = readPasswdCons(nil)
	} else {
		pwd, err = readPasswdReader(bufio.NewReader(in))
	}
	fmt.Println()
	return pwd, err
}

func createPasswd(in *os.File) (string, error) {
	fmt.Println("creating password for new RSK account")
	fmt.Println("WARNING: the account will be lost forever if you forget this password!!! Do you understand? (yes/[no])")

	var r *bufio.Reader
	var readPasswd func(*bufio.Reader) (string, error)
	if in == nil {
		r = bufio.NewReader(os.Stdin)
		readPasswd = readPasswdCons
	} else {
		r = bufio.NewReader(in)
		readPasswd = readPasswdReader
	}

	str, _ := r.ReadString('\n')
	if str != "yes\n" {
		return "", errors.New("must say yes")
	}
	fmt.Print("password: ")
	pwd1, err := readPasswd(r)
	fmt.Println()
	if err != nil {
		return "", err
	}

	fmt.Print("repeat password: ")
	pwd2, err := readPasswd(r)
	fmt.Println()
	if err != nil {
		return "", err
	}
	if pwd1 != pwd2 {
		return "", errors.New("passwords do not match")
	}
	return pwd1, nil
}

func readPasswdCons(_ *bufio.Reader) (string, error) {
	pass, err := term.ReadPassword(syscall.Stdin)
	return string(pass), err
}

func readPasswdReader(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(str, "\n"), nil
}
