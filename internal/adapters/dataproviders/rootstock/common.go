package rootstock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

const (
	rpcCallRetryMax   = 3
	rpcCallRetrySleep = 1 * time.Minute
)

var ErrShortRevertData = errors.New("revert data shorter than ABI selector")

var DefaultRetryParams = RetryParams{
	Retries: rpcCallRetryMax,
	Sleep:   rpcCallRetrySleep,
}

type RskClient struct {
	client RpcClientBinding
}

type RetryParams struct {
	Retries uint
	Sleep   time.Duration
}

func NewRskClient(client RpcClientBinding) *RskClient {
	return &RskClient{client: client}
}

func (c *RskClient) Rpc() RpcClientBinding {
	return c.client
}

func (c *RskClient) Shutdown(endChannel chan<- bool) {
	c.client.Close()
	endChannel <- true
	log.Debug("Disconnected from RSK node")
}

func (c *RskClient) CheckConnection(ctx context.Context) bool {
	_, err := c.client.ChainID(ctx)
	if err != nil {
		log.Error("Error checking RSK node connection: ", err)
	}
	return err == nil
}

type TransactionSigner interface {
	entities.Signer
	Address() common.Address
	Sign(common.Address, *types.Transaction) (*types.Transaction, error)
}

type RskSignerWallet interface {
	blockchain.RootstockWallet
	TransactionSigner
}

func ParseAddress(address *common.Address, textAddress string) error {
	if !common.IsHexAddress(textAddress) {
		return blockchain.InvalidAddressError
	}
	*address = common.HexToAddress(textAddress)
	return nil
}

func rskRetry[R any](retries uint, retrySleep time.Duration, call func() (R, error)) (R, error) {
	var result R
	var err error
	var i uint

	if retries == 0 {
		return call()
	}

	for i = 0; i < retries; i++ {
		result, err = call()
		if err == nil {
			return result, nil
		}
		time.Sleep(retrySleep)
	}
	return result, err
}

func awaitTx(client RpcClientBinding, miningTimeout time.Duration, logName string, txCall func() (*geth.Transaction, error)) (r *geth.Receipt, e error) {
	return AwaitTxWithCtx(client, miningTimeout, logName, context.Background(), txCall)
}

func AwaitTxWithCtx(client RpcClientBinding, miningTimeout time.Duration, logName string, ctx context.Context, txCall func() (*geth.Transaction, error)) (*geth.Receipt, error) {
	var tx *geth.Transaction
	var err error

	log.Infof("Executing %s transaction...", logName)
	deadline, ok := ctx.Deadline()
	if ok {
		log.Debugf("Waiting for transaction to be mined until %v...", deadline)
	}
	tx, err = txCall()
	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, errors.New("invalid transaction")
	}

	ctx, cancel := context.WithTimeout(ctx, miningTimeout)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil || receipt == nil {
		log.Infof("Error waiting for transaction %s (%s) to be mined: %v", logName, tx.Hash().String(), err)
		return nil, err
	}

	if receipt.Status == 1 {
		log.Infof("Transaction %s (%s) executed successfully", logName, tx.Hash().String())
	} else {
		log.Infof("Transaction %s (%s) reverted", logName, tx.Hash().String())
	}
	return receipt, nil
}

const revertSelectorSize = 4

// RevertKind selects how a caller must decode a revert payload.
type RevertKind int

const (
	// RevertNone means the call did not revert.
	RevertNone RevertKind = iota
	// RevertGeneric means the contract reverted with a reason string
	// produced by abi.UnpackRevert, including Solidity panic reasons.
	RevertGeneric
	// RevertCustom means the payload starts with a custom error selector.
	RevertCustom
)

// RevertPayload separates generic Solidity failures from contract-specific
// failures that need generated ABI decoding.
type RevertPayload struct {
	Kind RevertKind
	// Message holds the reason string when Kind is RevertGeneric.
	Message string
	// Data holds the selector and the encoded arguments when Kind is RevertCustom.
	Data []byte
}

// ParseRevert extracts and classifies the revert payload of a failed call. It
// returns ErrShortRevertData when the payload is too short to hold a selector.
// Callers decode a RevertCustom payload with the decoder they need: an ABI
// lookup for the error name, or a generated UnpackError for the typed error.
func ParseRevert(err error) (RevertPayload, error) {
	if err == nil {
		return RevertPayload{Kind: RevertNone}, nil
	}

	decoded, extractErr := revertDataBytes(err)
	if extractErr != nil {
		return RevertPayload{}, extractErr
	}

	if message, unpackErr := abi.UnpackRevert(decoded); unpackErr == nil {
		return RevertPayload{Kind: RevertGeneric, Message: message}, nil
	}

	if len(decoded) < revertSelectorSize {
		return RevertPayload{}, fmt.Errorf("%w: %w", ErrShortRevertData, err)
	}

	return RevertPayload{Kind: RevertCustom, Data: decoded}, nil
}

func ParseRevertReason(contractAbi *abi.ABI, err error) (*abi.Error, error) {
	payload, parseErr := ParseRevert(err)
	if parseErr != nil {
		return nil, parseErr
	}
	if payload.Kind == RevertNone {
		return nil, nil
	}
	if payload.Kind == RevertGeneric {
		return nil, fmt.Errorf("found generic error: %s", payload.Message)
	}

	var selectorBytes [revertSelectorSize]byte
	copy(selectorBytes[:], payload.Data[:revertSelectorSize])
	parsedError, abiErr := contractAbi.ErrorByID(selectorBytes)
	if abiErr != nil {
		return nil, fmt.Errorf("error decoding data using ABI: %w", abiErr)
	}
	return parsedError, nil
}

func revertDataBytes(err error) ([]byte, error) {
	const errorTemplate = "no data to recover in error: %w"
	if err == nil {
		return nil, nil
	}
	var dataError rpc.DataError
	if !errors.As(err, &dataError) {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	revertData, ok := dataError.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf(errorTemplate, dataError)
	}
	decoded, decodeErr := hexutil.Decode(revertData)
	if decodeErr != nil {
		return nil, fmt.Errorf("error decoding data: %w", decodeErr)
	}
	return decoded, nil
}
