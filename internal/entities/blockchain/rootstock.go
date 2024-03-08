package blockchain

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"math/big"
	"regexp"
	"strings"
)

var (
	rskAddressRegex       = regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	WaitingForBridgeError = errors.New("waiting for rootstock bridge")
	InvalidAddressError   = errors.New("invalid rootstock address")
)

type RskContracts struct {
	Bridge       RootstockBridge
	Lbc          LiquidityBridgeContract
	FeeCollector FeeCollector
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
}

type RootstockWallet interface {
	SendRbtc(ctx context.Context, config TransactionConfig, toAddress string) (string, error)
}

type FlyoverDerivationArgs struct {
	FedInfo              FederationInfo
	LbcAdress            []byte
	UserBtcRefundAddress []byte
	LpBtcAddress         []byte
	QuoteHash            []byte
}

type FlyoverDerivation struct {
	Address      string
	RedeemScript string
}

type RootstockBridge interface {
	GetAddress() string
	GetFedAddress() (string, error)
	GetMinimumLockTxValue() (*entities.Wei, error)
	GetFlyoverDerivationAddress(args FlyoverDerivationArgs) (FlyoverDerivation, error)
	GetRequiredTxConfirmations() uint64
	FetchFederationInfo() (FederationInfo, error)
}

type FederationInfo struct {
	FedSize              int64
	FedThreshold         int64
	PubKeys              []string
	FedAddress           string
	ActiveFedBlockHeight int64
	IrisActivationHeight int64
	ErpKeys              []string
}
