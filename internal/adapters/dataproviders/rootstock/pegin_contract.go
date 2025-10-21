package rootstock

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
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
	contract      PeginContractAdapter
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	abis          *FlyoverABIs
}

func NewPeginContractImpl(
	client *RskClient,
	address string,
	contract PeginContractAdapter,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	abis *FlyoverABIs,
) blockchain.PeginContract {
	return &peginContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
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
			return peginContract.contract.GetBalance(&bind.CallOpts{}, parsedAddress)
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
			return peginContract.contract.GetFeePercentage(&opts)
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
			return peginContract.contract.HashPegInQuote(&bind.CallOpts{}, parsedQuote)
		})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (peginContract *peginContractImpl) CallForUser(txConfig blockchain.TransactionConfig, peginQuote quote.PeginQuote) (string, error) {
	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return "", err
	}

	opts := &bind.TransactOpts{
		GasLimit: *txConfig.GasLimit,
		Value:    txConfig.Value.AsBigInt(),
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
	}

	receipt, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(peginContract.client, peginContract.miningTimeout, "CallForUser", func() (*geth.Transaction, error) {
				return peginContract.contract.CallForUser(opts, parsedQuote)
			})
		})

	if err != nil {
		return "", fmt.Errorf("call for user error: %w", err)
	} else if receipt == nil {
		return "", errors.New("call for user error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("call for user error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (peginContract *peginContractImpl) RegisterPegin(params blockchain.RegisterPeginParams) (string, error) {
	const (
		functionName          = "registerPegIn"
		waitingForBridgeError = "NotEnoughConfirmations"
	)
	var res []any
	var err error
	var parsedQuote bindings.QuotesPegInQuote
	if parsedQuote, err = parsePeginQuote(params.Quote); err != nil {
		return "", err
	}
	log.Infof("Executing RegisterPegIn with params: %s\n", params.String())
	revert := peginContract.contract.Caller().Call(
		&bind.CallOpts{}, &res, functionName,
		parsedQuote,
		params.QuoteSignature,
		params.BitcoinRawTransaction,
		params.PartialMerkleTree,
		params.BlockHeight,
	)

	parsedRevert, err := ParseRevertReason(peginContract.abis.PegIn, revert)
	if err != nil && parsedRevert == nil {
		return "", fmt.Errorf("error parsing registerPegIn result: %w", err)
	} else if parsedRevert != nil && strings.EqualFold(waitingForBridgeError, parsedRevert.Name) {
		log.Debugln("RegisterPegin: bridge failed to validate BTC transaction. retrying on next confirmation.")
		// allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		return "", blockchain.WaitingForBridgeError
	} else if parsedRevert != nil {
		return "", fmt.Errorf("registerPegIn reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:     peginContract.signer.Address(),
		Signer:   peginContract.signer.Sign,
		GasLimit: registerPeginGasLimit,
	}

	receipt, err := awaitTx(peginContract.client, peginContract.miningTimeout, "RegisterPegIn", func() (*geth.Transaction, error) {
		return peginContract.contract.RegisterPegIn(opts, parsedQuote, params.QuoteSignature,
			params.BitcoinRawTransaction, params.PartialMerkleTree, params.BlockHeight)
	})

	if err != nil {
		return "", fmt.Errorf("register pegin error: %w", err)
	} else if receipt == nil {
		return "", errors.New("register pegin error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("register pegin error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
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
