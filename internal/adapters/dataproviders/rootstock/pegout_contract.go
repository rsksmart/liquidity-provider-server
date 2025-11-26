package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	log "github.com/sirupsen/logrus"
)

type pegoutContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      PegoutContractAdapter
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	abis          *FlyoverABIs
}

func NewPegoutContractImpl(
	client *RskClient,
	address string,
	contract PegoutContractAdapter,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	abis *FlyoverABIs,
) blockchain.PegoutContract {
	return &pegoutContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
		abis:          abis,
	}
}

func (pegoutContract *pegoutContractImpl) GetAddress() string {
	return pegoutContract.address
}

func (pegoutContract *pegoutContractImpl) IsPegOutQuoteCompleted(quoteHash string) (bool, error) {
	var quoteHashBytes [32]byte
	opts := &bind.CallOpts{}
	hashBytesSlice, err := hex.DecodeString(quoteHash)
	if err != nil {
		return false, err
	} else if len(hashBytesSlice) != 32 {
		return false, errors.New("quote hash must be 32 bytes long")
	}
	copy(quoteHashBytes[:], hashBytesSlice)
	result, err := rskRetry(pegoutContract.retryParams.Retries, pegoutContract.retryParams.Sleep,
		func() (bool, error) {
			return pegoutContract.contract.IsQuoteCompleted(opts, quoteHashBytes)
		})
	if err != nil {
		return false, err
	}
	return result, nil
}

func (pegoutContract *pegoutContractImpl) ValidatePegout(quoteHash string, btcTx []byte) error {
	var quoteHashBytes [32]byte
	// Set From address to the LP provider address so msg.sender validation passes
	opts := &bind.CallOpts{
		From: pegoutContract.signer.Address(),
	}
	hashBytesSlice, err := hex.DecodeString(quoteHash)
	if err != nil {
		return err
	} else if len(hashBytesSlice) != 32 {
		return errors.New("quote hash must be 32 bytes long")
	}
	copy(quoteHashBytes[:], hashBytesSlice)
	_, err = pegoutContract.contract.ValidatePegout(opts, quoteHashBytes, btcTx)
	return err
}

func (pegoutContract *pegoutContractImpl) DaoFeePercentage() (uint64, error) {
	opts := bind.CallOpts{}
	amount, err := rskRetry(pegoutContract.retryParams.Retries, pegoutContract.retryParams.Sleep,
		func() (*big.Int, error) {
			return pegoutContract.contract.GetFeePercentage(&opts)
		})
	if err != nil {
		return 0, err
	}
	return amount.Uint64(), nil
}

func (pegoutContract *pegoutContractImpl) HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	parsedQuote, err := parsePegoutQuote(pegoutQuote)
	if err != nil {
		return "", err
	}

	results, err = rskRetry(pegoutContract.retryParams.Retries, pegoutContract.retryParams.Sleep,
		func() ([32]byte, error) {
			return pegoutContract.contract.HashPegOutQuote(&opts, parsedQuote)
		})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (pegoutContract *pegoutContractImpl) RefundUserPegOut(quoteHash string) (string, error) {
	// Validate the hash format
	hashBytesSlice, err := hex.DecodeString(quoteHash)
	if err != nil {
		return "", fmt.Errorf("invalid quote hash format: %w", err)
	}
	if len(hashBytesSlice) != 32 {
		return "", errors.New("quote hash must be 32 bytes long")
	}

	opts := &bind.TransactOpts{
		From:   pegoutContract.signer.Address(),
		Signer: pegoutContract.signer.Sign,
	}
	receipt, err := awaitTx(pegoutContract.client, pegoutContract.miningTimeout, "RefundUserPegOut", func() (*geth.Transaction, error) {
		return pegoutContract.contract.RefundUserPegOut(opts, common.HexToHash(quoteHash))
	})

	if err != nil {
		return "", fmt.Errorf("refund user peg out error: %w", err)
	} else if receipt == nil {
		return "", errors.New("refund user peg out error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("refund user peg out error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (pegoutContract *pegoutContractImpl) RefundPegout(txConfig blockchain.TransactionConfig, params blockchain.RefundPegoutParams) (string, error) {
	var res []any
	var err error
	const (
		functionName                = "refundPegOut"
		notEnoughConfirmationsError = "NotEnoughConfirmations"
		unableToGetConfirmations    = "UnableToGetConfirmations"
	)

	log.Infof("Executing RefundPegOut with params: %s", params.String())
	revert := pegoutContract.contract.Caller().Call(
		&bind.CallOpts{From: pegoutContract.signer.Address()},
		&res, functionName,
		params.QuoteHash,
		params.BtcRawTx,
		params.BtcBlockHeaderHash,
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
	parsedRevert, err := ParseRevertReason(pegoutContract.abis.PegOut, revert)
	if err != nil && parsedRevert == nil {
		return "", fmt.Errorf("error parsing refundPegout result: %w", err)
	} else if parsedRevert != nil && (strings.EqualFold(notEnoughConfirmationsError, parsedRevert.Name) || strings.EqualFold(unableToGetConfirmations, parsedRevert.Name)) {
		log.Debugln("RefundPegout: bridge failed to validate BTC transaction. retrying on next confirmation.")
		return "", blockchain.WaitingForBridgeError
	} else if parsedRevert != nil {
		return "", fmt.Errorf("refundPegout reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:     pegoutContract.signer.Address(),
		Signer:   pegoutContract.signer.Sign,
		GasLimit: *txConfig.GasLimit,
	}

	receipt, err := awaitTx(pegoutContract.client, pegoutContract.miningTimeout, "RefundPegOut", func() (*geth.Transaction, error) {
		return pegoutContract.contract.RefundPegOut(opts, params.QuoteHash, params.BtcRawTx,
			params.BtcBlockHeaderHash, params.MerkleBranchPath, params.MerkleBranchHashes)
	})

	if err != nil {
		return "", fmt.Errorf("refund pegout error: %w", err)
	} else if receipt == nil {
		return "", errors.New("refund pegout error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("refund pegout error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (pegoutContract *pegoutContractImpl) GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error) {
	var lbcEvent *bindings.IPegOutPegOutDeposit
	result := make([]quote.PegoutDeposit, 0)

	rawIterator, err := pegoutContract.contract.FilterPegOutDeposit(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	}, nil, nil, nil)
	// The adapter is to be able to mock the iterator in tests
	iterator := pegoutContract.contract.DepositEventIteratorAdapter(rawIterator)
	defer func() {
		if rawIterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing PegOutDeposit event iterator: ", err)
			}
		}
	}()
	if err != nil || rawIterator == nil {
		return nil, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Event()
		result = append(result, quote.PegoutDeposit{
			TxHash:      lbcEvent.Raw.TxHash.String(),
			QuoteHash:   hex.EncodeToString(lbcEvent.QuoteHash[:]),
			Amount:      entities.NewBigWei(lbcEvent.Amount),
			Timestamp:   time.Unix(lbcEvent.Timestamp.Int64(), 0),
			BlockNumber: lbcEvent.Raw.BlockNumber,
			From:        lbcEvent.Sender.String(),
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}

	return result, nil
}

// parsePegoutQuote parses a quote.PegoutQuote into a bindings.QuotesPegOutQuote. All BTC address fields support all address types.
func parsePegoutQuote(pegoutQuote quote.PegoutQuote) (bindings.QuotesPegOutQuote, error) {
	var parsedQuote bindings.QuotesPegOutQuote
	var err error

	if err = entities.ValidateStruct(pegoutQuote); err != nil {
		return bindings.QuotesPegOutQuote{}, err
	}

	if err = ParseAddress(&parsedQuote.LbcAddress, pegoutQuote.LbcAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing lbc address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.LpRskAddress, pegoutQuote.LpRskAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing liquidity provider rsk address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.RskRefundAddress, pegoutQuote.RskRefundAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing user rsk refund address: %w", err)
	}

	if parsedQuote.BtcRefundAddress, err = bitcoin.DecodeAddress(pegoutQuote.BtcRefundAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing user btc refund address: %w", err)
	}
	if parsedQuote.LpBtcAddress, err = bitcoin.DecodeAddress(pegoutQuote.LpBtcAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing liquidity provider btc address: %w", err)
	}
	if parsedQuote.DepositAddress, err = bitcoin.DecodeAddress(pegoutQuote.DepositAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing pegout deposit address: %w", err)
	}

	parsedQuote.CallFee = pegoutQuote.CallFee.AsBigInt()
	parsedQuote.PenaltyFee = pegoutQuote.PenaltyFee.AsBigInt()
	parsedQuote.Nonce = pegoutQuote.Nonce
	parsedQuote.Value = pegoutQuote.Value.AsBigInt()
	parsedQuote.AgreementTimestamp = pegoutQuote.AgreementTimestamp
	parsedQuote.DepositDateLimit = pegoutQuote.DepositDateLimit
	parsedQuote.DepositConfirmations = pegoutQuote.DepositConfirmations
	parsedQuote.TransferConfirmations = pegoutQuote.TransferConfirmations
	parsedQuote.TransferTime = pegoutQuote.TransferTime
	parsedQuote.ExpireDate = pegoutQuote.ExpireDate
	parsedQuote.ExpireBlock = pegoutQuote.ExpireBlock
	parsedQuote.ProductFeeAmount = pegoutQuote.ProductFeeAmount.AsBigInt()
	parsedQuote.GasFee = pegoutQuote.GasFee.AsBigInt()
	return parsedQuote, nil
}
