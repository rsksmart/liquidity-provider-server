package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitfirst "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	log "github.com/sirupsen/logrus"
)

const (
	// registerPeginGasLimit Fixed gas limit for registerPegin function, should change only if the function does
	registerPeginGasLimit = 2500000
	requestPegInGasLimit  = 2500000
)

type peginContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	binding       *bindings.PeginContract
	commitFirst   *commitfirst.PeginCommitFirstContract
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
		commitFirst:   commitfirst.NewPeginCommitFirstContract(),
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

func (peginContract *peginContractImpl) HashPeginQuote(peginQuote quote.PeginQuote) (string, error) {
	var results [32]byte

	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return "", err
	}

	results, err = rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() ([32]byte, error) {
			callData, dataErr := peginContract.binding.TryPackHashPegInQuote(parsedQuote)
			if dataErr != nil {
				return [32]byte{}, dataErr
			}
			return bind.Call(peginContract.contract, &bind.CallOpts{}, callData, peginContract.binding.UnpackHashPegInQuote)
		})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (peginContract *peginContractImpl) HashPeginQuoteEIP712(peginQuote quote.PeginQuote) ([32]byte, error) {
	var result [32]byte

	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return [32]byte{}, err
	}

	result, err = rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() ([32]byte, error) {
			callData, dataErr := peginContract.binding.TryPackHashPegInQuoteEIP712(parsedQuote)
			if dataErr != nil {
				return [32]byte{}, dataErr
			}
			return bind.Call(peginContract.contract, &bind.CallOpts{}, callData, peginContract.binding.UnpackHashPegInQuoteEIP712)
		})
	if err != nil {
		return [32]byte{}, err
	}
	return result, nil
}

func (peginContract *peginContractImpl) CallForUser(txConfig blockchain.TransactionConfig, peginQuote quote.PeginQuote) (blockchain.TransactionReceipt, error) {
	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}

	opts := &bind.TransactOpts{
		GasLimit: *txConfig.GasLimit,
		Value:    txConfig.Value.AsBigInt(),
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
	}

	var tx *geth.Transaction
	receipt, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(peginContract.client, peginContract.miningTimeout, "CallForUser", func() (*geth.Transaction, error) {
				var dataErr, txErr error
				callData, dataErr := peginContract.binding.TryPackCallForUser(parsedQuote)
				if dataErr != nil {
					return nil, dataErr
				}
				tx, txErr = bind.Transact(peginContract.contract, opts, callData)
				return tx, txErr
			})
		})

	if err != nil {
		return blockchain.TransactionReceipt{}, fmt.Errorf("call for user error: %w", err)
	} else if receipt == nil {
		return blockchain.TransactionReceipt{}, errors.New("call for user error: incomplete receipt")
	}

	// Fetch the transaction to get the "To" address and the Value
	toAddress := ""
	txValue := entities.NewWei(0)
	if tx != nil {
		if tx.To() != nil {
			toAddress = tx.To().String()
		}
		txValue = entities.NewBigWei(tx.Value())
	}

	transactionReceipt := blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              peginContract.signer.Address().String(),
		To:                toAddress,
		CumulativeGasUsed: new(big.Int).SetUint64(receipt.CumulativeGasUsed),
		GasUsed:           new(big.Int).SetUint64(receipt.GasUsed),
		Value:             txValue,
		GasPrice:          entities.NewWei(receipt.EffectiveGasPrice.Int64()),
	}

	// Return populated receipt even on revert, but with error
	if receipt.Status == 0 {
		return transactionReceipt, fmt.Errorf("call for user error: transaction reverted (%s)", receipt.TxHash.String())
	}

	return transactionReceipt, nil
}

// TODO: ignore cyclop and funlen added during the merge of the LBC split, the function should be refactored separately
// nolint:cyclop,funlen
func (peginContract *peginContractImpl) RegisterPegin(params blockchain.RegisterPeginParams) (blockchain.TransactionReceipt, error) {
	const (
		waitingForBridgeError = "NotEnoughConfirmations"
	)
	var err error
	var parsedQuote bindings.QuotesPegInQuote
	if parsedQuote, err = parsePeginQuote(params.Quote); err != nil {
		return blockchain.TransactionReceipt{}, err
	}
	log.Infof("Executing RegisterPegIn with params: %s\n", params.String())
	callData, dataErr := peginContract.binding.TryPackRegisterPegIn(parsedQuote, params.QuoteSignature, params.BitcoinRawTransaction, params.PartialMerkleTree, params.BlockHeight)
	if dataErr != nil {
		return blockchain.TransactionReceipt{}, dataErr
	}
	_, revert := peginContract.contract.CallRaw(&bind.CallOpts{}, callData)
	parsedRevert, err := ParseRevertReason(peginContract.abis.PegIn, revert)
	if err != nil && parsedRevert == nil {
		return blockchain.TransactionReceipt{}, fmt.Errorf("error parsing registerPegIn result: %w", err)
	} else if parsedRevert != nil && strings.EqualFold(waitingForBridgeError, parsedRevert.Name) {
		log.Debugln("RegisterPegin: bridge failed to validate BTC transaction. retrying on next confirmation.")
		// allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		return blockchain.TransactionReceipt{}, blockchain.WaitingForBridgeError
	} else if parsedRevert != nil {
		return blockchain.TransactionReceipt{}, fmt.Errorf("registerPegIn reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
		GasLimit: registerPeginGasLimit,
	}

	var tx *geth.Transaction
	receipt, err := awaitTx(peginContract.client, peginContract.miningTimeout, "RegisterPegIn", func() (*geth.Transaction, error) {
		var txErr error
		tx, txErr = bind.Transact(peginContract.contract, opts, callData)
		return tx, txErr
	})

	if err != nil {
		return blockchain.TransactionReceipt{}, fmt.Errorf("register pegin error: %w", err)
	} else if receipt == nil {
		return blockchain.TransactionReceipt{}, errors.New("register pegin error: incomplete receipt")
	}
	// Fetch the transaction to get the "To" address and Value
	toAddress := ""
	txValue := entities.NewWei(0)
	if tx != nil {
		if tx.To() != nil {
			toAddress = tx.To().String()
		}
		txValue = entities.NewBigWei(tx.Value())
	}
	transactionReceipt := blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              peginContract.signer.Address().String(),
		To:                toAddress,
		CumulativeGasUsed: new(big.Int).SetUint64(receipt.CumulativeGasUsed),
		GasUsed:           new(big.Int).SetUint64(receipt.GasUsed),
		Value:             txValue,
		GasPrice:          entities.NewWei(receipt.EffectiveGasPrice.Int64()),
	}
	if receipt.Status == 0 {
		return transactionReceipt, fmt.Errorf("register pegin error: transaction reverted (%s)", receipt.TxHash.String())
	}

	return transactionReceipt, nil
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

// parsePeginQuote parses a quote.PeginQuote into a bindings.QuotesPegInQuote. All BTC address fields support all address types
// except for FedBtcAddress which must be a P2SH address.
func parsePeginQuote(peginQuote quote.PeginQuote) (bindings.QuotesPegInQuote, error) {
	var decodedFederationAddress []byte
	var parsedQuote bindings.QuotesPegInQuote
	var err error

	if err = entities.ValidateStruct(peginQuote); err != nil {
		return bindings.QuotesPegInQuote{}, err
	}

	if decodedFederationAddress, err = bitcoin.DecodeAddressBase58(peginQuote.FedBtcAddress, false); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing federation address: %w", err)
	} else {
		copy(parsedQuote.FedBtcAddress[:], decodedFederationAddress)
	}
	if parsedQuote.BtcRefundAddress, err = bitcoin.DecodeAddress(peginQuote.BtcRefundAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing user btc refund address: %w", err)
	}
	if parsedQuote.LiquidityProviderBtcAddress, err = bitcoin.DecodeAddress(peginQuote.LpBtcAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing btc liquidity provider address: %w", err)
	}

	if err = ParseAddress(&parsedQuote.LbcAddress, peginQuote.LbcAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing lbc address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.LiquidityProviderRskAddress, peginQuote.LpRskAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing liquidity provider rsk address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.RskRefundAddress, peginQuote.RskRefundAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing user rsk refund address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.ContractAddress, peginQuote.ContractAddress); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing destination contract address: %w", err)
	}

	if parsedQuote.Data, err = blockchain.DecodeStringTrimPrefix(peginQuote.Data); err != nil {
		return bindings.QuotesPegInQuote{}, fmt.Errorf("error parsing data: %w", err)
	}

	chainId := new(big.Int)
	parsedQuote.CallFee = peginQuote.CallFee.AsBigInt()
	parsedQuote.PenaltyFee = peginQuote.PenaltyFee.AsBigInt()
	parsedQuote.GasLimit = peginQuote.GasLimit
	parsedQuote.Nonce = peginQuote.Nonce
	parsedQuote.Value = peginQuote.Value.AsBigInt()
	parsedQuote.AgreementTimestamp = peginQuote.AgreementTimestamp
	parsedQuote.CallTime = peginQuote.LpCallTime
	parsedQuote.DepositConfirmations = peginQuote.Confirmations
	parsedQuote.TimeForDeposit = peginQuote.TimeForDeposit
	parsedQuote.GasFee = peginQuote.GasFee.AsBigInt()
	parsedQuote.CallOnRegister = peginQuote.CallOnRegister
	parsedQuote.ChainId = chainId.SetUint64(peginQuote.ChainId)
	return parsedQuote, nil
}

func (peginContract *peginContractImpl) RequestPegIn(params blockchain.RequestPegInParams) (blockchain.RequestPegInResult, error) {
	parsedAddress, value, err := peginContract.prepareRequestPegIn(params)
	if err != nil {
		return blockchain.RequestPegInResult{}, err
	}
	callData, dataErr := peginContract.commitFirst.TryPackRequestPegIn(
		parsedAddress,
		params.BitcoinRawTx,
		[]byte{},
		params.BtcBlockHash,
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
	if dataErr != nil {
		return blockchain.RequestPegInResult{}, dataErr
	}
	if err = peginContract.preflightRequestPegIn(callData); err != nil {
		return blockchain.RequestPegInResult{}, err
	}
	return peginContract.submitRequestPegIn(callData, value)
}

func (peginContract *peginContractImpl) IdentifyRequestPegIn(params blockchain.RequestPegInParams) error {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, params.RskAddress); err != nil {
		return err
	}
	if err := rejectWitnessSerialized(params.BitcoinRawTx); err != nil {
		return err
	}
	callData, dataErr := peginContract.commitFirst.TryPackRequestPegIn(
		parsedAddress,
		params.BitcoinRawTx,
		[]byte{},
		params.BtcBlockHash,
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
	if dataErr != nil {
		return dataErr
	}
	return peginContract.preflightRequestPegIn(callData)
}

func (peginContract *peginContractImpl) UnpackPegInRequested(
	receipt blockchain.TransactionReceipt,
) (*blockchain.PegInRequestedEvent, error) {
	logs := make([]*geth.Log, 0, len(receipt.Logs))
	for _, eventLog := range receipt.Logs {
		copied := eventLog
		logs = append(logs, transactionLogToGeth(copied))
	}
	return unpackPegInRequested(peginContract.commitFirst, &geth.Receipt{Logs: logs})
}

func transactionLogToGeth(eventLog blockchain.TransactionLog) *geth.Log {
	topics := make([]common.Hash, len(eventLog.Topics))
	for i, topic := range eventLog.Topics {
		topics[i] = topic
	}
	return &geth.Log{
		Address:     common.HexToAddress(eventLog.Address),
		Topics:      topics,
		Data:        eventLog.Data,
		BlockNumber: eventLog.BlockNumber,
		TxHash:      common.HexToHash(eventLog.TxHash),
		TxIndex:     eventLog.TxIndex,
		BlockHash:   common.HexToHash(eventLog.BlockHash),
		Index:       eventLog.Index,
		Removed:     eventLog.Removed,
	}
}

func (peginContract *peginContractImpl) prepareRequestPegIn(params blockchain.RequestPegInParams) (common.Address, *entities.Wei, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, params.RskAddress); err != nil {
		return common.Address{}, nil, err
	}
	if err := rejectWitnessSerialized(params.BitcoinRawTx); err != nil {
		return common.Address{}, nil, err
	}
	hardPaused, err := peginContract.IsHardPaused()
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("hard-pause check: %w", err)
	}
	if hardPaused {
		return common.Address{}, nil, blockchain.ErrHardPaused
	}
	value, err := requestPegInValue(params.Amount, params.Fee)
	if err != nil {
		return common.Address{}, nil, err
	}
	return parsedAddress, value, nil
}

func (peginContract *peginContractImpl) submitRequestPegIn(callData []byte, value *entities.Wei) (blockchain.RequestPegInResult, error) {
	opts := &bind.TransactOpts{
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
		GasLimit: requestPegInGasLimit,
		Value:    value.AsBigInt(),
	}
	var tx *geth.Transaction
	receipt, err := awaitTx(peginContract.client, peginContract.miningTimeout, "RequestPegIn", func() (*geth.Transaction, error) {
		var txErr error
		tx, txErr = bind.Transact(peginContract.contract, opts, callData)
		return tx, txErr
	})
	if err != nil {
		return blockchain.RequestPegInResult{}, fmt.Errorf("request pegin error: %w", err)
	}
	if receipt == nil {
		return blockchain.RequestPegInResult{}, errors.New("request pegin error: incomplete receipt")
	}
	transactionReceipt := receiptFromTx(peginContract.signer.Address().String(), tx, receipt)
	if receipt.Status == 0 {
		return blockchain.RequestPegInResult{Receipt: transactionReceipt}, fmt.Errorf("request pegin error: transaction reverted (%s)", receipt.TxHash.String())
	}
	event, err := unpackPegInRequested(peginContract.commitFirst, receipt)
	if err != nil {
		return blockchain.RequestPegInResult{Receipt: transactionReceipt}, err
	}
	return blockchain.RequestPegInResult{Receipt: transactionReceipt, Event: event}, nil
}

func (peginContract *peginContractImpl) IsHardPaused() (bool, error) {
	registryAddr, err := peginContract.pauseRegistryAddress()
	if err != nil {
		return false, err
	}
	level, err := peginContract.pauseLevel(registryAddr)
	if err != nil {
		return false, err
	}
	return level >= blockchain.PauseLevelHard, nil
}

func (peginContract *peginContractImpl) pauseRegistryAddress() (common.Address, error) {
	callData, err := pauseRegistryABI.Pack("pauseRegistry")
	if err != nil {
		return common.Address{}, err
	}
	output, revert := peginContract.contract.CallRaw(&bind.CallOpts{}, callData)
	if revert != nil {
		return common.Address{}, fmt.Errorf("pauseRegistry call: %w", revert)
	}
	registryAddr, err := unpackABIAddress(pauseRegistryABI, "pauseRegistry", output)
	if err != nil {
		return common.Address{}, err
	}
	if registryAddr == (common.Address{}) {
		return common.Address{}, errors.New("pause registry address is zero")
	}
	return registryAddr, nil
}

func (peginContract *peginContractImpl) pauseLevel(registryAddr common.Address) (uint8, error) {
	levelData, err := pauseLevelABI.Pack("pauseLevel")
	if err != nil {
		return 0, err
	}
	levelOutput, err := peginContract.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &registryAddr,
		Data: levelData,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("pauseLevel call: %w", err)
	}
	return unpackABIUint8(pauseLevelABI, "pauseLevel", levelOutput)
}

func (peginContract *peginContractImpl) preflightRequestPegIn(callData []byte) error {
	_, revert := peginContract.contract.CallRaw(&bind.CallOpts{}, callData)
	if revert == nil {
		return nil
	}
	if classified := classifyRequestPegInRevert(peginContract.commitFirst, peginContract.abis.PegInCommitFirst, revert); classified != nil {
		return classified
	}
	return fmt.Errorf("error parsing requestPegIn result: %w", revert)
}

func classifyRequestPegInRevert(commitFirst *commitfirst.PeginCommitFirstContract, commitABI *abi.ABI, revert error) error {
	raw, rawErr := revertDataBytes(revert)
	if rawErr == nil {
		if typed := typedRequestPegInError(commitFirst, raw); typed != nil {
			return typed
		}
	}
	parsedRevert, parseErr := ParseRevertReason(commitABI, revert)
	if parsedRevert == nil {
		_ = parseErr
		return nil
	}
	return mapRequestPegInRevertName(parsedRevert.Name)
}

func typedRequestPegInError(commitFirst *commitfirst.PeginCommitFirstContract, raw []byte) error {
	if len(raw) < 4 {
		return nil
	}
	unpacked, err := commitFirst.UnpackError(raw)
	if err == nil {
		return mapUnpackedRequestPegInError(unpacked)
	}
	return nil
}

func mapUnpackedRequestPegInError(unpacked any) error {
	switch unpacked.(type) {
	case *commitfirst.PeginCommitFirstContractPegInAlreadyProcessed:
		return blockchain.ErrPegInAlreadyProcessed
	case *commitfirst.PeginCommitFirstContractAddressNotRegistered:
		return blockchain.ErrAddressNotRegistered
	case *commitfirst.PeginCommitFirstContractDepositOutputNotFound:
		return blockchain.ErrDepositOutputNotFound
	case *commitfirst.PeginCommitFirstContractInsufficientConfirmations:
		return blockchain.ErrInsufficientConfirmations
	case *commitfirst.PeginCommitFirstContractIncorrectFronting:
		return blockchain.ErrIncorrectFronting
	default:
		return nil
	}
}

func mapRequestPegInRevertName(name string) error {
	switch name {
	case "PegInAlreadyProcessed":
		return blockchain.ErrPegInAlreadyProcessed
	case "AddressNotRegistered":
		return blockchain.ErrAddressNotRegistered
	case "DepositOutputNotFound":
		return blockchain.ErrDepositOutputNotFound
	case "InsufficientConfirmations":
		return blockchain.ErrInsufficientConfirmations
	case "IncorrectFronting":
		return blockchain.ErrIncorrectFronting
	default:
		return fmt.Errorf("requestPegIn reverted with: %s", name)
	}
}

func rejectWitnessSerialized(rawTx []byte) error {
	if len(rawTx) < 6 {
		return blockchain.ErrWitnessSerializedTxNotAccepted
	}
	if rawTx[4] == 0x00 && rawTx[5] == 0x01 {
		return blockchain.ErrWitnessSerializedTxNotAccepted
	}
	return nil
}

func requestPegInValue(amount, fee *entities.Wei) (*entities.Wei, error) {
	if amount == nil || fee == nil {
		return nil, blockchain.ErrIncorrectFronting
	}
	if amount.Cmp(fee) < 0 {
		return nil, blockchain.ErrIncorrectFronting
	}
	return new(entities.Wei).Sub(amount, fee), nil
}

func unpackPegInRequested(commitFirst *commitfirst.PeginCommitFirstContract, receipt *geth.Receipt) (*blockchain.PegInRequestedEvent, error) {
	for _, eventLog := range receipt.Logs {
		if eventLog == nil {
			continue
		}
		unpacked, err := commitFirst.UnpackPegInRequestedEvent(eventLog)
		if err != nil {
			continue
		}
		return &blockchain.PegInRequestedEvent{
			PegInId:     unpacked.PegInId,
			Claimer:     unpacked.Claimer.Hex(),
			RskAddress:  unpacked.RskAddr.Hex(),
			Amount:      entities.NewBigWei(unpacked.Amount),
			NetToUser:   entities.NewBigWei(unpacked.NetToUser),
			CallSuccess: unpacked.CallSuccess,
		}, nil
	}
	return nil, errors.New("request pegin error: PegInRequested event not found")
}

func receiptFromTx(from string, tx *geth.Transaction, receipt *geth.Receipt) blockchain.TransactionReceipt {
	toAddress := ""
	txValue := entities.NewWei(0)
	if tx != nil {
		if tx.To() != nil {
			toAddress = tx.To().String()
		}
		txValue = entities.NewBigWei(tx.Value())
	}
	gasPrice := entities.NewWei(0)
	if receipt.EffectiveGasPrice != nil {
		gasPrice = entities.NewBigWei(receipt.EffectiveGasPrice)
	}
	return blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		Status:            receipt.Status,
		From:              from,
		To:                toAddress,
		CumulativeGasUsed: new(big.Int).SetUint64(receipt.CumulativeGasUsed),
		GasUsed:           new(big.Int).SetUint64(receipt.GasUsed),
		Value:             txValue,
		GasPrice:          gasPrice,
	}
}

func unpackABIAddress(parsed abi.ABI, name string, output []byte) (common.Address, error) {
	unpacked, err := parsed.Unpack(name, output)
	if err != nil {
		return common.Address{}, err
	}
	if len(unpacked) == 0 {
		return common.Address{}, errors.New("empty ABI output")
	}
	addr, ok := unpacked[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("unexpected ABI type %T", unpacked[0])
	}
	return addr, nil
}

func unpackABIUint8(parsed abi.ABI, name string, output []byte) (uint8, error) {
	unpacked, err := parsed.Unpack(name, output)
	if err != nil {
		return 0, err
	}
	if len(unpacked) == 0 {
		return 0, errors.New("empty ABI output")
	}
	level, ok := unpacked[0].(uint8)
	if !ok {
		return 0, fmt.Errorf("unexpected ABI type %T", unpacked[0])
	}
	return level, nil
}

var (
	pauseRegistryABI = mustParseJSONABI(`[{"inputs":[],"name":"pauseRegistry","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`)
	pauseLevelABI    = mustParseJSONABI(`[{"inputs":[],"name":"pauseLevel","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"}]`)
)

func mustParseJSONABI(raw string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(raw))
	if err != nil {
		panic(err)
	}
	return parsed
}
