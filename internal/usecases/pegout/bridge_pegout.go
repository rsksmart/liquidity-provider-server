package pegout

import (
	"context"
	"errors"
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

type BridgePegoutUseCase struct {
	quoteRepository quote.PegoutQuoteRepository
	pegoutProvider  liquidity_provider.PegoutLiquidityProvider
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
	rskWalletMutex  sync.Locker
}

func NewBridgePegoutUseCase(
	quoteRepository quote.PegoutQuoteRepository,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
	rskWalletMutex sync.Locker,
) *BridgePegoutUseCase {
	return &BridgePegoutUseCase{
		quoteRepository: quoteRepository,
		pegoutProvider:  pegoutProvider,
		rskWallet:       rskWallet,
		contracts:       contracts,
		rskWalletMutex:  rskWalletMutex,
	}
}

func (useCase *BridgePegoutUseCase) Run(ctx context.Context, watchedQuotes ...quote.WatchedPegoutQuote) error {
	var err error
	var balance, totalValue *entities.Wei

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

	requiredBalance := new(entities.Wei).Add(totalValue, entities.NewWei(BridgeConversionGasLimit*BridgeConversionGasPrice))
	if balance, err = useCase.rskWallet.GetBalance(ctx); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	} else if balance.Cmp(requiredBalance) < 0 {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.InsufficientAmountError)
	}

	config := blockchain.NewTransactionConfig(totalValue, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr := useCase.rskWallet.SendRbtc(ctx, config, useCase.contracts.Bridge.GetAddress())
	if txErr == nil {
		log.Debugf("%s: transaction sent to the bridge successfully (%s)", usecases.BridgePegoutId, receipt.TransactionHash)
	}

	err = useCase.updateQuotes(ctx, receipt, txErr, watchedQuotes)
	if err != nil {
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
