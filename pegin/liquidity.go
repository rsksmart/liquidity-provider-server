package pegin

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rsksmart/liquidity-provider-server/account"
	"math/big"
	"math/rand"
	"sort"
	"sync"
	"time"

	"bytes"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider/types"
)

type LiquidityProvider interface {
	Address() string
	GetQuote(*Quote, uint64, *types.Wei) (*Quote, error)
	SignQuote(hash []byte, depositAddr string, reqLiq *types.Wei) ([]byte, error)
	SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error)
	HasLiquidity(reqLiq *types.Wei) (bool, error)
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
	chainId    *big.Int
}

type ProviderConfig struct {
	BtcAddr             string         `env:"BTC_ADDR"`
	MaxConf             uint16         `env:"MAX_CONF"`
	Confirmations       map[int]uint16 `env:"CONFIRMATIONS"`
	TimeForDeposit      uint32         `env:"TIME_FOR_DEPOSIT"`
	CallTime            uint32         `env:"CALL_TIME"`
	PenaltyFee          *types.Wei     `env:"PENALTY_FEE"`
	Fee                 *types.Wei     `env:"FEE"`
	MinTransactionValue *big.Int       `env:"MIN_TRANSACTION_VALUE"`
	MaxTransactionValue *big.Int       `env:"MAX_TRANSACTION_VALUE"`
}

func NewLocalProvider(config ProviderConfig, repository LocalProviderRepository, accountProvider account.AccountProvider, chainId *big.Int) (*LocalProvider, error) {
	acc, err := accountProvider.GetAccount()

	if err != nil {
		return nil, err
	}
	lp := LocalProvider{
		account:    acc.Account,
		ks:         acc.Keystore,
		cfg:        config,
		repository: repository,
		chainId:    chainId,
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
	fee := lp.cfg.Fee
	res.CallFee = new(types.Wei).Add(callCost, fee)
	res.CallCost = callCost
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
	return lp.ks.SignTx(*lp.account, tx, lp.chainId)
}

func (lp *LocalProvider) HasLiquidity(reqLiq *types.Wei) (bool, error) {
	hasLiquidity, err := lp.repository.HasLiquidity(lp, reqLiq)
	if err != nil {
		return false, err
	}

	return hasLiquidity, err
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

func GetPeginProviderByAddress(liquidityProvider LiquidityProvider, addr string) LiquidityProvider {
	if liquidityProvider.Address() == addr {
		return liquidityProvider
	}

	return nil
}

func GetPeginProviderTransactOpts(liquidityProvider LiquidityProvider, address string) (*bind.TransactOpts, error) {
	lp := GetPeginProviderByAddress(liquidityProvider, address)
	if lp == nil {
		return nil, errors.New("missing liquidity provider")
	}
	return &bind.TransactOpts{
		From:   common.HexToAddress(address),
		Signer: lp.SignTx,
	}, nil
}
