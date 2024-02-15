package blockchain

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"math/big"
	"regexp"
)

var (
	btcTestnetP2PKHRegex = regexp.MustCompile("^[mn]([a-km-zA-HJ-NP-Z1-9]{25,34})$")
	btcMainnetP2PKHRegex = regexp.MustCompile("^[1]([a-km-zA-HJ-NP-Z1-9]{25,34})$")
	btcMainnetP2SHRegex  = regexp.MustCompile("^[3]([a-km-zA-HJ-NP-Z1-9]{33,34})$")
	btcTestnetP2SHRegex  = regexp.MustCompile("^[2]([a-km-zA-HJ-NP-Z1-9]{33,34})$")
)

// IsSupportedBtcAddress checks if flyover protocol supports the given address
// Currently the supported address types are P2PKH and P2SH
func IsSupportedBtcAddress(address string) bool {
	return isP2PKH(address) || isP2SH(address)
}

func isP2PKH(address string) bool {
	return btcTestnetP2PKHRegex.MatchString(address) || btcMainnetP2PKHRegex.MatchString(address)
}

func isP2SH(address string) bool {
	return btcTestnetP2SHRegex.MatchString(address) || btcMainnetP2SHRegex.MatchString(address)
}

type BitcoinWallet interface {
	EstimateTxFees(toAddress string, value *entities.Wei) (*entities.Wei, error)
	GetBalance() (*entities.Wei, error)
	SendWithOpReturn(address string, value *entities.Wei, opReturnContent []byte) (string, error)
	ImportAddress(address string) error
	GetTransactions(address string) ([]BitcoinTransactionInformation, error)
	Unlock() error
}

type BitcoinNetwork interface {
	DecodeAddress(address string, keepVersion bool) ([]byte, error)
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
}

type MerkleBranch struct {
	Hashes [][32]byte
	Path   *big.Int
}
