package pegout

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type AllAtOnceHandler struct {
	quoteRepository quote.PegoutQuoteRepository
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
}

func NewAllAtOnceHandler(
	quoteRepository quote.PegoutQuoteRepository,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
) *AllAtOnceHandler {
	return &AllAtOnceHandler{
		quoteRepository: quoteRepository,
		rskWallet:       rskWallet,
		contracts:       contracts,
	}
}

func (h *AllAtOnceHandler) Execute(
	ctx context.Context,
	pegoutConfig liquidity_provider.PegoutConfiguration,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	totalValue := calculateTotalToPegout(watchedQuotes)
	if pegoutConfig.BridgeTransactionMin.Cmp(totalValue) > 0 {
		log.Infof(
			"Refunded pegouts total value: %v out of %v. Threshold not met yet. Skipping transaction to the bridge.",
			totalValue.ToRbtc(),
			pegoutConfig.BridgeTransactionMin.ToRbtc(),
		)
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)
	}

	requiredBalance := new(entities.Wei).Add(totalValue, entities.NewWei(BridgeConversionGasLimit*BridgeConversionGasPrice))
	if err := checkBalance(ctx, h.rskWallet, requiredBalance); err != nil {
		return err
	}

	config := blockchain.NewTransactionConfig(totalValue, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr := h.rskWallet.SendRbtc(ctx, config, h.contracts.Bridge.GetAddress())
	if txErr == nil {
		log.Debugf("%s: transaction sent to the bridge successfully (%s)", usecases.BridgePegoutId, receipt.TransactionHash)
	}

	if err := updateQuotes(ctx, h.quoteRepository, receipt, txErr, watchedQuotes); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	return nil
}

func updateQuotes(
	ctx context.Context,
	quoteRepository quote.PegoutQuoteRepository,
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
			watchedQuote.RetainedQuote.BridgeRebalances = []quote.BridgeRebalanceAllocation{{
				TxHash:   receipt.TransactionHash,
				GasUsed:  receipt.GasUsed.Uint64(),
				GasPrice: receipt.GasPrice,
			}}
		}
		if txErr == nil {
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
		} else {
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxFailed
		}
		retainedQuotes = append(retainedQuotes, watchedQuote.RetainedQuote)
	}
	if updateErr := quoteRepository.UpdateRetainedQuotes(ctx, retainedQuotes); updateErr != nil {
		err = errors.Join(err, updateErr)
	}
	return err
}
