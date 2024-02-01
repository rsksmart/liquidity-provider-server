package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

const registerPeginGasLimit = 2500000

type liquidityBridgeContractImpl struct {
	client   *ethclient.Client
	address  string
	contract *bindings.LiquidityBridgeContract
	signer   TransactionSigner
}

func NewLiquidityBridgeContractImpl(
	client *RskClient,
	address string,
	contract *bindings.LiquidityBridgeContract,
	signer TransactionSigner,
) blockchain.LiquidityBridgeContract {
	return &liquidityBridgeContractImpl{
		client:   client.client,
		address:  address,
		contract: contract,
		signer:   signer,
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

	results, err = rskRetry(func() ([32]byte, error) {
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

	results, err = rskRetry(func() ([32]byte, error) {
		return lbc.contract.HashPegoutQuote(&opts, parsedQuote)
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(results[:]), nil
}

func (lbc *liquidityBridgeContractImpl) GetProviders() ([]entities.RegisteredLiquidityProvider, error) {
	var i, maxProviderId int64
	var providerType entities.ProviderType
	var providers []bindings.LiquidityBridgeContractLiquidityProvider
	var provider bindings.LiquidityBridgeContractLiquidityProvider

	opts := &bind.CallOpts{}
	maxId, err := rskRetry(func() (*big.Int, error) {
		return lbc.contract.GetProviderIds(opts)
	})
	if err != nil {
		return nil, err
	}

	maxProviderId = maxId.Int64()
	providerIds := make([]*big.Int, 0)

	for i = 1; i <= maxProviderId; i++ {
		providerIds = append(providerIds, big.NewInt(i))
	}

	providers, err = rskRetry(func() ([]bindings.LiquidityBridgeContractLiquidityProvider, error) {
		return lbc.contract.GetProviders(opts, providerIds)
	})
	if err != nil {
		return nil, err
	}
	parsedProviders := make([]entities.RegisteredLiquidityProvider, 0)
	for i = 0; i < maxProviderId+1; i++ {
		provider = providers[i]
		providerType = entities.ProviderType(provider.ProviderType)
		if !providerType.IsValid() {
			return nil, entities.InvalidProviderTypeError
		}
		parsedProviders = append(parsedProviders, entities.RegisteredLiquidityProvider{
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

func (lbc *liquidityBridgeContractImpl) ProviderResign() error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}
	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "Resign", func() (*geth.Transaction, error) {
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
	var parsedId *big.Int
	parsedId.SetUint64(id)
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "SetProviderStatus", func() (*geth.Transaction, error) {
			return lbc.contract.SetProviderStatus(opts, parsedId, newStatus)
		})
	})

	if err != nil {
		return err
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("resign transaction failed")
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
	collateral, err := rskRetry(func() (*big.Int, error) {
		return lbc.contract.GetCollateral(opts, parsedAddress)
	})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(collateral), nil
}

func (lbc *liquidityBridgeContractImpl) GetPegoutCollateral(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	collateral, err := rskRetry(func() (*big.Int, error) {
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
	collateral, err := rskRetry(func() (*big.Int, error) {
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

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "AddCollateral", func() (*geth.Transaction, error) {
			return lbc.contract.AddCollateral(opts)
		})
	})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return fmt.Errorf("error adding pegin collateral")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) AddPegoutCollateral(amount *entities.Wei) error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
		Value:  amount.AsBigInt(),
	}

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "AddPegoutCollateral", func() (*geth.Transaction, error) {
			return lbc.contract.AddPegoutCollateral(opts)
		})
	})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return fmt.Errorf("error adding pegout collateral")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) WithdrawCollateral() error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "WithdrawCollateral", func() (*geth.Transaction, error) {
			return lbc.contract.WithdrawCollateral(opts)
		})
	})

	if err != nil {
		return fmt.Errorf("withdraw pegin collateral error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return fmt.Errorf("withdraw pegin collateral error")
	}
	return nil
}

func (lbc *liquidityBridgeContractImpl) WithdrawPegoutCollateral() error {
	opts := &bind.TransactOpts{
		From:   lbc.signer.Address(),
		Signer: lbc.signer.Sign,
	}

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "WithdrawPegoutCollateral", func() (*geth.Transaction, error) {
			return lbc.contract.WithdrawPegoutCollateral(opts)
		})
	})

	if err != nil {
		return fmt.Errorf("withdraw pegout collateral error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return fmt.Errorf("withdraw pegout collateral error")
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
	balance, err := rskRetry(func() (*big.Int, error) {
		return lbc.contract.GetBalance(opts, parsedAddress)
	})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}

func (lbc *liquidityBridgeContractImpl) CallForUser(txConfig blockchain.TransactionConfig, peginQuote quote.PeginQuote) (string, error) {
	opts := &bind.TransactOpts{
		GasLimit: *txConfig.GasLimit,
		Value:    txConfig.Value.AsBigInt(),
		From:     lbc.signer.Address(),
		Signer:   lbc.signer.Sign,
	}

	parsedQuote, err := parsePeginQuote(peginQuote)
	if err != nil {
		return "", err
	}

	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "CallForUser", func() (*geth.Transaction, error) {
			return lbc.contract.CallForUser(opts, parsedQuote)
		})
	})

	if err != nil {
		return "", fmt.Errorf("call for user error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return "", errors.New("callfor user error: incomplete receipt")
	}
	return receipt.TxHash.String(), nil
}

func (lbc *liquidityBridgeContractImpl) RegisterPegin(params blockchain.RegisterPeginParams) (string, error) {
	var res []any
	var err error
	var parsedQuote bindings.QuotesPeginQuote
	lbcCaller := &bindings.LiquidityBridgeContractCallerRaw{Contract: &lbc.contract.LiquidityBridgeContractCaller}
	if parsedQuote, err = parsePeginQuote(params.Quote); err != nil {
		return "", err
	}
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

	receipt, err := awaitTx(lbc.client, "RegisterPegIn", func() (*geth.Transaction, error) {
		return lbc.contract.RegisterPegIn(opts, parsedQuote, params.QuoteSignature,
			params.BitcoinRawTransaction, params.PartialMerkleTree, params.BlockHeight)
	})

	if err != nil {
		return "", fmt.Errorf("register pegin error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return "", errors.New("register pegin error: incomplete receipt")
	}
	return receipt.TxHash.String(), nil
}

func (lbc *liquidityBridgeContractImpl) RefundPegout(txConfig blockchain.TransactionConfig, params blockchain.RefundPegoutParams) (string, error) {
	opts := &bind.TransactOpts{
		From:     lbc.signer.Address(),
		Signer:   lbc.signer.Sign,
		GasLimit: *txConfig.GasLimit,
	}

	log.Infof("Executing RefundPegOut with params: %s\n", params.String())
	receipt, err := awaitTx(lbc.client, "RefundPegOut", func() (*geth.Transaction, error) {
		return lbc.contract.RefundPegOut(opts, params.QuoteHash, params.BtcRawTx,
			params.BtcBlockHeaderHash, params.MerkleBranchPath, params.MerkleBranchHashes)
	})

	if err != nil && strings.Contains(err.Error(), "LBC049") {
		log.Debugln("RefundPegout: bridge failed to validate BTC transaction. retrying on next confirmation.")
		return "", blockchain.WaitingForBridgeError
	} else if err != nil {
		return "", fmt.Errorf("refund pegout error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return "", errors.New("refund pegout error: incomplete receipt")
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

	return rskRetry(func() (bool, error) {
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

	return rskRetry(func() (bool, error) {
		return lbc.contract.IsOperationalForPegout(opts, parsedAddress)
	})
}

func (lbc *liquidityBridgeContractImpl) GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error) {
	var lbcEvent *bindings.LiquidityBridgeContractPegOutDeposit
	result := make([]quote.PegoutDeposit, 0)

	iterator, err := lbc.contract.FilterPegOutDeposit(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	}, nil, nil)
	defer func() {
		if iterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing PegOutDeposit event iterator: ", err)
			}
		}
	}()
	if err != nil || iterator == nil {
		return result, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Event
		result = append(result, quote.PegoutDeposit{
			TxHash:      lbcEvent.Raw.TxHash.String(),
			QuoteHash:   hex.EncodeToString(lbcEvent.QuoteHash[:]),
			Amount:      entities.NewBigWei(lbcEvent.Amount),
			Timestamp:   time.Unix(lbcEvent.Timestamp.Int64(), 0),
			BlockNumber: lbcEvent.Raw.BlockNumber,
			From:        lbcEvent.Sender.String(),
		})
	}
	if iterator.Error() != nil {
		return nil, err
	}

	return result, nil
}

func (lbc *liquidityBridgeContractImpl) GetPeginPunishmentEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]entities.PunishmentEvent, error) {
	var lbcEvent *bindings.LiquidityBridgeContractPenalized
	result := make([]entities.PunishmentEvent, 0)

	iterator, err := lbc.contract.FilterPenalized(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	})
	defer func() {
		if iterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing Penalized event iterator: ", err)
			}
		}
	}()
	if err != nil || iterator == nil {
		return result, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Event
		result = append(result, entities.PunishmentEvent{
			LiquidityProvider: lbcEvent.LiquidityProvider.String(),
			Penalty:           entities.NewBigWei(lbcEvent.Penalty),
			QuoteHash:         hex.EncodeToString(lbcEvent.QuoteHash[:]),
		})
	}
	if iterator.Error() != nil {
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
	receipt, err := rskRetry(func() (*geth.Receipt, error) {
		return awaitTx(lbc.client, "Register", func() (*geth.Transaction, error) {
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
		return 0, fmt.Errorf("error registering provider: %w", err)
	}
	return registerEvent.Id.Int64(), nil
}

// TODO currently we only support P2PKH addresses (P2SH is allowed for federation address)
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
	if parsedQuote.BtcRefundAddress, err = bitcoin.DecodeAddressBase58OnlyLegacy(peginQuote.BtcRefundAddress, true); err != nil {
		return bindings.QuotesPeginQuote{}, fmt.Errorf("error parsing user btc refund address: %w", err)
	}
	if parsedQuote.LiquidityProviderBtcAddress, err = bitcoin.DecodeAddressBase58OnlyLegacy(peginQuote.LpBtcAddress, true); err != nil {
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
	return parsedQuote, nil
}

// TODO currently we only support P2PKH addresses
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

	if parsedQuote.BtcRefundAddress, err = bitcoin.DecodeAddressBase58OnlyLegacy(pegoutQuote.BtcRefundAddress, true); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing user btc refund address: %w", err)
	}
	if parsedQuote.LpBtcAddress, err = bitcoin.DecodeAddressBase58OnlyLegacy(pegoutQuote.LpBtcAddress, true); err != nil {
		return bindings.QuotesPegOutQuote{}, fmt.Errorf("error parsing liquidity provider btc address: %w", err)
	}
	if parsedQuote.DeposityAddress, err = bitcoin.DecodeAddressBase58OnlyLegacy(pegoutQuote.DepositAddress, true); err != nil {
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
