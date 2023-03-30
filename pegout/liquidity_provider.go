package pegout

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/account"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider/types"
)

type LocalProvider struct {
	mu         sync.Mutex
	account    *accounts.Account
	ks         *keystore.KeyStore
	cfg        *ProviderConfig
	repository LocalProviderRepository
}

type ProviderConfig struct {
	pegin.ProviderConfig
	DepositConfirmations  uint16 `env:"DEPOSIT_CONFIRMATIONS"`
	DepositDateLimit      uint32 `env:"DEPOSIT_DATE_LIMIT"`
	TransferConfirmations uint16 `env:"TRANSFER_CONFIRMATIONS"`
	TransferTime          uint32 `env:"TRANSFER_TIME"`
	ExpireDate            uint32 `env:"EXPIRED_DATE"`
	ExpireBlocks          uint32 `env:"EXPIRED_BLOCKS"`
}

type LocalProviderRepository interface {
	RetainPegOutQuote(rq *RetainedQuote) error
	HasRetainedPegOutQuote(hash string) (bool, error)
	HasLiquidityPegOut(satoshis uint64) (bool, error)
}

type LiquidityProvider interface {
	Address() string
	GetQuote(*Quote, uint64, uint64, *types.Wei) (*Quote, error)
	SignQuote(hash []byte, depositAddr string, satoshis uint64) ([]byte, error)
	SignTx(common.Address, *gethTypes.Transaction) (*gethTypes.Transaction, error)
}

func NewLocalProvider(config *ProviderConfig, repository LocalProviderRepository, accountProvider account.AccountProvider) (*LocalProvider, error) {
	acc, err := accountProvider.GetAccount()

	if err != nil {
		return nil, err
	}
	lp := LocalProvider{
		account:    acc.Account,
		ks:         acc.Keystore,
		cfg:        config,
		repository: repository,
	}
	return &lp, nil
}

func GetPegoutProviderByAddress(liquidityProviders []LiquidityProvider, addr string) LiquidityProvider {
	for _, p := range liquidityProviders {
		if p.Address() == addr {
			return p
		}
	}
	return nil
}

func (lp *LocalProvider) GetQuote(q *Quote, rskLastBlockNumber uint64, gas uint64, gasPrice *types.Wei) (*Quote, error) {
	res := *q
	res.LPRSKAddr = lp.account.Address.String()
	res.AgreementTimestamp = uint32(time.Now().Unix())
	res.Nonce = int64(rand.Int())
	res.DepositDateLimit = lp.cfg.DepositDateLimit
	res.TransferConfirmations = lp.cfg.TransferConfirmations
	res.TransferTime = lp.cfg.TransferTime
	res.ExpireDate = res.AgreementTimestamp + lp.cfg.ExpireDate
	res.ExpireBlock = lp.cfg.ExpireBlocks + uint32(rskLastBlockNumber)
	res.PenaltyFee = lp.cfg.PenaltyFee.Uint64()
	res.LpBTCAddr = lp.cfg.BtcAddr

	res.DepositConfirmations = lp.cfg.MaxConf
	for _, k := range sortedConfirmations(lp.cfg.Confirmations) {
		v := lp.cfg.Confirmations[k]

		if res.Value.Uint64() < uint64(k) {
			res.DepositConfirmations = v
			break
		}
	}
	callCost := new(types.Wei).Mul(types.NewUWei(gasPrice.Uint64()), types.NewUWei(gas))
	res.CallFee = new(types.Wei).Add(callCost, types.NewUWei(lp.cfg.CallFee.Uint64()))
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
		hasLiquidity, err := lp.repository.HasLiquidityPegOut(satoshis)
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
