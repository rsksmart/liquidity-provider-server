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
	RskChainHeightErrorTemplate = "error getting Rootstock chain height: %v"
)

const (
	RskZeroAddress = "0x0000000000000000000000000000000000000000"
)

var (
	rskAddressRegex       = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	WaitingForBridgeError = errors.New("waiting for rootstock bridge")
	InvalidAddressError   = errors.New("invalid rootstock address")
	ContractPausedError   = errors.New("contract is paused")
	TxFailedError         = errors.New("transaction failed")
	// AlreadyClaimedError signals that another LP already claimed/processed this peg-in.
	// It is the expected outcome of the first-mined-wins race in the commit-first peg-in
	// path (DoS-removal redesign) and must be treated as benign by claim callers.
	AlreadyClaimedError = errors.New("peg-in already claimed by another provider")
	// AddressNotRegisteredError signals the RSK address is not (yet) registered in the
	// PegInAddressRegistry, so a claim cannot proceed.
	AddressNotRegisteredError = errors.New("rsk address is not registered")
)

type RskContracts struct {
	Bridge               rootstock.Bridge
	PegIn                PeginContract
	PegOut               PegoutContract
	CollateralManagement CollateralManagementContract
	Discovery            DiscoveryContract
	// PegInAddressRegistry and FlyoverConfigurations back the commit-first peg-in path
	// (DoS-removal redesign, EPICs E1/E2/E5). They are optional in legacy deployments,
	// so consumers must nil-check before use.
	PegInAddressRegistry  PegInAddressRegistryContract
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
