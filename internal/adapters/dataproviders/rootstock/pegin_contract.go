package rootstock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
)

// registerPeginGasLimit Fixed gas limit for registerPegin function, should change only if the function does
const registerPeginGasLimit = 2500000

type peginContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	binding       *bindings.PeginContract
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

func (peginContract *peginContractImpl) DaoFeePercentage() (uint64, error) {
	opts := bind.CallOpts{}
	amount, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := peginContract.binding.TryPackGetFeePercentage()
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(peginContract.contract, &opts, callData, peginContract.binding.UnpackGetFeePercentage)
		})
	if err != nil {
		return 0, err
	}
	return amount.Uint64(), nil
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

	parsedQuote.CallFee = peginQuote.CallFee.AsBigInt()
	parsedQuote.PenaltyFee = peginQuote.PenaltyFee.AsBigInt()
	parsedQuote.GasLimit = peginQuote.GasLimit
	parsedQuote.Nonce = peginQuote.Nonce
	parsedQuote.Value = peginQuote.Value.AsBigInt()
	parsedQuote.AgreementTimestamp = peginQuote.AgreementTimestamp
	parsedQuote.CallTime = peginQuote.LpCallTime
	parsedQuote.DepositConfirmations = peginQuote.Confirmations
	parsedQuote.TimeForDeposit = peginQuote.TimeForDeposit
	parsedQuote.ProductFeeAmount = peginQuote.ProductFeeAmount.AsBigInt()
	parsedQuote.GasFee = peginQuote.GasFee.AsBigInt()
	parsedQuote.CallOnRegister = peginQuote.CallOnRegister
	return parsedQuote, nil
}
