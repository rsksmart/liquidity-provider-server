package pegout

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

const (
	// BridgeConversionGasLimit see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	BridgeConversionGasLimit = 100000
	// BridgeConversionGasPrice see https://dev.rootstock.io/rsk/rbtc/conversion/networks/
	BridgeConversionGasPrice = 60000000
)

type RebalanceStrategy string

const (
	AllAtOnce RebalanceStrategy = "ALL_AT_ONCE"
	UtxoSplit RebalanceStrategy = "UTXO_SPLIT"
)

func ParseRebalanceStrategy(s string) (RebalanceStrategy, error) {
	switch s {
	case string(AllAtOnce):
		return AllAtOnce, nil
	case string(UtxoSplit):
		return UtxoSplit, nil
	default:
		return "", fmt.Errorf("unknown rebalance strategy: %q", s)
	}
}

type BridgePegoutUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
	pegoutProvider  liquidity_provider.PegoutLiquidityProvider
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
	rskWalletMutex  sync.Locker
	strategy        RebalanceStrategy
}

func NewBridgePegoutUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
	rskWalletMutex sync.Locker,
	strategy RebalanceStrategy,
) *BridgePegoutUseCase {
	return &BridgePegoutUseCase{
		quoteRepository: quoteRepository,
		pegoutProvider:  pegoutProvider,
		rskWallet:       rskWallet,
		contracts:       contracts,
		rskWalletMutex:  rskWalletMutex,
		strategy:        strategy,
	}
}

func (useCase *BridgePegoutUseCase) Run(ctx context.Context, watchedQuotes ...quote.WatchedPegoutQuote) error {
	var err error
	var totalValue *entities.Wei

	totalValue, err = useCase.calculateTotalToPegout(watchedQuotes)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}

	pegoutConfig := useCase.pegoutProvider.PegoutConfiguration(ctx)
	if pegoutConfig.BridgeTransactionMin.Cmp(totalValue) > 0 {
		log.Infof(
			"Refunded pegouts total value: %v out of %v. Threshold not met yet. Skipping transaction to the bridge.",
			totalValue.ToRbtc(),
			pegoutConfig.BridgeTransactionMin.ToRbtc(),
		)
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)
	}

	useCase.rskWalletMutex.Lock()
	defer useCase.rskWalletMutex.Unlock()

	switch useCase.strategy {
	case UtxoSplit:
		return useCase.runUtxoSplit(ctx, totalValue, pegoutConfig, watchedQuotes)
	default:
		return useCase.runAllAtOnce(ctx, totalValue, watchedQuotes)
	}
}

func (useCase *BridgePegoutUseCase) checkBalance(ctx context.Context, requiredBalance *entities.Wei) error {
	balance, err := useCase.rskWallet.GetBalance(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	if balance.Cmp(requiredBalance) < 0 {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.InsufficientAmountError)
	}
	return nil
}

func (useCase *BridgePegoutUseCase) runAllAtOnce(ctx context.Context, totalValue *entities.Wei, watchedQuotes []quote.WatchedPegoutQuote) error {
	requiredBalance := new(entities.Wei).Add(totalValue, entities.NewWei(BridgeConversionGasLimit*BridgeConversionGasPrice))
	if err := useCase.checkBalance(ctx, requiredBalance); err != nil {
		return err
	}

	config := blockchain.NewTransactionConfig(totalValue, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr := useCase.rskWallet.SendRbtc(ctx, config, useCase.contracts.Bridge.GetAddress())
	if txErr == nil {
		log.Debugf("%s: transaction sent to the bridge successfully (%s)", usecases.BridgePegoutId, receipt.TransactionHash)
	}

	if err := useCase.updateQuotes(ctx, receipt, txErr, watchedQuotes); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	return nil
}

func (useCase *BridgePegoutUseCase) runUtxoSplit(
	ctx context.Context,
	totalValue *entities.Wei,
	pegoutConfig liquidity_provider.PegoutConfiguration,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	bridgeMin := pegoutConfig.BridgeTransactionMin
	numTxs, err := new(entities.Wei).Div(totalValue, bridgeMin)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	remainder := new(entities.Wei).Sub(totalValue, new(entities.Wei).Mul(numTxs, bridgeMin))

	gasPerTx := entities.NewWei(BridgeConversionGasLimit * BridgeConversionGasPrice)
	requiredBalance := new(entities.Wei).Add(totalValue, new(entities.Wei).Mul(numTxs, gasPerTx))
	if err := useCase.checkBalance(ctx, requiredBalance); err != nil {
		return err
	}

	bridgeAddress := useCase.contracts.Bridge.GetAddress()
	n := numTxs.Uint64()

	// First chunk absorbs the remainder (when N=1, firstChunk == totalValue)
	firstChunk := new(entities.Wei).Add(bridgeMin.Copy(), remainder)
	var receipt blockchain.TransactionReceipt
	var txErr error

	config := blockchain.NewTransactionConfig(firstChunk, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr = useCase.rskWallet.SendRbtc(ctx, config, bridgeAddress)
	if txErr == nil {
		log.Debugf("%s: split tx 1/%d sent to the bridge successfully (%s)", usecases.BridgePegoutId, n, receipt.TransactionHash)
	} else {
		if err := useCase.updateQuotes(ctx, receipt, txErr, watchedQuotes); err != nil {
			return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
		}
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, txErr)
	}

	for i := uint64(1); i < n; i++ {
		config = blockchain.NewTransactionConfig(bridgeMin.Copy(), BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
		receipt, txErr = useCase.rskWallet.SendRbtc(ctx, config, bridgeAddress)
		if txErr == nil {
			log.Debugf("%s: split tx %d/%d sent to the bridge successfully (%s)", usecases.BridgePegoutId, i+1, n, receipt.TransactionHash)
		} else {
			break
		}
	}

	if err := useCase.updateQuotes(ctx, receipt, txErr, watchedQuotes); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	return nil
}

func (useCase *BridgePegoutUseCase) updateQuotes(
	ctx context.Context,
	receipt blockchain.TransactionReceipt,
	txErr error,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	var err error
	retainedQuotes := make([]quote.RetainedPegoutQuote, 0)
	err = errors.Join(err, txErr)
	for _, watchedQuote := range watchedQuotes {
		if receipt.TransactionHash != "" {
			watchedQuote.RetainedQuote.BridgeRefundTxHash = receipt.TransactionHash
			watchedQuote.RetainedQuote.BridgeRefundGasUsed = receipt.GasUsed.Uint64()
			watchedQuote.RetainedQuote.BridgeRefundGasPrice = receipt.GasPrice
		}
		if txErr == nil {
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
		} else {
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxFailed
		}
		retainedQuotes = append(retainedQuotes, watchedQuote.RetainedQuote)
	}
	if updateErr := useCase.quoteRepository.UpdateRetainedQuotes(ctx, retainedQuotes); updateErr != nil {
		err = errors.Join(err, updateErr)
	}
	return err
}

func (useCase *BridgePegoutUseCase) calculateTotalToPegout(watchedQuotes []quote.WatchedPegoutQuote) (*entities.Wei, error) {
	totalValue := new(entities.Wei)
	for _, watchedQuote := range watchedQuotes {
		if watchedQuote.RetainedQuote.State != quote.PegoutStateRefundPegOutSucceeded {
			return nil, errors.New("not all quotes were refunded successfully")
		}
		if watchedQuote.PegoutQuote.Value != nil {
			totalValue.Add(totalValue, watchedQuote.PegoutQuote.Value)
		}
		if watchedQuote.PegoutQuote.CallFee != nil {
			totalValue.Add(totalValue, watchedQuote.PegoutQuote.CallFee)
		}
		if watchedQuote.PegoutQuote.GasFee != nil {
			totalValue.Add(totalValue, watchedQuote.PegoutQuote.GasFee)
		}
	}
	return totalValue, nil
}
