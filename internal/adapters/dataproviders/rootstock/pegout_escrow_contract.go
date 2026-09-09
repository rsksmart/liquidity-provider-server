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
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout_escrow"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	log "github.com/sirupsen/logrus"
)

var errRequestHashLength = errors.New("request hash must be 32 bytes long")

type pegOutEscrowContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	binding       *bindings.PegOutEscrowContract
	abis          *FlyoverABIs
}

func NewPegOutEscrowContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	binding *bindings.PegOutEscrowContract,
	abis *FlyoverABIs,
) blockchain.PegOutEscrowContract {
	return &pegOutEscrowContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
		binding:       binding,
		abis:          abis,
	}
}

func (escrow *pegOutEscrowContractImpl) GetAddress() string {
	return escrow.address
}

func (escrow *pegOutEscrowContractImpl) GetPegOutState(requestHash string) (blockchain.EscrowedPegOutState, error) {
	hash, err := parseRequestHash(requestHash)
	if err != nil {
		return blockchain.EscrowedPegOutStateNone, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(escrow.retryParams.Retries, escrow.retryParams.Sleep,
		func() (uint8, error) {
			callData, dataErr := escrow.binding.TryPackGetPegOutState(hash)
			if dataErr != nil {
				return 0, dataErr
			}
			return bind.Call(escrow.contract, opts, callData, escrow.binding.UnpackGetPegOutState)
		})
	if err != nil {
		return blockchain.EscrowedPegOutStateNone, err
	}
	return blockchain.EscrowedPegOutState(result), nil
}

func (escrow *pegOutEscrowContractImpl) GetPegOutQuote(requestHash string) (quote.PegoutQuote, error) {
	hash, err := parseRequestHash(requestHash)
	if err != nil {
		return quote.PegoutQuote{}, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(escrow.retryParams.Retries, escrow.retryParams.Sleep,
		func() (bindings.QuotesPegOutQuote, error) {
			callData, dataErr := escrow.binding.TryPackGetPegOutQuote(hash)
			if dataErr != nil {
				return bindings.QuotesPegOutQuote{}, dataErr
			}
			return bind.Call(escrow.contract, opts, callData, escrow.binding.UnpackGetPegOutQuote)
		})
	if err != nil {
		return quote.PegoutQuote{}, err
	}
	return mapEscrowPegOutQuote(result), nil
}

func (escrow *pegOutEscrowContractImpl) GetMaxMinerFee(requestHash string) (*entities.Wei, error) {
	hash, err := parseRequestHash(requestHash)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(escrow.retryParams.Retries, escrow.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := escrow.binding.TryPackGetMaxMinerFee(hash)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(escrow.contract, opts, callData, escrow.binding.UnpackGetMaxMinerFee)
		})
	if err != nil {
		return nil, err
	}
	return bigIntToWei(result), nil
}

func (escrow *pegOutEscrowContractImpl) RestrictedUntil(lpAddress string) (uint64, error) {
	var parsedAddress common.Address
	if err := ParseAddress(&parsedAddress, lpAddress); err != nil {
		return 0, err
	}
	opts := &bind.CallOpts{}
	result, err := rskRetry(escrow.retryParams.Retries, escrow.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := escrow.binding.TryPackRestrictedUntil(parsedAddress)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(escrow.contract, opts, callData, escrow.binding.UnpackRestrictedUntil)
		})
	if err != nil {
		return 0, err
	}
	return bigIntToUint64(result), nil
}

func (escrow *pegOutEscrowContractImpl) EstimateClaimPegOut(requestHash string, signature []byte) (*entities.Wei, error) {
	hash, err := parseRequestHash(requestHash)
	if err != nil {
		return nil, err
	}
	callData, dataErr := escrow.binding.TryPackClaimPegOut(hash, signature)
	if dataErr != nil {
		return nil, dataErr
	}
	destination := common.HexToAddress(escrow.address)
	tx := ethereum.CallMsg{
		From: escrow.signer.Address(),
		To:   &destination,
		Data: callData,
	}
	result, err := rskRetry(escrow.retryParams.Retries, escrow.retryParams.Sleep,
		func() (uint64, error) {
			return escrow.client.EstimateGas(context.Background(), tx)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewUWei(result), nil
}

func (escrow *pegOutEscrowContractImpl) ClaimPegOut(
	txConfig blockchain.TransactionConfig,
	requestHash string,
	signature []byte,
) (blockchain.TransactionReceipt, error) {
	hash, err := parseRequestHash(requestHash)
	if err != nil {
		return blockchain.TransactionReceipt{}, err
	}
	callData, dataErr := escrow.binding.TryPackClaimPegOut(hash, signature)
	if dataErr != nil {
		return blockchain.TransactionReceipt{}, dataErr
	}

	opts := &bind.TransactOpts{
		From:   escrow.signer.Address(),
		Signer: escrow.signer.Sign,
	}
	if txConfig.GasLimit != nil {
		opts.GasLimit = *txConfig.GasLimit
	}

	var tx *geth.Transaction
	receipt, err := awaitTx(escrow.client, escrow.miningTimeout, "ClaimPegOut", func() (*geth.Transaction, error) {
		var txErr error
		tx, txErr = bind.Transact(escrow.contract, opts, callData)
		return tx, txErr
	})
	if err != nil {
		return blockchain.TransactionReceipt{}, fmt.Errorf("claim peg out error: %w", err)
	}
	if receipt == nil {
		return blockchain.TransactionReceipt{}, errors.New("claim peg out error: incomplete receipt")
	}

	transactionReceipt := escrow.receiptFromGeth(receipt, tx)
	if receipt.Status == 0 {
		return transactionReceipt, fmt.Errorf("claim peg out error: transaction reverted (%s)", receipt.TxHash.String())
	}
	return transactionReceipt, nil
}

func (escrow *pegOutEscrowContractImpl) GetPegOutRequestedEvents(
	ctx context.Context,
	fromBlock uint64,
	toBlock *uint64,
) ([]blockchain.PegOutRequested, error) {
	events, err := filterBoundEvents(ctx, escrow.contract, fromBlock, toBlock, escrow.binding.UnpackPegOutRequestedEvent, "PegOutRequested")
	if err != nil {
		return nil, err
	}
	result := make([]blockchain.PegOutRequested, 0, len(events))
	for _, event := range events {
		result = append(result, blockchain.PegOutRequested{
			RequestHash:        hex.EncodeToString(event.RequestHash[:]),
			RefundAddress:      event.RefundAddress.String(),
			Amount:             bigIntToWei(event.Amount),
			DestinationAddress: event.DestinationAddress,
			TxHash:             event.Raw.TxHash.String(),
			BlockNumber:        event.Raw.BlockNumber,
		})
	}
	return result, nil
}

func (escrow *pegOutEscrowContractImpl) GetPegOutClaimedEvents(
	ctx context.Context,
	fromBlock uint64,
	toBlock *uint64,
) ([]blockchain.PegOutClaimed, error) {
	events, err := filterBoundEvents(ctx, escrow.contract, fromBlock, toBlock, escrow.binding.UnpackPegOutClaimedEvent, "PegOutClaimed")
	if err != nil {
		return nil, err
	}
	result := make([]blockchain.PegOutClaimed, 0, len(events))
	for _, event := range events {
		result = append(result, blockchain.PegOutClaimed{
			LpAddress:   event.LpAddress.String(),
			RequestHash: hex.EncodeToString(event.RequestHash[:]),
			TxHash:      event.Raw.TxHash.String(),
			BlockNumber: event.Raw.BlockNumber,
		})
	}
	return result, nil
}

func (escrow *pegOutEscrowContractImpl) GetPegOutCancelledEvents(
	ctx context.Context,
	fromBlock uint64,
	toBlock *uint64,
) ([]blockchain.PegOutCancelled, error) {
	events, err := filterBoundEvents(ctx, escrow.contract, fromBlock, toBlock, escrow.binding.UnpackPegOutCancelledEvent, "PegOutCancelled")
	if err != nil {
		return nil, err
	}
	result := make([]blockchain.PegOutCancelled, 0, len(events))
	for _, event := range events {
		result = append(result, blockchain.PegOutCancelled{
			RequestHash: hex.EncodeToString(event.RequestHash[:]),
			TxHash:      event.Raw.TxHash.String(),
			BlockNumber: event.Raw.BlockNumber,
		})
	}
	return result, nil
}

func (escrow *pegOutEscrowContractImpl) receiptFromGeth(receipt *geth.Receipt, tx *geth.Transaction) blockchain.TransactionReceipt {
	toAddress := ""
	txValue := entities.NewWei(0)
	if tx != nil && tx.To() != nil {
		toAddress = tx.To().String()
		txValue = entities.NewBigWei(tx.Value())
	}
	gasPrice := entities.NewWei(0)
	if receipt.EffectiveGasPrice != nil {
		gasPrice = entities.NewWei(receipt.EffectiveGasPrice.Int64())
	}
	return blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              escrow.signer.Address().String(),
		To:                toAddress,
		CumulativeGasUsed: new(big.Int).SetUint64(receipt.CumulativeGasUsed),
		GasUsed:           new(big.Int).SetUint64(receipt.GasUsed),
		Value:             txValue,
		GasPrice:          gasPrice,
		Logs:              convertReceiptLogs(receipt),
	}
}

func filterBoundEvents[E bind.ContractEvent](
	ctx context.Context,
	contract *bind.BoundContract,
	fromBlock uint64,
	toBlock *uint64,
	unpack func(*geth.Log) (*E, error),
	eventName string,
) ([]E, error) {
	result := make([]E, 0)
	iterator, err := bind.FilterEvents(
		contract,
		&bind.FilterOpts{Start: fromBlock, End: toBlock, Context: ctx},
		unpack,
	)
	defer func() {
		if iterator == nil {
			return
		}
		if iteratorError := iterator.Close(); iteratorError != nil {
			log.Error("Error closing ", eventName, " event iterator: ", iteratorError)
		}
	}()
	if err != nil || iterator == nil {
		return nil, err
	}
	for iterator.Next() {
		result = append(result, *iterator.Value())
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

func parseRequestHash(requestHash string) ([32]byte, error) {
	var hash [32]byte
	hashBytes, err := hex.DecodeString(strings.TrimPrefix(requestHash, "0x"))
	if err != nil {
		return hash, fmt.Errorf("invalid request hash format: %w", err)
	}
	if len(hashBytes) != 32 {
		return hash, errRequestHashLength
	}
	copy(hash[:], hashBytes)
	return hash, nil
}

func mapEscrowPegOutQuote(raw bindings.QuotesPegOutQuote) quote.PegoutQuote {
	return quote.PegoutQuote{
		LbcAddress:            raw.LbcAddress.String(),
		LpRskAddress:          raw.LpRskAddress.String(),
		BtcRefundAddress:      hex.EncodeToString(raw.BtcRefundAddress),
		RskRefundAddress:      raw.RskRefundAddress.String(),
		LpBtcAddress:          hex.EncodeToString(raw.LpBtcAddress),
		CallFee:               bigIntToWei(raw.CallFee),
		PenaltyFee:            bigIntToWei(raw.PenaltyFee),
		Nonce:                 raw.Nonce,
		DepositAddress:        hex.EncodeToString(raw.DepositAddress),
		Value:                 bigIntToWei(raw.Value),
		AgreementTimestamp:    raw.AgreementTimestamp,
		DepositDateLimit:      raw.DepositDateLimit,
		DepositConfirmations:  raw.DepositConfirmations,
		TransferConfirmations: raw.TransferConfirmations,
		TransferTime:          raw.TransferTime,
		ExpireDate:            raw.ExpireDate,
		ExpireBlock:           raw.ExpireBlock,
		GasFee:                bigIntToWei(raw.GasFee),
		ChainId:               bigIntToUint64(raw.ChainId),
	}
}
