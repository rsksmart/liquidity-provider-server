package pegout

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

type LocalProvider struct {
	mu         sync.Mutex
	account    *accounts.Account
	ks         *keystore.KeyStore
	cfg        ProviderConfig
	repository LocalProviderRepository
}

type ProviderConfig struct {
	providers.ProviderConfig
	DepositConfirmations  uint16
	DepositDateLimit      uint32
	TransferConfirmations uint16
	TransferTime          uint32
	ExpireDate            uint32
	ExpireBlocks          uint32
	CallFee               uint64
	PenaltyFee            uint64
}

type LocalProviderRepository interface {
	RetainPegOutQuote(rq *RetainedQuote) error
	HasRetainedPegOutQuote(hash string) (bool, error)
	HasLiquidityPegOut(lp LiquidityProvider, satoshis uint64) (bool, error)
}

type LiquidityProvider interface {
	Address() string
	GetQuote(*Quote) (*Quote, error)
	SignQuote(hash []byte, depositAddr string, satoshis uint64) ([]byte, error)
	SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error)
}

func NewLocalProvider(config ProviderConfig, repository LocalProviderRepository) (*LocalProvider, error) {
	if config.Keydir == "" {
		config.Keydir = "keystore"
	}
	if err := os.MkdirAll(config.Keydir, 0700); err != nil {
		return nil, err
	}
	var f *os.File
	if config.PwdFile != "" {
		var err error
		f, err = os.Open(config.PwdFile)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", config.PwdFile)
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)
	}

	ks := keystore.NewKeyStore(config.Keydir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := retrieveOrCreateAccount(ks, config.AccountNum, f)

	if err != nil {
		return nil, err
	}
	lp := LocalProvider{
		account:    acc,
		ks:         ks,
		cfg:        config,
		repository: repository,
	}
	return &lp, nil
}

func (lp *LocalProvider) GetQuote(q *Quote) (*Quote, error) {
	res := *q
	res.LPRSKAddr = lp.account.Address.String()
	res.AgreementTimestamp = uint32(time.Now().Unix())
	res.Nonce = int64(rand.Int())
	res.DepositDateLimit = lp.cfg.DepositDateLimit
	res.TransferConfirmations = lp.cfg.TransferConfirmations
	res.TransferTime = lp.cfg.TransferTime
	res.ExpireDate = lp.cfg.ExpireDate
	res.ExpireBlocks = lp.cfg.ExpireBlocks
	res.PenaltyFee = lp.cfg.PenaltyFee

	res.DepositConfirmations = lp.cfg.MaxConf
	for _, k := range sortedConfirmations(lp.cfg.Confirmations) {
		v := lp.cfg.Confirmations[k]

		if res.Value < uint64(k) {
			res.DepositConfirmations = v
			break
		}
	}

	res.Fee = lp.cfg.CallFee
	return &res, nil
}

func sortedConfirmations(m map[int]uint16) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func (lp *LocalProvider) Address() string {
	return lp.account.Address.String()
}

func (lp *LocalProvider) SignQuote(hash []byte, depositAddr string, satoshis uint64) ([]byte, error) {
	quoteHash := hex.EncodeToString(hash)

	var buf bytes.Buffer
	buf.WriteString("\x19Ethereum Signed Message:\n32")
	buf.Write(hash)
	fmt.Println("Liquidity provider")
	fmt.Println(lp.account)
	signB, err := lp.ks.SignHash(*lp.account, crypto.Keccak256(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	signB[len(signB)-1] += 27 // v must be 27 or 28

	lp.mu.Lock()
	defer lp.mu.Unlock()

	hasRq, err := lp.repository.HasRetainedPegOutQuote(quoteHash)
	if err != nil {
		return nil, err
	}
	if !hasRq {
		hasLiquidity, err := lp.repository.HasLiquidityPegOut(lp, satoshis)
		if err != nil {
			return nil, err
		}
		if !hasLiquidity {
			return nil, fmt.Errorf("not enough liquidity. required: %v", satoshis)
		}

		signature := hex.EncodeToString(signB)
		rq := RetainedQuote{
			QuoteHash:   quoteHash,
			DepositAddr: depositAddr,
			Signature:   signature,
			ReqLiq:      satoshis,
			State:       types.RQStateWaitingForDeposit,
		}
		err = lp.repository.RetainPegOutQuote(&rq)
		if err != nil {
			return nil, err
		}
	}

	return signB, nil
}

func (lp *LocalProvider) SignTx(address common.Address, tx *gethTypes.Transaction) (*gethTypes.Transaction, error) {
	if !bytes.Equal(address[:], lp.account.Address[:]) {
		return nil, fmt.Errorf("provider address %v is incorrect", address.Hash())
	}
	return lp.ks.SignTx(*lp.account, tx, lp.cfg.ChainId)
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
