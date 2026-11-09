package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirst "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
	"time"
)

const (
	requestPegInLogName          = "RequestPegIn"
	requestPegInGasEstimateParts = 10
	requestPegInGasBufferParts   = 2
)

var requestPegInRevertErrors = map[string]error{
	"PegInAlreadyProcessed":     blockchain.PegInAlreadyProcessedError,
	"PegInBelowMinimum":         blockchain.PegInBelowMinimumError,
	"DepositBelowBridgeMinimum": blockchain.DepositBelowBridgeMinimumError,
	"InsufficientConfirmations": blockchain.InsufficientConfirmationsError,
	"UnableToGetConfirmations":  blockchain.UnableToGetConfirmationsError,
	"IncorrectFronting":         blockchain.IncorrectFrontingError,
	"EnforcedPause":             blockchain.ContractPausedError,
}

type peginContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	binding       *bindings.PeginContract
	commitFirst   *commitFirst.PeginCommitFirstContract
	abis          *FlyoverABIs
}

func NewPeginContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	binding *bindings.PeginContract,
	commitFirstBinding *commitFirst.PeginCommitFirstContract,
	abis *FlyoverABIs,
) blockchain.PeginContract {
	return &peginContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
		binding:       binding,
		commitFirst:   commitFirstBinding,
		abis:          abis,
	}
}

func (peginContract *peginContractImpl) GetAddress() string {
	return peginContract.address
}

func (peginContract *peginContractImpl) GetBalance(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	balance, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := peginContract.binding.TryPackGetBalance(parsedAddress)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(peginContract.contract, &bind.CallOpts{}, callData, peginContract.binding.UnpackGetBalance)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}

// Withdraw withdraws the specified amount from the LP's balance in the pegin contract.
// It first performs a dry-run call to check for reverts before submitting the actual transaction.
func (peginContract *peginContractImpl) Withdraw(amount *entities.Wei) error {
	callData, dataErr := peginContract.binding.TryPackWithdraw(amount.AsBigInt())
	if dataErr != nil {
		return dataErr
	}
	_, revert := peginContract.contract.CallRaw(&bind.CallOpts{}, callData)
	parsedRevert, err := ParseRevertReason(peginContract.abis.Flyover, revert)
	if err != nil && parsedRevert == nil {
		return fmt.Errorf("error parsing withdraw result: %w", err)
	} else if parsedRevert != nil {
		return fmt.Errorf("withdraw reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:   peginContract.signer.Address(),
		Signer: peginContract.signer.Sign,
	}

	receipt, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(peginContract.client, peginContract.miningTimeout, "Withdraw", func() (*geth.Transaction, error) {
				return bind.Transact(peginContract.contract, opts, callData)
			})
		})

	if err != nil {
		return fmt.Errorf("withdraw error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("withdraw error: transaction failed")
	}
	return nil
}

func (peginContract *peginContractImpl) PausedStatus() (blockchain.PauseStatus, error) {
	opts := new(bind.CallOpts)
	result, err := rskRetry(
		peginContract.retryParams.Retries,
		peginContract.retryParams.Sleep,
		func() (bindings.PauseStatusOutput, error) {
			callData, dataErr := peginContract.binding.TryPackPauseStatus()
			if dataErr != nil {
				return bindings.PauseStatusOutput{}, dataErr
			}
			return bind.Call(peginContract.contract, opts, callData, peginContract.binding.UnpackPauseStatus)
		},
	)
	if err != nil {
		return blockchain.PauseStatus{}, err
	}
	return blockchain.PauseStatus{
		IsPaused: result.IsPaused,
		Reason:   result.Reason,
		Since:    result.Since,
	}, nil
}

func (peginContract *peginContractImpl) SimulateRequestPegIn(
	txConfig blockchain.TransactionConfig,
	params blockchain.RequestPegInParams,
) error {
	callData, err := peginContract.packRequestPegIn(params)
	if err != nil {
		return err
	}
	_, callErr := peginContract.client.CallContract(
		context.Background(),
		peginContract.requestPegInCallMsg(callData, txConfig.Value),
		nil,
	)
	return peginContract.requestPegInRevertError(callErr)
}

func (peginContract *peginContractImpl) EstimateRequestPegInGas(
	txConfig blockchain.TransactionConfig,
	params blockchain.RequestPegInParams,
) (uint64, error) {
	callData, err := peginContract.packRequestPegIn(params)
	if err != nil {
		return 0, err
	}
	estimated, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (uint64, error) {
			return peginContract.client.EstimateGas(
				context.Background(),
				peginContract.requestPegInCallMsg(callData, txConfig.Value),
			)
		})
	if err != nil {
		return 0, err
	}
	return estimated * (requestPegInGasEstimateParts + requestPegInGasBufferParts) / requestPegInGasEstimateParts, nil
}

func (peginContract *peginContractImpl) RequestPegIn(
	txConfig blockchain.TransactionConfig,
	params blockchain.RequestPegInParams,
) (blockchain.RequestPegInResult, error) {
	callData, err := peginContract.packRequestPegIn(params)
	if err != nil {
		return blockchain.RequestPegInResult{}, err
	}
	opts := &bind.TransactOpts{
		From:   peginContract.signer.Address(),
		Signer: peginContract.signer.Sign,
		Value:  txConfig.Value.AsBigInt(),
	}
	if txConfig.GasLimit != nil {
		opts.GasLimit = *txConfig.GasLimit
	}

	var tx *geth.Transaction
	receipt, err := awaitTx(peginContract.client, peginContract.miningTimeout, requestPegInLogName, func() (*geth.Transaction, error) {
		var txErr error
		tx, txErr = bind.Transact(peginContract.contract, opts, callData)
		return tx, txErr
	})
	if err != nil {
		return requestPegInResultForHash(tx), fmt.Errorf("request pegin error: %w", err)
	}
	if receipt == nil {
		return requestPegInResultForHash(tx), errors.New("request pegin error: incomplete receipt")
	}
	return peginContract.requestPegInResult(tx, receipt, callData, txConfig.Value)
}

func (peginContract *peginContractImpl) requestPegInResult(
	tx *geth.Transaction,
	receipt *geth.Receipt,
	callData []byte,
	value *entities.Wei,
) (blockchain.RequestPegInResult, error) {
	transactionReceipt, err := ParseReceipt(tx, receipt)
	if err != nil {
		return requestPegInResultForHash(tx), fmt.Errorf("request pegin error: %w", err)
	}
	result := blockchain.RequestPegInResult{Receipt: transactionReceipt}
	if receipt.Status == 0 {
		_, callErr := peginContract.client.CallContract(
			context.Background(),
			peginContract.requestPegInCallMsg(callData, value),
			receipt.BlockNumber,
		)
		revertErr := errors.Join(blockchain.TxFailedError, peginContract.requestPegInRevertError(callErr))
		return result, fmt.Errorf("request pegin transaction %s reverted: %w", transactionReceipt.TransactionHash, revertErr)
	}
	result.Event, err = peginContract.pegInRequestedEvent(receipt)
	return result, err
}

func (peginContract *peginContractImpl) pegInRequestedEvent(receipt *geth.Receipt) (blockchain.PegInRequestedEvent, error) {
	contractAddress := common.HexToAddress(peginContract.address)
	eventId := peginContract.abis.PegInCommitFirst.Events[commitFirst.PeginCommitFirstContractPegInRequestedEventName].ID
	for _, eventLog := range receipt.Logs {
		if eventLog != nil && eventLog.Address == contractAddress && len(eventLog.Topics) > 0 && eventLog.Topics[0] == eventId {
			return peginContract.unpackPegInRequested(eventLog)
		}
	}
	return blockchain.PegInRequestedEvent{}, fmt.Errorf(
		"request pegin error: %s event not found",
		commitFirst.PeginCommitFirstContractPegInRequestedEventName,
	)
}

func (peginContract *peginContractImpl) unpackPegInRequested(eventLog *geth.Log) (blockchain.PegInRequestedEvent, error) {
	event, err := peginContract.commitFirst.UnpackPegInRequestedEvent(eventLog)
	if err != nil {
		return blockchain.PegInRequestedEvent{}, fmt.Errorf("request pegin error: %w", err)
	}
	return blockchain.PegInRequestedEvent{
		PegInID:     hex.EncodeToString(event.PegInId[:]),
		Claimer:     event.Claimer.String(),
		RskAddress:  event.RskAddr.String(),
		Amount:      entities.NewBigWei(event.Amount),
		NetToUser:   entities.NewBigWei(event.NetToUser),
		CallSuccess: event.CallSuccess,
	}, nil
}

func (peginContract *peginContractImpl) packRequestPegIn(params blockchain.RequestPegInParams) ([]byte, error) {
	var rskAddress common.Address
	if err := ParseAddress(&rskAddress, params.RskAddress); err != nil {
		return nil, err
	}
	return peginContract.commitFirst.TryPackRequestPegIn(
		rskAddress,
		params.BtcRawTx,
		params.BtcBlockHeaderHash,
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
}

func (peginContract *peginContractImpl) requestPegInCallMsg(callData []byte, value *entities.Wei) ethereum.CallMsg {
	destination := common.HexToAddress(peginContract.address)
	return ethereum.CallMsg{
		From:  peginContract.signer.Address(),
		To:    &destination,
		Data:  callData,
		Value: value.AsBigInt(),
	}
}

func (peginContract *peginContractImpl) requestPegInRevertError(callErr error) error {
	if callErr == nil {
		return nil
	}
	var parseErr error
	for _, contractAbi := range []*abi.ABI{peginContract.abis.PegInCommitFirst, peginContract.abis.Flyover} {
		var parsedRevert *abi.Error
		if parsedRevert, parseErr = ParseRevertReason(contractAbi, callErr); parsedRevert != nil {
			return requestPegInRevertNameError(parsedRevert.Name)
		}
	}
	return fmt.Errorf("request pegin dry run error: %w", errors.Join(callErr, parseErr))
}

func requestPegInRevertNameError(name string) error {
	if mapped, found := requestPegInRevertErrors[name]; found {
		return mapped
	}
	return fmt.Errorf("request pegin reverted with: %s", name)
}

func requestPegInResultForHash(tx *geth.Transaction) blockchain.RequestPegInResult {
	if tx == nil {
		return blockchain.RequestPegInResult{}
	}
	return blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{TransactionHash: tx.Hash().String()},
	}
}
