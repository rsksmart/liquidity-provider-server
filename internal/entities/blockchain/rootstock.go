package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

const (
	RskZeroAddress     = "0x0000000000000000000000000000000000000000"
	SuccessfulTxStatus = 1
)

var (
	rskAddressRegex                   = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	WaitingForBridgeError             = errors.New("waiting for rootstock bridge")
	InvalidAddressError               = errors.New("invalid rootstock address")
	ContractPausedError               = errors.New("contract is paused")
	TxFailedError                     = errors.New("transaction failed")
	ErrPegInAlreadyProcessed          = errors.New("peg-in already processed")
	ErrAddressNotRegistered           = errors.New("address not registered")
	ErrDepositOutputNotFound          = errors.New("deposit output not found")
	ErrInsufficientConfirmations      = errors.New("insufficient confirmations")
	ErrIncorrectFronting              = errors.New("incorrect fronting")
	ErrWitnessSerializedTxNotAccepted = errors.New("witness-serialized tx not accepted")
	ErrTransactionReceiptNotFound     = errors.New("transaction receipt not found")
)

type RskContracts struct {
	Bridge                rootstock.Bridge
	PegIn                 PeginContract
	PegOut                PegoutContract
	CollateralManagement  CollateralManagementContract
	Discovery             DiscoveryContract
	PegInAddressRegistry  PegInAddressRegistryContract
	PauseRegistry         PauseRegistryContract
	FlyoverConfigurations FlyoverConfigurationsContract
}

func DecodeStringTrimPrefix(hexString string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(hexString, "0x"))
}
func IsRskAddress(address string) bool {
	return rskAddressRegex.MatchString(address)
}

// NormalizeRskAddress returns the canonical lowercase 0x-prefixed 40-hex form for RSK account addresses.
func NormalizeRskAddress(addr string) (string, error) {
	trimmed := strings.TrimSpace(addr)
	if !IsRskAddress(trimmed) {
		return "", fmt.Errorf("%w: %q", InvalidAddressError, addr)
	}
	return strings.ToLower(trimmed), nil
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
	GasPrice          *entities.Wei
	Status            uint64
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
	Hash       string
	ParentHash string
	Number     uint64
	Timestamp  time.Time
	Nonce      uint64
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
	ChainId(ctx context.Context) (uint64, error)
	PeerCount(ctx context.Context) (uint64, error)
}

type RootstockWallet interface {
	SendRbtc(ctx context.Context, config TransactionConfig, toAddress string) (TransactionReceipt, error)
	GetBalance(ctx context.Context) (*entities.Wei, error)
}

// RejectedPegoutReason is a display label decoded from the Bridge release_request_rejected event.
type RejectedPegoutReason string

const (
	RejectedPegoutReasonUnknown        RejectedPegoutReason = "unknown"
	RejectedPegoutReasonLowAmount      RejectedPegoutReason = "low_amount"
	RejectedPegoutReasonCallerContract RejectedPegoutReason = "caller_contract"
	RejectedPegoutReasonFeeAboveValue  RejectedPegoutReason = "fee_above_value"
)
