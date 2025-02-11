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
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	log "github.com/sirupsen/logrus"
)

// registerPeginGasLimit Fixed gas limit for registerPegin function, should change only if the function does
const registerPeginGasLimit = 2500000

type liquidityBridgeContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      LbcAdapter
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
}

func NewLiquidityBridgeContractImpl(
	client *RskClient,
	address string,
	contract LbcAdapter,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
) blockchain.LiquidityBridgeContract {
	return &liquidityBridgeContractImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		signer:      signer,
		retryParams: retryParams,
	}
}

func (lbc *liquidityBridgeContractImpl) GetAddress() string {
	return lbc.address
}

func (lbc *liquidityBridgeContractImpl) HashPeginQuote(peginQuote quote.PeginQuote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return "", err
	}

	results, err = rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() ([32]byte, error) {
			return lbc.contract.HashQuote(&opts, parsedQuote)
		})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (lbc *liquidityBridgeContractImpl) HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	parsedQuote, err := parsePegoutQuote(pegoutQuote)
	if err != nil {
		return "", err
	}

	results, err = rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() ([32]byte, error) {
			return lbc.contract.HashPegoutQuote(&opts, parsedQuote)
		})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (lbc *liquidityBridgeContractImpl) GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error) {
	var providerType liquidity_provider.ProviderType
	var providers []bindings.LiquidityBridgeContractLiquidityProvider

	opts := &bind.CallOpts{}
	providers, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() ([]bindings.LiquidityBridgeContractLiquidityProvider, error) {
			return lbc.contract.GetProviders(opts)
		})
	if err != nil {
		return nil, err
	}
	parsedProviders := make([]liquidity_provider.RegisteredLiquidityProvider, 0)
	for _, provider := range providers {
		providerType = liquidity_provider.ProviderType(provider.ProviderType)
		if !providerType.IsValid() {
			return nil, liquidity_provider.InvalidProviderTypeError
		}
		parsedProviders = append(parsedProviders, liquidity_provider.RegisteredLiquidityProvider{
			Id:           provider.Id.Uint64(),
			Address:      provider.Provider.String(),
			Name:         provider.Name,
			ApiBaseUrl:   provider.ApiBaseUrl,
			Status:       provider.Status,
			ProviderType: providerType,
		})
	}
	return parsedProviders, nil
}

func (lbc *liquidityBridgeContractImpl) GetProvider(address string) (liquidity_provider.RegisteredLiquidityProvider, error) {
	var providerType liquidity_provider.ProviderType
	const lbcProviderNotRegisteredError = "LBC001"

	if !common.IsHexAddress(address) {
		return liquidity_provider.RegisteredLiquidityProvider{}, blockchain.InvalidAddressError
	}

	opts := &bind.CallOpts{}
	provider, err := lbc.contract.GetProvider(opts, common.HexToAddress(address))
	if err != nil && err.Error() == lbcProviderNotRegisteredError {
		return liquidity_provider.RegisteredLiquidityProvider{}, liquidity_provider.ProviderNotFoundError
	} else if err != nil {
		return liquidity_provider.RegisteredLiquidityProvider{}, err
	}

	providerType = liquidity_provider.ProviderType(provider.ProviderType)
	if !providerType.IsValid() {
		return liquidity_provider.RegisteredLiquidityProvider{}, liquidity_provider.InvalidProviderTypeError
	}
	return liquidity_provider.RegisteredLiquidityProvider{
		Id:           provider.Id.Uint64(),
		Address:      provider.Provider.String(),
		Name:         provider.Name,
		ApiBaseUrl:   provider.ApiBaseUrl,
		Status:       provider.Status,
		ProviderType: providerType,
	}, nil
}

func (lbc *liquidityBridgeContractImpl) ProviderResign() error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}
	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "Resign", func() (*geth.Transaction, error) {
				return lbc.contract.Resign(opts)
			})
		})

	if err != nil {
		return err
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("resign transaction failed")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) SetProviderStatus(id uint64, newStatus bool) error {
	parsedId := new(big.Int)
	parsedId.SetUint64(id)
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}

	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "SetProviderStatus", func() (*geth.Transaction, error) {
				return lbc.contract.SetProviderStatus(opts, parsedId, newStatus)
			})
		})

	if err != nil {
		return err
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("setProviderStatus transaction failed")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) GetCollateral(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	collateral, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*big.Int, error) {
			return lbc.contract.GetCollateral(opts, parsedAddress)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(collateral), nil
}

func (lbc *liquidityBridgeContractImpl) IsPegOutQuoteCompleted(quoteHash string) (bool, error) {
	var quoteHashBytes [32]byte
	opts := &bind.CallOpts{}
	hashBytesSlice, err := hex.DecodeString(quoteHash)
	if err != nil {
		return false, err
	} else if len(hashBytesSlice) != 32 {
		return false, errors.New("quote hash must be 32 bytes long")
	}
	copy(quoteHashBytes[:], hashBytesSlice)
	result, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (bool, error) {
			return lbc.contract.IsPegOutQuoteCompleted(opts, quoteHashBytes)
		})
	if err != nil {
		return false, err
	}
	return result, nil
}

func (lbc *liquidityBridgeContractImpl) GetPegoutCollateral(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	collateral, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*big.Int, error) {
			return lbc.contract.GetPegoutCollateral(opts, parsedAddress)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(collateral), nil
}

func (lbc *liquidityBridgeContractImpl) GetMinimumCollateral() (*entities.Wei, error) {
	var err error
	opts := &bind.CallOpts{}
	collateral, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*big.Int, error) {
			return lbc.contract.GetMinCollateral(opts)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(collateral), nil
}

func (lbc *liquidityBridgeContractImpl) AddCollateral(amount *entities.Wei) error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
		Value:  amount.AsBigInt(),
	}

	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "AddCollateral", func() (*geth.Transaction, error) {
				return lbc.contract.AddCollateral(opts)
			})
		})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("error adding pegin collateral")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) AddPegoutCollateral(amount *entities.Wei) error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
		Value:  amount.AsBigInt(),
	}

	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "AddPegoutCollateral", func() (*geth.Transaction, error) {
				return lbc.contract.AddPegoutCollateral(opts)
			})
		})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("error adding pegout collateral")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) WithdrawCollateral() error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}

	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "WithdrawCollateral", func() (*geth.Transaction, error) {
				return lbc.contract.WithdrawCollateral(opts)
			})
		})

	if err != nil {
		return fmt.Errorf("withdraw pegin collateral error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("withdraw pegin collateral error")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) GetBalance(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	balance, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*big.Int, error) {
			return lbc.contract.GetBalance(opts, parsedAddress)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}

func (lbc *liquidityBridgeContractImpl) CallForUser(txConfig blockchain.TransactionConfig, peginQuote quote.PeginQuote) (string, error) {
	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return "", err
	}

	opts := &bind.TransactOpts{
		GasLimit: *txConfig.GasLimit,
		Value:    txConfig.Value.AsBigInt(),
		From:     lbc.signer.Address(),
		Signer:   lbc.signer.Sign,
	}

	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "CallForUser", func() (*geth.Transaction, error) {
				return lbc.contract.CallForUser(opts, parsedQuote)
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

func (lbc *liquidityBridgeContractImpl) RegisterPegin(params blockchain.RegisterPeginParams) (string, error) {
	var res []any
	var err error
	var parsedQuote bindings.QuotesPeginQuote
	if parsedQuote, err = parsePeginQuote(params.Quote); err != nil {
		return "", err
	}
	lbcCaller := lbc.contract.Caller()
	log.Infof("Executing RegisterPegIn with params: %s\n", params.String())
	err = lbcCaller.Call(
		&bind.CallOpts{}, &res, "registerPegIn",
		parsedQuote,
		params.QuoteSignature,
		params.BitcoinRawTransaction,
		params.PartialMerkleTree,
		params.BlockHeight,
	)
	if err != nil && strings.Contains(err.Error(), "LBC031") {
		log.Debugln("RegisterPegin: bridge failed to validate BTC transaction. retrying on next confirmation.")
		// allow retrying in case the bridge didn't acknowledge all required confirmations have occurred
		return "", blockchain.WaitingForBridgeError
	} else if err != nil {
		return "", err
	}

	opts := &bind.TransactOpts{
		From:     lbc.signer.Address(),
		Signer:   lbc.signer.Sign,
		GasLimit: registerPeginGasLimit,
	}

	receipt, err := awaitTx(lbc.client, lbc.miningTimeout, "RegisterPegIn", func() (*geth.Transaction, error) {
		return lbc.contract.RegisterPegIn(opts, parsedQuote, params.QuoteSignature,
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

func (lbc *liquidityBridgeContractImpl) RefundPegout(txConfig blockchain.TransactionConfig, params blockchain.RefundPegoutParams) (string, error) {
	var res []any
	var err error
	lbcCaller := lbc.contract.Caller()
	log.Infof("Executing RefundPegOut with params: %s", params.String())
	err = lbcCaller.Call(
		&bind.CallOpts{From: lbc.signer.Address()},
		&res, "refundPegOut",
		params.QuoteHash,
		params.BtcRawTx,
		params.BtcBlockHeaderHash,
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
	if err != nil && strings.Contains(err.Error(), "LBC049") {
		log.Debugln("RefundPegout: bridge failed to validate BTC transaction. retrying on next confirmation.")
		return "", blockchain.WaitingForBridgeError
	} else if err != nil {
		return "", err
	}

	opts := &bind.TransactOpts{
		From:     lbc.signer.Address(),
		Signer:   lbc.signer.Sign,
		GasLimit: *txConfig.GasLimit,
	}

	receipt, err := awaitTx(lbc.client, lbc.miningTimeout, "RefundPegOut", func() (*geth.Transaction, error) {
		return lbc.contract.RefundPegOut(opts, params.QuoteHash, params.BtcRawTx,
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

func (lbc *liquidityBridgeContractImpl) IsOperationalPegin(address string) (bool, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}

	if err = ParseAddress(&parsedAddress, address); err != nil {
		return false, err
	}

	return rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (bool, error) {
			return lbc.contract.IsOperational(opts, parsedAddress)
		})
}

func (lbc *liquidityBridgeContractImpl) IsOperationalPegout(address string) (bool, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}

	if err = ParseAddress(&parsedAddress, address); err != nil {
		return false, err
	}

	return rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (bool, error) {
			return lbc.contract.IsOperationalForPegout(opts, parsedAddress)
		})
}

func (lbc *liquidityBridgeContractImpl) GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error) {
	var lbcEvent *bindings.LiquidityBridgeContractPegOutDeposit
	result := make([]quote.PegoutDeposit, 0)

	rawIterator, err := lbc.contract.FilterPegOutDeposit(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	}, nil, nil)
	// The adapter is to be able to mock the iterator in tests
	iterator := lbc.contract.DepositEventIteratorAdapter(rawIterator)
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

func (lbc *liquidityBridgeContractImpl) GetPeginPunishmentEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]liquidity_provider.PunishmentEvent, error) {
	var lbcEvent *bindings.LiquidityBridgeContractPenalized
	result := make([]liquidity_provider.PunishmentEvent, 0)

	rawIterator, err := lbc.contract.FilterPenalized(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	})
	iterator := lbc.contract.PenalizedEventIteratorAdapter(rawIterator)
	defer func() {
		if rawIterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing Penalized event iterator: ", err)
			}
		}
	}()
	if err != nil || rawIterator == nil {
		return nil, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Event()
		result = append(result, liquidity_provider.PunishmentEvent{
			LiquidityProvider: lbcEvent.LiquidityProvider.String(),
			Penalty:           entities.NewBigWei(lbcEvent.Penalty),
			QuoteHash:         hex.EncodeToString(lbcEvent.QuoteHash[:]),
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

func (lbc *liquidityBridgeContractImpl) RegisterProvider(txConfig blockchain.TransactionConfig, params blockchain.ProviderRegistrationParams) (int64, error) {
	var err error

	opts := &bind.TransactOpts{
		Value:  txConfig.Value.AsBigInt(),
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}
	receipt, err := rskRetry(lbc.retryParams.Retries, lbc.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(lbc.client, lbc.miningTimeout, "Register", func() (*geth.Transaction, error) {
				return lbc.contract.Register(opts, params.Name, params.ApiBaseUrl, params.Status, string(params.Type))
			})
		})

	if err != nil {
		return 0, fmt.Errorf("error registering provider: %w", err)
	} else if receipt == nil || receipt.Status == 0 || len(receipt.Logs) == 0 {
		return 0, errors.New("error registering provider: incomplete receipt")
	}

	registerEvent, err := lbc.contract.ParseRegister(*receipt.Logs[0])
	if err != nil {
		return 0, fmt.Errorf("error parsing register event: %w", err)
	}
	return registerEvent.Id.Int64(), nil
}

func (lbc *liquidityBridgeContractImpl) UpdateProvider(name, url string) (string, error) {
	opts := &bind.TransactOpts{From: lbc.signer.Address(), Signer: lbc.signer.Sign}
	receipt, err := awaitTx(lbc.client, lbc.miningTimeout, "UpdateProvider", func() (*geth.Transaction, error) {
		return lbc.contract.UpdateProvider(opts, name, url)
	})

	if err != nil {
		return "", fmt.Errorf("update provider error: %w", err)
	} else if receipt == nil {
		return "", errors.New("update provider error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("update provider error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (lbc *liquidityBridgeContractImpl) RefundUserPegOut(quoteHash string) (string, error) {
	// Validate the hash format
	hashBytesSlice, err := hex.DecodeString(quoteHash)
	if err != nil {
		return "", fmt.Errorf("invalid quote hash format: %w", err)
	}
	if len(hashBytesSlice) != 32 {
		return "", errors.New("quote hash must be 32 bytes long")
	}

	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}
	receipt, err := awaitTx(lbc.client, lbc.miningTimeout, "RefundUserPegOut", func() (*geth.Transaction, error) {
		return lbc.contract.RefundUserPegOut(opts, common.HexToHash(quoteHash))
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

// parsePeginQuote parses a quote.PeginQuote into a bindings.QuotesPeginQuote. All BTC address fields support all address types
// except for FedBtcAddress which must be a P2SH address.
func parsePeginQuote(peginQuote quote.PeginQuote) (bindings.QuotesPeginQuote, error) {
	var decodedFederationAddress []byte
	var parsedQuote bindings.QuotesPeginQuote
	var err error

	if err = entities.ValidateStruct(peginQuote); err != nil {
		return bindings.QuotesPeginQuote{}, err
	}

	if decodedFederationAddress, err = bitcoin.DecodeAddressBase58(peginQuote.FedBtcAddress, false); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing federation address: %w", err)
	} else {
		copy(parsedQuote.FedBtcAddress[:], decodedFederationAddress)
	}
	if parsedQuote.BtcRefundAddress, err = bitcoin.DecodeAddress(peginQuote.BtcRefundAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing user btc refund address: %w", err)
	}
	if parsedQuote.LiquidityProviderBtcAddress, err = bitcoin.DecodeAddress(peginQuote.LpBtcAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing btc liquidity provider address: %w", err)
	}

	if err = ParseAddress(&parsedQuote.LbcAddress, peginQuote.LbcAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing lbc address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.LiquidityProviderRskAddress, peginQuote.LpRskAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing liquidity provider rsk address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.RskRefundAddress, peginQuote.RskRefundAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing user rsk refund address: %w", err)
	}
	if err = ParseAddress(&parsedQuote.ContractAddress, peginQuote.ContractAddress); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing destination contract address: %w", err)
	}

	if parsedQuote.Data, err = blockchain.DecodeStringTrimPrefix(peginQuote.Data); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing data: %w", err)
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
	parsedQuote.ProductFeeAmount = new(big.Int)
	parsedQuote.ProductFeeAmount.SetUint64(peginQuote.ProductFeeAmount)
	parsedQuote.GasFee = peginQuote.GasFee.AsBigInt()
	parsedQuote.CallOnRegister = peginQuote.CallOnRegister
	return parsedQuote, nil
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
	if parsedQuote.DeposityAddress, err = bitcoin.DecodeAddress(pegoutQuote.DepositAddress); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing pegout deposit address: %w", err)
	}

	parsedQuote.CallFee = pegoutQuote.CallFee.AsBigInt()
	parsedQuote.PenaltyFee = new(big.Int)
	parsedQuote.PenaltyFee.SetUint64(pegoutQuote.PenaltyFee)
	parsedQuote.Nonce = pegoutQuote.Nonce
	parsedQuote.Value = pegoutQuote.Value.AsBigInt()
	parsedQuote.AgreementTimestamp = pegoutQuote.AgreementTimestamp
	parsedQuote.DepositDateLimit = pegoutQuote.DepositDateLimit
	parsedQuote.DepositConfirmations = pegoutQuote.DepositConfirmations
	parsedQuote.TransferConfirmations = pegoutQuote.TransferConfirmations
	parsedQuote.TransferTime = pegoutQuote.TransferTime
	parsedQuote.ExpireDate = pegoutQuote.ExpireDate
	parsedQuote.ExpireBlock = pegoutQuote.ExpireBlock
	parsedQuote.ProductFeeAmount = new(big.Int)
	parsedQuote.ProductFeeAmount.SetUint64(pegoutQuote.ProductFeeAmount)
	parsedQuote.GasFee = pegoutQuote.GasFee.AsBigInt()
	return parsedQuote, nil
}
