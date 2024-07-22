package blockchain

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"math/big"
	"regexp"
	"time"
)

var (
	btcTestnetP2PKHRegex  = regexp.MustCompile("^[mn]([a-km-zA-HJ-NP-Z1-9]{25,34})$")
	btcMainnetP2PKHRegex  = regexp.MustCompile("^[1]([a-km-zA-HJ-NP-Z1-9]{25,34})$")
	btcMainnetP2SHRegex   = regexp.MustCompile("^[3]([a-km-zA-HJ-NP-Z1-9]{33,34})$")
	btcTestnetP2SHRegex   = regexp.MustCompile("^[2]([a-km-zA-HJ-NP-Z1-9]{33,34})$")
	btcMainnetP2WPKHRegex = regexp.MustCompile("^(bc1q)([ac-hj-np-z02-9]{38})$")
	btcTestnetP2WPKHRegex = regexp.MustCompile("^(tb1q)([ac-hj-np-z02-9]{38})$")
	btcRegtestP2WPKHRegex = regexp.MustCompile("^(bcrt1q)([ac-hj-np-z02-9]{38})$")
	btcMainnetP2WSHRegex  = regexp.MustCompile("^(bc1q)([ac-hj-np-z02-9]{58})$")
	btcTestnetP2WSHRegex  = regexp.MustCompile("^(tb1q)([ac-hj-np-z02-9]{58})$")
	btcRegtestP2WSHRegex  = regexp.MustCompile("^(bcrt1q)([ac-hj-np-z02-9]{58})$")
	btcMainnetP2TRRegex   = regexp.MustCompile("^(bc1p)([ac-hj-np-z02-9]{58})$")
	btcTestnetP2TRRegex   = regexp.MustCompile("^(tb1p)([ac-hj-np-z02-9]{58})$")
	btcRegtestP2TRRegex   = regexp.MustCompile("^(bcrt1p)([ac-hj-np-z02-9]{58})$")
)

var (
	BtcAddressInvalidNetworkError = errors.New("address network is not valid")
	BtcAddressNotSupportedError   = errors.New("btc address not supported")
)

const (
	BtcChainHeightErrorTemplate = "error getting Bitcoin chain height: %v"
	BtcTxInfoErrorTemplate      = "error getting Bitcoin transaction information (%s): %v"
)

// IsSupportedBtcAddress checks if flyover protocol supports the given address
func IsSupportedBtcAddress(address string) bool {
	return IsTestnetBtcAddress(address) || IsMainnetBtcAddress(address) || IsRegtestBtcAddress(address)
}

func IsBtcP2PKHAddress(address string) bool {
	return btcTestnetP2PKHRegex.MatchString(address) || btcMainnetP2PKHRegex.MatchString(address)
}

func IsBtcP2SHAddress(address string) bool {
	return btcMainnetP2SHRegex.MatchString(address) || btcTestnetP2SHRegex.MatchString(address)
}

func IsBtcP2WPKHAddress(address string) bool {
	return btcMainnetP2WPKHRegex.MatchString(address) || btcTestnetP2WPKHRegex.MatchString(address) || btcRegtestP2WPKHRegex.MatchString(address)
}

func IsBtcP2WSHAddress(address string) bool {
	return btcMainnetP2WSHRegex.MatchString(address) || btcTestnetP2WSHRegex.MatchString(address) || btcRegtestP2WSHRegex.MatchString(address)
}

func IsBtcP2TRAddress(address string) bool {
	return btcMainnetP2TRRegex.MatchString(address) || btcTestnetP2TRRegex.MatchString(address) || btcRegtestP2TRRegex.MatchString(address)
}

func IsRegtestBtcAddress(address string) bool {
	// only base58 addresses have the same structure in regtest and testnet
	return btcRegtestP2WPKHRegex.MatchString(address) ||
		btcRegtestP2WSHRegex.MatchString(address) ||
		btcRegtestP2TRRegex.MatchString(address) ||
		btcTestnetP2PKHRegex.MatchString(address) ||
		btcTestnetP2SHRegex.MatchString(address)
}

func IsTestnetBtcAddress(address string) bool {
	return btcTestnetP2PKHRegex.MatchString(address) ||
		btcTestnetP2SHRegex.MatchString(address) ||
		btcTestnetP2WPKHRegex.MatchString(address) ||
		btcTestnetP2WSHRegex.MatchString(address) ||
		btcTestnetP2TRRegex.MatchString(address)
}

func IsMainnetBtcAddress(address string) bool {
	return btcMainnetP2PKHRegex.MatchString(address) ||
		btcMainnetP2SHRegex.MatchString(address) ||
		btcMainnetP2WPKHRegex.MatchString(address) ||
		btcMainnetP2WSHRegex.MatchString(address) ||
		btcMainnetP2TRRegex.MatchString(address)
}

type BitcoinWallet interface {
	entities.Closeable
	EstimateTxFees(toAddress string, value *entities.Wei) (*entities.Wei, error)
	GetBalance() (*entities.Wei, error)
	SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (string, error)
	ImportAddress(address string) error
	GetTransactions(address string) ([]BitcoinTransactionInformation, error)
	Address() string
	Unlock() error
}

type BitcoinNetwork interface {
	ValidateAddress(address string) error
	DecodeAddress(address string) ([]byte, error)
	GetTransactionInfo(hash string) (BitcoinTransactionInformation, error)
	GetRawTransaction(hash string) ([]byte, error)
	GetPartialMerkleTree(hash string) ([]byte, error)
	GetHeight() (*big.Int, error)
	BuildMerkleBranch(txHash string) (MerkleBranch, error)
	GetTransactionBlockInfo(txHash string) (BitcoinBlockInformation, error)
}

type BitcoinTransactionInformation struct {
	Hash          string
	Confirmations uint64
	Outputs       map[string][]*entities.Wei
}

func (tx *BitcoinTransactionInformation) AmountToAddress(address string) *entities.Wei {
	total := new(entities.Wei)
	utxos, ok := tx.Outputs[address]
	if !ok {
		return entities.NewWei(0)
	}
	for _, utxo := range utxos {
		total.Add(total, utxo)
	}
	return total
}

type BitcoinBlockInformation struct {
	Hash   [32]byte
	Height *big.Int
	Time   time.Time
}

type MerkleBranch struct {
	Hashes [][32]byte
	Path   *big.Int
}
