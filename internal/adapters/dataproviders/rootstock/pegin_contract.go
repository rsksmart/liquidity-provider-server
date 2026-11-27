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
	requestPegInGasLimitPercent = 120
	pegInRequestedDataSize      = 96
	resolvePegInGasLimit        = 2500000
	pegInResolvedDataSize       = 128
)

type pegInFunction struct {
	logName      string
	errorPrefix  string
	revertErrors map[string]error
}

var requestPegInFunction = pegInFunction{
	logName:     "RequestPegIn",
	errorPrefix: "request pegin",
	revertErrors: map[string]error{
		"PegInAlreadyProcessed":     blockchain.PegInAlreadyProcessedError,
		"PegInBelowMinimum":         blockchain.PegInBelowMinimumError,
		"DepositBelowBridgeMinimum": blockchain.DepositBelowBridgeMinimumError,
		"InsufficientConfirmations": blockchain.InsufficientConfirmationsError,
		"UnableToGetConfirmations":  blockchain.UnableToGetConfirmationsError,
		"IncorrectFronting":         blockchain.IncorrectFrontingError,
		"EnforcedPause":             blockchain.ContractPausedError,
	},
}

var resolvePegInFunction = pegInFunction{
	logName:     "ResolvePegIn",
	errorPrefix: "resolve pegin",
	revertErrors: map[string]error{
		"PegInAlreadyProcessed":  blockchain.PegInAlreadyProcessedError,
		"BridgeValidationFailed": blockchain.WaitingForBridgeError,
		"UnexpectedBridgeError":  blockchain.BridgeRejectedPegInError,
		"EnforcedPause":          blockchain.ContractPausedError,
	},
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
		peginContract.pegInCallMsg(callData, txConfig.Value.AsBigInt()),
		nil,
	)
	return peginContract.pegInRevertError(requestPegInFunction, callErr)
}

func (peginContract *peginContractImpl) EstimateRequestPegInGas(
	txConfig blockchain.TransactionConfig,
	params blockchain.RequestPegInParams,
) (uint64, error) {
	callData, err := peginContract.packRequestPegIn(params)
	if err != nil {
		return 0, err
	}
	estimated, err := peginContract.client.EstimateGas(
		context.Background(),
		peginContract.pegInCallMsg(callData, txConfig.Value.AsBigInt()),
	)
	if err != nil {
		return 0, peginContract.pegInRevertError(requestPegInFunction, err)
	}
	return estimated * requestPegInGasLimitPercent / 100, nil
}

func (peginContract *peginContractImpl) RequestPegIn(
	txConfig blockchain.TransactionConfig,
	params blockchain.RequestPegInParams,
) (blockchain.TransactionReceipt, error) {
	callData, err := peginContract.packRequestPegIn(params)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}
	opts := &bind.TransactOpts{
		From:   peginContract.signer.Address(),
		Signer: peginContract.signer.Sign,
		Value:  txConfig.Value.AsBigInt(),
	}
	if txConfig.GasLimit != nil {
		opts.GasLimit = *txConfig.GasLimit
	}
	if txConfig.GasPrice != nil {
		opts.GasPrice = txConfig.GasPrice.AsBigInt()
	}
	return peginContract.sendPegInTx(requestPegInFunction, opts, callData)
}

func (peginContract *peginContractImpl) sendPegInTx(
	function pegInFunction,
	opts *bind.TransactOpts,
	callData []byte,
) (blockchain.TransactionReceipt, error) {
	var tx *geth.Transaction
	receipt, err := awaitTx(peginContract.client, peginContract.miningTimeout, function.logName, func() (*geth.Transaction, error) {
		var txErr error
		tx, txErr = bind.Transact(peginContract.contract, opts, callData)
		return tx, txErr
	})
	if err != nil {
		return receiptForHash(tx), fmt.Errorf("%s error: %w", function.errorPrefix, err)
	}
	if receipt == nil {
		return receiptForHash(tx), fmt.Errorf("%s error: incomplete receipt", function.errorPrefix)
	}
	return peginContract.minedReceipt(function, tx, receipt)
}

func (peginContract *peginContractImpl) minedReceipt(
	function pegInFunction,
	tx *geth.Transaction,
	receipt *geth.Receipt,
) (blockchain.TransactionReceipt, error) {
	transactionReceipt, err := ParseReceipt(tx, receipt)
	if err != nil {
		return receiptForHash(tx), fmt.Errorf("%s error: %w", function.errorPrefix, err)
	}
	if receipt.Status == 0 {
		return transactionReceipt, peginContract.minedRevertError(function, transactionReceipt.TransactionHash, receipt.BlockNumber, tx)
	}
	return transactionReceipt, nil
}

func (peginContract *peginContractImpl) minedRevertError(
	function pegInFunction,
	txHash string,
	block *big.Int,
	tx *geth.Transaction,
) error {
	_, callErr := peginContract.client.CallContract(
		context.Background(),
		peginContract.pegInCallMsg(tx.Data(), tx.Value()),
		block,
	)
	revertErr := errors.Join(blockchain.TxFailedError, peginContract.pegInRevertError(function, callErr))
	return fmt.Errorf("%s transaction %s reverted: %w", function.errorPrefix, txHash, revertErr)
}

func (peginContract *peginContractImpl) PegInRequestedEvent(receipt blockchain.TransactionReceipt) (blockchain.PegInRequestedEvent, error) {
	const eventName = commitFirst.PeginCommitFirstContractPegInRequestedEventName
	eventLog, err := peginContract.pegInEventLog(receipt, eventName, pegInRequestedDataSize)
	if err != nil {
		return blockchain.PegInRequestedEvent{}, err
	}
	event, err := peginContract.commitFirst.UnpackPegInRequestedEvent(eventLog)
	if err != nil {
		return blockchain.PegInRequestedEvent{}, fmt.Errorf("%s event in transaction %s: %w", eventName, receipt.TransactionHash, err)
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

func (peginContract *peginContractImpl) pegInEventLog(
	receipt blockchain.TransactionReceipt,
	eventName string,
	dataSize int,
) (*geth.Log, error) {
	eventId := peginContract.abis.PegInCommitFirst.Events[eventName].ID
	eventLog, found := findEventLog(receipt.Logs, peginContract.address, eventId)
	if !found {
		return nil, fmt.Errorf("%s event not found in transaction %s", eventName, receipt.TransactionHash)
	}
	if len(eventLog.Data) != dataSize {
		return nil, fmt.Errorf(
			"%s event in transaction %s has %d data bytes, expected %d",
			eventName, receipt.TransactionHash, len(eventLog.Data), dataSize,
		)
	}
	topics := make([]common.Hash, len(eventLog.Topics))
	for i, topic := range eventLog.Topics {
		topics[i] = topic
	}
	return &geth.Log{
		Address: common.HexToAddress(eventLog.Address),
		Topics:  topics,
		Data:    eventLog.Data,
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

func (peginContract *peginContractImpl) pegInCallMsg(callData []byte, value *big.Int) ethereum.CallMsg {
	destination := common.HexToAddress(peginContract.address)
	return ethereum.CallMsg{
		From:  peginContract.signer.Address(),
		To:    &destination,
		Data:  callData,
		Value: value,
	}
}

func (peginContract *peginContractImpl) pegInRevertError(function pegInFunction, callErr error) error {
	if callErr == nil {
		return nil
	}
	var parseErr error
	for _, contractAbi := range []*abi.ABI{peginContract.abis.PegInCommitFirst, peginContract.abis.Flyover} {
		var parsedRevert *abi.Error
		if parsedRevert, parseErr = ParseRevertReason(contractAbi, callErr); parsedRevert != nil {
			return function.revertNameError(parsedRevert.Name)
		}
	}
	return fmt.Errorf("%s call error: %w", function.errorPrefix, errors.Join(callErr, parseErr))
}

func (function pegInFunction) revertNameError(name string) error {
	if mapped, found := function.revertErrors[name]; found {
		return mapped
	}
	return fmt.Errorf("%s reverted with: %s", function.errorPrefix, name)
}

func receiptForHash(tx *geth.Transaction) blockchain.TransactionReceipt {
	if tx == nil {
		return blockchain.TransactionReceipt{}
	}
	return blockchain.TransactionReceipt{TransactionHash: tx.Hash().String()}
}

func (peginContract *peginContractImpl) SimulateResolvePegIn(params blockchain.ResolvePegInParams) error {
	callData, err := peginContract.packResolvePegIn(params)
	if err != nil {
		return err
	}
	_, callErr := peginContract.client.CallContract(context.Background(), peginContract.pegInCallMsg(callData, nil), nil)
	return peginContract.pegInRevertError(resolvePegInFunction, callErr)
}

func (peginContract *peginContractImpl) ResolvePegIn(params blockchain.ResolvePegInParams) (blockchain.TransactionReceipt, error) {
	callData, err := peginContract.packResolvePegIn(params)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}
	opts := &bind.TransactOpts{
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
		GasLimit: resolvePegInGasLimit,
	}
	return peginContract.sendPegInTx(resolvePegInFunction, opts, callData)
}

func (peginContract *peginContractImpl) PegInResolvedEvent(receipt blockchain.TransactionReceipt) (blockchain.PegInResolvedEvent, error) {
	const eventName = commitFirst.PeginCommitFirstContractPegInResolvedEventName
	eventLog, err := peginContract.pegInEventLog(receipt, eventName, pegInResolvedDataSize)
	if err != nil {
		return blockchain.PegInResolvedEvent{}, err
	}
	event, err := peginContract.commitFirst.UnpackPegInResolvedEvent(eventLog)
	if err != nil {
		return blockchain.PegInResolvedEvent{}, fmt.Errorf("%s event in transaction %s: %w", eventName, receipt.TransactionHash, err)
	}
	return blockchain.PegInResolvedEvent{
		PegInID:       hex.EncodeToString(event.PegInId[:]),
		Claimer:       event.Claimer.String(),
		Registrant:    event.Registrant.String(),
		Released:      entities.NewBigWei(event.Released),
		ClaimerPayout: entities.NewBigWei(event.ClaimerPayout),
		RegistrantFee: entities.NewBigWei(event.RegistrantFee),
		UserPayout:    entities.NewBigWei(event.UserPayout),
	}, nil
}

func (peginContract *peginContractImpl) packResolvePegIn(params blockchain.ResolvePegInParams) ([]byte, error) {
	var rskAddress common.Address
	if err := ParseAddress(&rskAddress, params.RskAddress); err != nil {
		return nil, err
	}
	if params.BlockHeight == nil {
		return nil, fmt.Errorf("%s error: missing block height", resolvePegInFunction.errorPrefix)
	}
	return peginContract.commitFirst.TryPackResolvePegIn(
		rskAddress,
		params.BtcRawTx,
		params.PartialMerkleTree,
		params.BlockHeight,
	)
}
