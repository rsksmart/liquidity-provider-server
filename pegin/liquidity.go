package pegin

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"bytes"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/term"
)

type LiquidityProvider interface {
	Address() string
	GetQuote(*Quote, uint64, *types.Wei) (*Quote, error)
	SignQuote(hash []byte, depositAddr string, reqLiq *types.Wei) ([]byte, error)
	SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error)
}

type LocalProviderRepository interface {
	RetainQuote(rq *types.RetainedQuote) error
	HasRetainedQuote(hash string) (bool, error)
	HasLiquidity(lp LiquidityProvider, wei *types.Wei) (bool, error)
}

type LocalProvider struct {
	mu         sync.Mutex
	account    *accounts.Account
	ks         *keystore.KeyStore
	cfg        ProviderConfig
	repository LocalProviderRepository
}

type ProviderConfig struct {
	Keydir         string         `env:"KEY_DIR"`
	BtcAddr        string         `env:"BTC_ADDR"`
	AccountNum     int            `env:"ACCOUNT_NUM"`
	PwdFile        string         `env:"PWD_FILE"`
	ChainId        *big.Int       `env:"CHAIN_ID"`
	MaxConf        uint16         `env:"MAX_CONF"`
	Confirmations  map[int]uint16 `env:"CONFIRMATIONS,delimiter=|"`
	TimeForDeposit uint32         `env:"TIME_FOR_DEPOSIT"`
	CallTime       uint32         `env:"CALL_TIME"`
	CallFee        *types.Wei     `env:"CALL_FEE"`
	PenaltyFee     *types.Wei     `env:"PENALTY_FEE"`
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

func (lp *LocalProvider) Address() string {
	return lp.account.Address.String()
}

func (lp *LocalProvider) GetQuote(q *Quote, gas uint64, gasPrice *types.Wei) (*Quote, error) {
	res := *q
	res.LPBTCAddr = lp.cfg.BtcAddr
	res.LPRSKAddr = lp.account.Address.String()
	res.AgreementTimestamp = uint32(time.Now().Unix())
	res.Nonce = int64(rand.Int())
	res.TimeForDeposit = lp.cfg.TimeForDeposit
	res.LpCallTime = lp.cfg.CallTime
	res.PenaltyFee = lp.cfg.PenaltyFee.Copy()

	res.Confirmations = lp.cfg.MaxConf
	for _, k := range sortedConfirmations(lp.cfg.Confirmations) {
		v := lp.cfg.Confirmations[k]

		if res.Value.AsBigInt().Uint64() < uint64(k) {
			res.Confirmations = v
			break
		}
	}
	callCost := new(types.Wei).Mul(gasPrice, types.NewUWei(gas))
	res.CallFee = new(types.Wei).Add(callCost, lp.cfg.CallFee)
	return &res, nil
}

func (lp *LocalProvider) SignQuote(hash []byte, depositAddr string, reqLiq *types.Wei) ([]byte, error) {
	quoteHash := hex.EncodeToString(hash)

	var buf bytes.Buffer
	buf.WriteString("\x19Ethereum Signed Message:\n32")
	buf.Write(hash)

	signB, err := lp.ks.SignHash(*lp.account, crypto.Keccak256(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	signB[len(signB)-1] += 27 // v must be 27 or 28

	lp.mu.Lock()
	defer lp.mu.Unlock()

	hasRq, err := lp.repository.HasRetainedQuote(quoteHash)
	if err != nil {
		return nil, err
	}
	if !hasRq {
		hasLiquidity, err := lp.repository.HasLiquidity(lp, reqLiq)
		if err != nil {
			return nil, err
		}
		if !hasLiquidity {
			return nil, fmt.Errorf("not enough liquidity. required: %v", reqLiq)
		}

		signature := hex.EncodeToString(signB)
		rq := types.RetainedQuote{
			QuoteHash:   quoteHash,
			DepositAddr: depositAddr,
			Signature:   signature,
			ReqLiq:      reqLiq.Copy(),
			State:       types.RQStateWaitingForDeposit,
		}
		err = lp.repository.RetainQuote(&rq)
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
