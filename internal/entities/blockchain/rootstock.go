package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"math/big"
	"regexp"
	"strings"
	"time"
)

const (
	RskChainHeightErrorTemplate = "error getting Rootstock chain height: %v"
)

var (
	rskAddressRegex       = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	WaitingForBridgeError = errors.New("waiting for rootstock bridge")
	InvalidAddressError   = errors.New("invalid rootstock address")
)

type RskContracts struct {
	Bridge               rootstock.Bridge
	PegIn                PeginContract
	PegOut               PegoutContract
	CollateralManagement CollateralManagementContract
	Discovery            DiscoveryContract
}

func DecodeStringTrimPrefix(hexString string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(hexString, "0x"))
}
func IsRskAddress(address string) bool {
	return rskAddressRegex.MatchString(address)
}

type TransactionConfig struct {
	Value    *entities.Wei
	GasLimit *uint64
	GasPrice *entities.Wei
}

type TransactionReceipt struct {
	TransactionHash   string
	BlockHash         string
	BlockNumber       uint64
	From              string
	To                string
	CumulativeGasUsed *big.Int
	GasUsed           *big.Int
	Value             *entities.Wei
	Logs              []TransactionLog
}

type TransactionLog struct {
	Address     string
	Topics      [][32]byte
	Data        []byte
	BlockNumber uint64
	TxHash      string
	TxIndex     uint
	BlockHash   string
	Index       uint
	Removed     bool
}

type ParsedLog[E any] struct {
	Log    E
	RawLog TransactionLog
}

type BlockInfo struct {
	Hash      string
	Number    uint64
	Timestamp time.Time
	Nonce     uint64
}

func NewTransactionConfig(value *entities.Wei, gasLimit uint64, gasPrice *entities.Wei) TransactionConfig {
	var gas *uint64
	if gasLimit != 0 {
		gas = &gasLimit
	}
	return TransactionConfig{Value: value, GasLimit: gas, GasPrice: gasPrice}
}

type RootstockRpcServer interface {
	EstimateGas(ctx context.Context, addr string, value *entities.Wei, data []byte) (*entities.Wei, error)
	GasPrice(ctx context.Context) (*entities.Wei, error)
	GetHeight(ctx context.Context) (uint64, error)
	GetTransactionReceipt(ctx context.Context, hash string) (TransactionReceipt, error)
	GetBalance(ctx context.Context, address string) (*entities.Wei, error)
	GetBlockByHash(ctx context.Context, hash string) (BlockInfo, error)
	GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (BlockInfo, error)
}

type RootstockWallet interface {
	SendRbtc(ctx context.Context, config TransactionConfig, toAddress string) (string, error)
	GetBalance(ctx context.Context) (*entities.Wei, error)
}
