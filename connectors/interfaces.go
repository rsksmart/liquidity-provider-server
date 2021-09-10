package connectors

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	"github.com/rsksmart/liquidity-provider/types"
	"math/big"
	"time"
)

type AddressWatcher interface {
	OnNewConfirmation(txHash string, confirmations int64, amount float64)
}

type BTCInterface interface {
	Connect(endpoint string, username string, password string) error
	AddAddressWatcher(address string, interval time.Duration, w AddressWatcher) error
	GetParams() chaincfg.Params
	RemoveAddressWatcher(address string)
	Close()
	SerializePMT(txHash string) ([]byte, error)
	SerializeTx(txHash string) ([]byte, error)
	GetBlockNumberByTx(txHash string) (int64, error)
	GetDerivedBitcoinAddress(userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error)
}

type RSKInterface interface {
	Connect(endpoint string) error
	Close()
	EstimateGas(addr string, value big.Int, data []byte) (uint64, error)
	GasPrice() (*big.Int, error)
	GetChainId() *big.Int
	HashQuote(q *types.Quote) (string, error)
	ParseQuote(q *types.Quote) (bindings.LiquidityBridgeContractQuote, error)
	RegisterPegIn(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote, signature []byte, btcRawTrx []byte, partialMerkleTree []byte, height *big.Int) (*gethTypes.Transaction, error)
	GetFedSize() (int, error)
	GetFedThreshold() (int, error)
	GetFedPublicKey(index int) (string, error)
	GetFedAddress() (string, error)
	GetActiveFederationCreationBlockHeight() (int, error)
	GetLBCAddress() string
	GetRequiredBridgeConfirmations() int64
	CallForUser(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote) (*gethTypes.Transaction, error)
}
