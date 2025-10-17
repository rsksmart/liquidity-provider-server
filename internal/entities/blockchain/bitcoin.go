package blockchain

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"math/big"
	"regexp"
	"strings"
	"time"
)

type BtcAddressType string

const (
	BtcAddressTypeP2PKH  BtcAddressType = "p2pkh"
	BtcAddressTypeP2SH   BtcAddressType = "p2sh"
	BtcAddressTypeP2WPKH BtcAddressType = "p2wpkh"
	BtcAddressTypeP2WSH  BtcAddressType = "p2wsh"
	BtcAddressTypeP2TR   BtcAddressType = "p2tr"
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

const (
	BitcoinMainnetP2PKHZeroAddress  = "1111111111111111111114oLvT2"
	BitcoinTestnetP2PKHZeroAddress  = "mfWxJ45yp2SFn7UciZyNpvDKrzbhyfKrY8"
	BitcoinMainnetP2SHZeroAddress   = "31h1vYVSYuKP6AhS86fbRdMw9XHieotbST"
	BitcoinTestnetP2SHZeroAddress   = "2MsFDzHRUAMpjHxKyoEHU3aMCMsVtMqs1PV"
	BitcoinMainnetP2WPKHZeroAddress = "bc1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq9e75rs"
	BitcoinTestnetP2WPKHZeroAddress = "tb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq0l98cr"
	BitcoinRegtestP2WPKHZeroAddress = "bcrt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqdku202"
	BitcoinMainnetP2WSHZeroAddress  = "bc1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqthqst8"
	BitcoinTestnetP2WSHZeroAddress  = "tb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqulkl3g"
	BitcoinRegtestP2WSHZeroAddress  = "bcrt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq3xueyj"
	BitcoinMainnetP2TRZeroAddress   = "bc1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpqqenm"
	BitcoinTestnetP2TRZeroAddress   = "tb1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqkgkkf5"
	BitcoinRegtestP2TRZeroAddress   = "bcrt1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqm3usuw"
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
	EstimateTxFees(toAddress string, value *entities.Wei) (BtcFeeEstimation, error)
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
	// GetCoinbaseInformation returns the coinbase transaction information of the block that includes txHash
	GetCoinbaseInformation(txHash string) (rootstock.BtcCoinbaseTransactionInformation, error)
	NetworkName() string
	GetBlockchainInfo() (BitcoinBlockchainInfo, error)
	GetZeroAddress(addressType BtcAddressType) (string, error)
}

type BitcoinTransactionInformation struct {
	Hash          string
	Confirmations uint64
	Outputs       map[string][]*entities.Wei
	HasWitness    bool
}

type BitcoinBlockchainInfo struct {
	NetworkName      string
	ValidatedBlocks  *big.Int
	ValidatedHeaders *big.Int
	BestBlockHash    string
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

func (tx *BitcoinTransactionInformation) UTXOsToAddress(address string) []*entities.Wei {
	utxos, ok := tx.Outputs[address]
	if !ok {
		return []*entities.Wei{}
	}
	return utxos
}

func BtcAddressTypeFromString(value string) (BtcAddressType, error) {
	value = strings.ToLower(value)
	switch value {
	case string(BtcAddressTypeP2PKH):
		return BtcAddressTypeP2PKH, nil
	case string(BtcAddressTypeP2SH):
		return BtcAddressTypeP2SH, nil
	case string(BtcAddressTypeP2WPKH):
		return BtcAddressTypeP2WPKH, nil
	case string(BtcAddressTypeP2WSH):
		return BtcAddressTypeP2WSH, nil
	case string(BtcAddressTypeP2TR):
		return BtcAddressTypeP2TR, nil
	default:
		return "", BtcAddressNotSupportedError
	}
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

type BtcFeeEstimation struct {
	Value   *entities.Wei
	FeeRate *utils.BigFloat
}
